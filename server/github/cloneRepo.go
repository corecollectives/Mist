package github

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CloneRepo(db *sql.DB, appId int64, userId int64, logFile *os.File) error {
	var repo, branch string
	var projectId int64
	var name string

	err := db.QueryRow(`
		SELECT git_repository, git_branch, project_id, name
		FROM apps WHERE id = ?
	`, appId).Scan(&repo, &branch, &projectId, &name)
	if err != nil {
		return fmt.Errorf("failed to fetch app: %w", err)
	}

	accessToken, err := GetGitHubAccessToken(db, int(userId))
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
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	_, err = os.Stat(path + "/.git")

	if err == nil {

		fmt.Println("Repository exists, pulling latest changes...")
		cmd := exec.Command("git", "-C", path, "pull", "origin", branch)
		output, err := cmd.CombinedOutput()
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if len(line) > 0 {
				fmt.Fprintf(logFile, "[GITHUB] %s\n", line)
			}
		}
		fmt.Println(string(output))
		if err != nil {
			return fmt.Errorf("error pulling repository: %v\n%s", err, string(output))
		}
		return nil
	}

	if !os.IsNotExist(err) {

		return fmt.Errorf("error checking repo directory: %w", err)
	}

	fmt.Println("Repository does not exist, cloning...")
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

	return nil
}
