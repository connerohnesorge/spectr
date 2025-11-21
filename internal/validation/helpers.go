// Package validation provides validation helper functions.
package validation

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/connerohnesorge/spectr/internal/discovery"
)

const (
	// ItemTypeChange represents a change item type
	ItemTypeChange = "change"
	// ItemTypeSpec represents a spec item type
	ItemTypeSpec = "spec"
	// SpectrDir is the base directory for spectr files
	SpectrDir = "spectr"
)

// ItemTypeInfo holds information about an item's type
type ItemTypeInfo struct {
	ItemType string
	IsChange bool
	IsSpec   bool
}

// DetermineItemType determines if an item is a change or spec
func DetermineItemType(
	projectPath, itemName string,
	typeFlag *string,
) (ItemTypeInfo, error) {
	changes, err := discovery.GetActiveChangeIDs(projectPath)
	if err != nil {
		return ItemTypeInfo{}, fmt.Errorf(
			"failed to discover changes: %w",
			err,
		)
	}

	specs, err := discovery.GetSpecIDs(projectPath)
	if err != nil {
		return ItemTypeInfo{}, fmt.Errorf(
			"failed to discover specs: %w",
			err,
		)
	}

	info := ItemTypeInfo{
		IsChange: containsString(changes, itemName),
		IsSpec:   containsString(specs, itemName),
	}

	// Handle explicit type flag
	if typeFlag != nil {
		info.ItemType = *typeFlag
		if info.ItemType == ItemTypeChange && !info.IsChange {
			return info, fmt.Errorf(
				"change '%s' not found",
				itemName,
			)
		}
		if info.ItemType == ItemTypeSpec && !info.IsSpec {
			return info, fmt.Errorf("spec '%s' not found", itemName)
		}

		return info, nil
	}

	// Auto-detect type
	if info.IsChange && info.IsSpec {
		return info, fmt.Errorf(
			"item '%s' exists as both change and spec, "+
				"use --type flag to disambiguate",
			itemName,
		)
	}
	if !info.IsChange && !info.IsSpec {
		return info, fmt.Errorf("item '%s' not found", itemName)
	}
	if info.IsChange {
		info.ItemType = ItemTypeChange
	} else {
		info.ItemType = ItemTypeSpec
	}

	return info, nil
}

// ValidateItemByType validates an item based on its type
func ValidateItemByType(
	validator *Validator,
	projectPath, itemName, itemType string,
) (*ValidationReport, error) {
	if itemType == ItemTypeChange {
		changePath := filepath.Join(
			projectPath,
			SpectrDir,
			"changes",
			itemName,
		)

		return validator.ValidateChange(changePath)
	}

	specPath := filepath.Join(
		projectPath,
		SpectrDir,
		"specs",
		itemName,
		"spec.md",
	)

	return validator.ValidateSpec(specPath)
}

// ValidateSingleItem validates a single item and returns the bulk result
func ValidateSingleItem(
	validator *Validator,
	item ValidationItem,
) (BulkResult, error) {
	var report *ValidationReport
	var err error

	if item.ItemType == ItemTypeChange {
		report, err = validator.ValidateChange(item.Path)
	} else {
		report, err = validator.ValidateSpec(item.Path)
	}

	if err != nil {
		return BulkResult{
			Name:  item.Name,
			Type:  item.ItemType,
			Valid: false,
			Error: err.Error(),
		}, err
	}

	return BulkResult{
		Name:   item.Name,
		Type:   item.ItemType,
		Valid:  report.Valid,
		Report: report,
	}, nil
}

// containsString checks if a slice contains a string
func containsString(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
