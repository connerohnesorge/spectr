package cmd

import (
	"reflect"
	"testing"
)

func TestListCmdStructure(t *testing.T) {
	cmd := &ListCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Specs field exists
	specsField := val.FieldByName("Specs")
	if !specsField.IsValid() {
		t.Error("ListCmd does not have Specs field")
	}

	// Check Long field exists
	longField := val.FieldByName("Long")
	if !longField.IsValid() {
		t.Error("ListCmd does not have Long field")
	}

	// Check JSON field exists
	jsonField := val.FieldByName("JSON")
	if !jsonField.IsValid() {
		t.Error("ListCmd does not have JSON field")
	}
}

func TestCLIHasListCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	listField := val.FieldByName("List")

	if !listField.IsValid() {
		t.Fatal("CLI struct does not have List field")
	}

	// Check the type
	if listField.Type().Name() != "ListCmd" {
		t.Errorf("List field type: got %s, want ListCmd", listField.Type().Name())
	}
}

// TestListCmd_HasLsAlias verifies that the list command has the "ls" alias.
// This ensures "spectr ls" works identically to "spectr list".
func TestListCmd_HasLsAlias(t *testing.T) {
	cliType := reflect.TypeOf(CLI{})
	listField, found := cliType.FieldByName("List")
	if !found {
		t.Fatal("CLI does not have List field")
	}

	// Get the aliases tag
	tag := listField.Tag
	aliases := tag.Get("aliases")
	if aliases == "" {
		t.Fatal("List field does not have aliases tag")
	}

	if aliases != "ls" {
		t.Errorf("List aliases = %q, want %q", aliases, "ls")
	}
}
