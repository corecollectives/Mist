package github

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

func CloneRepo(appId int64, logFile *os.File) error {
	log.Info().Int64("app_id", appId).Msg("Starting repository clone")

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
		log.Info().Str("path", path).Msg("Repository already exists, removing directory")

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove existing repository: %w", err)

		}
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	log.Info().Str("repo", repo).Str("branch", branch).Str("path", path).Msg("Cloning repository")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "clone", "--branch", branch, repoURL, path)
	output, err := cmd.CombinedOutput()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			fmt.Fprintf(logFile, "[GITHUB] %s\n", line)
		}
	}
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("git clone timed out after 10 minutes")
		}
		return fmt.Errorf("error cloning repository: %v\n%s", err, string(output))
	}

	log.Info().Int64("app_id", appId).Str("path", path).Msg("Repository cloned successfully")
	return nil
}
