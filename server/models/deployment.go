package models

import (
	"database/sql"
	"time"
)

type Deployment struct {
	ID            int64          `json:"id"`
	AppID         int64          `json:"app_id"`
	CommitHash    string         `json:"commit_hash"`
	CommitMessage string         `json:"commit_message"`
	TriggeredBy   *int64         `json:"triggered_by"`
	Logs          sql.NullString `json:"logs"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	FinishedAt    *time.Time     `json:"finished_at"`
}
