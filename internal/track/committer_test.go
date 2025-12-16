package track

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// initTestGitRepo creates a new git repository in the given directory
// and returns a cleanup function.
func initTestGitRepo(t *testing.T, dir string) {
	t.Helper()

	// Initialize the git repo
	cmd := exec.Command("git", "init", dir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	// Configure user email
	cmd = exec.Command(
		"git", "-C", dir,
		"config", "user.email", "test@test.com",
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git email: %v", err)
	}

	// Configure user name
	cmd = exec.Command(
		"git", "-C", dir,
		"config", "user.name", "Test",
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git name: %v", err)
	}
}

// createTestFile creates a file with the given content in the directory.
func createTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)

	// Create parent directories if needed
	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
}

// makeInitialCommit creates an initial commit in the repo.
func makeInitialCommit(t *testing.T, dir string) {
	t.Helper()

	// Create a dummy file
	createTestFile(t, dir, "README.md", "# Test Repo\n")

	// Add and commit
	cmd := exec.Command("git", "-C", dir, "add", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to stage files: %v", err)
	}

	cmd = exec.Command(
		"git", "-C", dir,
		"commit", "-m", "Initial commit",
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}
}

func TestNewCommitter(t *testing.T) {
	t.Run("creates committer with correct fields", func(t *testing.T) {
		changeID := "test-change-123"
		repoRoot := "/path/to/repo"

		c := NewCommitter(changeID, repoRoot)

		if c == nil {
			t.Fatal("NewCommitter returned nil")
		}
		if c.changeID != changeID {
			t.Errorf(
				"NewCommitter().changeID = %q, want %q",
				c.changeID, changeID,
			)
		}
		if c.repoRoot != repoRoot {
			t.Errorf(
				"NewCommitter().repoRoot = %q, want %q",
				c.repoRoot, repoRoot,
			)
		}
	})

	t.Run("creates independent instances", func(t *testing.T) {
		c1 := NewCommitter("change-1", "/repo1")
		c2 := NewCommitter("change-2", "/repo2")

		if c1 == c2 {
			t.Error("NewCommitter returned same instance for different calls")
		}
		if c1.changeID == c2.changeID {
			t.Error("Different committers have same changeID")
		}
	})
}

func TestAction_String(t *testing.T) {
	tests := []struct {
		name   string
		action Action
		want   string
	}{
		{
			name:   "ActionStart returns start",
			action: ActionStart,
			want:   "start",
		},
		{
			name:   "ActionComplete returns complete",
			action: ActionComplete,
			want:   "complete",
		},
		{
			name:   "unknown action returns update",
			action: Action(99),
			want:   "update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.action.String()
			if got != tt.want {
				t.Errorf(
					"Action.String() = %q, want %q",
					got, tt.want,
				)
			}
		})
	}
}

func TestCommitter_buildCommitMessage(t *testing.T) {
	tests := []struct {
		name       string
		changeID   string
		taskID     string
		action     Action
		wantPrefix string
		wantFooter string
	}{
		{
			name:       "start action message format",
			changeID:   "add-feature",
			taskID:     "1.1",
			action:     ActionStart,
			wantPrefix: "spectr(add-feature): start task 1.1",
			wantFooter: "[Automated by spectr track]",
		},
		{
			name:       "complete action message format",
			changeID:   "fix-bug",
			taskID:     "2.3",
			action:     ActionComplete,
			wantPrefix: "spectr(fix-bug): complete task 2.3",
			wantFooter: "[Automated by spectr track]",
		},
		{
			name:       "change ID with special characters",
			changeID:   "add-track-command",
			taskID:     "4.2",
			action:     ActionComplete,
			wantPrefix: "spectr(add-track-command): complete task 4.2",
			wantFooter: "[Automated by spectr track]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommitter(tt.changeID, "/repo")
			got := c.buildCommitMessage(tt.taskID, tt.action)

			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf(
					"buildCommitMessage() prefix = %q, want prefix %q",
					got, tt.wantPrefix,
				)
			}
			if !strings.HasSuffix(got, tt.wantFooter) {
				t.Errorf(
					"buildCommitMessage() suffix = %q, want suffix %q",
					got, tt.wantFooter,
				)
			}
			// Verify double newline separates message and footer
			if !strings.Contains(got, "\n\n") {
				t.Error("buildCommitMessage() missing double newline separator")
			}
		})
	}
}

func TestCommitter_filterTaskFiles(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{
			name:  "filters out tasks.json",
			files: []string{"src/main.go", "tasks.json", "README.md"},
			want:  []string{"src/main.go", "README.md"},
		},
		{
			name:  "filters out tasks.jsonc",
			files: []string{"src/main.go", "tasks.jsonc"},
			want:  []string{"src/main.go"},
		},
		{
			name:  "filters out tasks.md",
			files: []string{"tasks.md", "src/main.go"},
			want:  []string{"src/main.go"},
		},
		{
			name:  "filters all task files",
			files: []string{"tasks.json", "tasks.jsonc", "tasks.md"},
			want:  nil,
		},
		{
			name: "filters task files in subdirectories",
			files: []string{
				"spectr/changes/test/tasks.json",
				"spectr/changes/test/tasks.jsonc",
				"spectr/changes/test/tasks.md",
				"src/main.go",
			},
			want: []string{"src/main.go"},
		},
		{
			name:  "preserves files with similar names",
			files: []string{"tasks.json.bak", "my-tasks.json", "tasks.jsonc.old"},
			want:  []string{"tasks.json.bak", "my-tasks.json", "tasks.jsonc.old"},
		},
		{
			name:  "handles empty input",
			files: make([]string, 0),
			want:  nil,
		},
		{
			name:  "handles nil input",
			files: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterTaskFiles(tt.files)

			if len(got) != len(tt.want) {
				t.Errorf(
					"filterTaskFiles() len = %d, want %d",
					len(got), len(tt.want),
				)
			}
			for i, f := range got {
				if i >= len(tt.want) || f != tt.want[i] {
					t.Errorf(
						"filterTaskFiles()[%d] = %q, want %q",
						i, f, tt.want[i],
					)
				}
			}
		})
	}
}

func TestCommitter_Commit_NoFilesToStage(t *testing.T) {
	tempDir := t.TempDir()
	initTestGitRepo(t, tempDir)
	makeInitialCommit(t, tempDir)

	c := NewCommitter("test-change", tempDir)

	// Create only task files (which should be filtered out)
	createTestFile(t, tempDir, "tasks.json", `{"tasks":[]}`)
	createTestFile(t, tempDir, "tasks.jsonc", `{"tasks":[]}`)
	createTestFile(t, tempDir, "tasks.md", "# Tasks\n")

	result, err := c.Commit("1.1", ActionStart)

	if err != nil {
		t.Fatalf("Commit() error = %v, want nil", err)
	}
	if !result.NoFiles {
		t.Error("Commit().NoFiles = false, want true")
	}
	if result.CommitHash != "" {
		t.Errorf(
			"Commit().CommitHash = %q, want empty",
			result.CommitHash,
		)
	}
	// Message should still be set even when no files
	if result.Message == "" {
		t.Error("Commit().Message should not be empty")
	}
}

func TestCommitter_Commit_WithFiles(t *testing.T) {
	tempDir := t.TempDir()
	initTestGitRepo(t, tempDir)
	makeInitialCommit(t, tempDir)

	c := NewCommitter("test-change", tempDir)

	// Create a non-task file
	createTestFile(t, tempDir, "src/main.go", "package main\n")

	result, err := c.Commit("1.1", ActionStart)

	if err != nil {
		t.Fatalf("Commit() error = %v, want nil", err)
	}
	if result.NoFiles {
		t.Error("Commit().NoFiles = true, want false")
	}
	if result.CommitHash == "" {
		t.Error("Commit().CommitHash should not be empty")
	}
	// Verify commit hash format (40 hex characters)
	if len(result.CommitHash) != 40 {
		t.Errorf(
			"Commit().CommitHash length = %d, want 40",
			len(result.CommitHash),
		)
	}

	// Verify message format
	expectedPrefix := "spectr(test-change): start task 1.1"
	if !strings.HasPrefix(result.Message, expectedPrefix) {
		t.Errorf(
			"Commit().Message prefix = %q, want %q",
			result.Message, expectedPrefix,
		)
	}

	// Verify the commit was actually created
	cmd := exec.Command("git", "-C", tempDir, "log", "-1", "--format=%s")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to get commit log: %v", err)
	}
	commitSubject := strings.TrimSpace(string(output))
	if commitSubject != expectedPrefix {
		t.Errorf(
			"git log subject = %q, want %q",
			commitSubject, expectedPrefix,
		)
	}
}

func TestCommitter_Commit_WithMixedFiles(t *testing.T) {
	tempDir := t.TempDir()
	initTestGitRepo(t, tempDir)
	makeInitialCommit(t, tempDir)

	c := NewCommitter("fix-bug", tempDir)

	// Create both task files and regular files
	createTestFile(t, tempDir, "tasks.jsonc", `{"tasks":[]}`)
	createTestFile(t, tempDir, "src/fix.go", "package fix\n")
	createTestFile(t, tempDir, "tests/fix_test.go", "package fix\n")

	result, err := c.Commit("2.1", ActionComplete)

	if err != nil {
		t.Fatalf("Commit() error = %v, want nil", err)
	}
	if result.NoFiles {
		t.Error("Commit().NoFiles = true, want false (should stage non-task files)")
	}
	if result.CommitHash == "" {
		t.Error("Commit().CommitHash should not be empty")
	}

	// Verify only non-task files were staged
	cmd := exec.Command(
		"git", "-C", tempDir,
		"show", "--name-only", "--format=", "HEAD",
	)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to get committed files: %v", err)
	}

	committedFiles := strings.TrimSpace(string(output))
	if strings.Contains(committedFiles, "tasks.jsonc") {
		t.Error("tasks.jsonc should not be in commit")
	}
	if !strings.Contains(committedFiles, "src/fix.go") {
		t.Error("src/fix.go should be in commit")
	}
	if !strings.Contains(committedFiles, "tests/fix_test.go") {
		t.Error("tests/fix_test.go should be in commit")
	}
}

func TestCommitter_Commit_GitError(t *testing.T) {
	// Use a non-existent directory to trigger git error
	c := NewCommitter("test-change", "/nonexistent/repo/path")

	result, err := c.Commit("1.1", ActionStart)

	if err == nil {
		t.Fatal("Commit() expected error, got nil")
	}

	// Verify it's wrapped as GitCommitError
	var gitErr *specterrs.GitCommitError
	if !isGitCommitError(err, &gitErr) {
		t.Errorf(
			"Commit() error type = %T, want *specterrs.GitCommitError",
			err,
		)
	}

	// Result should indicate no files
	if result.CommitHash != "" {
		t.Errorf(
			"Commit().CommitHash = %q, want empty on error",
			result.CommitHash,
		)
	}
}

// isGitCommitError checks if the error is a GitCommitError and sets target.
func isGitCommitError(err error, target **specterrs.GitCommitError) bool {
	if e, ok := err.(*specterrs.GitCommitError); ok {
		*target = e

		return true
	}

	return false
}

func TestCommitter_Commit_ActionTypes(t *testing.T) {
	tests := []struct {
		name         string
		action       Action
		wantContains string
	}{
		{
			name:         "ActionStart in message",
			action:       ActionStart,
			wantContains: "start task",
		},
		{
			name:         "ActionComplete in message",
			action:       ActionComplete,
			wantContains: "complete task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			initTestGitRepo(t, tempDir)
			makeInitialCommit(t, tempDir)

			c := NewCommitter("test-change", tempDir)
			createTestFile(t, tempDir, "file.go", "package main\n")

			result, err := c.Commit("1.1", tt.action)
			if err != nil {
				t.Fatalf("Commit() error = %v", err)
			}

			if !strings.Contains(result.Message, tt.wantContains) {
				t.Errorf(
					"Commit().Message = %q, want to contain %q",
					result.Message, tt.wantContains,
				)
			}
		})
	}
}

func TestCommitter_getModifiedFiles(t *testing.T) {
	tempDir := t.TempDir()
	initTestGitRepo(t, tempDir)
	makeInitialCommit(t, tempDir)

	c := NewCommitter("test-change", tempDir)

	t.Run("detects new untracked files", func(t *testing.T) {
		createTestFile(t, tempDir, "new_file.go", "package main\n")

		files, err := c.getModifiedFiles()
		if err != nil {
			t.Fatalf("getModifiedFiles() error = %v", err)
		}

		found := false
		for _, f := range files {
			if f == "new_file.go" {
				found = true

				break
			}
		}
		if !found {
			t.Error("getModifiedFiles() should include untracked new_file.go")
		}

		// Cleanup
		_ = os.Remove(filepath.Join(tempDir, "new_file.go"))
	})

	t.Run("detects modified tracked files", func(t *testing.T) {
		// Modify the README that was created in initial commit
		createTestFile(t, tempDir, "README.md", "# Modified\n")

		files, err := c.getModifiedFiles()
		if err != nil {
			t.Fatalf("getModifiedFiles() error = %v", err)
		}

		found := false
		for _, f := range files {
			if f == "README.md" {
				found = true

				break
			}
		}
		if !found {
			t.Error("getModifiedFiles() should include modified README.md")
		}

		// Reset the file
		cmd := exec.Command("git", "-C", tempDir, "checkout", "README.md")
		_ = cmd.Run()
	})
}

func TestCommitter_stageFiles(t *testing.T) {
	tempDir := t.TempDir()
	initTestGitRepo(t, tempDir)
	makeInitialCommit(t, tempDir)

	c := NewCommitter("test-change", tempDir)

	t.Run("stages single file", func(t *testing.T) {
		createTestFile(t, tempDir, "stage_test.go", "package main\n")

		err := c.stageFiles([]string{"stage_test.go"})
		if err != nil {
			t.Fatalf("stageFiles() error = %v", err)
		}

		// Verify file is staged
		cmd := exec.Command("git", "-C", tempDir, "diff", "--cached", "--name-only")
		output, _ := cmd.Output()
		if !strings.Contains(string(output), "stage_test.go") {
			t.Error("stageFiles() should stage stage_test.go")
		}

		// Cleanup
		cmd = exec.Command("git", "-C", tempDir, "reset", "HEAD")
		_ = cmd.Run()
		_ = os.Remove(filepath.Join(tempDir, "stage_test.go"))
	})

	t.Run("stages multiple files", func(t *testing.T) {
		createTestFile(t, tempDir, "file1.go", "package main\n")
		createTestFile(t, tempDir, "file2.go", "package main\n")

		err := c.stageFiles([]string{"file1.go", "file2.go"})
		if err != nil {
			t.Fatalf("stageFiles() error = %v", err)
		}

		// Verify files are staged
		cmd := exec.Command("git", "-C", tempDir, "diff", "--cached", "--name-only")
		output, _ := cmd.Output()
		staged := string(output)
		if !strings.Contains(staged, "file1.go") {
			t.Error("stageFiles() should stage file1.go")
		}
		if !strings.Contains(staged, "file2.go") {
			t.Error("stageFiles() should stage file2.go")
		}

		// Cleanup
		cmd = exec.Command("git", "-C", tempDir, "reset", "HEAD")
		_ = cmd.Run()
		_ = os.Remove(filepath.Join(tempDir, "file1.go"))
		_ = os.Remove(filepath.Join(tempDir, "file2.go"))
	})
}

func TestCommitResult_Struct(t *testing.T) {
	t.Run("NoFiles result", func(t *testing.T) {
		result := CommitResult{
			NoFiles:    true,
			CommitHash: "",
			Message:    "spectr(test): start task 1.1",
		}

		if !result.NoFiles {
			t.Error("CommitResult.NoFiles should be true")
		}
		if result.CommitHash != "" {
			t.Error("CommitResult.CommitHash should be empty for NoFiles")
		}
	})

	t.Run("WithFiles result", func(t *testing.T) {
		result := CommitResult{
			NoFiles:    false,
			CommitHash: "abc123def456789012345678901234567890abcd",
			Message:    "spectr(test): complete task 1.1\n\n[Automated by spectr track]",
		}

		if result.NoFiles {
			t.Error("CommitResult.NoFiles should be false")
		}
		if result.CommitHash == "" {
			t.Error("CommitResult.CommitHash should not be empty")
		}
		if result.Message == "" {
			t.Error("CommitResult.Message should not be empty")
		}
	})
}
