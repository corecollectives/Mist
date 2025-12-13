package deployments

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
	"github.com/rs/zerolog/log"
)

func AddDeployHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AppId int `json:"appId"`
	}
	queue := queue.GetQueue()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.SendResponse(w, http.StatusBadRequest, false, nil, "invalid request body", err.Error())
		return
	}

	user, ok := middleware.GetUser(r)
	if !ok {
		handlers.SendResponse(w, http.StatusForbidden, false, nil, "unauthorized", "")
		return
	}
	userId := int64(user.ID)
	commit, err := github.GetLatestCommit(int64(req.AppId), userId)
	if err != nil {
		log.Error().Err(err).Int("app_id", req.AppId).Msg("Error getting latest commit")
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to get latest commit", err.Error())
		return
	}
	commitHash := commit.SHA

	commitMessage := commit.Message
	deployment := models.Deployment{
		AppID:         int64(req.AppId),
		CommitHash:    commitHash,
		CommitMessage: &commitMessage,
		Status:        models.DeploymentStatusPending,
	}
	err = deployment.CreateDeployment()

	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to insert deployment", err.Error())
		return
	}

	if err := queue.AddJob(int64(deployment.ID)); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to add job to queue", err.Error())
		return
	}

	log.Info().Int64("deployment_id", deployment.ID).Int("app_id", req.AppId).Msg("Deployment added to queue")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployment)

}
