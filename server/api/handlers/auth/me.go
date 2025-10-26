package auth

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
	"github.com/corecollectives/mist/models"
)

func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("mist_token")
	setupRequired, _ := utils.IsSetupRequired(h.DB)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"setupRequired": setupRequired,
				"user":          nil,
			},
			"message": "Not logged in",
			"error":   "missing auth cookie",
		})
		return
	}

	tokenStr := cookie.Value

	claims, err := middleware.VerifyJWT(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"setupRequired": setupRequired, "user": nil},
			"message": "Invalid token",
			"error":   err.Error(),
		})
		return
	}

	userId := claims.UserID
	var user models.User
	row := h.DB.QueryRow(
		"SELECT id, username, email, role FROM users WHERE id = ?",
		userId,
	)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"setupRequired": setupRequired, "user": nil},
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    map[string]interface{}{"setupRequired": setupRequired, "user": user},
		"message": "User fetched successfully",
		"error":   "",
	})
}
