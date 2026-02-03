package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/discovery"
)

// GetDiscoveredRoots returns all discovered spectr roots from the current
// working directory. It wraps discovery.FindSpectrRoots with cwd handling.
func GetDiscoveredRoots() ([]discovery.SpectrRoot, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get current directory: %w",
			err,
		)
	}

	roots, err := discovery.FindSpectrRoots(cwd)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover spectr roots: %w",
			err,
		)
	}

	return roots, nil
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
