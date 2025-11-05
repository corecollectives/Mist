package dockerdeploy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/corecollectives/mist/models"
)

func (d *Deployer) DeployApp(dep *models.Deployment, appContextPath, imageTag, containerName string) error {
	fmt.Println("[DeployApp] Starting deployment process...")

	logFileName := fmt.Sprintf("%s_build_logs", dep.CommitHash)
	logPath := filepath.Join(d.LogDirectory, logFileName)
	fmt.Printf("[DeployApp] Log path: %s\n", logPath)

	err := os.MkdirAll(d.LogDirectory, 0755)
	if err != nil {
		fmt.Printf("[DeployApp] Failed to create log directory %s: %v\n", d.LogDirectory, err)
		return fmt.Errorf("failed to create log directory %s: %w", d.LogDirectory, err)
	}

	logfile, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("[DeployApp] Failed to create logfile %s: %v\n", logPath, err)
		return err
	}
	defer logfile.Close()
	fmt.Println("[DeployApp] Log file created successfully")

	dep.Logs = logPath
	dep.Status = "building"

	fmt.Printf("[DeployApp] Updating deployment status to %s\n", dep.Status)
	d.UpdateDeployment(dep)

	fmt.Printf("[DeployApp] Building image %s from %s\n", imageTag, appContextPath)
	err = BuildImage(imageTag, appContextPath, logfile)
	if err != nil {
		fmt.Printf("[DeployApp] Build failed: %v\n", err)
		fmt.Fprintf(logfile, "Build failed: %v\n", err)

		dep.Status = "failed"
		d.UpdateDeployment(dep)
		return err
	}
	fmt.Println("[DeployApp] Image build successful")

	dep.Status = "deploying"
	fmt.Printf("[DeployApp] Updating deployment status to %s\n", dep.Status)
	d.UpdateDeployment(dep)

	fmt.Printf("[DeployApp] Stopping/removing container %s if it exists\n", containerName)
	err = StopRemoveContainer(containerName, logfile)
	if err != nil {
		fmt.Fprintf(logfile, "Error stopping/removing container: %v\n", err)
		fmt.Printf("[DeployApp] Warning: stop/remove container failed: %v\n", err)
	}

	fmt.Printf("[DeployApp] Running container %s from image %s\n", containerName, imageTag)
	err = RunContainer(imageTag, containerName, []string{"-p", "6124:6124"}, logfile)
	if err != nil {
		fmt.Printf("[DeployApp] Failed to run container: %v\n", err)
		fmt.Fprintf(logfile, "Run container failed: %v\n", err)

		dep.Status = "failed"
		d.UpdateDeployment(dep)
		return err
	}

	dep.Status = "success"
	now := time.Now()
	dep.FinishedAt = &now
	fmt.Printf("[DeployApp] Deployment finished successfully at %v\n", dep.FinishedAt)
	d.UpdateDeployment(dep)
	return nil
}

func (d *Deployer) UpdateDeployment(dep *models.Deployment) error {
	fmt.Printf("[UpdateDeployment] Updating deployment ID %d -> status=%s, logs=%s, finished_at=%v\n",
		dep.ID, dep.Status, dep.Logs, dep.FinishedAt)

	doIt, err := d.DB.Prepare("UPDATE deployments SET status=?, logs=?, finished_at=? WHERE id=?")
	if err != nil {
		fmt.Printf("[UpdateDeployment] Failed to prepare statement: %v\n", err)
		return err
	}
	defer doIt.Close()

	_, err = doIt.Exec(dep.Status, dep.Logs, dep.FinishedAt, dep.ID)
	if err != nil {
		fmt.Printf("[UpdateDeployment] Exec error: %v\n", err)
		return err
	}

	fmt.Println("[UpdateDeployment] Deployment updated successfully")
	return nil
}

func (d *Deployer) GetLogsPath(commitHash string) string {
	logfileName := fmt.Sprintf("%s_build_logs", commitHash)
	path := filepath.Join(d.LogDirectory, logfileName)
	fmt.Printf("[GetLogsPath] Returning log path for commit %s -> %s\n", commitHash, path)
	return path
}
