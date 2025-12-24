package providers

// InitResult contains the files created or modified by an initializer.
// This provides explicit tracking of what changed during initialization,
// which is useful for:
//   - Reporting what was done to the user
//   - Testing initialization behavior
//   - Working in non-git projects (where git status isn't available)
//   - Aggregating changes across multiple initializers
type InitResult struct {
	// CreatedFiles contains paths of files that were created
	// (did not exist before). Paths are relative to the filesystem
	// root (projectFs or globalFs).
	CreatedFiles []string

	// UpdatedFiles contains paths of files that were modified
	// (already existed). Paths are relative to the filesystem root
	// (projectFs or globalFs).
	UpdatedFiles []string
}

// IsEmpty returns true if this result contains no file changes.
// Used to determine if initialization actually did any work.
func (r InitResult) IsEmpty() bool {
	return len(r.CreatedFiles) == 0 &&
		len(r.UpdatedFiles) == 0
}

// Merge combines this result with another, concatenating file lists.
// Used to aggregate results from multiple initializers.
//
// Example:
//
//	result1 := InitResult{CreatedFiles: []string{"a.txt"}}
//	result2 := InitResult{UpdatedFiles: []string{"b.txt"}}
//	combined := result1.Merge(result2)
//	// combined.CreatedFiles = ["a.txt"]
//	// combined.UpdatedFiles = ["b.txt"]
func (r InitResult) Merge(
	other InitResult,
) InitResult {
	// Create new slices to avoid modifying original slices
	createdLen := len(
		r.CreatedFiles,
	) + len(
		other.CreatedFiles,
	)
	updatedLen := len(
		r.UpdatedFiles,
	) + len(
		other.UpdatedFiles,
	)

	created := make([]string, 0, createdLen)
	created = append(created, r.CreatedFiles...)
	created = append(
		created,
		other.CreatedFiles...)

	updated := make([]string, 0, updatedLen)
	updated = append(updated, r.UpdatedFiles...)
	updated = append(
		updated,
		other.UpdatedFiles...)

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}
}

// TotalFiles returns the total number of files affected (created + updated).
// Useful for summary reporting.
func (r InitResult) TotalFiles() int {
	return len(
		r.CreatedFiles,
	) + len(
		r.UpdatedFiles,
	)
}
