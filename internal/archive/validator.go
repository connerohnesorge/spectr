//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package archive

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/validation"
)

// ValidatePreArchive performs all pre-archive validation checks
// Returns validation report and error (error is for filesystem issues,
// not validation failures)
func ValidatePreArchive(
	changeDir string,
	strictMode bool,
) (*validation.ValidationReport, error) {
	// Derive spectrRoot from changeDir
	// changeDir format: /path/to/project/spectr/changes/<change-id>
	// spectrRoot should be: /path/to/project/spectr
	spectrRoot := filepath.Dir(filepath.Dir(changeDir))

	// Use existing change validation from validation package
	report, err := validation.ValidateChangeDeltaSpecs(changeDir, spectrRoot, strictMode)
	if err != nil {
		return nil, fmt.Errorf("validate change delta specs: %w", err)
	}

	return report, nil
}

// ValidatePostMerge validates a merged spec for correctness
// Ensures the merged spec has valid structure and no duplicate requirements
func ValidatePostMerge(mergedContent, _ string) error {
	// Write to temp file for validation
	tmpFile, err := os.CreateTemp("", "spec-*.md")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.WriteString(mergedContent); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	// Parse requirements from merged spec
	reqs, err := parsers.ParseRequirements(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("parse merged spec: %w", err)
	}

	// Check for duplicate requirement names (normalized)
	seen := make(map[string]bool)
	for _, req := range reqs {
		normalized := parsers.NormalizeRequirementName(req.Name)
		if seen[normalized] {
			return fmt.Errorf(
				"duplicate requirement name in merged spec: %q",
				req.Name,
			)
		}
		seen[normalized] = true
	}

	// Check that each requirement has at least one scenario
	for _, req := range reqs {
		scenarios := parsers.ParseScenarios(req.Raw)
		if len(scenarios) == 0 {
			return fmt.Errorf("requirement %q has no scenarios", req.Name)
		}
	}

	return nil
}

// ValidatePreMerge validates delta operations against base spec
// Checks that source requirements exist for MODIFIED/REMOVED/RENAMED
// and that target requirements don't exist for ADDED/RENAMED
//
// This is a wrapper around validation.ValidatePreMerge for backward compatibility.
//
//nolint:revive // specExists is a legitimate control parameter
func ValidatePreMerge(baseSpecPath string, deltaPlan *parsers.DeltaPlan, specExists bool) error {
	return validation.ValidatePreMerge(baseSpecPath, deltaPlan, specExists)
}

// CheckDuplicatesAndConflicts validates that there are no
// duplicate requirements within delta sections and no cross-section conflicts
func CheckDuplicatesAndConflicts(deltaPlan *parsers.DeltaPlan) error {
	// Check for duplicates within ADDED
	if err := checkDuplicatesInSection(
		deltaPlan.Added,
		"ADDED",
	); err != nil {
		return err
	}

	// Check for duplicates within MODIFIED
	if err := checkDuplicatesInSection(
		deltaPlan.Modified,
		"MODIFIED",
	); err != nil {
		return err
	}

	// Build normalized name sets
	nameSets := buildNameSets(deltaPlan)

	// Check for cross-section conflicts
	return checkCrossSectionConflicts(nameSets)
}

type nameSets struct {
	added       map[string]bool
	modified    map[string]bool
	removed     map[string]bool
	renamedFrom map[string]bool
	renamedTo   map[string]bool
}

func buildNameSets(deltaPlan *parsers.DeltaPlan) nameSets {
	sets := nameSets{
		added:       make(map[string]bool),
		modified:    make(map[string]bool),
		removed:     make(map[string]bool),
		renamedFrom: make(map[string]bool),
		renamedTo:   make(map[string]bool),
	}

	for _, req := range deltaPlan.Added {
		sets.added[parsers.NormalizeRequirementName(req.Name)] = true
	}

	for _, req := range deltaPlan.Modified {
		sets.modified[parsers.NormalizeRequirementName(req.Name)] = true
	}

	for _, name := range deltaPlan.Removed {
		sets.removed[parsers.NormalizeRequirementName(name)] = true
	}

	for _, op := range deltaPlan.Renamed {
		sets.renamedFrom[parsers.NormalizeRequirementName(op.From)] = true
		sets.renamedTo[parsers.NormalizeRequirementName(op.To)] = true
	}

	return sets
}

func checkCrossSectionConflicts(sets nameSets) error {
	// ADDED cannot conflict with MODIFIED, REMOVED, or RENAMED TO
	for name := range sets.added {
		if sets.modified[name] {
			return errors.New(
				"requirement appears in both ADDED and " +
					"MODIFIED sections",
			)
		}
		if sets.removed[name] {
			return errors.New(
				"requirement appears in both ADDED and " +
					"REMOVED sections",
			)
		}
		if sets.renamedTo[name] {
			return errors.New(
				"requirement appears in both ADDED and " +
					"RENAMED TO sections",
			)
		}
	}

	// MODIFIED cannot conflict with REMOVED or RENAMED FROM
	for name := range sets.modified {
		if sets.removed[name] {
			return errors.New(
				"requirement appears in both MODIFIED and " +
					"REMOVED sections",
			)
		}
		if sets.renamedFrom[name] {
			return errors.New(
				"requirement appears in both MODIFIED and " +
					"RENAMED FROM sections",
			)
		}
	}

	// REMOVED cannot conflict with RENAMED
	for name := range sets.removed {
		if sets.renamedFrom[name] {
			return errors.New(
				"requirement appears in both REMOVED and " +
					"RENAMED FROM sections",
			)
		}
	}

	return nil
}

// checkDuplicatesInSection checks for duplicate requirement names
// within a section
func checkDuplicatesInSection(
	reqs []parsers.RequirementBlock,
	sectionName string,
) error {
	seen := make(map[string]bool)
	for _, req := range reqs {
		normalized := parsers.NormalizeRequirementName(req.Name)
		if seen[normalized] {
			return fmt.Errorf(
				"duplicate requirement %q in %s section",
				req.Name,
				sectionName,
			)
		}
		seen[normalized] = true
	}

	return nil
}
