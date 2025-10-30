package models

import "time"

type Project struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Tags           []string  `json:"tags"`
	OwnerID        int64     `json:"ownerId"`
	Owner          *User     `json:"owner,omitempty"`
	ProjectMembers []User    `json:"projectMembers"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

