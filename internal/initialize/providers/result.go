package providers

// InitResult contains the files created or modified by an initializer.
type InitResult struct {
	CreatedFiles []string // files created by this initializer
	UpdatedFiles []string // files updated by this initializer
}

// ExecutionResult aggregates results from all initializers.
// Note: Error is returned separately from runInitializers(), not stored
// in this struct.
type ExecutionResult struct {
	CreatedFiles []string // all files created across all initializers
	UpdatedFiles []string // all files updated across all initializers
}

// AggregateResults combines multiple InitResult values into a single
// ExecutionResult.
func AggregateResults(results []InitResult) ExecutionResult {
	var created, updated []string
	for _, r := range results {
		created = append(created, r.CreatedFiles...)
		updated = append(updated, r.UpdatedFiles...)
	}

	return ExecutionResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}
}
