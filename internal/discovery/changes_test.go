package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

const testChangeIDRefactorTUI = "refactor-unified-interactive-tui"

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestGetActiveChanges(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create active changes
	testChanges := []string{
		"add-feature",
		"fix-bug",
		"update-docs",
	}
	for _, name := range testChanges {
		changeDir := filepath.Join(
			changesDir,
			name,
		)
		if err := os.MkdirAll(changeDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		proposalPath := filepath.Join(
			changeDir,
			"proposal.md",
		)
		if err := os.WriteFile(
			proposalPath,
			[]byte("# Test"),
			testFilePerm,
		); err != nil {
			t.Fatal(err)
		}
	}

	// Create archive directory (should be excluded)
	archiveDir := filepath.Join(
		changesDir,
		"archive",
		"old-change",
	)
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(archiveDir, "proposal.md"),
		[]byte("# Old"),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}

	// Create hidden directory (should be excluded)
	hiddenDir := filepath.Join(
		changesDir,
		".hidden",
	)
	if err := os.MkdirAll(hiddenDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(hiddenDir, "proposal.md"),
		[]byte("# Hidden"),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}

	// Create directory without proposal.md (should be excluded)
	emptyDir := filepath.Join(
		changesDir,
		"incomplete",
	)
	if err := os.MkdirAll(emptyDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test discovery
	changes, err := GetActiveChanges(tmpDir)
	if err != nil {
		t.Fatalf(
			"GetActiveChanges failed: %v",
			err,
		)
	}

	if len(changes) != len(testChanges) {
		t.Errorf(
			"Expected %d changes, got %d",
			len(testChanges),
			len(changes),
		)
	}

	// Verify all expected changes are found
	changeMap := make(map[string]bool)
	for _, c := range changes {
		changeMap[c] = true
	}
	for _, expected := range testChanges {
		if !changeMap[expected] {
			t.Errorf(
				"Expected change %s not found",
				expected,
			)
		}
	}

	// Verify archived and hidden changes are not included
	if changeMap["old-change"] {
		t.Error(
			"Archived change should not be included",
		)
	}
	if changeMap[".hidden"] {
		t.Error(
			"Hidden directory should not be included",
		)
	}
	if changeMap["incomplete"] {
		t.Error(
			"Incomplete change should not be included",
		)
	}
}

func TestGetActiveChanges_EmptyDirectory(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changes, err := GetActiveChanges(tmpDir)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if len(changes) != 0 {
		t.Errorf(
			"Expected empty result, got %d changes",
			len(changes),
		)
	}
}

func TestGetActiveChangeIDs_EmptyDirectory(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if len(changes) != 0 {
		t.Errorf(
			"Expected empty result, got %d changes",
			len(changes),
		)
	}
}

func TestGetActiveChangeIDs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create changes in non-alphabetical order to test sorting
	testChanges := []string{
		"zebra-feature",
		"add-feature",
		"middle-feature",
	}
	for _, name := range testChanges {
		createChangeDir(
			t,
			changesDir,
			name,
			"# Test",
		)
	}

	// Create archive directory (should be excluded)
	archiveDir := filepath.Join(
		changesDir,
		"archive",
	)
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	createChangeDir(
		t,
		archiveDir,
		"archived-change",
		"# Archived",
	)

	// Test GetActiveChangeIDs
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf(
			"GetActiveChangeIDs failed: %v",
			err,
		)
	}

	if len(changes) != len(testChanges) {
		t.Errorf(
			"Expected %d changes, got %d",
			len(testChanges),
			len(changes),
		)
	}

	// Verify sorting
	// (should be: add-feature, middle-feature, zebra-feature)
	expectedSorted := []string{
		"add-feature",
		"middle-feature",
		"zebra-feature",
	}
	verifyOrdering(
		t,
		changes,
		expectedSorted,
		"change",
	)

	// Verify archive is excluded
	verifyChangesExcluded(
		t,
		changes,
		[]string{"archived-change"},
	)
}

func TestResolveChangeID_ExactMatch(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		"add-feature",
		"# Test",
	)
	createChangeDir(
		t,
		changesDir,
		"add-feature-extended",
		"# Test",
	)

	result, err := ResolveChangeID(
		"add-feature",
		tmpDir,
	)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if result.ChangeID != "add-feature" {
		t.Errorf(
			"Expected 'add-feature', got '%s'",
			result.ChangeID,
		)
	}
	if result.PartialMatch {
		t.Error(
			"Expected PartialMatch to be false for exact match",
		)
	}
}

func TestResolveChangeID_UniquePrefixMatch(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		testChangeIDRefactorTUI,
		"# Test",
	)
	createChangeDir(
		t,
		changesDir,
		"add-feature",
		"# Test",
	)

	result, err := ResolveChangeID(
		"refactor",
		tmpDir,
	)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if result.ChangeID != testChangeIDRefactorTUI {
		t.Errorf(
			"Expected '%s', got '%s'",
			testChangeIDRefactorTUI,
			result.ChangeID,
		)
	}
	if !result.PartialMatch {
		t.Error(
			"Expected PartialMatch to be true for prefix match",
		)
	}
}

func TestResolveChangeID_UniqueSubstringMatch(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		testChangeIDRefactorTUI,
		"# Test",
	)
	createChangeDir(
		t,
		changesDir,
		"add-feature",
		"# Test",
	)

	result, err := ResolveChangeID(
		"unified",
		tmpDir,
	)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if result.ChangeID != testChangeIDRefactorTUI {
		t.Errorf(
			"Expected '%s', got '%s'",
			testChangeIDRefactorTUI,
			result.ChangeID,
		)
	}
	if !result.PartialMatch {
		t.Error(
			"Expected PartialMatch to be true for substring match",
		)
	}
}

func TestResolveChangeID_MultiplePrefixMatches(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		"add-feature",
		"# Test",
	)
	createChangeDir(
		t,
		changesDir,
		"add-hotkey",
		"# Test",
	)

	_, err := ResolveChangeID("add", tmpDir)
	if err == nil {
		t.Fatal(
			"Expected error for ambiguous prefix match",
		)
	}

	expectedMsg := "ambiguous ID 'add' matches multiple changes: add-feature, add-hotkey"
	if err.Error() != expectedMsg {
		t.Errorf(
			"Expected error message '%s', got '%s'",
			expectedMsg,
			err.Error(),
		)
	}
}

func TestResolveChangeID_MultipleSubstringMatches(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		"add-search-hotkey",
		"# Test",
	)
	createChangeDir(
		t,
		changesDir,
		"update-search-ui",
		"# Test",
	)

	_, err := ResolveChangeID("search", tmpDir)
	if err == nil {
		t.Fatal(
			"Expected error for ambiguous substring match",
		)
	}

	expectedMsg := "ambiguous ID 'search' matches multiple changes: add-search-hotkey, update-search-ui"
	if err.Error() != expectedMsg {
		t.Errorf(
			"Expected error message '%s', got '%s'",
			expectedMsg,
			err.Error(),
		)
	}
}

func TestResolveChangeID_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		"add-feature",
		"# Test",
	)

	_, err := ResolveChangeID(
		"nonexistent",
		tmpDir,
	)
	if err == nil {
		t.Fatal("Expected error for no match")
	}

	expectedMsg := "no change found matching 'nonexistent'"
	if err.Error() != expectedMsg {
		t.Errorf(
			"Expected error message '%s', got '%s'",
			expectedMsg,
			err.Error(),
		)
	}
}

func TestResolveChangeID_CaseInsensitive(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		testChangeIDRefactorTUI,
		"# Test",
	)

	// Test uppercase prefix
	result, err := ResolveChangeID(
		"REFACTOR",
		tmpDir,
	)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if result.ChangeID != testChangeIDRefactorTUI {
		t.Errorf(
			"Expected '%s', got '%s'",
			testChangeIDRefactorTUI,
			result.ChangeID,
		)
	}

	// Test mixed case substring
	result, err = ResolveChangeID(
		"Unified",
		tmpDir,
	)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if result.ChangeID != testChangeIDRefactorTUI {
		t.Errorf(
			"Expected '%s', got '%s'",
			testChangeIDRefactorTUI,
			result.ChangeID,
		)
	}
}

func TestResolveChangeID_PrefixPreferredOverSubstring(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(
		t,
		changesDir,
		"add-feature",
		"# Test",
	)
	createChangeDir(
		t,
		changesDir,
		"update-add-button",
		"# Test",
	)

	result, err := ResolveChangeID("add", tmpDir)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if result.ChangeID != "add-feature" {
		t.Errorf(
			"Expected 'add-feature' (prefix match), got '%s'",
			result.ChangeID,
		)
	}
	if !result.PartialMatch {
		t.Error(
			"Expected PartialMatch to be true",
		)
	}
}

func TestResolveChangeID_EmptyChanges(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveChangeID("anything", tmpDir)
	if err == nil {
		t.Fatal(
			"Expected error for empty changes directory",
		)
	}

	expectedMsg := "no change found matching 'anything'"
	if err.Error() != expectedMsg {
		t.Errorf(
			"Expected error message '%s', got '%s'",
			expectedMsg,
			err.Error(),
		)
	}
}

func TestExtractChangeIDFromArchivePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "date prefix",
			input:    "2024-01-15-feat-auth",
			expected: "feat-auth",
		},
		{
			name:     "no date prefix",
			input:    "feat-auth",
			expected: "feat-auth",
		},
		{
			name:     "date prefix with longer ID",
			input:    "2024-12-31-add-feature-x",
			expected: "add-feature-x",
		},
		{
			name:     "short name",
			input:    "abc",
			expected: "abc",
		},
		{
			name:     "exactly 11 chars but not date",
			input:    "abcdefghijk",
			expected: "abcdefghijk",
		},
		{
			name:     "looks like date but wrong format",
			input:    "2024-1-15-feat-auth",
			expected: "2024-1-15-feat-auth",
		},
		{
			name:     "date prefix with single char ID",
			input:    "2024-01-15-x",
			expected: "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractChangeIDFromArchivePath(tt.input)
			if result != tt.expected {
				t.Errorf(
					"ExtractChangeIDFromArchivePath(%q) = %q, want %q",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestIsChangeArchived(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"archive",
	)

	// Create archive structure
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create archived changes
	archivedChanges := []string{
		"2024-01-15-feat-auth",
		"2024-02-20-add-feature",
		"standalone-change",
	}
	for _, name := range archivedChanges {
		dir := filepath.Join(archiveDir, name)
		if err := os.MkdirAll(dir, testDirPerm); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		changeID string
		expected bool
	}{
		{
			name:     "archived with date prefix",
			changeID: "feat-auth",
			expected: true,
		},
		{
			name:     "archived with date prefix 2",
			changeID: "add-feature",
			expected: true,
		},
		{
			name:     "archived without date prefix",
			changeID: "standalone-change",
			expected: true,
		},
		{
			name:     "not archived",
			changeID: "nonexistent",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IsChangeArchived(tt.changeID, tmpDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf(
					"IsChangeArchived(%q) = %v, want %v",
					tt.changeID,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestIsChangeArchived_NoArchiveDir(t *testing.T) {
	tmpDir := t.TempDir()

	// No archive directory exists
	result, err := IsChangeArchived("anything", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result {
		t.Error("expected false when archive directory doesn't exist")
	}
}

func TestGetArchivedChangeIDs(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"archive",
	)

	// Create archive structure
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create archived changes
	archivedChanges := []string{
		"2024-01-15-feat-auth",
		"2024-02-20-add-feature",
		"standalone-change",
	}
	for _, name := range archivedChanges {
		dir := filepath.Join(archiveDir, name)
		if err := os.MkdirAll(dir, testDirPerm); err != nil {
			t.Fatal(err)
		}
	}

	// Create a hidden directory (should be excluded)
	hiddenDir := filepath.Join(archiveDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create a file (should be excluded)
	if err := os.WriteFile(
		filepath.Join(archiveDir, "readme.txt"),
		[]byte("test"),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}

	result, err := GetArchivedChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []string{
		"add-feature",
		"feat-auth",
		"standalone-change",
	}

	if len(result) != len(expectedIDs) {
		t.Fatalf(
			"expected %d IDs, got %d: %v",
			len(expectedIDs),
			len(result),
			result,
		)
	}

	for i, id := range expectedIDs {
		if result[i] != id {
			t.Errorf(
				"expected ID %q at index %d, got %q",
				id,
				i,
				result[i],
			)
		}
	}
}

func TestGetArchivedChangeIDs_NoArchiveDir(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := GetArchivedChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestGetChangeStatus(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	archiveDir := filepath.Join(changesDir, "archive")

	// Create directories
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create active change
	createChangeDir(t, changesDir, "active-change", "# Active")

	// Create archived change
	archivedDir := filepath.Join(archiveDir, "2024-01-15-archived-change")
	if err := os.MkdirAll(archivedDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		changeID string
		expected ChangeStatus
	}{
		{
			name:     "active change",
			changeID: "active-change",
			expected: ChangeStatusActive,
		},
		{
			name:     "archived change",
			changeID: "archived-change",
			expected: ChangeStatusArchived,
		},
		{
			name:     "unknown change",
			changeID: "nonexistent",
			expected: ChangeStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := GetChangeStatus(tt.changeID, tmpDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if status != tt.expected {
				t.Errorf(
					"GetChangeStatus(%q) = %v, want %v",
					tt.changeID,
					status,
					tt.expected,
				)
			}
		})
	}
}

func TestChangeStatus_String(t *testing.T) {
	tests := []struct {
		status   ChangeStatus
		expected string
	}{
		{ChangeStatusActive, "active"},
		{ChangeStatusArchived, "archived"},
		{ChangeStatusUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.status.String() != tt.expected {
				t.Errorf(
					"ChangeStatus.String() = %q, want %q",
					tt.status.String(),
					tt.expected,
				)
			}
		})
	}
}
