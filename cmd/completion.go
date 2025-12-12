// Package cmd provides command-line interface implementations.
// This file contains shell completion predictors for the spectr CLI.
// Predictors provide context-aware suggestions for tab completion in
// supported shells (bash, zsh, fish).
package cmd

import (
	"os"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/posener/complete"
)

// PredictChangeIDs returns a predictor that suggests active change IDs.
// It scans the spectr/changes/ directory for active changes, excluding
// the archive directory. Returns nil on error.
func PredictChangeIDs() complete.Predictor {
	return complete.PredictFunc(
		func(_ complete.Args) []string {
			projectPath, err := os.Getwd()
			if err != nil {
				return nil
			}

			changeIDs, err := discovery.GetActiveChangeIDs(
				projectPath,
			)
			if err != nil {
				return nil
			}

			return changeIDs
		},
	)
}

// PredictSpecIDs returns a predictor that suggests spec IDs.
// It scans the spectr/specs/ directory for specification directories.
// Returns nil on error.
func PredictSpecIDs() complete.Predictor {
	return complete.PredictFunc(
		func(_ complete.Args) []string {
			projectPath, err := os.Getwd()
			if err != nil {
				return nil
			}

			specIDs, err := discovery.GetSpecIDs(
				projectPath,
			)
			if err != nil {
				return nil
			}

			return specIDs
		},
	)
}

// PredictItemTypes returns a predictor that suggests item types.
// Valid types are "change" and "spec".
func PredictItemTypes() complete.Predictor {
	return complete.PredictSet("change", "spec")
}

// PredictItems returns a predictor that suggests both change and spec IDs.
// This is useful for commands that accept either type of item.
// Combines results from both changes and specs directories.
func PredictItems() complete.Predictor {
	return complete.PredictFunc(
		func(_ complete.Args) []string {
			projectPath, err := os.Getwd()
			if err != nil {
				return nil
			}

			var items []string

			changeIDs, err := discovery.GetActiveChangeIDs(
				projectPath,
			)
			if err == nil {
				items = append(
					items,
					changeIDs...)
			}

			specIDs, err := discovery.GetSpecIDs(
				projectPath,
			)
			if err == nil {
				items = append(items, specIDs...)
			}

			return items
		},
	)
}
