// Package theme provides color theming functionality for Spectr CLI.
package theme

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
)

// Theme defines a complete color palette for the Spectr CLI.
type Theme struct {
	Primary       lipgloss.Color // Main accent - headers, titles
	Secondary     lipgloss.Color // Secondary accent - cursors, selections
	Success       lipgloss.Color // Success states, checkmarks
	Error         lipgloss.Color // Errors
	Warning       lipgloss.Color // Caution indicators
	Muted         lipgloss.Color // Dim/subtle text
	Border        lipgloss.Color // Table borders, separators
	Header        lipgloss.Color // Section headers
	Selected      lipgloss.Color // Selected item foreground
	Highlight     lipgloss.Color // Selected item background
	GradientStart lipgloss.Color // ASCII art gradient start
	GradientEnd   lipgloss.Color // ASCII art gradient end
}

// Default theme matching current hardcoded colors in the codebase
var defaultTheme = &Theme{
	Primary:       lipgloss.Color("99"),  // Purple/violet for headers/titles
	Secondary:     lipgloss.Color("170"), // Pink for selections
	Success:       lipgloss.Color("42"),  // Green
	Error:         lipgloss.Color("196"), // Red
	Warning:       lipgloss.Color("3"),   // Yellow
	Muted:         lipgloss.Color("240"), // Dim gray
	Border:        lipgloss.Color("240"), // Dim gray
	Header:        lipgloss.Color("99"),  // Purple
	Selected:      lipgloss.Color("229"), // Light yellow foreground
	Highlight:     lipgloss.Color("57"),  // Purple background
	GradientStart: lipgloss.Color("99"),  // Purple
	GradientEnd:   lipgloss.Color("205"), // Pink
}

// Dark theme: high contrast on dark backgrounds, brighter colors
var darkTheme = &Theme{
	Primary:       lipgloss.Color("141"), // Bright purple
	Secondary:     lipgloss.Color("213"), // Bright pink
	Success:       lipgloss.Color("46"),  // Bright green
	Error:         lipgloss.Color("196"), // Bright red
	Warning:       lipgloss.Color("226"), // Bright yellow
	Muted:         lipgloss.Color("243"), // Medium gray
	Border:        lipgloss.Color("238"), // Dark gray border
	Header:        lipgloss.Color("141"), // Bright purple
	Selected:      lipgloss.Color("231"), // White foreground
	Highlight:     lipgloss.Color("61"),  // Bright purple background
	GradientStart: lipgloss.Color("141"), // Bright purple
	GradientEnd:   lipgloss.Color("213"), // Bright pink
}

// Light theme: optimized for light terminal backgrounds, darker accents
var lightTheme = &Theme{
	Primary:       lipgloss.Color("55"),  // Dark purple
	Secondary:     lipgloss.Color("125"), // Dark pink
	Success:       lipgloss.Color("28"),  // Dark green
	Error:         lipgloss.Color("160"), // Dark red
	Warning:       lipgloss.Color("136"), // Dark yellow/orange
	Muted:         lipgloss.Color("246"), // Light gray
	Border:        lipgloss.Color("250"), // Very light gray border
	Header:        lipgloss.Color("55"),  // Dark purple
	Selected:      lipgloss.Color("16"),  // Black foreground
	Highlight:     lipgloss.Color("189"), // Light purple background
	GradientStart: lipgloss.Color("55"),  // Dark purple
	GradientEnd:   lipgloss.Color("125"), // Dark pink
}

// Solarized theme: Solarized Dark palette colors
var solarizedTheme = &Theme{
	Primary:       lipgloss.Color("33"),  // Blue (base0)
	Secondary:     lipgloss.Color("125"), // Magenta
	Success:       lipgloss.Color("64"),  // Green
	Error:         lipgloss.Color("160"), // Red
	Warning:       lipgloss.Color("136"), // Yellow
	Muted:         lipgloss.Color("240"), // Base01
	Border:        lipgloss.Color("235"), // Base02
	Header:        lipgloss.Color("37"),  // Cyan
	Selected:      lipgloss.Color("230"), // Base3 (light)
	Highlight:     lipgloss.Color("235"), // Base02 (dark)
	GradientStart: lipgloss.Color("33"),  // Blue
	GradientEnd:   lipgloss.Color("125"), // Magenta
}

// Monokai theme: Monokai palette colors
var monokaiTheme = &Theme{
	Primary:       lipgloss.Color("141"), // Purple
	Secondary:     lipgloss.Color("197"), // Pink
	Success:       lipgloss.Color("148"), // Green
	Error:         lipgloss.Color("197"), // Pink/red
	Warning:       lipgloss.Color("208"), // Orange
	Muted:         lipgloss.Color("243"), // Gray
	Border:        lipgloss.Color("237"), // Dark gray
	Header:        lipgloss.Color("81"),  // Cyan/blue
	Selected:      lipgloss.Color("231"), // White
	Highlight:     lipgloss.Color("237"), // Dark gray background
	GradientStart: lipgloss.Color("141"), // Purple
	GradientEnd:   lipgloss.Color("197"), // Pink
}

// themes is the registry of all available themes
var themes = map[string]*Theme{
	"default":   defaultTheme,
	"dark":      darkTheme,
	"light":     lightTheme,
	"solarized": solarizedTheme,
	"monokai":   monokaiTheme,
}

// current holds the currently active theme
var current *Theme

// Get returns the theme with the given name.
// Returns an error if the theme does not exist.
func Get(name string) (*Theme, error) {
	theme, ok := themes[name]
	if !ok {
		return nil, fmt.Errorf("theme not found: %s", name)
	}

	return theme, nil
}

// Load loads the theme with the given name as the current theme.
// Returns an error if the theme does not exist.
func Load(name string) error {
	theme, err := Get(name)
	if err != nil {
		return err
	}
	current = theme

	return nil
}

// Current returns the currently active theme.
// If no theme has been loaded, returns the default theme.
func Current() *Theme {
	if current == nil {
		return defaultTheme
	}

	return current
}

// Available returns a sorted list of all available theme names.
func Available() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}
