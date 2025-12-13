package models

import (
	"time"
)

type Session struct {
	ID             string     `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"userId"`
	SessionData    *string    `db:"session_data" json:"sessionData,omitempty"` // JSON
	IPAddress      *string    `db:"ip_address" json:"ipAddress,omitempty"`
	UserAgent      *string    `db:"user_agent" json:"userAgent,omitempty"`
	DeviceType     *string    `db:"device_type" json:"deviceType,omitempty"`
	Browser        *string    `db:"browser" json:"browser,omitempty"`
	OS             *string    `db:"os" json:"os,omitempty"`
	Location       *string    `db:"location" json:"location,omitempty"`
	IsActive       bool       `db:"is_active" json:"isActive"`
	LastActivityAt time.Time  `db:"last_activity_at" json:"lastActivityAt"`
	RevokedAt      *time.Time `db:"revoked_at" json:"revokedAt,omitempty"`
	RevokedReason  *string    `db:"revoked_reason" json:"revokedReason,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	ExpiresAt      time.Time  `db:"expires_at" json:"expiresAt"`
}

func (s *Session) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":             s.ID,
		"userId":         s.UserID,
		"sessionData":    s.SessionData,
		"ipAddress":      s.IPAddress,
		"userAgent":      s.UserAgent,
		"deviceType":     s.DeviceType,
		"browser":        s.Browser,
		"os":             s.OS,
		"location":       s.Location,
		"isActive":       s.IsActive,
		"lastActivityAt": s.LastActivityAt,
		"revokedAt":      s.RevokedAt,
		"revokedReason":  s.RevokedReason,
		"createdAt":      s.CreatedAt,
		"expiresAt":      s.ExpiresAt,
	}
}

func (s *Session) InsertInDB() error {
	query := `
	INSERT INTO sessions (
		id, user_id, session_data, ip_address, user_agent,
		device_type, browser, os, location, expires_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING created_at, last_activity_at
	`
	err := db.QueryRow(query, s.ID, s.UserID, s.SessionData, s.IPAddress, s.UserAgent,
		s.DeviceType, s.Browser, s.OS, s.Location, s.ExpiresAt).Scan(&s.CreatedAt, &s.LastActivityAt)
	return err
}

func GetSessionByID(sessionID string) (*Session, error) {
	var session Session
	query := `
	SELECT id, user_id, session_data, ip_address, user_agent,
	       device_type, browser, os, location, is_active,
	       last_activity_at, revoked_at, revoked_reason,
	       created_at, expires_at
	FROM sessions
	WHERE id = ? AND is_active = 1
	`
	err := db.QueryRow(query, sessionID).Scan(
		&session.ID, &session.UserID, &session.SessionData, &session.IPAddress,
		&session.UserAgent, &session.DeviceType, &session.Browser, &session.OS,
		&session.Location, &session.IsActive, &session.LastActivityAt,
		&session.RevokedAt, &session.RevokedReason, &session.CreatedAt, &session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func GetSessionsByUserID(userID int64) ([]Session, error) {
	var sessions []Session
	query := `
	SELECT id, user_id, session_data, ip_address, user_agent,
	       device_type, browser, os, location, is_active,
	       last_activity_at, revoked_at, revoked_reason,
	       created_at, expires_at
	FROM sessions
	WHERE user_id = ? AND is_active = 1
	ORDER BY last_activity_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err := rows.Scan(
			&session.ID, &session.UserID, &session.SessionData, &session.IPAddress,
			&session.UserAgent, &session.DeviceType, &session.Browser, &session.OS,
			&session.Location, &session.IsActive, &session.LastActivityAt,
			&session.RevokedAt, &session.RevokedReason, &session.CreatedAt, &session.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

func (s *Session) UpdateActivity() error {
	query := `
	UPDATE sessions
	SET last_activity_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := db.Exec(query, s.ID)
	return err
}

func (s *Session) Revoke(reason string) error {
	query := `
	UPDATE sessions
	SET is_active = 0,
	    revoked_at = CURRENT_TIMESTAMP,
	    revoked_reason = ?
	WHERE id = ?
	`
	_, err := db.Exec(query, reason, s.ID)
	return err
}

func RevokeAllUserSessions(userID int64, reason string) error {
	query := `
	UPDATE sessions
	SET is_active = 0,
	    revoked_at = CURRENT_TIMESTAMP,
	    revoked_reason = ?
	WHERE user_id = ? AND is_active = 1
	`
	_, err := db.Exec(query, reason, userID)
	return err
}

func DeleteExpiredSessions() error {
	query := `
	DELETE FROM sessions
	WHERE expires_at < CURRENT_TIMESTAMP
	`
	_, err := db.Exec(query)
	return err
}
