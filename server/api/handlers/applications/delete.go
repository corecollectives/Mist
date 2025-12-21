package applications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64 `json:"appId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", fmt.Sprintf("Error fetching application: %v", err))
		return
	}
	if app == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", "No application with the given ID exists")
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, app.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this application", "Forbidden")
		return
	}

	containerName := docker.GetContainerName(app.Name, app.ID)

	if docker.ContainerExists(containerName) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		stopCmd := exec.CommandContext(ctx, "docker", "stop", containerName)
		if err := stopCmd.Run(); err != nil {
			log.Warn().Err(err).Str("container", containerName).Msg("Failed to stop container during app deletion")
		}

		removeCmd := exec.CommandContext(ctx, "docker", "rm", containerName)
		if err := removeCmd.Run(); err != nil {
			log.Warn().Err(err).Str("container", containerName).Msg("Failed to remove container during app deletion")
		}
	}

	imagePattern := fmt.Sprintf("mist-app-%d-", app.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	listImagesCmd := exec.CommandContext(ctx, "docker", "images", "-q", "--filter", fmt.Sprintf("reference=%s*", imagePattern))
	output, err := listImagesCmd.Output()
	if err == nil && len(output) > 0 {
		rmiCmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("docker images -q --filter 'reference=%s*' | xargs -r docker rmi -f", imagePattern))
		if err := rmiCmd.Run(); err != nil {
			log.Warn().Err(err).Int64("app_id", app.ID).Msg("Failed to remove Docker images during app deletion")
		}
	}

	appPath := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", app.ProjectID, app.Name)
	if _, err := os.Stat(appPath); err == nil {
		if err := os.RemoveAll(appPath); err != nil {
			log.Warn().Err(err).Str("path", appPath).Msg("Failed to remove app directory during deletion")
		}
	}

	logPath := constants.Constants["LogPath"].(string)
	logPattern := filepath.Join(logPath, fmt.Sprintf("*%d_build_logs", app.ID))
	matches, err := filepath.Glob(logPattern)
	if err == nil {
		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				log.Warn().Err(err).Str("log_file", match).Msg("Failed to remove log file during app deletion")
			}
		}
	}

	err = models.DeleteApplication(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete application from database", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "delete", "application", &req.AppID, map[string]interface{}{
		"app_name":   app.Name,
		"project_id": app.ProjectID,
	})

	handlers.SendResponse(w, http.StatusOK, true, nil, "Application deleted successfully", "")
}
