package providers

// Registration holds metadata for a provider registration.
type Registration struct {
	// ID is the unique provider identifier (kebab-case).
	ID string
	// Name is the human-readable provider name.
	Name string
	// Priority is the display order (lower = higher priority).
	Priority int
	// Provider is the provider implementation.
	Provider Provider
}
