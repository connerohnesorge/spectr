package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// GetBaseBranch determines the appropriate base branch to use for a PR.
// If preferredBase is provided and exists, it returns origin/<preferredBase>.
// Otherwise, it auto-detects origin/main or falls back to origin/master.
func GetBaseBranch(preferredBase string) (string, error) {
	if preferredBase != "" {
		// Check if the preferred base exists
		exists, err := remoteBranchExists(preferredBase)
		if err != nil {
			return "", fmt.Errorf(
				"failed to check if branch '%s' exists: %w",
				preferredBase,
				err,
			)
		}
		if exists {
			return fmt.Sprintf("origin/%s", preferredBase), nil
		}

		return "", fmt.Errorf(
			"specified base branch '%s' does not exist on origin",
			preferredBase,
		)
	}

	// Auto-detect: check for main first, then master
	mainExists, err := remoteBranchExists("main")
	if err != nil {
		return "", fmt.Errorf("failed to check for main branch: %w", err)
	}
	if mainExists {
		return "origin/main", nil
	}

	masterExists, err := remoteBranchExists("master")
	if err != nil {
		return "", fmt.Errorf("failed to check for master branch: %w", err)
	}
	if masterExists {
		return "origin/master", nil
	}

	return "", errors.New(
		"could not determine base branch: " +
			"neither 'main' nor 'master' found on origin",
	)
}

// remoteBranchExists checks if a branch exists on the origin remote.
func remoteBranchExists(branchName string) (bool, error) {
	cmd := exec.Command(gitCmd, "ls-remote", "--heads", "origin", branchName)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf(
				"git ls-remote failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return false, fmt.Errorf("failed to run git ls-remote: %w", err)
	}

	// If output contains the branch name, it exists
	return strings.TrimSpace(string(output)) != "", nil
}

// BranchExists checks if a branch exists on the origin remote.
func BranchExists(branchName string) (bool, error) {
	return remoteBranchExists(branchName)
}

// DeleteRemoteBranch deletes a branch from the origin remote.
func DeleteRemoteBranch(branchName string) error {
	cmd := exec.Command(gitCmd, "push", "origin", "--delete", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to delete remote branch '%s': %s",
			branchName,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// FetchOrigin fetches the latest refs from the origin remote.
func FetchOrigin() error {
	cmd := exec.Command(gitCmd, "fetch", "origin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to fetch from origin: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// GetRepoRoot returns the absolute path to the root of the git repository.
// Returns an error if not in a git repository.
func GetRepoRoot() (string, error) {
	cmd := exec.Command(gitCmd, "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if strings.Contains(stderr, "not a git repository") {
				return "", errors.New("not a git repository")
			}

			return "", fmt.Errorf("git rev-parse failed: %s", stderr)
		}

		return "", fmt.Errorf("failed to run git rev-parse: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// PathExistsOnRef checks if a given path exists on a specific git ref.
// The ref should be a full ref like "origin/main" or "origin/master".
// The path should be relative to the repository root.
func PathExistsOnRef(ref, path string) (bool, error) {
	cmd := exec.Command(gitCmd, "ls-tree", ref, path)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf(
				"git ls-tree failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return false, fmt.Errorf("failed to run git ls-tree: %w", err)
	}

	// If output is non-empty, the path exists on the ref
	return strings.TrimSpace(string(output)) != "", nil
}

// deleteBranch deletes a local branch, ignoring errors if not found.
func deleteBranch(branchName string) []string {
	cmd := exec.Command(gitCmd, "branch", "-D", branchName)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	outputStr := strings.TrimSpace(string(output))
	notFound := strings.Contains(outputStr, "not found")
	branchErr := strings.Contains(outputStr, "error: branch")
	if notFound || branchErr {
		return nil
	}

	return []string{fmt.Sprintf("failed to delete branch: %s", outputStr)}
}
