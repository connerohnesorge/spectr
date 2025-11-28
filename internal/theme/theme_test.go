package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestGet verifies the Get function retrieves themes correctly.
func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		themeName string
		wantTheme *Theme
		wantError bool
	}{
		{
			name:      "get default theme",
			themeName: "default",
			wantTheme: defaultTheme,
			wantError: false,
		},
		{
			name:      "get dark theme",
			themeName: "dark",
			wantTheme: darkTheme,
			wantError: false,
		},
		{
			name:      "get light theme",
			themeName: "light",
			wantTheme: lightTheme,
			wantError: false,
		},
		{
			name:      "get solarized theme",
			themeName: "solarized",
			wantTheme: solarizedTheme,
			wantError: false,
		},
		{
			name:      "get monokai theme",
			themeName: "monokai",
			wantTheme: monokaiTheme,
			wantError: false,
		},
		{
			name:      "get nonexistent theme",
			themeName: "nonexistent",
			wantTheme: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.themeName)
			if (err != nil) != tt.wantError {
				t.Errorf("Get(%q) error = %v, wantError %v", tt.themeName, err, tt.wantError)

				return
			}
			if got != tt.wantTheme {
				t.Errorf("Get(%q) = %v, want %v", tt.themeName, got, tt.wantTheme)
			}
		})
	}
}

// TestLoad verifies the Load function sets the current theme correctly.
func TestLoad(t *testing.T) {
	// Reset current theme before tests
	current = nil

	tests := []struct {
		name      string
		themeName string
		wantError bool
	}{
		{
			name:      "load default theme",
			themeName: "default",
			wantError: false,
		},
		{
			name:      "load dark theme",
			themeName: "dark",
			wantError: false,
		},
		{
			name:      "load nonexistent theme",
			themeName: "nonexistent",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Load(tt.themeName)
			if (err != nil) != tt.wantError {
				t.Errorf("Load(%q) error = %v, wantError %v", tt.themeName, err, tt.wantError)

				return
			}

			// If load succeeded, verify current theme is set
			if tt.wantError {
				return
			}

			expectedTheme, _ := Get(tt.themeName)
			if current != expectedTheme {
				t.Errorf(
					"After Load(%q), current = %v, want %v",
					tt.themeName,
					current,
					expectedTheme,
				)
			}
		})
	}

	// Reset current theme after tests
	current = nil
}

// TestCurrent verifies the Current function returns the correct theme.
func TestCurrent(t *testing.T) {
	// Reset current theme
	current = nil

	// Test 1: Current() returns defaultTheme when no theme has been loaded
	t.Run("returns default theme when none loaded", func(t *testing.T) {
		got := Current()
		if got != defaultTheme {
			t.Errorf("Current() = %v, want %v", got, defaultTheme)
		}
	})

	// Test 2: After Load("dark"), Current() returns the dark theme
	t.Run("returns dark theme after loading", func(t *testing.T) {
		err := Load("dark")
		if err != nil {
			t.Fatalf("Load(\"dark\") failed: %v", err)
		}

		got := Current()
		if got != darkTheme {
			t.Errorf("After Load(\"dark\"), Current() = %v, want %v", got, darkTheme)
		}
	})

	// Reset current theme after tests
	current = nil
}

// TestAvailable verifies the Available function returns all theme names sorted.
func TestAvailable(t *testing.T) {
	got := Available()

	// Expected themes in sorted order
	expected := []string{"dark", "default", "light", "monokai", "solarized"}

	// Test length
	if len(got) != len(expected) {
		t.Errorf("Available() returned %d themes, want %d", len(got), len(expected))
	}

	// Test contents and order
	for i, name := range expected {
		if i >= len(got) {
			t.Errorf("Available() missing theme at index %d: %s", i, name)

			continue
		}
		if got[i] != name {
			t.Errorf("Available()[%d] = %q, want %q", i, got[i], name)
		}
	}

	// Verify specific themes are present
	themeSet := make(map[string]bool)
	for _, name := range got {
		themeSet[name] = true
	}

	requiredThemes := []string{"default", "dark", "light", "monokai", "solarized"}
	for _, name := range requiredThemes {
		if !themeSet[name] {
			t.Errorf("Available() missing required theme: %q", name)
		}
	}
}

// TestDefaultThemeColors verifies the default theme has expected color values.
func TestDefaultThemeColors(t *testing.T) {
	tests := []struct {
		name  string
		got   lipgloss.Color
		want  lipgloss.Color
		field string
	}{
		{
			name:  "Primary color",
			got:   defaultTheme.Primary,
			want:  lipgloss.Color("99"),
			field: "Primary",
		},
		{
			name:  "Header color",
			got:   defaultTheme.Header,
			want:  lipgloss.Color("99"),
			field: "Header",
		},
		{
			name:  "Border color",
			got:   defaultTheme.Border,
			want:  lipgloss.Color("240"),
			field: "Border",
		},
		{
			name:  "Secondary color",
			got:   defaultTheme.Secondary,
			want:  lipgloss.Color("170"),
			field: "Secondary",
		},
		{
			name:  "Success color",
			got:   defaultTheme.Success,
			want:  lipgloss.Color("42"),
			field: "Success",
		},
		{
			name:  "Error color",
			got:   defaultTheme.Error,
			want:  lipgloss.Color("196"),
			field: "Error",
		},
		{
			name:  "Warning color",
			got:   defaultTheme.Warning,
			want:  lipgloss.Color("3"),
			field: "Warning",
		},
		{
			name:  "Muted color",
			got:   defaultTheme.Muted,
			want:  lipgloss.Color("240"),
			field: "Muted",
		},
		{
			name:  "Selected color",
			got:   defaultTheme.Selected,
			want:  lipgloss.Color("229"),
			field: "Selected",
		},
		{
			name:  "Highlight color",
			got:   defaultTheme.Highlight,
			want:  lipgloss.Color("57"),
			field: "Highlight",
		},
		{
			name:  "GradientStart color",
			got:   defaultTheme.GradientStart,
			want:  lipgloss.Color("99"),
			field: "GradientStart",
		},
		{
			name:  "GradientEnd color",
			got:   defaultTheme.GradientEnd,
			want:  lipgloss.Color("205"),
			field: "GradientEnd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("defaultTheme.%s = %q, want %q", tt.field, tt.got, tt.want)
			}
		})
	}
}
