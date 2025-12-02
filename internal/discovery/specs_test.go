package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSpecs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")

	// Create test structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create specs
	testSpecs := []string{"auth", "api", "database"}
	for _, name := range testSpecs {
		specDir := filepath.Join(specsDir, name)
		if err := os.MkdirAll(specDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		specPath := filepath.Join(specDir, "spec.md")
		if err := os.WriteFile(
			specPath,
			[]byte("# Test Spec"),
			testFilePerm,
		); err != nil {
			t.Fatal(err)
		}
	}

	// Create hidden directory (should be excluded)
	hiddenDir := filepath.Join(specsDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(hiddenDir, "spec.md"),
		[]byte("# Hidden"),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}

	// Create directory without spec.md (should be excluded)
	emptyDir := filepath.Join(specsDir, "incomplete")
	if err := os.MkdirAll(emptyDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test discovery
	specs, err := GetSpecs(tmpDir)
	if err != nil {
		t.Fatalf("GetSpecs failed: %v", err)
	}

	if len(specs) != len(testSpecs) {
		t.Errorf("Expected %d specs, got %d", len(testSpecs), len(specs))
	}

	// Verify all expected specs are found
	specMap := make(map[string]bool)
	for _, s := range specs {
		specMap[s] = true
	}
	for _, expected := range testSpecs {
		if !specMap[expected] {
			t.Errorf("Expected spec %s not found", expected)
		}
	}
}

func TestGetSpecs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	specs, err := GetSpecs(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("Expected empty result, got %d specs", len(specs))
	}
}

func TestGetSpecIDs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")

	// Create test structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create specs in non-alphabetical order to test sorting
	testSpecs := []string{"zebra-spec", "auth", "database"}
	for _, name := range testSpecs {
		specDir := filepath.Join(specsDir, name)
		if err := os.MkdirAll(specDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		specPath := filepath.Join(specDir, "spec.md")
		if err := os.WriteFile(
			specPath,
			[]byte("# Test Spec"),
			testFilePerm,
		); err != nil {
			t.Fatal(err)
		}
	}

	// Test GetSpecIDs
	specs, err := GetSpecIDs(tmpDir)
	if err != nil {
		t.Fatalf("GetSpecIDs failed: %v", err)
	}

	if len(specs) != len(testSpecs) {
		t.Errorf("Expected %d specs, got %d", len(testSpecs), len(specs))
	}

	// Verify sorting (should be: auth, database, zebra-spec)
	expectedSorted := []string{"auth", "database", "zebra-spec"}
	for i, expected := range expectedSorted {
		if i >= len(specs) {
			t.Error("Not enough specs returned")

			break
		}
		if specs[i] != expected {
			t.Errorf(
				"Expected spec[%d] to be %s, got %s",
				i,
				expected,
				specs[i],
			)
		}
	}
}

func TestGetSpecIDs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	specs, err := GetSpecIDs(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf("Expected empty result, got %d specs", len(specs))
	}
}

func TestGetActiveChangeIDs_MissingProposalMd(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

	// Create test structure
	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create a change directory WITHOUT proposal.md
	changeDir := filepath.Join(changesDir, "incomplete-change")
	if err := os.MkdirAll(changeDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test that it's excluded
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf("GetActiveChangeIDs failed: %v", err)
	}

	if len(changes) != 0 {
		t.Errorf(
			"Expected 0 changes (no proposal.md), got %d",
			len(changes),
		)
	}
}

func TestGetSpecIDs_MissingSpecMd(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")

	// Create test structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create a spec directory WITHOUT spec.md
	specDir := filepath.Join(specsDir, "incomplete-spec")
	if err := os.MkdirAll(specDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test that it's excluded
	specs, err := GetSpecIDs(tmpDir)
	if err != nil {
		t.Fatalf("GetSpecIDs failed: %v", err)
	}

	if len(specs) != 0 {
		t.Errorf("Expected 0 specs (no spec.md), got %d", len(specs))
	}
}

func TestResolveSpecID_ExactMatch(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "auth-service", "# Test")
	createSpecDir(t, specsDir, "auth-service-extended", "# Test")

	result, err := ResolveSpecID("auth-service", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.SpecID != "auth-service" {
		t.Errorf("Expected 'auth-service', got '%s'", result.SpecID)
	}
	if result.PartialMatch {
		t.Error("Expected PartialMatch to be false for exact match")
	}
}

func TestResolveSpecID_UniquePrefixMatch(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "database-schema-design", "# Test")
	createSpecDir(t, specsDir, "auth-service", "# Test")

	result, err := ResolveSpecID("database", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.SpecID != "database-schema-design" {
		t.Errorf("Expected 'database-schema-design', got '%s'", result.SpecID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true for prefix match")
	}
}

func TestResolveSpecID_UniqueSubstringMatch(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "database-schema-design", "# Test")
	createSpecDir(t, specsDir, "auth-service", "# Test")

	result, err := ResolveSpecID("schema", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.SpecID != "database-schema-design" {
		t.Errorf("Expected 'database-schema-design', got '%s'", result.SpecID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true for substring match")
	}
}

func TestResolveSpecID_MultiplePrefixMatches(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "api-gateway", "# Test")
	createSpecDir(t, specsDir, "api-versioning", "# Test")

	_, err := ResolveSpecID("api", tmpDir)
	if err == nil {
		t.Fatal("Expected error for ambiguous prefix match")
	}

	expectedMsg := "ambiguous ID 'api' matches multiple specs: api-gateway, api-versioning"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveSpecID_MultipleSubstringMatches(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "add-user-service", "# Test")
	createSpecDir(t, specsDir, "update-user-profile", "# Test")

	_, err := ResolveSpecID("user", tmpDir)
	if err == nil {
		t.Fatal("Expected error for ambiguous substring match")
	}

	expectedMsg := "ambiguous ID 'user' matches multiple specs: add-user-service, update-user-profile"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveSpecID_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "auth-service", "# Test")

	_, err := ResolveSpecID("nonexistent", tmpDir)
	if err == nil {
		t.Fatal("Expected error for no match")
	}

	expectedMsg := "no spec found matching 'nonexistent'"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestResolveSpecID_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "database-schema-design", "# Test")

	// Test uppercase prefix
	result, err := ResolveSpecID("DATABASE", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.SpecID != "database-schema-design" {
		t.Errorf("Expected 'database-schema-design', got '%s'", result.SpecID)
	}

	// Test mixed case substring
	result, err = ResolveSpecID("Schema", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.SpecID != "database-schema-design" {
		t.Errorf("Expected 'database-schema-design', got '%s'", result.SpecID)
	}
}

func TestResolveSpecID_PrefixPreferredOverSubstring(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	createSpecDir(t, specsDir, "api-gateway", "# Test")
	createSpecDir(t, specsDir, "rest-api-design", "# Test")

	result, err := ResolveSpecID("api", tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result.SpecID != "api-gateway" {
		t.Errorf("Expected 'api-gateway' (prefix match), got '%s'", result.SpecID)
	}
	if !result.PartialMatch {
		t.Error("Expected PartialMatch to be true")
	}
}

func TestResolveSpecID_EmptySpecs(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveSpecID("anything", tmpDir)
	if err == nil {
		t.Fatal("Expected error for empty specs directory")
	}

	expectedMsg := "no spec found matching 'anything'"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
