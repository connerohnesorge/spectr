// Package cmd provides the command-line interface for spectr.
// This file implements the pr command and its subcommands for creating
// pull requests from change proposals.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/pr"
)

// PRCmd represents the pr command with subcommands.
// It provides archive and proposal modes for creating PRs.
type PRCmd struct {
	Archive  PRArchiveCmd  `cmd:"" aliases:"a" help:"Archive and create PR"`
	Proposal PRProposalCmd `cmd:"" aliases:"p" help:"Create proposal PR"`
}

// PRArchiveCmd represents the pr archive subcommand.
// It archives a change proposal to the specs directory and creates a PR.
type PRArchiveCmd struct {
	ChangeID  string `arg:"" optional:"" predictor:"changeID" help:"Change ID"`
	Base      string `name:"base" short:"b" help:"Target branch for PR"`
	Draft     bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force     bool   `name:"force" short:"f" help:"Delete existing branch"`
	DryRun    bool   `name:"dry-run" help:"Preview without executing"`
	SkipSpecs bool   `name:"skip-specs" help:"Skip spec merging"`
}

// PRProposalCmd represents the pr proposal subcommand.
// It creates a PR for a change proposal without archiving it, useful for
// getting early feedback on a proposal before it's ready to be merged.
type PRProposalCmd struct {
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change ID"`
	Base     string `name:"base" short:"b" help:"Target branch for PR"`
	Draft    bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force    bool   `name:"force" short:"f" help:"Delete existing branch"`
	DryRun   bool   `name:"dry-run" help:"Preview without executing"`
	Yes      bool   `name:"yes" short:"y" help:"Skip prompts (keep change)"`
}

// Run executes the pr archive command. It resolves the change ID,
// creates the PR configuration, and executes the PR workflow.
func (c *PRArchiveCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	changeID, err := resolveOrSelectChangeID(c.ChangeID, projectRoot)
	if err != nil {
		if errors.Is(err, archive.ErrUserCancelled) {
			return nil
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

// Run executes the pr proposal command. It resolves the change ID,
// creates the PR configuration, executes the workflow, and optionally
// prompts for local change cleanup after successful PR creation.
func (c *PRProposalCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	var changeID string

	if c.ChangeID == "" {
		changeID, err = selectChangeForProposal(projectRoot, c.Base)
	} else {
		changeID, err = resolveOrSelectChangeID(c.ChangeID, projectRoot)
	}

	if err != nil {
		if errors.Is(err, archive.ErrUserCancelled) {
			return nil
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

	// Handle cleanup after successful PR (not in dry-run mode)
	if !c.DryRun {
		if c.Yes {
			keepLocalChange(changeID)
		} else {
			cleanupWithPrompt(projectRoot, changeID)
		}
	}

	return nil
}
