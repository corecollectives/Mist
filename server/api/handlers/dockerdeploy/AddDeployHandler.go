package dockerdeploy

import (
	"encoding/json"
	"net/http"

	"github.com/corecollectives/mist/api/utils"
	"github.com/corecollectives/mist/models"
)

func (d *Deployer) AddDeployHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AppId int `json:"appId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	commitHash, err := GetCommitHash(req.AppId)
	if err != nil {
		http.Error(w, "failed to get commit hash", http.StatusInternalServerError)
		return
	}
	commitMessage, err := GetCommitMessage(req.AppId)
	if err != nil {
		http.Error(w, "failed to get commit message", http.StatusInternalServerError)
		return
	}
	deploymentId := utils.GenerateRandomId()
	result, err := d.DB.Exec(
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

	if err := d.Queue.AddJob(int64(id)); err != nil {
		http.Error(w, "failed to add job to queue", http.StatusInternalServerError)
		return
	}

	var deployment models.Deployment
	row := d.DB.QueryRow(`SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at FROM deployments WHERE id = ?`, id)
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
