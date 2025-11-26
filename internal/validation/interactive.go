//nolint:revive // file-length-limit - interactive functions logically grouped
package validation

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/connerohnesorge/spectr/internal/tui"
	"github.com/mattn/go-isatty"
)

const (
	// Table column widths for validation items
	validationIDWidth   = 35
	validationTypeWidth = 10
	validationPathWidth = 55

	// Truncation settings
	validationPathTruncate = 53

	// Table height
	tableHeight = 10

	// Menu selection indices
	menuSelectionAll      = 0
	menuSelectionChanges  = 1
	menuSelectionSpecs    = 2
	menuSelectionPickItem = 3
)

// ellipsisMinLength is exported from tui package for backward compatibility.
const ellipsisMinLength = tui.EllipsisMinLength

// RunInteractiveValidation runs the interactive validation TUI
func RunInteractiveValidation(projectPath string, strict bool, jsonOutput bool) error {
	// Check if running in a TTY
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return fmt.Errorf("interactive mode requires a TTY")
	}

	// Show menu and get selection
	selection, err := showValidationMenu()
	if err != nil {
		return err
	}

	if selection < 0 {
		// User cancelled
		return nil
	}

	// Handle selection
	return handleMenuSelection(selection, projectPath, strict, jsonOutput)
}

// showValidationMenu displays the validation menu and returns the selection
func showValidationMenu() (int, error) {
	menu := tui.NewMenuPicker(tui.MenuConfig{
		Title: "Validation Menu",
		Choices: []string{
			"All (changes and specs)",
			"All changes",
			"All specs",
			"Pick specific item",
		},
	})

	return menu.Run()
}

// handleMenuSelection processes the menu selection
func handleMenuSelection(selection int, projectPath string, strict, jsonOutput bool) error {
	var items []ValidationItem
	var err error

	switch selection {
	case menuSelectionAll:
		items, err = GetAllItems(projectPath)
	case menuSelectionChanges:
		items, err = GetChangeItems(projectPath)
	case menuSelectionSpecs:
		items, err = GetSpecItems(projectPath)
	case menuSelectionPickItem:
		return runItemPicker(projectPath, strict, jsonOutput)
	default:
		return nil
	}

	if err != nil {
		return fmt.Errorf("error loading items: %w", err)
	}

	return runValidationAndPrint(items, strict, jsonOutput)
}

// runItemPicker shows the item picker and validates the selected item
func runItemPicker(projectPath string, strict, jsonOutput bool) error {
	items, err := GetAllItems(projectPath)
	if err != nil {
		return fmt.Errorf("error loading items: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No items to validate")

		return nil
	}

	// Build table rows
	columns := []table.Column{
		{Title: "Name", Width: validationIDWidth},
		{Title: "Type", Width: validationTypeWidth},
		{Title: "Path", Width: validationPathWidth},
	}

	rows := make([]table.Row, len(items))
	for i, item := range items {
		rows[i] = table.Row{
			item.Name,
			item.ItemType,
			tui.TruncateString(item.Path, validationPathTruncate),
		}
	}

	// Create picker with enter action for selection
	var selectedIdx = -1
	picker := tui.NewTablePicker(tui.TableConfig{
		Columns:     columns,
		Rows:        rows,
		Height:      tableHeight,
		ProjectPath: projectPath,
		Actions: map[string]tui.Action{
			"enter": {
				Key:         "enter",
				Description: "validate",
				Handler: func(row table.Row) (tea.Cmd, *tui.ActionResult) {
					if len(row) == 0 {
						return nil, nil
					}
					// Find the index by matching the name
					for i, item := range items {
						if item.Name == row[0] {
							selectedIdx = i

							break
						}
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
		return fmt.Errorf("error running item picker: %w", err)
	}

	if result == nil || result.Cancelled || selectedIdx < 0 {

		return nil
	}

	// Validate the selected item
	selectedItem := items[selectedIdx]

	return runValidationAndPrint([]ValidationItem{selectedItem}, strict, jsonOutput)
}

// runValidationAndPrint validates items and prints results
func runValidationAndPrint(items []ValidationItem, strict, jsonOutput bool) error {
	if len(items) == 0 {
		fmt.Println("No items to validate")

		return nil
	}

	validator := NewValidator(strict)
	results, _ := validateItems(validator, items)

	if jsonOutput {
		PrintBulkJSONResults(results)
	} else {
		PrintBulkHumanResults(results)
	}

	return nil
}

// validateItems validates a list of items and returns results
func validateItems(
	validator *Validator,
	items []ValidationItem,
) ([]BulkResult, bool) {
	results := make([]BulkResult, 0, len(items))
	hasFailures := false

	for _, item := range items {
		result, err := ValidateSingleItem(validator, item)
		results = append(results, result)

		if err != nil || !result.Valid {
			hasFailures = true
		}
	}

	return results, hasFailures
}

// truncateString truncates a string and adds ellipsis if needed.
// This is a thin wrapper around tui.TruncateString for backward compatibility.
func truncateString(s string, maxLen int) string {
	return tui.TruncateString(s, maxLen)
}

// applyTableStyles applies default styling to a table.
// This is a thin wrapper around tui.ApplyTableStyles for backward compatibility.
func applyTableStyles(t *table.Model) {
	tui.ApplyTableStyles(t)
}
