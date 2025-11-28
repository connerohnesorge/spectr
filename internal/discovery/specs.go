package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/connerohnesorge/spectr/internal/config"
)

// GetSpecsWithConfig finds all specs using the provided config
func GetSpecsWithConfig(cfg *config.Config) ([]string, error) {
	specsDir := cfg.SpecsPath()

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
// This function loads config from the project path.
// For direct config usage, use GetSpecsWithConfig.
func GetSpecs(projectPath string) ([]string, error) {
	cfg, err := config.LoadFromPath(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return GetSpecsWithConfig(cfg)
}

// GetSpecIDs returns a list of spec IDs (directory names under spectr/specs/)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetSpecIDs(projectRoot string) ([]string, error) {
	return GetSpecs(projectRoot)
}

// GetSpecIDsWithConfig returns a list of spec IDs using the provided config
func GetSpecIDsWithConfig(cfg *config.Config) ([]string, error) {
	return GetSpecsWithConfig(cfg)
}
