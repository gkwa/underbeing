package underbeing

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	optmod "github.com/taylormonacelli/underbeing/options"
)

func Main(opts *optmod.Options) int {
	err := run(opts)
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	return 0
}

func run(opts *optmod.Options) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	isGitRepo, err := isGitRepository(currentDir)
	if err != nil {
		return fmt.Errorf("failed to check Git repository: %w", err)
	}

	if isGitRepo {
		slog.Debug("the current directory is a Git-initialized directory", "path", currentDir)
	} else {
		return fmt.Errorf("the current directory is not a Git-initialized directory")
	}

	err = createOrUpdateGitHubRepo(opts.GithubUser, "your-repo-name")
	if err != nil {
		return fmt.Errorf("failed to create or update GitHub repository: %w", err)
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

func createOrUpdateGitHubRepo(username, repoName string) error {
	if username == "" {
		username = os.Getenv("GITHUB_USER")
	}

	exists, err := checkGitHubRepoExists(username, repoName)
	if err != nil {
		return fmt.Errorf("failed to check GitHub repository existence: %w", err)
	}

	if exists {
		fmt.Printf("the GitHub repository '%s/%s' already exists.\n", username, repoName)
	} else {
		args := []string{"repo", "create", username + "/" + repoName, "--public"}
		stdOut, stdErr, err := gh.Exec(args...)
		if err != nil {
			return fmt.Errorf("failed to execute gh command: %w\n%s", err, stdErr.String())
		}

		slog.Debug("repo create", "stdout", stdOut.String())
		slog.Error("repo create", "stderr", stdErr.String())
		fmt.Printf("GitHub repository '%s/%s' created successfully.\n", username, repoName)
	}

	return nil
}

func checkGitHubRepoExists(username, repoName string) (bool, error) {
	args := []string{"repo", "view", username + "/" + repoName}
	stdOut, _, err := gh.Exec(args...)

	if err != nil && stdOut.String() == "Error: Not Found\n" {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to execute gh command: %w", err)
	}

	return true, nil
}
