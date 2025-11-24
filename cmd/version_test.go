package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// TestVersionCmdStructure verifies that VersionCmd has the required fields.
func TestVersionCmdStructure(t *testing.T) {
	cmd := &VersionCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check Short field exists
	shortField := val.FieldByName("Short")
	if !shortField.IsValid() {
		t.Error("VersionCmd does not have Short field")
	}

	// Check JSON field exists
	jsonField := val.FieldByName("JSON")
	if !jsonField.IsValid() {
		t.Error("VersionCmd does not have JSON field")
	}
}

// TestCLIHasVersionCommand verifies that the CLI struct includes VersionCmd.
func TestCLIHasVersionCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	versionField := val.FieldByName("Version")

	if !versionField.IsValid() {
		t.Fatal("CLI struct does not have Version field")
	}

	// Check the type
	if versionField.Type().Name() != "VersionCmd" {
		t.Errorf("Version field type: got %s, want VersionCmd", versionField.Type().Name())
	}
}

// TestVersionCmdRunMethod verifies that VersionCmd has a Run() method.
func TestVersionCmdRunMethod(t *testing.T) {
	cmd := &VersionCmd{}
	val := reflect.ValueOf(cmd)

	// Check if Run method exists
	runMethod := val.MethodByName("Run")
	if !runMethod.IsValid() {
		t.Fatal("VersionCmd does not have Run method")
	}

	// Verify method signature: func() error
	methodType := runMethod.Type()
	if methodType.NumIn() != 0 {
		t.Errorf("Run method should have 0 input parameters, got %d", methodType.NumIn())
	}
	if methodType.NumOut() != 1 {
		t.Errorf("Run method should have 1 output parameter, got %d", methodType.NumOut())
	}
}

// TestVersionCmdRun tests the Run method with different flag combinations.
// This test verifies that the version command can execute successfully
// and produce output in various formats.
func TestVersionCmdRun(t *testing.T) {
	tests := []struct {
		name          string
		short         bool
		jsonFlag      bool
		expectContain []string
		expectJSON    bool
	}{
		{
			name:     "default output",
			short:    false,
			jsonFlag: false,
			expectContain: []string{
				"Spectr version",
				"Commit:",
				"Build date:",
				"Go version:",
				"OS/Arch:",
			},
			expectJSON: false,
		},
		{
			name:          "short output",
			short:         true,
			jsonFlag:      false,
			expectContain: nil, // Will check for single line
			expectJSON:    false,
		},
		{
			name:          "JSON output",
			short:         false,
			jsonFlag:      true,
			expectContain: nil,
			expectJSON:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run command
			cmd := &VersionCmd{
				Short: tt.short,
				JSON:  tt.jsonFlag,
			}
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			// Restore stdout and read output
			_ = w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			// Verify output
			if tt.expectJSON {
				// Verify JSON structure
				var result map[string]string
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
				}

				// Check required fields
				requiredFields := []string{"version", "commit", "date", "goVersion", "os", "arch"}
				for _, field := range requiredFields {
					if _, ok := result[field]; !ok {
						t.Errorf("JSON output missing field: %s", field)
					}
				}

				// Verify runtime values
				if result["goVersion"] != runtime.Version() {
					t.Errorf("goVersion = %v, want %v", result["goVersion"], runtime.Version())
				}
				if result["os"] != runtime.GOOS {
					t.Errorf("os = %v, want %v", result["os"], runtime.GOOS)
				}
				if result["arch"] != runtime.GOARCH {
					t.Errorf("arch = %v, want %v", result["arch"], runtime.GOARCH)
				}
			} else {
				// Verify text output contains expected strings
				for _, expected := range tt.expectContain {
					if !strings.Contains(output, expected) {
						t.Errorf("Output does not contain %q\nGot: %s", expected, output)
					}
				}

				// Verify short output is single line
				if !tt.short {
					return
				}

				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) != 1 {
					t.Errorf("Short output should be single line, got %d lines", len(lines))
				}
				// Verify output is not empty
				if strings.TrimSpace(output) == "" {
					t.Error("Short output should not be empty")
				}
			}
		})
	}
}

// TestVersionCmdRunExecution tests that the version command can execute
// without errors. This is a basic smoke test to ensure the command works.
func TestVersionCmdRunExecution(t *testing.T) {
	// Capture and discard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command with default flags
	cmd := &VersionCmd{}
	err := cmd.Run()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout
	_, _ = io.Copy(io.Discard, r)

	// Verify no error
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}
}

// TestVersionOutputFormats tests different output formats produce valid output.
func TestVersionOutputFormats(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *VersionCmd
		validate func(t *testing.T, output string)
	}{
		{
			name: "default format has multiple lines",
			cmd:  &VersionCmd{},
			validate: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) < 3 {
					t.Errorf("Default output should have at least 3 lines, got %d", len(lines))
				}
			},
		},
		{
			name: "short format is single line",
			cmd:  &VersionCmd{Short: true},
			validate: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) != 1 {
					t.Errorf("Short output should be exactly 1 line, got %d", len(lines))
				}
			},
		},
		{
			name: "JSON format is valid JSON",
			cmd:  &VersionCmd{JSON: true},
			validate: func(t *testing.T, output string) {
				var result map[string]string
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("JSON output is not valid: %v\nOutput: %s", err, output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run command
			err := tt.cmd.Run()
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			// Restore stdout and read output
			_ = w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			// Validate output
			tt.validate(t, output)
		})
	}
}
