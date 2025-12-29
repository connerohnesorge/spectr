package providers

import (
	"context"

	"github.com/spf13/afero"
)

// Initializer represents a unit of initialization work that creates or updates files.
// Each initializer is responsible for a specific aspect of provider configuration
// (e.g., creating directories, updating config files, creating slash commands).
//
// Initializers must be idempotent - they can be run multiple times safely.
type Initializer interface {
	// Init creates or updates files for this initializer.
	// Returns InitResult with created/updated file paths and error if initialization fails.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - projectFs: Filesystem rooted at project directory (for .claude/, CLAUDE.md, etc.)
	//   - homeFs: Filesystem rooted at user home directory (for ~/.codex/, etc.)
	//   - cfg: Configuration with SpectrDir and derived paths
	//   - tm: TemplateManager for rendering templates
	//
	// The initializer decides which filesystem to use based on its type
	// (e.g., DirectoryInitializer uses projectFs, HomeDirectoryInitializer uses homeFs).
	Init(
		ctx context.Context,
		projectFs, homeFs afero.Fs,
		cfg *Config,
		tm TemplateManager,
	) (InitResult, error)

	// IsSetup returns true if this initializer's artifacts already exist.
	//
	// Parameters:
	//   - projectFs: Filesystem rooted at project directory
	//   - homeFs: Filesystem rooted at user home directory
	//   - cfg: Configuration with SpectrDir and derived paths
	//
	// The initializer checks the appropriate filesystem based on its type.
	//
	// PURPOSE: Used by the setup wizard to show which providers are already configured.
	// NOT used to skip initializers during execution - Init() always runs (idempotent).
	IsSetup(projectFs, homeFs afero.Fs, cfg *Config) bool
}
