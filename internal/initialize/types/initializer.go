package types

import (
	"context"

	"github.com/spf13/afero"
)

// Initializer defines a unit of initialization logic.
// It is designed to be composable and testable.
type Initializer interface {
	// Init performs the initialization logic.
	// It receives the project filesystem, global filesystem, configuration, and template manager.
	// It returns an error if initialization fails.
	Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *Config, tm TemplateRenderer) error

	// IsSetup checks if the initialization has already been performed.
	// It is used to determine if the initializer should run.
	IsSetup(projectFs, globalFs afero.Fs, cfg *Config) (bool, error)

	// Path returns the unique path identifier for this initializer.
	// It is used for deduplication. Use "" if the initializer is not tied to a specific file.
	Path() string
}
