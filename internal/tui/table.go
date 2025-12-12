//nolint:revive // TUI code - interactive model patterns require specific structure
package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

const keyCtrlC = "ctrl+c"

const (
	// DefaultTableHeight is the default number of visible rows.
	DefaultTableHeight = 10
)

// TablePicker is a configurable table-based item selector.
type TablePicker struct {
	table       table.Model
	actions     map[string]Action
	result      *ActionResult
	quitting    bool
	helpText    string
	projectPath string
	footerExtra string
	err         error
	showHelp    bool
}

// NewTablePicker creates a new TablePicker with the given configuration.
func NewTablePicker(
	config TableConfig,
) *TablePicker {
	height := config.Height
	if height == 0 {
		height = DefaultTableHeight
	}

	t := table.New(
		table.WithColumns(config.Columns),
		table.WithRows(config.Rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)
	ApplyTableStyles(&t)

	picker := &TablePicker{
		table:       t,
		actions:     make(map[string]Action),
		projectPath: config.ProjectPath,
		footerExtra: config.FooterExtra,
	}

	// Register provided actions
	for key, action := range config.Actions {
		picker.actions[key] = action
	}

	// Add standard quit action if not overridden
	if _, exists := picker.actions["q"]; !exists {
		picker.actions["q"] = Action{
			Key:         "q",
			Description: "quit",
			Handler: func(_ table.Row) (tea.Cmd, *ActionResult) {
				return tea.Quit, &ActionResult{
					Quit:      true,
					Cancelled: true,
				}
			},
		}
	}

	// Add ctrl+c as quit
	if _, exists := picker.actions[keyCtrlC]; !exists {
		picker.actions[keyCtrlC] = Action{
			Key:         keyCtrlC,
			Description: "",
			Handler: func(_ table.Row) (tea.Cmd, *ActionResult) {
				return tea.Quit, &ActionResult{
					Quit:      true,
					Cancelled: true,
				}
			},
		}
	}

	// Generate help text if not provided
	if config.HelpText != "" {
		picker.helpText = config.HelpText
	} else {
		picker.helpText = picker.generateHelpText(len(config.Rows))
	}

	return picker
}

// WithAction adds a key action to the picker.
func (p *TablePicker) WithAction(
	key, description string,
	handler ActionHandler,
) *TablePicker {
	p.actions[key] = Action{
		Key:         key,
		Description: description,
		Handler:     handler,
	}
	// Regenerate help text
	p.helpText = p.generateHelpText(
		len(p.table.Rows()),
	)

	return p
}

// WithProjectPath sets the project path for file operations.
func (p *TablePicker) WithProjectPath(
	path string,
) *TablePicker {
	p.projectPath = path
	// Regenerate help text
	p.helpText = p.generateHelpText(
		len(p.table.Rows()),
	)

	return p
}

// generateHelpText generates help text from registered actions.
func (p *TablePicker) generateHelpText(
	rowCount int,
) string {
	var parts []string

	// Always add navigation first
	parts = append(parts, "↑/↓/j/k: navigate")

	// Sort action keys for consistent ordering
	keys := make([]string, 0, len(p.actions))
	for key := range p.actions {
		// Skip ctrl+c as it's implied by q
		if key == keyCtrlC {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Add action descriptions
	for _, key := range keys {
		action := p.actions[key]
		if action.Description != "" {
			parts = append(
				parts,
				fmt.Sprintf(
					"%s: %s",
					key,
					action.Description,
				),
			)
		}
	}

	helpLine := strings.Join(parts, " | ")

	// Add footer
	var footerParts []string
	footerParts = append(
		footerParts,
		fmt.Sprintf("showing: %d", rowCount),
	)
	if p.projectPath != "" {
		footerParts = append(
			footerParts,
			fmt.Sprintf(
				"project: %s",
				p.projectPath,
			),
		)
	}
	if p.footerExtra != "" {
		footerParts = append(
			footerParts,
			p.footerExtra,
		)
	}

	return helpLine + "\n" + strings.Join(
		footerParts,
		" | ",
	)
}

// generateMinimalFooter generates minimal footer with item count, project path, and help hint.
func (p *TablePicker) generateMinimalFooter(
	rowCount int,
) string {
	var parts []string
	parts = append(
		parts,
		fmt.Sprintf("showing: %d", rowCount),
	)
	if p.projectPath != "" {
		parts = append(
			parts,
			fmt.Sprintf(
				"project: %s",
				p.projectPath,
			),
		)
	}
	if p.footerExtra != "" {
		parts = append(parts, p.footerExtra)
	}
	parts = append(parts, "?: help")

	return strings.Join(parts, " | ")
}

// UpdateHelpText regenerates the help text with the current row count.
func (p *TablePicker) UpdateHelpText() {
	p.helpText = p.generateHelpText(
		len(p.table.Rows()),
	)
}

// SetRows updates the table rows and regenerates help text.
func (p *TablePicker) SetRows(rows []table.Row) {
	p.table.SetRows(rows)
	p.UpdateHelpText()
}

// SetFooterExtra sets additional footer text.
func (p *TablePicker) SetFooterExtra(
	extra string,
) {
	p.footerExtra = extra
	p.UpdateHelpText()
}

// Init implements tea.Model.
func (p *TablePicker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (p *TablePicker) Update(
	msg tea.Msg,
) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		keyStr := keyMsg.String()

		// Handle help toggle
		if keyStr == "?" {
			p.showHelp = !p.showHelp

			return p, nil
		}

		// Auto-hide help on navigation keys
		if keyStr == "up" || keyStr == "down" ||
			keyStr == "j" ||
			keyStr == "k" {
			p.showHelp = false
		}

		// Check if we have a handler for this key
		if action, exists := p.actions[keyStr]; exists {
			return p.handleAction(action)
		}
	}

	// Update table with key events (handles navigation)
	p.table, cmd = p.table.Update(msg)

	return p, cmd
}

// handleAction executes an action handler and processes the result.
func (p *TablePicker) handleAction(
	action Action,
) (tea.Model, tea.Cmd) {
	// Get selected row
	row := p.getSelectedRow()

	// Execute handler
	teaCmd, result := action.Handler(row)
	if result != nil {
		p.result = result
		if result.Quit {
			p.quitting = true
		}
		if result.Error != nil {
			p.err = result.Error
		}
	}

	if teaCmd != nil {
		return p, teaCmd
	}

	// If result indicated quit but no command, return quit
	if result != nil && result.Quit {
		return p, tea.Quit
	}

	return p, nil
}

// View implements tea.Model.
func (p *TablePicker) View() string {
	if p.quitting && p.result != nil {
		return p.renderQuitView()
	}

	// Display table with footer (minimal or full help based on showHelp state)
	var footer string
	if p.showHelp {
		footer = p.helpText
	} else {
		footer = p.generateMinimalFooter(len(p.table.Rows()))
	}

	view := p.table.View() + "\n" + footer + "\n"
	if p.err != nil {
		view += fmt.Sprintf(
			"\nError: %v\n",
			p.err,
		)
	}

	return view
}

// renderQuitView renders the view when quitting.
func (p *TablePicker) renderQuitView() string {
	if p.result == nil {
		return "Cancelled.\n"
	}

	if p.result.ArchiveRequested &&
		p.result.ID != "" {
		return fmt.Sprintf(
			"Archiving: %s\n",
			p.result.ID,
		)
	}

	if p.result.Copied && p.result.Error == nil {
		return fmt.Sprintf(
			"✓ Copied: %s\n",
			p.result.ID,
		)
	}

	if p.result.Error != nil {
		if p.result.ID != "" {
			return fmt.Sprintf(
				"Copied: %s\nError: %v\n",
				p.result.ID,
				p.result.Error,
			)
		}

		return fmt.Sprintf(
			"Error: %v\n",
			p.result.Error,
		)
	}

	if p.result.Cancelled {
		return "Cancelled.\n"
	}

	return ""
}

// getSelectedRow returns the currently selected row.
func (p *TablePicker) getSelectedRow() table.Row {
	cursor := p.table.Cursor()
	rows := p.table.Rows()

	if cursor < 0 || cursor >= len(rows) {
		return nil
	}

	return rows[cursor]
}

// Result returns the action result after the picker has quit.
func (p *TablePicker) Result() *ActionResult {
	return p.result
}

// Table returns the underlying table model for advanced use cases.
func (p *TablePicker) Table() *table.Model {
	return &p.table
}

// Run runs the TablePicker and returns the result.
func (p *TablePicker) Run() (*ActionResult, error) {
	prog := tea.NewProgram(p)
	finalModel, err := prog.Run()
	if err != nil {
		return nil, fmt.Errorf(
			"error running interactive mode: %w",
			err,
		)
	}

	if fm, ok := finalModel.(*TablePicker); ok {
		return fm.Result(), nil
	}

	return nil, nil
}
