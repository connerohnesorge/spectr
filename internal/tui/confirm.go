//nolint:revive // TUI code - interactive model patterns require specific structure
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmPrompt is a simple yes/no confirmation dialog.
// It provides two choices with navigation and selection support.
// By default, the "no/cancel" option is selected first for safety.
type ConfirmPrompt struct {
	// title is the prompt title/question displayed at the top.
	title string
	// choices contains the two options: [0] = cancel/no, [1] = confirm/yes.
	choices []string
	// cursor tracks the currently highlighted item (0 or 1).
	cursor int
	// confirmed indicates whether the user confirmed (selected yes).
	confirmed bool
	// cancelled indicates whether the prompt was quit/cancelled via q/Esc.
	cancelled bool
}

// NewConfirmPrompt creates a new ConfirmPrompt with the given configuration.
//
// The prompt displays two options:
//   - Index 0: Cancel/No option (default selection for safety)
//   - Index 1: Confirm/Yes option
//
// If CancelText is empty, defaults to "No".
// If ConfirmText is empty, defaults to "Yes".
// If DefaultConfirm is true, the cursor starts on the confirm option.
func NewConfirmPrompt(config ConfirmConfig) *ConfirmPrompt {
	cancelText := config.CancelText
	if cancelText == "" {
		cancelText = "No"
	}

	confirmText := config.ConfirmText
	if confirmText == "" {
		confirmText = "Yes"
	}

	cursor := 0 // Default to cancel/no for safety
	if config.DefaultConfirm {
		cursor = 1
	}

	return &ConfirmPrompt{
		title:   config.Title,
		choices: []string{cancelText, confirmText},
		cursor:  cursor,
	}
}

// Init implements tea.Model.
func (m *ConfirmPrompt) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *ConfirmPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "q", "ctrl+c", "esc":
		// Cancelled - treat as "no/keep"
		m.cancelled = true
		m.confirmed = false

		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}

	case "enter":
		// Index 0 = cancel/no, Index 1 = confirm/yes
		m.confirmed = m.cursor == 1

		return m, tea.Quit
	}

	return m, nil
}

// View implements tea.Model.
func (m *ConfirmPrompt) View() string {
	if m.cancelled {
		return ""
	}

	titleStyle := TitleStyle()
	choiceStyle := ChoiceStyle()
	selectedStyle := SelectedStyle()
	helpStyle := HelpStyle()

	s := titleStyle.Render(m.title) + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		if m.cursor == i {
			s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
		} else {
			s += choiceStyle.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
		}
	}

	s += "\n" + helpStyle.Render("Up/Down or j/k: navigate | Enter: select | q/Esc: cancel")

	return s
}

// Confirmed returns whether the user confirmed (selected yes).
func (m *ConfirmPrompt) Confirmed() bool {
	return m.confirmed
}

// Cancelled returns whether the prompt was cancelled via q/Esc.
func (m *ConfirmPrompt) Cancelled() bool {
	return m.cancelled
}

// Run runs the ConfirmPrompt and returns whether the user confirmed.
// Returns (false, nil) if cancelled or if "no" was selected.
// Returns (true, nil) if "yes" was selected.
// Returns (false, error) if there was an error running the prompt.
func (m *ConfirmPrompt) Run() (bool, error) {
	prog := tea.NewProgram(m)
	finalModel, err := prog.Run()
	if err != nil {
		return false, fmt.Errorf("error running confirm prompt: %w", err)
	}

	if fm, ok := finalModel.(*ConfirmPrompt); ok {
		// If cancelled (q/Esc), return false (same as selecting "no")
		if fm.cancelled {
			return false, nil
		}

		return fm.confirmed, nil
	}

	return false, nil
}
