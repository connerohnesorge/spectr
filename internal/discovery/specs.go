package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GetSpecs finds all specs in spectr/specs/ that contain spec.md
func GetSpecs(projectPath string) ([]string, error) {
	specsDir := filepath.Join(projectPath, "spectr", "specs")

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

// GetSpecIDs returns a list of spec IDs (directory names under spectr/specs/)
// Returns empty slice (not error) if the directory doesn't exist
// Results are sorted alphabetically for consistency
func GetSpecIDs(projectRoot string) ([]string, error) {
	return GetSpecs(projectRoot)
}

// SpecResolveResult contains the resolved spec ID and whether it was a partial
// match.
type SpecResolveResult struct {
	SpecID       string
	PartialMatch bool
}

// ResolveSpecID resolves a partial spec ID to a full spec ID.
// The resolution algorithm:
//  1. Exact match: returns immediately if partialID exactly matches a spec ID
//  2. Prefix match: finds all spec IDs that start with partialID
//     (case-insensitive)
//  3. Substring match: if no prefix matches, finds all spec IDs containing
//     partialID (case-insensitive)
//
// Returns an error if:
// - No specs match the partial ID
// - Multiple specs match the partial ID (ambiguous)
func ResolveSpecID(partialID, projectRoot string) (SpecResolveResult, error) {
	specs, err := GetSpecIDs(projectRoot)
	if err != nil {
		return SpecResolveResult{},
			fmt.Errorf("get specs: %w", err)
	}

	if len(specs) == 0 {
		return SpecResolveResult{},
			fmt.Errorf("no spec found matching '%s'", partialID)
	}

	partialLower := strings.ToLower(partialID)

	// Check for exact match first
	for _, spec := range specs {
		if spec == partialID {
			return SpecResolveResult{
				SpecID:       spec,
				PartialMatch: false,
			}, nil
		}
	}

	// Try prefix matching (case-insensitive)
	var prefixMatches []string
	for _, spec := range specs {
		if strings.HasPrefix(strings.ToLower(spec), partialLower) {
			prefixMatches = append(prefixMatches, spec)
		}
	}

	if len(prefixMatches) == 1 {
		return SpecResolveResult{
			SpecID:       prefixMatches[0],
			PartialMatch: true,
		}, nil
	}

	if len(prefixMatches) > 1 {
		return SpecResolveResult{}, fmt.Errorf(
			"ambiguous ID '%s' matches multiple specs: %s",
			partialID,
			strings.Join(prefixMatches, ", "),
		)
	}

	// Try substring matching (case-insensitive) as fallback
	var substringMatches []string
	for _, spec := range specs {
		if strings.Contains(strings.ToLower(spec), partialLower) {
			substringMatches = append(substringMatches, spec)
		}
	}

	if len(substringMatches) == 1 {
		return SpecResolveResult{
			SpecID:       substringMatches[0],
			PartialMatch: true,
		}, nil
	}

	if len(substringMatches) > 1 {
		return SpecResolveResult{}, fmt.Errorf(
			"ambiguous ID '%s' matches multiple specs: %s",
			partialID,
			strings.Join(substringMatches, ", "),
		)
	}

	return SpecResolveResult{},
		fmt.Errorf(
			"no spec found matching '%s'",
			partialID,
		)
}
