package dockerdeploy

import (
	"fmt"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
)

func (d *Deployer) DeployerMain(Id int64) (string, error) {

	// var id int64
	// err := d.DB.QueryRow("SELECT id FROM deployments WHERE id = ?", Id).Scan(&id)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("No deployment found with id:", Id)
	// } else if err != nil {
	// 	fmt.Println("Query error:", err)
	// } else {
	// 	fmt.Println("Deployment found with id:", id)
	// }
	// fmt.Println("Rows: ", rows)
	fmt.Println("Deploying deployment", Id)
	dep, err := d.loadDeployment(Id)
	// if err != nil {
	// 	http.Error(w, "Deployment not found", 404)
	// 	return
	// }

	var appId int64
	err = d.DB.QueryRow("SELECT app_id FROM deployments WHERE id = ?", Id).Scan(&appId)
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

	fmt.Println("App loaded: ", app)

	if err != nil {
		fmt.Printf("Error loading deployment: %v\n", err)

		return "", err
	}

	fmt.Printf("Loaded deployment: %+v\n", dep)
	appContextPath := constants.Constants["RootPath"] + "/" + fmt.Sprintf("projects/%d/apps/%s", app.ProjectID, app.Name)
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", Id)

	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		return "", err

	}

	go func() {
		err := d.DeployApp(dep, appContextPath, imageTag, containerName)
		if err != nil {
			// http.Error(w, err.Error(), 500)
			fmt.Printf("Error deploying app: %v\n", err)
			dep.Status = "failed"
			d.UpdateDeployment(dep)
		}
	}()

	return "Deployment started", nil
}
