package providers

// Registration contains provider metadata and implementation for
// registration.
//
// This struct is used at registration time to associate a provider
// implementation with its metadata (ID, name, display priority).
// By separating metadata from the Provider interface, we eliminate
// boilerplate and allow providers to focus purely on what they
// initialize.
//
// Fields:
//   - ID: Unique kebab-case identifier
//     (e.g., "claude-code", "cursor", "cline")
//   - Name: Human-readable display name
//     (e.g., "Claude Code", "Cursor", "Cline")
//   - Priority: Display order
//     (lower = higher priority, e.g., 1 = first)
//   - Provider: The provider implementation that returns
//     Initializers
//
// Example usage in a provider's init() function:
//
//	func init() {
//	    providers.Register(providers.Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	}
//
// The actual registry implementation (Register function,
// Registry type) will be implemented in task 4.x. This struct defines
// the data structure for registration.
type Registration struct {
	// ID is the unique kebab-case identifier for this provider.
	// Used for provider selection and internal lookups.
	// Examples: "claude-code", "cursor", "cline", "aider"
	ID string

	// Name is the human-readable display name for this provider.
	// Used in CLI output and documentation.
	// Examples: "Claude Code", "Cursor", "Cline", "Aider"
	Name string

	// Priority determines the display order when listing providers.
	// Lower values = higher priority (displayed first).
	// Example: 1 = first, 2 = second, etc.
	Priority int

	// Provider is the implementation that returns this provider's Initializers.
	// The Provider.Initializers() method will be called during initialization
	// to collect the list of tasks this provider needs to perform.
	Provider Provider
}
