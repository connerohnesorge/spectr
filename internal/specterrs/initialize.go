package specterrs

import (
	"fmt"
	"strings"
)

// EmptyPathError indicates a path cannot be empty.
type EmptyPathError struct {
	Operation string
}

func (*EmptyPathError) Error() string {
	return "path cannot be empty"
}

// WizardModelCastError indicates failed to cast the final model to WizardModel.
type WizardModelCastError struct {
	ActualType string
}

func (*WizardModelCastError) Error() string {
	return "failed to cast final model to WizardModel"
}

// InitializationCompletedWithErrorsError indicates initialization
// completed with errors.
type InitializationCompletedWithErrorsError struct {
	ErrorCount int
	Errors     []error
}

func (e *InitializationCompletedWithErrorsError) Error() string {
	if e.ErrorCount == 1 {
		return "initialization completed with 1 error"
	}

	return fmt.Sprintf("initialization completed with %d errors", e.ErrorCount)
}

func (e *InitializationCompletedWithErrorsError) Unwrap() []error {
	return e.Errors
}

// ErrorMessages returns a formatted string of all error messages.
func (e *InitializationCompletedWithErrorsError) ErrorMessages() string {
	msgs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}

	return strings.Join(msgs, "\n")
}
