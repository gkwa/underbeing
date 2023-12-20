package underbeing

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
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

	repoName := opts.RepoName
	if repoName == "" {
		// Use the current directory name as the repository name
		_, repoName = filepath.Split(currentDir)
	}

	slog.Debug("debug reponame", "repo", repoName)

	var username string
	username = opts.GithubUser
	if username == "" {
		username = os.Getenv("GITHUB_USER")
		if username == "" {
			return fmt.Errorf("username is empty")
		}
	}

	err = createOrUpdateGitHubRepo(username, repoName)
	if err != nil {
		slog.Error("createOrUpdateGitHubRepo", "error", err)
		return fmt.Errorf("failed to create or update GitHub repository: %w", err)
	}

	err = addGitRemote(username, repoName)
	if err != nil {
		slog.Error("addGitRemote", "error", err)
		return fmt.Errorf("failed to add Git remote: %w", err)
	}

	slog.Debug("check githubuser", "githubuser", username)

	err = pushToRemote(username, repoName)
	if err != nil {
		slog.Error("pushToRemote", "error", err)
		return fmt.Errorf("failed to push changes to remote: %w", err)
	}

	return nil
}

func pushToRemote(username, repoName string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := git.PlainOpen(currentDir)
	if err != nil {
		return fmt.Errorf("failed to open Git repository: %w", err)
	}

	// Get the current branch
	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	remoteURL, err := getRemoteOriginURL(currentDir)
	if err != nil {
		return err
	}

	// expect: master:master
	refspec := gitconfig.RefSpec(headRef.Name().Short() + ":" + headRef.Name().Short())
	slog.Debug("refspec debug", "refspec", refspec)

	remoteConfig := &gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	}

	// Set up the push options
	pushOptions := &git.PushOptions{
		RemoteName: remoteConfig.Name,
		RefSpecs:   []gitconfig.RefSpec{gitconfig.RefSpec(fmt.Sprintf("%s:%s", headRef.Name(), headRef.Name()))},
	}

	err = repo.Push(pushOptions)
	if err != nil {
		return fmt.Errorf("failed to push changes to remote: %w", err)
	}

	fmt.Printf("Changes pushed to remote 'origin' and upstream branch set to '%s'.\n", headRef.Name())
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
	exists, err := checkGitHubRepoExists(username, repoName)
	if err != nil {
		slog.Error("checkGitHubRepoExists", "error", err)
		return fmt.Errorf("failed to check GitHub repository existence: %w", err)
	}

	if exists {
		remoteURL, err := getRemoteOriginURL(".")
		if err != nil {
			return err
		}

		slog.Error("the GitHub repository already exists", "repo", remoteURL)
	} else {
		args := []string{"repo", "create", username + "/" + repoName, "--public"}
		stdOut, stdErr, err := gh.Exec(args...)
		if err != nil {
			slog.Error("repo create", "username", username, "repo", repoName)
			return fmt.Errorf("failed to execute gh command: %w\n%s", err, stdErr.String())
		}

		slog.Debug("repo create", "stdout", stdOut.String())
		slog.Debug("repo create", "stderr", stdErr.String())
		slog.Debug("gitHub repository created", "repo", username+"/"+repoName)
	}

	return nil
}

func checkGitHubRepoExists(username, repoName string) (bool, error) {
	args := []string{"repo", "view", username + "/" + repoName}
	stdOut, stdErr, err := gh.Exec(args...)

	if err != nil && stdOut.String() == "Error: Not Found\n" {
		return false, nil
	}

	slog.Debug("gh.Exec() returns stderr string", "stderr", stdErr.String())

	// FIXME: use github response code instead if there is one
	errStr := fmt.Sprintf("Could not resolve to a Repository with the name '%s'",
		fmt.Sprintf("%s/%s", username, repoName),
	)
	if strings.Contains(stdErr.String(), errStr) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to execute gh command: %w", err)
	}

	return true, nil
}

func addGitRemote(username, repoName string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := git.PlainOpen(currentDir)
	if err != nil {
		return fmt.Errorf("failed to open Git repository: %w", err)
	}

	remoteURL := fmt.Sprintf("git@github.com:%s/%s.git", username, repoName)
	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	})

	slog.Debug("remote url", "remoteURL", remoteURL)

	if err != nil {
		return fmt.Errorf("failed to add Git remote: %w", err)
	}

	fmt.Printf("Git remote 'origin' added successfully with URL: %s\n", remoteURL)
	return nil
}

func getRemoteOriginURL(dir string) (string, error) {
	currentDir := "."

	repo, err := git.PlainOpen(currentDir)
	if err != nil {
		return "", fmt.Errorf("failed to open Git repository: %w", err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return "", fmt.Errorf("failed to get remote 'origin': %w", err)
	}

	remoteURLs := remote.Config().URLs
	if len(remoteURLs) == 0 {
		return "", fmt.Errorf("remote 'origin' has no URLs configured")
	}

	remoteURL := remoteURLs[0]

	return remoteURL, nil
}
