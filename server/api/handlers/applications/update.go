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

	// ✅ Verify ownership
	var isOwner bool
	err := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM apps WHERE id = ? AND created_by = ?
		)
	`, req.AppID, userInfo.ID).Scan(&isOwner)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}
	if !isOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to update this app", "Forbidden")
		return
	}

	// ✅ Build dynamic update query
	setClauses := []string{}
	args := []interface{}{}

	if req.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, strings.TrimSpace(*req.Name))
	}
	if req.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, strings.TrimSpace(*req.Description))
	}
	if req.GitRepository != nil {
		setClauses = append(setClauses, "git_repository = ?")
		args = append(args, strings.TrimSpace(*req.GitRepository))
	}
	if req.GitBranch != nil {
		setClauses = append(setClauses, "git_branch = ?")
		args = append(args, strings.TrimSpace(*req.GitBranch))
	}
	if req.Port != nil {
		setClauses = append(setClauses, "port = ?")
		args = append(args, *req.Port)
	}
	if req.RootDirectory != nil {
		setClauses = append(setClauses, "root_directory = ?")
		args = append(args, strings.TrimSpace(*req.RootDirectory))
	}
	if req.DockerfilePath != nil {
		setClauses = append(setClauses, "dockerfile_path = ?")
		args = append(args, strings.TrimSpace(*req.DockerfilePath))
	}
	if req.DeploymentStrategy != nil {
		setClauses = append(setClauses, "deployment_strategy = ?")
		args = append(args, *req.DeploymentStrategy)
	}
	if req.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *req.Status)
	}

	// Always update timestamp
	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, time.Now())

	if len(setClauses) == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "No fields to update", "Nothing provided")
		return
	}

	query := `
		UPDATE apps
		SET ` + strings.Join(setClauses, ", ") + `
		WHERE id = ?
		RETURNING 
			id, project_id, created_by, name, description,
			git_provider_id, git_repository, git_branch, deployment_strategy,
			port, root_directory, build_command, start_command, dockerfile_path,
			healthcheck_path, healthcheck_interval, status, created_at, updated_at
	`
	args = append(args, req.AppID)

	var app models.App
	err = h.DB.QueryRow(query, args...).Scan(
		&app.ID,
		&app.ProjectID,
		&app.CreatedBy,
		&app.Name,
		&app.Description,
		&app.GitProviderID,
		&app.GitRepository,
		&app.GitBranch,
		&app.DeploymentStrategy,
		&app.Port,
		&app.RootDirectory,
		&app.BuildCommand,
		&app.StartCommand,
		&app.DockerfilePath,
		&app.HealthcheckPath,
		&app.HealthcheckInterval,
		&app.Status,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "App not found", "Invalid ID")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, app, "App updated successfully", "")
}
