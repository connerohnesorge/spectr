package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetSpecs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)

	// Create test structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create specs
	testSpecs := []string{
		"auth",
		"api",
		"database",
	}
	for _, name := range testSpecs {
		specDir := filepath.Join(specsDir, name)
		if err := os.MkdirAll(specDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		specPath := filepath.Join(
			specDir,
			"spec.md",
		)
		if err := os.WriteFile(
			specPath,
			[]byte("# Test Spec"),
			testFilePerm,
		); err != nil {
			t.Fatal(err)
		}
	}

	// Create hidden directory (should be excluded)
	hiddenDir := filepath.Join(
		specsDir,
		".hidden",
	)
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
	emptyDir := filepath.Join(
		specsDir,
		"incomplete",
	)
	if err := os.MkdirAll(emptyDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test discovery
	specs, err := GetSpecs(tmpDir)
	if err != nil {
		t.Fatalf("GetSpecs failed: %v", err)
	}

	if len(specs) != len(testSpecs) {
		t.Errorf(
			"Expected %d specs, got %d",
			len(testSpecs),
			len(specs),
		)
	}

	// Verify all expected specs are found
	specMap := make(map[string]bool)
	for _, s := range specs {
		specMap[s] = true
	}
	for _, expected := range testSpecs {
		if !specMap[expected] {
			t.Errorf(
				"Expected spec %s not found",
				expected,
			)
		}
	}
}

func TestGetSpecs_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	specs, err := GetSpecs(tmpDir)
	if err != nil {
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if len(specs) != 0 {
		t.Errorf(
			"Expected empty result, got %d specs",
			len(specs),
		)
	}
}

func TestGetSpecIDs(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)

	// Create test structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create specs in non-alphabetical order to test sorting
	testSpecs := []string{
		"zebra-spec",
		"auth",
		"database",
	}
	for _, name := range testSpecs {
		specDir := filepath.Join(specsDir, name)
		if err := os.MkdirAll(specDir, testDirPerm); err != nil {
			t.Fatal(err)
		}
		specPath := filepath.Join(
			specDir,
			"spec.md",
		)
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
		t.Errorf(
			"Expected %d specs, got %d",
			len(testSpecs),
			len(specs),
		)
	}

	// Verify sorting (should be: auth, database, zebra-spec)
	expectedSorted := []string{
		"auth",
		"database",
		"zebra-spec",
	}
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
		t.Fatalf(
			"Expected no error, got: %v",
			err,
		)
	}
	if len(specs) != 0 {
		t.Errorf(
			"Expected empty result, got %d specs",
			len(specs),
		)
	}
}

func TestGetActiveChangeIDs_MissingProposalMd(
	t *testing.T,
) {
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

	// Create a change directory WITHOUT proposal.md
	changeDir := filepath.Join(
		changesDir,
		"incomplete-change",
	)
	if err := os.MkdirAll(changeDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test that it's excluded
	changes, err := GetActiveChangeIDs(tmpDir)
	if err != nil {
		t.Fatalf(
			"GetActiveChangeIDs failed: %v",
			err,
		)
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
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)

	// Create test structure
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create a spec directory WITHOUT spec.md
	specDir := filepath.Join(
		specsDir,
		"incomplete-spec",
	)
	if err := os.MkdirAll(specDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Test that it's excluded
	specs, err := GetSpecIDs(tmpDir)
	if err != nil {
		t.Fatalf("GetSpecIDs failed: %v", err)
	}

	if len(specs) != 0 {
		t.Errorf(
			"Expected 0 specs (no spec.md), got %d",
			len(specs),
		)
	}
}
