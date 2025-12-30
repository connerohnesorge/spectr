package cmd

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"syscall"
	"testing"
)

const errorTypeName = "error"

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
		runType.Out(0).Name() != errorTypeName {
		t.Errorf(
			"Run method should return error, got %s",
			runType.Out(0).Name(),
		)
	}
}

func TestIsTTYError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "wrapped ENXIO from bubbletea (real error format)",
			err: fmt.Errorf("could not open a new TTY: %w",
				&os.PathError{
					Op:   "open",
					Path: "/dev/tty",
					Err:  syscall.ENXIO,
				}),
			expected: true,
		},
		{
			name:     "direct ENXIO syscall error",
			err:      syscall.ENXIO,
			expected: true,
		},
		{
			name:     "direct ENOTTY syscall error",
			err:      syscall.ENOTTY,
			expected: true,
		},
		{
			name: "PathError with /dev/tty path",
			err: &os.PathError{
				Op:   "open",
				Path: "/dev/tty",
				Err:  syscall.EACCES,
			},
			expected: true,
		},
		{
			name: "PathError with Windows CONIN$ path",
			err: &os.PathError{
				Op:   "open",
				Path: "CONIN$",
				Err:  syscall.EACCES,
			},
			expected: true,
		},
		{
			name: "deeply wrapped ENXIO error",
			err: fmt.Errorf("wizard failed: %w",
				fmt.Errorf("could not open a new TTY: %w",
					&os.PathError{
						Op:   "open",
						Path: "/dev/tty",
						Err:  syscall.ENXIO,
					})),
			expected: true,
		},
		{
			name: "PathError with non-TTY path",
			err: &os.PathError{
				Op:   "open",
				Path: "/etc/passwd",
				Err:  syscall.EACCES,
			},
			expected: false,
		},
		{
			name:     "regular error without TTY reference",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "string error mentioning tty (should NOT match)",
			err:      errors.New("pretty formatted output"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTTYError(tt.err)
			if result != tt.expected {
				t.Errorf(
					"isTTYError(%v) = %v, want %v",
					tt.err,
					result,
					tt.expected,
				)
			}
		})
	}
}
