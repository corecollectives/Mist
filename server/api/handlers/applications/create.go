package applications

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ProjectID   int64  `json:"projectId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.Name == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Application name is required", "Missing fields")
		return
	}

	if req.ProjectID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project ID is required", "Missing fields")
		return
	}

	// check if user is part of the project
	var isProjectMember bool
	err := h.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM project_members 
			WHERE project_id = ? AND user_id = ?
		)
	`, req.ProjectID, userInfo.ID).Scan(&isProjectMember)
	if err != nil {
		fmt.Println("Error checking project membership:", err)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", "Internal Server Error")
		return
	}
	if !isProjectMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this project", "Forbidden")
		return
	}

	var app models.App
	app.ID = utils.GenerateRandomId()

	query := `
INSERT INTO apps (
    id, name, description, project_id, created_by
) VALUES (?, ?, ?, ?, ?)
RETURNING 
    id,
    project_id,
    created_by,
    name,
    description,
    git_provider_id,
    git_repository,
    git_branch,
    deployment_strategy,
    port,
    root_directory,
    build_command,
    start_command,
    dockerfile_path,
    healthcheck_path,
    healthcheck_interval,
    status,
    created_at,
    updated_at
`

	err = h.DB.QueryRow(
		query,
		app.ID,
		req.Name,
		req.Description,
		req.ProjectID,
		userInfo.ID,
	).Scan(
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
		fmt.Println("Error creating application:", err.Error())
		handlers.SendResponse(
			w,
			http.StatusInternalServerError,
			false,
			nil,
			"Failed to create application",
			"Internal Server Error",
		)
		return
	}

	handlers.SendResponse(
		w,
		http.StatusCreated,
		true,
		app,
		"Application created successfully",
		"",
	)

}
