package providers

import (
	"context"

	"github.com/spf13/afero"
)

// TemplateManager is an interface for the template manager from
// internal/initialize.
// This avoids import cycles while allowing initializers to use
// template rendering.
// The concrete implementation is *initialize.TemplateManager.
type TemplateManager interface {
	// RenderAgents renders the agents template.
	// Will be deprecated once all code uses TemplateRef.
	RenderAgents(
		ctx TemplateContext,
	) (string, error)

	// RenderInstructionPointer renders the instruction pointer template.
	// Will be deprecated once all code uses TemplateRef.
	RenderInstructionPointer(
		ctx TemplateContext,
	) (string, error)

	// RenderSlashCommand renders a slash command template.
	// Will be deprecated once all code uses TemplateRef.
	RenderSlashCommand(
		commandType string,
		ctx TemplateContext,
	) (string, error)

	// InstructionPointer returns the instruction pointer template ref.
	// Type-safe accessor added in task 2.2-2.4.
	// These methods are defined here as interface placeholders to
	// avoid import cycles.
	// The actual return type (TemplateRef) is defined in
	// internal/initialize/templates package.
	// Callers should import that package to use Render() method.
	//
	// Example usage:
	//   import "github.com/.../spectr/internal/initialize/templates"
	//   ref := tm.InstructionPointer()
	//   content, err := ref.Render(ctx)
	//
	// Note: We cannot import templates package here due to import
	// cycle:
	//   providers -> templates -> providers (for TemplateContext)
	// So we use any as return type and document the actual type.
	InstructionPointer() any // returns templates.TemplateRef

	// Agents returns the agents template ref.
	Agents() any // returns templates.TemplateRef

	// Project returns the project template ref.
	Project() any // returns templates.TemplateRef

	// CIWorkflow returns the CI workflow template ref.
	CIWorkflow() any // returns templates.TemplateRef

	// SlashCommand returns the slash command template ref.
	// cmd: templates.SlashCommand, returns templates.TemplateRef
	SlashCommand(cmd any) any
}

// Initializer represents a unit of initialization work that creates or
// updates files.
// Initializers are composable, idempotent, and can be shared between
// providers.
//
// Key properties:
//   - Idempotent: Safe to run multiple times without adverse effects
//   - Deduplicatable: Multiple initializers with same Path() run once
//   - Type-ordered: DirectoryInitializer → ConfigFileInitializer →
//     SlashCommandsInitializer
//   - Filesystem-aware: Uses afero.Fs for testability
type Initializer interface {
	// Init creates or updates files. Returns result with file changes
	// and error if initialization fails.
	// Must be idempotent (safe to run multiple times).
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - fs: Filesystem (projectFs or globalFs based on IsGlobal())
	//   - cfg: Configuration containing directory paths
	//   - tm: Template manager for rendering templates
	//
	// Returns:
	//   - InitResult: Files created/updated by this initializer
	//   - error: If initialization fails (I/O, template rendering)
	Init(
		ctx context.Context,
		fs afero.Fs,
		cfg *Config,
		tm TemplateManager,
	) (InitResult, error)

	// IsSetup returns true if this initializer's artifacts already exist.
	// Used to determine if re-initialization is needed.
	//
	// Parameters:
	//   - fs: Filesystem to check (projectFs or globalFs based on IsGlobal())
	//   - cfg: Configuration containing directory paths
	//
	// Returns:
	//   - true if files/directories managed by this initializer exist
	//   - false if initialization is needed
	IsSetup(fs afero.Fs, cfg *Config) bool

	// Path returns the file/directory path this initializer manages.
	// Used for deduplication: same path = run once.
	//
	// Examples:
	//   - DirectoryInitializer: ".claude/commands/spectr"
	//   - ConfigFileInitializer: "CLAUDE.md"
	//   - SlashCommandsInitializer: ".claude/commands/spectr"
	//     (parent directory)
	//
	// Returns:
	//   - File or directory path relative to filesystem root
	Path() string

	// IsGlobal returns true if this initializer uses globalFs
	// instead of projectFs.
	//
	// Most initializers use projectFs (rooted at project directory):
	//   - CLAUDE.md
	//   - .claude/commands/spectr/
	//
	// Some use globalFs (rooted at user home directory):
	//   - ~/.config/aider/commands/
	//
	// Returns:
	//   - true: Use globalFs (e.g., ~/.config/)
	//   - false: Use projectFs (e.g., project-relative paths)
	IsGlobal() bool
}
