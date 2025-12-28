// Package cmd provides command-line interface implementations.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/specterrs"
	"github.com/connerohnesorge/spectr/internal/validation"
)

// ValidateCmd represents the validate command
type ValidateCmd struct {
	ItemName      *string `arg:"" optional:"" predictor:"item"`
	JSON          bool    `                                        name:"json"           help:"Output as JSON"`                      //nolint:lll,revive // Kong struct tag with alignment
	All           bool    `                                        name:"all"            help:"Validate all"`                        //nolint:lll,revive // Kong struct tag with alignment
	Changes       bool    `                                        name:"changes"        help:"Validate changes"`                    //nolint:lll,revive // Kong struct tag with alignment
	Specs         bool    `                                        name:"specs"          help:"Validate specs"`                      //nolint:lll,revive // Kong struct tag with alignment
	Type          *string `                   predictor:"itemType" name:"type"                                   enum:"change,spec"` //nolint:lll,revive // Kong struct tag with alignment
	NoInteractive bool    `                                        name:"no-interactive" help:"No prompts"`                          //nolint:lll,revive // Kong struct tag with alignment
}

// Run executes the validate command
func (c *ValidateCmd) Run() error {
	// Get current working directory
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(
			"failed to get current directory: %w",
			err,
		)
	}

	cfg, err := config.Load(projectPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Check if bulk validation flags are set
	if c.All || c.Changes || c.Specs {
		return c.runBulkValidation(projectPath, cfg.Dir)
	}

	// If no item name provided
	if c.ItemName == nil || *c.ItemName == "" {
		if c.NoInteractive {
			return getUsageError()
		}
		// Launch interactive mode
		return validation.RunInteractiveValidation(
			projectPath,
			c.JSON,
			cfg.Dir,
		)
	}

	// Direct validation
	return c.runDirectValidation(
		projectPath,
		*c.ItemName,
		cfg.Dir,
	)
}

// runDirectValidation validates a single item (change or spec)
func (c *ValidateCmd) runDirectValidation(
	projectPath, itemName, spectrDir string,
) error {
	// Normalize the item path to extract ID and infer type
	normalizedID, inferredType := discovery.NormalizeItemPath(
		itemName,
	)

	// Use inferred type if available, otherwise fall back to explicit type flag
	typeHint := c.Type
	if inferredType != "" {
		typeHint = &inferredType
	}

	// Determine item type
	info, err := validation.DetermineItemType(
		projectPath, normalizedID, typeHint, spectrDir,
	)
	if err != nil {
		return err
	}

	// Create validator and validate
	validator := validation.NewValidator()
	report, err := validation.ValidateItemByType(
		validator,
		projectPath,
		normalizedID,
		info.ItemType,
		spectrDir,
	)
	if err != nil {
		return fmt.Errorf(
			"validation failed: %w",
			err,
		)
	}

	// Print report
	if c.JSON {
		validation.PrintJSONReport(report)
	} else {
		validation.PrintHumanReport(normalizedID, report)
	}

	// Return error if validation failed
	if !report.Valid {
		return &specterrs.ValidationFailedError{
			ErrorCount:   report.Summary.Errors,
			WarningCount: report.Summary.Warnings,
		}
	}

	return nil
}

// runBulkValidation validates multiple items based on flags
func (c *ValidateCmd) runBulkValidation(
	projectPath, spectrDir string,
) error {
	validator := validation.NewValidator()

	// Determine what to validate
	items, err := c.getItemsToValidate(
		projectPath,
		spectrDir,
	)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return c.handleNoItems()
	}

	// Validate all items
	results, hasFailures := c.validateAllItems(
		validator,
		items,
	)

	// Print results
	if c.JSON {
		validation.PrintBulkJSONResults(results)
	} else {
		validation.PrintBulkHumanResults(results)
	}

	if hasFailures {
		return &specterrs.MultiValidationFailedError{
			ItemCount: len(items),
		}
	}

	return nil
}

// getItemsToValidate returns the items to validate based on flags
func (c *ValidateCmd) getItemsToValidate(
	projectPath, spectrDir string,
) ([]validation.ValidationItem, error) {
	switch {
	case c.All:
		return validation.GetAllItems(projectPath, spectrDir)
	case c.Changes:
		return validation.GetChangeItems(
			projectPath,
			spectrDir,
		)
	case c.Specs:
		return validation.GetSpecItems(
			projectPath,
			spectrDir,
		)
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
	results := make(
		[]validation.BulkResult,
		0,
		len(items),
	)
	hasFailures := false

	for _, item := range items {
		result, err := validation.ValidateSingleItem(
			validator,
			item,
		)
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
