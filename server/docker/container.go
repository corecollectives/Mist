package docker

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ContainerStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	State   string `json:"state"`
	Uptime  string `json:"uptime"`
	Healthy bool   `json:"healthy"`
}

func GetContainerStatus(containerName string) (*ContainerStatus, error) {
	if !ContainerExists(containerName) {
		return &ContainerStatus{
			Name:    containerName,
			Status:  "not_found",
			State:   "stopped",
			Uptime:  "N/A",
			Healthy: false,
		}, nil
	}

	cmd := exec.Command("docker", "inspect", containerName, "--format", "{{json .}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	var inspectData struct {
		State struct {
			Status  string `json:"Status"`
			Running bool   `json:"Running"`
			Paused  bool   `json:"Paused"`
			Health  *struct {
				Status string `json:"Status"`
			} `json:"Health"`
		} `json:"State"`
		Name string `json:"Name"`
	}

	if err := json.Unmarshal(output, &inspectData); err != nil {
		return nil, fmt.Errorf("failed to parse inspect output: %w", err)
	}

	uptimeCmd := exec.Command("docker", "inspect", containerName, "--format", "{{.State.StartedAt}}")
	uptimeOutput, err := uptimeCmd.Output()
	uptime := "N/A"
	if err == nil {
		uptime = strings.TrimSpace(string(uptimeOutput))
	}

	state := "stopped"
	if inspectData.State.Running {
		state = "running"
	} else if inspectData.State.Status == "exited" {
		state = "stopped"
	} else {
		state = inspectData.State.Status
	}

	healthy := true
	if inspectData.State.Health != nil {
		healthy = inspectData.State.Health.Status == "healthy"
	}

	return &ContainerStatus{
		Name:    strings.TrimPrefix(inspectData.Name, "/"),
		Status:  inspectData.State.Status,
		State:   state,
		Uptime:  uptime,
		Healthy: healthy,
	}, nil
}

func StopContainer(containerName string) error {
	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	cmd := exec.Command("docker", "stop", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

func StartContainer(containerName string) error {
	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	cmd := exec.Command("docker", "start", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func RestartContainer(containerName string) error {
	if !ContainerExists(containerName) {
		return fmt.Errorf("container %s does not exist", containerName)
	}

	cmd := exec.Command("docker", "restart", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	return nil
}

func GetContainerLogs(containerName string, tail int) (string, error) {
	if !ContainerExists(containerName) {
		return "", fmt.Errorf("container %s does not exist", containerName)
	}

	tailStr := fmt.Sprintf("%d", tail)
	cmd := exec.Command("docker", "logs", "--tail", tailStr, containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}

	return string(output), nil
}

func GetContainerName(appName string, appId int64) string {
	return fmt.Sprintf("app-%d", appId)
}
