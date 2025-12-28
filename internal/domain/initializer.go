package domain

import (
	"context"

	"github.com/spf13/afero"
)

// Initializer represents a single initialization step that creates or updates files.
// Implementations are composable and can be shared across providers.
type Initializer interface {
	// Init creates or updates files. Returns result with file changes and error if initialization fails.
	// Must be idempotent (safe to run multiple times).
	// Receives both filesystems - initializer decides which to use based on its configuration.
	Init(
		ctx context.Context,
		projectFs, homeFs afero.Fs,
		cfg *Config,
		tm any,
	) (ExecutionResult, error)

	// IsSetup returns true if this initializer's artifacts already exist.
	// Receives both filesystems - initializer checks the appropriate one.
	// PURPOSE: Used by the setup wizard to show which providers are already configured.
	// NOT used to skip initializers during execution - Init() always runs (idempotent).
	IsSetup(projectFs, homeFs afero.Fs, cfg *Config) bool
}
