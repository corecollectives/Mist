package store

import "database/sql"

func InitStore(db *sql.DB) error {
	return SetSetupRequired(db)
}
