// Package cmd provides command-line interface implementations for Spectr.
// This file contains the view command for displaying the project dashboard.
package cmd

import (
	"fmt"
	"os"

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
// It collects dashboard data from the project and formats the output
// based on the JSON flag (either human-readable text or JSON).
// Returns an error if the spectr directory is missing or if
// discovery/parsing fails.
func (c *ViewCmd) Run() error {
	// Get current working directory as the project path
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(
			"failed to get current directory: %w",
			err,
		)
	}

	// Collect dashboard data from the project
	data, err := view.CollectData(projectPath)
	if err != nil {
		// Handle missing spectr directory error
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"spectr directory not found: %w\n"+
					"Hint: Run 'spectr init' to initialize Spectr",
				err,
			)
		}
		// Handle other discovery/parsing failures
		return fmt.Errorf("failed to collect dashboard data: %w", err)
	}

	// Format and output the dashboard
	var output string
	if c.JSON {
		// JSON format for machine consumption
		output, err = view.FormatDashboardJSON(data)
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
	} else {
		// Human-readable text format with colors and progress bars
		output = view.FormatDashboardText(data)
	}

	// Print the formatted output
	fmt.Println(output)

	return nil
}
