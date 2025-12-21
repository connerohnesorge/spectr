package providers

import (
	"fmt"
	"sort"
	"sync"
)

var (
	// DefaultRegistry is the default global registry.
	DefaultRegistry = NewRegistry()
)

// Register adds a provider to the default registry.
func Register(r Registration) {
	if err := DefaultRegistry.Register(r); err != nil {
		panic(err)
	}
}

// Get retrieves a provider from the default registry.
func Get(id string) (Registration, bool) {
	return DefaultRegistry.Get(id)
}

// All returns all providers from the default registry.
func All() []Registration {
	return DefaultRegistry.All()
}

// Registry manages a collection of provider registrations.
type Registry struct {
	registrations map[string]Registration
	mu            sync.RWMutex
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		registrations: make(map[string]Registration),
	}
}

// Register adds a provider registration.
func (r *Registry) Register(reg Registration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf("provider %q already registered", reg.ID)
	}

	r.registrations[reg.ID] = reg
	return nil
}

// Get retrieves a provider registration by ID.
func (r *Registry) Get(id string) (Registration, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, exists := r.registrations[id]
	return reg, exists
}

// All returns all registrations sorted by priority.
func (r *Registry) All() []Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	regs := make([]Registration, 0, len(r.registrations))
	for _, reg := range r.registrations {
		regs = append(regs, reg)
	}

	sort.Slice(regs, func(i, j int) bool {
		if regs[i].Priority != regs[j].Priority {
			return regs[i].Priority < regs[j].Priority
		}
		return regs[i].ID < regs[j].ID
	})

	return regs
}

// IDs returns all registered IDs sorted by priority.
func (r *Registry) IDs() []string {
	regs := r.All()
	ids := make([]string, len(regs))
	for i, reg := range regs {
		ids[i] = reg.ID
	}
	return ids
}

// Count returns the number of registrations.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.registrations)
}
