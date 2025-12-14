// Package parsers provides functions for parsing delta specifications,
// requirements, and other structured spec documents.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package parsers

import (
	"bufio"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
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

// ParseDeltaSpec parses a delta spec file and extracts operations.
// Returns a DeltaPlan with ADDED, MODIFIED, REMOVED, and RENAMED reqs.
// Uses AST-based parsing via markdown.ParseDocument for accurate extraction.
func ParseDeltaSpec(
	filePath string,
) (*DeltaPlan, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		return nil, err
	}

	plan := &DeltaPlan{
		Added:    extractRequirementsFromDeltaSection(doc, "ADDED"),
		Modified: extractRequirementsFromDeltaSection(doc, "MODIFIED"),
		Removed:  extractRemovedFromDoc(doc),
		Renamed:  extractRenamedFromDoc(doc),
	}

	return plan, nil
}

// extractRequirementsFromDeltaSection extracts requirements from a delta section.
func extractRequirementsFromDeltaSection(
	doc *markdown.Document,
	deltaType string,
) []RequirementBlock {
	section := doc.GetDeltaSection(deltaType)
	if section == nil || section.Content == "" {
		return nil
	}

	// Parse the section content to extract requirements
	sectionDoc, err := markdown.ParseDocument([]byte(section.Content))
	if err != nil {
		return nil
	}

	names := sectionDoc.GetRequirementNames()
	requirements := make([]RequirementBlock, 0, len(names))

	for _, name := range names {
		req := sectionDoc.GetRequirement(name)
		if req == nil {
			continue
		}

		headerLine := "### Requirement: " + req.Name
		raw := headerLine + "\n" + req.Content

		requirements = append(requirements, RequirementBlock{
			HeaderLine: headerLine,
			Name:       req.Name,
			Raw:        raw,
		})
	}

	return requirements
}

// extractRemovedFromDoc extracts requirement names from REMOVED section.
func extractRemovedFromDoc(doc *markdown.Document) []string {
	section := doc.GetDeltaSection("REMOVED")
	if section == nil || section.Content == "" {
		return nil
	}

	// Parse the section content to extract requirement names
	sectionDoc, err := markdown.ParseDocument([]byte(section.Content))
	if err != nil {
		return nil
	}

	return sectionDoc.GetRequirementNames()
}

// extractRenamedFromDoc extracts FROM/TO pairs from RENAMED section.
func extractRenamedFromDoc(doc *markdown.Document) []RenameOp {
	section := doc.GetDeltaSection("RENAMED")
	if section == nil || section.Content == "" {
		return nil
	}

	// Parse FROM/TO pairs
	// Expected format:
	// - FROM: `### Requirement: Old Name`
	// - TO: `### Requirement: New Name`
	var renamed []RenameOp
	var currentFrom string

	scanner := bufio.NewScanner(
		strings.NewReader(section.Content),
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for FROM line (backtick format)
		if name, ok := markdown.MatchRenamedFrom(line); ok {
			currentFrom = strings.TrimSpace(name)

			continue
		}

		// Check for TO line (backtick format)
		if name, ok := markdown.MatchRenamedTo(line); ok {
			if currentFrom == "" {
				continue
			}
			renamed = append(renamed, RenameOp{
				From: currentFrom,
				To:   strings.TrimSpace(name),
			})
			currentFrom = ""
		}
	}

	return renamed
}

// HasDeltas returns true if the DeltaPlan has at least one operation
func (dp *DeltaPlan) HasDeltas() bool {
	hasAdded := len(dp.Added) > 0
	hasModified := len(dp.Modified) > 0
	hasRemoved := len(dp.Removed) > 0
	hasRenamed := len(dp.Renamed) > 0

	return hasAdded || hasModified ||
		hasRemoved ||
		hasRenamed
}

// CountOperations returns the total number of delta operations
func (dp *DeltaPlan) CountOperations() int {
	return len(
		dp.Added,
	) + len(
		dp.Modified,
	) + len(
		dp.Removed,
	) + len(
		dp.Renamed,
	)
}
