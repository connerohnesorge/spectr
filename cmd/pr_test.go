package cmd

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/git"
	"github.com/connerohnesorge/spectr/internal/pr"
)

// TestPRArchiveCmd_Struct tests the PRArchiveCmd struct construction.
func TestPRArchiveCmd_Struct(t *testing.T) {
	cmd := &PRArchiveCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check ChangeID field exists and is string
	changeIDField := val.FieldByName("ChangeID")
	if !changeIDField.IsValid() {
		t.Error("PRArchiveCmd does not have ChangeID field")
	}
	if changeIDField.Kind() != reflect.String {
		t.Errorf("ChangeID should be string, got %v", changeIDField.Kind())
	}

	// Check Base field exists
	baseField := val.FieldByName("Base")
	if !baseField.IsValid() {
		t.Error("PRArchiveCmd does not have Base field")
	}
	if baseField.Kind() != reflect.String {
		t.Errorf("Base should be string, got %v", baseField.Kind())
	}

	// Check Draft field exists
	draftField := val.FieldByName("Draft")
	if !draftField.IsValid() {
		t.Error("PRArchiveCmd does not have Draft field")
	}
	if draftField.Kind() != reflect.Bool {
		t.Errorf("Draft should be bool, got %v", draftField.Kind())
	}

	// Check Force field exists
	forceField := val.FieldByName("Force")
	if !forceField.IsValid() {
		t.Error("PRArchiveCmd does not have Force field")
	}
	if forceField.Kind() != reflect.Bool {
		t.Errorf("Force should be bool, got %v", forceField.Kind())
	}

	// Check DryRun field exists
	dryRunField := val.FieldByName("DryRun")
	if !dryRunField.IsValid() {
		t.Error("PRArchiveCmd does not have DryRun field")
	}
	if dryRunField.Kind() != reflect.Bool {
		t.Errorf("DryRun should be bool, got %v", dryRunField.Kind())
	}

	// Check SkipSpecs field exists (archive-specific)
	skipSpecsField := val.FieldByName("SkipSpecs")
	if !skipSpecsField.IsValid() {
		t.Error("PRArchiveCmd does not have SkipSpecs field")
	}
	if skipSpecsField.Kind() != reflect.Bool {
		t.Errorf("SkipSpecs should be bool, got %v", skipSpecsField.Kind())
	}
}

// TestPRArchiveCmd_DefaultValues tests that default values are correct.
func TestPRArchiveCmd_DefaultValues(t *testing.T) {
	cmd := &PRArchiveCmd{}

	// All fields should have zero values by default
	if cmd.ChangeID != "" {
		t.Errorf("ChangeID should default to empty string, got %q", cmd.ChangeID)
	}
	if cmd.Base != "" {
		t.Errorf("Base should default to empty string, got %q", cmd.Base)
	}
	if cmd.Draft {
		t.Errorf("Draft should default to false, got %v", cmd.Draft)
	}
	if cmd.Force {
		t.Errorf("Force should default to false, got %v", cmd.Force)
	}
	if cmd.DryRun {
		t.Errorf("DryRun should default to false, got %v", cmd.DryRun)
	}
	if cmd.SkipSpecs {
		t.Errorf("SkipSpecs should default to false, got %v", cmd.SkipSpecs)
	}
}

// TestPRArchiveCmd_SetFields tests setting all field values.
func TestPRArchiveCmd_SetFields(t *testing.T) {
	cmd := &PRArchiveCmd{
		ChangeID:  "test-change",
		Base:      "main",
		Draft:     true,
		Force:     true,
		DryRun:    true,
		SkipSpecs: true,
	}

	if cmd.ChangeID != "test-change" {
		t.Errorf("ChangeID = %q, want %q", cmd.ChangeID, "test-change")
	}
	if cmd.Base != "main" {
		t.Errorf("Base = %q, want %q", cmd.Base, "main")
	}
	if !cmd.Draft {
		t.Error("Draft should be true")
	}
	if !cmd.Force {
		t.Error("Force should be true")
	}
	if !cmd.DryRun {
		t.Error("DryRun should be true")
	}
	if !cmd.SkipSpecs {
		t.Error("SkipSpecs should be true")
	}
}

// TestPRNewCmd_Struct tests the PRNewCmd struct construction.
func TestPRNewCmd_Struct(t *testing.T) {
	cmd := &PRNewCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check ChangeID field exists
	changeIDField := val.FieldByName("ChangeID")
	if !changeIDField.IsValid() {
		t.Error("PRNewCmd does not have ChangeID field")
	}
	if changeIDField.Kind() != reflect.String {
		t.Errorf("ChangeID should be string, got %v", changeIDField.Kind())
	}

	// Check Base field exists
	baseField := val.FieldByName("Base")
	if !baseField.IsValid() {
		t.Error("PRNewCmd does not have Base field")
	}
	if baseField.Kind() != reflect.String {
		t.Errorf("Base should be string, got %v", baseField.Kind())
	}

	// Check Draft field exists
	draftField := val.FieldByName("Draft")
	if !draftField.IsValid() {
		t.Error("PRNewCmd does not have Draft field")
	}
	if draftField.Kind() != reflect.Bool {
		t.Errorf("Draft should be bool, got %v", draftField.Kind())
	}

	// Check Force field exists
	forceField := val.FieldByName("Force")
	if !forceField.IsValid() {
		t.Error("PRNewCmd does not have Force field")
	}
	if forceField.Kind() != reflect.Bool {
		t.Errorf("Force should be bool, got %v", forceField.Kind())
	}

	// Check DryRun field exists
	dryRunField := val.FieldByName("DryRun")
	if !dryRunField.IsValid() {
		t.Error("PRNewCmd does not have DryRun field")
	}
	if dryRunField.Kind() != reflect.Bool {
		t.Errorf("DryRun should be bool, got %v", dryRunField.Kind())
	}

	// Verify SkipSpecs is NOT in PRNewCmd (it's archive-specific)
	skipSpecsField := val.FieldByName("SkipSpecs")
	if skipSpecsField.IsValid() {
		t.Error("PRNewCmd should not have SkipSpecs field (archive-specific)")
	}
}

// TestPRNewCmd_DefaultValues tests that default values are correct.
func TestPRNewCmd_DefaultValues(t *testing.T) {
	cmd := &PRNewCmd{}

	// All fields should have zero values by default
	if cmd.ChangeID != "" {
		t.Errorf("ChangeID should default to empty string, got %q", cmd.ChangeID)
	}
	if cmd.Base != "" {
		t.Errorf("Base should default to empty string, got %q", cmd.Base)
	}
	if cmd.Draft {
		t.Errorf("Draft should default to false, got %v", cmd.Draft)
	}
	if cmd.Force {
		t.Errorf("Force should default to false, got %v", cmd.Force)
	}
	if cmd.DryRun {
		t.Errorf("DryRun should default to false, got %v", cmd.DryRun)
	}
}

// TestPRNewCmd_SetFields tests setting all field values.
func TestPRNewCmd_SetFields(t *testing.T) {
	cmd := &PRNewCmd{
		ChangeID: "new-proposal",
		Base:     "develop",
		Draft:    true,
		Force:    true,
		DryRun:   true,
	}

	if cmd.ChangeID != "new-proposal" {
		t.Errorf("ChangeID = %q, want %q", cmd.ChangeID, "new-proposal")
	}
	if cmd.Base != "develop" {
		t.Errorf("Base = %q, want %q", cmd.Base, "develop")
	}
	if !cmd.Draft {
		t.Error("Draft should be true")
	}
	if !cmd.Force {
		t.Error("Force should be true")
	}
	if !cmd.DryRun {
		t.Error("DryRun should be true")
	}
}

// TestPRCmd_Struct tests the parent PRCmd struct.
func TestPRCmd_Struct(t *testing.T) {
	cmd := &PRCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Archive subcommand exists
	archiveField := val.FieldByName("Archive")
	if !archiveField.IsValid() {
		t.Error("PRCmd does not have Archive field")
	}
	if archiveField.Type().Name() != "PRArchiveCmd" {
		t.Errorf("Archive field type = %s, want PRArchiveCmd", archiveField.Type().Name())
	}

	// Check New subcommand exists
	newField := val.FieldByName("New")
	if !newField.IsValid() {
		t.Error("PRCmd does not have New field")
	}
	if newField.Type().Name() != "PRNewCmd" {
		t.Errorf("New field type = %s, want PRNewCmd", newField.Type().Name())
	}
}

// TestPRArchiveCmd_HasRunMethod verifies the Run method exists.
func TestPRArchiveCmd_HasRunMethod(t *testing.T) {
	cmd := &PRArchiveCmd{}
	val := reflect.ValueOf(cmd)

	// Check that Run method exists
	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal("PRArchiveCmd does not have Run method")
	}

	// Check that Run returns error
	runType := runMethod.Type()
	if runType.NumOut() != 1 {
		t.Errorf("Run method should return 1 value, got %d", runType.NumOut())
	}

	if runType.NumOut() > 0 && runType.Out(0).Name() != "error" {
		t.Errorf("Run method should return error, got %s", runType.Out(0).Name())
	}
}

// TestPRNewCmd_HasRunMethod verifies the Run method exists.
func TestPRNewCmd_HasRunMethod(t *testing.T) {
	cmd := &PRNewCmd{}
	val := reflect.ValueOf(cmd)

	// Check that Run method exists
	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal("PRNewCmd does not have Run method")
	}

	// Check that Run returns error
	runType := runMethod.Type()
	if runType.NumOut() != 1 {
		t.Errorf("Run method should return 1 value, got %d", runType.NumOut())
	}

	if runType.NumOut() > 0 && runType.Out(0).Name() != "error" {
		t.Errorf("Run method should return error, got %s", runType.Out(0).Name())
	}
}

// TestPrintPRResult_Basic tests that printPRResult does not panic with basic input.
func TestPrintPRResult_Basic(t *testing.T) {
	// Capture stdout to prevent test output clutter
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		_ = w.Close()
		os.Stdout = oldStdout
		// Drain the pipe
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)

		// Check for panic recovery
		if rec := recover(); rec != nil {
			t.Errorf("printPRResult panicked: %v", rec)
		}
	}()

	result := &pr.PRResult{
		BranchName: "spectr/test-change",
	}

	// Should not panic
	printPRResult(result)
}

// TestPrintPRResult_WithPRURL tests output when PRURL is set.
func TestPrintPRResult_WithPRURL(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := &pr.PRResult{
		BranchName: "spectr/add-feature",
		PRURL:      "https://github.com/owner/repo/pull/123",
	}

	printPRResult(result)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains branch name
	if !bytes.Contains([]byte(output), []byte("spectr/add-feature")) {
		t.Errorf("Output should contain branch name, got: %s", output)
	}

	// Verify output contains PR URL
	if !bytes.Contains([]byte(output), []byte("https://github.com/owner/repo/pull/123")) {
		t.Errorf("Output should contain PR URL, got: %s", output)
	}

	// Verify output contains "PR created"
	if !bytes.Contains([]byte(output), []byte("PR created")) {
		t.Errorf("Output should contain 'PR created', got: %s", output)
	}
}

// TestPrintPRResult_WithManualURL tests output when ManualURL is set (no PRURL).
func TestPrintPRResult_WithManualURL(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := &pr.PRResult{
		BranchName: "spectr/add-feature",
		ManualURL:  "https://bitbucket.org/owner/repo/pull-requests/new",
	}

	printPRResult(result)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains branch name
	if !bytes.Contains([]byte(output), []byte("spectr/add-feature")) {
		t.Errorf("Output should contain branch name, got: %s", output)
	}

	// Verify output contains manual URL
	if !bytes.Contains(
		[]byte(output),
		[]byte("https://bitbucket.org/owner/repo/pull-requests/new"),
	) {
		t.Errorf("Output should contain manual URL, got: %s", output)
	}

	// Verify output contains "manually"
	if !bytes.Contains([]byte(output), []byte("manually")) {
		t.Errorf("Output should contain 'manually', got: %s", output)
	}
}

// TestPrintPRResult_WithArchivePath tests output when ArchivePath is set.
func TestPrintPRResult_WithArchivePath(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := &pr.PRResult{
		BranchName:  "spectr/add-feature",
		ArchivePath: "spectr/accepted/add-feature",
		PRURL:       "https://github.com/owner/repo/pull/456",
	}

	printPRResult(result)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains archive path
	if !bytes.Contains([]byte(output), []byte("spectr/accepted/add-feature")) {
		t.Errorf("Output should contain archive path, got: %s", output)
	}

	// Verify output contains "Archived to"
	if !bytes.Contains([]byte(output), []byte("Archived to")) {
		t.Errorf("Output should contain 'Archived to', got: %s", output)
	}
}

// TestPrintPRResult_WithCounts tests output when operation counts are present.
func TestPrintPRResult_WithCounts(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := &pr.PRResult{
		BranchName:  "spectr/add-feature",
		ArchivePath: "spectr/accepted/add-feature",
		Counts: archive.OperationCounts{
			Added:    3,
			Modified: 2,
			Removed:  1,
		},
		PRURL: "https://github.com/owner/repo/pull/789",
	}

	printPRResult(result)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains spec operations line
	if !bytes.Contains([]byte(output), []byte("Spec operations")) {
		t.Errorf("Output should contain 'Spec operations', got: %s", output)
	}

	// Verify counts are displayed (+3 ~2 -1)
	if !bytes.Contains([]byte(output), []byte("+3")) {
		t.Errorf("Output should contain '+3' for added, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("~2")) {
		t.Errorf("Output should contain '~2' for modified, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("-1")) {
		t.Errorf("Output should contain '-1' for removed, got: %s", output)
	}
}

// TestPrintPRResult_ZeroCounts tests that zero counts are not displayed.
func TestPrintPRResult_ZeroCounts(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := &pr.PRResult{
		BranchName: "spectr/add-feature",
		Counts: archive.OperationCounts{
			Added:    0,
			Modified: 0,
			Removed:  0,
		},
		PRURL: "https://github.com/owner/repo/pull/100",
	}

	printPRResult(result)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// With zero counts, "Spec operations" should NOT appear
	if bytes.Contains([]byte(output), []byte("Spec operations")) {
		t.Errorf("Output should NOT contain 'Spec operations' with zero counts, got: %s", output)
	}
}

// TestPrintPRResult_FullResult tests a complete result with all fields.
func TestPrintPRResult_FullResult(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := &pr.PRResult{
		PRURL:       "https://github.com/owner/repo/pull/42",
		BranchName:  "spectr/complete-test",
		ArchivePath: "spectr/accepted/complete-test",
		Counts: archive.OperationCounts{
			Added:    5,
			Modified: 3,
			Removed:  1,
			Renamed:  0,
		},
		Platform: git.PlatformGitHub,
	}

	printPRResult(result)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify all key elements are present
	expectedElements := []string{
		"spectr/complete-test",
		"spectr/accepted/complete-test",
		"https://github.com/owner/repo/pull/42",
		"Branch:",
		"Archived to:",
		"PR created:",
		"Spec operations:",
	}

	for _, expected := range expectedElements {
		if !bytes.Contains([]byte(output), []byte(expected)) {
			t.Errorf("Output should contain %q, got: %s", expected, output)
		}
	}
}

// TestCLIHasPRCommand verifies PR command is registered in CLI.
func TestCLIHasPRCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	prField := val.FieldByName("PR")

	if !prField.IsValid() {
		t.Fatal("CLI struct does not have PR field")
	}

	// Check the type
	if prField.Type().Name() != "PRCmd" {
		t.Errorf("PR field type: got %s, want PRCmd", prField.Type().Name())
	}
}
