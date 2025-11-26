package docker

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

func DeployApp(dep *models.Deployment, appContextPath, imageTag, containerName string, appId int64, createdBy int64, db *sql.DB, logfile *os.File, logger *utils.DeploymentLogger) error {

	logger.Info("Starting deployment process")

	// Update to building stage
	dep.Status = "building"
	dep.Stage = "building"
	dep.Progress = 50
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "building", "building", 50, nil)

	logger.Info("Building Docker image")
	if err := BuildImage(imageTag, appContextPath, logfile); err != nil {
		logger.Error(err, "Docker image build failed")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Build failed: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("build image failed: %w", err)
	}

	logger.Info("Docker image built successfully")

	// Update to deploying stage
	dep.Status = "deploying"
	dep.Stage = "deploying"
	dep.Progress = 80
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 80, nil)

	logger.Info("Stopping existing container if exists")
	err := StopRemoveContainer(containerName, logfile)
	if err != nil {
		logger.Error(err, "Failed to stop/remove existing container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to stop/remove container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("stop/remove container failed: %w", err)
	}

	logger.Info("Getting port and domain configuration")
	port, domain, err := GetPortAndDomainFromDeployment(dep.ID, db)
	if err != nil {
		logger.Error(err, "Failed to get port and domain")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to get port/domain: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("get port/domain failed: %w", err)
	}

	logger.InfoWithFields("Running container", map[string]interface{}{
		"domain": domain,
		"port":   port,
	})

	if err := RunContainer(imageTag, containerName, domain, port, logfile); err != nil {
		logger.Error(err, "Failed to run container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to run container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("run container failed: %w", err)
	}

	// Success
	dep.Status = "success"
	dep.Stage = "success"
	dep.Progress = 100
	now := time.Now()
	dep.FinishedAt = &now
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

	logger.InfoWithFields("Deployment succeeded", map[string]interface{}{
		"deployment_id": dep.ID,
		"container":     containerName,
	})

	return nil
}

func UpdateDeployment(dep *models.Deployment, db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE deployments SET status=?, stage=?, progress=?, logs=?, error_message=?, finished_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dep.Status, dep.Stage, dep.Progress, dep.Logs, dep.ErrorMessage, dep.FinishedAt, dep.ID)
	return err
}

func GetLogsPath(commitHash string, depId int64) string {
	return filepath.Join(constants.Constants["LogPath"], commitHash+strconv.FormatInt(depId, 10)+"_build_logs")
}

func GetPortAndDomainFromDeployment(deploymentID int64, db *sql.DB) (int, string, error) {
	appID, err := models.GetAppIDByDeploymentID(deploymentID)
	if err != nil {
		return 0, "", fmt.Errorf("get app ID failed: %w", err)
	}
	var port *int
	var domain string = "something.com"
	err = db.QueryRow(
		"SELECT port FROM apps WHERE id = ?",
		appID,
	).Scan(&port)
	if err != nil {
		return 0, "", fmt.Errorf("get port failed: %w", err)
	}
	// err = db.QueryRow(
	// 	"SELECT domain_name FROM domains WHERE id = ?",
	// 	appID,
	// ).Scan(&domain)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return 0, "", fmt.Errorf("port not configured")
	// 	}
	// 	return 0, "", fmt.Errorf("get domain failed: %w", err)
	// }

	return *port, domain, nil
}
