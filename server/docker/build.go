package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/corecollectives/mist/models"
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

func RunContainer(app *models.App, imageTag, containerName string, domains []string, Port int, envVars map[string]string, logfile *os.File) error {

	runArgs := []string{
		"run", "-d",
		"--name", containerName,
	}

	restartPolicy := string(app.RestartPolicy)
	if restartPolicy == "" {
		restartPolicy = "unless-stopped"
	}
	runArgs = append(runArgs, "--restart", restartPolicy)

	if app.CPULimit != nil && *app.CPULimit > 0 {
		runArgs = append(runArgs, "--cpus", fmt.Sprintf("%.2f", *app.CPULimit))
	}

	if app.MemoryLimit != nil && *app.MemoryLimit > 0 {
		runArgs = append(runArgs, "-m", fmt.Sprintf("%dm", *app.MemoryLimit))
	}

	for key, value := range envVars {
		runArgs = append(runArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	switch app.AppType {
	case models.AppTypeWeb:
		// Web apps: Always use Traefik for routing if domains exist, otherwise expose port
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

	case models.AppTypeService:
		runArgs = append(runArgs, "--network", "traefik-net")

	case models.AppTypeDatabase:
		runArgs = append(runArgs, "--network", "traefik-net")

	default:
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

func PullDockerImage(imageName string, logfile *os.File) error {
	pullCmd := exec.Command("docker", "pull", imageName)
	pullCmd.Stdout = logfile
	pullCmd.Stderr = logfile

	if err := pullCmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return fmt.Errorf("docker pull failed with exit code %d: %w", exitCode, err)
	}
	return nil
}
