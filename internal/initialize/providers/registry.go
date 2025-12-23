package providers

import (
	"fmt"
	"sort"
	"sync"
)

// Registry is the registry implementation that stores Registration
// structs.
// This separates provider metadata (ID, name, priority)
// from provider implementation,
// allowing for cleaner registration and reduced boilerplate.
//
// Registry is thread-safe and can be used concurrently.
type Registry struct {
	mu            sync.RWMutex
	registrations map[string]Registration
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		registrations: make(map[string]Registration),
	}
}

// Global default instance for package-level registration
var (
	defaultRegistry     *Registry
	defaultRegistryOnce sync.Once
)

// getDefaultRegistry returns the global default registry instance.
func getDefaultRegistry() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})

	return defaultRegistry
}

// Register adds a registration to the global registry.
// Returns an error if a provider with the same ID is already registered.
//
// Example usage:
//
//	err := providers.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
func Register(reg Registration) error {
	return getDefaultRegistry().Register(reg)
}

// Get retrieves a registration by provider ID from the global registry.
// Returns nil if the provider is not found.
func Get(id string) *Registration {
	return getDefaultRegistry().Get(id)
}

// All returns all registrations from the global registry sorted by priority.
// Lower priority values come first.
func All() []Registration {
	return getDefaultRegistry().All()
}

// IDs returns all provider IDs from the global registry sorted by priority.
func IDs() []string {
	return getDefaultRegistry().IDs()
}

// Count returns the number of registered providers in the global registry.
func Count() int {
	return getDefaultRegistry().Count()
}

// Reset clears the global registry. Only use in tests.
func Reset() {
	getDefaultRegistry().mu.Lock()
	defer getDefaultRegistry().mu.Unlock()
	getDefaultRegistry().registrations = make(map[string]Registration)
}

// Register adds a registration to this registry instance.
// Returns an error if a provider with the same ID is already registered.
func (r *Registry) Register(reg Registration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf("provider with ID %q already registered", reg.ID)
	}

	r.registrations[reg.ID] = reg

	return nil
}

// Get retrieves a registration by provider ID.
// Returns nil if the provider is not found.
func (r *Registry) Get(id string) *Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if reg, exists := r.registrations[id]; exists {
		return &reg
	}

	return nil
}

// All returns all registrations sorted by priority (lower = higher priority).
func (r *Registry) All() []Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	regs := make([]Registration, 0, len(r.registrations))
	for _, reg := range r.registrations {
		regs = append(regs, reg)
	}

	// Sort by priority (lower values first)
	sort.Slice(regs, func(i, j int) bool {
		return regs[i].Priority < regs[j].Priority
	})

	return regs
}

// IDs returns all provider IDs sorted by their registration's priority.
func (r *Registry) IDs() []string {
	regs := r.All()
	ids := make([]string, len(regs))
	for i, reg := range regs {
		ids[i] = reg.ID
	}

	return ids
}

// Count returns the number of registrations in this registry.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.registrations)
}
