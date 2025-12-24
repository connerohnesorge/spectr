package providers

import (
	"fmt"
	"sort"
	"sync"
)

// ProviderRegistry manages registration of providers using the new
// composable architecture.
//
// This registry stores Registration structs which contain provider
// metadata (ID, Name, Priority) and the Provider implementation.
// It provides thread-safe access to registered providers and supports
// priority-based sorting for consistent ordering across the codebase.
//
// Usage:
//
//	registry := NewProviderRegistry()
//	registry.RegisterProvider(Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//
// Thread Safety:
//
//	All methods are thread-safe and can be called concurrently from
//	multiple goroutines.
type ProviderRegistry struct {
	registrations map[string]Registration
	mu            sync.RWMutex
}

// NewProviderRegistry creates a new empty provider registry.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		registrations: make(
			map[string]Registration,
		),
	}
}

// RegisterProvider adds a registration to the registry.
//
// Returns an error if a provider with the same ID is already registered.
// This ensures uniqueness and prevents accidental duplicate registrations.
//
// Example:
//
//	err := registry.RegisterProvider(Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//	if err != nil {
//	    // Handle duplicate registration error
//	}
func (r *ProviderRegistry) RegisterProvider(
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

// Get retrieves a registration by ID.
//
// Returns the registration and true if found, or an empty
// registration and false if not found.
//
// Example:
//
//	reg, ok := registry.Get("claude-code")
//	if !ok {
//	    // Provider not found
//	}
func (r *ProviderRegistry) Get(
	id string,
) (Registration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, ok := r.registrations[id]

	return reg, ok
}

// All returns all registrations sorted by priority.
//
// Lower priority values are returned first (e.g., Priority: 1 comes
// before Priority: 2).
// This ensures consistent ordering when displaying providers to users.
//
// The returned slice is a copy, so modifications to it won't affect
// the registry.
func (r *ProviderRegistry) All() []Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	registrations := make(
		[]Registration,
		0,
		len(r.registrations),
	)
	for _, reg := range r.registrations {
		registrations = append(registrations, reg)
	}

	sort.Slice(
		registrations,
		func(i, j int) bool {
			return registrations[i].Priority < registrations[j].Priority
		},
	)

	return registrations
}

// IDs returns all provider IDs sorted by priority.
//
// This is a convenience method that returns just the IDs from All().
// Lower priority values are returned first.
func (r *ProviderRegistry) IDs() []string {
	registrations := r.All()
	ids := make([]string, len(registrations))
	for i, reg := range registrations {
		ids[i] = reg.ID
	}

	return ids
}

// Count returns the number of registered providers.
func (r *ProviderRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.registrations)
}

// Global registry instance and convenience functions

var globalProviderRegistry = NewProviderRegistry()

// RegisterProvider adds a registration to the global registry.
//
// Returns an error if a provider with the same ID is already registered.
// This is typically called from init() functions in provider files.
//
// Example:
//
//	func init() {
//	    err := providers.RegisterProvider(providers.Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	    if err != nil {
//	        panic(err) // init() failures are typically fatal
//	    }
//	}
func RegisterProvider(reg Registration) error {
	return globalProviderRegistry.RegisterProvider(
		reg,
	)
}

// GetProvider retrieves a registration by ID from the global registry.
//
// Returns the registration and true if found, or an empty
// registration and false if not found.
func GetProvider(id string) (Registration, bool) {
	return globalProviderRegistry.Get(id)
}

// AllProviders returns all registrations from the global registry
// sorted by priority.
//
// Lower priority values are returned first.
func AllProviders() []Registration {
	return globalProviderRegistry.All()
}

// ProviderIDs returns all provider IDs from the global registry
// sorted by priority.
//
// Lower priority values are returned first.
func ProviderIDs() []string {
	return globalProviderRegistry.IDs()
}

// ProviderCount returns the number of registered providers in the
// global registry.
func ProviderCount() int {
	return globalProviderRegistry.Count()
}

// ResetProviders clears the global registry. Only use in tests.
//
// This creates a completely new registry instance, which is useful
// for test isolation.
func ResetProviders() {
	globalProviderRegistry = NewProviderRegistry()
}
