// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the provider registry that works with the Registration
// struct and Provider interface. The registry accepts Registration structs
// which contain metadata (ID, Name, Priority) along with the Provider
// implementation.
//
// The registry is thread-safe using sync.RWMutex and provides both a global
// default registry instance and the ability to create independent registry
// instances for testing or other purposes.
//
// Key features:
//   - Thread-safe registration and retrieval using sync.RWMutex
//   - Validation of registrations before acceptance
//   - Duplicate ID rejection with clear error messages
//   - Priority-sorted retrieval (lower priority number = earlier in list)
//   - Both global convenience functions and instance methods
//
// Example usage:
//
//	// Using global registry
//	err := providers.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//
//	// Using instance registry
//	registry := providers.NewRegistry()
//	err := registry.Register(providers.Registration{...})
//
// See also:
//   - registration.go: Registration struct definition and validation
//   - provider.go: Provider interface
//
//nolint:revive // line-length-limit - registry documentation
package providers

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/spf13/afero"
)

// globalRegistry is the global provider registry for the Registration-based API.
var (
	globalRegistry     = NewRegistry()
	globalRegistryLock sync.RWMutex
)

// Registry provides an instance-based registry for providers using the
// Registration struct. It is thread-safe and supports validation, duplicate
// rejection, and priority-sorted retrieval.
//
// Use NewRegistry() to create a new instance, or use the global convenience
// functions (Register, Get, All, etc.) for the default global registry.
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

// Register adds a provider registration to this registry.
//
// The registration is validated before being accepted. Returns an error if:
//   - The registration fails validation (invalid ID, Name, Priority, nil Provider)
//   - A provider with the same ID is already registered
//
// This method is thread-safe.
func (r *Registry) Register(reg Registration) error {
	// Validate the registration first
	if err := reg.Validate(); err != nil {
		return fmt.Errorf("cannot register provider: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate ID
	if existing, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf(
			"cannot register provider %q: ID %q already registered (%q)",
			reg.Name, reg.ID, existing.Name,
		)
	}

	r.registrations[reg.ID] = reg

	return nil
}

// Get retrieves a registration by its ID.
//
// Returns the registration and true if found, or a zero Registration and false
// if not found. This method is thread-safe.
func (r *Registry) Get(id string) (*Registration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, exists := r.registrations[id]
	if !exists {
		return nil, false
	}

	return &reg, true
}

// All returns all registrations sorted by priority (lower priority number first).
//
// This maintains backwards-compatible behavior with the legacy registry.
// This method is thread-safe.
func (r *Registry) All() []Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	registrations := make([]Registration, 0, len(r.registrations))
	for _, reg := range r.registrations {
		registrations = append(registrations, reg)
	}

	// Sort by priority (ascending), then by ID (ascending) for deterministic order
	sort.Slice(registrations, func(i, j int) bool {
		if registrations[i].Priority != registrations[j].Priority {
			return registrations[i].Priority < registrations[j].Priority
		}

		return registrations[i].ID < registrations[j].ID
	})

	return registrations
}

// IDs returns all registered provider IDs sorted by priority.
//
// This is a convenience method that returns just the IDs from All().
// This method is thread-safe.
func (r *Registry) IDs() []string {
	registrations := r.All()
	ids := make([]string, len(registrations))
	for i, reg := range registrations {
		ids[i] = reg.ID
	}

	return ids
}

// Count returns the number of registered providers.
//
// This method is thread-safe.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.registrations)
}

// Reset clears all registrations from this registry.
//
// This is primarily useful for testing. This method is thread-safe.
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = make(map[string]Registration)
}

// -----------------------------------------------------------------------------
// Global Convenience Functions
// -----------------------------------------------------------------------------

// Register adds a provider registration to the global registry.
//
// This is a convenience function that calls Register on the global Registry.
// The registration is validated before being accepted. Returns an error if:
//   - The registration fails validation (invalid ID, Name, Priority, nil)
//   - A provider with the same ID is already registered
//
// This function is thread-safe.
func Register(reg Registration) error {
	globalRegistryLock.Lock()
	defer globalRegistryLock.Unlock()

	return globalRegistry.Register(reg)
}

// Get retrieves a registration by its ID from the global registry.
//
// This is a convenience function that calls Get on the global Registry.
// Returns the registration and true if found, or nil/false if not found.
// This function is thread-safe.
func Get(id string) (*Registration, bool) {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.Get(id)
}

// All returns all registrations from the global registry sorted by priority.
//
// This is a convenience function that calls All on the global Registry.
// Returns registrations sorted by priority (lower priority number first).
// This function is thread-safe.
func All() []Registration {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.All()
}

// IDs returns all registered provider IDs from the global registry sorted by priority.
//
// This is a convenience function that calls IDs on the global Registry.
// This function is thread-safe.
func IDs() []string {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.IDs()
}

// Count returns the number of providers in the global registry.
//
// This is a convenience function that calls Count on the global Registry.
// This function is thread-safe.
func Count() int {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.Count()
}

// Reset clears the global Registry. Only use in tests.
//
// This function is thread-safe.
func Reset() {
	globalRegistryLock.Lock()
	defer globalRegistryLock.Unlock()

	globalRegistry.Reset()
}

// -----------------------------------------------------------------------------
// Helper Functions for Provider Status
// -----------------------------------------------------------------------------

// IsProviderConfigured checks if a provider's initializers are all set up.
//
// This function creates filesystems for the project path and home directory,
// gets the initializers from the provider, and checks if each initializer
// reports IsSetup() == true.
//
// Returns true if all initializers are set up, false otherwise.
// Returns false if the project path doesn't exist or provider has no initializers.
func IsProviderConfigured(reg Registration, projectPath string) bool {
	if reg.Provider == nil {
		return false
	}

	// Create filesystems
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectPath)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	globalFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

	// Create config with default spectr dir
	cfg := NewConfig("spectr")

	// Get initializers from provider
	ctx := context.Background()
	initializers := reg.Provider.Initializers(ctx)

	if len(initializers) == 0 {
		return false
	}

	// Check if all initializers are set up
	for _, init := range initializers {
		if init == nil {
			continue
		}

		// Select appropriate filesystem based on IsGlobal()
		var fs afero.Fs
		if init.IsGlobal() {
			fs = globalFs
		} else {
			fs = projectFs
		}

		if !init.IsSetup(fs, cfg) {
			return false
		}
	}

	return true
}
