// Package providers defines interfaces for AI agent providers used in Spectr initialization.
package providers

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/connerohnesorge/spectr/internal/ralph"
)

// Ralpher is a type alias for ralph.Ralpher for convenience in the providers package.
// This allows providers to implement the ralph.Ralpher interface while keeping
// helper functions in this package.
type Ralpher = ralph.Ralpher

// IsBinaryAvailable checks if a CLI binary is available on the system PATH.
//
// This function uses exec.LookPath to determine if the specified binary can be
// found in any directory listed in the system's PATH environment variable.
//
// Parameters:
//   - binaryName: The name of the binary to check (e.g., "claude", "gemini", "bash")
//
// Returns:
//   - true if the binary is found on PATH
//   - false if the binary is not found or if binaryName is empty
//
// Example usage:
//
//	if IsBinaryAvailable("claude") {
//	    fmt.Println("Claude CLI is installed and available")
//	}
//
// Note: This function does not verify that the binary is executable or that it
// will work correctly. It only checks for existence on PATH.
func IsBinaryAvailable(binaryName string) bool {
	if binaryName == "" {
		return false
	}
	_, err := exec.LookPath(binaryName)

	return err == nil
}

// IsRalpherAvailable checks if a provider implements Ralpher and its binary is available.
//
// This helper function combines two checks:
//  1. Does the provider implement the Ralpher interface?
//  2. Is the provider's CLI binary available on the system PATH?
//
// Parameters:
//   - provider: Any provider instance (may or may not implement Ralpher)
//
// Returns:
//   - true if the provider implements Ralpher AND its binary is available
//   - false otherwise (provider doesn't implement Ralpher, or binary not found)
//
// Example usage:
//
//	provider := providers.NewClaudeProvider()
//	if IsRalpherAvailable(provider) {
//	    // Safe to use provider for ralph orchestration
//	    ralpher := provider.(Ralpher)
//	    // ...
//	} else {
//	    fmt.Println("Provider does not support ralph or binary not found")
//	}
//
// This is useful when iterating through multiple providers to find one that
// can be used for task orchestration:
//
//	for _, provider := range allProviders {
//	    if IsRalpherAvailable(provider) {
//	        return provider.(Ralpher)
//	    }
//	}
func IsRalpherAvailable(provider Provider) bool {
	ralpher, ok := provider.(Ralpher)
	if !ok {
		return false
	}

	return IsBinaryAvailable(ralpher.Binary())
}

// ValidateRalpher checks if a provider can be used for ralph orchestration and returns a descriptive error if not.
//
// This function provides detailed error messages that can be shown to users when
// a provider cannot be used for task orchestration.
//
// Parameters:
//   - provider: Any provider instance (may or may not implement Ralpher)
//
// Returns:
//   - nil if the provider implements Ralpher and its binary is available
//   - error describing the issue if the provider cannot be used:
//   - Provider doesn't implement Ralpher interface
//   - Binary not found on PATH
//
// Example usage:
//
//	provider := providers.NewClaudeProvider()
//	if err := ValidateRalpher(provider); err != nil {
//	    return fmt.Errorf("cannot use provider for ralph: %w", err)
//	}
//	ralpher := provider.(Ralpher)
//	// Safe to use ralpher for orchestration
//
// This is particularly useful in CLI commands that need to provide helpful
// error messages to users:
//
//	func Run(providerID string) error {
//	    provider := getProvider(providerID)
//	    if err := ValidateRalpher(provider); err != nil {
//	        return fmt.Errorf("provider %s: %w", providerID, err)
//	    }
//	    // ...
//	}
func ValidateRalpher(provider Provider) error {
	ralpher, ok := provider.(Ralpher)
	if !ok {
		return errors.New("provider does not implement Ralpher interface")
	}

	binaryName := ralpher.Binary()
	if binaryName == "" {
		return errors.New("provider Binary() returned empty string")
	}

	if !IsBinaryAvailable(binaryName) {
		return fmt.Errorf("binary %q not found on PATH", binaryName)
	}

	return nil
}
