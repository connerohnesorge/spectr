// Package archive provides command structures and execution logic
// for archiving completed changes.
package archive

import "fmt"

// ArchiveCmd represents the archive command configuration
type ArchiveCmd struct {
	ChangeID   string `arg:"" optional:"" predictor:"changeID"`
	Yes        bool   `                                        name:"yes"         short:"y" help:"Skip confirmation"` //nolint:lll,revive
	SkipSpecs  bool   `                                        name:"skip-specs"            help:"Skip spec updates"` //nolint:lll,revive
	NoValidate bool   `                                        name:"no-validate"           help:"Skip validation"`   //nolint:lll,revive
}

// Run executes the archive command
func (c *ArchiveCmd) Run() error {
	// Pass empty string to use current working directory
	// Result is discarded for CLI usage - already prints to terminal
	_, err := Archive(c, "")
	if err != nil {
		return fmt.Errorf(
			"archive failed: %w",
			err,
		)
	}

	return nil
}
