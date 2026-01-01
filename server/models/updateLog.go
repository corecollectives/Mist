package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type UpdateLog struct {
	ID           int64      `json:"id"`
	VersionFrom  string     `json:"versionFrom"`
	VersionTo    string     `json:"versionTo"`
	Status       string     `json:"status"` // in_progress, success, failed
	Logs         string     `json:"logs"`
	ErrorMessage *string    `json:"errorMessage"`
	StartedBy    int64      `json:"startedBy"`
	StartedAt    time.Time  `json:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt"`
	Username     string     `json:"username"`
}

func CreateUpdateLog(versionFrom, versionTo string, startedBy int64) (*UpdateLog, error) {
	query := `
		INSERT INTO update_logs (version_from, version_to, status, logs, started_by, started_at)
		VALUES (?, ?, 'in_progress', '', ?, CURRENT_TIMESTAMP)
		RETURNING id, version_from, version_to, status, logs, error_message, started_by, started_at, completed_at
	`

	updateLog := &UpdateLog{}
	err := db.QueryRow(query, versionFrom, versionTo, startedBy).Scan(
		&updateLog.ID,
		&updateLog.VersionFrom,
		&updateLog.VersionTo,
		&updateLog.Status,
		&updateLog.Logs,
		&updateLog.ErrorMessage,
		&updateLog.StartedBy,
		&updateLog.StartedAt,
		&updateLog.CompletedAt,
	)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create update log")
		return nil, err
	}

	log.Info().
		Int64("update_log_id", updateLog.ID).
		Str("from", versionFrom).
		Str("to", versionTo).
		Int64("started_by", startedBy).
		Msg("Update log created")

	return updateLog, nil
}

func UpdateUpdateLogStatus(id int64, status string, logs string, errorMessage *string) error {
	query := `
		UPDATE update_logs
		SET status = ?, logs = ?, error_message = ?, completed_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := db.Exec(query, status, logs, errorMessage, id)
	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to update log status")
		return err
	}

	log.Info().
		Int64("update_log_id", id).
		Str("status", status).
		Msg("Update log status updated")

	return nil
}

func AppendUpdateLog(id int64, logLine string) error {
	query := `
		UPDATE update_logs
		SET logs = logs || ?
		WHERE id = ?
	`

	_, err := db.Exec(query, logLine+"\n", id)
	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to append log line")
		return err
	}

	return nil
}

func GetUpdateLogs(limit int) ([]UpdateLog, error) {
	query := `
		SELECT 
			ul.id, ul.version_from, ul.version_to, ul.status, 
			ul.logs, ul.error_message, ul.started_by, ul.started_at, 
			ul.completed_at, u.username
		FROM update_logs ul
		LEFT JOIN users u ON ul.started_by = u.id
		ORDER BY ul.started_at DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query update logs")
		return nil, err
	}
	defer rows.Close()

	var logs []UpdateLog
	for rows.Next() {
		var updateLog UpdateLog
		err := rows.Scan(
			&updateLog.ID,
			&updateLog.VersionFrom,
			&updateLog.VersionTo,
			&updateLog.Status,
			&updateLog.Logs,
			&updateLog.ErrorMessage,
			&updateLog.StartedBy,
			&updateLog.StartedAt,
			&updateLog.CompletedAt,
			&updateLog.Username,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan update log row")
			return nil, err
		}
		logs = append(logs, updateLog)
	}

	return logs, nil
}

func GetUpdateLogByID(id int64) (*UpdateLog, error) {
	query := `
		SELECT 
			ul.id, ul.version_from, ul.version_to, ul.status, 
			ul.logs, ul.error_message, ul.started_by, ul.started_at, 
			ul.completed_at, u.username
		FROM update_logs ul
		LEFT JOIN users u ON ul.started_by = u.id
		WHERE ul.id = ?
	`

	updateLog := &UpdateLog{}
	err := db.QueryRow(query, id).Scan(
		&updateLog.ID,
		&updateLog.VersionFrom,
		&updateLog.VersionTo,
		&updateLog.Status,
		&updateLog.Logs,
		&updateLog.ErrorMessage,
		&updateLog.StartedBy,
		&updateLog.StartedAt,
		&updateLog.CompletedAt,
		&updateLog.Username,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to get update log by ID")
		return nil, err
	}

	return updateLog, nil
}

func GetUpdateLogsAsString() (string, error) {
	logs, err := GetUpdateLogs(10)
	if err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "No update history available", nil
	}

	var builder strings.Builder
	builder.WriteString("Recent Update History:\n")
	builder.WriteString("======================\n\n")

	for _, log := range logs {
		builder.WriteString("Version: ")
		builder.WriteString(log.VersionFrom)
		builder.WriteString(" → ")
		builder.WriteString(log.VersionTo)
		builder.WriteString("\n")
		builder.WriteString("Status: ")
		builder.WriteString(log.Status)
		builder.WriteString("\n")
		builder.WriteString("Started: ")
		builder.WriteString(log.StartedAt.Format("2006-01-02 15:04:05"))
		builder.WriteString(" by ")
		builder.WriteString(log.Username)
		builder.WriteString("\n")
		if log.CompletedAt != nil {
			builder.WriteString("Completed: ")
			builder.WriteString(log.CompletedAt.Format("2006-01-02 15:04:05"))
			builder.WriteString("\n")
		}
		if log.ErrorMessage != nil && *log.ErrorMessage != "" {
			builder.WriteString("Error: ")
			builder.WriteString(*log.ErrorMessage)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
