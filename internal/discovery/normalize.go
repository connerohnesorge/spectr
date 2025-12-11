package discovery

import (
	"path/filepath"
	"strings"
)

const pathSeparator = "/"

// NormalizeItemPath extracts the item ID and infers type from a path argument.
// Returns the normalized ID and inferred type:
// "change", "spec", or "" if not inferred from path.
//
// Path patterns recognized:
//   - spectr/changes/<id> -> id, "change"
//   - spectr/changes/<id>/specs/foo/spec.md -> id, "change"
//   - spectr/specs/<id> -> id, "spec"
//   - spectr/specs/<id>/spec.md -> id, "spec"
//   - /absolute/path/spectr/changes/<id> -> id, "change"
//   - /absolute/path/spectr/specs/<id> -> id, "spec"
//   - simple-id -> simple-id, "" (no inferred type)
func NormalizeItemPath(input string) (id, inferredType string) {
	if input == "" {
		return "", ""
	}

	// Normalize path separators to forward slashes for consistent parsing
	normalized := filepath.ToSlash(input)

	// Look for spectr/changes/<id> pattern
	if idx := strings.Index(normalized, "spectr/changes/"); idx != -1 {
		// Extract everything after "spectr/changes/"
		remainder := normalized[idx+len("spectr/changes/"):]
		// The ID is the first path component after spectr/changes/
		changeID := extractFirstPathComponent(remainder)
		if changeID != "" && changeID != "archive" {
			return changeID, "change"
		}
	}

	// Look for spectr/specs/<id> pattern
	if idx := strings.Index(normalized, "spectr/specs/"); idx != -1 {
		// Extract everything after "spectr/specs/"
		remainder := normalized[idx+len("spectr/specs/"):]
		// The ID is the first path component after spectr/specs/
		specID := extractFirstPathComponent(remainder)
		if specID != "" {
			return specID, "spec"
		}
	}

	// No recognized pattern - return input as-is with no inferred type
	// Clean up: remove trailing slashes, get base name if it looks like a path
	cleaned := strings.TrimSuffix(normalized, pathSeparator)
	if cleaned == "" {
		return input, ""
	}

	// If there are no slashes, it's a simple ID
	if !strings.Contains(cleaned, pathSeparator) {
		return cleaned, ""
	}

	// If it looks like a path but doesn't match our patterns,
	// return the original input as-is
	return input, ""
}

// extractFirstPathComponent extracts the first path component from a string.
// For "my-change/specs/foo/spec.md" returns "my-change"
// For "my-change" returns "my-change"
// For "" returns ""
func extractFirstPathComponent(path string) string {
	trimmed := strings.TrimPrefix(path, pathSeparator)
	trimmed = strings.TrimSuffix(trimmed, pathSeparator)

	if trimmed == "" {
		return ""
	}

	if idx := strings.Index(trimmed, pathSeparator); idx != -1 {
		return trimmed[:idx]
	}

	return trimmed
}
