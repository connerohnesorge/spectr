package initializers

// InitResult contains the files created or modified by an initializer.
// This provides explicit tracking of what changed during initialization,
// which is useful for reporting to users and for testing.
type InitResult struct {
	// CreatedFiles contains paths of files that were newly created
	CreatedFiles []string

	// UpdatedFiles contains paths of files that were modified
	// (existed before but were updated)
	UpdatedFiles []string
}

// Merge combines this result with another result, returning a new result
// containing the union of both file lists.
//
// Example:
//
//	r1 := InitResult{CreatedFiles: []string{"a.txt"}}
//	r2 := InitResult{CreatedFiles: []string{"b.txt"}}
//	merged := r1.Merge(r2)
//	// merged.CreatedFiles == []string{"a.txt", "b.txt"}
func (r InitResult) Merge(
	other InitResult,
) InitResult {
	created := make(
		[]string,
		0,
		len(
			r.CreatedFiles,
		)+len(
			other.CreatedFiles,
		),
	)
	created = append(created, r.CreatedFiles...)
	created = append(
		created,
		other.CreatedFiles...)

	updated := make(
		[]string,
		0,
		len(
			r.UpdatedFiles,
		)+len(
			other.UpdatedFiles,
		),
	)
	updated = append(updated, r.UpdatedFiles...)
	updated = append(
		updated,
		other.UpdatedFiles...)

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}
}

// IsEmpty returns true if no files were created or updated.
func (r InitResult) IsEmpty() bool {
	return len(r.CreatedFiles) == 0 &&
		len(r.UpdatedFiles) == 0
}

// TotalFiles returns the total number of files affected (created + updated).
func (r InitResult) TotalFiles() int {
	return len(
		r.CreatedFiles,
	) + len(
		r.UpdatedFiles,
	)
}

// AddCreated adds a file path to the CreatedFiles list.
// Returns a new InitResult with the file added.
func (r InitResult) AddCreated(
	path string,
) InitResult {
	created := make(
		[]string,
		0,
		len(r.CreatedFiles)+1,
	)
	created = append(created, r.CreatedFiles...)
	created = append(created, path)

	updated := make([]string, len(r.UpdatedFiles))
	copy(updated, r.UpdatedFiles)

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}
}

// AddUpdated adds a file path to the UpdatedFiles list.
// Returns a new InitResult with the file added.
func (r InitResult) AddUpdated(
	path string,
) InitResult {
	created := make([]string, len(r.CreatedFiles))
	copy(created, r.CreatedFiles)

	updated := make(
		[]string,
		0,
		len(r.UpdatedFiles)+1,
	)
	updated = append(updated, r.UpdatedFiles...)
	updated = append(updated, path)

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}
}
