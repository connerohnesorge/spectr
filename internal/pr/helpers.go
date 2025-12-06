// Package pr provides pull request creation and management functionality.
// This file contains helper functions for the PR workflow, including
// CLI tool verification, archive execution, and git operations.
package pr

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
)

// dirPerm is the default permission for created directories (rwxr-xr-x).
const dirPerm = 0755

// checkCLITool verifies that the required CLI tool is installed.
// Returns an error with installation suggestions if not found.
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
// Supports gh (GitHub CLI), glab (GitLab CLI), and tea (Gitea CLI).
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
// It copies the change files to the worktree and executes the archive command.
func executeArchiveInWorktree(
	config PRConfig,
	worktreePath string,
) (archive.ArchiveResult, error) {
	fmt.Println("Running archive operation in worktree...")

	if err := copyChangeToWorktree(config, worktreePath); err != nil {
		return archive.ArchiveResult{},
			fmt.Errorf("copy change to worktree: %w", err)
	}

	archiveCmd := &archive.ArchiveCmd{
		ChangeID:  config.ChangeID,
		Yes:       true,
		SkipSpecs: config.SkipSpecs,
	}

	result, err := archive.Archive(archiveCmd, worktreePath)
	if err != nil {
		return archive.ArchiveResult{}, err
	}

	return result, nil
}

// stageAndCommit stages the spectr/ directory and creates a commit.
// It runs git add and git commit within the worktree directory.
func stageAndCommit(worktreePath, commitMsg string) error {
	fmt.Println("Staging changes...")

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
// It runs git push -u origin within the worktree directory.
func pushBranch(worktreePath, branchName string) error {
	fmt.Printf("Pushing branch: %s\n", branchName)

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
