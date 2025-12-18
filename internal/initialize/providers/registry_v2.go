// Package providers defines the core interfaces for the provider architecture.
//
// This file contains the new registry (RegistryV2) that uses the Registration
// struct for metadata. This replaces the old registry which stored Provider
// objects directly.
//
// # Migration Note
//
// During the migration period, this is named RegistryV2 to avoid conflicts
// with the existing Registry. After migration is complete (tasks 7.x), the
// old registry will be removed and RegistryV2 will be renamed to Registry.
package providers

import (
	"fmt"
	"maps"
	"sort"
	"sync"
)

// registryV2 is the global provider registry for the new architecture.
var (
	registryV2     = make(map[string]Registration)
	registryLockV2 sync.RWMutex
)

// RegisterV2 adds a provider registration to the global registry.
//
// This is typically called from init() in each provider file.
// Returns an error if:
//   - The registration is invalid (empty ID, Name, or nil Provider)
//   - A provider with the same ID is already registered
//
// Example:
//
//	func init() {
//	    err := providers.RegisterV2(Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	    if err != nil {
//	        panic(err)
//	    }
//	}
func RegisterV2(r Registration) error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("invalid registration: %w", err)
	}

	registryLockV2.Lock()
	defer registryLockV2.Unlock()

	if _, exists := registryV2[r.ID]; exists {
		return fmt.Errorf("provider %q: %w", r.ID, ErrDuplicateID)
	}

	registryV2[r.ID] = r

	return nil
}

// GetV2 retrieves a registration by its ID from the global registry.
//
// Returns nil if no provider with the given ID is registered.
func GetV2(id string) *Registration {
	registryLockV2.RLock()
	defer registryLockV2.RUnlock()

	if r, exists := registryV2[id]; exists {
		return &r
	}

	return nil
}

// AllV2 returns all registrations from the global registry sorted by priority.
//
// Lower priority values appear first in the returned slice.
// This maintains backwards-compatible behavior with the old registry.
func AllV2() []Registration {
	registryLockV2.RLock()
	defer registryLockV2.RUnlock()

	registrations := make(
		[]Registration,
		0,
		len(registryV2),
	)
	for _, r := range registryV2 {
		registrations = append(registrations, r)
	}

	sort.Slice(registrations, func(i, j int) bool {
		return registrations[i].Priority < registrations[j].Priority
	})

	return registrations
}

// IDsV2 returns all registered provider IDs from the global registry
// sorted by priority.
//
// Lower priority values appear first in the returned slice.
func IDsV2() []string {
	registrations := AllV2()
	ids := make([]string, len(registrations))
	for i, r := range registrations {
		ids[i] = r.ID
	}

	return ids
}

// CountV2 returns the number of providers in the global registry.
func CountV2() int {
	registryLockV2.RLock()
	defer registryLockV2.RUnlock()

	return len(registryV2)
}

// ResetV2 clears the global registry.
//
// This function should only be used in tests.
func ResetV2() {
	registryLockV2.Lock()
	defer registryLockV2.Unlock()

	registryV2 = make(map[string]Registration)
}

// RegistryV2 provides an instance-based registry for cases where
// global state is not desired (e.g., testing).
//
// Unlike the global registry functions, RegistryV2 instances are
// not thread-safe. If you need thread-safety, use the global
// registry functions (RegisterV2, GetV2, AllV2, etc.).
//
// # Migration Note
//
// This type will be renamed to Registry after the old Registry type
// is removed in task 7.x.
type RegistryV2 struct {
	registrations map[string]Registration
}

// NewRegistryV2 creates a new empty registry.
//
// Example:
//
//	r := providers.NewRegistryV2()
//	err := r.Register(Registration{
//	    ID:       "test-provider",
//	    Name:     "Test Provider",
//	    Priority: 1,
//	    Provider: &TestProvider{},
//	})
func NewRegistryV2() *RegistryV2 {
	return &RegistryV2{
		registrations: make(map[string]Registration),
	}
}

// NewRegistryV2FromGlobal creates a registry populated with all globally
// registered providers.
//
// This is useful for testing when you want to start with the global
// providers and potentially add or override some.
func NewRegistryV2FromGlobal() *RegistryV2 {
	r := NewRegistryV2()

	registryLockV2.RLock()
	defer registryLockV2.RUnlock()

	maps.Copy(r.registrations, registryV2)

	return r
}

// Register adds a provider registration to this registry.
//
// Returns an error if:
//   - The registration is invalid (empty ID, Name, or nil Provider)
//   - A provider with the same ID is already registered
func (r *RegistryV2) Register(reg Registration) error {
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
func (r *RegistryV2) Get(id string) *Registration {
	if reg, exists := r.registrations[id]; exists {
		return &reg
	}

	return nil
}

// All returns all registrations in this registry sorted by priority.
//
// Lower priority values appear first in the returned slice.
func (r *RegistryV2) All() []Registration {
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
func (r *RegistryV2) IDs() []string {
	registrations := r.All()
	ids := make([]string, len(registrations))
	for i, reg := range registrations {
		ids[i] = reg.ID
	}

	return ids
}

// Count returns the number of providers in this registry.
func (r *RegistryV2) Count() int {
	return len(r.registrations)
}
