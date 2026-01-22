package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/sync"
	kongcompletion "github.com/jotaen/kong-completion"
)

// CLI represents the root command structure for Kong
type CLI struct {
	// Global flags (apply to all commands)
	NoSync  bool `help:"Skip automatic task sync" name:"no-sync" short:"S"` //nolint:lll,revive // Kong struct tag
	Verbose bool `help:"Enable verbose output"    name:"verbose" short:"v"` //nolint:lll,revive // Kong struct tag

	// Commands
	Init       InitCmd                   `cmd:"" help:"Initialize Spectr"`                 //nolint:lll,revive // Kong struct tag with alignment
	List       ListCmd                   `cmd:"" help:"List items"           aliases:"ls"` //nolint:lll,revive // Kong struct tag with alignment
	Validate   ValidateCmd               `cmd:"" help:"Validate items"`                    //nolint:lll,revive // Kong struct tag with alignment
	Accept     AcceptCmd                 `cmd:"" help:"Accept tasks.md"`                   //nolint:lll,revive // Kong struct tag with alignment
	Archive    archive.ArchiveCmd        `cmd:"" help:"Archive a change"`                  //nolint:lll,revive // Kong struct tag with alignment
	PR         PRCmd                     `cmd:"" help:"Create pull requests"`              //nolint:lll,revive // Kong struct tag with alignment
	View       ViewCmd                   `cmd:"" help:"Display dashboard"`                 //nolint:lll,revive // Kong struct tag with alignment
	Ralph      RalphCmd                  `cmd:"" help:"Orchestrate task execution"`        //nolint:lll,revive // Kong struct tag with alignment
	Version    VersionCmd                `cmd:"" help:"Show version info"`                 //nolint:lll,revive // Kong struct tag with alignment
	Completion kongcompletion.Completion `cmd:"" help:"Generate completions"`              //nolint:lll,revive // Kong struct tag with alignment
}

// AfterApply is called by Kong after parsing flags but before running the command.
// It synchronizes task statuses from tasks.jsonc to tasks.md for all active changes.
func (c *CLI) AfterApply() error {
	if c.NoSync {
		return nil
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		// Log error but don't block command
		fmt.Fprintf(
			os.Stderr,
			"sync: failed to get working directory: %v\n",
			err,
		)

		return nil
	}

	// Check if spectr/ directory exists (not initialized = skip)
	spectrDir := filepath.Join(
		projectRoot,
		"spectr",
	)
	if _, err := os.Stat(spectrDir); os.IsNotExist(
		err,
	) {
		return nil
	}

	return sync.SyncAllActiveChanges(
		projectRoot,
		c.Verbose,
	)
}
