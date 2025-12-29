package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const testOriginMain = "origin/main"

// isGitAvailable checks if git command is available.
func isGitAvailable() bool {
	_, err := exec.LookPath("git")

	return err == nil
}

// isInGitRepo checks if we're currently in a git repository.
func isInGitRepo() bool {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--git-dir",
	)

	return cmd.Run() == nil
}

// hasOriginRemote checks if the origin remote is configured.
func hasOriginRemote() bool {
	cmd := exec.Command(
		"git",
		"remote",
		"get-url",
		"origin",
	)

	return cmd.Run() == nil
}

func TestGetRepoRoot(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}

	root, err := GetRepoRoot()
	if err != nil {
		t.Fatalf("GetRepoRoot() error = %v", err)
	}

	// Root should be a non-empty absolute path
	if root == "" {
		t.Error(
			"GetRepoRoot() returned empty string",
		)
	}

	if !filepath.IsAbs(root) {
		t.Errorf(
			"GetRepoRoot() returned non-absolute path: %s",
			root,
		)
	}

	// The root should contain a .git directory
	gitDir := filepath.Join(root, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		t.Errorf(
			"GetRepoRoot() returned path without .git: %s, error: %v",
			root,
			err,
		)
	} else if !info.IsDir() {
		// .git could be a file (gitdir reference for worktrees)
		// This is fine, just check it exists
		if !info.Mode().IsRegular() {
			t.Errorf(".git is neither a directory nor a file: %s", gitDir)
		}
	}
}

func TestGetRepoRoot_NotInRepo(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}

	// Create a temp directory that is not a git repo
	tempDir, err := os.MkdirTemp(
		"",
		"not-a-repo-*",
	)
	if err != nil {
		t.Fatalf(
			"failed to create temp dir: %v",
			err,
		)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Change to the temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf(
			"failed to change to temp dir: %v",
			err,
		)
	}

	_, err = GetRepoRoot()
	if err == nil {
		t.Error(
			"GetRepoRoot() expected error in non-git directory, got nil",
		)
	}
	if !strings.Contains(
		err.Error(),
		"not a git repository",
	) {
		t.Errorf(
			"GetRepoRoot() error = %v, want error containing 'not a git repository'",
			err,
		)
	}
}

func TestGetBaseBranch(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	t.Run("auto-detect", func(t *testing.T) {
		branch, err := GetBaseBranch("")
		if err != nil {
			t.Fatalf(
				"GetBaseBranch(\"\") error = %v",
				err,
			)
		}

		// Should return origin/main or origin/master
		if branch != testOriginMain &&
			branch != "origin/master" {
			t.Errorf(
				"GetBaseBranch(\"\") = %v, want origin/main or origin/master",
				branch,
			)
		}
	})

	t.Run(
		"with valid branch",
		func(t *testing.T) {
			// First detect what branch exists
			baseBranch, err := GetBaseBranch("")
			if err != nil {
				t.Skip(
					"could not auto-detect base branch",
				)
			}

			// Extract just the branch name (remove origin/)
			branchName := strings.TrimPrefix(
				baseBranch,
				"origin/",
			)

			// Now test with the explicit branch name
			branch, err := GetBaseBranch(
				branchName,
			)
			if err != nil {
				t.Fatalf(
					"GetBaseBranch(%q) error = %v",
					branchName,
					err,
				)
			}

			if branch != baseBranch {
				t.Errorf(
					"GetBaseBranch(%q) = %v, want %v",
					branchName,
					branch,
					baseBranch,
				)
			}
		},
	)

	t.Run(
		"with invalid branch",
		func(t *testing.T) {
			_, err := GetBaseBranch(
				"nonexistent-branch-12345-xyz",
			)
			if err == nil {
				t.Error(
					"GetBaseBranch(\"nonexistent-branch-12345-xyz\") expected error, got nil",
				)
			}
			if !strings.Contains(
				err.Error(),
				"does not exist",
			) {
				t.Errorf(
					"GetBaseBranch error = %v, want error containing 'does not exist'",
					err,
				)
			}
		},
	)
}

func TestBranchExists(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	t.Run("existing branch", func(t *testing.T) {
		// First detect what branch exists
		baseBranch, err := GetBaseBranch("")
		if err != nil {
			t.Skip(
				"could not auto-detect base branch",
			)
		}

		// Extract just the branch name (remove origin/)
		branchName := strings.TrimPrefix(
			baseBranch,
			"origin/",
		)

		exists, err := BranchExists(branchName)
		if err != nil {
			t.Fatalf(
				"BranchExists(%q) error = %v",
				branchName,
				err,
			)
		}
		if !exists {
			t.Errorf(
				"BranchExists(%q) = false, want true",
				branchName,
			)
		}
	})

	t.Run(
		"non-existent branch",
		func(t *testing.T) {
			exists, err := BranchExists(
				"nonexistent-branch-12345-xyz",
			)
			if err != nil {
				t.Fatalf(
					"BranchExists(\"nonexistent-branch-12345-xyz\") error = %v",
					err,
				)
			}
			if exists {
				t.Error(
					"BranchExists(\"nonexistent-branch-12345-xyz\") = true, want false",
				)
			}
		},
	)
}

func TestCreateWorktree_ValidationErrors(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}

	tests := []struct {
		name    string
		config  WorktreeConfig
		wantErr string
	}{
		{
			name: "empty branch name",
			config: WorktreeConfig{
				BranchName: "",
				BaseBranch: "origin/main",
			},
			wantErr: "branch name is required",
		},
		{
			name: "empty base branch",
			config: WorktreeConfig{
				BranchName: "test-branch",
				BaseBranch: "",
			},
			wantErr: "base branch is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateWorktree(tt.config)
			if err == nil {
				t.Fatal(
					"CreateWorktree() expected error, got nil",
				)
			}
			if !strings.Contains(
				err.Error(),
				tt.wantErr,
			) {
				t.Errorf(
					"CreateWorktree() error = %v, want error containing %q",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestCreateWorktree_InvalidBaseBranch(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}

	config := WorktreeConfig{
		BranchName: "test-branch-invalid-base",
		BaseBranch: "origin/nonexistent-branch-12345",
	}

	_, err := CreateWorktree(config)
	if err == nil {
		t.Fatal(
			"CreateWorktree() with invalid base branch expected error, got nil",
		)
	}

	// Error should mention the base branch not being found
	if !strings.Contains(
		err.Error(),
		"not found",
	) &&
		!strings.Contains(
			err.Error(),
			"failed to create worktree",
		) {
		t.Errorf(
			"CreateWorktree() error = %v, want error about base branch not found",
			err,
		)
	}
}

func TestCreateWorktree_And_Cleanup_Integration(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	// First detect what base branch exists
	baseBranch, err := GetBaseBranch("")
	if err != nil {
		t.Skipf(
			"could not auto-detect base branch: %v",
			err,
		)
	}

	// Create a unique branch name for testing
	suffix, err := randomHex(4)
	if err != nil {
		t.Fatalf(
			"failed to generate random suffix: %v",
			err,
		)
	}
	testBranchName := "spectr-test-worktree-" + suffix

	config := WorktreeConfig{
		BranchName: testBranchName,
		BaseBranch: baseBranch,
	}

	// Create the worktree
	info, err := CreateWorktree(config)
	if err != nil {
		t.Fatalf(
			"CreateWorktree() error = %v",
			err,
		)
	}

	// Verify the worktree info
	if info == nil {
		t.Fatal(
			"CreateWorktree() returned nil info",
		)
	}
	if info.Path == "" {
		t.Error(
			"CreateWorktree() returned empty path",
		)
	}
	if info.BranchName != testBranchName {
		t.Errorf(
			"CreateWorktree() BranchName = %v, want %v",
			info.BranchName,
			testBranchName,
		)
	}
	if !info.TempDir {
		t.Error(
			"CreateWorktree() TempDir = false, want true",
		)
	}

	// Verify the worktree directory exists
	if _, err := os.Stat(info.Path); os.IsNotExist(
		err,
	) {
		t.Errorf(
			"worktree directory does not exist: %s",
			info.Path,
		)
	}

	// Verify it's a valid git worktree (has .git file)
	gitPath := filepath.Join(info.Path, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(
		err,
	) {
		t.Errorf(
			"worktree missing .git: %s",
			gitPath,
		)
	}

	// Cleanup the worktree
	err = CleanupWorktree(info)
	if err != nil {
		t.Errorf(
			"CleanupWorktree() error = %v",
			err,
		)
	}

	// Verify the worktree directory is removed
	if _, err := os.Stat(info.Path); !os.IsNotExist(
		err,
	) {
		t.Errorf(
			"worktree directory still exists after cleanup: %s",
			info.Path,
		)
		// Clean up manually if test failed
		_ = os.RemoveAll(info.Path)
	}

	// Verify cleanup is idempotent (calling again should not error)
	err = CleanupWorktree(info)
	if err != nil {
		t.Errorf(
			"CleanupWorktree() second call error = %v, want nil",
			err,
		)
	}
}

func TestCleanupWorktree_NilInfo(t *testing.T) {
	// Cleanup with nil info should not error
	err := CleanupWorktree(nil)
	if err != nil {
		t.Errorf(
			"CleanupWorktree(nil) error = %v, want nil",
			err,
		)
	}
}

func TestCleanupWorktree_EmptyInfo(t *testing.T) {
	// Cleanup with empty info should not error
	err := CleanupWorktree(&WorktreeInfo{})
	if err != nil {
		t.Errorf(
			"CleanupWorktree(&WorktreeInfo{}) error = %v, want nil",
			err,
		)
	}
}

func TestRandomHex(t *testing.T) {
	tests := []struct {
		name       string
		n          int
		wantHexLen int
	}{
		{"1 byte", 1, 2},
		{"4 bytes", 4, 8},
		{"8 bytes", 8, 16},
		{"16 bytes", 16, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := randomHex(tt.n)
			if err != nil {
				t.Fatalf(
					"randomHex(%d) error = %v",
					tt.n,
					err,
				)
			}

			if len(result) != tt.wantHexLen {
				t.Errorf(
					"randomHex(%d) length = %d, want %d",
					tt.n,
					len(result),
					tt.wantHexLen,
				)
			}

			// Verify it's valid hex
			for _, c := range result {
				isDigit := c >= '0' && c <= '9'
				isHexLetter := c >= 'a' &&
					c <= 'f'
				if !isDigit && !isHexLetter {
					t.Errorf(
						"randomHex(%d) contains non-hex character: %c",
						tt.n,
						c,
					)
				}
			}
		})
	}
}

func TestRandomHex_Uniqueness(t *testing.T) {
	// Generate multiple values and verify they're unique
	const iterations = 100
	seen := make(map[string]bool)
	for range iterations {
		result, err := randomHex(4)
		if err != nil {
			t.Fatalf(
				"randomHex(4) error = %v",
				err,
			)
		}
		if seen[result] {
			t.Errorf(
				"randomHex(4) generated duplicate value: %s",
				result,
			)
		}
		seen[result] = true
	}
}

func TestWorktreeConfig_Struct(t *testing.T) {
	// Test that WorktreeConfig can be constructed properly
	config := WorktreeConfig{
		BranchName: "test-branch",
		BaseBranch: testOriginMain,
	}

	if config.BranchName != "test-branch" {
		t.Errorf(
			"WorktreeConfig.BranchName = %v, want test-branch",
			config.BranchName,
		)
	}
	if config.BaseBranch != testOriginMain {
		t.Errorf(
			"WorktreeConfig.BaseBranch = %v, want origin/main",
			config.BaseBranch,
		)
	}
}

func TestWorktreeInfo_Struct(t *testing.T) {
	// Test that WorktreeInfo can be constructed properly
	info := WorktreeInfo{
		Path:       "/tmp/test-worktree",
		BranchName: "test-branch",
		TempDir:    true,
	}

	if info.Path != "/tmp/test-worktree" {
		t.Errorf(
			"WorktreeInfo.Path = %v, want /tmp/test-worktree",
			info.Path,
		)
	}
	if info.BranchName != "test-branch" {
		t.Errorf(
			"WorktreeInfo.BranchName = %v, want test-branch",
			info.BranchName,
		)
	}
	if !info.TempDir {
		t.Error(
			"WorktreeInfo.TempDir = false, want true",
		)
	}
}

func TestCreateWorktree_UniquePathGeneration(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	// First detect what base branch exists
	baseBranch, err := GetBaseBranch("")
	if err != nil {
		t.Skipf(
			"could not auto-detect base branch: %v",
			err,
		)
	}

	// Create two worktrees with different branch names
	suffix1, _ := randomHex(4)
	suffix2, _ := randomHex(4)
	branch1 := "spectr-test-unique-1-" + suffix1
	branch2 := "spectr-test-unique-2-" + suffix2

	config1 := WorktreeConfig{
		BranchName: branch1,
		BaseBranch: baseBranch,
	}
	config2 := WorktreeConfig{
		BranchName: branch2,
		BaseBranch: baseBranch,
	}

	info1, err := CreateWorktree(config1)
	if err != nil {
		t.Fatalf(
			"CreateWorktree(config1) error = %v",
			err,
		)
	}
	defer func() { _ = CleanupWorktree(info1) }()

	info2, err := CreateWorktree(config2)
	if err != nil {
		t.Fatalf(
			"CreateWorktree(config2) error = %v",
			err,
		)
	}
	defer func() { _ = CleanupWorktree(info2) }()

	// Verify paths are different
	if info1.Path == info2.Path {
		t.Error(
			"CreateWorktree generated same path for two worktrees",
		)
	}

	// Both should contain "spectr-pr-" in the path
	if !strings.Contains(
		info1.Path,
		"spectr-pr-",
	) {
		t.Errorf(
			"CreateWorktree path does not contain 'spectr-pr-': %s",
			info1.Path,
		)
	}
	if !strings.Contains(
		info2.Path,
		"spectr-pr-",
	) {
		t.Errorf(
			"CreateWorktree path does not contain 'spectr-pr-': %s",
			info2.Path,
		)
	}
}

func TestPathExistsOnRef_ExistingPath(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	// Change to the repo root for consistent path resolution
	repoRoot, err := GetRepoRoot()
	if err != nil {
		t.Fatalf("GetRepoRoot() error = %v", err)
	}
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf(
			"failed to change to repo root: %v",
			err,
		)
	}

	// Auto-detect the base branch
	baseBranch, err := GetBaseBranch("")
	if err != nil {
		t.Skipf(
			"could not auto-detect base branch: %v",
			err,
		)
	}

	// Test with known paths that should exist in the repository
	tests := []struct {
		name string
		path string
	}{
		{"cmd directory", "cmd"},
		{"internal directory", "internal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := PathExistsOnRef(
				baseBranch,
				tt.path,
			)
			if err != nil {
				t.Fatalf(
					"PathExistsOnRef(%q, %q) error = %v",
					baseBranch,
					tt.path,
					err,
				)
			}
			if !exists {
				t.Errorf(
					"PathExistsOnRef(%q, %q) = false, want true",
					baseBranch,
					tt.path,
				)
			}
		})
	}
}

func TestPathExistsOnRef_NonExistingPath(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	// Change to the repo root for consistent path resolution
	repoRoot, err := GetRepoRoot()
	if err != nil {
		t.Fatalf("GetRepoRoot() error = %v", err)
	}
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf(
			"failed to change to repo root: %v",
			err,
		)
	}

	// Auto-detect the base branch
	baseBranch, err := GetBaseBranch("")
	if err != nil {
		t.Skipf(
			"could not auto-detect base branch: %v",
			err,
		)
	}

	// Test with a path that should not exist
	nonExistentPath := "nonexistent-path-xyz-12345"
	exists, err := PathExistsOnRef(
		baseBranch,
		nonExistentPath,
	)
	if err != nil {
		t.Fatalf(
			"PathExistsOnRef(%q, %q) error = %v",
			baseBranch,
			nonExistentPath,
			err,
		)
	}
	if exists {
		t.Errorf(
			"PathExistsOnRef(%q, %q) = true, want false",
			baseBranch,
			nonExistentPath,
		)
	}
}

func TestPathExistsOnRef_SubDirectoryPath(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	// Change to the repo root for consistent path resolution
	repoRoot, err := GetRepoRoot()
	if err != nil {
		t.Fatalf("GetRepoRoot() error = %v", err)
	}
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf(
			"failed to change to repo root: %v",
			err,
		)
	}

	// Auto-detect the base branch
	baseBranch, err := GetBaseBranch("")
	if err != nil {
		t.Skipf(
			"could not auto-detect base branch: %v",
			err,
		)
	}

	// Test with deeper paths that should exist in the repository
	tests := []struct {
		name string
		path string
	}{
		{
			"internal/git subdirectory",
			"internal/git",
		},
		{"cmd subdirectory", "cmd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := PathExistsOnRef(
				baseBranch,
				tt.path,
			)
			if err != nil {
				t.Fatalf(
					"PathExistsOnRef(%q, %q) error = %v",
					baseBranch,
					tt.path,
					err,
				)
			}
			if !exists {
				t.Errorf(
					"PathExistsOnRef(%q, %q) = false, want true",
					baseBranch,
					tt.path,
				)
			}
		})
	}
}
