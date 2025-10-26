package utils

import (
	"database/sql"
	"fmt"
)

func IsSetupRequired(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check users count: %w", err)
	}
	return count == 0, nil
}
