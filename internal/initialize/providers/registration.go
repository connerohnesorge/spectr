package providers

// Registration contains provider metadata and implementation.
type Registration struct {
	ID       string   // Unique identifier (kebab-case, e.g., "claude-code")
	Name     string   // Human-readable name (e.g., "Claude Code")
	Priority int      // Display order (lower = higher priority)
	Provider Provider // Implementation
}
