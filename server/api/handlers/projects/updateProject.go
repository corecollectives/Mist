package projects

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	projectIDStr := r.URL.Query().Get("id")
	if projectIDStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing project ID", "project id is required")
		return
	}
	projectId, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid project ID", err.Error())
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
	}

	var ownerId int64
	err = h.DB.QueryRow("SELECT owner_id FROM projects WHERE id=?", projectId).Scan(&ownerId)
	if err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "no such project")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", err.Error())
		return
	}

	if ownerId != userData.ID {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}

	tagsString := ""
	if len(input.Tags) > 0 {
		for i, tag := range input.Tags {
			if i > 0 {
				tagsString += ","
			}
			tagsString += tag
		}
	}

	_, err = h.DB.Exec("UPDATE projects SET name = ?, description = ?, tags = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", input.Name, input.Description, tagsString, projectId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update project", err.Error())
		return
	}

	// Fetch the updated project with owner information
	var project models.Project
	var tags string = ""
	var ownerUsername, ownerEmail, ownerRole string
	var ownerCreatedAt, ownerUpdatedAt time.Time

	err = h.DB.QueryRow(`
		SELECT p.id, p.name, p.description, p.tags, p.owner_id, p.created_at, p.updated_at,
		       u.username, u.email, u.role, u.created_at, u.updated_at
		FROM projects p
		JOIN users u ON p.owner_id = u.id
		WHERE p.id = ?
	`, projectId).Scan(
		&project.ID, &project.Name, &project.Description, &tags, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt,
		&ownerUsername, &ownerEmail, &ownerRole, &ownerCreatedAt, &ownerUpdatedAt,
	)

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch updated project", err.Error())
		return
	}

	if tags != "" {
		// project.Tags = strings.Split(tags, ",")
	}

	// Populate the owner field
	project.Owner = &models.User{
		ID:        project.OwnerID,
		Username:  ownerUsername,
		Email:     ownerEmail,
		Role:      ownerRole,
		CreatedAt: ownerCreatedAt,
		UpdatedAt: ownerUpdatedAt,
	}

	handlers.SendResponse(w, http.StatusOK, true, project, "Project updated successfully", "")

}
