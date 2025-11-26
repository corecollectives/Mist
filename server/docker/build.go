package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func BuildImage(imageTag, contextPath string, envVars map[string]string, logfile *os.File) error {
	buildArgs := []string{"build", "-t", imageTag}

	for key, value := range envVars {
		buildArgs = append(buildArgs, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	buildArgs = append(buildArgs, contextPath)

	fmt.Println(buildArgs)
	cmd := exec.Command("docker", buildArgs...)
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
		return nil
	}

	stopCmd := exec.Command("docker", "stop", containerName)
	stopCmd.Stdout = logfile
	stopCmd.Stderr = logfile
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}

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

func RunContainer(imageTag, containerName string, domains []string, Port int, envVars map[string]string, logfile *os.File) error {

	runArgs := []string{
		"run", "-d",
		"--name", containerName,
	}

	for key, value := range envVars {
		runArgs = append(runArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	if len(domains) > 0 {
		runArgs = append(runArgs,
			"--network", "traefik-net",
			"-l", "traefik.enable=true",
		)

		var hostRules []string
		for _, domain := range domains {
			hostRules = append(hostRules, fmt.Sprintf("Host(`%s`)", domain))
		}
		hostRule := strings.Join(hostRules, " || ")

		runArgs = append(runArgs,
			"-l", fmt.Sprintf("traefik.http.routers.%s.rule=%s", containerName, hostRule),
			"-l", fmt.Sprintf("traefik.http.routers.%s.entrypoints=web", containerName),
			"-l", fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", containerName, Port),
		)
	} else {
		runArgs = append(runArgs,
			"-p", fmt.Sprintf("%d:%d", Port, Port),
		)
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
