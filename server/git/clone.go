package git

import (
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

func CloneRepo(url string, branch string, logFile *os.File, path string) error {
	_, err := git.PlainClone(path, &git.CloneOptions{
		URL: url,
		// Progress:      logFile,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		return err
	}

	return nil
}
