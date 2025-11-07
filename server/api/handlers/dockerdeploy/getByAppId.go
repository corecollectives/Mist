package dockerdeploy

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
)

func (d *Deployer) GetByApplicationID(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	var req struct {
		AppID int64 `json:"appId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.AppID == 0 {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "App ID is required", "Missing fields")
		return
	}

	app, err := models.GetApplicationByID(req.AppID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get application", err.Error())
		return
	}

	deployments, err := models.GetDeploymentsByAppID(app.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get deployments", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, deployments, "Deployments retrieved successfully", "")
}
