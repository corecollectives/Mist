package applications

import (
	"net/http"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
)

func StopContainerHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to control this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)

	if err := docker.StopContainer(containerName); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to stop container", err.Error())
		return
	}

	app.Status = models.StatusStopped
	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app status", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "stop", "container", &appId, map[string]interface{}{
		"app_name":       app.Name,
		"container_name": containerName,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"message": "Container stopped successfully",
	}, "Container stopped successfully", "")
}

func StartContainerHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to control this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)

	if err := docker.StartContainer(containerName); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to start container", err.Error())
		return
	}

	app.Status = models.StatusRunning
	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app status", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "start", "container", &appId, map[string]interface{}{
		"app_name":       app.Name,
		"container_name": containerName,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"message": "Container started successfully",
	}, "Container started successfully", "")
}

func RestartContainerHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to control this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)

	if err := docker.RestartContainer(containerName); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to restart container", err.Error())
		return
	}

	app.Status = models.StatusRunning
	if err := app.UpdateApplication(); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update app status", err.Error())
		return
	}

	models.LogUserAudit(userInfo.ID, "restart", "container", &appId, map[string]interface{}{
		"app_name":       app.Name,
		"container_name": containerName,
	})

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"message": "Container restarted successfully",
	}, "Container restarted successfully", "")
}

func GetContainerStatusHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to view this application", "Forbidden")
		return
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)

	status, err := docker.GetContainerStatus(containerName)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get container status", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, status, "Container status retrieved successfully", "")
}

func GetContainerLogsHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Not logged in", "Unauthorized")
		return
	}
	appIdStr := r.URL.Query().Get("appId")
	if appIdStr == "" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "appId is required", "")
		return
	}

	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid appId", "")
		return
	}

	isApplicationOwner, err := models.IsUserApplicationOwner(userInfo.ID, appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to verify application ownership", err.Error())
		return
	}
	if !isApplicationOwner {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "You do not have permission to view this application", "Forbidden")
		return
	}

	tailStr := r.URL.Query().Get("tail")
	tail := 100
	if tailStr != "" {
		if parsedTail, err := strconv.Atoi(tailStr); err == nil && parsedTail > 0 {
			tail = parsedTail
		}
	}

	app, err := models.GetApplicationByID(appId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "Application not found", err.Error())
		return
	}

	containerName := docker.GetContainerName(app.Name, appId)

	logs, err := docker.GetContainerLogs(containerName, tail)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to get container logs", err.Error())
		return
	}

	handlers.SendResponse(w, http.StatusOK, true, map[string]any{
		"logs": logs,
	}, "Container logs retrieved successfully", "")
}
