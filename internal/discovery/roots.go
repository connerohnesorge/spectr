package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	// spectrDirName is the standard name for spectr directories.
	spectrDirName = "spectr"

	// gitDirName is the standard name for git directories.
	gitDirName = ".git"

	// maxDiscoveryDepth limits how deep downward discovery will traverse.
	// Set to 5 to balance discovery coverage with performance.
	// Nested repos beyond this depth are unlikely in practice.
	maxDiscoveryDepth = 5
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
// stopping at git boundaries, and also searches downward from cwd (or git root)
// to find nested spectr/ directories in subdirectories.
func findSpectrRootsFromCwd(cwd string) ([]SpectrRoot, error) {
	var roots []SpectrRoot

	// Make cwd absolute
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Find the git root first to establish the boundary
	gitRoot := findGitRoot(absCwd)

	// 1. Upward discovery: Walk up from cwd to git root (or filesystem root if no git)
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

	// 2. Downward discovery: Search for nested spectr/ directories from cwd
	roots = appendDownwardRoots(roots, absCwd, gitRoot)

	// 3. Deduplicate roots (upward and downward may find same directories)
	roots = deduplicateRoots(roots)

	// 4. Sort by distance from cwd (closest first)
	roots = sortRootsByDistance(roots, absCwd)

	return roots, nil
}

// deduplicateRoots removes duplicate SpectrRoot entries based on their Path field.
// Preserves the order of first occurrence.
func deduplicateRoots(roots []SpectrRoot) []SpectrRoot {
	if len(roots) == 0 {
		return roots
	}

	seen := make(map[string]bool)
	result := make([]SpectrRoot, 0, len(roots))

	for _, root := range roots {
		if !seen[root.Path] {
			seen[root.Path] = true
			result = append(result, root)
		}
	}

	return result
}

// sortRootsByDistance sorts SpectrRoot entries by their relative path length from cwd.
// Roots closer to cwd (shorter relative paths) appear first.
// This ensures that when discovering both upward and downward, the closest root is prioritized.
func sortRootsByDistance(roots []SpectrRoot, cwd string) []SpectrRoot {
	if len(roots) <= 1 {
		return roots
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]SpectrRoot, len(roots))
	copy(sorted, roots)

	// Sort by the length of the relative path
	// filepath.Rel returns the shortest path, so shorter = closer
	sort.Slice(sorted, func(i, j int) bool {
		relI, errI := filepath.Rel(cwd, sorted[i].Path)
		relJ, errJ := filepath.Rel(cwd, sorted[j].Path)

		// If there's an error calculating relative path, fall back to comparing paths
		if errI != nil || errJ != nil {
			return sorted[i].Path < sorted[j].Path
		}

		// Compare by number of path separators (shorter = closer)
		// "." has 0 separators (closest)
		// ".." has 1 separator
		// "../.." has 2 separators, etc.
		depthI := len(filepath.SplitList(relI))
		depthJ := len(filepath.SplitList(relJ))

		if depthI == depthJ {
			// If same depth, sort alphabetically for consistency
			return relI < relJ
		}

		return depthI < depthJ
	})

	return sorted
}
