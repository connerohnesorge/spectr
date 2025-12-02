//nolint:revive // file-length-limit - logically cohesive workflow functions
package pr

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/git"
)

// File permission constants.
const (
	dirPerm  = 0o755
	filePerm = 0o644
	gitCmd   = "git"
)

// Config contains configuration for PR workflow execution.
type Config struct {
	// ChangeID is the identifier of the change to process.
	ChangeID string
	// BaseBranch is the target branch for the PR (auto-detected if empty).
	BaseBranch string
	// Draft creates the PR as a draft when true.
	Draft bool
	// Force overwrites existing branch if true.
	Force bool
	// DryRun prints what would be done without executing.
	DryRun bool
	// SkipSpecs skips spec updates during archive (archive only).
	SkipSpecs bool
	// WorkingDir is the working directory for operations.
	WorkingDir string
}

// Result contains the outcome of a PR workflow execution.
type Result struct {
	// PRURL is the URL of the created pull request.
	PRURL string
	// BranchName is the name of the created branch.
	BranchName string
	// ArchivePath is the path to the archived change (archive mode only).
	ArchivePath string
	// OperationCounts tracks spec operations applied (archive mode only).
	OperationCounts archive.OperationCounts
}

// Error definitions for PR workflow.
var (
	// ErrNoCLI is returned when the required CLI tool is not installed.
	ErrNoCLI = errors.New("required CLI tool is not installed")
	// ErrNotAuthenticated is returned when CLI tool is not authenticated.
	ErrNotAuthenticated = errors.New("CLI tool is not authenticated")
	// ErrBranchExists is returned when the target branch already exists.
	ErrBranchExists = errors.New(
		"branch already exists (use --force to overwrite)",
	)
	// ErrChangeNotFound is returned when the specified change does not exist.
	ErrChangeNotFound = errors.New("change not found")
	// ErrNoPRSupport is returned when platform doesn't support CLI PR creation.
	ErrNoPRSupport = errors.New("platform does not support CLI PR creation")
)

// ExecuteArchivePR executes the full archive PR workflow.
// It creates a worktree, runs spectr archive, commits, pushes, and creates PR.
//
//nolint:revive // function-length - workflow function with many sequential steps
func ExecuteArchivePR(cfg Config) (*Result, error) {
	platformInfo, baseBranch, branchName, err := prepareArchiveWorkflow(cfg)
	if err != nil {
		return nil, err
	}

	// Dry run mode - print what would be done
	if cfg.DryRun {
		return executeDryRunArchive(cfg, baseBranch, branchName, platformInfo)
	}

	return executeArchiveInWorktree(cfg, platformInfo, baseBranch, branchName)
}

// prepareArchiveWorkflow validates prerequisites and prepares workflow config.
func prepareArchiveWorkflow(
	cfg Config,
) (git.PlatformInfo, string, string, error) {
	// Validate prerequisites
	if err := validatePrereqs(); err != nil {
		return git.PlatformInfo{}, "", "",
			fmt.Errorf("prerequisites check failed: %w", err)
	}

	// Detect platform from origin URL
	remoteURL, err := git.GetOriginRemote()
	if err != nil {
		return git.PlatformInfo{}, "", "",
			fmt.Errorf("get origin remote: %w", err)
	}
	platformInfo := git.DetectPlatform(remoteURL)

	// Validate CLI availability for platform
	err = validatePlatformCLI(platformInfo)
	if err != nil {
		return git.PlatformInfo{}, "", "", err
	}

	// Get or validate base branch
	baseBranch, err := resolveBaseBranch(cfg.BaseBranch)
	if err != nil {
		return git.PlatformInfo{}, "", "", err
	}

	// Generate branch name and check existence
	branchName := fmt.Sprintf("spectr/archive/%s", cfg.ChangeID)
	if err := handleExistingBranch(branchName, cfg.Force, cfg.DryRun); err != nil {
		return git.PlatformInfo{}, "", "", err
	}

	return platformInfo, baseBranch, branchName, nil
}

// executeArchiveInWorktree runs archive workflow in a git worktree.
func executeArchiveInWorktree(
	cfg Config,
	platformInfo git.PlatformInfo,
	baseBranch, branchName string,
) (*Result, error) {
	worktree, err := git.CreateWorktree(baseBranch, branchName)
	if err != nil {
		return nil, fmt.Errorf("create worktree: %w", err)
	}
	defer cleanupWorktreeSafe(worktree.Path)

	// Execute spectr archive in worktree
	archiveCmd := buildArchiveCommand(cfg.ChangeID, cfg.SkipSpecs)
	archiveCmd.Dir = worktree.Path
	output, err := archiveCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"spectr archive failed: %s: %w", string(output), err,
		)
	}

	// Stage and commit
	err = stageSpectrDir(worktree.Path)
	if err != nil {
		return nil, fmt.Errorf("stage changes: %w", err)
	}

	archivePath := fmt.Sprintf("spectr/changes/archive/%s", cfg.ChangeID)
	commitData := ArchiveCommitData{
		ChangeID:    cfg.ChangeID,
		ArchivePath: archivePath,
	}
	commitMsg, err := RenderArchiveCommitMessage(commitData)
	if err != nil {
		return nil, fmt.Errorf("render commit message: %w", err)
	}

	err = createCommit(worktree.Path, commitMsg)
	if err != nil {
		return nil, fmt.Errorf("create commit: %w", err)
	}

	// Push and create PR
	err = pushBranch(worktree.Path, branchName)
	if err != nil {
		return nil, fmt.Errorf("push branch: %w", err)
	}

	prURL, err := createArchivePR(cfg, platformInfo, archivePath, baseBranch)
	if err != nil {
		return nil, err
	}

	return &Result{
		PRURL:       prURL,
		BranchName:  branchName,
		ArchivePath: archivePath,
	}, nil
}

// createArchivePR creates the PR for archive workflow.
func createArchivePR(
	cfg Config,
	platformInfo git.PlatformInfo,
	archivePath, baseBranch string,
) (string, error) {
	prTitle := fmt.Sprintf("spectr(archive): %s", cfg.ChangeID)
	prBodyData := ArchivePRBodyData{
		ChangeID:     cfg.ChangeID,
		ArchivePath:  archivePath,
		Capabilities: nil,
	}
	prBody, err := RenderArchivePRBody(prBodyData)
	if err != nil {
		return "", fmt.Errorf("render PR body: %w", err)
	}

	prURL, err := createPR(platformInfo, prTitle, prBody, baseBranch, cfg.Draft)
	if err != nil {
		return "", fmt.Errorf("create PR: %w", err)
	}

	return prURL, nil
}

// ExecuteNewPR executes the full new proposal PR workflow.
// It creates a worktree, copies the change, commits, pushes, and creates PR.
//
//nolint:revive // function-length - workflow function with many sequential steps
func ExecuteNewPR(cfg Config) (*Result, error) {
	platformInfo, baseBranch, branchName, changePath, err := prepareNewWorkflow(cfg)
	if err != nil {
		return nil, err
	}

	// Dry run mode - print what would be done
	if cfg.DryRun {
		return executeDryRunNew(cfg, baseBranch, branchName, platformInfo)
	}

	return executeNewInWorktree(
		cfg, platformInfo, baseBranch, branchName, changePath,
	)
}

// prepareNewWorkflow validates prerequisites and prepares workflow config.
func prepareNewWorkflow(
	cfg Config,
) (git.PlatformInfo, string, string, string, error) {
	// Validate prerequisites
	if err := validatePrereqs(); err != nil {
		return git.PlatformInfo{}, "", "", "",
			fmt.Errorf("prerequisites check failed: %w", err)
	}

	// Detect platform from origin URL
	remoteURL, err := git.GetOriginRemote()
	if err != nil {
		return git.PlatformInfo{}, "", "", "",
			fmt.Errorf("get origin remote: %w", err)
	}
	platformInfo := git.DetectPlatform(remoteURL)

	// Validate CLI availability for platform
	if err := validatePlatformCLI(platformInfo); err != nil {
		return git.PlatformInfo{}, "", "", "", err
	}

	// Validate change exists
	changePath, err := resolveChangePath(cfg)
	if err != nil {
		return git.PlatformInfo{}, "", "", "", err
	}

	// Get or validate base branch
	baseBranch, err := resolveBaseBranch(cfg.BaseBranch)
	if err != nil {
		return git.PlatformInfo{}, "", "", "", err
	}

	// Generate branch name and check existence
	branchName := fmt.Sprintf("spectr/proposal/%s", cfg.ChangeID)
	if err := handleExistingBranch(branchName, cfg.Force, cfg.DryRun); err != nil {
		return git.PlatformInfo{}, "", "", "", err
	}

	return platformInfo, baseBranch, branchName, changePath, nil
}

// executeNewInWorktree runs new proposal workflow in a git worktree.
func executeNewInWorktree(
	cfg Config,
	platformInfo git.PlatformInfo,
	baseBranch, branchName, changePath string,
) (*Result, error) {
	worktree, err := git.CreateWorktree(baseBranch, branchName)
	if err != nil {
		return nil, fmt.Errorf("create worktree: %w", err)
	}
	defer cleanupWorktreeSafe(worktree.Path)

	// Copy change directory to worktree
	destPath := filepath.Join(
		worktree.Path, "spectr", "changes", cfg.ChangeID,
	)
	err = copyDir(changePath, destPath)
	if err != nil {
		return nil, fmt.Errorf("copy change directory: %w", err)
	}

	// Stage and commit
	err = stageSpectrDir(worktree.Path)
	if err != nil {
		return nil, fmt.Errorf("stage changes: %w", err)
	}

	proposalPath := fmt.Sprintf("spectr/changes/%s", cfg.ChangeID)
	commitData := NewCommitData{
		ChangeID:     cfg.ChangeID,
		ProposalPath: proposalPath,
	}
	commitMsg, err := RenderNewCommitMessage(commitData)
	if err != nil {
		return nil, fmt.Errorf("render commit message: %w", err)
	}

	err = createCommit(worktree.Path, commitMsg)
	if err != nil {
		return nil, fmt.Errorf("create commit: %w", err)
	}

	// Push and create PR
	err = pushBranch(worktree.Path, branchName)
	if err != nil {
		return nil, fmt.Errorf("push branch: %w", err)
	}

	prURL, err := createNewPR(cfg, platformInfo, proposalPath, baseBranch)
	if err != nil {
		return nil, err
	}

	return &Result{
		PRURL:      prURL,
		BranchName: branchName,
	}, nil
}

// createNewPR creates the PR for new proposal workflow.
func createNewPR(
	cfg Config,
	platformInfo git.PlatformInfo,
	proposalPath, baseBranch string,
) (string, error) {
	prTitle := fmt.Sprintf("spectr(proposal): %s", cfg.ChangeID)
	prBodyData := NewPRBodyData{
		ChangeID:     cfg.ChangeID,
		ProposalPath: proposalPath,
	}
	prBody, err := RenderNewPRBody(prBodyData)
	if err != nil {
		return "", fmt.Errorf("render PR body: %w", err)
	}

	prURL, err := createPR(platformInfo, prTitle, prBody, baseBranch, cfg.Draft)
	if err != nil {
		return "", fmt.Errorf("create PR: %w", err)
	}

	return prURL, nil
}

// validatePrereqs checks git version and origin remote existence.
func validatePrereqs() error {
	if err := git.CheckGitVersion(); err != nil {
		return err
	}

	_, err := git.GetOriginRemote()

	return err
}

// validatePlatformCLI validates CLI availability for the platform.
func validatePlatformCLI(platform git.PlatformInfo) error {
	if platform.Platform == git.PlatformBitbucket ||
		platform.Platform == git.PlatformUnknown {
		return nil
	}

	if err := checkCLIAvailable(platform.CLITool); err != nil {
		return err
	}

	return checkCLIAuthenticated(platform.CLITool)
}

// resolveBaseBranch gets or validates the base branch.
func resolveBaseBranch(configured string) (string, error) {
	if configured != "" {
		return configured, nil
	}

	baseBranch, err := git.GetBaseBranch()
	if err != nil {
		return "", fmt.Errorf("detect base branch: %w", err)
	}

	return baseBranch, nil
}

// resolveChangePath validates and returns the change path.
func resolveChangePath(cfg Config) (string, error) {
	workingDir := cfg.WorkingDir
	if workingDir == "" {
		var err error
		workingDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
	}

	changePath := filepath.Join(workingDir, "spectr", "changes", cfg.ChangeID)
	if _, err := os.Stat(changePath); os.IsNotExist(err) {
		return "", ErrChangeNotFound
	}

	return changePath, nil
}

// handleExistingBranch checks if branch exists and handles accordingly.
func handleExistingBranch(branchName string, force, dryRun bool) error {
	exists, err := git.CheckBranchExists(branchName)
	if err != nil {
		return fmt.Errorf("check branch exists: %w", err)
	}

	if !exists {
		return nil
	}

	if !force {
		return ErrBranchExists
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would delete existing branch: %s\n", branchName)

		return nil
	}

	return deleteBranch(branchName)
}

// cleanupWorktreeSafe removes worktree with error logging.
func cleanupWorktreeSafe(path string) {
	if err := git.CleanupWorktree(path); err != nil {
		fmt.Fprintf(
			os.Stderr, "Warning: failed to cleanup worktree: %v\n", err,
		)
	}
}

// checkCLIAvailable checks if the specified CLI tool is installed.
func checkCLIAvailable(tool string) error {
	if tool == "" {
		return ErrNoCLI
	}

	_, err := exec.LookPath(tool)
	if err != nil {
		return fmt.Errorf("%w: %s not found in PATH", ErrNoCLI, tool)
	}

	return nil
}

// checkCLIAuthenticated checks if the CLI tool is authenticated.
func checkCLIAuthenticated(tool string) error {
	var cmd *exec.Cmd

	switch tool {
	case "gh":
		cmd = exec.Command("gh", "auth", "status")
	case "glab":
		cmd = exec.Command("glab", "auth", "status")
	case "tea":
		cmd = exec.Command("tea", "login", "list")
	default:
		// For unknown tools, skip auth check
		return nil
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"%w: %s is not authenticated", ErrNotAuthenticated, tool,
		)
	}

	return nil
}

// createPR creates a PR using the platform-specific CLI.
//
//nolint:revive // function-length - switch for different platforms
func createPR(
	platform git.PlatformInfo,
	title, body, baseBranch string,
	draft bool,
) (string, error) {
	// Write body to temp file for CLI
	bodyFile, err := os.CreateTemp("", "spectr-pr-body-*.md")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer func() { _ = os.Remove(bodyFile.Name()) }()

	_, err = bodyFile.WriteString(body)
	if err != nil {
		return "", fmt.Errorf("write body file: %w", err)
	}
	err = bodyFile.Close()
	if err != nil {
		return "", fmt.Errorf("close body file: %w", err)
	}

	cmd, err := buildPRCommand(platform, title, bodyFile.Name(), baseBranch, draft)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"CLI command failed: %s: %w", string(output), err,
		)
	}

	// Parse PR URL from output
	prURL := parsePRURL(string(output))
	if prURL == "" {
		// If we can't parse URL, return the full output
		return strings.TrimSpace(string(output)), nil
	}

	return prURL, nil
}

// buildPRCommand constructs the platform-specific PR creation command.
func buildPRCommand(
	platform git.PlatformInfo,
	title, bodyFilePath, baseBranch string,
	draft bool,
) (*exec.Cmd, error) {
	switch platform.Platform {
	case git.PlatformGitHub:
		return buildGitHubPRCommand(title, bodyFilePath, baseBranch, draft), nil

	case git.PlatformGitLab:
		return buildGitLabMRCommand(title, bodyFilePath, baseBranch, draft)

	case git.PlatformGitea:
		return buildGiteaPRCommand(title, bodyFilePath, baseBranch)

	case git.PlatformBitbucket, git.PlatformUnknown:
		return nil, ErrNoPRSupport
	}

	return nil, ErrNoPRSupport
}

// buildGitHubPRCommand constructs gh pr create command.
func buildGitHubPRCommand(
	title, bodyFilePath, baseBranch string,
	draft bool,
) *exec.Cmd {
	args := []string{
		"pr", "create",
		"--title", title,
		"--body-file", bodyFilePath,
		"--base", baseBranch,
	}
	if draft {
		args = append(args, "--draft")
	}

	return exec.Command("gh", args...)
}

// buildGitLabMRCommand constructs glab mr create command.
func buildGitLabMRCommand(
	title, bodyFilePath, baseBranch string,
	draft bool,
) (*exec.Cmd, error) {
	bodyContent, err := os.ReadFile(bodyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read body file: %w", err)
	}
	args := []string{
		"mr", "create",
		"--title", title,
		"--description", string(bodyContent),
		"--target-branch", baseBranch,
	}
	if draft {
		args = append(args, "--draft")
	}

	return exec.Command("glab", args...), nil
}

// buildGiteaPRCommand constructs tea pr create command.
func buildGiteaPRCommand(
	title, bodyFilePath, baseBranch string,
) (*exec.Cmd, error) {
	bodyContent, err := os.ReadFile(bodyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read body file: %w", err)
	}
	args := []string{
		"pr", "create",
		"--title", title,
		"--description", string(bodyContent),
		"--base", baseBranch,
	}

	return exec.Command("tea", args...), nil
}

// deleteBranch deletes a remote branch.
func deleteBranch(branch string) error {
	cmd := exec.Command(
		gitCmd, "push", "origin", "--delete", branch,
	) //nolint:gosec
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to delete remote branch %q (branch may still exist locally): %w\nOutput: %s",
			branch, err, strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// stageSpectrDir stages the spectr/ directory for commit.
func stageSpectrDir(worktreePath string) error {
	cmd := exec.Command(gitCmd, "add", "spectr/")
	cmd.Dir = worktreePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add: %s: %w", string(output), err)
	}

	return nil
}

// createCommit creates a git commit with the specified message.
func createCommit(worktreePath, message string) error {
	cmd := exec.Command(gitCmd, "commit", "-m", message)
	cmd.Dir = worktreePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit: %s: %w", string(output), err)
	}

	return nil
}

// pushBranch pushes a branch to origin.
func pushBranch(worktreePath, branch string) error {
	cmd := exec.Command(
		gitCmd, "push", "-u", "origin", branch,
	) //nolint:gosec
	cmd.Dir = worktreePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push: %s: %w", string(output), err)
	}

	return nil
}

// buildArchiveCommand builds the spectr archive command.
func buildArchiveCommand(changeID string, skipSpecs bool) *exec.Cmd {
	args := []string{"archive", changeID, "--yes"}
	if skipSpecs {
		args = append(args, "--skip-specs")
	}

	// Use the current executable to ensure we're running the same version
	executable, err := os.Executable()
	if err != nil {
		// Fallback to spectr in PATH
		executable = "spectr"
	}

	return exec.Command(executable, args...) //nolint:gosec
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(dst, dirPerm); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Handle symlinks
		if entry.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(srcPath)
			if err != nil {
				return fmt.Errorf("read symlink %s: %w", srcPath, err)
			}
			if err := os.Symlink(target, dstPath); err != nil {
				return fmt.Errorf("create symlink %s: %w", dstPath, err)
			}

			continue
		}

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, content, filePerm)
}

// parsePRURL extracts a URL from CLI output.
func parsePRURL(output string) string {
	// Common patterns for PR URLs
	patterns := []string{
		`https://github\.com/[^\s]+/pull/\d+`,
		`https://gitlab\.com/[^\s]+/-/merge_requests/\d+`,
		`https://[^\s]+/pulls/\d+`, // Gitea pattern
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(output); match != "" {
			return match
		}
	}

	return ""
}

// executeDryRunArchive prints what would be done for archive workflow.
func executeDryRunArchive(
	cfg Config,
	baseBranch, branchName string,
	platform git.PlatformInfo,
) (*Result, error) {
	fmt.Println("[DRY RUN] Archive PR workflow:")
	fmt.Printf(
		"  Platform: %s (CLI: %s)\n", platform.Platform, platform.CLITool,
	)
	fmt.Printf("  Base branch: %s\n", baseBranch)
	fmt.Printf("  New branch: %s\n", branchName)
	fmt.Printf("  Change ID: %s\n", cfg.ChangeID)
	fmt.Println()
	fmt.Println("Would execute:")
	fmt.Println("  1. Create worktree based on origin/" + baseBranch)
	fmt.Printf("  2. Run: spectr archive %s --yes", cfg.ChangeID)
	if cfg.SkipSpecs {
		fmt.Print(" --skip-specs")
	}
	fmt.Println()
	fmt.Println("  3. Stage spectr/ directory")
	fmt.Println("  4. Create commit with archive message")
	fmt.Println("  5. Push branch to origin")
	fmt.Printf("  6. Create PR via %s\n", platform.CLITool)
	fmt.Println("  7. Cleanup worktree")

	return &Result{
		PRURL:       "[dry-run]",
		BranchName:  branchName,
		ArchivePath: fmt.Sprintf("spectr/changes/archive/%s", cfg.ChangeID),
	}, nil
}

// executeDryRunNew prints what would be done for new proposal workflow.
func executeDryRunNew(
	cfg Config,
	baseBranch, branchName string,
	platform git.PlatformInfo,
) (*Result, error) {
	fmt.Println("[DRY RUN] New proposal PR workflow:")
	fmt.Printf(
		"  Platform: %s (CLI: %s)\n", platform.Platform, platform.CLITool,
	)
	fmt.Printf("  Base branch: %s\n", baseBranch)
	fmt.Printf("  New branch: %s\n", branchName)
	fmt.Printf("  Change ID: %s\n", cfg.ChangeID)
	fmt.Println()
	fmt.Println("Would execute:")
	fmt.Println("  1. Create worktree based on origin/" + baseBranch)
	fmt.Printf("  2. Copy spectr/changes/%s to worktree\n", cfg.ChangeID)
	fmt.Println("  3. Stage spectr/ directory")
	fmt.Println("  4. Create commit with proposal message")
	fmt.Println("  5. Push branch to origin")
	fmt.Printf("  6. Create PR via %s\n", platform.CLITool)
	fmt.Println("  7. Cleanup worktree")

	return &Result{
		PRURL:      "[dry-run]",
		BranchName: branchName,
	}, nil
}
