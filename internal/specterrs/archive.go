package specterrs

import "fmt"

// UserCancelledError indicates the user cancelled a selection.
type UserCancelledError struct {
	Operation string
}

func (*UserCancelledError) Error() string {
	return "user cancelled selection"
}

// ArchiveCancelledError indicates an archive operation was cancelled.
type ArchiveCancelledError struct {
	Reason string
}

func (*ArchiveCancelledError) Error() string {
	return "archive cancelled"
}

// ValidationRequiredError indicates validation errors must be fixed
// before proceeding.
type ValidationRequiredError struct {
	Operation string
}

func (e *ValidationRequiredError) Error() string {
	return fmt.Sprintf(
		"validation errors must be fixed before %s",
		e.Operation,
	)
}

// DeltaConflictError indicates a requirement appears in multiple sections.
type DeltaConflictError struct {
	Section1        string
	Section2        string
	RequirementName string
}

func (e *DeltaConflictError) Error() string {
	return fmt.Sprintf(
		"requirement %q appears in both %s and %s sections",
		e.RequirementName,
		e.Section1,
		e.Section2,
	)
}

// DuplicateRequirementError indicates a duplicate requirement in a section.
type DuplicateRequirementError struct {
	RequirementName string
	SectionName     string
}

func (e *DuplicateRequirementError) Error() string {
	return fmt.Sprintf(
		"duplicate requirement %q in %s section",
		e.RequirementName,
		e.SectionName,
	)
}

// IncompleteTasksError indicates an archive was cancelled due to
// incomplete tasks.
type IncompleteTasksError struct {
	Total     int
	Completed int
}

func (*IncompleteTasksError) Error() string {
	return "archive cancelled due to incomplete tasks"
}
