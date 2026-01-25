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

// ValidFrontmatterKeys defines all known valid frontmatter field names.
// Used to validate overrides and catch typos.
//
// Claude: https://code.claude.com/docs/en/slash-commands
// OpenCode: https://opencode.ai/docs/commands/
var ValidFrontmatterKeys = map[string]bool{
	"description":   true,
	"allowed-tools": true,
	"subtask":       true,
	"context":       false, // Claude Code: "fork" runs in forked sub-agent context
	"agent":         true,  // Agent routing (e.g., "plan" for planning subagent)
}

// ValidateFrontmatterOverride checks that all keys in an override are known valid keys.
// Returns an error listing any unknown keys found.
// This helps catch typos like "contxt" instead of "context".
func ValidateFrontmatterOverride(
	overrides *FrontmatterOverride,
) error {
	if overrides == nil {
		return nil
	}

	var unknownKeys []string

	// Check Set keys
	for k := range overrides.Set {
		if !ValidFrontmatterKeys[k] {
			unknownKeys = append(unknownKeys, k)
		}
	}

	// Check Remove keys
	for _, k := range overrides.Remove {
		if !ValidFrontmatterKeys[k] {
			unknownKeys = append(unknownKeys, k)
		}
	}

	if len(unknownKeys) > 0 {
		sort.Strings(unknownKeys)

		return fmt.Errorf(
			"unknown frontmatter keys: %v",
			unknownKeys,
		)
	}

	return nil
}

// BaseSlashCommandFrontmatter defines default frontmatter for each slash command.
// Templates (.tmpl files) contain only body content; frontmatter is data.
var BaseSlashCommandFrontmatter = map[SlashCommand]map[string]any{
	SlashProposal: {
		"description":   "Proposal Creation Guide (project)",
		"allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"subtask":       false,
	},
	SlashApply: {
		"description":   "Change Proposal Application/Acceptance Process (project)",
		"allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"subtask":       false,
	},
	SlashNext: {
		"description":   "Spectr: Next Task Execution",
		"allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"subtask":       false,
	},
}

// GetBaseFrontmatter returns a copy of the base frontmatter for a command.
// Returns empty map if command not found.
func GetBaseFrontmatter(
	cmd SlashCommand,
) map[string]any {
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
	case []int:
		return copyIntSlice(val)
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

// copyIntSlice creates a copy of a []int.
func copyIntSlice(src []int) []int {
	if src == nil {
		return nil
	}

	dst := make([]int, len(src))
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
//
// Keys are sorted alphabetically to ensure deterministic output.
func RenderFrontmatter(
	fm map[string]any,
	body string,
) (string, error) {
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

		// Build yaml.Node with explicit key ordering
		// Go maps don't preserve insertion order, so we use yaml.Node
		// to guarantee deterministic output
		mapNode := &yaml.Node{
			Kind: yaml.MappingNode,
		}
		for _, k := range keys {
			// Add key node
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: k,
			}
			// Add value node
			valNode := &yaml.Node{}
			if err := valNode.Encode(fm[k]); err != nil {
				return "", fmt.Errorf(
					"failed to encode frontmatter field %q: %w",
					k,
					err,
				)
			}
			mapNode.Content = append(
				mapNode.Content,
				keyNode,
				valNode,
			)
		}

		// Encode the ordered node
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(
			0,
		) // No extra indentation
		if err := encoder.Encode(mapNode); err != nil {
			return "", fmt.Errorf(
				"failed to encode frontmatter as YAML: %w",
				err,
			)
		}
		if err := encoder.Close(); err != nil {
			return "", fmt.Errorf(
				"failed to close YAML encoder: %w",
				err,
			)
		}
	}

	// Write closing fence
	buf.WriteString("---\n")

	// Append body content
	buf.WriteString(body)

	return buf.String(), nil
}
