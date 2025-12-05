//nolint:revive // TUI code - interactive model patterns require specific structure
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmConfig holds configuration for ConfirmPicker.
type ConfirmConfig struct {
	// Question is the confirmation question to display.
	Question string

	// DefaultYes sets the default selection to Yes (default is No for safety).
	DefaultYes bool
}

// ConfirmPicker is a Yes/No confirmation prompt.
// It provides a simple binary choice with keyboard navigation.
type ConfirmPicker struct {
	// question is the confirmation prompt displayed at the top.
	question string
	// cursor tracks the currently highlighted option (0 = Yes, 1 = No).
	cursor int
	// confirmed indicates whether Yes was selected.
	confirmed bool
	// cancelled indicates whether the user quit/cancelled.
	cancelled bool
}

// NewConfirmPicker creates a new ConfirmPicker with the given configuration.
//
// By default, the cursor starts on "No" for safety. Set DefaultYes to true
// in the config to start on "Yes" instead.
func NewConfirmPicker(config ConfirmConfig) *ConfirmPicker {
	cursor := 1 // Default to "No" (index 1) for safety
	if config.DefaultYes {
		cursor = 0 // "Yes" is at index 0
	}

	return &ConfirmPicker{
		question: config.Question,
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
		m.cancelled = true

		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < 1 {
			m.cursor++
		}

	case "enter":
		m.confirmed = m.cursor == 0 // Yes is at index 0

		return m, tea.Quit
	}

	return m, nil
}

// View implements tea.Model.
func (m *ConfirmPicker) View() string {
	if m.cancelled {
		return ""
	}

	titleStyle := TitleStyle()
	choiceStyle := ChoiceStyle()
	selectedStyle := SelectedStyle()
	helpStyle := HelpStyle()

	s := titleStyle.Render(m.question) + "\n\n"

	choices := []string{"Yes", "No"}
	for i, choice := range choices {
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

	s += "\n" + helpStyle.Render("↑/↓ or j/k: navigate | Enter: confirm | q/esc: cancel")

	return s
}

// Confirmed returns whether Yes was selected.
func (m *ConfirmPicker) Confirmed() bool {
	return m.confirmed
}

// Cancelled returns whether the prompt was cancelled.
func (m *ConfirmPicker) Cancelled() bool {
	return m.cancelled
}

// Run runs the ConfirmPicker and returns the result.
//
// Returns:
//   - confirmed: true if user selected Yes, false if selected No
//   - cancelled: true if user quit/cancelled (q, esc, ctrl+c)
//   - err: any error from running the TUI
func (m *ConfirmPicker) Run() (confirmed bool, cancelled bool, err error) {
	prog := tea.NewProgram(m)
	finalModel, err := prog.Run()
	if err != nil {
		return false, false, fmt.Errorf("error running confirm prompt: %w", err)
	}

	if fm, ok := finalModel.(*ConfirmPicker); ok {
		return fm.confirmed, fm.cancelled, nil
	}

	return false, true, nil
}
