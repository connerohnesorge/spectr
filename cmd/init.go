package cmd

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	initpkg "github.com/connerohnesorge/spectr/internal/init"
)

// InitCmd wraps the init package's InitCmd type to add Run method
type InitCmd struct {
	initpkg.InitCmd
}

// Run executes the init command
func (c *InitCmd) Run() error {
	// Determine project path - positional arg takes precedence over flag
	projectPath := c.Path
	if projectPath == "" {
		projectPath = c.PathFlag
	}
	if projectPath == "" {
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf(
				"failed to get current directory: %w",
				err,
			)
		}
	}

	// Expand and validate path
	expandedPath, err := initpkg.ExpandPath(projectPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Update the command with resolved path
	c.Path = expandedPath

	// Check if already initialized
	if c.NonInteractive && initpkg.IsSpectrInitialized(expandedPath) {
		return fmt.Errorf(
			"init failed: Spectr is already initialized in %s",
			expandedPath,
		)
	}

	// Non-interactive mode
	if c.NonInteractive {
		return runNonInteractiveInit(c)
	}

	// Interactive mode (TUI wizard)
	return runInteractiveInit(c)
}

func runInteractiveInit(c *InitCmd) error {
	model, err := initpkg.NewWizardModel(&c.InitCmd)
	if err != nil {
		return fmt.Errorf("failed to create wizard: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("wizard failed: %w", err)
	}

	// Check if there were errors during execution
	wizardModel, ok := finalModel.(initpkg.WizardModel)
	if !ok {
		return errors.New("failed to cast final model to WizardModel")
	}
	err = wizardModel.GetError()
	if err != nil {
		return err
	}

	return nil
}

func runNonInteractiveInit(c *InitCmd) error {
	// Get registry
	registry := initpkg.NewRegistry()

	// Handle "all" special case
	selectedTools := c.Tools
	if len(c.Tools) == 1 && c.Tools[0] == "all" {
		allToolIDs := registry.ListTools()
		selectedTools = make([]string, len(allToolIDs))
		for i, toolID := range allToolIDs {
			selectedTools[i] = string(toolID)
		}
	}

	// Validate tool IDs
	for _, id := range selectedTools {
		if _, err := registry.GetTool(initpkg.ToolID(id)); err != nil {
			return fmt.Errorf("invalid tool ID: %s", id)
		}
	}

	// Create executor and run
	executor, err := initpkg.NewInitExecutor(&c.InitCmd)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	result, err := executor.Execute(selectedTools)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	return printInitResults(c.Path, result)
}

func printInitResults(
	projectPath string,
	result *initpkg.ExecutionResult,
) error {
	fmt.Println("Spectr initialized successfully!")
	fmt.Printf("Project: %s\n\n", projectPath)

	if len(result.CreatedFiles) > 0 {
		fmt.Println("Created files:")
		for _, file := range result.CreatedFiles {
			fmt.Printf("  ✓ %s\n", file)
		}
		fmt.Println()
	}

	if len(result.UpdatedFiles) > 0 {
		fmt.Println("Updated files:")
		for _, file := range result.UpdatedFiles {
			fmt.Printf("  ✓ %s\n", file)
		}
		fmt.Println()
	}

	if len(result.Errors) > 0 {
		fmt.Println("Errors:")
		for _, e := range result.Errors {
			fmt.Printf("  ✗ %s\n", e)
		}

		return errors.New("initialization completed with errors")
	}

	fmt.Print(initpkg.FormatNextStepsMessage())

	return nil
}
