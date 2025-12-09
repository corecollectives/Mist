package applications

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID              int64   `json:"appId"`
		Name               *string `json:"name"`
		Description        *string `json:"description"`
		GitRepository      *string `json:"gitRepository"`
		GitBranch          *string `json:"gitBranch"`
		Port               *int    `json:"port"`
		RootDirectory      *string `json:"rootDirectory"`
		DockerfilePath     *string `json:"dockerfilePath"`
		DeploymentStrategy *string `json:"deploymentStrategy"`
		Status             *string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to update this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}
	if app == nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", "No application with the given ID exists")
		return
	}

	if req.Name != nil {
		app.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		app.Description = &trimmed
	}
	if req.GitRepository != nil {
		trimmed := strings.TrimSpace(*req.GitRepository)
		app.GitRepository = &trimmed
	}
	if req.GitBranch != nil {
		app.GitBranch = strings.TrimSpace(*req.GitBranch)
	}
	if req.Port != nil {
		port := int64(*req.Port)
		app.Port = &port
	}
	if req.RootDirectory != nil {
		app.RootDirectory = strings.TrimSpace(*req.RootDirectory)
	}
	if req.DockerfilePath != nil {
		trimmed := strings.TrimSpace(*req.DockerfilePath)
		app.DockerfilePath = &trimmed
	}
	if req.DeploymentStrategy != nil {
		app.DeploymentStrategy = models.DeploymentStrategy(strings.TrimSpace(*req.DeploymentStrategy))
	}
	if req.Status != nil {
		app.Status = models.AppStatus(strings.TrimSpace(*req.Status))
	}

	app.UpdatedAt = time.Now()

	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update application", err.Error())
		return
	}

	changes := make(map[string]interface{})
	if req.Name != nil {
		changes["name"] = *req.Name
	}
	if req.GitRepository != nil {
		changes["git_repository"] = *req.GitRepository
	}
	if req.GitBranch != nil {
		changes["git_branch"] = *req.GitBranch
	}
	if req.Port != nil {
		changes["port"] = *req.Port
	}
	if req.Status != nil {
		changes["status"] = *req.Status
	}
	models.LogUserAudit(userInfo.ID, "update", "application", &app.ID, map[string]interface{}{
		"changes": changes,
	})

	handlers.SendResponse(w, http.StatusOK, true, app.ToJson(), "Application updated successfully", "")
}
