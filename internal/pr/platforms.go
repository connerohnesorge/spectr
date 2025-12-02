// Package pr provides platform-specific PR creation functions.
// This file handles GitHub (gh), GitLab (glab), Gitea (tea), and Bitbucket.
//
//nolint:revive // file-length-limit - platform handlers are cohesive
package pr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/connerohnesorge/spectr/internal/git"
)

// prCreateArgs bundles arguments for PR creation.
// This reduces the number of function parameters.
type prCreateArgs struct {
	branchName   string
	baseBranch   string
	title        string
	body         string
	draft        bool
	worktreePath string
}

// prResult holds the result of a PR creation operation.
// It contains the PR URL and optionally a manual URL for Bitbucket.
type prResult struct {
	prURL     string
	manualURL string
}

// createPRInput bundles all input parameters for createPR.
type createPRInput struct {
	platform     git.PlatformInfo
	branchName   string
	baseBranch   string
	title        string
	body         string
	draft        bool
	worktreePath string
}

// createPR creates a pull request using the appropriate platform CLI.
// It returns the PR URL and optionally a manual URL for Bitbucket.
// The platform parameter determines which CLI tool to use.
//
//nolint:revive // argument-limit - kept for API compatibility
func createPR(
	platform git.PlatformInfo,
	branchName, baseBranch, title, body string,
	draft bool,
	worktreePath string,
) (prURL, manualURL string, err error) {
	input := createPRInput{
		platform:     platform,
		branchName:   branchName,
		baseBranch:   baseBranch,
		title:        title,
		body:         body,
		draft:        draft,
		worktreePath: worktreePath,
	}

	return doCreatePR(input)
}

// doCreatePR is the internal implementation of createPR.
func doCreatePR(input createPRInput) (prURL, manualURL string, err error) {
	args := prCreateArgs{
		branchName:   input.branchName,
		baseBranch:   input.baseBranch,
		title:        input.title,
		body:         input.body,
		draft:        input.draft,
		worktreePath: input.worktreePath,
	}

	result, err := createPRForPlatform(input.platform, args)
	if err != nil {
		return "", "", err
	}

	return result.prURL, result.manualURL, nil
}

// createPRForPlatform dispatches to the platform-specific PR creator.
// Returns an error for unknown or unsupported platforms.
func createPRForPlatform(
	platform git.PlatformInfo,
	args prCreateArgs,
) (*prResult, error) {
	switch platform.Platform {
	case git.PlatformGitHub:
		return createGitHubPR(args)

	case git.PlatformGitLab:
		return createGitLabMR(args)

	case git.PlatformGitea:
		return createGiteaPR(args)

	case git.PlatformBitbucket:
		return createBitbucketPR(platform, args)

	case git.PlatformUnknown:
		return nil, errors.New("unknown platform; please create PR manually")
	}

	return nil, fmt.Errorf(
		"unsupported platform '%s'; please create PR manually",
		platform.Platform,
	)
}

// createGitHubPR creates a GitHub pull request using the gh CLI.
// It writes the PR body to a temp file and uses --body-file.
func createGitHubPR(args prCreateArgs) (*prResult, error) {
	fmt.Println("Creating GitHub pull request...")

	// Write PR body to temp file for gh CLI
	bodyFile, err := writeTempBodyFile(args.body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(bodyFile) }()

	cmdArgs := []string{
		"pr", "create",
		"--title", args.title,
		"--body-file", bodyFile,
		"--base", args.baseBranch,
	}

	if args.draft {
		cmdArgs = append(cmdArgs, "--draft")
	}

	cmd := exec.Command("gh", cmdArgs...)
	cmd.Dir = args.worktreePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"gh pr create failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return &prResult{
		prURL: strings.TrimSpace(string(output)),
	}, nil
}

// createGitLabMR creates a GitLab merge request using the glab CLI.
// GitLab uses --description instead of --body-file.
func createGitLabMR(args prCreateArgs) (*prResult, error) {
	fmt.Println("Creating GitLab merge request...")

	cmdArgs := []string{
		"mr", "create",
		"--title", args.title,
		"--description", args.body,
		"--target-branch", args.baseBranch,
	}

	if args.draft {
		cmdArgs = append(cmdArgs, "--draft")
	}

	cmd := exec.Command("glab", cmdArgs...)
	cmd.Dir = args.worktreePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"glab mr create failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return &prResult{
		prURL: extractURLFromOutput(string(output)),
	}, nil
}

// createGiteaPR creates a Gitea pull request using the tea CLI.
// Gitea uses --description and --base for the target branch.
func createGiteaPR(args prCreateArgs) (*prResult, error) {
	fmt.Println("Creating Gitea pull request...")

	cmdArgs := []string{
		"pr", "create",
		"--title", args.title,
		"--description", args.body,
		"--base", args.baseBranch,
	}

	cmd := exec.Command("tea", cmdArgs...)
	cmd.Dir = args.worktreePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"tea pr create failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return &prResult{
		prURL: extractURLFromOutput(string(output)),
	}, nil
}

// createBitbucketPR handles Bitbucket which has no standard CLI.
// It returns a manual URL for the user to create the PR.
func createBitbucketPR(
	platform git.PlatformInfo,
	args prCreateArgs,
) (*prResult, error) {
	manualURL := fmt.Sprintf(
		"%s/pull-requests/new?source=%s&dest=%s",
		platform.RepoURL,
		args.branchName,
		args.baseBranch,
	)

	fmt.Println()
	fmt.Println("PR creation not automated for Bitbucket.")
	fmt.Printf("Create manually at: %s\n", manualURL)

	return &prResult{
		manualURL: manualURL,
	}, nil
}

// writeTempBodyFile writes the PR body to a temporary file.
// The caller is responsible for removing the file when done.
func writeTempBodyFile(body string) (string, error) {
	bodyFile, err := os.CreateTemp("", "spectr-pr-body-*.md")
	if err != nil {
		return "", fmt.Errorf("create temp file for PR body: %w", err)
	}

	if _, err := bodyFile.WriteString(body); err != nil {
		_ = bodyFile.Close()
		_ = os.Remove(bodyFile.Name())

		return "", fmt.Errorf("write PR body: %w", err)
	}

	if err := bodyFile.Close(); err != nil {
		_ = os.Remove(bodyFile.Name())

		return "", fmt.Errorf("close PR body file: %w", err)
	}

	return bodyFile.Name(), nil
}

// extractURLFromOutput attempts to extract a URL from command output.
// It looks for lines starting with http:// or https://.
// Returns the full output if no URL is found.
func extractURLFromOutput(output string) string {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "http://") {
			return line
		}

		if strings.HasPrefix(line, "https://") {
			return line
		}
	}

	// Return the full output if no URL found
	return strings.TrimSpace(output)
}
