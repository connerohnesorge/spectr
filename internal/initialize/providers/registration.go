// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the Registration struct which holds metadata for a
// provider. The Registration struct separates concerns: providers implement
// behavior through the Provider interface, while metadata (ID, Name,
// Priority) is provided at registration time via this struct.
//
// This design means providers don't need to know their own ID, name, or
// priority. That information is a concern of the registry, not the provider.
//
// Example usage:
//
//	providers.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
//
// See also:
//   - provider_new.go: Provider interface for provider behavior
//   - registry.go: Registry that accepts Registration structs
package providers

import (
	"fmt"
	"regexp"
)

// kebabCaseRegex matches valid kebab-case identifiers.
// Valid examples: "claude-code", "cursor", "open-code-ai"
// Invalid examples: "ClaudeCode", "claude_code", "claude--code", "-claude"
var kebabCaseRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

// Registration holds the metadata for a provider at registration time.
//
// The Registration struct is used when registering a provider with the
// registry. It separates provider metadata from the provider's behavior
// (the Provider interface). This design allows:
//
//   - Providers to focus solely on returning their initializers
//   - Registry to manage identification and priority ordering
//   - Clean separation between "what a provider does" and "how it's identified"
//
// Fields:
//
//   - ID: A unique kebab-case identifier (e.g., "claude-code", "cursor").
//     Used for programmatic access and command-line selection.
//     Must be lowercase, alphanumeric with hyphens, starting with a letter.
//
//   - Name: A human-readable display name (e.g., "Claude Code", "Cursor").
//     Used in user-facing output like provider lists and confirmation prompts.
//     Must be non-empty.
//
//   - Priority: Determines the default ordering when multiple providers are
//     available. Lower numbers indicate higher priority. Must be >= 0.
//     Priority 0 is reserved for the most commonly used providers.
//
//   - Provider: The Provider implementation that returns initializers.
//     Must not be nil.
//
//nolint:revive // line-length-limit - struct documentation
type Registration struct {
	// ID is a unique kebab-case identifier for the provider.
	// Examples: "claude-code", "cursor", "gemini-cli"
	//
	// Constraints:
	//   - Must be kebab-case (lowercase letters, numbers, hyphens)
	//   - Must start with a lowercase letter
	//   - Must not contain consecutive hyphens
	//   - Must not start or end with a hyphen
	ID string

	// Name is the human-readable display name for the provider.
	// Examples: "Claude Code", "Cursor", "Gemini CLI"
	//
	// Constraints:
	//   - Must be non-empty
	Name string

	// Priority determines the default ordering in provider lists.
	// Lower numbers appear first and are considered "higher priority".
	//
	// Constraints:
	//   - Must be >= 0
	//
	// Guidelines:
	//   - Priority 0-9: Most popular/commonly used providers
	//   - Priority 10-99: Standard providers
	//   - Priority 100+: Specialized or less common providers
	Priority int

	// Provider is the Provider implementation that returns initializers.
	//
	// Constraints:
	//   - Must not be nil
	Provider Provider
}

// Validate checks that all Registration fields meet their constraints.
// Returns an error describing all validation failures, or nil if valid.
//
// Validation rules:
//   - ID must be valid kebab-case (lowercase, alphanumeric, hyphens only)
//   - Name must be non-empty
//   - Priority must be >= 0
//   - Provider must not be nil
func (r Registration) Validate() error {
	var errs []string

	// Validate ID
	if r.ID == "" {
		errs = append(errs, "ID is required")
	} else if !kebabCaseRegex.MatchString(r.ID) {
		errs = append(errs, fmt.Sprintf(
			"ID %q must be kebab-case (lowercase, hyphens, no consecutive)",
			r.ID,
		))
	}

	// Validate Name
	if r.Name == "" {
		errs = append(errs, "Name is required")
	}

	// Validate Priority
	if r.Priority < 0 {
		errs = append(
			errs,
			fmt.Sprintf("Priority must be >= 0, got %d", r.Priority),
		)
	}

	// Validate Provider
	if r.Provider == nil {
		errs = append(errs, "Provider is required")
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid registration: %v", errs)
	}

	return nil
}

// IsValid returns true if the Registration passes all validation checks.
// This is a convenience method; use Validate() to get detailed error messages.
func (r Registration) IsValid() bool {
	return r.Validate() == nil
}
