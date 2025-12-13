package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type DeploymentStatus string

const (
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusBuilding   DeploymentStatus = "building"
	DeploymentStatusDeploying  DeploymentStatus = "deploying"
	DeploymentStatusSuccess    DeploymentStatus = "success"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusStopped    DeploymentStatus = "stopped"
	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)

type Deployment struct {
	ID               int64            `db:"id" json:"id"`
	AppID            int64            `db:"app_id" json:"app_id"`
	CommitHash       string           `db:"commit_hash" json:"commit_hash"`
	CommitMessage    *string          `db:"commit_message" json:"commit_message,omitempty"`
	CommitAuthor     *string          `db:"commit_author" json:"commit_author,omitempty"`
	TriggeredBy      *int64           `db:"triggered_by" json:"triggered_by,omitempty"`
	DeploymentNumber *int             `db:"deployment_number" json:"deployment_number,omitempty"`
	ContainerID      *string          `db:"container_id" json:"container_id,omitempty"`
	ContainerName    *string          `db:"container_name" json:"container_name,omitempty"`
	ImageTag         *string          `db:"image_tag" json:"image_tag,omitempty"`
	Logs             *string          `db:"logs" json:"logs,omitempty"`
	BuildLogsPath    *string          `db:"build_logs_path" json:"build_logs_path,omitempty"`
	Status           DeploymentStatus `db:"status" json:"status"`
	Stage            string           `db:"stage" json:"stage"`
	Progress         int              `db:"progress" json:"progress"`
	ErrorMessage     *string          `db:"error_message" json:"error_message,omitempty"`
	CreatedAt        time.Time        `db:"created_at" json:"created_at"`
	StartedAt        *time.Time       `db:"started_at" json:"started_at,omitempty"`
	FinishedAt       *time.Time       `db:"finished_at" json:"finished_at,omitempty"`
	Duration         *int             `db:"duration" json:"duration,omitempty"`
	IsActive         bool             `db:"is_active" json:"is_active"`
	RolledBackFrom   *int64           `db:"rolled_back_from" json:"rolled_back_from,omitempty"`
}

func (d *Deployment) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":               d.ID,
		"appId":            d.AppID,
		"commitHash":       d.CommitHash,
		"commitMessage":    d.CommitMessage,
		"commitAuthor":     d.CommitAuthor,
		"triggeredBy":      d.TriggeredBy,
		"deploymentNumber": d.DeploymentNumber,
		"containerId":      d.ContainerID,
		"containerName":    d.ContainerName,
		"imageTag":         d.ImageTag,
		"logs":             d.Logs,
		"buildLogsPath":    d.BuildLogsPath,
		"status":           d.Status,
		"stage":            d.Stage,
		"progress":         d.Progress,
		"errorMessage":     d.ErrorMessage,
		"createdAt":        d.CreatedAt,
		"startedAt":        d.StartedAt,
		"finishedAt":       d.FinishedAt,
		"duration":         d.Duration,
		"isActive":         d.IsActive,
		"rolledBackFrom":   d.RolledBackFrom,
	}
}

func GetDeploymentsByAppID(appID int64) ([]Deployment, error) {
	query := `
	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by, 
	       deployment_number, container_id, container_name, image_tag,
	       logs, build_logs_path, status, stage, progress, error_message, 
	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
	FROM deployments
	WHERE app_id = ?
	ORDER BY created_at DESC
	`

	rows, err := db.Query(query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []Deployment
	for rows.Next() {
		var d Deployment
		if err := rows.Scan(
			&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
			&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
			&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
			&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
			&d.IsActive, &d.RolledBackFrom,
		); err != nil {
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
	d.ID = id

	// Get next deployment number for this app
	var maxDeploymentNum int
	err := db.QueryRow(`SELECT COALESCE(MAX(deployment_number), 0) FROM deployments WHERE app_id = ?`, d.AppID).Scan(&maxDeploymentNum)
	if err == nil {
		deploymentNum := maxDeploymentNum + 1
		d.DeploymentNumber = &deploymentNum
	}

	query := `
	INSERT INTO deployments (
		id, app_id, commit_hash, commit_message, commit_author, triggered_by,
		deployment_number, status, stage, progress
	) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', 'pending', 0)
	RETURNING created_at
	`
	err = db.QueryRow(query, d.ID, d.AppID, d.CommitHash, d.CommitMessage,
		d.CommitAuthor, d.TriggeredBy, d.DeploymentNumber).Scan(&d.CreatedAt)
	if err != nil {
		return err
	}

	d.Status = DeploymentStatusPending
	d.Stage = "pending"
	d.Progress = 0
	d.IsActive = false
	return nil
}

func GetDeploymentByID(depID int64) (*Deployment, error) {
	query := `
	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
	       deployment_number, container_id, container_name, image_tag,
	       logs, build_logs_path, status, stage, progress, error_message,
	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
	FROM deployments
	WHERE id = ?
	`

	var d Deployment
	if err := db.QueryRow(query, depID).Scan(
		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
		&d.IsActive, &d.RolledBackFrom,
	); err != nil {
		return nil, err
	}

	return &d, nil
}

func GetActiveDeploymentByAppID(appID int64) (*Deployment, error) {
	query := `
	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
	       deployment_number, container_id, container_name, image_tag,
	       logs, build_logs_path, status, stage, progress, error_message,
	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
	FROM deployments
	WHERE app_id = ? AND is_active = 1
	LIMIT 1
	`

	var d Deployment
	if err := db.QueryRow(query, appID).Scan(
		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
		&d.IsActive, &d.RolledBackFrom,
	); err != nil {
		return nil, err
	}

	return &d, nil
}

func GetCommitHashByDeploymentID(depID int64) (string, error) {
	var commitHash string
	err := db.QueryRow(`SELECT commit_hash FROM deployments WHERE id = ?`, depID).Scan(&commitHash)
	if err != nil {
		return "", err
	}
	return commitHash, nil
}

func UpdateDeploymentStatus(depID int64, status, stage string, progress int, errorMsg *string) error {
	query := `
	UPDATE deployments 
	SET status = ?, stage = ?, progress = ?, error_message = ?, 
	    finished_at = CASE WHEN ? IN ('success', 'failed', 'stopped') THEN CURRENT_TIMESTAMP ELSE finished_at END,
	    duration = CASE WHEN ? IN ('success', 'failed', 'stopped') THEN 
	                  CAST((julianday(CURRENT_TIMESTAMP) - julianday(started_at)) * 86400 AS INTEGER)
	               ELSE duration END
	WHERE id = ?
	`
	_, err := db.Exec(query, status, stage, progress, errorMsg, status, status, depID)
	return err
}

func MarkDeploymentStarted(depID int64) error {
	query := `UPDATE deployments SET started_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, depID)
	return err
}

// MarkDeploymentActive marks a deployment as the active one (deactivates others)
func MarkDeploymentActive(depID int64, appID int64) error {
	// First, deactivate all deployments for this app
	_, err := db.Exec(`UPDATE deployments SET is_active = 0 WHERE app_id = ?`, appID)
	if err != nil {
		return err
	}

	// Then, activate the specified deployment
	_, err = db.Exec(`UPDATE deployments SET is_active = 1 WHERE id = ?`, depID)
	return err
}

// UpdateContainerInfo updates container and image information
func UpdateContainerInfo(depID int64, containerID, containerName, imageTag string) error {
	query := `
	UPDATE deployments 
	SET container_id = ?, container_name = ?, image_tag = ?
	WHERE id = ?
	`
	_, err := db.Exec(query, containerID, containerName, imageTag, depID)
	return err
}
