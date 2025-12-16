// Package initialize provides utilities for initializing Spectr
// in a project directory.
//
// This file contains helper functions for the executor to reduce file length
// and improve code organization.
package initialize

import (
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

// ExecuteOptions holds configuration options for the Execute method.
// This replaces boolean flag parameters to avoid control coupling.
type ExecuteOptions struct {
	// CIWorkflowEnabled specifies whether to create the CI workflow file.
	CIWorkflowEnabled bool
}

// NewExecuteOptions creates ExecuteOptions with default values.
func NewExecuteOptions() ExecuteOptions {
	return ExecuteOptions{
		CIWorkflowEnabled: false,
	}
}

// WithCIWorkflow returns a copy of options with CI workflow enabled.
func (o ExecuteOptions) WithCIWorkflow(enabled bool) ExecuteOptions {
	o.CIWorkflowEnabled = enabled

	return o
}

// FileTrackingResult indicates whether a file was created or updated.
type FileTrackingResult int

const (
	// FileCreated indicates a new file was created.
	FileCreated FileTrackingResult = iota
	// FileUpdated indicates an existing file was updated.
	FileUpdated
)

// trackFile adds the file path to the appropriate result list based on
// whether it was created or updated.
func trackFile(
	result *ExecutionResult,
	path string,
	tracking FileTrackingResult,
) {
	if tracking == FileUpdated {
		result.UpdatedFiles = append(result.UpdatedFiles, path)

		return
	}
	result.CreatedFiles = append(result.CreatedFiles, path)
}

// trackInitializerByKey parses the initializer key and tracks appropriate
// file changes. This extracts the file tracking logic from
// trackInitializerFiles.
func trackInitializerByKey(
	key string,
	tracking FileTrackingResult,
	result *ExecutionResult,
) {
	switch {
	case len(key) > keyPrefixDir && key[:keyPrefixDir] == "dir:":
		trackDirectoryInitializer(key, tracking, result)
	case len(key) > keyPrefixConfig && key[:keyPrefixConfig] == "config:":
		trackConfigInitializer(key, tracking, result)
	case len(key) > keyPrefixSlashCmds &&
		key[:keyPrefixSlashCmds] == "slashcmds:":
		trackSlashCmdsInitializer(key, tracking, result)
	}
}

// trackDirectoryInitializer tracks directory initializer file changes.
func trackDirectoryInitializer(
	key string,
	tracking FileTrackingResult,
	result *ExecutionResult,
) {
	path := key[keyPrefixDir:] + "/"
	trackFile(result, path, tracking)
}

// trackConfigInitializer tracks config file initializer changes.
func trackConfigInitializer(
	key string,
	tracking FileTrackingResult,
	result *ExecutionResult,
) {
	path := key[keyPrefixConfig:]
	trackFile(result, path, tracking)
}

// trackSlashCmdsInitializer tracks slash commands initializer file changes.
func trackSlashCmdsInitializer(
	key string,
	tracking FileTrackingResult,
	result *ExecutionResult,
) {
	parts := parseSlashCmdsKey(key)
	if parts.dir == "" {
		return
	}
	proposal := filepath.Join(parts.dir, "proposal"+parts.ext)
	apply := filepath.Join(parts.dir, "apply"+parts.ext)
	trackFile(result, proposal, tracking)
	trackFile(result, apply, tracking)
}

// slashCmdsKeyParts holds parsed parts from a slash commands key.
type slashCmdsKeyParts struct {
	dir string
	ext string
}

// parseSlashCmdsKey parses a slash commands key into its parts.
// Format: "slashcmds:<dir>:<ext>:<format>"
func parseSlashCmdsKey(key string) slashCmdsKeyParts {
	// Remove "slashcmds:" prefix
	rest := key[keyPrefixSlashCmds:]

	// Find the last two colons (extension and format)
	lastColon := -1
	secondLastColon := -1
	for i := len(rest) - 1; i >= 0; i-- {
		if rest[i] != ':' {
			continue
		}
		if lastColon == -1 {
			lastColon = i

			continue
		}
		secondLastColon = i

		break
	}

	if secondLastColon == -1 {
		return slashCmdsKeyParts{}
	}

	return slashCmdsKeyParts{
		dir: rest[:secondLastColon],
		ext: rest[secondLastColon+1 : lastColon],
	}
}

// trackInitializerFiles tracks which files were created/updated by an
// initializer. This provides feedback to the user about what changed.
// The tracking parameter indicates whether the file was created or updated.
func trackInitializerFiles(
	init providers.Initializer,
	tracking FileTrackingResult,
	result *ExecutionResult,
) {
	keyer, ok := init.(providers.Keyer)
	if !ok {
		return
	}

	key := keyer.Key()
	trackInitializerByKey(key, tracking, result)
}
