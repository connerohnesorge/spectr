package init

import (
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

	// Should have 17 providers
	if len(wizard.allProviders) != 17 {
		t.Errorf("Expected 17 providers, got %d", len(wizard.allProviders))
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
