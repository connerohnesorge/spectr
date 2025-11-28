package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/connerohnesorge/spectr/internal/config"
)

// GetActiveChangesWithConfig finds all active changes using the provided
// config, excluding archive directory.
func GetActiveChangesWithConfig(cfg *config.Config) ([]string, error) {
	changesDir := cfg.ChangesPath()

	// Check if changes directory exists
	_, err := os.Stat(changesDir)
	if os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read changes directory: %w", err)
	}

	var changes []string
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip archive directory
		if entry.Name() == "archive" {
			continue
		}

		// Check if proposal.md exists
		proposalPath := filepath.Join(changesDir, entry.Name(), "proposal.md")
		_, err = os.Stat(proposalPath)
		if err == nil {
			changes = append(changes, entry.Name())
		}
	}

	// Sort alphabetically for consistency
	sort.Strings(changes)

	return changes, nil
}

// GetActiveChanges finds all active changes in spectr/changes/,
// excluding archive directory. This function loads config from the project
// path. For direct config usage, use GetActiveChangesWithConfig.
func GetActiveChanges(projectPath string) ([]string, error) {
	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return GetActiveChangesWithConfig(cfg)
}

// GetActiveChangeIDs returns a list of active change IDs
// (directory names under spectr/changes/, excluding archive/)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetActiveChangeIDs(projectRoot string) ([]string, error) {
	return GetActiveChanges(projectRoot)
}

// GetActiveChangeIDsWithConfig returns a list of active change IDs
// using the provided config.
func GetActiveChangeIDsWithConfig(cfg *config.Config) ([]string, error) {
	return GetActiveChangesWithConfig(cfg)
}
