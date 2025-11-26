package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/charmbracelet/bubbles/table"
)

// TestRunInteractiveValidation_NotTTY tests error when not in a TTY
func TestRunInteractiveValidation_NotTTY(t *testing.T) {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe (not a TTY)
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	defer func() { _ = r.Close() }()
	defer func() { _ = w.Close() }()

	// Replace stdout with pipe
	os.Stdout = w

	// Restore stdout after test
	defer func() {
		os.Stdout = oldStdout
	}()

	// This should return an error because pipe is not a TTY
	err = RunInteractiveValidation("/test", false, false)

	// Restore stdout before assertions
	os.Stdout = oldStdout

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interactive mode requires a TTY")
}

// TestValidateItems tests the validateItems helper function
func TestValidateItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"test-spec"})
	createValidSpec(t, tmpDir, "test-spec")

	specPath := filepath.Join(tmpDir, SpectrDir, "specs", "test-spec", "spec.md")
	items := []ValidationItem{
		{
			Name:     "test-spec",
			ItemType: ItemTypeSpec,
			Path:     specPath,
		},
	}

	validator := NewValidator(false)
	results, hasFailures := validateItems(validator, items)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "test-spec", results[0].Name)
	assert.False(t, hasFailures)
	assert.True(t, results[0].Valid)
}

// TestValidateItems_MultipleItems tests validating multiple items
func TestValidateItems_MultipleItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"spec1", "spec2"})
	createValidSpec(t, tmpDir, "spec1")
	createValidSpec(t, tmpDir, "spec2")

	items := []ValidationItem{
		{
			Name:     "spec1",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "spec1", "spec.md"),
		},
		{
			Name:     "spec2",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "spec2", "spec.md"),
		},
	}

	validator := NewValidator(false)
	results, hasFailures := validateItems(validator, items)

	assert.Equal(t, 2, len(results))
	assert.False(t, hasFailures)
	assert.True(t, results[0].Valid)
	assert.True(t, results[1].Valid)
}

// TestValidateItems_WithFailures tests validation with some failures
func TestValidateItems_WithFailures(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"good-spec", "bad-spec"})
	createValidSpec(t, tmpDir, "good-spec")

	// Create invalid spec
	badSpecDir := filepath.Join(tmpDir, SpectrDir, "specs", "bad-spec")
	err := os.MkdirAll(badSpecDir, testDirPerm)
	assert.NoError(t, err)
	badSpecPath := filepath.Join(badSpecDir, "spec.md")
	err = os.WriteFile(badSpecPath, []byte("# Bad\nInvalid content"), testFilePerm)
	assert.NoError(t, err)

	items := []ValidationItem{
		{
			Name:     "good-spec",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "good-spec", "spec.md"),
		},
		{
			Name:     "bad-spec",
			ItemType: ItemTypeSpec,
			Path:     badSpecPath,
		},
	}

	validator := NewValidator(true) // strict mode
	results, hasFailures := validateItems(validator, items)

	assert.Equal(t, 2, len(results))
	assert.True(t, hasFailures)
	assert.True(t, results[0].Valid)  // good-spec should be valid
	assert.False(t, results[1].Valid) // bad-spec should be invalid
}

// TestValidateItems_EmptyList tests validation with empty item list
func TestValidateItems_EmptyList(t *testing.T) {
	validator := NewValidator(false)
	results, hasFailures := validateItems(validator, make([]ValidationItem, 0))

	assert.Equal(t, 0, len(results))
	assert.False(t, hasFailures)
}

// TestTruncateString tests the truncateString helper
func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "shorter than max",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "equal to max",
			input:    "exact",
			maxLen:   5,
			expected: "exact",
		},
		{
			name:     "longer than max",
			input:    "this is a very long string",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "maxLen less than ellipsis length",
			input:    "test",
			maxLen:   2,
			expected: "te",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   5,
			expected: "",
		},
		{
			name:     "path truncation",
			input:    "/very/long/path/to/some/file.md",
			maxLen:   20,
			expected: "/very/long/path/t...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
			assert.True(t, len(result) <= tt.maxLen)
		})
	}
}

// TestApplyTableStyles tests the applyTableStyles helper
func TestApplyTableStyles(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"test"}}
	tbl := table.New(table.WithColumns(columns), table.WithRows(rows))

	// This should not panic
	assert.NotPanics(t, func() {
		applyTableStyles(&tbl)
	})
}

// TestConstants tests that constants are defined correctly
func TestConstants(t *testing.T) {
	assert.Equal(t, 35, validationIDWidth)
	assert.Equal(t, 10, validationTypeWidth)
	assert.Equal(t, 55, validationPathWidth)
	assert.Equal(t, 53, validationPathTruncate)
	assert.Equal(t, 3, ellipsisMinLength)
	assert.Equal(t, 10, tableHeight)
}

// TestTruncateString_ExactBoundary tests truncation at exact boundary
func TestTruncateString_ExactBoundary(t *testing.T) {
	// String exactly at truncation point
	input := strings.Repeat("a", validationPathTruncate)
	result := truncateString(input, validationPathTruncate)
	assert.Equal(t, input, result)

	// String one character over
	input = strings.Repeat("a", validationPathTruncate+1)
	result = truncateString(input, validationPathTruncate)
	assert.Equal(t, validationPathTruncate, len(result))
	assert.True(t, strings.HasSuffix(result, "..."))
}

// TestMenuSelectionConstants tests menu selection indices
func TestMenuSelectionConstants(t *testing.T) {
	assert.Equal(t, 0, menuSelectionAll)
	assert.Equal(t, 1, menuSelectionChanges)
	assert.Equal(t, 2, menuSelectionSpecs)
	assert.Equal(t, 3, menuSelectionPickItem)
}

// TestHandleMenuSelection_All tests handling "All" selection
func TestHandleMenuSelection_All(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"test-spec"})
	createValidSpec(t, tmpDir, "test-spec")

	// This should not error
	err := handleMenuSelection(menuSelectionAll, tmpDir, false, false)
	assert.NoError(t, err)
}

// TestHandleMenuSelection_Changes tests handling "Changes" selection
func TestHandleMenuSelection_Changes(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, []string{"test-change"}, nil)

	// This should not error (even with empty changes)
	err := handleMenuSelection(menuSelectionChanges, tmpDir, false, false)
	assert.NoError(t, err)
}

// TestHandleMenuSelection_Specs tests handling "Specs" selection
func TestHandleMenuSelection_Specs(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"test-spec"})
	createValidSpec(t, tmpDir, "test-spec")

	// This should not error
	err := handleMenuSelection(menuSelectionSpecs, tmpDir, false, false)
	assert.NoError(t, err)
}

// TestHandleMenuSelection_Invalid tests handling invalid selection
func TestHandleMenuSelection_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, nil)

	// Invalid selection should return nil (no error)
	err := handleMenuSelection(999, tmpDir, false, false)
	assert.NoError(t, err)
}

// TestRunValidationAndPrint_Empty tests printing with no items
func TestRunValidationAndPrint_Empty(t *testing.T) {
	// Should not error with empty items
	err := runValidationAndPrint(make([]ValidationItem, 0), false, false)
	assert.NoError(t, err)
}

// TestRunValidationAndPrint_JSON tests JSON output
func TestRunValidationAndPrint_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, []string{"test-spec"})
	createValidSpec(t, tmpDir, "test-spec")

	items := []ValidationItem{
		{
			Name:     "test-spec",
			ItemType: ItemTypeSpec,
			Path:     filepath.Join(tmpDir, SpectrDir, "specs", "test-spec", "spec.md"),
		},
	}

	// Should not error
	err := runValidationAndPrint(items, false, true)
	assert.NoError(t, err)
}
