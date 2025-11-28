// Package validation provides functions for collecting validation items.
package validation

import (
	"fmt"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/discovery"
)

// ValidationItem represents an item to validate
type ValidationItem struct {
	Name     string
	ItemType string // "change" or "spec"
	Path     string
}

// CreateValidationItems creates validation items from IDs and item type.
// The projectPath parameter is intentionally unused for now but kept for
// potential future use in path construction.
func CreateValidationItems(
	_ string,
	ids []string,
	itemType, basePath string,
) []ValidationItem {
	items := make([]ValidationItem, 0, len(ids))
	for _, id := range ids {
		var path string
		if itemType == ItemTypeSpec {
			path = filepath.Join(basePath, id, "spec.md")
		} else {
			path = filepath.Join(basePath, id)
		}
		items = append(items, ValidationItem{
			Name:     id,
			ItemType: itemType,
			Path:     path,
		})
	}

	return items
}

// GetAllItemsWithConfig returns all changes and specs using the provided
// config.
func GetAllItemsWithConfig(
	cfg *config.Config,
) ([]ValidationItem, error) {
	changes, err := GetChangeItemsWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	specs, err := GetSpecItemsWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	return append(changes, specs...), nil
}

// GetAllItems returns all changes and specs from the project path.
// Deprecated: Use GetAllItemsWithConfig for projects with custom root
// directories.
func GetAllItems(
	projectPath string,
) ([]ValidationItem, error) {
	cfg := &config.Config{
		RootDir:     config.DefaultRootDir,
		ProjectRoot: projectPath,
	}

	return GetAllItemsWithConfig(cfg)
}

// GetChangeItemsWithConfig returns all changes using the provided config.
func GetChangeItemsWithConfig(
	cfg *config.Config,
) ([]ValidationItem, error) {
	changeIDs, err := discovery.GetActiveChangesWithConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover changes: %w",
			err,
		)
	}

	return CreateValidationItems(
		cfg.ProjectRoot,
		changeIDs,
		ItemTypeChange,
		cfg.ChangesDir(),
	), nil
}

// GetChangeItems returns all changes from the project path.
// Deprecated: Use GetChangeItemsWithConfig for projects with custom root
// directories.
func GetChangeItems(
	projectPath string,
) ([]ValidationItem, error) {
	cfg := &config.Config{
		RootDir:     config.DefaultRootDir,
		ProjectRoot: projectPath,
	}

	return GetChangeItemsWithConfig(cfg)
}

// GetSpecItemsWithConfig returns all specs using the provided config.
func GetSpecItemsWithConfig(
	cfg *config.Config,
) ([]ValidationItem, error) {
	specIDs, err := discovery.GetSpecIDsWithConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to discover specs: %w", err)
	}

	return CreateValidationItems(
		cfg.ProjectRoot,
		specIDs,
		ItemTypeSpec,
		cfg.SpecsDir(),
	), nil
}

// GetSpecItems returns all specs from the project path.
// Deprecated: Use GetSpecItemsWithConfig for projects with custom root
// directories.
func GetSpecItems(
	projectPath string,
) ([]ValidationItem, error) {
	cfg := &config.Config{
		RootDir:     config.DefaultRootDir,
		ProjectRoot: projectPath,
	}

	return GetSpecItemsWithConfig(cfg)
}
