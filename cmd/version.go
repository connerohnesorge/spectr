// Package cmd provides command-line interface implementations for Spectr.
// This file contains the version command for displaying version and build
// information.
package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/connerohnesorge/spectr/internal/version"
)

// VersionCmd represents the version command which displays version and
// build information. It supports multiple output formats: human-readable
// (default), short (version only), and JSON (machine-readable).
//
// Examples:
//
//	spectr version              # Show full version information
//	spectr version --short      # Show only version number
//	spectr version --json       # Output as JSON
type VersionCmd struct {
	// Short outputs only the version number for use in scripts
	Short bool `name:"short" help:"Output only the version number"`
	// JSON outputs version information in JSON format
	JSON bool `name:"json" help:"Output as JSON"`
}

// Run executes the version command.
// It reads the embedded VERSION file, parses it, and displays the information
// in the requested format (default, short, or JSON).
//
// If the VERSION file is missing or cannot be parsed, it falls back to
// development defaults (version="dev", commit="unknown", date="unknown").
func (c *VersionCmd) Run() error {
	// Get version information from embedded file
	info := version.Get()

	// Handle short output - version number only
	if c.Short {
		fmt.Println(info.Version)

		return nil
	}

	// Handle JSON output - machine-readable format
	if c.JSON {
		return c.outputJSON(info)
	}

	// Default human-readable output
	c.outputDefault(info)

	return nil
}

// outputDefault displays version information in human-readable format.
// This includes version number, commit hash, build date, Go version,
// OS, and architecture.
func (*VersionCmd) outputDefault(info version.Info) {
	fmt.Printf("Spectr version %s\n", info.Version)
	fmt.Printf("  Commit:      %s\n", info.Commit)
	fmt.Printf("  Build date:  %s\n", info.Date)
	fmt.Printf("  Go version:  %s\n", runtime.Version())
	fmt.Printf("  OS/Arch:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// outputJSON displays version information in JSON format.
// The JSON includes all fields from the VERSION file plus runtime
// information (Go version, OS, and architecture).
func (*VersionCmd) outputJSON(info version.Info) error {
	// Create output structure with all fields
	output := map[string]string{
		"version":   info.Version,
		"commit":    info.Commit,
		"date":      info.Date,
		"goVersion": runtime.Version(),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
	}

	// Marshal to JSON with indentation for readability
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Display JSON output
	fmt.Println(string(jsonOutput))

	return nil
}
