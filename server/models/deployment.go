package models

import (
	"database/sql"
	"time"
)

type Deployment struct {
	ID            int64          `json:"id"`
	AppID         int64          `json:"app_id"`
	CommitHash    string         `json:"commit_hash"`
	CommitMessage string         `json:"commit_message"`
	TriggeredBy   *int64         `json:"triggered_by"`
	Logs          sql.NullString `json:"logs"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	FinishedAt    *time.Time     `json:"finished_at"`
}

func GetDeploymentsByAppID(appID int64) ([]Deployment, error) {
	query := `SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at
			  FROM deployments
			  WHERE app_id = ?
			  ORDER BY created_at DESC`

	rows, err := db.Query(query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []Deployment
	for rows.Next() {
		var d Deployment
		if err := rows.Scan(&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.TriggeredBy, &d.Logs, &d.Status, &d.CreatedAt, &d.FinishedAt); err != nil {
			return nil, err
		}
		deployments = append(deployments, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deployments, nil
}
