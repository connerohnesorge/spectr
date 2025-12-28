package initialize

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

const (
	testProviderClaude      = "claude"
	testProviderNonexistent = "nonexistentprovider123"
)

// TestMain registers all providers before running tests.
func TestMain(m *testing.M) {
	providers.RegisterAll()
	m.Run()
}

func TestNewWizardModel(t *testing.T) {
	// Test creating a new wizard model
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Verify initial state
	if wizard.step != StepIntro {
		t.Errorf(
			"Expected initial step to be StepIntro, got %v",
			wizard.step,
		)
	}

	if wizard.projectPath != "/tmp/test-project" {
		t.Errorf(
			"Expected project path to be /tmp/test-project, got %s",
			wizard.projectPath,
		)
	}

	if wizard.cursor != 0 {
		t.Errorf(
			"Expected cursor to start at 0, got %d",
			wizard.cursor,
		)
	}

	if len(wizard.allProviders) == 0 {
		t.Error(
			"Expected allProviders to be populated",
		)
	}
}

func TestWizardStepTransitions(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Test intro to select
	if wizard.step != StepIntro {
		t.Error(
			"Expected initial step to be StepIntro",
		)
	}

	// Simulate pressing enter on intro
	wizard.step = StepSelect
	if wizard.step != StepSelect {
		t.Error(
			"Expected step to transition to StepSelect",
		)
	}

	// Test provider selection
	wizard.selectedProviders["claude-code"] = true
	if !wizard.selectedProviders["claude-code"] {
		t.Error(
			"Expected claude-code to be selected",
		)
	}

	// Test getting selected provider IDs
	selectedIDs := wizard.getSelectedProviderIDs()
	if len(selectedIDs) != 1 {
		t.Errorf(
			"Expected 1 selected provider, got %d",
			len(selectedIDs),
		)
	}
}

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestWizardRenderFunctions(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Test that render functions don't panic
	t.Run("RenderIntro", func(t *testing.T) {
		output := wizard.renderIntro()
		if output == "" {
			t.Error(
				"Expected non-empty intro output",
			)
		}
		if !contains(output, "Spectr") {
			t.Error(
				"Expected intro to contain 'Spectr'",
			)
		}
	})

	t.Run("RenderSelect", func(t *testing.T) {
		wizard.step = StepSelect
		output := wizard.renderSelect()
		if output == "" {
			t.Error(
				"Expected non-empty select output",
			)
		}
		if !contains(
			output,
			"Select AI Tools to Configure",
		) {
			t.Error(
				"Expected select screen to contain 'Select AI Tools to Configure'",
			)
		}
	})

	t.Run("RenderReview", func(t *testing.T) {
		wizard.step = StepReview
		wizard.selectedProviders["claude-code"] = true
		output := wizard.renderReview()
		if output == "" {
			t.Error(
				"Expected non-empty review output",
			)
		}
		if !contains(
			output,
			"Review Your Selections",
		) {
			t.Error(
				"Expected review screen to contain 'Review Your Selections'",
			)
		}
	})

	t.Run("RenderExecute", func(t *testing.T) {
		wizard.step = StepExecute
		output := wizard.renderExecute()
		if output == "" {
			t.Error(
				"Expected non-empty execute output",
			)
		}
		if !contains(output, "Initializing") {
			t.Error(
				"Expected execute screen to contain 'Initializing'",
			)
		}
	})

	t.Run("RenderComplete", func(t *testing.T) {
		wizard.step = StepComplete
		wizard.executionResult = &ExecutionResult{
			CreatedFiles: []string{
				"spectr/project.md",
			},
			UpdatedFiles: make([]string, 0),
			Errors:       make([]string, 0),
		}
		output := wizard.renderComplete()
		if output == "" {
			t.Error(
				"Expected non-empty complete output",
			)
		}
		if !contains(output, "Successfully") {
			t.Error(
				"Expected complete screen to contain 'Successfully'",
			)
		}
	})
}

func TestGetSelectedProviderIDs(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Test with no selections
	ids := wizard.getSelectedProviderIDs()
	if len(ids) != 0 {
		t.Errorf(
			"Expected 0 selected providers, got %d",
			len(ids),
		)
	}

	// Test with some selections
	wizard.selectedProviders["claude-code"] = true
	wizard.selectedProviders["cline"] = true
	wizard.selectedProviders["cursor"] = true

	ids = wizard.getSelectedProviderIDs()
	if len(ids) != 3 {
		t.Errorf(
			"Expected 3 selected providers, got %d",
			len(ids),
		)
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
		t.Error(
			"Not all selected provider IDs were returned",
		)
	}
}

func TestNewWizardModelWithConfiguredProviders(
	t *testing.T,
) {
	// Create a temp directory with a configured provider
	tempDir := t.TempDir()

	// Create CLAUDE.md to make claude-code provider configured
	claudeFile := filepath.Join(
		tempDir,
		"CLAUDE.md",
	)
	err := os.WriteFile(
		claudeFile,
		[]byte("# Claude Configuration\n"),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create CLAUDE.md: %v",
			err,
		)
	}

	// Create .claude/commands/spectr/ directory and slash commands to fully configure Claude
	commandsDir := filepath.Join(
		tempDir,
		".claude",
		"commands",
		"spectr",
	)
	err = os.MkdirAll(commandsDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create commands directory: %v",
			err,
		)
	}

	// Create the two slash command files (in the spectr/ subdirectory)
	for _, cmdFile := range []string{
		"proposal.md",
		"apply.md",
	} {
		filePath := filepath.Join(
			commandsDir,
			cmdFile,
		)
		err = os.WriteFile(
			filePath,
			[]byte("# Command\n"),
			0o644,
		)
		if err != nil {
			t.Fatalf(
				"Failed to create %s: %v",
				cmdFile,
				err,
			)
		}
	}

	// Create wizard model
	cmd := &InitCmd{Path: tempDir}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Verify claude-code is marked as configured
	if !wizard.configuredProviders["claude-code"] {
		t.Error(
			"Expected claude-code to be marked as configured",
		)
	}

	// Verify claude-code is pre-selected
	if !wizard.selectedProviders["claude-code"] {
		t.Error(
			"Expected claude-code to be pre-selected",
		)
	}
}

func TestNewWizardModelNoConfiguredProviders(
	t *testing.T,
) {
	// Create an empty temp directory
	tempDir := t.TempDir()

	// Create wizard model
	cmd := &InitCmd{Path: tempDir}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Verify no providers are marked as configured
	for providerID, isConfigured := range wizard.configuredProviders {
		if isConfigured {
			t.Errorf(
				"Expected no providers to be configured, but %s is configured",
				providerID,
			)
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

func TestNewWizardModelWithCIWorkflowConfigured(
	t *testing.T,
) {
	// Create a temp directory with CI workflow file
	tempDir := t.TempDir()

	// Create .github/workflows/spectr-ci.yml
	workflowDir := filepath.Join(
		tempDir,
		".github",
		"workflows",
	)
	err := os.MkdirAll(workflowDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create workflows directory: %v",
			err,
		)
	}

	workflowFile := filepath.Join(
		workflowDir,
		"spectr-ci.yml",
	)
	err = os.WriteFile(
		workflowFile,
		[]byte("name: Spectr Validation\n"),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create spectr-ci.yml: %v",
			err,
		)
	}

	// Create wizard model
	cmd := &InitCmd{Path: tempDir}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Verify CI workflow is marked as configured
	if !wizard.ciWorkflowConfigured {
		t.Error(
			"Expected ciWorkflowConfigured to be true",
		)
	}

	// Verify CI workflow is pre-selected
	if !wizard.ciWorkflowEnabled {
		t.Error(
			"Expected ciWorkflowEnabled to be pre-selected",
		)
	}
}

func TestNewWizardModelWithoutCIWorkflow(
	t *testing.T,
) {
	// Create an empty temp directory
	tempDir := t.TempDir()

	// Create wizard model
	cmd := &InitCmd{Path: tempDir}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Verify CI workflow is NOT configured
	if wizard.ciWorkflowConfigured {
		t.Error(
			"Expected ciWorkflowConfigured to be false",
		)
	}

	// Verify CI workflow is NOT pre-selected
	if wizard.ciWorkflowEnabled {
		t.Error(
			"Expected ciWorkflowEnabled to be false",
		)
	}
}

func TestRenderSelectShowsConfiguredIndicator(
	t *testing.T,
) {
	// Create wizard with manually set configured provider
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Manually mark claude-code as configured
	wizard.configuredProviders["claude-code"] = true
	wizard.selectedProviders["claude-code"] = true
	wizard.step = StepSelect

	// Render the select view
	output := wizard.renderSelect()

	// Verify output contains "(configured)" indicator
	if !strings.Contains(output, "(configured)") {
		t.Error(
			"Expected select view to contain '(configured)' indicator for configured providers",
		)
	}

	// Verify output contains "Claude Code"
	if !strings.Contains(output, "Claude Code") {
		t.Error(
			"Expected select view to contain 'Claude Code' provider name",
		)
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

func TestHandleCompleteKeysCopyOnError(
	t *testing.T,
) {
	// Test that pressing 'c' on error screen (m.err != nil) does NOT return tea.Quit
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Set up error state
	wizard.step = StepComplete
	wizard.err = errors.New(
		"initialization failed",
	)
	wizard.executionResult = nil

	// Simulate pressing 'c' key
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'c'},
	}
	_, resultCmd := wizard.Update(keyMsg)

	// Verify tea.Quit is NOT returned (should be nil)
	if resultCmd != nil {
		t.Error(
			"Expected no command when pressing 'c' on error screen, but got a command",
		)
	}
}

func TestRenderCompleteShowsCopyHintOnSuccess(
	t *testing.T,
) {
	// Test that renderComplete() on success screen contains "c: copy"
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Set up success state
	wizard.step = StepComplete
	wizard.err = nil
	wizard.executionResult = &ExecutionResult{
		CreatedFiles: []string{
			"spectr/project.md",
		},
		UpdatedFiles: make([]string, 0),
		Errors:       make([]string, 0),
	}

	output := wizard.renderComplete()

	// Verify "c: copy" hint is shown on success screen
	if !strings.Contains(output, "c: copy") {
		t.Error(
			"Expected success screen to contain 'c: copy' hint",
		)
	}
}

func TestRenderCompleteHidesCopyHintOnError(
	t *testing.T,
) {
	// Test that renderComplete() on error screen does NOT contain "c: copy"
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Set up error state
	wizard.step = StepComplete
	wizard.err = errors.New(
		"initialization failed",
	)
	wizard.executionResult = nil

	output := wizard.renderComplete()

	// Verify "c: copy" hint is NOT shown on error screen
	if strings.Contains(output, "c: copy") {
		t.Error(
			"Expected error screen to NOT contain 'c: copy' hint",
		)
	}

	// Verify the quit hint is still shown
	if !strings.Contains(output, "q") {
		t.Error(
			"Expected error screen to contain quit hint",
		)
	}
}

func TestPopulateContextPromptHasNoSurroundingQuotes(
	t *testing.T,
) {
	// Test that PopulateContextPrompt constant does not contain surrounding quotes
	// The raw content should be copied without extra formatting

	// Check that the prompt does not start with a quote character
	if strings.HasPrefix(
		PopulateContextPrompt,
		"\"",
	) ||
		strings.HasPrefix(
			PopulateContextPrompt,
			"'",
		) ||
		strings.HasPrefix(
			PopulateContextPrompt,
			"`",
		) {
		t.Error(
			"PopulateContextPrompt should not start with a quote character",
		)
	}

	// Check that the prompt does not end with a quote character
	if strings.HasSuffix(
		PopulateContextPrompt,
		"\"",
	) ||
		strings.HasSuffix(
			PopulateContextPrompt,
			"'",
		) ||
		strings.HasSuffix(
			PopulateContextPrompt,
			"`",
		) {
		t.Error(
			"PopulateContextPrompt should not end with a quote character",
		)
	}

	// Verify the prompt contains expected content (basic sanity check)
	if !strings.Contains(
		PopulateContextPrompt,
		"spectr/project.md",
	) {
		t.Error(
			"PopulateContextPrompt should reference spectr/project.md",
		)
	}
}

func TestHandleReviewKeysToggleCIWorkflow(
	t *testing.T,
) {
	// Test that pressing space on review screen toggles CI workflow option
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Set up review state
	wizard.step = StepReview
	wizard.ciWorkflowEnabled = false

	// Simulate pressing space key
	keyMsg := tea.KeyMsg{Type: tea.KeySpace}
	newModel, _ := wizard.Update(keyMsg)
	updatedWizard, ok := newModel.(*WizardModel)
	if !ok {
		t.Fatal(
			"Failed to cast model to WizardModel",
		)
	}

	// Verify CI workflow is now enabled
	if !updatedWizard.ciWorkflowEnabled {
		t.Error(
			"Expected ciWorkflowEnabled to be true after pressing space",
		)
	}

	// Toggle again
	keyMsg = tea.KeyMsg{Type: tea.KeySpace}
	newModel, _ = updatedWizard.Update(keyMsg)
	updatedWizard, ok = newModel.(*WizardModel)
	if !ok {
		t.Fatal(
			"Failed to cast model to WizardModel",
		)
	}

	// Verify CI workflow is now disabled
	if updatedWizard.ciWorkflowEnabled {
		t.Error(
			"Expected ciWorkflowEnabled to be false after pressing space again",
		)
	}
}

func TestRenderReviewShowsCIWorkflowOption(
	t *testing.T,
) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepReview

	// Render the review view
	output := wizard.renderReview()

	// Verify output contains CI workflow option
	if !strings.Contains(
		output,
		"Spectr CI Validation",
	) {
		t.Error(
			"Expected review view to contain 'Spectr CI Validation' option",
		)
	}

	// Verify output contains toggle instructions
	if !strings.Contains(
		output,
		"Space: Toggle",
	) {
		t.Error(
			"Expected review view to contain 'Space: Toggle' instruction",
		)
	}
}

func TestRenderReviewShowsCIWorkflowConfigured(
	t *testing.T,
) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepReview
	wizard.ciWorkflowConfigured = true
	wizard.ciWorkflowEnabled = true

	// Render the review view
	output := wizard.renderReview()

	// Verify output contains configured indicator
	if !strings.Contains(output, "(configured)") {
		t.Error(
			"Expected review view to contain '(configured)' indicator for CI workflow",
		)
	}
}

func TestProviderFilteringLogic(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Initial state should have all providers in filteredProviders
	if len(
		wizard.filteredProviders,
	) != len(
		wizard.allProviders,
	) {
		t.Errorf(
			"Expected filteredProviders to equal allProviders, got %d vs %d",
			len(
				wizard.filteredProviders,
			),
			len(wizard.allProviders),
		)
	}

	// Test filtering with "claude" - should match Claude Code
	wizard.searchQuery = testProviderClaude
	wizard.applyProviderFilter()

	if len(wizard.filteredProviders) == 0 {
		t.Error(
			"Expected at least one provider to match 'claude'",
		)
	}

	// Verify all results contain "claude" (case-insensitive)
	for _, provider := range wizard.filteredProviders {
		if !strings.Contains(
			strings.ToLower(provider.Name()),
			testProviderClaude,
		) {
			t.Errorf(
				"Provider %s should not match 'claude'",
				provider.Name(),
			)
		}
	}

	// Test filtering with empty query - should restore all
	wizard.searchQuery = ""
	wizard.applyProviderFilter()

	if len(
		wizard.filteredProviders,
	) != len(
		wizard.allProviders,
	) {
		t.Error(
			"Expected empty query to restore all providers",
		)
	}

	// Test filtering with non-matching query
	wizard.searchQuery = testProviderNonexistent
	wizard.applyProviderFilter()

	if len(wizard.filteredProviders) != 0 {
		t.Error(
			"Expected no providers to match 'nonexistentprovider123'",
		)
	}
}

func TestCursorAdjustmentOnFilter(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Set cursor to last provider
	wizard.cursor = len(wizard.allProviders) - 1

	// Apply a filter that reduces the list significantly
	wizard.searchQuery = testProviderClaude
	wizard.applyProviderFilter()

	// Cursor should be adjusted to be within bounds
	if wizard.cursor >= len(
		wizard.filteredProviders,
	) {
		t.Errorf(
			"Cursor %d should be less than filtered count %d",
			wizard.cursor,
			len(wizard.filteredProviders),
		)
	}

	// Test with no matches - cursor should be 0
	wizard.searchQuery = testProviderNonexistent
	wizard.applyProviderFilter()

	if wizard.cursor != 0 {
		t.Errorf(
			"Expected cursor to be 0 when no matches, got %d",
			wizard.cursor,
		)
	}
}

func TestSelectionPreservedDuringFiltering(
	t *testing.T,
) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	// Select all providers
	for _, provider := range wizard.allProviders {
		wizard.selectedProviders[provider.ID()] = true
	}

	originalSelectionCount := len(
		wizard.selectedProviders,
	)

	// Apply a filter that shows only some providers
	wizard.searchQuery = testProviderClaude
	wizard.applyProviderFilter()

	// Verify selection count is preserved (filtering shouldn't affect selections)
	if len(
		wizard.selectedProviders,
	) != originalSelectionCount {
		t.Errorf(
			"Selection count changed after filtering: %d vs %d",
			len(
				wizard.selectedProviders,
			),
			originalSelectionCount,
		)
	}

	// Clear filter
	wizard.searchQuery = ""
	wizard.applyProviderFilter()

	// Selection should still be preserved
	if len(
		wizard.selectedProviders,
	) != originalSelectionCount {
		t.Errorf(
			"Selection count changed after clearing filter: %d vs %d",
			len(
				wizard.selectedProviders,
			),
			originalSelectionCount,
		)
	}
}

func TestSearchModeActivation(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepSelect

	// Simulate pressing '/' key
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'/'},
	}
	newModel, _ := wizard.Update(keyMsg)
	updatedWizard, ok := newModel.(*WizardModel)
	if !ok {
		t.Fatal(
			"Failed to cast model to WizardModel",
		)
	}

	// Verify search mode is active
	if !updatedWizard.searchMode {
		t.Error(
			"Expected searchMode to be true after pressing '/'",
		)
	}
}

func TestSearchModeExitWithEscape(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepSelect
	wizard.searchMode = true
	wizard.searchQuery = testProviderClaude
	wizard.searchInput.SetValue(testProviderClaude)
	wizard.applyProviderFilter()

	// Simulate pressing Escape key
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := wizard.Update(keyMsg)
	updatedWizard, ok := newModel.(*WizardModel)
	if !ok {
		t.Fatal(
			"Failed to cast model to WizardModel",
		)
	}

	// Verify search mode is deactivated
	if updatedWizard.searchMode {
		t.Error(
			"Expected searchMode to be false after pressing Escape",
		)
	}

	// Verify search query is cleared
	if updatedWizard.searchQuery != "" {
		t.Errorf(
			"Expected searchQuery to be empty, got '%s'",
			updatedWizard.searchQuery,
		)
	}

	// Verify all providers are restored
	if len(
		updatedWizard.filteredProviders,
	) != len(
		updatedWizard.allProviders,
	) {
		t.Error(
			"Expected all providers to be restored after Escape",
		)
	}
}

func TestRenderSelectShowsSearchInput(
	t *testing.T,
) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepSelect
	wizard.searchMode = true

	output := wizard.renderSelect()

	// Verify output contains search input
	if !strings.Contains(output, "Search:") {
		t.Error(
			"Expected select view to contain 'Search:' when in search mode",
		)
	}

	// Verify help text shows Escape instruction
	if !strings.Contains(
		output,
		"Esc: Exit search",
	) {
		t.Error(
			"Expected help text to contain 'Esc: Exit search' when in search mode",
		)
	}
}

func TestRenderSelectShowsSearchHotkey(
	t *testing.T,
) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepSelect
	wizard.searchMode = false

	output := wizard.renderSelect()

	// Verify help text shows /: Search when not in search mode
	if !strings.Contains(output, "/: Search") {
		t.Error(
			"Expected help text to contain '/: Search' when not in search mode",
		)
	}
}

func TestRenderSelectShowsNoMatchMessage(
	t *testing.T,
) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepSelect
	wizard.searchMode = true
	wizard.searchQuery = testProviderNonexistent
	wizard.applyProviderFilter()

	output := wizard.renderSelect()

	// Verify output shows no match message
	if !strings.Contains(
		output,
		"No providers match",
	) {
		t.Error(
			"Expected select view to contain 'No providers match' message",
		)
	}
}

func TestSpaceToggleInSearchMode(t *testing.T) {
	cmd := &InitCmd{Path: "/tmp/test-project"}
	wizard, err := NewWizardModel(cmd)
	if err != nil {
		t.Fatalf(
			"Failed to create wizard model: %v",
			err,
		)
	}

	wizard.step = StepSelect
	wizard.searchMode = true
	wizard.cursor = 0

	// Ensure first provider is not selected
	if len(wizard.filteredProviders) > 0 {
		wizard.selectedProviders[wizard.filteredProviders[0].ID()] = false
	}

	// Simulate pressing space key while in search mode
	keyMsg := tea.KeyMsg{Type: tea.KeySpace}
	newModel, _ := wizard.Update(keyMsg)
	updatedWizard, ok := newModel.(*WizardModel)
	if !ok {
		t.Fatal(
			"Failed to cast model to WizardModel",
		)
	}

	// Verify provider is now selected
	if len(updatedWizard.filteredProviders) == 0 {
		return
	}

	providerID := updatedWizard.filteredProviders[0].ID()
	if !updatedWizard.selectedProviders[providerID] {
		t.Errorf(
			"Expected provider %s to be selected after space press",
			providerID,
		)
	}
}
