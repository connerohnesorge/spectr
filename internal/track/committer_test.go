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
func (m *MockGitExecutor) Status(
	_ string,
) (string, error) {
	m.StatusCalls++

	return m.StatusOutput, m.StatusError
}

// Add implements GitExecutor.Add.
func (m *MockGitExecutor) Add(
	_ string,
	files []string,
) error {
	m.AddCalls++
	m.AddedFiles = append(m.AddedFiles, files...)

	return m.AddError
}

// Commit implements GitExecutor.Commit.
func (m *MockGitExecutor) Commit(
	_, message string,
) error {
	m.CommitCalls++
	m.CommitMessages = append(
		m.CommitMessages,
		message,
	)

	return m.CommitError
}

// RevParse implements GitExecutor.RevParse.
func (m *MockGitExecutor) RevParse(
	_, _ string,
) (string, error) {
	m.RevParseCalls++

	return m.RevParseOutput, m.RevParseError
}

// DiffNumstat implements GitExecutor.DiffNumstat.
func (m *MockGitExecutor) DiffNumstat(
	_ string,
	_ []string,
) (string, error) {
	m.DiffNumstatCalls++

	return m.DiffNumstatOutput, m.DiffNumstatError
}

func TestNewCommitter(t *testing.T) {
	t.Run(
		"creates committer with correct fields",
		func(t *testing.T) {
			changeID := "test-change-123"
			repoRoot := "/path/to/repo"

			c := NewCommitter(
				changeID,
				repoRoot,
				false,
			)

			if c == nil {
				t.Fatal(
					"NewCommitter returned nil",
				)
			}
			if c.changeID != changeID {
				t.Errorf(
					"NewCommitter().changeID = %q, want %q",
					c.changeID,
					changeID,
				)
			}
			if c.repoRoot != repoRoot {
				t.Errorf(
					"NewCommitter().repoRoot = %q, want %q",
					c.repoRoot,
					repoRoot,
				)
			}
			if c.gitExecutor == nil {
				t.Error(
					"NewCommitter().gitExecutor should not be nil",
				)
			}
			if c.includeBinaries {
				t.Error(
					"NewCommitter().includeBinaries should be false",
				)
			}
		},
	)

	t.Run(
		"creates independent instances",
		func(t *testing.T) {
			c1 := NewCommitter(
				"change-1",
				"/repo1",
				false,
			)
			c2 := NewCommitter(
				"change-2",
				"/repo2",
				true,
			)

			if c1 == c2 {
				t.Error(
					"NewCommitter returned same instance for different calls",
				)
			}
			if c1.changeID == c2.changeID {
				t.Error(
					"Different committers have same changeID",
				)
			}
		},
	)

	t.Run(
		"includes binaries when flag is true",
		func(t *testing.T) {
			c := NewCommitter(
				"test-change",
				"/repo",
				true,
			)

			if !c.includeBinaries {
				t.Error(
					"NewCommitter().includeBinaries should be true",
				)
			}
		},
	)
}

func TestNewCommitterWithExecutor(t *testing.T) {
	t.Run(
		"uses provided executor",
		func(t *testing.T) {
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			if c.gitExecutor == nil {
				t.Error(
					"NewCommitterWithExecutor did not set executor",
				)
			}
		},
	)
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
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCommitter_buildCommitMessage(
	t *testing.T,
) {
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
			c := NewCommitterWithExecutor(
				tt.changeID,
				"/repo",
				false,
				mock,
			)
			got := c.buildCommitMessage(
				tt.taskID,
				tt.action,
			)

			if !strings.HasPrefix(
				got,
				tt.wantPrefix,
			) {
				t.Errorf(
					"buildCommitMessage() prefix = %q, want prefix %q",
					got,
					tt.wantPrefix,
				)
			}
			if !strings.HasSuffix(
				got,
				tt.wantFooter,
			) {
				t.Errorf(
					"buildCommitMessage() suffix = %q, want suffix %q",
					got,
					tt.wantFooter,
				)
			}
			// Verify double newline separates message and footer
			if !strings.Contains(got, "\n\n") {
				t.Error(
					"buildCommitMessage() missing double newline separator",
				)
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
			name: "filters out tasks.json",
			files: []string{
				"src/main.go",
				"tasks.json",
				"README.md",
			},
			want: []string{
				"src/main.go",
				"README.md",
			},
		},
		{
			name: "filters out tasks.jsonc",
			files: []string{
				"src/main.go",
				"tasks.jsonc",
			},
			want: []string{"src/main.go"},
		},
		{
			name: "filters out tasks.md",
			files: []string{
				"tasks.md",
				"src/main.go",
			},
			want: []string{"src/main.go"},
		},
		{
			name: "filters all task files",
			files: []string{
				"tasks.json",
				"tasks.jsonc",
				"tasks.md",
			},
			want: nil,
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
			name: "preserves files with similar names",
			files: []string{
				"tasks.json.bak",
				"my-tasks.json",
				"tasks.jsonc.old",
			},
			want: []string{
				"tasks.json.bak",
				"my-tasks.json",
				"tasks.jsonc.old",
			},
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
					len(got),
					len(tt.want),
				)
			}
			for i, f := range got {
				if i >= len(tt.want) ||
					f != tt.want[i] {
					t.Errorf(
						"filterTaskFiles()[%d] = %q, want %q",
						i,
						f,
						tt.want[i],
					)
				}
			}
		})
	}
}

func TestCommitter_Commit_NoFilesToStage(
	t *testing.T,
) {
	// Mock returns only task files (which should be filtered out)
	mock := &MockGitExecutor{
		StatusOutput: "?? tasks.json\n?? tasks.jsonc\n?? tasks.md\n",
	}
	c := NewCommitterWithExecutor(
		"test-change",
		"/repo",
		false,
		mock,
	)

	result, err := c.Commit("1.1", ActionStart)
	if err != nil {
		t.Fatalf(
			"Commit() error = %v, want nil",
			err,
		)
	}
	if !result.NoFiles {
		t.Error(
			"Commit().NoFiles = false, want true",
		)
	}
	if result.CommitHash != "" {
		t.Errorf(
			"Commit().CommitHash = %q, want empty",
			result.CommitHash,
		)
	}
	// Message should still be set even when no files
	if result.Message == "" {
		t.Error(
			"Commit().Message should not be empty",
		)
	}
	// Verify Status was called but Add and Commit were not
	if mock.StatusCalls != 1 {
		t.Errorf(
			"Status should be called once, got %d",
			mock.StatusCalls,
		)
	}
	if mock.AddCalls != 0 {
		t.Errorf(
			"Add should not be called when no files, got %d",
			mock.AddCalls,
		)
	}
	if mock.CommitCalls != 0 {
		t.Errorf(
			"Commit should not be called when no files, got %d",
			mock.CommitCalls,
		)
	}
}

func TestCommitter_Commit_WithFiles(
	t *testing.T,
) {
	// Mock returns a non-task file
	mock := &MockGitExecutor{
		StatusOutput:   "?? src/main.go\n",
		RevParseOutput: "abc123def456789012345678901234567890abcd",
	}
	c := NewCommitterWithExecutor(
		"test-change",
		"/repo",
		false,
		mock,
	)

	result, err := c.Commit("1.1", ActionStart)
	if err != nil {
		t.Fatalf(
			"Commit() error = %v, want nil",
			err,
		)
	}
	if result.NoFiles {
		t.Error(
			"Commit().NoFiles = true, want false",
		)
	}
	if result.CommitHash == "" {
		t.Error(
			"Commit().CommitHash should not be empty",
		)
	}
	// Verify commit hash
	if result.CommitHash != "abc123def456789012345678901234567890abcd" {
		t.Errorf(
			"Commit().CommitHash = %q, want %q",
			result.CommitHash,
			"abc123def456789012345678901234567890abcd",
		)
	}

	// Verify message format
	expectedPrefix := "spectr(test-change): start task 1.1"
	if !strings.HasPrefix(
		result.Message,
		expectedPrefix,
	) {
		t.Errorf(
			"Commit().Message prefix = %q, want %q",
			result.Message,
			expectedPrefix,
		)
	}

	// Verify the mock was called correctly
	if mock.StatusCalls != 1 {
		t.Errorf(
			"Status should be called once, got %d",
			mock.StatusCalls,
		)
	}
	if mock.AddCalls != 1 {
		t.Errorf(
			"Add should be called once, got %d",
			mock.AddCalls,
		)
	}
	if mock.CommitCalls != 1 {
		t.Errorf(
			"Commit should be called once, got %d",
			mock.CommitCalls,
		)
	}
	if mock.RevParseCalls != 1 {
		t.Errorf(
			"RevParse should be called once, got %d",
			mock.RevParseCalls,
		)
	}

	// Verify the correct file was staged
	if len(mock.AddedFiles) != 1 ||
		mock.AddedFiles[0] != "src/main.go" {
		t.Errorf(
			"AddedFiles = %v, want [src/main.go]",
			mock.AddedFiles,
		)
	}
}

func TestCommitter_Commit_WithMixedFiles(
	t *testing.T,
) {
	// Mock returns both task files and regular files
	mock := &MockGitExecutor{
		StatusOutput:   "?? tasks.jsonc\n?? src/fix.go\n?? tests/fix_test.go\n",
		RevParseOutput: "def456789012345678901234567890abcdef1234",
	}
	c := NewCommitterWithExecutor(
		"fix-bug",
		"/repo",
		false,
		mock,
	)

	result, err := c.Commit("2.1", ActionComplete)
	if err != nil {
		t.Fatalf(
			"Commit() error = %v, want nil",
			err,
		)
	}
	if result.NoFiles {
		t.Error(
			"Commit().NoFiles = true, want false (should stage non-task files)",
		)
	}
	if result.CommitHash == "" {
		t.Error(
			"Commit().CommitHash should not be empty",
		)
	}

	// Verify only non-task files were staged
	expectedFiles := []string{
		"src/fix.go",
		"tests/fix_test.go",
	}
	if len(
		mock.AddedFiles,
	) != len(
		expectedFiles,
	) {
		t.Errorf(
			"AddedFiles count = %d, want %d",
			len(
				mock.AddedFiles,
			),
			len(expectedFiles),
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
			t.Errorf(
				"Expected file %q to be staged",
				expected,
			)
		}
	}
	// Verify task file was NOT staged
	for _, added := range mock.AddedFiles {
		if added == "tasks.jsonc" {
			t.Error(
				"tasks.jsonc should not be staged",
			)
		}
	}
}

func TestCommitter_Commit_WithModifiedFiles(
	t *testing.T,
) {
	// Mock returns modified files (unstaged and staged)
	mock := &MockGitExecutor{
		StatusOutput:   " M src/modified.go\nM  src/staged.go\n",
		RevParseOutput: "abc123def456789012345678901234567890abcd",
	}
	c := NewCommitterWithExecutor(
		"test-change",
		"/repo",
		false,
		mock,
	)

	result, err := c.Commit("1.1", ActionStart)
	if err != nil {
		t.Fatalf(
			"Commit() error = %v, want nil",
			err,
		)
	}
	if result.NoFiles {
		t.Error(
			"Commit().NoFiles = true, want false (should stage modified files)",
		)
	}
	if result.CommitHash == "" {
		t.Error(
			"Commit().CommitHash should not be empty",
		)
	}

	// Verify both modified files were staged
	expectedFiles := []string{
		"src/modified.go",
		"src/staged.go",
	}
	if len(
		mock.AddedFiles,
	) != len(
		expectedFiles,
	) {
		t.Errorf(
			"AddedFiles count = %d, want %d",
			len(
				mock.AddedFiles,
			),
			len(expectedFiles),
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
			t.Errorf(
				"Expected modified file %q to be staged",
				expected,
			)
		}
	}
}

func TestCommitter_Commit_WithMixedNewAndModified(
	t *testing.T,
) {
	// Mock returns both new (untracked) and modified files
	mock := &MockGitExecutor{
		StatusOutput:   "?? src/new.go\n M src/modified.go\nM  src/staged.go\n",
		RevParseOutput: "def456789012345678901234567890abcdef1234",
	}
	c := NewCommitterWithExecutor(
		"test-change",
		"/repo",
		false,
		mock,
	)

	result, err := c.Commit("1.1", ActionComplete)
	if err != nil {
		t.Fatalf(
			"Commit() error = %v, want nil",
			err,
		)
	}
	if result.NoFiles {
		t.Error(
			"Commit().NoFiles = true, want false",
		)
	}
	if result.CommitHash == "" {
		t.Error(
			"Commit().CommitHash should not be empty",
		)
	}

	// Verify all three files were staged (new + modified)
	expectedFiles := []string{
		"src/new.go",
		"src/modified.go",
		"src/staged.go",
	}
	if len(
		mock.AddedFiles,
	) != len(
		expectedFiles,
	) {
		t.Errorf(
			"AddedFiles count = %d, want %d",
			len(
				mock.AddedFiles,
			),
			len(expectedFiles),
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
			t.Errorf(
				"Expected file %q to be staged",
				expected,
			)
		}
	}
}

func TestCommitter_Commit_GitError(t *testing.T) {
	t.Run("status error", func(t *testing.T) {
		mock := &MockGitExecutor{
			StatusError: errors.New(
				"git status failed: not a git repository",
			),
		}
		c := NewCommitterWithExecutor(
			"test-change",
			"/nonexistent",
			false,
			mock,
		)

		result, err := c.Commit(
			"1.1",
			ActionStart,
		)

		if err == nil {
			t.Fatal(
				"Commit() expected error, got nil",
			)
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
			AddError: errors.New(
				"git add failed: permission denied",
			),
		}
		c := NewCommitterWithExecutor(
			"test-change",
			"/repo",
			false,
			mock,
		)

		result, err := c.Commit(
			"1.1",
			ActionStart,
		)

		if err == nil {
			t.Fatal(
				"Commit() expected error, got nil",
			)
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
			CommitError: errors.New(
				"git commit failed: nothing to commit",
			),
		}
		c := NewCommitterWithExecutor(
			"test-change",
			"/repo",
			false,
			mock,
		)

		result, err := c.Commit(
			"1.1",
			ActionStart,
		)

		if err == nil {
			t.Fatal(
				"Commit() expected error, got nil",
			)
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
			StatusOutput: "?? src/main.go\n",
			RevParseError: errors.New(
				"git rev-parse failed: ambiguous argument",
			),
		}
		c := NewCommitterWithExecutor(
			"test-change",
			"/repo",
			false,
			mock,
		)

		result, err := c.Commit(
			"1.1",
			ActionStart,
		)

		if err == nil {
			t.Fatal(
				"Commit() expected error, got nil",
			)
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
func isGitCommitError(
	err error,
	target **specterrs.GitCommitError,
) bool {
	return errors.As(err, target)
}

func TestCommitter_Commit_ActionTypes(
	t *testing.T,
) {
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
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				tt.action,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			if !strings.Contains(
				result.Message,
				tt.wantContains,
			) {
				t.Errorf(
					"Commit().Message = %q, want to contain %q",
					result.Message,
					tt.wantContains,
				)
			}

			// Verify the commit message was passed to the executor
			if len(mock.CommitMessages) != 1 {
				t.Fatalf(
					"Expected 1 commit message, got %d",
					len(mock.CommitMessages),
				)
			}
			if !strings.Contains(
				mock.CommitMessages[0],
				tt.wantContains,
			) {
				t.Errorf(
					"Commit message passed to executor = %q, want to contain %q",
					mock.CommitMessages[0],
					tt.wantContains,
				)
			}
		})
	}
}

func TestCommitter_getModifiedFiles(
	t *testing.T,
) {
	t.Run(
		"parses untracked files",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "?? new_file.go\n?? another.go\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files, err := c.getModifiedFiles()
			if err != nil {
				t.Fatalf(
					"getModifiedFiles() error = %v",
					err,
				)
			}

			if len(files) != 2 {
				t.Errorf(
					"getModifiedFiles() returned %d files, want 2",
					len(files),
				)
			}

			found := false
			for _, f := range files {
				if f == "new_file.go" {
					found = true

					break
				}
			}
			if !found {
				t.Error(
					"getModifiedFiles() should include new_file.go",
				)
			}
		},
	)

	t.Run(
		"parses modified files",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: " M README.md\nM  staged.go\nMM both.go\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files, err := c.getModifiedFiles()
			if err != nil {
				t.Fatalf(
					"getModifiedFiles() error = %v",
					err,
				)
			}

			if len(files) != 3 {
				t.Errorf(
					"getModifiedFiles() returned %d files, want 3",
					len(files),
				)
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
					t.Errorf(
						"getModifiedFiles() should include %s",
						expected,
					)
				}
			}
		},
	)

	t.Run(
		"excludes deleted files",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: " D deleted.go\nD  staged_delete.go\n?? new.go\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files, err := c.getModifiedFiles()
			if err != nil {
				t.Fatalf(
					"getModifiedFiles() error = %v",
					err,
				)
			}

			if len(files) != 1 {
				t.Errorf(
					"getModifiedFiles() returned %d files, want 1",
					len(files),
				)
			}

			for _, f := range files {
				if f == "deleted.go" ||
					f == "staged_delete.go" {
					t.Errorf(
						"getModifiedFiles() should not include deleted file %s",
						f,
					)
				}
			}
		},
	)

	t.Run(
		"handles empty status",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files, err := c.getModifiedFiles()
			if err != nil {
				t.Fatalf(
					"getModifiedFiles() error = %v",
					err,
				)
			}

			if len(files) != 0 {
				t.Errorf(
					"getModifiedFiles() returned %d files, want 0",
					len(files),
				)
			}
		},
	)

	t.Run(
		"returns error on status failure",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusError: errors.New(
					"git status failed",
				),
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			_, err := c.getModifiedFiles()
			if err == nil {
				t.Fatal(
					"getModifiedFiles() expected error, got nil",
				)
			}
		},
	)
}

func TestCommitter_stageFiles(t *testing.T) {
	t.Run(
		"stages single file",
		func(t *testing.T) {
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			err := c.stageFiles(
				[]string{"stage_test.go"},
			)
			if err != nil {
				t.Fatalf(
					"stageFiles() error = %v",
					err,
				)
			}

			// Verify file was passed to Add
			if len(mock.AddedFiles) != 1 {
				t.Errorf(
					"AddedFiles count = %d, want 1",
					len(mock.AddedFiles),
				)
			}
			if mock.AddedFiles[0] != "stage_test.go" {
				t.Errorf(
					"AddedFiles[0] = %q, want %q",
					mock.AddedFiles[0],
					"stage_test.go",
				)
			}
		},
	)

	t.Run(
		"stages multiple files",
		func(t *testing.T) {
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			err := c.stageFiles(
				[]string{"file1.go", "file2.go"},
			)
			if err != nil {
				t.Fatalf(
					"stageFiles() error = %v",
					err,
				)
			}

			// Verify files were passed to Add
			if len(mock.AddedFiles) != 2 {
				t.Errorf(
					"AddedFiles count = %d, want 2",
					len(mock.AddedFiles),
				)
			}
		},
	)

	t.Run(
		"returns error on add failure",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				AddError: errors.New(
					"git add failed",
				),
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			err := c.stageFiles(
				[]string{"file.go"},
			)
			if err == nil {
				t.Fatal(
					"stageFiles() expected error, got nil",
				)
			}
		},
	)
}

func TestCommitResult_Struct(t *testing.T) {
	t.Run("NoFiles result", func(t *testing.T) {
		result := CommitResult{
			NoFiles:    true,
			CommitHash: "",
			Message:    "spectr(test): start task 1.1",
		}

		if !result.NoFiles {
			t.Error(
				"CommitResult.NoFiles should be true",
			)
		}
		if result.CommitHash != "" {
			t.Error(
				"CommitResult.CommitHash should be empty for NoFiles",
			)
		}
	})

	t.Run("WithFiles result", func(t *testing.T) {
		result := CommitResult{
			NoFiles:    false,
			CommitHash: "abc123def456789012345678901234567890abcd",
			Message:    "spectr(test): complete task 1.1\n\n[Automated by spectr track]",
		}

		if result.NoFiles {
			t.Error(
				"CommitResult.NoFiles should be false",
			)
		}
		if result.CommitHash == "" {
			t.Error(
				"CommitResult.CommitHash should not be empty",
			)
		}
		if result.Message == "" {
			t.Error(
				"CommitResult.Message should not be empty",
			)
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
			want: []string{
				"file1.go",
				"file2.go",
			},
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
			want: []string{
				"new.go",
				"modified.go",
				"added.go",
			},
		},
		{
			name:   "empty output",
			output: "",
			want:   nil,
		},
		{
			name:   "files in subdirectories",
			output: "?? src/pkg/file.go\n M internal/track/committer.go\n",
			want: []string{
				"src/pkg/file.go",
				"internal/track/committer.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGitStatus(tt.output)

			if len(got) != len(tt.want) {
				t.Errorf(
					"parseGitStatus() returned %d files, want %d",
					len(got),
					len(tt.want),
				)

				return
			}

			for i, f := range got {
				if f != tt.want[i] {
					t.Errorf(
						"parseGitStatus()[%d] = %q, want %q",
						i,
						f,
						tt.want[i],
					)
				}
			}
		})
	}
}

// ============================================================================
// Task 2.1: Unit tests for binary detection (parseBinaryFilesFromNumstat)
// ============================================================================

func TestParseBinaryFilesFromNumstat(
	t *testing.T,
) {
	tests := []struct {
		name  string
		input string
		want  map[string]bool
	}{
		{
			name:  "empty input",
			input: "",
			want:  make(map[string]bool),
		},
		{
			name:  "single text file",
			input: "10\t5\tsrc/main.go\n",
			want:  make(map[string]bool),
		},
		{
			name:  "single binary file",
			input: "-\t-\timage.png\n",
			want: map[string]bool{
				"image.png": true,
			},
		},
		{
			name:  "multiple text files",
			input: "10\t5\tsrc/main.go\n20\t3\tsrc/utils.go\n1\t1\tREADME.md\n",
			want:  make(map[string]bool),
		},
		{
			name:  "multiple binary files",
			input: "-\t-\timage.png\n-\t-\tlogo.jpg\n-\t-\tapp.exe\n",
			want: map[string]bool{
				"image.png": true,
				"logo.jpg":  true,
				"app.exe":   true,
			},
		},
		{
			name:  "mixed binary and text files",
			input: "10\t5\tsrc/main.go\n-\t-\timage.png\n20\t3\tsrc/utils.go\n-\t-\tlogo.jpg\n",
			want: map[string]bool{
				"image.png": true,
				"logo.jpg":  true,
			},
		},
		{
			name:  "whitespace handling",
			input: "  \n10\t5\tsrc/main.go\n  -\t-\timage.png  \n\n",
			want: map[string]bool{
				"image.png": true,
			},
		},
		{
			name:  "file with spaces in name",
			input: "-\t-\tpath/to/my file.png\n",
			want: map[string]bool{
				"path/to/my file.png": true,
			},
		},
		{
			name:  "file in subdirectory",
			input: "-\t-\tassets/images/logo.png\n10\t5\tsrc/pkg/file.go\n",
			want: map[string]bool{
				"assets/images/logo.png": true,
			},
		},
		{
			name:  "zero additions and deletions is not binary",
			input: "0\t0\tempty.txt\n",
			want:  make(map[string]bool),
		},
		{
			name:  "malformed line - too few parts",
			input: "-\t-\n10\n",
			want:  make(map[string]bool),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBinaryFilesFromNumstat(
				tt.input,
			)

			if len(got) != len(tt.want) {
				t.Errorf(
					"parseBinaryFilesFromNumstat() returned %d entries, want %d",
					len(got),
					len(tt.want),
				)

				return
			}

			for file, wantBinary := range tt.want {
				if got[file] != wantBinary {
					t.Errorf(
						"parseBinaryFilesFromNumstat()[%q] = %v, want %v",
						file,
						got[file],
						wantBinary,
					)
				}
			}

			// Verify no unexpected files in result
			for file := range got {
				if !tt.want[file] {
					t.Errorf(
						"parseBinaryFilesFromNumstat() unexpected file %q",
						file,
					)
				}
			}
		})
	}
}

func TestCommitter_getBinaryFiles(t *testing.T) {
	t.Run(
		"empty file list returns empty map",
		func(t *testing.T) {
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.getBinaryFiles(nil)
			if err != nil {
				t.Fatalf(
					"getBinaryFiles() error = %v",
					err,
				)
			}

			if len(result) != 0 {
				t.Errorf(
					"getBinaryFiles() returned %d entries, want 0",
					len(result),
				)
			}

			// DiffNumstat should not be called for empty list
			if mock.DiffNumstatCalls != 0 {
				t.Errorf(
					"DiffNumstat should not be called, got %d calls",
					mock.DiffNumstatCalls,
				)
			}
		},
	)

	t.Run(
		"detects binary files from numstat output",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				DiffNumstatOutput: "10\t5\tmain.go\n-\t-\timage.png\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.getBinaryFiles(
				[]string{"main.go", "image.png"},
			)
			if err != nil {
				t.Fatalf(
					"getBinaryFiles() error = %v",
					err,
				)
			}

			if !result["image.png"] {
				t.Error(
					"getBinaryFiles() should detect image.png as binary",
				)
			}
			if result["main.go"] {
				t.Error(
					"getBinaryFiles() should not detect main.go as binary",
				)
			}

			if mock.DiffNumstatCalls != 1 {
				t.Errorf(
					"DiffNumstat should be called once, got %d",
					mock.DiffNumstatCalls,
				)
			}
		},
	)

	t.Run(
		"returns error on DiffNumstat failure",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				DiffNumstatError: errors.New(
					"git diff failed",
				),
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			_, err := c.getBinaryFiles(
				[]string{"file.go"},
			)
			if err == nil {
				t.Fatal(
					"getBinaryFiles() expected error, got nil",
				)
			}
		},
	)
}

// ============================================================================
// Task 2.2: Unit tests for filtering with IncludeBinaries=false
// ============================================================================

func TestCommitter_filterFiles_ExcludeBinaries(
	t *testing.T,
) {
	t.Run(
		"excludes binary files when IncludeBinaries is false",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				DiffNumstatOutput: "10\t5\tmain.go\n-\t-\timage.png\n20\t3\tutils.go\n-\t-\tlogo.jpg\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files := []string{
				"main.go",
				"image.png",
				"utils.go",
				"logo.jpg",
			}
			filtered, skipped, err := c.filterFiles(
				files,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() error = %v",
					err,
				)
			}

			// Verify binary files are excluded
			expectedFiltered := []string{
				"main.go",
				"utils.go",
			}
			if len(
				filtered,
			) != len(
				expectedFiltered,
			) {
				t.Errorf(
					"filterFiles() filtered %d files, want %d",
					len(filtered),
					len(expectedFiltered),
				)
			}
			for i, f := range filtered {
				if f != expectedFiltered[i] {
					t.Errorf(
						"filterFiles() filtered[%d] = %q, want %q",
						i,
						f,
						expectedFiltered[i],
					)
				}
			}

			// Verify binary files are in skipped list
			expectedSkipped := []string{
				"image.png",
				"logo.jpg",
			}
			if len(
				skipped,
			) != len(
				expectedSkipped,
			) {
				t.Errorf(
					"filterFiles() skipped %d files, want %d",
					len(skipped),
					len(expectedSkipped),
				)
			}
			for i, f := range skipped {
				if f != expectedSkipped[i] {
					t.Errorf(
						"filterFiles() skipped[%d] = %q, want %q",
						i,
						f,
						expectedSkipped[i],
					)
				}
			}
		},
	)

	t.Run(
		"empty input returns empty results",
		func(t *testing.T) {
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			filtered, skipped, err := c.filterFiles(
				nil,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() error = %v",
					err,
				)
			}

			if len(filtered) != 0 {
				t.Errorf(
					"filterFiles() filtered should be empty, got %v",
					filtered,
				)
			}
			if skipped != nil {
				t.Errorf(
					"filterFiles() skipped should be nil, got %v",
					skipped,
				)
			}
		},
	)

	t.Run(
		"all binary files returns empty filtered list",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				DiffNumstatOutput: "-\t-\timage.png\n-\t-\tlogo.jpg\n-\t-\tapp.exe\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files := []string{
				"image.png",
				"logo.jpg",
				"app.exe",
			}
			filtered, skipped, err := c.filterFiles(
				files,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() error = %v",
					err,
				)
			}

			if len(filtered) != 0 {
				t.Errorf(
					"filterFiles() should return empty filtered, got %v",
					filtered,
				)
			}
			if len(skipped) != 3 {
				t.Errorf(
					"filterFiles() should return 3 skipped, got %d",
					len(skipped),
				)
			}
		},
	)

	t.Run(
		"no binary files returns all files",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				DiffNumstatOutput: "10\t5\tmain.go\n20\t3\tutils.go\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files := []string{
				"main.go",
				"utils.go",
			}
			filtered, skipped, err := c.filterFiles(
				files,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() error = %v",
					err,
				)
			}

			if len(filtered) != 2 {
				t.Errorf(
					"filterFiles() should return 2 filtered, got %d",
					len(filtered),
				)
			}
			if len(skipped) != 0 {
				t.Errorf(
					"filterFiles() should return empty skipped, got %v",
					skipped,
				)
			}
		},
	)

	t.Run(
		"DiffNumstat error returns all files gracefully",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				DiffNumstatError: errors.New(
					"git diff failed",
				),
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			files := []string{
				"main.go",
				"image.png",
			}
			filtered, skipped, err := c.filterFiles(
				files,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() should not error on DiffNumstat failure, got %v",
					err,
				)
			}

			// Should return all files when binary detection fails
			if len(filtered) != 2 {
				t.Errorf(
					"filterFiles() should return all files on error, got %d",
					len(filtered),
				)
			}
			if skipped != nil {
				t.Errorf(
					"filterFiles() skipped should be nil on error, got %v",
					skipped,
				)
			}
		},
	)
}

func TestCommitter_Commit_SkipsBinaries(
	t *testing.T,
) {
	t.Run(
		"commit excludes binary files and populates SkippedBinaries",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput:      "?? main.go\n?? image.png\n?? utils.go\n",
				DiffNumstatOutput: "10\t5\tmain.go\n-\t-\timage.png\n20\t3\tutils.go\n",
				RevParseOutput:    "abc123def456789012345678901234567890abcd",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionStart,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// Verify binary files were skipped
			if len(result.SkippedBinaries) != 1 {
				t.Errorf(
					"Commit().SkippedBinaries = %v, want 1 file",
					result.SkippedBinaries,
				)
			}
			if len(result.SkippedBinaries) > 0 &&
				result.SkippedBinaries[0] != "image.png" {
				t.Errorf(
					"Commit().SkippedBinaries[0] = %q, want image.png",
					result.SkippedBinaries[0],
				)
			}

			// Verify only non-binary files were staged
			expectedFiles := []string{
				"main.go",
				"utils.go",
			}
			if len(
				mock.AddedFiles,
			) != len(
				expectedFiles,
			) {
				t.Errorf(
					"AddedFiles = %v, want %v",
					mock.AddedFiles,
					expectedFiles,
				)
			}
			for _, f := range mock.AddedFiles {
				if f == "image.png" {
					t.Error(
						"image.png should not be staged",
					)
				}
			}
		},
	)

	t.Run(
		"commit with only binary files returns NoFiles=true",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput:      "?? image.png\n?? logo.jpg\n",
				DiffNumstatOutput: "-\t-\timage.png\n-\t-\tlogo.jpg\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionStart,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			if !result.NoFiles {
				t.Error(
					"Commit().NoFiles should be true when only binary files present",
				)
			}
			if len(result.SkippedBinaries) != 2 {
				t.Errorf(
					"Commit().SkippedBinaries should have 2 files, got %d",
					len(result.SkippedBinaries),
				)
			}
			if mock.AddCalls != 0 {
				t.Errorf(
					"Add should not be called, got %d calls",
					mock.AddCalls,
				)
			}
			if mock.CommitCalls != 0 {
				t.Errorf(
					"Commit should not be called, got %d calls",
					mock.CommitCalls,
				)
			}
		},
	)
}

// ============================================================================
// Task 2.3: Unit tests for including binaries with IncludeBinaries=true
// ============================================================================

func TestCommitter_filterFiles_IncludeBinaries(
	t *testing.T,
) {
	t.Run(
		"includes all files when IncludeBinaries is true",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				// DiffNumstat should not be called when IncludeBinaries is true
				DiffNumstatOutput: "-\t-\timage.png\n",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				true,
				mock,
			)

			files := []string{
				"main.go",
				"image.png",
				"utils.go",
			}
			filtered, skipped, err := c.filterFiles(
				files,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() error = %v",
					err,
				)
			}

			// All files should be included
			if len(filtered) != 3 {
				t.Errorf(
					"filterFiles() should return all 3 files, got %d",
					len(filtered),
				)
			}
			for i, f := range files {
				if filtered[i] != f {
					t.Errorf(
						"filterFiles() filtered[%d] = %q, want %q",
						i,
						filtered[i],
						f,
					)
				}
			}

			// No files should be skipped
			if skipped != nil {
				t.Errorf(
					"filterFiles() skipped should be nil, got %v",
					skipped,
				)
			}

			// DiffNumstat should not be called
			if mock.DiffNumstatCalls != 0 {
				t.Errorf(
					"DiffNumstat should not be called when IncludeBinaries=true, got %d calls",
					mock.DiffNumstatCalls,
				)
			}
		},
	)

	t.Run(
		"empty input with IncludeBinaries=true",
		func(t *testing.T) {
			mock := &MockGitExecutor{}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				true,
				mock,
			)

			filtered, skipped, err := c.filterFiles(
				nil,
			)
			if err != nil {
				t.Fatalf(
					"filterFiles() error = %v",
					err,
				)
			}

			if len(filtered) != 0 {
				t.Errorf(
					"filterFiles() should return empty, got %v",
					filtered,
				)
			}
			if skipped != nil {
				t.Errorf(
					"filterFiles() skipped should be nil, got %v",
					skipped,
				)
			}
		},
	)
}

func TestCommitter_Commit_IncludesBinaries(
	t *testing.T,
) {
	t.Run(
		"commit includes binary files when IncludeBinaries=true",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput:   "?? main.go\n?? image.png\n?? utils.go\n",
				RevParseOutput: "abc123def456789012345678901234567890abcd",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				true,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionStart,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// SkippedBinaries should be empty
			if len(result.SkippedBinaries) != 0 {
				t.Errorf(
					"Commit().SkippedBinaries should be empty, got %v",
					result.SkippedBinaries,
				)
			}

			// All files should be staged
			if len(mock.AddedFiles) != 3 {
				t.Errorf(
					"AddedFiles should have 3 files, got %d",
					len(mock.AddedFiles),
				)
			}

			// Verify image.png was staged
			foundImage := false
			for _, f := range mock.AddedFiles {
				if f == "image.png" {
					foundImage = true

					break
				}
			}
			if !foundImage {
				t.Error(
					"image.png should be staged when IncludeBinaries=true",
				)
			}

			// DiffNumstat should not be called
			if mock.DiffNumstatCalls != 0 {
				t.Errorf(
					"DiffNumstat should not be called when IncludeBinaries=true, got %d calls",
					mock.DiffNumstatCalls,
				)
			}
		},
	)

	t.Run(
		"commit with only binary files succeeds when IncludeBinaries=true",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput:   "?? image.png\n?? logo.jpg\n",
				RevParseOutput: "def456789012345678901234567890abcdef1234",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				true,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionComplete,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			if result.NoFiles {
				t.Error(
					"Commit().NoFiles should be false when IncludeBinaries=true",
				)
			}
			if len(result.SkippedBinaries) != 0 {
				t.Errorf(
					"Commit().SkippedBinaries should be empty, got %v",
					result.SkippedBinaries,
				)
			}
			if len(mock.AddedFiles) != 2 {
				t.Errorf(
					"AddedFiles should have 2 files, got %d",
					len(mock.AddedFiles),
				)
			}
			if mock.CommitCalls != 1 {
				t.Errorf(
					"Commit should be called once, got %d",
					mock.CommitCalls,
				)
			}
		},
	)
}

// ============================================================================
// Task 2.4: Integration test for track command with binary files
// ============================================================================

func TestCommitter_Commit_IntegrationWithBinaryFiles(
	t *testing.T,
) {
	t.Run(
		"full commit flow excludes binaries with IncludeBinaries=false",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "?? src/main.go\n" +
					"?? assets/logo.png\n" +
					"?? internal/utils.go\n" +
					"?? docs/diagram.jpg\n" +
					"?? build/app.exe\n",
				DiffNumstatOutput: "100\t50\tsrc/main.go\n" +
					"-\t-\tassets/logo.png\n" +
					"30\t10\tinternal/utils.go\n" +
					"-\t-\tdocs/diagram.jpg\n" +
					"-\t-\tbuild/app.exe\n",
				RevParseOutput: "abc123def456789012345678901234567890abcd",
			}
			c := NewCommitterWithExecutor(
				"feature-x",
				"/repo",
				false,
				mock,
			)

			result, err := c.Commit(
				"2.1",
				ActionComplete,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// Verify commit was created
			if result.NoFiles {
				t.Error(
					"Commit should have created a commit",
				)
			}
			if result.CommitHash != "abc123def456789012345678901234567890abcd" {
				t.Errorf(
					"CommitHash = %q, want abc123...",
					result.CommitHash,
				)
			}

			// Verify message format
			expectedPrefix := "spectr(feature-x): complete task 2.1"
			if !strings.HasPrefix(
				result.Message,
				expectedPrefix,
			) {
				t.Errorf(
					"Message should start with %q, got %q",
					expectedPrefix,
					result.Message,
				)
			}

			// Verify 3 binary files were skipped
			if len(result.SkippedBinaries) != 3 {
				t.Errorf(
					"SkippedBinaries should have 3 files, got %d: %v",
					len(result.SkippedBinaries),
					result.SkippedBinaries,
				)
			}
			expectedSkipped := map[string]bool{
				"assets/logo.png":  true,
				"docs/diagram.jpg": true,
				"build/app.exe":    true,
			}
			for _, f := range result.SkippedBinaries {
				if !expectedSkipped[f] {
					t.Errorf(
						"Unexpected skipped file: %q",
						f,
					)
				}
			}

			// Verify only text files were staged
			if len(mock.AddedFiles) != 2 {
				t.Errorf(
					"AddedFiles should have 2 files, got %d: %v",
					len(mock.AddedFiles),
					mock.AddedFiles,
				)
			}
			for _, f := range mock.AddedFiles {
				if expectedSkipped[f] {
					t.Errorf(
						"Binary file %q should not be staged",
						f,
					)
				}
			}
		},
	)

	t.Run(
		"full commit flow includes all files with IncludeBinaries=true",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "?? src/main.go\n" +
					"?? assets/logo.png\n" +
					"?? internal/utils.go\n" +
					"?? docs/diagram.jpg\n",
				RevParseOutput: "def456789012345678901234567890abcdef1234",
			}
			c := NewCommitterWithExecutor(
				"feature-y",
				"/repo",
				true,
				mock,
			)

			result, err := c.Commit(
				"1.3",
				ActionStart,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// Verify commit was created
			if result.NoFiles {
				t.Error(
					"Commit should have created a commit",
				)
			}

			// No files should be skipped
			if len(result.SkippedBinaries) != 0 {
				t.Errorf(
					"SkippedBinaries should be empty, got %v",
					result.SkippedBinaries,
				)
			}

			// All 4 files should be staged
			if len(mock.AddedFiles) != 4 {
				t.Errorf(
					"AddedFiles should have 4 files, got %d",
					len(mock.AddedFiles),
				)
			}

			// DiffNumstat should not be called
			if mock.DiffNumstatCalls != 0 {
				t.Errorf(
					"DiffNumstat should not be called, got %d",
					mock.DiffNumstatCalls,
				)
			}
		},
	)

	t.Run(
		"mixed task files and binary files with IncludeBinaries=false",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "?? src/main.go\n" +
					"?? tasks.jsonc\n" +
					"?? image.png\n" +
					"?? tasks.md\n",
				DiffNumstatOutput: "50\t20\tsrc/main.go\n-\t-\timage.png\n",
				RevParseOutput:    "abc123def456789012345678901234567890abcd",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionStart,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// Only main.go should be staged (task files filtered, binary skipped)
			if len(mock.AddedFiles) != 1 ||
				mock.AddedFiles[0] != "src/main.go" {
				t.Errorf(
					"AddedFiles should be [src/main.go], got %v",
					mock.AddedFiles,
				)
			}

			// Binary file should be in SkippedBinaries
			if len(result.SkippedBinaries) != 1 ||
				result.SkippedBinaries[0] != "image.png" {
				t.Errorf(
					"SkippedBinaries should be [image.png], got %v",
					result.SkippedBinaries,
				)
			}
		},
	)
}

// ============================================================================
// Task 2.5: Test with various binary file types
// ============================================================================

func TestParseBinaryFilesFromNumstat_VariousFileTypes(
	t *testing.T,
) {
	tests := []struct {
		name         string
		input        string
		wantBinaries []string
		wantText     []string
	}{
		{
			name: "image files",
			input: "-\t-\tlogo.png\n" +
				"-\t-\tphoto.jpg\n" +
				"-\t-\ticon.gif\n" +
				"-\t-\tbackground.bmp\n" +
				"-\t-\tvector.svg\n" +
				"-\t-\timage.webp\n" +
				"-\t-\tpicture.tiff\n",
			wantBinaries: []string{
				"logo.png", "photo.jpg", "icon.gif", "background.bmp",
				"vector.svg", "image.webp", "picture.tiff",
			},
			wantText: nil,
		},
		{
			name: "executable files",
			input: "-\t-\tapp.exe\n" +
				"-\t-\tprogram.dll\n" +
				"-\t-\tlibrary.so\n" +
				"-\t-\tmodule.dylib\n" +
				"-\t-\tbinary\n",
			wantBinaries: []string{
				"app.exe", "program.dll", "library.so", "module.dylib", "binary",
			},
			wantText: nil,
		},
		{
			name: "archive files",
			input: "-\t-\tpackage.zip\n" +
				"-\t-\tarchive.tar.gz\n" +
				"-\t-\tbackup.rar\n" +
				"-\t-\tcompressed.7z\n" +
				"-\t-\tbundle.tar\n",
			wantBinaries: []string{
				"package.zip", "archive.tar.gz", "backup.rar",
				"compressed.7z", "bundle.tar",
			},
			wantText: nil,
		},
		{
			name: "document files - binary",
			input: "-\t-\tdocument.pdf\n" +
				"-\t-\tspreadsheet.xlsx\n" +
				"-\t-\tpresentation.pptx\n" +
				"-\t-\tword.docx\n",
			wantBinaries: []string{
				"document.pdf", "spreadsheet.xlsx",
				"presentation.pptx", "word.docx",
			},
			wantText: nil,
		},
		{
			name: "font files",
			input: "-\t-\tfont.ttf\n" +
				"-\t-\tfont.otf\n" +
				"-\t-\tfont.woff\n" +
				"-\t-\tfont.woff2\n" +
				"-\t-\tfont.eot\n",
			wantBinaries: []string{
				"font.ttf", "font.otf", "font.woff", "font.woff2", "font.eot",
			},
			wantText: nil,
		},
		{
			name: "mixed binary and source files",
			input: "100\t50\tsrc/main.go\n" +
				"-\t-\tassets/logo.png\n" +
				"200\t100\tinternal/handler.go\n" +
				"-\t-\tbuild/app.exe\n" +
				"50\t25\tREADME.md\n" +
				"-\t-\tdocs/architecture.pdf\n" +
				"30\t10\tMakefile\n" +
				"-\t-\tvendor/lib.so\n",
			wantBinaries: []string{
				"assets/logo.png", "build/app.exe",
				"docs/architecture.pdf", "vendor/lib.so",
			},
			wantText: []string{
				"src/main.go", "internal/handler.go", "README.md", "Makefile",
			},
		},
		{
			name: "database and data files",
			input: "-\t-\tdata.db\n" +
				"-\t-\tcache.sqlite\n" +
				"-\t-\tindex.idx\n",
			wantBinaries: []string{
				"data.db",
				"cache.sqlite",
				"index.idx",
			},
			wantText: nil,
		},
		{
			name: "files in nested directories",
			input: "-\t-\tproject/assets/images/icons/logo.png\n" +
				"-\t-\tproject/build/bin/release/app.exe\n" +
				"50\t25\tproject/src/pkg/handler/handler.go\n",
			wantBinaries: []string{
				"project/assets/images/icons/logo.png",
				"project/build/bin/release/app.exe",
			},
			wantText: []string{
				"project/src/pkg/handler/handler.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBinaryFilesFromNumstat(
				tt.input,
			)

			// Check all expected binary files are detected
			for _, binary := range tt.wantBinaries {
				if !result[binary] {
					t.Errorf(
						"Expected %q to be detected as binary",
						binary,
					)
				}
			}

			// Check text files are NOT detected as binary
			for _, text := range tt.wantText {
				if result[text] {
					t.Errorf(
						"Expected %q to NOT be detected as binary",
						text,
					)
				}
			}

			// Verify count matches
			if len(
				result,
			) != len(
				tt.wantBinaries,
			) {
				t.Errorf(
					"Detected %d binaries, want %d",
					len(result),
					len(tt.wantBinaries),
				)
			}
		})
	}
}

func TestCommitter_Commit_VariousBinaryTypes(
	t *testing.T,
) {
	t.Run(
		"skips various binary file types",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "?? main.go\n" +
					"?? logo.png\n" +
					"?? photo.jpg\n" +
					"?? icon.gif\n" +
					"?? app.exe\n" +
					"?? lib.dll\n" +
					"?? module.so\n" +
					"?? archive.zip\n" +
					"?? backup.tar.gz\n" +
					"?? doc.pdf\n",
				DiffNumstatOutput: "100\t50\tmain.go\n" +
					"-\t-\tlogo.png\n" +
					"-\t-\tphoto.jpg\n" +
					"-\t-\ticon.gif\n" +
					"-\t-\tapp.exe\n" +
					"-\t-\tlib.dll\n" +
					"-\t-\tmodule.so\n" +
					"-\t-\tarchive.zip\n" +
					"-\t-\tbackup.tar.gz\n" +
					"-\t-\tdoc.pdf\n",
				RevParseOutput: "abc123def456789012345678901234567890abcd",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				false,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionStart,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// Should skip 9 binary files
			if len(result.SkippedBinaries) != 9 {
				t.Errorf(
					"SkippedBinaries should have 9 files, got %d: %v",
					len(
						result.SkippedBinaries,
					),
					result.SkippedBinaries,
				)
			}

			// Only main.go should be staged
			if len(mock.AddedFiles) != 1 ||
				mock.AddedFiles[0] != "main.go" {
				t.Errorf(
					"AddedFiles should be [main.go], got %v",
					mock.AddedFiles,
				)
			}

			// Verify specific binary types are in skipped list
			expectedSkipped := map[string]bool{
				"logo.png":      true, // image
				"photo.jpg":     true, // image
				"icon.gif":      true, // image
				"app.exe":       true, // executable
				"lib.dll":       true, // library
				"module.so":     true, // shared object
				"archive.zip":   true, // archive
				"backup.tar.gz": true, // archive
				"doc.pdf":       true, // document
			}
			for _, f := range result.SkippedBinaries {
				if !expectedSkipped[f] {
					t.Errorf(
						"Unexpected skipped file: %q",
						f,
					)
				}
			}
		},
	)

	t.Run(
		"includes all binary types when IncludeBinaries=true",
		func(t *testing.T) {
			mock := &MockGitExecutor{
				StatusOutput: "?? main.go\n" +
					"?? logo.png\n" +
					"?? app.exe\n" +
					"?? archive.zip\n" +
					"?? doc.pdf\n",
				RevParseOutput: "def456789012345678901234567890abcdef1234",
			}
			c := NewCommitterWithExecutor(
				"test-change",
				"/repo",
				true,
				mock,
			)

			result, err := c.Commit(
				"1.1",
				ActionComplete,
			)
			if err != nil {
				t.Fatalf(
					"Commit() error = %v",
					err,
				)
			}

			// No files should be skipped
			if len(result.SkippedBinaries) != 0 {
				t.Errorf(
					"SkippedBinaries should be empty, got %v",
					result.SkippedBinaries,
				)
			}

			// All 5 files should be staged
			if len(mock.AddedFiles) != 5 {
				t.Errorf(
					"AddedFiles should have 5 files, got %d",
					len(mock.AddedFiles),
				)
			}
		},
	)
}
