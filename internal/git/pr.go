package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// PROptions contains options for creating a pull request
type PROptions struct {
	Title  string
	Body   string
	Branch string
}

// CreatePR creates a pull request using the appropriate platform CLI tool
func CreatePR(platform Platform, opts PROptions) (string, error) {
	tool, err := GetCLITool(platform)
	if err != nil {
		return "", err
	}

	// Check if tool is installed
	if err := CheckCLIToolInstalled(platform); err != nil {
		return "", err
	}

	// Create PR based on platform
	switch platform {
	case PlatformGitHub:
		return createGitHubPR(tool, opts)
	case PlatformGitLab:
		return createGitLabPR(tool, opts)
	case PlatformGitea:
		return createGiteaPR(tool, opts)
	case PlatformUnknown:
		fallthrough
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
}

// createGitHubPR creates a GitHub pull request using gh CLI
func createGitHubPR(tool string, opts PROptions) (string, error) {
	cmd := exec.Command(tool, "pr", "create",
		"--title", opts.Title,
		"--body", opts.Body,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf(
			"create GitHub PR: %v\nOutput: %s",
			err,
			string(output),
		)

		return "", fmt.Errorf("%s", msg)
	}

	// Extract PR URL from output (gh prints the URL)
	url := extractURL(string(output))

	return url, nil
}

// createGitLabPR creates a GitLab merge request using glab CLI
func createGitLabPR(tool string, opts PROptions) (string, error) {
	cmd := exec.Command(tool, "mr", "create",
		"--title", opts.Title,
		"--description", opts.Body,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf(
			"create GitLab MR: %v\nOutput: %s",
			err,
			string(output),
		)

		return "", fmt.Errorf("%s", msg)
	}

	// Extract MR URL from output
	url := extractURL(string(output))

	return url, nil
}

// createGiteaPR creates a Gitea pull request using tea CLI
func createGiteaPR(tool string, opts PROptions) (string, error) {
	cmd := exec.Command(
		tool,
		"pr",
		"create",
		"--title",
		opts.Title,
		"--description",
		opts.Body,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf(
			"create Gitea PR: %v\nOutput: %s",
			err,
			string(output),
		)

		return "", fmt.Errorf("%s", msg)
	}

	// Extract PR URL from output
	url := extractURL(string(output))

	return url, nil
}

// extractURL extracts a URL from CLI tool output.
func extractURL(output string) string {
	lines := strings.SplitSeq(output, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)

		// Check if line starts with URL
		if isURLStart(line) {
			return line
		}

		// Try to find URL embedded in line
		url := findURLInLine(line)
		if url != "" {
			return url
		}
	}

	return strings.TrimSpace(output) // Fallback to full output
}

// isURLStart checks if a string starts with http:// or https://
func isURLStart(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// findURLInLine extracts a URL from a line that contains one embedded.
func findURLInLine(line string) string {
	httpPresent := strings.Contains(line, "http://")
	httpsPresent := strings.Contains(line, "https://")
	if !httpPresent && !httpsPresent {
		return ""
	}

	parts := strings.FieldsSeq(line)
	for part := range parts {
		if isURLStart(part) {
			return part
		}
	}

	return ""
}
