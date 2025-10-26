package projects

import (
	"encoding/json"
	"net/http"

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

	userRole := userData.Role
	if userRole != "owner" || userRole != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

	res, err := tx.Exec(`INSERT INTO projects(name,description,owner_id) VALUES(?,?,?)`, input.Name, input.Description, userData.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create project", err.Error())
		return
	}
	projectId, _ := res.LastInsertId()
	_, err = tx.Exec(`INSERT INTO project_members(user_id,project_id,role) VALUES(?,?,?)`, userData.ID, projectId, "owner")
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to add project owner", err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to commit transaction", err.Error())
		return
	}

	var project models.Project
	err = h.DB.QueryRow(`SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = ?`, projectId).
		Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch created project", err.Error())
		return
	}
	project.ProjectMembers = append(project.ProjectMembers, *userData)
	handlers.SendResponse(w, http.StatusCreated, true, project, "Project created successfully", "")
}
