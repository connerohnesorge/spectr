package specterrs

import "fmt"

// NoTasksFileError indicates tasks.jsonc was not found for a change.
type NoTasksFileError struct {
	ChangeID string
}

func (e *NoTasksFileError) Error() string {
	return fmt.Sprintf("tasks file not found for change %q", e.ChangeID)
}

// TasksAlreadyCompleteError indicates all tasks are already completed.
type TasksAlreadyCompleteError struct {
	ChangeID string
}

func (e *TasksAlreadyCompleteError) Error() string {
	return fmt.Sprintf("all tasks already completed for change %q", e.ChangeID)
}

// TrackInterruptedError indicates tracking was stopped by user interrupt.
type TrackInterruptedError struct{}

func (*TrackInterruptedError) Error() string {
	return "tracking stopped by user interrupt"
}

// GitCommitError indicates a git commit operation failed.
type GitCommitError struct {
	Err error
}

func (*GitCommitError) Error() string {
	return "git commit failed"
}

func (e *GitCommitError) Unwrap() error {
	return e.Err
}
