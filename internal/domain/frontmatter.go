// Package domain contains shared domain types used across the Spectr codebase.
package domain

import (
	"bytes"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

// FrontmatterOverride specifies modifications to slash command frontmatter.
// Set operations are applied first, then Remove operations.
type FrontmatterOverride struct {
	// Set contains fields to add or modify.
	// Values must be YAML-serializable (string, bool, int, []string, map, etc.)
	Set map[string]any

	// Remove contains field names to delete.
	// Applied after Set to allow field replacement.
	Remove []string
}

// BaseSlashCommandFrontmatter defines default frontmatter for each slash command.
// Templates (.tmpl files) contain only body content; frontmatter is data.
var BaseSlashCommandFrontmatter = map[SlashCommand]map[string]any{
	SlashProposal: {
		"description":   "Proposal Creation Guide (project)",
		"allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"agent":         "plan",
		"subtask":       false,
	},
	SlashApply: {
		"description":   "Change Proposal Application/Acceptance Process (project)",
		"allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"subtask":       false,
	},
}

// GetBaseFrontmatter returns a copy of the base frontmatter for a command.
// Returns empty map if command not found.
func GetBaseFrontmatter(cmd SlashCommand) map[string]any {
	base, ok := BaseSlashCommandFrontmatter[cmd]
	if !ok {
		return make(map[string]any)
	}

	// Return deep copy to prevent mutation
	return copyMap(base)
}

// copyMap creates a deep copy of a map[string]any.
// It handles nested maps and slices.
func copyMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = copyValue(v)
	}

	return dst
}

// copyValue creates a deep copy of an any value.
// It handles maps, slices, and primitive types.
func copyValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return copyMap(val)
	case []any:
		return copySlice(val)
	case []string:
		return copyStringSlice(val)
	default:
		// Primitive types (string, bool, int, etc.) are copied by value
		return v
	}
}

// copySlice creates a deep copy of a []any.
func copySlice(src []any) []any {
	if src == nil {
		return nil
	}

	dst := make([]any, len(src))
	for i, v := range src {
		dst[i] = copyValue(v)
	}

	return dst
}

// copyStringSlice creates a copy of a []string.
func copyStringSlice(src []string) []string {
	if src == nil {
		return nil
	}

	dst := make([]string, len(src))
	copy(dst, src)

	return dst
}

// ApplyFrontmatterOverrides applies Set and Remove operations to a frontmatter map.
// Returns a new map; does not mutate the input.
//
// Operation order:
//  1. Copy base map
//  2. Apply all Set operations (merge/overwrite)
//  3. Apply all Remove operations (delete keys)
func ApplyFrontmatterOverrides(
	base map[string]any,
	overrides *FrontmatterOverride,
) map[string]any {
	// Copy base to avoid mutation
	result := copyMap(base)
	if result == nil {
		result = make(map[string]any)
	}

	// If no overrides, return copy of base
	if overrides == nil {
		return result
	}

	// Apply Set operations
	for k, v := range overrides.Set {
		result[k] = copyValue(v)
	}

	// Apply Remove operations (after Set)
	for _, k := range overrides.Remove {
		delete(result, k)
	}

	return result
}

// RenderFrontmatter serializes a frontmatter map to YAML and combines with body.
// Returns markdown with YAML frontmatter block.
//
// Output format:
//
//	---
//	key: value
//	---
//	Body content
func RenderFrontmatter(fm map[string]any, body string) (string, error) {
	var buf bytes.Buffer

	// Write opening fence
	buf.WriteString("---\n")

	// Only render YAML if there are fields
	if len(fm) > 0 {
		// Sort keys for deterministic output
		keys := make([]string, 0, len(fm))
		for k := range fm {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Create ordered map for YAML rendering
		orderedFm := make(map[string]any, len(fm))
		for _, k := range keys {
			orderedFm[k] = fm[k]
		}

		// Encode YAML
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(0) // No extra indentation
		if err := encoder.Encode(orderedFm); err != nil {
			return "", fmt.Errorf("failed to encode frontmatter as YAML: %w", err)
		}
		if err := encoder.Close(); err != nil {
			return "", fmt.Errorf("failed to close YAML encoder: %w", err)
		}
	}

	// Write closing fence
	buf.WriteString("---\n")

	// Append body content
	buf.WriteString(body)

	return buf.String(), nil
}
