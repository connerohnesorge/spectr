// Package pr provides helper functions for the pull request workflow.
// These helpers handle CLI tool verification, file operations for worktrees,
// git staging/committing, and spec deletion operations.
package pr

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/git"
)

// dirPerm is the default permission for created directories.
const dirPerm = 0755

// checkCLITool verifies that the required CLI tool is installed.
func checkCLITool(tool string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		suggestions := getCLIInstallSuggestion(tool)

		return fmt.Errorf(
			"CLI tool '%s' not found in PATH; %s",
			tool,
			suggestions,
		)
	}

	return nil
}

// getCLIInstallSuggestion returns installation suggestions for a CLI tool.
func getCLIInstallSuggestion(tool string) string {
	switch tool {
	case "gh":
		return "Install: brew install gh or see https://cli.github.com"
	case "glab":
		return "Install: brew install glab or see " +
			"https://gitlab.com/gitlab-org/cli"
	case "tea":
		return "Install: brew install tea or see https://gitea.com/gitea/tea"
	default:
		return "please install the required CLI tool"
	}
}

// executeArchiveInWorktree runs the archive workflow within the worktree.
func executeArchiveInWorktree(
	config PRConfig,
	worktreePath string,
) (archive.ArchiveResult, error) {
	fmt.Println("Running archive operation in worktree...")

	// Copy the change from source to worktree first
	if err := copyChangeToWorktree(config, worktreePath); err != nil {
		return archive.ArchiveResult{},
			fmt.Errorf("copy change to worktree: %w", err)
	}

	// Create archive command
	archiveCmd := &archive.ArchiveCmd{
		ChangeID:  config.ChangeID,
		Yes:       true, // Non-interactive
		SkipSpecs: config.SkipSpecs,
	}

	// Execute archive within the worktree and capture results.
	// The ArchiveResult contains path, operation counts, and capabilities.
	result, err := archive.Archive(archiveCmd, worktreePath)
	if err != nil {
		return archive.ArchiveResult{}, err
	}

	return result, nil
}

// copyChangeToWorktree copies the change directory from source to worktree.
func copyChangeToWorktree(config PRConfig, worktreePath string) error {
	projectRoot := config.ProjectRoot
	if projectRoot == "" {
		var err error
		projectRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf("get repo root: %w", err)
		}
	}

	sourceDir := filepath.Join(
		projectRoot, "spectr", "changes", config.ChangeID,
	)
	targetDir := filepath.Join(
		worktreePath, "spectr", "changes", config.ChangeID,
	)

	fmt.Printf("Copying change to worktree: %s\n", config.ChangeID)

	// Create target directory structure
	if err := os.MkdirAll(filepath.Dir(targetDir), dirPerm); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	// Copy directory recursively
	if err := copyDir(sourceDir, targetDir); err != nil {
		return fmt.Errorf("copy directory: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory from src to dst.
// It preserves the original directory permissions and handles nested
// directories by recursively calling itself for subdirectories.
func copyDir(src, dst string) error {
	// Get source directory info for permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read all entries in source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry (file or directory)
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy regular files
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst.
// It preserves the original file permissions from the source file.
func copyFile(src, dst string) error {
	// Open source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file with same permissions
	dstFile, err := os.OpenFile(
		dst,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		srcInfo.Mode(),
	)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	// Copy file contents
	_, err = io.Copy(dstFile, srcFile)

	return err
}

// stageAndCommit stages the spectr/ directory and creates a commit.
// It runs git add on the spectr/ directory and then creates a commit
// with the provided message. Both operations are executed in the
// specified worktree path.
func stageAndCommit(worktreePath, commitMsg string) error {
	fmt.Println("Staging changes...")

	// Stage all changes in the spectr directory
	addCmd := exec.Command("git", "add", "spectr/")
	addCmd.Dir = worktreePath

	output, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git add failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	fmt.Println("Creating commit...")

	// Create commit with the provided message
	commitCmd := exec.Command("git", "commit", "-m", commitMsg)
	commitCmd.Dir = worktreePath

	output, err = commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git commit failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// pushBranch pushes the branch to origin with upstream tracking.
// It uses the -u flag to set up the remote tracking branch for
// future push/pull operations.
func pushBranch(worktreePath, branchName string) error {
	fmt.Printf("Pushing branch: %s\n", branchName)

	// Push with upstream tracking enabled
	cmd := exec.Command("git", "push", "-u", "origin", branchName)
	cmd.Dir = worktreePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git push failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// executeDeleteInWorktree removes the spec directory from the worktree.
// This is used in the PR delete workflow to remove a spec before
// creating a pull request for the deletion.
func executeDeleteInWorktree(specID, worktreePath string) error {
	// Build path to spec directory
	specDir := filepath.Join(worktreePath, "spectr", "specs", specID)

	fmt.Printf("Removing spec directory: spectr/specs/%s\n", specID)

	// Remove the entire directory tree
	if err := os.RemoveAll(specDir); err != nil {
		return fmt.Errorf("remove spec directory: %w", err)
	}

	return nil
}
