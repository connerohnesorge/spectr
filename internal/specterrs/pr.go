package specterrs

import "fmt"

// UnknownPlatformError indicates an unknown git platform.
type UnknownPlatformError struct {
	Platform string
	RepoURL  string
}

func (*UnknownPlatformError) Error() string {
	return "unknown platform; please create PR manually"
}

// PRPrerequisiteError indicates a PR prerequisite check failed.
type PRPrerequisiteError struct {
	Check   string
	Details string
	Err     error
}

func (e *PRPrerequisiteError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf(
			"PR prerequisite failed (%s): %s",
			e.Check,
			e.Details,
		)
	}

	return fmt.Sprintf(
		"PR prerequisite failed (%s)",
		e.Check,
	)
}

func (e *PRPrerequisiteError) Unwrap() error {
	return e.Err
}
