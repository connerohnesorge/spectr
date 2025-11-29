// Package config provides user configuration types and defaults for Spectr CLI.
package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const maxANSIColorCode = 255

var errEmptyColorValue = errors.New("color value cannot be empty")

// Config represents the full user configuration for Spectr.
// This structure maps to the config.yaml file located at
// $XDG_CONFIG_HOME/spectr/config.yaml.
type Config struct {
	// Theme contains color customization settings for TUI elements.
	// If nil or partially set, defaults are used for missing values.
	Theme Theme `yaml:"theme"`
}

// Theme contains color settings for TUI elements.
// Each field accepts either:
//   - ANSI 256 color codes as strings (e.g., "240" for gray)
//   - Hex color codes (e.g., "#FF5733" for orange-red)
//
// All fields are pointers to distinguish between "unset" (nil) and
// "set to empty" ("").
// When a field is nil, the default color is used.
type Theme struct {
	// Accent is the primary accent color used for highlights and emphasis.
	Accent *string `yaml:"accent,omitempty"`

	// Error is the color used for error messages and warnings.
	Error *string `yaml:"error,omitempty"`

	// Success is the color used for success messages and confirmations.
	Success *string `yaml:"success,omitempty"`

	// Border is the color used for table borders and separators.
	Border *string `yaml:"border,omitempty"`

	// Help is the color used for help text and secondary information.
	Help *string `yaml:"help,omitempty"`

	// Selected is the foreground color for selected/highlighted items.
	Selected *string `yaml:"selected,omitempty"`

	// Highlight is the background color for selected/highlighted items.
	Highlight *string `yaml:"highlight,omitempty"`

	// Header is the color used for headers and titles.
	Header *string `yaml:"header,omitempty"`
}

// DefaultTheme returns a Theme with default colors matching the current
// hardcoded values used throughout the TUI components.
//
// These defaults are:
//   - border: "240" (gray)
//   - header: "99" (purple)
//   - selected: "229" (light yellow)
//   - highlight: "57" (blue)
//   - help: "240" (gray)
//   - accent: "99" (purple, matching header for consistency)
//   - error: "1" (red)
//   - success: "2" (green)
func DefaultTheme() Theme {
	return Theme{
		Accent:    stringPtr("99"),
		Error:     stringPtr("1"),
		Success:   stringPtr("2"),
		Border:    stringPtr("240"),
		Help:      stringPtr("240"),
		Selected:  stringPtr("229"),
		Highlight: stringPtr("57"),
		Header:    stringPtr("99"),
	}
}

// stringPtr is a helper function to create a string pointer.
func stringPtr(s string) *string {
	return &s
}

// ValidateColor validates a single color value.
// It accepts:
//   - ANSI 256 color codes: strings representing integers 0-255
//     (e.g., "0", "240", "255")
//   - Hex color codes: strings like "#RRGGBB" or "#RGB"
//     (3 or 6 hex digits, case-insensitive)
//
// Returns an error with a descriptive message if the format is invalid.
func ValidateColor(value string) error {
	if value == "" {
		return errEmptyColorValue
	}

	// Check for hex color format: #RGB or #RRGGBB
	hexPattern := regexp.MustCompile(`^#[0-9A-Fa-f]{3}$|^#[0-9A-Fa-f]{6}$`)
	if hexPattern.MatchString(value) {
		return nil
	}

	// Check for ANSI 256 color code: 0-maxANSIColorCode
	num, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf(
			"invalid color format: must be ANSI 256 code (0-%d) "+
				"or hex code (#RGB or #RRGGBB), got %q",
			maxANSIColorCode,
			value,
		)
	}
	if num < 0 || num > maxANSIColorCode {
		return fmt.Errorf(
			"ANSI color code must be between 0 and %d, got %d",
			maxANSIColorCode,
			num,
		)
	}

	return nil
}

// fieldValidationError is a helper struct to track which field
// had a validation error.
type fieldValidationError struct {
	fieldName string
	err       error
}

// ValidateTheme validates all set theme colors.
// It checks each non-nil field with ValidateColor and returns a slice
// of errors (one per invalid field) with field names included in the
// error messages. Returns an empty slice if all colors are valid.
func ValidateTheme(theme Theme) []error {
	var errs []error

	// Helper to validate a field if it's set
	validateField := func(fieldName string, value *string) {
		if value == nil {
			return
		}
		if err := ValidateColor(*value); err != nil {
			errs = append(
				errs,
				fmt.Errorf("theme.%s: %w", fieldName, err),
			)
		}
	}

	validateField("accent", theme.Accent)
	validateField("error", theme.Error)
	validateField("success", theme.Success)
	validateField("border", theme.Border)
	validateField("help", theme.Help)
	validateField("selected", theme.Selected)
	validateField("highlight", theme.Highlight)
	validateField("header", theme.Header)

	return errs
}

// validateThemeWithFieldTracking validates theme and returns
// field-specific errors. This is used internally by LoadFromPath to
// reset invalid fields to defaults.
func validateThemeWithFieldTracking(
	theme Theme,
) []fieldValidationError {
	var errs []fieldValidationError

	// Helper to validate a field if it's set
	validateField := func(fieldName string, value *string) {
		if value == nil {
			return
		}
		if err := ValidateColor(*value); err != nil {
			errs = append(errs, fieldValidationError{
				fieldName: fieldName,
				err:       fmt.Errorf("theme.%s: %w", fieldName, err),
			})
		}
	}

	validateField("accent", theme.Accent)
	validateField("error", theme.Error)
	validateField("success", theme.Success)
	validateField("border", theme.Border)
	validateField("help", theme.Help)
	validateField("selected", theme.Selected)
	validateField("highlight", theme.Highlight)
	validateField("header", theme.Header)

	return errs
}
