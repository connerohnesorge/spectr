package providers

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/ralph"
)

const testOSWindows = "windows"

// TestIsBinaryAvailable tests the IsBinaryAvailable function with various inputs.
func TestIsBinaryAvailable(t *testing.T) {
	tests := []struct {
		name       string
		binaryName string
		want       bool
	}{
		{
			name:       "empty string",
			binaryName: "",
			want:       false,
		},
		{
			name:       "non-existent binary",
			binaryName: "this-binary-definitely-does-not-exist-12345",
			want:       false,
		},
		{
			name:       "sh exists on unix",
			binaryName: "sh",
			want:       runtime.GOOS != testOSWindows,
		},
		{
			name:       "bash exists on unix",
			binaryName: "bash",
			want:       runtime.GOOS != testOSWindows,
		},
		{
			name:       "ls exists on unix",
			binaryName: "ls",
			want:       runtime.GOOS != testOSWindows,
		},
		{
			name:       "go should exist in test environment",
			binaryName: "go",
			want:       true, // go must be available to run tests
		},
		{
			name:       "special characters",
			binaryName: "binary-with-!@#$%",
			want:       false,
		},
		{
			name:       "path separator in name",
			binaryName: "/usr/bin/sh",
			want:       false, // Should not find absolute paths
		},
		{
			name:       "binary with extension",
			binaryName: "binary.exe",
			want:       false,
		},
		{
			name:       "whitespace",
			binaryName: "   ",
			want:       false,
		},
		{
			name:       "very long name",
			binaryName: strings.Repeat("a", 1000),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBinaryAvailable(tt.binaryName)
			if got != tt.want {
				t.Errorf("IsBinaryAvailable(%q) = %v, want %v", tt.binaryName, got, tt.want)
			}
		})
	}
}

// TestIsBinaryAvailable_ActualBinaries tests with binaries that should actually exist.
func TestIsBinaryAvailable_ActualBinaries(t *testing.T) {
	// Use exec.LookPath to find a binary that actually exists on the system
	// Then verify IsBinaryAvailable returns true for it
	commonBinaries := []string{"sh", "bash", "ls", "cat", "echo", "pwd"}

	if runtime.GOOS == testOSWindows {
		commonBinaries = []string{"cmd", "powershell"}
	}

	foundOne := false
	for _, binary := range commonBinaries {
		path, err := exec.LookPath(binary)
		if err != nil || path == "" {
			continue
		}

		foundOne = true
		if !IsBinaryAvailable(binary) {
			t.Errorf(
				"IsBinaryAvailable(%q) = false, but exec.LookPath found it at %s",
				binary,
				path,
			)
		}
	}

	if !foundOne {
		t.Skip("No common binaries found on PATH to test with")
	}
}

// testProviderWithoutRalpher is a test provider that doesn't implement Ralpher.
type testProviderWithoutRalpher struct{}

func (testProviderWithoutRalpher) Initializers(context.Context, TemplateManager) []Initializer {
	return nil
}

// testProviderWithRalpher is a test provider that implements Ralpher.
type testProviderWithRalpher struct {
	binary string
}

func (testProviderWithRalpher) Initializers(context.Context, TemplateManager) []Initializer {
	return nil
}

func (m testProviderWithRalpher) Binary() string {
	return m.binary
}

//
//nolint:revive,gocritic // task parameter unused but required by interface
func (m testProviderWithRalpher) InvokeTask(
	ctx context.Context,
	_ *ralph.Task,
	_ string,
) (*exec.Cmd, error) {
	return exec.CommandContext(ctx, m.binary), nil
}

// TestIsRalpherAvailable tests the IsRalpherAvailable function.
func TestIsRalpherAvailable(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		want     bool
	}{
		{
			name:     "provider without Ralpher interface",
			provider: &testProviderWithoutRalpher{},
			want:     false,
		},
		{
			name:     "ralpher with non-existent binary",
			provider: &testProviderWithRalpher{binary: "this-binary-does-not-exist-xyz"},
			want:     false,
		},
		{
			name:     "ralpher with empty binary name",
			provider: &testProviderWithRalpher{binary: ""},
			want:     false,
		},
		{
			name:     "ralpher with existing binary",
			provider: &testProviderWithRalpher{binary: "go"}, // go must exist in test environment
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRalpherAvailable(tt.provider)
			if got != tt.want {
				t.Errorf("IsRalpherAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsRalpherAvailable_ClaudeProvider tests with the actual ClaudeProvider.
func TestIsRalpherAvailable_ClaudeProvider(t *testing.T) {
	provider := &ClaudeProvider{}

	// Verify ClaudeProvider implements Ralpher
	_, ok := any(provider).(Ralpher)
	if !ok {
		t.Fatal("ClaudeProvider does not implement Ralpher interface")
	}

	// Check if claude binary is available
	hasClaudeBinary := IsBinaryAvailable("claude")
	result := IsRalpherAvailable(provider)

	if result != hasClaudeBinary {
		t.Errorf(
			"IsRalpherAvailable(ClaudeProvider) = %v, but IsBinaryAvailable(\"claude\") = %v (should match)",
			result,
			hasClaudeBinary,
		)
	}
}

// TestValidateRalpher tests the ValidateRalpher function.
func TestValidateRalpher(t *testing.T) {
	tests := []struct {
		name        string
		provider    Provider
		wantErr     bool
		errContains string
	}{
		{
			name:        "provider without Ralpher interface",
			provider:    &testProviderWithoutRalpher{},
			wantErr:     true,
			errContains: "does not implement Ralpher interface",
		},
		{
			name:        "ralpher with empty binary name",
			provider:    &testProviderWithRalpher{binary: ""},
			wantErr:     true,
			errContains: "returned empty string",
		},
		{
			name:        "ralpher with non-existent binary",
			provider:    &testProviderWithRalpher{binary: "this-binary-does-not-exist-xyz"},
			wantErr:     true,
			errContains: "not found on PATH",
		},
		{
			name:        "ralpher with existing binary",
			provider:    &testProviderWithRalpher{binary: "go"}, // go must exist
			wantErr:     false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRalpher(tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRalpher() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if err == nil || tt.errContains == "" {
				return
			}

			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf(
					"ValidateRalpher() error = %v, want error containing %q",
					err,
					tt.errContains,
				)
			}
		})
	}
}

// TestValidateRalpher_ClaudeProvider tests ValidateRalpher with the actual ClaudeProvider.
func TestValidateRalpher_ClaudeProvider(t *testing.T) {
	provider := &ClaudeProvider{}

	err := ValidateRalpher(provider)

	// If claude binary is available, validation should pass
	// If not available, validation should fail with specific error
	hasClaudeBinary := IsBinaryAvailable("claude")

	if hasClaudeBinary && err != nil {
		t.Errorf("ValidateRalpher(ClaudeProvider) failed but claude binary is available: %v", err)
	}

	if !hasClaudeBinary && err == nil {
		t.Error("ValidateRalpher(ClaudeProvider) succeeded but claude binary is not available")
	}

	if err != nil && !strings.Contains(err.Error(), "not found on PATH") {
		t.Errorf("ValidateRalpher error should mention PATH, got: %v", err)
	}
}

// TestClaudeProvider_ImplementsRalpher verifies ClaudeProvider correctly implements Ralpher.
func TestClaudeProvider_ImplementsRalpher(t *testing.T) {
	provider := &ClaudeProvider{}

	// Type assertion should succeed
	ralpher, ok := any(provider).(Ralpher)
	if !ok {
		t.Fatal("ClaudeProvider does not implement Ralpher interface")
	}

	// Binary() should return "claude"
	if binary := ralpher.Binary(); binary != "claude" {
		t.Errorf("ClaudeProvider.Binary() = %q, want \"claude\"", binary)
	}

	// InvokeTask should return a valid command
	ctx := context.Background()
	task := ralph.Task{
		ID:          "1.1",
		Section:     "Test",
		Description: "Test task",
		Status:      "pending",
	}
	prompt := "Test prompt"

	cmd, err := ralpher.InvokeTask(ctx, &task, prompt)
	if err != nil {
		t.Fatalf("ClaudeProvider.InvokeTask() error = %v, want nil", err)
	}
	if cmd == nil {
		t.Fatal("ClaudeProvider.InvokeTask() returned nil command")
	}

	// Verify command is configured correctly
	if cmd.Path == "" && len(cmd.Args) == 0 {
		t.Error("ClaudeProvider.InvokeTask() returned command with no path or args")
	}
}

// TestEdgeCases_BinaryNames tests edge cases in binary name handling.
func TestEdgeCases_BinaryNames(t *testing.T) {
	tests := []struct {
		name       string
		binaryName string
	}{
		{"null bytes", "binary\x00name"},
		{"newlines", "binary\nname"},
		{"tabs", "binary\tname"},
		{"unicode", "binary-λ-test"},
		{"leading dash", "-binary"},
		{"trailing dash", "binary-"},
		{"double dash", "binary--name"},
		{"dot prefix", ".binary"},
		{"dot suffix", "binary."},
		{"multiple dots", "binary.name.test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic and should return false for weird names
			result := IsBinaryAvailable(tt.binaryName)
			if !result {
				return
			}

			// If it somehow returns true, verify with exec.LookPath
			_, err := exec.LookPath(tt.binaryName)
			if err != nil {
				t.Errorf(
					"IsBinaryAvailable(%q) = true, but exec.LookPath failed: %v",
					tt.binaryName,
					err,
				)
			}
		})
	}
}

// TestConcurrency_IsBinaryAvailable tests concurrent calls to IsBinaryAvailable.
//
//nolint:revive // t unused but required by test signature
func TestConcurrency_IsBinaryAvailable(_ *testing.T) {
	// Run multiple goroutines checking binary availability concurrently
	// This ensures thread safety of the function
	done := make(chan bool)
	binaries := []string{"go", "this-does-not-exist", "sh", "", "another-missing"}

	for range 10 {
		go func() {
			for range 100 {
				for _, binary := range binaries {
					_ = IsBinaryAvailable(binary)
				}
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}
}

// BenchmarkIsBinaryAvailable benchmarks the IsBinaryAvailable function.
func BenchmarkIsBinaryAvailable(b *testing.B) {
	benchmarks := []struct {
		name   string
		binary string
	}{
		{"existing", "go"},
		{"missing", "this-binary-does-not-exist"},
		{"empty", ""},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for range b.N {
				IsBinaryAvailable(bm.binary)
			}
		})
	}
}
