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

	// Step 1: Check if user has access (either member or owner)
	var hasAccess bool
	err = h.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM projects p
			LEFT JOIN project_members pm ON pm.project_id = p.id
			WHERE p.id = ? AND (p.owner_id = ? OR pm.user_id = ?)
		)
	`, projectId, userID, userID).Scan(&hasAccess)

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}

	if !hasAccess {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You are not a member of this project", "")
		return
	}

	// Step 2: Fetch project info + all members
	rows, err := h.DB.Query(`
		SELECT 
			p.id, p.name, p.description, p.tags, p.owner_id, p.created_at, p.updated_at,
			o.id, o.username, o.email, o.password_hash, o.role, o.created_at, o.updated_at,
			u.id, u.username, u.email, u.password_hash, u.role, u.created_at, u.updated_at
		FROM projects p
		JOIN users o ON o.id = p.owner_id
		LEFT JOIN project_members pm ON pm.project_id = p.id
		LEFT JOIN users u ON u.id = pm.user_id
		WHERE p.id = ?
		ORDER BY u.id
	`, projectId)
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
		var owner models.User
		var pID int64
		var pName, pDescription string
		var pOwnerID int64
		var pCreated, pUpdated time.Time
		var tags sql.NullString
		var memberID sql.NullInt64
		var memberUsername, memberEmail, memberPassword, memberRole sql.NullString
		var memberCreated, memberUpdated sql.NullTime

		if err := rows.Scan(
			&pID, &pName, &pDescription, &tags, &pOwnerID, &pCreated, &pUpdated,
			&owner.ID, &owner.Username, &owner.Email, &owner.PasswordHash, &owner.Role, &owner.CreatedAt, &owner.UpdatedAt,
			&memberID, &memberUsername, &memberEmail, &memberPassword, &memberRole, &memberCreated, &memberUpdated,
		); err != nil {
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
			project.Owner = &owner
			firstRow = false
		}

		if memberID.Valid && !membersMap[memberID.Int64] {
			member.ID = memberID.Int64
			member.Username = memberUsername.String
			member.Email = memberEmail.String
			member.PasswordHash = memberPassword.String
			member.Role = memberRole.String
			member.CreatedAt = memberCreated.Time
			member.UpdatedAt = memberUpdated.Time
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
