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
	return cmd.Run()
}

func StopRemoveContainer(containerName string, logfile *os.File) error {
	stopCmd := exec.Command("docker", "stop", containerName)
	stopCmd.Stdout = logfile
	stopCmd.Stderr = logfile
	ifExists := ContainerExists(containerName)
	fmt.Println("Container", containerName, "exists:", ifExists)
	if !ifExists {
		return nil
	}
	removeCmd := exec.Command("docker", "rm", containerName)
	removeCmd.Stdout = logfile
	removeCmd.Stderr = logfile
	if err := stopCmd.Run(); err != nil {
		fmt.Println("Warning: failed to stop container:", err.Error())
		return err
	}
	if err := removeCmd.Run(); err != nil {
		fmt.Println("Warning: failed to remove container:", err.Error())
		return err
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

	return cmd.Run()
}
