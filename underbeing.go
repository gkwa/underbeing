package underbeing

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Main() int {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	return 0
}

func run() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	isGitRepo, err := isGitRepository(currentDir)
	if err != nil {
		return fmt.Errorf("failed to check Git repository: %w", err)
	}

	if isGitRepo {
		fmt.Println("The current directory is a Git-initialized directory.")
	} else {
		fmt.Println("The current directory is not a Git-initialized directory.")
	}

	return nil
}

func isGitRepository(dir string) (bool, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return false, nil
		}
		return false, fmt.Errorf("failed to open Git repository: %w", err)
	}

	_, err = repo.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	return true, nil
}
