package cmd

import (
	"errors"
	"reflect"
	"testing"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

func TestListCmdStructure(t *testing.T) {
	cmd := &ListCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Specs field exists
	specsField := val.FieldByName("Specs")
	if !specsField.IsValid() {
		t.Error(
			"ListCmd does not have Specs field",
		)
	}

	// Check Long field exists
	longField := val.FieldByName("Long")
	if !longField.IsValid() {
		t.Error(
			"ListCmd does not have Long field",
		)
	}

	// Check JSON field exists
	jsonField := val.FieldByName("JSON")
	if !jsonField.IsValid() {
		t.Error(
			"ListCmd does not have JSON field",
		)
	}

	// Check Stdout field exists
	stdoutField := val.FieldByName("Stdout")
	if !stdoutField.IsValid() {
		t.Error(
			"ListCmd does not have Stdout field",
		)
	}
}

func TestCLIHasListCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	listField := val.FieldByName("List")

	if !listField.IsValid() {
		t.Fatal(
			"CLI struct does not have List field",
		)
	}

	// Check the type
	if listField.Type().Name() != "ListCmd" {
		t.Errorf(
			"List field type: got %s, want ListCmd",
			listField.Type().Name(),
		)
	}
}

// TestListCmd_HasLsAlias verifies that the list command has the "ls" alias.
// This ensures "spectr ls" works identically to "spectr list".
func TestListCmd_HasLsAlias(t *testing.T) {
	cliType := reflect.TypeOf(CLI{})
	listField, found := cliType.FieldByName(
		"List",
	)
	if !found {
		t.Fatal("CLI does not have List field")
	}

	// Get the aliases tag
	tag := listField.Tag
	aliases := tag.Get("aliases")
	if aliases == "" {
		t.Fatal(
			"List field does not have aliases tag",
		)
	}

	if aliases != "ls" {
		t.Errorf(
			"List aliases = %q, want %q",
			aliases,
			"ls",
		)
	}
}

// TestListCmd_StdoutRequiresInteractive verifies that --stdout returns an error
// when used without --interactive (-I).
func TestListCmd_StdoutRequiresInteractive(
	t *testing.T,
) {
	cmd := &ListCmd{
		Stdout:      true,
		Interactive: false,
	}

	err := cmd.Run()
	if err == nil {
		t.Error(
			"Expected error when --stdout used without --interactive",
		)
	}

	// Check it's the correct error type
	var reqErr *specterrs.RequiresFlagError
	if !errors.As(err, &reqErr) {
		t.Errorf(
			"Expected RequiresFlagError, got: %T",
			err,
		)
	}
}

// TestListCmd_StdoutIncompatibleWithJSON verifies that --stdout returns an error
// when used together with --json since they are mutually exclusive.
func TestListCmd_StdoutIncompatibleWithJSON(
	t *testing.T,
) {
	cmd := &ListCmd{
		Stdout:      true,
		Interactive: true,
		JSON:        true,
	}

	err := cmd.Run()
	if err == nil {
		t.Error(
			"Expected error when --stdout used with --json",
		)
	}

	// Check it's the correct error type
	var incompatErr *specterrs.IncompatibleFlagsError
	if !errors.As(err, &incompatErr) {
		t.Errorf(
			"Expected IncompatibleFlagsError, got: %T",
			err,
		)
	}
}
