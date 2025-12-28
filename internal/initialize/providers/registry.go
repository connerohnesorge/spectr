package providers

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// registry is the package-level registry of all providers.
// Maps provider ID to Registration.
var (
	registry     = make(map[string]Registration)
	registryLock sync.RWMutex
)

// RegisterProvider registers a provider with its metadata.
// Returns an error if:
// - ID is empty
// - Provider is nil
// - A provider with the same ID is already registered
func RegisterProvider(reg Registration) error {
	registryLock.Lock()
	defer registryLock.Unlock()

	// Validate ID
	if reg.ID == "" {
		return errors.New("provider ID is required")
	}

	// Validate Provider
	if reg.Provider == nil {
		return errors.New("provider implementation is required")
	}

	// Check for duplicates
	if _, exists := registry[reg.ID]; exists {
		return fmt.Errorf("provider %q already registered", reg.ID)
	}

	// Register
	registry[reg.ID] = reg

	return nil
}

// RegisteredProviders returns all registered providers sorted by
// Priority (lower first).
func RegisteredProviders() []Registration {
	registryLock.RLock()
	defer registryLock.RUnlock()

	result := make([]Registration, 0, len(registry))
	for _, reg := range registry {
		result = append(result, reg)
	}

	// Sort by Priority (lower priority number = higher priority)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

// Get retrieves a provider registration by ID.
// Returns the Registration and true if found, or an empty
// Registration and false if not found.
func Get(id string) (Registration, bool) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	reg, exists := registry[id]

	return reg, exists
}

// Count returns the number of registered providers.
func Count() int {
	registryLock.RLock()
	defer registryLock.RUnlock()

	return len(registry)
}

// Reset clears all registered providers.
// This is primarily for testing purposes.
func Reset() {
	registryLock.Lock()
	defer registryLock.Unlock()

	registry = make(map[string]Registration)
}

// Priority constants for built-in providers (lower = higher priority).
const (
	PriorityClaudeCode  = 1
	PriorityGemini      = 2
	PriorityCostrict    = 3
	PriorityQoder       = 4
	PriorityQwen        = 5
	PriorityAntigravity = 6
	PriorityCline       = 7
	PriorityCursor      = 8
	PriorityCodex       = 9
	PriorityAider       = 10
	PriorityWindsurf    = 11
	PriorityKilocode    = 12
	PriorityContinue    = 13
	PriorityCrush       = 14
	PriorityOpencode    = 15
)

// RegisterAllProviders registers all built-in providers.
// This should be called once at application startup.
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
			Name:     "Costrict",
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
		{
			ID:       "opencode",
			Name:     "OpenCode",
			Priority: PriorityOpencode,
			Provider: &OpencodeProvider{},
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			return fmt.Errorf("failed to register %s provider: %w", reg.ID, err)
		}
	}

	return nil
}
