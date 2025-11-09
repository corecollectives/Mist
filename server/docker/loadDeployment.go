package docker

import (
	"database/sql"

	"github.com/corecollectives/mist/models"
)

func LoadDeployment(depId int64, db *sql.DB) (*models.Deployment, error) {
	row := db.QueryRow("SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at FROM deployments WHERE id = ?", depId)
	dep := &models.Deployment{}
	var triggeredBy sql.NullInt64
	var finishedAt sql.NullTime
	var logs sql.NullString
	err := row.Scan(&dep.ID, &dep.AppID, &dep.CommitHash, &dep.CommitMessage, &triggeredBy, &logs, &dep.Status, &dep.CreatedAt, &finishedAt)
	if err != nil {
		return nil, err
	}
	if triggeredBy.Valid {
		dep.TriggeredBy = &triggeredBy.Int64
	}
	if finishedAt.Valid {
		dep.FinishedAt = &finishedAt.Time
	}
	return dep, nil

}
