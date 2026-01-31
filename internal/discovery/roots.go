package discovery

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// spectrDirName is the standard name for spectr directories.
	spectrDirName = "spectr"
)

// SpectrRoot represents a discovered spectr/ directory with its location context.
type SpectrRoot struct {
	// Path is the absolute path to the directory containing spectr/
	// (e.g., /home/user/mono/project)
	Path string

	// RelativeTo is the path relative to the current working directory
	// (e.g., "../project" or ".")
	RelativeTo string

	// GitRoot is the absolute path to the parent .git directory
	// (e.g., /home/user/mono)
	GitRoot string
}

// SpectrDir returns the absolute path to the spectr/ directory.
func (r SpectrRoot) SpectrDir() string {
	return filepath.Join(r.Path, spectrDirName)
}

// ChangesDir returns the absolute path to the spectr/changes/ directory.
func (r SpectrRoot) ChangesDir() string {
	return filepath.Join(r.Path, "spectr", "changes")
}

// SpecsDir returns the absolute path to the spectr/specs/ directory.
func (r SpectrRoot) SpecsDir() string {
	return filepath.Join(r.Path, "spectr", "specs")
}

// DisplayName returns a human-readable name for the root.
// If RelativeTo is ".", returns the directory name; otherwise returns RelativeTo.
func (r SpectrRoot) DisplayName() string {
	if r.RelativeTo == "." {
		return filepath.Base(r.Path)
	}

	return r.RelativeTo
}

// FindSpectrRoots discovers all spectr/ directories by walking up the directory
// tree from the given current working directory, stopping at git repository
// boundaries.
//
// If SPECTR_ROOT environment variable is set, it returns only that root
// (and validates it exists).
//
// The function returns roots in order from closest to cwd to furthest.
func FindSpectrRoots(cwd string) ([]SpectrRoot, error) {
	// Check for SPECTR_ROOT env var override
	if envRoot := os.Getenv("SPECTR_ROOT"); envRoot != "" {
		return findSpectrRootFromEnv(envRoot, cwd)
	}

	return findSpectrRootsFromCwd(cwd)
}

// findSpectrRootFromEnv handles the SPECTR_ROOT environment variable case.
func findSpectrRootFromEnv(envRoot, cwd string) ([]SpectrRoot, error) {
	// Make path absolute if relative
	absPath := envRoot
	if !filepath.IsAbs(envRoot) {
		absPath = filepath.Join(cwd, envRoot)
	}
	absPath = filepath.Clean(absPath)

	// Validate spectr/ directory exists
	spectrDir := filepath.Join(absPath, spectrDirName)
	info, err := os.Stat(spectrDir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf(
			"SPECTR_ROOT path does not contain spectr/ directory: %s",
			absPath,
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"failed to check SPECTR_ROOT path: %w",
			err,
		)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf(
			"SPECTR_ROOT path does not contain spectr/ directory: %s",
			absPath,
		)
	}

	// Calculate relative path from cwd
	relPath, err := filepath.Rel(cwd, absPath)
	if err != nil {
		relPath = absPath // Fallback to absolute if rel fails
	}

	// Find git root for this path
	gitRoot := findGitRoot(absPath)

	return []SpectrRoot{
		{
			Path:       absPath,
			RelativeTo: relPath,
			GitRoot:    gitRoot,
		},
	}, nil
}

// findSpectrRootsFromCwd walks up from cwd to find all spectr/ directories,
// stopping at git boundaries.
func findSpectrRootsFromCwd(cwd string) ([]SpectrRoot, error) {
	var roots []SpectrRoot

	// Make cwd absolute
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Find the git root first to establish the boundary
	gitRoot := findGitRoot(absCwd)

	// Walk up from cwd to git root (or filesystem root if no git)
	current := absCwd
	for {
		// Check if spectr/ directory exists at this level
		spectrDir := filepath.Join(current, spectrDirName)
		info, err := os.Stat(spectrDir)
		if err == nil && info.IsDir() {
			// Calculate relative path from original cwd
			relPath, relErr := filepath.Rel(absCwd, current)
			if relErr != nil {
				relPath = current // Fallback to absolute
			}

			roots = append(roots, SpectrRoot{
				Path:       current,
				RelativeTo: relPath,
				GitRoot:    gitRoot,
			})
		}

		// Stop if we've reached the git root
		if gitRoot != "" && current == gitRoot {
			break
		}

		// Stop if we've reached the filesystem root
		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		current = parent
	}

	return roots, nil
}

// findGitRoot walks up from the given path to find the nearest .git directory.
// Returns empty string if no git root is found.
func findGitRoot(startPath string) string {
	current := startPath
	for {
		gitDir := filepath.Join(current, ".git")
		info, err := os.Stat(gitDir)
		if err == nil && info.IsDir() {
			return current
		}

		// Also check for git worktree files (where .git is a file, not dir)
		if err == nil && !info.IsDir() {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root without finding .git
			return ""
		}

		current = parent
	}
}
