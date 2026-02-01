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
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/connerohnesorge/spectr/internal/specterrs"
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

	// Line number column width
	lineNumberColumnWidth = 3

	// Item type string constants
	itemTypeAll    = "all"
	itemTypeChange = "change"
	itemTypeSpec   = "spec"

	// Type display strings
	typeDisplayChange = "CHANGE"
	typeDisplaySpec   = "SPEC"

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

	// Width breakpoint thresholds for responsive column layout.
	// breakpointFull: all columns shown at default widths
	breakpointFull = 110
	// breakpointMedium: title column narrowed, all columns still visible
	breakpointMedium = 90
	// breakpointNarrow: low-priority columns hidden
	breakpointNarrow = 70
	// breakpointHideTitle: threshold below which Title column is hidden (changes view)
	breakpointHideTitle = 80
)

// LineNumberMode controls how line numbers are displayed in the interactive list.
type LineNumberMode int

const (
	// LineNumberOff - no line numbers displayed
	LineNumberOff LineNumberMode = iota
	// LineNumberRelative - show relative distance from cursor
	LineNumberRelative
	// LineNumberHybrid - cursor shows absolute, others show relative
	LineNumberHybrid
)

// ColumnPriority defines the priority level for table columns.
// Higher priority columns are shown first when space is limited.
type ColumnPriority int

const (
	// ColumnPriorityEssential - always shown (ID column)
	ColumnPriorityEssential ColumnPriority = iota
	// ColumnPriorityHigh - always shown, width may be adjusted (Title, Type)
	ColumnPriorityHigh
	// ColumnPriorityMedium - hidden below narrow breakpoint (Deltas, Requirements)
	ColumnPriorityMedium
	// ColumnPriorityLow - hidden below medium breakpoint (Tasks, Details)
	ColumnPriorityLow
)

// calculateChangesColumns returns the appropriate columns for the changes view
// based on the current terminal width. Column visibility and widths are
// adjusted according to the priority system:
//   - ID: Essential (always shown)
//   - Tasks: High (shown until very narrow widths)
//   - Deltas: Medium (hidden below 70 columns)
//   - Title: Low (hidden below 80 columns to prioritize Tasks)
//
// Column order is always: ID | Title | Deltas | Tasks (when visible)
func calculateChangesColumns(
	width int,
	lineNumberMode LineNumberMode,
) []table.Column {
	// Calculate available width for content (accounting for table borders/padding)
	// Table has approximately 4 chars of padding/borders per column
	const paddingPerColumn = 4

	// Prepend line number column if enabled
	var lineNumCol []table.Column
	if lineNumberMode != LineNumberOff {
		lineNumCol = []table.Column{
			{Title: "Ln", Width: lineNumberColumnWidth},
		}
	}

	switch {
	case width >= breakpointFull:
		// Full width (110+): all 4 columns at default widths
		cols := []table.Column{
			{
				Title: columnTitleID,
				Width: changeIDWidth,
			},
			{
				Title: columnTitleTitle,
				Width: changeTitleWidth,
			},
			{
				Title: columnTitleDeltas,
				Width: changeDeltaWidth,
			},
			{
				Title: columnTitleTasks,
				Width: changeTasksWidth,
			},
		}
		return append(lineNumCol, cols...)

	case width >= breakpointMedium:
		// Medium width (90-109): all 4 columns visible, Title narrowed
		titleWidth := max(
			width-changeIDWidth-changeDeltaWidth-
				changeTasksWidth-(paddingPerColumn*4),
			20,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: changeIDWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
			{
				Title: columnTitleDeltas,
				Width: changeDeltaWidth,
			},
			{
				Title: columnTitleTasks,
				Width: changeTasksWidth,
			},
		}

	case width >= breakpointHideTitle:
		// Width 80-89: all 4 columns, Title very narrow (15-20 chars)
		titleWidth := max(
			width-changeIDWidth-changeDeltaWidth-
				changeTasksWidth-(paddingPerColumn*4),
			15,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: changeIDWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
			{
				Title: columnTitleDeltas,
				Width: changeDeltaWidth,
			},
			{
				Title: columnTitleTasks,
				Width: changeTasksWidth,
			},
		}

	case width >= breakpointNarrow:
		// Width 70-79: 3 columns - ID, Deltas, Tasks (Title hidden)
		return []table.Column{
			{
				Title: columnTitleID,
				Width: changeIDWidth,
			},
			{
				Title: columnTitleDeltas,
				Width: changeDeltaWidth,
			},
			{
				Title: columnTitleTasks,
				Width: changeTasksWidth,
			},
		}

	default:
		// Minimal width (<70): 2 columns - ID, Tasks only
		idWidth := 20
		tasksWidth := max(
			width-idWidth-(paddingPerColumn*2),
			10,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: idWidth,
			},
			{
				Title: columnTitleTasks,
				Width: tasksWidth,
			},
		}
	}
}

// calculateSpecsColumns returns the appropriate columns for the specs view
// based on the current terminal width. Column visibility and widths are
// adjusted according to the priority system:
//   - ID: Essential (always shown)
//   - Title: High (always shown, width adjustable)
//   - Requirements: Medium (width reduced or hidden below 70 columns)
func calculateSpecsColumns(
	width int,
) []table.Column {
	const paddingPerColumn = 4

	switch {
	case width >= breakpointFull:
		// Full width: all columns at default widths
		return []table.Column{
			{
				Title: columnTitleID,
				Width: specIDWidth,
			},
			{
				Title: columnTitleTitle,
				Width: specTitleWidth,
			},
			{
				Title: columnTitleRequirements,
				Width: specRequirementsWidth,
			},
		}

	case width >= breakpointMedium:
		// Medium width: all columns visible, Title narrowed
		titleWidth := max(
			width-specIDWidth-specRequirementsWidth-
				(paddingPerColumn*3),
			25,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: specIDWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
			{
				Title: columnTitleRequirements,
				Width: specRequirementsWidth,
			},
		}

	case width >= breakpointNarrow:
		// Narrow width: Requirements column narrowed significantly
		reqWidth := 8
		titleWidth := max(
			width-specIDWidth-reqWidth-(paddingPerColumn*3),
			20,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: specIDWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
			{
				Title: columnTitleRequirements,
				Width: reqWidth,
			},
		}

	default:
		// Minimal width: hide Requirements, compress ID and Title
		idWidth := 25
		titleWidth := max(
			width-idWidth-(paddingPerColumn*2),
			15,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: idWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
		}
	}
}

// calculateUnifiedColumns returns the appropriate columns for the unified view
// based on the current terminal width. Column visibility and widths are
// adjusted according to the priority system:
//   - ID: Essential (always shown)
//   - Type: High (always shown at fixed 8 width)
//   - Title: High (width adjustable)
//   - Details: Low (hidden below 90 columns)
func calculateUnifiedColumns(
	width int,
) []table.Column {
	const paddingPerColumn = 4

	switch {
	case width >= breakpointFull:
		// Full width: all columns at default widths
		return []table.Column{
			{
				Title: columnTitleID,
				Width: unifiedIDWidth,
			},
			{
				Title: columnTitleType,
				Width: unifiedTypeWidth,
			},
			{
				Title: columnTitleTitle,
				Width: unifiedTitleWidth,
			},
			{
				Title: columnTitleDetails,
				Width: unifiedDetailsWidth,
			},
		}

	case width >= breakpointMedium:
		// Medium width: all columns visible, Title narrowed
		titleWidth := max(
			width-unifiedIDWidth-unifiedTypeWidth-
				unifiedDetailsWidth-(paddingPerColumn*4),
			25,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: unifiedIDWidth,
			},
			{
				Title: columnTitleType,
				Width: unifiedTypeWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
			{
				Title: columnTitleDetails,
				Width: unifiedDetailsWidth,
			},
		}

	case width >= breakpointNarrow:
		// Narrow width: hide Details column
		titleWidth := max(
			width-unifiedIDWidth-unifiedTypeWidth-
				(paddingPerColumn*3),
			20,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: unifiedIDWidth,
			},
			{
				Title: columnTitleType,
				Width: unifiedTypeWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
		}

	default:
		// Minimal width: hide Details, compress ID and Title
		idWidth := 20
		titleWidth := max(
			width-idWidth-unifiedTypeWidth-(paddingPerColumn*3),
			15,
		)

		return []table.Column{
			{
				Title: columnTitleID,
				Width: idWidth,
			},
			{
				Title: columnTitleType,
				Width: unifiedTypeWidth,
			},
			{
				Title: columnTitleTitle,
				Width: titleWidth,
			},
		}
	}
}

// calculateTitleTruncate returns the appropriate title truncation limit
// based on the view type and terminal width. The truncation is set slightly
// below the column width to account for any table rendering overhead.
func calculateTitleTruncate(
	viewType string,
	width int,
) int {
	const truncateBuffer = 2 // Leave 2 chars buffer for clean truncation

	switch viewType {
	case itemTypeChange:
		cols := calculateChangesColumns(width, LineNumberOff)
		// Find Title column (may not be present at narrow widths)
		for _, col := range cols {
			if col.Title == columnTitleTitle {
				return col.Width - truncateBuffer
			}
		}
		// Title not present, return default
		return changeTitleTruncate

	case itemTypeSpec:
		cols := calculateSpecsColumns(width)
		// Title is always the second column
		if len(cols) >= 2 {
			return cols[1].Width - truncateBuffer
		}

		return specTitleTruncate

	case itemTypeAll:
		cols := calculateUnifiedColumns(width)
		// Title is the third column in unified view
		if len(cols) >= 3 {
			return cols[2].Width - truncateBuffer
		}

		return unifiedTitleTruncate

	default:
		return 38 // Default fallback
	}
}

// hasHiddenColumns returns true if any columns are hidden at the current width
func hasHiddenColumns(
	viewType string,
	width int,
) bool {
	switch viewType {
	case itemTypeChange:
		return len(
			calculateChangesColumns(width, LineNumberOff),
		) < 4 // Full has 4 columns
	case itemTypeSpec:
		return len(
			calculateSpecsColumns(width),
		) < 3 // Full has 3 columns
	case itemTypeAll:
		return len(
			calculateUnifiedColumns(width),
		) < 4 // Full has 4 columns
	default:
		return false
	}
}

// buildChangesRows creates table rows for changes data with the given
// title truncation and column set. The column set determines which fields
// are included in each row to match the visible columns.
func buildChangesRows(
	changes []ChangeInfo,
	titleTruncate int,
	numColumns int,
	lineNumberMode LineNumberMode,
	cursor int,
) []table.Row {
	rows := make([]table.Row, len(changes))

	// Calculate effective number of data columns (excluding line number column)
	dataColumns := numColumns
	if lineNumberMode != LineNumberOff {
		dataColumns = numColumns - 1
	}

	for i, change := range changes {
		tasksStatus := fmt.Sprintf("%d/%d",
			change.TaskStatus.Completed,
			change.TaskStatus.Total)

		// Calculate line number
		var lineNumStr string
		if lineNumberMode != LineNumberOff {
			lineNum := calculateLineNumberValue(i, cursor, lineNumberMode)
			lineNumStr = fmt.Sprintf("%d", lineNum)
		}

		switch dataColumns {
		case 4:
			// Full: ID, Title, Deltas, Tasks
			row := table.Row{
				change.ID,
				tui.TruncateString(
					change.Title,
					titleTruncate,
				),
				fmt.Sprintf(
					"%d",
					change.DeltaCount,
				),
				tasksStatus,
			}
			if lineNumberMode != LineNumberOff {
				rows[i] = append(table.Row{lineNumStr}, row...)
			} else {
				rows[i] = row
			}
		case 3:
			// 3 columns without Title: ID, Deltas, Tasks
			row := table.Row{
				change.ID,
				fmt.Sprintf(
					"%d",
					change.DeltaCount,
				),
				tasksStatus,
			}
			if lineNumberMode != LineNumberOff {
				rows[i] = append(table.Row{lineNumStr}, row...)
			} else {
				rows[i] = row
			}
		default:
			// Minimal 2 columns: ID, Tasks only
			rows[i] = table.Row{
				change.ID,
				tasksStatus,
			}
		}
	}

	return rows
}

// buildSpecsRows creates table rows for specs data with the given
// title truncation and number of columns.
func buildSpecsRows(
	specs []SpecInfo,
	titleTruncate int,
	numColumns int,
) []table.Row {
	rows := make([]table.Row, len(specs))
	for i, spec := range specs {
		switch numColumns {
		case 3:
			// Full: ID, Title, Requirements
			rows[i] = table.Row{
				spec.ID,
				tui.TruncateString(
					spec.Title,
					titleTruncate,
				),
				fmt.Sprintf(
					"%d",
					spec.RequirementCount,
				),
			}
		default:
			// Minimal: ID, Title only
			rows[i] = table.Row{
				spec.ID,
				tui.TruncateString(
					spec.Title,
					titleTruncate,
				),
			}
		}
	}

	return rows
}

// buildUnifiedRows creates table rows for unified (all items) view with the
// given title truncation and number of columns.
func buildUnifiedRows(
	items ItemList,
	titleTruncate int,
	numColumns int,
) []table.Row {
	rows := make([]table.Row, len(items))
	for i, item := range items {
		var typeStr, details string
		switch item.Type {
		case ItemTypeChange:
			typeStr = typeDisplayChange
			if item.Change != nil {
				details = fmt.Sprintf(
					"Tasks: %d/%d 🔺 %d",
					item.Change.TaskStatus.Completed,
					item.Change.TaskStatus.Total,
					item.Change.DeltaCount,
				)
			}
		case ItemTypeSpec:
			typeStr = typeDisplaySpec
			if item.Spec != nil {
				details = fmt.Sprintf(
					"Reqs: %d",
					item.Spec.RequirementCount,
				)
			}
		}

		switch numColumns {
		case 4:
			// Full: ID, Type, Title, Details
			rows[i] = table.Row{
				item.ID(),
				typeStr,
				tui.TruncateString(
					item.Title(),
					titleTruncate,
				),
				details,
			}
		default:
			// Narrow/Minimal: ID, Type, Title (no Details)
			rows[i] = table.Row{
				item.ID(),
				typeStr,
				tui.TruncateString(
					item.Title(),
					titleTruncate,
				),
			}
		}
	}

	return rows
}

// interactiveModel represents the bubbletea model for interactive table
type interactiveModel struct {
	table            table.Model
	selectedID       string
	copied           bool
	quitting         bool
	archiveRequested bool
	prRequested      bool // true when P (pr) hotkey was pressed
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
	stdoutMode       bool        // true: prints ID to stdout instead of clipboard
	terminalWidth    int         // current terminal width for responsive columns
	// Source data for rebuilding rows on resize
	changesData      []ChangeInfo         // original changes data for changes/archive views
	specsData        []SpecInfo           // original specs data for specs view
	countPrefixState tui.CountPrefixState // vim-style count prefix state
	lineNumberMode   LineNumberMode       // line number display mode (off, relative, hybrid)
}

// Init initializes the model
func (interactiveModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *interactiveModel) Update(
	msg tea.Msg,
) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		keyStr := typedMsg.String()

		// Handle count prefix (before search mode)
		// Count prefix mode and search mode are mutually exclusive
		if !m.searchMode {
			count, isNavKey, handled := m.countPrefixState.HandleKey(typedMsg)
			if handled {
				if isNavKey {
					cursor := m.table.Cursor()
					rowCount := len(m.table.Rows())
					// Apply counted navigation with SetCursor
					switch keyStr {
					case "up", "k":
						newCursor := cursor - count
						if newCursor < 0 {
							newCursor = 0
						}
						m.table.SetCursor(newCursor)
					case "down", "j":
						newCursor := cursor + count
						if newCursor >= rowCount {
							newCursor = rowCount - 1
						}
						m.table.SetCursor(newCursor)
					}
					m.showHelp = false

					return m, nil
				}
				// Key was handled (digit or ESC) but not a nav key
				return m, nil
			}
		}

		// Handle search mode input
		if m.searchMode {
			var handled bool
			cmd, handled = m.handleSearchModeInput(typedMsg)
			if handled {
				return m, cmd
			}
		}

		switch keyStr {
		case "q", "ctrl+c":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			m.handleEnter()

			return m, tea.Quit

		case "e":
			return m.handleEdit()

		case "t":
			// Toggle filter type in unified mode
			if m.itemType == itemTypeAll {
				m.handleToggleFilter()

				return m, nil
			}

		case "a":
			return m.handleArchive()

		case "P":
			return m.handlePR()

		case "/":
			m.toggleSearchMode()

			return m, nil

		case "#":
			m.cycleLineNumberMode()

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

	case tea.WindowSizeMsg:
		// Store terminal width for responsive column calculations
		m.terminalWidth = typedMsg.Width
		// Trigger table rebuild to apply new column widths
		m.rebuildTableForWidth()

		return m, nil
	}

	prevCursor := m.table.Cursor()
	m.table, cmd = m.table.Update(msg)

	if m.lineNumberMode != LineNumberOff && m.table.Cursor() != prevCursor {
		m.updateLineNumbers()
	}

	return m, cmd
}

// handleEnter handles the enter key press for copying selected path
// or selecting in selection mode
func (m *interactiveModel) handleEnter() {
	cursor := m.table.Cursor()
	if cursor < 0 ||
		cursor >= len(m.table.Rows()) {
		return
	}

	row := m.table.Rows()[cursor]
	if len(row) == 0 {
		return
	}

	// ID is in first column for all modes
	itemID := row[0]
	m.selectedID = itemID

	// In selection mode, just select without copying to clipboard
	if m.selectionMode {
		return
	}

	// In stdout mode, just set selectedID and return (no clipboard copy)
	if m.stdoutMode {
		return
	}

	// Build the full path to copy to clipboard
	copyPath := m.buildCopyPath(itemID, row)

	// Copy to clipboard
	m.copied = true
	err := tui.CopyToClipboard(copyPath)
	if err != nil {
		m.err = err
	}
}

// buildCopyPath builds the path to copy for the selected item.
// Returns path relative to cwd (e.g., "spectr/changes/<id>/proposal.md")
func (m *interactiveModel) buildCopyPath(itemID string, row table.Row) string {
	// Determine the item type
	var itemType, rootPath string

	switch m.itemType {
	case itemTypeChange:
		// Find the change in changesData to get root path
		for _, change := range m.changesData {
			if change.ID == itemID {
				rootPath = change.RootPath

				break
			}
		}

		return buildChangePath(rootPath, itemID)

	case itemTypeSpec:
		// Find the spec in specsData to get root path
		for _, spec := range m.specsData {
			if spec.ID == itemID {
				rootPath = spec.RootPath

				break
			}
		}

		return buildSpecPath(rootPath, itemID)

	case itemTypeAll:
		// In unified mode, check the type column (second column)
		if len(row) > 1 {
			itemType = row[1]
		}

		// Find the item in allItems to get root path
		for i := range m.allItems {
			if m.allItems[i].ID() == itemID {
				rootPath = m.allItems[i].RootPath()

				break
			}
		}

		if itemType == typeDisplaySpec {
			return buildSpecPath(rootPath, itemID)
		}

		return buildChangePath(rootPath, itemID)
	}

	// Fallback to just the ID
	return itemID
}

// buildChangePath builds the path for a change.
func buildChangePath(rootPath, changeID string) string {
	// If rootPath is "." or empty, use current directory
	if rootPath == "" || rootPath == "." {
		return fmt.Sprintf("spectr/changes/%s", changeID)
	}
	// Otherwise prefix with root path
	return fmt.Sprintf("%s/spectr/changes/%s", rootPath, changeID)
}

// buildSpecPath builds the path for a spec.
func buildSpecPath(rootPath, specID string) string {
	// If rootPath is "." or empty, use current directory
	if rootPath == "" || rootPath == "." {
		return fmt.Sprintf("spectr/specs/%s", specID)
	}
	// Otherwise prefix with root path
	return fmt.Sprintf("%s/spectr/specs/%s", rootPath, specID)
}

// handleEdit handles the 'e' key press for opening file in editor
func (m *interactiveModel) handleEdit() (tea.Model, tea.Cmd) {
	// Get the selected row
	cursor := m.table.Cursor()
	if cursor < 0 ||
		cursor >= len(m.table.Rows()) {
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
		if itemTypeStr == typeDisplaySpec {
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
		m.err = &specterrs.EditorNotSetError{
			Operation: "edit",
		}

		return m, nil
	}

	// Construct file path based on type
	filePath := m.getEditFilePath(
		itemID,
		editItemType,
	)

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(
		err,
	) {
		m.err = fmt.Errorf(
			"file not found: %s",
			filePath,
		)

		return m, nil
	}

	// Launch editor - use tea.ExecProcess to handle editor lifecycle
	c := exec.Command(
		editor,
		filePath,
	) //nolint:gosec // G204: User controls EDITOR env var, intentional for opening editor

	return m, tea.ExecProcess(
		c,
		func(err error) tea.Msg {
			return editorFinishedMsg{err: err}
		},
	)
}

// handleToggleFilter toggles between showing all items,
// only changes, and only specs
func (m *interactiveModel) handleToggleFilter() {
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
	rebuildUnifiedTable(m)
}

// handleArchive handles the 'a' key press for archiving a change
func (m *interactiveModel) handleArchive() (tea.Model, tea.Cmd) {
	cursor := m.table.Cursor()
	if cursor < 0 ||
		cursor >= len(m.table.Rows()) {
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
		if len(row) > 1 && row[1] == typeDisplayChange {
			m.selectedID = row[0]
			m.archiveRequested = true

			return m, tea.Quit
		}
		// Not a change, do nothing
		return m, nil
	}

	return m, nil
}

// handlePR handles the 'P' key press for entering PR mode
func (m *interactiveModel) handlePR() (tea.Model, tea.Cmd) {
	// Only available in changes mode
	if m.itemType != itemTypeChange {
		return m, nil
	}

	// Get selected row
	cursor := m.table.Cursor()
	if cursor < 0 ||
		cursor >= len(m.table.Rows()) {
		return m, nil
	}

	row := m.table.Rows()[cursor]
	if len(row) == 0 {
		return m, nil
	}

	m.selectedID = row[0]
	m.prRequested = true

	return m, tea.Quit
}

// rebuildUnifiedTable rebuilds the table based on current filter
// and terminal width (for responsive columns)
func rebuildUnifiedTable(
	m *interactiveModel,
) {
	var items ItemList
	if m.filterType == nil {
		items = m.allItems
	} else {
		items = m.allItems.FilterByType(*m.filterType)
	}

	// Use responsive columns based on terminal width
	// If width is 0 (unknown), use full-width defaults
	width := m.terminalWidth
	if width == 0 {
		width = breakpointFull
	}
	columns := calculateUnifiedColumns(width)
	titleTruncate := calculateTitleTruncate(
		itemTypeAll,
		width,
	)

	rows := buildUnifiedRows(
		items,
		titleTruncate,
		len(columns),
	)

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
		"↑/↓/j/k: navigate (try 9j) | Enter: copy ID | e: edit | "+
			"a: archive | t: filter (%s) | #: line numbers | /: search | q: quit",
		filterDesc,
	)
	m.minimalFooter = fmt.Sprintf(
		"showing: %d | project: %s | ?: help",
		len(rows),
		m.projectPath,
	)
	// Preserve showHelp state (no reset needed here)
}

// applyFilter filters the table rows based on the search query.
// Memory optimization: reuses pre-allocated filteredRows buffer to reduce
// GC pressure during interactive search (which triggers on every keystroke).
func (m *interactiveModel) applyFilter() {
	if len(m.allRows) == 0 {
		return
	}

	query := strings.ToLower(m.searchQuery)

	if query == "" {
		// No filter - show all rows directly
		m.table.SetRows(m.allRows)
	} else {
		// Reuse the filteredRows buffer to avoid allocations on every
		// keystroke
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
func rowMatchesQuery(
	row table.Row,
	query string,
) bool {
	for _, col := range row {
		if strings.Contains(
			strings.ToLower(col),
			query,
		) {
			return true
		}
	}

	return false
}

// handleSearchModeInput handles keyboard input when in search mode
// Returns (cmd, handled) - handled is true if the input was processed
func (m *interactiveModel) handleSearchModeInput(
	keyMsg tea.KeyMsg,
) (tea.Cmd, bool) {
	var cmd tea.Cmd

	//nolint:exhaustive // Only handling specific keys
	switch keyMsg.Type {
	case tea.KeyEsc:
		// Exit search mode, clear query and restore all rows
		m.searchMode = false
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m.applyFilter()

		return nil, true
	case tea.KeyEnter:
		// Exit search mode but keep filter applied
		m.searchMode = false

		return nil, true
	default:
		// Update text input and filter
		m.searchInput, cmd = m.searchInput.Update(
			keyMsg,
		)
		m.searchQuery = m.searchInput.Value()
		m.applyFilter()

		return cmd, true
	}
}

// toggleSearchMode toggles search mode on/off
func (m *interactiveModel) toggleSearchMode() {
	if m.searchMode {
		m.searchMode = false
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m.applyFilter()
	} else {
		m.searchMode = true
		m.searchInput.Focus()
	}
}

func (m *interactiveModel) cycleLineNumberMode() {
	switch m.lineNumberMode {
	case LineNumberOff:
		m.lineNumberMode = LineNumberRelative
	case LineNumberRelative:
		m.lineNumberMode = LineNumberHybrid
	case LineNumberHybrid:
		m.lineNumberMode = LineNumberOff
	}

	// Trigger table rebuild to add/remove line number column
	m.rebuildTableForWidth()
}

func (m *interactiveModel) renderLineNumbers() string {
	if m.lineNumberMode == LineNumberOff {
		return ""
	}

	rows := m.table.Rows()
	if len(rows) == 0 {
		return ""
	}

	cursor := m.table.Cursor()
	var result strings.Builder

	for i := range rows {
		num := m.calculateLineNumber(i, cursor)
		if i == cursor {
			result.WriteString(tui.CurrentLineNumberStyle().Render(fmt.Sprintf("%d", num)))
		} else {
			result.WriteString(tui.LineNumberStyle().Render(fmt.Sprintf("%d", num)))
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (m *interactiveModel) calculateLineNumber(rowIdx, cursorIdx int) int {
	switch m.lineNumberMode {
	case LineNumberRelative:
		return abs(rowIdx - cursorIdx)
	case LineNumberHybrid:
		if rowIdx == cursorIdx {
			return cursorIdx + 1
		}
		return abs(rowIdx - cursorIdx)
	default:
		return 0
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// calculateLineNumberValue returns the display value for a line number
func calculateLineNumberValue(rowIdx, cursorIdx int, mode LineNumberMode) int {
	switch mode {
	case LineNumberRelative:
		return abs(rowIdx - cursorIdx)
	case LineNumberHybrid:
		if rowIdx == cursorIdx {
			return cursorIdx + 1
		}
		return abs(rowIdx - cursorIdx)
	default:
		return 0
	}
}

// rebuildTableForWidth rebuilds the table with columns adjusted for the
// current terminal width. This is called when the terminal is resized.
// It preserves cursor position and search state during the rebuild.
func (m *interactiveModel) rebuildTableForWidth() {
	// Store current cursor position to restore after rebuild
	currentCursor := m.table.Cursor()

	// Use responsive columns based on terminal width
	// If width is 0 (unknown), use full-width defaults
	width := m.terminalWidth
	if width == 0 {
		width = breakpointFull
	}

	// Rebuild table based on item type
	switch m.itemType {
	case itemTypeAll:
		rebuildUnifiedTable(m)
	case itemTypeChange:
		m.rebuildChangesTable(width)
	case itemTypeSpec:
		m.rebuildSpecsTable(width)
	}

	// Restore cursor position (bounded to row count)
	rowCount := len(m.table.Rows())
	if rowCount > 0 {
		if currentCursor >= rowCount {
			m.table.SetCursor(rowCount - 1)
		} else {
			m.table.SetCursor(currentCursor)
		}
	}

	// Re-apply search filter if active
	if m.searchQuery != "" {
		m.applyFilter()
	}
}

// rebuildChangesTable rebuilds the changes table with responsive columns
func (m *interactiveModel) rebuildChangesTable(
	width int,
) {
	if len(m.changesData) == 0 {
		return
	}

	columns := calculateChangesColumns(width, m.lineNumberMode)
	titleTruncate := calculateTitleTruncate(
		itemTypeChange,
		width,
	)
	rows := buildChangesRows(
		m.changesData,
		titleTruncate,
		len(columns),
		m.lineNumberMode,
		m.table.Cursor(),
	)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)
	tui.ApplyTableStyles(&t)

	m.table = t
	m.allRows = rows
}

// rebuildSpecsTable rebuilds the specs table with responsive columns
func (m *interactiveModel) rebuildSpecsTable(
	width int,
) {
	if len(m.specsData) == 0 {
		return
	}

	columns := calculateSpecsColumns(width)
	titleTruncate := calculateTitleTruncate(
		itemTypeSpec,
		width,
	)
	rows := buildSpecsRows(
		m.specsData,
		titleTruncate,
		len(columns),
	)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
	)
	tui.ApplyTableStyles(&t)

	m.table = t
	m.allRows = rows
}

// getEditFilePath returns the file path to edit based on item type
func (m *interactiveModel) getEditFilePath(
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
func (m *interactiveModel) View() string {
	if m.quitting {
		if m.archiveRequested &&
			m.selectedID != "" {
			return fmt.Sprintf(
				"Archiving: %s\n",
				m.selectedID,
			)
		}

		if m.prRequested && m.selectedID != "" {
			return fmt.Sprintf(
				"PR mode: %s\n",
				m.selectedID,
			)
		}

		// In stdout mode, output just the ID for piping
		if m.stdoutMode && m.selectedID != "" {
			return m.selectedID + "\n"
		}

		// In selection mode, just show selected ID without clipboard
		// message.
		if m.selectionMode && m.selectedID != "" {
			return fmt.Sprintf(
				"Selected: %s\n",
				m.selectedID,
			)
		}

		if m.copied && m.err == nil {
			return fmt.Sprintf(
				"✓ Copied: %s\n",
				m.selectedID,
			)
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
		view = fmt.Sprintf(
			"Search: %s\n\n",
			m.searchInput.View(),
		)
	}

	// Choose which footer to display based on showHelp state
	footer := m.minimalFooter
	if m.showHelp {
		footer = m.helpText
	}

	// Append hidden columns hint if columns are hidden due to narrow terminal
	if m.terminalWidth > 0 &&
		hasHiddenColumns(
			m.itemType,
			m.terminalWidth,
		) {
		footer += " | (some columns hidden)"
	}

	if m.countPrefixState.IsActive() {
		footer += fmt.Sprintf(" | count: %s_", m.countPrefixState.String())
	}

	if m.lineNumberMode != LineNumberOff {
		modeStr := "rel"
		if m.lineNumberMode == LineNumberHybrid {
			modeStr = "hyb"
		}
		footer += fmt.Sprintf(" | ln: %s", modeStr)
	}

	view += m.table.View() + "\n" + footer + "\n"

	// Display error message if present, but keep TUI active
	if m.err != nil {
		view += fmt.Sprintf(
			"\nError: %v\n",
			m.err,
		)
	}

	return view
}

// RunInteractiveChanges runs the interactive table for changes.
// Returns (archiveID, prID, error):
//   - archiveID is set if archive was requested via 'a' key
//   - prID is set if PR mode was requested via 'P' key
//   - Both are empty if user quit or cancelled
func RunInteractiveChanges(
	changes []ChangeInfo,
	projectPath string,
	stdoutMode bool,
) (archiveID, prID string, err error) {
	if len(changes) == 0 {
		return "", "", nil
	}

	// Use default full-width columns initially (terminalWidth=0 means unknown)
	//
	// WindowSizeMsg will trigger a rebuild with correct responsive columns
	columns := calculateChangesColumns(
		breakpointFull,
		LineNumberOff,
	)
	titleTruncate := calculateTitleTruncate(
		itemTypeChange,
		breakpointFull,
	)

	rows := buildChangesRows(
		changes,
		titleTruncate,
		len(columns),
		LineNumberOff,
		0,
	)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	tui.ApplyTableStyles(&t)

	m := &interactiveModel{
		table:         t,
		itemType:      itemTypeChange,
		projectPath:   projectPath,
		searchInput:   newTextInput(),
		allRows:       rows,
		terminalWidth: 0,          // Will be set by WindowSizeMsg
		changesData:   changes,    // Store for rebuild on resize
		stdoutMode:    stdoutMode, // Output to stdout instead of clipboard
		helpText: "↑/↓/j/k: navigate (try 9j) | Enter: copy ID | e: edit | " +
			"a: archive | P: pr | #: line numbers | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, runErr := p.Run()
	if runErr != nil {
		return "", "", fmt.Errorf(
			errInteractiveModeFormat,
			runErr,
		)
	}

	// Check if there was an error during execution
	if fm, ok := finalModel.(*interactiveModel); ok {
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
		if fm.archiveRequested &&
			fm.selectedID != "" {
			return fm.selectedID, "", nil
		}

		// Return PR ID if PR was requested
		if fm.prRequested && fm.selectedID != "" {
			return "", fm.selectedID, nil
		}
	}

	return "", "", nil
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

	// Use default full-width columns initially (terminalWidth=0 means unknown)
	// WindowSizeMsg will trigger a rebuild with correct responsive columns
	columns := calculateChangesColumns(
		breakpointFull,
		LineNumberOff,
	)
	titleTruncate := calculateTitleTruncate(
		itemTypeChange,
		breakpointFull,
	)

	rows := buildChangesRows(
		changes,
		titleTruncate,
		len(columns),
		LineNumberOff,
		0,
	)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	tui.ApplyTableStyles(&t)

	m := &interactiveModel{
		table:         t,
		itemType:      itemTypeChange,
		projectPath:   projectPath,
		searchInput:   newTextInput(),
		allRows:       rows,
		terminalWidth: 0,       // Will be set by WindowSizeMsg
		changesData:   changes, // Store for rebuild on resize
		selectionMode: true,    // Enter selects without copying
		helpText:      "↑/↓/j/k: navigate (try 9j) | Enter: select | #: line numbers | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf(
			errInteractiveModeFormat,
			err,
		)
	}

	// Check if there was an error during execution
	if fm, ok := finalModel.(*interactiveModel); ok {
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
func RunInteractiveSpecs(
	specs []SpecInfo,
	projectPath string,
	stdoutMode bool,
) error {
	if len(specs) == 0 {
		return nil
	}

	// Use default full-width columns initially (terminalWidth=0 means unknown)
	//
	// WindowSizeMsg will trigger a rebuild with correct responsive columns
	columns := calculateSpecsColumns(
		breakpointFull,
	)
	titleTruncate := calculateTitleTruncate(
		itemTypeSpec,
		breakpointFull,
	)

	rows := buildSpecsRows(
		specs,
		titleTruncate,
		len(columns),
	)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
	)

	tui.ApplyTableStyles(&t)

	m := &interactiveModel{
		table:         t,
		itemType:      itemTypeSpec,
		projectPath:   projectPath,
		searchInput:   newTextInput(),
		allRows:       rows,
		terminalWidth: 0,          // Will be set by WindowSizeMsg
		specsData:     specs,      // Store for rebuild on resize
		stdoutMode:    stdoutMode, // Output to stdout instead of clipboard
		helpText: "↑/↓/j/k: navigate (try 9j) | Enter: copy ID | e: edit | " +
			"#: line numbers | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(specs),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf(
			errInteractiveModeFormat,
			err,
		)
	}

	// Check if there was an error during execution
	fm, ok := finalModel.(*interactiveModel)
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
func RunInteractiveAll(
	items ItemList,
	projectPath string,
	stdoutMode bool,
) error {
	if len(items) == 0 {
		return nil
	}

	// Use default full-width columns initially (terminalWidth=0 means unknown)
	// WindowSizeMsg will trigger a rebuild with correct responsive columns
	columns := calculateUnifiedColumns(
		breakpointFull,
	)
	titleTruncate := calculateTitleTruncate(
		itemTypeAll,
		breakpointFull,
	)

	rows := buildUnifiedRows(
		items,
		titleTruncate,
		len(columns),
	)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	tui.ApplyTableStyles(&t)

	m := &interactiveModel{
		table:         t,
		itemType:      itemTypeAll,
		projectPath:   projectPath,
		allItems:      items,
		filterType:    nil, // Start with all items visible
		searchInput:   newTextInput(),
		allRows:       rows,
		terminalWidth: 0,          // Will be set by WindowSizeMsg
		stdoutMode:    stdoutMode, // Output to stdout instead of clipboard
		helpText: "↑/↓/j/k: navigate (try 9j) | Enter: copy ID | e: edit | " +
			"a: archive | t: filter (all) | #: line numbers | /: search | q: quit",
		minimalFooter: fmt.Sprintf(
			"showing: %d | project: %s | ?: help",
			len(rows),
			projectPath,
		),
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf(
			errInteractiveModeFormat,
			err,
		)
	}

	// Check if there was an error during execution
	if fm, ok := finalModel.(*interactiveModel); ok &&
		fm.err != nil {
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
