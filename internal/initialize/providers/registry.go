package providers

import (
	"errors"
	"fmt"
	"sort"
)

const (
	// Priority values for built-in providers.
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

// registry is the global provider registry mapping ID to Registration.
var registry = make(map[string]Registration)

// RegisterProvider registers a provider with metadata. Returns an error if
// the registration is invalid or if a provider with the same ID is already
// registered.
func RegisterProvider(reg Registration) error {
	if reg.ID == "" {
		return errors.New("provider ID is required")
	}
	if reg.Provider == nil {
		return errors.New("provider implementation is required")
	}
	if _, exists := registry[reg.ID]; exists {
		return fmt.Errorf("provider %q already registered", reg.ID)
	}
	registry[reg.ID] = reg

	return nil
}

// RegisteredProviders returns all registered providers sorted by priority
// (lower first).
func RegisteredProviders() []Registration {
	result := make([]Registration, 0, len(registry))
	for _, reg := range registry {
		result = append(result, reg)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

// Get retrieves a provider registration by ID. Returns the registration
// and true if found, or an empty registration and false if not found.
func Get(id string) (Registration, bool) {
	reg, ok := registry[id]

	return reg, ok
}

// Count returns the number of registered providers.
func Count() int {
	return len(registry)
}

// IDs returns all registered provider IDs in priority order.
func IDs() []string {
	registrations := RegisteredProviders()
	ids := make([]string, len(registrations))
	for i, reg := range registrations {
		ids[i] = reg.ID
	}

	return ids
}

// ResetRegistry clears all registered providers. Only for testing.
func ResetRegistry() {
	registry = make(map[string]Registration)
}

// RegisterAllProviders registers all built-in providers.
// Called once at application startup (e.g., from cmd/init.go).
// Returns error if any registration fails.
func RegisterAllProviders() error {
	providers := []Registration{
		{
			ID: "claude-code", Name: "Claude Code",
			Priority: priorityClaudeCode, Provider: &ClaudeProvider{},
		},
		{
			ID: "gemini", Name: "Gemini CLI",
			Priority: priorityGemini, Provider: &GeminiProvider{},
		},
		{
			ID: "costrict", Name: "Costrict",
			Priority: priorityCostrict, Provider: &CostrictProvider{},
		},
		{
			ID: "qoder", Name: "Qoder",
			Priority: priorityQoder, Provider: &QoderProvider{},
		},
		{
			ID: "qwen", Name: "Qwen Code",
			Priority: priorityQwen, Provider: &QwenProvider{},
		},
		{
			ID: "antigravity", Name: "Antigravity",
			Priority: priorityAntigravity, Provider: &AntigravityProvider{},
		},
		{
			ID: "cline", Name: "Cline",
			Priority: priorityCline, Provider: &ClineProvider{},
		},
		{
			ID: "cursor", Name: "Cursor",
			Priority: priorityCursor, Provider: &CursorProvider{},
		},
		{
			ID: "codex", Name: "Codex CLI",
			Priority: priorityCodex, Provider: &CodexProvider{},
		},
		{
			ID: "aider", Name: "Aider",
			Priority: priorityAider, Provider: &AiderProvider{},
		},
		{
			ID: "windsurf", Name: "Windsurf",
			Priority: priorityWindsurf, Provider: &WindsurfProvider{},
		},
		{
			ID: "kilocode", Name: "Kilocode",
			Priority: priorityKilocode, Provider: &KilocodeProvider{},
		},
		{
			ID: "continue", Name: "Continue",
			Priority: priorityContinue, Provider: &ContinueProvider{},
		},
		{
			ID: "crush", Name: "Crush",
			Priority: priorityCrush, Provider: &CrushProvider{},
		},
		{
			ID: "opencode", Name: "OpenCode",
			Priority: priorityOpencode, Provider: &OpencodeProvider{},
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			// Successfully registered providers remain (no rollback)
			return fmt.Errorf(
				"failed to register %s provider: %w",
				reg.ID, err,
			)
		}
	}

	return nil
}
