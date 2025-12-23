package providers

// Registration holds provider metadata and the provider instance.
// This separates provider identity from implementation, reducing boilerplate.
//
// Example usage:
//
//	providers.Register(providers.Registration{
//	    ID:       "claude-code",
//	    Name:     "Claude Code",
//	    Priority: 1,
//	    Provider: &ClaudeProvider{},
//	})
type Registration struct {
	// ID is the unique provider identifier (kebab-case).
	// Example: "claude-code", "gemini", "cline"
	ID string

	// Name is the human-readable provider name for display.
	// Example: "Claude Code", "Gemini CLI", "Cline"
	Name string

	// Priority determines the display order (lower = higher priority).
	// Claude Code should be 1, other major tools 2-10, etc.
	Priority int

	// Provider is the provider instance that returns initializers.
	Provider Provider
}
