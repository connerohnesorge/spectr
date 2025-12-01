package git

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// MinGitVersion is the minimum required git version for worktree support.
const MinGitVersion = "2.5"

// Constants for version parsing.
const (
	randomIDBytes        = 4
	versionMatchGroups   = 3
	minMajorVersion      = 2
	minMinorVersion      = 5
	gitCommand           = "git"
	worktreePathTemplate = "spectr-worktree-%s-%s"
)

// WorktreeInfo contains information about a created git worktree.
type WorktreeInfo struct {
	// Path is the filesystem path to the worktree directory.
	Path string
	// Branch is the name of the branch checked out in the worktree.
	Branch string
}

// CreateWorktree creates a new git worktree with a new branch.
// The worktree is created in a temporary directory with a UUID suffix.
// The new branch is based on origin/<baseBranch>.
func CreateWorktree(baseBranch, newBranch string) (WorktreeInfo, error) {
	if baseBranch == "" {
		return WorktreeInfo{}, errors.New("baseBranch cannot be empty")
	}
	if newBranch == "" {
		return WorktreeInfo{}, errors.New("newBranch cannot be empty")
	}

	// Generate unique worktree path using random bytes
	randomBytes := make([]byte, randomIDBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return WorktreeInfo{}, fmt.Errorf(
			"failed to generate random ID: %w", err,
		)
	}
	id := hex.EncodeToString(randomBytes)
	worktreePath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf(worktreePathTemplate, newBranch, id),
	)

	// Create worktree with new branch based on origin/<baseBranch>
	// Use -b to create new branch, specify start point as origin/<baseBranch>
	startPoint := fmt.Sprintf("origin/%s", baseBranch)
	cmd := exec.Command( //nolint:gosec
		gitCommand, "worktree", "add",
		"-b", newBranch, worktreePath, startPoint,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return WorktreeInfo{}, fmt.Errorf(
			"failed to create worktree: %s: %w",
			strings.TrimSpace(string(output)), err,
		)
	}

	return WorktreeInfo{
		Path:   worktreePath,
		Branch: newBranch,
	}, nil
}

// CleanupWorktree removes a git worktree safely.
// It uses 'git worktree remove' to ensure proper cleanup.
func CleanupWorktree(path string) error {
	if path == "" {
		return errors.New("worktree path cannot be empty")
	}

	// Use --force to handle uncommitted changes if needed
	cmd := exec.Command(gitCommand, "worktree", "remove", "--force", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If worktree doesn't exist, try removing the directory directly
		if strings.Contains(string(output), "is not a working tree") {
			if removeErr := os.RemoveAll(path); removeErr != nil {
				return fmt.Errorf(
					"failed to remove directory %s: %w", path, removeErr,
				)
			}

			return nil
		}

		return fmt.Errorf(
			"failed to remove worktree: %s: %w",
			strings.TrimSpace(string(output)), err,
		)
	}

	return nil
}

// ErrNoBaseBranch is returned when neither main nor master branch exists.
var ErrNoBaseBranch = errors.New(
	"neither origin/main nor origin/master exists; " +
		"please specify a base branch",
)

// GetBaseBranch determines the default base branch for the repository.
// It checks for origin/main first, then falls back to origin/master.
// Returns an error if neither branch exists on the remote.
func GetBaseBranch() (string, error) {
	// Check if origin/main exists
	cmd := exec.Command(gitCommand, "rev-parse", "--verify", "origin/main")
	if err := cmd.Run(); err == nil {
		return "main", nil
	}

	// Check if origin/master exists
	cmd = exec.Command(gitCommand, "rev-parse", "--verify", "origin/master")
	if err := cmd.Run(); err == nil {
		return "master", nil
	}

	return "", ErrNoBaseBranch
}

// gitVersionPattern matches git version strings like "git version 2.39.3"
var gitVersionPattern = regexp.MustCompile(`git version (\d+)\.(\d+)`)

// CheckGitVersion verifies that git is installed and meets minimum version.
// Returns an error if git is not found or version is less than 2.5.
func CheckGitVersion() error {
	cmd := exec.Command(gitCommand, "--version")
	output, err := cmd.Output()
	if err != nil {
		return errors.New("git is not installed or not in PATH")
	}

	matches := gitVersionPattern.FindStringSubmatch(string(output))
	if len(matches) != versionMatchGroups {
		return fmt.Errorf(
			"unable to parse git version from: %s",
			strings.TrimSpace(string(output)),
		)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("invalid git major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return fmt.Errorf("invalid git minor version: %s", matches[2])
	}

	// Minimum version is 2.5 (worktree support)
	belowMinor := major == minMajorVersion && minor < minMinorVersion
	if major < minMajorVersion || belowMinor {
		return fmt.Errorf(
			"git version %d.%d is below minimum required version %s",
			major, minor, MinGitVersion,
		)
	}

	return nil
}

// CheckBranchExists checks if a branch exists on the remote.
// It uses 'git ls-remote' to check for the branch.
func CheckBranchExists(branch string) (bool, error) {
	if branch == "" {
		return false, errors.New("branch name cannot be empty")
	}

	cmd := exec.Command(gitCommand, "ls-remote", "--heads", "origin", branch)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check remote branches: %w", err)
	}

	// If output is not empty, the branch exists
	return strings.TrimSpace(string(output)) != "", nil
}
