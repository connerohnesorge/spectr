package pr

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrEmptyChangeID is returned when an empty changeID is provided.
var ErrEmptyChangeID = errors.New("changeID cannot be empty")

// RemoveChangeDirectory removes a change proposal directory from the project.
// It validates that the path is within the expected spectr/changes/ directory
// before performing the removal for safety.
func RemoveChangeDirectory(projectRoot, changeID string) error {
	if changeID == "" {
		return ErrEmptyChangeID
	}

	// Construct the expected change directory path
	changesDir := filepath.Join(projectRoot, "spectr", "changes")
	changeDir := filepath.Join(changesDir, changeID)

	// Validate the path is within the expected directory structure
	// by checking that the cleaned path starts with the changes directory
	cleanedChangeDir := filepath.Clean(changeDir)
	cleanedChangesDir := filepath.Clean(changesDir)
	sep := string(filepath.Separator)

	if !strings.HasPrefix(cleanedChangeDir, cleanedChangesDir+sep) {
		return fmt.Errorf(
			"invalid change directory path: %s is not within %s",
			cleanedChangeDir,
			cleanedChangesDir,
		)
	}

	// Verify the directory exists before attempting removal
	info, err := os.Stat(changeDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"change directory does not exist: %s",
				changeDir,
			)
		}

		return fmt.Errorf("failed to access change directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", changeDir)
	}

	// Remove the directory recursively
	if err := os.RemoveAll(changeDir); err != nil {
		return fmt.Errorf(
			"failed to remove change directory %s: %w",
			changeDir,
			err,
		)
	}

	return nil
}
