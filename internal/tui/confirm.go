//nolint:revive // TUI code - interactive model patterns require specific structure
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmPicker is a simple Yes/No confirmation menu.
// It provides a two-option selection with navigation and confirmation support.
type ConfirmPicker struct {
	// question is the confirmation prompt displayed at the top.
	question string
	// choices contains the Yes/No options.
	choices []string
	// cursor tracks the currently highlighted item (0=Yes, 1=No).
	cursor int
	// confirmed indicates whether "Yes" was selected.
	confirmed bool
	// quitting indicates whether the menu was quit/cancelled.
	quitting bool
}

// NewConfirmPicker creates a new ConfirmPicker with the given configuration.
//
// Default behavior:
//   - If config.DefaultNo is true (default), cursor starts on "No" for safety
//   - If config.DefaultNo is false, cursor starts on "Yes"
//
// This default ensures users must explicitly navigate to "Yes" to confirm
// potentially destructive operations.
func NewConfirmPicker(config ConfirmConfig) *ConfirmPicker {
	choices := []string{"Yes", "No"}

	// Default to "No" (index 1) for safety unless explicitly overridden
	cursor := 1
	if !config.DefaultNo {
		cursor = 0
	}

	return &ConfirmPicker{
		question: config.Question,
		choices:  choices,
		cursor:   cursor,
	}
}

// Init implements tea.Model.
func (m *ConfirmPicker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *ConfirmPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "q", "ctrl+c", "esc":
		m.quitting = true
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
		// cursor 0 = Yes (confirmed), cursor 1 = No (not confirmed)
		m.confirmed = m.cursor == 0

		return m, tea.Quit
	}

	return m, nil
}

// View implements tea.Model.
func (m *ConfirmPicker) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := TitleStyle()
	choiceStyle := ChoiceStyle()
	selectedStyle := SelectedStyle()
	helpStyle := HelpStyle()

	s := titleStyle.Render(m.question) + "\n\n"

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

	s += "\n" + helpStyle.Render("^/v or j/k: navigate | Enter: confirm | q/esc: cancel")

	return s
}

// Confirmed returns whether "Yes" was selected.
func (m *ConfirmPicker) Confirmed() bool {
	return m.confirmed
}

// Quitting returns whether the menu was quit/cancelled.
func (m *ConfirmPicker) Quitting() bool {
	return m.quitting
}

// Run runs the ConfirmPicker and returns (true, nil) for Yes, (false, nil) for No/cancelled.
// Returns (false, error) if there was an error running the TUI.
func (m *ConfirmPicker) Run() (bool, error) {
	prog := tea.NewProgram(m)
	finalModel, err := prog.Run()
	if err != nil {
		return false, fmt.Errorf("error running confirmation: %w", err)
	}

	if fm, ok := finalModel.(*ConfirmPicker); ok {
		if fm.quitting {
			return false, nil
		}

		return fm.confirmed, nil
	}

	return false, nil
}
