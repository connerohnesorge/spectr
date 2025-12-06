//nolint:revive // file-length-limit - PR subcommands are logically cohesive
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/git"
	"github.com/connerohnesorge/spectr/internal/list"
	"github.com/connerohnesorge/spectr/internal/pr"
)

// PRCmd represents the pr command with subcommands.
type PRCmd struct {
	Archive  PRArchiveCmd  `cmd:"" aliases:"a" help:"Archive and create PR"`
	Proposal PRProposalCmd `cmd:"" aliases:"p" help:"Create proposal PR"`
	Remove PRRemoveCmd `cmd:"" name:"rm" aliases:"r,remove" help:"Remove via PR"`
}

// PRArchiveCmd represents the pr archive subcommand.
type PRArchiveCmd struct {
	ChangeID  string `arg:"" optional:"" predictor:"changeID" help:"Change ID"`
	Base      string `name:"base" short:"b" help:"Target branch for PR"`
	Draft     bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force     bool   `name:"force" short:"f" help:"Delete existing branch"`
	DryRun    bool   `name:"dry-run" help:"Preview without executing"`
	SkipSpecs bool   `name:"skip-specs" help:"Skip spec merging"`
}

// PRProposalCmd represents the pr proposal subcommand.
type PRProposalCmd struct {
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change ID"`
	Base     string `name:"base" short:"b" help:"Target branch for PR"`
	Draft    bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force    bool   `name:"force" short:"f" help:"Delete existing branch"`
	DryRun   bool   `name:"dry-run" help:"Preview without executing"`
}

// PRRemoveCmd represents the pr remove subcommand.
type PRRemoveCmd struct {
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change ID"`
	Base     string `name:"base" short:"b" help:"Target branch for PR"`
	Draft    bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force    bool   `name:"force" short:"f" help:"Delete existing branch"`
	DryRun   bool   `name:"dry-run" help:"Preview without executing"`
}

// Run executes the pr remove command.
func (c *PRRemoveCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	changeID, err := resolveOrSelectChangeID(c.ChangeID, projectRoot)
	if err != nil {
		if errors.Is(err, archive.ErrUserCancelled) {
			return nil // User cancelled, exit gracefully
		}

		return err
	}

	config := pr.PRConfig{
		ChangeID:    changeID,
		Mode:        pr.ModeRemove,
		BaseBranch:  c.Base,
		Draft:       c.Draft,
		Force:       c.Force,
		DryRun:      c.DryRun,
		ProjectRoot: projectRoot,
	}

	result, err := pr.ExecutePR(config)
	if err != nil {
		return fmt.Errorf("pr remove failed: %w", err)
	}

	printPRResult(result)

	return nil
}

// Run executes the pr archive command.
func (c *PRArchiveCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	changeID, err := resolveOrSelectChangeID(c.ChangeID, projectRoot)
	if err != nil {
		if errors.Is(err, archive.ErrUserCancelled) {
			return nil // User cancelled, exit gracefully
		}

		return err
	}

	config := pr.PRConfig{
		ChangeID:    changeID,
		Mode:        pr.ModeArchive,
		BaseBranch:  c.Base,
		Draft:       c.Draft,
		Force:       c.Force,
		DryRun:      c.DryRun,
		SkipSpecs:   c.SkipSpecs,
		ProjectRoot: projectRoot,
	}

	result, err := pr.ExecutePR(config)
	if err != nil {
		return fmt.Errorf("pr archive failed: %w", err)
	}

	printPRResult(result)

	return nil
}

// Run executes the pr proposal command.
func (c *PRProposalCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	var changeID string

	// For proposal command without explicit ID, filter to unmerged changes only
	if c.ChangeID == "" {
		changeID, err = selectChangeForProposal(projectRoot, c.Base)
	} else {
		// Explicit ID provided - resolve without filtering
		changeID, err = resolveOrSelectChangeID(c.ChangeID, projectRoot)
	}

	if err != nil {
		if errors.Is(err, archive.ErrUserCancelled) {
			return nil // User cancelled, exit gracefully
		}

		return err
	}

	config := pr.PRConfig{
		ChangeID:    changeID,
		Mode:        pr.ModeProposal,
		BaseBranch:  c.Base,
		Draft:       c.Draft,
		Force:       c.Force,
		DryRun:      c.DryRun,
		ProjectRoot: projectRoot,
	}

	result, err := pr.ExecutePR(config)
	if err != nil {
		return fmt.Errorf("pr proposal failed: %w", err)
	}

	printPRResult(result)

	return nil
}

// resolveOrSelectChangeID resolves a partial change ID or prompts for
// interactive selection if no ID is provided.
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

// selectChangeInteractive prompts the user to select a change interactively.
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
// filtering out changes that already exist on the base branch.
// This ensures only unmerged proposals are shown for PR creation.
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

	// Fetch origin to ensure refs are current
	if err := git.FetchOrigin(); err != nil {
		return "", fmt.Errorf("fetch origin: %w", err)
	}

	// Determine the base branch ref
	ref, err := git.GetBaseBranch(baseBranch)
	if err != nil {
		return "", fmt.Errorf("get base branch: %w", err)
	}

	// Filter to only show changes not already on the base branch
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

// printPRResult displays the result of a PR operation.
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
