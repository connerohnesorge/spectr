package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDeltaSpec_Added(t *testing.T) {
	content := `# Delta Spec

## ADDED Requirements

### Requirement: New Feature
The system SHALL support new functionality.

#### Scenario: Basic usage
- **WHEN** user performs action
- **THEN** feature responds

### Requirement: Another Feature
The system SHALL support another feature.

#### Scenario: Another case
- **WHEN** something happens
- **THEN** result occurs
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if len(plan.Added) != 2 {
		t.Errorf(
			"Expected 2 added requirements, got %d",
			len(plan.Added),
		)
	}

	if len(plan.Added) > 0 {
		if plan.Added[0].Name != "New Feature" {
			t.Errorf(
				"Expected first requirement name 'New Feature', got %q",
				plan.Added[0].Name,
			)
		}
	}

	if len(plan.Added) <= 1 {
		return
	}

	if plan.Added[1].Name != "Another Feature" {
		t.Errorf(
			"Expected second requirement name 'Another Feature', got %q",
			plan.Added[1].Name,
		)
	}
}

func TestParseDeltaSpec_Modified(t *testing.T) {
	content := `# Delta Spec

## MODIFIED Requirements

### Requirement: Updated Feature
The system SHALL have updated behavior.

#### Scenario: Modified scenario
- **WHEN** modified action
- **THEN** modified result
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if len(plan.Modified) != 1 {
		t.Errorf(
			"Expected 1 modified requirement, got %d",
			len(plan.Modified),
		)
	}

	if len(plan.Modified) == 0 {
		return
	}

	if plan.Modified[0].Name != "Updated Feature" {
		t.Errorf(
			"Expected requirement name 'Updated Feature', got %q",
			plan.Modified[0].Name,
		)
	}
}

func TestParseDeltaSpec_Removed(t *testing.T) {
	content := `# Delta Spec

## REMOVED Requirements

### Requirement: Deprecated Feature
**Reason**: No longer needed
**Migration**: Use new feature instead

### Requirement: Old Feature
**Reason**: Replaced by better implementation
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if len(plan.Removed) != 2 {
		t.Errorf(
			"Expected 2 removed requirements, got %d",
			len(plan.Removed),
		)
	}

	if len(plan.Removed) > 0 {
		if plan.Removed[0] != "Deprecated Feature" {
			t.Errorf(
				"Expected first removed requirement 'Deprecated Feature', got %q",
				plan.Removed[0],
			)
		}
	}

	if len(plan.Removed) <= 1 {
		return
	}

	if plan.Removed[1] != "Old Feature" {
		t.Errorf(
			"Expected second removed requirement 'Old Feature', got %q",
			plan.Removed[1],
		)
	}
}

func TestParseDeltaSpec_Renamed(t *testing.T) {
	content := `# Delta Spec

## RENAMED Requirements

- FROM: ` + "`" + `### Requirement: Old Name` + "`" + `
- TO: ` + "`" + `### Requirement: New Name` + "`" + `

- FROM: ` + "`" + `### Requirement: Another Old Name` + "`" + `
- TO: ` + "`" + `### Requirement: Another New Name` + "`" + `
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if len(plan.Renamed) != 2 {
		t.Errorf(
			"Expected 2 renamed requirements, got %d",
			len(plan.Renamed),
		)
	}

	if len(plan.Renamed) > 0 {
		if plan.Renamed[0].From != "Old Name" {
			t.Errorf(
				"Expected first rename from 'Old Name', got %q",
				plan.Renamed[0].From,
			)
		}
		if plan.Renamed[0].To != "New Name" {
			t.Errorf(
				"Expected first rename to 'New Name', got %q",
				plan.Renamed[0].To,
			)
		}
	}

	if len(plan.Renamed) <= 1 {
		return
	}

	if plan.Renamed[1].From != "Another Old Name" {
		t.Errorf(
			"Expected second rename from 'Another Old Name', got %q",
			plan.Renamed[1].From,
		)
	}
	if plan.Renamed[1].To != "Another New Name" {
		t.Errorf(
			"Expected second rename to 'Another New Name', got %q",
			plan.Renamed[1].To,
		)
	}
}

func TestParseDeltaSpec_AllOperations(
	t *testing.T,
) {
	content := `# Delta Spec

## ADDED Requirements

### Requirement: New Feature
Content here.

#### Scenario: Test
- **WHEN** something
- **THEN** result

## MODIFIED Requirements

### Requirement: Updated Feature
Modified content.

#### Scenario: Modified test
- **WHEN** action
- **THEN** new result

## REMOVED Requirements

### Requirement: Old Feature
**Reason**: Deprecated

## RENAMED Requirements

- FROM: ` + "`" + `### Requirement: Previous Name` + "`" + `
- TO: ` + "`" + `### Requirement: Current Name` + "`" + `
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if len(plan.Added) != 1 {
		t.Errorf(
			"Expected 1 added requirement, got %d",
			len(plan.Added),
		)
	}
	if len(plan.Modified) != 1 {
		t.Errorf(
			"Expected 1 modified requirement, got %d",
			len(plan.Modified),
		)
	}
	if len(plan.Removed) != 1 {
		t.Errorf(
			"Expected 1 removed requirement, got %d",
			len(plan.Removed),
		)
	}
	if len(plan.Renamed) != 1 {
		t.Errorf(
			"Expected 1 renamed requirement, got %d",
			len(plan.Renamed),
		)
	}

	if !plan.HasDeltas() {
		t.Error(
			"HasDeltas should return true when operations exist",
		)
	}

	expectedCount := 4
	if plan.CountOperations() != expectedCount {
		t.Errorf(
			"Expected %d total operations, got %d",
			expectedCount,
			plan.CountOperations(),
		)
	}
}

func TestParseDeltaSpec_NoDeltas(t *testing.T) {
	content := `# Delta Spec

## Purpose
This is just a regular spec without deltas.

## Requirements

### Requirement: Regular Requirement
This is not a delta operation.
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if plan.HasDeltas() {
		t.Error(
			"HasDeltas should return false when no operations exist",
		)
	}

	if plan.CountOperations() != 0 {
		t.Errorf(
			"Expected 0 operations, got %d",
			plan.CountOperations(),
		)
	}
}

func TestParseDeltaSpec_EmptySections(
	t *testing.T,
) {
	content := `# Delta Spec

## ADDED Requirements

## MODIFIED Requirements

## REMOVED Requirements

## RENAMED Requirements
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "spec.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	plan, err := ParseDeltaSpec(filePath)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if len(plan.Added) != 0 {
		t.Errorf(
			"Expected 0 added requirements, got %d",
			len(plan.Added),
		)
	}
	if len(plan.Modified) != 0 {
		t.Errorf(
			"Expected 0 modified requirements, got %d",
			len(plan.Modified),
		)
	}
	if len(plan.Removed) != 0 {
		t.Errorf(
			"Expected 0 removed requirements, got %d",
			len(plan.Removed),
		)
	}
	if len(plan.Renamed) != 0 {
		t.Errorf(
			"Expected 0 renamed requirements, got %d",
			len(plan.Renamed),
		)
	}
}

func TestParseDeltaSpec_MissingFile(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(
		tmpDir,
		"nonexistent.md",
	)

	_, err := ParseDeltaSpec(filePath)
	if err == nil {
		t.Error(
			"Expected error for missing file, got nil",
		)
	}
}
