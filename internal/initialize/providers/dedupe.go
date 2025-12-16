// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file provides deduplication logic for initializers. When multiple
// providers return the same initializer (same type and configuration), we
// only run it once.
package providers

import "fmt"

// Keyer is an optional interface for initializers that support deduplication.
// Initializers implementing this interface can be deduplicated based on their
// key.
//
// The key should be deterministic and based on the initializer's type and
// configuration. Two initializers with the same key are considered duplicates,
// and only the first occurrence is kept when deduplicating.
//
// Example key formats:
//   - "dir:.claude/commands/spectr" for a DirectoryInitializer
//   - "config:CLAUDE.md" for a ConfigFileInitializer
//   - "slashcmds:.claude/commands/spectr:.md:0" for SlashCommandsInitializer
type Keyer interface {
	// Key returns a unique key for deduplication.
	Key() string
}

// DedupeInitializers removes duplicate initializers from a slice.
// Deduplication is based on the Key() method for initializers implementing
// Keyer. For initializers that don't implement Keyer, they are never
// deduplicated (each is unique).
//
// The function preserves order - the first occurrence of each key is kept.
//
// Example usage:
//
//	allInitializers := collectFromProviders(providers)
//	dedupedInitializers := DedupeInitializers(allInitializers)
//	for _, init := range dedupedInitializers {
//	    init.Init(ctx, fs, cfg)
//	}
func DedupeInitializers(
	all []Initializer,
) []Initializer {
	if len(all) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var result []Initializer

	for _, init := range all {
		key := initializerKey(init)
		if !seen[key] {
			seen[key] = true
			result = append(result, init)
		}
	}

	return result
}

// initializerKey returns a unique key for an initializer.
// If the initializer implements Keyer, its Key() method is used.
// Otherwise, a fallback key based on the pointer address is generated,
// ensuring the initializer is never considered a duplicate.
func initializerKey(init Initializer) string {
	if keyer, ok := init.(Keyer); ok {
		return keyer.Key()
	}
	// Fallback: use pointer address to ensure uniqueness
	// Each non-Keyer initializer is considered unique
	return fmt.Sprintf("ptr:%p", init)
}
