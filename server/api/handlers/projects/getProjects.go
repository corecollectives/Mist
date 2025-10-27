package projects

import (
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	userId := userData.ID

	rows, err := h.DB.Query(`
		SELECT 
			p.id, p.name, p.description, p.tags, p.owner_id, p.created_at, p.updated_at,
			owner.id, owner.username, owner.email, owner.role, owner.created_at, owner.updated_at,
			u.id, u.username, u.email, u.role, u.created_at, u.updated_at

		FROM projects p
		LEFT JOIN users owner ON owner.id = p.owner_id        -- join owner here
		JOIN project_members pm ON pm.project_id = p.id
		JOIN users u ON u.id = pm.user_id                     -- project members
		WHERE pm.user_id = ? OR p.owner_id = ?
		ORDER BY p.updated_at DESC;
	`, userId, userId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}
	defer rows.Close()

	projectMap := make(map[int]*models.Project)

	for rows.Next() {
		var p models.Project
		var owner models.User
		var member models.User
		var tags string

		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &tags, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt,
			&owner.ID, &owner.Username, &owner.Email, &owner.Role, &owner.CreatedAt, &owner.UpdatedAt,
			&member.ID, &member.Username, &member.Email, &member.Role, &member.CreatedAt, &member.UpdatedAt,
		); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to scan row", err.Error())
			return
		}

		if tags != "" {
			p.Tags = strings.Split(tags, ",")
		}

		if _, exists := projectMap[int(p.ID)]; !exists {
			p.ProjectMembers = []models.User{}
			p.Owner = &owner // ✅ Set owner once
			projectMap[int(p.ID)] = &p
		}

		projectMap[int(p.ID)].ProjectMembers = append(projectMap[int(p.ID)].ProjectMembers, member)
	}

	projects := make([]models.Project, 0, len(projectMap))
	for _, p := range projectMap {
		projects = append(projects, *p)
	}

	handlers.SendResponse(w, http.StatusOK, true, projects, "Projects retrieved successfully", "")
}
