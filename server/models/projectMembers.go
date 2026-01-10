package models

import "time"

type ProjectMembers struct{
	ID		int64  `json:"id"`
	ProjectID	int64  `json:"project_id"`
	USERID		int64  `json:"user_id"`
	AddedAt	time.Time  `json:"added_at"`
}