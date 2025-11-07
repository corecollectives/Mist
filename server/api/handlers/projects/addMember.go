package projects

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
)

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	if userInfo.Role != "owner" && userInfo.Role != "admin" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Not authorized", "Forbidden")
		return
	}

	var req struct {
		ProjectID int `json:"projectId"`
		UserID    int `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", err.Error())
		return
	}

	if req.ProjectID == 0 || req.UserID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project ID and User ID are required", "Missing fields")
		return
	}

	// check if project exists
	var projectExists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id=$1)", req.ProjectID).Scan(&projectExists)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Internal server error", err.Error())
		return
	}
	if !projectExists {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Project not found", "Invalid project ID")
		return
	}

	// check if user exists
	var userExists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", req.UserID).Scan(&userExists)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Internal server error", err.Error())
		return
	}
	if !userExists {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "User not found", "Invalid user ID")
		return
	}

	// check if user is already a member of the project
	var isMember bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM project_members WHERE project_id=$1 AND user_id=$2)", req.ProjectID, req.UserID).Scan(&isMember)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Internal server error", err.Error())
		return
	}
	if isMember {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "User is already a member of the project", "Duplicate member")
		return
	}

	// add user to project members
	_, err = h.DB.Exec("INSERT INTO project_members(project_id, user_id) VALUES($1, $2)", req.ProjectID, req.UserID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to add member to project", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, nil, "Member added to project successfully", "")
}
