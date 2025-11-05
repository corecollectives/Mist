package models

import "time"

type DeploymentStrategy string
type AppStatus string

const (
	DeploymentAuto   DeploymentStrategy = "auto"
	DeploymentManual DeploymentStrategy = "manual"

	StatusStopped  AppStatus = "stopped"
	StatusRunning  AppStatus = "running"
	StatusError    AppStatus = "error"
	StatusBuilding AppStatus = "building"
)

type App struct {
	ID                  int64              `db:"id" json:"id"`
	ProjectID           int64              `db:"project_id" json:"project_id"`
	CreatedBy           *int64             `db:"created_by" json:"created_by"`
	Name                string             `db:"name" json:"name"`
	Description         *string            `db:"description" json:"description,omitempty"`
	GitProviderID       *int64             `db:"git_provider_id" json:"git_provider_id,omitempty"`
	GitRepository       *string            `db:"git_repository" json:"git_repository,omitempty"`
	GitBranch           *string            `db:"git_branch" json:"git_branch,omitempty"`
	DeploymentStrategy  DeploymentStrategy `db:"deployment_strategy" json:"deployment_strategy"`
	Port                *int               `db:"port" json:"port,omitempty"`
	RootDirectory       *string            `db:"root_directory" json:"root_directory,omitempty"`
	BuildCommand        *string            `db:"build_command" json:"build_command,omitempty"`
	StartCommand        *string            `db:"start_command" json:"start_command,omitempty"`
	DockerfilePath      *string            `db:"dockerfile_path" json:"dockerfile_path,omitempty"`
	HealthcheckPath     *string            `db:"healthcheck_path" json:"healthcheck_path,omitempty"`
	HealthcheckInterval int                `db:"healthcheck_interval" json:"healthcheck_interval"`
	Status              AppStatus          `db:"status" json:"status"`
	CreatedAt           time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `db:"updated_at" json:"updated_at"`
}
