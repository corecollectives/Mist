package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
