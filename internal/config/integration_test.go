package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/discovery"
)

// setupTestProject creates a test project structure without config file
func setupTestProject(t *testing.T, rootDir string) string {
	t.Helper()

	return createTestProjectStructure(t, rootDir)
}

// setupTestProjectWithConfig creates a test project structure with config file
func setupTestProjectWithConfig(t *testing.T, rootDir string) string {
	t.Helper()

	projectPath := createTestProjectStructure(t, rootDir)
	createSpectrYaml(t, projectPath, rootDir)

	return projectPath
}

// createTestProjectStructure creates the common project structure
func createTestProjectStructure(t *testing.T, rootDir string) string {
	t.Helper()

	// Create temporary project directory
	projectPath := t.TempDir()

	// Create root directory (e.g., "spectr" or "myspecs")
	spectrRoot := filepath.Join(projectPath, rootDir)
	if err := os.MkdirAll(spectrRoot, 0755); err != nil {
		t.Fatalf("failed to create root directory: %v", err)
	}

	// Create specs directory structure
	specsDir := filepath.Join(spectrRoot, "specs")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("failed to create specs directory: %v", err)
	}

	// Create changes directory structure
	changesDir := filepath.Join(spectrRoot, "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("failed to create changes directory: %v", err)
	}

	// Create a sample spec
	authSpecDir := filepath.Join(specsDir, "auth")
	if err := os.MkdirAll(authSpecDir, 0755); err != nil {
		t.Fatalf("failed to create auth spec directory: %v", err)
	}
	specContent := `# Auth Specification

## Requirements

### Requirement: User Authentication
The system SHALL authenticate users.

#### Scenario: Valid credentials
- **WHEN** valid credentials are provided
- **THEN** user is authenticated
`
	specPath := filepath.Join(authSpecDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("failed to create spec.md: %v", err)
	}

	// Create a sample change
	changeDir := filepath.Join(changesDir, "add-2fa")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("failed to create change directory: %v", err)
	}
	proposalContent := `# Change: Add Two-Factor Authentication

## Why
Enhance security with 2FA.

## What Changes
- Add 2FA support

## Impact
- Affects auth spec
`
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(proposalContent), 0644); err != nil {
		t.Fatalf("failed to create proposal.md: %v", err)
	}

	return projectPath
}

// createSpectrYaml creates a spectr.yaml with custom root_dir
func createSpectrYaml(t *testing.T, projectPath, rootDir string) {
	t.Helper()

	configContent := "root_dir: " + rootDir + "\n"
	configPath := filepath.Join(projectPath, config.ConfigFileName)

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create spectr.yaml: %v", err)
	}
}

func TestIntegration_DefaultBehavior(t *testing.T) {
	// Create a project with standard spectr/ directory structure
	// No spectr.yaml file
	projectPath := setupTestProject(t, config.DefaultRootDir)

	// Load config
	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		t.Fatalf("LoadFromPath failed: %v", err)
	}

	// Verify default settings
	if cfg.RootDir != config.DefaultRootDir {
		t.Errorf("expected RootDir=%q, got %q", config.DefaultRootDir, cfg.RootDir)
	}
	if cfg.ProjectRoot != projectPath {
		t.Errorf("expected ProjectRoot=%q, got %q", projectPath, cfg.ProjectRoot)
	}

	// Verify discovery works
	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}
	if len(specs) != 1 || specs[0] != "auth" {
		t.Errorf("expected specs=[auth], got %v", specs)
	}

	// Verify changes discovery works
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetActiveChangesWithConfig failed: %v", err)
	}
	if len(changes) != 1 || changes[0] != "add-2fa" {
		t.Errorf("expected changes=[add-2fa], got %v", changes)
	}

	// Verify paths are correct
	expectedRootPath := filepath.Join(projectPath, config.DefaultRootDir)
	if cfg.RootPath() != expectedRootPath {
		t.Errorf("expected RootPath=%q, got %q", expectedRootPath, cfg.RootPath())
	}

	expectedSpecsPath := filepath.Join(expectedRootPath, "specs")
	if cfg.SpecsPath() != expectedSpecsPath {
		t.Errorf("expected SpecsPath=%q, got %q", expectedSpecsPath, cfg.SpecsPath())
	}

	expectedChangesPath := filepath.Join(expectedRootPath, "changes")
	if cfg.ChangesPath() != expectedChangesPath {
		t.Errorf(
			"expected ChangesPath=%q, got %q",
			expectedChangesPath,
			cfg.ChangesPath(),
		)
	}
}

func TestIntegration_CustomRootDir(t *testing.T) {
	// Create spectr.yaml with root_dir: myspecs
	customRootDir := "myspecs"
	projectPath := setupTestProjectWithConfig(t, customRootDir)

	// Load config
	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		t.Fatalf("LoadFromPath failed: %v", err)
	}

	// Verify custom settings
	if cfg.RootDir != customRootDir {
		t.Errorf("expected RootDir=%q, got %q", customRootDir, cfg.RootDir)
	}
	if cfg.ProjectRoot != projectPath {
		t.Errorf("expected ProjectRoot=%q, got %q", projectPath, cfg.ProjectRoot)
	}

	// Verify discovery uses custom directory
	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}
	if len(specs) != 1 || specs[0] != "auth" {
		t.Errorf("expected specs=[auth], got %v", specs)
	}

	// Verify changes discovery uses custom directory
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetActiveChangesWithConfig failed: %v", err)
	}
	if len(changes) != 1 || changes[0] != "add-2fa" {
		t.Errorf("expected changes=[add-2fa], got %v", changes)
	}

	// Verify paths use custom root
	expectedRootPath := filepath.Join(projectPath, customRootDir)
	if cfg.RootPath() != expectedRootPath {
		t.Errorf("expected RootPath=%q, got %q", expectedRootPath, cfg.RootPath())
	}

	expectedSpecsPath := filepath.Join(expectedRootPath, "specs")
	if cfg.SpecsPath() != expectedSpecsPath {
		t.Errorf("expected SpecsPath=%q, got %q", expectedSpecsPath, cfg.SpecsPath())
	}

	expectedChangesPath := filepath.Join(expectedRootPath, "changes")
	if cfg.ChangesPath() != expectedChangesPath {
		t.Errorf(
			"expected ChangesPath=%q, got %q",
			expectedChangesPath,
			cfg.ChangesPath(),
		)
	}
}

func TestIntegration_ConfigDiscoveryFromSubdirectory(t *testing.T) {
	// Create project structure with spectr.yaml at root
	customRootDir := "docs"
	projectPath := setupTestProjectWithConfig(t, customRootDir)

	// Create a subdirectory to run from
	subDir := filepath.Join(projectPath, "src", "components")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	// Load config from subdirectory
	cfg, err := config.LoadFromPath(subDir)
	if err != nil {
		t.Fatalf("LoadFromPath from subdirectory failed: %v", err)
	}

	// Verify config was found and ProjectRoot points to project root
	if cfg.RootDir != customRootDir {
		t.Errorf("expected RootDir=%q, got %q", customRootDir, cfg.RootDir)
	}
	if cfg.ProjectRoot != projectPath {
		t.Errorf(
			"expected ProjectRoot=%q (project root), got %q",
			projectPath,
			cfg.ProjectRoot,
		)
	}

	// Verify discovery works from subdirectory
	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}
	if len(specs) != 1 || specs[0] != "auth" {
		t.Errorf("expected specs=[auth], got %v", specs)
	}

	// Verify changes discovery works from subdirectory
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetActiveChangesWithConfig failed: %v", err)
	}
	if len(changes) != 1 || changes[0] != "add-2fa" {
		t.Errorf("expected changes=[add-2fa], got %v", changes)
	}
}

func TestIntegration_BackwardCompatibility(t *testing.T) {
	// Create a project without spectr.yaml (backward compatibility)
	projectPath := setupTestProject(t, config.DefaultRootDir)

	// Load config should return defaults
	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		t.Fatalf("LoadFromPath failed: %v", err)
	}

	// Verify default behavior is preserved
	if cfg.RootDir != config.DefaultRootDir {
		t.Errorf("expected RootDir=%q, got %q", config.DefaultRootDir, cfg.RootDir)
	}

	// Verify all functionality works as before
	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}
	if len(specs) != 1 || specs[0] != "auth" {
		t.Errorf("expected specs=[auth], got %v", specs)
	}

	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetActiveChangesWithConfig failed: %v", err)
	}
	if len(changes) != 1 || changes[0] != "add-2fa" {
		t.Errorf("expected changes=[add-2fa], got %v", changes)
	}

	// Now add a spectr.yaml with default settings
	createSpectrYaml(t, projectPath, config.DefaultRootDir)

	// Load again
	cfg2, err := config.LoadFromPath(projectPath)
	if err != nil {
		t.Fatalf("LoadFromPath with config file failed: %v", err)
	}

	// Verify behavior is identical
	if cfg2.RootDir != cfg.RootDir {
		t.Errorf(
			"expected same RootDir after adding config, got %q vs %q",
			cfg2.RootDir,
			cfg.RootDir,
		)
	}

	specs2, err := discovery.GetSpecsWithConfig(cfg2)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig with config file failed: %v", err)
	}
	if len(specs2) != len(specs) {
		t.Errorf(
			"expected same specs count after adding config, got %d vs %d",
			len(specs2),
			len(specs),
		)
	}
}

func TestIntegration_MultipleCustomDirectories(t *testing.T) {
	// Test with various custom directory names
	customDirs := []string{"specs", "documentation", "my-specs", "proj_docs"}

	for _, customDir := range customDirs {
		t.Run(customDir, func(t *testing.T) {
			projectPath := setupTestProjectWithConfig(t, customDir)

			cfg, err := config.LoadFromPath(projectPath)
			if err != nil {
				t.Fatalf("LoadFromPath failed for %q: %v", customDir, err)
			}

			if cfg.RootDir != customDir {
				t.Errorf("expected RootDir=%q, got %q", customDir, cfg.RootDir)
			}

			// Verify discovery works
			specs, err := discovery.GetSpecsWithConfig(cfg)
			if err != nil {
				t.Fatalf("GetSpecsWithConfig failed for %q: %v", customDir, err)
			}
			if len(specs) != 1 || specs[0] != "auth" {
				t.Errorf("expected specs=[auth] for %q, got %v", customDir, specs)
			}

			changes, err := discovery.GetActiveChangesWithConfig(cfg)
			if err != nil {
				t.Fatalf(
					"GetActiveChangesWithConfig failed for %q: %v",
					customDir,
					err,
				)
			}
			if len(changes) != 1 || changes[0] != "add-2fa" {
				t.Errorf(
					"expected changes=[add-2fa] for %q, got %v",
					customDir,
					changes,
				)
			}
		})
	}
}

func TestIntegration_EmptyDirectories(t *testing.T) {
	// Test behavior when directories are empty
	projectPath := t.TempDir()
	customRootDir := "myspecs"

	// Create empty directory structure
	spectrRoot := filepath.Join(projectPath, customRootDir)
	if err := os.MkdirAll(filepath.Join(spectrRoot, "specs"), 0755); err != nil {
		t.Fatalf("failed to create specs directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(spectrRoot, "changes"), 0755); err != nil {
		t.Fatalf("failed to create changes directory: %v", err)
	}

	createSpectrYaml(t, projectPath, customRootDir)

	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		t.Fatalf("LoadFromPath failed: %v", err)
	}

	// Verify empty results
	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("expected empty specs, got %v", specs)
	}

	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetActiveChangesWithConfig failed: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("expected empty changes, got %v", changes)
	}
}

func TestIntegration_MissingDirectories(t *testing.T) {
	// Test behavior when spectr directories don't exist yet
	projectPath := t.TempDir()
	customRootDir := "docs"

	// Only create config, no directories
	createSpectrYaml(t, projectPath, customRootDir)

	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		t.Fatalf("LoadFromPath failed: %v", err)
	}

	// Verify discovery handles missing directories gracefully
	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf(
			"GetSpecsWithConfig should not error on missing directory: %v",
			err,
		)
	}
	if len(specs) != 0 {
		t.Errorf("expected empty specs for missing directory, got %v", specs)
	}

	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf(
			"GetActiveChangesWithConfig should not error on missing directory: %v",
			err,
		)
	}
	if len(changes) != 0 {
		t.Errorf("expected empty changes for missing directory, got %v", changes)
	}
}
