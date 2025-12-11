//nolint:revive // test file
package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNewTablePicker(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Name", Width: 30},
	}
	rows := []table.Row{
		{"1", "First"},
		{"2", "Second"},
	}

	config := TableConfig{
		Columns:     columns,
		Rows:        rows,
		Height:      5,
		ProjectPath: "/test/path",
	}

	picker := NewTablePicker(config)

	if picker == nil {
		t.Fatal("NewTablePicker returned nil")
	}

	if picker.projectPath != "/test/path" {
		t.Errorf(
			"projectPath = %q, want %q",
			picker.projectPath,
			"/test/path",
		)
	}

	// Should have default q and ctrl+c actions
	if _, exists := picker.actions["q"]; !exists {
		t.Error("Expected default 'q' action")
	}
	if _, exists := picker.actions["ctrl+c"]; !exists {
		t.Error(
			"Expected default 'ctrl+c' action",
		)
	}
}

func TestTablePicker_WithAction(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: 20},
	}
	rows := []table.Row{{"1"}}
	config := TableConfig{
		Columns: columns,
		Rows:    rows,
	}

	called := false
	picker := NewTablePicker(config).
		WithAction("e", "edit", func(_ table.Row) (tea.Cmd, *ActionResult) {
			called = true

			return nil, nil
		})

	if _, exists := picker.actions["e"]; !exists {
		t.Error(
			"Expected 'e' action to be registered",
		)
	}

	// Execute the action
	action := picker.actions["e"]
	action.Handler(nil)

	if !called {
		t.Error(
			"Expected action handler to be called",
		)
	}
}

func TestTablePicker_generateHelpText(
	t *testing.T,
) {
	columns := []table.Column{
		{Title: "ID", Width: 20},
	}
	rows := []table.Row{{"1"}, {"2"}, {"3"}}

	config := TableConfig{
		Columns:     columns,
		Rows:        rows,
		ProjectPath: "/test",
	}

	picker := NewTablePicker(config).
		WithAction("e", "edit", func(_ table.Row) (tea.Cmd, *ActionResult) { return nil, nil }).
		WithAction("a", "archive", func(_ table.Row) (tea.Cmd, *ActionResult) { return nil, nil })

	helpText := picker.helpText

	// Check that help text contains expected elements
	if !strings.Contains(helpText, "navigate") {
		t.Error(
			"Help text should contain 'navigate'",
		)
	}
	if !strings.Contains(helpText, "e: edit") {
		t.Error(
			"Help text should contain 'e: edit'",
		)
	}
	if !strings.Contains(helpText, "a: archive") {
		t.Error(
			"Help text should contain 'a: archive'",
		)
	}
	if !strings.Contains(helpText, "q: quit") {
		t.Error(
			"Help text should contain 'q: quit'",
		)
	}
	if !strings.Contains(helpText, "showing: 3") {
		t.Error(
			"Help text should contain 'showing: 3'",
		)
	}
	if !strings.Contains(
		helpText,
		"project: /test",
	) {
		t.Error(
			"Help text should contain 'project: /test'",
		)
	}
}

func TestTablePicker_SetRows(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: 20},
	}
	initialRows := []table.Row{{"1"}}

	config := TableConfig{
		Columns: columns,
		Rows:    initialRows,
	}

	picker := NewTablePicker(config)

	// Verify initial state
	if !strings.Contains(
		picker.helpText,
		"showing: 1",
	) {
		t.Error(
			"Initial help text should show 1 row",
		)
	}

	// Update rows
	newRows := []table.Row{
		{"1"},
		{"2"},
		{"3"},
		{"4"},
	}
	picker.SetRows(newRows)

	// Verify help text updated
	if !strings.Contains(
		picker.helpText,
		"showing: 4",
	) {
		t.Errorf(
			"Help text should show 4 rows after update, got: %s",
			picker.helpText,
		)
	}
}

func TestTablePicker_ActionResult(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: 20},
	}
	rows := []table.Row{{"test-id"}}

	config := TableConfig{
		Columns: columns,
		Rows:    rows,
		Actions: map[string]Action{
			"enter": {
				Key:         "enter",
				Description: "select",
				Handler: func(row table.Row) (tea.Cmd, *ActionResult) {
					id := ""
					if len(row) > 0 {
						id = row[0]
					}

					return tea.Quit, &ActionResult{
						ID:     id,
						Quit:   true,
						Copied: true,
					}
				},
			},
		},
	}

	picker := NewTablePicker(config)

	// The first call returns the string representation
	model, _ := picker.Update(
		tea.KeyMsg{Type: tea.KeyEnter},
	)
	_ = model.(*TablePicker)

	// Check if result was set (this tests the action was called but won't work
	// because tea.KeyMsg for enter is special)
	// Let's test with a regular key
	picker2 := NewTablePicker(TableConfig{
		Columns: columns,
		Rows:    rows,
		Actions: map[string]Action{
			"x": {
				Key:         "x",
				Description: "test",
				Handler: func(row table.Row) (tea.Cmd, *ActionResult) {
					return tea.Quit, &ActionResult{
						ID:   "test",
						Quit: true,
					}
				},
			},
		},
	})

	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'x'},
	}
	model, _ = picker2.Update(keyMsg)
	updatedPicker := model.(*TablePicker)

	result := updatedPicker.Result()
	if result == nil {
		t.Fatal("Expected result after action")
	}
	if result.ID != "test" {
		t.Errorf(
			"Result ID = %q, want %q",
			result.ID,
			"test",
		)
	}
	if !result.Quit {
		t.Error("Expected Quit to be true")
	}
}

func TestTablePicker_View(t *testing.T) {
	columns := []table.Column{
		{Title: "ID", Width: 20},
	}
	rows := []table.Row{{"1"}}

	config := TableConfig{
		Columns:     columns,
		Rows:        rows,
		ProjectPath: "/test",
	}

	picker := NewTablePicker(config)
	view := picker.View()

	// View should contain table and minimal footer (showHelp=false by default)
	if !strings.Contains(view, "ID") {
		t.Error(
			"View should contain column header",
		)
	}
	if !strings.Contains(view, "?: help") {
		t.Error(
			"View should contain minimal footer with help hint",
		)
	}
	if strings.Contains(view, "navigate") {
		t.Error(
			"View should not contain full help by default",
		)
	}

	// When showHelp is toggled, should show full help
	picker.showHelp = true
	viewWithHelp := picker.View()
	if !strings.Contains(
		viewWithHelp,
		"navigate",
	) {
		t.Error(
			"View should contain full help when showHelp is true",
		)
	}
}

func TestTablePicker_QuitView(t *testing.T) {
	tests := []struct {
		name     string
		result   *ActionResult
		expected string
	}{
		{
			name: "cancelled",
			result: &ActionResult{
				Cancelled: true,
			},
			expected: "Cancelled.\n",
		},
		{
			name: "copied",
			result: &ActionResult{
				ID:     "test-id",
				Copied: true,
			},
			expected: "✓ Copied: test-id\n",
		},
		{
			name: "archive requested",
			result: &ActionResult{
				ID:               "change-1",
				ArchiveRequested: true,
			},
			expected: "Archiving: change-1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns := []table.Column{
				{Title: "ID", Width: 20},
			}
			rows := []table.Row{{"1"}}

			picker := NewTablePicker(
				TableConfig{
					Columns: columns,
					Rows:    rows,
				},
			)
			picker.quitting = true
			picker.result = tt.result

			view := picker.View()
			if view != tt.expected {
				t.Errorf(
					"View() = %q, want %q",
					view,
					tt.expected,
				)
			}
		})
	}
}
