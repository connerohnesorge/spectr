// Package init provides filesystem operations and configuration utilities
// for initializing Spectr projects.
package init

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// File permissions
	dirPerms  = 0755
	filePerms = 0644

	// Error messages
	errEmptyPath = "path cannot be empty"
)

// ExpandPath expands a path that may contain ~ for home directory
// or relative paths.
// It handles:
// - Home directory expansion (~/)
// - Relative paths (converts to absolute)
// - Absolute paths (returns as-is)
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
// Directories are created with dirPerms permissions (rwxr-xr-x).
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
	err = os.MkdirAll(expandedPath, dirPerms)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory %s: %w",
			expandedPath,
			err,
		)
	}

	return nil
}

// WriteFile writes content to a file with filePerms permissions
// (rw-r--r--). It creates parent directories if needed.
// Returns an error if the file already exists to prevent accidental
// overwrites.
//
// Use BackupFile before calling WriteFile if you want to preserve
// existing files.
func WriteFile(path string, content []byte) error {
	if path == "" {
		return errors.New(errEmptyPath)
	}

	// Expand path to handle ~ and relative paths
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	// Check if file already exists (conflict detection)
	if FileExists(expandedPath) {
		return fmt.Errorf("file already exists: %s", expandedPath)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(expandedPath)
	err = EnsureDir(dir)
	if err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write file
	err = os.WriteFile(expandedPath, content, filePerms)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", expandedPath, err)
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

// IsSpectrInitialized checks if Spectr is already initialized in the project.
// It looks for the "spectr/project.md" file in the project directory.
//
// Returns true if the project is already initialized, false otherwise.
func IsSpectrInitialized(projectPath string) bool {
	if projectPath == "" {
		return false
	}

	// Expand path to handle ~ and relative paths
	expandedPath, err := ExpandPath(projectPath)
	if err != nil {
		return false
	}

	// Check for spectr/project.md
	projectFile := filepath.Join(expandedPath, "spectr", "project.md")

	return FileExists(projectFile)
}

// BackupFile creates a backup of an existing file with a timestamp
// suffix. The backup file name format is:
// original_name.backup.YYYYMMDD_HHMMSS.NNNNNNNNN
// If the file doesn't exist, no backup is created and no error is
// returned.
//
// Returns an error if the backup operation fails.
func BackupFile(path string) error {
	if path == "" {
		return errors.New(errEmptyPath)
	}

	// Expand path to handle ~ and relative paths
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	// Only backup if file exists
	if !FileExists(expandedPath) {
		return nil
	}

	// Read existing file content
	content, err := os.ReadFile(expandedPath)
	if err != nil {
		return fmt.Errorf("failed to read file for backup: %w", err)
	}

	// Generate backup filename with timestamp
	// (including nanoseconds for uniqueness)
	timestamp := time.Now().Format("20060102_150405.000000000")
	backupPath := fmt.Sprintf("%s.backup.%s", expandedPath, timestamp)

	// Write backup file
	err = os.WriteFile(backupPath, content, filePerms)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	return nil
}
