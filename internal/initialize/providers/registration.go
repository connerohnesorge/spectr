// Package providers defines the core interfaces for the provider architecture.
//
// This file contains the Registration struct, which holds metadata for
// registering providers with the registry. The Registration struct separates
// provider metadata (ID, Name, Priority) from the provider implementation,
// following the design principle that providers don't need to know their
// own registry metadata.
package providers

// Registration holds the metadata required to register a provider.
//
// Registration separates provider metadata from the provider implementation.
// This design recognizes that metadata like ID, Name, and Priority are
// registry concerns, not provider concerns. Providers only need to know
// how to create their initializers.
//
// # Fields
//
// ID is a unique kebab-case identifier for the provider (e.g., "claude-code",
// "cursor", "gemini"). This ID is used:
//   - For programmatic selection of providers
//   - As a stable key that won't change if the display name changes
//   - In configuration files and command-line arguments
//
// Name is a human-readable display name (e.g., "Claude Code", "Cursor").
// This is shown to users in selection menus and status output.
//
// Priority determines display and execution order. Lower values appear first
// in provider lists. When users select multiple providers, they are processed
// in priority order.
//
// Provider is the actual provider implementation that returns initializers.
// This can be either a struct implementing ProviderV2 or a ProviderV2Func.
//
// # Example: Struct Provider
//
//	type ClaudeProvider struct{}
//
//	func (p *ClaudeProvider) Initializers(ctx context.Context) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".claude/commands/spectr"),
//	        NewConfigFileInitializer("CLAUDE.md", InstructionTemplate),
//	        NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown),
//	    }
//	}
//
//	func init() {
//	    providers.Register(Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	}
//
// # Example: Function Provider
//
// For simple providers, you can use ProviderV2Func instead of creating a struct:
//
//	func init() {
//	    providers.Register(Registration{
//	        ID:       "simple-tool",
//	        Name:     "Simple Tool",
//	        Priority: 50,
//	        Provider: ProviderV2Func(func(ctx context.Context) []Initializer {
//	            return []Initializer{
//	                NewConfigFileInitializer("SIMPLE.md", template),
//	            }
//	        }),
//	    })
//	}
//
// # Priority Guidelines
//
// Priority values are used to sort providers for display and processing.
// Lower values appear first. Suggested ranges:
//
//	1-10:   Primary CLI tools (Claude Code, Cursor)
//	11-20:  Popular IDE extensions (Cline, Continue)
//	21-50:  Specialized tools
//	51-99:  Experimental or less common tools
//
// # Validation
//
// The registry validates registrations and rejects:
//   - Empty ID
//   - Empty Name
//   - Nil Provider
//   - Duplicate ID (already registered)
//
// # Migration Note
//
// This struct is part of the new provider architecture (v2). After the
// migration is complete (tasks 5.x-7.x), the Registration struct will
// be the standard way to register all providers.
type Registration struct {
	// ID is the unique identifier for this provider.
	//
	// The ID must be:
	//   - Non-empty
	//   - Unique across all registered providers
	//   - Kebab-case (e.g., "claude-code", not "ClaudeCode" or "claude_code")
	//   - Stable (should not change between versions)
	//
	// The ID is used for:
	//   - Programmatic provider selection
	//   - Configuration file references
	//   - Command-line arguments
	//
	// Examples: "claude-code", "cursor", "gemini", "cline", "aider"
	ID string

	// Name is the human-readable display name for this provider.
	//
	// The Name should be:
	//   - Non-empty
	//   - User-friendly (proper capitalization, spaces)
	//   - The official name of the tool
	//
	// The Name is used in:
	//   - Provider selection menus
	//   - Status and progress output
	//   - Documentation
	//
	// Examples: "Claude Code", "Cursor", "Gemini CLI", "Cline", "Aider"
	Name string

	// Priority determines the display and processing order.
	//
	// Lower values appear first in provider lists and are processed first.
	// This allows primary tools to be shown prominently while keeping
	// specialized tools accessible but not cluttering the interface.
	//
	// Suggested ranges:
	//   1-10:   Primary CLI tools (Claude Code, Cursor, Gemini)
	//   11-20:  Popular IDE extensions (Cline, Continue, Windsurf)
	//   21-50:  Specialized tools (Codex, Costrict, Qoder)
	//   51-99:  Experimental or less common tools
	//
	// Priority 0 is valid and will sort before priority 1.
	Priority int

	// Provider is the provider implementation.
	//
	// Provider must implement the ProviderV2 interface, which has a single
	// method: Initializers(ctx context.Context) []Initializer
	//
	// The Provider can be:
	//   - A pointer to a struct implementing ProviderV2
	//   - A ProviderV2Func for simple function-based providers
	//
	// Provider must not be nil.
	Provider ProviderV2
}

// Validate checks that the Registration has all required fields.
//
// This method is called by the registry during registration to ensure
// that all required fields are present and valid.
//
// Returns nil if the registration is valid, or an error describing
// the validation failure.
//
// Validation rules:
//   - ID must not be empty
//   - Name must not be empty
//   - Provider must not be nil
//
// Note: Priority is always valid (including 0 and negative values).
// Duplicate ID checking is handled by the registry, not this method.
func (r Registration) Validate() error {
	if r.ID == "" {
		return ErrEmptyID
	}
	if r.Name == "" {
		return ErrEmptyName
	}
	if r.Provider == nil {
		return ErrNilProvider
	}
	return nil
}

// RegistrationError represents an error during provider registration.
type RegistrationError string

func (e RegistrationError) Error() string {
	return string(e)
}

// Registration validation errors.
const (
	// ErrEmptyID indicates that the Registration.ID field is empty.
	ErrEmptyID RegistrationError = "registration ID must not be empty"

	// ErrEmptyName indicates that the Registration.Name field is empty.
	ErrEmptyName RegistrationError = "registration Name must not be empty"

	// ErrNilProvider indicates that the Registration.Provider field is nil.
	ErrNilProvider RegistrationError = "registration Provider must not be nil"

	// ErrDuplicateID indicates that a provider with the same ID is already registered.
	// This error is returned by the registry, not by Registration.Validate().
	ErrDuplicateID RegistrationError = "provider with this ID is already registered"
)
