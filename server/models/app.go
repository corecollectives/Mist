package models

import "time"





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