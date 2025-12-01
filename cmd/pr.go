package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/list"
	"github.com/connerohnesorge/spectr/internal/pr"
)

// PRCmd represents the pr command with its subcommands
type PRCmd struct {
	Archive PRArchiveCmd `cmd:"" help:"Archive a change and create a PR"`
	New     PRNewCmd     `cmd:"" help:"Create a PR for a change proposal"`
}

// PRArchiveCmd creates a PR from an archived change
//
//nolint:revive // line-length-limit - struct tags cannot be split
type PRArchiveCmd struct {
	ChangeID  string `arg:"" optional:"" predictor:"changeID" help:"Change ID to archive (supports partial matching)"`
	Base      string `name:"base" short:"b" help:"Base branch for PR (default: auto-detect main/master)"`
	Draft     bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force     bool   `name:"force" short:"f" help:"Force overwrite existing branch"`
	DryRun    bool   `name:"dry-run" help:"Show what would be done without executing"`
	SkipSpecs bool   `name:"skip-specs" help:"Skip spec merging during archive"`
}

// PRNewCmd creates a PR for a change proposal without archiving
//
//nolint:revive // line-length-limit - struct tags cannot be split
type PRNewCmd struct {
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change ID (supports partial matching)"`
	Base     string `name:"base" short:"b" help:"Base branch for PR (default: auto-detect main/master)"`
	Draft    bool   `name:"draft" short:"d" help:"Create as draft PR"`
	Force    bool   `name:"force" short:"f" help:"Force overwrite existing branch"`
	DryRun   bool   `name:"dry-run" help:"Show what would be done without executing"`
}

// Run executes the pr archive command
func (c *PRArchiveCmd) Run() error {
	changeID, err := resolveChangeIDForPR(c.ChangeID)
	if err != nil {
		return err
	}

	cfg := pr.Config{
		ChangeID:   changeID,
		BaseBranch: c.Base,
		Draft:      c.Draft,
		Force:      c.Force,
		DryRun:     c.DryRun,
		SkipSpecs:  c.SkipSpecs,
	}

	result, err := pr.ExecuteArchivePR(cfg)
	if err != nil {
		return fmt.Errorf("archive PR failed: %w", err)
	}

	if !c.DryRun {
		fmt.Printf("\nPR created: %s\n", result.PRURL)
	}

	return nil
}

// Run executes the pr new command
func (c *PRNewCmd) Run() error {
	changeID, err := resolveChangeIDForPR(c.ChangeID)
	if err != nil {
		return err
	}

	cfg := pr.Config{
		ChangeID:   changeID,
		BaseBranch: c.Base,
		Draft:      c.Draft,
		Force:      c.Force,
		DryRun:     c.DryRun,
	}

	result, err := pr.ExecuteNewPR(cfg)
	if err != nil {
		return fmt.Errorf("new PR failed: %w", err)
	}

	if !c.DryRun {
		fmt.Printf("\nPR created: %s\n", result.PRURL)
	}

	return nil
}

// resolveChangeIDForPR resolves the change ID.
// Uses interactive selection if not provided.
func resolveChangeIDForPR(changeID string) (string, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	// Check if spectr directory exists
	spectrRoot := filepath.Join(projectRoot, "spectr")
	if _, err := os.Stat(spectrRoot); os.IsNotExist(err) {
		return "", fmt.Errorf("spectr directory not found in %s", projectRoot)
	}

	// If no change ID provided, use interactive selection
	if changeID == "" {
		return selectChangeInteractiveForPR(projectRoot)
	}

	// Resolve partial ID to full change ID
	result, err := discovery.ResolveChangeID(changeID, projectRoot)
	if err != nil {
		return "", err
	}

	if result.PartialMatch {
		fmt.Printf("Resolved '%s' -> '%s'\n\n", changeID, result.ChangeID)
	}

	return result.ChangeID, nil
}

// selectChangeInteractiveForPR uses the interactive table for change selection
func selectChangeInteractiveForPR(projectRoot string) (string, error) {
	lister := list.NewLister(projectRoot)
	changes, err := lister.ListChanges()
	if err != nil {
		return "", fmt.Errorf("list changes: %w", err)
	}

	if len(changes) == 0 {
		return "", errors.New("no changes found")
	}

	selectedID, err := list.RunInteractiveArchive(changes, projectRoot)
	if err != nil {
		return "", fmt.Errorf("interactive selection: %w", err)
	}

	if selectedID == "" {
		return "", errors.New("no change selected")
	}

	return selectedID, nil
}
