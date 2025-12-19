// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the new registry (RegistryV2) that works with the
// Registration struct and ProviderV2 interface. Unlike the legacy registry
// that accepted Provider interfaces directly, this registry accepts
// Registration structs which contain metadata (ID, Name, Priority) along
// with the ProviderV2 implementation.
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
//	err := providers.RegisterV2(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//
//	// Using instance registry
//	registry := providers.NewRegistryV2()
//	err := registry.Register(providers.Registration{...})
//
// See also:
//   - registration.go: Registration struct definition and validation
//   - provider_new.go: ProviderV2 interface
//   - registry.go: Legacy registry (for backwards compatibility)
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

// registryV2 is the global provider registry for the new Registration-based API.
var (
	globalRegistryV2     = NewRegistryV2()
	globalRegistryV2Lock sync.RWMutex
)

// RegistryV2 provides an instance-based registry for providers using the new
// Registration struct. It is thread-safe and supports validation, duplicate
// rejection, and priority-sorted retrieval.
//
// Use NewRegistryV2() to create a new instance, or use the global convenience
// functions (RegisterV2, GetV2, AllV2, etc.) for the default global registry.
type RegistryV2 struct {
	mu            sync.RWMutex
	registrations map[string]Registration
}

// NewRegistryV2 creates a new empty registry.
func NewRegistryV2() *RegistryV2 {
	return &RegistryV2{
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
func (r *RegistryV2) Register(reg Registration) error {
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
func (r *RegistryV2) Get(id string) (*Registration, bool) {
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
func (r *RegistryV2) All() []Registration {
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
func (r *RegistryV2) IDs() []string {
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
func (r *RegistryV2) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.registrations)
}

// Reset clears all registrations from this registry.
//
// This is primarily useful for testing. This method is thread-safe.
func (r *RegistryV2) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = make(map[string]Registration)
}

// -----------------------------------------------------------------------------
// Global Convenience Functions
// -----------------------------------------------------------------------------

// RegisterV2 adds a provider registration to the global registry.
//
// This is a convenience function that calls Register on the global RegistryV2.
// The registration is validated before being accepted. Returns an error if:
//   - The registration fails validation (invalid ID, Name, Priority, nil)
//   - A provider with the same ID is already registered
//
// This function is thread-safe.
func RegisterV2(reg Registration) error {
	globalRegistryV2Lock.Lock()
	defer globalRegistryV2Lock.Unlock()

	return globalRegistryV2.Register(reg)
}

// GetV2 retrieves a registration by its ID from the global registry.
//
// This is a convenience function that calls Get on the global RegistryV2.
// Returns the registration and true if found, or nil/false if not found.
// This function is thread-safe.
func GetV2(id string) (*Registration, bool) {
	globalRegistryV2Lock.RLock()
	defer globalRegistryV2Lock.RUnlock()

	return globalRegistryV2.Get(id)
}

// AllV2 returns all registrations from the global registry sorted by priority.
//
// This is a convenience function that calls All on the global RegistryV2.
// Returns registrations sorted by priority (lower priority number first).
// This function is thread-safe.
func AllV2() []Registration {
	globalRegistryV2Lock.RLock()
	defer globalRegistryV2Lock.RUnlock()

	return globalRegistryV2.All()
}

// IDsV2 returns all registered provider IDs from the global registry sorted by priority.
//
// This is a convenience function that calls IDs on the global RegistryV2.
// This function is thread-safe.
func IDsV2() []string {
	globalRegistryV2Lock.RLock()
	defer globalRegistryV2Lock.RUnlock()

	return globalRegistryV2.IDs()
}

// CountV2 returns the number of providers in the global registry.
//
// This is a convenience function that calls Count on the global RegistryV2.
// This function is thread-safe.
func CountV2() int {
	globalRegistryV2Lock.RLock()
	defer globalRegistryV2Lock.RUnlock()

	return globalRegistryV2.Count()
}

// ResetV2 clears the global RegistryV2. Only use in tests.
//
// This function is thread-safe.
func ResetV2() {
	globalRegistryV2Lock.Lock()
	defer globalRegistryV2Lock.Unlock()

	globalRegistryV2.Reset()
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
