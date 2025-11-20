package cmd

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	initpkg "github.com/conneroisu/spectr/internal/init"
)

// CLI represents the root command structure for Kong
type CLI struct {
	Init     InitCmd     `cmd:"" help:"Initialize Spectr in a project"`
	List     ListCmd     `cmd:"" help:"List changes or specifications"`
	Validate ValidateCmd `cmd:"" help:"Validate changes or specifications"`
	Archive  ArchiveCmd  `cmd:"" help:"Archive a completed change"`
	View     ViewCmd     `cmd:"" help:"Display project dashboard with overview"`
}

// InitCmd represents the init command with all its flags
type InitCmd struct {
	Path           string   `arg:"" optional:"" help:"Project path"`
	PathFlag       string   `name:"path" short:"p" help:"Alt project path"`
	Tools          []string `name:"tools" short:"t" help:"Tools list"`
	NonInteractive bool     `name:"non-interactive" help:"Non-interactive"`
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
	projectPath = expandedPath

	// Check if already initialized
	if c.NonInteractive && initpkg.IsSpectrInitialized(projectPath) {
		return fmt.Errorf(
			"init failed: Spectr is already initialized in %s",
			projectPath,
		)
	}

	// Non-interactive mode
	if c.NonInteractive {
		return runNonInteractiveInit(projectPath, c.Tools)
	}

	// Interactive mode (TUI wizard)
	return runInteractiveInit(projectPath)
}

func runInteractiveInit(projectPath string) error {
	model, err := initpkg.NewWizardModel(projectPath)
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

func runNonInteractiveInit(projectPath string, toolIDs []string) error {
	// Get registry
	registry := initpkg.NewRegistryFromProviders()

	// Handle "all" special case
	selectedTools := toolIDs
	if len(toolIDs) == 1 && toolIDs[0] == "all" {
		selectedTools = registry.ListTools()
	}

	// Validate tool IDs
	for _, id := range selectedTools {
		if _, err := registry.GetTool(id); err != nil {
			return fmt.Errorf("invalid tool ID: %s", id)
		}
	}

	// Create executor and run
	executor, err := initpkg.NewInitExecutor(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	result, err := executor.Execute(selectedTools)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	return printInitResults(projectPath, result)
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
