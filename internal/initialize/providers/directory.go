package providers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/templates"
	"github.com/spf13/afero"
)

// DirectoryInitializer creates directories in the project filesystem.
// Uses MkdirAll for silent success if directory already exists.
type DirectoryInitializer struct {
	paths []string // directories to create (relative to project root)
}

// NewDirectoryInitializer creates a new DirectoryInitializer for project
// filesystem.
func NewDirectoryInitializer(paths ...string) *DirectoryInitializer {
	return &DirectoryInitializer{
		paths: paths,
	}
}

// Init creates directories in the project filesystem.
// Returns silent success if directories already exist (MkdirAll behavior).
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (d *DirectoryInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	_ *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	var created []string

	for _, path := range d.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(projectFs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf("failed to check directory %s: %w", path, err)
		}

		// Create directory with MkdirAll (silent success if exists)
		if err := projectFs.MkdirAll(path, 0755); err != nil {
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

// IsSetup returns true if all directories exist in the project filesystem.
func (d *DirectoryInitializer) IsSetup(projectFs, _ afero.Fs, _ *Config) bool {
	for _, path := range d.paths {
		exists, err := afero.DirExists(projectFs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// dedupeKey returns a unique key for deduplication.
// Format: "DirectoryInitializer:<path1>:<path2>:..."
func (d *DirectoryInitializer) dedupeKey() string {
	// Normalize paths and sort for consistent keys
	normalized := make([]string, len(d.paths))
	for i, path := range d.paths {
		normalized[i] = filepath.Clean(path)
	}

	key := "DirectoryInitializer"
	for _, path := range normalized {
		key += ":" + path
	}

	return key
}

// HomeDirectoryInitializer creates directories in the home filesystem.
// Uses MkdirAll for silent success if directory already exists.
type HomeDirectoryInitializer struct {
	paths []string // directories to create (relative to home directory)
}

// NewHomeDirectoryInitializer creates a new HomeDirectoryInitializer for
// home filesystem.
func NewHomeDirectoryInitializer(
	paths ...string,
) *HomeDirectoryInitializer {
	return &HomeDirectoryInitializer{
		paths: paths,
	}
}

// Init creates directories in the home filesystem.
// Returns silent success if directories already exist (MkdirAll behavior).
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (h *HomeDirectoryInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	_ *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	var created []string

	for _, path := range h.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(homeFs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf("failed to check directory %s: %w", path, err)
		}

		// Create directory with MkdirAll (silent success if exists)
		if err := homeFs.MkdirAll(path, 0755); err != nil {
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

// IsSetup returns true if all directories exist in the home filesystem.
func (h *HomeDirectoryInitializer) IsSetup(_, homeFs afero.Fs, _ *Config) bool {
	for _, path := range h.paths {
		exists, err := afero.DirExists(homeFs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// dedupeKey returns a unique key for deduplication.
// Format: "HomeDirectoryInitializer:<path1>:<path2>:..."
func (h *HomeDirectoryInitializer) dedupeKey() string {
	// Normalize paths and sort for consistent keys
	normalized := make([]string, len(h.paths))
	for i, path := range h.paths {
		normalized[i] = filepath.Clean(path)
	}

	key := "HomeDirectoryInitializer"
	for _, path := range normalized {
		key += ":" + path
	}

	return key
}
