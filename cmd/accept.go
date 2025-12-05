package cmd

import (
	"fmt"

	"github.com/connerohnesorge/spectr/internal/accept"
)

// AcceptCmd wraps the accept package's AcceptCmd type to add Run method.
type AcceptCmd struct {
	accept.AcceptCmd
}

// Run executes the accept command.
func (c *AcceptCmd) Run() error {
	err := accept.Accept(&c.AcceptCmd, "")
	if err != nil {
		return fmt.Errorf("accept failed: %w", err)
	}

	return nil
}
