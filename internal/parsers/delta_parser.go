// Package parsers provides functions for parsing delta specifications,
// requirements, and other structured spec documents.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package parsers

import (
	"os"
	"regexp"
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

// ParseDeltaSpec parses a delta spec file and extracts operations
// Returns a DeltaPlan with ADDED, MODIFIED, REMOVED, and RENAMED reqs
func ParseDeltaSpec(
	filePath string,
) (*DeltaPlan, error) {
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

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		// Return empty plan if content is empty or invalid
		return plan, nil
	}

	// Parse each section using the parsed document
	plan.Added = parseDeltaSectionFromDoc(doc, "ADDED")
	plan.Modified = parseDeltaSectionFromDoc(doc, "MODIFIED")
	plan.Removed = parseRemovedSectionFromDoc(doc)
	plan.Renamed = parseRenamedSection(string(content)) // Keep regex for complex FROM/TO parsing

	return plan, nil
}

// parseDeltaSectionFromDoc extracts requirements from a delta section using markdown package
func parseDeltaSectionFromDoc(
	doc *markdown.Document,
	sectionType string,
) []RequirementBlock {
	var requirements []RequirementBlock

	// Find the delta section header
	sectionHeader := sectionType + " Requirements"
	var inSection bool
	var sectionLevel int

	for i, header := range doc.Headers {
		// Check if this is the start of our target section
		if header.Level == 2 && header.Text == sectionHeader {
			inSection = true
			sectionLevel = header.Level

			continue
		}

		// If we're in the section and hit another H2, stop
		if inSection && header.Level <= sectionLevel {
			break
		}

		// Process requirements within the section
		if inSection && header.Level == 3 && strings.HasPrefix(header.Text, "Requirement:") {
			name := strings.TrimPrefix(header.Text, "Requirement:")
			name = strings.TrimSpace(name)

			headerLine := "### " + header.Text

			// Get section content
			sectionContent := ""
			if section, ok := doc.Sections[header.Text]; ok {
				sectionContent = section.Content
			}

			// Build raw content
			raw := headerLine + "\n"
			if sectionContent != "" {
				raw += sectionContent + "\n"
			}

			requirements = append(requirements, RequirementBlock{
				HeaderLine: headerLine,
				Name:       name,
				Raw:        raw,
			})

			// Skip any nested headers that belong to this requirement
			for j := i + 1; j < len(doc.Headers); j++ {
				if doc.Headers[j].Level <= 3 {
					break
				}
			}
		}
	}

	return requirements
}

// parseRemovedSectionFromDoc extracts requirement names from REMOVED section using markdown package
func parseRemovedSectionFromDoc(
	doc *markdown.Document,
) []string {
	var removed []string

	// Find the REMOVED section header
	var inSection bool

	for _, header := range doc.Headers {
		// Check if this is the REMOVED section
		if header.Level == 2 && header.Text == "REMOVED Requirements" {
			inSection = true

			continue
		}

		// If we're in the section and hit another H2, stop
		if inSection && header.Level == 2 {
			break
		}

		// Process requirements within the section
		if inSection && header.Level == 3 && strings.HasPrefix(header.Text, "Requirement:") {
			name := strings.TrimPrefix(header.Text, "Requirement:")
			name = strings.TrimSpace(name)
			removed = append(removed, name)
		}
	}

	return removed
}

// parseRenamedSection extracts FROM/TO pairs from RENAMED section
// Note: This uses regex because the FROM/TO format uses backticks and
// complex inline formatting that is better parsed with regex patterns.
func parseRenamedSection(
	content string,
) []RenameOp {
	// Find the RENAMED section header
	sectionPattern := regexp.MustCompile(
		`(?m)^##\s+RENAMED\s+Requirements\s*$`,
	)
	matches := sectionPattern.FindStringIndex(
		content,
	)
	if matches == nil {
		return []RenameOp{}
	}

	// Extract content from this section until next ## header or end of file
	sectionStart := matches[1]
	nextSectionPattern := regexp.MustCompile(
		`(?m)^##\s+`,
	)
	nextMatches := nextSectionPattern.FindStringIndex(
		content[sectionStart:],
	)

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

	lines := strings.Split(sectionContent, "\n")

	// Count FROM lines for pre-allocation
	fromCount := 0
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if fromPattern.MatchString(line) {
			fromCount++
		}
	}

	renamed := make([]RenameOp, 0, fromCount)
	var currentFrom string

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)

		// Check for FROM line
		if fromMatches := fromPattern.FindStringSubmatch(line); len(
			fromMatches,
		) > 1 {
			currentFrom = strings.TrimSpace(
				fromMatches[1],
			)

			continue
		}

		// Check for TO line
		toMatches := toPattern.FindStringSubmatch(
			line,
		)
		if len(toMatches) <= 1 ||
			currentFrom == "" {
			continue
		}

		renamed = append(renamed, RenameOp{
			From: currentFrom,
			To:   strings.TrimSpace(toMatches[1]),
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
