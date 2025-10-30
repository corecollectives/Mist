package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/store"
)

func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("mist_token")
	setupRequired := store.IsSetupRequired()
	if err != nil {
		handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": nil}, "No auth cookie", "")
		return
	}

	tokenStr := cookie.Value

	claims, err := middleware.VerifyJWT(tokenStr)
	if err != nil {
		handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": nil}, "Invalid token", "")
		return
	}

	userId := claims.UserID
	var user models.User
	row := h.DB.QueryRow(
		"SELECT id, username, email, role FROM users WHERE id = ?",
		userId,
	)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": nil}, "User not found", "")
			return
		}
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch user", "database scan error")
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, map[string]interface{}{"setupRequired": setupRequired, "user": user}, "User fetched successfully", "")
}
