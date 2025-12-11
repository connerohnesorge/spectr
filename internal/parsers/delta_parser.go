// Package parsers provides functions for parsing delta specifications,
// requirements, and other structured spec documents.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package parsers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	plan := &DeltaPlan{
		Added:    make([]RequirementBlock, 0),
		Modified: make([]RequirementBlock, 0),
		Removed:  make([]string, 0),
		Renamed:  make([]RenameOp, 0),
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse each section
	plan.Added = parseDeltaSection(
		string(content),
		"ADDED",
	)
	plan.Modified = parseDeltaSection(
		string(content),
		"MODIFIED",
	)
	plan.Removed = parseRemovedSection(
		string(content),
	)
	plan.Renamed = parseRenamedSection(
		string(content),
	)

	return plan, nil
}

// parseDeltaSection extracts requirements from a delta section
func parseDeltaSection(
	content, sectionType string,
) []RequirementBlock {
	sectionContent := extractSectionContent(
		content,
		sectionType,
	)
	if sectionContent == "" {
		return nil
	}

	return parseRequirementsFromSection(
		sectionContent,
	)
}

// extractSectionContent extracts content from a section header
func extractSectionContent(
	content, sectionType string,
) string {
	pattern := fmt.Sprintf(
		`(?m)^##\s+%s\s+Requirements\s*$`,
		sectionType,
	)
	sectionPattern := regexp.MustCompile(pattern)
	matches := sectionPattern.FindStringIndex(
		content,
	)
	if matches == nil {
		return ""
	}

	sectionStart := matches[1]
	nextSectionPattern := regexp.MustCompile(
		`(?m)^##\s+`,
	)
	nextMatches := nextSectionPattern.FindStringIndex(
		content[sectionStart:],
	)

	if nextMatches != nil {
		return content[sectionStart : sectionStart+nextMatches[0]]
	}

	return content[sectionStart:]
}

// parseRequirementsFromSection parses requirement blocks from content
func parseRequirementsFromSection(
	sectionContent string,
) []RequirementBlock {
	var requirements []RequirementBlock
	var currentReq *RequirementBlock

	reqPattern := regexp.MustCompile(
		`^###\s+Requirement:\s*(.+)$`,
	)
	h3Pattern := regexp.MustCompile(`^###\s+`)

	scanner := bufio.NewScanner(
		strings.NewReader(sectionContent),
	)
	for scanner.Scan() {
		line := scanner.Text()

		if matches := reqPattern.FindStringSubmatch(line); len(
			matches,
		) > 1 {
			currentReq = saveAndStartNewRequirement(
				&requirements,
				currentReq,
				line,
				matches[1],
			)

			continue
		}

		if isNonRequirementH3(
			line,
			h3Pattern,
			reqPattern,
		) {
			currentReq = saveCurrentRequirement(
				&requirements,
				currentReq,
			)

			continue
		}

		appendLineToRequirement(currentReq, line)
	}

	// Save the last requirement
	saveCurrentRequirement(
		&requirements,
		currentReq,
	)

	return requirements
}

// saveAndStartNewRequirement saves current req and starts a new one
func saveAndStartNewRequirement(
	requirements *[]RequirementBlock,
	currentReq *RequirementBlock,
	line, name string,
) *RequirementBlock {
	if currentReq != nil {
		*requirements = append(
			*requirements,
			*currentReq,
		)
	}

	return &RequirementBlock{
		HeaderLine: line,
		Name:       strings.TrimSpace(name),
		Raw:        line + "\n",
	}
}

// saveCurrentRequirement saves the current requirement if it exists
func saveCurrentRequirement(
	requirements *[]RequirementBlock,
	currentReq *RequirementBlock,
) *RequirementBlock {
	if currentReq != nil {
		*requirements = append(
			*requirements,
			*currentReq,
		)
	}

	return nil
}

// isNonRequirementH3 checks if line is an H3 but not a requirement
func isNonRequirementH3(
	line string,
	h3Pattern, reqPattern *regexp.Regexp,
) bool {
	return h3Pattern.MatchString(line) &&
		!reqPattern.MatchString(line)
}

// appendLineToRequirement appends a line to the current requirement
func appendLineToRequirement(
	currentReq *RequirementBlock,
	line string,
) {
	if currentReq != nil {
		currentReq.Raw += line + "\n"
	}
}

// parseRemovedSection extracts requirement names from REMOVED section
func parseRemovedSection(
	content string,
) []string {
	var removed []string

	// Find the REMOVED section header
	sectionPattern := regexp.MustCompile(
		`(?m)^##\s+REMOVED\s+Requirements\s*$`,
	)
	matches := sectionPattern.FindStringIndex(
		content,
	)
	if matches == nil {
		return removed
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

	// Parse requirement headers within this section
	reqPattern := regexp.MustCompile(
		`^###\s+Requirement:\s*(.+)$`,
	)

	scanner := bufio.NewScanner(
		strings.NewReader(sectionContent),
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := reqPattern.FindStringSubmatch(line); len(
			matches,
		) > 1 {
			removed = append(
				removed,
				strings.TrimSpace(matches[1]),
			)
		}
	}

	return removed
}

// parseRenamedSection extracts FROM/TO pairs from RENAMED section
func parseRenamedSection(
	content string,
) []RenameOp {
	var renamed []RenameOp

	// Find the RENAMED section header
	sectionPattern := regexp.MustCompile(
		`(?m)^##\s+RENAMED\s+Requirements\s*$`,
	)
	matches := sectionPattern.FindStringIndex(
		content,
	)
	if matches == nil {
		return renamed
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

	var currentFrom string
	scanner := bufio.NewScanner(
		strings.NewReader(sectionContent),
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for FROM line
		if matches := fromPattern.FindStringSubmatch(line); len(
			matches,
		) > 1 {
			currentFrom = strings.TrimSpace(
				matches[1],
			)

			continue
		}

		// Check for TO line
		matches := toPattern.FindStringSubmatch(
			line,
		)
		if len(matches) <= 1 ||
			currentFrom == "" {
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
