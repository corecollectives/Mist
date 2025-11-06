package dockerdeploy

import (
	"fmt"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
)

func (d *Deployer) DeployerMain(Id int64) (string, error) {
	dep, err := d.loadDeployment(Id)
	if err != nil {
		return "", err
	}

	var appId int64
	err = d.DB.QueryRow("SELECT app_id FROM deployments WHERE id = ?", Id).Scan(&appId)
	if err != nil {
		return "", err
	}

	var app models.App
	err = d.DB.QueryRow(`
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

	err = github.CloneRepo(d.DB, app.ID, app.CreatedBy)
	appContextPath := constants.Constants["RootPath"] + "/" + fmt.Sprintf("projects/%d/apps/%s", app.ProjectID, app.Name)
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", app.ID)

	go func() {
		err := d.DeployApp(dep, appContextPath, imageTag, containerName)
		if err != nil {
			dep.Status = "failed"
			d.UpdateDeployment(dep)
		}
	}()

	return "Deployment started", nil
}
