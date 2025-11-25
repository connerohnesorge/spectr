package init

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestApplyGradient(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		colorA lipgloss.Color
		colorB lipgloss.Color
	}{
		{
			name:   "single line",
			text:   "HELLO",
			colorA: lipgloss.Color("#FF0000"),
			colorB: lipgloss.Color("#0000FF"),
		},
		{
			name:   "multi-line ascii art",
			text:   "███\n███\n███",
			colorA: lipgloss.Color("#FF00FF"),
			colorB: lipgloss.Color("#00FFFF"),
		},
		{
			name:   "empty string",
			text:   "",
			colorA: lipgloss.Color("#FF0000"),
			colorB: lipgloss.Color("#0000FF"),
		},
		{
			name:   "single character",
			text:   "X",
			colorA: lipgloss.Color("#FFFF00"),
			colorB: lipgloss.Color("#00FF00"),
		},
		{
			name:   "text with empty lines",
			text:   "TOP\n\nBOTTOM",
			colorA: lipgloss.Color("#FF0000"),
			colorB: lipgloss.Color("#0000FF"),
		},
		{
			name:   "ANSI 256 colors",
			text:   "TEST",
			colorA: lipgloss.Color("205"),
			colorB: lipgloss.Color("99"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyGradient(tt.text, tt.colorA, tt.colorB)

			// Basic sanity checks
			if tt.text == "" {
				if result != "" {
					t.Errorf("expected empty result for empty input, got: %q", result)
				}

				return
			}

			// Strip ANSI codes to get plain text (lipgloss may not output ANSI in test env)
			plainResult := stripAnsiCodes(result)
			if plainResult == "" {
				plainResult = result // No ANSI codes, use as-is
			}

			// Result should contain the original text structure (same number of lines)
			originalLines := strings.Split(tt.text, "\n")
			resultLines := strings.Split(plainResult, "\n")

			if len(resultLines) != len(originalLines) {
				t.Errorf("expected %d lines, got %d", len(originalLines), len(resultLines))
			}

			// Verify the text content is preserved
			originalText := strings.ReplaceAll(tt.text, "\n", "")
			resultText := strings.ReplaceAll(plainResult, "\n", "")

			if originalText != resultText {
				t.Errorf("expected text content %q, got %q", originalText, resultText)
			}

			// The function should always return something (even if just the original text)
			if result == "" && tt.text != "" {
				t.Error("expected non-empty result for non-empty input")
			}
		})
	}
}

func TestApplyGradientInvalidColors(t *testing.T) {
	// Test with invalid color codes - should fallback gracefully
	result := applyGradient("TEST", lipgloss.Color("invalid"), lipgloss.Color("#0000FF"))
	if result != "TEST" {
		t.Errorf("expected fallback to original text for invalid color, got: %q", result)
	}

	result = applyGradient("TEST", lipgloss.Color("#FF0000"), lipgloss.Color("also-invalid"))
	if result != "TEST" {
		t.Errorf("expected fallback to original text for invalid color, got: %q", result)
	}
}

func TestApplyGradientPreservesStructure(t *testing.T) {
	// Test that newlines and empty lines are preserved
	input := "LINE1\n\nLINE3\n\n\nLINE6"
	result := applyGradient(input, lipgloss.Color("#FF0000"), lipgloss.Color("#0000FF"))

	plainResult := stripAnsiCodes(result)
	if plainResult == "" {
		plainResult = result
	}

	if plainResult != input {
		t.Errorf("expected structure preserved.\nInput:  %q\nResult: %q", input, plainResult)
	}
}

// stripAnsiCodes removes ANSI escape sequences from a string
func stripAnsiCodes(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true

			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}

			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}
