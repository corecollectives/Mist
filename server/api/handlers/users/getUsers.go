package users

import (
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT id, username, email, role, created_at, updated_at
		FROM users
		ORDER BY id;
	`)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to scan row", err.Error())
			return
		}
		users = append(users, u)
	}

	handlers.SendResponse(w, http.StatusOK, true, users, "Users retrieved successfully", "")
}
