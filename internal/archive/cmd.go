// Package archive provides command structures and execution logic
// for archiving completed changes.
package archive

import "fmt"

// ArchiveCmd represents the archive command configuration
type ArchiveCmd struct {
	ChangeID   string `arg:"" optional:"" help:"Change ID to archive"`
	Yes        bool   `name:"yes" short:"y" help:"Skip confirmation"`
	SkipSpecs  bool   `name:"skip-specs" help:"Skip spec updates"`
	NoValidate bool   `name:"no-validate" help:"Skip validation"`
	PR         bool   `name:"pr" help:"Create PR after archive"`
}

// Run executes the archive command
func (c *ArchiveCmd) Run() error {
	// Pass empty string to use current working directory
	err := Archive(c, "")
	if err != nil {
		return fmt.Errorf("archive failed: %w", err)
	}

	return nil
}
