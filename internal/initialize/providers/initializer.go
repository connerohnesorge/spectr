package providers

import (
	"context"

	"github.com/spf13/afero"
)

// TemplateManager provides template rendering capabilities.
// This interface is satisfied by the TemplateManager in the parent
// initialize package.
type TemplateManager interface {
	// RenderAgents renders the AGENTS.md template content.
	RenderAgents(ctx TemplateContext) (string, error)

	// RenderInstructionPointer renders a short pointer template that directs
	// AI assistants to read spectr/AGENTS.md for full instructions.
	RenderInstructionPointer(ctx TemplateContext) (string, error)

	// RenderSlashCommand renders a slash command template (proposal or apply).
	RenderSlashCommand(command string, ctx TemplateContext) (string, error)
}

// Initializer represents a single initialization step that creates
// or updates files.
// Initializers are composable units that can be shared across providers.
//
// Example initializers:
//   - DirectoryInitializer: Creates directories
//   - ConfigFileInitializer: Creates/updates instruction files with markers
//   - SlashCommandsInitializer: Creates slash command files
//
// Initializers must be idempotent (safe to run multiple times).
type Initializer interface {
	// Init creates or updates files.
	// Returns result with file changes and error if initialization fails.
	// Must be idempotent (safe to run multiple times).
	//
	// Parameters:
	//   - ctx: Context for cancellation and deadlines
	//   - fs: Filesystem abstraction
	//        (project-relative or global, based on IsGlobal())
	//   - cfg: Configuration with SpectrDir and derived paths
	//   - tm: TemplateManager for rendering templates
	//
	// Returns:
	//   - InitResult: Contains lists of created and updated files
	//   - error: Non-nil if initialization fails
	Init(
		ctx context.Context,
		fs afero.Fs,
		cfg *Config,
		tm TemplateManager,
	) (InitResult, error)

	// IsSetup returns true if this initializer's artifacts already exist.
	// Used to determine if the provider is already configured.
	//
	// Parameters:
	//   - fs: Filesystem abstraction (same as Init)
	//   - cfg: Configuration with SpectrDir and derived paths
	IsSetup(fs afero.Fs, cfg *Config) bool

	// Path returns the file/directory path this initializer manages.
	// Used for deduplication: same path = run once.
	//
	// Example paths:
	//   - ".claude/commands/spectr"         (directory)
	//   - "CLAUDE.md"                        (config file)
	//   - ".claude/commands/spectr/proposal.md" (slash command)
	Path() string

	// IsGlobal returns true if this initializer uses globalFs
	// instead of projectFs.
	// Most initializers use projectFs (false).
	// Some tools like Aider use ~/.config/ (true).
	IsGlobal() bool
}
