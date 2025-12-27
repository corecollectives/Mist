package updates

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

const updateLockFile = "/var/lib/mist/update.lock"

func TriggerUpdate(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userInfo.ID).Msg("Failed to verify user role for update trigger")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		log.Warn().Int64("user_id", userInfo.ID).Str("role", role).Msg("Non-owner attempted to trigger update")
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can trigger updates", "Forbidden")
		return
	}

	if _, err := os.Stat(updateLockFile); err == nil {
		log.Warn().Int64("user_id", userInfo.ID).Msg("Update already in progress")
		handlers.SendResponse(w, http.StatusConflict, false, nil, "Another update is already in progress", "Please wait for it to complete")
		return
	}

	inProgressLogs, err := models.GetUpdateLogs(1)
	if err == nil && len(inProgressLogs) > 0 {
		if inProgressLogs[0].Status == "in_progress" {
			if time.Since(inProgressLogs[0].StartedAt) < 10*time.Minute {
				log.Warn().Int64("user_id", userInfo.ID).Int64("existing_log_id", inProgressLogs[0].ID).Msg("Recent update still in progress")
				handlers.SendResponse(w, http.StatusConflict, false, nil, "An update is currently in progress", "Please wait for it to complete")
				return
			} else {
				log.Warn().Int64("stale_log_id", inProgressLogs[0].ID).Msg("Found stale in-progress update, continuing")
			}
		}
	}

	currentVersion, err := models.GetSystemSetting("version")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current version before update")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get current version", err.Error())
		return
	}
	if currentVersion == "" {
		currentVersion = "1.0.0"
	}

	targetVersion := "unknown"
	resp, err := http.Get("https://api.github.com/repos/corecollectives/mist/releases/latest")
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var release GithubRelease
			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(&release); err == nil {
				targetVersion = strings.TrimPrefix(release.TagName, "v")
			}
		}
	}

	updateLog, err := models.CreateUpdateLog(currentVersion, targetVersion, userInfo.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create update log")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create update log", err.Error())
		return
	}

	log.Info().
		Int64("update_log_id", updateLog.ID).
		Str("from_version", currentVersion).
		Str("to_version", targetVersion).
		Int64("user_id", userInfo.ID).
		Msg("Update triggered")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error().Msg("Streaming not supported for update")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Streaming not supported", "")
		return
	}

	var logBuilder strings.Builder

	cmd := exec.Command("/bin/bash", "/opt/mist/install.sh")

	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errMsg := "Error setting up stdout pipe: " + err.Error() + "\n"
		w.Write([]byte(errMsg))
		logBuilder.WriteString(errMsg)
		models.UpdateUpdateLogStatus(updateLog.ID, "failed", logBuilder.String(), &errMsg)
		log.Error().Err(err).Int64("update_log_id", updateLog.ID).Msg("Failed to set up stdout pipe")
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		errMsg := "Error setting up stderr pipe: " + err.Error() + "\n"
		w.Write([]byte(errMsg))
		logBuilder.WriteString(errMsg)
		models.UpdateUpdateLogStatus(updateLog.ID, "failed", logBuilder.String(), &errMsg)
		log.Error().Err(err).Int64("update_log_id", updateLog.ID).Msg("Failed to set up stderr pipe")
		return
	}

	if err := cmd.Start(); err != nil {
		errMsg := "Error starting update: " + err.Error() + "\n"
		w.Write([]byte(errMsg))
		logBuilder.WriteString(errMsg)
		models.UpdateUpdateLogStatus(updateLog.ID, "failed", logBuilder.String(), &errMsg)
		log.Error().Err(err).Int64("update_log_id", updateLog.ID).Msg("Failed to start update process")
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			w.Write([]byte(line))
			flusher.Flush()
			logBuilder.WriteString(line)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			w.Write([]byte(line))
			flusher.Flush()
			logBuilder.WriteString(line)
		}
	}()

	if err := cmd.Wait(); err != nil {
		errMsg := "\n❌ Update failed: " + err.Error() + "\n"
		w.Write([]byte(errMsg))
		logBuilder.WriteString(errMsg)

		finalLog := logBuilder.String()
		models.UpdateUpdateLogStatus(updateLog.ID, "failed", finalLog, &errMsg)

		log.Error().
			Err(err).
			Int64("update_log_id", updateLog.ID).
			Str("from_version", currentVersion).
			Str("to_version", targetVersion).
			Msg("Update failed")

		dummyID := int64(1)
		models.LogUserAudit(userInfo.ID, "update", "system", &dummyID, map[string]any{
			"status":        "failed",
			"error":         err.Error(),
			"update_log_id": updateLog.ID,
		})
		return
	}

	w.Write([]byte("\n✅ Update completed successfully\n"))
	logBuilder.WriteString("\n✅ Update completed successfully\n")
	flusher.Flush()

	newVersion := targetVersion
	resp2, err := http.Get("https://api.github.com/repos/corecollectives/mist/releases/latest")
	if err == nil {
		defer resp2.Body.Close()
		if resp2.StatusCode == http.StatusOK {
			var release GithubRelease
			decoder := json.NewDecoder(resp2.Body)
			if err := decoder.Decode(&release); err == nil {
				newVersion = strings.TrimPrefix(release.TagName, "v")
				if newVersion != "" {
					models.SetSystemSetting("version", newVersion)
					log.Info().Str("new_version", newVersion).Msg("Version updated in database")
				}
			}
		}
	}

	finalLog := logBuilder.String()
	models.UpdateUpdateLogStatus(updateLog.ID, "success", finalLog, nil)

	log.Info().
		Int64("update_log_id", updateLog.ID).
		Str("from_version", currentVersion).
		Str("to_version", newVersion).
		Int64("user_id", userInfo.ID).
		Msg("Update completed successfully")

	dummyID := int64(1)
	models.LogUserAudit(userInfo.ID, "update", "system", &dummyID, map[string]any{
		"status":        "success",
		"from_version":  currentVersion,
		"to_version":    newVersion,
		"update_log_id": updateLog.ID,
	})
}
