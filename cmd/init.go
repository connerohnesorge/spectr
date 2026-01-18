package cmd

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/connerohnesorge/spectr/internal/initialize"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// InitCmd wraps the initialize package's InitCmd type to add Run method
type InitCmd struct {
	initialize.InitCmd
}

// Run executes the init command
func (c *InitCmd) Run() error {
	// Register all providers with error handling
	if err := providers.RegisterAllProviders(); err != nil {
		return fmt.Errorf(
			"failed to register providers: %w",
			err,
		)
	}

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
	expandedPath, err := initialize.ExpandPath(
		projectPath,
	)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Update the command with resolved path
	c.Path = expandedPath

	// Check if already initialized
	if c.NonInteractive &&
		initialize.IsSpectrInitialized(
			expandedPath,
		) {
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

// isTTYError checks if an error is related to TTY unavailability.
// TTY errors occur when the Bubbletea TUI framework cannot access /dev/tty,
// typically in CI environments, Docker containers, or piped commands.
//
// Uses idiomatic Go error handling with errors.Is/errors.As to check for:
//   - syscall.ENXIO: "No such device or address" (most common in CI/Docker)
//   - syscall.ENOTTY: "Inappropriate ioctl for device"
//   - os.PathError with TTY-related paths (/dev/tty on Unix, CONIN$ on Windows)
func isTTYError(err error) bool {
	if err == nil {
		return false
	}

	// Check underlying syscall errors (works through wrapped errors)
	if errors.Is(err, syscall.ENXIO) ||
		errors.Is(err, syscall.ENOTTY) {
		return true
	}

	// Check for PathError with TTY-related paths
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		if pathErr.Path == "/dev/tty" ||
			pathErr.Path == "CONIN$" {
			return true
		}
	}

	return false
}

func runInteractiveInit(c *InitCmd) error {
	model, err := initialize.NewWizardModel(
		&c.InitCmd,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create wizard: %w",
			err,
		)
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	finalModel, err := p.Run()
	if err != nil {
		// Check if this is a TTY-related error
		if isTTYError(err) {
			return fmt.Errorf(
				"wizard failed: %w\n\nHint: It looks like this command is running without a TTY (terminal).\nTry using non-interactive mode instead:\n  spectr init --non-interactive --tools <tool1,tool2>",
				err,
			)
		}

		return fmt.Errorf(
			"wizard failed: %w",
			err,
		)
	}

	// Check if there were errors during execution
	wizardModel, ok := finalModel.(*initialize.WizardModel)
	if !ok {
		return &specterrs.WizardModelCastError{}
	}
	err = wizardModel.GetError()
	if err != nil {
		return err
	}

	return nil
}

func runNonInteractiveInit(c *InitCmd) error {
	// Handle "all" special case
	selectedProviders := c.Tools
	if len(c.Tools) == 1 && c.Tools[0] == "all" {
		// Get all provider IDs from registered providers
		allProviders := providers.RegisteredProviders()
		selectedProviders = make(
			[]string,
			len(allProviders),
		)
		for i, reg := range allProviders {
			selectedProviders[i] = reg.ID
		}
	}

	// Validate provider IDs
	for _, id := range selectedProviders {
		if _, exists := providers.Get(id); !exists {
			return fmt.Errorf(
				"invalid provider ID: %s",
				id,
			)
		}
	}

	// Create executor and run
	executor, err := initialize.NewInitExecutor(
		&c.InitCmd,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create executor: %w",
			err,
		)
	}

	result, err := executor.Execute(
		selectedProviders,
		c.CIWorkflow,
	)
	if err != nil {
		return fmt.Errorf(
			"initialization failed: %w",
			err,
		)
	}

	return printInitResults(c.Path, result)
}

func printInitResults(
	projectPath string,
	result *initialize.ExecutionResult,
) error {
	fmt.Println(
		"Spectr initialized successfully!",
	)
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

		// Convert string errors to error type
		errs := make([]error, len(result.Errors))
		for i, e := range result.Errors {
			errs[i] = errors.New(e)
		}

		return &specterrs.InitializationCompletedWithErrorsError{
			ErrorCount: len(result.Errors),
			Errors:     errs,
		}
	}

	fmt.Print(initialize.FormatNextStepsMessage())

	return nil
}
