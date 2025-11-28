// Package tui provides shared TUI components for Spectr CLI interactive modes.
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/connerohnesorge/spectr/internal/theme"
)

// ApplyTableStyles applies the default Spectr styling to a table.
func ApplyTableStyles(t *table.Model) {
	th := theme.Current()
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(th.Border).
		BorderBottom(true).
		Bold(true).
		Foreground(th.Header)
	s.Selected = s.Selected.
		Foreground(th.Selected).
		Background(th.Highlight).
		Bold(true)

	t.SetStyles(s)
}

// TitleStyle returns the style for titles.
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Current().Header).
		MarginBottom(1)
}

// HelpStyle returns the style for help text.
func HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Current().Muted).
		MarginTop(1)
}

// SelectedStyle returns the style for selected items.
func SelectedStyle() lipgloss.Style {
	th := theme.Current()

	return lipgloss.NewStyle().
		Foreground(th.Selected).
		Background(th.Highlight).
		Bold(true).
		PaddingLeft(2)
}

// ChoiceStyle returns the style for unselected choices.
func ChoiceStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		PaddingLeft(2)
}
