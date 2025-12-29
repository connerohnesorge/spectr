package providers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

// DirectoryInitializer creates directories in the project filesystem.
// Multiple paths can be created in a single initializer.
// Uses MkdirAll semantics - creates parent directories if needed, succeeds silently if directory exists. //nolint:lll
type DirectoryInitializer struct {
	paths []string // Relative paths from project root
}

// NewDirectoryInitializer creates an initializer for project directories.
// Paths are relative to the project filesystem root.
// Example: NewDirectoryInitializer(".claude/commands/spectr")
func NewDirectoryInitializer(paths ...string) Initializer {
	return &DirectoryInitializer{paths: paths}
}

// Init creates directories in the project filesystem.
// Returns created directories in InitResult (silent success if directory already exists).
//
//nolint:revive // Init signature is defined by Initializer interface
func (d *DirectoryInitializer) Init(
	ctx context.Context,
	projectFs, homeFs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	var created []string

	for _, path := range d.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(projectFs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf("failed to check directory %s: %w", path, err)
		}

		// Create directory with parents if needed
		if err := projectFs.MkdirAll(path, 0o755); err != nil {
			return InitResult{}, fmt.Errorf("failed to create directory %s: %w", path, err)
		}

		// Only report as created if it didn't exist before
		if !exists {
			created = append(created, path)
		}
	}

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: nil,
	}, nil
}

// IsSetup checks if all directories exist in the project filesystem.
func (d *DirectoryInitializer) IsSetup(
	projectFs, _ afero.Fs,
	_ *Config,
) bool { //nolint:lll // Function signature defined by Initializer interface
	for _, path := range d.paths {
		exists, err := afero.DirExists(projectFs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// dedupeKey returns a unique key for deduplication.
// Uses type name + normalized paths to prevent duplicate directory creation.
func (d *DirectoryInitializer) dedupeKey() string {
	if len(d.paths) == 0 {
		return "DirectoryInitializer:"
	}
	// Normalize paths and join with separator
	normalized := make([]string, 0, len(d.paths))
	for _, p := range d.paths {
		normalized = append(normalized, filepath.Clean(p))
	}
	// Use first path for key (most initializers have single path)
	return fmt.Sprintf("DirectoryInitializer:%s", normalized[0])
}

// HomeDirectoryInitializer creates directories in the home filesystem.
// Multiple paths can be created in a single initializer.
// Uses MkdirAll semantics - creates parent directories if needed, succeeds silently if directory exists. //nolint:lll
type HomeDirectoryInitializer struct {
	paths []string // Relative paths from home directory
}

// NewHomeDirectoryInitializer creates an initializer for home directories.
// Paths are relative to the home filesystem root (~/).
// Example: NewHomeDirectoryInitializer(".codex/prompts")
func NewHomeDirectoryInitializer(paths ...string) Initializer {
	return &HomeDirectoryInitializer{paths: paths}
}

// Init creates directories in the home filesystem.
// Returns created directories in InitResult (silent success if directory already exists).
//
//nolint:revive // Init signature is defined by Initializer interface
func (h *HomeDirectoryInitializer) Init(
	ctx context.Context,
	projectFs, homeFs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	var created []string

	for _, path := range h.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(homeFs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf("failed to check directory %s: %w", path, err)
		}

		// Create directory with parents if needed
		if err := homeFs.MkdirAll(path, 0o755); err != nil {
			return InitResult{}, fmt.Errorf("failed to create directory %s: %w", path, err)
		}

		// Only report as created if it didn't exist before
		if !exists {
			created = append(created, path)
		}
	}

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: nil,
	}, nil
}

// IsSetup checks if all directories exist in the home filesystem.
func (h *HomeDirectoryInitializer) IsSetup(
	_, homeFs afero.Fs,
	_ *Config,
) bool { //nolint:lll // Function signature defined by Initializer interface
	for _, path := range h.paths {
		exists, err := afero.DirExists(homeFs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// dedupeKey returns a unique key for deduplication.
// Uses type name + normalized paths to prevent duplicate directory creation.
// Separate type ensures home and project directories are not deduplicated against each other.
func (h *HomeDirectoryInitializer) dedupeKey() string {
	if len(h.paths) == 0 {
		return "HomeDirectoryInitializer:"
	}
	// Normalize paths and join with separator
	normalized := make([]string, 0, len(h.paths))
	for _, p := range h.paths {
		normalized = append(normalized, filepath.Clean(p))
	}
	// Use first path for key (most initializers have single path)
	return fmt.Sprintf("HomeDirectoryInitializer:%s", normalized[0])
}
