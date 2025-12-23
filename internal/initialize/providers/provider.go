// Package providers implements the provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
//nolint:revive // File length acceptable for provider interface definition
package providers

import "strings"

// CommandFormat specifies the format for slash command files.
type CommandFormat int

const (
	// FormatMarkdown uses markdown files with
	// YAML frontmatter (Claude, Cline, etc.)
	FormatMarkdown CommandFormat = iota
	// FormatTOML uses TOML files (Gemini CLI)
	FormatTOML
)

// TemplateContext holds path-related template variables for dynamic directory names.
// This struct is defined in the providers package to avoid import cycles.
type TemplateContext struct {
	// BaseDir is the base directory for spectr files (default: "spectr")
	BaseDir string
	// SpecsDir is the directory for spec files (default: "spectr/specs")
	SpecsDir string
	// ChangesDir is the directory for change proposals (default: "spectr/changes")
	ChangesDir string
	// ProjectFile is the path to the project configuration file (default: "spectr/project.md")
	ProjectFile string
	// AgentsFile is the path to the agents file (default: "spectr/AGENTS.md")
	AgentsFile string
}

// DefaultTemplateContext returns a TemplateContext with default values.
func DefaultTemplateContext() TemplateContext {
	return TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}
}

// Constants used by initializers
const (
	// File and directory permission constants.
	filePerm = 0o644
	dirPerm  = 0o755

	// Marker constants for managing config file updates.
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
)

// findMarkerIndex finds the index of a marker in content, starting from offset.
func findMarkerIndex(content, marker string, offset int) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}
