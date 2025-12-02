package git

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Platform represents a Git hosting platform.
type Platform string

const (
	// PlatformGitHub represents GitHub.
	PlatformGitHub Platform = "github"
	// PlatformGitLab represents GitLab.
	PlatformGitLab Platform = "gitlab"
	// PlatformGitea represents Gitea or Forgejo.
	PlatformGitea Platform = "gitea"
	// PlatformBitbucket represents Bitbucket.
	PlatformBitbucket Platform = "bitbucket"
	// PlatformUnknown represents an unknown or unsupported platform.
	PlatformUnknown Platform = "unknown"
)

// PlatformInfo contains information about a Git hosting platform.
type PlatformInfo struct {
	Platform Platform // The detected platform
	CLITool  string   // CLI tool for PR creation: "gh", "glab", "tea", or ""
	RepoURL  string   // Base URL for generating manual PR URLs
	Owner    string   // Repository owner/organization
	Repo     string   // Repository name
}

// DetectPlatform parses a Git remote URL and returns platform information.
// Supports HTTPS, SSH, and git:// protocol URLs.
func DetectPlatform(remoteURL string) (PlatformInfo, error) {
	if remoteURL == "" {
		return PlatformInfo{}, errors.New("empty remote URL")
	}

	// Normalize the URL to extract host and path
	host, path, err := parseRemoteURL(remoteURL)
	if err != nil {
		return PlatformInfo{}, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	// Extract owner and repo from path
	owner, repo := extractOwnerRepo(path)
	if owner == "" || repo == "" {
		return PlatformInfo{},
			fmt.Errorf("failed to extract owner/repo from URL path: %s", path)
	}

	// Detect platform based on hostname
	platform := detectPlatformFromHost(host)

	// Build repo URL for manual PR generation
	repoURL := buildRepoURL(host, owner, repo)

	// Determine CLI tool based on platform
	cliTool := getCLITool(platform)

	return PlatformInfo{
		Platform: platform,
		CLITool:  cliTool,
		RepoURL:  repoURL,
		Owner:    owner,
		Repo:     repo,
	}, nil
}

// GetOriginURL retrieves the URL for the 'origin' remote.
// Returns an error if not in a git repository or if no origin remote exists.
func GetOriginURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"failed to get origin URL: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git command: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// parseRemoteURL extracts the host and path from various URL formats.
// Supports:
//   - HTTPS: https://github.com/owner/repo.git
//   - SSH: git@github.com:owner/repo.git
//   - SSH with protocol: ssh://git@github.com/owner/repo.git
//   - Git protocol: git://github.com/owner/repo.git
func parseRemoteURL(url string) (host, path string, err error) {
	// SSH format: git@host:path
	sshPattern := regexp.MustCompile(`^(?:[\w-]+@)?([^:]+):(.+)$`)
	if matches := sshPattern.FindStringSubmatch(url); matches != nil {
		// Check it's not a protocol URL (has ://)
		if !strings.Contains(url, "://") {
			return matches[1], matches[2], nil
		}
	}

	// URL format: protocol://[user@]host/path
	urlPattern := regexp.MustCompile(
		`^(?:https?|ssh|git)://(?:[\w-]+@)?([^/]+)/(.+)$`,
	)
	if matches := urlPattern.FindStringSubmatch(url); matches != nil {
		return matches[1], matches[2], nil
	}

	return "", "", fmt.Errorf("unrecognized URL format: %s", url)
}

// extractOwnerRepo extracts the owner and repository name from a URL path.
// Handles paths like "owner/repo.git", "owner/repo", or
// "group/subgroup/repo.git".
func extractOwnerRepo(urlPath string) (owner, repo string) {
	// Remove .git suffix if present
	cleanPath := strings.TrimSuffix(urlPath, ".git")

	// Split by /
	parts := strings.Split(cleanPath, "/")
	if len(parts) < 2 {
		return "", ""
	}

	// For nested groups (GitLab), join all parts except last as owner
	repo = parts[len(parts)-1]
	owner = strings.Join(parts[:len(parts)-1], "/")

	return owner, repo
}

// detectPlatformFromHost determines the platform based on the hostname.
func detectPlatformFromHost(host string) Platform {
	hostLower := strings.ToLower(host)

	// GitHub
	if hostLower == "github.com" || strings.HasPrefix(hostLower, "github.") {
		return PlatformGitHub
	}

	// GitLab (includes self-hosted)
	if hostLower == "gitlab.com" || strings.Contains(hostLower, "gitlab") {
		return PlatformGitLab
	}

	// Gitea or Forgejo
	isGitea := strings.Contains(hostLower, "gitea")
	isForgejo := strings.Contains(hostLower, "forgejo")
	if isGitea || isForgejo {
		return PlatformGitea
	}

	// Bitbucket
	isBitbucketOrg := hostLower == "bitbucket.org"
	hasBitbucket := strings.Contains(hostLower, "bitbucket")
	if isBitbucketOrg || hasBitbucket {
		return PlatformBitbucket
	}

	return PlatformUnknown
}

// buildRepoURL constructs a web URL for the repository.
func buildRepoURL(host, owner, repo string) string {
	// Use HTTPS for the web URL
	return fmt.Sprintf("https://%s/%s/%s", host, owner, repo)
}

// getCLITool returns the CLI tool name for a platform.
func getCLITool(platform Platform) string {
	switch platform {
	case PlatformGitHub:
		return "gh"
	case PlatformGitLab:
		return "glab"
	case PlatformGitea:
		return "tea"
	case PlatformBitbucket, PlatformUnknown:
		// Bitbucket has no standard CLI tool
		// Unknown platforms also have no CLI tool
		return ""
	default:
		return ""
	}
}
