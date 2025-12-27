package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/templates"
	"github.com/spf13/afero"
)

// Initializer represents a unit of initialization work (e.g., creating a
// directory, updating a config file, or creating slash commands).
//
// Initializers are composable, type-safe, and can be shared across providers.
type Initializer interface {
	// Init creates or updates files. Returns result with file changes and
	// error if initialization fails. Must be idempotent (safe to run multiple
	// times).
	// Receives both filesystems - initializer decides which to use based on
	// its type.
	Init(
		ctx context.Context,
		projectFs, homeFs afero.Fs,
		cfg *Config,
		tm *templates.TemplateManager,
	) (InitResult, error)

	// IsSetup returns true if this initializer's artifacts already exist.
	// Receives both filesystems - initializer checks the appropriate one.
	// PURPOSE: Used by the setup wizard to show which providers are already
	// configured.
	// NOT used to skip initializers during execution - Init() always runs
	// (idempotent).
	IsSetup(projectFs, homeFs afero.Fs, cfg *Config) bool
}
