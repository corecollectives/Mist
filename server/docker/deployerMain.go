package docker

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
)

func DeployerMain(Id int64, db *sql.DB, logFile *os.File) (string, error) {
	dep, err := LoadDeployment(Id, db)
	if err != nil {
		return "", err
	}

	var appId int64
	err = db.QueryRow("SELECT app_id FROM deployments WHERE id = ?", Id).Scan(&appId)
	if err != nil {
		return "", err
	}

	var app models.App
	err = db.QueryRow(`
		SELECT 
			id,
			project_id,
			created_by,
			name,
			description,
			git_repository,
			git_branch,
			deployment_strategy,
			root_directory,
			healthcheck_interval,
			status,
			created_at,
			updated_at
		FROM apps
		WHERE id = ?
	`, appId).Scan(
		&app.ID,
		&app.ProjectID,
		&app.CreatedBy,
		&app.Name,
		&app.Description,
		&app.GitRepository,
		&app.GitBranch,
		&app.DeploymentStrategy,
		&app.RootDirectory,
		&app.HealthcheckInterval,
		&app.Status,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		return "", err
	}

	appContextPath := constants.Constants["RootPath"] + "/" + fmt.Sprintf("projects/%d/apps/%s", app.ProjectID, app.Name)
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", app.ID)

	err = DeployApp(dep, appContextPath, imageTag, containerName, app.ID, app.CreatedBy, db, logFile)
	if err != nil {
		println("Deployment error:", err.Error())
		dep.Status = "failed"
		UpdateDeployment(dep, db)
	}

	return "Deployment started", nil
}
