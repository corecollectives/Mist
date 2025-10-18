package models

import "time"

type Domain struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"appId"`
	Domain    string    `json:"domain"`
	SslStatus string    `json:"sslStatus"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
