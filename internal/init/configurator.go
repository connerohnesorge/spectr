// Package init provides initialization functionality for Spectr projects.
//
// This package orchestrates the initialization wizard, executes tool
// configuration, and manages the initialization state. Actual tool
// implementations are in the internal/providers package.
package init

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

// Configurator is an alias for providerkit.Provider, maintained for
// backward compatibility with existing code in this package.
//
// The Provider interface defines the contract for all AI tool configurators.
// See providerkit.Provider for full documentation.
type Configurator = providerkit.Provider
