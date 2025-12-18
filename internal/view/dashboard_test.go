package view

import (
	"os"
	"path/filepath"
	"testing"
)

//nolint:revive // cognitive-complexity justified for comprehensive testing
func TestCollectData(t *testing.T) {
	// Use the current project as test data
	projectPath := filepath.Join("..", "..")

	// Verify the project exists
	if _, err := os.Stat(filepath.Join(projectPath, "spectr")); err != nil {
		t.Skipf(
			"Skipping test: spectr directory not found at %s",
			projectPath,
		)
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
		t.Errorf(
			"TotalSpecs should be >= 0, got %d",
			data.Summary.TotalSpecs,
		)
	}
	if data.Summary.TotalRequirements < 0 {
		t.Errorf(
			"TotalRequirements should be >= 0, got %d",
			data.Summary.TotalRequirements,
		)
	}
	if data.Summary.ActiveChanges < 0 {
		t.Errorf(
			"ActiveChanges should be >= 0, got %d",
			data.Summary.ActiveChanges,
		)
	}
	if data.Summary.CompletedChanges < 0 {
		t.Errorf(
			"CompletedChanges should be >= 0, got %d",
			data.Summary.CompletedChanges,
		)
	}

	// Verify counts match slice lengths
	if data.Summary.ActiveChanges != len(
		data.ActiveChanges,
	) {
		t.Errorf(
			"ActiveChanges count mismatch: summary=%d, slice=%d",
			data.Summary.ActiveChanges,
			len(data.ActiveChanges),
		)
	}
	if data.Summary.CompletedChanges != len(
		data.CompletedChanges,
	) {
		t.Errorf(
			"CompletedChanges count mismatch: summary=%d, slice=%d",
			data.Summary.CompletedChanges,
			len(data.CompletedChanges),
		)
	}
	if data.Summary.TotalSpecs != len(
		data.Specs,
	) {
		t.Errorf(
			"TotalSpecs count mismatch: summary=%d, slice=%d",
			data.Summary.TotalSpecs,
			len(data.Specs),
		)
	}

	// Verify sorting of active changes (ascending percentage)
	for i := 1; i < len(data.ActiveChanges); i++ {
		prev := data.ActiveChanges[i-1]
		curr := data.ActiveChanges[i]
		if prev.Progress.Percentage > curr.Progress.Percentage {
			t.Errorf(
				"ActiveChanges not sorted by percentage: %s (%d%%) > %s (%d%%)",
				prev.ID,
				prev.Progress.Percentage,
				curr.ID,
				curr.Progress.Percentage,
			)
		}
	}

	// Verify sorting of specs (descending requirement count)
	for i := 1; i < len(data.Specs); i++ {
		prev := data.Specs[i-1]
		curr := data.Specs[i]
		if prev.RequirementCount < curr.RequirementCount {
			t.Errorf(
				"Specs not sorted by requirement count: %s (%d) < %s (%d)",
				prev.ID,
				prev.RequirementCount,
				curr.ID,
				curr.RequirementCount,
			)
		}
	}

	// Verify sorting of completed changes (alphabetical)
	for i := 1; i < len(data.CompletedChanges); i++ {
		prev := data.CompletedChanges[i-1]
		curr := data.CompletedChanges[i]
		if prev.ID > curr.ID {
			t.Errorf(
				"CompletedChanges not sorted alphabetically: %s > %s",
				prev.ID,
				curr.ID,
			)
		}
	}

	t.Log(
		"Dashboard data collected successfully:",
	)
	t.Logf(
		"  Specs: %d (with %d requirements)",
		data.Summary.TotalSpecs,
		data.Summary.TotalRequirements,
	)
	t.Logf(
		"  Active Changes: %d",
		data.Summary.ActiveChanges,
	)
	for _, change := range data.ActiveChanges {
		t.Logf(
			"    - %s: %d/%d tasks (%d%%)",
			change.ID,
			change.Progress.Completed,
			change.Progress.Total,
			change.Progress.Percentage,
		)
	}
	t.Logf(
		"  Completed Changes: %d",
		data.Summary.CompletedChanges,
	)
	for _, change := range data.CompletedChanges {
		t.Logf("    - %s", change.ID)
	}
	t.Logf(
		"  Tasks: %d/%d completed",
		data.Summary.CompletedTasks,
		data.Summary.TotalTasks,
	)
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
		{
			"Two thirds",
			2,
			3,
			67,
		}, // 66.67% rounds to 67
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentage(
				tt.completed,
				tt.total,
			)
			if result != tt.expected {
				t.Errorf(
					"calculatePercentage(%d, %d) = %d; want %d",
					tt.completed,
					tt.total,
					result,
					tt.expected,
				)
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
		t.Fatalf(
			"Failed to create spectr directory: %v",
			err,
		)
	}

	// Create empty subdirectories
	if err := os.MkdirAll(filepath.Join(spectrDir, "changes"), 0755); err != nil {
		t.Fatalf(
			"Failed to create changes directory: %v",
			err,
		)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf(
			"Failed to create specs directory: %v",
			err,
		)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf(
			"CollectData failed on empty project: %v",
			err,
		)
	}

	// Verify all metrics are zero
	if data.Summary.TotalSpecs != 0 {
		t.Errorf(
			"Expected TotalSpecs=0, got %d",
			data.Summary.TotalSpecs,
		)
	}
	if data.Summary.TotalRequirements != 0 {
		t.Errorf(
			"Expected TotalRequirements=0, got %d",
			data.Summary.TotalRequirements,
		)
	}
	if data.Summary.ActiveChanges != 0 {
		t.Errorf(
			"Expected ActiveChanges=0, got %d",
			data.Summary.ActiveChanges,
		)
	}
	if data.Summary.CompletedChanges != 0 {
		t.Errorf(
			"Expected CompletedChanges=0, got %d",
			data.Summary.CompletedChanges,
		)
	}
	if data.Summary.TotalTasks != 0 {
		t.Errorf(
			"Expected TotalTasks=0, got %d",
			data.Summary.TotalTasks,
		)
	}
	if data.Summary.CompletedTasks != 0 {
		t.Errorf(
			"Expected CompletedTasks=0, got %d",
			data.Summary.CompletedTasks,
		)
	}

	// Verify slices are initialized but empty
	if len(data.ActiveChanges) != 0 {
		t.Errorf(
			"Expected empty ActiveChanges, got %d items",
			len(data.ActiveChanges),
		)
	}
	if len(data.CompletedChanges) != 0 {
		t.Errorf(
			"Expected empty CompletedChanges, got %d items",
			len(data.CompletedChanges),
		)
	}
	if len(data.Specs) != 0 {
		t.Errorf(
			"Expected empty Specs, got %d items",
			len(data.Specs),
		)
	}

	t.Log(
		"Empty project correctly handled with all zero metrics",
	)
}

// TestCollectData_OnlyActiveChanges tests a project with only active changes
func TestCollectData_OnlyActiveChanges(
	t *testing.T,
) {
	// Create a temporary directory with only active changes
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(
		spectrDir,
		"changes",
	)
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf(
			"Failed to create changes directory: %v",
			err,
		)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf(
			"Failed to create specs directory: %v",
			err,
		)
	}

	// Create a test change with tasks
	changeDir := filepath.Join(
		changesDir,
		"test-change",
	)
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf(
			"Failed to create change directory: %v",
			err,
		)
	}

	// Write a proposal.md
	proposalContent := "# Test Change\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf(
			"Failed to write proposal.md: %v",
			err,
		)
	}

	// Write tasks.md with some tasks
	tasksContent := "## 1. Implementation\n- [x] 1.1 Task one\n- [ ] 1.2 Task two\n- [ ] 1.3 Task three\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf(
			"Failed to write tasks.md: %v",
			err,
		)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify metrics
	if data.Summary.ActiveChanges != 1 {
		t.Errorf(
			"Expected ActiveChanges=1, got %d",
			data.Summary.ActiveChanges,
		)
	}
	if data.Summary.CompletedChanges != 0 {
		t.Errorf(
			"Expected CompletedChanges=0, got %d",
			data.Summary.CompletedChanges,
		)
	}
	if data.Summary.TotalSpecs != 0 {
		t.Errorf(
			"Expected TotalSpecs=0, got %d",
			data.Summary.TotalSpecs,
		)
	}
	if data.Summary.TotalTasks != 3 {
		t.Errorf(
			"Expected TotalTasks=3, got %d",
			data.Summary.TotalTasks,
		)
	}
	if data.Summary.CompletedTasks != 1 {
		t.Errorf(
			"Expected CompletedTasks=1, got %d",
			data.Summary.CompletedTasks,
		)
	}

	// Verify slices
	if len(data.ActiveChanges) != 1 {
		t.Fatalf(
			"Expected 1 active change, got %d",
			len(data.ActiveChanges),
		)
	}
	if data.ActiveChanges[0].ID != "test-change" {
		t.Errorf(
			"Expected change ID='test-change', got '%s'",
			data.ActiveChanges[0].ID,
		)
	}
	if len(data.CompletedChanges) != 0 {
		t.Errorf(
			"Expected 0 completed changes, got %d",
			len(data.CompletedChanges),
		)
	}
	if len(data.Specs) != 0 {
		t.Errorf(
			"Expected 0 specs, got %d",
			len(data.Specs),
		)
	}

	t.Log(
		"Project with only active changes correctly handled",
	)
}

// TestCollectData_OnlyCompletedChanges tests a project with only completed changes
func TestCollectData_OnlyCompletedChanges(
	t *testing.T,
) {
	// Create a temporary directory with a change that has all tasks completed
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(
		spectrDir,
		"changes",
	)
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf(
			"Failed to create changes directory: %v",
			err,
		)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf(
			"Failed to create specs directory: %v",
			err,
		)
	}

	// Create a test change with all tasks completed
	changeDir := filepath.Join(
		changesDir,
		"test-completed",
	)
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf(
			"Failed to create change directory: %v",
			err,
		)
	}

	// Write a proposal.md
	proposalContent := "# Test Completed\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf(
			"Failed to write proposal.md: %v",
			err,
		)
	}

	// Write tasks.md with all tasks completed
	tasksContent := "## 1. Implementation\n- [x] 1.1 Task one\n- [x] 1.2 Task two\n- [x] 1.3 Task three\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf(
			"Failed to write tasks.md: %v",
			err,
		)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify metrics
	if data.Summary.ActiveChanges != 0 {
		t.Errorf(
			"Expected ActiveChanges=0, got %d",
			data.Summary.ActiveChanges,
		)
	}
	if data.Summary.CompletedChanges != 1 {
		t.Errorf(
			"Expected CompletedChanges=1, got %d",
			data.Summary.CompletedChanges,
		)
	}
	if data.Summary.TotalSpecs != 0 {
		t.Errorf(
			"Expected TotalSpecs=0, got %d",
			data.Summary.TotalSpecs,
		)
	}
	// Completed changes don't add to the summary task counts (see dashboard.go line 80-81)
	if data.Summary.TotalTasks != 0 {
		t.Errorf(
			"Expected TotalTasks=0, got %d",
			data.Summary.TotalTasks,
		)
	}
	if data.Summary.CompletedTasks != 0 {
		t.Errorf(
			"Expected CompletedTasks=0, got %d",
			data.Summary.CompletedTasks,
		)
	}

	// Verify slices
	if len(data.ActiveChanges) != 0 {
		t.Errorf(
			"Expected 0 active changes, got %d",
			len(data.ActiveChanges),
		)
	}
	if len(data.CompletedChanges) != 1 {
		t.Fatalf(
			"Expected 1 completed change, got %d",
			len(data.CompletedChanges),
		)
	}
	if data.CompletedChanges[0].ID != "test-completed" {
		t.Errorf(
			"Expected change ID='test-completed', got '%s'",
			data.CompletedChanges[0].ID,
		)
	}
	if len(data.Specs) != 0 {
		t.Errorf(
			"Expected 0 specs, got %d",
			len(data.Specs),
		)
	}

	t.Log(
		"Project with only completed changes correctly handled",
	)
}

// TestCollectData_ChangeWithZeroTasks tests a change with no tasks
func TestCollectData_ChangeWithZeroTasks(
	t *testing.T,
) {
	// Create a temporary directory with a change that has no tasks
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr")
	changesDir := filepath.Join(
		spectrDir,
		"changes",
	)
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf(
			"Failed to create changes directory: %v",
			err,
		)
	}
	if err := os.MkdirAll(filepath.Join(spectrDir, "specs"), 0755); err != nil {
		t.Fatalf(
			"Failed to create specs directory: %v",
			err,
		)
	}

	// Create a test change with no tasks
	changeDir := filepath.Join(
		changesDir,
		"test-no-tasks",
	)
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf(
			"Failed to create change directory: %v",
			err,
		)
	}

	// Write a proposal.md
	proposalContent := "# Test No Tasks\n\n## Why\nTest\n\n## What Changes\n- Test\n"
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		t.Fatalf(
			"Failed to write proposal.md: %v",
			err,
		)
	}

	// Write tasks.md with no actual tasks (just a header)
	tasksContent := "## 1. Implementation\n\n(No tasks yet)\n"
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		t.Fatalf(
			"Failed to write tasks.md: %v",
			err,
		)
	}

	// Collect data
	data, err := CollectData(tempDir)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Verify metrics
	// Changes with zero tasks are categorized as "completed" (see dashboard.go line 57-59)
	if data.Summary.ActiveChanges != 0 {
		t.Errorf(
			"Expected ActiveChanges=0, got %d",
			data.Summary.ActiveChanges,
		)
	}
	if data.Summary.CompletedChanges != 1 {
		t.Errorf(
			"Expected CompletedChanges=1, got %d",
			data.Summary.CompletedChanges,
		)
	}
	if data.Summary.TotalTasks != 0 {
		t.Errorf(
			"Expected TotalTasks=0, got %d",
			data.Summary.TotalTasks,
		)
	}
	if data.Summary.CompletedTasks != 0 {
		t.Errorf(
			"Expected CompletedTasks=0, got %d",
			data.Summary.CompletedTasks,
		)
	}

	// Verify the change is in completed list
	if len(data.ActiveChanges) != 0 {
		t.Errorf(
			"Expected 0 active changes, got %d",
			len(data.ActiveChanges),
		)
	}
	if len(data.CompletedChanges) != 1 {
		t.Fatalf(
			"Expected 1 completed change, got %d",
			len(data.CompletedChanges),
		)
	}
	if data.CompletedChanges[0].ID != "test-no-tasks" {
		t.Errorf(
			"Expected change ID='test-no-tasks', got '%s'",
			data.CompletedChanges[0].ID,
		)
	}

	t.Log(
		"Change with zero tasks correctly categorized as completed",
	)
}

func TestCalculatePercentage_EdgeCases(
	t *testing.T,
) {
	tests := []struct {
		name      string
		completed int
		total     int
		expected  int
	}{
		// Edge cases for rounding (values now ROUND to nearest integer)
		{"1 of 16 (6.25%)", 1, 16, 6},
		{"5 of 16 (31.25%)", 5, 16, 31},
		{"8 of 16 (exactly 50%)", 8, 16, 50},
		{
			"15 of 16 (93.75%)",
			15,
			16,
			94,
		}, // rounds up
		{"16 of 16 (exactly 100%)", 16, 16, 100},

		// More rounding edge cases
		{"1 of 7 (14.29%)", 1, 7, 14},
		{
			"3 of 7 (42.86%)",
			3,
			7,
			43,
		}, // rounds up
		{
			"6 of 7 (85.71%)",
			6,
			7,
			86,
		}, // rounds up

		// Large numbers
		{"500 of 1000", 500, 1000, 50},
		{
			"999 of 1000",
			999,
			1000,
			100,
		}, // 99.9% rounds to 100
		{"1000 of 1000", 1000, 1000, 100},

		// Small edge cases
		{"0 of 1", 0, 1, 0},
		{"1 of 1", 1, 1, 100},
		{"0 of 2", 0, 2, 0},
		{"1 of 2", 1, 2, 50},
		{"2 of 2", 2, 2, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentage(
				tt.completed,
				tt.total,
			)
			if result != tt.expected {
				t.Errorf(
					"calculatePercentage(%d, %d) = %d; want %d (%.2f%%)",
					tt.completed,
					tt.total,
					result,
					tt.expected,
					float64(
						tt.completed,
					)/float64(
						tt.total,
					)*100,
				)
			}
		})
	}
}

func TestCalculatePercentage_FixtureLike(
	t *testing.T,
) {
	// Test with values matching the 16-task fixture
	// Values should be ROUNDED to nearest integer (not truncated)
	testCases := []struct {
		completed int
		expected  int
	}{
		{0, 0},    // All pending
		{1, 6},    // 6.25% -> rounds to 6
		{2, 13},   // 12.5% -> rounds to 13
		{3, 19},   // 18.75% -> rounds to 19
		{4, 25},   // 25% -> exactly 25
		{5, 31},   // 31.25% -> rounds to 31
		{6, 38},   // 37.5% -> rounds to 38
		{7, 44},   // 43.75% -> rounds to 44
		{8, 50},   // 50% -> exactly 50
		{9, 56},   // 56.25% -> rounds to 56
		{10, 63},  // 62.5% -> rounds to 63
		{11, 69},  // 68.75% -> rounds to 69
		{12, 75},  // 75% -> exactly 75
		{13, 81},  // 81.25% -> rounds to 81
		{14, 88},  // 87.5% -> rounds to 88
		{15, 94},  // 93.75% -> rounds to 94
		{16, 100}, // 100%
	}

	total := 16
	for _, tc := range testCases {
		result := calculatePercentage(
			tc.completed,
			total,
		)
		if result != tc.expected {
			t.Errorf(
				"calculatePercentage(%d, %d) = %d; want %d",
				tc.completed,
				total,
				result,
				tc.expected,
			)
		}
	}
}
