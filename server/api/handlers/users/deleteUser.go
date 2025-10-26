package users

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/store"
)

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	if userData.Role != "owner" && userData.Role != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}
	userIDParam := r.URL.Query().Get("id")
	if userIDParam == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User ID is required", "Missing 'id' parameter")
		return
	}

	if userIDParam == strconv.FormatInt(userData.ID, 10) {
		_, err := h.DB.Exec(`DELETE FROM users WHERE id = ?`, userIDParam)
		if err != nil {
			handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete user", err.Error())
			return
		}
		http.Redirect(w, r, "/api/auth/logout", http.StatusSeeOther)
		store.SetSetupRequired(h.DB)
		return
	}
	userToDeleteRole := ""
	err := h.DB.QueryRow(`SELECT role FROM users WHERE id = ?`, userIDParam).Scan(&userToDeleteRole)
	if err == sql.ErrNoRows {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", "No such user")
		return
	} else if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve user role", err.Error())
		return
	}

	if userData.Role == "admin" && userToDeleteRole == "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized to delete owner", "Forbidden")
		return
	}
	if userData.Role == "admin" && userToDeleteRole == "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized to delete admin", "Forbidden")
		return
	}

	_, err = h.DB.Exec(`DELETE FROM users WHERE id = ?`, userIDParam)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to delete user", err.Error())
		return
	}
	store.SetSetupRequired(h.DB)
	handlers.SendResponse(w, http.StatusOK, true, nil, "User deleted successfully", "")

}
