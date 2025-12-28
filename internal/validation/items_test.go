package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
)

const testProjectPath = "/test/project"

// TestCreateValidationItems_Changes tests creating validation items for changes
func TestCreateValidationItems_Changes(
	t *testing.T,
) {
	projectPath := testProjectPath
	basePath := "/test/project/spectr/changes"
	changeIDs := []string{
		"add-feature",
		"fix-bug",
		"update-docs",
	}

	items := CreateValidationItems(
		projectPath,
		changeIDs,
		ItemTypeChange,
		basePath,
	)

	assert.Equal(t, len(changeIDs), len(items))
	for i, id := range changeIDs {
		assert.Equal(t, id, items[i].Name)
		assert.Equal(
			t,
			ItemTypeChange,
			items[i].ItemType,
		)
		expectedPath := filepath.Join(
			basePath,
			id,
		)
		assert.Equal(
			t,
			expectedPath,
			items[i].Path,
		)
	}
}

// TestCreateValidationItems_Specs tests creating validation items for specs
func TestCreateValidationItems_Specs(
	t *testing.T,
) {
	projectPath := testProjectPath
	basePath := "/test/project/spectr/specs"
	specIDs := []string{
		"user-auth",
		"payment",
		"notifications",
	}

	items := CreateValidationItems(
		projectPath,
		specIDs,
		ItemTypeSpec,
		basePath,
	)

	assert.Equal(t, len(specIDs), len(items))
	for i, id := range specIDs {
		assert.Equal(t, id, items[i].Name)
		assert.Equal(
			t,
			ItemTypeSpec,
			items[i].ItemType,
		)
		expectedPath := filepath.Join(
			basePath,
			id,
			"spec.md",
		)
		assert.Equal(
			t,
			expectedPath,
			items[i].Path,
		)
	}
}

// TestCreateValidationItems_EmptyList tests creating items from empty list
func TestCreateValidationItems_EmptyList(
	t *testing.T,
) {
	projectPath := testProjectPath
	basePath := "/test/project/spectr/specs"

	items := CreateValidationItems(
		projectPath,
		make([]string, 0),
		ItemTypeSpec,
		basePath,
	)

	assert.Equal(t, 0, len(items))
}

// TestCreateValidationItems_PathConstruction tests correct path construction
func TestCreateValidationItems_PathConstruction(
	t *testing.T,
) {
	tests := []struct {
		name         string
		itemType     string
		id           string
		basePath     string
		expectedPath string
	}{
		{
			name:         "change path",
			itemType:     ItemTypeChange,
			id:           "add-feature",
			basePath:     "/project/spectr/changes",
			expectedPath: "/project/spectr/changes/add-feature",
		},
		{
			name:         "spec path",
			itemType:     ItemTypeSpec,
			id:           "user-auth",
			basePath:     "/project/spectr/specs",
			expectedPath: "/project/spectr/specs/user-auth/spec.md",
		},
		{
			name:         "change with complex id",
			itemType:     ItemTypeChange,
			id:           "add-multi-factor-auth",
			basePath:     "/project/spectr/changes",
			expectedPath: "/project/spectr/changes/add-multi-factor-auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := CreateValidationItems(
				"/project",
				[]string{tt.id},
				tt.itemType,
				tt.basePath,
			)
			assert.Equal(t, 1, len(items))
			assert.Equal(
				t,
				tt.expectedPath,
				items[0].Path,
			)
		})
	}
}

// TestGetAllItems tests getting all items (changes + specs)
func TestGetAllItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir,
		[]string{"add-feature", "fix-bug"},
		[]string{"user-auth", "payment"})

	items, err := GetAllItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 4, len(items))

	// Verify we got both changes and specs
	changeCount := 0
	specCount := 0
	for _, item := range items {
		switch item.ItemType {
		case ItemTypeChange:
			changeCount++
		case ItemTypeSpec:
			specCount++
		}
	}

	assert.Equal(t, 2, changeCount)
	assert.Equal(t, 2, specCount)
}

// TestGetAllItems_EmptyProject tests getting items from empty project
func TestGetAllItems_EmptyProject(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir, nil, nil)

	items, err := GetAllItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

// TestGetAllItems_OnlyChanges tests project with only changes
func TestGetAllItems_OnlyChanges(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{"add-feature", "fix-bug"},
		nil,
	)

	items, err := GetAllItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))
	for _, item := range items {
		assert.Equal(
			t,
			ItemTypeChange,
			item.ItemType,
		)
	}
}

// TestGetAllItems_OnlySpecs tests project with only specs
func TestGetAllItems_OnlySpecs(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		nil,
		[]string{"user-auth", "payment"},
	)

	items, err := GetAllItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))
	for _, item := range items {
		assert.Equal(
			t,
			ItemTypeSpec,
			item.ItemType,
		)
	}
}

// TestGetAllItems_DiscoveryError tests handling of discovery errors
func TestGetAllItems_DiscoveryError(
	t *testing.T,
) {
	// Use a path that will cause discovery to fail
	_, err := GetAllItems("/nonexistent/path", "spectr")

	// Should not error because discovery returns empty slice for nonexistent dirs
	// This is by design in the discovery package
	assert.NoError(t, err)
}

// TestGetChangeItems tests getting only change items
func TestGetChangeItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{
			"add-feature",
			"fix-bug",
			"update-docs",
		},
		[]string{"user-auth", "payment"},
	)

	items, err := GetChangeItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 3, len(items))

	expectedNames := map[string]bool{
		"add-feature": true,
		"fix-bug":     true,
		"update-docs": true,
	}

	for _, item := range items {
		assert.Equal(
			t,
			ItemTypeChange,
			item.ItemType,
		)
		assert.True(
			t,
			expectedNames[item.Name],
			"unexpected change: %s",
			item.Name,
		)
		expectedPath := filepath.Join(
			tmpDir,
			"spectr",
			"changes",
			item.Name,
		)
		assert.Equal(t, expectedPath, item.Path)
	}
}

// TestGetChangeItems_Empty tests getting changes from project without changes
func TestGetChangeItems_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		nil,
		[]string{"user-auth"},
	)

	items, err := GetChangeItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

// TestGetChangeItems_PathConstruction tests change item path construction
func TestGetChangeItems_PathConstruction(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{"add-feature"},
		nil,
	)

	items, err := GetChangeItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))

	expectedPath := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"add-feature",
	)
	assert.Equal(t, expectedPath, items[0].Path)
}

// TestGetSpecItems tests getting only spec items
func TestGetSpecItems(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{"add-feature"},
		[]string{
			"user-auth",
			"payment",
			"notifications",
		},
	)

	items, err := GetSpecItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 3, len(items))

	expectedNames := map[string]bool{
		"user-auth":     true,
		"payment":       true,
		"notifications": true,
	}

	for _, item := range items {
		assert.Equal(
			t,
			ItemTypeSpec,
			item.ItemType,
		)
		assert.True(
			t,
			expectedNames[item.Name],
			"unexpected spec: %s",
			item.Name,
		)
		expectedPath := filepath.Join(
			tmpDir,
			"spectr",
			"specs",
			item.Name,
			"spec.md",
		)
		assert.Equal(t, expectedPath, item.Path)
	}
}

// TestGetSpecItems_Empty tests getting specs from project without specs
func TestGetSpecItems_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{"add-feature"},
		nil,
	)

	items, err := GetSpecItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

// TestGetSpecItems_PathConstruction tests spec item path construction
func TestGetSpecItems_PathConstruction(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		nil,
		[]string{"user-auth"},
	)

	items, err := GetSpecItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))

	expectedPath := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
		"user-auth",
		"spec.md",
	)
	assert.Equal(t, expectedPath, items[0].Path)
}

// TestValidationItem_Structure tests the ValidationItem struct fields
func TestValidationItem_Structure(t *testing.T) {
	item := ValidationItem{
		Name:     "test-item",
		ItemType: ItemTypeChange,
		Path:     "/path/to/item",
	}

	assert.Equal(t, "test-item", item.Name)
	assert.Equal(t, ItemTypeChange, item.ItemType)
	assert.Equal(t, "/path/to/item", item.Path)
}

// TestGetChangeItems_DiscoveryFailure tests handling when discovery fails
func TestGetChangeItems_DiscoveryFailure(
	t *testing.T,
) {
	// Discovery should handle non-existent paths gracefully
	items, err := GetChangeItems(
		"/nonexistent/path",
		"spectr",
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

// TestGetSpecItems_DiscoveryFailure tests handling when discovery fails
func TestGetSpecItems_DiscoveryFailure(
	t *testing.T,
) {
	// Discovery should handle non-existent paths gracefully
	items, err := GetSpecItems(
		"/nonexistent/path",
		"spectr",
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(items))
}

// TestGetAllItems_OrderConsistency tests that items are returned in consistent order
func TestGetAllItems_OrderConsistency(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{
			"b-change",
			"a-change",
			"c-change",
		},
		[]string{"z-spec", "x-spec", "y-spec"},
	)

	// Get items multiple times
	items1, err1 := GetAllItems(tmpDir, "spectr")
	items2, err2 := GetAllItems(tmpDir, "spectr")

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, len(items1), len(items2))

	// Results should be in same order (changes are sorted, specs are sorted)
	for i := range items1 {
		assert.Equal(
			t,
			items1[i].Name,
			items2[i].Name,
		)
		assert.Equal(
			t,
			items1[i].ItemType,
			items2[i].ItemType,
		)
	}
}

// TestCreateValidationItems_PreservesOrder tests that item order matches input order
func TestCreateValidationItems_PreservesOrder(
	t *testing.T,
) {
	ids := []string{"third", "first", "second"}
	items := CreateValidationItems(
		"/project",
		ids,
		ItemTypeChange,
		"/base",
	)

	for i, id := range ids {
		assert.Equal(t, id, items[i].Name)
	}
}

// Test edge cases

// TestGetChangeItems_WithArchive tests that archived changes are excluded
func TestGetChangeItems_WithArchive(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{"active-change"},
		nil,
	)

	// Create archived change
	archiveDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"archive",
		"old-change",
	)
	err := os.MkdirAll(archiveDir, testDirPerm)
	assert.NoError(t, err)
	err = os.WriteFile(
		filepath.Join(archiveDir, "proposal.md"),
		[]byte("# Old"),
		testFilePerm,
	)
	assert.NoError(t, err)

	items, err := GetChangeItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(
		t,
		"active-change",
		items[0].Name,
	)
}

// TestGetChangeItems_WithHidden tests that hidden directories are excluded
func TestGetChangeItems_WithHidden(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		[]string{"visible-change"},
		nil,
	)

	// Create hidden directory
	hiddenDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		".hidden",
	)
	err := os.MkdirAll(hiddenDir, testDirPerm)
	assert.NoError(t, err)
	err = os.WriteFile(
		filepath.Join(hiddenDir, "proposal.md"),
		[]byte("# Hidden"),
		testFilePerm,
	)
	assert.NoError(t, err)

	items, err := GetChangeItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(
		t,
		"visible-change",
		items[0].Name,
	)
}

// TestGetSpecItems_WithHidden tests that hidden spec directories are excluded
func TestGetSpecItems_WithHidden(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(
		t,
		tmpDir,
		nil,
		[]string{"visible-spec"},
	)

	// Create hidden directory
	hiddenDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
		".hidden",
	)
	err := os.MkdirAll(hiddenDir, testDirPerm)
	assert.NoError(t, err)
	err = os.WriteFile(
		filepath.Join(hiddenDir, "spec.md"),
		[]byte("# Hidden"),
		testFilePerm,
	)
	assert.NoError(t, err)

	items, err := GetSpecItems(tmpDir, "spectr")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(items))
	assert.Equal(t, "visible-spec", items[0].Name)
}
