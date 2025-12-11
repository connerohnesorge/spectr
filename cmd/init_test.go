package cmd

import (
	"reflect"
	"testing"
)

func TestCLIStructure(t *testing.T) {
	// Check that CLI struct has Init field
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	initField := val.FieldByName("Init")

	if !initField.IsValid() {
		t.Fatal(
			"CLI struct does not have Init field",
		)
	}

	// Check the type
	if initField.Type().Name() != "InitCmd" {
		t.Errorf(
			"Init field type: got %s, want InitCmd",
			initField.Type().Name(),
		)
	}
}

func TestInitCmdStructure(t *testing.T) {
	cmd := &InitCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Path field exists
	pathField := val.FieldByName("Path")
	if !pathField.IsValid() {
		t.Error(
			"InitCmd does not have Path field",
		)
	}

	// Check PathFlag field exists
	pathFlagField := val.FieldByName("PathFlag")
	if !pathFlagField.IsValid() {
		t.Error(
			"InitCmd does not have PathFlag field",
		)
	}

	// Check Tools field exists
	toolsField := val.FieldByName("Tools")
	if !toolsField.IsValid() {
		t.Error(
			"InitCmd does not have Tools field",
		)
	}

	// Check NonInteractive field exists
	nonInteractiveField := val.FieldByName(
		"NonInteractive",
	)
	if !nonInteractiveField.IsValid() {
		t.Error(
			"InitCmd does not have NonInteractive field",
		)
	}
}

func TestInitCmdHasRunMethod(t *testing.T) {
	cmd := &InitCmd{}
	val := reflect.ValueOf(cmd)

	// Check that Run method exists
	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal(
			"InitCmd does not have Run method",
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
