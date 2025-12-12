package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GetActiveChanges finds all active changes in spectr/changes/,
// excluding archive directory
func GetActiveChanges(
	projectPath string,
) ([]string, error) {
	changesDir := filepath.Join(
		projectPath,
		"spectr",
		"changes",
	)

	// Check if changes directory exists
	_, err := os.Stat(changesDir)
	if os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read changes directory: %w",
			err,
		)
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
		proposalPath := filepath.Join(
			changesDir,
			entry.Name(),
			"proposal.md",
		)
		_, err = os.Stat(proposalPath)
		if err == nil {
			changes = append(
				changes,
				entry.Name(),
			)
		}
	}

	// Sort alphabetically for consistency
	sort.Strings(changes)

	return changes, nil
}

// GetActiveChangeIDs returns a list of active change IDs
// (directory names under spectr/changes/, excluding archive/)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetActiveChangeIDs(
	projectRoot string,
) ([]string, error) {
	return GetActiveChanges(projectRoot)
}

// ResolveResult contains the resolved change ID and whether it was a partial
// match.
type ResolveResult struct {
	ChangeID     string
	PartialMatch bool
}

// ResolveChangeID resolves a partial change ID to a full change ID.
//
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
func ResolveChangeID(
	partialID, projectRoot string,
) (ResolveResult, error) {
	changes, err := GetActiveChangeIDs(
		projectRoot,
	)
	if err != nil {
		return ResolveResult{},
			fmt.Errorf(
				"get active changes: %w",
				err,
			)
	}

	if len(changes) == 0 {
		return ResolveResult{},
			fmt.Errorf(
				"no change found matching '%s'",
				partialID,
			)
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
		if strings.HasPrefix(
			strings.ToLower(change),
			partialLower,
		) {
			prefixMatches = append(
				prefixMatches,
				change,
			)
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
		if strings.Contains(
			strings.ToLower(change),
			partialLower,
		) {
			substringMatches = append(
				substringMatches,
				change,
			)
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
