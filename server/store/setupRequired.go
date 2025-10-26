package store

import (
	"database/sql"
)

var SetupRequired bool = true

func SetSetupRequired(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return err
	}
	SetupRequired = count == 0
	return nil
}
