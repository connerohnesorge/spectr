package providers

// Registration contains provider metadata and implementation.
// This struct is used during provider registration to associate
// a provider implementation with its identifying information.
type Registration struct {
	// ID is the unique provider identifier (kebab-case).
	// Examples: "claude-code", "gemini", "cline"
	// Must be unique across all registered providers.
	ID string

	// Name is the human-readable provider name for display.
	// Examples: "Claude Code", "Gemini CLI", "Cline"
	// Used in UI elements like the setup wizard.
	Name string

	// Priority determines the display order (lower = higher priority).
	// Claude Code should be 1, other major tools 2-10, etc.
	// Used to sort providers in lists and determine which provider wins
	// during initializer deduplication (higher priority = kept first).
	Priority int

	// Provider is the implementation that returns initializers.
	// Must not be nil.
	Provider Provider
}
