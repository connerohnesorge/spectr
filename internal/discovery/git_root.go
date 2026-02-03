package discovery

import (
	"os"
	"path/filepath"
	"sync"
)

// gitRootCache caches git root lookups to avoid repeated filesystem traversals.
var (
	gitRootCache   = make(map[string]string)
	gitRootCacheMu sync.RWMutex
)

// hasGitAtLevel checks if a .git directory (or file for worktrees) exists at the given path.
// This is used to validate that a spectr/ directory belongs to a real git repository.
func hasGitAtLevel(path string) bool {
	gitDir := filepath.Join(path, gitDirName)
	_, err := os.Stat(gitDir)

	return err == nil
}

// findGitRoot walks up from the given path to find the nearest .git directory.
// Returns empty string if no git root is found.
// Results are cached for performance.
func findGitRoot(startPath string) string {
	// Check cache first
	gitRootCacheMu.RLock()
	if cached, ok := gitRootCache[startPath]; ok {
		gitRootCacheMu.RUnlock()

		return cached
	}
	gitRootCacheMu.RUnlock()

	result := findGitRootUncached(startPath)

	// Cache the result
	gitRootCacheMu.Lock()
	gitRootCache[startPath] = result
	gitRootCacheMu.Unlock()

	return result
}

// findGitRootUncached is the uncached implementation of findGitRoot.
func findGitRootUncached(startPath string) string {
	current := startPath
	for {
		gitDir := filepath.Join(current, gitDirName)
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
