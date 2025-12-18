// Package cmd provides CLI command implementations for spectr.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/connerohnesorge/spectr/internal/git"
	"github.com/connerohnesorge/spectr/internal/specterrs"
	"github.com/connerohnesorge/spectr/internal/track"
)

// TrackCmd represents the track command for automatic git commits.
// It watches the tasks.jsonc file for status changes and automatically
// creates commits when tasks transition to "in_progress" or "completed".
type TrackCmd struct {
	// ChangeID is the optional change identifier to track.
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change ID"` //nolint:lll,revive
	// NoInteractive disables interactive prompts for change selection.
	NoInteractive bool `                                        help:"Disable prompts"                 name:"no-interactive"` //nolint:lll,revive
	// IncludeBinaries enables inclusion of binary files in commits.
	// By default, binary files are excluded from automated commits.
	IncludeBinaries bool `                                        help:"Include binary files in commits" name:"include-binaries"` //nolint:lll,revive
}

// Run executes the track command. It resolves the change ID,
// locates the tasks.jsonc file, and starts the tracking loop.
func (c *TrackCmd) Run() error {
	changeID, projectRoot, err := c.resolveChangeID()
	if err != nil {
		return err
	}

	// User cancelled interactive selection
	if changeID == "" {
		return nil
	}

	return c.runTracker(changeID, projectRoot)
}

// resolveChangeID resolves the change ID from command arguments or
// prompts for interactive selection if no ID was provided.
//
//nolint:revive // confusing-results is acceptable here
func (c *TrackCmd) resolveChangeID() (string, string, error) {
	// Check for --no-interactive flag before requiring git repository.
	// This provides a better error message for the user.
	if c.ChangeID == "" && c.NoInteractive {
		return "", "", errors.New(
			"change ID required when --no-interactive is set",
		)
	}

	// Use git repository root to ensure correct paths for git operations.
	// This is necessary because git status returns paths relative to the
	// repo root, and git add expects paths relative to the repo root.
	projectRoot, err := git.GetRepoRoot()
	if err != nil {
		return "", "", fmt.Errorf(
			"get git repository root: %w",
			err,
		)
	}

	var changeID string
	if c.ChangeID == "" {
		// No change ID provided - prompt for selection
		changeID, err = selectChangeInteractive(
			projectRoot,
		)
	} else {
		// Resolve the provided change ID (may be partial)
		changeID, err = resolveOrSelectChangeID(c.ChangeID, projectRoot)
	}

	if err != nil {
		var userCancelledErr *specterrs.UserCancelledError
		if errors.As(err, &userCancelledErr) {
			return "", projectRoot, nil
		}

		return "", "", err
	}

	return changeID, projectRoot, nil
}

// runTracker creates and runs the tracker for the specified change.
// It watches the tasks.jsonc file and creates commits on status changes.
// The tracker runs until all tasks complete, an error occurs, or the
// user interrupts with Ctrl+C.
func (c *TrackCmd) runTracker(
	changeID, projectRoot string,
) error {
	// Build the path to the tasks.jsonc file
	tasksPath := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		changeID,
		"tasks.jsonc",
	)

	// Verify the tasks file exists
	if _, err := os.Stat(tasksPath); os.IsNotExist(
		err,
	) {
		return &specterrs.NoTasksFileError{
			ChangeID: changeID,
		}
	}

	// Create the tracker configuration
	config := track.Config{
		ChangeID:        changeID,
		TasksPath:       tasksPath,
		RepoRoot:        projectRoot,
		Writer:          os.Stdout,
		IncludeBinaries: c.IncludeBinaries,
	}

	// Create and start the tracker
	tracker, err := track.New(config)
	if err != nil {
		return fmt.Errorf(
			"create tracker: %w",
			err,
		)
	}
	defer func() { _ = tracker.Close() }()

	// Set up context for graceful shutdown on Ctrl+C
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	return handleTrackerResult(
		tracker.Run(ctx),
		changeID,
	)
}

// handleTrackerResult processes the result from tracker.Run and
// converts tracker errors into appropriate user-facing messages.
func handleTrackerResult(
	err error,
	changeID string,
) error {
	if err == nil {
		return nil
	}

	// Check for "all tasks complete" - not an error condition
	var tasksCompleteErr *specterrs.TasksAlreadyCompleteError
	if errors.As(err, &tasksCompleteErr) {
		fmt.Printf(
			"All tasks already completed for change %q\n",
			changeID,
		)

		return nil
	}

	// Check for user interrupt - graceful exit
	var interruptedErr *specterrs.TrackInterruptedError
	if errors.As(err, &interruptedErr) {
		fmt.Println("\nTracking stopped")

		return nil
	}

	// Check for git commit failure - propagate error
	var gitErr *specterrs.GitCommitError
	if errors.As(err, &gitErr) {
		return fmt.Errorf(
			"git commit failed: %w",
			gitErr.Unwrap(),
		)
	}

	return err
}
