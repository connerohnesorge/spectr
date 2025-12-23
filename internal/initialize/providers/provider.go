// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider returns a list of Initializers that configure the tool
// for use with spectr.
//
// # Adding a New Provider
//
// To add a new AI CLI provider, create a new file
// (e.g., providers/mytool.go) with:
//
// Example:
//
//	package providers
//
//	func init() {
//		err := Register(Registration{
//			ID:       "mytool",
//			Name:     "MyTool",
//			Priority: 100,
//			Provider: &MyToolProvider{},
//		})
//		if err != nil {
//			panic("failed to register mytool provider: " + err.Error())
//		}
//	}
//
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) Initializers(ctx context.Context) []Initializer {
//		return []Initializer{
//			NewDirectoryInitializer(false, ".mytool/commands/spectr"),
//			NewConfigFileInitializer("MYTOOL.md", "instruction-pointer", false),
//			NewSlashCommandsInitializer(
//				".mytool/commands/spectr",
//				".md",
//				FormatMarkdown,
//				StandardFrontmatter(),
//				false,
//			),
//		}
//	}
//
// See also:
//   - provider_new.go: Provider interface definition
//   - registration.go: Registration struct for provider metadata
//   - initializer.go: Initializer interface for file operations
//   - builtins.go: Built-in initializer implementations
//
//nolint:revive // line-length-limit - provider documentation
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

// TemplateContext holds path-related template variables for dynamic
// directory names. This struct is defined in the providers package
// to avoid import cycles.
type TemplateContext struct {
	// BaseDir is the base directory for spectr files (default: "spectr")
	BaseDir string
	// SpecsDir is the directory for spec files (default: "spectr/specs")
	SpecsDir string
	// ChangesDir is the directory for change proposals
	ChangesDir string
	// ProjectFile is the path to the project configuration file
	ProjectFile string
	// AgentsFile is the path to the agents file
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

// TemplateRenderer provides template rendering capabilities.
//
// This interface allows providers to render templates without depending on the
// full TemplateManager.
type TemplateRenderer interface {
	// RenderAgents renders the AGENTS.md template content.
	RenderAgents(
		ctx TemplateContext,
	) (string, error)
	// RenderInstructionPointer renders a short pointer template that directs
	// AI assistants to read spectr/AGENTS.md for full instructions.
	RenderInstructionPointer(
		ctx TemplateContext,
	) (string, error)
	// RenderSlashCommand renders a slash command template (proposal or apply).
	RenderSlashCommand(
		command string,
		ctx TemplateContext,
	) (string, error)
}
