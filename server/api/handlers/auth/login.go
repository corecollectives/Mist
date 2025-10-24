package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var cred struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	db := h.DB
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	cred.Email = strings.ToLower(strings.TrimSpace(cred.Email))

	if cred.Email == "" || cred.Password == "" {
		ErrorResponse(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	var user models.User
	err := db.QueryRowContext(r.Context(),
		`SELECT id, username, email, password_hash, role FROM users WHERE LOWER(email) = ?`,
		cred.Email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err == sql.ErrNoRows {
		ErrorResponse(w, "User not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("db query error: %v", err)
		ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cred.Password)); err != nil {
		ErrorResponse(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		ErrorResponse(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
