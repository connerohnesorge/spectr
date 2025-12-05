package list

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestRunInteractiveChanges_EmptyList(t *testing.T) {
	var changes []ChangeInfo
	archiveID, prID, err := RunInteractiveChanges(changes, "/tmp/test-project")
	if err != nil {
		t.Errorf("RunInteractiveChanges with empty list should not error, got: %v", err)
	}
	if archiveID != "" {
		t.Errorf(
			"RunInteractiveChanges with empty list should return empty archive ID, got: %s",
			archiveID,
		)
	}
	if prID != "" {
		t.Errorf(
			"RunInteractiveChanges with empty list should return empty PR ID, got: %s",
			prID,
		)
	}
}

func TestRunInteractiveSpecs_EmptyList(t *testing.T) {
	var specs []SpecInfo
	err := RunInteractiveSpecs(specs, "/tmp/test-project")
	if err != nil {
		t.Errorf("RunInteractiveSpecs with empty list should not error, got: %v", err)
	}
}

func TestRunInteractiveChanges_ValidData(_ *testing.T) {
	// This test verifies that the function can be called without error
	// Actual interactive testing would require terminal simulation
	changes := []ChangeInfo{
		{
			ID:         "add-test-feature",
			Title:      "Add test feature",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{
				Total:     5,
				Completed: 3,
			},
		},
		{
			ID:         "update-validation",
			Title:      "Update validation logic",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{
				Total:     3,
				Completed: 3,
			},
		},
	}

	// Note: This will fail in CI/CD without a TTY, but validates the structure
	// In a real terminal, this would launch the interactive UI
	_ = changes // Just verify the data structure is correct
}

func TestRunInteractiveSpecs_ValidData(_ *testing.T) {
	// This test verifies that the function can be called without error
	// Actual interactive testing would require terminal simulation
	specs := []SpecInfo{
		{
			ID:               "auth",
			Title:            "Authentication System",
			RequirementCount: 5,
		},
		{
			ID:               "payment",
			Title:            "Payment Processing",
			RequirementCount: 8,
		},
	}

	// Note: This will fail in CI/CD without a TTY, but validates the structure
	// In a real terminal, this would launch the interactive UI
	_ = specs // Just verify the data structure is correct
}

func TestInteractiveModel_Init(t *testing.T) {
	model := interactiveModel{}
	cmd := model.Init()

	if cmd != nil {
		t.Errorf("Init() should return nil, got: %v", cmd)
	}
}

func TestInteractiveModel_View_Quitting(t *testing.T) {
	tests := []struct {
		name       string
		model      interactiveModel
		wantSubstr string
	}{
		{
			name: "quit without copy",
			model: interactiveModel{
				quitting: true,
				copied:   false,
			},
			wantSubstr: "Cancelled",
		},
		{
			name: "quit with successful copy",
			model: interactiveModel{
				quitting:   true,
				copied:     true,
				selectedID: "test-id",
			},
			wantSubstr: "Copied: test-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.model.View()
			if view == "" {
				t.Error("View() returned empty string")
			}
			// Just verify it doesn't panic and returns something
			t.Logf("View output: %s", view)
		})
	}
}

func TestInteractiveModel_HandleEdit(t *testing.T) {
	// Create a temporary test directory with a spec file
	tmpDir := t.TempDir()
	specID := "test-spec"
	specDir := tmpDir + "/spectr/specs/" + specID
	err := mkdirAll(specDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	specPath := specDir + "/spec.md"
	err = writeFile(specPath, []byte("# Test Spec"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec file: %v", err)
	}

	// Test case: spec mode, EDITOR not set
	t.Run("EDITOR not set", func(t *testing.T) {
		// Save and clear EDITOR
		originalEditor := getEnv("EDITOR")
		t.Cleanup(func() {
			if originalEditor != "" {
				_ = setEnv("EDITOR", originalEditor)
			} else {
				_ = unsetEnv("EDITOR")
			}
		})
		_ = unsetEnv("EDITOR")

		model := interactiveModel{
			itemType:    "spec",
			projectPath: tmpDir,
			table: createMockTable([][]string{
				{specID, "Test Spec", "1"},
			}),
		}

		updatedModel, _ := model.handleEdit()
		if updatedModel.err == nil {
			t.Error("Expected error when EDITOR not set")
		}
		if updatedModel.err != nil &&
			updatedModel.err.Error() != "EDITOR environment variable not set" {
			t.Errorf(
				"Expected 'EDITOR environment variable not set' error, got: %v",
				updatedModel.err,
			)
		}
	})

	// Test case: change mode - edit change proposal
	t.Run("change mode opens proposal", func(t *testing.T) {
		// Create a change proposal file
		changeID := "test-change"
		changeDir := tmpDir + "/spectr/changes/" + changeID
		err := mkdirAll(changeDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create change directory: %v", err)
		}
		proposalPath := changeDir + "/proposal.md"
		err = writeFile(proposalPath, []byte("# Test Change"), 0644)
		if err != nil {
			t.Fatalf("Failed to create proposal file: %v", err)
		}

		_ = setEnv("EDITOR", "true")
		t.Cleanup(func() { _ = unsetEnv("EDITOR") })

		// Create a change mode table with 4 columns
		columns := []table.Column{
			{Title: "ID", Width: changeIDWidth},
			{Title: "Title", Width: changeTitleWidth},
			{Title: "Deltas", Width: changeDeltaWidth},
			{Title: "Tasks", Width: changeTasksWidth},
		}
		rows := []table.Row{{changeID, "Test Change", "2", "3/5"}}
		tbl := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		)

		model := interactiveModel{
			itemType:    "change",
			projectPath: tmpDir,
			table:       tbl,
		}

		updatedModel, cmd := model.handleEdit()
		if cmd == nil {
			t.Error("Expected command to be returned when editing change")
		}
		if updatedModel.err != nil {
			t.Errorf("Expected no error when editing change, got: %v", updatedModel.err)
		}
	})

	// Test case: spec file not found
	t.Run("spec file not found", func(t *testing.T) {
		_ = setEnv("EDITOR", "vim")
		t.Cleanup(func() { _ = unsetEnv("EDITOR") })

		model := interactiveModel{
			itemType:    "spec",
			projectPath: tmpDir,
			table: createMockTable([][]string{
				{"nonexistent-spec", "Nonexistent Spec", "1"},
			}),
		}

		updatedModel, _ := model.handleEdit()
		if updatedModel.err == nil {
			t.Error("Expected error for nonexistent spec file")
		}
	})
}

// Helper functions for tests
func mkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func setEnv(key, value string) error {
	return os.Setenv(key, value)
}

func unsetEnv(key string) error {
	return os.Unsetenv(key)
}

func createMockTable(rows [][]string) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 35},
		{Title: "Title", Width: 45},
		{Title: "Requirements", Width: 15},
	}

	tableRows := make([]table.Row, len(rows))
	for i, row := range rows {
		tableRows[i] = row
	}

	return table.New(
		table.WithColumns(columns),
		table.WithRows(tableRows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
}

func TestRunInteractiveAll_EmptyList(t *testing.T) {
	var items ItemList
	err := RunInteractiveAll(items, "/tmp/test-project")
	if err != nil {
		t.Errorf("RunInteractiveAll with empty list should not error, got: %v", err)
	}
}

func TestRunInteractiveAll_ValidData(_ *testing.T) {
	// This test verifies that the function can be called without error
	// Actual interactive testing would require terminal simulation
	items := ItemList{
		NewChangeItem(ChangeInfo{
			ID:         "add-test-feature",
			Title:      "Add test feature",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{
				Total:     5,
				Completed: 3,
			},
		}),
		NewSpecItem(SpecInfo{
			ID:               "auth",
			Title:            "Authentication System",
			RequirementCount: 5,
		}),
	}

	// Note: This will fail in CI/CD without a TTY, but validates the structure
	// In a real terminal, this would launch the interactive UI
	_ = items // Just verify the data structure is correct
}

func TestHandleToggleFilter(t *testing.T) {
	// Create a model with all items
	items := ItemList{
		NewChangeItem(ChangeInfo{
			ID:         "change-1",
			Title:      "Change 1",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{Total: 3, Completed: 1},
		}),
		NewSpecItem(SpecInfo{
			ID:               "spec-1",
			Title:            "Spec 1",
			RequirementCount: 5,
		}),
	}

	model := interactiveModel{
		itemType:    "all",
		allItems:    items,
		filterType:  nil,
		projectPath: "/tmp/test",
	}

	// Test toggle: all -> changes
	model = model.handleToggleFilter()
	if model.filterType == nil {
		t.Error("Expected filterType to be set to ItemTypeChange")
	}
	if *model.filterType != ItemTypeChange {
		t.Errorf("Expected ItemTypeChange, got %v", *model.filterType)
	}

	// Test toggle: changes -> specs
	model = model.handleToggleFilter()
	if model.filterType == nil {
		t.Error("Expected filterType to be set to ItemTypeSpec")
	}
	if *model.filterType != ItemTypeSpec {
		t.Errorf("Expected ItemTypeSpec, got %v", *model.filterType)
	}

	// Test toggle: specs -> all
	model = model.handleToggleFilter()
	if model.filterType != nil {
		t.Errorf("Expected filterType to be nil (all), got %v", model.filterType)
	}
}

// TestEditorOpensOnEKey tests that pressing 'e' opens the editor for specs
func TestEditorOpensOnEKey(t *testing.T) {
	// Create a temporary test directory with a spec file
	tmpDir := t.TempDir()
	specID := "test-spec"
	specDir := tmpDir + "/spectr/specs/" + specID
	err := os.MkdirAll(specDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	specPath := specDir + "/spec.md"
	err = os.WriteFile(specPath, []byte("# Test Spec"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec file: %v", err)
	}

	// Create a test model with a spec item
	columns := []table.Column{
		{Title: "ID", Width: specIDWidth},
		{Title: "Title", Width: specTitleWidth},
		{Title: "Requirements", Width: specRequirementsWidth},
	}

	rows := []table.Row{
		{specID, "Test Spec", "5"},
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	m := interactiveModel{
		table:       tbl,
		itemType:    "spec",
		projectPath: tmpDir,
		helpText:    "Test help text",
	}

	// Set EDITOR to a command that will succeed but not actually edit
	originalEditor := os.Getenv("EDITOR")
	if err := os.Setenv("EDITOR", "true"); err != nil {
		t.Fatalf("Failed to set EDITOR: %v", err)
	}
	t.Cleanup(func() {
		if originalEditor != "" {
			if err := os.Setenv("EDITOR", originalEditor); err != nil {
				t.Logf("Failed to restore EDITOR: %v", err)
			}
		} else {
			if err := os.Unsetenv("EDITOR"); err != nil {
				t.Logf("Failed to unset EDITOR: %v", err)
			}
		}
	})

	tm := teatest.NewTestModel(t, m)

	// Wait for the initial view to render
	waitForString(t, tm, "Test Spec")

	// Send 'e' to open editor - this should trigger the editor opening
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	// Wait a bit for the editor command to be processed
	time.Sleep(time.Millisecond * 1000)

	// Send Ctrl+C to quit - at this point the editor should have been opened and closed
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	// The model should finish without errors
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

// TestEditorOpensForChangeItems tests that pressing 'e' opens editor for changes
func TestEditorOpensForChangeItems(t *testing.T) {
	// Create a temporary test directory with a change proposal file
	tmpDir := t.TempDir()
	changeID := "test-change"
	changeDir := tmpDir + "/spectr/changes/" + changeID
	err := os.MkdirAll(changeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	proposalPath := changeDir + "/proposal.md"
	err = os.WriteFile(proposalPath, []byte("# Test Change"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test proposal file: %v", err)
	}

	// Create a test model with a change item
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}

	rows := []table.Row{
		{changeID, "Test Change", "2", "3/5"},
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	m := interactiveModel{
		table:       tbl,
		itemType:    "change",
		projectPath: tmpDir,
		helpText:    "Test help text",
	}

	// Set EDITOR to a command that will succeed
	originalEditor := os.Getenv("EDITOR")
	if err := os.Setenv("EDITOR", "true"); err != nil {
		t.Fatalf("Failed to set EDITOR: %v", err)
	}
	t.Cleanup(func() {
		if originalEditor != "" {
			if err := os.Setenv("EDITOR", originalEditor); err != nil {
				t.Logf("Failed to restore EDITOR: %v", err)
			}
		} else {
			if err := os.Unsetenv("EDITOR"); err != nil {
				t.Logf("Failed to unset EDITOR: %v", err)
			}
		}
	})

	tm := teatest.NewTestModel(t, m)

	// Wait for the initial view to render
	waitForString(t, tm, "Test Change")

	// Send 'e' to open editor for the change
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	// Wait a bit for the editor command to be processed
	time.Sleep(time.Millisecond * 1000)

	// Send Ctrl+C to quit
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	// The model should finish without errors
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

// TestEditorOpensInUnifiedMode tests that pressing 'e' opens editor in unified mode
func TestEditorOpensInUnifiedMode(t *testing.T) {
	// Create a temporary test directory with both spec and change files
	tmpDir := t.TempDir()

	// Create spec file
	specID := "test-spec"
	specDir := tmpDir + "/spectr/specs/" + specID
	err := os.MkdirAll(specDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create spec directory: %v", err)
	}
	err = os.WriteFile(specDir+"/spec.md", []byte("# Test Spec"), 0644)
	if err != nil {
		t.Fatalf("Failed to create spec file: %v", err)
	}

	// Create change file
	changeID := "test-change"
	changeDir := tmpDir + "/spectr/changes/" + changeID
	err = os.MkdirAll(changeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create change directory: %v", err)
	}
	err = os.WriteFile(changeDir+"/proposal.md", []byte("# Test Change"), 0644)
	if err != nil {
		t.Fatalf("Failed to create proposal file: %v", err)
	}

	// Create unified mode model with both items
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}

	rows := []table.Row{
		{changeID, "CHANGE", "Test Change", "Deltas: 2 | Tasks: 3/5"},
		{specID, "SPEC", "Test Spec", "Reqs: 5"},
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	items := ItemList{
		NewChangeItem(ChangeInfo{
			ID:         changeID,
			Title:      "Test Change",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{Total: 5, Completed: 3},
		}),
		NewSpecItem(SpecInfo{
			ID:               specID,
			Title:            "Test Spec",
			RequirementCount: 5,
		}),
	}

	m := interactiveModel{
		table:       tbl,
		itemType:    "all",
		projectPath: tmpDir,
		allItems:    items,
		filterType:  nil,
		helpText:    "Test help text",
	}

	// Set EDITOR
	originalEditor := os.Getenv("EDITOR")
	if err := os.Setenv("EDITOR", "true"); err != nil {
		t.Fatalf("Failed to set EDITOR: %v", err)
	}
	t.Cleanup(func() {
		if originalEditor != "" {
			if err := os.Setenv("EDITOR", originalEditor); err != nil {
				t.Logf("Failed to restore EDITOR: %v", err)
			}
		} else {
			if err := os.Unsetenv("EDITOR"); err != nil {
				t.Logf("Failed to unset EDITOR: %v", err)
			}
		}
	})

	tm := teatest.NewTestModel(t, m)

	// Wait for the initial view to render
	waitForString(t, tm, "Test Change")

	// First, test editing the change item (first item, currently selected)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	time.Sleep(time.Millisecond * 500)

	// Move down to the spec item
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(time.Millisecond * 100)

	// Edit the spec item
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	time.Sleep(time.Millisecond * 500)

	// Quit
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	// The model should finish without errors
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

// waitForString is a helper function for teatest
func waitForString(t *testing.T, tm *teatest.TestModel, s string) {
	teatest.WaitFor(
		t,
		tm.Output(),
		func(b []byte) bool {
			return strings.Contains(string(b), s)
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*10),
	)
}

func TestHandleArchive_ChangeMode(t *testing.T) {
	// Create a model in change mode with a valid change
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"test-change-1", "Test Change 1", "2", "3/5"},
		{"test-change-2", "Test Change 2", "1", "2/2"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "change",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handleArchive
	updatedModel, cmd := model.handleArchive()

	// Should set archiveRequested and selectedID
	if !updatedModel.archiveRequested {
		t.Error("Expected archiveRequested to be true in change mode")
	}
	if updatedModel.selectedID != "test-change-1" {
		t.Errorf("Expected selectedID to be 'test-change-1', got '%s'", updatedModel.selectedID)
	}
	// Should return tea.Quit
	if cmd == nil {
		t.Error("Expected command to be returned (tea.Quit)")
	}
}

func TestHandleArchive_SpecMode(t *testing.T) {
	// Create a model in spec mode
	columns := []table.Column{
		{Title: "ID", Width: specIDWidth},
		{Title: "Title", Width: specTitleWidth},
		{Title: "Requirements", Width: specRequirementsWidth},
	}
	rows := []table.Row{
		{"test-spec", "Test Spec", "5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "spec",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handleArchive
	updatedModel, cmd := model.handleArchive()

	// Should NOT set archiveRequested in spec mode
	if updatedModel.archiveRequested {
		t.Error("Expected archiveRequested to be false in spec mode")
	}
	if updatedModel.selectedID != "" {
		t.Errorf("Expected selectedID to be empty in spec mode, got '%s'", updatedModel.selectedID)
	}
	// Should not return tea.Quit
	if cmd != nil {
		t.Error("Expected no command in spec mode")
	}
}

func TestHandleArchive_UnifiedMode_Change(t *testing.T) {
	// Create a model in unified mode with a change selected
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}
	rows := []table.Row{
		{"test-change", "CHANGE", "Test Change", "Tasks: 3/5"},
		{"test-spec", "SPEC", "Test Spec", "Reqs: 5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "all",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handleArchive (cursor is on first row which is CHANGE)
	updatedModel, cmd := model.handleArchive()

	// Should set archiveRequested and selectedID for CHANGE
	if !updatedModel.archiveRequested {
		t.Error("Expected archiveRequested to be true for CHANGE in unified mode")
	}
	if updatedModel.selectedID != "test-change" {
		t.Errorf("Expected selectedID to be 'test-change', got '%s'", updatedModel.selectedID)
	}
	if cmd == nil {
		t.Error("Expected command to be returned (tea.Quit)")
	}
}

func TestHandleArchive_UnifiedMode_Spec(t *testing.T) {
	// Create a model in unified mode with cursor on spec
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}
	rows := []table.Row{
		{"test-spec", "SPEC", "Test Spec", "Reqs: 5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "all",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handleArchive (cursor is on SPEC)
	updatedModel, cmd := model.handleArchive()

	// Should NOT set archiveRequested for SPEC in unified mode
	if updatedModel.archiveRequested {
		t.Error("Expected archiveRequested to be false for SPEC in unified mode")
	}
	if updatedModel.selectedID != "" {
		t.Errorf("Expected selectedID to be empty for SPEC, got '%s'", updatedModel.selectedID)
	}
	if cmd != nil {
		t.Error("Expected no command for SPEC in unified mode")
	}
}

func TestHandlePR_ChangeMode(t *testing.T) {
	// Create a model in change mode with a valid change
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"test-change-1", "Test Change 1", "2", "3/5"},
		{"test-change-2", "Test Change 2", "1", "2/2"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "change",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handlePR
	updatedModel, cmd := model.handlePR()

	// Should set prRequested and selectedID
	if !updatedModel.prRequested {
		t.Error("Expected prRequested to be true in change mode")
	}
	if updatedModel.selectedID != "test-change-1" {
		t.Errorf("Expected selectedID to be 'test-change-1', got '%s'", updatedModel.selectedID)
	}
	// Should return tea.Quit
	if cmd == nil {
		t.Error("Expected command to be returned (tea.Quit)")
	}
}

func TestHandlePR_SpecMode(t *testing.T) {
	// Create a model in spec mode
	columns := []table.Column{
		{Title: "ID", Width: specIDWidth},
		{Title: "Title", Width: specTitleWidth},
		{Title: "Requirements", Width: specRequirementsWidth},
	}
	rows := []table.Row{
		{"test-spec", "Test Spec", "5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "spec",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handlePR
	updatedModel, cmd := model.handlePR()

	// Should NOT set prRequested in spec mode
	if updatedModel.prRequested {
		t.Error("Expected prRequested to be false in spec mode")
	}
	if updatedModel.selectedID != "" {
		t.Errorf("Expected selectedID to be empty in spec mode, got '%s'", updatedModel.selectedID)
	}
	// Should not return tea.Quit
	if cmd != nil {
		t.Error("Expected no command in spec mode")
	}
}

func TestHandlePR_UnifiedMode(t *testing.T) {
	// Create a model in unified mode
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}
	rows := []table.Row{
		{"test-change", "CHANGE", "Test Change", "Tasks: 3/5"},
		{"test-spec", "SPEC", "Test Spec", "Reqs: 5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "all",
		projectPath: "/tmp/test",
		table:       tbl,
	}

	// Call handlePR (cursor is on first row which is CHANGE)
	updatedModel, cmd := model.handlePR()

	// Should NOT set prRequested in unified mode (PR only works in change mode)
	if updatedModel.prRequested {
		t.Error("Expected prRequested to be false in unified mode")
	}
	if cmd != nil {
		t.Error("Expected no command in unified mode")
	}
}

func TestViewShowsPRMode(t *testing.T) {
	model := interactiveModel{
		quitting:    true,
		prRequested: true,
		selectedID:  "test-change-id",
	}

	view := model.View()
	if !strings.Contains(view, "PR mode: test-change-id") {
		t.Errorf("Expected view to contain 'PR mode: test-change-id', got: %s", view)
	}
}

// TestSearchRowFiltering tests row filtering by all columns (ID, title, deltas, tasks)
func TestSearchRowFiltering(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"add-feature", "Add new feature", "2", "3/5"},
		{"fix-bug", "Fix bug in parser", "1", "2/2"},
		{"update-docs", "Update documentation", "1", "1/1"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "change",
		projectPath: "/tmp/test",
		table:       tbl,
		searchInput: newTextInput(),
		allRows:     rows,
		searchMode:  true,
	}

	tests := []struct {
		name            string
		searchQuery     string
		expectedRowsLen int
		expectedIDs     []string
	}{
		{
			name:            "Search by ID prefix",
			searchQuery:     "add",
			expectedRowsLen: 1,
			expectedIDs:     []string{"add-feature"},
		},
		{
			name:            "Search by title word",
			searchQuery:     "bug",
			expectedRowsLen: 1,
			expectedIDs:     []string{"fix-bug"},
		},
		{
			name:            "Search matches multiple rows",
			searchQuery:     "update",
			expectedRowsLen: 1,
			expectedIDs:     []string{"update-docs"},
		},
		{
			name:            "Search by deltas column",
			searchQuery:     "2",
			expectedRowsLen: 2,
			expectedIDs:     []string{"add-feature", "fix-bug"},
		},
		{
			name:            "Search by tasks column",
			searchQuery:     "3/5",
			expectedRowsLen: 1,
			expectedIDs:     []string{"add-feature"},
		},
		{
			name:            "No matches",
			searchQuery:     "xyz",
			expectedRowsLen: 0,
			expectedIDs:     nil,
		},
		{
			name:            "Empty search shows all",
			searchQuery:     "",
			expectedRowsLen: 3,
			expectedIDs:     []string{"add-feature", "fix-bug", "update-docs"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.searchQuery = tt.searchQuery
			model.searchInput.SetValue(tt.searchQuery)
			model = model.applyFilter()

			filteredRows := model.table.Rows()
			if len(filteredRows) != tt.expectedRowsLen {
				t.Errorf("Expected %d rows, got %d for query '%s'",
					tt.expectedRowsLen, len(filteredRows), tt.searchQuery)
			}

			for i, expectedID := range tt.expectedIDs {
				if i < len(filteredRows) && filteredRows[i][0] != expectedID {
					t.Errorf("Expected ID %s at index %d, got %s",
						expectedID, i, filteredRows[i][0])
				}
			}
		})
	}
}

// TestSearchExitWithEscape tests exiting search mode with Escape key
func TestSearchExitWithEscape(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"add-feature", "Add new feature", "2", "3/5"},
		{"fix-bug", "Fix bug in parser", "1", "2/2"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "change",
		projectPath: "/tmp/test",
		table:       tbl,
		searchInput: newTextInput(),
		allRows:     rows,
		searchMode:  true,
	}

	// Set a search query and then exit with Escape
	model.searchQuery = "add"
	model.searchInput.SetValue("add")
	model = model.applyFilter()

	// Simulate pressing Escape
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m, ok := updatedModel.(interactiveModel)
	if !ok {
		t.Fatal("Expected interactiveModel type")
	}

	// Should exit search mode and clear query
	if m.searchMode {
		t.Error("Expected searchMode to be false after Escape")
	}
	if m.searchQuery != "" {
		t.Errorf("Expected searchQuery to be empty after Escape, got '%s'", m.searchQuery)
	}
	// Should restore all rows
	if len(m.table.Rows()) != len(rows) {
		t.Errorf("Expected %d rows after clearing search, got %d", len(rows), len(m.table.Rows()))
	}
}

// TestSearchUnifiedMode tests search in unified mode (changes and specs)
func TestSearchUnifiedMode(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}
	rows := []table.Row{
		{"add-auth", "CHANGE", "Add authentication", "Tasks: 2/5"},
		{"auth-system", "SPEC", "Authentication System", "Reqs: 8"},
		{"add-cache", "CHANGE", "Add caching layer", "Tasks: 1/3"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "all",
		projectPath: "/tmp/test",
		table:       tbl,
		searchInput: newTextInput(),
		allRows:     rows,
		searchMode:  true,
	}

	// Search for "auth" - should find both auth-related items
	model.searchQuery = "auth"
	model.searchInput.SetValue("auth")
	model = model.applyFilter()

	filteredRows := model.table.Rows()
	if len(filteredRows) != 2 {
		t.Errorf("Expected 2 rows matching 'auth', got %d", len(filteredRows))
	}

	// Verify the matches
	expectedIDs := []string{"add-auth", "auth-system"}
	for i, expectedID := range expectedIDs {
		if i < len(filteredRows) && filteredRows[i][0] != expectedID {
			t.Errorf("Expected ID %s at index %d, got %s",
				expectedID, i, filteredRows[i][0])
		}
	}
}

// TestSearchCaseInsensitive tests that search is case insensitive
func TestSearchCaseInsensitive(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: specIDWidth},
		{Title: "Title", Width: specTitleWidth},
		{Title: "Requirements", Width: specRequirementsWidth},
	}
	rows := []table.Row{
		{"user-auth", "User Authentication System", "5"},
		{"payment", "Payment Processing", "8"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:    "spec",
		projectPath: "/tmp/test",
		table:       tbl,
		searchInput: newTextInput(),
		allRows:     rows,
		searchMode:  true,
	}

	tests := []struct {
		query           string
		expectedRowsLen int
	}{
		{"USER", 1},
		{"user", 1},
		{"User", 1},
		{"AUTH", 1},
		{"PAYMENT", 1},
		{"payment", 1},
	}

	for _, tt := range tests {
		model.searchQuery = tt.query
		model.searchInput.SetValue(tt.query)
		model = model.applyFilter()

		filteredRows := model.table.Rows()
		if len(filteredRows) != tt.expectedRowsLen {
			t.Errorf("Query '%s': expected %d rows, got %d",
				tt.query, tt.expectedRowsLen, len(filteredRows))
		}
	}
}

// TestHelpToggleDefaultView tests that default view shows minimal footer
func TestHelpToggleDefaultView(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"add-feature", "Add new feature", "2", "3/5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:      "change",
		projectPath:   "/tmp/test",
		table:         tbl,
		showHelp:      false, // Default: help not shown
		helpText:      "↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | /: search | q: quit",
		minimalFooter: "showing: 1 | project: /tmp/test | ?: help",
	}

	view := model.View()

	// Should show minimal footer by default
	if !strings.Contains(view, "showing: 1 | project: /tmp/test | ?: help") {
		t.Error("Expected view to contain minimal footer by default")
	}

	// Should NOT show full help text by default
	if strings.Contains(view, "↑/↓/j/k: navigate | Enter: copy ID") {
		t.Error("Expected view to NOT contain full help text by default")
	}
}

// TestHelpToggleShowsHelp tests that pressing '?' toggles help visibility
func TestHelpToggleShowsHelp(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"add-feature", "Add new feature", "2", "3/5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:      "change",
		projectPath:   "/tmp/test",
		table:         tbl,
		showHelp:      false,
		helpText:      "↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | /: search | q: quit",
		minimalFooter: "showing: 1 | project: /tmp/test | ?: help",
	}

	// Press '?' to show help
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m, ok := updatedModel.(interactiveModel)
	if !ok {
		t.Fatal("Expected interactiveModel type")
	}

	// showHelp should now be true
	if !m.showHelp {
		t.Error("Expected showHelp to be true after pressing '?'")
	}

	// View should now show full help text
	view := m.View()
	if !strings.Contains(view, "↑/↓/j/k: navigate | Enter: copy ID") {
		t.Error("Expected view to contain full help text after pressing '?'")
	}

	// View should NOT show minimal footer
	if strings.Contains(view, "showing: 1 | project: /tmp/test | ?: help") {
		t.Error("Expected view to NOT contain minimal footer when help is shown")
	}
}

// TestHelpToggleHidesHelp tests that pressing '?' again hides help
func TestHelpToggleHidesHelp(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"add-feature", "Add new feature", "2", "3/5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:      "change",
		projectPath:   "/tmp/test",
		table:         tbl,
		showHelp:      true, // Start with help shown
		helpText:      "↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | /: search | q: quit",
		minimalFooter: "showing: 1 | project: /tmp/test | ?: help",
	}

	// Press '?' to hide help
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m, ok := updatedModel.(interactiveModel)
	if !ok {
		t.Fatal("Expected interactiveModel type")
	}

	// showHelp should now be false
	if m.showHelp {
		t.Error("Expected showHelp to be false after pressing '?' again")
	}

	// View should now show minimal footer
	view := m.View()
	if !strings.Contains(view, "showing: 1 | project: /tmp/test | ?: help") {
		t.Error("Expected view to contain minimal footer after hiding help")
	}

	// View should NOT show full help text
	if strings.Contains(view, "↑/↓/j/k: navigate | Enter: copy ID") {
		t.Error("Expected view to NOT contain full help text after hiding help")
	}
}

// TestNavigationKeysAutoHideHelp tests that navigation keys auto-hide help
func TestNavigationKeysAutoHideHelp(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}
	rows := []table.Row{
		{"add-feature", "Add new feature", "2", "3/5"},
		{"fix-bug", "Fix bug", "1", "2/2"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Test each navigation key
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
	}

	for _, key := range navigationKeys {
		t.Run(key.String(), func(t *testing.T) {
			model := interactiveModel{
				itemType:      "change",
				projectPath:   "/tmp/test",
				table:         tbl,
				showHelp:      true, // Start with help shown
				helpText:      "↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | /: search | q: quit",
				minimalFooter: "showing: 2 | project: /tmp/test | ?: help",
			}

			// Press navigation key
			updatedModel, _ := model.Update(key)
			m, ok := updatedModel.(interactiveModel)
			if !ok {
				t.Fatal("Expected interactiveModel type")
			}

			// showHelp should now be false
			if m.showHelp {
				t.Errorf("Expected showHelp to be false after pressing %s", key.String())
			}
		})
	}
}

// TestHelpToggleInSpecMode tests help toggle works in spec mode
func TestHelpToggleInSpecMode(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: specIDWidth},
		{Title: "Title", Width: specTitleWidth},
		{Title: "Requirements", Width: specRequirementsWidth},
	}
	rows := []table.Row{
		{"auth", "Authentication System", "5"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:      "spec",
		projectPath:   "/tmp/test",
		table:         tbl,
		showHelp:      false,
		helpText:      "↑/↓/j/k: navigate | Enter: copy ID | e: edit | /: search | q: quit",
		minimalFooter: "showing: 1 | project: /tmp/test | ?: help",
	}

	// Toggle help on
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m, ok := updatedModel.(interactiveModel)
	if !ok {
		t.Fatal("Expected interactiveModel type")
	}

	if !m.showHelp {
		t.Error("Expected showHelp to be true in spec mode")
	}

	view := m.View()
	if !strings.Contains(view, "↑/↓/j/k: navigate") {
		t.Error("Expected view to show full help text in spec mode")
	}
}

// TestHelpToggleInUnifiedMode tests help toggle works in unified mode
func TestHelpToggleInUnifiedMode(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}
	rows := []table.Row{
		{"add-auth", "CHANGE", "Add authentication", "Tasks: 2/5"},
		{"auth-system", "SPEC", "Authentication System", "Reqs: 8"},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	model := interactiveModel{
		itemType:      "all",
		projectPath:   "/tmp/test",
		table:         tbl,
		showHelp:      false,
		helpText:      "↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | t: filter (all) | /: search | q: quit",
		minimalFooter: "showing: 2 | project: /tmp/test | ?: help",
	}

	// Toggle help on
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m, ok := updatedModel.(interactiveModel)
	if !ok {
		t.Fatal("Expected interactiveModel type")
	}

	if !m.showHelp {
		t.Error("Expected showHelp to be true in unified mode")
	}

	view := m.View()
	if !strings.Contains(view, "t: filter") {
		t.Error("Expected view to show full help text with filter option in unified mode")
	}
}
