package cmd

import (
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/list"
	"github.com/connerohnesorge/spectr/internal/pr"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// ListCmd represents the list command which displays changes or specs.
// It supports multiple output formats: text, long (detailed), JSON, and
// interactive table mode with clipboard support.
type ListCmd struct {
	// Specs determines whether to list specifications instead of changes
	Specs bool `name:"specs" help:"List specifications instead of changes"` //nolint:lll,revive
	// All determines whether to list both changes and specs in unified mode
	All bool `name:"all"   help:"List both changes and specs in unified mode"` //nolint:lll,revive

	// Long enables detailed output with titles and counts
	Long bool `name:"long" help:"Show detailed output with titles and counts"` //nolint:lll,revive

	// JSON enables JSON output format
	JSON bool `name:"json" help:"Output as JSON"` //nolint:lll,revive

	// Interactive enables interactive table mode with clipboard
	Interactive bool `name:"interactive" help:"Interactive mode" short:"I"` //nolint:lll,revive

	// Stdout prints selected ID to stdout instead of clipboard.
	// Requires -I (interactive mode).
	Stdout bool `name:"stdout" help:"Print ID to stdout (requires -I)"` //nolint:lll,revive
}

// Run executes the list command.
// It validates flags, determines the project path, and delegates to
// either listSpecs, listChanges, or listAll based on the flags.
func (c *ListCmd) Run() error {
	// Validate flags - interactive and JSON are mutually exclusive
	if c.Interactive && c.JSON {
		return &specterrs.IncompatibleFlagsError{
			Flag1: "--interactive",
			Flag2: "--json",
		}
	}

	// Validate flags - all and specs are mutually exclusive
	if c.All && c.Specs {
		return &specterrs.IncompatibleFlagsError{
			Flag1: "--all",
			Flag2: "--specs",
		}
	}

	// Validate flags - stdout requires interactive mode
	if c.Stdout && !c.Interactive {
		return &specterrs.RequiresFlagError{
			Flag:         "--stdout",
			RequiredFlag: "--interactive",
		}
	}

	// Validate flags - stdout and JSON are mutually exclusive
	if c.Stdout && c.JSON {
		return &specterrs.IncompatibleFlagsError{
			Flag1: "--stdout",
			Flag2: "--json",
		}
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
func (c *ListCmd) listChanges(
	lister *list.Lister,
	projectPath string,
) error {
	// Retrieve all changes from the project
	changes, err := lister.ListChanges()
	if err != nil {
		return fmt.Errorf(
			"failed to list changes: %w",
			err,
		)
	}

	// Handle interactive mode - shows a navigable table
	if c.Interactive {
		if len(changes) == 0 {
			fmt.Println("No changes found.")

			return nil
		}

		archiveID, prID, err := list.RunInteractiveChanges(
			changes,
			projectPath,
			c.Stdout,
		)
		if err != nil {
			return err
		}

		// If an archive was requested, run the archive workflow
		if archiveID != "" {
			return c.runArchiveWorkflow(
				archiveID,
				projectPath,
			)
		}

		// If PR mode was requested, run the PR workflow
		if prID != "" {
			return c.runPRWorkflow(
				prID,
				projectPath,
			)
		}

		return nil
	}

	// Format output based on flags
	var output string
	switch {
	case c.JSON:
		// JSON format for machine consumption
		var err error
		output, err = list.FormatChangesJSON(
			changes,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to format JSON: %w",
				err,
			)
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

// runArchiveWorkflow executes the archive workflow for a change.
func (*ListCmd) runArchiveWorkflow(
	changeID, projectPath string,
) error {
	// Create archive command with the selected change ID
	archiveCmd := &archive.ArchiveCmd{
		ChangeID: changeID,
		// Skip confirmation since user already selected in interactive mode
		Yes: true,
	}

	// Run the archive workflow
	// Result is discarded for interactive usage - already prints to terminal
	if _, err := archive.Archive(archiveCmd, projectPath); err != nil {
		return fmt.Errorf(
			"archive workflow failed: %w",
			err,
		)
	}

	return nil
}

// runPRWorkflow executes the PR proposal workflow for a change.
func (*ListCmd) runPRWorkflow(
	changeID, projectPath string,
) error {
	config := pr.PRConfig{
		ChangeID:    changeID,
		Mode:        pr.ModeProposal,
		ProjectRoot: projectPath,
	}

	result, err := pr.ExecutePR(config)
	if err != nil {
		return fmt.Errorf(
			"pr workflow failed: %w",
			err,
		)
	}

	// Print the PR result
	fmt.Println()
	fmt.Printf("Branch: %s\n", result.BranchName)

	if result.PRURL != "" {
		fmt.Printf(
			"\nPR created: %s\n",
			result.PRURL,
		)
	} else if result.ManualURL != "" {
		fmt.Printf("\nCreate PR manually: %s\n", result.ManualURL)
	}

	return nil
}

// listSpecs retrieves and displays specifications in the requested format.
// It handles interactive mode, JSON, long, and default text formats.
func (c *ListCmd) listSpecs(
	lister *list.Lister,
	projectPath string,
) error {
	// Retrieve all specifications from the project
	specs, err := lister.ListSpecs()
	if err != nil {
		return fmt.Errorf(
			"failed to list specs: %w",
			err,
		)
	}

	// Handle interactive mode - shows a navigable table
	if c.Interactive {
		if len(specs) == 0 {
			fmt.Println("No specs found.")

			return nil
		}

		return list.RunInteractiveSpecs(
			specs,
			projectPath,
			c.Stdout,
		)
	}

	// Format output based on flags
	var output string
	switch {
	case c.JSON:
		// JSON format for machine consumption
		var err error
		output, err = list.FormatSpecsJSON(specs)
		if err != nil {
			return fmt.Errorf(
				"failed to format JSON: %w",
				err,
			)
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
func (c *ListCmd) listAll(
	lister *list.Lister,
	projectPath string,
) error {
	// Retrieve all items (changes and specs) from the project
	items, err := lister.ListAll(nil)
	if err != nil {
		return fmt.Errorf(
			"failed to list all items: %w",
			err,
		)
	}

	// Handle interactive mode - shows a unified navigable table
	if c.Interactive {
		if len(items) == 0 {
			fmt.Println("No items found.")

			return nil
		}

		return list.RunInteractiveAll(
			items,
			projectPath,
			c.Stdout,
		)
	}

	// Format output based on flags
	var output string
	switch {
	case c.JSON:
		// JSON format for machine consumption
		var err error
		output, err = list.FormatAllJSON(items)
		if err != nil {
			return fmt.Errorf(
				"failed to format JSON: %w",
				err,
			)
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
