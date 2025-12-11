// Package list provides functionality for listing and formatting changes and
// specifications. This file handles unified formatting for displaying both
// changes and specs together.
package list

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// FormatAllText formats all items (changes and specs) as text with type
// indicators. Items are sorted by ID and displayed with their type and
// relevant details (tasks for changes, requirements for specs).
func FormatAllText(items ItemList) string {
	if len(items) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID() < items[j].ID()
	})

	// Find the longest ID for alignment
	maxIDLen := 0
	for _, item := range items {
		if len(item.ID()) > maxIDLen {
			maxIDLen = len(item.ID())
		}
	}

	// Build lines for each item with type indicator
	lines := make([]string, 0, len(items))
	for _, item := range items {
		var typeIndicator, details string
		switch item.Type {
		case ItemTypeChange:
			typeIndicator = "[CHANGE]"
			if item.Change != nil {
				details = fmt.Sprintf(
					"%d/%d tasks",
					item.Change.TaskStatus.Completed,
					item.Change.TaskStatus.Total,
				)
			}
		case ItemTypeSpec:
			typeIndicator = "[SPEC]  "
			if item.Spec != nil {
				details = fmt.Sprintf(
					"%d requirements",
					item.Spec.RequirementCount,
				)
			}
		}

		line := fmt.Sprintf("%-*s  %s  %s",
			maxIDLen,
			item.ID(),
			typeIndicator,
			details,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, lineSeparator)
}

// FormatAllLong formats all items with detailed information. Shows ID, type,
// title, and relevant counts (deltas/tasks for changes, requirements for
// specs). Items are sorted by ID.
func FormatAllLong(items ItemList) string {
	if len(items) == 0 {
		return noItemsFoundMsg
	}

	// Sort by ID
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID() < items[j].ID()
	})

	var lines []string
	for _, item := range items {
		var line string
		switch item.Type {
		case ItemTypeChange:
			if item.Change != nil {
				line = fmt.Sprintf(
					"%s [CHANGE]: %s [deltas %d] [tasks %d/%d]",
					item.Change.ID,
					item.Change.Title,
					item.Change.DeltaCount,
					item.Change.TaskStatus.Completed,
					item.Change.TaskStatus.Total,
				)
			}
		case ItemTypeSpec:
			if item.Spec != nil {
				line = fmt.Sprintf(
					"%s [SPEC]: %s [requirements %d]",
					item.Spec.ID,
					item.Spec.Title,
					item.Spec.RequirementCount,
				)
			}
		}
		if line != "" {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, lineSeparator)
}

// FormatAllJSON formats all items as a JSON array. Each item includes its type
// (change or spec) along with all relevant fields. Items are sorted by ID
// before serialization. Returns an empty JSON array if no items are provided.
func FormatAllJSON(
	items ItemList,
) (string, error) {
	if len(items) == 0 {
		return "[]", nil
	}

	// Sort by ID for consistent output
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID() < items[j].ID()
	})

	// Marshal items to JSON with indentation for readability
	data, err := json.MarshalIndent(
		items,
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
