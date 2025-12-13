package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type ApiToken struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"userId"`
	Name        string     `db:"name" json:"name"`
	TokenHash   string     `db:"token_hash" json:"-"` // Never expose in JSON
	TokenPrefix string     `db:"token_prefix" json:"tokenPrefix"`
	Scopes      *string    `db:"scopes" json:"scopes,omitempty"` // JSON array
	LastUsedAt  *time.Time `db:"last_used_at" json:"lastUsedAt,omitempty"`
	LastUsedIP  *string    `db:"last_used_ip" json:"lastUsedIp,omitempty"`
	UsageCount  int        `db:"usage_count" json:"usageCount"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expiresAt,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	RevokedAt   *time.Time `db:"revoked_at" json:"revokedAt,omitempty"`
}

func (t *ApiToken) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"userId":      t.UserID,
		"name":        t.Name,
		"tokenPrefix": t.TokenPrefix,
		"scopes":      t.Scopes,
		"lastUsedAt":  t.LastUsedAt,
		"lastUsedIp":  t.LastUsedIP,
		"usageCount":  t.UsageCount,
		"expiresAt":   t.ExpiresAt,
		"createdAt":   t.CreatedAt,
		"revokedAt":   t.RevokedAt,
	}
}

func (t *ApiToken) InsertInDB() error {
	id := utils.GenerateRandomId()
	t.ID = id
	query := `
	INSERT INTO api_tokens (
		id, user_id, name, token_hash, token_prefix, scopes, expires_at
	) VALUES (?, ?, ?, ?, ?, ?, ?)
	RETURNING created_at
	`
	err := db.QueryRow(query, t.ID, t.UserID, t.Name, t.TokenHash, t.TokenPrefix, t.Scopes, t.ExpiresAt).Scan(&t.CreatedAt)
	return err
}

func GetApiTokensByUserID(userID int64) ([]ApiToken, error) {
	var tokens []ApiToken
	query := `
	SELECT id, user_id, name, token_hash, token_prefix, scopes,
	       last_used_at, last_used_ip, usage_count, expires_at,
	       created_at, revoked_at
	FROM api_tokens
	WHERE user_id = ? AND revoked_at IS NULL
	ORDER BY created_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var token ApiToken
		err := rows.Scan(
			&token.ID, &token.UserID, &token.Name, &token.TokenHash, &token.TokenPrefix,
			&token.Scopes, &token.LastUsedAt, &token.LastUsedIP, &token.UsageCount,
			&token.ExpiresAt, &token.CreatedAt, &token.RevokedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func GetApiTokenByHash(tokenHash string) (*ApiToken, error) {
	var token ApiToken
	query := `
	SELECT id, user_id, name, token_hash, token_prefix, scopes,
	       last_used_at, last_used_ip, usage_count, expires_at,
	       created_at, revoked_at
	FROM api_tokens
	WHERE token_hash = ? AND revoked_at IS NULL
	`
	err := db.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.Name, &token.TokenHash, &token.TokenPrefix,
		&token.Scopes, &token.LastUsedAt, &token.LastUsedIP, &token.UsageCount,
		&token.ExpiresAt, &token.CreatedAt, &token.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *ApiToken) UpdateUsage(ipAddress string) error {
	query := `
	UPDATE api_tokens
	SET last_used_at = CURRENT_TIMESTAMP,
	    last_used_ip = ?,
	    usage_count = usage_count + 1
	WHERE id = ?
	`
	_, err := db.Exec(query, ipAddress, t.ID)
	return err
}

func (t *ApiToken) Revoke() error {
	query := `
	UPDATE api_tokens
	SET revoked_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := db.Exec(query, t.ID)
	return err
}

func DeleteApiToken(tokenID int64) error {
	query := `DELETE FROM api_tokens WHERE id = ?`
	_, err := db.Exec(query, tokenID)
	return err
}
