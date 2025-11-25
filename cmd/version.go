// Package cmd provides command-line interface implementations for Spectr.
// This file contains the version command for displaying build information.
package cmd

import (
	"fmt"

	"github.com/connerohnesorge/spectr/internal/version"
)

// VersionCmd represents the version command which displays build information
// including version number, git commit hash, and build date.
//
// Output formats:
//   - Default: Multi-line formatted output with version, commit, and date
//   - --short: Version number only (e.g., "v0.1.0")
//   - --json: Machine-readable JSON for automation and scripting
//
// Examples:
//
//	spectr version              # Full build information
//	spectr version --short      # Version number only
//	spectr version --json       # JSON format
type VersionCmd struct {
	// JSON enables JSON output format for scripting and automation.
	// When enabled, outputs structured data with version, commit, date.
	JSON bool `kong:"help='Output in JSON format for scripting'"`

	// Short enables minimal output showing only the version number.
	// Useful for scripts that need to parse or compare version numbers.
	Short bool `kong:"help='Output version number only'"`
}

// Run executes the version command.
// It retrieves build information and formats the output based on the flags:
// JSON flag takes precedence over Short flag if both are set.
// Returns an error if JSON marshaling fails, nil otherwise.
func (c *VersionCmd) Run() error {
	// Get build information
	info := version.GetBuildInfo()

	// Format output based on flags
	switch {
	case c.JSON:
		// JSON format for machine consumption
		jsonBytes, err := info.JSON()
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonBytes))
	case c.Short:
		// Short format: version number only
		fmt.Println(info.Short())
	default:
		// Default format: multi-line with all build info
		fmt.Println(info.String())
	}

	return nil
}
