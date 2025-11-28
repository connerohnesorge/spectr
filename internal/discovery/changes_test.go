package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
)

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestGetActiveChanges(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create active changes
	testChanges := []string{"add-feature", "fix-bug", "update-docs"}
	for _, name := range testChanges {
		changeDir := filepath.Join(changesDir, name)
		if err := os.MkdirAll(changeDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		proposalPath := filepath.Join(changeDir, "proposal.md")
		if err := os.WriteFile(
			proposalPath,
			[]byte("# Test"),
			testFilePerm,
		); err != nil {
			t.Fatal(err)
		}
	}

	// Create archive directory (should be excluded)
	archiveDir := filepath.Join(changesDir, "archive", "old-change")
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
	hiddenDir := filepath.Join(changesDir, ".hidden")
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
	emptyDir := filepath.Join(changesDir, "incomplete")
	if err := os.MkdirAll(emptyDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test discovery
	changes, err := GetActiveChanges(tmpDir)
	if err != nil {
		t.Fatalf("GetActiveChanges failed: %v", err)
	}

	if len(changes) != len(testChanges) {
		t.Errorf("Expected %d changes, got %d", len(testChanges), len(changes))
	}

	// Verify all expected changes are found
	changeMap := make(map[string]bool)
	for _, c := range changes {
		changeMap[c] = true
	}
	for _, expected := range testChanges {
		if !changeMap[expected] {
			t.Errorf("Expected change %s not found", expected)
		}
	}

	// Verify archived and hidden changes are not included
	if changeMap["old-change"] {
		t.Error("Archived change should not be included")
	}
	if changeMap[".hidden"] {
		t.Error("Hidden directory should not be included")
	}
	if changeMap["incomplete"] {
		t.Error("Incomplete change should not be included")
	}
}

func TestGetActiveChangesWithConfig(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create active changes
	testChanges := []string{"add-feature", "fix-bug"}
	for _, name := range testChanges {
		createChangeDir(t, changesDir, name, "# Test")
	}

	// Create config
	cfg := &config.Config{
		RootDir:     "spectr",
		ProjectRoot: tmpDir,
	}

	// Test discovery
	changes, err := GetActiveChangesWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetActiveChangesWithConfig failed: %v", err)
	}

	if len(changes) != len(testChanges) {
		t.Errorf("Expected %d changes, got %d", len(testChanges), len(changes))
	}

	// Verify all expected changes are found
	changeMap := make(map[string]bool)
	for _, c := range changes {
		changeMap[c] = true
	}
	for _, expected := range testChanges {
		if !changeMap[expected] {
			t.Errorf("Expected change %s not found", expected)
		}
	}
}

func TestGetActiveChanges_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	changes, err := GetActiveChanges(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("Expected empty result, got %d changes", len(changes))
	}
}
func TestGetActiveChangeIDs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("Expected empty result, got %d changes", len(changes))
	}
}
func TestGetActiveChangeIDs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

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
		createChangeDir(t, changesDir, name, "# Test")
	}

	// Create archive directory (should be excluded)
	archiveDir := filepath.Join(changesDir, "archive")
	if err := os.MkdirAll(archiveDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	createChangeDir(t, archiveDir, "archived-change", "# Archived")

	// Test GetActiveChangeIDs
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("GetActiveChangeIDs failed: %v", err)
	}

	if len(changes) != len(testChanges) {
		t.Errorf("Expected %d changes, got %d", len(testChanges), len(changes))
	}

	// Verify sorting
	// (should be: add-feature, middle-feature, zebra-feature)
	expectedSorted := []string{
		"add-feature",
		"middle-feature",
		"zebra-feature",
	}
	verifyOrdering(t, changes, expectedSorted, "change")

	// Verify archive is excluded
	verifyChangesExcluded(t, changes, []string{"archived-change"})
}

func TestResolveChangeID_ExactMatch(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "add-feature", "# Test")
	createChangeDir(t, changesDir, "add-feature-extended", "# Test")

	result, err := ResolveChangeID("add-feature", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "add-feature" {
		t.Errorf("Expected 'add-feature', got '%s'", result.ChangeID)
	}
	if result.PartialMatch {
		t.Error("Expected PartialMatch to be false for exact match")
	}
}

func TestResolveChangeID_UniquePrefixMatch(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "refactor-unified-interactive-tui", "# Test")
	createChangeDir(t, changesDir, "add-feature", "# Test")

	result, err := ResolveChangeID("refactor", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "refactor-unified-interactive-tui" {
		t.Errorf("Expected 'refactor-unified-interactive-tui', got '%s'", result.ChangeID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true for prefix match")
	}
}

func TestResolveChangeID_UniqueSubstringMatch(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "refactor-unified-interactive-tui", "# Test")
	createChangeDir(t, changesDir, "add-feature", "# Test")

	result, err := ResolveChangeID("unified", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "refactor-unified-interactive-tui" {
		t.Errorf("Expected 'refactor-unified-interactive-tui', got '%s'", result.ChangeID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true for substring match")
	}
}

func TestResolveChangeID_MultiplePrefixMatches(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "add-feature", "# Test")
	createChangeDir(t, changesDir, "add-hotkey", "# Test")

	_, err := ResolveChangeID("add", tmpDir)
	if err == nil {
		t.Fatal("Expected error for ambiguous prefix match")
	}

	expectedMsg := "ambiguous ID 'add' matches multiple changes: add-feature, add-hotkey"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveChangeID_MultipleSubstringMatches(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "add-search-hotkey", "# Test")
	createChangeDir(t, changesDir, "update-search-ui", "# Test")

	_, err := ResolveChangeID("search", tmpDir)
	if err == nil {
		t.Fatal("Expected error for ambiguous substring match")
	}

	expectedMsg := "ambiguous ID 'search' matches multiple changes: add-search-hotkey, update-search-ui"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveChangeID_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "add-feature", "# Test")

	_, err := ResolveChangeID("nonexistent", tmpDir)
	if err == nil {
		t.Fatal("Expected error for no match")
	}

	expectedMsg := "no change found matching 'nonexistent'"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveChangeID_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "refactor-unified-interactive-tui", "# Test")

	// Test uppercase prefix
	result, err := ResolveChangeID("REFACTOR", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "refactor-unified-interactive-tui" {
		t.Errorf("Expected 'refactor-unified-interactive-tui', got '%s'", result.ChangeID)
	}

	// Test mixed case substring
	result, err = ResolveChangeID("Unified", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "refactor-unified-interactive-tui" {
		t.Errorf("Expected 'refactor-unified-interactive-tui', got '%s'", result.ChangeID)
	}
}

func TestResolveChangeID_PrefixPreferredOverSubstring(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "add-feature", "# Test")
	createChangeDir(t, changesDir, "update-add-button", "# Test")

	result, err := ResolveChangeID("add", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "add-feature" {
		t.Errorf("Expected 'add-feature' (prefix match), got '%s'", result.ChangeID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true")
	}
}

func TestResolveChangeID_EmptyChanges(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveChangeID("anything", tmpDir)
	if err == nil {
		t.Fatal("Expected error for empty changes directory")
	}

	expectedMsg := "no change found matching 'anything'"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveChangeIDWithConfig_CustomRootDir(t *testing.T) {
	tmpDir := t.TempDir()
	customRoot := "my-specs"
	changesDir := filepath.Join(tmpDir, customRoot, "changes")
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createChangeDir(t, changesDir, "add-custom-feature", "# Test")
	createChangeDir(t, changesDir, "update-custom-ui", "# Test")

	cfg := &config.Config{
		RootDir:     customRoot,
		ProjectRoot: tmpDir,
	}

	// Test exact match
	result, err := ResolveChangeIDWithConfig("add-custom-feature", cfg)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "add-custom-feature" {
		t.Errorf("Expected 'add-custom-feature', got '%s'", result.ChangeID)
	}
	if result.PartialMatch {
		t.Error("Expected PartialMatch to be false for exact match")
	}

	// Test prefix match
	result, err = ResolveChangeIDWithConfig("update", cfg)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.ChangeID != "update-custom-ui" {
		t.Errorf("Expected 'update-custom-ui', got '%s'", result.ChangeID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true for prefix match")
	}
}
