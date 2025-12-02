package git

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

// Platform represents the git hosting platform type.
type Platform int

const (
	// PlatformUnknown indicates an unrecognized git hosting platform.
	PlatformUnknown Platform = iota
	// PlatformGitHub indicates GitHub hosting.
	PlatformGitHub
	// PlatformGitLab indicates GitLab hosting.
	PlatformGitLab
	// PlatformGitea indicates Gitea or Forgejo hosting.
	PlatformGitea
	// PlatformBitbucket indicates Bitbucket hosting.
	PlatformBitbucket
)

// urlMatchGroups is the expected number of match groups for URL patterns.
const urlMatchGroups = 3

// String returns the string representation of the platform.
func (p Platform) String() string {
	switch p {
	case PlatformUnknown:
		return "Unknown"
	case PlatformGitHub:
		return "GitHub"
	case PlatformGitLab:
		return "GitLab"
	case PlatformGitea:
		return "Gitea"
	case PlatformBitbucket:
		return "Bitbucket"
	}

	return "Unknown"
}

// PlatformInfo contains information about the detected git platform.
type PlatformInfo struct {
	// Platform is the detected hosting platform type.
	Platform Platform
	// CLITool is the recommended CLI tool for this platform.
	// Examples: "gh" for GitHub, "glab" for GitLab.
	CLITool string
	// RepoURL is the normalized repository URL.
	RepoURL string
}

// sshURLPattern matches SSH-style git URLs like git@github.com:org/repo.git
var sshURLPattern = regexp.MustCompile(`^git@([^:]+):(.+?)(?:\.git)?$`)

// httpsURLPattern matches HTTPS git URLs like https://github.com/org/repo.git
var httpsURLPattern = regexp.MustCompile(`^https?://([^/]+)/(.+?)(?:\.git)?$`)

// DetectPlatform analyzes a git remote URL and returns platform information.
// It supports SSH (git@host:path) and HTTPS (https://host/path) URL formats.
func DetectPlatform(remoteURL string) PlatformInfo {
	url := strings.TrimSpace(remoteURL)
	if url == "" {
		return PlatformInfo{Platform: PlatformUnknown}
	}

	var host string

	// Try SSH URL format first
	sshMatches := sshURLPattern.FindStringSubmatch(url)
	httpsMatches := httpsURLPattern.FindStringSubmatch(url)

	switch {
	case len(sshMatches) == urlMatchGroups:
		host = strings.ToLower(sshMatches[1])
	case len(httpsMatches) == urlMatchGroups:
		host = strings.ToLower(httpsMatches[1])
	default:
		return PlatformInfo{Platform: PlatformUnknown, RepoURL: url}
	}

	info := PlatformInfo{RepoURL: url}

	// Detect platform based on host.
	// Uses strict matching to avoid false positives.
	switch {
	case host == "github.com" ||
		strings.HasSuffix(host, ".github.com") ||
		strings.HasPrefix(host, "github."):
		info.Platform = PlatformGitHub
		info.CLITool = "gh"
	case host == "gitlab.com" ||
		strings.HasSuffix(host, ".gitlab.com") ||
		strings.HasPrefix(host, "gitlab."):
		info.Platform = PlatformGitLab
		info.CLITool = "glab"
	case strings.HasPrefix(host, "gitea.") ||
		strings.HasPrefix(host, "forgejo."):
		info.Platform = PlatformGitea
		info.CLITool = "tea"
	case host == "bitbucket.org" ||
		strings.HasSuffix(host, ".bitbucket.org") ||
		strings.HasPrefix(host, "bitbucket."):
		info.Platform = PlatformBitbucket
		info.CLITool = "bb"
	default:
		info.Platform = PlatformUnknown
	}

	return info
}

// ErrNoOriginRemote is returned when no origin remote is configured.
var ErrNoOriginRemote = errors.New(
	"failed to get origin remote URL: " +
		"no origin remote configured or not a git repository",
)

// GetOriginRemote returns the URL of the origin remote.
// It executes 'git remote get-url origin' and returns the result.
func GetOriginRemote() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", ErrNoOriginRemote
	}

	return strings.TrimSpace(string(output)), nil
}
