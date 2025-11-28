package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Create a temporary directory with no config file
	tmpDir := t.TempDir()

	// Change to temp directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	_ = os.Chdir(tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf("expected RootDir=%q, got %q", DefaultRootDir, cfg.RootDir)
	}

	if cfg.ProjectRoot != tmpDir {
		t.Errorf("expected ProjectRoot=%q, got %q", tmpDir, cfg.ProjectRoot)
	}

	if cfg.Theme != "default" {
		t.Errorf("expected Theme=%q, got %q", "default", cfg.Theme)
	}
}

func TestLoadFromPath_Defaults(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf("expected RootDir=%q, got %q", DefaultRootDir, cfg.RootDir)
	}

	if cfg.ProjectRoot != tmpDir {
		t.Errorf("expected ProjectRoot=%q, got %q", tmpDir, cfg.ProjectRoot)
	}

	if cfg.Theme != "default" {
		t.Errorf("expected Theme=%q, got %q", "default", cfg.Theme)
	}
}

func TestLoadFromPath_CustomRootDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create spectr.yaml with custom root_dir
	configContent := `root_dir: specs`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.RootDir != "specs" {
		t.Errorf("expected RootDir=%q, got %q", "specs", cfg.RootDir)
	}

	if cfg.ProjectRoot != tmpDir {
		t.Errorf("expected ProjectRoot=%q, got %q", tmpDir, cfg.ProjectRoot)
	}
}

func TestLoadFromPath_DiscoveryFromSubdirectory(t *testing.T) {
	// Create directory structure:
	// tmpDir/
	//   spectr.yaml
	//   subdir/
	//     nested/
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir", "nested")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectories: %v", err)
	}

	// Create config in root
	configContent := `root_dir: custom-specs`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Load from nested subdirectory
	cfg, err := LoadFromPath(subDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.RootDir != "custom-specs" {
		t.Errorf("expected RootDir=%q, got %q", "custom-specs", cfg.RootDir)
	}

	if cfg.ProjectRoot != tmpDir {
		t.Errorf("expected ProjectRoot=%q, got %q", tmpDir, cfg.ProjectRoot)
	}
}

func TestLoadFromPath_MultipleConfigs_NearestWins(t *testing.T) {
	// Create directory structure:
	// tmpDir/
	//   spectr.yaml (root_dir: outer)
	//   subdir/
	//     spectr.yaml (root_dir: inner)
	//     nested/
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	nestedDir := filepath.Join(subDir, "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectories: %v", err)
	}

	// Create outer config
	outerConfig := `root_dir: outer`
	if err := os.WriteFile(filepath.Join(tmpDir, ConfigFileName), []byte(outerConfig), 0644); err != nil {
		t.Fatalf("failed to write outer config: %v", err)
	}

	// Create inner config
	innerConfig := `root_dir: inner`
	if err := os.WriteFile(filepath.Join(subDir, ConfigFileName), []byte(innerConfig), 0644); err != nil {
		t.Fatalf("failed to write inner config: %v", err)
	}

	// Load from nested - should find inner config first
	cfg, err := LoadFromPath(nestedDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.RootDir != "inner" {
		t.Errorf("expected RootDir=%q (nearest config), got %q", "inner", cfg.RootDir)
	}

	if cfg.ProjectRoot != subDir {
		t.Errorf("expected ProjectRoot=%q, got %q", subDir, cfg.ProjectRoot)
	}
}

func TestLoadFromPath_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid YAML
	invalidYAML := `
root_dir: test
  invalid_indent: value
    more_bad: stuff
`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadFromPath(tmpDir)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}

	if !strings.Contains(err.Error(), "YAML") {
		t.Errorf("expected error to mention YAML, got: %v", err)
	}
}

func TestValidate_InvalidCharacters(t *testing.T) {
	tests := []struct {
		name      string
		rootDir   string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid simple name",
			rootDir:   "spectr",
			wantError: false,
		},
		{
			name:      "valid with dash",
			rootDir:   "my-specs",
			wantError: false,
		},
		{
			name:      "valid with underscore",
			rootDir:   "my_specs",
			wantError: false,
		},
		{
			name:      "invalid with forward slash",
			rootDir:   "path/to/specs",
			wantError: true,
			errorMsg:  "invalid characters: /",
		},
		{
			name:      "invalid with backslash",
			rootDir:   "path\\to\\specs",
			wantError: true,
			errorMsg:  "invalid characters: \\",
		},
		{
			name:      "invalid with double dot",
			rootDir:   "../specs",
			wantError: true,
			errorMsg:  "invalid characters",
		},
		{
			name:      "invalid with asterisk",
			rootDir:   "spec*",
			wantError: true,
			errorMsg:  "invalid characters: *",
		},
		{
			name:      "invalid with multiple bad chars",
			rootDir:   "../path/to/*",
			wantError: true,
			errorMsg:  "invalid characters",
		},
		{
			name:      "invalid empty string",
			rootDir:   "",
			wantError: true,
			errorMsg:  "cannot be empty",
		},
		{
			name:      "invalid hidden directory",
			rootDir:   ".specs",
			wantError: true,
			errorMsg:  "cannot start with '.'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				RootDir:     tt.rootDir,
				ProjectRoot: "/tmp/test",
				Theme:       "default",
			}

			err := cfg.validate()

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLoadFromPath_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config with invalid root_dir
	invalidConfig := `root_dir: ../invalid`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadFromPath(tmpDir)
	if err == nil {
		t.Fatal("expected error for invalid root_dir, got nil")
	}

	if !strings.Contains(err.Error(), "invalid configuration") {
		t.Errorf("expected error to mention invalid configuration, got: %v", err)
	}

	if !strings.Contains(err.Error(), "invalid characters") {
		t.Errorf("expected error to mention invalid characters, got: %v", err)
	}
}

func TestConfig_HelperMethods(t *testing.T) {
	cfg := &Config{
		RootDir:     "specs",
		ProjectRoot: "/home/user/project",
		Theme:       "default",
	}

	rootPath := cfg.RootPath()
	expectedRoot := filepath.Join("/home/user/project", "specs")
	if rootPath != expectedRoot {
		t.Errorf("RootPath() = %q, want %q", rootPath, expectedRoot)
	}

	specsPath := cfg.SpecsPath()
	expectedSpecs := filepath.Join("/home/user/project", "specs", "specs")
	if specsPath != expectedSpecs {
		t.Errorf("SpecsPath() = %q, want %q", specsPath, expectedSpecs)
	}

	changesPath := cfg.ChangesPath()
	expectedChanges := filepath.Join("/home/user/project", "specs", "changes")
	if changesPath != expectedChanges {
		t.Errorf("ChangesPath() = %q, want %q", changesPath, expectedChanges)
	}
}

func TestLoadFromPath_EmptyConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create empty config file (should use defaults)
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf(
			"expected RootDir=%q (default for empty config), got %q",
			DefaultRootDir,
			cfg.RootDir,
		)
	}
}

func TestLoadFromPath_CommentOnlyConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config with only comments
	configContent := `
# This is a comment
# root_dir: commented-out
`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.RootDir != DefaultRootDir {
		t.Errorf(
			"expected RootDir=%q (default for comment-only config), got %q",
			DefaultRootDir,
			cfg.RootDir,
		)
	}
}

func TestLoadFromPath_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config
	configContent := `root_dir: specs`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Change to parent directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	_ = os.Chdir(filepath.Dir(tmpDir))

	// Load using relative path
	relPath, _ := filepath.Rel(filepath.Dir(tmpDir), tmpDir)
	cfg, err := LoadFromPath(relPath)
	if err != nil {
		t.Fatalf("LoadFromPath() with relative path failed: %v", err)
	}

	if cfg.RootDir != "specs" {
		t.Errorf("expected RootDir=%q, got %q", "specs", cfg.RootDir)
	}

	// ProjectRoot should be absolute
	if !filepath.IsAbs(cfg.ProjectRoot) {
		t.Errorf("expected ProjectRoot to be absolute, got %q", cfg.ProjectRoot)
	}

	expectedRoot, _ := filepath.Abs(tmpDir)
	if cfg.ProjectRoot != expectedRoot {
		t.Errorf("expected ProjectRoot=%q, got %q", expectedRoot, cfg.ProjectRoot)
	}
}

func TestLoadFromPath_DefaultTheme(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a config file with just root_dir (no theme specified)
	configContent := `root_dir: spectr`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.Theme != "default" {
		t.Errorf("expected Theme=%q (default when not specified), got %q", "default", cfg.Theme)
	}
}

func TestLoadFromPath_CustomTheme(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a config file with a valid custom theme
	configContent := `root_dir: spectr
theme: monokai`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromPath() failed: %v", err)
	}

	if cfg.Theme != "monokai" {
		t.Errorf("expected Theme=%q, got %q", "monokai", cfg.Theme)
	}
}

func TestLoadFromPath_InvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a config file with an invalid theme
	configContent := `root_dir: spectr
theme: nonexistent`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadFromPath(tmpDir)
	if err == nil {
		t.Fatal("expected error for invalid theme, got nil")
	}

	// Verify error message contains helpful information
	if !strings.Contains(err.Error(), "invalid theme 'nonexistent'") {
		t.Errorf("expected error to contain \"invalid theme 'nonexistent'\", got: %v", err)
	}

	if !strings.Contains(err.Error(), "available themes:") {
		t.Errorf("expected error to contain \"available themes:\", got: %v", err)
	}
}
