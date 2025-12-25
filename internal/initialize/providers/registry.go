package providers

import (
	"fmt"
	"sort"
	"sync"
)

// Global registry instance
var (
	globalRegistry     = NewRegistry()
	globalRegistryLock sync.RWMutex
)

// Registry stores provider registrations and provides access methods.
// Each registration contains a provider implementation and its metadata
// (ID, Name, Priority).
type Registry struct {
	registrations map[string]Registration
	mu            sync.RWMutex
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		registrations: make(
			map[string]Registration,
		),
	}
}

// RegisterProvider registers a provider with its metadata.
// Returns an error if:
// - The registration is invalid (missing ID or Provider)
// - A provider with the same ID is already registered
//
// This method is thread-safe.
func RegisterProvider(reg Registration) error {
	globalRegistryLock.Lock()
	defer globalRegistryLock.Unlock()

	return globalRegistry.registerProvider(reg)
}

// registerProvider is the internal implementation.
// Not thread-safe, caller must lock.
func (r *Registry) registerProvider(
	reg Registration,
) error {
	// Validate the registration
	if err := reg.Validate(); err != nil {
		return err
	}

	// Check for duplicate ID
	if _, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf(
			"provider %q already registered",
			reg.ID,
		)
	}

	r.registrations[reg.ID] = reg

	return nil
}

// GetRegistration retrieves a provider registration by ID.
// Returns nil if the provider is not found.
//
// This is the NEW API - use this in new code.
func GetRegistration(id string) *Registration {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.GetRegistration(id)
}

// GetRegistration retrieves a provider registration by ID from this registry.
// Returns nil if the provider is not found.
func (r *Registry) GetRegistration(
	id string,
) *Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if reg, exists := r.registrations[id]; exists {
		return &reg
	}

	return nil
}

// AllRegistrations returns all provider registrations sorted by priority.
//
// This is the NEW API - use this in new code.
func AllRegistrations() []Registration {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.AllRegistrations()
}

// AllRegistrations returns all provider registrations sorted by priority.
func (r *Registry) AllRegistrations() []Registration {
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

	// Sort by priority (lower = higher priority)
	sort.Slice(
		registrations,
		func(i, j int) bool {
			return registrations[i].Priority < registrations[j].Priority
		},
	)

	return registrations
}

// Get retrieves a legacy provider by ID.
//
// Deprecated: Use GetRegistration() in new code.
func Get(_ string) Provider {
	// Providers now use Registration-based API.
	return nil
}

// All returns all legacy providers.
//
// Deprecated: Use AllRegistrations() in new code.
func All() []Provider {
	// Providers now use Registration-based API.
	return nil
}

// IDs returns all provider IDs sorted by priority.
//
// For new code, prefer using AllRegistrations().
func IDs() []string {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.IDs()
}

// Count returns the number of registered providers.
//
// This is a global function that operates on the global registry.
func Count() int {
	globalRegistryLock.RLock()
	defer globalRegistryLock.RUnlock()

	return globalRegistry.Count()
}

// Count returns the number of registered providers in this registry.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.registrations)
}

// IDs returns all provider IDs from this registry sorted by priority.
func (r *Registry) IDs() []string {
	registrations := r.AllRegistrations()
	ids := make([]string, len(registrations))
	for i, reg := range registrations {
		ids[i] = reg.ID
	}

	return ids
}

// Get retrieves a provider registration by ID from this registry.
//
// Deprecated: Use GetRegistration() for clarity.
func (r *Registry) Get(id string) *Registration {
	return r.GetRegistration(id)
}

// All returns all provider registrations from this registry.
//
// Deprecated: Use AllRegistrations() for clarity.
func (r *Registry) All() []Registration {
	return r.AllRegistrations()
}

// Reset clears the global registry. Only use in tests.
func Reset() {
	globalRegistryLock.Lock()
	defer globalRegistryLock.Unlock()

	globalRegistry = NewRegistry()
}

// RegisterAllProviders registers all built-in providers.
// This should be called once at application startup.
//
// Returns an error if any registration fails.
func RegisterAllProviders() error {
	providers := []Registration{
		{
			ID:       "claude-code",
			Name:     "Claude Code",
			Priority: PriorityClaudeCode,
			Provider: &ClaudeProvider{},
		},
		{
			ID:       "gemini",
			Name:     "Gemini CLI",
			Priority: PriorityGemini,
			Provider: &GeminiProvider{},
		},
		{
			ID:       "costrict",
			Name:     "CoStrict",
			Priority: PriorityCostrict,
			Provider: &CostrictProvider{},
		},
		{
			ID:       "qoder",
			Name:     "Qoder",
			Priority: PriorityQoder,
			Provider: &QoderProvider{},
		},
		{
			ID:       "qwen",
			Name:     "Qwen Code",
			Priority: PriorityQwen,
			Provider: &QwenProvider{},
		},
		{
			ID:       "antigravity",
			Name:     "Antigravity",
			Priority: PriorityAntigravity,
			Provider: &AntigravityProvider{},
		},
		{
			ID:       "cline",
			Name:     "Cline",
			Priority: PriorityCline,
			Provider: &ClineProvider{},
		},
		{
			ID:       "cursor",
			Name:     "Cursor",
			Priority: PriorityCursor,
			Provider: &CursorProvider{},
		},
		{
			ID:       "codex",
			Name:     "Codex CLI",
			Priority: PriorityCodex,
			Provider: &CodexProvider{},
		},
		{
			ID:       "opencode",
			Name:     "OpenCode",
			Priority: PriorityOpencode,
			Provider: &OpencodeProvider{},
		},
		{
			ID:       "aider",
			Name:     "Aider",
			Priority: PriorityAider,
			Provider: &AiderProvider{},
		},
		{
			ID:       "windsurf",
			Name:     "Windsurf",
			Priority: PriorityWindsurf,
			Provider: &WindsurfProvider{},
		},
		{
			ID:       "kilocode",
			Name:     "Kilocode",
			Priority: PriorityKilocode,
			Provider: &KilocodeProvider{},
		},
		{
			ID:       "continue",
			Name:     "Continue",
			Priority: PriorityContinue,
			Provider: &ContinueProvider{},
		},
		{
			ID:       "crush",
			Name:     "Crush",
			Priority: PriorityCrush,
			Provider: &CrushProvider{},
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			return fmt.Errorf(
				"failed to register %s provider: %w",
				reg.ID,
				err,
			)
		}
	}

	return nil
}
