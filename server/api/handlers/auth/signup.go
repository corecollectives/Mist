package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
	"github.com/corecollectives/mist/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	db := h.DB

	// only first user i.e. the owner should be able to sign up rest users will be created by the owner later but not throug signup
	row := db.QueryRow("SELECT COUNT(*) FROM users")
	var count int
	if err := row.Scan(&count); err != nil {
		log.Printf("Error checking user count: %v", err)
		utils.ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		utils.ErrorResponse(w, "Signup is disabled after initial setup", http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		utils.ErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	fmt.Println(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.ErrorResponse(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	result, err := db.ExecContext(r.Context(),
		`INSERT INTO users (username,email,password_hash,role,created_at,updated_at) VALUES (?,?,?,?,?,?)`,
		req.Username, req.Email, string(hashedPassword), "owner", time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		utils.ErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	userId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		utils.ErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:       userId,
		Username: req.Username,
		Email:    req.Email,
		Role:     "owner",
	}

	token, err := middleware.GenerateJWT(userId, user.Email, user.Role)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		utils.ErrorResponse(w, "Error generating token", http.StatusInternalServerError)
		return
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    user,
		"message": "User created successfully",
		"error":   "",
	})
}
