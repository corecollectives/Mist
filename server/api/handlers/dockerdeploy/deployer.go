package dockerdeploy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/corecollectives/mist/models"
)

// func InitDeployer(db *sql.DB, logDir string) {
// 	deployer = &Deployer{DB: db, LogDirectory: logDir}
// }

func (d *Deployer) DeployApp(dep *models.Deployment, appConntextPath, imageTag, containerName string) error {
	logFileName := fmt.Sprintf("%s_build_logs", dep.CommitHash)
	logPath := filepath.Join(d.LogDirectory, logFileName)
	err := os.MkdirAll(d.LogDirectory, 0755)
	if err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", d.LogDirectory, err)
	}
	logfile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer logfile.Close()
	dep.Logs = logPath

	dep.Status = "building"

	d.UpdateDeployment(dep)

	err = BuildImage(imageTag, appConntextPath, logfile)
	if err != nil {
		dep.Status = "failed"
		d.UpdateDeployment(dep)
		return err
	}

	dep.Status = "deploying"
	d.UpdateDeployment(dep)

	err = StopRemoveContainer(containerName, logfile)
	if err != nil {
		fmt.Fprintf(logfile, "Error failed stop/remove: %v\n", err)
	}

	err = RunContainer(imageTag, containerName, []string{"-p", "6124:6124"}, logfile)
	if err != nil {
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
	doIt, err := d.DB.Prepare("UPDATE deployments SET status=?, logs=?, finished_at=? WHERE id=?")
	if err != nil {
		return nil
	}
	_, err = doIt.Exec(dep.Status, dep.Logs, dep.FinishedAt, dep.ID)
	return err
}

func (d *Deployer) GetLogsPath(CommitHash string) string {
	logfileName := fmt.Sprintf("%s_build_logs", CommitHash)
	return filepath.Join(d.LogDirectory, logfileName)
}
