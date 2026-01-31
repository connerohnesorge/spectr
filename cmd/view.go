// Package cmd provides command-line interface implementations for Spectr.
// This file contains the view command for displaying the project dashboard.
package cmd

import (
	"fmt"

	"github.com/connerohnesorge/spectr/internal/view"
)

// ViewCmd represents the view command which displays a comprehensive
// project dashboard including summary metrics, active changes, completed
// changes, and specifications.
//
// The dashboard provides an at-a-glance overview of the entire project state:
//   - Summary metrics: total specs, requirements, changes, and task progress
//   - Active changes: changes in progress with visual progress bars
//   - Completed changes: changes with all tasks complete
//   - Specifications: all specs with requirement counts
//
// Output formats:
//   - Default: Colored terminal output with Unicode box-drawing characters
//   - --json: Machine-readable JSON for automation and scripting
//
// The terminal output uses lipgloss for styling and requires a terminal
// with Unicode support for optimal display. All modern terminal emulators
// (iTerm2, GNOME Terminal, Windows Terminal, Terminal.app) are supported.
type ViewCmd struct {
	// JSON enables JSON output format for scripting and automation.
	// When enabled, outputs structured data matching the schema defined
	// in the view command design specification.
	JSON bool `kong:"help='Output in JSON format for scripting'"`
}

// Run executes the view command.
// It collects dashboard data from all discovered spectr roots and formats the output
// based on the JSON flag (either human-readable text or JSON).
// Returns an error if no spectr directories are found or if
// discovery/parsing fails.
func (c *ViewCmd) Run() error {
	// Discover all spectr roots
	roots, err := GetDiscoveredRoots()
	if err != nil {
		return fmt.Errorf(
			"failed to discover spectr roots: %w",
			err,
		)
	}

	// If no roots found, show empty dashboard (graceful degradation)
	if len(roots) == 0 {
		// Return empty dashboard data
		data := &view.DashboardData{
			Summary:          view.SummaryMetrics{},
			ActiveChanges:    make([]view.ChangeProgress, 0),
			CompletedChanges: make([]view.CompletedChange, 0),
			Specs:            make([]view.SpecInfo, 0),
		}

		output := view.FormatDashboardText(data)
		fmt.Println(output)

		return nil
	}

	// Collect dashboard data from all roots
	data, err := view.CollectDataMultiRoot(roots)
	if err != nil {
		// Handle other discovery/parsing failures
		return fmt.Errorf(
			"failed to collect dashboard data: %w",
			err,
		)
	}

	// Format and output the dashboard
	var output string
	if c.JSON {
		// JSON format for machine consumption
		output, err = view.FormatDashboardJSON(
			data,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to format JSON: %w",
				err,
			)
		}
	} else {
		// Human-readable text format with colors and progress bars
		output = view.FormatDashboardText(data)
	}

	// Print the formatted output
	fmt.Println(output)

	return nil
}
