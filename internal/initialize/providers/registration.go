package providers

import "errors"

// Error constants for registration validation.
var (
	// ErrProviderIDRequired is returned when a registration has an empty ID.
	ErrProviderIDRequired = errors.New(
		"provider ID is required",
	)

	// ErrProviderRequired is returned when a registration has a nil Provider.
	ErrProviderRequired = errors.New(
		"provider implementation is required",
	)
)

// Registration contains provider metadata and implementation.
//
// This separates provider identity (ID, Name, Priority) from
// initialization logic. Provider metadata is stored at registration time,
// allowing the Provider interface to focus on returning initializers.
type Registration struct {
	// ID is the unique provider identifier (kebab-case).
	// Example: "claude-code", "gemini", "cline"
	// Must be non-empty and unique across all registrations.
	ID string

	// Name is the human-readable provider name for display in UI.
	// Example: "Claude Code", "Gemini CLI", "Cline"
	Name string

	// Priority determines the display order in provider lists.
	// Lower values appear first (higher priority).
	// Example: Claude Code = 1, Gemini = 2, Cursor = 3, etc.
	Priority int

	// Provider is the implementation that returns initializers.
	Provider Provider
}

// Validate checks if the registration is valid.
// Returns an error if any required fields are missing or invalid.
func (r *Registration) Validate() error {
	if r.ID == "" {
		return ErrProviderIDRequired
	}
	if r.Provider == nil {
		return ErrProviderRequired
	}

	return nil
}
