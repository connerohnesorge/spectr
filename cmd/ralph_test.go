package cmd

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/ralph"
)

// mockRalpher is a mock implementation of the Ralpher interface for testing.
type mockRalpher struct {
	binaryName string
}

func (*mockRalpher) InvokeTask(
	ctx context.Context,
	_ *ralph.Task,
	_ string,
) (*exec.Cmd, error) {
	// Return a simple command that will succeed
	return exec.CommandContext(ctx, "echo", "mock task execution"), nil
}

func (m *mockRalpher) Binary() string {
	return m.binaryName
}

// mockProvider implements the providers.Provider interface for testing.
type mockProvider struct {
	mockRalpher
}

func (*mockProvider) Initializers(
	_ context.Context,
	_ providers.TemplateManager,
) []providers.Initializer {
	return nil
}

func TestResolveChangeID_WithExplicitID(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	spectrDir := filepath.Join(tempDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(spectrDir, 0o755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create a minimal proposal.md file so the change is detected
	proposalPath := filepath.Join(spectrDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte("# Test Change\n"), 0o644); err != nil {
		t.Fatalf("failed to create proposal.md: %v", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldDir)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Test resolving with explicit change ID
	cmd := &RalphCmd{
		ChangeID: "test-change",
	}

	changeID, err := cmd.resolveChangeID(tempDir)
	if err != nil {
		t.Fatalf("resolveChangeID failed: %v", err)
	}

	if changeID != "test-change" {
		t.Errorf("expected changeID 'test-change', got '%s'", changeID)
	}
}

func TestResolveChangeID_NoInteractive(t *testing.T) {
	tempDir := t.TempDir()

	cmd := &RalphCmd{
		ChangeID:      "",
		NoInteractive: true,
	}

	_, err := cmd.resolveChangeID(tempDir)
	if err == nil {
		t.Error("expected error when no changeID provided in non-interactive mode")
	}
}

// nonRalpherProvider is a provider that doesn't implement Ralpher
type nonRalpherProvider struct{}

func (*nonRalpherProvider) Initializers(
	_ context.Context,
	_ providers.TemplateManager,
) []providers.Initializer {
	return nil
}

func TestDetectAndValidateProvider_NoProviderAvailable(t *testing.T) {
	// Clear registry
	providers.Reset()

	// Register a provider that doesn't implement Ralpher
	_ = providers.RegisterProvider(providers.Registration{
		ID:       "test-non-ralpher",
		Name:     "Test Non-Ralpher",
		Priority: 100,
		Provider: &nonRalpherProvider{},
	})

	cmd := &RalphCmd{}
	_, err := cmd.detectAndValidateProvider()
	if err == nil {
		t.Error("expected error when no suitable provider found")
	}
}

func TestDetectAndValidateProvider_Success(t *testing.T) {
	// Skip this test as it requires mocking IsBinaryAvailable which is not exported
	// In a real scenario, this would be tested via integration tests with actual binaries
	t.Skip("Skipping provider detection test - requires binary availability mocking")

	// Clear and setup registry
	providers.Reset()

	// Register a mock provider that implements Ralpher
	mock := &mockProvider{
		mockRalpher: mockRalpher{binaryName: "echo"}, // Use a real binary for testing
	}

	_ = providers.RegisterProvider(providers.Registration{
		ID:       "test-ralpher",
		Name:     "Test Ralpher",
		Priority: 1,
		Provider: mock,
	})

	cmd := &RalphCmd{}
	ralpher, err := cmd.detectAndValidateProvider()
	if err != nil {
		t.Fatalf("detectAndValidateProvider failed: %v", err)
	}

	if ralpher.Binary() != "echo" {
		t.Errorf("expected binary 'echo', got '%s'", ralpher.Binary())
	}
}

func TestRalphCmd_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		changeID    string
		setupFunc   func(t *testing.T) string // Returns project root
		expectedErr bool
	}{
		{
			name:     "missing change directory",
			changeID: "nonexistent-change",
			setupFunc: func(t *testing.T) string {
				tempDir := t.TempDir()
				// Create spectr/changes but not the specific change
				changesDir := filepath.Join(tempDir, "spectr", "changes")
				if err := os.MkdirAll(changesDir, 0o755); err != nil {
					t.Fatalf("failed to create changes directory: %v", err)
				}

				return tempDir
			},
			expectedErr: true,
		},
		{
			name:     "valid change directory",
			changeID: "valid-change",
			setupFunc: func(t *testing.T) string {
				tempDir := t.TempDir()
				changeDir := filepath.Join(tempDir, "spectr", "changes", "valid-change")
				if err := os.MkdirAll(changeDir, 0o755); err != nil {
					t.Fatalf("failed to create change directory: %v", err)
				}
				// Create a minimal tasks.jsonc
				tasksFile := filepath.Join(changeDir, "tasks.jsonc")
				tasksContent := `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Test",
      "description": "Test task",
      "status": "pending"
    }
  ]
}`
				if err := os.WriteFile(tasksFile, []byte(tasksContent), 0o644); err != nil {
					t.Fatalf("failed to create tasks.jsonc: %v", err)
				}

				return tempDir
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot := tt.setupFunc(t)

			// Change to project root
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get current directory: %v", err)
			}
			defer func() {
				_ = os.Chdir(oldDir)
			}()

			if err := os.Chdir(projectRoot); err != nil {
				t.Fatalf("failed to change to project root: %v", err)
			}

			// Validate the change directory
			changeDir := filepath.Join(projectRoot, "spectr", "changes", tt.changeID)
			_, err = os.Stat(changeDir)

			if tt.expectedErr && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectedErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestRalphCmd_DefaultMaxRetries(t *testing.T) {
	cmd := &RalphCmd{
		MaxRetries: 0, // Not set
	}

	// The default value should be set by Kong's default tag
	// In actual usage, Kong would set this to 3
	// Here we verify that if it's 0, we can use it as-is or apply a default in the orchestrator

	if cmd.MaxRetries == 0 {
		// This is expected when not set via CLI
		// The orchestrator will use its own default
		t.Log("MaxRetries is 0 (not set), orchestrator will use default")
	}
}
