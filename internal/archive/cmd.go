// Package archive provides command structures and execution logic
// for archiving completed changes.
package archive

import (
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/config"
)

// ArchiveCmd represents the archive command configuration
type ArchiveCmd struct {
	ChangeID   string `arg:"" optional:"" predictor:"changeID"`
	Yes        bool   `                                        name:"yes"         short:"y" help:"Skip confirmation"` //nolint:lll,revive // Kong struct tag with alignment
	SkipSpecs  bool   `                                        name:"skip-specs"            help:"Skip spec updates"` //nolint:lll,revive // Kong struct tag with alignment
	NoValidate bool   `                                        name:"no-validate"           help:"Skip validation"`   //nolint:lll,revive // Kong struct tag with alignment
}

// Run executes the archive command
func (c *ArchiveCmd) Run() error {
	// Load config from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Pass empty string to use current working directory
	// Result is discarded for CLI usage - already prints to terminal
	_, err = Archive(c, "", cfg.Dir)
	if err != nil {
		return fmt.Errorf(
			"archive failed: %w",
			err,
		)
	}

	return nil
}
