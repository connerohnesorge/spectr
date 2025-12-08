//nolint:revive // test file
package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewConfirmPrompt(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		config := ConfirmConfig{
			Title: "Test Prompt",
		}

		prompt := NewConfirmPrompt(config)

		if prompt == nil {
			t.Fatal("NewConfirmPrompt returned nil")
		}

		if prompt.title != "Test Prompt" {
			t.Errorf("title = %q, want %q", prompt.title, "Test Prompt")
		}

		// Check default CancelText is "No"
		if prompt.choices[0] != "No" {
			t.Errorf("default CancelText = %q, want %q", prompt.choices[0], "No")
		}

		// Check default ConfirmText is "Yes"
		if prompt.choices[1] != "Yes" {
			t.Errorf("default ConfirmText = %q, want %q", prompt.choices[1], "Yes")
		}

		// Check default cursor is 0 (cancel/no option)
		if prompt.cursor != 0 {
			t.Errorf("default cursor = %d, want %d", prompt.cursor, 0)
		}
	})

	t.Run("custom texts preserved", func(t *testing.T) {
		config := ConfirmConfig{
			Title:       "Custom Prompt",
			CancelText:  "Keep",
			ConfirmText: "Delete",
		}

		prompt := NewConfirmPrompt(config)

		if prompt.choices[0] != "Keep" {
			t.Errorf("CancelText = %q, want %q", prompt.choices[0], "Keep")
		}

		if prompt.choices[1] != "Delete" {
			t.Errorf("ConfirmText = %q, want %q", prompt.choices[1], "Delete")
		}
	})

	t.Run("DefaultConfirm moves cursor to index 1", func(t *testing.T) {
		config := ConfirmConfig{
			Title:          "Confirm Default",
			DefaultConfirm: true,
		}

		prompt := NewConfirmPrompt(config)

		if prompt.cursor != 1 {
			t.Errorf("cursor with DefaultConfirm = %d, want %d", prompt.cursor, 1)
		}
	})
}

func TestConfirmPrompt_Navigation(t *testing.T) {
	config := ConfirmConfig{
		Title: "Test Navigation",
	}

	t.Run("down arrow moves cursor down", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)

		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		model, _ := prompt.Update(downMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.cursor != 1 {
			t.Errorf("After down, cursor = %d, want %d", prompt.cursor, 1)
		}
	})

	t.Run("j key moves cursor down", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)

		jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
		model, _ := prompt.Update(jMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.cursor != 1 {
			t.Errorf("After 'j', cursor = %d, want %d", prompt.cursor, 1)
		}
	})

	t.Run("up arrow moves cursor up", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		prompt.cursor = 1 // Start at bottom

		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		model, _ := prompt.Update(upMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.cursor != 0 {
			t.Errorf("After up, cursor = %d, want %d", prompt.cursor, 0)
		}
	})

	t.Run("k key moves cursor up", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		prompt.cursor = 1 // Start at bottom

		kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
		model, _ := prompt.Update(kMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.cursor != 0 {
			t.Errorf("After 'k', cursor = %d, want %d", prompt.cursor, 0)
		}
	})
}

func TestConfirmPrompt_Selection(t *testing.T) {
	config := ConfirmConfig{
		Title: "Test Selection",
	}

	t.Run("enter on index 0 sets confirmed=false", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		// cursor starts at 0

		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		model, cmd := prompt.Update(enterMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.confirmed {
			t.Error("Expected confirmed to be false when selecting index 0")
		}

		if cmd == nil {
			t.Error("Expected tea.Quit command")
		}
	})

	t.Run("enter on index 1 sets confirmed=true", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		prompt.cursor = 1 // Move to confirm option

		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		model, cmd := prompt.Update(enterMsg)
		prompt = model.(*ConfirmPrompt)

		if !prompt.confirmed {
			t.Error("Expected confirmed to be true when selecting index 1")
		}

		if cmd == nil {
			t.Error("Expected tea.Quit command")
		}
	})

	t.Run("Confirmed and Cancelled return correct values", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		prompt.cursor = 1

		enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
		model, _ := prompt.Update(enterMsg)
		prompt = model.(*ConfirmPrompt)

		if !prompt.Confirmed() {
			t.Error("Expected Confirmed() to return true")
		}

		if prompt.Cancelled() {
			t.Error("Expected Cancelled() to return false")
		}
	})
}

func TestConfirmPrompt_Cancel(t *testing.T) {
	config := ConfirmConfig{
		Title: "Test Cancel",
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
			prompt := NewConfirmPrompt(config)
			model, cmd := prompt.Update(tt.key)
			prompt = model.(*ConfirmPrompt)

			if !prompt.cancelled {
				t.Errorf("%s: Expected cancelled to be true", tt.name)
			}

			if prompt.confirmed {
				t.Errorf("%s: Expected confirmed to be false", tt.name)
			}

			if cmd == nil {
				t.Errorf("%s: Expected tea.Quit command", tt.name)
			}
		})
	}
}

func TestConfirmPrompt_View(t *testing.T) {
	config := ConfirmConfig{
		Title:       "Delete this file?",
		CancelText:  "Keep",
		ConfirmText: "Delete",
	}

	prompt := NewConfirmPrompt(config)
	view := prompt.View()

	t.Run("view contains title", func(t *testing.T) {
		if !strings.Contains(view, "Delete this file?") {
			t.Error("View should contain title")
		}
	})

	t.Run("view contains both choices", func(t *testing.T) {
		if !strings.Contains(view, "Keep") {
			t.Error("View should contain CancelText 'Keep'")
		}
		if !strings.Contains(view, "Delete") {
			t.Error("View should contain ConfirmText 'Delete'")
		}
	})

	t.Run("view contains navigation help", func(t *testing.T) {
		if !strings.Contains(view, "navigate") {
			t.Error("View should contain navigation help")
		}
		if !strings.Contains(view, "Enter") {
			t.Error("View should contain Enter help")
		}
		if !strings.Contains(view, "cancel") {
			t.Error("View should contain cancel help")
		}
	})

	t.Run("cancelled view is empty", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		prompt.cancelled = true

		view := prompt.View()
		if view != "" {
			t.Errorf("Cancelled view should be empty, got %q", view)
		}
	})
}

func TestConfirmPrompt_Boundaries(t *testing.T) {
	config := ConfirmConfig{
		Title: "Test Boundaries",
	}

	t.Run("cannot move cursor above 0", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		// cursor starts at 0

		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		model, _ := prompt.Update(upMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.cursor != 0 {
			t.Errorf("Cursor should stay at 0, got %d", prompt.cursor)
		}
	})

	t.Run("cannot move cursor below 1", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)
		prompt.cursor = 1 // Move to bottom

		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		model, _ := prompt.Update(downMsg)
		prompt = model.(*ConfirmPrompt)

		if prompt.cursor != 1 {
			t.Errorf("Cursor should stay at 1, got %d", prompt.cursor)
		}
	})

	t.Run("multiple boundary violations stay at bounds", func(t *testing.T) {
		prompt := NewConfirmPrompt(config)

		// Try going up multiple times from 0
		upMsg := tea.KeyMsg{Type: tea.KeyUp}
		for range 5 {
			model, _ := prompt.Update(upMsg)
			prompt = model.(*ConfirmPrompt)
		}

		if prompt.cursor != 0 {
			t.Errorf("Cursor should remain at 0 after multiple ups, got %d", prompt.cursor)
		}

		// Go to bottom and try going down multiple times
		downMsg := tea.KeyMsg{Type: tea.KeyDown}
		model, _ := prompt.Update(downMsg)
		prompt = model.(*ConfirmPrompt)

		for range 5 {
			model, _ = prompt.Update(downMsg)
			prompt = model.(*ConfirmPrompt)
		}

		if prompt.cursor != 1 {
			t.Errorf("Cursor should remain at 1 after multiple downs, got %d", prompt.cursor)
		}
	})
}

func TestConfirmPrompt_Init(t *testing.T) {
	config := ConfirmConfig{
		Title: "Test Init",
	}

	prompt := NewConfirmPrompt(config)
	cmd := prompt.Init()

	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestConfirmPrompt_NonKeyMsg(t *testing.T) {
	config := ConfirmConfig{
		Title: "Test Non-Key",
	}

	prompt := NewConfirmPrompt(config)
	initialCursor := prompt.cursor

	// Send a non-KeyMsg (e.g., a custom message type)
	type customMsg struct{}
	model, cmd := prompt.Update(customMsg{})
	prompt = model.(*ConfirmPrompt)

	if prompt.cursor != initialCursor {
		t.Error("Non-key message should not change cursor")
	}

	if cmd != nil {
		t.Error("Non-key message should return nil command")
	}
}
