package cmd

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/connerohnesorge/spectr/internal/discovery"
)

// discoveryCache caches the results of GetDiscoveredRoots to avoid
// redundant filesystem traversals within a single command execution.
var (
	cachedRoots      []discovery.SpectrRoot
	errCachedRoots   error
	cachedCwd        string
	discoveryCacheMu sync.Mutex
)

// GetDiscoveredRoots returns all discovered spectr roots from the current
// working directory. It wraps discovery.FindSpectrRoots with cwd handling.
// Results are cached per working directory to avoid redundant
// filesystem traversals within a single command execution.
func GetDiscoveredRoots() ([]discovery.SpectrRoot, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get current directory: %w",
			err,
		)
	}

	discoveryCacheMu.Lock()
	defer discoveryCacheMu.Unlock()

	// If cwd changed, invalidate cache
	if cachedCwd != cwd {
		cachedRoots = nil
		errCachedRoots = nil
		cachedCwd = cwd
	}

	// Return cached result if available
	if cachedRoots != nil || errCachedRoots != nil {
		return cachedRoots, errCachedRoots
	}

	// Perform discovery
	roots, err := discovery.FindSpectrRoots(cwd)
	if err != nil {
		errCachedRoots = fmt.Errorf(
			"failed to discover spectr roots: %w",
			err,
		)

		return nil, errCachedRoots
	}

	cachedRoots = roots

	return cachedRoots, nil
}

// ResetDiscoveryCache clears the discovery cache. This is primarily
// intended for testing where the working directory changes.
func ResetDiscoveryCache() {
	discoveryCacheMu.Lock()
	defer discoveryCacheMu.Unlock()
	cachedRoots = nil
	errCachedRoots = nil
	cachedCwd = ""
}

// GetSingleRoot returns the first discovered root, or an error if no roots
// are found. This is useful for commands that operate on a single root.
func GetSingleRoot() (discovery.SpectrRoot, error) {
	roots, err := GetDiscoveredRoots()
	if err != nil {
		return discovery.SpectrRoot{}, err
	}

	if len(roots) == 0 {
		return discovery.SpectrRoot{}, errors.New(
			"no spectr directory found\nHint: Run 'spectr init' to initialize Spectr",
		)
	}

	return roots[0], nil
}

// HasMultipleRoots returns true if there are multiple discovered roots.
func HasMultipleRoots(roots []discovery.SpectrRoot) bool {
	return len(roots) > 1
}
