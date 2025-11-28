package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
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

// TestGetSpecsWithConfig_CustomRootDir tests discovery with custom root_dir
func TestGetSpecsWithConfig_CustomRootDir(t *testing.T) {
	tmpDir := t.TempDir()
	customRoot := "my-docs"
	specsDir := filepath.Join(tmpDir, customRoot, "specs")

	// Create test structure with custom root
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create test specs
	testSpecs := []string{"spec-a", "spec-b"}
	for _, name := range testSpecs {
		specDir := filepath.Join(specsDir, name)
		if err := os.MkdirAll(specDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		specPath := filepath.Join(specDir, "spec.md")
		if err := os.WriteFile(specPath, []byte("# Test Spec"), testFilePerm); err != nil {
			t.Fatal(err)
		}
	}

	// Create config with custom root
	cfg := &config.Config{
		RootDir:     customRoot,
		ProjectRoot: tmpDir,
	}

	// Test discovery with config
	specs, err := GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}

	if len(specs) != len(testSpecs) {
		t.Errorf("Expected %d specs, got %d", len(testSpecs), len(specs))
	}

	// Verify all specs found
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

// TestGetSpecsWithConfig_FromConfigFile tests discovery with spectr.yaml
func TestGetSpecsWithConfig_FromConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	customRoot := "documentation"

	// Create spectr.yaml with custom root_dir
	configContent := "root_dir: " + customRoot + "\n"
	configPath := filepath.Join(tmpDir, "spectr.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), testFilePerm); err != nil {
		t.Fatal(err)
	}

	// Create specs directory with custom root
	specsDir := filepath.Join(tmpDir, customRoot, "specs")
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create test spec
	specDir := filepath.Join(specsDir, "test-spec")
	if err := os.MkdirAll(specDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Test Spec"), testFilePerm); err != nil {
		t.Fatal(err)
	}

	// Load config and test discovery
	cfg, err := config.LoadFromPath(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	specs, err := GetSpecsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecsWithConfig failed: %v", err)
	}

	if len(specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(specs))
	}

	if len(specs) > 0 && specs[0] != "test-spec" {
		t.Errorf("Expected spec 'test-spec', got '%s'", specs[0])
	}
}

// TestGetSpecs_BackwardCompatibility ensures existing code still works
func TestGetSpecs_BackwardCompatibility(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")

	// Create standard structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create test spec
	specDir := filepath.Join(specsDir, "legacy-spec")
	if err := os.MkdirAll(specDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Legacy"), testFilePerm); err != nil {
		t.Fatal(err)
	}

	// Test using old API (should still work)
	specs, err := GetSpecs(tmpDir)
	if err != nil {
		t.Fatalf("GetSpecs failed: %v", err)
	}

	if len(specs) != 1 || specs[0] != "legacy-spec" {
		t.Errorf("Backward compatibility broken: expected ['legacy-spec'], got %v", specs)
	}
}

// TestGetSpecIDsWithConfig tests the config-based ID function
func TestGetSpecIDsWithConfig(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "spectr", "specs")

	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create test spec
	specDir := filepath.Join(specsDir, "id-test")
	if err := os.MkdirAll(specDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join(specDir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Test"), testFilePerm); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		RootDir:     "spectr",
		ProjectRoot: tmpDir,
	}

	ids, err := GetSpecIDsWithConfig(cfg)
	if err != nil {
		t.Fatalf("GetSpecIDsWithConfig failed: %v", err)
	}

	if len(ids) != 1 || ids[0] != "id-test" {
		t.Errorf("Expected ['id-test'], got %v", ids)
	}
}
