// Package cmd provides the command-line interface for spectr.
// This file handles cleanup of local change proposals after PR creation.
package cmd

import (
	"fmt"

	"github.com/connerohnesorge/spectr/internal/pr"
	"github.com/connerohnesorge/spectr/internal/tui"
)

// promptForCleanup shows a TUI prompt asking if local change should be removed.
// Returns true if user confirms removal, false otherwise.
func promptForCleanup() bool {
	picker := tui.NewConfirmPicker(tui.ConfirmConfig{
		Question:  "Remove local change proposal from spectr/changes/?",
		DefaultNo: true,
	})

	confirmed, err := picker.Run()
	if err != nil {
		return false
	}

	return confirmed
}

// removeLocalChange removes the local change directory and prints status.
func removeLocalChange(projectRoot, changeID string) {
	if err := pr.RemoveChangeDirectory(projectRoot, changeID); err != nil {
		fmt.Printf("Failed to remove: %v\n", err)

		return
	}

	fmt.Printf("Removed local change: %s\n", changeID)
}

// keepLocalChange prints a message indicating the local change was kept.
func keepLocalChange(changeID string) {
	fmt.Printf("Kept local change: %s\n", changeID)
}

// cleanupWithPrompt shows a confirmation prompt and removes if confirmed.
// It handles the interactive cleanup flow after successful PR creation.
func cleanupWithPrompt(projectRoot, changeID string) {
	if promptForCleanup() {
		removeLocalChange(projectRoot, changeID)
	} else {
		keepLocalChange(changeID)
	}
}
