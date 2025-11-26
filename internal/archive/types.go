package archive

import "errors"

// ErrUserCancelled is returned when the user cancels an interactive selection.
var ErrUserCancelled = errors.New("user cancelled selection")

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

// Add increments the total operation count
func (oc *OperationCounts) Total() int {
	return oc.Added + oc.Modified + oc.Removed + oc.Renamed
}
