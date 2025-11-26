// Package tui provides shared TUI components for Spectr CLI interactive modes.
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Color constants used across TUI components
const (
	ColorBorder    = "240"
	ColorHeader    = "99"
	ColorSelected  = "229"
	ColorHighlight = "57"
	ColorHelp      = "240"
)

// ApplyTableStyles applies the default Spectr styling to a table.
func ApplyTableStyles(t *table.Model) {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorBorder)).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(ColorHeader))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(ColorSelected)).
		Background(lipgloss.Color(ColorHighlight)).
		Bold(true)

	t.SetStyles(s)
}

// TitleStyle returns the style for titles.
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorHeader)).
		MarginBottom(1)
}

// HelpStyle returns the style for help text.
func HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorHelp)).
		MarginTop(1)
}

// SelectedStyle returns the style for selected items.
func SelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSelected)).
		Background(lipgloss.Color(ColorHighlight)).
		Bold(true).
		PaddingLeft(2)
}

// ChoiceStyle returns the style for unselected choices.
func ChoiceStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		PaddingLeft(2)
}
