package models

import "time"

type Registries struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"projectId"`
	RegistryURL string    `json:"registryUrl"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"createdAt"`
}
