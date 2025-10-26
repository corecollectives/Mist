package users

import (
	"database/sql"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	userIDParam := r.URL.Query().Get("id")
	if userIDParam == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User ID is required", "Missing 'id' parameter")
		return
	}

	var user models.User
	err := h.DB.QueryRow(`
        SELECT id, username, email, role, created_at, updated_at
        FROM users
        WHERE id = ?`, userIDParam,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", "No user exists with the given ID")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database query failed", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, user, "User retrieved successfully", "")
}
