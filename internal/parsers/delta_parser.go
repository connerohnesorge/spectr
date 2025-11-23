// Package parsers provides functions for parsing delta specifications,
// requirements, and other structured spec documents.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package parsers

import (
	"fmt"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parser"
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

	// Parse the document using the new parser
	doc, err := parser.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse delta spec: %w", err)
	}

	// Extract delta operations
	deltas, err := parser.ExtractDeltas(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract deltas: %w", err)
	}

	// Convert to DeltaPlan format
	return convertToDeltaPlan(deltas), nil
}

// convertToDeltaPlan converts parser.DeltaSpec to DeltaPlan format
func convertToDeltaPlan(deltas *parser.DeltaSpec) *DeltaPlan {
	plan := &DeltaPlan{
		Added:    convertRequirements(deltas.Added),
		Modified: convertRequirements(deltas.Modified),
		Removed:  extractRemovedNames(deltas.Removed),
		Renamed:  convertRenamedOps(deltas.Renamed),
	}

	return plan
}

// convertRequirements converts parser.Requirement to RequirementBlock
func convertRequirements(reqs []parser.Requirement) []RequirementBlock {
	blocks := make([]RequirementBlock, 0, len(reqs))
	for _, req := range reqs {
		blocks = append(blocks, RequirementBlock{
			HeaderLine: "### Requirement: " + req.Name,
			Name:       req.Name,
			Raw:        buildRawContent(req),
		})
	}

	return blocks
}

// buildRawContent reconstructs the raw markdown from a parser.Requirement
func buildRawContent(req parser.Requirement) string {
	var parts []string

	// Add header line
	parts = append(parts, "### Requirement: "+req.Name)

	// Add content (which includes scenarios)
	if req.Content != "" {
		parts = append(parts, req.Content)
	}

	return strings.Join(parts, "\n") + "\n"
}

// extractRemovedNames extracts requirement names from removed requirements
func extractRemovedNames(removed []parser.Requirement) []string {
	names := make([]string, 0, len(removed))
	for _, req := range removed {
		names = append(names, req.Name)
	}

	return names
}

// convertRenamedOps converts parser.RenamedRequirement to RenameOp
func convertRenamedOps(renamed []parser.RenamedRequirement) []RenameOp {
	ops := make([]RenameOp, 0, len(renamed))
	for _, r := range renamed {
		ops = append(ops, RenameOp{
			From: r.From,
			To:   r.To,
		})
	}

	return ops
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
