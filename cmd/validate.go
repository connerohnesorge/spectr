// Package cmd provides command-line interface implementations.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/validation"
)

// ValidateCmd represents the validate command
type ValidateCmd struct {
	ItemName      *string `arg:"" optional:"" help:"Item to validate"`
	Strict        bool    `name:"strict" help:"Treat warnings as errors"`
	JSON          bool    `name:"json" help:"Output as JSON"`
	All           bool    `name:"all" help:"Validate all"`
	Changes       bool    `name:"changes" help:"Validate changes"`
	Specs         bool    `name:"specs" help:"Validate specs"`
	Type          *string `name:"type" enum:"change,spec" help:"Item type"`
	NoInteractive bool    `name:"no-interactive" help:"No prompts"`
}

// Run executes the validate command
func (c *ValidateCmd) Run() error {
	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if bulk validation flags are set
	if c.All || c.Changes || c.Specs {
		return c.runBulkValidation(projectPath)
	}

	// If no item name provided
	if c.ItemName == nil || *c.ItemName == "" {
		if c.NoInteractive {
			return getUsageError()
		}
		// Launch interactive mode
		return validation.RunInteractiveValidation(
			projectPath, c.Strict, c.JSON,
		)
	}

	// Direct validation
	return c.runDirectValidation(projectPath, *c.ItemName)
}

// runDirectValidation validates a single item (change or spec)
func (c *ValidateCmd) runDirectValidation(
	projectPath, itemName string,
) error {
	// Determine item type
	info, err := validation.DetermineItemType(projectPath, itemName, c.Type)
	if err != nil {
		return err
	}

	// Create validator and validate
	validator := validation.NewValidator(c.Strict)
	report, err := validation.ValidateItemByType(
		validator,
		projectPath,
		itemName,
		info.ItemType,
	)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Print report
	if c.JSON {
		validation.PrintJSONReport(report)
	} else {
		validation.PrintHumanReport(itemName, report)
	}

	// Return error if validation failed
	if !report.Valid {
		return errors.New("validation failed")
	}

	return nil
}

// runBulkValidation validates multiple items based on flags
func (c *ValidateCmd) runBulkValidation(projectPath string) error {
	validator := validation.NewValidator(c.Strict)

	// Determine what to validate
	items, err := c.getItemsToValidate(projectPath)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return c.handleNoItems()
	}

	// Validate all items
	results, hasFailures := c.validateAllItems(validator, items)

	// Print results
	if c.JSON {
		validation.PrintBulkJSONResults(results)
	} else {
		validation.PrintBulkHumanResults(results)
	}

	if hasFailures {
		return errors.New("validation failed for one or more items")
	}

	return nil
}

// getItemsToValidate returns the items to validate based on flags
func (c *ValidateCmd) getItemsToValidate(
	projectPath string,
) ([]validation.ValidationItem, error) {
	switch {
	case c.All:
		return validation.GetAllItems(projectPath)
	case c.Changes:
		return validation.GetChangeItems(projectPath)
	case c.Specs:
		return validation.GetSpecItems(projectPath)
	default:
		return nil, nil
	}
}

// handleNoItems handles the case when there are no items to validate
func (c *ValidateCmd) handleNoItems() error {
	if c.JSON {
		fmt.Println("[]")
	} else {
		fmt.Println("No items to validate")
	}

	return nil
}

// validateAllItems validates all items and returns results
func (*ValidateCmd) validateAllItems(
	validator *validation.Validator,
	items []validation.ValidationItem,
) ([]validation.BulkResult, bool) {
	results := make([]validation.BulkResult, 0, len(items))
	hasFailures := false

	for _, item := range items {
		result, err := validation.ValidateSingleItem(validator, item)
		results = append(results, result)

		if err != nil || !result.Valid {
			hasFailures = true
		}
	}

	return results, hasFailures
}

// getUsageError returns the usage error message
func getUsageError() error {
	return errors.New(
		"usage: spectr validate <item-name> [flags]\n" +
			"       spectr validate --all\n" +
			"       spectr validate --changes\n" +
			"       spectr validate --specs",
	)
}
