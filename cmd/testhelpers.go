package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// FindProjectRoot finds the project root by walking up from the current directory
// until it finds a go.mod file. This is more reliable than assuming relative paths.
func FindProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		// Check if go.mod exists in this directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod
			t.Fatalf("Could not find project root (go.mod) from %s", dir)
		}
		dir = parent
	}
}

// GetSpectrDir returns the path to the spectr directory relative to project root
func GetSpectrDir(t *testing.T) string {
	projectRoot := FindProjectRoot(t)
	spectrDir := filepath.Join(projectRoot, "spectr")
	if _, err := os.Stat(spectrDir); err != nil {
		t.Skipf("spectr directory not found at %s", spectrDir)
	}

	return spectrDir
}

// GetTestDataDir returns the path to the testdata directory
func GetTestDataDir(t *testing.T) string {
	projectRoot := FindProjectRoot(t)

	return filepath.Join(projectRoot, "testdata")
}
