package providers

import (
	"context"

	"github.com/spf13/afero"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// Initializer represents a single initialization step.
//
// Each initializer is responsible for one aspect of provider configuration
// (creating a directory, writing a config file, or creating slash commands).
//
// Initializers are designed to be:
// - Composable: Multiple initializers can be combined to configure a provider
// - Reusable: The same initializer can be shared across providers
// - Testable: Small, focused units that are easy to test in isolation
// - Idempotent: Safe to run multiple times without unwanted side effects
type Initializer interface {
	// Init creates or updates files managed by this initializer.
	// Returns InitResult containing the files that were created or updated.
	// Returns an error if initialization fails.
	//
	// Parameters:
	//   - ctx: Context for cancellation and deadlines
	//   - fs: Filesystem abstraction (either projectFs or globalFs)
	//   - cfg: Provider configuration (SpectrDir and derived paths)
	//   - tm: Template manager for rendering templates
	//
	// Must be idempotent: running Init() multiple times should be safe.
	//
	// Note: tm is typed as 'any' to avoid import cycles. The concrete type is
	// *initialize.TemplateManager from internal/initialize/templates.go.
	// The executor will pass the correct type.
	Init(
		ctx context.Context,
		fs afero.Fs,
		cfg *initializers.Config,
		tm any,
	) (initializers.InitResult, error)

	// IsSetup returns true if this initializer's artifacts already exist.
	// Used to determine if initialization is needed.
	//
	// Parameters:
	//   - fs: Filesystem abstraction (either projectFs or globalFs)
	//   - cfg: Provider configuration (SpectrDir and derived paths)
	IsSetup(
		fs afero.Fs,
		cfg *initializers.Config,
	) bool

	// Path returns the primary file or directory path this initializer manages.
	// Used for deduplication: if multiple initializers return the same path,
	// only one will be executed.
	//
	// Example paths:
	//   - ".claude/commands/spectr" (directory)
	//   - "CLAUDE.md" (config file)
	//   - ".claude/commands/spectr/proposal.md" (slash command)
	Path() string

	// IsGlobal returns true if this initializer operates globally
	// (e.g., ~/.config/tool/) instead of on the project filesystem.
	//
	// Most initializers return false (operate on project files).
	// Return true for tools that use global configuration directories.
	IsGlobal() bool
}
