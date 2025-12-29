package providers

// InitResult contains the files created or modified by an initializer.
// This provides explicit change tracking for each initialization step.
type InitResult struct {
	// CreatedFiles is the list of files that were newly created.
	// Paths should be relative to the filesystem root (projectFs or homeFs).
	CreatedFiles []string

	// UpdatedFiles is the list of files that already existed and were modified.
	// Paths should be relative to the filesystem root (projectFs or homeFs).
	UpdatedFiles []string
}

// ExecutionResult aggregates results from all initializers.
// This is the final result returned after running all initializers.
//
// Note: Error is returned separately from the executor function,
// not stored in this struct. This allows partial results to be
// returned even when an error occurs (fail-fast behavior).
type ExecutionResult struct {
	// CreatedFiles is the list of all files created
	// across all initializers.
	CreatedFiles []string

	// UpdatedFiles is the list of all files updated
	// across all initializers.
	UpdatedFiles []string
}
