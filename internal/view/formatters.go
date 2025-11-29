// Package view provides dashboard functionality for displaying
// a comprehensive project overview including specs, changes, and tasks.
package view

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/connerohnesorge/spectr/internal/tui"
)

const (
	// Dashboard title
	dashboardTitle = "Spectr Dashboard"

	// Separator characters
	doubleLineSeparator = "═══════════════════════════════════════" +
		"═════════════════════════"
	singleLineSeparator = "───────────────────────────────────────" +
		"─────────────────────────"

	// Section headers
	summaryHeader          = "Summary:"
	activeChangesHeader    = "Active Changes"
	completedChangesHeader = "Completed Changes"
	specsHeader            = "Specifications"

	// Footer hint
	footerHint = "Use spectr list --changes or " +
		"spectr list --specs for detailed views"

	// Indicators
	summaryBullet      = "●"
	activeChangeCircle = "◉"
	completedCheckmark = "✓"
	specSquare         = "▪"
	indentation        = "  "
	// Fixed width for change IDs in active changes section
	changeIDWidth = 28
	// Fixed width for spec IDs in specs section
	specIDWidth = 28

	// Percentage calculation constant
	percentageMultiplier = 100
	// Newline separator
	newline = "\n"
)

// Style functions that use centralized theme configuration.
// These are functions rather than variables to ensure colors are evaluated
// at runtime after theme initialization.

// headerStyle returns the style for section headers (bold, accent color)
func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(tui.GetAccentColor()))
}

// summaryBulletStyle returns the style for summary bullets (accent color)
func summaryBulletStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(tui.GetAccentColor()))
}

// activeChangeStyle returns the style for active change indicators
// (header color for standout)
func activeChangeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(tui.GetHeaderColor()))
}

// completedChangeStyle returns the style for completed change
// indicators (success/green)
func completedChangeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(tui.GetSuccessColor()))
}

// specStyle returns the style for spec indicators (accent color)
func specStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(tui.GetAccentColor()))
}

// footerStyle returns the style for footer hints (help/dim)
func footerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(tui.GetHelpColor()))
}

// FormatDashboardText formats the dashboard data as
// human-readable terminal output with colored sections,
// progress bars, and visual indicators.
//
// The output follows the format specification in design.md:
// - Dashboard title with double-line separator
// - Summary section with bullet points and metrics
// - Active changes with progress bars
// - Completed changes with checkmarks
// - Specifications with requirement counts
// - Footer with navigation hints
//
// Empty sections (e.g., no completed changes) are automatically hidden.
func FormatDashboardText(data *DashboardData) string {
	var sections []string

	// Header: Dashboard title
	sections = append(sections, dashboardTitle)
	sections = append(sections, "")
	sections = append(sections, doubleLineSeparator)
	sections = append(sections, "")

	// Section 1: Summary
	sections = append(sections, formatSummarySection(data.Summary))
	sections = append(sections, "")

	// Section 2: Active Changes (only if there are active changes)
	if len(data.ActiveChanges) > 0 {
		sections = append(sections,
			formatActiveChangesSection(data.ActiveChanges))
		sections = append(sections, "")
	}

	// Section 3: Completed Changes (only if there are completed changes)
	if len(data.CompletedChanges) > 0 {
		sections = append(sections,
			formatCompletedChangesSection(data.CompletedChanges))
		sections = append(sections, "")
	}

	// Section 4: Specifications (only if there are specs)
	if len(data.Specs) > 0 {
		sections = append(sections, formatSpecsSection(data.Specs))
		sections = append(sections, "")
	}

	// Footer: Double-line separator and hints
	sections = append(sections, doubleLineSeparator)
	sections = append(sections, "")
	sections = append(sections, footerStyle().Render(footerHint))

	return strings.Join(sections, "\n")
}

// formatSummarySection creates the summary section with aggregate metrics
func formatSummarySection(summary SummaryMetrics) string {
	var lines []string

	lines = append(lines, summaryHeader)

	// Specifications: X specs, Y requirements
	specsLine := fmt.Sprintf("%s %s Specifications: %d specs, %d requirements",
		indentation,
		summaryBulletStyle().Render(summaryBullet),
		summary.TotalSpecs,
		summary.TotalRequirements,
	)
	lines = append(lines, specsLine)

	// Active Changes: X in progress
	activeLine := fmt.Sprintf(
		"%s %s Active Changes: %d in progress",
		indentation,
		summaryBulletStyle().Render(summaryBullet),
		summary.ActiveChanges,
	)
	lines = append(lines, activeLine)

	// Completed Changes: X
	completedLine := fmt.Sprintf("%s %s Completed Changes: %d",
		indentation,
		summaryBulletStyle().Render(summaryBullet),
		summary.CompletedChanges,
	)
	lines = append(lines, completedLine)

	// Task Progress: X/Y (Z% complete)
	taskPercentage := 0
	if summary.TotalTasks > 0 {
		taskPercentage = (summary.CompletedTasks *
			percentageMultiplier) / summary.TotalTasks
	}
	taskLine := fmt.Sprintf(
		"%s %s Task Progress: %d/%d (%d%% complete)",
		indentation,
		summaryBulletStyle().Render(summaryBullet),
		summary.CompletedTasks,
		summary.TotalTasks,
		taskPercentage,
	)
	lines = append(lines, taskLine)

	return strings.Join(lines, newline)
}

// formatActiveChangesSection creates the active changes section
// with progress bars
func formatActiveChangesSection(changes []ChangeProgress) string {
	var lines []string

	// Section header
	lines = append(lines, headerStyle().Render(activeChangesHeader))
	lines = append(lines, singleLineSeparator)

	// Each active change: ◉ id [progress bar] percentage%
	for _, change := range changes {
		// Render progress bar using progress.RenderBar()
		progressBar := RenderBar(change.Progress.Completed,
			change.Progress.Total)

		// Format: "  ◉ change-id              [████████░░░░░░░░░░░░] 37%"
		line := fmt.Sprintf("%s %s %-*s %s",
			indentation,
			activeChangeStyle().Render(activeChangeCircle),
			changeIDWidth,
			change.ID,
			progressBar,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, newline)
}

// formatCompletedChangesSection creates the completed changes
// section with checkmarks
func formatCompletedChangesSection(changes []CompletedChange) string {
	var lines []string

	// Section header
	lines = append(lines, headerStyle().Render(completedChangesHeader))
	lines = append(lines, singleLineSeparator)

	// Each completed change: ✓ id
	for _, change := range changes {
		line := fmt.Sprintf("%s %s %s",
			indentation,
			completedChangeStyle().Render(completedCheckmark),
			change.ID,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, newline)
}

// formatSpecsSection creates the specifications section with requirement counts
func formatSpecsSection(specs []SpecInfo) string {
	var lines []string

	// Section header
	lines = append(lines, headerStyle().Render(specsHeader))
	lines = append(lines, singleLineSeparator)

	// Each spec: ▪ id                  X requirements
	for _, spec := range specs {
		line := fmt.Sprintf("%s %s %-*s %d requirements",
			indentation,
			specStyle().Render(specSquare),
			specIDWidth,
			spec.ID,
			spec.RequirementCount,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, newline)
}

// FormatDashboardJSON formats the dashboard data as
// machine-readable JSON output.
//
// The output structure matches the schema defined in design.md:
//   - summary: Aggregate metrics
//     (totalSpecs, totalRequirements, activeChanges, etc.)
//   - activeChanges: Array of changes in progress with
//     task completion metrics
//   - completedChanges: Array of completed changes
//   - specs: Array of specifications with requirement counts
//
// Arrays are pre-sorted by the CollectData() function,
// ensuring consistent output.
// The JSON is formatted with indentation for human readability.
//
// Returns the formatted JSON string and any marshaling error.
func FormatDashboardJSON(data *DashboardData) (string, error) {
	// Marshal with indentation for readability
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf(
			"failed to marshal dashboard data to JSON: %w",
			err,
		)
	}

	return string(jsonBytes), nil
}
