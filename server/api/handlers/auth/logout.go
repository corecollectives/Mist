package auth

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "mist_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production (HTTPS)
		MaxAge:   -1,    // expire immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    nil,
		"message": "Logged out successfully",
		"error":   "",
	})
}
