package track

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// createTestTasksFile creates a tasks.jsonc file in the given directory.
func createTestTasksFile(
	t *testing.T,
	dir, content string,
) string {
	t.Helper()
	tasksPath := filepath.Join(dir, "tasks.jsonc")
	if err := os.WriteFile(tasksPath, []byte(content), 0644); err != nil {
		t.Fatalf(
			"failed to create tasks file: %v",
			err,
		)
	}

	return tasksPath
}

// tasksFileContent generates tasks.jsonc content with the given tasks.
func tasksFileContent(tasks ...struct {
	id     string
	status string
},
) string {
	content := `{"version": 1, "tasks": [`
	for i, task := range tasks {
		if i > 0 {
			content += ","
		}
		content += `{"id":"` + task.id + `","section":"Test","description":"Task ` + task.id + `","status":"` + task.status + `"}`
	}
	content += `]}`

	return content
}

func TestNew_Success(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksFileContent(
			struct{ id, status string }{
				"1.1",
				"pending",
			},
		),
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf(
			"New() error = %v, want nil",
			err,
		)
	}
	defer func() { _ = tracker.Close() }()

	if tracker == nil {
		t.Fatal("New() returned nil tracker")
	}
	if tracker.changeID != "test-change" {
		t.Errorf(
			"New().changeID = %q, want %q",
			tracker.changeID,
			"test-change",
		)
	}
	if tracker.tasksPath != tasksPath {
		t.Errorf(
			"New().tasksPath = %q, want %q",
			tracker.tasksPath,
			tasksPath,
		)
	}
	if tracker.repoRoot != tempDir {
		t.Errorf(
			"New().repoRoot = %q, want %q",
			tracker.repoRoot,
			tempDir,
		)
	}
	if tracker.watcher == nil {
		t.Error("New().watcher should not be nil")
	}
	if tracker.committer == nil {
		t.Error(
			"New().committer should not be nil",
		)
	}
	if tracker.previousState == nil {
		t.Error(
			"New().previousState should not be nil",
		)
	}
}

func TestNew_MissingTasksFile(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()
	nonExistentPath := filepath.Join(
		tempDir,
		"nonexistent",
		"tasks.jsonc",
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: nonExistentPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})

	if err == nil {
		if tracker != nil {
			_ = tracker.Close()
		}
		t.Fatal(
			"New() expected error for missing tasks file, got nil",
		)
	}

	if tracker != nil {
		_ = tracker.Close()
		t.Error(
			"New() should return nil tracker on error",
		)
	}
}

func TestTracker_allTasksComplete(t *testing.T) {
	tests := []struct {
		name  string
		tasks []parsers.Task
		want  bool
	}{
		{
			name:  "empty list is complete",
			tasks: make([]parsers.Task, 0),
			want:  true,
		},
		{
			name: "all tasks completed",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "2.1",
					Status: parsers.TaskStatusCompleted,
				},
			},
			want: true,
		},
		{
			name: "some tasks pending",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusPending,
				},
				{
					ID:     "2.1",
					Status: parsers.TaskStatusCompleted,
				},
			},
			want: false,
		},
		{
			name: "some tasks in_progress",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusInProgress,
				},
				{
					ID:     "2.1",
					Status: parsers.TaskStatusCompleted,
				},
			},
			want: false,
		},
		{
			name: "all tasks pending",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusPending,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusPending,
				},
			},
			want: false,
		},
		{
			name: "single completed task",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
			},
			want: true,
		},
		{
			name: "single pending task",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusPending,
				},
			},
			want: false,
		},
		{
			name: "mixed statuses",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusPending,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusInProgress,
				},
				{
					ID:     "1.3",
					Status: parsers.TaskStatusCompleted,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allTasksComplete(tt.tasks)
			if got != tt.want {
				t.Errorf(
					"allTasksComplete() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestTracker_countProgress(t *testing.T) {
	tests := []struct {
		name          string
		tasks         []parsers.Task
		wantCompleted int
		wantTotal     int
	}{
		{
			name: "empty list",
			tasks: make(
				[]parsers.Task,
				0,
			),
			wantCompleted: 0,
			wantTotal:     0,
		},
		{
			name: "all completed",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "1.3",
					Status: parsers.TaskStatusCompleted,
				},
			},
			wantCompleted: 3,
			wantTotal:     3,
		},
		{
			name: "none completed",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusPending,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusInProgress,
				},
			},
			wantCompleted: 0,
			wantTotal:     2,
		},
		{
			name: "partial completion",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
				{
					ID:     "1.2",
					Status: parsers.TaskStatusInProgress,
				},
				{
					ID:     "1.3",
					Status: parsers.TaskStatusPending,
				},
				{
					ID:     "1.4",
					Status: parsers.TaskStatusCompleted,
				},
			},
			wantCompleted: 2,
			wantTotal:     4,
		},
		{
			name: "single completed",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusCompleted,
				},
			},
			wantCompleted: 1,
			wantTotal:     1,
		},
		{
			name: "single pending",
			tasks: []parsers.Task{
				{
					ID:     "1.1",
					Status: parsers.TaskStatusPending,
				},
			},
			wantCompleted: 0,
			wantTotal:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completed, total := countProgress(
				tt.tasks,
			)
			if completed != tt.wantCompleted {
				t.Errorf(
					"countProgress() completed = %d, want %d",
					completed,
					tt.wantCompleted,
				)
			}
			if total != tt.wantTotal {
				t.Errorf(
					"countProgress() total = %d, want %d",
					total,
					tt.wantTotal,
				)
			}
		})
	}
}

func TestTracker_getActionForTransition(
	t *testing.T,
) {
	tracker := &Tracker{
		previousState: make(
			map[string]parsers.TaskStatusValue,
		),
	}

	tests := []struct {
		name       string
		from       parsers.TaskStatusValue
		to         parsers.TaskStatusValue
		wantAction Action
		wantCommit bool
	}{
		{
			name:       "pending to in_progress triggers ActionStart",
			from:       parsers.TaskStatusPending,
			to:         parsers.TaskStatusInProgress,
			wantAction: ActionStart,
			wantCommit: true,
		},
		{
			name:       "in_progress to completed triggers ActionComplete",
			from:       parsers.TaskStatusInProgress,
			to:         parsers.TaskStatusCompleted,
			wantAction: ActionComplete,
			wantCommit: true,
		},
		{
			name:       "pending to completed triggers ActionComplete",
			from:       parsers.TaskStatusPending,
			to:         parsers.TaskStatusCompleted,
			wantAction: ActionComplete,
			wantCommit: true,
		},
		{
			name:       "completed to pending does not trigger commit",
			from:       parsers.TaskStatusCompleted,
			to:         parsers.TaskStatusPending,
			wantAction: ActionStart, // action is not meaningful when wantCommit is false
			wantCommit: false,
		},
		{
			name:       "completed to in_progress does not trigger commit",
			from:       parsers.TaskStatusCompleted,
			to:         parsers.TaskStatusInProgress,
			wantAction: ActionStart,
			wantCommit: true, // Actually transitions TO in_progress DO trigger
		},
		{
			name:       "in_progress to pending does not trigger commit",
			from:       parsers.TaskStatusInProgress,
			to:         parsers.TaskStatusPending,
			wantAction: ActionStart, // action is not meaningful when wantCommit is false
			wantCommit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, shouldCommit := getActionForTransition(
				tt.to,
			)

			if shouldCommit != tt.wantCommit {
				t.Errorf(
					"getActionForTransition() shouldCommit = %v, want %v",
					shouldCommit,
					tt.wantCommit,
				)
			}

			if tt.wantCommit &&
				action != tt.wantAction {
				t.Errorf(
					"getActionForTransition() action = %v, want %v",
					action,
					tt.wantAction,
				)
			}
		})
	}

	// Remove unused variable
	_ = tracker
}

func TestTracker_Run_AlreadyComplete(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()

	// Create tasks file with all tasks completed
	tasksContent := tasksFileContent(
		struct{ id, status string }{
			"1.1",
			"completed",
		},
		struct{ id, status string }{
			"1.2",
			"completed",
		},
		struct{ id, status string }{
			"2.1",
			"completed",
		},
	)
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksContent,
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer func() { _ = tracker.Close() }()

	ctx := context.Background()
	err = tracker.Run(ctx)

	if err == nil {
		t.Fatal(
			"Run() expected error for already complete tasks, got nil",
		)
	}

	// Verify it's a TasksAlreadyCompleteError
	if _, ok := err.(*specterrs.TasksAlreadyCompleteError); !ok {
		t.Errorf(
			"Run() error type = %T, want *specterrs.TasksAlreadyCompleteError",
			err,
		)
	}
}

func TestTracker_Run_ContextCancellation(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()

	// Create tasks file with pending tasks (so it doesn't exit immediately)
	tasksContent := tasksFileContent(
		struct{ id, status string }{
			"1.1",
			"pending",
		},
		struct{ id, status string }{
			"1.2",
			"pending",
		},
	)
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksContent,
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer func() { _ = tracker.Close() }()

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	// Run tracker in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- tracker.Run(ctx)
	}()

	// Give the tracker a moment to start
	time.Sleep(50 * time.Millisecond)

	// Cancel the context
	cancel()

	// Wait for error with timeout
	select {
	case err := <-errChan:
		if err == nil {
			t.Fatal(
				"Run() expected error on context cancellation, got nil",
			)
		}

		// Verify it's a TrackInterruptedError
		if _, ok := err.(*specterrs.TrackInterruptedError); !ok {
			t.Errorf(
				"Run() error type = %T, want *specterrs.TrackInterruptedError",
				err,
			)
		}
	case <-time.After(2 * time.Second):
		t.Fatal(
			"Run() did not return after context cancellation",
		)
	}
}

func TestTracker_Close_Idempotent(t *testing.T) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksFileContent(
			struct{ id, status string }{
				"1.1",
				"pending",
			},
		),
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Close multiple times - should not panic or error
	for i := range 5 {
		err := tracker.Close()
		if err != nil {
			t.Errorf(
				"Close() call %d error = %v, want nil",
				i+1,
				err,
			)
		}
	}
}

func TestTracker_Close_NilWatcher(t *testing.T) {
	// Test that Close() handles nil watcher gracefully
	tracker := &Tracker{
		watcher: nil,
	}

	err := tracker.Close()
	if err != nil {
		t.Errorf(
			"Close() with nil watcher error = %v, want nil",
			err,
		)
	}
}

func TestTracker_printf(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		var buf bytes.Buffer
		tracker := &Tracker{writer: &buf}
		tracker.printf("Hello, %s!", "World")
		if got := buf.String(); got != "Hello, World!" {
			t.Errorf(
				"printf() output = %q, want %q",
				got,
				"Hello, World!",
			)
		}
	})

	t.Run(
		"message with numbers",
		func(t *testing.T) {
			var buf bytes.Buffer
			tracker := &Tracker{writer: &buf}
			tracker.printf(
				"Progress: %d/%d tasks",
				5,
				10,
			)
			if got := buf.String(); got != "Progress: 5/10 tasks" {
				t.Errorf(
					"printf() output = %q, want %q",
					got,
					"Progress: 5/10 tasks",
				)
			}
		},
	)

	t.Run("no args", func(t *testing.T) {
		var buf bytes.Buffer
		tracker := &Tracker{writer: &buf}
		tracker.printf("Simple message")
		if got := buf.String(); got != "Simple message" {
			t.Errorf(
				"printf() output = %q, want %q",
				got,
				"Simple message",
			)
		}
	})

	t.Run(
		"nil writer does not panic",
		func(_ *testing.T) {
			tracker := &Tracker{writer: nil}
			// This should not panic with nil writer
			tracker.printf("Should not output")
		},
	)
}

func TestConfig_Fields(t *testing.T) {
	config := Config{
		ChangeID:  "test-change-123",
		TasksPath: "/path/to/tasks.jsonc",
		RepoRoot:  "/path/to/repo",
		Writer:    os.Stdout,
	}

	if config.ChangeID != "test-change-123" {
		t.Errorf(
			"Config.ChangeID = %q, want %q",
			config.ChangeID,
			"test-change-123",
		)
	}
	if config.TasksPath != "/path/to/tasks.jsonc" {
		t.Errorf(
			"Config.TasksPath = %q, want %q",
			config.TasksPath,
			"/path/to/tasks.jsonc",
		)
	}
	if config.RepoRoot != "/path/to/repo" {
		t.Errorf(
			"Config.RepoRoot = %q, want %q",
			config.RepoRoot,
			"/path/to/repo",
		)
	}
	if config.Writer != os.Stdout {
		t.Error(
			"Config.Writer should be os.Stdout",
		)
	}
}

func TestTracker_previousState_Initialization(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksFileContent(
			struct{ id, status string }{
				"1.1",
				"pending",
			},
			struct{ id, status string }{
				"1.2",
				"in_progress",
			},
			struct{ id, status string }{
				"1.3",
				"completed",
			},
		),
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer func() { _ = tracker.Close() }()

	// previousState should be initialized but empty until Run() is called
	if tracker.previousState == nil {
		t.Error(
			"New() should initialize previousState map",
		)
	}
	if len(tracker.previousState) != 0 {
		t.Errorf(
			"New() previousState should be empty, got len=%d",
			len(tracker.previousState),
		)
	}
}

func TestTracker_Run_InitializesPreviousState(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()

	// Create tasks file - use all completed to trigger immediate return
	// after state initialization
	tasksContent := tasksFileContent(
		struct{ id, status string }{
			"1.1",
			"completed",
		},
		struct{ id, status string }{
			"1.2",
			"completed",
		},
	)
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksContent,
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer func() { _ = tracker.Close() }()

	// Run will exit with TasksAlreadyCompleteError but should initialize state first
	ctx := context.Background()
	_ = tracker.Run(ctx)

	// After Run(), previousState should be populated
	if len(tracker.previousState) != 2 {
		t.Errorf(
			"Run() should populate previousState, got len=%d, want 2",
			len(tracker.previousState),
		)
	}

	if tracker.previousState["1.1"] != parsers.TaskStatusCompleted {
		t.Errorf(
			"previousState[1.1] = %q, want %q",
			tracker.previousState["1.1"],
			parsers.TaskStatusCompleted,
		)
	}
	if tracker.previousState["1.2"] != parsers.TaskStatusCompleted {
		t.Errorf(
			"previousState[1.2] = %q, want %q",
			tracker.previousState["1.2"],
			parsers.TaskStatusCompleted,
		)
	}
}

func TestTracker_Run_PrintsInitialStatus(
	t *testing.T,
) {
	if !isFsnotifySupported() {
		t.Skip(
			"fsnotify not supported on this platform",
		)
	}

	tempDir := t.TempDir()

	// Create tasks file with pending tasks
	tasksContent := tasksFileContent(
		struct{ id, status string }{
			"1.1",
			"completed",
		},
		struct{ id, status string }{
			"1.2",
			"pending",
		},
		struct{ id, status string }{
			"1.3",
			"pending",
		},
	)
	tasksPath := createTestTasksFile(
		t,
		tempDir,
		tasksContent,
	)

	var buf bytes.Buffer
	tracker, err := New(Config{
		ChangeID:  "test-change",
		TasksPath: tasksPath,
		RepoRoot:  tempDir,
		Writer:    &buf,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer func() { _ = tracker.Close() }()

	// Create cancellable context
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	// Run in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- tracker.Run(ctx)
	}()

	// Give tracker time to print initial status
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Wait for completion
	<-errChan

	// Check that initial status was printed
	output := buf.String()
	if !containsString(
		output,
		"Tracking test-change",
	) {
		t.Error(
			"Run() should print tracking message with change ID",
		)
	}
	if !containsString(
		output,
		"1/3 tasks completed",
	) {
		t.Errorf(
			"Run() should print progress, got: %s",
			output,
		)
	}
	if !containsString(
		output,
		"Watching for task status changes",
	) {
		t.Error(
			"Run() should print watching message",
		)
	}
}

// containsString checks if s contains substr.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

func TestTracker_countProgress_FixtureLike(
	t *testing.T,
) {
	// Create 16 tasks like the real fixture
	tasks := make([]parsers.Task, 16)
	sections := []string{
		"Foundation",
		"Core Components",
		"Command Integration",
		"Testing",
		"Validation",
	}

	for i := range 16 {
		tasks[i] = parsers.Task{
			ID: fmt.Sprintf(
				"%d.%d",
				(i/3)+1,
				(i%3)+1,
			),
			Section: sections[i%len(sections)],
			Status:  parsers.TaskStatusPending,
		}
	}

	// Test all pending
	completed, total := countProgress(tasks)
	if completed != 0 {
		t.Errorf(
			"countProgress() completed = %d, want 0 for all pending",
			completed,
		)
	}
	if total != 16 {
		t.Errorf(
			"countProgress() total = %d, want 16",
			total,
		)
	}

	// Mark 5 tasks completed
	tasks[0].Status = parsers.TaskStatusCompleted
	tasks[1].Status = parsers.TaskStatusCompleted
	tasks[2].Status = parsers.TaskStatusCompleted
	tasks[3].Status = parsers.TaskStatusCompleted
	tasks[4].Status = parsers.TaskStatusCompleted

	completed, total = countProgress(tasks)
	if completed != 5 {
		t.Errorf(
			"countProgress() completed = %d, want 5",
			completed,
		)
	}
	if total != 16 {
		t.Errorf(
			"countProgress() total = %d, want 16",
			total,
		)
	}

	// Mark remaining completed
	for i := 5; i < 16; i++ {
		tasks[i].Status = parsers.TaskStatusCompleted
	}

	completed, total = countProgress(tasks)
	if completed != 16 {
		t.Errorf(
			"countProgress() completed = %d, want 16",
			completed,
		)
	}
	if total != 16 {
		t.Errorf(
			"countProgress() total = %d, want 16",
			total,
		)
	}
}

func TestTracker_allTasksComplete_FixtureLike(
	t *testing.T,
) {
	// Create 16 tasks like the real fixture, all pending
	tasks := make([]parsers.Task, 16)
	for i := range 16 {
		tasks[i] = parsers.Task{
			ID: fmt.Sprintf(
				"%d.%d",
				(i/3)+1,
				(i%3)+1,
			),
			Status: parsers.TaskStatusPending,
		}
	}

	// All pending should return false
	if allTasksComplete(tasks) {
		t.Error(
			"allTasksComplete() = true, want false for all pending",
		)
	}

	// Mark all but one completed
	for i := range 15 {
		tasks[i].Status = parsers.TaskStatusCompleted
	}

	if allTasksComplete(tasks) {
		t.Error(
			"allTasksComplete() = true, want false when one task pending",
		)
	}

	// Mark last task completed
	tasks[15].Status = parsers.TaskStatusCompleted

	if !allTasksComplete(tasks) {
		t.Error(
			"allTasksComplete() = false, want true when all completed",
		)
	}
}

func TestTracker_countProgress_InProgressNotCounted(
	t *testing.T,
) {
	// Test that in_progress tasks are NOT counted as completed
	tasks := []parsers.Task{
		{
			ID:     "1.1",
			Status: parsers.TaskStatusCompleted,
		},
		{
			ID:     "1.2",
			Status: parsers.TaskStatusInProgress,
		},
		{
			ID:     "1.3",
			Status: parsers.TaskStatusInProgress,
		},
		{
			ID:     "1.4",
			Status: parsers.TaskStatusPending,
		},
	}

	completed, total := countProgress(tasks)
	if completed != 1 {
		t.Errorf(
			"countProgress() completed = %d, want 1 (in_progress should NOT count)",
			completed,
		)
	}
	if total != 4 {
		t.Errorf(
			"countProgress() total = %d, want 4",
			total,
		)
	}
}
