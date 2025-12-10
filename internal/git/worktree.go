package git

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

const (
	// randomSuffixBytes is the number of random bytes for worktree
	// directory suffix.
	randomSuffixBytes = 4
	// gitCmd is the git command name.
	gitCmd = "git"
)

// WorktreeConfig contains configuration for creating a git worktree.
type WorktreeConfig struct {
	BranchName string // Branch name to create (e.g., "spectr/change-id")
	BaseBranch string // Base branch to start from (e.g., "origin/main")
}

// WorktreeInfo contains information about a created worktree.
type WorktreeInfo struct {
	Path       string // Absolute path to worktree directory
	BranchName string // The branch created
	TempDir    bool   // Whether this is in temp directory
}

// CreateWorktree creates a new git worktree in a temporary directory.
// It creates a new branch based on the specified base branch.
func CreateWorktree(config WorktreeConfig) (*WorktreeInfo, error) {
	if config.BranchName == "" {
		return nil, &specterrs.BranchNameRequiredError{}
	}
	if config.BaseBranch == "" {
		return nil, &specterrs.BaseBranchRequiredError{}
	}
	if _, err := GetRepoRoot(); err != nil {
		return nil, &specterrs.NotInGitRepositoryError{}
	}

	suffix, err := randomHex(randomSuffixBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random suffix: %w", err)
	}
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("spectr-pr-%s", suffix))

	cmd := exec.Command(
		gitCmd, "worktree", "add", tempDir,
		"-b", config.BranchName, config.BaseBranch,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, handleWorktreeError(output, config)
	}

	return &WorktreeInfo{
		Path:       tempDir,
		BranchName: config.BranchName,
		TempDir:    true,
	}, nil
}

// handleWorktreeError processes worktree creation errors.
func handleWorktreeError(output []byte, config WorktreeConfig) error {
	outputStr := strings.TrimSpace(string(output))
	if strings.Contains(outputStr, "already exists") {
		return fmt.Errorf(
			"branch '%s' already exists: %s", config.BranchName, outputStr,
		)
	}
	if strings.Contains(outputStr, "not a valid branch name") ||
		strings.Contains(outputStr, "invalid reference") {
		return fmt.Errorf(
			"base branch '%s' not found: %s", config.BaseBranch, outputStr,
		)
	}
	if strings.Contains(outputStr, "not a git repository") {
		return &specterrs.NotInGitRepositoryError{}
	}

	return fmt.Errorf("failed to create worktree: %s", outputStr)
}

// CleanupWorktree removes a git worktree and its associated branch.
// It is safe to call multiple times.
func CleanupWorktree(info *WorktreeInfo) error {
	if info == nil {
		return nil
	}

	var errs []string
	if info.Path != "" {
		errs = append(errs, removeWorktree(info.Path)...)
	}
	if info.BranchName != "" {
		errs = append(errs, deleteBranch(info.BranchName)...)
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errs, "; "))
	}

	return nil
}

// removeWorktree removes a worktree, ignoring errors if already removed.
func removeWorktree(path string) []string {
	cmd := exec.Command(gitCmd, "worktree", "remove", path, "--force")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	outputStr := strings.TrimSpace(string(output))
	isNotWT := strings.Contains(outputStr, "is not a working tree")
	noFile := strings.Contains(outputStr, "No such file or directory")
	notExist := strings.Contains(outputStr, "does not exist")
	if isNotWT || noFile || notExist {
		return nil
	}

	return []string{fmt.Sprintf("failed to remove worktree: %s", outputStr)}
}

// randomHex generates a random hex string of the specified byte length.
func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
