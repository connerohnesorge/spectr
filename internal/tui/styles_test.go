package tui

import (
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
)

// TestGetBorderColor verifies GetBorderColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetBorderColor(t *testing.T) {
	color := GetBorderColor()
	if color == "" {
		t.Error("GetBorderColor() returned empty string")
	}
	// Should match default when using default config
	// The init() already loaded config, so we verify it returns a valid value
	if color != ColorBorder && color != *config.DefaultTheme().Border {
		t.Errorf(
			"GetBorderColor() = %q, expected default %q or %q",
			color,
			ColorBorder,
			*config.DefaultTheme().Border,
		)
	}
}

// TestGetHeaderColor verifies GetHeaderColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetHeaderColor(t *testing.T) {
	color := GetHeaderColor()
	if color == "" {
		t.Error("GetHeaderColor() returned empty string")
	}
	// Should match default when using default config
	if color != ColorHeader && color != *config.DefaultTheme().Header {
		t.Errorf(
			"GetHeaderColor() = %q, expected default %q or %q",
			color,
			ColorHeader,
			*config.DefaultTheme().Header,
		)
	}
}

// TestGetSelectedColor verifies GetSelectedColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetSelectedColor(t *testing.T) {
	color := GetSelectedColor()
	if color == "" {
		t.Error("GetSelectedColor() returned empty string")
	}
	// Should match default when using default config
	if color != ColorSelected && color != *config.DefaultTheme().Selected {
		t.Errorf(
			"GetSelectedColor() = %q, expected default %q or %q",
			color,
			ColorSelected,
			*config.DefaultTheme().Selected,
		)
	}
}

// TestGetHighlightColor verifies GetHighlightColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetHighlightColor(t *testing.T) {
	color := GetHighlightColor()
	if color == "" {
		t.Error("GetHighlightColor() returned empty string")
	}
	// Should match default when using default config
	if color != ColorHighlight && color != *config.DefaultTheme().Highlight {
		t.Errorf(
			"GetHighlightColor() = %q, expected default %q or %q",
			color,
			ColorHighlight,
			*config.DefaultTheme().Highlight,
		)
	}
}

// TestGetHelpColor verifies GetHelpColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetHelpColor(t *testing.T) {
	color := GetHelpColor()
	if color == "" {
		t.Error("GetHelpColor() returned empty string")
	}
	// Should match default when using default config
	if color != ColorHelp && color != *config.DefaultTheme().Help {
		t.Errorf(
			"GetHelpColor() = %q, expected default %q or %q",
			color,
			ColorHelp,
			*config.DefaultTheme().Help,
		)
	}
}

// TestGetAccentColor verifies GetAccentColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetAccentColor(t *testing.T) {
	color := GetAccentColor()
	if color == "" {
		t.Error("GetAccentColor() returned empty string")
	}
	// Should match default when using default config (or fallback to header)
	defaultAccent := *config.DefaultTheme().Accent
	if color != defaultAccent && color != ColorHeader {
		t.Errorf(
			"GetAccentColor() = %q, expected default %q or fallback %q",
			color,
			defaultAccent,
			ColorHeader,
		)
	}
}

// TestGetErrorColor verifies GetErrorColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetErrorColor(t *testing.T) {
	color := GetErrorColor()
	if color == "" {
		t.Error("GetErrorColor() returned empty string")
	}
	// Should match default when using default config
	defaultError := *config.DefaultTheme().Error
	if color != defaultError && color != "1" {
		t.Errorf(
			"GetErrorColor() = %q, expected default %q or fallback %q",
			color,
			defaultError,
			"1",
		)
	}
}

// TestGetSuccessColor verifies GetSuccessColor returns a non-empty string
// and matches the default value when no user config exists.
func TestGetSuccessColor(t *testing.T) {
	color := GetSuccessColor()
	if color == "" {
		t.Error("GetSuccessColor() returned empty string")
	}
	// Should match default when using default config
	defaultSuccess := *config.DefaultTheme().Success
	if color != defaultSuccess && color != "2" {
		t.Errorf(
			"GetSuccessColor() = %q, expected default %q or fallback %q",
			color,
			defaultSuccess,
			"2",
		)
	}
}

// TestDefaultValuesMatch verifies all getters return values matching
// the const defaults when no config exists (backward compatibility).
func TestDefaultValuesMatch(t *testing.T) {
	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"BorderColor", GetBorderColor, ColorBorder},
		{"HeaderColor", GetHeaderColor, ColorHeader},
		{"SelectedColor", GetSelectedColor, ColorSelected},
		{"HighlightColor", GetHighlightColor, ColorHighlight},
		{"HelpColor", GetHelpColor, ColorHelp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getter()
			if got == "" {
				t.Errorf("%s getter returned empty string", tt.name)
			}
			// When using default config, values should match the
			// deprecated constants or the DefaultTheme values
			// (which are the same)
			if got == tt.expected {
				return
			}

			// This is acceptable if it matches the DefaultTheme value instead
			defaultTheme := config.DefaultTheme()
			var defaultValue string
			switch tt.name {
			case "BorderColor":
				defaultValue = *defaultTheme.Border
			case "HeaderColor":
				defaultValue = *defaultTheme.Header
			case "SelectedColor":
				defaultValue = *defaultTheme.Selected
			case "HighlightColor":
				defaultValue = *defaultTheme.Highlight
			case "HelpColor":
				defaultValue = *defaultTheme.Help
			}
			if got != defaultValue {
				t.Errorf(
					"%s = %q, expected const %q or default theme %q",
					tt.name,
					got,
					tt.expected,
					defaultValue,
				)
			}
		})
	}
}

// TestColorGettersWithNilTheme verifies that getters handle nil theme gracefully.
func TestColorGettersWithNilTheme(t *testing.T) {
	// Save original theme
	original := currentTheme

	// Test with nil theme
	currentTheme = nil

	tests := []struct {
		name     string
		getter   func() string
		fallback string
	}{
		{"GetBorderColor", GetBorderColor, ColorBorder},
		{"GetHeaderColor", GetHeaderColor, ColorHeader},
		{"GetSelectedColor", GetSelectedColor, ColorSelected},
		{"GetHighlightColor", GetHighlightColor, ColorHighlight},
		{"GetHelpColor", GetHelpColor, ColorHelp},
		{"GetAccentColor", GetAccentColor, ColorHeader},
		{"GetErrorColor", GetErrorColor, "1"},
		{"GetSuccessColor", GetSuccessColor, "2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getter()
			if got != tt.fallback {
				t.Errorf("%s with nil theme = %q, expected fallback %q", tt.name, got, tt.fallback)
			}
		})
	}

	// Restore original theme
	currentTheme = original
}

// TestColorGettersWithNilFields verifies that getters handle nil theme fields gracefully.
func TestColorGettersWithNilFields(t *testing.T) {
	// Save original theme
	original := currentTheme

	// Create theme with nil fields
	emptyTheme := &config.Theme{}
	currentTheme = emptyTheme

	tests := []struct {
		name     string
		getter   func() string
		fallback string
	}{
		{"GetBorderColor", GetBorderColor, ColorBorder},
		{"GetHeaderColor", GetHeaderColor, ColorHeader},
		{"GetSelectedColor", GetSelectedColor, ColorSelected},
		{"GetHighlightColor", GetHighlightColor, ColorHighlight},
		{"GetHelpColor", GetHelpColor, ColorHelp},
		{"GetAccentColor", GetAccentColor, ColorHeader},
		{"GetErrorColor", GetErrorColor, "1"},
		{"GetSuccessColor", GetSuccessColor, "2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getter()
			if got != tt.fallback {
				t.Errorf("%s with nil fields = %q, expected fallback %q", tt.name, got, tt.fallback)
			}
		})
	}

	// Restore original theme
	currentTheme = original
}

// TestColorGettersWithCustomTheme verifies that getters return custom theme values.
func TestColorGettersWithCustomTheme(t *testing.T) {
	// Save original theme
	original := currentTheme

	// Create custom theme
	customBorder := "100"
	customHeader := "200"
	customSelected := "150"
	customHighlight := "75"
	customHelp := "125"
	customAccent := "180"
	customError := "9"
	customSuccess := "10"

	customTheme := &config.Theme{
		Border:    &customBorder,
		Header:    &customHeader,
		Selected:  &customSelected,
		Highlight: &customHighlight,
		Help:      &customHelp,
		Accent:    &customAccent,
		Error:     &customError,
		Success:   &customSuccess,
	}
	currentTheme = customTheme

	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"GetBorderColor", GetBorderColor, customBorder},
		{"GetHeaderColor", GetHeaderColor, customHeader},
		{"GetSelectedColor", GetSelectedColor, customSelected},
		{"GetHighlightColor", GetHighlightColor, customHighlight},
		{"GetHelpColor", GetHelpColor, customHelp},
		{"GetAccentColor", GetAccentColor, customAccent},
		{"GetErrorColor", GetErrorColor, customError},
		{"GetSuccessColor", GetSuccessColor, customSuccess},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getter()
			if got != tt.expected {
				t.Errorf("%s with custom theme = %q, expected %q", tt.name, got, tt.expected)
			}
		})
	}

	// Restore original theme
	currentTheme = original
}
