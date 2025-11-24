// Package version provides version information embedded from the VERSION
// file. This package handles embedding and parsing the VERSION file which
// is generated during build by GoReleaser or Nix.
//
// The VERSION file must be located at the project root. Go's embed directive
// requires the file to be within or below the package directory, so we use
// a workaround: the build system (GoReleaser/Nix) copies VERSION to this
// directory before building, and we embed it from here.
package version

import (
	_ "embed"
	"encoding/json"
	"errors"
)

const (
	// DefaultVersion is used when version is not specified.
	DefaultVersion = "dev"
	// DefaultValue is used for unknown commit/date values.
	DefaultValue = "unknown"
)

// versionFile holds the embedded VERSION file content.
// The VERSION file is copied here during build by GoReleaser or Nix and
// contains version metadata in JSON format: version, commit, date.
//
//go:embed VERSION
var versionFile []byte

// Info represents the structure of the VERSION file.
// It contains metadata about the build that is embedded at compile time.
type Info struct {
	Version string `json:"version"` // Version number (e.g., "v1.2.3" or "dev")
	Commit  string `json:"commit"`  // Git commit hash (short form, 7 chars)
	Date    string `json:"date"`    // Build date in ISO 8601 format
}

// Get reads and parses the embedded VERSION file.
// Returns version info with defaults applied for empty fields.
// If parsing fails, returns development defaults.
func Get() Info {
	info, err := parse()
	if err != nil {
		// Fall back to development defaults if parsing fails
		return Info{
			Version: DefaultVersion,
			Commit:  DefaultValue,
			Date:    DefaultValue,
		}
	}

	return info
}

// parse reads and parses the embedded VERSION file.
// Returns an error if the file is empty or contains invalid JSON.
func parse() (Info, error) {
	var info Info

	// Check if VERSION file is empty
	if len(versionFile) == 0 {
		return info, errors.New("VERSION file is empty")
	}

	// Parse JSON from embedded file
	if err := json.Unmarshal(versionFile, &info); err != nil {
		return info, err
	}

	// Apply defaults for empty fields
	if info.Version == "" {
		info.Version = DefaultVersion
	}
	if info.Commit == "" {
		info.Commit = DefaultValue
	}
	if info.Date == "" {
		info.Date = DefaultValue
	}

	return info, nil
}
