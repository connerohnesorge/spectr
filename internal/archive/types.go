package archive

// SpecUpdate represents a spec file to update during archive
type SpecUpdate struct {
	Source string // Path to delta spec in change
	Target string // Path to main spec in spectr/specs
	Exists bool   // Does target spec already exist?
}

// OperationCounts tracks the number of each delta operation applied
type OperationCounts struct {
	Added    int
	Modified int
	Removed  int
	Renamed  int
}

// Total adds up all counts for delta operations.
func (oc *OperationCounts) Total() int {
	return oc.Added + oc.Modified + oc.Removed + oc.Renamed
}

// ArchiveResult contains the results of an archive operation.
type ArchiveResult struct {
	// ArchivePath is the relative path to the archived change
	// (e.g., "spectr/changes/archive/2025-12-02-change-id/")
	ArchivePath string

	// Counts tracks the number of each delta operation applied
	Counts OperationCounts

	// Capabilities lists the updated capability names (e.g., ["auth", "api"])
	Capabilities []string
}
