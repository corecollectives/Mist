package queue

import "database/sql"

func InitQueue(db *sql.DB) *Queue {
	q := NewQueue(5, db)
	return q
}
