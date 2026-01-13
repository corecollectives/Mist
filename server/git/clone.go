package git

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/corecollectives/mist/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/rs/zerolog/log"
)

func CloneRepo(url string, branch string, logFile *os.File, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GetConfig().Git.GitCloneTimeout)*time.Minute)
	defer cancel()

	_, err := fmt.Fprintf(logFile, "[GIT]: Cloning into %s", path)
	if err != nil {
		log.Warn().Msg("error logging into log file")
	}
	_, err = git.PlainClone(path, &git.CloneOptions{
		URL: url,
		// Progress:      logFile,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("git clone timed out after 10 minutes")
		}
		return err
	}

	return nil
}
