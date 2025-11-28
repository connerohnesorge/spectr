//go:build integration

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/list"
	"github.com/connerohnesorge/spectr/internal/view"
)

// TestIntegration_CustomRootDirectory tests that all modules work together
// with a custom root directory configured via spectr.yaml
func TestIntegration_CustomRootDirectory(t *testing.T) {
	// Setup: Create temp project with spectr.yaml and custom root "my-specs"
	projectRoot := t.TempDir()

	// Create spectr.yaml with custom root_dir
	configContent := "root_dir: my-specs\n"
	configPath := filepath.Join(projectRoot, "spectr.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create directory structure with custom root
	mySpecsDir := filepath.Join(projectRoot, "my-specs")
	changesDir := filepath.Join(mySpecsDir, "changes")
	specsDir := filepath.Join(mySpecsDir, "specs")

	// Create a test change
	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	proposalContent := `# Change: Test Change

## Why
Testing custom root directory support.

## What Changes
- Add integration tests
`
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to create proposal.md: %v", err)
	}

	tasksContent := `## Implementation
- [ ] Task 1
- [x] Task 2
- [ ] Task 3
`
	tasksPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to create tasks.md: %v", err)
	}

	// Create a test spec
	specDir := filepath.Join(specsDir, "test-spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("Failed to create spec directory: %v", err)
	}

	specContent := `# Specification: Test Spec

## Requirements

### Requirement: Test Requirement 1
The system SHALL support custom root directories.

#### Scenario: Config loaded correctly
- **WHEN** config file is present
- **THEN** custom root directory is used

### Requirement: Test Requirement 2
All modules SHALL work with custom root directories.

#### Scenario: Discovery works
- **WHEN** using custom root
- **THEN** changes and specs are discovered
`
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create spec.md: %v", err)
	}

	// Test 1: Load config from project root
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}

	// Verify config loaded correctly
	if cfg.RootDir != "my-specs" {
		t.Errorf("Expected RootDir='my-specs', got %q", cfg.RootDir)
	}

	if cfg.ProjectRoot != projectRoot {
		t.Errorf("Expected ProjectRoot=%q, got %q", projectRoot, cfg.ProjectRoot)
	}

	if cfg.ConfigPath != configPath {
		t.Errorf("Expected ConfigPath=%q, got %q", configPath, cfg.ConfigPath)
	}

	// Test 2: Discovery works with custom config
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetActiveChangesWithConfig() failed: %v", err)
	}

	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}

	if len(changes) > 0 && changes[0] != "test-change" {
		t.Errorf("Expected change ID 'test-change', got %q", changes[0])
	}

	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetSpecsWithConfig() failed: %v", err)
	}

	if len(specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(specs))
	}

	if len(specs) > 0 && specs[0] != "test-spec" {
		t.Errorf("Expected spec ID 'test-spec', got %q", specs[0])
	}

	// Test 3: List module works with custom config
	lister := list.NewListerWithConfig(cfg)

	changeInfos, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("lister.ListChanges() failed: %v", err)
	}

	if len(changeInfos) != 1 {
		t.Errorf("Expected 1 change info, got %d", len(changeInfos))
	}

	if len(changeInfos) > 0 {
		change := changeInfos[0]
		if change.ID != "test-change" {
			t.Errorf("Expected change ID 'test-change', got %q", change.ID)
		}
		if change.Title != "Test Change" {
			t.Errorf("Expected title 'Test Change', got %q", change.Title)
		}
		if change.TaskStatus.Total != 3 {
			t.Errorf("Expected 3 total tasks, got %d", change.TaskStatus.Total)
		}
		if change.TaskStatus.Completed != 1 {
			t.Errorf("Expected 1 completed task, got %d", change.TaskStatus.Completed)
		}
	}

	specInfos, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("lister.ListSpecs() failed: %v", err)
	}

	if len(specInfos) != 1 {
		t.Errorf("Expected 1 spec info, got %d", len(specInfos))
	}

	if len(specInfos) > 0 {
		spec := specInfos[0]
		if spec.ID != "test-spec" {
			t.Errorf("Expected spec ID 'test-spec', got %q", spec.ID)
		}
		if spec.Title != "Specification: Test Spec" {
			t.Errorf("Expected title 'Specification: Test Spec', got %q", spec.Title)
		}
		if spec.RequirementCount != 2 {
			t.Errorf("Expected 2 requirements, got %d", spec.RequirementCount)
		}
	}

	// Test 4: View module works with custom config
	data, err := view.CollectDataWithConfig(cfg)
	if err != nil {
		t.Fatalf("view.CollectDataWithConfig() failed: %v", err)
	}

	totalChanges := data.Summary.ActiveChanges + data.Summary.CompletedChanges
	if totalChanges != 1 {
		t.Errorf("Expected 1 total change, got %d", totalChanges)
	}

	if data.Summary.ActiveChanges != 1 {
		t.Errorf("Expected 1 active change, got %d", data.Summary.ActiveChanges)
	}

	if data.Summary.TotalSpecs != 1 {
		t.Errorf("Expected 1 spec, got %d", data.Summary.TotalSpecs)
	}

	if data.Summary.TotalRequirements != 2 {
		t.Errorf("Expected 2 requirements, got %d", data.Summary.TotalRequirements)
	}

	if data.Summary.TotalTasks != 3 {
		t.Errorf("Expected 3 total tasks, got %d", data.Summary.TotalTasks)
	}

	if data.Summary.CompletedTasks != 1 {
		t.Errorf("Expected 1 completed task, got %d", data.Summary.CompletedTasks)
	}
}

// TestIntegration_BackwardCompatibility tests that projects without spectr.yaml
// still work correctly using the default "spectr/" directory
func TestIntegration_BackwardCompatibility(t *testing.T) {
	// Setup: Create temp project WITHOUT spectr.yaml, using default "spectr/" dir
	projectRoot := t.TempDir()

	// Create directory structure with default root
	spectrDir := filepath.Join(projectRoot, "spectr")
	changesDir := filepath.Join(spectrDir, "changes")
	specsDir := filepath.Join(spectrDir, "specs")

	// Create a test change
	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	proposalContent := `# Change: Legacy Test

## Why
Testing backward compatibility.
`
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to create proposal.md: %v", err)
	}

	tasksContent := `## Tasks
- [x] Task 1
`
	tasksPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to create tasks.md: %v", err)
	}

	// Create a test spec
	specDir := filepath.Join(specsDir, "legacy-spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("Failed to create spec directory: %v", err)
	}

	specContent := `# Specification: Legacy Spec

## Requirements

### Requirement: Backward Compatibility
The system SHALL work without config file.

#### Scenario: Default directory used
- **WHEN** no config file exists
- **THEN** default spectr/ directory is used
`
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create spec.md: %v", err)
	}

	// Load config - should find no file and use defaults
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}

	if cfg.RootDir != "spectr" {
		t.Errorf("Expected default RootDir='spectr', got %q", cfg.RootDir)
	}

	if cfg.ConfigPath != "" {
		t.Errorf("Expected empty ConfigPath for default config, got %q", cfg.ConfigPath)
	}

	if cfg.ProjectRoot != projectRoot {
		t.Errorf("Expected ProjectRoot=%q, got %q", projectRoot, cfg.ProjectRoot)
	}

	// Verify all modules work with default config
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetActiveChangesWithConfig() failed: %v", err)
	}

	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}

	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetSpecsWithConfig() failed: %v", err)
	}

	if len(specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(specs))
	}

	// Verify list module works
	lister := list.NewListerWithConfig(cfg)

	changeInfos, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("lister.ListChanges() failed: %v", err)
	}

	if len(changeInfos) != 1 {
		t.Errorf("Expected 1 change info, got %d", len(changeInfos))
	}

	specInfos, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("lister.ListSpecs() failed: %v", err)
	}

	if len(specInfos) != 1 {
		t.Errorf("Expected 1 spec info, got %d", len(specInfos))
	}

	// Verify view module works
	data, err := view.CollectDataWithConfig(cfg)
	if err != nil {
		t.Fatalf("view.CollectDataWithConfig() failed: %v", err)
	}

	totalChanges := data.Summary.ActiveChanges + data.Summary.CompletedChanges
	if totalChanges != 1 {
		t.Errorf("Expected 1 total change (active or completed), got %d", totalChanges)
	}

	if data.Summary.TotalSpecs != 1 {
		t.Errorf("Expected 1 spec, got %d", data.Summary.TotalSpecs)
	}
}

// TestIntegration_NestedDirectoryDiscovery tests that config is discovered
// correctly when running from deeply nested subdirectories
func TestIntegration_NestedDirectoryDiscovery(t *testing.T) {
	// Setup: Create temp project with spectr.yaml at root
	projectRoot := t.TempDir()

	// Create spectr.yaml with custom root_dir
	configContent := "root_dir: custom-specs\n"
	configPath := filepath.Join(projectRoot, "spectr.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create directory structure
	customSpecsDir := filepath.Join(projectRoot, "custom-specs")
	changesDir := filepath.Join(customSpecsDir, "changes")
	specsDir := filepath.Join(customSpecsDir, "specs")

	// Create a change and spec for verification
	changeDir := filepath.Join(changesDir, "nested-test")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	proposalContent := `# Change: Nested Test

## Why
Testing nested directory discovery.
`
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to create proposal.md: %v", err)
	}

	specDir := filepath.Join(specsDir, "nested-spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("Failed to create spec directory: %v", err)
	}

	specContent := `# Specification: Nested Spec

## Requirements

### Requirement: Discovery from nested directories
The system SHALL discover config from parent directories.

#### Scenario: Nested execution
- **WHEN** running from nested directory
- **THEN** config is found in parent
`
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to create spec.md: %v", err)
	}

	// Create deeply nested subdirectory (simulating working from deep in project)
	nestedDir := filepath.Join(projectRoot, "src", "components", "nested", "deep")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	// Test: Load config from deeply nested directory
	cfg, err := config.Load(nestedDir)
	if err != nil {
		t.Fatalf("config.Load() from nested dir failed: %v", err)
	}

	// Verify config was discovered correctly
	if cfg.ProjectRoot != projectRoot {
		t.Errorf("Expected ProjectRoot=%q, got %q", projectRoot, cfg.ProjectRoot)
	}

	if cfg.RootDir != "custom-specs" {
		t.Errorf("Expected RootDir='custom-specs', got %q", cfg.RootDir)
	}

	if cfg.ConfigPath != configPath {
		t.Errorf("Expected ConfigPath=%q, got %q", configPath, cfg.ConfigPath)
	}

	// Verify all modules work from nested directory context
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetActiveChangesWithConfig() failed: %v", err)
	}

	if len(changes) != 1 {
		t.Errorf("Expected 1 change from nested dir, got %d", len(changes))
	}

	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetSpecsWithConfig() failed: %v", err)
	}

	if len(specs) != 1 {
		t.Errorf("Expected 1 spec from nested dir, got %d", len(specs))
	}

	// Verify lister works
	lister := list.NewListerWithConfig(cfg)

	changeInfos, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("lister.ListChanges() from nested dir failed: %v", err)
	}

	if len(changeInfos) != 1 {
		t.Errorf("Expected 1 change info from nested dir, got %d", len(changeInfos))
	}

	specInfos, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("lister.ListSpecs() from nested dir failed: %v", err)
	}

	if len(specInfos) != 1 {
		t.Errorf("Expected 1 spec info from nested dir, got %d", len(specInfos))
	}

	// Verify view works
	data, err := view.CollectDataWithConfig(cfg)
	if err != nil {
		t.Fatalf("view.CollectDataWithConfig() from nested dir failed: %v", err)
	}

	totalChanges := data.Summary.ActiveChanges + data.Summary.CompletedChanges
	if totalChanges != 1 {
		t.Errorf("Expected 1 total change from nested dir, got %d", totalChanges)
	}

	if data.Summary.TotalSpecs != 1 {
		t.Errorf("Expected 1 spec from nested dir, got %d", data.Summary.TotalSpecs)
	}
}

// TestIntegration_MultipleChangesAndSpecs tests integration with multiple
// changes and specs to ensure proper counting and aggregation
func TestIntegration_MultipleChangesAndSpecs(t *testing.T) {
	projectRoot := t.TempDir()

	// Create config
	configContent := "root_dir: docs\n"
	configPath := filepath.Join(projectRoot, "spectr.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	docsDir := filepath.Join(projectRoot, "docs")
	changesDir := filepath.Join(docsDir, "changes")
	specsDir := filepath.Join(docsDir, "specs")

	// Create multiple changes
	for i, changeID := range []string{"change-a", "change-b", "change-c"} {
		changeDir := filepath.Join(changesDir, changeID)
		if err := os.MkdirAll(changeDir, 0755); err != nil {
			t.Fatalf("Failed to create change directory: %v", err)
		}

		proposalContent := "# Change: Change " + string(rune('A'+i)) + "\n\n## Why\nTest change.\n"
		if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
			t.Fatalf("Failed to create proposal.md: %v", err)
		}

		tasksContent := "## Tasks\n- [ ] Task 1\n- [ ] Task 2\n"
		if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
			t.Fatalf("Failed to create tasks.md: %v", err)
		}
	}

	// Create multiple specs
	for i, specID := range []string{"spec-x", "spec-y"} {
		specDir := filepath.Join(specsDir, specID)
		if err := os.MkdirAll(specDir, 0755); err != nil {
			t.Fatalf("Failed to create spec directory: %v", err)
		}

		specContent := "# Specification: Spec " + string(rune('X'+i)) + "\n\n## Requirements\n\n" +
			"### Requirement: Req 1\n" +
			"Description.\n\n" +
			"#### Scenario: Test\n" +
			"- **WHEN** something\n" +
			"- **THEN** result\n"

		if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
			t.Fatalf("Failed to create spec.md: %v", err)
		}
	}

	// Load config
	cfg, err := config.Load(projectRoot)
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}

	// Test discovery
	changes, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetActiveChangesWithConfig() failed: %v", err)
	}

	if len(changes) != 3 {
		t.Errorf("Expected 3 changes, got %d", len(changes))
	}

	specs, err := discovery.GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("discovery.GetSpecsWithConfig() failed: %v", err)
	}

	if len(specs) != 2 {
		t.Errorf("Expected 2 specs, got %d", len(specs))
	}

	// Test lister
	lister := list.NewListerWithConfig(cfg)

	changeInfos, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("lister.ListChanges() failed: %v", err)
	}

	if len(changeInfos) != 3 {
		t.Errorf("Expected 3 change infos, got %d", len(changeInfos))
	}

	specInfos, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("lister.ListSpecs() failed: %v", err)
	}

	if len(specInfos) != 2 {
		t.Errorf("Expected 2 spec infos, got %d", len(specInfos))
	}

	// Test view
	data, err := view.CollectDataWithConfig(cfg)
	if err != nil {
		t.Fatalf("view.CollectDataWithConfig() failed: %v", err)
	}

	if data.Summary.ActiveChanges != 3 {
		t.Errorf("Expected 3 active changes, got %d", data.Summary.ActiveChanges)
	}

	if data.Summary.TotalSpecs != 2 {
		t.Errorf("Expected 2 specs, got %d", data.Summary.TotalSpecs)
	}

	if data.Summary.TotalTasks != 6 { // 3 changes × 2 tasks each
		t.Errorf("Expected 6 total tasks, got %d", data.Summary.TotalTasks)
	}
}
