package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func GetSystemVersion(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Admin access required", "Forbidden")
		return
	}

	version, err := models.GetSystemVersion()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get version", err.Error())
		return
	}

	buildDate, _ := models.GetSystemInfo("build_date")

	response := map[string]interface{}{
		"version":   version,
		"buildDate": buildDate,
	}

	handlers.SendResponse(w, http.StatusOK, true, response, "Version retrieved", "")
}

func CheckForUpdates(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Admin access required", "Forbidden")
		return
	}

	resp, err := http.Get("https://api.github.com/repos/corecollectives/mist/releases/latest")
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to check for updates", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{
			"hasUpdate": false,
			"message":   "No releases available yet",
		}, "No updates available", "")
		return
	}

	if resp.StatusCode != 200 {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch release info", fmt.Sprintf("GitHub API returned status: %d", resp.StatusCode))
		return
	}

	var release models.GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to parse release info", err.Error())
		return
	}

	currentVersion, err := models.GetSystemVersion()
	if err != nil {
		currentVersion = "1.0.0"
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion = strings.TrimPrefix(currentVersion, "v")

	hasUpdate := latestVersion != currentVersion

	response := map[string]interface{}{
		"hasUpdate":      hasUpdate,
		"currentVersion": currentVersion,
		"latestVersion":  latestVersion,
		"release":        release,
	}

	handlers.SendResponse(w, http.StatusOK, true, response, "Update check completed", "")
}

func TriggerUpdate(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Admin access required", "Forbidden")
		return
	}

	var req struct {
		Version string `json:"version"`
		Branch  string `json:"branch"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.Version == "" {
		req.Version = "latest"
	}

	if req.Branch == "" {
		req.Branch = "main"
	}

	currentVersion, err := models.GetSystemVersion()
	if err != nil {
		currentVersion = "unknown"
	}

	updateID, err := models.CreateUpdateHistory(currentVersion, req.Version, userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create update record", err.Error())
		return
	}

	scriptPath := "/opt/mist/scripts/update.sh"

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = filepath.Join("scripts", "update.sh")
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			errorMsg := "Update script not found"
			models.UpdateUpdateHistoryStatus(updateID, "failed", &errorMsg)
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Update script not found", err.Error())
			return
		}
	}

	go func() {
		cmd := exec.Command("bash", scriptPath, req.Version, req.Branch)
		cmd.Dir = "/opt/mist"

		models.UpdateUpdateHistoryStatus(updateID, "downloading", nil)

		output, err := cmd.CombinedOutput()
		if err != nil {
			errorMsg := fmt.Sprintf("Update failed: %s\n%s", err.Error(), string(output))
			models.UpdateUpdateHistoryStatus(updateID, "failed", &errorMsg)
			return
		}

		models.SetSystemVersion(req.Version)
		models.UpdateUpdateHistoryStatus(updateID, "success", nil)
	}()

	time.Sleep(500 * time.Millisecond)

	response := map[string]interface{}{
		"updateId": updateID,
		"message":  "Update started. The system will restart shortly.",
	}

	handlers.SendResponse(w, http.StatusOK, true, response, "Update initiated", "")
}

func GetUpdateHistory(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Admin access required", "Forbidden")
		return
	}

	history, err := models.GetUpdateHistory(50)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get update history", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, history, "Update history retrieved", "")
}

func GetUpdateStatus(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Admin access required", "Forbidden")
		return
	}

	update, err := models.GetLatestUpdateHistory()
	if err != nil {
		handlers.SendResponse(w, http.StatusOK, true, nil, "No updates found", "")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, update, "Update status retrieved", "")
}

func GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userInfo.Role == "user" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Admin access required", "Forbidden")
		return
	}

	cmd := exec.Command("systemctl", "is-active", "mist")
	output, err := cmd.Output()
	serviceActive := err == nil && strings.TrimSpace(string(output)) == "active"

	var diskFree, diskTotal uint64
	cmd = exec.Command("df", "-B1", "/opt/mist")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				fmt.Sscanf(fields[1], "%d", &diskTotal)
				fmt.Sscanf(fields[3], "%d", &diskFree)
			}
		}
	}

	var uptime string
	cmd = exec.Command("systemctl", "show", "mist", "--property=ActiveEnterTimestamp", "--value")
	if output, err := cmd.Output(); err == nil {
		uptime = strings.TrimSpace(string(output))
	}

	health := map[string]interface{}{
		"serviceActive": serviceActive,
		"diskFree":      diskFree,
		"diskTotal":     diskTotal,
		"uptime":        uptime,
	}

	handlers.SendResponse(w, http.StatusOK, true, health, "System health retrieved", "")
}
