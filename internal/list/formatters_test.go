package list

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestFormatChangesText(t *testing.T) {
	tests := []struct {
		name     string
		changes  []ChangeInfo
		expected []string
	}{
		{
			name:     "Empty list",
			changes:  make([]ChangeInfo, 0),
			expected: []string{"No items found"},
		},
		{
			name: "Single change",
			changes: []ChangeInfo{
				{
					ID:         "add-feature",
					Title:      "Add Feature",
					DeltaCount: 2,
					TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
				},
			},
			expected: []string{"add-feature  3/5 tasks"},
		},
		{
			name: "Multiple changes sorted",
			changes: []ChangeInfo{
				{
					ID:         "update-docs",
					Title:      "Update Docs",
					DeltaCount: 1,
					TaskStatus: parsers.TaskStatus{Total: 2, Completed: 1},
				},
				{
					ID:         "add-feature",
					Title:      "Add Feature",
					DeltaCount: 2,
					TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
				},
				{
					ID:         "fix-bug",
					Title:      "Fix Bug",
					DeltaCount: 1,
					TaskStatus: parsers.TaskStatus{Total: 3, Completed: 3},
				},
			},
			expected: []string{
				"add-feature  3/5 tasks",
				"fix-bug      3/3 tasks",
				"update-docs  1/2 tasks",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatChangesText(tt.changes)
			lines := strings.Split(result, "\n")
			if len(lines) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d", len(tt.expected), len(lines))
			}
			for i, expected := range tt.expected {
				if i < len(lines) && lines[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
				}
			}
		})
	}
}

func TestFormatChangesLong(t *testing.T) {
	changes := []ChangeInfo{
		{
			ID:         "update-docs",
			Title:      "Update Docs",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{Total: 2, Completed: 1},
		},
		{
			ID:         "add-feature",
			Title:      "Add Feature",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
		},
	}

	result := FormatChangesLong(changes)
	lines := strings.Split(result, "\n")

	// Should be sorted alphabetically
	if !strings.Contains(lines[0], "add-feature") {
		t.Error("First line should contain add-feature")
	}
	if !strings.Contains(lines[1], "update-docs") {
		t.Error("Second line should contain update-docs")
	}

	// Check format includes all components
	if !strings.Contains(lines[0], "Add Feature") {
		t.Error("Should contain title")
	}
	if !strings.Contains(lines[0], "[deltas 2]") {
		t.Error("Should contain delta count")
	}
	if !strings.Contains(lines[0], "[tasks 3/5]") {
		t.Error("Should contain task status")
	}
}

func TestFormatChangesJSON(t *testing.T) {
	changes := []ChangeInfo{
		{
			ID:         "update-docs",
			Title:      "Update Docs",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{Total: 2, Completed: 1},
		},
		{
			ID:         "add-feature",
			Title:      "Add Feature",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
		},
	}

	result, err := FormatChangesJSON(changes)
	if err != nil {
		t.Fatalf("FormatChangesJSON failed: %v", err)
	}

	// Parse JSON to verify structure
	var parsed []ChangeInfo
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check sorting
	if parsed[0].ID != "add-feature" {
		t.Error("First item should be add-feature (sorted)")
	}
	if parsed[1].ID != "update-docs" {
		t.Error("Second item should be update-docs (sorted)")
	}

	// Verify data integrity
	if parsed[0].Title != "Add Feature" || parsed[0].DeltaCount != 2 {
		t.Error("Data mismatch in JSON output")
	}
}

func TestFormatChangesJSON_Empty(t *testing.T) {
	result, err := FormatChangesJSON(make([]ChangeInfo, 0))
	if err != nil {
		t.Fatalf("FormatChangesJSON failed: %v", err)
	}
	if result != "[]" {
		t.Errorf("Expected '[]', got %q", result)
	}
}

func TestFormatSpecsText(t *testing.T) {
	tests := []struct {
		name     string
		specs    []SpecInfo
		expected []string
	}{
		{
			name:     "Empty list",
			specs:    make([]SpecInfo, 0),
			expected: []string{"No items found"},
		},
		{
			name: "Single spec",
			specs: []SpecInfo{
				{ID: "auth", Title: "Authentication", RequirementCount: 5},
			},
			expected: []string{"auth"},
		},
		{
			name: "Multiple specs sorted",
			specs: []SpecInfo{
				{ID: "database", Title: "Database", RequirementCount: 8},
				{ID: "api", Title: "API", RequirementCount: 12},
				{ID: "auth", Title: "Authentication", RequirementCount: 5},
			},
			expected: []string{"api", "auth", "database"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSpecsText(tt.specs)
			lines := strings.Split(result, "\n")
			if len(lines) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d", len(tt.expected), len(lines))
			}
			for i, expected := range tt.expected {
				if i < len(lines) && lines[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
				}
			}
		})
	}
}

func TestFormatSpecsLong(t *testing.T) {
	specs := []SpecInfo{
		{ID: "database", Title: "Database", RequirementCount: 8},
		{ID: "api", Title: "API", RequirementCount: 12},
	}

	result := FormatSpecsLong(specs)
	lines := strings.Split(result, "\n")

	// Should be sorted alphabetically
	if !strings.Contains(lines[0], "api") {
		t.Error("First line should contain api")
	}
	if !strings.Contains(lines[1], "database") {
		t.Error("Second line should contain database")
	}

	// Check format includes all components
	if !strings.Contains(lines[0], "API") {
		t.Error("Should contain title")
	}
	if !strings.Contains(lines[0], "[requirements 12]") {
		t.Error("Should contain requirement count")
	}
}

func TestFormatSpecsJSON(t *testing.T) {
	specs := []SpecInfo{
		{ID: "database", Title: "Database", RequirementCount: 8},
		{ID: "api", Title: "API", RequirementCount: 12},
	}

	result, err := FormatSpecsJSON(specs)
	if err != nil {
		t.Fatalf("FormatSpecsJSON failed: %v", err)
	}

	// Parse JSON to verify structure
	var parsed []SpecInfo
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check sorting
	if parsed[0].ID != "api" {
		t.Error("First item should be api (sorted)")
	}
	if parsed[1].ID != "database" {
		t.Error("Second item should be database (sorted)")
	}

	// Verify data integrity
	if parsed[0].Title != "API" || parsed[0].RequirementCount != 12 {
		t.Error("Data mismatch in JSON output")
	}
}

func TestFormatSpecsJSON_Empty(t *testing.T) {
	result, err := FormatSpecsJSON(make([]SpecInfo, 0))
	if err != nil {
		t.Fatalf("FormatSpecsJSON failed: %v", err)
	}
	if result != "[]" {
		t.Errorf("Expected '[]', got %q", result)
	}
}

func TestFormatAllText(t *testing.T) {
	tests := []struct {
		name     string
		items    ItemList
		expected []string
	}{
		{
			name:     "Empty list",
			items:    ItemList{},
			expected: []string{"No items found"},
		},
		{
			name: "Mixed items sorted",
			items: ItemList{
				NewChangeItem(ChangeInfo{
					ID:         "update-docs",
					Title:      "Update Docs",
					DeltaCount: 1,
					TaskStatus: parsers.TaskStatus{Total: 2, Completed: 1},
				}),
				NewSpecItem(SpecInfo{
					ID:               "api",
					Title:            "API",
					RequirementCount: 12,
				}),
				NewChangeItem(ChangeInfo{
					ID:         "add-feature",
					Title:      "Add Feature",
					DeltaCount: 2,
					TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
				}),
			},
			expected: []string{
				"add-feature  [CHANGE]  3/5 tasks",
				"api          [SPEC]    12 requirements",
				"update-docs  [CHANGE]  1/2 tasks",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatAllText(tt.items)
			lines := strings.Split(result, "\n")
			if len(lines) != len(tt.expected) {
				t.Errorf("Expected %d lines, got %d", len(tt.expected), len(lines))
			}
			for i, expected := range tt.expected {
				if i < len(lines) && lines[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
				}
			}
		})
	}
}

func TestFormatAllLong(t *testing.T) {
	items := ItemList{
		NewChangeItem(ChangeInfo{
			ID:         "add-feature",
			Title:      "Add Feature",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
		}),
		NewSpecItem(SpecInfo{
			ID:               "api",
			Title:            "API",
			RequirementCount: 12,
		}),
	}

	result := FormatAllLong(items)
	lines := strings.Split(result, "\n")

	// Should be sorted alphabetically
	if !strings.Contains(lines[0], "add-feature") {
		t.Error("First line should contain add-feature")
	}
	if !strings.Contains(lines[1], "api") {
		t.Error("Second line should contain api")
	}

	// Check format includes type indicators
	if !strings.Contains(lines[0], "[CHANGE]") {
		t.Error("Should contain [CHANGE] indicator")
	}
	if !strings.Contains(lines[1], "[SPEC]") {
		t.Error("Should contain [SPEC] indicator")
	}

	// Check format includes all components for change
	if !strings.Contains(lines[0], "Add Feature") {
		t.Error("Should contain change title")
	}
	if !strings.Contains(lines[0], "[deltas 2]") {
		t.Error("Should contain delta count")
	}
	if !strings.Contains(lines[0], "[tasks 3/5]") {
		t.Error("Should contain task status")
	}

	// Check format includes all components for spec
	if !strings.Contains(lines[1], "API") {
		t.Error("Should contain spec title")
	}
	if !strings.Contains(lines[1], "[requirements 12]") {
		t.Error("Should contain requirement count")
	}
}

func TestFormatAllJSON(t *testing.T) {
	items := ItemList{
		NewChangeItem(ChangeInfo{
			ID:         "update-docs",
			Title:      "Update Docs",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{Total: 2, Completed: 1},
		}),
		NewSpecItem(SpecInfo{
			ID:               "api",
			Title:            "API",
			RequirementCount: 12,
		}),
	}

	result, err := FormatAllJSON(items)
	if err != nil {
		t.Fatalf("FormatAllJSON failed: %v", err)
	}

	// Parse JSON to verify structure
	var parsed ItemList
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check sorting
	if parsed[0].ID() != "api" {
		t.Error("First item should be api (sorted)")
	}
	if parsed[1].ID() != "update-docs" {
		t.Error("Second item should be update-docs (sorted)")
	}

	// Verify data integrity
	if parsed[0].Type != ItemTypeSpec || parsed[0].Spec == nil {
		t.Error("First item should be spec type")
	}
	if parsed[1].Type != ItemTypeChange || parsed[1].Change == nil {
		t.Error("Second item should be change type")
	}
}

func TestFormatAllJSON_Empty(t *testing.T) {
	result, err := FormatAllJSON(ItemList{})
	if err != nil {
		t.Fatalf("FormatAllJSON failed: %v", err)
	}
	if result != "[]" {
		t.Errorf("Expected '[]', got %q", result)
	}
}
