package view

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
)

//nolint:revive // cognitive-complexity justified for comprehensive testing
func TestCollectData(t *testing.T) {
	// Use the current project as test data
	projectPath := filepath.Join("..", "..")

	// Verify the project exists
	if _, err := os.Stat(filepath.Join(projectPath, "spectr")); err != nil {
		t.Skipf("Skipping test: spectr directory not found at %s", projectPath)
	}

	data, err := CollectData(projectPath)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify data structure is populated
	if data == nil {
		t.Fatal("CollectData returned nil data")
	}

	// Basic validation - summary metrics should be non-negative
	if data.Summary.TotalSpecs < 0 {
		t.Errorf("TotalSpecs should be >= 0, got %d", data.Summary.TotalSpecs)
	}
	if data.Summary.TotalRequirements < 0 {
		t.Errorf("TotalRequirements should be >= 0, got %d", data.Summary.TotalRequirements)
	}
	if data.Summary.ActiveChanges < 0 {
		t.Errorf("ActiveChanges should be >= 0, got %d", data.Summary.ActiveChanges)
	}
	if data.Summary.CompletedChanges < 0 {
		t.Errorf("CompletedChanges should be >= 0, got %d", data.Summary.CompletedChanges)
	}

	// Verify counts match slice lengths
	if data.Summary.ActiveChanges != len(data.ActiveChanges) {
		t.Errorf("ActiveChanges count mismatch: summary=%d, slice=%d",
			data.Summary.ActiveChanges, len(data.ActiveChanges))
	}
	if data.Summary.CompletedChanges != len(data.CompletedChanges) {
		t.Errorf("CompletedChanges count mismatch: summary=%d, slice=%d",
			data.Summary.CompletedChanges, len(data.CompletedChanges))
	}
	if data.Summary.TotalSpecs != len(data.Specs) {
		t.Errorf("TotalSpecs count mismatch: summary=%d, slice=%d",
			data.Summary.TotalSpecs, len(data.Specs))
	}

	// Verify sorting of active changes (ascending percentage)
	for i := 1; i < len(data.ActiveChanges); i++ {
		prev := data.ActiveChanges[i-1]
		curr := data.ActiveChanges[i]
		if prev.Progress.Percentage > curr.Progress.Percentage {
			t.Errorf("ActiveChanges not sorted by percentage: %s (%d%%) > %s (%d%%)",
				prev.ID, prev.Progress.Percentage, curr.ID, curr.Progress.Percentage)
		}
	}

	// Verify sorting of specs (descending requirement count)
	for i := 1; i < len(data.Specs); i++ {
		prev := data.Specs[i-1]
		curr := data.Specs[i]
		if prev.RequirementCount < curr.RequirementCount {
			t.Errorf("Specs not sorted by requirement count: %s (%d) < %s (%d)",
				prev.ID, prev.RequirementCount, curr.ID, curr.RequirementCount)
		}
	}

	// Verify sorting of completed changes (alphabetical)
	for i := 1; i < len(data.CompletedChanges); i++ {
		prev := data.CompletedChanges[i-1]
		curr := data.CompletedChanges[i]
		if prev.ID > curr.ID {
			t.Errorf("CompletedChanges not sorted alphabetically: %s > %s",
				prev.ID, curr.ID)
		}
	}

	t.Log("Dashboard data collected successfully:")
	t.Logf(
		"  Specs: %d (with %d requirements)",
		data.Summary.TotalSpecs,
		data.Summary.TotalRequirements,
	)
	t.Logf("  Active Changes: %d", data.Summary.ActiveChanges)
	for _, change := range data.ActiveChanges {
		t.Logf(
			"    - %s: %d/%d tasks (%d%%)",
			change.ID,
			change.Progress.Completed,
			change.Progress.Total,
			change.Progress.Percentage,
		)
	}
	t.Logf("  Completed Changes: %d", data.Summary.CompletedChanges)
	for _, change := range data.CompletedChanges {
		t.Logf("    - %s", change.ID)
	}
	t.Logf("  Tasks: %d/%d completed", data.Summary.CompletedTasks, data.Summary.TotalTasks)
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		name      string
		completed int
		total     int
		expected  int
	}{
		{"Zero total", 0, 0, 0},
		{"Zero completed", 0, 10, 0},
		{"Half completed", 5, 10, 50},
		{"All completed", 10, 10, 100},
		{"One third", 1, 3, 33},
		{"Two thirds", 2, 3, 66},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentage(tt.completed, tt.total)
			if result != tt.expected {
				t.Errorf("calculatePercentage(%d, %d) = %d; want %d",
					tt.completed, tt.total, result, tt.expected)
			}
		})
	}
}

// TestCollectData_EmptyProject tests collecting data from a project with no changes or specs
func TestCollectData_EmptyProject(t *testing.T) {
	// Create a temporary directory with only a spectr directory (no changes, no specs)
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	if err := os.MkdirAll(spectrDir, 0755); err != nil {
		t.Fatalf("Failed to create spectr directory: %v", err)
	}

	// Create empty subdirectories
	if err := os.MkdirAll(filepath.Join(spectrDir, "changes"), 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed on empty project: %v", err)
	}

	// Verify all metrics are zero
	if data.Summary.TotalSpecs != 0 {
		t.Errorf("Expected TotalSpecs=0, got %d", data.Summary.TotalSpecs)
	}
	if data.Summary.TotalRequirements != 0 {
		t.Errorf("Expected TotalRequirements=0, got %d", data.Summary.TotalRequirements)
	}
	if data.Summary.ActiveChanges != 0 {
		t.Errorf("Expected ActiveChanges=0, got %d", data.Summary.ActiveChanges)
	}
	if data.Summary.CompletedChanges != 0 {
		t.Errorf("Expected CompletedChanges=0, got %d", data.Summary.CompletedChanges)
	}
	if data.Summary.TotalTasks != 0 {
		t.Errorf("Expected TotalTasks=0, got %d", data.Summary.TotalTasks)
	}
	if data.Summary.CompletedTasks != 0 {
		t.Errorf("Expected CompletedTasks=0, got %d", data.Summary.CompletedTasks)
	}

	// Verify slices are initialized but empty
	if len(data.ActiveChanges) != 0 {
		t.Errorf("Expected empty ActiveChanges, got %d items", len(data.ActiveChanges))
	}
	if len(data.CompletedChanges) != 0 {
		t.Errorf("Expected empty CompletedChanges, got %d items", len(data.CompletedChanges))
	}
	if len(data.Specs) != 0 {
		t.Errorf("Expected empty Specs, got %d items", len(data.Specs))
	}

	t.Log("Empty project correctly handled with all zero metrics")
}

// TestCollectData_OnlyActiveChanges tests a project with only active changes
func TestCollectData_OnlyActiveChanges(t *testing.T) {
	// Create a temporary directory with only active changes
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(spectrDir, "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create a test change with tasks
	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	// Write a proposal.md
	proposalContent := "# Test Change\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to write proposal.md: %v", err)
	}

	// Write tasks.md with some tasks
	tasksContent := "## 1. Implementation\n- [x] 1.1 Task one\n- [ ] 1.2 Task two\n- [ ] 1.3 Task three\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify metrics
	if data.Summary.ActiveChanges != 1 {
		t.Errorf("Expected ActiveChanges=1, got %d", data.Summary.ActiveChanges)
	}
	if data.Summary.CompletedChanges != 0 {
		t.Errorf("Expected CompletedChanges=0, got %d", data.Summary.CompletedChanges)
	}
	if data.Summary.TotalSpecs != 0 {
		t.Errorf("Expected TotalSpecs=0, got %d", data.Summary.TotalSpecs)
	}
	if data.Summary.TotalTasks != 3 {
		t.Errorf("Expected TotalTasks=3, got %d", data.Summary.TotalTasks)
	}
	if data.Summary.CompletedTasks != 1 {
		t.Errorf("Expected CompletedTasks=1, got %d", data.Summary.CompletedTasks)
	}

	// Verify slices
	if len(data.ActiveChanges) != 1 {
		t.Fatalf("Expected 1 active change, got %d", len(data.ActiveChanges))
	}
	if data.ActiveChanges[0].ID != "test-change" {
		t.Errorf("Expected change ID='test-change', got '%s'", data.ActiveChanges[0].ID)
	}
	if len(data.CompletedChanges) != 0 {
		t.Errorf("Expected 0 completed changes, got %d", len(data.CompletedChanges))
	}
	if len(data.Specs) != 0 {
		t.Errorf("Expected 0 specs, got %d", len(data.Specs))
	}

	t.Log("Project with only active changes correctly handled")
}

// TestCollectData_OnlyCompletedChanges tests a project with only completed changes
func TestCollectData_OnlyCompletedChanges(t *testing.T) {
	// Create a temporary directory with a change that has all tasks completed
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(spectrDir, "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create a test change with all tasks completed
	changeDir := filepath.Join(changesDir, "test-completed")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	// Write a proposal.md
	proposalContent := "# Test Completed\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to write proposal.md: %v", err)
	}

	// Write tasks.md with all tasks completed
	tasksContent := "## 1. Implementation\n- [x] 1.1 Task one\n- [x] 1.2 Task two\n- [x] 1.3 Task three\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify metrics
	if data.Summary.ActiveChanges != 0 {
		t.Errorf("Expected ActiveChanges=0, got %d", data.Summary.ActiveChanges)
	}
	if data.Summary.CompletedChanges != 1 {
		t.Errorf("Expected CompletedChanges=1, got %d", data.Summary.CompletedChanges)
	}
	if data.Summary.TotalSpecs != 0 {
		t.Errorf("Expected TotalSpecs=0, got %d", data.Summary.TotalSpecs)
	}
	// Completed changes don't add to the summary task counts (see dashboard.go line 80-81)
	if data.Summary.TotalTasks != 0 {
		t.Errorf("Expected TotalTasks=0, got %d", data.Summary.TotalTasks)
	}
	if data.Summary.CompletedTasks != 0 {
		t.Errorf("Expected CompletedTasks=0, got %d", data.Summary.CompletedTasks)
	}

	// Verify slices
	if len(data.ActiveChanges) != 0 {
		t.Errorf("Expected 0 active changes, got %d", len(data.ActiveChanges))
	}
	if len(data.CompletedChanges) != 1 {
		t.Fatalf("Expected 1 completed change, got %d", len(data.CompletedChanges))
	}
	if data.CompletedChanges[0].ID != "test-completed" {
		t.Errorf("Expected change ID='test-completed', got '%s'",
			data.CompletedChanges[0].ID)
	}
	if len(data.Specs) != 0 {
		t.Errorf("Expected 0 specs, got %d", len(data.Specs))
	}

	t.Log("Project with only completed changes correctly handled")
}

// TestCollectData_ChangeWithZeroTasks tests a change with no tasks
func TestCollectData_ChangeWithZeroTasks(t *testing.T) {
	// Create a temporary directory with a change that has no tasks
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(spectrDir, "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create a test change with no tasks
	changeDir := filepath.Join(changesDir, "test-no-tasks")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	// Write a proposal.md
	proposalContent := "# Test No Tasks\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to write proposal.md: %v", err)
	}

	// Write tasks.md with no actual tasks (just a header)
	tasksContent := "## 1. Implementation\n\n(No tasks yet)\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify metrics
	// Changes with zero tasks are categorized as "completed" (see dashboard.go line 57-59)
	if data.Summary.ActiveChanges != 0 {
		t.Errorf("Expected ActiveChanges=0, got %d", data.Summary.ActiveChanges)
	}
	if data.Summary.CompletedChanges != 1 {
		t.Errorf("Expected CompletedChanges=1, got %d", data.Summary.CompletedChanges)
	}
	if data.Summary.TotalTasks != 0 {
		t.Errorf("Expected TotalTasks=0, got %d", data.Summary.TotalTasks)
	}
	if data.Summary.CompletedTasks != 0 {
		t.Errorf("Expected CompletedTasks=0, got %d", data.Summary.CompletedTasks)
	}

	// Verify the change is in completed list
	if len(data.ActiveChanges) != 0 {
		t.Errorf("Expected 0 active changes, got %d", len(data.ActiveChanges))
	}
	if len(data.CompletedChanges) != 1 {
		t.Fatalf("Expected 1 completed change, got %d", len(data.CompletedChanges))
	}
	if data.CompletedChanges[0].ID != "test-no-tasks" {
		t.Errorf("Expected change ID='test-no-tasks', got '%s'",
			data.CompletedChanges[0].ID)
	}

	t.Log("Change with zero tasks correctly categorized as completed")
}

// TestCollectDataWithConfig tests the new config-based API
func TestCollectDataWithConfig(t *testing.T) {
	// Create a temporary directory with custom root directory
	tempDir := t.TempDir()
	customRoot := "my-spectr"

	// Create directory structure with custom root
	spectrDir := filepath.Join(tempDir, customRoot)
	changesDir := filepath.Join(spectrDir, "changes")
	specsDir := filepath.Join(spectrDir, "specs")

	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create a test change
	changeDir := filepath.Join(changesDir, "test-change")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	proposalContent := "# Test Change\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to write proposal.md: %v", err)
	}

	tasksContent := "## 1. Implementation\n- [x] 1.1 Task one\n- [ ] 1.2 Task two\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}

	// Create a test spec
	specDir := filepath.Join(specsDir, "test-spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("Failed to create spec directory: %v", err)
	}

	specContent := `# Test Spec

## Requirements

### Requirement: Feature One
Test requirement one

### Requirement: Feature Two
Test requirement two
`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to write spec.md: %v", err)
	}

	// Create config with custom root
	cfg := &config.Config{
		RootDir:     customRoot,
		ProjectRoot: tempDir,
	}

	// Collect data using the new API
	data, err := CollectDataWithConfig(cfg)
	if err != nil {
		t.Fatalf("CollectDataWithConfig failed: %v", err)
	}

	// Verify data was collected correctly
	if data == nil {
		t.Fatal("CollectDataWithConfig returned nil data")
	}

	// Verify change was found
	if data.Summary.ActiveChanges != 1 {
		t.Errorf("Expected 1 active change, got %d", data.Summary.ActiveChanges)
	}
	if len(data.ActiveChanges) != 1 {
		t.Fatalf("Expected 1 active change in slice, got %d", len(data.ActiveChanges))
	}
	if data.ActiveChanges[0].ID != "test-change" {
		t.Errorf("Expected change ID 'test-change', got '%s'", data.ActiveChanges[0].ID)
	}
	if data.ActiveChanges[0].Progress.Total != 2 {
		t.Errorf("Expected 2 total tasks, got %d", data.ActiveChanges[0].Progress.Total)
	}
	if data.ActiveChanges[0].Progress.Completed != 1 {
		t.Errorf("Expected 1 completed task, got %d", data.ActiveChanges[0].Progress.Completed)
	}

	// Verify spec was found
	if data.Summary.TotalSpecs != 1 {
		t.Errorf("Expected 1 spec, got %d", data.Summary.TotalSpecs)
	}
	if len(data.Specs) != 1 {
		t.Fatalf("Expected 1 spec in slice, got %d", len(data.Specs))
	}
	if data.Specs[0].ID != "test-spec" {
		t.Errorf("Expected spec ID 'test-spec', got '%s'", data.Specs[0].ID)
	}
	if data.Specs[0].RequirementCount != 2 {
		t.Errorf("Expected 2 requirements, got %d", data.Specs[0].RequirementCount)
	}

	t.Log("CollectDataWithConfig correctly uses custom root directory")
}

// TestCollectData_BackwardCompatibility verifies that the deprecated CollectData
// function still works and delegates to CollectDataWithConfig correctly
func TestCollectData_BackwardCompatibility(t *testing.T) {
	// Create a temporary directory with standard "spectr" root
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(spectrDir, "changes")
	specsDir := filepath.Join(spectrDir, "specs")

	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}

	// Create a test change
	changeDir := filepath.Join(changesDir, "compat-test")
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	proposalContent := "# Compat Test\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf("Failed to write proposal.md: %v", err)
	}

	tasksContent := "## 1. Implementation\n- [x] 1.1 Done\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf("Failed to write tasks.md: %v", err)
	}

	// Use the deprecated API
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify it still works
	if data == nil {
		t.Fatal("CollectData returned nil data")
	}

	if len(data.CompletedChanges) != 1 {
		t.Errorf("Expected 1 completed change, got %d", len(data.CompletedChanges))
	}
	if len(data.CompletedChanges) > 0 && data.CompletedChanges[0].ID != "compat-test" {
		t.Errorf("Expected change ID 'compat-test', got '%s'", data.CompletedChanges[0].ID)
	}

	t.Log("Deprecated CollectData function maintains backward compatibility")
}
