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
	// countPrefixState manages vim-style count prefix navigation (e.g., "9j").
	countPrefixState CountPrefixState
}

// NewMenuPicker creates a new MenuPicker with the given configuration.
//
// Selection behavior:
//   - If config.SelectHandler is provided, it will be called on selection
//   - If config.SelectHandler is nil, the menu will quit on selection (default)
//
// This default allows callers to create simple menus without specifying a handler
// when they only need to retrieve the selected index via Run().
func NewMenuPicker(
	config MenuConfig,
) *MenuPicker {
	return &MenuPicker{
		title:   config.Title,
		choices: config.Choices,
		cursor:  0,
		selectHandler: func(index int) (tea.Model, tea.Cmd) {
			// Default: call custom handler if provided, otherwise just quit.
			// This enables simple usage patterns where caller only needs the index.
			if config.SelectHandler != nil {
				return nil, config.SelectHandler(
					index,
				)
			}

			return nil, tea.Quit
		},
	}
}

// WithSelectHandler sets the handler for menu selection.
// The handler receives the selected index and can return a new model to transition to.
func (m *MenuPicker) WithSelectHandler(
	handler func(index int) (tea.Model, tea.Cmd),
) *MenuPicker {
	m.selectHandler = handler

	return m
}

// Init implements tea.Model.
func (m *MenuPicker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *MenuPicker) Update(
	msg tea.Msg,
) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	keyStr := keyMsg.String()

	// Handle quit keys (but check count prefix state for ESC)
	if keyStr == "q" || keyStr == "ctrl+c" {
		m.quitting = true

		return m, tea.Quit
	}

	// Handle ESC - quit only if count prefix is not active
	if keyStr == keyEsc {
		if m.countPrefixState.IsActive() {
			// Reset count prefix and continue
			m.countPrefixState.Reset()

			return m, nil
		}

		// No count prefix active, so quit
		m.quitting = true

		return m, tea.Quit
	}

	// Handle count prefix navigation
	count, isNavKey, handled := m.countPrefixState.HandleKey(keyMsg)
	if handled && isNavKey {
		// Apply counted navigation with boundary checks
		switch keyStr {
		case keyUp, "k":
			m.cursor = maxInt(0, m.cursor-count)
		case keyDown, "j":
			m.cursor = minInt(len(m.choices)-1, m.cursor+count)
		}

		return m, nil
	}

	// Handle enter key for selection
	if keyStr == "enter" {
		m.selected = m.cursor

		return m.handleSelection()
	}

	return m, nil
}

// maxInt returns the larger of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// minInt returns the smaller of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// handleSelection processes the menu selection.
func (m *MenuPicker) handleSelection() (tea.Model, tea.Cmd) {
	if m.selectHandler != nil {
		newModel, cmd := m.selectHandler(
			m.selected,
		)
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
			s += selectedStyle.Render(
				fmt.Sprintf(
					"%s %s",
					cursor,
					choice,
				),
			) + "\n"
		} else {
			s += choiceStyle.Render(fmt.Sprintf("%s %s", cursor, choice)) + "\n"
		}
	}

	helpText := "↑/↓ or j/k: navigate | Enter: select | q: quit"
	if m.countPrefixState.IsActive() {
		helpText += fmt.Sprintf(" | count: %s_", m.countPrefixState.String())
	}

	s += "\n" + helpStyle.Render(helpText)

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
		return -1, fmt.Errorf(
			"error running menu: %w",
			err,
		)
	}

	if fm, ok := finalModel.(*MenuPicker); ok {
		if fm.quitting {
			return -1, nil
		}

		return fm.selected, nil
	}

	return -1, nil
}
