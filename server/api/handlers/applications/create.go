package applications

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func CreateApplication(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ProjectID   int64  `json:"projectId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.Name == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Application name is required", "Missing fields")
		return
	}

	if req.ProjectID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Project ID is required", "Missing fields")
		return
	}

	isUserMember, err := models.HasUserAccessToProject(userInfo.ID, req.ProjectID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify project access", err.Error())
		return
	}
	if !isUserMember {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have access to this project", "Forbidden")
		return
	}

	app := models.App{
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		ProjectID:   req.ProjectID,
		CreatedBy:   userInfo.ID,
	}
	if err := app.InsertInDB(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create application", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, app.ToJson(), "Application created successfully", "")

}
