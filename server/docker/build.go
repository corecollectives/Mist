package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func BuildImage(imageTag, contextPath string, logfile *os.File) error {
	cmd := exec.Command("docker", "build", "-t", imageTag, contextPath)
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	if err := cmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return fmt.Errorf("docker build failed with exit code %d: %w", exitCode, err)
	}
	return nil
}

func StopRemoveContainer(containerName string, logfile *os.File) error {
	ifExists := ContainerExists(containerName)
	if !ifExists {
		// Container doesn't exist, nothing to stop/remove
		return nil
	}

	// Stop container
	stopCmd := exec.Command("docker", "stop", containerName)
	stopCmd.Stdout = logfile
	stopCmd.Stderr = logfile
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}

	// Remove container
	removeCmd := exec.Command("docker", "rm", containerName)
	removeCmd.Stdout = logfile
	removeCmd.Stderr = logfile
	if err := removeCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerName, err)
	}

	return nil
}

func ContainerExists(name string) bool {
	cmd := exec.Command("docker", "inspect", name)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if strings.Contains(string(output), "No such object") {
			return false
		}
		return false
	}

	return true
}

func RunContainer(imageTag, containerName string, domain string, Port int, logfile *os.File) error {

	runArgs := []string{
		"run", "-d",
		"--network", "traefik-net",
		"-l", "traefik.enable=true",
		"-l", fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s`)", containerName, domain),
		"-l", fmt.Sprintf("traefik.http.routers.%s.entrypoints=web", containerName),
		"-l", fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", containerName, Port),
		"--name", containerName,
	}
	runArgs = append(runArgs, imageTag)

	cmd := exec.Command("docker", runArgs...)
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	if err := cmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return fmt.Errorf("docker run failed with exit code %d: %w", exitCode, err)
	}

	return nil
}
