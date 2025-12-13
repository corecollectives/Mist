package docker

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

func DeployerMain(Id int64, db *sql.DB, logFile *os.File, logger *utils.DeploymentLogger) (string, error) {
	dep, err := LoadDeployment(Id, db)
	if err != nil {
		logger.Error(err, "Failed to load deployment")
		return "", fmt.Errorf("failed to load deployment: %w", err)
	}

	var appId int64
	err = db.QueryRow("SELECT app_id FROM deployments WHERE id = ?", Id).Scan(&appId)
	if err != nil {
		logger.Error(err, "Failed to get app_id")
		return "", fmt.Errorf("failed to get app_id: %w", err)
	}

	appPtr, err := models.GetApplicationByID(appId)
	if err != nil {
		logger.Error(err, "Failed to get app details")
		return "", fmt.Errorf("failed to get app details: %w", err)
	}
	app := *appPtr

	logger.InfoWithFields("App details loaded", map[string]interface{}{
		"app_name":   app.Name,
		"project_id": app.ProjectID,
		"app_type":   app.AppType,
	})

	appContextPath := constants.Constants["RootPath"] + "/" + fmt.Sprintf("projects/%d/apps/%s", app.ProjectID, app.Name)
	imageTag := dep.CommitHash
	containerName := fmt.Sprintf("app-%d", app.ID)

	err = DeployApp(dep, &app, appContextPath, imageTag, containerName, db, logFile, logger)
	if err != nil {
		logger.Error(err, "DeployApp failed")
		dep.Status = "failed"
		dep.Stage = "failed"
		errMsg := err.Error()
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		return "", err
	}

	logger.Info("Deployment completed successfully")
	return "Deployment started", nil
}
