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
// config.
// It searches in the changes directory specified by the config,
// excluding the archive directory.
func GetActiveChangesWithConfig(cfg *config.Config) ([]string, error) {
	changesDir := cfg.ChangesDir()

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
// excluding archive directory.
// Deprecated: Use GetActiveChangesWithConfig for projects with custom root
// directories.
func GetActiveChanges(projectPath string) ([]string, error) {
	cfg := &config.Config{
		RootDir:     config.DefaultRootDir,
		ProjectRoot: projectPath,
	}

	return GetActiveChangesWithConfig(cfg)
}

// GetActiveChangeIDsWithConfig returns a list of active change IDs using
// the provided config (directory names under the configured changes directory,
// excluding archive/).
// Returns empty slice (not error) if the directory doesn't exist.
// Results are sorted alphabetically for consistency.
func GetActiveChangeIDsWithConfig(cfg *config.Config) ([]string, error) {
	return GetActiveChangesWithConfig(cfg)
}

// GetActiveChangeIDs returns a list of active change IDs
// (directory names under spectr/changes/, excluding archive/).
// Returns empty slice (not error) if the directory doesn't exist.
// Results are sorted alphabetically for consistency.
// Deprecated: Use GetActiveChangeIDsWithConfig for projects with custom root
// directories.
func GetActiveChangeIDs(projectRoot string) ([]string, error) {
	return GetActiveChanges(projectRoot)
}

// ResolveResult contains the resolved change ID and whether it was a partial
// match.
type ResolveResult struct {
	ChangeID     string
	PartialMatch bool
}

// ResolveChangeIDWithConfig resolves a partial change ID to a full change ID
// using the provided config.
// The resolution algorithm:
//  1. Exact match: returns immediately if partialID exactly matches a change ID
//  2. Prefix match: finds all change IDs that start with partialID
//     (case-insensitive)
//  3. Substring match: if no prefix matches, finds all change IDs containing
//     partialID (case-insensitive)
//
// Returns an error if:
// - No changes match the partial ID
// - Multiple changes match the partial ID (ambiguous)
func ResolveChangeIDWithConfig(
	partialID string,
	cfg *config.Config,
) (ResolveResult, error) {
	changes, err := GetActiveChangeIDsWithConfig(cfg)
	if err != nil {
		return ResolveResult{},
			fmt.Errorf("get active changes: %w", err)
	}

	if len(changes) == 0 {
		return ResolveResult{},
			fmt.Errorf("no change found matching '%s'", partialID)
	}

	partialLower := strings.ToLower(partialID)

	// Check for exact match first
	for _, change := range changes {
		if change == partialID {
			return ResolveResult{
				ChangeID:     change,
				PartialMatch: false,
			}, nil
		}
	}

	// Try prefix matching (case-insensitive)
	var prefixMatches []string
	for _, change := range changes {
		if strings.HasPrefix(strings.ToLower(change), partialLower) {
			prefixMatches = append(prefixMatches, change)
		}
	}

	if len(prefixMatches) == 1 {
		return ResolveResult{
			ChangeID:     prefixMatches[0],
			PartialMatch: true,
		}, nil
	}

	if len(prefixMatches) > 1 {
		return ResolveResult{}, fmt.Errorf(
			"ambiguous ID '%s' matches multiple changes: %s",
			partialID,
			strings.Join(prefixMatches, ", "),
		)
	}

	// Try substring matching (case-insensitive) as fallback
	var substringMatches []string
	for _, change := range changes {
		if strings.Contains(strings.ToLower(change), partialLower) {
			substringMatches = append(substringMatches, change)
		}
	}

	if len(substringMatches) == 1 {
		return ResolveResult{
			ChangeID:     substringMatches[0],
			PartialMatch: true,
		}, nil
	}

	if len(substringMatches) > 1 {
		return ResolveResult{}, fmt.Errorf(
			"ambiguous ID '%s' matches multiple changes: %s",
			partialID,
			strings.Join(substringMatches, ", "),
		)
	}

	return ResolveResult{},
		fmt.Errorf(
			"no change found matching '%s'",
			partialID,
		)
}

// ResolveChangeID resolves a partial change ID to a full change ID.
// The resolution algorithm:
//  1. Exact match: returns immediately if partialID exactly matches a change ID
//  2. Prefix match: finds all change IDs that start with partialID
//     (case-insensitive)
//  3. Substring match: if no prefix matches, finds all change IDs containing
//     partialID (case-insensitive)
//
// Returns an error if:
// - No changes match the partial ID
// - Multiple changes match the partial ID (ambiguous)
// Deprecated: Use ResolveChangeIDWithConfig for projects with custom root
// directories.
func ResolveChangeID(partialID, projectRoot string) (ResolveResult, error) {
	cfg := &config.Config{
		RootDir:     config.DefaultRootDir,
		ProjectRoot: projectRoot,
	}

	return ResolveChangeIDWithConfig(partialID, cfg)
}
