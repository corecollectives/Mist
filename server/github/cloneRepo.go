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

	// Build repo URL
	repoURL := fmt.Sprintf("https://github.com/%s.git", repo)
	if accessToken != "" {
		repoURL = fmt.Sprintf(
			"https://x-access-token:%s@github.com/%s.git",
			accessToken, repo,
		)
	}

	// Filesystem path
	path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", projectId, name)


	// ✅ If repo exists → delete it completely
	if _, err := os.Stat(path + "/.git"); err == nil {
		fmt.Println("Repository exists → removing directory...")

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove existing repository: %w", err)

		}
	}

	// ✅ Ensure directory exists
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// ✅ Clone fresh repo
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

	return nil
}
