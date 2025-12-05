package pr

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRemoveChangeDirectory_Success tests successful removal of a change directory.
func TestRemoveChangeDirectory_Success(t *testing.T) {
	// Create temp directory structure:
	// tmpDir/spectr/changes/test-change/proposal.md
	tmpDir := t.TempDir()
	changeID := "test-change"
	changeDir := filepath.Join(tmpDir, "spectr", "changes", changeID)

	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}

	// Create a file inside the change directory
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte("# Test Proposal"), 0644); err != nil {
		t.Fatalf("Failed to create proposal.md: %v", err)
	}

	// Verify directory exists before removal
	if _, err := os.Stat(changeDir); os.IsNotExist(err) {
		t.Fatal("Change directory should exist before removal")
	}

	// Call RemoveChangeDirectory
	err := RemoveChangeDirectory(tmpDir, changeID)
	if err != nil {
		t.Fatalf("RemoveChangeDirectory() error = %v", err)
	}

	// Verify directory was removed
	if _, err := os.Stat(changeDir); !os.IsNotExist(err) {
		t.Error("Change directory should not exist after removal")
	}

	// Verify parent directories still exist
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if _, err := os.Stat(changesDir); os.IsNotExist(err) {
		t.Error("Parent changes directory should still exist")
	}
}

// TestRemoveChangeDirectory_EmptyChangeID tests that empty changeID returns an error.
func TestRemoveChangeDirectory_EmptyChangeID(t *testing.T) {
	tmpDir := t.TempDir()

	err := RemoveChangeDirectory(tmpDir, "")
	if err == nil {
		t.Error("RemoveChangeDirectory() expected error for empty changeID, got nil")
	}

	// Verify the error message mentions that changeID cannot be empty
	expectedMsg := "changeID cannot be empty"
	if err.Error() != expectedMsg {
		t.Errorf("RemoveChangeDirectory() error = %q, want %q", err.Error(), expectedMsg)
	}
}

// TestRemoveChangeDirectory_PathTraversal tests rejection of path traversal attempts.
func TestRemoveChangeDirectory_PathTraversal(t *testing.T) {
	tests := []struct {
		name     string
		changeID string
	}{
		{
			name:     "parent directory traversal",
			changeID: "../something",
		},
		{
			name:     "double parent traversal",
			changeID: "../../etc",
		},
		{
			name:     "traversal with normal path",
			changeID: "valid/../../../outside",
		},
		{
			name:     "hidden traversal in middle",
			changeID: "foo/../../../bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create the spectr/changes directory so paths can be resolved
			changesDir := filepath.Join(tmpDir, "spectr", "changes")
			if err := os.MkdirAll(changesDir, 0755); err != nil {
				t.Fatalf("Failed to create changes directory: %v", err)
			}

			err := RemoveChangeDirectory(tmpDir, tt.changeID)
			if err == nil {
				t.Errorf(
					"RemoveChangeDirectory() expected error for path traversal %q, got nil",
					tt.changeID,
				)
			}

			// Verify the error message mentions invalid path
			if err == nil {
				return
			}

			errStr := err.Error()
			if errStr == "" {
				t.Errorf("Expected non-empty error message for path traversal %q", tt.changeID)
			}
		})
	}
}

// TestRemoveChangeDirectory_NonexistentDirectory tests that removing a
// nonexistent directory succeeds (os.RemoveAll behavior).
func TestRemoveChangeDirectory_NonexistentDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the parent directories but NOT the change directory itself
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}

	// The change directory does not exist
	changeID := "nonexistent-change"
	changeDir := filepath.Join(changesDir, changeID)

	// Verify directory does not exist
	if _, err := os.Stat(changeDir); !os.IsNotExist(err) {
		t.Fatal("Change directory should not exist before test")
	}

	// Should succeed even though directory doesn't exist (os.RemoveAll behavior)
	err := RemoveChangeDirectory(tmpDir, changeID)
	if err != nil {
		t.Errorf("RemoveChangeDirectory() error = %v, expected nil for nonexistent directory", err)
	}
}

// TestRemoveChangeDirectory_NestedContents tests removal of directory with nested contents.
func TestRemoveChangeDirectory_NestedContents(t *testing.T) {
	tmpDir := t.TempDir()
	changeID := "nested-content-change"
	changeDir := filepath.Join(tmpDir, "spectr", "changes", changeID)

	// Create nested directory structure
	nestedDirs := []string{
		filepath.Join(changeDir, "specs"),
		filepath.Join(changeDir, "specs", "subdir"),
		filepath.Join(changeDir, "docs"),
	}

	for _, dir := range nestedDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create files at various levels
	files := map[string]string{
		"proposal.md":             "# Proposal",
		"tasks.md":                "# Tasks",
		"specs/spec1.yaml":        "spec: 1",
		"specs/subdir/spec2.yaml": "spec: 2",
		"docs/readme.md":          "# Readme",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(changeDir, relPath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", relPath, err)
		}
	}

	// Call RemoveChangeDirectory
	err := RemoveChangeDirectory(tmpDir, changeID)
	if err != nil {
		t.Fatalf("RemoveChangeDirectory() error = %v", err)
	}

	// Verify entire directory tree was removed
	if _, err := os.Stat(changeDir); !os.IsNotExist(err) {
		t.Error("Change directory should not exist after removal")
	}
}

// TestRemoveChangeDirectory_ValidChangeIDs tests various valid change ID formats.
func TestRemoveChangeDirectory_ValidChangeIDs(t *testing.T) {
	tests := []struct {
		name     string
		changeID string
	}{
		{
			name:     "simple name",
			changeID: "feature",
		},
		{
			name:     "hyphenated name",
			changeID: "add-new-feature",
		},
		{
			name:     "underscored name",
			changeID: "add_new_feature",
		},
		{
			name:     "with numbers",
			changeID: "feature-123",
		},
		{
			name:     "long name",
			changeID: "this-is-a-very-long-change-id-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			changeDir := filepath.Join(tmpDir, "spectr", "changes", tt.changeID)

			// Create the change directory
			if err := os.MkdirAll(changeDir, 0755); err != nil {
				t.Fatalf("Failed to create change directory: %v", err)
			}

			// Create a test file
			testFile := filepath.Join(changeDir, "test.txt")
			if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Call RemoveChangeDirectory
			err := RemoveChangeDirectory(tmpDir, tt.changeID)
			if err != nil {
				t.Errorf("RemoveChangeDirectory() error = %v for changeID %q", err, tt.changeID)
			}

			// Verify directory was removed
			if _, err := os.Stat(changeDir); !os.IsNotExist(err) {
				t.Errorf(
					"Change directory should not exist after removal for changeID %q",
					tt.changeID,
				)
			}
		})
	}
}
