package models

import "time"



type EnvVariable struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"appId"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}