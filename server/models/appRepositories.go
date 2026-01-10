package models

import "time"

type AppRepositories struct {
	ID           int64     `json:"id"`
	AppID        int64     `json:"app_id"`
	SourceType   string    `json:"source_type"`
	SourceID     int64     `json:"source_id"`
	RepoFullName string    `json:"repo_full_name"`
	RepoURL      string    `json:"repo_url"`
	Branch       string    `json:"branch"`
	WebhookID    int64     `json:"webhook_id"`
	AutoDeploy   bool      `json:"auto_deploy"`
	LastSyncedAt time.Time `json:"last_synced_at"`
}
