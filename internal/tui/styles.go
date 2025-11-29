// Package tui provides shared TUI components for Spectr CLI interactive modes.
package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/connerohnesorge/spectr/internal/config"
)

// Color constants used across TUI components
// Deprecated: Use getter functions (GetBorderColor, GetHeaderColor,
// etc.) instead. These constants are kept for backwards compatibility
// but will be removed in a future version.
const (
	ColorBorder    = "240"
	ColorHeader    = "99"
	ColorSelected  = "229"
	ColorHighlight = "57"
	ColorHelp      = "240"
)

// currentTheme holds the loaded theme configuration.
// It is initialized in the init() function and used by all getter functions.
var currentTheme *config.Theme

// init loads the theme configuration from the user's config file.
// If loading fails, it falls back to the default theme.
func init() {
	cfg, err := config.Load()
	if err != nil || cfg == nil {
		// Fall back to defaults if config loading fails
		defaultTheme := config.DefaultTheme()
		currentTheme = &defaultTheme
	} else {
		currentTheme = &cfg.Theme
	}
}

// GetBorderColor returns the configured border color.
func GetBorderColor() string {
	if currentTheme == nil || currentTheme.Border == nil {
		return ColorBorder
	}

	return *currentTheme.Border
}

// GetHeaderColor returns the configured header color.
func GetHeaderColor() string {
	if currentTheme == nil || currentTheme.Header == nil {
		return ColorHeader
	}

	return *currentTheme.Header
}

// GetSelectedColor returns the configured selected item foreground color.
func GetSelectedColor() string {
	if currentTheme == nil || currentTheme.Selected == nil {
		return ColorSelected
	}

	return *currentTheme.Selected
}

// GetHighlightColor returns the configured highlight
// (selected background) color.
func GetHighlightColor() string {
	if currentTheme == nil || currentTheme.Highlight == nil {
		return ColorHighlight
	}

	return *currentTheme.Highlight
}

// GetHelpColor returns the configured help text color.
func GetHelpColor() string {
	if currentTheme == nil || currentTheme.Help == nil {
		return ColorHelp
	}

	return *currentTheme.Help
}

// GetAccentColor returns the configured accent color.
func GetAccentColor() string {
	if currentTheme == nil || currentTheme.Accent == nil {
		return ColorHeader // Fall back to header color for accent
	}

	return *currentTheme.Accent
}

// GetErrorColor returns the configured error color.
func GetErrorColor() string {
	if currentTheme == nil || currentTheme.Error == nil {
		return "1" // Default red
	}

	return *currentTheme.Error
}

// GetSuccessColor returns the configured success color.
func GetSuccessColor() string {
	if currentTheme == nil || currentTheme.Success == nil {
		return "2" // Default green
	}

	return *currentTheme.Success
}

// ApplyTableStyles applies the default Spectr styling to a table.
func ApplyTableStyles(t *table.Model) {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(GetBorderColor())).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(GetHeaderColor()))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(GetSelectedColor())).
		Background(lipgloss.Color(GetHighlightColor())).
		Bold(true)

	t.SetStyles(s)
}

// TitleStyle returns the style for titles.
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(GetHeaderColor())).
		MarginBottom(1)
}

// HelpStyle returns the style for help text.
func HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(GetHelpColor())).
		MarginTop(1)
}

// SelectedStyle returns the style for selected items.
func SelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(GetSelectedColor())).
		Background(lipgloss.Color(GetHighlightColor())).
		Bold(true).
		PaddingLeft(2)
}

// ChoiceStyle returns the style for unselected choices.
func ChoiceStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		PaddingLeft(2)
}
