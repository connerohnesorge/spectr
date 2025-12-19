package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// MockGitExecutor is a mock implementation of GitExecutor for testing.
type MockGitExecutor struct {
	RevParseFunc        func(repoPath, ref string) (string, error)
	StatusPorcelainFunc func(repoPath string) (string, error)
	DiffNameOnlyFunc    func(repoPath, fromRef, toRef string) (string, error)
	StashCreateFunc     func(repoPath string) (string, error)
	IsGitRepoFunc       func(path string) bool
}

func (m *MockGitExecutor) RevParse(repoPath, ref string) (string, error) {
	if m.RevParseFunc != nil {
		return m.RevParseFunc(repoPath, ref)
	}

	return "abc123", nil
}

func (m *MockGitExecutor) StatusPorcelain(repoPath string) (string, error) {
	if m.StatusPorcelainFunc != nil {
		return m.StatusPorcelainFunc(repoPath)
	}

	return "", nil
}

func (m *MockGitExecutor) DiffNameOnly(repoPath, fromRef, toRef string) (string, error) {
	if m.DiffNameOnlyFunc != nil {
		return m.DiffNameOnlyFunc(repoPath, fromRef, toRef)
	}

	return "", nil
}

func (m *MockGitExecutor) StashCreate(repoPath string) (string, error) {
	if m.StashCreateFunc != nil {
		return m.StashCreateFunc(repoPath)
	}

	return "", nil
}

func (m *MockGitExecutor) IsGitRepo(path string) bool {
	if m.IsGitRepoFunc != nil {
		return m.IsGitRepoFunc(path)
	}

	return true
}

// TestIsGitRepo_RealDirectory tests IsGitRepo with actual directories.
func TestIsGitRepo_RealDirectory(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()

	// Test 1: Non-git directory should return false
	nonGitDir := filepath.Join(tmpDir, "not-a-repo")
	if err := os.MkdirAll(nonGitDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}
	if IsGitRepo(nonGitDir) {
		t.Error("IsGitRepo should return false for non-git directory")
	}

	// Test 2: Directory with .git should return true
	gitDir := filepath.Join(tmpDir, "git-repo")
	if err := os.MkdirAll(filepath.Join(gitDir, ".git"), 0755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}
	if !IsGitRepo(gitDir) {
		t.Error("IsGitRepo should return true for directory with .git")
	}

	// Test 3: Actual git init should work
	realGitDir := filepath.Join(tmpDir, "real-git-repo")
	if err := os.MkdirAll(realGitDir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}
	cmd := exec.Command("git", "init", "-q", realGitDir)
	if err := cmd.Run(); err != nil {
		t.Skipf("git init failed, skipping: %v", err)
	}
	if !IsGitRepo(realGitDir) {
		t.Error("IsGitRepo should return true for initialized git repo")
	}
}

// TestNewChangeDetector tests the constructor.
func TestNewChangeDetector(t *testing.T) {
	detector := NewChangeDetector("/test/path")
	if detector == nil {
		t.Fatal("NewChangeDetector returned nil")
	}
	if detector.RepoPath() != "/test/path" {
		t.Errorf("expected repoPath '/test/path', got %q", detector.RepoPath())
	}
}

// TestSnapshot_Success tests successful snapshot creation.
func TestSnapshot_Success(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		RevParseFunc: func(_, _ string) (string, error) {
			return "abc123def456", nil
		},
		StashCreateFunc: func(_ string) (string, error) {
			return "stash@{0}", nil
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return "?? untracked.txt\n M modified.txt\n", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	snapshot, err := detector.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	// Verify snapshot contains expected components
	if !strings.Contains(snapshot, "HEAD:abc123def456") {
		t.Errorf("snapshot should contain HEAD commit, got: %s", snapshot)
	}
	if !strings.Contains(snapshot, "STASH:stash@{0}") {
		t.Errorf("snapshot should contain stash ref, got: %s", snapshot)
	}
	if !strings.Contains(snapshot, "UNTRACKED:untracked.txt") {
		t.Errorf("snapshot should contain untracked files, got: %s", snapshot)
	}
}

// TestSnapshot_NotGitRepo tests snapshot on non-git directory.
func TestSnapshot_NotGitRepo(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return false
		},
	}

	detector := NewChangeDetectorWithExecutor("/not/a/repo", mock)
	_, err := detector.Snapshot()
	if err == nil {
		t.Fatal("expected error for non-git repo")
	}
	if !strings.Contains(err.Error(), ErrNotGitRepo) {
		t.Errorf("expected ErrNotGitRepo, got: %v", err)
	}
}

// TestSnapshot_EmptyRepo tests snapshot on empty repo (no commits).
func TestSnapshot_EmptyRepo(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		RevParseFunc: func(_, _ string) (string, error) {
			return "", &exec.ExitError{}
		},
		StashCreateFunc: func(_ string) (string, error) {
			return "", nil
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return "", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	snapshot, err := detector.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot should succeed on empty repo: %v", err)
	}
	if !strings.Contains(snapshot, "HEAD:") {
		t.Errorf("snapshot should have HEAD component (empty), got: %s", snapshot)
	}
}

// TestChangedFiles_ModifiedFiles tests detection of modified tracked files.
func TestChangedFiles_ModifiedFiles(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return " M file1.txt\n M file2.txt\n", nil
		},
		DiffNameOnlyFunc: func(_, _, _ string) (string, error) {
			return "file1.txt\nfile2.txt\n", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	beforeSnapshot := "HEAD:abc123|STASH:|UNTRACKED:"

	files, err := detector.ChangedFiles(beforeSnapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	sort.Strings(files)
	expected := []string{"file1.txt", "file2.txt"}
	sort.Strings(expected)

	if len(files) != len(expected) {
		t.Errorf("expected %d files, got %d: %v", len(expected), len(files), files)
	}
	for i, f := range expected {
		if i < len(files) && files[i] != f {
			t.Errorf("expected file %q, got %q", f, files[i])
		}
	}
}

// TestChangedFiles_UntrackedFiles tests detection of new untracked files.
func TestChangedFiles_UntrackedFiles(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return "?? new_file.txt\n?? another_new.txt\n", nil
		},
		DiffNameOnlyFunc: func(_, _, _ string) (string, error) {
			return "", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	// Before snapshot had no untracked files
	beforeSnapshot := "HEAD:abc123|STASH:|UNTRACKED:"

	files, err := detector.ChangedFiles(beforeSnapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	sort.Strings(files)
	expected := []string{"another_new.txt", "new_file.txt"}
	sort.Strings(expected)

	if len(files) != len(expected) {
		t.Errorf("expected %d files, got %d: %v", len(expected), len(files), files)
	}
}

// TestChangedFiles_ExcludesPreExistingUntracked tests that pre-existing
// untracked files are not reported as changed.
func TestChangedFiles_ExcludesPreExistingUntracked(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return "?? existing.txt\n?? new_file.txt\n", nil
		},
		DiffNameOnlyFunc: func(_, _, _ string) (string, error) {
			return "", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	// Before snapshot already had existing.txt as untracked
	beforeSnapshot := "HEAD:abc123|STASH:|UNTRACKED:existing.txt"

	files, err := detector.ChangedFiles(beforeSnapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	// Should only contain new_file.txt, not existing.txt
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d: %v", len(files), files)
	}
	if len(files) > 0 && files[0] != "new_file.txt" {
		t.Errorf("expected new_file.txt, got %s", files[0])
	}
}

// TestChangedFiles_MixedChanges tests detection with modified and new files.
func TestChangedFiles_MixedChanges(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return " M modified.go\nA  added.go\n?? untracked.go\n", nil
		},
		DiffNameOnlyFunc: func(_, _, _ string) (string, error) {
			return "modified.go\n", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	beforeSnapshot := "HEAD:abc123|STASH:|UNTRACKED:"

	files, err := detector.ChangedFiles(beforeSnapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	sort.Strings(files)
	expected := []string{"added.go", "modified.go", "untracked.go"}
	sort.Strings(expected)

	if len(files) != len(expected) {
		t.Errorf("expected %d files, got %d: %v", len(expected), len(files), files)
	}
}

// TestChangedFiles_ExcludesDeleted tests that deleted files are excluded.
func TestChangedFiles_ExcludesDeleted(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return " M modified.go\n D deleted.go\nD  deleted2.go\n", nil
		},
		DiffNameOnlyFunc: func(_, _, _ string) (string, error) {
			return "modified.go\ndeleted.go\ndeleted2.go\n", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	beforeSnapshot := "HEAD:abc123|STASH:|UNTRACKED:"

	files, err := detector.ChangedFiles(beforeSnapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	// Should only contain modified.go, not deleted files
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d: %v", len(files), files)
	}
	if len(files) > 0 && files[0] != "modified.go" {
		t.Errorf("expected modified.go, got %s", files[0])
	}
}

// TestChangedFiles_NotGitRepo tests error handling for non-git repos.
func TestChangedFiles_NotGitRepo(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return false
		},
	}

	detector := NewChangeDetectorWithExecutor("/not/a/repo", mock)
	_, err := detector.ChangedFiles("HEAD:|STASH:|UNTRACKED:")
	if err == nil {
		t.Fatal("expected error for non-git repo")
	}
	if !strings.Contains(err.Error(), ErrNotGitRepo) {
		t.Errorf("expected ErrNotGitRepo, got: %v", err)
	}
}

// TestChangedFiles_RenamedFiles tests handling of renamed files.
func TestChangedFiles_RenamedFiles(t *testing.T) {
	mock := &MockGitExecutor{
		IsGitRepoFunc: func(_ string) bool {
			return true
		},
		StatusPorcelainFunc: func(_ string) (string, error) {
			return "R  old.txt -> new.txt\n", nil
		},
		DiffNameOnlyFunc: func(_, _, _ string) (string, error) {
			return "", nil
		},
	}

	detector := NewChangeDetectorWithExecutor("/test/repo", mock)
	beforeSnapshot := "HEAD:abc123|STASH:|UNTRACKED:"

	files, err := detector.ChangedFiles(beforeSnapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	// Should contain the new name
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d: %v", len(files), files)
	}
	if len(files) > 0 && files[0] != "new.txt" {
		t.Errorf("expected new.txt, got %s", files[0])
	}
}

// TestParseSnapshot tests the parseSnapshot function.
func TestParseSnapshot(t *testing.T) {
	tests := []struct {
		name              string
		snapshot          string
		expectedHead      string
		expectedStash     string
		expectedUntracked []string
	}{
		{
			name:              "full snapshot",
			snapshot:          "HEAD:abc123|STASH:stash@{0}|UNTRACKED:file1.txt,file2.txt",
			expectedHead:      "abc123",
			expectedStash:     "stash@{0}",
			expectedUntracked: []string{"file1.txt", "file2.txt"},
		},
		{
			name:              "empty snapshot",
			snapshot:          "HEAD:|STASH:|UNTRACKED:",
			expectedHead:      "",
			expectedStash:     "",
			expectedUntracked: nil,
		},
		{
			name:              "head only",
			snapshot:          "HEAD:abc123|STASH:|UNTRACKED:",
			expectedHead:      "abc123",
			expectedStash:     "",
			expectedUntracked: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, stash, untracked := parseSnapshot(tt.snapshot)
			if head != tt.expectedHead {
				t.Errorf("head: expected %q, got %q", tt.expectedHead, head)
			}
			if stash != tt.expectedStash {
				t.Errorf("stash: expected %q, got %q", tt.expectedStash, stash)
			}
			if len(untracked) != len(tt.expectedUntracked) {
				t.Errorf("untracked: expected %v, got %v", tt.expectedUntracked, untracked)
			}
		})
	}
}

// TestParseUntrackedFiles tests the parseUntrackedFiles function.
func TestParseUntrackedFiles(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected []string
	}{
		{
			name:     "single untracked",
			status:   "?? untracked.txt\n",
			expected: []string{"untracked.txt"},
		},
		{
			name:     "multiple untracked",
			status:   "?? file1.txt\n?? file2.txt\n",
			expected: []string{"file1.txt", "file2.txt"},
		},
		{
			name:     "mixed status",
			status:   " M modified.txt\n?? untracked.txt\nA  added.txt\n",
			expected: []string{"untracked.txt"},
		},
		{
			name:     "empty",
			status:   "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseUntrackedFiles(tt.status)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestParseChangedFilesFromStatus tests the parseChangedFilesFromStatus function.
func TestParseChangedFilesFromStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected []string
	}{
		{
			name:     "modified files",
			status:   " M file1.txt\nM  file2.txt\nMM file3.txt\n",
			expected: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name:     "added files",
			status:   "A  added.txt\n",
			expected: []string{"added.txt"},
		},
		{
			name:     "untracked files",
			status:   "?? untracked.txt\n",
			expected: []string{"untracked.txt"},
		},
		{
			name:     "deleted files excluded",
			status:   " D deleted.txt\nD  deleted2.txt\n",
			expected: nil,
		},
		{
			name:     "renamed file",
			status:   "R  old.txt -> new.txt\n",
			expected: []string{"new.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseChangedFilesFromStatus(tt.status)
			sort.Strings(result)
			sort.Strings(tt.expected)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)

				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, result)

					break
				}
			}
		})
	}
}

// TestIntegration_RealGitRepo tests the full flow with a real git repository.
func TestIntegration_RealGitRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping integration test")
	}

	// Create a temp directory and init git repo
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init", "-q", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config email failed: %v", err)
	}
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config name failed: %v", err)
	}

	// Create initial file and commit
	initialFile := filepath.Join(tmpDir, "initial.txt")
	if err := os.WriteFile(initialFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}
	cmd = exec.Command("git", "-C", tmpDir, "add", "initial.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("git add failed: %v", err)
	}
	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}

	// Create detector and take snapshot
	detector := NewChangeDetector(tmpDir)
	snapshot, err := detector.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	// Make changes: modify existing file and create new file
	if err := os.WriteFile(initialFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}
	newFile := filepath.Join(tmpDir, "new_file.txt")
	if err := os.WriteFile(newFile, []byte("new content"), 0644); err != nil {
		t.Fatalf("failed to create new file: %v", err)
	}

	// Detect changes
	files, err := detector.ChangedFiles(snapshot)
	if err != nil {
		t.Fatalf("ChangedFiles failed: %v", err)
	}

	// Should detect both modified and new file
	sort.Strings(files)
	if len(files) != 2 {
		t.Errorf("expected 2 changed files, got %d: %v", len(files), files)
	}

	foundInitial := false
	foundNew := false
	for _, f := range files {
		if f == "initial.txt" {
			foundInitial = true
		}
		if f == "new_file.txt" {
			foundNew = true
		}
	}

	if !foundInitial {
		t.Error("should have detected initial.txt as modified")
	}
	if !foundNew {
		t.Error("should have detected new_file.txt as new")
	}
}

// TestRealGitExecutor tests the RealGitExecutor implementation.
func TestRealGitExecutor(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping RealGitExecutor test")
	}

	// Create a temp directory
	tmpDir := t.TempDir()

	executor := &RealGitExecutor{}

	// Test IsGitRepo on non-git directory
	if executor.IsGitRepo(tmpDir) {
		t.Error("IsGitRepo should return false for non-git directory")
	}

	// Initialize git repo
	cmd := exec.Command("git", "init", "-q", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// Test IsGitRepo on git directory
	if !executor.IsGitRepo(tmpDir) {
		t.Error("IsGitRepo should return true for git directory")
	}

	// Test StatusPorcelain
	status, err := executor.StatusPorcelain(tmpDir)
	if err != nil {
		t.Errorf("StatusPorcelain failed: %v", err)
	}
	// Empty repo should have no status output
	if status != "" {
		t.Logf("StatusPorcelain returned: %q", status)
	}
}
