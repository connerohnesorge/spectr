package specterrs

import "fmt"

// MissingChangeIDError indicates that a change ID is required but was not
// provided, and interactive mode is disabled.
type MissingChangeIDError struct{}

func (*MissingChangeIDError) Error() string {
	return "usage: spectr accept <change-id> [flags]\n" +
		"       spectr accept <change-id> --dry-run"
}

// NoValidTasksError indicates that tasks.md has content but no valid tasks
// were found during parsing.
type NoValidTasksError struct {
	TasksMdPath string
	FileSize    int64
}

const (
	expectedTaskFormat = `- [ ] N.N Task, - [ ] N. Task, or - [ ] Task`
)

// Error implements the error interface on NoValidTasksError.
func (e *NoValidTasksError) Error() string {
	return fmt.Sprintf(
		"tasks.md has content (%d bytes) but no valid tasks found; "+
			"expected format: "+expectedTaskFormat,
		e.FileSize,
	)
}
