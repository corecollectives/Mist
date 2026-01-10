package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID           int64     `json:"id"`
	UserID       *int64    `json:"userId"`
	Username     *string   `json:"username"`
	Email        *string   `json:"email"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resourceType"`
	ResourceID   *int64    `json:"resourceId"`
	ResourceName *string   `json:"resourceName"`
	Details      *string   `json:"details"`
	IPAddress    *string   `json:"ipAddress"`
	UserAgent    *string   `json:"userAgent"`
	TriggerType  string    `json:"triggerType"` // "user", "webhook", "system"
	CreatedAt    time.Time `json:"createdAt"`
}

type AuditLogDetails struct {
	Before interface{} `json:"before,omitempty"`
	After  interface{} `json:"after,omitempty"`
	Reason string      `json:"reason,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
}

func (a *AuditLog) Create() error {
	query := `
		INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	return db.QueryRow(query, a.UserID, a.Action, a.ResourceType, a.ResourceID, a.Details).
		Scan(&a.ID, &a.CreatedAt)
}

func GetAllAuditLogs(limit, offset int) ([]AuditLog, error) {
	query := `
		SELECT 
			al.id, 
			al.user_id, 
			u.username,
			u.email,
			al.action, 
			al.resource_type, 
			al.resource_id, 
			al.details, 
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		ORDER BY al.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Username,
			&log.Email,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if log.UserID == nil {
			log.TriggerType = "system"
			if log.Details != nil && (*log.Details != "") {
				var detailsMap map[string]interface{}
				if err := json.Unmarshal([]byte(*log.Details), &detailsMap); err == nil {
					if triggerType, ok := detailsMap["trigger_type"].(string); ok {
						log.TriggerType = triggerType
					}
				}
			}
		} else {
			log.TriggerType = "user"
		}

		logs = append(logs, log)
	}
	return logs, nil
}

func GetAuditLogsByUser(userID int64, limit, offset int) ([]AuditLog, error) {
	query := `
		SELECT 
			al.id, 
			al.user_id, 
			u.username,
			u.email,
			al.action, 
			al.resource_type, 
			al.resource_id, SELECT
			al.details, 
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.user_id = ?
		ORDER BY al.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Username,
			&log.Email,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		log.TriggerType = "user"
		logs = append(logs, log)
	}
	return logs, nil
}

func GetAuditLogsByResource(resourceType string, resourceID int64, limit, offset int) ([]AuditLog, error) {
	query := `
		SELECT 
			al.id, 
			al.user_id, 
			u.username,
			u.email,
			al.action, 
			al.resource_type, 
			al.resource_id, 
			al.details, 
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.resource_type = ? AND al.resource_id = ?
		ORDER BY al.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, resourceType, resourceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Username,
			&log.Email,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if log.UserID == nil {
			log.TriggerType = "system"
			if log.Details != nil && (*log.Details != "") {
				var detailsMap map[string]interface{}
				if err := json.Unmarshal([]byte(*log.Details), &detailsMap); err == nil {
					if triggerType, ok := detailsMap["trigger_type"].(string); ok {
						log.TriggerType = triggerType
					}
				}
			}
		} else {
			log.TriggerType = "user"
		}

		logs = append(logs, log)
	}
	return logs, nil
}

func GetAuditLogsCount() (int, error) {
	query := `SELECT COUNT(*) FROM audit_logs`
	var count int
	err := db.QueryRow(query).Scan(&count)
	return count, err
}

func LogAudit(userID *int64, action, resourceType string, resourceID *int64, details interface{}) error {
	var detailsJSON *string
	if details != nil {
		jsonBytes, err := json.Marshal(details)
		if err != nil {
			return err
		}
		jsonStr := string(jsonBytes)
		detailsJSON = &jsonStr
	}

	log := &AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      detailsJSON,
	}
	return log.Create()
}

func LogUserAudit(userID int64, action, resourceType string, resourceID *int64, details interface{}) error {
	return LogAudit(&userID, action, resourceType, resourceID, details)
}

func LogWebhookAudit(action, resourceType string, resourceID *int64, details map[string]interface{}) error {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["trigger_type"] = "webhook"
	return LogAudit(nil, action, resourceType, resourceID, details)
}

func LogSystemAudit(action, resourceType string, resourceID *int64, details interface{}) error {
	detailsMap := make(map[string]interface{})
	detailsMap["trigger_type"] = "system"
	if details != nil {
		detailsMap["data"] = details
	}
	return LogAudit(nil, action, resourceType, resourceID, detailsMap)
}

func GetAuditLogsByResourceType(resourceType string, limit, offset int) ([]AuditLog, error) {
	query := `
		SELECT 
			al.id, 
			al.user_id, 
			u.username,
			u.email,
			al.action, 
			al.resource_type, 
			al.resource_id, 
			al.details, 
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.resource_type = ?
		ORDER BY al.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, resourceType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Username,
			&log.Email,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if log.UserID == nil {
			log.TriggerType = "system"
			if log.Details != nil && (*log.Details != "") {
				var detailsMap map[string]interface{}
				if err := json.Unmarshal([]byte(*log.Details), &detailsMap); err == nil {
					if triggerType, ok := detailsMap["trigger_type"].(string); ok {
						log.TriggerType = triggerType
					}
				}
			}
		} else {
			log.TriggerType = "user"
		}

		logs = append(logs, log)
	}
	return logs, nil
}

func GetAuditLogByID(id int64) (*AuditLog, error) {
	query := `
		SELECT 
			al.id, 
			al.user_id, 
			u.username,
			u.email,
			al.action, 
			al.resource_type, 
			al.resource_id, 
			al.details, 
			al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.id = ?
	`

	var log AuditLog
	err := db.QueryRow(query, id).Scan(
		&log.ID,
		&log.UserID,
		&log.Username,
		&log.Email,
		&log.Action,
		&log.ResourceType,
		&log.ResourceID,
		&log.Details,
		&log.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if log.UserID == nil {
		log.TriggerType = "system"
		if log.Details != nil && (*log.Details != "") {
			var detailsMap map[string]interface{}
			if err := json.Unmarshal([]byte(*log.Details), &detailsMap); err == nil {
				if triggerType, ok := detailsMap["trigger_type"].(string); ok {
					log.TriggerType = triggerType
				}
			}
		}
	} else {
		log.TriggerType = "user"
	}

	return &log, nil
}
