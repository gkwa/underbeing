package underbeing

import (
	"flag"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	optmod "github.com/taylormonacelli/underbeing/options"
)

func Main(options *optmod.Options) int {
	flag.Parse()

	if err := run(options); err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	return 0
}

func run(options *optmod.Options) error {
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

	err = createGitHubRepo(options.GithubUser, "your-repo-name")
	if err != nil {
		return fmt.Errorf("failed to create GitHub repository: %w", err)
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

func createGitHubRepo(username, repoName string) error {
	if username == "" {
		username = os.Getenv("GITHUB_USER")
	}

	args := []string{"repo", "create", username + "/" + repoName, "--public"}
	stdOut, stdErr, err := gh.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute gh command: %w\n%s", err, stdErr.String())
	}

	fmt.Println(stdOut.String())
	fmt.Println(stdErr.String())

	return nil
}
