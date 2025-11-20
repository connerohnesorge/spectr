// Package git provides utilities for git operations used in the
// archive workflow. It includes functions for repository validation,
// branch management, commits, and push operations to support
// automated git workflows.
package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const gitCommand = "git"

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() error {
	cmd := exec.Command(gitCommand, "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return errors.New(
			"not in a git repository. Initialize git with 'git init'",
		)
	}

	return nil
}

// HasOriginRemote checks if the origin remote is configured
func HasOriginRemote() error {
	cmd := exec.Command(gitCommand, "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return errors.New(
			"no 'origin' remote configured. " +
				"Run 'git remote add origin <url>'",
		)
	}
	if strings.TrimSpace(string(output)) == "" {
		return errors.New("origin remote URL is empty")
	}

	return nil
}

// CreateBranch creates a new git branch with the given name
func CreateBranch(branchName string) error {
	cmd := exec.Command(gitCommand, "checkout", "-b", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create branch: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// StageFiles stages files for commit
func StageFiles(paths []string) error {
	args := append([]string{"add"}, paths...)
	cmd := exec.Command(gitCommand, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stage files: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Commit creates a git commit with the given message
func Commit(message string) error {
	// Use heredoc-style message passing via stdin
	cmd := exec.Command(gitCommand, "commit", "-F", "-")
	cmd.Stdin = strings.NewReader(message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("commit: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Push pushes the current branch to origin
func Push(branchName string) error {
	cmd := exec.Command(gitCommand, "push", "-u", "origin", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("push branch: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetCurrentBranch returns the name of the current git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command(gitCommand, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CheckoutBranch switches to the specified git branch
func CheckoutBranch(branchName string) error {
	cmd := exec.Command(gitCommand, "checkout", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "checkout branch: %w\nOutput: %s"

		return fmt.Errorf(msg, err, string(output))
	}

	return nil
}

// BranchExists checks if a branch exists locally
func BranchExists(branchName string) bool {
	cmd := exec.Command(gitCommand, "rev-parse", "--verify", branchName)
	err := cmd.Run()

	return err == nil
}

// RestorePath restores a file or directory from HEAD
func RestorePath(path string) error {
	cmd := exec.Command(gitCommand, "checkout", "HEAD", "--", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("restore path: %w\nOutput: %s", err, string(output))
	}

	return nil
}
