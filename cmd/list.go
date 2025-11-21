// Package cmd provides command-line interface implementations for Spectr.
// This file contains the list command for displaying changes and specs.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/list"
)

// ListCmd represents the list command which displays changes or specs.
// It supports multiple output formats: text, long (detailed), JSON, and
// interactive table mode with clipboard support.
type ListCmd struct {
	// Specs determines whether to list specifications instead of changes
	Specs bool `name:"specs" help:"List specifications instead of changes"`
	// All determines whether to list both changes and specs in unified mode
	All bool `name:"all" help:"List both changes and specs in unified mode"`
	// Long enables detailed output with titles and counts
	Long bool `name:"long" help:"Show detailed output with titles and counts"`
	// JSON enables JSON output format
	JSON bool `name:"json" help:"Output as JSON"`
	// Interactive enables interactive table mode with clipboard
	Interactive bool `short:"I" name:"interactive" help:"Interactive mode"`
}

// Run executes the list command.
// It validates flags, determines the project path, and delegates to
// either listSpecs, listChanges, or listAll based on the flags.
func (c *ListCmd) Run() error {
	// Validate flags - interactive and JSON are mutually exclusive
	if c.Interactive && c.JSON {
		return errors.New("cannot use --interactive with --json")
	}

	// Validate flags - all and specs are mutually exclusive
	if c.All && c.Specs {
		return errors.New("cannot use --all with --specs")
	}

	// Get current working directory as the project path
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(
			"failed to get current directory: %w",
			err,
		)
	}

	// Create lister instance for the project
	lister := list.NewLister(projectPath)

	// Route to appropriate listing function
	if c.All {
		return c.listAll(lister, projectPath)
	}
	if c.Specs {
		return c.listSpecs(lister, projectPath)
	}

	return c.listChanges(lister, projectPath)
}

// listChanges retrieves and displays changes in the requested format.
// It handles interactive mode, JSON, long, and default text formats.
func (c *ListCmd) listChanges(lister *list.Lister, projectPath string) error {
	// Retrieve all changes from the project
	changes, err := lister.ListChanges()
	if err != nil {
		return fmt.Errorf("failed to list changes: %w", err)
	}

	// Handle interactive mode - shows a navigable table
	if c.Interactive {
		if len(changes) == 0 {
			fmt.Println("No changes found.")

			return nil
		}

		return list.RunInteractiveChanges(changes, projectPath)
	}

	// Format output based on flags
	var output string
	switch {
	case c.JSON:
		// JSON format for machine consumption
		var err error
		output, err = list.FormatChangesJSON(changes)
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
	case c.Long:
		// Long format with detailed information
		output = list.FormatChangesLong(changes)
	default:
		// Default text format - simple ID list
		output = list.FormatChangesText(changes)
	}

	// Display the formatted output
	fmt.Println(output)

	return nil
}

// listSpecs retrieves and displays specifications in the requested format.
// It handles interactive mode, JSON, long, and default text formats.
func (c *ListCmd) listSpecs(lister *list.Lister, projectPath string) error {
	// Retrieve all specifications from the project
	specs, err := lister.ListSpecs()
	if err != nil {
		return fmt.Errorf("failed to list specs: %w", err)
	}

	// Handle interactive mode - shows a navigable table
	if c.Interactive {
		if len(specs) == 0 {
			fmt.Println("No specs found.")

			return nil
		}

		return list.RunInteractiveSpecs(specs, projectPath)
	}

	// Format output based on flags
	var output string
	switch {
	case c.JSON:
		// JSON format for machine consumption
		var err error
		output, err = list.FormatSpecsJSON(specs)
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
	case c.Long:
		// Long format with detailed information
		output = list.FormatSpecsLong(specs)
	default:
		// Default text format - simple ID list
		output = list.FormatSpecsText(specs)
	}

	// Display the formatted output
	fmt.Println(output)

	return nil
}

// listAll retrieves and displays both changes and specs in unified format.
// It handles interactive mode, JSON, long, and default text formats.
func (c *ListCmd) listAll(lister *list.Lister, projectPath string) error {
	// Retrieve all items (changes and specs) from the project
	items, err := lister.ListAll(nil)
	if err != nil {
		return fmt.Errorf("failed to list all items: %w", err)
	}

	// Handle interactive mode - shows a unified navigable table
	if c.Interactive {
		if len(items) == 0 {
			fmt.Println("No items found.")

			return nil
		}

		return list.RunInteractiveAll(items, projectPath)
	}

	// Format output based on flags
	var output string
	switch {
	case c.JSON:
		// JSON format for machine consumption
		var err error
		output, err = list.FormatAllJSON(items)
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
	case c.Long:
		// Long format with detailed information
		output = list.FormatAllLong(items)
	default:
		// Default text format - simple ID list with type indicators
		output = list.FormatAllText(items)
	}

	// Display the formatted output
	fmt.Println(output)

	return nil
}
