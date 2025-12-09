package queue

import (
	"database/sql"
	"fmt"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/fs"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/utils"
)

func (q *Queue) HandleWork(id int64, db *sql.DB) {
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("panic during deployment: %v", r)
			models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		}
	}()

	appId, err := models.GetAppIDByDeploymentID(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get app ID: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	dep, err := docker.LoadDeployment(id, db)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load deployment: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logger := utils.NewDeploymentLogger(id, appId, dep.CommitHash)
	logger.Info("Starting deployment processing")

	if err := models.MarkDeploymentStarted(id); err != nil {
		logger.Error(err, "Failed to mark deployment as started")
		errMsg := fmt.Sprintf("Failed to update deployment start time: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logFile, _, err := fs.CreateDockerBuildLogFile(id)
	if err != nil {
		logger.Error(err, "Failed to create log file")
		errMsg := fmt.Sprintf("Failed to create log file: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}
	defer logFile.Close()

	logger.Info("Cloning repository")
	models.UpdateDeploymentStatus(id, "cloning", "cloning", 20, nil)

	err = github.CloneRepo(appId, logFile)
	if err != nil {
		logger.Error(err, "Failed to clone repository")
		errMsg := fmt.Sprintf("Failed to clone repository: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logger.Info("Repository cloned successfully")

	_, err = docker.DeployerMain(id, db, logFile, logger)
	if err != nil {
		logger.Error(err, "Deployment failed")
		errMsg := fmt.Sprintf("Deployment failed: %v", err)
		models.UpdateDeploymentStatus(id, "failed", "failed", 0, &errMsg)
		return
	}

	logger.Info("Deployment completed successfully")
}
