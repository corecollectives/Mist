package models

import (
	"time"
)

type SystemInfo struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateHistory struct {
	ID                int64      `json:"id"`
	FromVersion       string     `json:"fromVersion"`
	ToVersion         string     `json:"toVersion"`
	Status            string     `json:"status"`
	StartedAt         time.Time  `json:"startedAt"`
	CompletedAt       *time.Time `json:"completedAt"`
	ErrorMessage      *string    `json:"errorMessage"`
	RollbackAvailable bool       `json:"rollbackAvailable"`
	InitiatedBy       *int64     `json:"initiatedBy"`
}

type GithubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	TarballURL  string    `json:"tarball_url"`
	ZipballURL  string    `json:"zipball_url"`
	HTMLURL     string    `json:"html_url"`
}

func GetSystemVersion() (string, error) {
	query := `SELECT value FROM system_info WHERE key = 'version'`
	var version string
	err := db.QueryRow(query).Scan(&version)
	if err != nil {
		return "1.0.0", err
	}
	return version, nil
}

func SetSystemVersion(version string) error {
	query := `
		INSERT INTO system_info (key, value, updated_at) 
		VALUES ('version', ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.Exec(query, version, version)
	return err
}

func GetSystemInfo(key string) (string, error) {
	query := `SELECT value FROM system_info WHERE key = ?`
	var value string
	err := db.QueryRow(query).Scan(&value)
	return value, err
}

func SetSystemInfo(key, value string) error {
	query := `
		INSERT INTO system_info (key, value, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.Exec(query, key, value, value)
	return err
}

func CreateUpdateHistory(fromVersion, toVersion string, initiatedBy int64) (int64, error) {
	query := `
		INSERT INTO update_history (from_version, to_version, status, initiated_by)
		VALUES (?, ?, 'pending', ?)
		RETURNING id
	`
	var id int64
	err := db.QueryRow(query, fromVersion, toVersion, initiatedBy).Scan(&id)
	return id, err
}

func UpdateUpdateHistoryStatus(id int64, status string, errorMsg *string) error {
	var query string
	var err error

	if errorMsg != nil {
		query = `
			UPDATE update_history 
			SET status = ?, error_message = ?, completed_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`
		_, err = db.Exec(query, status, errorMsg, id)
	} else {
		query = `
			UPDATE update_history 
			SET status = ?, completed_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`
		_, err = db.Exec(query, status, id)
	}

	return err
}

func GetUpdateHistory(limit int) ([]UpdateHistory, error) {
	query := `
		SELECT id, from_version, to_version, status, started_at, 
		       completed_at, error_message, rollback_available, initiated_by
		FROM update_history
		ORDER BY started_at DESC
		LIMIT ?
	`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []UpdateHistory
	for rows.Next() {
		var h UpdateHistory
		err := rows.Scan(
			&h.ID, &h.FromVersion, &h.ToVersion, &h.Status, &h.StartedAt,
			&h.CompletedAt, &h.ErrorMessage, &h.RollbackAvailable, &h.InitiatedBy,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func GetLatestUpdateHistory() (*UpdateHistory, error) {
	query := `
		SELECT id, from_version, to_version, status, started_at, 
		       completed_at, error_message, rollback_available, initiated_by
		FROM update_history
		ORDER BY started_at DESC
		LIMIT 1
	`
	var h UpdateHistory
	err := db.QueryRow(query).Scan(
		&h.ID, &h.FromVersion, &h.ToVersion, &h.Status, &h.StartedAt,
		&h.CompletedAt, &h.ErrorMessage, &h.RollbackAvailable, &h.InitiatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
