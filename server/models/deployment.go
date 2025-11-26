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
	Stage         string     `json:"stage"`
	Progress      int        `json:"progress"`
	ErrorMessage  *string    `json:"error_message"`
	CreatedAt     time.Time  `json:"created_at"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	Duration      *int       `json:"duration"`
}

func GetDeploymentsByAppID(appID int64) ([]Deployment, error) {
	query := `SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, stage, progress, error_message, created_at, started_at, finished_at, duration
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
		if err := rows.Scan(&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.TriggeredBy, &d.Logs, &d.Status, &d.Stage, &d.Progress, &d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration); err != nil {
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
	query := `INSERT INTO deployments (id, app_id, commit_hash, commit_message, status, stage, progress)
			  VALUES (?, ?, ?, ?, 'pending', 'pending', 0)`
	result, err := db.Exec(query, id, d.AppID, d.CommitHash, d.CommitMessage)
	if err != nil {
		return err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	d.Status = "pending"
	d.Stage = "pending"
	d.Progress = 0
	return nil
}

func GetDeploymentByID(depID int64) (*Deployment, error) {
	query := `SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, stage, progress, error_message, created_at, started_at, finished_at, duration
			  FROM deployments
			  WHERE id = ?`

	row := db.QueryRow(query, depID)
	var d Deployment
	if err := row.Scan(&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.TriggeredBy, &d.Logs, &d.Status, &d.Stage, &d.Progress, &d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration); err != nil {
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

// UpdateDeploymentStatus updates the deployment status, stage, progress and error message
func UpdateDeploymentStatus(depID int64, status, stage string, progress int, errorMsg *string) error {
	query := `UPDATE deployments 
			  SET status = ?, stage = ?, progress = ?, error_message = ?, 
			      finished_at = CASE WHEN ? IN ('success', 'failed') THEN ? ELSE finished_at END,
			      duration = CASE WHEN ? IN ('success', 'failed') THEN TIMESTAMPDIFF(SECOND, started_at, ?) ELSE duration END
			  WHERE id = ?`

	now := time.Now()
	_, err := db.Exec(query, status, stage, progress, errorMsg, status, now, status, now, depID)
	return err
}

// MarkDeploymentStarted marks the deployment as started
func MarkDeploymentStarted(depID int64) error {
	query := `UPDATE deployments SET started_at = ? WHERE id = ?`
	_, err := db.Exec(query, time.Now(), depID)
	return err
}
