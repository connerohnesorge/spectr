// Package sync provides automatic synchronization of task statuses
// from tasks.jsonc (source of truth) back to tasks.md (human-readable format).
package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/discovery"
)

// SyncAllActiveChanges synchronizes task statuses from tasks.jsonc to tasks.md
// for all active changes in the project.
// Returns nil on success; errors are logged but don't block execution.
//
//nolint:revive // verbose flag is intentional for CLI integration
func SyncAllActiveChanges(projectRoot string, verbose bool) error {
	changeIDs, err := discovery.GetActiveChanges(projectRoot)
	if err != nil {
		// Directory doesn't exist or other error - skip silently
		return nil
	}

	var totalSynced int
	for _, id := range changeIDs {
		changeDir := filepath.Join(projectRoot, "spectr", "changes", id)

		synced, err := SyncTasksToMarkdown(changeDir)
		if err != nil {
			// Log error but continue with other changes
			fmt.Fprintf(os.Stderr, "sync: %s: %v\n", id, err)

			continue
		}

		if verbose && synced > 0 {
			fmt.Printf("Synced %d task statuses in %s\n", synced, id)
		}
		totalSynced += synced
	}

	return nil
}
