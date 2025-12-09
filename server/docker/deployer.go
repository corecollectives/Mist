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

	logger.Info("Getting port, domains, and environment variables")
	port, domains, envVars, err := GetDeploymentConfig(dep.ID, db)
	if err != nil {
		logger.Error(err, "Failed to get deployment configuration")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to get deployment config: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		return fmt.Errorf("get deployment config failed: %w", err)
	}

	logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"domains": domains,
		"port":    port,
		"envVars": len(envVars),
	})

	dep.Status = "building"
	dep.Stage = "building"
	dep.Progress = 50
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "building", "building", 50, nil)

	logger.Info("Building Docker image with environment variables")
	if err := BuildImage(imageTag, appContextPath, envVars, logfile); err != nil {
		logger.Error(err, "Docker image build failed")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Build failed: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateAppStatus(appId, "error", db)
		return fmt.Errorf("build image failed: %w", err)
	}

	logger.Info("Docker image built successfully")

	dep.Status = "deploying"
	dep.Stage = "deploying"
	dep.Progress = 80
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "deploying", "deploying", 80, nil)

	logger.Info("Stopping existing container if exists")
	err = StopRemoveContainer(containerName, logfile)
	if err != nil {
		logger.Error(err, "Failed to stop/remove existing container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to stop/remove container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateAppStatus(appId, "error", db)
		return fmt.Errorf("stop/remove container failed: %w", err)
	}

	logger.InfoWithFields("Running container", map[string]interface{}{
		"domains": domains,
		"port":    port,
		"envVars": len(envVars),
	})

	if err := RunContainer(imageTag, containerName, domains, port, envVars, logfile); err != nil {
		logger.Error(err, "Failed to run container")
		dep.Status = "failed"
		dep.Stage = "failed"
		dep.Progress = 0
		errMsg := fmt.Sprintf("Failed to run container: %v", err)
		dep.ErrorMessage = &errMsg
		UpdateDeployment(dep, db)
		models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errMsg)
		UpdateAppStatus(appId, "error", db)
		return fmt.Errorf("run container failed: %w", err)
	}

	dep.Status = "success"
	dep.Stage = "success"
	dep.Progress = 100
	now := time.Now()
	dep.FinishedAt = &now
	UpdateDeployment(dep, db)
	models.UpdateDeploymentStatus(dep.ID, "success", "success", 100, nil)

	logger.Info("Updating app status to running")
	err = UpdateAppStatus(appId, "running", db)
	if err != nil {
		logger.Error(err, "Failed to update app status (non-fatal)")
	}

	logger.InfoWithFields("Deployment succeeded", map[string]interface{}{
		"deployment_id": dep.ID,
		"container":     containerName,
		"app_status":    "running",
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
	return filepath.Join(constants.Constants["LogPath"].(string), commitHash+strconv.FormatInt(depId, 10)+"_build_logs")
}

func UpdateAppStatus(appID int64, status string, db *sql.DB) error {
	query := `UPDATE apps SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, status, appID)
	return err
}

func GetDeploymentConfig(deploymentID int64, db *sql.DB) (int, []string, map[string]string, error) {
	appID, err := models.GetAppIDByDeploymentID(deploymentID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get app ID failed: %w", err)
	}

	var port *int
	err = db.QueryRow("SELECT port FROM apps WHERE id = ?", appID).Scan(&port)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("get port failed: %w", err)
	}
	if port == nil {
		defaultPort := 3000
		port = &defaultPort
	}

	domains, err := models.GetDomainsByAppID(appID)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, nil, fmt.Errorf("get domains failed: %w", err)
	}

	var domainStrings []string
	for _, d := range domains {
		domainStrings = append(domainStrings, d.Domain)
	}

	envs, err := models.GetEnvVariablesByAppID(appID)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, nil, fmt.Errorf("get env variables failed: %w", err)
	}

	envMap := make(map[string]string)
	for _, env := range envs {
		envMap[env.Key] = env.Value
	}

	return *port, domainStrings, envMap, nil
}
