// Package view provides dashboard functionality for displaying
// a comprehensive project overview including specs, changes, and tasks.
package view

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/connerohnesorge/spectr/internal/theme"
)

const (
	// progressBarWidth is the fixed width of the progress bar
	// in characters
	progressBarWidth = 20

	// filledChar is the Unicode filled block character used
	// for completed portion
	filledChar = "█"

	// emptyChar is the Unicode light shade block character used
	// for remaining portion
	emptyChar = "░"
)

// filledStyle returns a style with the success color from the current theme
// for the filled portion of the progress bar
func filledStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Success)
}

// emptyStyle returns a style with the muted color from the current theme
// for the empty portion of the progress bar
func emptyStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.Current().Muted)
}

// RenderBar creates a visual progress bar string with the format:
// [████████████░░░░░░░░] 60%
// It uses a fixed width of 20 characters with filled (█) and
// empty (░) block characters.
// The filled portion is colored green and the empty portion
// is dim gray.
//
// Parameters:
//   - completed: number of completed tasks
//   - total: total number of tasks
//
// Returns a formatted string containing the progress bar and percentage.
//
// Edge case: If total is 0, renders an empty bar
// [░░░░░░░░░░░░░░░░░░░░] with dim styling.
func RenderBar(completed, total int) string {
	// Handle edge case of zero total tasks
	if total == 0 {
		emptyBar := ""
		for range progressBarWidth {
			emptyBar += emptyChar
		}

		return fmt.Sprintf("[%s] 0%%", emptyStyle().Render(emptyBar))
	}

	// Calculate percentage (0-100)
	percentage := int(
		math.Round((float64(completed) / float64(total)) * 100),
	)

	// Calculate filled width (0-20 characters)
	ratio := float64(completed) / float64(total)
	filledWidth := int(math.Round(ratio * float64(progressBarWidth)))

	// Ensure filled width is within bounds
	if filledWidth < 0 {
		filledWidth = 0
	}
	if filledWidth > progressBarWidth {
		filledWidth = progressBarWidth
	}

	// Build filled portion
	filledPortion := ""
	for range filledWidth {
		filledPortion += filledChar
	}

	// Build empty portion
	emptyPortion := ""
	emptyWidth := progressBarWidth - filledWidth
	for range emptyWidth {
		emptyPortion += emptyChar
	}

	// Apply styling and combine
	styledFilled := filledStyle().Render(filledPortion)
	styledEmpty := emptyStyle().Render(emptyPortion)

	return fmt.Sprintf("[%s%s] %d%%", styledFilled, styledEmpty, percentage)
}
