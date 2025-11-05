package queuehandlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/utils"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//below 2 things are hardcoded currently add the getcommithash and msg function along with error on lhs of :=
	commitHash := 2332
	// if err != nil {
	// 	http.Error(w, "failed to get commit hash", http.StatusInternalServerError)
	// 	return
	// }
	commitMessage := "hello hi"
	// if err != nil {
	// 	http.Error(w, "failed to get commit message", http.StatusInternalServerError)
	// 	return
	// }
	deploymentId := utils.GenerateRandomId()
	result, err := q.DB.Exec(
		`INSERT INTO deployments (id,app_id, commit_hash, commit_message, status) VALUES (?,?, ?, ?, 'pending')`,
		deploymentId, req.AppId, commitHash, commitMessage,
	)
	if err != nil {
		http.Error(w, "failed to insert deployment", http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "failed to get inserted id", http.StatusInternalServerError)
		return
	}

	if err := q.Queue.AddJob(int64(id)); err != nil {
		http.Error(w, "failed to add job to queue", http.StatusInternalServerError)
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
		http.Error(w, "failed to fetch deployment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deployment)

}
