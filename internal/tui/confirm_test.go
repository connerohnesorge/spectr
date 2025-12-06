//nolint:revive // test file
package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewConfirmPicker_DefaultNo(t *testing.T) {
	// When DefaultNo is true (or not specified), cursor should start on "No" (index 1)
	config := ConfirmConfig{
		Question:  "Are you sure?",
		DefaultNo: true,
	}

	picker := NewConfirmPicker(config)

	if picker == nil {
		t.Fatal("NewConfirmPicker returned nil")
	}

	if picker.cursor != 1 {
		t.Errorf("cursor = %d, want %d (No)", picker.cursor, 1)
	}

	if picker.question != "Are you sure?" {
		t.Errorf("question = %q, want %q", picker.question, "Are you sure?")
	}

	// Test with zero value (default behavior)
	configZero := ConfirmConfig{
		Question: "Test?",
		// DefaultNo not set - defaults to false in Go, but NewConfirmPicker treats it as cursor=0
	}

	pickerZero := NewConfirmPicker(configZero)

	// When DefaultNo is false, cursor should be 0 (Yes)
	if pickerZero.cursor != 0 {
		t.Errorf("zero value config: cursor = %d, want %d (Yes)", pickerZero.cursor, 0)
	}
}

func TestNewConfirmPicker_DefaultYes(t *testing.T) {
	// When DefaultNo is false, cursor should start on "Yes" (index 0)
	config := ConfirmConfig{
		Question:  "Continue?",
		DefaultNo: false,
	}

	picker := NewConfirmPicker(config)

	if picker == nil {
		t.Fatal("NewConfirmPicker returned nil")
	}

	if picker.cursor != 0 {
		t.Errorf("cursor = %d, want %d (Yes)", picker.cursor, 0)
	}

	if len(picker.choices) != 2 {
		t.Errorf("choices length = %d, want %d", len(picker.choices), 2)
	}

	if picker.choices[0] != "Yes" {
		t.Errorf("choices[0] = %q, want %q", picker.choices[0], "Yes")
	}

	if picker.choices[1] != "No" {
		t.Errorf("choices[1] = %q, want %q", picker.choices[1], "No")
	}
}

func TestConfirmPicker_ConfirmYes(t *testing.T) {
	config := ConfirmConfig{
		Question:  "Confirm?",
		DefaultNo: false, // Start on Yes
	}

	picker := NewConfirmPicker(config)

	// Cursor is on Yes (index 0), press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := picker.Update(enterMsg)
	picker = model.(*ConfirmPicker)

	if !picker.confirmed {
		t.Error("Expected confirmed to be true after Enter on Yes")
	}

	if picker.quitting {
		t.Error("Expected quitting to be false")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command")
	}

	if !picker.Confirmed() {
		t.Error("Confirmed() should return true")
	}
}

func TestConfirmPicker_ConfirmNo(t *testing.T) {
	config := ConfirmConfig{
		Question:  "Confirm?",
		DefaultNo: true, // Start on No
	}

	picker := NewConfirmPicker(config)

	// Cursor is on No (index 1), press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := picker.Update(enterMsg)
	picker = model.(*ConfirmPicker)

	if picker.confirmed {
		t.Error("Expected confirmed to be false after Enter on No")
	}

	if picker.quitting {
		t.Error("Expected quitting to be false")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command")
	}

	if picker.Confirmed() {
		t.Error("Confirmed() should return false")
	}
}

func TestConfirmPicker_QuitBehavior(t *testing.T) {
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

			if !picker.quitting {
				t.Error("Expected quitting to be true")
			}

			if picker.confirmed {
				t.Error("Expected confirmed to be false after quit")
			}

			if cmd == nil {
				t.Error("Expected tea.Quit command")
			}

			if !picker.Quitting() {
				t.Error("Quitting() should return true")
			}

			if picker.Confirmed() {
				t.Error("Confirmed() should return false after quit")
			}
		})
	}
}

func TestConfirmPicker_NavigationUp(t *testing.T) {
	config := ConfirmConfig{
		Question:  "Test?",
		DefaultNo: true, // Start on No (index 1)
	}

	picker := NewConfirmPicker(config)

	// Test up arrow navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ := picker.Update(upMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Errorf("After up, cursor = %d, want %d (Yes)", picker.cursor, 0)
	}

	// Test boundary - can't go above 0
	model, _ = picker.Update(upMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Errorf("After up at boundary, cursor = %d, want %d", picker.cursor, 0)
	}

	// Test vim-style k navigation
	picker = NewConfirmPicker(config) // Reset to start on No
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	model, _ = picker.Update(kMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Errorf("After 'k', cursor = %d, want %d (Yes)", picker.cursor, 0)
	}
}

func TestConfirmPicker_NavigationDown(t *testing.T) {
	config := ConfirmConfig{
		Question:  "Test?",
		DefaultNo: false, // Start on Yes (index 0)
	}

	picker := NewConfirmPicker(config)

	// Test down arrow navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ := picker.Update(downMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("After down, cursor = %d, want %d (No)", picker.cursor, 1)
	}

	// Test boundary - can't go below 1 (last index)
	model, _ = picker.Update(downMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("After down at boundary, cursor = %d, want %d", picker.cursor, 1)
	}

	// Test vim-style j navigation
	picker = NewConfirmPicker(config) // Reset to start on Yes
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	model, _ = picker.Update(jMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Errorf("After 'j', cursor = %d, want %d (No)", picker.cursor, 1)
	}
}

func TestConfirmPicker_ViewOutput(t *testing.T) {
	config := ConfirmConfig{
		Question:  "Are you sure you want to proceed?",
		DefaultNo: false, // Start on Yes
	}

	picker := NewConfirmPicker(config)
	view := picker.View()

	// Check question is displayed
	if !strings.Contains(view, "Are you sure you want to proceed?") {
		t.Error("View should contain the question")
	}

	// Check choices are displayed
	if !strings.Contains(view, "Yes") {
		t.Error("View should contain Yes option")
	}
	if !strings.Contains(view, "No") {
		t.Error("View should contain No option")
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

	// Check cursor indicator
	if !strings.Contains(view, ">") {
		t.Error("View should contain cursor indicator")
	}
}

func TestConfirmPicker_QuitView(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	picker.quitting = true

	view := picker.View()
	if view != "" {
		t.Errorf("Quit view should be empty, got %q", view)
	}
}

func TestConfirmPicker_Init(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	cmd := picker.Init()

	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestConfirmPicker_NonKeyMsg(t *testing.T) {
	config := ConfirmConfig{
		Question: "Test?",
	}

	picker := NewConfirmPicker(config)
	initialCursor := picker.cursor

	// Send a non-KeyMsg (e.g., WindowSizeMsg)
	model, cmd := picker.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	picker = model.(*ConfirmPicker)

	if picker.cursor != initialCursor {
		t.Error("Non-KeyMsg should not change cursor")
	}

	if cmd != nil {
		t.Error("Non-KeyMsg should return nil command")
	}
}

func TestConfirmPicker_NavigateAndConfirm(t *testing.T) {
	// Test navigating from No to Yes and confirming
	config := ConfirmConfig{
		Question:  "Test?",
		DefaultNo: true, // Start on No
	}

	picker := NewConfirmPicker(config)

	// Navigate up to Yes
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ := picker.Update(upMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 0 {
		t.Fatalf("Expected cursor to be on Yes (0), got %d", picker.cursor)
	}

	// Confirm
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = picker.Update(enterMsg)
	picker = model.(*ConfirmPicker)

	if !picker.confirmed {
		t.Error("Expected confirmed to be true after navigating to Yes and pressing Enter")
	}
}

func TestConfirmPicker_NavigateAndDecline(t *testing.T) {
	// Test navigating from Yes to No and declining
	config := ConfirmConfig{
		Question:  "Test?",
		DefaultNo: false, // Start on Yes
	}

	picker := NewConfirmPicker(config)

	// Navigate down to No
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ := picker.Update(downMsg)
	picker = model.(*ConfirmPicker)

	if picker.cursor != 1 {
		t.Fatalf("Expected cursor to be on No (1), got %d", picker.cursor)
	}

	// Confirm selection (which is No)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = picker.Update(enterMsg)
	picker = model.(*ConfirmPicker)

	if picker.confirmed {
		t.Error("Expected confirmed to be false after navigating to No and pressing Enter")
	}
}
