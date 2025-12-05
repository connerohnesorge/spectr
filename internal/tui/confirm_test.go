//nolint:revive // test file
package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewConfirmPicker(t *testing.T) {
	config := ConfirmConfig{
		Question: "Are you sure?",
	}

	picker := NewConfirmPicker(config)

	if picker == nil {
		t.Fatal("NewConfirmPicker returned nil")
	}

	if picker.question != "Are you sure?" {
		t.Errorf("question = %q, want %q", picker.question, "Are you sure?")
	}

	// Default should be "No" (index 1) for safety
	if picker.cursor != 1 {
		t.Errorf("initial cursor = %d, want %d (No)", picker.cursor, 1)
	}

	if picker.confirmed {
		t.Error("confirmed should be false initially")
	}

	if picker.cancelled {
		t.Error("cancelled should be false initially")
	}
}

func TestNewConfirmPicker_DefaultYes(t *testing.T) {
	config := ConfirmConfig{
		Question:   "Are you sure?",
		DefaultYes: true,
	}

	picker := NewConfirmPicker(config)

	if picker == nil {
		t.Fatal("NewConfirmPicker returned nil")
	}

	// With DefaultYes, cursor should be on "Yes" (index 0)
	if picker.cursor != 0 {
		t.Errorf("initial cursor = %d, want %d (Yes)", picker.cursor, 0)
	}
}

func TestConfirmPicker_Navigation(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	// Starts at cursor = 1 (No)

	// Test up navigation (should move to Yes at index 0)
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ := picker.Update(upMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Errorf("After up, cursor = %d, want %d", picker.cursor, 0)
	}

	// Test down navigation (should move back to No at index 1)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ = picker.Update(downMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("After down, cursor = %d, want %d", picker.cursor, 1)
	}

	// Test vim-style navigation: k (up)
	model, _ = picker.Update(upMsg)
	picker = model.(*ConfirmPicker)

	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	// Reset to index 1 first
	picker.cursor = 1
	model, _ = picker.Update(kMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Errorf("After 'k', cursor = %d, want %d", picker.cursor, 0)
	}

	// Test vim-style navigation: j (down)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	model, _ = picker.Update(jMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("After 'j', cursor = %d, want %d", picker.cursor, 1)
	}
}

func TestConfirmPicker_Boundaries(t *testing.T) {
	config := ConfirmConfig{
		Question:   "Test?",
		DefaultYes: true, // Start at Yes (index 0)
	}

	picker := NewConfirmPicker(config)

	// Try to go up when at top (Yes = 0)
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ := picker.Update(upMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Errorf("Cursor should stay at 0, got %d", picker.cursor)
	}

	// Go to bottom (No = 1)
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ = picker.Update(downMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("Cursor should be at 1, got %d", picker.cursor)
	}

	// Try to go down past end
	model, _ = picker.Update(downMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("Cursor should stay at 1, got %d", picker.cursor)
	}
}

func TestConfirmPicker_SelectYes(t *testing.T) {
	config := ConfirmConfig{
		Question:   "Test?",
		DefaultYes: true, // Start at Yes (index 0)
	}

	picker := NewConfirmPicker(config)

	// Cursor is at 0 (Yes), press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := picker.Update(enterMsg)
	picker = model.(*ConfirmPicker)

	if !picker.confirmed {
		t.Error("Expected confirmed to be true when selecting Yes")
	}

	if picker.Confirmed() != true {
		t.Error("Expected Confirmed() to return true")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command")
	}
}

func TestConfirmPicker_SelectNo(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
		// Default: cursor at No (index 1)
	}

	picker := NewConfirmPicker(config)

	// Cursor is at 1 (No), press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := picker.Update(enterMsg)
	picker = model.(*ConfirmPicker)

	if picker.confirmed {
		t.Error("Expected confirmed to be false when selecting No")
	}

	if picker.Confirmed() != false {
		t.Error("Expected Confirmed() to return false")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command")
	}
}

func TestConfirmPicker_Cancel(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"q key", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}},
		{"esc", tea.KeyMsg{Type: tea.KeyEsc}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			picker := NewConfirmPicker(config)
			model, cmd := picker.Update(tt.key)
			picker = model.(*ConfirmPicker)

			if !picker.cancelled {
				t.Error("Expected cancelled to be true")
			}

			if picker.Cancelled() != true {
				t.Error("Expected Cancelled() to return true")
			}

			if cmd == nil {
				t.Error("Expected tea.Quit command")
			}
		})
	}
}

func TestConfirmPicker_View(t *testing.T) {
	config := ConfirmConfig{
		Question: "Delete all files?",
	}

	picker := NewConfirmPicker(config)
	view := picker.View()

	// Check question
	if !strings.Contains(view, "Delete all files?") {
		t.Error("View should contain question")
	}

	// Check choices
	if !strings.Contains(view, "Yes") {
		t.Error("View should contain 'Yes' choice")
	}
	if !strings.Contains(view, "No") {
		t.Error("View should contain 'No' choice")
	}

	// Check help text
	if !strings.Contains(view, "navigate") {
		t.Error("View should contain navigation help")
	}
	if !strings.Contains(view, "Enter") {
		t.Error("View should contain Enter help")
	}
	if !strings.Contains(view, "cancel") {
		t.Error("View should contain cancel help")
	}
}

func TestConfirmPicker_CancelledView(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	picker.cancelled = true

	view := picker.View()
	if view != "" {
		t.Errorf("Cancelled view should be empty, got %q", view)
	}
}

func TestConfirmPicker_Init(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	cmd := picker.Init()

	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestConfirmPicker_NonKeyMsg(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	initialCursor := picker.cursor

	// Send a non-key message (e.g., window size)
	model, cmd := picker.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	picker = model.(*ConfirmPicker)

	if picker.cursor != initialCursor {
		t.Error("Cursor should not change for non-key messages")
	}

	if cmd != nil {
		t.Error("Should return nil cmd for non-key messages")
	}
}
