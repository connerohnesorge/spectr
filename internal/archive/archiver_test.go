package archive

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testDirPerm  = 0o755
	testFilePerm = 0o644
)

// setupTestProject creates a minimal spectr project structure for testing
func setupTestProject(t *testing.T, tmpDir string, changes []string) {
	t.Helper()

	// Create spectr directory structure
	spectrDir := filepath.Join(tmpDir, "spectr")
	changesDir := filepath.Join(spectrDir, "changes")
	specsDir := filepath.Join(spectrDir, "specs")

	if err := os.MkdirAll(changesDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(specsDir, testDirPerm); err != nil {
		t.Fatal(err)
	}

	// Create project.md
	projectContent := "# Test Project\n"
	if err := os.WriteFile(filepath.Join(spectrDir, "project.md"), []byte(projectContent), testFilePerm); err != nil {
		t.Fatal(err)
	}

	// Create each change with minimal required files
	for _, changeName := range changes {
		changeDir := filepath.Join(changesDir, changeName)
		changeSpecsDir := filepath.Join(changeDir, "specs", "test-feature")

		if err := os.MkdirAll(changeSpecsDir, testDirPerm); err != nil {
			t.Fatal(err)
		}

		// Create proposal.md
		proposalContent := "# Change: " + changeName + "\n\n## Why\nTest change.\n\n## What Changes\n- Test\n\n## Impact\n- specs: test-feature\n"
		if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), testFilePerm); err != nil {
			t.Fatal(err)
		}

		// Create tasks.md with completed tasks
		tasksContent := "## Tasks\n- [x] Task 1\n"
		if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), testFilePerm); err != nil {
			t.Fatal(err)
		}

		// Create delta spec
		deltaSpec := `## ADDED Requirements

### Requirement: Test Feature
The system SHALL provide test functionality.

#### Scenario: Basic test
- **WHEN** test is run
- **THEN** it passes
`
		if err := os.WriteFile(filepath.Join(changeSpecsDir, "spec.md"), []byte(deltaSpec), testFilePerm); err != nil {
			t.Fatal(err)
		}
	}
}

func TestArchive_PartialIDPrefix(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup project with a change that has a long name
	setupTestProject(t, tmpDir, []string{"refactor-unified-interactive-tui"})

	// Capture stdout to verify the resolution message
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &ArchiveCmd{
		ChangeID: "refactor", // Partial prefix
		Yes:      true,       // Skip confirmations
	}

	err := Archive(cmd, tmpDir)

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	// Verify the resolution message was printed
	expectedMsg := "Resolved 'refactor' -> 'refactor-unified-interactive-tui'"
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("Expected output to contain '%s', got: %s", expectedMsg, output)
	}

	// Verify the change was archived
	archiveDir := filepath.Join(tmpDir, "spectr", "changes", "archive")
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		t.Fatalf("Failed to read archive dir: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 archived change, got %d", len(entries))
	}

	// Verify the archive name contains the original change ID
	if !strings.Contains(entries[0].Name(), "refactor-unified-interactive-tui") {
		t.Errorf("Archive name should contain original ID, got: %s", entries[0].Name())
	}
}

func TestArchive_PartialIDSubstring(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup project with a change
	setupTestProject(t, tmpDir, []string{"refactor-unified-interactive-tui"})

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &ArchiveCmd{
		ChangeID: "unified", // Substring match (not prefix)
		Yes:      true,
	}

	err := Archive(cmd, tmpDir)

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	expectedMsg := "Resolved 'unified' -> 'refactor-unified-interactive-tui'"
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("Expected output to contain '%s', got: %s", expectedMsg, output)
	}
}

func TestArchive_PartialIDAmbiguous(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup project with multiple changes that match "add"
	setupTestProject(t, tmpDir, []string{"add-feature", "add-hotkey"})

	cmd := &ArchiveCmd{
		ChangeID: "add",
		Yes:      true,
	}

	err := Archive(cmd, tmpDir)

	if err == nil {
		t.Fatal("Expected error for ambiguous partial ID")
	}

	if !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("Expected ambiguous error, got: %v", err)
	}
}

func TestArchive_PartialIDNoMatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup project with a change that doesn't match
	setupTestProject(t, tmpDir, []string{"add-feature"})

	cmd := &ArchiveCmd{
		ChangeID: "nonexistent",
		Yes:      true,
	}

	err := Archive(cmd, tmpDir)

	if err == nil {
		t.Fatal("Expected error for non-matching partial ID")
	}

	expectedMsg := "no change found matching 'nonexistent'"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected '%s' error, got: %v", expectedMsg, err)
	}
}

func TestArchive_ExactIDNoResolutionMessage(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup project
	setupTestProject(t, tmpDir, []string{"add-feature"})

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &ArchiveCmd{
		ChangeID: "add-feature", // Exact match
		Yes:      true,
	}

	err := Archive(cmd, tmpDir)

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	// Should NOT contain resolution message for exact match
	if strings.Contains(output, "Resolved") {
		t.Error("Should not show resolution message for exact match")
	}
}

func TestArchive_CaseInsensitiveMatch(t *testing.T) {
	tmpDir := t.TempDir()

	setupTestProject(t, tmpDir, []string{"refactor-unified-tui"})

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &ArchiveCmd{
		ChangeID: "REFACTOR", // Uppercase
		Yes:      true,
	}

	err := Archive(cmd, tmpDir)

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	expectedMsg := "Resolved 'REFACTOR' -> 'refactor-unified-tui'"
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("Expected case-insensitive match message, got: %s", output)
	}
}
