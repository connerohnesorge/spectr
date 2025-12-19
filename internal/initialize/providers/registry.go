// Package providers defines the core interfaces for the provider architecture.
//
// This file contains the provider registry that uses the Registration struct
// for metadata. Providers are registered with their metadata (ID, Name, Priority)
// and can be retrieved by ID or as a sorted list.
//
// # Global Registry
//
// The global registry provides thread-safe registration and retrieval of providers:
//
//	// Register a provider in init()
//	func init() {
//	    err := providers.Register(Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	    if err != nil {
//	        panic(err)
//	    }
//	}
//
//	// Get a specific provider
//	reg := providers.Get("claude-code")
//
//	// Get all providers sorted by priority
//	allRegs := providers.All()
//
// # Instance-Based Registry
//
// For testing or isolated use cases, create a Registry instance:
//
//	r := providers.NewRegistry()
//	err := r.Register(Registration{...})
package providers

import (
	"fmt"
	"maps"
	"sort"
	"sync"
)

// registry is the global provider registry.
var (
	registry     = make(map[string]Registration)
	registryLock sync.RWMutex
)

// Register adds a provider registration to the global registry.
//
// This is typically called from init() in each provider file.
// Returns an error if:
//   - The registration is invalid (empty ID, Name, or nil Provider)
//   - A provider with the same ID is already registered
//
// Example:
//
//	func init() {
//	    err := providers.Register(Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	    if err != nil {
//	        panic(err)
//	    }
//	}
func Register(r Registration) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("invalid registration: %w", err)
	}

	registryLock.Lock()
	defer registryLock.Unlock()

	if _, exists := registry[r.ID]; exists {
		return fmt.Errorf("provider %q: %w", r.ID, ErrDuplicateID)
	}

	registry[r.ID] = r

	return nil
}

// Get retrieves a registration by its ID from the global registry.
//
// Returns nil if no provider with the given ID is registered.
func Get(id string) *Registration {
	registryLock.RLock()
	defer registryLock.RUnlock()

	if r, exists := registry[id]; exists {
		return &r
	}

	return nil
}

// All returns all registrations from the global registry sorted by priority.
//
// Lower priority values appear first in the returned slice.
func All() []Registration {
	registryLock.RLock()
	defer registryLock.RUnlock()

	registrations := make(
		[]Registration,
		0,
		len(registry),
	)
	for _, r := range registry {
		registrations = append(registrations, r)
	}

	sort.Slice(registrations, func(i, j int) bool {
		return registrations[i].Priority < registrations[j].Priority
	})

	return registrations
}

// IDs returns all registered provider IDs from the global registry
// sorted by priority.
//
// Lower priority values appear first in the returned slice.
func IDs() []string {
	registrations := All()
	ids := make([]string, len(registrations))
	for i, r := range registrations {
		ids[i] = r.ID
	}

	return ids
}

// Count returns the number of providers in the global registry.
func Count() int {
	registryLock.RLock()
	defer registryLock.RUnlock()

	return len(registry)
}

// Reset clears the global registry.
//
// This function should only be used in tests.
func Reset() {
	registryLock.Lock()
	defer registryLock.Unlock()

	registry = make(map[string]Registration)
}

// Registry provides an instance-based registry for cases where
// global state is not desired (e.g., testing).
//
// Unlike the global registry functions, Registry instances are
// not thread-safe. If you need thread-safety, use the global
// registry functions (Register, Get, All, etc.).
type Registry struct {
	registrations map[string]Registration
}

// NewRegistry creates a new empty registry.
//
// Example:
//
//	r := providers.NewRegistry()
//	err := r.Register(Registration{
//	    ID:       "test-provider",
//	    Name:     "Test Provider",
//	    Priority: 1,
//	    Provider: &TestProvider{},
//	})
func NewRegistry() *Registry {
	return &Registry{
		registrations: make(map[string]Registration),
	}
}

// NewRegistryFromGlobal creates a registry populated with all globally
// registered providers.
//
// This is useful for testing when you want to start with the global
// providers and potentially add or override some.
func NewRegistryFromGlobal() *Registry {
	r := NewRegistry()

	registryLock.RLock()
	defer registryLock.RUnlock()

	maps.Copy(r.registrations, registry)

	return r
}

// Register adds a provider registration to this registry.
//
// Returns an error if:
//   - The registration is invalid (empty ID, Name, or nil Provider)
//   - A provider with the same ID is already registered
func (r *Registry) Register(reg Registration) error {
	if err := reg.Validate(); err != nil {
		return fmt.Errorf("invalid registration: %w", err)
	}

	if _, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf("provider %q: %w", reg.ID, ErrDuplicateID)
	}

	r.registrations[reg.ID] = reg

	return nil
}

// Get retrieves a registration by its ID.
//
// Returns nil if no provider with the given ID is registered.
func (r *Registry) Get(id string) *Registration {
	if reg, exists := r.registrations[id]; exists {
		return &reg
	}

	return nil
}

// All returns all registrations in this registry sorted by priority.
//
// Lower priority values appear first in the returned slice.
func (r *Registry) All() []Registration {
	registrations := make(
		[]Registration,
		0,
		len(r.registrations),
	)
	for _, reg := range r.registrations {
		registrations = append(registrations, reg)
	}

	sort.Slice(registrations, func(i, j int) bool {
		return registrations[i].Priority < registrations[j].Priority
	})

	return registrations
}

// IDs returns all provider IDs in this registry sorted by priority.
//
// Lower priority values appear first in the returned slice.
func (r *Registry) IDs() []string {
	registrations := r.All()
	ids := make([]string, len(registrations))
	for i, reg := range registrations {
		ids[i] = reg.ID
	}

	return ids
}

// Count returns the number of providers in this registry.
func (r *Registry) Count() int {
	return len(r.registrations)
}
