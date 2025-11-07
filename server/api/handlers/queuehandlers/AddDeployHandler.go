package queuehandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/api/middleware"
	"github.com/corecollectives/mist/api/utils"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
)

type QueueHelper struct {
	DB           *sql.DB
	LogDirectory string
	Queue        *queue.Queue
}

func (q *QueueHelper) AddDeployHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AppId int `json:"appId"`
	}
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
	commit, err := github.GetLatestCommit(q.DB, int64(req.AppId), userId)
	if err != nil {
		fmt.Println("Error getting latest commit:", err.Error())
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to get latest commit", err.Error())
		return
	}
	commitHash := commit.SHA

	commitMessage := commit.Message
	deploymentId := utils.GenerateRandomId()
	result, err := q.DB.Exec(
		`INSERT INTO deployments (id,app_id, commit_hash, commit_message, status) VALUES (?,?, ?, ?, 'pending')`,
		deploymentId, req.AppId, commitHash, commitMessage,
	)
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to insert deployment", err.Error())
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to get inserted id", err.Error())
		return
	}

	if err := q.Queue.AddJob(int64(id)); err != nil {
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to add job to queue", err.Error())
		return
	}

	println("Deployment added to queue with ID:", id)

	var deployment models.Deployment
	row := q.DB.QueryRow(`SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at FROM deployments WHERE id = ?`, id)
	err = row.Scan(
		&deployment.ID,
		&deployment.AppID,
		&deployment.CommitHash,
		&deployment.CommitMessage,
		&deployment.TriggeredBy,
		&deployment.Logs,
		&deployment.Status,
		&deployment.CreatedAt,
		&deployment.FinishedAt,
	)
	if err != nil {
		fmt.Println("Error fetching deployment:", err.Error())
		handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "failed to fetch deployment", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployment)

}
