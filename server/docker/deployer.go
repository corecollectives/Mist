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
)

func DeployApp(dep *models.Deployment, appContextPath, imageTag, containerName string, appId int64, createdBy int64, db *sql.DB, logfile *os.File) error {

	fmt.Println("deploying", dep.ID)
	dep.Status = "building"
	UpdateDeployment(dep, db)

	if err := BuildImage(imageTag, appContextPath, logfile); err != nil {
		println("Error: failed to build image:", err.Error())
		dep.Status = "failed"
		UpdateDeployment(dep, db)
		return err
	}

	dep.Status = "deploying"
	UpdateDeployment(dep, db)

	err := StopRemoveContainer(containerName, logfile)
	if err != nil {
		fmt.Println("Warning: failed to stop/remove existing container:", err.Error())
		dep.Status = "failed"
		UpdateDeployment(dep, db)
		return err
	}

	if err := RunContainer(imageTag, containerName, []string{"-p", "6124:6124"}, logfile); err != nil {
		dep.Status = "failed"
		UpdateDeployment(dep, db)
		return err
	}

	dep.Status = "success"
	now := time.Now()
	dep.FinishedAt = &now
	UpdateDeployment(dep, db)
	println("Deployment", dep.ID, "succeeded")

	return nil
}

func UpdateDeployment(dep *models.Deployment, db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE deployments SET status=?, logs=?, finished_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dep.Status, dep.Logs, dep.FinishedAt, dep.ID)
	return err
}

func GetLogsPath(commitHash string, depId int64) string {
	return filepath.Join(constants.Constants["LogPath"], commitHash+strconv.FormatInt(depId, 10)+"_build_logs")
}
