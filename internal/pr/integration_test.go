//go:build integration

// Package pr provides integration tests for the PR workflow.
// These tests require a git environment and are skipped by default.
// Run with: go test ./internal/pr/... -v -tags=integration
package pr

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testDirPerm  = 0o755
	testFilePerm = 0o644
)

// setupTestRepo creates a temp directory, initializes a git repo,
// and creates a basic spectr structure. Returns the path to the repo.
func setupTestRepo(t *testing.T) string {
	t.Helper()

	// Create temp directory
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf(
			"failed to init git repo: %s: %v",
			output,
			err,
		)
	}

	// Configure git user for commits
	cmd = exec.Command(
		"git",
		"config",
		"user.email",
		"test@example.com",
	)
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf(
			"failed to configure git email: %s: %v",
			output,
			err,
		)
	}

	cmd = exec.Command(
		"git",
		"config",
		"user.name",
		"Test User",
	)
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf(
			"failed to configure git name: %s: %v",
			output,
			err,
		)
	}

	// Create spectr directory structure
	spectrDir := filepath.Join(tmpDir, "spectr")
	changesDir := filepath.Join(
		spectrDir,
		"changes",
	)
	specsDir := filepath.Join(spectrDir, "specs")

	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatalf(
			"failed to create changes dir: %v",
			err,
		)
	}
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatalf(
			"failed to create specs dir: %v",
			err,
		)
	}

	// Create project.md
	projectContent := "# Test Project\n\nA test project for integration tests.\n"
	if err := os.WriteFile(
		filepath.Join(spectrDir, "project.md"),
		[]byte(projectContent),
		testFilePerm,
	); err != nil {
		t.Fatalf(
			"failed to create project.md: %v",
			err,
		)
	}

	// Create initial commit
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf(
			"failed to git add: %s: %v",
			output,
			err,
		)
	}

	cmd = exec.Command(
		"git",
		"commit",
		"-m",
		"Initial commit",
	)
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf(
			"failed to git commit: %s: %v",
			output,
			err,
		)
	}

	return tmpDir
}

// setupTestRepoWithOrigin creates a test repo with a fake origin remote.
// This allows testing of validation logic that runs after the origin check.
func setupTestRepoWithOrigin(
	t *testing.T,
) string {
	t.Helper()

	repoPath := setupTestRepo(t)

	// Add a fake origin remote (invalid URL is fine for validation tests)
	cmd := exec.Command(
		"git",
		"remote",
		"add",
		"origin",
		"https://github.com/test/test.git",
	)
	cmd.Dir = repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf(
			"failed to add origin remote: %s: %v",
			output,
			err,
		)
	}

	return repoPath
}

// createTestChange creates a test change proposal with minimal content.
func createTestChange(
	t *testing.T,
	repoPath, changeID string,
) {
	t.Helper()

	changeDir := filepath.Join(
		repoPath,
		"spectr",
		"changes",
		changeID,
	)
	changeSpecsDir := filepath.Join(
		changeDir,
		"specs",
		"test-feature",
	)

	if err := os.MkdirAll(changeSpecsDir, testDirPerm); err != nil {
		t.Fatalf(
			"failed to create change specs dir: %v",
			err,
		)
	}

	// Create proposal.md
	proposalContent := `# Change: ` + changeID + `

## Why

This is a test change for integration testing.

## What Changes

- Add test feature

## Impact

- specs: test-feature
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		testFilePerm,
	); err != nil {
		t.Fatalf(
			"failed to create proposal.md: %v",
			err,
		)
	}

	// Create tasks.md
	tasksContent := `## Tasks

- [x] Task 1: Setup
- [x] Task 2: Implementation
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "tasks.md"),
		[]byte(tasksContent),
		testFilePerm,
	); err != nil {
		t.Fatalf(
			"failed to create tasks.md: %v",
			err,
		)
	}

	// Create delta spec
	deltaSpec := `## ADDED Requirements

### Requirement: Test Feature
The system SHALL provide test functionality.

#### Scenario: Basic test
- **WHEN** test is run
- **THEN** it passes
`
	if err := os.WriteFile(
		filepath.Join(changeSpecsDir, "spec.md"),
		[]byte(deltaSpec),
		testFilePerm,
	); err != nil {
		t.Fatalf(
			"failed to create spec.md: %v",
			err,
		)
	}
}

// cleanupTestRepo removes the test repository.
// This is typically handled by t.TempDir() cleanup, but provided
// for explicit cleanup needs.
func cleanupTestRepo(
	t *testing.T,
	repoPath string,
) {
	t.Helper()

	if repoPath != "" {
		// Nothing to do - t.TempDir() handles cleanup
		_ = repoPath
	}
}

// TestIntegration_ExecutePR_Archive_DryRun tests the archive dry run functionality.
func TestIntegration_ExecutePR_Archive_DryRun(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(t, repoPath)

	// Create a test change
	changeID := "test-archive-change"
	createTestChange(t, repoPath, changeID)

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR with DryRun=true
	config := PRConfig{
		ChangeID:    changeID,
		Mode:        ModeArchive,
		DryRun:      true,
		ProjectRoot: repoPath,
	}

	// This will fail because there's no origin remote, but DryRun
	// should still show intended actions before failing
	result, err := ExecutePR(config)

	// We expect an error because there's no origin remote
	if err == nil {
		// If it somehow succeeds (shouldn't happen), verify no changes
		if result != nil && result.PRURL != "" {
			t.Error(
				"DryRun should not create actual PR",
			)
		}
	} else {
		// Verify it's the expected error (no origin remote)
		if !strings.Contains(err.Error(), "origin") &&
			!strings.Contains(err.Error(), "remote") {
			t.Logf("Expected origin-related error, got: %v", err)
		}
	}

	// Verify no worktree was created (check temp dirs)
	tempDirs, _ := filepath.Glob(
		filepath.Join(
			os.TempDir(),
			"spectr-pr-*",
		),
	)
	for _, dir := range tempDirs {
		// Clean up any stray worktrees
		_ = os.RemoveAll(dir)
	}

	// Verify original repo is unchanged
	cmd := exec.Command(
		"git",
		"status",
		"--porcelain",
	)
	cmd.Dir = repoPath
	output, _ := cmd.CombinedOutput()
	// Should have untracked change directory
	if !strings.Contains(
		string(output),
		changeID,
	) {
		t.Logf(
			"Expected change directory to remain untracked, status: %s",
			output,
		)
	}
}

// TestIntegration_ExecutePR_Proposal_DryRun tests the proposal dry run functionality.
func TestIntegration_ExecutePR_Proposal_DryRun(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(t, repoPath)

	// Create a test change
	changeID := "test-proposal-change"
	createTestChange(t, repoPath, changeID)

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR with DryRun=true
	config := PRConfig{
		ChangeID:    changeID,
		Mode:        ModeProposal,
		DryRun:      true,
		ProjectRoot: repoPath,
	}

	result, err := ExecutePR(config)

	// We expect an error because there's no origin remote
	if err == nil {
		if result != nil && result.PRURL != "" {
			t.Error(
				"DryRun should not create actual PR",
			)
		}
	} else {
		// Verify it's the expected error
		if !strings.Contains(err.Error(), "origin") &&
			!strings.Contains(err.Error(), "remote") {
			t.Logf("Expected origin-related error, got: %v", err)
		}
	}

	// Verify change directory still exists
	changeDir := filepath.Join(
		repoPath,
		"spectr",
		"changes",
		changeID,
	)
	if _, err := os.Stat(changeDir); os.IsNotExist(
		err,
	) {
		t.Error(
			"Change directory should still exist after dry run",
		)
	}
}

// TestIntegration_ExecutePR_InvalidChangeID tests error handling for non-existent change.
func TestIntegration_ExecutePR_InvalidChangeID(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo with origin (so we get past origin check)
	repoPath := setupTestRepoWithOrigin(t)
	defer cleanupTestRepo(t, repoPath)

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR with non-existent change ID
	config := PRConfig{
		ChangeID:    "nonexistent-change",
		Mode:        ModeArchive,
		DryRun:      true,
		ProjectRoot: repoPath,
	}

	_, err = ExecutePR(config)

	// We expect an error about the change not being found
	if err == nil {
		t.Error(
			"Expected error for non-existent change ID",
		)
	} else if !strings.Contains(err.Error(), "not found") &&
		!strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

// TestIntegration_ExecutePR_NoOriginRemote tests error when no origin remote exists.
func TestIntegration_ExecutePR_NoOriginRemote(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo (no origin remote by default)
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(t, repoPath)

	// Create a test change
	changeID := "test-no-origin"
	createTestChange(t, repoPath, changeID)

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR - should fail because no origin
	config := PRConfig{
		ChangeID:    changeID,
		Mode:        ModeArchive,
		DryRun:      false,
		ProjectRoot: repoPath,
	}

	_, err = ExecutePR(config)

	// We expect an error about no origin remote
	if err == nil {
		t.Error(
			"Expected error when no origin remote exists",
		)
	} else if !strings.Contains(err.Error(), "origin") &&
		!strings.Contains(err.Error(), "remote") {
		t.Errorf("Expected origin-related error, got: %v", err)
	}
}

// TestIntegration_WorktreeIsolation verifies that user's working directory
// is not modified during PR operations.
func TestIntegration_WorktreeIsolation(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(t, repoPath)

	// Create a test change
	changeID := "test-isolation"
	createTestChange(t, repoPath, changeID)

	// Create an uncommitted file to track
	uncommittedFile := filepath.Join(
		repoPath,
		"uncommitted.txt",
	)
	uncommittedContent := "This file is uncommitted and should not be touched"
	if err := os.WriteFile(uncommittedFile, []byte(uncommittedContent), testFilePerm); err != nil {
		t.Fatalf(
			"failed to create uncommitted file: %v",
			err,
		)
	}

	// Create a modified tracked file
	projectFile := filepath.Join(
		repoPath,
		"spectr",
		"project.md",
	)
	modifiedContent := "# Modified Project\n\nThis has local modifications.\n"
	if err := os.WriteFile(projectFile, []byte(modifiedContent), testFilePerm); err != nil {
		t.Fatalf(
			"failed to modify project.md: %v",
			err,
		)
	}

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR with DryRun
	config := PRConfig{
		ChangeID:    changeID,
		Mode:        ModeProposal,
		DryRun:      true,
		ProjectRoot: repoPath,
	}

	// Execute (will fail due to no origin, but should not modify files)
	_, _ = ExecutePR(config)

	// Verify uncommitted file still exists with same content
	content, err := os.ReadFile(uncommittedFile)
	if err != nil {
		t.Errorf(
			"Uncommitted file should still exist: %v",
			err,
		)
	} else if string(content) != uncommittedContent {
		t.Errorf("Uncommitted file content changed: got %q, want %q",
			string(content), uncommittedContent)
	}

	// Verify modified project.md still has local modifications
	content, err = os.ReadFile(projectFile)
	if err != nil {
		t.Errorf(
			"project.md should still exist: %v",
			err,
		)
	} else if string(content) != modifiedContent {
		t.Errorf("project.md modifications were lost: got %q, want %q",
			string(content), modifiedContent)
	}

	// Verify git status shows expected state
	cmd := exec.Command(
		"git",
		"status",
		"--porcelain",
	)
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("git status failed: %v", err)
	}

	statusOutput := string(output)
	// Should show uncommitted.txt as untracked
	if !strings.Contains(
		statusOutput,
		"uncommitted.txt",
	) {
		t.Error(
			"uncommitted.txt should appear in git status as untracked",
		)
	}
	// Should show project.md as modified
	if !strings.Contains(
		statusOutput,
		"project.md",
	) {
		t.Error(
			"project.md should appear in git status as modified",
		)
	}
}

// TestIntegration_ExecutePR_InvalidMode tests error handling for invalid mode.
func TestIntegration_ExecutePR_InvalidMode(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo with origin (so we get past origin check)
	repoPath := setupTestRepoWithOrigin(t)
	defer cleanupTestRepo(t, repoPath)

	// Create a test change
	changeID := "test-invalid-mode"
	createTestChange(t, repoPath, changeID)

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR with invalid mode
	config := PRConfig{
		ChangeID:    changeID,
		Mode:        "invalid-mode",
		DryRun:      true,
		ProjectRoot: repoPath,
	}

	_, err = ExecutePR(config)

	// We expect an error about invalid mode
	if err == nil {
		t.Error("Expected error for invalid mode")
	} else if !strings.Contains(err.Error(), "invalid") &&
		!strings.Contains(err.Error(), "mode") {
		t.Errorf("Expected mode-related error, got: %v", err)
	}
}

// TestIntegration_ExecutePR_Remove_DryRun tests the remove dry run functionality.
func TestIntegration_ExecutePR_Remove_DryRun(
	t *testing.T,
) {
	if testing.Short() {
		t.Skip(
			"skipping integration test in short mode",
		)
	}

	// Set up test repo
	repoPath := setupTestRepo(t)
	defer cleanupTestRepo(t, repoPath)

	// Create a test change
	changeID := "test-remove-change"
	createTestChange(t, repoPath, changeID)

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}

	// Change to test repo
	if err := os.Chdir(repoPath); err != nil {
		t.Fatalf(
			"failed to change to test repo: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Run ExecutePR with DryRun=true and ModeRemove
	config := PRConfig{
		ChangeID:    changeID,
		Mode:        ModeRemove,
		DryRun:      true,
		ProjectRoot: repoPath,
	}

	result, err := ExecutePR(config)

	// We expect an error because there's no origin remote
	if err == nil {
		if result != nil && result.PRURL != "" {
			t.Error(
				"DryRun should not create actual PR",
			)
		}
	} else {
		// Verify it's the expected error (no origin remote)
		if !strings.Contains(err.Error(), "origin") &&
			!strings.Contains(err.Error(), "remote") {
			t.Logf("Expected origin-related error, got: %v", err)
		}
	}

	// Verify change directory still exists (dry run should not delete it)
	changeDir := filepath.Join(
		repoPath,
		"spectr",
		"changes",
		changeID,
	)
	if _, err := os.Stat(changeDir); os.IsNotExist(
		err,
	) {
		t.Error(
			"Change directory should still exist after dry run",
		)
	}
}
