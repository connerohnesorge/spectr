package providers

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// registry is the package-level map storing all registered providers.
var registry = make(map[string]Registration)

// registryLock protects concurrent access to the registry.
var registryLock sync.RWMutex

// RegisterProvider registers a provider and returns an error if registration fails.
// Validation:
//   - ID must be non-empty
//   - Provider must be non-nil
//   - ID must not already be registered (returns error, does not panic)
func RegisterProvider(reg Registration) error {
	if reg.ID == "" {
		return errors.New("provider ID is required")
	}
	if reg.Provider == nil {
		return errors.New("provider implementation is required")
	}

	registryLock.Lock()
	defer registryLock.Unlock()

	if _, exists := registry[reg.ID]; exists {
		return fmt.Errorf("provider %q already registered", reg.ID)
	}

	registry[reg.ID] = reg

	return nil
}

// RegisteredProviders returns all registered providers sorted by priority (lower first).
func RegisteredProviders() []Registration {
	registryLock.RLock()
	defer registryLock.RUnlock()

	result := make([]Registration, 0, len(registry))
	for _, reg := range registry {
		result = append(result, reg)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

// Get retrieves a provider registration by its ID.
// Returns the Registration and true if found, or zero Registration and false if not found.
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

// ResetRegistry clears the global registry. Only use in tests.
func ResetRegistry() {
	registryLock.Lock()
	defer registryLock.Unlock()

	registry = make(map[string]Registration)
}

// Provider priority constants.
const (
	priorityClaudeCode  = 1
	priorityGemini      = 2
	priorityCostrict    = 3
	priorityQoder       = 4
	priorityQwen        = 5
	priorityAntigravity = 6
	priorityCline       = 7
	priorityCursor      = 8
	priorityCodex       = 9
	priorityAider       = 10
	priorityWindsurf    = 11
	priorityKilocode    = 12
	priorityContinue    = 13
	priorityCrush       = 14
	priorityOpencode    = 15
)

// RegisterAllProviders registers all built-in providers.
// Called once at application startup (e.g., from main() or cmd/root.go).
// Returns error if any registration fails.
//
// Provider priorities (1-15):
//   - Priority 1: claude-code (Claude Code)
//   - Priority 2: gemini (Gemini CLI)
//   - Priority 3: costrict (CoStrict)
//   - Priority 4: qoder (Qoder)
//   - Priority 5: qwen (Qwen Code)
//   - Priority 6: antigravity (Antigravity)
//   - Priority 7: cline (Cline)
//   - Priority 8: cursor (Cursor)
//   - Priority 9: codex (Codex CLI)
//   - Priority 10: aider (Aider)
//   - Priority 11: windsurf (Windsurf)
//   - Priority 12: kilocode (Kilocode)
//   - Priority 13: continue (Continue)
//   - Priority 14: crush (Crush)
//   - Priority 15: opencode (OpenCode)
func RegisterAllProviders() error {
	providers := []Registration{
		{
			ID:       "claude-code",
			Name:     "Claude Code",
			Priority: priorityClaudeCode,
			Provider: &ClaudeProvider{},
		},
		{ID: "gemini", Name: "Gemini CLI", Priority: priorityGemini, Provider: &GeminiProvider{}},
		{
			ID:       "costrict",
			Name:     "CoStrict",
			Priority: priorityCostrict,
			Provider: &CostrictProvider{},
		},
		{ID: "qoder", Name: "Qoder", Priority: priorityQoder, Provider: &QoderProvider{}},
		{ID: "qwen", Name: "Qwen Code", Priority: priorityQwen, Provider: &QwenProvider{}},
		{
			ID:       "antigravity",
			Name:     "Antigravity",
			Priority: priorityAntigravity,
			Provider: &AntigravityProvider{},
		},
		{ID: "cline", Name: "Cline", Priority: priorityCline, Provider: &ClineProvider{}},
		{ID: "cursor", Name: "Cursor", Priority: priorityCursor, Provider: &CursorProvider{}},
		{ID: "codex", Name: "Codex CLI", Priority: priorityCodex, Provider: &CodexProvider{}},
		{ID: "aider", Name: "Aider", Priority: priorityAider, Provider: &AiderProvider{}},
		{
			ID:       "windsurf",
			Name:     "Windsurf",
			Priority: priorityWindsurf,
			Provider: &WindsurfProvider{},
		},
		{
			ID:       "kilocode",
			Name:     "Kilocode",
			Priority: priorityKilocode,
			Provider: &KilocodeProvider{},
		},
		{
			ID:       "continue",
			Name:     "Continue",
			Priority: priorityContinue,
			Provider: &ContinueProvider{},
		},
		{ID: "crush", Name: "Crush", Priority: priorityCrush, Provider: &CrushProvider{}},
		{
			ID:       "opencode",
			Name:     "OpenCode",
			Priority: priorityOpencode,
			Provider: &OpencodeProvider{},
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			// Note: Successfully registered providers remain registered (no rollback)
			return fmt.Errorf("failed to register %s provider: %w", reg.ID, err)
		}
	}

	return nil
}
