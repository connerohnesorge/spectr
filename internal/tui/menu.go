//nolint:revive // TUI code - interactive model patterns require specific structure
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuPicker is a simple menu-based selector.
// It provides a vertical list of choices with navigation and selection support.
type MenuPicker struct {
	// title is the menu title displayed at the top.
	title string
	// choices contains the menu options to display.
	choices []string
	// cursor tracks the currently highlighted item.
	cursor int
	// selected stores the index of the selected item after selection.
	selected int
	// quitting indicates whether the menu was quit/cancelled.
	quitting bool
	// selectHandler is called when an item is selected.
	selectHandler func(index int) (tea.Model, tea.Cmd)
}

// NewMenuPicker creates a new MenuPicker with the given configuration.
func NewMenuPicker(config MenuConfig) *MenuPicker {
	return &MenuPicker{
		title:   config.Title,
		choices: config.Choices,
		cursor:  0,
		selectHandler: func(index int) (tea.Model, tea.Cmd) {
			if config.SelectHandler != nil {
				return nil, config.SelectHandler(index)
			}

			return nil, tea.Quit
		},
	}
}

// WithSelectHandler sets the handler for menu selection.
// The handler receives the selected index and can return a new model to transition to.
func (m *MenuPicker) WithSelectHandler(handler func(index int) (tea.Model, tea.Cmd)) *MenuPicker {
	m.selectHandler = handler

	return m
}

// Init implements tea.Model.
func (m *MenuPicker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *MenuPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "q", "ctrl+c", "esc":
		m.quitting = true

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
		m.selected = m.cursor

		return m.handleSelection()
	}

	return m, nil
}

// handleSelection processes the menu selection.
func (m *MenuPicker) handleSelection() (tea.Model, tea.Cmd) {
	if m.selectHandler != nil {
		newModel, cmd := m.selectHandler(m.selected)
		if newModel != nil {
			return newModel, cmd
		}

		return m, cmd
	}

	return m, tea.Quit
}

// View implements tea.Model.
func (m *MenuPicker) View() string {
	if m.quitting {
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

	s += "\n" + helpStyle.Render("↑/↓ or j/k: navigate | Enter: select | q: quit")

	return s
}

// Selected returns the index of the selected choice.
func (m *MenuPicker) Selected() int {
	return m.selected
}

// Quitting returns whether the menu was quit.
func (m *MenuPicker) Quitting() bool {
	return m.quitting
}

// Run runs the MenuPicker and returns the selected index, or -1 if cancelled.
func (m *MenuPicker) Run() (int, error) {
	prog := tea.NewProgram(m)
	finalModel, err := prog.Run()
	if err != nil {
		return -1, fmt.Errorf("error running menu: %w", err)
	}

	if fm, ok := finalModel.(*MenuPicker); ok {
		if fm.quitting {
			return -1, nil
		}

		return fm.selected, nil
	}

	return -1, nil
}
