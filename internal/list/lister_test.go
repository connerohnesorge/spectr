package list

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
)

func TestListChanges(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test change with all components
	changeDir := filepath.Join(changesDir, "add-feature")
	if err := os.MkdirAll(filepath.Join(changeDir, "specs", "test-spec"), 0755); err != nil {
		t.Fatal(err)
	}

	// Write proposal.md
	proposalContent := `# Change: Add Amazing Feature

More details here.`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Write tasks.md
	tasksContent := `## Tasks
- [x] Task 1
- [ ] Task 2
- [x] Task 3`
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Write spec delta
	specContent := `## ADDED Requirements
### Requirement: New Feature

## MODIFIED Requirements
### Requirement: Updated Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "specs", "test-spec", "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test listing
	lister := NewLister(tmpDir)
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	change := changes[0]
	if change.ID != "add-feature" {
		t.Errorf("Expected ID 'add-feature', got %q", change.ID)
	}
	if change.Title != "Add Amazing Feature" {
		t.Errorf("Expected title 'Add Amazing Feature', got %q", change.Title)
	}
	if change.DeltaCount != 2 {
		t.Errorf("Expected delta count 2, got %d", change.DeltaCount)
	}
	if change.TaskStatus.Total != 3 {
		t.Errorf("Expected 3 total tasks, got %d", change.TaskStatus.Total)
	}
	if change.TaskStatus.Completed != 2 {
		t.Errorf("Expected 2 completed tasks, got %d", change.TaskStatus.Completed)
	}
}

func TestListChanges_NoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	lister := NewLister(tmpDir)
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("Expected empty list, got %d changes", len(changes))
	}
}

func TestListChanges_FallbackTitle(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write proposal.md without H1 heading
	proposalContent := `Some content without heading`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	lister := NewLister(tmpDir)
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	// Should fall back to ID as title
	if changes[0].Title != "test-change" {
		t.Errorf("Expected fallback title 'test-change', got %q", changes[0].Title)
	}
}

func TestListSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	specDir := filepath.Join(specsDir, "authentication")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write spec.md
	specContent := `# Authentication

### Requirement: User Login
Login feature

### Requirement: Password Reset
Reset feature

### Requirement: Two-Factor Auth
2FA feature`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test listing
	lister := NewLister(tmpDir)
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}

	if len(specs) != 1 {
		t.Fatalf("Expected 1 spec, got %d", len(specs))
	}

	spec := specs[0]
	if spec.ID != "authentication" {
		t.Errorf("Expected ID 'authentication', got %q", spec.ID)
	}
	if spec.Title != "Authentication" {
		t.Errorf("Expected title 'Authentication', got %q", spec.Title)
	}
	if spec.RequirementCount != 3 {
		t.Errorf("Expected 3 requirements, got %d", spec.RequirementCount)
	}
}

func TestListSpecs_NoSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	lister := NewLister(tmpDir)
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("Expected empty list, got %d specs", len(specs))
	}
}

func TestListSpecs_FallbackTitle(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	specDir := filepath.Join(specsDir, "test-spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Write spec.md without H1 heading
	specContent := `Some content without heading

### Requirement: Feature`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	lister := NewLister(tmpDir)
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}

	if len(specs) != 1 {
		t.Fatalf("Expected 1 spec, got %d", len(specs))
	}

	// Should fall back to ID as title
	if specs[0].Title != "test-spec" {
		t.Errorf("Expected fallback title 'test-spec', got %q", specs[0].Title)
	}
}

func TestListAll(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a change
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	changeDir := filepath.Join(changesDir, "add-feature")
	if err := os.MkdirAll(filepath.Join(changeDir, "specs", "test-spec"), 0755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Change: Add Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	tasksContent := `## Tasks
- [x] Task 1`
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatal(err)
	}

	specDeltaContent := `## ADDED Requirements
### Requirement: New Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "specs", "test-spec", "spec.md"), []byte(specDeltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a spec
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	specDir := filepath.Join(specsDir, "authentication")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Authentication

### Requirement: User Login`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test listing all items
	lister := NewLister(tmpDir)
	items, err := lister.ListAll(nil)
	if err != nil {
		t.Fatalf("ListAll failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// Verify sorting by ID (add-feature comes before authentication)
	if items[0].ID() != "add-feature" {
		t.Errorf("Expected first item to be 'add-feature', got %q", items[0].ID())
	}
	if items[0].Type != ItemTypeChange {
		t.Errorf("Expected first item to be a change, got %v", items[0].Type)
	}

	if items[1].ID() != "authentication" {
		t.Errorf("Expected second item to be 'authentication', got %q", items[1].ID())
	}
	if items[1].Type != ItemTypeSpec {
		t.Errorf("Expected second item to be a spec, got %v", items[1].Type)
	}
}

func TestListAll_FilterByType(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a change
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	changeDir := filepath.Join(changesDir, "add-feature")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Change: Add Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a spec
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	specDir := filepath.Join(specsDir, "authentication")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Authentication

### Requirement: User Login`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	lister := NewLister(tmpDir)

	// Test filtering for changes only
	changeType := ItemTypeChange
	items, err := lister.ListAll(&ListAllOptions{
		FilterType: &changeType,
		SortByID:   true,
	})
	if err != nil {
		t.Fatalf("ListAll with change filter failed: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item (change), got %d", len(items))
	}

	if items[0].Type != ItemTypeChange {
		t.Errorf("Expected change item, got %v", items[0].Type)
	}

	// Test filtering for specs only
	specType := ItemTypeSpec
	items, err = lister.ListAll(&ListAllOptions{
		FilterType: &specType,
		SortByID:   true,
	})
	if err != nil {
		t.Fatalf("ListAll with spec filter failed: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Expected 1 item (spec), got %d", len(items))
	}

	if items[0].Type != ItemTypeSpec {
		t.Errorf("Expected spec item, got %v", items[0].Type)
	}
}

func TestListAll_NoSorting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple changes with IDs that sort differently
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	for _, id := range []string{"zebra-change", "alpha-change"} {
		changeDir := filepath.Join(changesDir, id)
		if err := os.MkdirAll(changeDir, 0755); err != nil {
			t.Fatal(err)
		}

		content := `# Change: Test`
		if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	lister := NewLister(tmpDir)

	// Test with sorting disabled
	items, err := lister.ListAll(&ListAllOptions{
		SortByID: false,
	})
	if err != nil {
		t.Fatalf("ListAll without sorting failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}

	// Without sorting, order depends on filesystem readdir order
	// We just verify all items are present
	ids := make(map[string]bool)
	for _, item := range items {
		ids[item.ID()] = true
	}

	if !ids["zebra-change"] || !ids["alpha-change"] {
		t.Error("Expected both zebra-change and alpha-change to be present")
	}
}

func TestListAll_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	lister := NewLister(tmpDir)

	items, err := lister.ListAll(nil)
	if err != nil {
		t.Fatalf("ListAll on empty directory failed: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected empty list, got %d items", len(items))
	}
}

func TestNewListerWithConfig(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		RootDir:     "spectr",
		ProjectRoot: tmpDir,
	}

	lister := NewListerWithConfig(cfg)

	if lister.projectPath != tmpDir {
		t.Errorf("Expected projectPath to be %q, got %q", tmpDir, lister.projectPath)
	}

	if lister.cfg != cfg {
		t.Error("Expected lister.cfg to match the provided config")
	}
}

func TestNewListerWithConfig_CustomRootDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up custom root directory structure
	customRoot := "custom-specs"
	changesDir := filepath.Join(tmpDir, customRoot, "changes")
	specsDir := filepath.Join(tmpDir, customRoot, "specs")

	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a change in custom directory
	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Change: Test Custom Root`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a spec in custom directory
	specDir := filepath.Join(specsDir, "test-spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Test Spec

### Requirement: Test Requirement`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create lister with custom config
	cfg := &config.Config{
		RootDir:     customRoot,
		ProjectRoot: tmpDir,
	}

	lister := NewListerWithConfig(cfg)

	// Test listing changes
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].ID != "test-change" {
		t.Errorf("Expected change ID 'test-change', got %q", changes[0].ID)
	}

	if changes[0].Title != "Test Custom Root" {
		t.Errorf("Expected title 'Test Custom Root', got %q", changes[0].Title)
	}

	// Test listing specs
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}

	if len(specs) != 1 {
		t.Fatalf("Expected 1 spec, got %d", len(specs))
	}

	if specs[0].ID != "test-spec" {
		t.Errorf("Expected spec ID 'test-spec', got %q", specs[0].ID)
	}

	if specs[0].Title != "Test Spec" {
		t.Errorf("Expected title 'Test Spec', got %q", specs[0].Title)
	}

	// Test listing all
	items, err := lister.ListAll(nil)
	if err != nil {
		t.Fatalf("ListAll failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(items))
	}
}

func TestNewLister_BackwardCompatibility(t *testing.T) {
	tmpDir := t.TempDir()

	// Create standard spectr directory structure
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatal(err)
	}

	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Change: Test Backward Compatibility`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Use deprecated NewLister function
	lister := NewLister(tmpDir)

	// Should still work with default "spectr" directory
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].ID != "test-change" {
		t.Errorf("Expected change ID 'test-change', got %q", changes[0].ID)
	}

	// Verify internal config uses defaults
	if lister.cfg.RootDir != config.DefaultRootDir {
		t.Errorf("Expected RootDir to be %q, got %q", config.DefaultRootDir, lister.cfg.RootDir)
	}
}
