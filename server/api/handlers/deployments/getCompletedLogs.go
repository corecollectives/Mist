package deployments

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

type GetDeploymentLogsResponse struct {
	Deployment *models.Deployment `json:"deployment"`
	Logs       string             `json:"logs"`
}

func GetCompletedDeploymentLogsHandler(w http.ResponseWriter, r *http.Request) {
	depIdstr := r.URL.Query().Get("id")
	depId, err := strconv.ParseInt(depIdstr, 10, 64)
	if err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "invalid deployment id", err.Error())
		return
	}

	dep, err := models.GetDeploymentByID(depId)
	if err != nil {
		handlers.SendResponse(w, http.StatusNotFound, false, nil, "deployment not found", err.Error())
		return
	}

	if dep.Status != "success" && dep.Status != "failed" {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "deployment is still in progress, use WebSocket endpoint", "")
		return
	}

	logPath := docker.GetLogsPath(dep.CommitHash, depId)
	logContent := ""

	if _, err := os.Stat(logPath); err == nil {
		file, err := os.Open(logPath)
		if err != nil {
			log.Error().Err(err).Int64("deployment_id", depId).Msg("Failed to open log file")
		} else {
			defer file.Close()
			content, err := io.ReadAll(file)
			if err != nil {
				log.Error().Err(err).Int64("deployment_id", depId).Msg("Failed to read log file")
			} else {
				logContent = string(content)
			}
		}
	}

	response := GetDeploymentLogsResponse{
		Deployment: dep,
		Logs:       logContent,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "Deployment logs retrieved successfully",
	})
}
