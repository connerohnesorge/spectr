// Package pr provides PR workflow orchestration for creating pull requests
// from spectr changes using git worktrees for isolation.
//
//nolint:revive // file-length-limit - workflow orchestration is logically cohesive
package pr

import (
	"fmt"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/git"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// PRConfig contains configuration for the PR workflow.
type PRConfig struct {
	ChangeID    string // The change ID to create PR for
	Mode        string // "archive" or "proposal"
	BaseBranch  string // Target branch for PR (optional, auto-detect if empty)
	Draft       bool   // Create as draft PR
	Force       bool   // Delete existing remote branch if present
	DryRun      bool   // Show what would be done without executing
	SkipSpecs   bool   // For archive mode: pass --skip-specs to archive command
	ProjectRoot string // Project root directory (for source change)
}

// PRResult contains the result of the PR workflow.
type PRResult struct {
	PRURL        string                  // URL of created PR
	BranchName   string                  // Branch name created
	ArchivePath  string                  // Archive path (archive mode only)
	Counts       archive.OperationCounts // Operation counts (archive mode only)
	Capabilities []string                // Updated capabilities (archive mode only)
	Platform     git.Platform            // Detected platform
	ManualURL    string                  // Manual PR creation URL (Bitbucket)
}

// ExecutePR orchestrates the complete PR workflow:
// 1. Validate prerequisites
// 2. Create worktree on new branch
// 3. Execute archive or copy operation
// 4. Stage, commit, push
// 5. Create PR
// 6. Cleanup worktree
func ExecutePR(
	config PRConfig,
) (*PRResult, error) {
	// Validate prerequisites
	if err := validatePrerequisites(config); err != nil {
		return nil, fmt.Errorf(
			"prerequisite check failed: %w",
			err,
		)
	}

	// Prepare workflow context
	ctx, err := prepareWorkflowContext(config)
	if err != nil {
		return nil, err
	}

	if config.DryRun {
		return executeDryRun(config, ctx)
	}

	return executeWorkflow(config, ctx)
}

// workflowContext holds prepared context for the PR workflow.
type workflowContext struct {
	platformInfo git.PlatformInfo
	baseBranch   string
	branchName   string
}

// prepareWorkflowContext prepares the context needed for the workflow.
func prepareWorkflowContext(
	config PRConfig,
) (*workflowContext, error) {
	// Get origin URL and detect platform
	originURL, err := git.GetOriginURL()
	if err != nil {
		return nil, fmt.Errorf(
			"get origin URL: %w",
			err,
		)
	}

	platformInfo, err := git.DetectPlatform(
		originURL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"detect platform: %w",
			err,
		)
	}

	// Check CLI tool availability (skip for Bitbucket which has no CLI)
	if platformInfo.CLITool != "" {
		if err := checkCLITool(platformInfo.CLITool); err != nil {
			return nil, err
		}
	}

	// Get base branch (auto-detect or use provided)
	baseBranch, err := git.GetBaseBranch(
		config.BaseBranch,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"determine base branch: %w",
			err,
		)
	}

	// Generate mode-specific branch name:
	// - archive mode: spectr/archive/<change-id>
	// - proposal mode: spectr/proposal/<change-id>
	// - remove mode: spectr/remove/<change-id>
	var branchPrefix string
	switch config.Mode {
	case ModeArchive:
		branchPrefix = "spectr/archive"
	case ModeProposal:
		branchPrefix = "spectr/proposal"
	case ModeRemove:
		branchPrefix = "spectr/remove"
	default:
		branchPrefix = "spectr"
	}
	branchName := fmt.Sprintf(
		"%s/%s",
		branchPrefix,
		config.ChangeID,
	)

	// Handle existing branch
	if err := handleExistingBranch(config, branchName); err != nil {
		return nil, err
	}

	// Fetch origin to ensure refs are up to date
	if err := fetchOrigin(config); err != nil {
		return nil, err
	}

	return &workflowContext{
		platformInfo: platformInfo,
		baseBranch:   baseBranch,
		branchName:   branchName,
	}, nil
}

// handleExistingBranch handles the case where the branch already exists.
func handleExistingBranch(
	config PRConfig,
	branchName string,
) error {
	exists, err := git.BranchExists(branchName)
	if err != nil {
		return fmt.Errorf(
			"check branch existence: %w",
			err,
		)
	}

	if !exists {
		return nil
	}

	if !config.Force {
		return fmt.Errorf(
			"branch '%s' already exists on remote; use --force to delete",
			branchName,
		)
	}

	if config.DryRun {
		fmt.Printf(
			"[dry-run] Would delete remote branch: %s\n",
			branchName,
		)

		return nil
	}

	fmt.Printf(
		"Deleting existing remote branch: %s\n",
		branchName,
	)

	return git.DeleteRemoteBranch(branchName)
}

// fetchOrigin fetches the origin remote.
func fetchOrigin(config PRConfig) error {
	if config.DryRun {
		fmt.Println(
			"[dry-run] Would fetch origin",
		)

		return nil
	}

	fmt.Println("Fetching origin...")

	if err := git.FetchOrigin(); err != nil {
		return fmt.Errorf("fetch origin: %w", err)
	}

	return nil
}

// executeWorkflow executes the main PR workflow.
func executeWorkflow(
	config PRConfig,
	ctx *workflowContext,
) (*PRResult, error) {
	result := &PRResult{
		BranchName: ctx.branchName,
		Platform:   ctx.platformInfo.Platform,
	}

	// Create worktree
	fmt.Printf(
		"Creating worktree on branch: %s (based on %s)\n",
		ctx.branchName,
		ctx.baseBranch,
	)

	worktreeInfo, err := git.CreateWorktree(
		git.WorktreeConfig{
			BranchName: ctx.branchName,
			BaseBranch: ctx.baseBranch,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create worktree: %w",
			err,
		)
	}

	// Ensure cleanup happens
	defer cleanupWorktree(worktreeInfo)

	// Execute operation in worktree
	if err := executeOperation(config, worktreeInfo.Path, result); err != nil {
		return nil, err
	}

	// Stage, commit, and push
	if err := commitAndPush(config, ctx, result, worktreeInfo.Path); err != nil {
		return nil, err
	}

	// Create PR
	result, err = createPRAndFinalize(
		config,
		ctx,
		result,
		worktreeInfo.Path,
	)
	if err != nil {
		return nil, err
	}

	// Clean up local change directory for archive and remove modes
	// (not for proposal mode, as the user may still be working on the proposal)
	if config.Mode == ModeArchive ||
		config.Mode == ModeRemove {
		fmt.Printf(
			"Cleaning up local change directory: spectr/changes/%s/\n",
			config.ChangeID,
		)
		if err := cleanupLocalChange(config); err != nil {
			fmt.Printf(
				"Warning: local cleanup failed: %v\n",
				err,
			)
		}
	}

	return result, nil
}

// cleanupWorktree cleans up the worktree.
func cleanupWorktree(info *git.WorktreeInfo) {
	fmt.Println("Cleaning up worktree...")

	if err := git.CleanupWorktree(info); err != nil {
		fmt.Printf(
			"Warning: worktree cleanup failed: %v\n",
			err,
		)
	}
}

// executeOperation executes the archive or copy operation.
func executeOperation(
	config PRConfig,
	worktreePath string,
	result *PRResult,
) error {
	switch config.Mode {
	case ModeArchive:
		archiveResult, err := executeArchiveInWorktree(
			config,
			worktreePath,
		)
		if err != nil {
			return fmt.Errorf(
				"archive operation failed: %w",
				err,
			)
		}
		result.ArchivePath = archiveResult.ArchivePath
		result.Counts = archiveResult.Counts
		result.Capabilities = archiveResult.Capabilities

	case ModeProposal:
		if err := copyChangeToWorktree(config, worktreePath); err != nil {
			return fmt.Errorf(
				"copy operation failed: %w",
				err,
			)
		}

	case ModeRemove:
		// Copy change to worktree first so git can track the deletion
		if err := copyChangeToWorktree(config, worktreePath); err != nil {
			return fmt.Errorf(
				"copy operation failed: %w",
				err,
			)
		}
		// Remove the change directory
		if err := removeChangeInWorktree(config, worktreePath); err != nil {
			return fmt.Errorf(
				"remove operation failed: %w",
				err,
			)
		}

	default:
		return fmt.Errorf(
			"unknown mode: %s",
			config.Mode,
		)
	}

	return nil
}

// commitAndPush stages, commits, and pushes the changes.
func commitAndPush(
	config PRConfig,
	ctx *workflowContext,
	result *PRResult,
	worktreePath string,
) error {
	// Generate commit message
	commitData := CommitTemplateData{
		ChangeID:    config.ChangeID,
		ArchivePath: result.ArchivePath,
		Mode:        config.Mode,
		Counts:      result.Counts,
	}

	commitMsg, err := RenderCommitMessage(
		commitData,
	)
	if err != nil {
		return fmt.Errorf(
			"render commit message: %w",
			err,
		)
	}

	// Stage and commit
	if err := stageAndCommit(worktreePath, commitMsg); err != nil {
		return fmt.Errorf(
			"stage and commit: %w",
			err,
		)
	}

	// Push branch
	if err := pushBranch(worktreePath, ctx.branchName); err != nil {
		return fmt.Errorf("push branch: %w", err)
	}

	return nil
}

// createPRAndFinalize creates the PR and finalizes the result.
func createPRAndFinalize(
	config PRConfig,
	ctx *workflowContext,
	result *PRResult,
	worktreePath string,
) (*PRResult, error) {
	// Generate PR body
	prData := PRTemplateData{
		ChangeID:     config.ChangeID,
		ArchivePath:  result.ArchivePath,
		Capabilities: result.Capabilities,
		Mode:         config.Mode,
		Counts:       result.Counts,
	}

	prBody, err := RenderPRBody(prData)
	if err != nil {
		return nil, fmt.Errorf(
			"render PR body: %w",
			err,
		)
	}

	prTitle := GetPRTitle(
		config.ChangeID,
		config.Mode,
	)
	baseBranchName := strings.TrimPrefix(
		ctx.baseBranch,
		"origin/",
	)

	// Create PR
	prURL, manualURL, err := createPR(
		ctx.platformInfo,
		ctx.branchName,
		baseBranchName,
		prTitle,
		prBody,
		config.Draft,
		worktreePath,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create PR: %w",
			err,
		)
	}

	result.PRURL = prURL
	result.ManualURL = manualURL

	return result, nil
}

// validatePrerequisites checks all prerequisites before starting the workflow.
func validatePrerequisites(
	config PRConfig,
) error {
	// Check we're in a git repository
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf(
			"not in a git repository: %w",
			err,
		)
	}

	// Use provided project root or detected repo root
	projectRoot := config.ProjectRoot
	if projectRoot == "" {
		projectRoot = repoRoot
	}

	// Check origin remote exists
	_, err = git.GetOriginURL()
	if err != nil {
		return &specterrs.PRPrerequisiteError{
			Check:   "origin remote",
			Details: "no 'origin' remote found; configure remote before creating PR",
			Err:     err,
		}
	}

	// Check change exists
	changes, err := discovery.GetActiveChangeIDs(
		projectRoot,
	)
	if err != nil {
		return fmt.Errorf("list changes: %w", err)
	}

	found := false
	for _, change := range changes {
		if change == config.ChangeID {
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf(
			"change '%s' not found in spectr/changes/",
			config.ChangeID,
		)
	}

	// Check mode is valid
	if config.Mode != ModeArchive &&
		config.Mode != ModeProposal &&
		config.Mode != ModeRemove {
		return fmt.Errorf(
			"invalid mode '%s'; must be 'archive', 'proposal', or 'remove'",
			config.Mode,
		)
	}

	return nil
}
