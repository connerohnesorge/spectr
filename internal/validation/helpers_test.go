package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/connerohnesorge/spectr/internal/config"
)

const (
	testDirPerm  = 0755
	testFilePerm = 0644
)

// TestDetermineItemType_ChangeOnly tests when item exists only as a change
func TestDetermineItemType_ChangeOnly(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, nil)

	info, err := DetermineItemType(tmpDir, "add-feature", nil)
	assert.NoError(t, err)
	assert.True(t, info.IsChange)
	assert.False(t, info.IsSpec)
	assert.Equal(t, ItemTypeChange, info.ItemType)
}

// TestDetermineItemType_SpecOnly tests when item exists only as a spec
func TestDetermineItemType_SpecOnly(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"user-auth"})

	info, err := DetermineItemType(tmpDir, "user-auth", nil)
	assert.NoError(t, err)
	assert.False(t, info.IsChange)
	assert.True(t, info.IsSpec)
	assert.Equal(t, ItemTypeSpec, info.ItemType)
}

// TestDetermineItemType_BothExists tests when item exists as both change and spec
func TestDetermineItemType_BothExists(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"user-auth"}, []string{"user-auth"})

	// Without type flag should error
	_, err := DetermineItemType(tmpDir, "user-auth", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exists as both change and spec")
	assert.Contains(t, err.Error(), "use --type flag")
}

// TestDetermineItemType_NotFound tests when item doesn't exist
func TestDetermineItemType_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, []string{"user-auth"})

	_, err := DetermineItemType(tmpDir, "nonexistent", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestDetermineItemType_ExplicitTypeChange tests explicit --type=change flag
func TestDetermineItemType_ExplicitTypeChange(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, []string{"user-auth"})

	typeFlag := ItemTypeChange
	info, err := DetermineItemType(tmpDir, "add-feature", &typeFlag)
	assert.NoError(t, err)
	assert.Equal(t, ItemTypeChange, info.ItemType)
}

// TestDetermineItemType_ExplicitTypeSpec tests explicit --type=spec flag
func TestDetermineItemType_ExplicitTypeSpec(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, []string{"user-auth"})

	typeFlag := ItemTypeSpec
	info, err := DetermineItemType(tmpDir, "user-auth", &typeFlag)
	assert.NoError(t, err)
	assert.Equal(t, ItemTypeSpec, info.ItemType)
}

// TestDetermineItemType_ExplicitTypeOverridesBoth tests that explicit type flag
// disambiguates when item exists as both
func TestDetermineItemType_ExplicitTypeOverridesBoth(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"user-auth"}, []string{"user-auth"})

	// With type=change flag should succeed
	typeFlag := ItemTypeChange
	info, err := DetermineItemType(tmpDir, "user-auth", &typeFlag)
	assert.NoError(t, err)
	assert.Equal(t, ItemTypeChange, info.ItemType)

	// With type=spec flag should succeed
	typeFlag = ItemTypeSpec
	info, err = DetermineItemType(tmpDir, "user-auth", &typeFlag)
	assert.NoError(t, err)
	assert.Equal(t, ItemTypeSpec, info.ItemType)
}

// TestDetermineItemType_ExplicitTypeNotFound tests explicit type flag when item not found
func TestDetermineItemType_ExplicitTypeNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, []string{"user-auth"})

	// Request spec that doesn't exist
	typeFlag := ItemTypeSpec
	_, err := DetermineItemType(tmpDir, "nonexistent", &typeFlag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spec 'nonexistent' not found")

	// Request change that doesn't exist
	typeFlag = ItemTypeChange
	_, err = DetermineItemType(tmpDir, "nonexistent", &typeFlag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "change 'nonexistent' not found")
}

// TestDetermineItemType_ExplicitTypeWrongType tests explicit type flag when item exists as different type
func TestDetermineItemType_ExplicitTypeWrongType(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, []string{"user-auth"})

	// Item exists as change but request as spec
	typeFlag := ItemTypeSpec
	_, err := DetermineItemType(tmpDir, "add-feature", &typeFlag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spec 'add-feature' not found")

	// Item exists as spec but request as change
	typeFlag = ItemTypeChange
	_, err = DetermineItemType(tmpDir, "user-auth", &typeFlag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "change 'user-auth' not found")
}

// TestDetermineItemType_DiscoveryFailure tests handling of discovery errors
func TestDetermineItemType_DiscoveryFailure(t *testing.T) {
	// Use non-existent directory to trigger discovery failure
	_, err := DetermineItemType("/nonexistent/path", "test", nil)
	assert.Error(t, err)
}

// TestValidateItemByType_Change tests validating a change by type
func TestValidateItemByType_Change(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, nil)
	createValidChange(t, tmpDir, "add-feature")

	validator := NewValidator(false)
	report, err := ValidateItemByType(validator, tmpDir, "add-feature", ItemTypeChange)

	assert.NoError(t, err)
	assert.NotZero(t, report)
	// The validity depends on the change content, just ensure we got a report
}

// TestValidateItemByType_Spec tests validating a spec by type
func TestValidateItemByType_Spec(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"user-auth"})
	createValidSpec(t, tmpDir, "user-auth")

	validator := NewValidator(false)
	report, err := ValidateItemByType(validator, tmpDir, "user-auth", ItemTypeSpec)

	assert.NoError(t, err)
	assert.NotZero(t, report)
	assert.True(t, report.Valid)
}

// TestValidateItemByType_InvalidChange tests validating an invalid change
func TestValidateItemByType_InvalidChange(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"bad-change"}, nil)
	createInvalidChange(t, tmpDir, "bad-change")

	validator := NewValidator(true) // strict mode
	report, err := ValidateItemByType(validator, tmpDir, "bad-change", ItemTypeChange)

	// Should get an error or invalid report
	if err == nil {
		assert.False(t, report.Valid)
	}
}

// TestValidateSingleItem_Change tests validating a single change item
func TestValidateSingleItem_Change(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"add-feature"}, nil)
	changePath := createValidChange(t, tmpDir, "add-feature")

	item := ValidationItem{
		Name:     "add-feature",
		ItemType: ItemTypeChange,
		Path:     changePath,
	}

	validator := NewValidator(false)
	result, err := ValidateSingleItem(validator, item)

	assert.NoError(t, err)
	assert.Equal(t, "add-feature", result.Name)
	assert.Equal(t, ItemTypeChange, result.Type)
	assert.NotZero(t, result.Report)
}

// TestValidateSingleItem_Spec tests validating a single spec item
func TestValidateSingleItem_Spec(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"user-auth"})
	specPath := createValidSpec(t, tmpDir, "user-auth")

	item := ValidationItem{
		Name:     "user-auth",
		ItemType: ItemTypeSpec,
		Path:     specPath,
	}

	validator := NewValidator(false)
	result, err := ValidateSingleItem(validator, item)

	assert.NoError(t, err)
	assert.Equal(t, "user-auth", result.Name)
	assert.Equal(t, ItemTypeSpec, result.Type)
	assert.True(t, result.Valid)
	assert.NotZero(t, result.Report)
}

// TestValidateSingleItem_ValidationError tests handling of validation errors
func TestValidateSingleItem_ValidationError(t *testing.T) {
	item := ValidationItem{
		Name:     "nonexistent",
		ItemType: ItemTypeSpec,
		Path:     "/nonexistent/path/spec.md",
	}

	validator := NewValidator(false)
	result, err := ValidateSingleItem(validator, item)

	assert.Error(t, err)
	assert.Equal(t, "nonexistent", result.Name)
	assert.Equal(t, ItemTypeSpec, result.Type)
	assert.False(t, result.Valid)
	assert.NotZero(t, result.Error)
}

// TestValidateSingleItem_InvalidContent tests validation of invalid content
func TestValidateSingleItem_InvalidContent(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"bad-spec"})

	// Create invalid spec
	specDir := filepath.Join(tmpDir, config.DefaultRootDir, "specs", "bad-spec")
	specPath := filepath.Join(specDir, "spec.md")
	err := os.WriteFile(specPath, []byte("# Bad Spec\nNo proper content"), testFilePerm)
	assert.NoError(t, err)

	item := ValidationItem{
		Name:     "bad-spec",
		ItemType: ItemTypeSpec,
		Path:     specPath,
	}

	validator := NewValidator(true) // strict mode
	result, err := ValidateSingleItem(validator, item)

	// Should succeed in running validation but report should be invalid
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.NotZero(t, result.Report)
}

// TestContainsString tests the containsString helper
func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{"found at start", slice, "apple", true},
		{"found in middle", slice, "banana", true},
		{"found at end", slice, "cherry", true},
		{"not found", slice, "grape", false},
		{"empty string", slice, "", false},
		{"empty slice", make([]string, 0), "apple", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsString(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWithCustomRootDir tests that validation respects custom root_dir from config
func TestWithCustomRootDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create custom root dir "specs-custom" instead of "spectr"
	customRootDir := "specs-custom"

	// Create config file
	configContent := "root_dir: " + customRootDir + "\n"
	configPath := filepath.Join(tmpDir, "spectr.yaml")
	err := os.WriteFile(configPath, []byte(configContent), testFilePerm)
	assert.NoError(t, err)

	// Create project structure with custom root
	changesDir := filepath.Join(tmpDir, customRootDir, "changes", "test-change")
	err = os.MkdirAll(changesDir, testDirPerm)
	assert.NoError(t, err)

	proposalPath := filepath.Join(changesDir, "proposal.md")
	err = os.WriteFile(proposalPath, []byte("# Test Proposal"), testFilePerm)
	assert.NoError(t, err)

	specsDir := filepath.Join(tmpDir, customRootDir, "specs", "test-spec")
	err = os.MkdirAll(specsDir, testDirPerm)
	assert.NoError(t, err)

	specPath := filepath.Join(specsDir, "spec.md")
	specContent := `# Test Specification

## Purpose
This specification describes test functionality.

## Requirements

### Requirement: Test Feature
The system SHALL provide test functionality.

#### Scenario: Feature works
- **WHEN** user activates feature
- **THEN** feature is enabled
`
	err = os.WriteFile(specPath, []byte(specContent), testFilePerm)
	assert.NoError(t, err)

	// Test DetermineItemType with custom root
	info, err := DetermineItemType(tmpDir, "test-change", nil)
	assert.NoError(t, err)
	assert.Equal(t, ItemTypeChange, info.ItemType)
	assert.True(t, info.IsChange)

	info, err = DetermineItemType(tmpDir, "test-spec", nil)
	assert.NoError(t, err)
	assert.Equal(t, ItemTypeSpec, info.ItemType)
	assert.True(t, info.IsSpec)

	// Test ValidateItemByType with custom root
	validator := NewValidator(false)
	report, err := ValidateItemByType(validator, tmpDir, "test-spec", ItemTypeSpec)
	assert.NoError(t, err)
	assert.NotZero(t, report)
	assert.True(t, report.Valid)
}

// Helper functions for test setup

// setupTestProject creates a basic spectr project structure
func setupTestProject(t *testing.T, tmpDir string, changes, specs []string) {
	t.Helper()

	// Create changes directory
	changesDir := filepath.Join(tmpDir, config.DefaultRootDir, "changes")
	err := os.MkdirAll(changesDir, testDirPerm)
	assert.NoError(t, err)

	// Create change directories with proposal.md
	for _, change := range changes {
		changeDir := filepath.Join(changesDir, change)
		err := os.MkdirAll(changeDir, testDirPerm)
		assert.NoError(t, err)

		proposalPath := filepath.Join(changeDir, "proposal.md")
		err = os.WriteFile(proposalPath, []byte("# Test Proposal"), testFilePerm)
		assert.NoError(t, err)
	}

	// Create specs directory
	specsDir := filepath.Join(tmpDir, config.DefaultRootDir, "specs")
	err = os.MkdirAll(specsDir, testDirPerm)
	assert.NoError(t, err)

	// Create spec directories with spec.md
	for _, spec := range specs {
		specDir := filepath.Join(specsDir, spec)
		err := os.MkdirAll(specDir, testDirPerm)
		assert.NoError(t, err)

		specPath := filepath.Join(specDir, "spec.md")
		err = os.WriteFile(specPath, []byte("# Test Spec"), testFilePerm)
		assert.NoError(t, err)
	}
}

// createValidChange creates a valid change for testing
func createValidChange(t *testing.T, tmpDir, changeName string) string {
	t.Helper()

	changeDir := filepath.Join(tmpDir, config.DefaultRootDir, "changes", changeName)

	// Create proposal.md
	proposalContent := `# Change: Add Feature

## Why
This adds a new feature for testing.

## What Changes
- Add new functionality

## Impact
- Affects specs: user-auth
`
	proposalPath := filepath.Join(changeDir, "proposal.md")
	err := os.WriteFile(proposalPath, []byte(proposalContent), testFilePerm)
	assert.NoError(t, err)

	// Create tasks.md
	tasksContent := `## 1. Implementation
- [ ] Task 1
- [ ] Task 2
`
	tasksPath := filepath.Join(changeDir, "tasks.md")
	err = os.WriteFile(tasksPath, []byte(tasksContent), testFilePerm)
	assert.NoError(t, err)

	// Create a valid delta spec
	specsDir := filepath.Join(changeDir, "specs", "test-spec")
	err = os.MkdirAll(specsDir, testDirPerm)
	assert.NoError(t, err)

	deltaContent := `# Test Specification

## ADDED Requirements

### Requirement: New Feature
The system SHALL provide the new feature.

#### Scenario: Feature works
- **WHEN** user activates feature
- **THEN** feature is enabled
`
	deltaPath := filepath.Join(specsDir, "spec.md")
	err = os.WriteFile(deltaPath, []byte(deltaContent), testFilePerm)
	assert.NoError(t, err)

	return changeDir
}

// createInvalidChange creates an invalid change for testing
func createInvalidChange(t *testing.T, tmpDir, changeName string) string {
	t.Helper()

	changeDir := filepath.Join(tmpDir, config.DefaultRootDir, "changes", changeName)

	// Create minimal proposal.md that might fail validation
	proposalPath := filepath.Join(changeDir, "proposal.md")
	err := os.WriteFile(proposalPath, []byte("# Bad Change"), testFilePerm)
	assert.NoError(t, err)

	return changeDir
}

// createValidSpec creates a valid spec for testing
func createValidSpec(t *testing.T, tmpDir, specName string) string {
	t.Helper()

	specDir := filepath.Join(tmpDir, config.DefaultRootDir, "specs", specName)
	err := os.MkdirAll(specDir, testDirPerm)
	assert.NoError(t, err)

	specContent := `# Test Specification

## Purpose
This specification describes the test functionality for validation purposes. It contains sufficient detail to meet the minimum purpose requirements.

## Requirements

### Requirement: User Authentication
The system SHALL provide user authentication functionality.

#### Scenario: Successful login
- **WHEN** user provides valid credentials
- **THEN** user is authenticated and session is created

#### Scenario: Failed login
- **WHEN** user provides invalid credentials
- **THEN** authentication fails and error message is displayed
`
	specPath := filepath.Join(specDir, "spec.md")
	err = os.WriteFile(specPath, []byte(specContent), testFilePerm)
	assert.NoError(t, err)

	return specPath
}
