package models

import "time"

type GitHubApp struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	AppID               int64     `json:"app_id"`
	ClientID            string    `json:"client_id"`
	ClientSecret        string    `json:"-"`
	WebhookSecret       string    `json:"-"`
	PrivateKeyEncrypted string    `json:"-"`
	Slug                string    `json:"slug"`
	CreatedAt           time.Time `json:"created_at"`
}
