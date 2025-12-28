package view

import (
	"strings"
	"testing"
)

//nolint:revive // cognitive-complexity justified for comprehensive testing
func TestRenderBar(t *testing.T) {
	tests := []struct {
		name       string
		completed  int
		total      int
		wantFilled int // expected number of filled characters
		wantPerc   string
	}{
		{
			name:       "0% complete",
			completed:  0,
			total:      10,
			wantFilled: 0,
			wantPerc:   "0%",
		},
		{
			name:       "25% complete",
			completed:  25,
			total:      100,
			wantFilled: 5, // 0.25 * 20 = 5
			wantPerc:   "25%",
		},
		{
			name:       "50% complete",
			completed:  5,
			total:      10,
			wantFilled: 10, // 0.5 * 20 = 10
			wantPerc:   "50%",
		},
		{
			name:       "75% complete",
			completed:  75,
			total:      100,
			wantFilled: 15,
			wantPerc:   "75%",
		},
		{
			name:       "100% complete",
			completed:  10,
			total:      10,
			wantFilled: 20, // 1.0 * 20 = 20
			wantPerc:   "100%",
		},
		{
			name:       "37% complete (example from design)",
			completed:  37,
			total:      100,
			wantFilled: 7, // 0.37 * 20 = 7.4, rounds to 7
			wantPerc:   "37%",
		},
		{
			name:       "60% complete (example from design)",
			completed:  60,
			total:      100,
			wantFilled: 12, // 0.6 * 20 = 12
			wantPerc:   "60%",
		},
		{
			name:       "zero total edge case",
			completed:  0,
			total:      0,
			wantFilled: 0,
			wantPerc:   "0%",
		},
		{
			name:       "rounding test - 3/8 = 37.5% rounds to 38%",
			completed:  3,
			total:      8,
			wantFilled: 8, // 0.375 * 20 = 7.5, rounds to 8
			wantPerc:   "38%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderBar(
				tt.completed,
				tt.total,
			)

			// Strip ANSI color codes for testing
			// (lipgloss adds escape codes that we can't easily predict)
			// So we'll just verify the bar structure and percentage

			// Check that result contains the percentage
			if !strings.Contains(
				result,
				tt.wantPerc,
			) {
				t.Errorf(
					"RenderBar(%d, %d) percentage = got result not containing %q, result = %q",
					tt.completed,
					tt.total,
					tt.wantPerc,
					result,
				)
			}

			// Check that result contains brackets
			if !strings.HasPrefix(result, "[") {
				t.Errorf(
					"RenderBar(%d, %d) should start with '[', got %q",
					tt.completed,
					tt.total,
					result,
				)
			}

			// Count filled and empty characters (ignoring ANSI codes)
			// This is a simple approach - count the actual unicode characters
			filledCount := strings.Count(
				result,
				filledChar,
			)
			emptyCount := strings.Count(
				result,
				emptyChar,
			)

			if filledCount != tt.wantFilled {
				t.Errorf(
					"RenderBar(%d, %d) filled characters = %d, want %d (result: %q)",
					tt.completed,
					tt.total,
					filledCount,
					tt.wantFilled,
					result,
				)
			}

			expectedEmpty := progressBarWidth - tt.wantFilled
			if emptyCount != expectedEmpty {
				t.Errorf(
					"RenderBar(%d, %d) empty characters = %d, want %d (result: %q)",
					tt.completed,
					tt.total,
					emptyCount,
					expectedEmpty,
					result,
				)
			}

			// Total bar width should always be progressBarWidth (20)
			totalChars := filledCount + emptyCount
			if totalChars != progressBarWidth {
				t.Errorf(
					"RenderBar(%d, %d) total bar characters = %d, want %d",
					tt.completed,
					tt.total,
					totalChars,
					progressBarWidth,
				)
			}
		})
	}
}

func TestRenderBarEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		completed int
		total     int
	}{
		{"negative completed", -1, 10},
		{"completed exceeds total", 15, 10},
		{"both zero", 0, 0},
		{"large numbers", 999999, 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := RenderBar(
				tt.completed,
				tt.total,
			)

			// Should return a non-empty string
			if result == "" {
				t.Errorf(
					"RenderBar(%d, %d) returned empty string",
					tt.completed,
					tt.total,
				)
			}

			// Should contain brackets and percentage
			if !strings.Contains(result, "[") ||
				!strings.Contains(result, "%") {
				t.Errorf(
					"RenderBar(%d, %d) missing expected format, got %q",
					tt.completed,
					tt.total,
					result,
				)
			}
		})
	}
}
