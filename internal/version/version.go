// Package version provides build information for the spectr binary.
package version

import (
	"encoding/json"
	"fmt"
)

// Build information variables set via ldflags during compilation.
// Example: go build -ldflags
// "-X github.com/connerohnesorge/spectr/internal/version.Version=v0.1.0"
var (
	// Version is the semantic version of the build.
	Version = "dev"

	// Commit is the git commit hash of the build.
	Commit = "unknown"

	// Date is the timestamp when the binary was built.
	Date = "unknown"
)

// BuildInfo contains version and build metadata.
type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// GetBuildInfo returns the current build information.
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

// String returns a formatted multi-line representation of build info.
func (b BuildInfo) String() string {
	return fmt.Sprintf(
		"Version:  %s\nCommit:   %s\nDate:     %s",
		b.Version,
		b.Commit,
		b.Date,
	)
}

// JSON returns the build info as JSON bytes.
func (b BuildInfo) JSON() ([]byte, error) {
	return json.Marshal(b)
}

// Short returns just the version string.
func (b BuildInfo) Short() string {
	return b.Version
}
