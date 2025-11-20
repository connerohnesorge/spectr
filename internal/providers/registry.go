// Package providers implements the provider registry and individual
// AI tool providers.
//
// This package serves as the central registry for all AI tool providers
// supported by Spectr. Providers self-register via init() functions,
// enabling automatic discovery without hardcoded switch statements.
package providers

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// Priority constants for providers.
// Config providers use 1-100, slash providers use 101-200.
const (
	// Config provider priorities (1-10 reserved)
	PriorityClaudeCode  = 1
	PriorityCline       = 2
	PriorityCostrict    = 3
	PriorityQoder       = 4
	PriorityCodeBuddy   = 5
	PriorityQwen        = 6
	PriorityAntigravity = 7

	// Slash provider priorities (101-200 reserved)
	PriorityClaudeSlash      = 101
	PriorityClineSlash       = 102
	PriorityKilocodeSlash    = 103
	PriorityQoderSlash       = 104
	PriorityCursorSlash      = 105
	PriorityAiderSlash       = 106
	PriorityContinueSlash    = 107
	PriorityCopilotSlash     = 108
	PriorityMentatSlash      = 109
	PriorityTabnineSlash     = 110
	PrioritySmolSlash        = 111
	PriorityCostrictSlash    = 112
	PriorityCodeBuddySlash   = 113
	PriorityQwenSlash        = 114
	PriorityAntigravitySlash = 115
)

// ProviderType represents the type of provider configuration
type ProviderType string

const (
	// TypeConfig represents providers that create single config files
	TypeConfig ProviderType = "config"
	// TypeSlash represents providers that create slash command files
	TypeSlash ProviderType = "slash"
)

// ProviderMetadata contains metadata about a provider for display
// and organization
type ProviderMetadata struct {
	// ID is the unique identifier for the provider
	// (e.g., "claude-code", "cursor")
	ID string
	// Name is the human-readable name shown to users
	Name string
	// Type indicates whether this is a config or slash command provider
	Type ProviderType
	// Priority determines display order (lower numbers first)
	Priority int
	// FilePaths are the relative paths to files this provider
	// creates/updates. For config providers: single file like ["CLAUDE.md"]
	// For slash providers: multiple files like
	// [".claude/commands/spectr/proposal.md", ...]
	FilePaths []string
	// AutoInstallSlashID is the slash provider ID to auto-install when
	// this config provider is selected. Only used for config providers
	// (e.g., "claude-code" auto-installs "claude")
	AutoInstallSlashID string
}

// ProviderFactory is a function that creates a new instance of a provider
type ProviderFactory func() providerkit.Provider

// ProviderRegistration combines metadata with the provider factory
type ProviderRegistration struct {
	Metadata ProviderMetadata
	Factory  ProviderFactory
}

// Registry is the global provider registry
type Registry struct {
	mu        sync.RWMutex
	providers map[string]*ProviderRegistration
	// Maps config provider ID to slash provider ID
	configToSlash map[string]string
}

var globalRegistry = &Registry{
	providers:     make(map[string]*ProviderRegistration),
	configToSlash: make(map[string]string),
}

// Register registers a provider with its metadata and factory function.
// This should be called from init() functions in individual provider files.
//
// Returns an error if:
//   - A provider with the same ID is already registered
//   - Metadata is invalid (empty ID, name, or file paths)
func Register(metadata ProviderMetadata, factory ProviderFactory) error {
	if err := validateMetadata(metadata); err != nil {
		return fmt.Errorf(
			"invalid metadata for provider %s: %w",
			metadata.ID,
			err,
		)
	}

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if _, exists := globalRegistry.providers[metadata.ID]; exists {
		return fmt.Errorf("provider %s is already registered", metadata.ID)
	}

	globalRegistry.providers[metadata.ID] = &ProviderRegistration{
		Metadata: metadata,
		Factory:  factory,
	}

	// Track auto-install relationship
	if metadata.Type == TypeConfig && metadata.AutoInstallSlashID != "" {
		globalRegistry.configToSlash[metadata.ID] = metadata.AutoInstallSlashID
	}

	return nil
}

// MustRegister is like Register but panics on error.
// Use this in init() functions where registration should never fail.
func MustRegister(metadata ProviderMetadata, factory ProviderFactory) {
	if err := Register(metadata, factory); err != nil {
		panic(fmt.Sprintf("failed to register provider: %v", err))
	}
}

// GetProvider retrieves a provider by ID and creates a new instance.
// Returns nil if the provider is not found.
func GetProvider(id string) providerkit.Provider {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	reg, exists := globalRegistry.providers[id]
	if !exists {
		return nil
	}

	return reg.Factory()
}

// GetMetadata retrieves metadata for a provider by ID.
// Returns an error if the provider is not found.
func GetMetadata(id string) (ProviderMetadata, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	reg, exists := globalRegistry.providers[id]
	if !exists {
		return ProviderMetadata{}, fmt.Errorf("provider %s not found", id)
	}

	return reg.Metadata, nil
}

// ListProviders returns all registered providers sorted by priority.
// Lower priority numbers come first.
func ListProviders() []ProviderMetadata {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	providers := make([]ProviderMetadata, 0, len(globalRegistry.providers))
	for _, reg := range globalRegistry.providers {
		providers = append(providers, reg.Metadata)
	}

	// Sort by priority (lower first), then by name for stable ordering
	sort.Slice(providers, func(i, j int) bool {
		if providers[i].Priority != providers[j].Priority {
			return providers[i].Priority < providers[j].Priority
		}

		return providers[i].Name < providers[j].Name
	})

	return providers
}

// ListProvidersByType returns all registered providers of a specific type,
// sorted by priority.
func ListProvidersByType(providerType ProviderType) []ProviderMetadata {
	all := ListProviders()
	filtered := make([]ProviderMetadata, 0)

	for _, p := range all {
		if p.Type == providerType {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

// GetSlashProviderForConfig returns the slash provider ID that should be
// auto-installed when the given config provider is selected.
// Returns empty string and false if no auto-install relationship exists.
func GetSlashProviderForConfig(configID string) (string, bool) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	slashID, exists := globalRegistry.configToSlash[configID]

	return slashID, exists
}

// ProviderExists checks if a provider with the given ID is registered.
func ProviderExists(id string) bool {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	_, exists := globalRegistry.providers[id]

	return exists
}

// validateMetadata ensures provider metadata is valid
func validateMetadata(m ProviderMetadata) error {
	if m.ID == "" {
		return errors.New("provider ID cannot be empty")
	}
	if m.Name == "" {
		return errors.New("provider name cannot be empty")
	}
	if len(m.FilePaths) == 0 {
		return errors.New("provider must specify at least one file path")
	}
	if m.Type != TypeConfig && m.Type != TypeSlash {
		return fmt.Errorf(
			"provider type must be %q or %q",
			TypeConfig,
			TypeSlash,
		)
	}

	return nil
}

// Helper constructors for common metadata patterns

// ConfigParams holds parameters for creating config provider metadata
type ConfigParams struct {
	ID, Name, ConfigFilePath, SlashID string
	Priority                          int
}

// NewConfigMetadata creates metadata for a config-based provider
func NewConfigMetadata(params ConfigParams) ProviderMetadata {
	return ProviderMetadata{
		ID:                 params.ID,
		Name:               params.Name,
		Type:               TypeConfig,
		Priority:           params.Priority,
		FilePaths:          []string{params.ConfigFilePath},
		AutoInstallSlashID: params.SlashID,
	}
}

// NewSlashMetadata creates metadata for a slash command provider
func NewSlashMetadata(
	id, name string, filePaths []string, priority int,
) ProviderMetadata {
	return ProviderMetadata{
		ID:        id,
		Name:      name,
		Type:      TypeSlash,
		Priority:  priority,
		FilePaths: filePaths,
	}
}
