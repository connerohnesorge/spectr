package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/connerohnesorge/spectr/internal/config"
)

// GetSpecsWithConfig finds all specs using the provided config.
// It searches in the specs directory specified by the config
// for directories that contain spec.md files.
func GetSpecsWithConfig(cfg *config.Config) ([]string, error) {
	specsDir := cfg.SpecsDir()

	// Check if specs directory exists
	if _, err := os.Stat(specsDir); os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read specs directory: %w", err)
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
		specPath := filepath.Join(specsDir, entry.Name(), "spec.md")
		if _, err := os.Stat(specPath); err == nil {
			specs = append(specs, entry.Name())
		}
	}

	// Sort alphabetically for consistency
	sort.Strings(specs)

	return specs, nil
}

// GetSpecs finds all specs in spectr/specs/ that contain spec.md.
// Deprecated: Use GetSpecsWithConfig for projects with custom root directories.
func GetSpecs(projectPath string) ([]string, error) {
	cfg := &config.Config{
		RootDir:     config.DefaultRootDir,
		ProjectRoot: projectPath,
	}

	return GetSpecsWithConfig(cfg)
}

// GetSpecIDsWithConfig returns a list of spec IDs using the provided config
// (directory names under the configured specs directory)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetSpecIDsWithConfig(cfg *config.Config) ([]string, error) {
	return GetSpecsWithConfig(cfg)
}

// GetSpecIDs returns a list of spec IDs (directory names under spectr/specs/).
// Returns empty slice (not error) if the directory doesn't exist.
// Results are sorted alphabetically for consistency.
// Deprecated: Use GetSpecIDsWithConfig for projects with custom root
// directories.
func GetSpecIDs(projectRoot string) ([]string, error) {
	return GetSpecs(projectRoot)
}
