//nolint:revive // cognitive-complexity and nesting are acceptable for event loops
package track

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// Tracker coordinates file watching and git commits for task status changes.
// It monitors the tasks.jsonc file and creates commits when tasks transition
// to "in_progress" or "completed" status.
type Tracker struct {
	changeID      string
	tasksPath     string
	repoRoot      string
	watcher       *Watcher
	committer     *Committer
	writer        io.Writer
	previousState map[string]parsers.TaskStatusValue
}

// Config holds configuration for creating a new Tracker.
type Config struct {
	// ChangeID is the identifier for the change being tracked.
	ChangeID string
	// TasksPath is the absolute path to the tasks.jsonc file.
	TasksPath string
	// RepoRoot is the root directory of the git repository.
	RepoRoot string
	// Writer is used for progress output (e.g., os.Stdout).
	Writer io.Writer
	// IncludeBinaries controls whether binary files are included in commits.
	// When false (default), binary files are excluded from automated commits.
	IncludeBinaries bool
}

// New creates a new Tracker with the specified configuration.
// Returns an error if the tasks file cannot be watched.
func New(config Config) (*Tracker, error) {
	watcher, err := NewWatcher(config.TasksPath)
	if err != nil {
		return nil, err
	}

	committer := NewCommitter(
		config.ChangeID,
		config.RepoRoot,
		config.IncludeBinaries,
	)

	return &Tracker{
		changeID:  config.ChangeID,
		tasksPath: config.TasksPath,
		repoRoot:  config.RepoRoot,
		watcher:   watcher,
		committer: committer,
		writer:    config.Writer,
		previousState: make(
			map[string]parsers.TaskStatusValue,
		),
	}, nil
}

// Run starts the tracking event loop. It watches for task status changes
// and creates commits when tasks transition to "in_progress" or "completed".
//
// The loop exits when:
//   - All tasks are completed
//   - The context is cancelled (e.g., Ctrl+C)
//   - A git commit operation fails
//
// Returns:
//   - TasksAlreadyCompleteError if all tasks are already complete
//   - TrackInterruptedError if cancelled via context
//   - GitCommitError if a git operation fails
func (t *Tracker) Run(ctx context.Context) error {
	tasksFile, err := parsers.ReadTasksJson(
		t.tasksPath,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to read tasks file: %w",
			err,
		)
	}

	for _, task := range tasksFile.Tasks {
		t.previousState[task.ID] = task.Status
	}

	if allTasksComplete(tasksFile.Tasks) {
		return &specterrs.TasksAlreadyCompleteError{
			ChangeID: t.changeID,
		}
	}

	completed, total := countProgress(
		tasksFile.Tasks,
	)
	t.printf(
		"Tracking %s: %d/%d tasks completed\n",
		t.changeID, completed, total,
	)
	t.printf(
		"Watching for task status changes... (Ctrl+C to stop)\n\n",
	)

	return t.eventLoop(ctx)
}

// eventLoop is the main event loop that processes file changes.
func (t *Tracker) eventLoop(
	ctx context.Context,
) error {
	for {
		select {
		case <-ctx.Done():
			return &specterrs.TrackInterruptedError{}

		case <-t.watcher.Events():
			if err := t.handleFileChange(); err != nil {
				if _, ok := err.(*specterrs.GitCommitError); ok {
					return err
				}
				t.printf("Warning: %v\n", err)
			}

			if t.checkAllComplete() {
				return nil
			}

		case err := <-t.watcher.Errors():
			t.printf(
				"Warning: watcher error: %v\n",
				err,
			)
		}
	}
}

// checkAllComplete checks if all tasks are complete and prints status.
func (t *Tracker) checkAllComplete() bool {
	tasksFile, err := parsers.ReadTasksJson(
		t.tasksPath,
	)
	if err != nil {
		return false
	}

	if !allTasksComplete(tasksFile.Tasks) {
		return false
	}

	completed, total := countProgress(
		tasksFile.Tasks,
	)
	t.printf(
		"\nAll tasks completed! (%d/%d)\n",
		completed,
		total,
	)

	return true
}

// Close stops the tracker and releases resources.
func (t *Tracker) Close() error {
	if t.watcher != nil {
		return t.watcher.Close()
	}

	return nil
}

// handleFileChange processes a file change event by reloading tasks
// and committing for any status transitions.
func (t *Tracker) handleFileChange() error {
	tasksFile, err := parsers.ReadTasksJson(
		t.tasksPath,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to read tasks file: %w",
			err,
		)
	}

	for _, task := range tasksFile.Tasks {
		if err := t.processTaskTransition(task); err != nil {
			return err
		}
	}

	return nil
}

// processTaskTransition checks and commits a single task's transition.
func (t *Tracker) processTaskTransition(
	task parsers.Task,
) error {
	prevStatus, exists := t.previousState[task.ID]
	if !exists {
		t.previousState[task.ID] = task.Status

		return nil
	}

	if prevStatus == task.Status {
		return nil
	}

	action, shouldCommit := getActionForTransition(
		task.Status,
	)
	if shouldCommit {
		if err := t.commitTransition(task.ID, action); err != nil {
			return err
		}
	}

	t.previousState[task.ID] = task.Status

	return nil
}

// getActionForTransition determines the commit action for a status transition.
// Returns the action and whether a commit should be created.
func getActionForTransition(
	to parsers.TaskStatusValue,
) (Action, bool) {
	switch to {
	case parsers.TaskStatusInProgress:
		return ActionStart, true
	case parsers.TaskStatusCompleted:
		return ActionComplete, true
	case parsers.TaskStatusPending:
		return ActionStart, false
	}

	return ActionStart, false
}

// commitTransition creates a commit for the task status transition.
func (t *Tracker) commitTransition(
	taskID string,
	action Action,
) error {
	result, err := t.committer.Commit(
		taskID,
		action,
	)
	if err != nil {
		t.printf(
			"Error: failed to commit for task %s: %v\n",
			taskID,
			err,
		)

		return err
	}

	if result.NoFiles {
		t.printf(
			"  Task %s: %s (no files to commit)\n",
			taskID,
			action.String(),
		)
	} else {
		hash := result.CommitHash[:7]
		t.printf("  Task %s: %s [%s]\n", taskID, action.String(), hash)
	}

	// Display warning about skipped binary files
	if len(result.SkippedBinaries) > 0 {
		t.printf(
			"  Warning: Skipped binary files: %s\n",
			strings.Join(
				result.SkippedBinaries,
				", ",
			),
		)
	}

	return nil
}

// allTasksComplete checks if all tasks have completed status.
func allTasksComplete(tasks []parsers.Task) bool {
	if len(tasks) == 0 {
		return true
	}

	for _, task := range tasks {
		if task.Status != parsers.TaskStatusCompleted {
			return false
		}
	}

	return true
}

// countProgress returns the number of completed tasks and total tasks.
func countProgress(
	tasks []parsers.Task,
) (completed, total int) {
	total = len(tasks)
	for _, task := range tasks {
		if task.Status == parsers.TaskStatusCompleted {
			completed++
		}
	}

	return completed, total
}

// printf writes formatted output to the tracker's writer.
func (t *Tracker) printf(
	format string,
	args ...any,
) {
	if t.writer != nil {
		_, _ = fmt.Fprintf(
			t.writer,
			format,
			args...)
	}
}
