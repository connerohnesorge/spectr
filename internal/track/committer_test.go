package track

import (
	"errors"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// MockGitExecutor implements GitExecutor for testing.
type MockGitExecutor struct {
	// StatusOutput is the output returned by Status.
	StatusOutput string
	// StatusError is the error returned by Status.
	StatusError error
	// AddError is the error returned by Add.
	AddError error
	// CommitError is the error returned by Commit.
	CommitError error
	// RevParseOutput is the output returned by RevParse.
	RevParseOutput string
	// RevParseError is the error returned by RevParse.
	RevParseError error
	// DiffNumstatOutput is the output returned by DiffNumstat.
	DiffNumstatOutput string
	// DiffNumstatError is the error returned by DiffNumstat.
	DiffNumstatError error

	// AddedFiles records the files passed to Add.
	AddedFiles []string
	// CommitMessages records the messages passed to Commit.
	CommitMessages []string
	// StatusCalls counts the number of times Status was called.
	StatusCalls int
	// AddCalls counts the number of times Add was called.
	AddCalls int
	// CommitCalls counts the number of times Commit was called.
	CommitCalls int
	// RevParseCalls counts the number of times RevParse was called.
	RevParseCalls int
	// DiffNumstatCalls counts the number of times DiffNumstat was called.
	DiffNumstatCalls int
}

// Status implements GitExecutor.Status.
func (m *MockGitExecutor) Status(_ string) (string, error) {
	m.StatusCalls++

	return m.StatusOutput, m.StatusError
}

// Add implements GitExecutor.Add.
func (m *MockGitExecutor) Add(_ string, files []string) error {
	m.AddCalls++
	m.AddedFiles = append(m.AddedFiles, files...)

	return m.AddError
}

// Commit implements GitExecutor.Commit.
func (m *MockGitExecutor) Commit(_, message string) error {
	m.CommitCalls++
	m.CommitMessages = append(m.CommitMessages, message)

	return m.CommitError
}

// RevParse implements GitExecutor.RevParse.
func (m *MockGitExecutor) RevParse(_, _ string) (string, error) {
	m.RevParseCalls++

	return m.RevParseOutput, m.RevParseError
}

// DiffNumstat implements GitExecutor.DiffNumstat.
func (m *MockGitExecutor) DiffNumstat(_ string, _ []string) (string, error) {
	m.DiffNumstatCalls++

	return m.DiffNumstatOutput, m.DiffNumstatError
}

func TestNewCommitter(t *testing.T) {
	t.Run("creates committer with correct fields", func(t *testing.T) {
		changeID := "test-change-123"
		repoRoot := "/path/to/repo"

		c := NewCommitter(changeID, repoRoot, false)

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
		if c.gitExecutor == nil {
			t.Error("NewCommitter().gitExecutor should not be nil")
		}
		if c.includeBinaries {
			t.Error("NewCommitter().includeBinaries should be false")
		}
	})

	t.Run("creates independent instances", func(t *testing.T) {
		c1 := NewCommitter("change-1", "/repo1", false)
		c2 := NewCommitter("change-2", "/repo2", true)

		if c1 == c2 {
			t.Error("NewCommitter returned same instance for different calls")
		}
		if c1.changeID == c2.changeID {
			t.Error("Different committers have same changeID")
		}
	})

	t.Run("includes binaries when flag is true", func(t *testing.T) {
		c := NewCommitter("test-change", "/repo", true)

		if !c.includeBinaries {
			t.Error("NewCommitter().includeBinaries should be true")
		}
	})
}

func TestNewCommitterWithExecutor(t *testing.T) {
	t.Run("uses provided executor", func(t *testing.T) {
		mock := &MockGitExecutor{}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		if c.gitExecutor == nil {
			t.Error("NewCommitterWithExecutor did not set executor")
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
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(tt.changeID, "/repo", false, mock)
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
	// Mock returns only task files (which should be filtered out)
	mock := &MockGitExecutor{
		StatusOutput: "?? tasks.json\n?? tasks.jsonc\n?? tasks.md\n",
	}
	c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

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
	// Verify Status was called but Add and Commit were not
	if mock.StatusCalls != 1 {
		t.Errorf("Status should be called once, got %d", mock.StatusCalls)
	}
	if mock.AddCalls != 0 {
		t.Errorf("Add should not be called when no files, got %d", mock.AddCalls)
	}
	if mock.CommitCalls != 0 {
		t.Errorf(
			"Commit should not be called when no files, got %d",
			mock.CommitCalls,
		)
	}
}

func TestCommitter_Commit_WithFiles(t *testing.T) {
	// Mock returns a non-task file
	mock := &MockGitExecutor{
		StatusOutput:   "?? src/main.go\n",
		RevParseOutput: "abc123def456789012345678901234567890abcd",
	}
	c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

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
	// Verify commit hash
	if result.CommitHash != "abc123def456789012345678901234567890abcd" {
		t.Errorf(
			"Commit().CommitHash = %q, want %q",
			result.CommitHash, "abc123def456789012345678901234567890abcd",
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

	// Verify the mock was called correctly
	if mock.StatusCalls != 1 {
		t.Errorf("Status should be called once, got %d", mock.StatusCalls)
	}
	if mock.AddCalls != 1 {
		t.Errorf("Add should be called once, got %d", mock.AddCalls)
	}
	if mock.CommitCalls != 1 {
		t.Errorf("Commit should be called once, got %d", mock.CommitCalls)
	}
	if mock.RevParseCalls != 1 {
		t.Errorf("RevParse should be called once, got %d", mock.RevParseCalls)
	}

	// Verify the correct file was staged
	if len(mock.AddedFiles) != 1 || mock.AddedFiles[0] != "src/main.go" {
		t.Errorf("AddedFiles = %v, want [src/main.go]", mock.AddedFiles)
	}
}

func TestCommitter_Commit_WithMixedFiles(t *testing.T) {
	// Mock returns both task files and regular files
	mock := &MockGitExecutor{
		StatusOutput:   "?? tasks.jsonc\n?? src/fix.go\n?? tests/fix_test.go\n",
		RevParseOutput: "def456789012345678901234567890abcdef1234",
	}
	c := NewCommitterWithExecutor("fix-bug", "/repo", false, mock)

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
	expectedFiles := []string{"src/fix.go", "tests/fix_test.go"}
	if len(mock.AddedFiles) != len(expectedFiles) {
		t.Errorf(
			"AddedFiles count = %d, want %d",
			len(mock.AddedFiles), len(expectedFiles),
		)
	}
	for _, expected := range expectedFiles {
		found := false
		for _, added := range mock.AddedFiles {
			if added == expected {
				found = true

				break
			}
		}
		if !found {
			t.Errorf("Expected file %q to be staged", expected)
		}
	}
	// Verify task file was NOT staged
	for _, added := range mock.AddedFiles {
		if added == "tasks.jsonc" {
			t.Error("tasks.jsonc should not be staged")
		}
	}
}

func TestCommitter_Commit_WithModifiedFiles(t *testing.T) {
	// Mock returns modified files (unstaged and staged)
	mock := &MockGitExecutor{
		StatusOutput:   " M src/modified.go\nM  src/staged.go\n",
		RevParseOutput: "abc123def456789012345678901234567890abcd",
	}
	c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

	result, err := c.Commit("1.1", ActionStart)

	if err != nil {
		t.Fatalf("Commit() error = %v, want nil", err)
	}
	if result.NoFiles {
		t.Error("Commit().NoFiles = true, want false (should stage modified files)")
	}
	if result.CommitHash == "" {
		t.Error("Commit().CommitHash should not be empty")
	}

	// Verify both modified files were staged
	expectedFiles := []string{"src/modified.go", "src/staged.go"}
	if len(mock.AddedFiles) != len(expectedFiles) {
		t.Errorf(
			"AddedFiles count = %d, want %d",
			len(mock.AddedFiles), len(expectedFiles),
		)
	}
	for _, expected := range expectedFiles {
		found := false
		for _, added := range mock.AddedFiles {
			if added == expected {
				found = true

				break
			}
		}
		if !found {
			t.Errorf("Expected modified file %q to be staged", expected)
		}
	}
}

func TestCommitter_Commit_WithMixedNewAndModified(t *testing.T) {
	// Mock returns both new (untracked) and modified files
	mock := &MockGitExecutor{
		StatusOutput:   "?? src/new.go\n M src/modified.go\nM  src/staged.go\n",
		RevParseOutput: "def456789012345678901234567890abcdef1234",
	}
	c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

	result, err := c.Commit("1.1", ActionComplete)

	if err != nil {
		t.Fatalf("Commit() error = %v, want nil", err)
	}
	if result.NoFiles {
		t.Error("Commit().NoFiles = true, want false")
	}
	if result.CommitHash == "" {
		t.Error("Commit().CommitHash should not be empty")
	}

	// Verify all three files were staged (new + modified)
	expectedFiles := []string{"src/new.go", "src/modified.go", "src/staged.go"}
	if len(mock.AddedFiles) != len(expectedFiles) {
		t.Errorf(
			"AddedFiles count = %d, want %d",
			len(mock.AddedFiles), len(expectedFiles),
		)
	}
	for _, expected := range expectedFiles {
		found := false
		for _, added := range mock.AddedFiles {
			if added == expected {
				found = true

				break
			}
		}
		if !found {
			t.Errorf("Expected file %q to be staged", expected)
		}
	}
}

func TestCommitter_Commit_GitError(t *testing.T) {
	t.Run("status error", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusError: errors.New("git status failed: not a git repository"),
		}
		c := NewCommitterWithExecutor("test-change", "/nonexistent", false, mock)

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
	})

	t.Run("add error", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput: "?? src/main.go\n",
			AddError:     errors.New("git add failed: permission denied"),
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		result, err := c.Commit("1.1", ActionStart)

		if err == nil {
			t.Fatal("Commit() expected error, got nil")
		}

		var gitErr *specterrs.GitCommitError
		if !isGitCommitError(err, &gitErr) {
			t.Errorf(
				"Commit() error type = %T, want *specterrs.GitCommitError",
				err,
			)
		}

		if result.CommitHash != "" {
			t.Errorf(
				"Commit().CommitHash = %q, want empty on error",
				result.CommitHash,
			)
		}
	})

	t.Run("commit error", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput: "?? src/main.go\n",
			CommitError:  errors.New("git commit failed: nothing to commit"),
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		result, err := c.Commit("1.1", ActionStart)

		if err == nil {
			t.Fatal("Commit() expected error, got nil")
		}

		var gitErr *specterrs.GitCommitError
		if !isGitCommitError(err, &gitErr) {
			t.Errorf(
				"Commit() error type = %T, want *specterrs.GitCommitError",
				err,
			)
		}

		if result.CommitHash != "" {
			t.Errorf(
				"Commit().CommitHash = %q, want empty on error",
				result.CommitHash,
			)
		}
	})

	t.Run("rev-parse error", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput:  "?? src/main.go\n",
			RevParseError: errors.New("git rev-parse failed: ambiguous argument"),
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		result, err := c.Commit("1.1", ActionStart)

		if err == nil {
			t.Fatal("Commit() expected error, got nil")
		}

		var gitErr *specterrs.GitCommitError
		if !isGitCommitError(err, &gitErr) {
			t.Errorf(
				"Commit() error type = %T, want *specterrs.GitCommitError",
				err,
			)
		}

		if result.CommitHash != "" {
			t.Errorf(
				"Commit().CommitHash = %q, want empty on error",
				result.CommitHash,
			)
		}
	})
}

// isGitCommitError checks if the error is a GitCommitError and sets target.
func isGitCommitError(err error, target **specterrs.GitCommitError) bool {
	return errors.As(err, target)
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
			mock := &MockGitExecutor{
				StatusOutput:   "?? file.go\n",
				RevParseOutput: "abc123def456789012345678901234567890abcd",
			}
			c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

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

			// Verify the commit message was passed to the executor
			if len(mock.CommitMessages) != 1 {
				t.Fatalf("Expected 1 commit message, got %d", len(mock.CommitMessages))
			}
			if !strings.Contains(mock.CommitMessages[0], tt.wantContains) {
				t.Errorf(
					"Commit message passed to executor = %q, want to contain %q",
					mock.CommitMessages[0], tt.wantContains,
				)
			}
		})
	}
}

func TestCommitter_getModifiedFiles(t *testing.T) {
	t.Run("parses untracked files", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput: "?? new_file.go\n?? another.go\n",
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		files, err := c.getModifiedFiles()
		if err != nil {
			t.Fatalf("getModifiedFiles() error = %v", err)
		}

		if len(files) != 2 {
			t.Errorf("getModifiedFiles() returned %d files, want 2", len(files))
		}

		found := false
		for _, f := range files {
			if f == "new_file.go" {
				found = true

				break
			}
		}
		if !found {
			t.Error("getModifiedFiles() should include new_file.go")
		}
	})

	t.Run("parses modified files", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput: " M README.md\nM  staged.go\nMM both.go\n",
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		files, err := c.getModifiedFiles()
		if err != nil {
			t.Fatalf("getModifiedFiles() error = %v", err)
		}

		if len(files) != 3 {
			t.Errorf("getModifiedFiles() returned %d files, want 3", len(files))
		}

		for _, expected := range []string{"README.md", "staged.go", "both.go"} {
			found := false
			for _, f := range files {
				if f == expected {
					found = true

					break
				}
			}
			if !found {
				t.Errorf("getModifiedFiles() should include %s", expected)
			}
		}
	})

	t.Run("excludes deleted files", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput: " D deleted.go\nD  staged_delete.go\n?? new.go\n",
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		files, err := c.getModifiedFiles()
		if err != nil {
			t.Fatalf("getModifiedFiles() error = %v", err)
		}

		if len(files) != 1 {
			t.Errorf("getModifiedFiles() returned %d files, want 1", len(files))
		}

		for _, f := range files {
			if f == "deleted.go" || f == "staged_delete.go" {
				t.Errorf("getModifiedFiles() should not include deleted file %s", f)
			}
		}
	})

	t.Run("handles empty status", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusOutput: "",
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		files, err := c.getModifiedFiles()
		if err != nil {
			t.Fatalf("getModifiedFiles() error = %v", err)
		}

		if len(files) != 0 {
			t.Errorf("getModifiedFiles() returned %d files, want 0", len(files))
		}
	})

	t.Run("returns error on status failure", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusError: errors.New("git status failed"),
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		_, err := c.getModifiedFiles()
		if err == nil {
			t.Fatal("getModifiedFiles() expected error, got nil")
		}
	})
}

func TestCommitter_stageFiles(t *testing.T) {
	t.Run("stages single file", func(t *testing.T) {
		mock := &MockGitExecutor{}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		err := c.stageFiles([]string{"stage_test.go"})
		if err != nil {
			t.Fatalf("stageFiles() error = %v", err)
		}

		// Verify file was passed to Add
		if len(mock.AddedFiles) != 1 {
			t.Errorf("AddedFiles count = %d, want 1", len(mock.AddedFiles))
		}
		if mock.AddedFiles[0] != "stage_test.go" {
			t.Errorf(
				"AddedFiles[0] = %q, want %q",
				mock.AddedFiles[0], "stage_test.go",
			)
		}
	})

	t.Run("stages multiple files", func(t *testing.T) {
		mock := &MockGitExecutor{}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		err := c.stageFiles([]string{"file1.go", "file2.go"})
		if err != nil {
			t.Fatalf("stageFiles() error = %v", err)
		}

		// Verify files were passed to Add
		if len(mock.AddedFiles) != 2 {
			t.Errorf("AddedFiles count = %d, want 2", len(mock.AddedFiles))
		}
	})

	t.Run("returns error on add failure", func(t *testing.T) {
		mock := &MockGitExecutor{
			AddError: errors.New("git add failed"),
		}
		c := NewCommitterWithExecutor("test-change", "/repo", false, mock)

		err := c.stageFiles([]string{"file.go"})
		if err == nil {
			t.Fatal("stageFiles() expected error, got nil")
		}
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

func TestParseGitStatus(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []string
	}{
		{
			name:   "untracked files",
			output: "?? file1.go\n?? file2.go\n",
			want:   []string{"file1.go", "file2.go"},
		},
		{
			name:   "modified files - unstaged",
			output: " M file.go\n",
			want:   []string{"file.go"},
		},
		{
			name:   "modified files - staged",
			output: "M  file.go\n",
			want:   []string{"file.go"},
		},
		{
			name:   "added files",
			output: "A  file.go\n",
			want:   []string{"file.go"},
		},
		{
			name:   "renamed files",
			output: "R  old.go -> new.go\n",
			want:   []string{"old.go -> new.go"},
		},
		{
			name:   "deleted files - excluded",
			output: " D deleted.go\nD  staged_delete.go\n",
			want:   nil,
		},
		{
			name:   "mixed status",
			output: "?? new.go\n M modified.go\nD  deleted.go\nA  added.go\n",
			want:   []string{"new.go", "modified.go", "added.go"},
		},
		{
			name:   "empty output",
			output: "",
			want:   nil,
		},
		{
			name:   "files in subdirectories",
			output: "?? src/pkg/file.go\n M internal/track/committer.go\n",
			want:   []string{"src/pkg/file.go", "internal/track/committer.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGitStatus(tt.output)

			if len(got) != len(tt.want) {
				t.Errorf(
					"parseGitStatus() returned %d files, want %d",
					len(got), len(tt.want),
				)

				return
			}

			for i, f := range got {
				if f != tt.want[i] {
					t.Errorf(
						"parseGitStatus()[%d] = %q, want %q",
						i, f, tt.want[i],
					)
				}
			}
		})
	}
}
