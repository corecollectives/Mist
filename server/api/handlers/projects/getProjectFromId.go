package projects

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetProjectFromId(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	userID := userData.ID

	projectIDStr := r.URL.Query().Get("id")
	if projectIDStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Missing project ID", "no id provided")
		return
	}

	projectId, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid project ID", err.Error())
		return
	}

	row := h.DB.QueryRow(`SELECT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at
            FROM projects p
            JOIN project_members pm ON pm.project_id = p.id
            WHERE p.id = ? AND pm.user_id = ?`, projectId, userID)

	var project models.Project

	if err := row.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", err.Error())
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}
	memberRows, err := h.DB.Query(`
            SELECT u.id, u.username, u.email, u.password_hash, u.role, u.created_at, u.updated_at
            FROM users u
            JOIN project_members pm ON pm.user_id = u.id
            WHERE pm.project_id = ?`, project.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed (members)", err.Error())
		return
	}
	defer memberRows.Close()

	for memberRows.Next() {
		var member models.User
		if err := memberRows.Scan(&member.ID, &member.Username, &member.Email, &member.PasswordHash, &member.Role, &member.CreatedAt, &member.UpdatedAt); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to scan data", err.Error())
		}
		project.ProjectMembers = append(project.ProjectMembers, member)
	}
	if err := memberRows.Err(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to scan data", err.Error())
	}
	handlers.SendResponse(w, http.StatusOK, true, project, "Project found", "")
}
