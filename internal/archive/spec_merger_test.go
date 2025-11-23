package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeSpec_AddedOnly_NewSpec(t *testing.T) {
	tmpDir := t.TempDir()

	// Create delta spec with ADDED requirements
	deltaContent := `# Delta Spec

## ADDED Requirements

### Requirement: New Feature
The system SHALL support new functionality.

#### Scenario: Basic usage
- **WHEN** user performs action
- **THEN** feature responds
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	err := os.WriteFile(deltaPath, []byte(deltaContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Merge with non-existent base spec
	basePath := filepath.Join(tmpDir, "base.md")
	merged, counts, err := MergeSpec(basePath, deltaPath, false)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	if counts.Added != 1 {
		t.Errorf("Expected 1 added requirement, got %d", counts.Added)
	}

	if !strings.Contains(merged, "# ") {
		t.Error("Merged spec should contain H1 header")
	}
	if !strings.Contains(merged, "## Requirements") {
		t.Error("Merged spec should contain Requirements section")
	}
	if !strings.Contains(merged, "### Requirement: New Feature") {
		t.Error("Merged spec should contain new requirement")
	}
}

func TestMergeSpec_PreservesContentWhenNoRequirementsSection(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec without a Requirements section
	baseContent := `# Test Spec

## Purpose
Original purpose text.

## Notes
Additional notes.
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create delta spec with ADDED requirements
	deltaContent := `# Delta Spec

## ADDED Requirements

### Requirement: Added Feature
The system SHALL add a new capability.

#### Scenario: Added behavior
- **WHEN** a new action occurs
- **THEN** the new capability responds
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	merged, _, err := MergeSpec(basePath, deltaPath, true)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	if !strings.Contains(merged, "## Purpose\nOriginal purpose text.") {
		t.Error("Merged spec should preserve purpose section when missing requirements")
	}
	if !strings.Contains(merged, "## Notes\nAdditional notes.") {
		t.Error("Merged spec should preserve following sections when missing requirements")
	}
	if !strings.Contains(merged, "### Requirement: Added Feature") {
		t.Error("Merged spec should include added requirement")
	}

	purposeIdx := strings.Index(merged, "## Purpose")
	notesIdx := strings.Index(merged, "## Notes")
	reqIdx := strings.Index(merged, "## Requirements")

	if purposeIdx == -1 || notesIdx == -1 || reqIdx == -1 {
		t.Fatal("Merged spec missing expected sections")
	}

	if purposeIdx >= notesIdx || notesIdx >= reqIdx {
		t.Error("Requirements section should be appended after original content when missing")
	}
}

func TestMergeSpec_ModifiedOnly_ExistingSpec(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec
	baseContent := `# Test Spec

## Requirements

### Requirement: Existing Feature
The system SHALL have original behavior.

#### Scenario: Original scenario
- **WHEN** original action
- **THEN** original result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create delta spec with MODIFIED requirement
	deltaContent := `# Delta Spec

## MODIFIED Requirements

### Requirement: Existing Feature
The system SHALL have updated behavior.

#### Scenario: Updated scenario
- **WHEN** updated action
- **THEN** updated result
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	merged, counts, err := MergeSpec(basePath, deltaPath, true)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	if counts.Modified != 1 {
		t.Errorf("Expected 1 modified requirement, got %d", counts.Modified)
	}

	if !strings.Contains(merged, "updated behavior") {
		t.Error("Merged spec should contain updated content")
	}
	if strings.Contains(merged, "original behavior") {
		t.Error("Merged spec should not contain original content")
	}
}

func TestMergeSpec_RemovedOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec with two requirements
	baseContent := `# Test Spec

## Requirements

### Requirement: Keep This
Should remain.

#### Scenario: Test
- **WHEN** something
- **THEN** result

### Requirement: Remove This
Should be removed.

#### Scenario: Test
- **WHEN** something
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create delta spec with REMOVED requirement
	deltaContent := `# Delta Spec

## REMOVED Requirements

### Requirement: Remove This
**Reason**: No longer needed
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	merged, counts, err := MergeSpec(basePath, deltaPath, true)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	if counts.Removed != 1 {
		t.Errorf("Expected 1 removed requirement, got %d", counts.Removed)
	}

	if !strings.Contains(merged, "Keep This") {
		t.Error("Merged spec should contain kept requirement")
	}
	if strings.Contains(merged, "Remove This") {
		t.Error("Merged spec should not contain removed requirement")
	}
}

func TestMergeSpec_RenamedOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec
	baseContent := `# Test Spec

## Requirements

### Requirement: Old Name
Content here.

#### Scenario: Test
- **WHEN** something
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create delta spec with RENAMED requirement
	deltaContent := `# Delta Spec

## RENAMED Requirements

- FROM: ` + "`" + `### Requirement: Old Name` + "`" + `
- TO: ` + "`" + `### Requirement: New Name` + "`" + `
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	merged, counts, err := MergeSpec(basePath, deltaPath, true)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	if counts.Renamed != 1 {
		t.Errorf("Expected 1 renamed requirement, got %d", counts.Renamed)
	}

	if !strings.Contains(merged, "### Requirement: New Name") {
		t.Error("Merged spec should contain new name")
	}
	if strings.Contains(merged, "### Requirement: Old Name") {
		t.Error("Merged spec should not contain old name")
	}
}

func TestMergeSpec_AllOperations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec with multiple requirements
	baseContent := `# Test Spec

## Requirements

### Requirement: Keep Unchanged
This stays the same.

#### Scenario: Test
- **WHEN** action
- **THEN** result

### Requirement: Modify This
Original content.

#### Scenario: Original
- **WHEN** original
- **THEN** result

### Requirement: Remove This
Will be removed.

#### Scenario: Test
- **WHEN** action
- **THEN** result

### Requirement: Rename This
Will be renamed.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create delta spec with all operation types
	deltaContent := `# Delta Spec

## ADDED Requirements

### Requirement: Brand New
The system SHALL support brand new feature.

#### Scenario: New scenario
- **WHEN** new action
- **THEN** new result

## MODIFIED Requirements

### Requirement: Modify This
Updated content.

#### Scenario: Updated
- **WHEN** updated
- **THEN** new result

## REMOVED Requirements

### Requirement: Remove This
**Reason**: No longer needed

## RENAMED Requirements

- FROM: ` + "`" + `### Requirement: Rename This` + "`" + `
- TO: ` + "`" + `### Requirement: Renamed Feature` + "`" + `
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	merged, counts, err := MergeSpec(basePath, deltaPath, true)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	if counts.Added != 1 {
		t.Errorf("Expected 1 added, got %d", counts.Added)
	}
	if counts.Modified != 1 {
		t.Errorf("Expected 1 modified, got %d", counts.Modified)
	}
	if counts.Removed != 1 {
		t.Errorf("Expected 1 removed, got %d", counts.Removed)
	}
	if counts.Renamed != 1 {
		t.Errorf("Expected 1 renamed, got %d", counts.Renamed)
	}

	// Verify content
	if !strings.Contains(merged, "Keep Unchanged") {
		t.Error("Should contain unchanged requirement")
	}
	if !strings.Contains(merged, "Updated content") {
		t.Error("Should contain updated requirement")
	}
	if strings.Contains(merged, "Remove This") {
		t.Error("Should not contain removed requirement")
	}
	if !strings.Contains(merged, "Renamed Feature") {
		t.Error("Should contain renamed requirement")
	}
	if !strings.Contains(merged, "Brand New") {
		t.Error("Should contain added requirement")
	}
}

func TestMergeSpec_PreservesOrder(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec with specific order
	baseContent := `# Test Spec

## Requirements

### Requirement: First
Content A.

#### Scenario: Test
- **WHEN** action
- **THEN** result

### Requirement: Second
Content B.

#### Scenario: Test
- **WHEN** action
- **THEN** result

### Requirement: Third
Content C.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Modify middle requirement
	deltaContent := `# Delta Spec

## MODIFIED Requirements

### Requirement: Second
Updated content B.

#### Scenario: Updated test
- **WHEN** updated action
- **THEN** updated result
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	merged, _, err := MergeSpec(basePath, deltaPath, true)
	if err != nil {
		t.Fatalf("MergeSpec failed: %v", err)
	}

	// Check order is preserved
	firstIdx := strings.Index(merged, "### Requirement: First")
	secondIdx := strings.Index(merged, "### Requirement: Second")
	thirdIdx := strings.Index(merged, "### Requirement: Third")

	if firstIdx < 0 || secondIdx < 0 || thirdIdx < 0 {
		t.Fatal("Missing requirements in merged spec")
	}

	if firstIdx >= secondIdx || secondIdx >= thirdIdx {
		t.Error("Requirement order was not preserved")
	}
}

func TestMergeSpec_ErrorOnNewSpecWithModified(t *testing.T) {
	tmpDir := t.TempDir()

	// Create delta spec with MODIFIED (not allowed for new specs)
	deltaContent := `# Delta Spec

## MODIFIED Requirements

### Requirement: Something
Content.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	basePath := filepath.Join(tmpDir, "base.md")
	_, _, err := MergeSpec(basePath, deltaPath, false)
	if err == nil {
		t.Error("Expected error for MODIFIED on new spec")
	}
	if !strings.Contains(err.Error(), "only ADDED requirements are allowed") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestMergeSpec_ErrorOnNoDeltas(t *testing.T) {
	tmpDir := t.TempDir()

	// Create delta spec without any delta sections
	deltaContent := `# Delta Spec

## Purpose
Just a regular spec.
`
	deltaPath := filepath.Join(tmpDir, "delta.md")
	if err := os.WriteFile(deltaPath, []byte(deltaContent), 0644); err != nil {
		t.Fatal(err)
	}

	basePath := filepath.Join(tmpDir, "base.md")
	_, _, err := MergeSpec(basePath, deltaPath, false)
	if err == nil {
		t.Error("Expected error for spec with no deltas")
	}
}

func TestFormatCapabilityName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"archive-workflow", "Archive Workflow"},
		{"cli-framework", "Cli Framework"},
		{"single", "Single"},
		{"multi-word-name", "Multi Word Name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatCapabilityName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
