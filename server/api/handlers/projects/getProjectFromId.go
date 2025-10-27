package projects

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	rows, err := h.DB.Query(`
		SELECT 
			p.id, p.name, p.description,p.tags, p.owner_id, p.created_at, p.updated_at,
			u.id, u.username, u.email, u.password_hash, u.role, u.created_at, u.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		JOIN users u ON u.id = pm.user_id
		WHERE p.id = ? AND pm.user_id = ?
		ORDER BY u.id
	`, projectId, userID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}
	defer rows.Close()

	var project models.Project

	membersMap := make(map[int64]bool)
	firstRow := true

	for rows.Next() {
		var member models.User
		var pID int64
		var pName, pDescription string
		var pOwnerID int64
		var pCreated, pUpdated time.Time
		var tags sql.NullString

		if err := rows.Scan(&pID, &pName, &pDescription, &tags, &pOwnerID, &pCreated, &pUpdated,
			&member.ID, &member.Username, &member.Email, &member.PasswordHash, &member.Role, &member.CreatedAt, &member.UpdatedAt); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to scan data", err.Error())
			return
		}

		if firstRow {
			project.ID = pID
			project.Name = pName
			project.Description = pDescription
			project.OwnerID = pOwnerID
			project.CreatedAt = pCreated
			project.UpdatedAt = pUpdated
			if tags.Valid && tags.String != "" {
				project.Tags = strings.Split(tags.String, ",")
			}
			firstRow = false
		}

		if !membersMap[member.ID] {
			project.ProjectMembers = append(project.ProjectMembers, member)
			membersMap[member.ID] = true
		}
	}

	if firstRow {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, project, "Project found", "")
}
