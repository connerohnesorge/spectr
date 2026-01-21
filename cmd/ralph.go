// Package cmd provides command-line interface implementations.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/ralph"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// RalphCmd represents the ralph command for orchestrated task execution.
// It automates the execution of tasks from a change proposal by invoking
// an AI agent CLI (Claude, Gemini, etc.) for each task with appropriate context.
type RalphCmd struct {
	// ChangeID is the change proposal identifier to orchestrate
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change proposal ID to orchestrate"`

	// Interactive enables interactive task selection mode
	Interactive bool `short:"i" help:"Interactive task selection mode"`

	// MaxRetries is the maximum number of retry attempts per task
	MaxRetries int `default:"3" help:"Maximum retry attempts per task" name:"max-retries"`

	// NoInteractive disables interactive prompts
	NoInteractive bool `help:"Disable prompts" name:"no-interactive"`
}

// Run executes the ralph command.
// It validates the change ID, detects the configured provider,
// verifies the provider supports Ralpher, creates the orchestrator,
// starts the TUI, and runs the orchestration.
//
//nolint:revive // function-length: orchestration requires cohesive setup steps
func (c *RalphCmd) Run() error {
	// Register all providers
	if err := providers.RegisterAllProviders(); err != nil {
		return fmt.Errorf("failed to register providers: %w", err)
	}

	// Get project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Resolve change ID
	changeID, err := c.resolveChangeID(projectRoot)
	if err != nil {
		var userCancelledErr *specterrs.UserCancelledError
		if errors.As(err, &userCancelledErr) {
			return nil
		}

		return err
	}

	// Validate change directory exists
	changeDir := filepath.Join(projectRoot, "spectr", "changes", changeID)
	if _, err := os.Stat(changeDir); os.IsNotExist(err) {
		return fmt.Errorf(
			"change directory not found: %s\n\n"+
				"Available changes (use 'spectr list' to see all):\n"+
				"Run 'spectr list' to view all active changes",
			changeDir,
		)
	}

	// Detect provider and verify Ralpher support
	provider, err := c.detectAndValidateProvider()
	if err != nil {
		return err
	}

	// Parse task graph to get tasks for TUI
	graph, err := ralph.ParseTaskGraph(changeDir)
	if err != nil {
		return fmt.Errorf(
			"failed to parse task graph: %w\n\n"+
				"Possible causes:\n"+
				"  - tasks.jsonc file is missing (run 'spectr accept %s')\n"+
				"  - tasks.jsonc has invalid JSON syntax\n"+
				"  - task IDs don't follow expected format (e.g., '1.1', '2.3')\n\n"+
				"Try validating the change first: spectr validate %s",
			err,
			changeID,
			changeID,
		)
	}

	// Create orchestrator with callbacks
	orchestrator, err := c.createOrchestrator(changeID, changeDir, provider)
	if err != nil {
		return fmt.Errorf(
			"failed to create orchestrator: %w\n\n"+
				"This is an internal error. Please check:\n"+
				"  - The change directory exists and is readable\n"+
				"  - The provider is properly configured\n"+
				"  - tasks.jsonc is valid JSON",
			err,
		)
	}

	// Create TUI model
	tuiModel := c.createTUIModel(changeID, graph, orchestrator)

	// Run TUI program
	p := tea.NewProgram(
		tuiModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start orchestrator in background goroutine
	go func() {
		if err := orchestrator.Run(); err != nil {
			// Send error to TUI
			p.Send(ralph.TaskFailMsg{
				TaskID: orchestrator.GetSession().CurrentTaskID,
				Error:  err,
			})
		}
	}()

	// Run TUI
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf(
			"TUI failed: %w\n\n"+
				"If you're experiencing terminal issues:\n"+
				"  - Ensure your terminal supports PTY (most Unix terminals do)\n"+
				"  - Try running with TERM=xterm-256color\n"+
				"  - Check terminal dimensions with 'tput cols' and 'tput lines'\n"+
				"  - Verify terminal emulator compatibility",
			err,
		)
	}

	// Check if orchestration had errors
	_, ok := finalModel.(*ralph.TUIModel)
	if !ok {
		return errors.New("unexpected model type")
	}

	// If user quit, ensure orchestrator is stopped
	if err := orchestrator.Stop(); err != nil {
		return fmt.Errorf(
			"failed to stop orchestrator: %w\n\n"+
				"The orchestration session may have been saved but cleanup failed.\n"+
				"You can try resuming with: spectr ralph %s",
			err,
			changeID,
		)
	}

	return nil
}

// resolveChangeID resolves the change ID from argument or interactive selection.
func (c *RalphCmd) resolveChangeID(projectRoot string) (string, error) {
	if c.ChangeID != "" {
		// Normalize path to extract change ID
		normalizedID, _ := discovery.NormalizeItemPath(c.ChangeID)

		result, err := discovery.ResolveChangeID(normalizedID, projectRoot)
		if err != nil {
			return "", err
		}

		if result.PartialMatch {
			fmt.Printf("Resolved '%s' -> '%s'\n\n", c.ChangeID, result.ChangeID)
		}

		return result.ChangeID, nil
	}

	if c.NoInteractive {
		return "", &specterrs.MissingChangeIDError{}
	}

	return selectChangeInteractive(projectRoot)
}

// detectAndValidateProvider detects the configured provider and validates it supports Ralpher.
func (*RalphCmd) detectAndValidateProvider() (providers.Ralpher, error) {
	// For now, we'll use a heuristic: detect which provider is available
	// In the future, this could read from a config file
	allProviders := providers.RegisteredProviders()

	for _, reg := range allProviders {
		if !providers.IsRalpherAvailable(reg.Provider) {
			continue
		}

		ralpher, ok := reg.Provider.(providers.Ralpher)
		if !ok {
			continue
		}

		fmt.Printf("Using provider: %s (%s)\n", reg.Name, ralpher.Binary())

		return ralpher, nil
	}

	// No provider found
	return nil, errors.New(
		"no suitable provider found for ralph orchestration\n\n" +
			"Ralph requires an AI agent CLI that supports task orchestration.\n" +
			"Please install one of the following:\n" +
			"  - Claude Code (https://claude.ai/code)\n" +
			"  - Gemini CLI (https://ai.google.dev/gemini-api/docs/cli)\n\n" +
			"After installation, ensure the CLI is available in your PATH",
	)
}

// createOrchestrator creates a new orchestrator with TUI callbacks.
func (c *RalphCmd) createOrchestrator(
	changeID, changeDir string,
	provider providers.Ralpher,
) (*ralph.Orchestrator, error) {
	config := ralph.OrchestratorConfig{
		ChangeID:   changeID,
		ChangeDir:  changeDir,
		Provider:   provider,
		MaxRetries: c.MaxRetries,
		// Callbacks are set in createTUIModel
	}

	orchestrator, err := ralph.NewOrchestrator(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create orchestrator: %w", err)
	}

	return orchestrator, nil
}

// createTUIModel creates a TUI model with orchestrator callbacks.
func (c *RalphCmd) createTUIModel(
	changeID string,
	graph *ralph.TaskGraph,
	orchestrator *ralph.Orchestrator,
) *ralph.TUIModel {
	// Convert task map to slice for TUI
	tasks := make([]*ralph.Task, 0, len(graph.Tasks))
	for _, task := range graph.Tasks {
		tasks = append(tasks, task)
	}

	// Create TUI config
	config := ralph.TUIConfig{
		ChangeID:    changeID,
		Tasks:       tasks,
		Interactive: c.Interactive,
		OnUserAction: func(_ ralph.UserAction) tea.Cmd {
			// This will be called when user makes a decision on failure
			return nil
		},
		OnQuit: func() tea.Cmd {
			// Stop orchestrator when user quits
			_ = orchestrator.Stop()

			return tea.Quit
		},
	}

	return ralph.NewTUIModel(&config)
}
