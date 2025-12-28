package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GetSpecs finds all specs in {spectrDir}/specs/ that contain spec.md
func GetSpecs(
	projectPath string,
	spectrDir string,
) ([]string, error) {
	specsDir := filepath.Join(
		projectPath,
		spectrDir,
		"specs",
	)

	// Check if specs directory exists
	if _, err := os.Stat(specsDir); os.IsNotExist(
		err,
	) {
		return make([]string, 0), nil
	}

	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read specs directory: %w",
			err,
		)
	}

	var specs []string
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Check if spec.md exists
		specPath := filepath.Join(
			specsDir,
			entry.Name(),
			"spec.md",
		)
		if _, err := os.Stat(specPath); err == nil {
			specs = append(specs, entry.Name())
		}
	}

	// Sort alphabetically for consistency
	sort.Strings(specs)

	return specs, nil
}

// GetSpecIDs returns a list of spec IDs (directory names under {spectrDir}/specs/)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetSpecIDs(
	projectRoot string,
	spectrDir string,
) ([]string, error) {
	return GetSpecs(projectRoot, spectrDir)
}
