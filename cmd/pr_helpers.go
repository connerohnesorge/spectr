// Package cmd provides the command-line interface for spectr.
// This file contains helper functions for the PR command, including
// change ID resolution and interactive selection utilities.
package cmd

import (
	"fmt"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/git"
	"github.com/connerohnesorge/spectr/internal/list"
	"github.com/connerohnesorge/spectr/internal/pr"
)

// resolveOrSelectChangeID resolves a partial change ID or prompts for
// interactive selection if no ID is provided. It uses fuzzy matching
// to allow users to specify partial IDs for convenience.
func resolveOrSelectChangeID(changeID, projectRoot string) (string, error) {
	if changeID == "" {
		return selectChangeInteractive(projectRoot)
	}

	result, err := discovery.ResolveChangeID(changeID, projectRoot)
	if err != nil {
		return "", err
	}

	if result.PartialMatch {
		fmt.Printf("Resolved '%s' -> '%s'\n\n", changeID, result.ChangeID)
	}

	return result.ChangeID, nil
}

// selectChangeInteractive prompts user to select a change interactively.
// It displays a list of all available changes and allows the user to
// navigate and select one using a TUI interface.
func selectChangeInteractive(projectRoot string) (string, error) {
	lister := list.NewLister(projectRoot)

	changes, err := lister.ListChanges()
	if err != nil {
		return "", fmt.Errorf("list changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("No changes found.")

		return "", archive.ErrUserCancelled
	}

	selectedID, err := list.RunInteractiveArchive(changes, projectRoot)
	if err != nil {
		return "", fmt.Errorf("interactive selection: %w", err)
	}

	if selectedID == "" {
		return "", archive.ErrUserCancelled
	}

	return selectedID, nil
}

// selectChangeForProposal prompts the user to select a change interactively,
// filtering out changes that already exist on the base branch. This ensures
// that users only see proposals that haven't been merged yet, making the
// selection more focused and relevant for PR creation.
func selectChangeForProposal(projectRoot, baseBranch string) (string, error) {
	lister := list.NewLister(projectRoot)

	changes, err := lister.ListChanges()
	if err != nil {
		return "", fmt.Errorf("list changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("No changes found.")

		return "", archive.ErrUserCancelled
	}

	if err := git.FetchOrigin(); err != nil {
		return "", fmt.Errorf("fetch origin: %w", err)
	}

	ref, err := git.GetBaseBranch(baseBranch)
	if err != nil {
		return "", fmt.Errorf("get base branch: %w", err)
	}

	unmergedChanges, err := list.FilterChangesNotOnRef(changes, ref)
	if err != nil {
		return "", fmt.Errorf("filter changes: %w", err)
	}

	if len(unmergedChanges) == 0 {
		fmt.Println("No unmerged proposals found. " +
			"All changes already exist on main.")

		return "", archive.ErrUserCancelled
	}

	selectedID, err := list.RunInteractiveArchive(unmergedChanges, projectRoot)
	if err != nil {
		return "", fmt.Errorf("interactive selection: %w", err)
	}

	if selectedID == "" {
		return "", archive.ErrUserCancelled
	}

	return selectedID, nil
}

// printPRResult displays the result of a PR operation to the console.
// It shows the branch name, archive path (if applicable), spec operations
// counts, and either the PR URL or a manual URL for creating the PR.
func printPRResult(result *pr.PRResult) {
	fmt.Println()
	fmt.Printf("Branch: %s\n", result.BranchName)

	if result.ArchivePath != "" {
		fmt.Printf("Archived to: %s\n", result.ArchivePath)
	}

	if result.Counts.Total() > 0 {
		fmt.Printf("Spec operations: +%d ~%d -%d\n",
			result.Counts.Added,
			result.Counts.Modified,
			result.Counts.Removed)
	}

	if result.PRURL != "" {
		fmt.Printf("\nPR created: %s\n", result.PRURL)
	} else if result.ManualURL != "" {
		fmt.Printf("\nCreate PR manually: %s\n", result.ManualURL)
	}
}
