package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

func BuildImage(imageTag, contextPath string, envVars map[string]string, logfile *os.File) error {
	buildArgs := []string{"build", "-t", imageTag}

	for key, value := range envVars {
		buildArgs = append(buildArgs, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	buildArgs = append(buildArgs, contextPath)

	log.Debug().Strs("build_args", buildArgs).Str("image_tag", imageTag).Msg("Building Docker image")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", buildArgs...)
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker build timed out after 15 minutes")
		}
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	stopCmd := exec.CommandContext(ctx, "docker", "stop", containerName)
	stopCmd.Stdout = logfile
	stopCmd.Stderr = logfile
	if err := stopCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker stop timed out after 2 minutes for container %s", containerName)
		}
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}

	removeCmd := exec.CommandContext(ctx, "docker", "rm", containerName)
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

	// Add volumes from the volumes table (user-configurable)
	volumes, err := models.GetVolumesByAppID(app.ID)
	if err == nil {
		for _, vol := range volumes {
			volumeArg := fmt.Sprintf("%s:%s", vol.HostPath, vol.ContainerPath)
			if vol.ReadOnly {
				volumeArg += ":ro"
			}
			runArgs = append(runArgs, "-v", volumeArg)
		}
	}

	for key, value := range envVars {
		runArgs = append(runArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	switch app.AppType {
	case models.AppTypeWeb:
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
				"-l", fmt.Sprintf("traefik.http.routers.%s.entrypoints=websecure", containerName),
				"-l", fmt.Sprintf("traefik.http.routers.%s.tls=true", containerName),
				"-l", fmt.Sprintf("traefik.http.routers.%s.tls.certresolver=le", containerName),
				"-l", fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%d", containerName, Port),
			)

			runArgs = append(runArgs,

				"-l", fmt.Sprintf("traefik.http.routers.%s-http.rule=%s", containerName, hostRule),
				"-l", fmt.Sprintf("traefik.http.routers.%s-http.entrypoints=web", containerName),
				"-l", fmt.Sprintf("traefik.http.routers.%s-http.middlewares=%s-https-redirect", containerName, containerName),

				"-l", fmt.Sprintf("traefik.http.middlewares.%s-https-redirect.redirectscheme.scheme=https", containerName),
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", runArgs...)
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker run timed out after 5 minutes")
		}
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return fmt.Errorf("docker run failed with exit code %d: %w", exitCode, err)
	}

	return nil
}

func PullDockerImage(imageName string, logfile *os.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	pullCmd := exec.CommandContext(ctx, "docker", "pull", imageName)
	pullCmd.Stdout = logfile
	pullCmd.Stderr = logfile

	if err := pullCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docker pull timed out after 15 minutes for image %s", imageName)
		}
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return fmt.Errorf("docker pull failed with exit code %d: %w", exitCode, err)
	}
	return nil
}
