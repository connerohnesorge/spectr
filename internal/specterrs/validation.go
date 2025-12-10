package specterrs

import "fmt"

// ValidationFailedError indicates validation failed for a single item.
type ValidationFailedError struct {
	ItemCount    int
	ErrorCount   int
	WarningCount int
}

func (*ValidationFailedError) Error() string {
	return "validation failed"
}

// MultiValidationFailedError indicates validation failed for multiple items.
type MultiValidationFailedError struct {
	ItemCount int
}

func (*MultiValidationFailedError) Error() string {
	return "validation failed for one or more items"
}

// DeltaSpecParseError indicates a delta spec failed to parse.
type DeltaSpecParseError struct {
	SpecPath string
	Line     int
	Err      error
}

func (e *DeltaSpecParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf(
			"failed to parse delta spec %s at line %d: %v",
			e.SpecPath,
			e.Line,
			e.Err,
		)
	}

	return fmt.Sprintf("failed to parse delta spec %s: %v", e.SpecPath, e.Err)
}

func (e *DeltaSpecParseError) Unwrap() error {
	return e.Err
}
