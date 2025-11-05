package github

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
)

func CloneRepo(db *sql.DB, appId int64, userId int64) error {
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
		repoURL = fmt.Sprintf("https://x-access-token:%s@github.com/%s.git", accessToken, repo)
	}

	path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", projectId, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "--branch", branch, repoURL, path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error cloning repository: %v\n%s", err, string(output))
	}

	return nil
}

