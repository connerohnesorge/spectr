package providers

import (
	"fmt"
	"maps"
	"sort"
	"sync"
)

// registry is the global provider registry.
var (
	registry     = make(map[string]Provider)
	registryLock sync.RWMutex
)

// Register adds a provider to the global registry.
// This is typically called from init() in each provider file.
// Panics if a provider with the same ID is already registered.
func Register(p Provider) {
	registryLock.Lock()
	defer registryLock.Unlock()

	if _, exists := registry[p.ID()]; exists {
		panic(
			fmt.Sprintf(
				"provider %q already registered",
				p.ID(),
			),
		)
	}

	registry[p.ID()] = p
}

// Get retrieves a provider by its ID.
// Returns nil if the provider is not found.
func Get(id string) Provider {
	registryLock.RLock()
	defer registryLock.RUnlock()

	return registry[id]
}

// All returns all registered providers sorted by priority.
func All() []Provider {
	registryLock.RLock()
	defer registryLock.RUnlock()

	providers := make(
		[]Provider,
		0,
		len(registry),
	)
	for _, p := range registry {
		providers = append(providers, p)
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Priority() < providers[j].Priority()
	})

	return providers
}

// IDs returns all registered provider IDs sorted by priority.
func IDs() []string {
	providers := All()
	ids := make([]string, len(providers))
	for i, p := range providers {
		ids[i] = p.ID()
	}

	return ids
}

// Count returns the number of registered providers.
func Count() int {
	registryLock.RLock()
	defer registryLock.RUnlock()

	return len(registry)
}

// WithConfigFile returns all providers that have an instruction file,
// sorted by priority.
func WithConfigFile() []Provider {
	registryLock.RLock()
	defer registryLock.RUnlock()

	var providers []Provider
	for _, p := range registry {
		if p.HasConfigFile() {
			providers = append(providers, p)
		}
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Priority() < providers[j].Priority()
	})

	return providers
}

// WithSlashCommands returns all providers that have slash commands,
// sorted by priority.
func WithSlashCommands() []Provider {
	registryLock.RLock()
	defer registryLock.RUnlock()

	var providers []Provider
	for _, p := range registry {
		if p.HasSlashCommands() {
			providers = append(providers, p)
		}
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Priority() < providers[j].Priority()
	})

	return providers
}

// Reset clears the global registry. Only use in tests.
func Reset() {
	registryLock.Lock()
	defer registryLock.Unlock()

	registry = make(map[string]Provider)
}

// Registry provides an instance-based registry for cases where
// global state is not desired (e.g., testing).
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// NewRegistryFromGlobal creates a registry populated with all globally
// ) registered providers.
func NewRegistryFromGlobal() *Registry {
	r := NewRegistry()

	registryLock.RLock()
	defer registryLock.RUnlock()

	maps.Copy(r.providers, registry)

	return r
}

// Register adds a provider to this registry.
func (r *Registry) Register(p Provider) error {
	if _, exists := r.providers[p.ID()]; exists {
		return fmt.Errorf(
			"provider %q already registered",
			p.ID(),
		)
	}

	r.providers[p.ID()] = p

	return nil
}

// Get retrieves a provider by its ID.
func (r *Registry) Get(id string) Provider {
	return r.providers[id]
}

// All returns all providers in this registry sorted by priority.
func (r *Registry) All() []Provider {
	providers := make(
		[]Provider,
		0,
		len(r.providers),
	)
	for _, p := range r.providers {
		providers = append(providers, p)
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Priority() < providers[j].Priority()
	})

	return providers
}

// IDs returns all provider IDs in this registry sorted by priority.
func (r *Registry) IDs() []string {
	providers := r.All()
	ids := make([]string, len(providers))
	for i, p := range providers {
		ids[i] = p.ID()
	}

	return ids
}

// Count returns the number of providers in this registry.
func (r *Registry) Count() int {
	return len(r.providers)
}
