package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"golang.org/x/crypto/bcrypt"
)

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	db := h.DB
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		ErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		ErrorResponse(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	now := time.Now()

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_, err = db.ExecContext(r.Context(), `INSERT INTO users (username,email,password_hash,role,created_at,updated_at) VALUES (?,?,?,?,?,?)`,
		user.Username, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		ErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})

}

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

func (h *Handler) DoesUserExist(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("email")))
	if email == "" {
		ErrorResponse(w, "Email is required", http.StatusBadRequest)
		return
	}
	db := h.DB
	var exists bool
	err := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email)=?)`, email).Scan(&exists)
	if err != nil {
		log.Printf("DB query error: %v", err)
		ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}
