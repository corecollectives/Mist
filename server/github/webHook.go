package github

import (
	"errors"
	"strings"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
)

func CreateDeploymentFromGithubPushEvent(evt PushEvent) (int64, error) {
	repoName := evt.Repository.FullName
	branch := evt.Ref
	commit := evt.After

	branch = strings.TrimPrefix(branch, "refs/heads/")

	log.Info().
		Str("repo", repoName).
		Str("branch", branch).
		Str("commit", commit).
		Msg("Push event received")

	appID, err := models.FindApplicationIDByGitRepoAndBranch(repoName, branch)

	if err != nil {
		log.Error().Err(err).
			Str("repo", repoName).
			Str("branch", branch).
			Msg("Error finding application")
		return 0, err
	}

	if appID == 0 {
		log.Warn().
			Str("repo", repoName).
			Str("branch", branch).
			Msg("No application found for this repository and branch")
		return 0, errors.New("no application found for this repository and branch")
	}

	commitMsg := evt.HeadCommit.Message
	deployment := models.Deployment{
		AppID:         appID,
		CommitHash:    commit,
		CommitMessage: &commitMsg,
	}

	if err := deployment.CreateDeployment(); err != nil {
		log.Error().Err(err).
			Int64("app_id", appID).
			Msg("Error creating deployment")
		return 0, err
	}

	log.Info().
		Int64("deployment_id", deployment.ID).
		Int64("app_id", appID).
		Msg("Deployment created from GitHub webhook")

	models.LogWebhookAudit("create", "deployment", &deployment.ID, map[string]interface{}{
		"app_id":         appID,
		"commit_hash":    commit,
		"commit_message": evt.HeadCommit.Message,
		"repository":     repoName,
		"branch":         branch,
		"pusher":         evt.Pusher.Name,
	})

	return deployment.ID, nil

}
