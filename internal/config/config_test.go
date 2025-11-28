package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// Create a temporary directory with no config file
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf("Expected RootDir=%q, got %q", DefaultRootDir, cfg.RootDir)
	}

	if cfg.ConfigPath != "" {
		t.Errorf("Expected empty ConfigPath for default config, got %q", cfg.ConfigPath)
	}

	// ProjectRoot should be the temp directory (absolute)
	absPath, _ := filepath.Abs(tmpDir)
	if cfg.ProjectRoot != absPath {
		t.Errorf("Expected ProjectRoot=%q, got %q", absPath, cfg.ProjectRoot)
	}
}

func TestLoad_CustomRootDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create spectr.yaml with custom root_dir
	configContent := "root_dir: my-specs\n"
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != "my-specs" {
		t.Errorf("Expected RootDir=%q, got %q", "my-specs", cfg.RootDir)
	}

	if cfg.ConfigPath != configPath {
		t.Errorf("Expected ConfigPath=%q, got %q", configPath, cfg.ConfigPath)
	}

	expectedRoot := filepath.Join(tmpDir, "my-specs")
	if cfg.SpectrRoot() != expectedRoot {
		t.Errorf("Expected SpectrRoot=%q, got %q", expectedRoot, cfg.SpectrRoot())
	}
}

func TestLoad_DiscoveryFromNestedDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}

	// Create config at root
	configContent := "root_dir: custom-root\n"
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Load from nested directory
	cfg, err := Load(nestedDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != "custom-root" {
		t.Errorf("Expected RootDir=%q, got %q", "custom-root", cfg.RootDir)
	}

	// ProjectRoot should be tmpDir (where config was found)
	if cfg.ProjectRoot != tmpDir {
		t.Errorf("Expected ProjectRoot=%q, got %q", tmpDir, cfg.ProjectRoot)
	}

	if cfg.ConfigPath != configPath {
		t.Errorf("Expected ConfigPath=%q, got %q", configPath, cfg.ConfigPath)
	}
}

func TestLoad_NearestConfigWins(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory
	nestedDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Create config at root
	rootConfig := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(rootConfig, []byte("root_dir: root-config\n"), 0644); err != nil {
		t.Fatalf("Failed to create root config: %v", err)
	}

	// Create config in subdir (should win)
	nestedConfig := filepath.Join(nestedDir, ConfigFileName)
	if err := os.WriteFile(nestedConfig, []byte("root_dir: nested-config\n"), 0644); err != nil {
		t.Fatalf("Failed to create nested config: %v", err)
	}

	// Load from nested directory - should use nearest config
	cfg, err := Load(nestedDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != "nested-config" {
		t.Errorf(
			"Expected nearest config to win with RootDir=%q, got %q",
			"nested-config",
			cfg.RootDir,
		)
	}

	if cfg.ProjectRoot != nestedDir {
		t.Errorf("Expected ProjectRoot=%q, got %q", nestedDir, cfg.ProjectRoot)
	}
}

func TestLoad_InvalidRootDir_PathSeparator(t *testing.T) {
	tests := []struct {
		name    string
		rootDir string
	}{
		{"forward slash", "path/to/specs"},
		{"backward slash", "path\\to\\specs"},
		{"double dots", "../specs"},
		{"asterisk", "specs*"},
		{"question mark", "specs?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			configContent := "root_dir: " + tt.rootDir + "\n"
			configPath := filepath.Join(tmpDir, ConfigFileName)
			if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
				t.Fatalf("Failed to create config file: %v", err)
			}

			_, err := Load(tmpDir)
			if err == nil {
				t.Errorf("Expected error for invalid root_dir %q, got nil", tt.rootDir)
			}

			if !strings.Contains(err.Error(), ErrInvalidRootDir.Error()) {
				t.Errorf(
					"Expected error to contain %q, got %q",
					ErrInvalidRootDir.Error(),
					err.Error(),
				)
			}

			// Just verify some invalid character was mentioned - don't check which one
			// since strings like "../" contain multiple invalid characters
			if !strings.Contains(err.Error(), "found") {
				t.Errorf("Expected error to mention which character was found, got %q", err.Error())
			}
		})
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create malformed YAML
	configContent := "root_dir: [\ninvalid yaml\n"
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	_, err := Load(tmpDir)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	// Error should mention YAML or syntax
	errMsg := strings.ToLower(err.Error())
	if !strings.Contains(errMsg, "yaml") && !strings.Contains(errMsg, "syntax") {
		t.Errorf("Expected YAML/syntax error, got: %v", err)
	}
}

func TestLoad_EmptyRootDir_UsesDefault(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config with empty root_dir
	configContent := "root_dir: \n"
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf("Expected empty root_dir to use default %q, got %q", DefaultRootDir, cfg.RootDir)
	}
}

func TestLoad_NoRootDirField_UsesDefault(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config without root_dir field
	configContent := "# Just a comment\n"
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf("Expected missing root_dir to use default %q, got %q", DefaultRootDir, cfg.RootDir)
	}
}

func TestConfig_HelperMethods(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &Config{
		RootDir:     "my-spectr",
		ProjectRoot: tmpDir,
		ConfigPath:  filepath.Join(tmpDir, ConfigFileName),
	}

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "SpectrRoot",
			method:   cfg.SpectrRoot,
			expected: filepath.Join(tmpDir, "my-spectr"),
		},
		{
			name:     "ChangesDir",
			method:   cfg.ChangesDir,
			expected: filepath.Join(tmpDir, "my-spectr", "changes"),
		},
		{
			name:     "SpecsDir",
			method:   cfg.SpecsDir,
			expected: filepath.Join(tmpDir, "my-spectr", "specs"),
		},
		{
			name:     "ArchiveDir",
			method:   cfg.ArchiveDir,
			expected: filepath.Join(tmpDir, "my-spectr", "changes", "archive"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method()
			if got != tt.expected {
				t.Errorf("%s() = %q, want %q", tt.name, got, tt.expected)
			}
		})
	}
}

func TestValidateRootDir(t *testing.T) {
	tests := []struct {
		name      string
		rootDir   string
		wantError bool
	}{
		{"valid simple name", "spectr", false},
		{"valid with dash", "my-spectr", false},
		{"valid with underscore", "my_spectr", false},
		{"empty (will use default)", "", false},
		{"invalid forward slash", "path/to", true},
		{"invalid backward slash", "path\\to", true},
		{"invalid double dots", "../specs", true},
		{"invalid asterisk", "specs*", true},
		{"invalid question mark", "specs?", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRootDir(tt.rootDir)
			if tt.wantError && err == nil {
				t.Errorf("validateRootDir(%q) expected error, got nil", tt.rootDir)
			}
			if !tt.wantError && err != nil {
				t.Errorf("validateRootDir(%q) unexpected error: %v", tt.rootDir, err)
			}
		})
	}
}
