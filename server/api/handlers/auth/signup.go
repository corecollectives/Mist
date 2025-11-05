package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/store"
)

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

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

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "owner",
	}

	err := user.SetPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to process password", "Internal Server Error")
		return
	}
	err = user.Create()
	if err != nil {
		println(err.Error())
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create user", "Internal Server Error")
		return
	}
	fmt.Println(user)

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
