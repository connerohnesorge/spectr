// Package parsers provides functions for parsing delta specifications,
// requirements, and other structured spec documents.
package parsers

import (
	"os"

	"github.com/connerohnesorge/spectr/internal/mdparser"
)

// DeltaPlan represents all delta operations for a spec
type DeltaPlan struct {
	Added    []RequirementBlock
	Modified []RequirementBlock
	Removed  []string // Just requirement names
	Renamed  []RenameOp
}

// RenameOp represents a requirement rename operation
type RenameOp struct {
	From string
	To   string
}

// ParseDeltaSpec parses a delta spec file and extracts operations
// Returns a DeltaPlan with ADDED, MODIFIED, REMOVED, and RENAMED reqs
func ParseDeltaSpec(filePath string) (*DeltaPlan, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return nil, err
	}

	plan := &DeltaPlan{
		Added:    make([]RequirementBlock, 0),
		Modified: make([]RequirementBlock, 0),
		Removed:  make([]string, 0),
		Renamed:  make([]RenameOp, 0),
	}

	// Extract delta sections using mdparser extractor
	deltaSections, err := ExtractDeltaSections(doc)
	if err != nil {
		return nil, err
	}

	// Populate plan from extracted sections
	if added, ok := deltaSections["ADDED"]; ok {
		plan.Added = added
	}
	if modified, ok := deltaSections["MODIFIED"]; ok {
		plan.Modified = modified
	}
	if removed, ok := deltaSections["REMOVED"]; ok {
		// Extract just the names from removed requirements
		for _, req := range removed {
			plan.Removed = append(plan.Removed, req.Name)
		}
	}

	// Extract renamed requirements
	renamed, err := ExtractRenamedRequirements(doc)
	if err != nil {
		return nil, err
	}
	for _, r := range renamed {
		plan.Renamed = append(plan.Renamed, RenameOp(r))
	}

	return plan, nil
}

// HasDeltas returns true if the DeltaPlan has at least one operation
func (dp *DeltaPlan) HasDeltas() bool {
	hasAdded := len(dp.Added) > 0
	hasModified := len(dp.Modified) > 0
	hasRemoved := len(dp.Removed) > 0
	hasRenamed := len(dp.Renamed) > 0

	return hasAdded || hasModified || hasRemoved || hasRenamed
}

// CountOperations returns the total number of delta operations
func (dp *DeltaPlan) CountOperations() int {
	return len(dp.Added) + len(dp.Modified) + len(dp.Removed) + len(dp.Renamed)
}
