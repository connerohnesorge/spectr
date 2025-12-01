package git

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

// setupTestRepo creates a temporary git repository for testing.
// Returns the path to the repo and a cleanup function.
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git name: %v", err)
	}

	// Create initial commit
	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test"), testFilePerm); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Save original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	// Change to test repo directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to test directory: %v", err)
	}

	cleanup := func() {
		_ = os.Chdir(originalDir)
	}

	return tmpDir, cleanup
}

// setupTestRepoWithRemote creates a test repo with a simulated remote.
func setupTestRepoWithRemote(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	remoteDir := filepath.Join(tmpDir, "remote.git")
	workDir := filepath.Join(tmpDir, "work")

	// Create bare remote repo
	cmd := exec.Command("git", "init", "--bare", remoteDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init bare repo: %v", err)
	}

	// Create working repo
	if err := os.MkdirAll(workDir, testDirPerm); err != nil {
		t.Fatalf("failed to create work dir: %v", err)
	}

	cmd = exec.Command("git", "init")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init work repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git name: %v", err)
	}

	// Add remote
	cmd = exec.Command("git", "remote", "add", "origin", remoteDir)
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add remote: %v", err)
	}

	// Create initial commit
	testFile := filepath.Join(workDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test"), testFilePerm); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to add files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	// Push to remote (creates origin/main or origin/master)
	cmd = exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to push: %v", err)
	}

	// Fetch to update remote refs
	cmd = exec.Command("git", "fetch", "origin")
	cmd.Dir = workDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to fetch: %v", err)
	}

	// Save original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	// Change to work directory
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("failed to change to work directory: %v", err)
	}

	cleanup := func() {
		_ = os.Chdir(originalDir)
	}

	return workDir, cleanup
}

func TestCheckGitVersion(t *testing.T) {
	// This test requires git to be installed
	err := CheckGitVersion()
	if err != nil {
		// Check if git is just not installed vs version too old
		if strings.Contains(err.Error(), "not installed") {
			t.Skip("git is not installed, skipping version test")
		}
		t.Errorf("CheckGitVersion() error = %v", err)
	}
}

func TestCheckGitVersion_ParseVersion(t *testing.T) {
	// Test version parsing with various formats
	tests := []struct {
		name    string
		version string
		valid   bool
	}{
		{"modern git", "git version 2.39.3", true},
		{"old git", "git version 2.4.0", false},
		{"very old git", "git version 1.9.0", false},
		{"minimum version", "git version 2.5.0", true},
		{"git for Windows", "git version 2.39.3.windows.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := gitVersionPattern.FindStringSubmatch(tt.version)
			if len(matches) != versionMatchGroups {
				t.Errorf("failed to parse version string: %s", tt.version)

				return
			}

			// Verify parsing works correctly
			if matches[1] == "" || matches[2] == "" {
				t.Errorf("empty version components for: %s", tt.version)
			}
		})
	}
}

func TestGetBaseBranch(t *testing.T) {
	_, cleanup := setupTestRepoWithRemote(t)
	defer cleanup()

	branch, err := GetBaseBranch()
	if err != nil {
		t.Errorf("GetBaseBranch() error = %v", err)

		return
	}

	// Should return either "main" or "master"
	if branch != "main" && branch != "master" {
		t.Errorf("GetBaseBranch() = %v, want 'main' or 'master'", branch)
	}
}

func TestGetBaseBranch_NoRemote(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	_, err := GetBaseBranch()
	if err == nil {
		t.Error("GetBaseBranch() expected error for repo without remote, got nil")
	}
}

func TestCreateWorktree(t *testing.T) {
	_, cleanup := setupTestRepoWithRemote(t)
	defer cleanup()

	baseBranch, err := GetBaseBranch()
	if err != nil {
		t.Fatalf("GetBaseBranch() error = %v", err)
	}

	info, err := CreateWorktree(baseBranch, "test-feature-branch")
	if err != nil {
		t.Fatalf("CreateWorktree() error = %v", err)
	}

	// Verify worktree was created
	if _, err := os.Stat(info.Path); os.IsNotExist(err) {
		t.Errorf("worktree path does not exist: %s", info.Path)
	}

	// Verify branch name
	if info.Branch != "test-feature-branch" {
		t.Errorf("CreateWorktree() branch = %v, want %v", info.Branch, "test-feature-branch")
	}

	// Cleanup
	if err := CleanupWorktree(info.Path); err != nil {
		t.Errorf("CleanupWorktree() error = %v", err)
	}
}

func TestCreateWorktree_EmptyBaseBranch(t *testing.T) {
	_, err := CreateWorktree("", "new-branch")
	if err == nil {
		t.Error("CreateWorktree() expected error for empty baseBranch, got nil")
	}
}

func TestCreateWorktree_EmptyNewBranch(t *testing.T) {
	_, err := CreateWorktree("main", "")
	if err == nil {
		t.Error("CreateWorktree() expected error for empty newBranch, got nil")
	}
}

func TestCleanupWorktree(t *testing.T) {
	_, cleanup := setupTestRepoWithRemote(t)
	defer cleanup()

	baseBranch, err := GetBaseBranch()
	if err != nil {
		t.Fatalf("GetBaseBranch() error = %v", err)
	}

	info, err := CreateWorktree(baseBranch, "cleanup-test-branch")
	if err != nil {
		t.Fatalf("CreateWorktree() error = %v", err)
	}

	// Cleanup should succeed
	if err := CleanupWorktree(info.Path); err != nil {
		t.Errorf("CleanupWorktree() error = %v", err)
	}

	// Verify worktree was removed
	if _, err := os.Stat(info.Path); !os.IsNotExist(err) {
		t.Errorf("worktree path still exists after cleanup: %s", info.Path)
	}
}

func TestCleanupWorktree_EmptyPath(t *testing.T) {
	err := CleanupWorktree("")
	if err == nil {
		t.Error("CleanupWorktree() expected error for empty path, got nil")
	}
}

func TestCleanupWorktree_NonExistentPath(t *testing.T) {
	_, cleanup := setupTestRepoWithRemote(t)
	defer cleanup()

	// Should handle non-existent path gracefully
	err := CleanupWorktree("/nonexistent/path/that/does/not/exist")
	// This may or may not error depending on git behavior
	// The important thing is it doesn't panic
	_ = err
}

func TestCheckBranchExists(t *testing.T) {
	_, cleanup := setupTestRepoWithRemote(t)
	defer cleanup()

	baseBranch, err := GetBaseBranch()
	if err != nil {
		t.Fatalf("GetBaseBranch() error = %v", err)
	}

	// Should find the base branch
	exists, err := CheckBranchExists(baseBranch)
	if err != nil {
		t.Fatalf("CheckBranchExists() error = %v", err)
	}
	if !exists {
		t.Errorf("CheckBranchExists(%q) = false, want true", baseBranch)
	}
}

func TestCheckBranchExists_NonExistent(t *testing.T) {
	_, cleanup := setupTestRepoWithRemote(t)
	defer cleanup()

	exists, err := CheckBranchExists("nonexistent-branch-12345")
	if err != nil {
		t.Fatalf("CheckBranchExists() error = %v", err)
	}
	if exists {
		t.Error("CheckBranchExists() = true for nonexistent branch, want false")
	}
}

func TestCheckBranchExists_EmptyBranch(t *testing.T) {
	_, err := CheckBranchExists("")
	if err == nil {
		t.Error("CheckBranchExists() expected error for empty branch, got nil")
	}
}

func TestWorktreeInfo(t *testing.T) {
	info := WorktreeInfo{
		Path:   "/tmp/test-worktree",
		Branch: "feature-branch",
	}

	if info.Path != "/tmp/test-worktree" {
		t.Errorf("WorktreeInfo.Path = %v, want %v", info.Path, "/tmp/test-worktree")
	}
	if info.Branch != "feature-branch" {
		t.Errorf("WorktreeInfo.Branch = %v, want %v", info.Branch, "feature-branch")
	}
}
