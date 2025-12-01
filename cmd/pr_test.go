package cmd

import (
	"reflect"
	"testing"
)

func TestCLIHasPRField(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	prField := val.FieldByName("PR")

	if !prField.IsValid() {
		t.Fatal("CLI struct does not have PR field")
	}

	if prField.Type().Name() != "PRCmd" {
		t.Errorf("PR field type: got %s, want PRCmd", prField.Type().Name())
	}
}

func TestPRCmdStructure(t *testing.T) {
	cmd := &PRCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Archive subcommand exists
	archiveField := val.FieldByName("Archive")
	if !archiveField.IsValid() {
		t.Error("PRCmd does not have Archive field")
	}
	if archiveField.Type().Name() != "PRArchiveCmd" {
		t.Errorf("Archive field type: got %s, want PRArchiveCmd", archiveField.Type().Name())
	}

	// Check New subcommand exists
	newField := val.FieldByName("New")
	if !newField.IsValid() {
		t.Error("PRCmd does not have New field")
	}
	if newField.Type().Name() != "PRNewCmd" {
		t.Errorf("New field type: got %s, want PRNewCmd", newField.Type().Name())
	}
}

func TestPRArchiveCmdStructure(t *testing.T) {
	cmd := &PRArchiveCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check ChangeID field
	changeIDField := val.FieldByName("ChangeID")
	if !changeIDField.IsValid() {
		t.Error("PRArchiveCmd does not have ChangeID field")
	}

	// Check Base field
	baseField := val.FieldByName("Base")
	if !baseField.IsValid() {
		t.Error("PRArchiveCmd does not have Base field")
	}

	// Check Draft field
	draftField := val.FieldByName("Draft")
	if !draftField.IsValid() {
		t.Error("PRArchiveCmd does not have Draft field")
	}

	// Check Force field
	forceField := val.FieldByName("Force")
	if !forceField.IsValid() {
		t.Error("PRArchiveCmd does not have Force field")
	}

	// Check DryRun field
	dryRunField := val.FieldByName("DryRun")
	if !dryRunField.IsValid() {
		t.Error("PRArchiveCmd does not have DryRun field")
	}

	// Check SkipSpecs field
	skipSpecsField := val.FieldByName("SkipSpecs")
	if !skipSpecsField.IsValid() {
		t.Error("PRArchiveCmd does not have SkipSpecs field")
	}
}

func TestPRNewCmdStructure(t *testing.T) {
	cmd := &PRNewCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check ChangeID field
	changeIDField := val.FieldByName("ChangeID")
	if !changeIDField.IsValid() {
		t.Error("PRNewCmd does not have ChangeID field")
	}

	// Check Base field
	baseField := val.FieldByName("Base")
	if !baseField.IsValid() {
		t.Error("PRNewCmd does not have Base field")
	}

	// Check Draft field
	draftField := val.FieldByName("Draft")
	if !draftField.IsValid() {
		t.Error("PRNewCmd does not have Draft field")
	}

	// Check Force field
	forceField := val.FieldByName("Force")
	if !forceField.IsValid() {
		t.Error("PRNewCmd does not have Force field")
	}

	// Check DryRun field
	dryRunField := val.FieldByName("DryRun")
	if !dryRunField.IsValid() {
		t.Error("PRNewCmd does not have DryRun field")
	}
}

func TestPRArchiveCmdHasRunMethod(t *testing.T) {
	cmd := &PRArchiveCmd{}
	val := reflect.ValueOf(cmd)

	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal("PRArchiveCmd does not have Run method")
	}

	runType := runMethod.Type()
	if runType.NumOut() != 1 {
		t.Errorf("Run method should return 1 value, got %d", runType.NumOut())
	}

	if runType.NumOut() > 0 && runType.Out(0).Name() != "error" {
		t.Errorf("Run method should return error, got %s", runType.Out(0).Name())
	}
}

func TestPRNewCmdHasRunMethod(t *testing.T) {
	cmd := &PRNewCmd{}
	val := reflect.ValueOf(cmd)

	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal("PRNewCmd does not have Run method")
	}

	runType := runMethod.Type()
	if runType.NumOut() != 1 {
		t.Errorf("Run method should return 1 value, got %d", runType.NumOut())
	}

	if runType.NumOut() > 0 && runType.Out(0).Name() != "error" {
		t.Errorf("Run method should return error, got %s", runType.Out(0).Name())
	}
}

func TestPRArchiveCmdTags(t *testing.T) {
	cmdType := reflect.TypeOf(PRArchiveCmd{})

	tests := []struct {
		fieldName   string
		expectedTag string
	}{
		{
			"ChangeID",
			`arg:"" optional:"" predictor:"changeID" help:"Change ID to archive (supports partial matching)"`,
		},
		{
			"Base",
			`name:"base" short:"b" help:"Base branch for PR (default: auto-detect main/master)"`,
		},
		{"Draft", `name:"draft" short:"d" help:"Create as draft PR"`},
		{"Force", `name:"force" short:"f" help:"Force overwrite existing branch"`},
		{"DryRun", `name:"dry-run" help:"Show what would be done without executing"`},
		{"SkipSpecs", `name:"skip-specs" help:"Skip spec merging during archive"`},
	}

	for _, tc := range tests {
		field, found := cmdType.FieldByName(tc.fieldName)
		if !found {
			t.Errorf("Field %s not found", tc.fieldName)

			continue
		}

		tag := string(field.Tag)
		if tag != tc.expectedTag {
			t.Errorf(
				"Field %s tag mismatch:\n  got:  %s\n  want: %s",
				tc.fieldName,
				tag,
				tc.expectedTag,
			)
		}
	}
}

func TestPRNewCmdTags(t *testing.T) {
	cmdType := reflect.TypeOf(PRNewCmd{})

	tests := []struct {
		fieldName   string
		expectedTag string
	}{
		{
			"ChangeID",
			`arg:"" optional:"" predictor:"changeID" help:"Change ID (supports partial matching)"`,
		},
		{
			"Base",
			`name:"base" short:"b" help:"Base branch for PR (default: auto-detect main/master)"`,
		},
		{"Draft", `name:"draft" short:"d" help:"Create as draft PR"`},
		{"Force", `name:"force" short:"f" help:"Force overwrite existing branch"`},
		{"DryRun", `name:"dry-run" help:"Show what would be done without executing"`},
	}

	for _, tc := range tests {
		field, found := cmdType.FieldByName(tc.fieldName)
		if !found {
			t.Errorf("Field %s not found", tc.fieldName)

			continue
		}

		tag := string(field.Tag)
		if tag != tc.expectedTag {
			t.Errorf(
				"Field %s tag mismatch:\n  got:  %s\n  want: %s",
				tc.fieldName,
				tag,
				tc.expectedTag,
			)
		}
	}
}
