// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Providers declare their file initializers as a composable list via the
// Initializers() method. Each initializer is responsible for a single file.
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
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) ID() string       { return "mytool" }
//	func (p *MyToolProvider) Name() string     { return "MyTool" }
//	func (p *MyToolProvider) Priority() int    { return 100 }
//
//	func (p *MyToolProvider) Initializers() []FileInitializer {
//		return []FileInitializer{
//			NewInstructionFileInitializer("MYTOOL.md"),
//			NewMarkdownSlashCommandInitializer(
//				".mytool/commands/spectr/proposal.md",
//				"proposal",
//				FrontmatterProposal,
//			),
//			NewMarkdownSlashCommandInitializer(
//				".mytool/commands/spectr/apply.md",
//				"apply",
//				FrontmatterApply,
//			),
//		}
//	}
//
//	func (p *MyToolProvider) IsConfigured(projectPath string) bool {
//		return AreInitializersConfigured(p.Initializers(), projectPath)
//	}
//
//	func (p *MyToolProvider) GetFilePaths() []string {
//		return GetInitializerPaths(p.Initializers())
//	}
//
// Helper functions ConfigureInitializers(), AreInitializersConfigured(),
// and GetInitializerPaths() provide common functionality.
package providers

// TemplateContext holds path-related template variables for
// dynamic directory names.
// This struct is defined in the providers package to avoid import cycles.
type TemplateContext struct {
	// BaseDir is the base directory for spectr files (default: "spectr")
	BaseDir string
	// SpecsDir is the directory for spec files (default: "spectr/specs")
	SpecsDir string
	// ChangesDir is the directory for change proposals.
	// Default: "spectr/changes"
	ChangesDir string
	// ProjectFile is the path to the project configuration file.
	// Default: "spectr/project.md"
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

// Provider represents an AI CLI tool (Claude Code, Gemini, Cline, etc.).
// Each provider declares its file initializers as a composable list.
type Provider interface {
	// ID returns the unique provider identifier (kebab-case).
	// Example: "claude-code", "gemini", "cline"
	ID() string

	// Name returns the human-readable provider name for display.
	// Example: "Claude Code", "Gemini CLI", "Cline"
	Name() string

	// Priority returns the display order (lower = higher priority).
	// Claude Code should be 1, other major tools 2-10, etc.
	Priority() int

	// Initializers returns the list of file initializers for this provider.
	// Each initializer handles a single file (instruction file,
	// slash command, etc.).
	// Order in the list implies configuration sequence.
	Initializers() []FileInitializer

	// IsConfigured checks if the provider is fully configured.
	// Returns true if all initializers report configured.
	// Typically implemented as:
	// AreInitializersConfigured(p.Initializers(), projectPath)
	IsConfigured(projectPath string) bool

	// GetFilePaths returns the file paths that this provider creates/updates.
	// Typically implemented as: GetInitializerPaths(p.Initializers())
	GetFilePaths() []string
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
