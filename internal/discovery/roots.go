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

	// maxDiscoveryDepth limits how deep downward discovery will traverse.
	maxDiscoveryDepth = 10
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

// appendDownwardRoots performs downward discovery and appends results to roots.
// Downward discovery happens when:
// a) We're NOT inside a git repository (gitRoot is empty), OR
// b) We ARE at the git root itself (to find nested subprojects in monorepos)
// This enables monorepo support where the root contains subprojects with
// their own .git and spectr/ directories.
func appendDownwardRoots(existingRoots []SpectrRoot, absCwd, gitRoot string) []SpectrRoot {
	if gitRoot != "" && absCwd != gitRoot {
		return existingRoots
	}

	downwardRoots, err := findSpectrRootsDownward(absCwd, absCwd, maxDiscoveryDepth)
	// Ignore downward discovery errors - upward discovery already succeeded
	if err == nil {
		return append(existingRoots, downwardRoots...)
	}

	return existingRoots
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

// shouldSkipDirectory returns true if the directory should be skipped during downward discovery.
func shouldSkipDirectory(dirName string) bool {
	skipDirs := []string{".git", "node_modules", "vendor", "target", "dist", "build"}
	for _, skip := range skipDirs {
		if dirName == skip {
			return true
		}
	}

	return false
}

// calculateDepth computes the depth of a directory relative to the start path.
func calculateDepth(path, absStartPath string, depthMap map[string]int) int {
	parent := filepath.Dir(path)
	if depth, ok := depthMap[parent]; ok {
		return depth + 1
	}

	// Fallback: calculate depth from path segments
	relPath, relErr := filepath.Rel(absStartPath, path)
	if relErr == nil {
		return len(filepath.SplitList(relPath))
	}

	return 0
}

// addSpectrRootIfExists checks if a directory contains a spectr/ subdirectory
// and adds it to the roots slice if it does.
func addSpectrRootIfExists(path, cwd string, roots *[]SpectrRoot) {
	spectrDir := filepath.Join(path, spectrDirName)
	info, statErr := os.Stat(spectrDir)
	if statErr != nil || !info.IsDir() {
		return
	}

	// Found a spectr/ directory!
	// Calculate relative path from original cwd
	relPath, relErr := filepath.Rel(cwd, path)
	if relErr != nil {
		relPath = path // Fallback to absolute
	}

	// Find git root for this spectr root
	gitRoot := findGitRoot(path)

	*roots = append(*roots, SpectrRoot{
		Path:       path,
		RelativeTo: relPath,
		GitRoot:    gitRoot,
	})
}

// shouldSkipGitBoundary checks if a directory contains a .git subdirectory
// and should not be descended into (unless it's the start path).
func shouldSkipGitBoundary(path, absStartPath string) bool {
	if path == absStartPath {
		return false // Don't skip the start path itself
	}

	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	// If .git exists (as dir or file for worktrees), skip descending
	return err == nil && (info.IsDir() || !info.IsDir())
}

// downwardContext holds the context for downward directory traversal.
type downwardContext struct {
	absStartPath string
	cwd          string
	depthMap     map[string]int
	maxDepth     int
	roots        *[]SpectrRoot
}

// processDownwardDirectory handles a single directory during downward discovery.
// Returns filepath.SkipDir if the directory should not be descended into.
func processDownwardDirectory(path string, d os.DirEntry, ctx *downwardContext) error {
	// Only process directories
	if !d.IsDir() {
		return nil
	}

	// Calculate and store current depth
	currentDepth := calculateDepth(path, ctx.absStartPath, ctx.depthMap)
	ctx.depthMap[path] = currentDepth

	// Stop descending if we've hit max depth
	if currentDepth > ctx.maxDepth {
		return filepath.SkipDir
	}

	// Skip descending into common non-project directories
	if shouldSkipDirectory(d.Name()) {
		return filepath.SkipDir
	}

	// Check if this directory contains a spectr/ subdirectory and add it if so
	addSpectrRootIfExists(path, ctx.cwd, ctx.roots)

	// Check if we should skip descending into this directory (git boundary)
	if shouldSkipGitBoundary(path, ctx.absStartPath) {
		return filepath.SkipDir
	}

	return nil
}

// findSpectrRootsDownward searches for spectr/ directories in subdirectories,
// descending from startPath up to maxDepth levels deep. It discovers nested
// repositories (directories with .git) and their spectr/ directories.
//
// This complements upward discovery to support mono-repo structures where
// multiple nested projects each have their own .git and spectr/ directories.
//
// The function:
// - Uses filepath.WalkDir for efficient traversal
// - Tracks depth with configurable limit (prevents excessive traversal)
// - Finds all spectr/ directories in subdirectories
// - Creates SpectrRoot entries with Path, RelativeTo (from cwd), and GitRoot
// - Skips descending into .git/, node_modules/, vendor/, target/, dist/, build/
// - Includes directories that CONTAIN .git (nested repos are discovered)
// - Handles permission errors gracefully (continues search)
// - Continues searching after finding spectr/ (doesn't stop at first match)
func findSpectrRootsDownward(startPath, cwd string, maxDepth int) ([]SpectrRoot, error) {
	var roots []SpectrRoot
	absStartPath, err := filepath.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create context for traversal
	ctx := &downwardContext{
		absStartPath: absStartPath,
		cwd:          cwd,
		depthMap:     map[string]int{absStartPath: 0},
		maxDepth:     maxDepth,
		roots:        &roots,
	}

	err = filepath.WalkDir(absStartPath, func(path string, d os.DirEntry, err error) error {
		// Handle permission errors gracefully - continue walking
		if err != nil {
			// Skip directories we can't read
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}

			return nil // Continue for non-directory errors
		}

		return processDownwardDirectory(path, d, ctx)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory tree: %w", err)
	}

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
