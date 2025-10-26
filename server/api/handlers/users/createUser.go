package users

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userData, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	if userData.Role != "admin" && userData.Role != "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Insufficient permissions", "Forbidden")
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "All fields are required", "Missing fields")
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid role", "Role must be one of: admin, user")
		return
	}

	if req.Role == "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Cannot create another owner", "Forbidden")
		return
	}

	user, err := utils.InsertUserInDb(h.DB, req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create user", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusCreated, true, user, "User created successfully", "")
}
