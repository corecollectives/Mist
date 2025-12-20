package settings

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

func GetSystemSettings(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can view system settings", "Forbidden")
		return
	}

	settings, err := models.GetSystemSettings()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve system settings", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, settings, "System settings retrieved successfully", "")
}

func UpdateSystemSettings(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}

	role, err := models.GetUserRole(userInfo.ID)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify user role", err.Error())
		return
	}
	if role != "owner" {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "Only owners can update system settings", "Forbidden")
		return
	}

	var req struct {
		WildcardDomain *string `json:"wildcardDomain"`
		MistAppName    string  `json:"mistAppName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body", "Could not parse JSON")
		return
	}

	if req.MistAppName == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Mist app name is required", "Missing fields")
		return
	}

	if req.WildcardDomain != nil && *req.WildcardDomain == "" {
		req.WildcardDomain = nil
	}

	settings, err := models.UpdateSystemSettings(req.WildcardDomain, req.MistAppName)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update system settings", err.Error())
		return
	}

	if err := utils.GenerateDynamicConfig(settings.WildcardDomain, settings.MistAppName); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to generate Traefik configuration", err.Error())
		return
	}

	dummyID := int64(1)
	models.LogUserAudit(userInfo.ID, "update", "system_settings", &dummyID, map[string]any{
		"wildcardDomain": req.WildcardDomain,
		"mistAppName":    req.MistAppName,
	})

	handlers.SendResponse(w, http.StatusOK, true, settings, "System settings updated successfully", "")
}
