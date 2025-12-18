package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// MockCommandRunner is a test double for CommandRunner that returns
// predefined responses for git commands.
type MockCommandRunner struct {
	// Responses maps command strings to their outputs
	// Key format: "command arg1 arg2"
	Responses map[string]MockResponse
	// CallLog records all commands that were called
	CallLog []string
}

// MockResponse holds the output and error for a mocked command.
type MockResponse struct {
	Output []byte
	Err    error
}

// Run implements CommandRunner for testing.
func (m *MockCommandRunner) Run(
	dir string,
	name string,
	args ...string,
) ([]byte, error) {
	cmdKey := name + " " + strings.Join(args, " ")
	m.CallLog = append(m.CallLog, cmdKey)

	if resp, ok := m.Responses[cmdKey]; ok {
		return resp.Output, resp.Err
	}

	// Default: return empty output with no error
	return []byte{}, nil
}

func TestNewChangeDetector(t *testing.T) {
	detector := NewChangeDetector("/test/path")

	if detector == nil {
		t.Fatal("NewChangeDetector returned nil")
	}

	if detector.repoPath != "/test/path" {
		t.Errorf(
			"repoPath = %q, want %q",
			detector.repoPath,
			"/test/path",
		)
	}

	if detector.runner == nil {
		t.Error("runner is nil, want DefaultCommandRunner")
	}
}

func TestNewChangeDetectorWithRunner(t *testing.T) {
	mockRunner := &MockCommandRunner{}
	detector := NewChangeDetectorWithRunner("/test/path", mockRunner)

	if detector == nil {
		t.Fatal("NewChangeDetectorWithRunner returned nil")
	}

	if detector.repoPath != "/test/path" {
		t.Errorf(
			"repoPath = %q, want %q",
			detector.repoPath,
			"/test/path",
		)
	}

	if detector.runner != mockRunner {
		t.Error("runner is not the provided mock runner")
	}
}

func TestIsGitRepo_InRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Test with the current repository (this test runs in a git repo)
	// Get the repo root first
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		t.Skip("not running in a git repository")
	}
	repoRoot := strings.TrimSpace(string(output))

	result := IsGitRepo(repoRoot)
	if !result {
		t.Errorf("IsGitRepo(%q) = false, want true", repoRoot)
	}
}

func TestIsGitRepo_NotInRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Create a temp directory that is not a git repo
	tempDir, err := os.MkdirTemp("", "not-a-repo-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	result := IsGitRepo(tempDir)
	if result {
		t.Errorf("IsGitRepo(%q) = true, want false", tempDir)
	}
}

func TestSnapshot_Success(t *testing.T) {
	mockRunner := &MockCommandRunner{
		Responses: map[string]MockResponse{
			"git rev-parse HEAD": {
				Output: []byte("abc123def456\n"),
				Err:    nil,
			},
			"git status --porcelain": {
				Output: []byte(""),
				Err:    nil,
			},
		},
	}

	// We need to use a real git repo for IsGitRepo check
	// So let's skip this and test with integration tests
	t.Skip("Snapshot requires real git repo for IsGitRepo check")

	detector := NewChangeDetectorWithRunner("/test/repo", mockRunner)
	snapshot, err := detector.Snapshot()

	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}

	if !strings.HasPrefix(snapshot, "abc123def456:") {
		t.Errorf("Snapshot() = %q, want prefix 'abc123def456:'", snapshot)
	}
}

func TestChangedFiles_ParsesStatusOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected []string
	}{
		{
			name:     "empty output",
			output:   "",
			expected: nil,
		},
		{
			name:     "single untracked file",
			output:   "?? newfile.txt\n",
			expected: []string{"newfile.txt"},
		},
		{
			name:     "single modified file",
			output:   " M modified.txt\n",
			expected: []string{"modified.txt"},
		},
		{
			name:     "staged file",
			output:   "A  added.txt\n",
			expected: []string{"added.txt"},
		},
		{
			name: "multiple files",
			output: `?? untracked.txt
 M modified.txt
A  added.txt
`,
			expected: []string{"untracked.txt", "modified.txt", "added.txt"},
		},
		{
			name:     "renamed file",
			output:   "R  old.txt -> new.txt\n",
			expected: []string{"new.txt"},
		},
		{
			name:     "path with spaces",
			output:   "?? \"path with spaces/file.txt\"\n",
			expected: []string{"path with spaces/file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStatusOutput([]byte(tt.output))

			if len(result) != len(tt.expected) {
				t.Errorf(
					"parseStatusOutput() returned %d files, want %d",
					len(result),
					len(tt.expected),
				)
				t.Errorf("got: %v", result)
				t.Errorf("want: %v", tt.expected)
				return
			}

			for i, file := range result {
				if file != tt.expected[i] {
					t.Errorf(
						"parseStatusOutput()[%d] = %q, want %q",
						i,
						file,
						tt.expected[i],
					)
				}
			}
		})
	}
}

func TestCountStatusLines(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected int
	}{
		{
			name:     "empty output",
			output:   "",
			expected: 0,
		},
		{
			name:     "single line",
			output:   "?? file.txt\n",
			expected: 1,
		},
		{
			name: "multiple lines",
			output: `?? file1.txt
 M file2.txt
A  file3.txt
`,
			expected: 3,
		},
		{
			name:     "blank lines",
			output:   "\n\n",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countStatusLines([]byte(tt.output))
			if result != tt.expected {
				t.Errorf(
					"countStatusLines() = %d, want %d",
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestErrNotGitRepo_Message(t *testing.T) {
	expected := "spectr init requires a git repository. Run 'git init' first."
	if ErrNotGitRepo.Error() != expected {
		t.Errorf(
			"ErrNotGitRepo.Error() = %q, want %q",
			ErrNotGitRepo.Error(),
			expected,
		)
	}
}

// Integration tests that use a real git repository
func TestChangeDetector_Integration(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Create a temporary git repository
	tempDir, err := os.MkdirTemp("", "git-detector-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git name: %v", err)
	}

	// Create initial file and commit
	initialFile := filepath.Join(tempDir, "initial.txt")
	if err := os.WriteFile(initialFile, []byte("initial"), 0644); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git commit: %v", err)
	}

	t.Run("IsGitRepo returns true", func(t *testing.T) {
		if !IsGitRepo(tempDir) {
			t.Error("IsGitRepo() = false, want true")
		}
	})

	t.Run("Snapshot captures state", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)
		snapshot, err := detector.Snapshot()
		if err != nil {
			t.Fatalf("Snapshot() error = %v", err)
		}
		if snapshot == "" {
			t.Error("Snapshot() returned empty string")
		}
		// Snapshot format: "commit_hash:count"
		if !strings.Contains(snapshot, ":") {
			t.Errorf("Snapshot() = %q, want format 'hash:count'", snapshot)
		}
	})

	t.Run("IsClean returns true for clean tree", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)
		clean, err := detector.IsClean()
		if err != nil {
			t.Fatalf("IsClean() error = %v", err)
		}
		if !clean {
			t.Error("IsClean() = false, want true")
		}
	})

	t.Run("GetRepoRoot returns correct path", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)
		root, err := detector.GetRepoRoot()
		if err != nil {
			t.Fatalf("GetRepoRoot() error = %v", err)
		}
		// Compare canonical paths
		expectedRoot, _ := filepath.EvalSymlinks(tempDir)
		actualRoot, _ := filepath.EvalSymlinks(root)
		if actualRoot != expectedRoot {
			t.Errorf("GetRepoRoot() = %q, want %q", actualRoot, expectedRoot)
		}
	})

	t.Run("ChangedFiles detects new files", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)

		// Take snapshot before changes
		snapshot, err := detector.Snapshot()
		if err != nil {
			t.Fatalf("Snapshot() error = %v", err)
		}

		// Create a new file
		newFile := filepath.Join(tempDir, "newfile.txt")
		if err := os.WriteFile(newFile, []byte("new content"), 0644); err != nil {
			t.Fatalf("failed to create new file: %v", err)
		}
		defer os.Remove(newFile)

		// Get changed files
		files, err := detector.ChangedFiles(snapshot)
		if err != nil {
			t.Fatalf("ChangedFiles() error = %v", err)
		}

		found := false
		for _, f := range files {
			if f == "newfile.txt" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ChangedFiles() = %v, want to contain 'newfile.txt'", files)
		}
	})

	t.Run("UntrackedFiles returns untracked files", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)

		// Create an untracked file
		untrackedFile := filepath.Join(tempDir, "untracked.txt")
		if err := os.WriteFile(untrackedFile, []byte("untracked"), 0644); err != nil {
			t.Fatalf("failed to create untracked file: %v", err)
		}
		defer os.Remove(untrackedFile)

		files, err := detector.UntrackedFiles()
		if err != nil {
			t.Fatalf("UntrackedFiles() error = %v", err)
		}

		found := false
		for _, f := range files {
			if f == "untracked.txt" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("UntrackedFiles() = %v, want to contain 'untracked.txt'", files)
		}
	})

	t.Run("ModifiedFiles returns modified tracked files", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)

		// Modify the initial file
		if err := os.WriteFile(initialFile, []byte("modified"), 0644); err != nil {
			t.Fatalf("failed to modify file: %v", err)
		}
		defer os.WriteFile(initialFile, []byte("initial"), 0644) // restore

		files, err := detector.ModifiedFiles()
		if err != nil {
			t.Fatalf("ModifiedFiles() error = %v", err)
		}

		found := false
		for _, f := range files {
			if f == "initial.txt" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ModifiedFiles() = %v, want to contain 'initial.txt'", files)
		}
	})

	t.Run("IsClean returns false for dirty tree", func(t *testing.T) {
		detector := NewChangeDetector(tempDir)

		// Modify a file
		if err := os.WriteFile(initialFile, []byte("dirty"), 0644); err != nil {
			t.Fatalf("failed to modify file: %v", err)
		}
		defer os.WriteFile(initialFile, []byte("initial"), 0644)

		clean, err := detector.IsClean()
		if err != nil {
			t.Fatalf("IsClean() error = %v", err)
		}
		if clean {
			t.Error("IsClean() = true, want false")
		}
	})
}

func TestChangeDetector_NotGitRepo(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Create a temp directory that is not a git repo
	tempDir, err := os.MkdirTemp("", "not-a-repo-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	detector := NewChangeDetector(tempDir)

	t.Run("Snapshot returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.Snapshot()
		if err != ErrNotGitRepo {
			t.Errorf("Snapshot() error = %v, want ErrNotGitRepo", err)
		}
	})

	t.Run("ChangedFiles returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.ChangedFiles("abc:0")
		if err != ErrNotGitRepo {
			t.Errorf("ChangedFiles() error = %v, want ErrNotGitRepo", err)
		}
	})

	t.Run("GetRepoRoot returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.GetRepoRoot()
		if err != ErrNotGitRepo {
			t.Errorf("GetRepoRoot() error = %v, want ErrNotGitRepo", err)
		}
	})

	t.Run("IsClean returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.IsClean()
		if err != ErrNotGitRepo {
			t.Errorf("IsClean() error = %v, want ErrNotGitRepo", err)
		}
	})

	t.Run("UntrackedFiles returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.UntrackedFiles()
		if err != ErrNotGitRepo {
			t.Errorf("UntrackedFiles() error = %v, want ErrNotGitRepo", err)
		}
	})

	t.Run("ModifiedFiles returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.ModifiedFiles()
		if err != ErrNotGitRepo {
			t.Errorf("ModifiedFiles() error = %v, want ErrNotGitRepo", err)
		}
	})

	t.Run("DiffFiles returns ErrNotGitRepo", func(t *testing.T) {
		_, err := detector.DiffFiles("HEAD")
		if err != ErrNotGitRepo {
			t.Errorf("DiffFiles() error = %v, want ErrNotGitRepo", err)
		}
	})
}

func TestChangedFiles_InvalidSnapshot(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Get the repo root
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		t.Skip("not running in a git repository")
	}
	repoRoot := strings.TrimSpace(string(output))

	detector := NewChangeDetector(repoRoot)

	_, err = detector.ChangedFiles("invalid-snapshot-format")
	if err == nil {
		t.Error("ChangedFiles() with invalid snapshot expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid snapshot format") {
		t.Errorf(
			"ChangedFiles() error = %v, want error containing 'invalid snapshot format'",
			err,
		)
	}
}

func TestAbsolutePath(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Get the repo root
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		t.Skip("not running in a git repository")
	}
	repoRoot := strings.TrimSpace(string(output))

	detector := NewChangeDetector(repoRoot)

	absPath, err := detector.AbsolutePath("some/relative/path.txt")
	if err != nil {
		t.Fatalf("AbsolutePath() error = %v", err)
	}

	expected := filepath.Join(repoRoot, "some/relative/path.txt")
	if absPath != expected {
		t.Errorf("AbsolutePath() = %q, want %q", absPath, expected)
	}
}

func TestDiffFiles_Integration(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	// Create a temporary git repository
	tempDir, err := os.MkdirTemp("", "git-diff-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	cmd.Run()

	// Create initial file and commit
	initialFile := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(initialFile, []byte("v1"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tempDir
	cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "v1")
	cmd.Dir = tempDir
	cmd.Run()

	// Get the first commit hash
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = tempDir
	output, _ := cmd.Output()
	firstCommit := strings.TrimSpace(string(output))

	// Modify file
	if err := os.WriteFile(initialFile, []byte("v2"), 0644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	detector := NewChangeDetector(tempDir)
	files, err := detector.DiffFiles(firstCommit)
	if err != nil {
		t.Fatalf("DiffFiles() error = %v", err)
	}

	found := false
	for _, f := range files {
		if f == "file.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("DiffFiles() = %v, want to contain 'file.txt'", files)
	}
}

// Test helper function for verbose test output
func TestMockCommandRunner(t *testing.T) {
	mock := &MockCommandRunner{
		Responses: map[string]MockResponse{
			"git status --porcelain": {
				Output: []byte("?? test.txt\n"),
				Err:    nil,
			},
			"git rev-parse HEAD": {
				Output: []byte("abc123\n"),
				Err:    nil,
			},
		},
	}

	// Test running a known command
	output, err := mock.Run("/test", "git", "status", "--porcelain")
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if string(output) != "?? test.txt\n" {
		t.Errorf("Run() output = %q, want %q", string(output), "?? test.txt\n")
	}

	// Test running another command
	output, err = mock.Run("/test", "git", "rev-parse", "HEAD")
	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if string(output) != "abc123\n" {
		t.Errorf("Run() output = %q, want %q", string(output), "abc123\n")
	}

	// Verify call log
	if len(mock.CallLog) != 2 {
		t.Errorf("CallLog length = %d, want 2", len(mock.CallLog))
	}
	if mock.CallLog[0] != "git status --porcelain" {
		t.Errorf("CallLog[0] = %q, want 'git status --porcelain'", mock.CallLog[0])
	}
	if mock.CallLog[1] != "git rev-parse HEAD" {
		t.Errorf("CallLog[1] = %q, want 'git rev-parse HEAD'", mock.CallLog[1])
	}
}

func TestMockCommandRunner_Error(t *testing.T) {
	expectedErr := fmt.Errorf("command failed")
	mock := &MockCommandRunner{
		Responses: map[string]MockResponse{
			"git status": {
				Output: []byte("error output"),
				Err:    expectedErr,
			},
		},
	}

	output, err := mock.Run("/test", "git", "status")
	if err != expectedErr {
		t.Errorf("Run() error = %v, want %v", err, expectedErr)
	}
	if string(output) != "error output" {
		t.Errorf("Run() output = %q, want 'error output'", string(output))
	}
}

func TestMockCommandRunner_UnknownCommand(t *testing.T) {
	mock := &MockCommandRunner{
		Responses: map[string]MockResponse{},
	}

	// Unknown commands return empty output with no error
	output, err := mock.Run("/test", "git", "unknown", "command")
	if err != nil {
		t.Errorf("Run() unexpected error for unknown command: %v", err)
	}
	if len(output) != 0 {
		t.Errorf("Run() output = %q, want empty", string(output))
	}
}
