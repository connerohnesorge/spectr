package pr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRemoveChangeDirectory_Success tests successfully removing a valid change directory.
func TestRemoveChangeDirectory_Success(t *testing.T) {
	// Create a temporary directory to act as the project root
	projectRoot := t.TempDir()

	// Create the spectr/changes/test-change directory structure
	changeID := "test-change"
	changeDir := filepath.Join(projectRoot, "spectr", "changes", changeID)
	if err := os.MkdirAll(changeDir, 0755); err != nil {
		t.Fatalf("Failed to create test change directory: %v", err)
	}

	// Create some files inside the change directory
	testFile := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(testFile, []byte("# Test Proposal"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a subdirectory with files
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		t.Fatalf("Failed to create specs directory: %v", err)
	}
	specFile := filepath.Join(specsDir, "test.spec.md")
	if err := os.WriteFile(specFile, []byte("# Test Spec"), 0644); err != nil {
		t.Fatalf("Failed to create spec file: %v", err)
	}

	// Verify the directory exists before removal
	if _, err := os.Stat(changeDir); os.IsNotExist(err) {
		t.Fatal("Change directory should exist before removal")
	}

	// Call RemoveChangeDirectory
	err := RemoveChangeDirectory(projectRoot, changeID)
	if err != nil {
		t.Fatalf("RemoveChangeDirectory() error = %v, want nil", err)
	}

	// Verify the directory no longer exists
	if _, err := os.Stat(changeDir); !os.IsNotExist(err) {
		t.Error("RemoveChangeDirectory() directory should not exist after removal")
	}

	// Verify the parent directories still exist
	changesDir := filepath.Join(projectRoot, "spectr", "changes")
	if _, err := os.Stat(changesDir); os.IsNotExist(err) {
		t.Error("Parent changes directory should still exist")
	}
}

// TestRemoveChangeDirectory_EmptyChangeID tests that an empty changeID returns an error.
func TestRemoveChangeDirectory_EmptyChangeID(t *testing.T) {
	projectRoot := t.TempDir()

	err := RemoveChangeDirectory(projectRoot, "")
	if err == nil {
		t.Fatal("RemoveChangeDirectory() expected error for empty changeID, got nil")
	}

	if !strings.Contains(err.Error(), "changeID cannot be empty") {
		t.Errorf(
			"RemoveChangeDirectory() error should mention 'changeID cannot be empty', got: %v",
			err,
		)
	}
}

// TestRemoveChangeDirectory_NonExistentDirectory tests that a non-existent directory returns an error.
func TestRemoveChangeDirectory_NonExistentDirectory(t *testing.T) {
	projectRoot := t.TempDir()

	// Create the spectr/changes directory but NOT the specific change directory
	changesDir := filepath.Join(projectRoot, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}

	err := RemoveChangeDirectory(projectRoot, "non-existent-change")
	if err == nil {
		t.Fatal("RemoveChangeDirectory() expected error for non-existent directory, got nil")
	}

	if !strings.Contains(err.Error(), "change directory does not exist") {
		t.Errorf(
			"RemoveChangeDirectory() error should mention 'change directory does not exist', got: %v",
			err,
		)
	}
}

// TestRemoveChangeDirectory_PathNotDirectory tests that a file path (not directory) returns an error.
func TestRemoveChangeDirectory_PathNotDirectory(t *testing.T) {
	projectRoot := t.TempDir()

	// Create the spectr/changes directory
	changesDir := filepath.Join(projectRoot, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}

	// Create a FILE instead of a directory at the change path
	changeID := "file-not-dir"
	changePath := filepath.Join(changesDir, changeID)
	if err := os.WriteFile(changePath, []byte("I am a file, not a directory"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err := RemoveChangeDirectory(projectRoot, changeID)
	if err == nil {
		t.Fatal("RemoveChangeDirectory() expected error for file path, got nil")
	}

	if !strings.Contains(err.Error(), "path is not a directory") {
		t.Errorf(
			"RemoveChangeDirectory() error should mention 'path is not a directory', got: %v",
			err,
		)
	}
}

// TestRemoveChangeDirectory_PathTraversalPrevention tests that path traversal attempts return an error.
func TestRemoveChangeDirectory_PathTraversalPrevention(t *testing.T) {
	projectRoot := t.TempDir()

	// Create the spectr/changes directory
	changesDir := filepath.Join(projectRoot, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}

	// Create a directory that a path traversal attack might try to access
	targetDir := filepath.Join(projectRoot, "sensitive-data")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("Failed to create target directory: %v", err)
	}
	sensitiveFile := filepath.Join(targetDir, "secret.txt")
	if err := os.WriteFile(sensitiveFile, []byte("secret data"), 0644); err != nil {
		t.Fatalf("Failed to create sensitive file: %v", err)
	}

	tests := []struct {
		name     string
		changeID string
	}{
		{
			name:     "parent directory traversal",
			changeID: "../sensitive-data",
		},
		{
			name:     "double parent traversal",
			changeID: "../../sensitive-data",
		},
		{
			name:     "traversal with subdirectory",
			changeID: "../changes/../sensitive-data",
		},
		{
			name:     "absolute path injection",
			changeID: "/etc/passwd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RemoveChangeDirectory(projectRoot, tt.changeID)
			if err == nil {
				t.Errorf(
					"RemoveChangeDirectory() expected error for path traversal %q, got nil",
					tt.changeID,
				)

				return
			}

			// The error should either indicate path traversal prevention or directory not existing
			// Both are acceptable as they prevent the attack
			validError := strings.Contains(err.Error(), "invalid change directory path") ||
				strings.Contains(err.Error(), "change directory does not exist") ||
				strings.Contains(err.Error(), "is not within")

			if !validError {
				t.Errorf(
					"RemoveChangeDirectory() error should prevent path traversal for %q, got: %v",
					tt.changeID,
					err,
				)
			}
		})
	}

	// Verify the sensitive directory was NOT removed
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		t.Error("Path traversal prevention failed: sensitive directory was removed")
	}

	// Verify the sensitive file still exists
	if _, err := os.Stat(sensitiveFile); os.IsNotExist(err) {
		t.Error("Path traversal prevention failed: sensitive file was removed")
	}
}

// TestRemoveChangeDirectory_NestedPathTraversal tests more complex path traversal scenarios.
func TestRemoveChangeDirectory_NestedPathTraversal(t *testing.T) {
	projectRoot := t.TempDir()

	// Create a more complex directory structure
	changesDir := filepath.Join(projectRoot, "spectr", "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatalf("Failed to create changes directory: %v", err)
	}

	// Create a legitimate change directory
	legitimateChange := filepath.Join(changesDir, "legitimate-change")
	if err := os.MkdirAll(legitimateChange, 0755); err != nil {
		t.Fatalf("Failed to create legitimate change directory: %v", err)
	}

	// Create a sibling directory that should not be accessible
	siblingDir := filepath.Join(projectRoot, "spectr", "archives")
	if err := os.MkdirAll(siblingDir, 0755); err != nil {
		t.Fatalf("Failed to create sibling directory: %v", err)
	}

	tests := []struct {
		name     string
		changeID string
	}{
		{
			name:     "sibling directory access",
			changeID: "../archives",
		},
		{
			name:     "current directory",
			changeID: ".",
		},
		{
			name:     "parent from nested",
			changeID: "legitimate-change/../../archives",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RemoveChangeDirectory(projectRoot, tt.changeID)
			// Should either error or not remove the sibling
			if err != nil {
				return
			}

			// If no error, verify sibling still exists
			if _, statErr := os.Stat(siblingDir); os.IsNotExist(statErr) {
				t.Errorf(
					"RemoveChangeDirectory(%q) should not remove sibling directory",
					tt.changeID,
				)
			}
		})
	}

	// Verify sibling directory still exists
	if _, err := os.Stat(siblingDir); os.IsNotExist(err) {
		t.Error("Sibling directory should still exist after all tests")
	}
}
