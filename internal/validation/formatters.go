// Package validation provides validation result printing functions.
// This file contains functions for outputting validation results
// in both JSON and human-readable formats.
package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// Color constants for validation output styling
const (
	ColorError   = "1" // Red
	ColorWarning = "3" // Yellow
)

var (
	// errorStyle styles [ERROR] labels in red
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)).
			Bold(true)
	// warningStyle styles [WARNING] labels in yellow
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)).
			Bold(true)
)

// isTTY returns true if stdout is a terminal
func isTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// formatLevel formats a validation level with color styling if in a TTY
func formatLevel(level ValidationLevel) string {
	label := fmt.Sprintf("[%s]", level)
	if !isTTY() {
		return label
	}
	switch level {
	case LevelError:
		return errorStyle.Render(label)
	case LevelWarning:
		return warningStyle.Render(label)
	case LevelInfo:
		return label
	}

	return label
}

// ToRelativePath converts an absolute path to a path relative to the
// spectr/ directory. It removes everything up to and including the
// spectr/ prefix.
// Example: /home/user/project/spectr/changes/foo/spec.md -> changes/foo/spec.md
func ToRelativePath(absPath string) string {
	// Clean the path first
	cleanPath := filepath.Clean(absPath)

	// Look for spectr/ in the path
	const spectrDir = "spectr" + string(filepath.Separator)
	if idx := strings.Index(cleanPath, spectrDir); idx >= 0 {
		return cleanPath[idx+len(spectrDir):]
	}

	// If no spectr/ found, just return the base name
	return filepath.Base(cleanPath)
}

// BulkResult represents the result of validating a single item
type BulkResult struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Valid  bool              `json:"valid"`
	Report *ValidationReport `json:"report,omitempty"`
	Error  string            `json:"error,omitempty"`
}

// PrintJSONReport prints a single validation report as JSON
func PrintJSONReport(
	report *ValidationReport,
) {
	data, err := json.MarshalIndent(
		report,
		"",
		"  ",
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error marshaling JSON: %v\n",
			err,
		)

		return
	}
	fmt.Println(string(data))
}

// PrintHumanReport prints a single validation report in human format
func PrintHumanReport(
	itemName string,
	report *ValidationReport,
) {
	if report.Valid {
		fmt.Printf("✓ %s valid\n", itemName)

		return
	}

	issueCount := len(report.Issues)
	fmt.Printf(
		"✗ %s has %d issue(s):\n",
		itemName,
		issueCount,
	)

	for _, issue := range report.Issues {
		fmt.Printf(
			"  [%s] %s: %s\n",
			issue.Level,
			issue.Path,
			issue.Message,
		)
	}
}

// PrintBulkJSONResults prints bulk validation results as JSON
func PrintBulkJSONResults(results []BulkResult) {
	data, err := json.MarshalIndent(
		results,
		"",
		"  ",
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error marshaling JSON: %v\n",
			err,
		)

		return
	}
	fmt.Println(string(data))
}

// PrintBulkHumanResults prints bulk validation results in human format
// with improved formatting: visual separation between failed items,
// relative paths, issues grouped by file, colored error/warning labels,
// enhanced summary, and type indicators.
func PrintBulkHumanResults(results []BulkResult) {
	passCount := 0
	failCount := 0
	errorCount := 0
	warningCount := 0
	isFirstFailed := true

	for _, result := range results {
		if result.Valid {
			fmt.Printf(
				"✓ %s (%s)\n",
				result.Name,
				result.Type,
			)
			passCount++
		} else {
			// Add blank line before each failed item (except the first)
			if !isFirstFailed {
				fmt.Println()
			}
			isFirstFailed = false

			if result.Error != "" {
				fmt.Printf(
					"✗ %s (%s): %s\n",
					result.Name,
					result.Type,
					result.Error,
				)
				errorCount++
			} else {
				issueCount := len(result.Report.Issues)
				fmt.Printf(
					"✗ %s (%s) has %d issue(s):\n",
					result.Name,
					result.Type,
					issueCount,
				)
				// Count and print issues grouped by file
				printGroupedIssues(
					result.Report.Issues, &errorCount, &warningCount,
				)
			}
			failCount++
		}
	}

	// Print enhanced summary
	printSummary(summaryParams{
		passCount:    passCount,
		failCount:    failCount,
		errorCount:   errorCount,
		warningCount: warningCount,
		total:        len(results),
	})
}

// printGroupedIssues prints issues grouped by their file path with indentation
func printGroupedIssues(
	issues []ValidationIssue,
	errorCount, warningCount *int,
) {
	// Group issues by file path
	grouped := make(map[string][]ValidationIssue)
	var order []string // Preserve order of first occurrence

	for _, issue := range issues {
		relPath := ToRelativePath(issue.Path)
		if _, exists := grouped[relPath]; !exists {
			order = append(order, relPath)
		}
		grouped[relPath] = append(
			grouped[relPath],
			issue,
		)

		// Count errors and warnings
		switch issue.Level {
		case LevelError:
			*errorCount++
		case LevelWarning:
			*warningCount++
		case LevelInfo:
			// Info level issues are not counted in error/warning totals
		}
	}

	// Print issues grouped by file
	for _, path := range order {
		fileIssues := grouped[path]
		if len(fileIssues) == 1 {
			// Single issue: print on one line
			issue := fileIssues[0]
			fmt.Printf("  %s %s: %s\n",
				formatLevel(issue.Level),
				path,
				issue.Message,
			)
		} else {
			// Multiple issues: print file header then indented issues
			fmt.Printf("  %s:\n", path)
			for _, issue := range fileIssues {
				fmt.Printf("    %s %s\n",
					formatLevel(issue.Level),
					issue.Message,
				)
			}
		}
	}
}

// summaryParams holds parameters for printing the validation summary
type summaryParams struct {
	passCount    int
	failCount    int
	errorCount   int
	warningCount int
	total        int
}

// printSummary prints the validation summary line
func printSummary(p summaryParams) {
	if p.failCount > 0 {
		fmt.Printf(
			"\n%d passed, %d failed (%d errors, %d warnings), %d total\n",
			p.passCount,
			p.failCount,
			p.errorCount,
			p.warningCount,
			p.total,
		)
	} else {
		fmt.Printf(
			"\n%d passed, %d failed, %d total\n",
			p.passCount,
			p.failCount,
			p.total,
		)
	}
}
