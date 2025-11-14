package github

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/corecollectives/mist/models"
)

func CloneRepo(appId int64, logFile *os.File) error {
	println("Cloning repository for app ID:", appId)

	repo, branch, projectId, name, err := models.GetAppRepoInfo(appId)
	if err != nil {
		return fmt.Errorf("failed to fetch app: %w", err)
	}

	userId, err := models.GetUserIDByAppID(appId)
	if err != nil {
		return fmt.Errorf("failed to get user id by app id: %w", err)
	}
	accessToken, err := GetGitHubAccessToken(int(*userId))
	if err != nil {
		return fmt.Errorf("failed to get github access token: %w", err)
	}

	repoURL := fmt.Sprintf("https://github.com/%s.git", repo)
	if accessToken != "" {
		repoURL = fmt.Sprintf(
			"https://x-access-token:%s@github.com/%s.git",
			accessToken, repo,
		)
	}

	path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", projectId, name)

	if _, err := os.Stat(path + "/.git"); err == nil {
		fmt.Println("Repository exists → removing directory...")

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove existing repository: %w", err)

		}
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Println("Cloning repository...")
	cmd := exec.Command("git", "clone", "--branch", branch, repoURL, path)
	output, err := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Fprintf(logFile, "[GITHUB] %s\n", line)
		}
	}
	fmt.Println(string(output))
	if err != nil {
		return fmt.Errorf("error cloning repository: %v\n%s", err, string(output))
	}

	println("Repository cloned successfully to", appId)
	return nil
}
