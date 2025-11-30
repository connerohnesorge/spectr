package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWizardModel(t *testing.T) {
	// Test creating a new wizard model
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Verify initial state
	if wizard.step != StepIntro {
		t.Errorf("Expected initial step to be StepIntro, got %v", wizard.step)
	}

	if wizard.projectPath != "/tmp/test-project" {
		t.Errorf("Expected project path to be /tmp/test-project, got %s", wizard.projectPath)
	}

	if wizard.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", wizard.cursor)
	}

	if len(wizard.allProviders) == 0 {
		t.Error("Expected allProviders to be populated")
	}
}

func TestWizardStepTransitions(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Test intro to select
	if wizard.step != StepIntro {
		t.Error("Expected initial step to be StepIntro")
	}

	// Simulate pressing enter on intro
	wizard.step = StepSelect
	if wizard.step != StepSelect {
		t.Error("Expected step to transition to StepSelect")
	}

	// Test provider selection
	wizard.selectedProviders["claude-code"] = true
	if !wizard.selectedProviders["claude-code"] {
		t.Error("Expected claude-code to be selected")
	}

	// Test getting selected provider IDs
	selectedIDs := wizard.getSelectedProviderIDs()
	if len(selectedIDs) != 1 {
		t.Errorf("Expected 1 selected provider, got %d", len(selectedIDs))
	}
}

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestWizardRenderFunctions(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Test that render functions don't panic
	t.Run("RenderIntro", func(t *testing.T) {
		output := wizard.renderIntro()
		if output == "" {
			t.Error("Expected non-empty intro output")
		}
		if !contains(output, "Spectr") {
			t.Error("Expected intro to contain 'Spectr'")
		}
	})

	t.Run("RenderSelect", func(t *testing.T) {
		wizard.step = StepSelect
		output := wizard.renderSelect()
		if output == "" {
			t.Error("Expected non-empty select output")
		}
		if !contains(output, "Select AI Tools to Configure") {
			t.Error("Expected select screen to contain 'Select AI Tools to Configure'")
		}
	})

	t.Run("RenderReview", func(t *testing.T) {
		wizard.step = StepReview
		wizard.selectedProviders["claude-code"] = true
		output := wizard.renderReview()
		if output == "" {
			t.Error("Expected non-empty review output")
		}
		if !contains(output, "Review Your Selections") {
			t.Error("Expected review screen to contain 'Review Your Selections'")
		}
	})

	t.Run("RenderExecute", func(t *testing.T) {
		wizard.step = StepExecute
		output := wizard.renderExecute()
		if output == "" {
			t.Error("Expected non-empty execute output")
		}
		if !contains(output, "Initializing") {
			t.Error("Expected execute screen to contain 'Initializing'")
		}
	})

	t.Run("RenderComplete", func(t *testing.T) {
		wizard.step = StepComplete
		wizard.executionResult = &ExecutionResult{
			CreatedFiles: []string{"spectr/project.md"},
			UpdatedFiles: make([]string, 0),
			Errors:       make([]string, 0),
		}
		output := wizard.renderComplete()
		if output == "" {
			t.Error("Expected non-empty complete output")
		}
		if !contains(output, "Successfully") {
			t.Error("Expected complete screen to contain 'Successfully'")
		}
	})
}

func TestGetSelectedProviderIDs(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Test with no selections
	ids := wizard.getSelectedProviderIDs()
	if len(ids) != 0 {
		t.Errorf("Expected 0 selected providers, got %d", len(ids))
	}

	// Test with some selections
	wizard.selectedProviders["claude-code"] = true
	wizard.selectedProviders["cline"] = true
	wizard.selectedProviders["cursor"] = true

	ids = wizard.getSelectedProviderIDs()
	if len(ids) != 3 {
		t.Errorf("Expected 3 selected providers, got %d", len(ids))
	}

	// Verify all selected IDs are present
	hasClaudeCode := false
	hasCline := false
	hasCursor := false
	for _, id := range ids {
		switch id {
		case "claude-code":
			hasClaudeCode = true
		case "cline":
			hasCline = true
		case "cursor":
			hasCursor = true
		}
	}

	if !hasClaudeCode || !hasCline || !hasCursor {
		t.Error("Not all selected provider IDs were returned")
	}
}

func TestNewWizardModelWithConfiguredProviders(t *testing.T) {
	// Create a temp directory with a configured provider
	tempDir := t.TempDir()

	// Create CLAUDE.md to make claude-code provider configured
	claudeFile := filepath.Join(tempDir, "CLAUDE.md")
	err := os.WriteFile(claudeFile, []byte("# Claude Configuration\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create CLAUDE.md: %v", err)
	}

	// Create .claude/commands/spectr/ directory and slash commands to fully configure Claude
	commandsDir := filepath.Join(tempDir, ".claude", "commands", "spectr")
	err = os.MkdirAll(commandsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create commands directory: %v", err)
	}

	// Create the three slash command files (in the spectr/ subdirectory)
	for _, cmdFile := range []string{
		"proposal.md",
		"sync.md",
		"apply.md",
	} {
		filePath := filepath.Join(commandsDir, cmdFile)
		err = os.WriteFile(filePath, []byte("# Command\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", cmdFile, err)
		}
	}

	// Create wizard model
	cmd := &InitCmd{Path: tempDir}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Verify claude-code is marked as configured
	if !wizard.configuredProviders["claude-code"] {
		t.Error("Expected claude-code to be marked as configured")
	}

	// Verify claude-code is pre-selected
	if !wizard.selectedProviders["claude-code"] {
		t.Error("Expected claude-code to be pre-selected")
	}
}

func TestNewWizardModelNoConfiguredProviders(t *testing.T) {
	// Create an empty temp directory
	tempDir := t.TempDir()

	// Create wizard model
	cmd := &InitCmd{Path: tempDir}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Verify no providers are marked as configured
	for providerID, isConfigured := range wizard.configuredProviders {
		if isConfigured {
			t.Errorf("Expected no providers to be configured, but %s is configured", providerID)
		}
	}

	// Verify no providers are pre-selected
	if len(wizard.selectedProviders) != 0 {
		t.Errorf(
			"Expected no providers to be selected, but got %d selected",
			len(wizard.selectedProviders),
		)
	}
}

func TestRenderSelectShowsConfiguredIndicator(t *testing.T) {
	// Create wizard with manually set configured provider
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf("Failed to create wizard model: %v", err)
	}

	// Manually mark claude-code as configured
	wizard.configuredProviders["claude-code"] = true
	wizard.selectedProviders["claude-code"] = true
	wizard.step = StepSelect

	// Render the select view
	output := wizard.renderSelect()

	// Verify output contains "(configured)" indicator
	if !strings.Contains(output, "(configured)") {
		t.Error("Expected select view to contain '(configured)' indicator for configured providers")
	}

	// Verify output contains "Claude Code"
	if !strings.Contains(output, "Claude Code") {
		t.Error("Expected select view to contain 'Claude Code' provider name")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
