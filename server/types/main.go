package types

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Project struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	OwnerID        int64     `json:"ownerId"`
	ProjectMembers []User    `json:"projectMembers"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type App struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectID   int64     `json:"projectId"`
	Port        int       `json:"port"`
	GitRepo     string    `json:"gitRepo"`
	RootDir     string    `json:"rootDir"`
	BuildCmd    string    `json:"buildCmd"`
	StartCmd    string    `json:"startCmd"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Deployment struct {
	ID         int64     `json:"id"`
	AppID      int64     `json:"appId"`
	CommitSHA  string    `json:"commitSha"`
	Logs       string    `json:"logs"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	FinishedAt time.Time `json:"finishedAt"`
}

type EnvVariable struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"appId"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Domain struct {
	ID        int64     `json:"id"`
	AppID     int64     `json:"appId"`
	Domain    string    `json:"domain"`
	SslStatus string    `json:"sslStatus"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
