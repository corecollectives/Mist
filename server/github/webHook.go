package github

import (
	"errors"
	"fmt"
	"strings"

	"github.com/corecollectives/mist/models"
)

func HandlePushEvent(evt PushEvent) (int64, error) {
	repoName := evt.Repository.FullName
	branch := evt.Ref
	commit := evt.After

	branch = strings.TrimPrefix(branch, "refs/heads/")

	fmt.Printf("Push event received for repo: %s, branch: %s, commit: %s\n", repoName, branch, commit)

	appID, err := models.FindApplicationIDByGitRepoAndBranch(repoName, branch)

	if err != nil {
		fmt.Printf("Error finding application: %v\n", err)
		return 0, err

	}

	if appID == 0 {
		fmt.Println("No application found for this repository and branch.")
		return 0, errors.New("no application found for this repository and branch")
	}

	deployment := models.Deployment{
		AppID:         appID,
		CommitHash:    commit,
		CommitMessage: evt.HeadCommit.Message,
	}

	if err := deployment.CreateDeployment(); err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return 0, err
	}

	fmt.Printf("Deployment created with ID: %d for App ID: %d\n", deployment.ID, appID)
	return deployment.ID, nil

}
