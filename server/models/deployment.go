package models

import "time"



type Deployment struct {
	ID         int64     `json:"id"`
	AppID      int64     `json:"appId"`
	CommitSHA  string    `json:"commitSha"`
	Logs       string    `json:"logs"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	FinishedAt time.Time `json:"finishedAt"`
}