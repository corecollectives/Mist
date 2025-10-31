package projects

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userData.Role != "owner" && userData.Role != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}

	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if input.Name == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project name is required", "name field is empty")
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Internal server error", err.Error())
		return
	}
	defer tx.Rollback()

	var project models.Project
	tagsString := ""
	if len(input.Tags) > 0 {
		for i, tag := range input.Tags {
			if i > 0 {
				tagsString += ","
			}
			tagsString += tag
		}
	}

	var tags string = ""
	err = tx.QueryRow(`
		INSERT INTO projects(name, description,tags, owner_id)
		VALUES (?, ?, ?,?)
		RETURNING id, name, description, tags, owner_id, created_at, updated_at
	`, input.Name, input.Description, tagsString, userData.ID).
		Scan(&project.ID, &project.Name, &project.Description, &tags, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)

	if tags != "" {
		project.Tags = strings.Split(tags, ",")
	}
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create project", err.Error())
		return
	}

	_, err = tx.Exec(`INSERT INTO project_members(user_id, project_id) VALUES(?, ?)`, userData.ID, project.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to add project owner", err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to commit transaction", err.Error())
		return
	}

	project.ProjectMembers = []models.User{{
		ID:        userData.ID,
		Username:  userData.Username,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}}

	handlers.SendResponse(w, http.StatusCreated, true, project, "Project created successfully", "")
}
