// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the provider interfaces for the architecture.
// The design reduces provider boilerplate by separating concerns:
//   - Provider: Returns a list of initializers
//   - Initializer: Handles a single initialization step
//   - Registration: Contains provider metadata (ID, Name, Priority)
//   - Config: Holds configuration passed to initializers
package providers

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/spf13/afero"
)

// CommandFormat specifies the format for slash command files.
type CommandFormat int

const (
	// FormatMarkdown uses markdown files with
	// YAML frontmatter (Claude, Cline, etc.)
	FormatMarkdown CommandFormat = iota
	// FormatTOML uses TOML files (Gemini CLI)
	FormatTOML
)

// TemplateContext holds path-related template variables for dynamic
// directory names. This struct is defined in the providers package
// to avoid import cycles.
type TemplateContext struct {
	// BaseDir is the base directory for spectr files (default: "spectr")
	BaseDir string
	// SpecsDir is the directory for spec files (default: "spectr/specs")
	SpecsDir string
	// ChangesDir is the directory for change proposals
	// (default: "spectr/changes")
	ChangesDir string
}

// DefaultTemplateContext returns a TemplateContext with default values.
func DefaultTemplateContext() TemplateContext {
	return TemplateContext{
		BaseDir:    "spectr",
		SpecsDir:   "spectr/specs",
		ChangesDir: "spectr/changes",
	}
}

// TemplateRenderer provides template rendering capabilities.
//
// This interface allows providers to render templates without depending
// on the full TemplateManager.
type TemplateRenderer interface {
	// RenderAgents renders the AGENTS.md template content.
	RenderAgents(ctx TemplateContext) (string, error)
	// RenderInstructionPointer renders a short pointer template that
	// directs AI assistants to read spectr/AGENTS.md for full instructions.
	RenderInstructionPointer(ctx TemplateContext) (string, error)
	// RenderSlashCommand renders a slash command template
	// (proposal or apply).
	RenderSlashCommand(command string, ctx TemplateContext) (string, error)
}

// NewProvider is the new provider interface that returns a list of
// initializers. Providers no longer contain metadata (ID, Name, Priority) -
// that lives in Registration.
//
// This interface is intentionally minimal. Providers compose behavior by
// returning appropriate initializers for their requirements.
//
// Example implementation:
//
//	type ClaudeProvider struct{}
//
//	func (p *ClaudeProvider) Initializers(
//	    ctx context.Context,
//	) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".claude/commands/spectr"),
//	        NewConfigFileInitializer("CLAUDE.md", template),
//	        NewSlashCommandsInitializer(
//	            ".claude/commands/spectr", ".md", FormatMarkdown,
//	        ),
//	    }
//	}
type NewProvider interface {
	// Initializers returns the list of initializers needed to configure
	// this provider. The context can be used for cancellation or deadline
	// propagation.
	Initializers(
		ctx context.Context,
	) []Initializer
}

// Initializer represents a single initialization step. Each initializer
// handles one specific task (create directory, write config file, etc.).
//
// Initializers must be idempotent - running Init multiple times should
// produce the same result as running it once.
type Initializer interface {
	// Init performs the initialization step. The filesystem (fs) is rooted
	// at the project directory, so all paths are project-relative. Returns
	// an error if initialization fails.
	Init(
		ctx context.Context,
		fs afero.Fs,
		cfg *Config,
	) error

	// IsSetup returns true if this initializer's work is already complete.
	// This is used to determine if initialization is needed.
	IsSetup(fs afero.Fs, cfg *Config) bool
}

// Config holds configuration passed to initializers. This struct contains
// values that initializers need to know about the project structure.
type Config struct {
	// SpectrDir is the spectr directory relative to the project root.
	// Example: "spectr" (the default)
	SpectrDir string
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		SpectrDir: "spectr",
	}
}

// Registration holds provider metadata for the registry.
// This separates the "what does this provider do" (Provider interface)
// from the "how is it identified and ordered" (Registration).
type Registration struct {
	// ID is the unique provider identifier (kebab-case).
	// Example: "claude-code", "gemini-cli", "cline"
	ID string

	// Name is the human-readable provider name for display.
	// Example: "Claude Code", "Gemini CLI", "Cline"
	Name string

	// Priority is the display/processing order (lower = higher priority).
	// Claude Code should be 1, other major tools 2-10, etc.
	Priority int

	// Provider is the provider implementation.
	Provider NewProvider
}

// ProviderRegistry is the instance-only registry for the provider architecture.
// It stores Registration structs (which include provider metadata and the
// provider itself).
//
// This registry has NO global state - create instances with CreateRegistry().
// This design improves testability by eliminating shared state between tests.
//
// Thread-safety: All methods are safe for concurrent access.
type ProviderRegistry struct {
	mu            sync.RWMutex
	registrations map[string]Registration
}

// CreateRegistry creates a new empty registry instance.
// This is the primary constructor for the provider registry.
func CreateRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		registrations: make(map[string]Registration),
	}
}

// Register adds a provider registration to the registry.
// Returns an error if a registration with the same ID already exists.
func (r *ProviderRegistry) Register(reg Registration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registrations[reg.ID]; exists {
		return fmt.Errorf("provider %q already registered", reg.ID)
	}

	r.registrations[reg.ID] = reg

	return nil
}

// All returns all registrations sorted by priority (lower = higher priority).
// Returns an empty slice if no providers are registered.
func (r *ProviderRegistry) All() []Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Registration, 0, len(r.registrations))
	for _, reg := range r.registrations {
		result = append(result, reg)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result
}

// Get retrieves a registration by its provider ID.
// Returns nil if no registration with that ID exists.
func (r *ProviderRegistry) Get(id string) *Registration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if reg, exists := r.registrations[id]; exists {
		return &reg
	}

	return nil
}

// IDs returns all registered provider IDs sorted by priority.
// Returns an empty slice if no providers are registered.
func (r *ProviderRegistry) IDs() []string {
	all := r.All()
	ids := make([]string, len(all))
	for i, reg := range all {
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
