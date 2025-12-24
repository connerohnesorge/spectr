// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider handles both its instruction file (e.g., CLAUDE.md) and slash
// commands (e.g., .claude/commands/) in a single implementation.
//
// # Adding a New Provider
//
// To add a new AI CLI provider, create a new file
// (e.g., providers/mytools.go) with:
//
// Example:
//
//	package providers
//
//	func init() {
//		Register(&MyToolProvider{})
//	}
//
//	type MyToolProvider struct {
//		BaseProvider
//	}
//
//	func NewMyToolProvider() *MyToolProvider {
//		proposalPath, applyPath := StandardCommandPaths(
//			".mytool/commands", ".md",
//		)
//
//		return &MyToolProvider{
//			BaseProvider: BaseProvider{
//				id:            "mytool",
//				name:          "MyTool",
//				priority:      100,
//				// Empty if no instruction file
//				configFile:    "MYTOOL.md",
//				// Empty if no slash commands
//				proposalPath:  proposalPath,
//				applyPath:     applyPath,
//				commandFormat: FormatMarkdown,
//				frontmatter: map[string]string{
//					"proposal":
//					"---\ndescription: Scaffold a new Spectr change.\n---",
//					"apply":
//					"---\ndescription: Implement an \n---",
//				},
//			},
//		}
//	}
//
// The BaseProvider handles all common logic.
// Override Configure() only for special formats.
//
//nolint:revive // File length acceptable for provider interface definition
package providers

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
