package providers

import "github.com/connerohnesorge/spectr/internal/initialize/types"

// Provider defines the new provider interface for Spectr.
// Unlike the legacy interface, this one is minimal and composable.
type Provider interface {
	// Initializers returns the list of initializers for this provider.
	// These initializers will be executed during spectr init.
	Initializers() []types.Initializer
}
