package domain

// ExecutionResult contains results from initialization.
// Note: Error is returned separately, not stored in this struct.
type ExecutionResult struct {
	CreatedFiles []string // All files created
	UpdatedFiles []string // All files updated
}
