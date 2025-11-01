package docker

import (
	"os"
	"os/exec"
)

func BuildImage(imageTag, contextPath string, logfile *os.File) error {
	cmd := exec.Command("docker", "build", "-t", imageTag, contextPath)
	cmd.Stdout = logfile
	cmd.Stderr = logfile
	return cmd.Run()
}

func StopRemoveContainer(containerName string, logfile *os.File) error {
	stopCmd := exec.Command("docker", "stop", containerName)
	stopCmd.Stdout = logfile
	stopCmd.Stderr = logfile
	stopCmd.Run()

	removeCmd := exec.Command("docker", "rm", containerName)
	removeCmd.Stdout = logfile
	removeCmd.Stderr = logfile
	return removeCmd.Run()
}

func RunContainer(imageTag, containerName string, args []string, logfile *os.File) error {
	runArgs := []string{"run","-d"}
	runArgs = append(runArgs, args...)
	runArgs = append(runArgs, "--name", containerName)
	runArgs = append(runArgs, imageTag)
	cmd := exec.Command("docker", runArgs...)
	cmd.Stdout = logfile
	cmd.Stderr = logfile
	return cmd.Run()
}
