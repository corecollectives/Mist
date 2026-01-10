package models

import "time"

type Cron struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"app_id"`
	Name      string    `json:"name"`
	Schedule  string    `json:"schedule"`
	Command   string    `json:"command"`
	LastRun   time.Time `json:"last_run"`
	NextRun   time.Time `json:"next_run"`
	Enable    bool      `json:"enable"`
	CreatedAt time.Time `json:"created_at"`
}
