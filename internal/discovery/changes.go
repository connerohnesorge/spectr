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

// ChangeStatus represents the status of a change proposal.
type ChangeStatus int

const (
	// ChangeStatusUnknown indicates the change ID was not found anywhere.
	ChangeStatusUnknown ChangeStatus = iota
	// ChangeStatusActive indicates the change exists in spectr/changes/ (not archived).
	ChangeStatusActive
	// ChangeStatusArchived indicates the change exists in spectr/changes/archive/.
	ChangeStatusArchived
)

// String returns a human-readable string for the change status.
func (s ChangeStatus) String() string {
	switch s {
	case ChangeStatusActive:
		return "active"
	case ChangeStatusArchived:
		return "archived"
	case ChangeStatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// IsChangeArchived checks if a change with the given ID exists in the archive.
// It searches spectr/changes/archive/ for directories matching the pattern
// YYYY-MM-DD-<changeID> or just <changeID>.
func IsChangeArchived(
	changeID, projectRoot string,
) (bool, error) {
	archiveDir := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		"archive",
	)

	// Check if archive directory exists
	_, err := os.Stat(archiveDir)
	if os.IsNotExist(err) {
		return false, nil
	}

	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return false, fmt.Errorf(
			"failed to read archive directory: %w",
			err,
		)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		extractedID := ExtractChangeIDFromArchivePath(entry.Name())
		if extractedID == changeID {
			return true, nil
		}
	}

	return false, nil
}

// GetChangeStatus returns the status of a change proposal.
// It checks if the change exists as active or archived.
func GetChangeStatus(
	changeID, projectRoot string,
) (ChangeStatus, error) {
	// Check if archived first
	archived, err := IsChangeArchived(changeID, projectRoot)
	if err != nil {
		return ChangeStatusUnknown, err
	}
	if archived {
		return ChangeStatusArchived, nil
	}

	// Check if active
	activeChanges, err := GetActiveChangeIDs(projectRoot)
	if err != nil {
		return ChangeStatusUnknown, err
	}

	for _, id := range activeChanges {
		if id == changeID {
			return ChangeStatusActive, nil
		}
	}

	return ChangeStatusUnknown, nil
}

// GetArchivedChangeIDs returns a list of all archived change IDs.
// It extracts the change ID portion from archive directory names,
// which may have date prefixes (e.g., "2024-01-15-feat-auth" -> "feat-auth").
func GetArchivedChangeIDs(
	projectRoot string,
) ([]string, error) {
	archiveDir := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		"archive",
	)

	// Check if archive directory exists
	_, err := os.Stat(archiveDir)
	if os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read archive directory: %w",
			err,
		)
	}

	var changeIDs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Extract change ID from directory name
		changeID := ExtractChangeIDFromArchivePath(entry.Name())
		if changeID != "" {
			changeIDs = append(changeIDs, changeID)
		}
	}

	// Sort for consistency
	sort.Strings(changeIDs)

	return changeIDs, nil
}

// ExtractChangeIDFromArchivePath extracts the change ID from an archive
// directory name. Archive directories may have a date prefix in the format
// YYYY-MM-DD-<changeID>, or just be the change ID directly.
//
// Examples:
//   - "2024-01-15-feat-auth" -> "feat-auth"
//   - "feat-auth" -> "feat-auth"
//   - "2024-01-15-add-feature-x" -> "add-feature-x"
func ExtractChangeIDFromArchivePath(dirName string) string {
	// Try to match date prefix pattern: YYYY-MM-DD-
	// Date format: 4 digits, dash, 2 digits, dash, 2 digits, dash
	if len(dirName) > 11 {
		prefix := dirName[:11]
		// Check if it looks like a date prefix (YYYY-MM-DD-)
		if len(prefix) == 11 &&
			prefix[4] == '-' &&
			prefix[7] == '-' &&
			prefix[10] == '-' &&
			isDigits(prefix[0:4]) &&
			isDigits(prefix[5:7]) &&
			isDigits(prefix[8:10]) {
			return dirName[11:]
		}
	}

	// No date prefix, return as-is
	return dirName
}

// isDigits returns true if all characters in s are digits.
func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}
