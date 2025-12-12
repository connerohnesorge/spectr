//nolint:revive // test file
package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMenuPicker(t *testing.T) {
	choices := []string{
		"Option 1",
		"Option 2",
		"Option 3",
	}

	config := MenuConfig{
		Title:   "Test Menu",
		Choices: choices,
	}

	menu := NewMenuPicker(config)

	if menu == nil {
		t.Fatal("NewMenuPicker returned nil")
	}

	if menu.title != "Test Menu" {
		t.Errorf(
			"title = %q, want %q",
			menu.title,
			"Test Menu",
		)
	}

	if len(menu.choices) != 3 {
		t.Errorf(
			"choices length = %d, want %d",
			len(menu.choices),
			3,
		)
	}

	if menu.cursor != 0 {
		t.Errorf(
			"initial cursor = %d, want %d",
			menu.cursor,
			0,
		)
	}
}

func TestMenuPicker_Navigation(t *testing.T) {
	config := MenuConfig{
		Title:   "Test",
		Choices: []string{"A", "B", "C"},
	}

	menu := NewMenuPicker(config)

	// Test down navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ := menu.Update(downMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 1 {
		t.Errorf(
			"After down, cursor = %d, want %d",
			menu.cursor,
			1,
		)
	}

	// Test up navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ = menu.Update(upMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 0 {
		t.Errorf(
			"After up, cursor = %d, want %d",
			menu.cursor,
			0,
		)
	}

	// Test vim-style navigation
	jMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'j'},
	}
	model, _ = menu.Update(jMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 1 {
		t.Errorf(
			"After 'j', cursor = %d, want %d",
			menu.cursor,
			1,
		)
	}

	kMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'k'},
	}
	model, _ = menu.Update(kMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 0 {
		t.Errorf(
			"After 'k', cursor = %d, want %d",
			menu.cursor,
			0,
		)
	}
}

func TestMenuPicker_Boundaries(t *testing.T) {
	config := MenuConfig{
		Title:   "Test",
		Choices: []string{"A", "B"},
	}

	menu := NewMenuPicker(config)

	// Try to go up when at top
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ := menu.Update(upMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 0 {
		t.Errorf(
			"Cursor should stay at 0, got %d",
			menu.cursor,
		)
	}

	// Go to bottom
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ = menu.Update(downMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 1 {
		t.Errorf(
			"Cursor should be at 1, got %d",
			menu.cursor,
		)
	}

	// Try to go down past end
	model, _ = menu.Update(downMsg)
	menu = model.(*MenuPicker)

	if menu.cursor != 1 {
		t.Errorf(
			"Cursor should stay at 1, got %d",
			menu.cursor,
		)
	}
}

func TestMenuPicker_Quit(t *testing.T) {
	config := MenuConfig{
		Title:   "Test",
		Choices: []string{"A"},
	}

	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{
			"q key",
			tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'q'},
			},
		},
		{
			"ctrl+c",
			tea.KeyMsg{Type: tea.KeyCtrlC},
		},
		{"esc", tea.KeyMsg{Type: tea.KeyEsc}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := NewMenuPicker(config)
			model, cmd := menu.Update(tt.key)
			menu = model.(*MenuPicker)

			if !menu.quitting {
				t.Error(
					"Expected quitting to be true",
				)
			}

			if cmd == nil {
				t.Error(
					"Expected tea.Quit command",
				)
			}
		})
	}
}

func TestMenuPicker_Selection(t *testing.T) {
	handlerCalled := false
	selectedIdx := -1

	config := MenuConfig{
		Title:   "Test",
		Choices: []string{"A", "B", "C"},
		SelectHandler: func(index int) tea.Cmd {
			handlerCalled = true
			selectedIdx = index

			return tea.Quit
		},
	}

	menu := NewMenuPicker(config)

	// Move to second option
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ := menu.Update(downMsg)
	menu = model.(*MenuPicker)

	// Select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = menu.Update(enterMsg)
	menu = model.(*MenuPicker)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	if selectedIdx != 1 {
		t.Errorf(
			"Expected selectedIdx = 1, got %d",
			selectedIdx,
		)
	}

	if menu.Selected() != 1 {
		t.Errorf(
			"Expected Selected() = 1, got %d",
			menu.Selected(),
		)
	}
}

func TestMenuPicker_View(t *testing.T) {
	config := MenuConfig{
		Title:   "My Menu",
		Choices: []string{"Option A", "Option B"},
	}

	menu := NewMenuPicker(config)
	view := menu.View()

	// Check title
	if !strings.Contains(view, "My Menu") {
		t.Error("View should contain title")
	}

	// Check choices
	if !strings.Contains(view, "Option A") {
		t.Error(
			"View should contain first choice",
		)
	}
	if !strings.Contains(view, "Option B") {
		t.Error(
			"View should contain second choice",
		)
	}

	// Check help text
	if !strings.Contains(view, "navigate") {
		t.Error(
			"View should contain navigation help",
		)
	}
	if !strings.Contains(view, "Enter") {
		t.Error("View should contain Enter help")
	}
}

func TestMenuPicker_QuitView(t *testing.T) {
	config := MenuConfig{
		Title:   "Test",
		Choices: []string{"A"},
	}

	menu := NewMenuPicker(config)
	menu.quitting = true

	view := menu.View()
	if view != "" {
		t.Errorf(
			"Quit view should be empty, got %q",
			view,
		)
	}
}

func TestMenuPicker_WithSelectHandler(
	t *testing.T,
) {
	config := MenuConfig{
		Title:   "Test",
		Choices: []string{"A"},
	}

	called := false
	menu := NewMenuPicker(
		config,
	).WithSelectHandler(func(idx int) (tea.Model, tea.Cmd) {
		called = true

		return nil, tea.Quit
	})

	// Select
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	menu.Update(enterMsg)

	if !called {
		t.Error(
			"Expected custom handler to be called",
		)
	}
}
