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
	Graph      GraphCmd                  `cmd:"" help:"Show dependency graph"`             //nolint:lll,revive // Kong struct tag with alignment
	Hooks      HooksCmd                  `cmd:"" help:"Process hook events"`               //nolint:lll,revive // Kong struct tag with alignment
	PR         PRCmd                     `cmd:"" help:"Create pull requests"`              //nolint:lll,revive // Kong struct tag with alignment
	View       ViewCmd                   `cmd:"" help:"Display dashboard"`                 //nolint:lll,revive // Kong struct tag with alignment
	Version    VersionCmd                `cmd:"" help:"Show version info"`                 //nolint:lll,revive // Kong struct tag with alignment
	Completion kongcompletion.Completion `cmd:"" help:"Generate completions"`              //nolint:lll,revive // Kong struct tag with alignment
}

// AfterApply is called by Kong after parsing flags but before running the command.
// It synchronizes task statuses from tasks.jsonc to tasks.md for all active changes
// across all discovered spectr roots.
func (c *CLI) AfterApply() error {
	if c.NoSync {
		return nil
	}

	// Discover all spectr roots
	roots, err := GetDiscoveredRoots()
	if err != nil {
		// Log error but don't block command
		fmt.Fprintf(
			os.Stderr,
			"sync: failed to discover spectr roots: %v\n",
			err,
		)

		return nil
	}

	// If no roots found, skip sync (not initialized)
	if len(roots) == 0 {
		return nil
	}

	// Sync all active changes across all discovered roots
	for _, root := range roots {
		// Check if spectr/ directory exists for this root
		spectrDir := filepath.Join(root.Path, "spectr")
		if _, statErr := os.Stat(spectrDir); os.IsNotExist(statErr) {
			continue
		}

		syncErr := sync.SyncAllActiveChanges(root.Path, c.Verbose)
		if syncErr == nil {
			continue
		}

		// Log error but continue with other roots
		if c.Verbose {
			fmt.Fprintf(
				os.Stderr,
				"sync: failed for %s: %v\n",
				root.RelativeTo,
				syncErr,
			)
		}
	}

	return nil
}
