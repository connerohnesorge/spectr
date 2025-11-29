package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultTheme verifies that DefaultTheme returns correct default values.
func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()

	tests := []struct {
		name     string
		field    *string
		expected string
	}{
		{"Accent", theme.Accent, "99"},
		{"Error", theme.Error, "1"},
		{"Success", theme.Success, "2"},
		{"Border", theme.Border, "240"},
		{"Help", theme.Help, "240"},
		{"Selected", theme.Selected, "229"},
		{"Highlight", theme.Highlight, "57"},
		{"Header", theme.Header, "99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field == nil {
				t.Errorf("%s is nil, expected non-nil pointer", tt.name)

				return
			}
			if *tt.field != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, *tt.field, tt.expected)
			}
		})
	}
}

// TestValidateColor verifies color validation logic.
func TestValidateColor(t *testing.T) {
	tests := []struct {
		name    string
		color   string
		wantErr bool
	}{
		// Valid ANSI codes
		{"Valid ANSI 0", "0", false},
		{"Valid ANSI 240", "240", false},
		{"Valid ANSI 255", "255", false},
		{"Valid ANSI 1", "1", false},
		{"Valid ANSI 99", "99", false},

		// Invalid ANSI codes
		{"Invalid ANSI -1", "-1", true},
		{"Invalid ANSI 256", "256", true},
		{"Invalid ANSI 999", "999", true},

		// Valid hex codes
		{"Valid hex 6 digits uppercase", "#FF5733", false},
		{"Valid hex 3 digits uppercase", "#FFF", false},
		{"Valid hex 6 digits lowercase", "#abc123", false},
		{"Valid hex 3 digits mixed", "#ABC", false},
		{"Valid hex 6 digits mixed", "#AbC123", false},

		// Invalid values
		{"Invalid string", "not-a-color", true},
		{"Empty string", "", true},
		{"Invalid hex letters", "#GGG", true},
		{"Invalid hex 5 digits", "#12345", true},
		{"Invalid hex 7 digits", "#1234567", true},
		{"Hex without hash", "FF5733", true},
		{"Invalid hex 4 digits", "#ABCD", true},
		{"Invalid hex 2 digits", "#AB", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColor(tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColor(%q) error = %v, wantErr %v", tt.color, err, tt.wantErr)
			}
		})
	}
}

// TestValidateTheme verifies theme validation logic.
func TestValidateTheme(t *testing.T) {
	tests := []struct {
		name      string
		theme     Theme
		wantErrs  int
		checkMsgs []string
	}{
		{
			name:     "All valid colors",
			theme:    DefaultTheme(),
			wantErrs: 0,
		},
		{
			name: "All nil fields",
			theme: Theme{
				Accent:    nil,
				Error:     nil,
				Success:   nil,
				Border:    nil,
				Help:      nil,
				Selected:  nil,
				Highlight: nil,
				Header:    nil,
			},
			wantErrs: 0,
		},
		{
			name: "One invalid color",
			theme: Theme{
				Accent: stringPtr("99"),
				Error:  stringPtr("invalid"),
			},
			wantErrs:  1,
			checkMsgs: []string{"theme.error"},
		},
		{
			name: "Multiple invalid colors",
			theme: Theme{
				Accent:  stringPtr("999"),
				Error:   stringPtr("invalid"),
				Success: stringPtr("#GGGGGG"),
			},
			wantErrs:  3,
			checkMsgs: []string{"theme.accent", "theme.error", "theme.success"},
		},
		{
			name: "Valid hex colors",
			theme: Theme{
				Accent:    stringPtr("#FF5733"),
				Error:     stringPtr("#F00"),
				Success:   stringPtr("#00FF00"),
				Border:    stringPtr("#abc"),
				Help:      stringPtr("#123456"),
				Selected:  stringPtr("#ABC123"),
				Highlight: stringPtr("#fff"),
				Header:    stringPtr("#000000"),
			},
			wantErrs: 0,
		},
		{
			name: "Mixed valid and invalid",
			theme: Theme{
				Accent:  stringPtr("99"),
				Error:   stringPtr("#GGG"),
				Success: stringPtr("2"),
				Border:  stringPtr("256"),
			},
			wantErrs:  2,
			checkMsgs: []string{"theme.error", "theme.border"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateTheme(tt.theme)
			if len(errs) != tt.wantErrs {
				t.Errorf("ValidateTheme() returned %d errors, want %d", len(errs), tt.wantErrs)
				for _, err := range errs {
					t.Logf("  error: %v", err)
				}
			}

			// Check that specific error messages are present
			for _, msg := range tt.checkMsgs {
				found := false
				for _, err := range errs {
					if err != nil && contains(err.Error(), msg) {
						found = true

						break
					}
				}
				if !found {
					t.Errorf("Expected error containing %q, but not found in errors", msg)
				}
			}
		})
	}
}

// TestLoadFromPath_FileDoesNotExist verifies default config is returned when file doesn't exist.
func TestLoadFromPath_FileDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, "config.yaml")

	cfg, err := LoadFromPath(nonExistentPath)
	if err != nil {
		t.Fatalf("LoadFromPath() with non-existent file returned error: %v", err)
	}

	defaults := DefaultTheme()
	if cfg == nil {
		t.Fatal("LoadFromPath() returned nil config")
	}

	// Verify defaults are returned
	if cfg.Theme.Accent == nil || *cfg.Theme.Accent != *defaults.Accent {
		t.Errorf("Expected default Accent, got %v", cfg.Theme.Accent)
	}
}

// TestLoadFromPath_ValidYAMLFull verifies loading a complete valid YAML config.
func TestLoadFromPath_ValidYAMLFull(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `theme:
  accent: "200"
  error: "1"
  success: "2"
  border: "240"
  help: "8"
  selected: "229"
  highlight: "57"
  header: "99"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadFromPath() returned nil config")
	}

	// Verify loaded values
	if cfg.Theme.Accent == nil || *cfg.Theme.Accent != "200" {
		t.Errorf("Accent = %v, want '200'", cfg.Theme.Accent)
	}
	if cfg.Theme.Help == nil || *cfg.Theme.Help != "8" {
		t.Errorf("Help = %v, want '8'", cfg.Theme.Help)
	}
}

// TestLoadFromPath_ValidYAMLPartial verifies partial config merges with defaults.
func TestLoadFromPath_ValidYAMLPartial(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `theme:
  accent: "#FF5733"
  error: "9"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadFromPath() returned nil config")
	}

	defaults := DefaultTheme()

	// Verify custom values
	if cfg.Theme.Accent == nil || *cfg.Theme.Accent != "#FF5733" {
		t.Errorf("Accent = %v, want '#FF5733'", cfg.Theme.Accent)
	}
	if cfg.Theme.Error == nil || *cfg.Theme.Error != "9" {
		t.Errorf("Error = %v, want '9'", cfg.Theme.Error)
	}

	// Verify defaults for unspecified fields
	if cfg.Theme.Success == nil || *cfg.Theme.Success != *defaults.Success {
		t.Errorf("Success = %v, want default %v", cfg.Theme.Success, *defaults.Success)
	}
	if cfg.Theme.Border == nil || *cfg.Theme.Border != *defaults.Border {
		t.Errorf("Border = %v, want default %v", cfg.Theme.Border, *defaults.Border)
	}
}

// TestLoadFromPath_InvalidYAML verifies that invalid YAML returns defaults with no error.
func TestLoadFromPath_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `theme:
  accent: "99"
  this is not valid yaml [[[
`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() with invalid YAML returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadFromPath() returned nil config")
	}

	// Verify defaults are returned
	defaults := DefaultTheme()
	if cfg.Theme.Accent == nil || *cfg.Theme.Accent != *defaults.Accent {
		t.Errorf("Expected default Accent after invalid YAML, got %v", cfg.Theme.Accent)
	}
}

// TestLoadFromPath_InvalidColors verifies that invalid colors are reset to defaults.
func TestLoadFromPath_InvalidColors(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `theme:
  accent: "999"
  error: "invalid"
  success: "2"
  border: "#GGG"
  help: "240"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadFromPath() returned nil config")
	}

	defaults := DefaultTheme()

	// Invalid fields should be reset to defaults
	if cfg.Theme.Accent == nil || *cfg.Theme.Accent != *defaults.Accent {
		t.Errorf("Accent (invalid) = %v, want default %v", cfg.Theme.Accent, *defaults.Accent)
	}
	if cfg.Theme.Error == nil || *cfg.Theme.Error != *defaults.Error {
		t.Errorf("Error (invalid) = %v, want default %v", cfg.Theme.Error, *defaults.Error)
	}
	if cfg.Theme.Border == nil || *cfg.Theme.Border != *defaults.Border {
		t.Errorf("Border (invalid) = %v, want default %v", cfg.Theme.Border, *defaults.Border)
	}

	// Valid fields should be preserved
	if cfg.Theme.Success == nil || *cfg.Theme.Success != "2" {
		t.Errorf("Success (valid) = %v, want '2'", cfg.Theme.Success)
	}
	if cfg.Theme.Help == nil || *cfg.Theme.Help != "240" {
		t.Errorf("Help (valid) = %v, want '240'", cfg.Theme.Help)
	}
}

// TestConfigPath verifies ConfigPath follows XDG spec.
func TestConfigPath(t *testing.T) {
	// Save original environment
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if err := os.Setenv("XDG_CONFIG_HOME", originalXDG); err != nil {
			t.Errorf("Failed to restore XDG_CONFIG_HOME: %v", err)
		}
	}()

	t.Run("XDG_CONFIG_HOME set", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.Setenv("XDG_CONFIG_HOME", tmpDir); err != nil {
			t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
		}

		path := ConfigPath()
		expected := filepath.Join(tmpDir, "spectr", "config.yaml")

		if path != expected {
			t.Errorf(
				"ConfigPath() with XDG_CONFIG_HOME = %q, got %q, want %q",
				tmpDir,
				path,
				expected,
			)
		}
	})

	t.Run("XDG_CONFIG_HOME not set", func(t *testing.T) {
		if err := os.Unsetenv("XDG_CONFIG_HOME"); err != nil {
			t.Fatalf("Failed to unset XDG_CONFIG_HOME: %v", err)
		}

		path := ConfigPath()

		// Should contain .config/spectr/config.yaml
		if !contains(path, filepath.Join(".config", "spectr", "config.yaml")) {
			t.Errorf(
				"ConfigPath() without XDG_CONFIG_HOME = %q, expected to contain '.config/spectr/config.yaml'",
				path,
			)
		}
	})
}

// TestLoad verifies Load uses the default config path.
func TestLoad(t *testing.T) {
	// This test is tricky because Load() uses ConfigPath() which depends on environment.
	// We'll just verify it doesn't crash and returns a valid config.
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Should have all theme fields set (either from file or defaults)
	if cfg.Theme.Accent == nil {
		t.Error("Load() returned config with nil Accent")
	}
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
