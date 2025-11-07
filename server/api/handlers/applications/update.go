package applications

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
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
		app.Description = sql.NullString{String: strings.TrimSpace(*req.Description), Valid: true}
	}
	if req.GitRepository != nil {
		app.GitRepository = sql.NullString{String: strings.TrimSpace(*req.GitRepository), Valid: true}
	}
	if req.GitBranch != nil {
		app.GitBranch = strings.TrimSpace(*req.GitBranch)
	}
	if req.Port != nil {
		app.Port = sql.NullInt64{Int64: int64(*req.Port), Valid: true}
	}
	if req.RootDirectory != nil {
		app.RootDirectory = strings.TrimSpace(*req.RootDirectory)
	}
	if req.DockerfilePath != nil {
		app.DockerfilePath = sql.NullString{String: strings.TrimSpace(*req.DockerfilePath), Valid: true}
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

	handlers.SendResponse(w, http.StatusOK, true, app.ToJson(), "Application updated successfully", "")
}
