package discovery

import (
	"fmt"
	"os"
	"path/filepath"
)

// downwardContext holds the context for downward directory traversal.
type downwardContext struct {
	absStartPath string
	cwd          string
	depthMap     map[string]int
	maxDepth     int
	roots        *[]SpectrRoot
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
// Only directories that also have a .git at the same level are considered valid.
func addSpectrRootIfExists(path, cwd string, roots *[]SpectrRoot) {
	spectrDir := filepath.Join(path, spectrDirName)
	info, statErr := os.Stat(spectrDir)
	if statErr != nil || !info.IsDir() {
		return
	}

	// Only add as root if .git exists at same level
	// This prevents test fixtures and example directories from being discovered
	if !hasGitAtLevel(path) {
		return
	}

	// Found a valid spectr/ directory with .git at same level!
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

	gitDir := filepath.Join(path, gitDirName)
	info, err := os.Stat(gitDir)
	// If .git exists (as dir or file for worktrees), skip descending
	return err == nil && (info.IsDir() || !info.IsDir())
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
