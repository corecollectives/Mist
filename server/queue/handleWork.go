package queue

import (
	"database/sql"

	"github.com/corecollectives/mist/docker"
	"github.com/corecollectives/mist/fs"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
)

func (q *Queue) HandleWork(id int64, db *sql.DB) {
	logFile, _, err := fs.CreateDockerBuildLogFile(id)
	appId, err := models.GetAppIDByDeploymentID(id)
	if err != nil {
		return
	}
	err = github.CloneRepo(appId, logFile)
	if err != nil {
		return
	}

	_, err = docker.DeployerMain(id, db, logFile)
	if err != nil {
		return
	}

}
