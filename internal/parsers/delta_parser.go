// Package parsers provides functions for parsing delta specifications,
// requirements, and other structured spec documents.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
	bf "github.com/russross/blackfriday/v2"
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

	plan := &DeltaPlan{
		Added:    make([]RequirementBlock, 0),
		Modified: make([]RequirementBlock, 0),
		Removed:  make([]string, 0),
		Renamed:  make([]RenameOp, 0),
	}

	// Parse using markdown package
	node := markdown.Parse(content)

	// Parse each section
	plan.Added = parseDeltaSection(node, "ADDED")
	plan.Modified = parseDeltaSection(node, "MODIFIED")
	plan.Removed = parseRemovedSection(node)
	plan.Renamed = parseRenamedSection(string(content))

	return plan, nil
}

// parseDeltaSection extracts requirements from a delta section
func parseDeltaSection(node *bf.Node, sectionType string) []RequirementBlock {
	sectionContent := markdown.ExtractDeltaSectionContent(node, sectionType)
	if sectionContent == "" {
		return nil
	}

	return parseRequirementsFromSection(sectionContent)
}

// parseRequirementsFromSection parses requirement blocks from content
func parseRequirementsFromSection(sectionContent string) []RequirementBlock {
	mdRequirements := markdown.ExtractRequirementsFromContent(sectionContent)

	// Convert markdown.RequirementBlock to parsers.RequirementBlock
	requirements := make([]RequirementBlock, len(mdRequirements))
	for i, mdReq := range mdRequirements {
		requirements[i] = RequirementBlock{
			HeaderLine: mdReq.HeaderLine,
			Name:       mdReq.Name,
			Raw:        mdReq.Raw,
		}
	}

	return requirements
}

// parseRemovedSection extracts requirement names from REMOVED section
func parseRemovedSection(node *bf.Node) []string {
	sectionContent := markdown.ExtractDeltaSectionContent(node, "REMOVED")
	if sectionContent == "" {
		return nil
	}

	// Extract requirement names from the section content
	mdRequirements := markdown.ExtractRequirementsFromContent(sectionContent)

	removed := make([]string, len(mdRequirements))
	for i, mdReq := range mdRequirements {
		removed[i] = mdReq.Name
	}

	return removed
}

// parseRenamedSection extracts FROM/TO pairs from RENAMED section
func parseRenamedSection(content string) []RenameOp {
	var renamed []RenameOp

	// Find the RENAMED section header
	sectionPattern := regexp.MustCompile(`(?m)^##\s+RENAMED\s+Requirements\s*$`)
	matches := sectionPattern.FindStringIndex(content)
	if matches == nil {
		return renamed
	}

	// Extract content from this section until next ## header or end of file
	sectionStart := matches[1]
	nextSectionPattern := regexp.MustCompile(`(?m)^##\s+`)
	nextMatches := nextSectionPattern.FindStringIndex(content[sectionStart:])

	var sectionContent string
	if nextMatches != nil {
		sectionContent = content[sectionStart : sectionStart+nextMatches[0]]
	} else {
		sectionContent = content[sectionStart:]
	}

	// Parse FROM/TO pairs
	// Expected format:
	// - FROM: `### Requirement: Old Name`
	// - TO: `### Requirement: New Name`
	fromPattern := regexp.MustCompile(
		`^-\s*FROM:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `\s*$`,
	)
	toPattern := regexp.MustCompile(
		`^-\s*TO:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `\s*$`,
	)

	var currentFrom string
	scanner := bufio.NewScanner(strings.NewReader(sectionContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for FROM line
		if matches := fromPattern.FindStringSubmatch(line); len(matches) > 1 {
			currentFrom = strings.TrimSpace(matches[1])

			continue
		}

		// Check for TO line
		matches := toPattern.FindStringSubmatch(line)
		if len(matches) <= 1 || currentFrom == "" {
			continue
		}

		renamed = append(renamed, RenameOp{
			From: currentFrom,
			To:   strings.TrimSpace(matches[1]),
		})
		currentFrom = ""
	}

	return renamed
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
