package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type Deployment struct {
	ID            int64      `json:"id"`
	AppID         int64      `json:"app_id"`
	CommitHash    string     `json:"commit_hash"`
	CommitMessage string     `json:"commit_message"`
	TriggeredBy   *int64     `json:"triggered_by"`
	Logs          *string    `json:"logs"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	FinishedAt    *time.Time `json:"finished_at"`
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

func (d *Deployment) CreateDeployment() error {
	id := utils.GenerateRandomId()
	query := `INSERT INTO deployments (id, app_id, commit_hash, commit_message,status)
			  VALUES (?, ?, ?, ?, 'pending')`
	result, err := db.Exec(query, id, d.AppID, d.CommitHash, d.CommitMessage)
	if err != nil {
		return err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

func GetDeploymentByID(depID int64) (*Deployment, error) {
	query := `SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at
			  FROM deployments
			  WHERE id = ?`

	row := db.QueryRow(query, depID)
	var d Deployment
	if err := row.Scan(&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.TriggeredBy, &d.Logs, &d.Status, &d.CreatedAt, &d.FinishedAt); err != nil {
		return nil, err
	}

	return &d, nil
}

func GetCommitHashByDeploymentID(depID int64) (string, error) {
	query := `SELECT commit_hash
			  FROM deployments
			  WHERE id = ?`

	var commitHash string
	err := db.QueryRow(query, depID).Scan(&commitHash)
	if err != nil {
		return "", err
	}

	return commitHash, nil
}
