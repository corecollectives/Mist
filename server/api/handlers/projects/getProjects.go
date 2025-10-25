package projects

import (
	"net/http"

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
	projectRows, err := h.DB.Query(`
            SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at
            FROM projects p
            JOIN project_members pm ON pm.project_id = p.id
            WHERE pm.user_id = ?`, userId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}
	defer projectRows.Close()

	var projects []models.Project

	for projectRows.Next() {
		var p models.Project
		if err := projectRows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to scan data", err.Error())
		}

		memberRows, err := h.DB.Query(`
                SELECT u.id, u.username, u.email, u.password_hash, u.role, u.created_at, u.updated_at
                FROM users u
                JOIN project_members pm ON pm.user_id = u.id
                WHERE pm.project_id = ?`, p.ID)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch project members", err.Error())
			return
		}
		var members []models.User
		for memberRows.Next() {
			var member models.User
			if err := memberRows.Scan(&member.ID, &member.Username, &member.Email, &member.PasswordHash, &member.Role, &member.CreatedAt, &member.UpdatedAt); err != nil {
				memberRows.Close()
				handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to scan member data", err.Error())
				return
			}
			members = append(members, member)
		}
		memberRows.Close()
		if err := memberRows.Err(); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "error itereating members", err.Error())
		}
		p.ProjectMembers = members
		projects = append(projects, p)

		if err := projectRows.Err(); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Error iterating projects", err.Error())
			return
		}

		if projects == nil {
			projects = []models.Project{}
		}

		handlers.SendResponse(w, http.StatusOK, true, projects, "Projects retrieved successfully", "")
	}

}
