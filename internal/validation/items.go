// Package validation provides functions for collecting validation items.
package validation

import (
	"fmt"
	"path/filepath"

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

// GetAllItems returns all changes and specs from the project path.
func GetAllItems(
	projectPath string,
) ([]ValidationItem, error) {
	changes, err := GetChangeItems(projectPath)
	if err != nil {
		return nil, err
	}

	specs, err := GetSpecItems(projectPath)
	if err != nil {
		return nil, err
	}

	return append(changes, specs...), nil
}

// GetChangeItems returns all changes from the project path.
func GetChangeItems(
	projectPath string,
) ([]ValidationItem, error) {
	changeIDs, err := discovery.GetActiveChangeIDs(projectPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover changes: %w",
			err,
		)
	}

	basePath := filepath.Join(projectPath, SpectrDir, "changes")

	return CreateValidationItems(
		projectPath,
		changeIDs,
		ItemTypeChange,
		basePath,
	), nil
}

// GetSpecItems returns all specs from the project path.
func GetSpecItems(
	projectPath string,
) ([]ValidationItem, error) {
	specIDs, err := discovery.GetSpecIDs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover specs: %w", err)
	}

	basePath := filepath.Join(projectPath, SpectrDir, "specs")

	return CreateValidationItems(
		projectPath,
		specIDs,
		ItemTypeSpec,
		basePath,
	), nil
}
