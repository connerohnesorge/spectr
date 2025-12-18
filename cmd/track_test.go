package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// TestTrackCmd_Fields verifies struct has correct fields (ChangeID, NoInteractive).
func TestTrackCmd_Fields(t *testing.T) {
	cmd := &TrackCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check ChangeID field exists and is string
	changeIDField := val.FieldByName("ChangeID")
	if !changeIDField.IsValid() {
		t.Error(
			"TrackCmd does not have ChangeID field",
		)
	}
	if changeIDField.Kind() != reflect.String {
		t.Errorf(
			"ChangeID should be string, got %v",
			changeIDField.Kind(),
		)
	}

	// Check NoInteractive field exists and is bool
	noInteractiveField := val.FieldByName(
		"NoInteractive",
	)
	if !noInteractiveField.IsValid() {
		t.Error(
			"TrackCmd does not have NoInteractive field",
		)
	}
	if noInteractiveField.Kind() != reflect.Bool {
		t.Errorf(
			"NoInteractive should be bool, got %v",
			noInteractiveField.Kind(),
		)
	}
}

// TestTrackCmd_DefaultValues tests that default values are correct.
func TestTrackCmd_DefaultValues(t *testing.T) {
	cmd := &TrackCmd{}

	if cmd.ChangeID != "" {
		t.Errorf(
			"ChangeID should default to empty string, got %q",
			cmd.ChangeID,
		)
	}
	if cmd.NoInteractive {
		t.Errorf(
			"NoInteractive should default to false, got %v",
			cmd.NoInteractive,
		)
	}
}

// TestTrackCmd_SetFields tests setting all field values.
func TestTrackCmd_SetFields(t *testing.T) {
	cmd := &TrackCmd{
		ChangeID:      "test-change",
		NoInteractive: true,
	}

	if cmd.ChangeID != "test-change" {
		t.Errorf(
			"ChangeID = %q, want %q",
			cmd.ChangeID,
			"test-change",
		)
	}
	if !cmd.NoInteractive {
		t.Error("NoInteractive should be true")
	}
}

// TestTrackCmd_HasRunMethod verifies the Run method exists.
func TestTrackCmd_HasRunMethod(t *testing.T) {
	cmd := &TrackCmd{}
	val := reflect.ValueOf(cmd)

	// Check that Run method exists
	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal(
			"TrackCmd does not have Run method",
		)
	}

	// Check that Run returns error
	runType := runMethod.Type()
	if runType.NumOut() != 1 {
		t.Errorf(
			"Run method should return 1 value, got %d",
			runType.NumOut(),
		)
	}

	if runType.NumOut() > 0 &&
		runType.Out(0).Name() != "error" {
		t.Errorf(
			"Run method should return error, got %s",
			runType.Out(0).Name(),
		)
	}
}

// TestCLIHasTrackCommand verifies Track command is registered in CLI.
func TestCLIHasTrackCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	trackField := val.FieldByName("Track")

	if !trackField.IsValid() {
		t.Fatal(
			"CLI struct does not have Track field",
		)
	}

	// Check the type
	if trackField.Type().Name() != "TrackCmd" {
		t.Errorf(
			"Track field type: got %s, want TrackCmd",
			trackField.Type().Name(),
		)
	}
}

// TestTrackCmd_Run_NoInteractiveWithoutChangeID verifies error when
// --no-interactive is set without a change ID.
func TestTrackCmd_Run_NoInteractiveWithoutChangeID(
	t *testing.T,
) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"Failed to get working directory: %v",
			err,
		)
	}
	defer func() {
		cErr := os.Chdir(originalWd)
		if cErr != nil {
			t.Logf(
				"Warning: Failed to restore working directory: %v",
				cErr,
			)
		}
	}()

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf(
			"Failed to change to temp directory: %v",
			err,
		)
	}

	// Run with NoInteractive but no ChangeID
	cmd := &TrackCmd{
		ChangeID:      "",
		NoInteractive: true,
	}

	err = cmd.Run()

	// Should return an error
	if err == nil {
		t.Fatal(
			"Expected error when --no-interactive is set without change ID",
		)
	}

	// Verify the error message
	expectedMsg := "change ID required when --no-interactive is set"
	if err.Error() != expectedMsg {
		t.Errorf(
			"Error message = %q, want %q",
			err.Error(),
			expectedMsg,
		)
	}
}

// TestTrackCmd_Run_NonExistentChange verifies error for non-existent change ID.
func TestTrackCmd_Run_NonExistentChange(
	t *testing.T,
) {
	// Create a temporary directory with spectr structure but no changes
	tempDir := t.TempDir()

	// Create spectr/changes directory
	changesDir := filepath.Join(
		tempDir,
		"spectr",
		"changes",
	)
	err := os.MkdirAll(changesDir, 0755)
	if err != nil {
		t.Fatalf(
			"Failed to create changes directory: %v",
			err,
		)
	}

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"Failed to get working directory: %v",
			err,
		)
	}
	defer func() {
		cErr := os.Chdir(originalWd)
		if cErr != nil {
			t.Logf(
				"Warning: Failed to restore working directory: %v",
				cErr,
			)
		}
	}()

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf(
			"Failed to change to temp directory: %v",
			err,
		)
	}

	// Run with a non-existent change ID
	cmd := &TrackCmd{
		ChangeID:      "non-existent-change",
		NoInteractive: true,
	}

	err = cmd.Run()

	// Should return an error (either change not found or no tasks file)
	if err == nil {
		t.Fatal(
			"Expected error when change does not exist",
		)
	}

	// The error should be NoTasksFileError since the change directory doesn't exist
	var noTasksErr *specterrs.NoTasksFileError
	if !errors.As(err, &noTasksErr) {
		t.Logf(
			"Got error type %T: %v",
			err,
			err,
		)
	}
}

// TestTrackCmd_Run_NoTasksFile verifies NoTasksFileError when tasks.jsonc
// doesn't exist for a change.
func TestTrackCmd_Run_NoTasksFile(t *testing.T) {
	// Check if git is available (skip test in nix builds where git may not be present)
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip(
			"git binary not available, skipping test",
		)
	}

	// Create a temporary directory with spectr structure
	tempDir := t.TempDir()

	// Initialize git repository (required since track uses git.GetRepoRoot)
	// This test will be skipped if git is not available (checked above)
	gitInit := exec.Command("git", "init", "-q")
	gitInit.Dir = tempDir
	if output, err := gitInit.CombinedOutput(); err != nil {
		t.Fatalf(
			"Failed to initialize git repo: %v: %s",
			err,
			output,
		)
	}

	// Create the change directory without tasks.jsonc
	changeDir := filepath.Join(
		tempDir,
		"spectr",
		"changes",
		"test-change",
	)
	err := os.MkdirAll(changeDir, 0755)
	if err != nil {
		t.Fatalf(
			"Failed to create change directory: %v",
			err,
		)
	}

	// Create a minimal proposal.md to make it a valid change
	proposalPath := filepath.Join(
		changeDir,
		"proposal.md",
	)
	proposalContent := `# Test Change

## Overview
This is a test proposal.

## Tasks
- [ ] 1.1 Test task
`
	err = os.WriteFile(
		proposalPath,
		[]byte(proposalContent),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create proposal.md: %v",
			err,
		)
	}

	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"Failed to get working directory: %v",
			err,
		)
	}
	defer func() {
		cErr := os.Chdir(originalWd)
		if cErr != nil {
			t.Logf(
				"Warning: Failed to restore working directory: %v",
				cErr,
			)
		}
	}()

	// Change to temp directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf(
			"Failed to change to temp directory: %v",
			err,
		)
	}

	// Run with a change that exists but has no tasks.jsonc
	cmd := &TrackCmd{
		ChangeID:      "test-change",
		NoInteractive: true,
	}

	err = cmd.Run()

	// Should return NoTasksFileError
	if err == nil {
		t.Fatal(
			"Expected error when tasks.jsonc does not exist",
		)
	}

	var noTasksErr *specterrs.NoTasksFileError
	if !errors.As(err, &noTasksErr) {
		t.Errorf(
			"Expected NoTasksFileError, got %T: %v",
			err,
			err,
		)
	}

	// Verify the error contains the change ID
	if noTasksErr != nil &&
		noTasksErr.ChangeID != "test-change" {
		t.Errorf(
			"NoTasksFileError.ChangeID = %q, want %q",
			noTasksErr.ChangeID,
			"test-change",
		)
	}
}

// TestTrackCmd_FieldTags verifies the struct tags are correct.
func TestTrackCmd_FieldTags(t *testing.T) {
	trackCmdType := reflect.TypeOf(TrackCmd{})

	// Check ChangeID field tag
	changeIDField, found := trackCmdType.FieldByName(
		"ChangeID",
	)
	if !found {
		t.Fatal(
			"TrackCmd does not have ChangeID field",
		)
	}

	// Verify it's an optional arg - Kong uses `arg:""` for positional args
	tag := string(changeIDField.Tag)
	if !containsSubstring(tag, "arg:") {
		t.Error(
			"ChangeID field should have arg tag",
		)
	}

	// Verify it's marked as optional
	if !containsSubstring(tag, "optional:") {
		t.Error(
			"ChangeID field should have optional tag",
		)
	}

	// Check NoInteractive field tag
	noInteractiveField, found := trackCmdType.FieldByName(
		"NoInteractive",
	)
	if !found {
		t.Fatal(
			"TrackCmd does not have NoInteractive field",
		)
	}

	if nameTag := noInteractiveField.Tag.Get("name"); nameTag != "no-interactive" {
		t.Errorf(
			"NoInteractive name tag = %q, want %q",
			nameTag,
			"no-interactive",
		)
	}

	if helpTag := noInteractiveField.Tag.Get("help"); helpTag == "" {
		t.Error(
			"NoInteractive field should have help tag",
		)
	}
}

// containsSubstring checks if a string contains a substring.
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
