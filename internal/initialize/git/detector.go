// Package git provides git-based change detection for spectr initialization.
//
// This package implements the git integration for detecting file changes
// during spectr init. It uses git diff to determine which files were
// modified, created, or deleted after initialization completes.
//
// The ChangeDetector type captures the state of the working tree before
// initialization and compares it after to produce a list of changed files.
//
//nolint:revive // file-length-limit - git operations require cohesive logic
package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// gitCmd is the git command name.
	gitCmd = "git"

	// gitRepoFlag is the flag for specifying git repository root.
	gitRepoFlag = "-C"

	// minPorcelainLineLen is the minimum length of a git status porcelain line.
	// Format: "XY filename" where X and Y are status codes.
	minPorcelainLineLen = 3

	// ErrNotGitRepo is the error message for non-git repositories.
	ErrNotGitRepo = "spectr init requires a git repository"

	// newline is the newline character for splitting git output.
	newline = "\n"
)

// GitExecutor abstracts git operations for testing.
type GitExecutor interface {
	// RevParse runs `git rev-parse` and returns the result.
	RevParse(repoPath, ref string) (string, error)

	// StatusPorcelain runs `git status --porcelain` and returns the output.
	StatusPorcelain(repoPath string) (string, error)

	// DiffNameOnly runs `git diff --name-only` between two refs.
	DiffNameOnly(repoPath, fromRef, toRef string) (string, error)

	// StashCreate runs `git stash create` and returns the stash ref.
	StashCreate(repoPath string) (string, error)

	// IsGitRepo checks if the path is inside a git repository.
	IsGitRepo(path string) bool
}

// RealGitExecutor implements GitExecutor using actual git commands.
type RealGitExecutor struct{}

// RevParse runs `git rev-parse` and returns the result.
func (*RealGitExecutor) RevParse(repoPath, ref string) (string, error) {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, repoPath,
		"rev-parse",
		ref,
	)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git rev-parse failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git rev-parse: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// StatusPorcelain runs `git status --porcelain` and returns the output.
func (*RealGitExecutor) StatusPorcelain(repoPath string) (string, error) {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, repoPath,
		"status",
		"--porcelain",
	)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git status failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git status: %w", err)
	}

	return string(output), nil
}

// DiffNameOnly runs `git diff --name-only` between two refs.
// If toRef is empty, compares fromRef to working tree.
func (*RealGitExecutor) DiffNameOnly(
	repoPath, fromRef, toRef string,
) (string, error) {
	args := []string{gitRepoFlag, repoPath, "diff", "--name-only"}
	if fromRef != "" {
		args = append(args, fromRef)
	}
	if toRef != "" {
		args = append(args, toRef)
	}

	cmd := exec.Command(gitCmd, args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git diff failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git diff: %w", err)
	}

	return string(output), nil
}

// StashCreate runs `git stash create` and returns the stash ref.
// Returns empty string if there's nothing to stash (clean working tree).
func (*RealGitExecutor) StashCreate(repoPath string) (string, error) {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, repoPath,
		"stash",
		"create",
	)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git stash create failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git stash create: %w", err)
	}

	// Returns empty string if nothing to stash
	return strings.TrimSpace(string(output)), nil
}

// IsGitRepo checks if the path is inside a git repository.
func (*RealGitExecutor) IsGitRepo(path string) bool {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, path,
		"rev-parse",
		"--git-dir",
	)
	err := cmd.Run()

	return err == nil
}

// ChangeDetector tracks file changes in a git repository during initialization.
// It captures the state before initialization and compares after to determine
// which files were changed.
type ChangeDetector struct {
	repoPath string
	executor GitExecutor
}

// NewChangeDetector creates a new ChangeDetector for the given repository.
func NewChangeDetector(repoPath string) *ChangeDetector {
	return &ChangeDetector{
		repoPath: repoPath,
		executor: &RealGitExecutor{},
	}
}

// NewChangeDetectorWithExecutor creates a new ChangeDetector with a custom
// GitExecutor. This is primarily used for testing with mock implementations.
func NewChangeDetectorWithExecutor(
	repoPath string,
	executor GitExecutor,
) *ChangeDetector {
	return &ChangeDetector{
		repoPath: repoPath,
		executor: executor,
	}
}

// IsGitRepo checks if the specified path is inside a git repository.
// This is a standalone function for early validation before creating a
// ChangeDetector.
func IsGitRepo(path string) bool {
	// First check for .git directory (handles most cases efficiently)
	gitDir := filepath.Join(path, ".git")
	if info, err := os.Stat(gitDir); err == nil {
		// .git can be a directory (normal repo) or file (worktree/submodule)
		return info.IsDir() || info.Mode().IsRegular()
	}

	// Fall back to git rev-parse for worktrees and other edge cases
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, path,
		"rev-parse",
		"--git-dir",
	)
	err := cmd.Run()

	return err == nil
}

// Snapshot captures the current state of the git working tree.
// It returns a reference string that can be used with ChangedFiles to
// determine what files changed.
//
// The snapshot captures:
//   - Current HEAD commit
//   - Dirty working tree state (via git stash create)
//   - List of untracked files
//
// Returns an error if not in a git repository or git operations fail.
func (d *ChangeDetector) Snapshot() (string, error) {
	if !d.executor.IsGitRepo(d.repoPath) {
		return "", errors.New(ErrNotGitRepo)
	}

	// Get current HEAD commit
	head, err := d.executor.RevParse(d.repoPath, "HEAD")
	if err != nil {
		// Repository might be empty (no commits yet)
		head = ""
	}

	// Try to create a stash of current changes (doesn't modify working tree)
	stashRef, err := d.executor.StashCreate(d.repoPath)
	if err != nil {
		// Stash create can fail if there are issues, but we can continue
		// with just HEAD tracking
		stashRef = ""
	}

	// Get list of untracked files
	status, err := d.executor.StatusPorcelain(d.repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}

	untrackedFiles := parseUntrackedFiles(status)

	// Create a snapshot string that encodes all state
	// Format: "HEAD:<commit>|STASH:<stash>|UNTRACKED:<file1>,<file2>,..."
	snapshot := fmt.Sprintf(
		"HEAD:%s|STASH:%s|UNTRACKED:%s",
		head,
		stashRef,
		strings.Join(untrackedFiles, ","),
	)

	return snapshot, nil
}

// ChangedFiles returns a list of files that changed since the snapshot.
// This includes:
//   - Modified files (tracked files that changed)
//   - New files (untracked files that appeared)
//   - Deleted files are NOT included (only existing files)
//
// The snapshot parameter should be a string returned by Snapshot().
// Returns an error if git operations fail.
//
//nolint:revive // function-length - complex git diff logic requires this length
func (d *ChangeDetector) ChangedFiles(beforeSnapshot string) ([]string, error) {
	if !d.executor.IsGitRepo(d.repoPath) {
		return nil, errors.New(ErrNotGitRepo)
	}

	// Parse the before snapshot
	beforeHead, beforeStash, beforeUntracked := parseSnapshot(beforeSnapshot)

	// Get current state
	currentStatus, err := d.executor.StatusPorcelain(d.repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	// Build set of files that were untracked before (to exclude them)
	beforeUntrackedSet := make(map[string]bool)
	for _, f := range beforeUntracked {
		beforeUntrackedSet[f] = true
	}

	// Build set of deleted files from current status (to exclude from results)
	deletedFiles := parseDeletedFilesFromStatus(currentStatus)

	// Collect all changed files
	changedFiles := make(map[string]bool)

	// 1. Get modified tracked files via git diff (filter out deleted)
	if beforeHead != "" {
		diffOutput, err := d.executor.DiffNameOnly(d.repoPath, beforeHead, "")
		if err != nil {
			// If diff fails (e.g., reference doesn't exist), continue with
			// status-based detection
			_ = err
		} else {
			for _, file := range parseLines(diffOutput) {
				if file != "" && !deletedFiles[file] {
					changedFiles[file] = true
				}
			}
		}
	}

	// 2. If we had a stash, compare against it (filter out deleted)
	if beforeStash != "" {
		stashDiff, err := d.executor.DiffNameOnly(d.repoPath, beforeStash, "")
		if err == nil {
			for _, file := range parseLines(stashDiff) {
				if file != "" && !deletedFiles[file] {
					changedFiles[file] = true
				}
			}
		}
	}

	// 3. Get currently modified/added files from status (not untracked, not deleted)
	// parseChangedFilesFromStatus already excludes deleted files
	currentModified := parseModifiedFilesFromStatus(currentStatus)
	for _, file := range currentModified {
		changedFiles[file] = true
	}

	// 4. Find new untracked files (files that weren't untracked before)
	currentUntracked := parseUntrackedFiles(currentStatus)
	for _, file := range currentUntracked {
		if !beforeUntrackedSet[file] {
			changedFiles[file] = true
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(changedFiles))
	for file := range changedFiles {
		result = append(result, file)
	}

	return result, nil
}

// RepoPath returns the repository path this detector is configured for.
func (d *ChangeDetector) RepoPath() string {
	return d.repoPath
}

// parseSnapshot parses a snapshot string into its components.
func parseSnapshot(snapshot string) (head, stash string, untracked []string) {
	parts := strings.Split(snapshot, "|")
	for _, part := range parts {
		switch {
		case strings.HasPrefix(part, "HEAD:"):
			head = strings.TrimPrefix(part, "HEAD:")
		case strings.HasPrefix(part, "STASH:"):
			stash = strings.TrimPrefix(part, "STASH:")
		case strings.HasPrefix(part, "UNTRACKED:"):
			files := strings.TrimPrefix(part, "UNTRACKED:")
			if files != "" {
				untracked = strings.Split(files, ",")
			}
		}
	}

	return head, stash, untracked
}

// parseUntrackedFiles extracts untracked files from git status output.
func parseUntrackedFiles(status string) []string {
	var files []string
	for _, line := range strings.Split(status, newline) {
		if len(line) < minPorcelainLineLen {
			continue
		}
		// Untracked files have "??" as the status - skip others
		if line[0] != '?' || line[1] != '?' {
			continue
		}
		filename := strings.TrimSpace(line[minPorcelainLineLen:])
		if filename != "" {
			files = append(files, filename)
		}
	}

	return files
}

// parseChangedFilesFromStatus extracts all changed files from git status
// --porcelain output. This includes modified, added, and untracked files.
// Deleted files are excluded.
//
//nolint:revive // cognitive-complexity - git status parsing logic
func parseChangedFilesFromStatus(status string) []string {
	var files []string
	for _, line := range strings.Split(status, newline) {
		if len(line) < minPorcelainLineLen {
			continue
		}

		indexStatus := line[0]
		workTreeStatus := line[1]

		// Skip deleted files
		if indexStatus == 'D' || workTreeStatus == 'D' {
			continue
		}

		// Skip unchanged files
		if indexStatus == ' ' && workTreeStatus == ' ' {
			continue
		}

		// Include: modified, added, renamed, copied, untracked
		filename := extractFilename(line)
		if filename != "" {
			files = append(files, filename)
		}
	}

	return files
}

// extractFilename extracts the filename from a git status line,
// handling renamed files (e.g., "R  old -> new").
func extractFilename(line string) string {
	filename := strings.TrimSpace(line[minPorcelainLineLen:])
	if filename == "" {
		return ""
	}
	// Handle renamed files: "R  old -> new"
	if idx := strings.Index(filename, " -> "); idx >= 0 {
		parts := strings.Split(filename, " -> ")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	return filename
}

// parseModifiedFilesFromStatus extracts modified/added files from git status
// --porcelain output. This excludes untracked files and deleted files.
//
//nolint:revive // cognitive-complexity - git status parsing logic
func parseModifiedFilesFromStatus(status string) []string {
	var files []string
	for _, line := range strings.Split(status, newline) {
		if len(line) < minPorcelainLineLen {
			continue
		}

		indexStatus := line[0]
		workTreeStatus := line[1]

		// Skip untracked files (they're handled separately)
		if indexStatus == '?' && workTreeStatus == '?' {
			continue
		}

		// Skip deleted files
		if indexStatus == 'D' || workTreeStatus == 'D' {
			continue
		}

		// Skip unchanged files
		if indexStatus == ' ' && workTreeStatus == ' ' {
			continue
		}

		// Include modified, added, renamed, copied
		filename := extractFilename(line)
		if filename != "" {
			files = append(files, filename)
		}
	}

	return files
}

// parseDeletedFilesFromStatus extracts deleted files from git status
// --porcelain output.
func parseDeletedFilesFromStatus(status string) map[string]bool {
	deleted := make(map[string]bool)
	for _, line := range strings.Split(status, newline) {
		if len(line) < minPorcelainLineLen {
			continue
		}

		indexStatus := line[0]
		workTreeStatus := line[1]

		// Skip non-deleted files
		if indexStatus != 'D' && workTreeStatus != 'D' {
			continue
		}

		// Deleted files have D in either position
		filename := strings.TrimSpace(line[minPorcelainLineLen:])
		if filename != "" {
			deleted[filename] = true
		}
	}

	return deleted
}

// parseLines splits output into non-empty lines.
func parseLines(output string) []string {
	var lines []string
	for _, line := range strings.Split(output, newline) {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}
