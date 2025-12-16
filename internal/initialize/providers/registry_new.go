// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the new instance-only registry for the redesigned
// architecture.
//
// This registry has NO global state.
//
// All registration and lookup happens on explicit registry instances.
//
// Usage:
//
//	reg := providers.CreateRegistry()
//	reg.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//	all := reg.All() // Returns registrations sorted by priority
package providers

import (
	"fmt"
	"sort"
	"sync"
)

// ProviderRegistry is the new instance-only registry for the redesigned
// provider architecture.
//
// It stores Registration structs (which include provider metadata and the
// provider itself).
//
// This registry has NO global state - create instances with CreateRegistry().
// This design improves testability by eliminating shared state between tests.
//
// Thread-safety: All methods are safe for concurrent access.
type ProviderRegistry struct {
	mu            sync.RWMutex
	registrations map[string]Registration
}

// CreateRegistry creates a new empty registry instance.
// This is the primary constructor for the new provider registry.
//
// Example:
//
//	reg := providers.CreateRegistry()
//	reg.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
func CreateRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		registrations: make(
			map[string]Registration,
		),
	}
}

// Register adds a provider registration to the registry.
// Returns an error if a registration with the same ID already exists.
//
// Unlike the old global Register() function which panics on duplicates,
// this method returns an error for explicit error handling.
//
// Example:
//
//	err := reg.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//	if err != nil {
//	    // Handle duplicate registration
//	}
func (r *ProviderRegistry) Register(
	reg Registration,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf(
			"provider %q already registered",
			reg.ID,
		)
	}

	r.registrations[reg.ID] = reg

	return nil
}

// All returns all registrations sorted by priority (lower = higher priority).
// Returns an empty slice if no providers are registered.
//
// Example:
//
//	all := reg.All()
//	for _, r := range all {
//	    fmt.Printf("%s (priority %d)\n", r.Name, r.Priority)
//	}
func (r *ProviderRegistry) All() []Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(
		[]Registration,
		0,
		len(r.registrations),
	)
	for _, reg := range r.registrations {
		result = append(result, reg)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

// Get retrieves a registration by its provider ID.
// Returns nil if no registration with that ID exists.
//
// Example:
//
//	if claude := reg.Get("claude-code"); claude != nil {
//	    fmt.Println(claude.Name) // "Claude Code"
//	}
func (r *ProviderRegistry) Get(
	id string,
) *Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if reg, exists := r.registrations[id]; exists {
		return &reg
	}

	return nil
}

// IDs returns all registered provider IDs sorted by priority
// (lower = higher priority).
// Returns an empty slice if no providers are registered.
//
// Example:
//
//	ids := reg.IDs() // e.g., ["claude-code", "cursor", "cline"]
func (r *ProviderRegistry) IDs() []string {
	all := r.All()
	ids := make([]string, len(all))
	for i, reg := range all {
		ids[i] = reg.ID
	}

	return ids
}

// Count returns the number of registered providers.
//
// Example:
//
//	if reg.Count() == 0 {
//	    fmt.Println("No providers registered")
//	}
func (r *ProviderRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.registrations)
}
