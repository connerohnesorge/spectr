package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestValidatePostMerge_Success(t *testing.T) {
	validSpec := `# Test Spec

## Requirements

### Requirement: Feature One
Content here.

#### Scenario: Test case
- **WHEN** action
- **THEN** result

### Requirement: Feature Two
More content.

#### Scenario: Another test
- **WHEN** another action
- **THEN** another result
`

	err := ValidatePostMerge(validSpec, "test-path")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidatePostMerge_DuplicateRequirements(t *testing.T) {
	invalidSpec := `# Test Spec

## Requirements

### Requirement: Duplicate Name
Content here.

#### Scenario: Test
- **WHEN** action
- **THEN** result

### Requirement: Duplicate Name
Different content.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`

	err := ValidatePostMerge(invalidSpec, "test-path")
	if err == nil {
		t.Error("Expected error for duplicate requirements")
	}
	if !strings.Contains(err.Error(), "duplicate requirement") {
		t.Errorf("Expected duplicate error, got: %v", err)
	}
}

func TestValidatePostMerge_MissingScenario(t *testing.T) {
	invalidSpec := `# Test Spec

## Requirements

### Requirement: No Scenarios
This requirement has no scenarios.
`

	err := ValidatePostMerge(invalidSpec, "test-path")
	if err == nil {
		t.Error("Expected error for missing scenarios")
	}
	if !strings.Contains(err.Error(), "no scenarios") {
		t.Errorf("Expected scenario error, got: %v", err)
	}
}

func TestValidatePreMerge_NewSpec_OnlyAddedAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "base.md")

	// Delta with MODIFIED (not allowed for new spec)
	deltaPlan := &parsers.DeltaPlan{
		Modified: []parsers.RequirementBlock{
			{Name: "Something", Raw: "### Requirement: Something\n"},
		},
	}

	err := ValidatePreMerge(basePath, deltaPlan, false)
	if err == nil {
		t.Error("Expected error for MODIFIED on new spec")
	}
	if !strings.Contains(err.Error(), "only ADDED requirements") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestValidatePreMerge_ModifiedRequirementExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec
	baseContent := `# Test Spec

## Requirements

### Requirement: Existing Feature
Content.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Valid: Modify existing requirement
	deltaPlan := &parsers.DeltaPlan{
		Modified: []parsers.RequirementBlock{
			{
				Name: "Existing Feature",
				Raw:  "### Requirement: Existing Feature\nUpdated content.\n",
			},
		},
	}

	err := ValidatePreMerge(basePath, deltaPlan, true)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidatePreMerge_ModifiedRequirementDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec
	baseContent := `# Test Spec

## Requirements

### Requirement: Existing Feature
Content.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Invalid: Try to modify non-existent requirement
	deltaPlan := &parsers.DeltaPlan{
		Modified: []parsers.RequirementBlock{
			{Name: "Nonexistent Feature", Raw: "### Requirement: Nonexistent Feature\n"},
		},
	}

	err := ValidatePreMerge(basePath, deltaPlan, true)
	if err == nil {
		t.Error("Expected error for non-existent MODIFIED requirement")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected 'does not exist' error, got: %v", err)
	}
}

func TestValidatePreMerge_AddedRequirementAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec
	baseContent := `# Test Spec

## Requirements

### Requirement: Existing Feature
Content.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Invalid: Try to add requirement that already exists
	deltaPlan := &parsers.DeltaPlan{
		Added: []parsers.RequirementBlock{
			{Name: "Existing Feature", Raw: "### Requirement: Existing Feature\n"},
		},
	}

	err := ValidatePreMerge(basePath, deltaPlan, true)
	if err == nil {
		t.Error("Expected error for ADDED requirement that already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

func TestValidatePreMerge_RenamedRequirements(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base spec
	baseContent := `# Test Spec

## Requirements

### Requirement: Old Name
Content.

#### Scenario: Test
- **WHEN** action
- **THEN** result
`
	basePath := filepath.Join(tmpDir, "base.md")
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Valid: Rename existing requirement
	deltaPlan := &parsers.DeltaPlan{
		Renamed: []parsers.RenameOp{
			{From: "Old Name", To: "New Name"},
		},
	}

	err := ValidatePreMerge(basePath, deltaPlan, true)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCheckDuplicatesAndConflicts_NoDuplicates(t *testing.T) {
	deltaPlan := &parsers.DeltaPlan{
		Added: []parsers.RequirementBlock{
			{Name: "Feature One", Raw: "content"},
			{Name: "Feature Two", Raw: "content"},
		},
		Modified: []parsers.RequirementBlock{
			{Name: "Feature Three", Raw: "content"},
		},
	}

	err := CheckDuplicatesAndConflicts(deltaPlan)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCheckDuplicatesAndConflicts_DuplicateInAdded(t *testing.T) {
	deltaPlan := &parsers.DeltaPlan{
		Added: []parsers.RequirementBlock{
			{Name: "Duplicate", Raw: "content"},
			{Name: "Duplicate", Raw: "content"},
		},
	}

	err := CheckDuplicatesAndConflicts(deltaPlan)
	if err == nil {
		t.Error("Expected error for duplicate in ADDED section")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("Expected duplicate error, got: %v", err)
	}
}

func TestCheckDuplicatesAndConflicts_CrossSectionConflict(t *testing.T) {
	deltaPlan := &parsers.DeltaPlan{
		Added: []parsers.RequirementBlock{
			{Name: "Conflicting Feature", Raw: "content"},
		},
		Modified: []parsers.RequirementBlock{
			{Name: "Conflicting Feature", Raw: "content"},
		},
	}

	err := CheckDuplicatesAndConflicts(deltaPlan)
	if err == nil {
		t.Error("Expected error for cross-section conflict")
	}
	if !strings.Contains(err.Error(), "ADDED and MODIFIED") {
		t.Errorf("Expected cross-section conflict error, got: %v", err)
	}
}

func TestCheckDuplicatesAndConflicts_ModifiedAndRemoved(t *testing.T) {
	deltaPlan := &parsers.DeltaPlan{
		Modified: []parsers.RequirementBlock{
			{Name: "Conflicting Feature", Raw: "content"},
		},
		Removed: []string{"Conflicting Feature"},
	}

	err := CheckDuplicatesAndConflicts(deltaPlan)
	if err == nil {
		t.Error("Expected error for MODIFIED and REMOVED conflict")
	}
	if !strings.Contains(err.Error(), "MODIFIED and REMOVED") {
		t.Errorf("Expected cross-section conflict error, got: %v", err)
	}
}

func TestCheckDuplicatesAndConflicts_AddedAndRenamedTo(t *testing.T) {
	deltaPlan := &parsers.DeltaPlan{
		Added: []parsers.RequirementBlock{
			{Name: "New Name", Raw: "content"},
		},
		Renamed: []parsers.RenameOp{
			{From: "Old Name", To: "New Name"},
		},
	}

	err := CheckDuplicatesAndConflicts(deltaPlan)
	if err == nil {
		t.Error("Expected error for ADDED and RENAMED TO conflict")
	}
	if !strings.Contains(err.Error(), "ADDED and RENAMED TO") {
		t.Errorf("Expected cross-section conflict error, got: %v", err)
	}
}
