package providerkit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// File permissions
	defaultDirPerm  = 0755
	defaultFilePerm = 0644

	// Error messages
	errEmptyPath = "path cannot be empty"
)

// ExpandPath expands a path that may contain ~ for home directory
// or relative paths.
//
// It handles:
//   - Home directory expansion (~/)
//   - Relative paths (converts to absolute)
//   - Absolute paths (returns as-is)
//
// Returns an absolute path or an error if expansion fails.
func ExpandPath(path string) (string, error) {
	// Handle empty path
	if path == "" {
		return "", errors.New(errEmptyPath)
	}

	expandedPath := path

	// Handle home directory expansion
	if path == "~" || strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}

		if path == "~" {
			return homeDir, nil
		}

		// Replace ~ with home directory
		expandedPath = filepath.Join(homeDir, path[2:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return absPath, nil
}

// EnsureDir creates a directory and all parent directories if they don't exist.
// It is idempotent - no error is returned if the directory already exists.
// Directories are created with 0755 permissions (rwxr-xr-x).
//
// Returns an error if directory creation fails.
func EnsureDir(path string) error {
	if path == "" {
		return errors.New(errEmptyPath)
	}

	// Expand path to handle ~ and relative paths
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	// Create directory with parent directories, idempotent
	err = os.MkdirAll(expandedPath, defaultDirPerm)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory %s: %w",
			expandedPath,
			err,
		)
	}

	return nil
}

// FileExists checks if a file or directory exists at the given path.
// Returns false if the path doesn't exist or if there's an error checking.
func FileExists(path string) bool {
	if path == "" {
		return false
	}

	// Expand path to handle ~ and relative paths
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(expandedPath)

	return err == nil
}
