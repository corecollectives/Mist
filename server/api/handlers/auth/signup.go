package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
	"github.com/corecollectives/mist/store"
)

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	db := h.DB
	setupRequired := store.IsSetupRequired()
	if !setupRequired {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Sign up not allowed", "Only first user can sign up")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "All fields are required", "Missing fields")
		return
	}

	user, err := utils.InsertUserInDb(db, req.Username, req.Email, req.Password, "owner")
	if err != nil {
		log.Printf("db insert error: %v", err)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Database error", "Internal Server Error")
		return
	}

	token, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to generate token", "Internal Server Error")
		return
	}
	if setupRequired {
		store.SetSetupRequired(false)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "mist_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24 * 30,
	})
	handlers.SendResponse(w, http.StatusCreated, true, user, "User signed up successfully", "")

}
