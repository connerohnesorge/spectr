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

	// Path constants
	currentDirPath = "."
)

// FormatMode represents the display mode for multi-root formatting.
type FormatMode int

const (
	// FormatModeSingle indicates single-root mode (no prefix).
	FormatModeSingle FormatMode = iota
	// FormatModeMulti indicates multi-root mode (show root prefix).
	FormatModeMulti
)

// NewFormatMode creates a FormatMode based on whether there are multiple roots.
//
//nolint:revive // flag-parameter: hasMultipleRoots intentionally controls mode selection
func NewFormatMode(hasMultipleRoots bool) FormatMode {
	if hasMultipleRoots {
		return FormatModeMulti
	}

	return FormatModeSingle
}

// IsMulti returns true if this is multi-root mode.
func (m FormatMode) IsMulti() bool {
	return m == FormatModeMulti
}

// FormatItemWithRoot formats an item ID with a root prefix when in multi-root mode.
// For single-root or when rootPath is "." or empty, returns just the ID.
func FormatItemWithRoot(rootPath, id string, mode FormatMode) string {
	if rootPath == currentDirPath || rootPath == "" || !mode.IsMulti() {
		return id
	}

	return fmt.Sprintf("[%s] %s", rootPath, id)
}

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

// Multi-root formatting functions

// FormatChangesTextMulti formats changes with optional root prefix for multi-root scenarios.
func FormatChangesTextMulti(
	changes []ChangeInfo,
	mode FormatMode,
) string {
	if len(changes) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ID < changes[j].ID
	})

	// Find the longest ID and root for alignment
	maxIDLen := 0
	maxRootLen := 0
	for _, change := range changes {
		if len(change.ID) > maxIDLen {
			maxIDLen = len(change.ID)
		}
		if len(change.RootPath) > maxRootLen && mode.IsMulti() {
			maxRootLen = len(change.RootPath)
		}
	}

	lines := make([]string, 0, len(changes))
	for _, change := range changes {
		var line string
		if change.RootPath != currentDirPath && change.RootPath != "" && mode.IsMulti() {
			line = fmt.Sprintf("[%-*s] %-*s  %d/%d tasks",
				maxRootLen,
				change.RootPath,
				maxIDLen,
				change.ID,
				change.TaskStatus.Completed,
				change.TaskStatus.Total,
			)
		} else {
			line = fmt.Sprintf("%-*s  %d/%d tasks",
				maxIDLen,
				change.ID,
				change.TaskStatus.Completed,
				change.TaskStatus.Total,
			)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatChangesLongMulti formats changes with detailed information and optional root prefix.
func FormatChangesLongMulti(
	changes []ChangeInfo,
	mode FormatMode,
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
		var line string
		if change.RootPath != currentDirPath && change.RootPath != "" && mode.IsMulti() {
			line = fmt.Sprintf(
				"[%s] %s: %s [deltas %d] [tasks %d/%d]",
				change.RootPath,
				change.ID,
				change.Title,
				change.DeltaCount,
				change.TaskStatus.Completed,
				change.TaskStatus.Total,
			)
		} else {
			line = fmt.Sprintf(
				"%s: %s [deltas %d] [tasks %d/%d]",
				change.ID,
				change.Title,
				change.DeltaCount,
				change.TaskStatus.Completed,
				change.TaskStatus.Total,
			)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatSpecsTextMulti formats specs with optional root prefix for multi-root scenarios.
func FormatSpecsTextMulti(
	specs []SpecInfo,
	mode FormatMode,
) string {
	if len(specs) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ID < specs[j].ID
	})

	lines := make([]string, 0, len(specs))
	for _, spec := range specs {
		var line string
		if spec.RootPath != currentDirPath && spec.RootPath != "" && mode.IsMulti() {
			line = fmt.Sprintf("[%s] %s", spec.RootPath, spec.ID)
		} else {
			line = spec.ID
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatSpecsLongMulti formats specs with detailed information and optional root prefix.
func FormatSpecsLongMulti(
	specs []SpecInfo,
	mode FormatMode,
) string {
	if len(specs) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ID < specs[j].ID
	})

	lines := make([]string, 0, len(specs))
	for _, spec := range specs {
		var line string
		if spec.RootPath != currentDirPath && spec.RootPath != "" && mode.IsMulti() {
			line = fmt.Sprintf(
				"[%s] %s: %s [requirements %d]",
				spec.RootPath,
				spec.ID,
				spec.Title,
				spec.RequirementCount,
			)
		} else {
			line = fmt.Sprintf(
				"%s: %s [requirements %d]",
				spec.ID,
				spec.Title,
				spec.RequirementCount,
			)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}
