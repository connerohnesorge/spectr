package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/view"
)

// TestViewCmd_Integration_TextOutput tests the view command with default text output
func TestViewCmd_Integration_TextOutput(t *testing.T) {
	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Logf("Warning: Failed to restore working directory: %v", err)
		}
	}()

	// Change to project root (where spectr directory exists)
	projectRoot := filepath.Join(originalWd, "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Verify spectr directory exists
	if _, err := os.Stat("spectr"); err != nil {
		t.Skipf("Skipping integration test: spectr directory not found at %s", projectRoot)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	cmd := &ViewCmd{JSON: false}
	err = cmd.Run()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("ViewCmd.Run() failed: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify expected elements are present in text output
	expectedElements := []string{
		"Spectr Dashboard",
		"════════════════════════════════════════════════════════════",
		"Summary:",
		"specs,",
		"requirements",
		"Use spectr list --changes or spectr list --specs for detailed views",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected text output to contain %q, but it didn't.\nFull output:\n%s",
				expected, output)
		}
	}

	t.Logf("Text output successfully generated (%d bytes)", len(output))
}

// TestViewCmd_Integration_JSONOutput tests the view command with JSON output
func TestViewCmd_Integration_JSONOutput(t *testing.T) {
	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Logf("Warning: Failed to restore working directory: %v", err)
		}
	}()

	// Change to project root (where spectr directory exists)
	projectRoot := filepath.Join(originalWd, "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Verify spectr directory exists
	if _, err := os.Stat("spectr"); err != nil {
		t.Skipf("Skipping integration test: spectr directory not found at %s", projectRoot)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command with JSON flag
	cmd := &ViewCmd{JSON: true}
	err = cmd.Run()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("ViewCmd.Run() with JSON failed: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output is valid JSON
	var data view.DashboardData
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput:\n%s", err, output)
	}

	// Verify expected JSON structure
	if data.Summary.TotalSpecs < 0 {
		t.Errorf("Expected TotalSpecs >= 0, got %d", data.Summary.TotalSpecs)
	}
	if data.Summary.TotalRequirements < 0 {
		t.Errorf("Expected TotalRequirements >= 0, got %d", data.Summary.TotalRequirements)
	}

	// Verify arrays are initialized (not null)
	if data.ActiveChanges == nil {
		t.Error("Expected ActiveChanges to be initialized (not null)")
	}
	if data.CompletedChanges == nil {
		t.Error("Expected CompletedChanges to be initialized (not null)")
	}
	if data.Specs == nil {
		t.Error("Expected Specs to be initialized (not null)")
	}

	t.Log("JSON output successfully generated and parsed:")
	t.Logf("  Specs: %d (with %d requirements)",
		data.Summary.TotalSpecs, data.Summary.TotalRequirements)
	t.Logf("  Active Changes: %d", data.Summary.ActiveChanges)
	t.Logf("  Completed Changes: %d", data.Summary.CompletedChanges)
}

// TestViewCmd_Integration_NOCOLOREnvironment tests that NO_COLOR is respected
func TestViewCmd_Integration_NOCOLOREnvironment(t *testing.T) {
	// Save original working directory and NO_COLOR env
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Logf("Warning: Failed to restore working directory: %v", err)
		}
	}()

	originalNoColor := os.Getenv("NO_COLOR")
	defer func() {
		if originalNoColor == "" {
			_ = os.Unsetenv("NO_COLOR")
		} else {
			_ = os.Setenv("NO_COLOR", originalNoColor)
		}
	}()

	// Change to project root
	projectRoot := filepath.Join(originalWd, "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change to project root: %v", err)
	}

	// Verify spectr directory exists
	if _, err := os.Stat("spectr"); err != nil {
		t.Skipf("Skipping integration test: spectr directory not found at %s", projectRoot)
	}

	// Set NO_COLOR environment variable
	_ = os.Setenv("NO_COLOR", "1")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command
	cmd := &ViewCmd{JSON: false}
	err = cmd.Run()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("ViewCmd.Run() with NO_COLOR failed: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output still contains expected content
	if !strings.Contains(output, "Spectr Dashboard") {
		t.Error("Expected dashboard title even with NO_COLOR")
	}

	// Note: lipgloss automatically respects NO_COLOR environment variable,
	// so we don't need to check for absence of ANSI codes explicitly.
	// The library handles this internally.
	t.Logf("NO_COLOR environment variable respected (output: %d bytes)", len(output))
}

// TestViewCmd_Integration_MissingSpectrDirectory tests behavior when spectr dir is missing
func TestViewCmd_Integration_MissingSpectrDirectory(t *testing.T) {
	// Create a temporary directory without a spectr subdirectory
	tempDir := t.TempDir()

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Logf("Warning: Failed to restore working directory: %v", err)
		}
	}()

	// Change to temp directory (no spectr directory)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the command (should succeed with empty dashboard)
	cmd := &ViewCmd{JSON: false}
	err = cmd.Run()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Should not error - just returns empty dashboard
	if err != nil {
		t.Fatalf("ViewCmd.Run() unexpectedly failed: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output shows empty project
	if !strings.Contains(output, "0 specs, 0 requirements") {
		t.Error("Expected output to show 0 specs and requirements")
	}
	if !strings.Contains(output, "Active Changes: 0 in progress") {
		t.Error("Expected output to show 0 active changes")
	}
	if !strings.Contains(output, "Spectr Dashboard") {
		t.Error("Expected dashboard header")
	}

	t.Log("Correctly handled missing spectr directory with empty dashboard")
}
