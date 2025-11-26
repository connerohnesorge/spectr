//nolint:revive // file-length-limit - interactive functions logically grouped
package list

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/connerohnesorge/spectr/internal/tui"
)

const (
	// Table column widths for changes view
	changeIDWidth    = 30
	changeTitleWidth = 40
	changeDeltaWidth = 10
	changeTasksWidth = 15

	// Table column widths for specs view
	specIDWidth           = 35
	specTitleWidth        = 45
	specRequirementsWidth = 15

	// Table column widths for unified view
	unifiedIDWidth      = 30
	unifiedTypeWidth    = 8
	unifiedTitleWidth   = 40
	unifiedDetailsWidth = 20

	// Truncation settings
	changeTitleTruncate  = 38
	specTitleTruncate    = 43
	unifiedTitleTruncate = 38

	// Table height
	tableHeight = 10
)

// interactiveModel represents the bubbletea model for interactive table
type interactiveModel struct {
	table            table.Model
	selectedID       string
	copied           bool
	quitting         bool
	archiveRequested bool
	err              error
	helpText         string
	itemType         string    // "spec", "change", or "all"
	projectPath      string    // root directory of the project
	allItems         ItemList  // all items when in unified mode
	filterType       *ItemType // current filter when in unified mode (nil = show all)
	searchMode       bool      // whether search mode is active
	searchQuery      string    // current search query
	searchInput      textinput.Model
	allRows          []table.Row // stores all rows for filtering
}

// Init initializes the model
func (interactiveModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m interactiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input
		if m.searchMode {
			//nolint:exhaustive // Only handling specific keys, default handles the rest
			switch msg.Type {
			case tea.KeyEsc:
				// Exit search mode, clear query and restore all rows
				m.searchMode = false
				m.searchQuery = ""
				m.searchInput.SetValue("")
				m = m.applyFilter()

				return m, nil
			case tea.KeyEnter:
				// Exit search mode but keep filter applied
				m.searchMode = false

				return m, nil
			default:
				// Update text input and filter
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.searchQuery = m.searchInput.Value()
				m = m.applyFilter()

				return m, cmd
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			m = m.handleEnter()

			return m, tea.Quit

		case "e":
			return m.handleEdit()

		case "t":
			// Toggle filter type in unified mode
			if m.itemType == "all" {
				m = m.handleToggleFilter()

				return m, nil
			}

		case "a":
			return m.handleArchive()

		case "/":
			// Toggle search mode
			if m.searchMode {
				// Exit search mode if already active
				m.searchMode = false
				m.searchQuery = ""
				m.searchInput.SetValue("")
				m = m.applyFilter()
			} else {
				// Enter search mode
				m.searchMode = true
				m.searchInput.Focus()
			}

			return m, nil

		case "esc":
			// Exit search mode if active
			if m.searchMode {
				m.searchMode = false
				m.searchQuery = ""
				m.searchInput.SetValue("")
				m = m.applyFilter()

				return m, nil
			}
		}

	case editorFinishedMsg:
		if msg.err != nil {
			m.err = fmt.Errorf("editor error: %w", msg.err)
			m.quitting = true

			return m, tea.Quit
		}
		// Continue in TUI on success
		return m, nil
	}

	// Update table with key events
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

// handleEnter handles the enter key press for copying selected ID
func (m interactiveModel) handleEnter() interactiveModel {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.table.Rows()) {
		return m
	}

	row := m.table.Rows()[cursor]
	if len(row) == 0 {
		return m
	}

	// ID is in first column for all modes
	m.selectedID = row[0]
	m.copied = true

	// Copy to clipboard using shared helper
	err := tui.CopyToClipboard(m.selectedID)
	if err != nil {
		m.err = err
	}

	return m
}

// handleEdit handles the 'e' key press for opening file in editor
func (m interactiveModel) handleEdit() (interactiveModel, tea.Cmd) {
	// Get the selected row
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.table.Rows()) {
		return m, nil
	}

	row := m.table.Rows()[cursor]
	if len(row) == 0 {
		return m, nil
	}

	var itemID string
	var isSpec bool

	// Determine item type and ID based on mode
	switch m.itemType {
	case "all":
		// In unified mode, need to check the item type
		itemID = row[0]
		itemTypeStr := row[1] // Type is second column in unified mode
		isSpec = itemTypeStr == "SPEC"
	case "spec":
		// In spec-only mode
		itemID = row[0]
		isSpec = true
	case "change":
		// In change-only mode
		itemID = row[0]
		isSpec = false
	default:
		// Unknown mode, no editing allowed
		return m, nil
	}

	// Check if EDITOR is set
	editor := os.Getenv("EDITOR")
	if editor == "" {
		m.err = fmt.Errorf("EDITOR environment variable not set")

		return m, nil
	}

	// Construct file path based on type
	var filePath string
	if isSpec {
		filePath = fmt.Sprintf("%s/spectr/specs/%s/spec.md", m.projectPath, itemID)
	} else {
		filePath = fmt.Sprintf("%s/spectr/changes/%s/proposal.md", m.projectPath, itemID)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		m.err = fmt.Errorf("file not found: %s", filePath)

		return m, nil
	}

	// Launch editor - use tea.ExecProcess to handle editor lifecycle
	c := exec.Command(editor, filePath) //nolint:gosec

	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err}
	})
}

// handleToggleFilter toggles between showing all items, only changes, and only specs
func (m interactiveModel) handleToggleFilter() interactiveModel {
	// Cycle through filter states: all -> changes -> specs -> all
	if m.filterType == nil {
		// Currently showing all, switch to changes only
		changeType := ItemTypeChange
		m.filterType = &changeType
	} else {
		switch *m.filterType {
		case ItemTypeChange:
			// Currently showing changes, switch to specs only
			specType := ItemTypeSpec
			m.filterType = &specType
		case ItemTypeSpec:
			// Currently showing specs, switch back to all
			m.filterType = nil
		}
	}

	// Rebuild the table with the new filter
	m = rebuildUnifiedTable(m)

	return m
}

// handleArchive handles the 'a' key press for archiving a change
func (m interactiveModel) handleArchive() (interactiveModel, tea.Cmd) {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.table.Rows()) {
		return m, nil
	}

	row := m.table.Rows()[cursor]
	if len(row) == 0 {
		return m, nil
	}

	// Determine if item is a change based on mode
	switch m.itemType {
	case "spec":
		// Can't archive specs
		return m, nil
	case "change":
		// In change mode, all items are changes
		m.selectedID = row[0]
		m.archiveRequested = true

		return m, tea.Quit
	case "all":
		// In unified mode, check the type column
		if len(row) > 1 && row[1] == "CHANGE" {
			m.selectedID = row[0]
			m.archiveRequested = true

			return m, tea.Quit
		}
		// Not a change, do nothing
		return m, nil
	}

	return m, nil
}

// rebuildUnifiedTable rebuilds the table based on current filter
func rebuildUnifiedTable(m interactiveModel) interactiveModel {
	var items ItemList
	if m.filterType == nil {
		items = m.allItems
	} else {
		items = m.allItems.FilterByType(*m.filterType)
	}

	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}

	rows := make([]table.Row, len(items))
	for i, item := range items {
		var typeStr, details string
		switch item.Type {
		case ItemTypeChange:
			typeStr = "CHANGE"
			if item.Change != nil {
				details = fmt.Sprintf("Tasks: %d/%d 🔺 %d",
					item.Change.TaskStatus.Completed,
					item.Change.TaskStatus.Total,
					item.Change.DeltaCount,
				)
			}
		case ItemTypeSpec:
			typeStr = "SPEC"
			if item.Spec != nil {
				details = fmt.Sprintf("Reqs: %d", item.Spec.RequirementCount)
			}
		}

		rows[i] = table.Row{
			item.ID(),
			typeStr,
			tui.TruncateString(item.Title(), unifiedTitleTruncate),
			details,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)
	tui.ApplyTableStyles(&t)

	m.table = t
	m.allRows = rows // Update allRows for search filtering

	// Update help text
	filterDesc := "all"
	if m.filterType != nil {
		filterDesc = m.filterType.String() + "s"
	}
	m.helpText = fmt.Sprintf(
		"↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | t: filter (%s) | /: search | q: quit\nshowing: %d | project: %s",
		filterDesc,
		len(rows),
		m.projectPath,
	)

	return m
}

// applyFilter filters the table rows based on the search query
func (m interactiveModel) applyFilter() interactiveModel {
	if len(m.allRows) == 0 {
		return m
	}

	var filteredRows []table.Row
	query := strings.ToLower(m.searchQuery)

	if query == "" {
		filteredRows = m.allRows
	} else {
		for _, row := range m.allRows {
			// Match against ID (first column) and title (second or third column depending on mode)
			if len(row) > 0 {
				id := strings.ToLower(row[0])
				var title string
				if m.itemType == "all" && len(row) > 2 {
					title = strings.ToLower(row[2]) // Title is third column in unified mode
				} else if len(row) > 1 {
					title = strings.ToLower(row[1]) // Title is second column otherwise
				}

				if strings.Contains(id, query) || strings.Contains(title, query) {
					filteredRows = append(filteredRows, row)
				}
			}
		}
	}

	m.table.SetRows(filteredRows)

	return m
}

// newTextInput creates a new text input for search
func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 50
	ti.Width = 30

	return ti
}

// editorFinishedMsg is sent when the editor finishes
type editorFinishedMsg struct {
	err error
}

// View renders the model
func (m interactiveModel) View() string {
	if m.quitting {
		if m.archiveRequested && m.selectedID != "" {
			return fmt.Sprintf("Archiving: %s\n", m.selectedID)
		}

		if m.copied && m.err == nil {
			return fmt.Sprintf("✓ Copied: %s\n", m.selectedID)
		} else if m.err != nil {
			return fmt.Sprintf(
				"Copied: %s\nError: %v\n",
				m.selectedID,
				m.err,
			)
		}

		return "Cancelled.\n"
	}

	// Display search input if search mode is active
	var view string
	if m.searchMode {
		view = fmt.Sprintf("Search: %s\n\n", m.searchInput.View())
	}

	view += m.table.View() + "\n" + m.helpText + "\n"

	// Display error message if present, but keep TUI active
	if m.err != nil {
		view += fmt.Sprintf("\nError: %v\n", m.err)
	}

	return view
}

// RunInteractiveChanges runs the interactive table for changes.
// Returns the change ID if archive was requested, empty string otherwise.
func RunInteractiveChanges(changes []ChangeInfo, projectPath string) (string, error) {
	if len(changes) == 0 {
		return "", nil
	}

	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}

	rows := make([]table.Row, len(changes))
	for i, change := range changes {
		tasksStatus := fmt.Sprintf("%d/%d",
			change.TaskStatus.Completed,
			change.TaskStatus.Total)

		rows[i] = table.Row{
			change.ID,
			tui.TruncateString(change.Title, changeTitleTruncate),
			fmt.Sprintf("%d", change.DeltaCount),
			tasksStatus,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	tui.ApplyTableStyles(&t)

	m := interactiveModel{
		table:       t,
		itemType:    "change",
		projectPath: projectPath,
		searchInput: newTextInput(),
		allRows:     rows,
		helpText: fmt.Sprintf(
			"↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | /: search | q: quit\nshowing: %d | project: %s",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("error running interactive mode: %w", err)
	}

	// Check if there was an error during execution
	if fm, ok := finalModel.(interactiveModel); ok {
		if fm.err != nil {
			// Don't return error, just warn - clipboard failure shouldn't
			// stop the command
			fmt.Fprintf(
				os.Stderr,
				"Warning: clipboard operation failed: %v\n",
				fm.err,
			)
		}

		// Return archive ID if archive was requested
		if fm.archiveRequested && fm.selectedID != "" {
			return fm.selectedID, nil
		}
	}

	return "", nil
}

// RunInteractiveArchive runs the interactive table for archive selection
// Returns the selected change ID or empty string if cancelled
func RunInteractiveArchive(changes []ChangeInfo, projectPath string) (string, error) {
	if len(changes) == 0 {
		return "", nil
	}

	columns := []table.Column{
		{Title: "ID", Width: changeIDWidth},
		{Title: "Title", Width: changeTitleWidth},
		{Title: "Deltas", Width: changeDeltaWidth},
		{Title: "Tasks", Width: changeTasksWidth},
	}

	rows := make([]table.Row, len(changes))
	for i, change := range changes {
		tasksStatus := fmt.Sprintf("%d/%d",
			change.TaskStatus.Completed,
			change.TaskStatus.Total)

		rows[i] = table.Row{
			change.ID,
			tui.TruncateString(change.Title, changeTitleTruncate),
			fmt.Sprintf("%d", change.DeltaCount),
			tasksStatus,
		}
	}

	// Create picker with enter action for selection
	picker := tui.NewTablePicker(tui.TableConfig{
		Columns:     columns,
		Rows:        rows,
		Height:      tableHeight,
		ProjectPath: projectPath,
		Actions: map[string]tui.Action{
			"enter": {
				Key:         "enter",
				Description: "select",
				Handler: func(row table.Row) (tea.Cmd, *tui.ActionResult) {
					if len(row) == 0 {
						return nil, nil
					}

					return tea.Quit, &tui.ActionResult{
						ID:   row[0],
						Quit: true,
					}
				},
			},
		},
	})

	result, err := picker.Run()
	if err != nil {
		return "", err
	}

	if result == nil || result.Cancelled {
		return "", nil
	}

	return result.ID, nil
}

// RunInteractiveSpecs runs the interactive table for specs
func RunInteractiveSpecs(specs []SpecInfo, projectPath string) error {
	if len(specs) == 0 {
		return nil
	}

	columns := []table.Column{
		{Title: "ID", Width: specIDWidth},
		{Title: "Title", Width: specTitleWidth},
		{Title: "Requirements", Width: specRequirementsWidth},
	}

	rows := make([]table.Row, len(specs))
	for i, spec := range specs {
		rows[i] = table.Row{
			spec.ID,
			tui.TruncateString(spec.Title, specTitleTruncate),
			fmt.Sprintf("%d", spec.RequirementCount),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
	)

	tui.ApplyTableStyles(&t)

	m := interactiveModel{
		table:       t,
		itemType:    "spec",
		projectPath: projectPath,
		searchInput: newTextInput(),
		allRows:     rows,
		helpText: fmt.Sprintf(
			"↑/↓/j/k: navigate | Enter: copy ID | e: edit | /: search | q: quit\nshowing: %d | project: %s",
			len(specs),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive mode: %w", err)
	}

	// Check if there was an error during execution
	fm, ok := finalModel.(interactiveModel)
	if ok && fm.err != nil {
		// Don't return error, just warn - clipboard failure shouldn't
		// stop the command.
		fmt.Fprintf(
			os.Stderr,
			"Warning: clipboard operation failed: %v\n",
			fm.err,
		)
	}

	return nil
}

// RunInteractiveAll runs the interactive table for all items (changes and specs)
func RunInteractiveAll(items ItemList, projectPath string) error {
	if len(items) == 0 {
		return nil
	}

	// Build initial table with all items
	columns := []table.Column{
		{Title: "ID", Width: unifiedIDWidth},
		{Title: "Type", Width: unifiedTypeWidth},
		{Title: "Title", Width: unifiedTitleWidth},
		{Title: "Details", Width: unifiedDetailsWidth},
	}

	rows := make([]table.Row, len(items))
	for i, item := range items {
		var typeStr, details string
		switch item.Type {
		case ItemTypeChange:
			typeStr = "CHANGE"
			if item.Change != nil {
				details = fmt.Sprintf("Tasks: %d/%d 🔺 %d",
					item.Change.TaskStatus.Completed,
					item.Change.TaskStatus.Total,
					item.Change.DeltaCount,
				)
			}
		case ItemTypeSpec:
			typeStr = "SPEC"
			if item.Spec != nil {
				details = fmt.Sprintf("Reqs: %d", item.Spec.RequirementCount)
			}
		}

		rows[i] = table.Row{
			item.ID(),
			typeStr,
			tui.TruncateString(item.Title(), unifiedTitleTruncate),
			details,
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	tui.ApplyTableStyles(&t)

	m := interactiveModel{
		table:       t,
		itemType:    "all",
		projectPath: projectPath,
		allItems:    items,
		filterType:  nil, // Start with all items visible
		searchInput: newTextInput(),
		allRows:     rows,
		helpText: fmt.Sprintf(
			"↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | t: filter (all) | /: search | q: quit\nshowing: %d | project: %s",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive mode: %w", err)
	}

	// Check if there was an error during execution
	if fm, ok := finalModel.(interactiveModel); ok && fm.err != nil {
		// Don't return error, just warn - clipboard failure shouldn't
		// stop the command
		fmt.Fprintf(
			os.Stderr,
			"Warning: clipboard operation failed: %v\n",
			fm.err,
		)
	}

	return nil
}
