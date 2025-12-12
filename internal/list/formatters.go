// Package list provides functionality for listing and formatting
// changes and specifications in various output formats.
package list

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	// Common messages
	noItemsFoundMsg = "No items found"
	lineSeparator   = "\n"
)

// FormatChangesText formats changes as simple text list (IDs only)
func FormatChangesText(
	changes []ChangeInfo,
) string {
	if len(changes) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ID < changes[j].ID
	})

	// Find the longest ID for alignment
	maxIDLen := 0
	for _, change := range changes {
		if len(change.ID) > maxIDLen {
			maxIDLen = len(change.ID)
		}
	}

	lines := make([]string, 0, len(changes))
	for _, change := range changes {
		line := fmt.Sprintf("%-*s  %d/%d tasks",
			maxIDLen,
			change.ID,
			change.TaskStatus.Completed,
			change.TaskStatus.Total,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatChangesLong formats changes with detailed information
func FormatChangesLong(
	changes []ChangeInfo,
) string {
	if len(changes) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ID < changes[j].ID
	})

	lines := make([]string, 0, len(changes))
	for _, change := range changes {
		line := fmt.Sprintf(
			"%s: %s [deltas %d] [tasks %d/%d]",
			change.ID,
			change.Title,
			change.DeltaCount,
			change.TaskStatus.Completed,
			change.TaskStatus.Total,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatChangesJSON formats changes as JSON array
func FormatChangesJSON(
	changes []ChangeInfo,
) (string, error) {
	if len(changes) == 0 {
		return "[]", nil
	}

	// Sort by ID
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ID < changes[j].ID
	})

	data, err := json.MarshalIndent(
		changes,
		"",
		"  ",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to marshal JSON: %w",
			err,
		)
	}

	return string(data), nil
}

// FormatSpecsText formats specs as simple text list (IDs only)
func FormatSpecsText(specs []SpecInfo) string {
	if len(specs) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ID < specs[j].ID
	})

	lines := make([]string, 0, len(specs))
	for _, spec := range specs {
		lines = append(lines, spec.ID)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatSpecsLong formats specs with detailed information
func FormatSpecsLong(specs []SpecInfo) string {
	if len(specs) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ID < specs[j].ID
	})

	lines := make([]string, 0, len(specs))
	for _, spec := range specs {
		line := fmt.Sprintf(
			"%s: %s [requirements %d]",
			spec.ID,
			spec.Title,
			spec.RequirementCount,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatSpecsJSON formats specs as JSON array
func FormatSpecsJSON(
	specs []SpecInfo,
) (string, error) {
	if len(specs) == 0 {
		return "[]", nil
	}

	// Sort by ID
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ID < specs[j].ID
	})

	data, err := json.MarshalIndent(
		specs,
		"",
		"  ",
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to marshal JSON: %w",
			err,
		)
	}

	return string(data), nil
}
