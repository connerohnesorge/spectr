// Package list provides interactive TUI components for displaying and
// filtering changes and specifications.
//
// Memory Strategy:
// The interactive model maintains two row slices to optimize memory usage
// during filtering operations:
//   - allRows: master list of all rows (persisted for the lifetime of the
//     model)
//   - filteredRows: reusable buffer for filtered results (reduces GC
//     pressure)
//
// The filteredRows buffer is pre-allocated and reused on every filter
// operation (triggered by each keystroke during search). This avoids repeated
// allocations and reduces garbage collection overhead during interactive use.
//
// Trade-off: We store one extra slice header (24 bytes + capacity) to
// eliminate N temporary allocations during search, where N is the number of
// keystrokes. For typical use cases (10-100 items), this adds <1KB overhead
// while significantly reducing GC pressure during interactive filtering.

//nolint:revive // file-length-limit: interactive functions logically grouped
package list

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	// Item type string constants
	itemTypeAll    = "all"
	itemTypeChange = "change"
	itemTypeSpec   = "spec"

	// Column title constants
	columnTitleID           = "ID"
	columnTitleTitle        = "Title"
	columnTitleType         = "Type"
	columnTitleDeltas       = "Deltas"
	columnTitleTasks        = "Tasks"
	columnTitleDetails      = "Details"
	columnTitleRequirements = "Requirements"

	// Text input settings
	searchInputCharLimit = 50
	searchInputWidth     = 30

	// Error message format
	errInteractiveModeFormat = "error running interactive mode: %w"
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
	minimalFooter    string
	showHelp         bool
	itemType         string    // "spec", "change", or "all"
	projectPath      string    // root directory of the project
	allItems         ItemList  // all items when in unified mode
	filterType       *ItemType // current filter in unified mode (nil = all)
	searchMode       bool      // whether search mode is active
	searchQuery      string    // current search query
	searchInput      textinput.Model
	allRows          []table.Row // stores all rows for filtering
	filteredRows     []table.Row // pre-allocated buffer for GC
	selectionMode    bool        // true: Enter selects without copying
	confirmDelete    bool        // whether we're in delete confirmation mode
	deleteTarget     string      // the spec ID being deleted
}

// Init initializes the model
func (interactiveModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m interactiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		// Handle delete confirmation mode input
		if m.confirmDelete {
			var handled bool
			m, cmd, handled = m.handleDeleteConfirmation(typedMsg)
			if handled {
				return m, cmd
			}
		}

		// Handle search mode input
		if m.searchMode {
			var handled bool
			m, cmd, handled = m.handleSearchModeInput(typedMsg)
			if handled {
				return m, cmd
			}
		}

		switch typedMsg.String() {
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
			if m.itemType == itemTypeAll {
				m = m.handleToggleFilter()

				return m, nil
			}

		case "a":
			return m.handleArchive()

		case "d":
			return m.handleDelete()

		case "/":
			m = m.toggleSearchMode()

			return m, nil

		case "?":
			// Toggle help display
			m.showHelp = !m.showHelp

			return m, nil

		case "up", "down", "j", "k":
			// Auto-hide help on navigation keys
			m.showHelp = false
		}

	case editorFinishedMsg:
		if typedMsg.err != nil {
			m.err = fmt.Errorf("editor error: %w", typedMsg.err)
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
// or selecting in selection mode
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

	// In selection mode, just select without copying to clipboard
	if m.selectionMode {
		return m
	}

	// Otherwise, copy to clipboard
	m.copied = true
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
	var editItemType string

	// Determine item type and ID based on mode
	switch m.itemType {
	case itemTypeAll:
		// In unified mode, need to check the item type
		itemID = row[0]
		itemTypeStr := row[1] // Type is second column in unified mode
		if itemTypeStr == "SPEC" {
			editItemType = itemTypeSpec
		} else {
			editItemType = itemTypeChange
		}
	case itemTypeSpec:
		// In spec-only mode
		itemID = row[0]
		editItemType = itemTypeSpec
	case itemTypeChange:
		// In change-only mode
		itemID = row[0]
		editItemType = itemTypeChange
	default:
		// Unknown mode, no editing allowed
		return m, nil
	}

	// Check if EDITOR is set
	editor := os.Getenv("EDITOR")
	if editor == "" {
		m.err = errors.New("EDITOR environment variable not set")

		return m, nil
	}

	// Construct file path based on type
	filePath := m.getEditFilePath(itemID, editItemType)

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

// handleToggleFilter toggles between showing all items,
// only changes, and only specs
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
	case itemTypeSpec:
		// Can't archive specs
		return m, nil
	case itemTypeChange:
		// In change mode, all items are changes
		m.selectedID = row[0]
		m.archiveRequested = true

		return m, tea.Quit
	case itemTypeAll:
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

// handleDelete handles the 'd' key press for deleting a spec
func (m interactiveModel) handleDelete() (interactiveModel, tea.Cmd) {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.table.Rows()) {
		return m, nil
	}

	row := m.table.Rows()[cursor]
	if len(row) == 0 {
		return m, nil
	}

	// Determine if item is a spec based on mode
	switch m.itemType {
	case itemTypeChange:
		// Can't delete changes, show message
		m.err = errors.New("cannot delete changes; use archive instead")

		return m, nil
	case itemTypeSpec:
		// In spec mode, all items are specs
		m.deleteTarget = row[0]
		m.confirmDelete = true

		return m, nil
	case itemTypeAll:
		// In unified mode, check the type column
		if len(row) > 1 && row[1] == "SPEC" {
			m.deleteTarget = row[0]
			m.confirmDelete = true

			return m, nil
		}
		// Not a spec, show message
		m.err = errors.New("cannot delete changes; use archive instead")

		return m, nil
	}

	return m, nil
}

// handleDeleteConfirmation handles keyboard input when in delete confirmation mode
// Returns (model, cmd, handled) - handled is true if the input was processed
func (m interactiveModel) handleDeleteConfirmation(
	keyMsg tea.KeyMsg,
) (interactiveModel, tea.Cmd, bool) {
	switch keyMsg.String() {
	case "y", "Y":
		// Perform deletion
		err := m.deleteSpecFolder(m.deleteTarget)
		if err != nil {
			m.err = err
			m.confirmDelete = false
			m.deleteTarget = ""

			return m, nil, true
		}

		// Deletion successful - rebuild table
		m = m.removeDeletedSpec()
		m.err = nil // Clear any previous errors
		m.confirmDelete = false
		m.deleteTarget = ""

		return m, nil, true
	case "n", "N", "esc":
		// Cancel deletion
		m.err = errors.New("cancelled")
		m.confirmDelete = false
		m.deleteTarget = ""

		return m, nil, true
	default:
		// Any other key cancels
		m.err = errors.New("cancelled")
		m.confirmDelete = false
		m.deleteTarget = ""

		return m, nil, true
	}
}

// deleteSpecFolder deletes the spec folder from disk
func (m interactiveModel) deleteSpecFolder(specID string) error {
	// Validate specID against path traversal attacks
	if strings.Contains(specID, "..") || strings.Contains(specID, "/") ||
		strings.Contains(specID, "\\") {
		return fmt.Errorf("invalid spec ID: %s", specID)
	}

	path := filepath.Join(m.projectPath, "spectr", "specs", specID)

	// Use RemoveAll directly - it's idempotent and avoids race conditions
	// between existence check and removal
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete spec folder: %w", err)
	}

	return nil
}

// removeDeletedSpec removes the deleted spec from the table and rebuilds
func (m interactiveModel) removeDeletedSpec() interactiveModel {
	cursor := m.table.Cursor()

	// Remove from allRows
	newRows := make([]table.Row, 0, len(m.allRows)-1)
	for _, row := range m.allRows {
		if len(row) > 0 && row[0] != m.deleteTarget {
			newRows = append(newRows, row)
		}
	}
	m.allRows = newRows

	// If in unified mode, also update allItems
	if m.itemType == itemTypeAll {
		newItems := make(ItemList, 0, len(m.allItems)-1)
		for _, item := range m.allItems {
			if item.ID() != m.deleteTarget {
				newItems = append(newItems, item)
			}
		}
		m.allItems = newItems

		// Rebuild the unified table
		m = rebuildUnifiedTable(m)
	} else {
		// Update help text footer for spec-only mode
		m.minimalFooter = fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(m.allRows),
			m.projectPath,
		)
	}

	// Re-apply search filter if active to maintain filtered view
	m = m.applyFilter()

	// Update footer to reflect actual visible count after filtering
	visibleCount := len(m.table.Rows())
	m.minimalFooter = fmt.Sprintf(
		"showing: %d | project: %s | ?: help",
		visibleCount,
		m.projectPath,
	)

	// Adjust cursor if needed
	if visibleCount > 0 {
		if cursor >= visibleCount {
			m.table.SetCursor(visibleCount - 1)
		}
	} else {
		m.table.SetCursor(0)
	}

	return m
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
		{Title: columnTitleID, Width: unifiedIDWidth},
		{Title: columnTitleType, Width: unifiedTypeWidth},
		{Title: columnTitleTitle, Width: unifiedTitleWidth},
		{Title: columnTitleDetails, Width: unifiedDetailsWidth},
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

	// Update help text and minimal footer
	filterDesc := "all"
	if m.filterType != nil {
		filterDesc = m.filterType.String() + "s"
	}
	m.helpText = fmt.Sprintf(
		"↑/↓/j/k: navigate | Enter: copy ID | e: edit | "+
			"a: archive | d: delete (specs) | t: filter (%s) | /: search | q: quit",
		filterDesc,
	)
	m.minimalFooter = fmt.Sprintf(
		"showing: %d | project: %s | ?: help",
		len(rows),
		m.projectPath,
	)
	// Preserve showHelp state (no reset needed here)

	return m
}

// applyFilter filters the table rows based on the search query.
// Memory optimization: reuses pre-allocated filteredRows buffer to reduce
// GC pressure during interactive search (which triggers on every keystroke).
func (m interactiveModel) applyFilter() interactiveModel {
	if len(m.allRows) == 0 {
		return m
	}

	query := strings.ToLower(m.searchQuery)

	if query == "" {
		// No filter - show all rows directly
		m.table.SetRows(m.allRows)
	} else {
		// Reuse the filteredRows buffer to avoid allocations on every keystroke
		m.filteredRows = m.filteredRows[:0] // reset length but keep capacity
		for _, row := range m.allRows {
			if rowMatchesQuery(row, query) {
				m.filteredRows = append(m.filteredRows, row)
			}
		}
		m.table.SetRows(m.filteredRows)
	}

	// Ensure cursor is within bounds after filtering
	rowCount := len(m.table.Rows())
	if rowCount > 0 {
		cursor := m.table.Cursor()
		if cursor >= rowCount {
			m.table.SetCursor(rowCount - 1)
		}
	} else {
		m.table.SetCursor(0)
	}

	return m
}

// newTextInput creates a new text input for search
func newTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = searchInputCharLimit
	ti.Width = searchInputWidth

	return ti
}

// editorFinishedMsg is sent when the editor finishes
type editorFinishedMsg struct {
	err error
}

// rowMatchesQuery checks if any column in the row contains the query
func rowMatchesQuery(row table.Row, query string) bool {
	for _, col := range row {
		if strings.Contains(strings.ToLower(col), query) {
			return true
		}
	}

	return false
}

// handleSearchModeInput handles keyboard input when in search mode
// Returns (model, cmd, handled) - handled is true if the input was processed
func (m interactiveModel) handleSearchModeInput(
	keyMsg tea.KeyMsg,
) (interactiveModel, tea.Cmd, bool) {
	var cmd tea.Cmd

	//nolint:exhaustive // Only handling specific keys
	switch keyMsg.Type {
	case tea.KeyEsc:
		// Exit search mode, clear query and restore all rows
		m.searchMode = false
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m = m.applyFilter()

		return m, nil, true
	case tea.KeyEnter:
		// Exit search mode but keep filter applied
		m.searchMode = false

		return m, nil, true
	default:
		// Update text input and filter
		m.searchInput, cmd = m.searchInput.Update(keyMsg)
		m.searchQuery = m.searchInput.Value()
		m = m.applyFilter()

		return m, cmd, true
	}
}

// toggleSearchMode toggles search mode on/off
func (m interactiveModel) toggleSearchMode() interactiveModel {
	if m.searchMode {
		m.searchMode = false
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m = m.applyFilter()
	} else {
		m.searchMode = true
		m.searchInput.Focus()
	}

	return m
}

// getEditFilePath returns the file path to edit based on item type
func (m interactiveModel) getEditFilePath(
	itemID string,
	itemType string,
) string {
	if itemType == itemTypeSpec {
		return fmt.Sprintf(
			"%s/spectr/specs/%s/spec.md",
			m.projectPath, itemID,
		)
	}

	return fmt.Sprintf(
		"%s/spectr/changes/%s/proposal.md",
		m.projectPath, itemID,
	)
}

// View renders the model
func (m interactiveModel) View() string {
	if m.quitting {
		if m.archiveRequested && m.selectedID != "" {
			return fmt.Sprintf("Archiving: %s\n", m.selectedID)
		}

		// In selection mode, just show selected ID without clipboard message
		if m.selectionMode && m.selectedID != "" {
			return fmt.Sprintf("Selected: %s\n", m.selectedID)
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

	// Display delete confirmation prompt if active
	if m.confirmDelete {
		return fmt.Sprintf(
			"%s\n\nDelete spec '%s'? This will remove the entire folder. (y/N)\n",
			m.table.View(),
			m.deleteTarget,
		)
	}

	// Display search input if search mode is active
	var view string
	if m.searchMode {
		view = fmt.Sprintf("Search: %s\n\n", m.searchInput.View())
	}

	// Choose which footer to display based on showHelp state
	footer := m.minimalFooter
	if m.showHelp {
		footer = m.helpText
	}

	view += m.table.View() + "\n" + footer + "\n"

	// Display error message if present, but keep TUI active
	if m.err != nil {
		view += fmt.Sprintf("\nError: %v\n", m.err)
	}

	return view
}

// RunInteractiveChanges runs the interactive table for changes.
// Returns the change ID if archive was requested, empty string
// otherwise.
func RunInteractiveChanges(
	changes []ChangeInfo,
	projectPath string,
) (string, error) {
	if len(changes) == 0 {
		return "", nil
	}

	columns := []table.Column{
		{Title: columnTitleID, Width: changeIDWidth},
		{Title: columnTitleTitle, Width: changeTitleWidth},
		{Title: columnTitleDeltas, Width: changeDeltaWidth},
		{Title: columnTitleTasks, Width: changeTasksWidth},
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
		itemType:    itemTypeChange,
		projectPath: projectPath,
		searchInput: newTextInput(),
		allRows:     rows,
		helpText: "↑/↓/j/k: navigate | Enter: copy ID | e: edit | " +
			"a: archive | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf(errInteractiveModeFormat, err)
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

// RunInteractiveArchive runs the interactive table for
// archive selection. Returns the selected change ID or empty
// string if cancelled
func RunInteractiveArchive(
	changes []ChangeInfo,
	projectPath string,
) (string, error) {
	if len(changes) == 0 {
		return "", nil
	}

	columns := []table.Column{
		{Title: columnTitleID, Width: changeIDWidth},
		{Title: columnTitleTitle, Width: changeTitleWidth},
		{Title: columnTitleDeltas, Width: changeDeltaWidth},
		{Title: columnTitleTasks, Width: changeTasksWidth},
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
		table:         t,
		itemType:      itemTypeChange,
		projectPath:   projectPath,
		searchInput:   newTextInput(),
		allRows:       rows,
		selectionMode: true, // Enter selects without copying
		helpText:      "↑/↓/j/k: navigate | Enter: select | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf(errInteractiveModeFormat, err)
	}

	// Check if there was an error during execution
	if fm, ok := finalModel.(interactiveModel); ok {
		if fm.err != nil {
			// Don't return error, just warn
			fmt.Fprintf(
				os.Stderr,
				"Warning: operation failed: %v\n",
				fm.err,
			)
		}

		// Return selected ID if one was selected
		if fm.selectedID != "" {
			return fm.selectedID, nil
		}
	}

	// Cancelled
	return "", nil
}

// RunInteractiveSpecs runs the interactive table for specs
func RunInteractiveSpecs(specs []SpecInfo, projectPath string) error {
	if len(specs) == 0 {
		return nil
	}

	columns := []table.Column{
		{Title: columnTitleID, Width: specIDWidth},
		{Title: columnTitleTitle, Width: specTitleWidth},
		{Title: columnTitleRequirements, Width: specRequirementsWidth},
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
		itemType:    itemTypeSpec,
		projectPath: projectPath,
		searchInput: newTextInput(),
		allRows:     rows,
		helpText: "↑/↓/j/k: navigate | Enter: copy ID | e: edit | " +
			"d: delete | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(specs),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf(errInteractiveModeFormat, err)
	}

	// Check if there was an error during execution
	fm, ok := finalModel.(interactiveModel)
	if ok && fm.err != nil {
		// Don't return error, just warn - clipboard failure
		// shouldn't stop the command.
		fmt.Fprintf(
			os.Stderr,
			"Warning: clipboard operation failed: %v\n",
			fm.err,
		)
	}

	return nil
}

// RunInteractiveAll runs the interactive table for all items
// (changes and specs)
func RunInteractiveAll(items ItemList, projectPath string) error {
	if len(items) == 0 {
		return nil
	}

	// Build initial table with all items
	columns := []table.Column{
		{Title: columnTitleID, Width: unifiedIDWidth},
		{Title: columnTitleType, Width: unifiedTypeWidth},
		{Title: columnTitleTitle, Width: unifiedTitleWidth},
		{Title: columnTitleDetails, Width: unifiedDetailsWidth},
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
		itemType:    itemTypeAll,
		projectPath: projectPath,
		allItems:    items,
		filterType:  nil, // Start with all items visible
		searchInput: newTextInput(),
		allRows:     rows,
		helpText: "↑/↓/j/k: navigate | Enter: copy ID | e: edit | " +
			"a: archive | d: delete (specs) | t: filter (all) | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf(errInteractiveModeFormat, err)
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
