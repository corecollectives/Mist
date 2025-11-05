package applications

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetApplicationById(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		AppID int64 `json:"appId"`
	}
	// parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	// check if user is part of the project
	var isProjectMember bool
	err := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM project_members pm
			JOIN apps a ON pm.project_id = a.project_id
			WHERE a.id = ? AND pm.user_id = ?
		)
	`, req.AppID, userInfo.ID).Scan(&isProjectMember)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", "Internal Server Error")
		return
	}

	if !isProjectMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this application", "Forbidden")
		return
	}

	// query app by ID
	var app models.App
	err = h.DB.QueryRow(`
		SELECT 
			id, project_id, created_by, name, description, git_provider_id, git_repository, git_branch,
			deployment_strategy, port, root_directory, build_command, start_command, dockerfile_path,
			healthcheck_path, healthcheck_interval, status, created_at, updated_at
		FROM apps
		WHERE id = ?
	`, req.AppID).Scan(
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
		fmt.Println("Error querying app by ID:", err)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", "Internal Server Error")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, app, "Application retrieved successfully", "")

}
