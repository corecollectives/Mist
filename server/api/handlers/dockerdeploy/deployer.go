package dockerdeploy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/corecollectives/mist/models"
)

func (d *Deployer) DeployApp(dep *models.Deployment, appContextPath, imageTag, containerName string) error {
	logFileName := imageTag + "_build_logs"
	logPath := filepath.Join(d.LogDirectory, logFileName)

	if err := os.MkdirAll(d.LogDirectory, 0755); err != nil {
		return err
	}

	logfile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer logfile.Close()

	dep.Logs.String = logPath
	dep.Status = "building"
	d.UpdateDeployment(dep)

	if err := BuildImage(imageTag, appContextPath, logfile); err != nil {
		dep.Status = "failed"
		d.UpdateDeployment(dep)
		return err
	}

	dep.Status = "deploying"
	d.UpdateDeployment(dep)

	err = StopRemoveContainer(containerName, logfile)
	if err != nil {
		fmt.Println("Warning: failed to stop/remove existing container:", err.Error())
		dep.Status = "failed"
		d.UpdateDeployment(dep)
		return err
	}

	if err := RunContainer(imageTag, containerName, []string{"-p", "6124:6124"}, logfile); err != nil {
		dep.Status = "failed"
		d.UpdateDeployment(dep)
		return err
	}

	dep.Status = "success"
	now := time.Now()
	dep.FinishedAt = &now
	d.UpdateDeployment(dep)

	return nil
}

func (d *Deployer) UpdateDeployment(dep *models.Deployment) error {
	stmt, err := d.DB.Prepare("UPDATE deployments SET status=?, logs=?, finished_at=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dep.Status, dep.Logs, dep.FinishedAt, dep.ID)
	return err
}

func (d *Deployer) GetLogsPath(commitHash string) string {
	return filepath.Join(d.LogDirectory, commitHash+"_build_logs")
}
