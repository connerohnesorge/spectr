package providers

// InitResult contains the files created or modified by an
// initializer.
// This provides explicit change tracking without relying on git or
// filesystem state.
type InitResult struct {
	// CreatedFiles lists the files that were created by this initializer.
	// Paths are relative to the filesystem root (projectFs or globalFs).
	CreatedFiles []string

	// UpdatedFiles lists the files that were updated by this initializer.
	// Paths are relative to the filesystem root (projectFs or globalFs).
	UpdatedFiles []string
}

// Merge combines multiple InitResults into a single result.
// Used to aggregate results from multiple initializers.
func (r InitResult) Merge(other InitResult) InitResult {
	return InitResult{
		CreatedFiles: append(r.CreatedFiles, other.CreatedFiles...),
		UpdatedFiles: append(r.UpdatedFiles, other.UpdatedFiles...),
	}
}

// HasChanges returns true if any files were created or updated.
func (r InitResult) HasChanges() bool {
	return len(r.CreatedFiles) > 0 || len(r.UpdatedFiles) > 0
}
