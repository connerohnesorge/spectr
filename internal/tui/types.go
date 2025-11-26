package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// ActionHandler is a function that handles a key action on a selected row.
// It receives the selected row and returns a tea.Cmd and an optional result.
type ActionHandler func(row table.Row) (tea.Cmd, *ActionResult)

// Action represents a configurable key action for TablePicker.
type Action struct {
	Key         string        // The key binding (e.g., "e", "a", "enter")
	Description string        // Human-readable description for help text
	Handler     ActionHandler // The function to call when key is pressed
}

// ActionResult represents the result of an action.
type ActionResult struct {
	// ID is the selected item ID (usually from first column).
	ID string

	// Quit indicates the TUI should exit.
	Quit bool

	// Cancelled indicates the user cancelled.
	Cancelled bool

	// Copied indicates an ID was copied to clipboard.
	Copied bool

	// ArchiveRequested indicates archive action was requested.
	ArchiveRequested bool

	// Error contains any error from the action.
	Error error

	// Custom allows actions to pass custom data.
	Custom any
}

// TableConfig holds configuration for TablePicker.
type TableConfig struct {
	// Columns defines the table columns.
	Columns []table.Column

	// Rows contains the table data.
	Rows []table.Row

	// Height is the visible height of the table.
	Height int

	// Actions is a map of key bindings to actions.
	Actions map[string]Action

	// HelpText is custom help text; if empty, generated from actions.
	HelpText string

	// ProjectPath is the project root for file operations.
	ProjectPath string

	// FooterExtra is additional text to show after help.
	FooterExtra string
}

// MenuConfig holds configuration for MenuPicker.
type MenuConfig struct {
	// Title is the menu title.
	Title string

	// Choices are the menu options.
	Choices []string

	// SelectHandler is called when an option is selected.
	// Receives the index of the selected choice.
	SelectHandler func(index int) tea.Cmd
}
