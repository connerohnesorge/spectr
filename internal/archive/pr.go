package archive

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/git"
)

// PRContext holds the context needed for PR creation
type PRContext struct {
	ChangeID     string
	ArchiveName  string
	SkipSpecs    bool
	OpCounts     OperationCounts
	Capabilities []string
	SpectrRoot   string
}

// createPR orchestrates the PR creation workflow after successful archive
func createPR(ctx PRContext) error {
	fmt.Println("\nCreating pull request...")

	platform, err := validateAndDetectPlatform()
	if err != nil {
		return err
	}

	baseBranchName := fmt.Sprintf("archive-%s", ctx.ChangeID)
	branchName := git.GenerateUniqueBranchName(baseBranchName)

	// Create temporary worktree directory
	tempPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("spectr-archive-%s", ctx.ChangeID),
	)

	// Ensure cleanup on exit
	defer func() {
		if err := git.RemoveWorktree(tempPath); err != nil {
			msg := "\nWarning: Failed to remove worktree: %v\n"
			fmt.Fprintf(os.Stderr, msg, err)
		}
	}()

	// Create worktree with new branch
	if err := git.CreateWorktree(tempPath, branchName); err != nil {
		return fmt.Errorf("create worktree: %w", err)
	}

	fmt.Printf("Created worktree at: %s\n", tempPath)

	// Run archive operations in the worktree
	archiveCmd := &ArchiveCmd{
		ChangeID:  ctx.ChangeID,
		SkipSpecs: ctx.SkipSpecs,
		Yes:       true,  // Non-interactive mode for worktree operations
		PR:        false, // Prevent recursive PR creation
	}

	if err := Archive(archiveCmd, tempPath); err != nil {
		return fmt.Errorf("archive in worktree: %w", err)
	}

	// Stage and commit in worktree
	err = prepareBranchAndCommit(ctx, tempPath)
	if err != nil {
		return err
	}

	// Push and create PR from worktree
	prURL, err := pushAndCreatePR(ctx, platform, branchName, tempPath)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Pull request created: %s\n", prURL)

	return nil
}

// validateAndDetectPlatform validates git environment and detects platform.
// It checks git repository status, detects the platform, and verifies CLI tool
// installation. Returns platform and error if any validation fails.
func validateAndDetectPlatform() (git.Platform, error) {
	if err := validateGitEnvironment(); err != nil {
		return git.PlatformUnknown, err
	}

	platform, remoteURL, err := git.DetectPlatform()
	if err != nil {
		msg := "detect git platform: %w. " +
			"Archive completed successfully. Create PR manually"

		return git.PlatformUnknown, fmt.Errorf(msg, err)
	}

	if platform == git.PlatformUnknown {
		msg := "could not detect git hosting platform. Remote URL: %s. " +
			"Archive completed successfully. Create PR manually"

		return git.PlatformUnknown, fmt.Errorf(msg, remoteURL)
	}

	if err := git.CheckCLIToolInstalled(platform); err != nil {
		msg := "%w. Archive completed successfully. Create PR manually"

		return git.PlatformUnknown, fmt.Errorf(msg, err)
	}

	return platform, nil
}

// prepareBranchAndCommit stages files and commits in the worktree.
// It stages the archive files and specs, then commits with an automatically
// generated message. The branch is already created by CreateWorktree.
func prepareBranchAndCommit(
	ctx PRContext,
	workingDir string,
) error {
	if err := stageArchiveFiles(ctx, workingDir); err != nil {
		return err
	}

	commitMsg := buildCommitMessage(ctx)
	if err := commitInWorktree(commitMsg, workingDir); err != nil {
		msg := "%v. Archive completed. Branch created. " +
			"Commit manually and push"

		return fmt.Errorf(msg, err)
	}

	return nil
}

// pushAndCreatePR pushes the branch and creates a pull request.
// It pushes the archive branch to origin and creates a PR/MR using the
// platform-specific CLI tool (gh, glab, or tea). Returns PR URL on success.
func pushAndCreatePR(
	ctx PRContext,
	platform git.Platform,
	branchName,
	workingDir string,
) (string, error) {
	if err := pushFromWorktree(branchName, workingDir); err != nil {
		msg := "%v. Archive completed. Branch created and committed. " +
			"Push manually"

		return "", fmt.Errorf(msg, err)
	}

	// Build PR options with title and body
	prOpts := git.PROptions{
		Title:  buildPRTitle(ctx.ChangeID),
		Body:   buildPRBody(ctx),
		Branch: branchName,
	}

	prURL, err := git.CreatePR(platform, prOpts)
	if err != nil {
		msg := "%v. Archive completed. Branch pushed. Create PR manually"

		return "", fmt.Errorf(msg, err)
	}

	return prURL, nil
}

// validateGitEnvironment checks if git is properly configured
func validateGitEnvironment() error {
	if err := git.IsGitRepository(); err != nil {
		return err
	}

	return git.HasOriginRemote()
}

// stageArchiveFiles stages the archived directory and updated specs
func stageArchiveFiles(ctx PRContext, workingDir string) error {
	// Construct paths relative to the worktree's spectr root
	worktreeSpectrRoot := filepath.Join(workingDir, "spectr")

	paths := []string{
		filepath.Join(
			worktreeSpectrRoot,
			"changes",
			"archive",
			ctx.ArchiveName,
		),
	}

	// Add specs directory if specs were updated
	if !ctx.SkipSpecs {
		paths = append(paths, filepath.Join(worktreeSpectrRoot, "specs"))
	}

	if err := stageFilesInWorktree(paths, workingDir); err != nil {
		return err
	}

	fmt.Println("Staged files for commit")

	return nil
}

// stageFilesInWorktree stages files in the specified worktree directory
func stageFilesInWorktree(paths []string, workingDir string) error {
	args := append([]string{"-C", workingDir, "add"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "stage files in worktree: %w\nOutput: %s"

		return fmt.Errorf(msg, err, string(output))
	}

	return nil
}

// commitInWorktree creates a commit in the specified worktree directory
func commitInWorktree(message, workingDir string) error {
	cmd := exec.Command("git", "-C", workingDir, "commit", "-F", "-")
	cmd.Stdin = strings.NewReader(message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "commit in worktree: %w\nOutput: %s"

		return fmt.Errorf(msg, err, string(output))
	}

	return nil
}

// pushFromWorktree pushes the branch from the specified worktree directory
func pushFromWorktree(branchName, workingDir string) error {
	cmd := exec.Command(
		"git",
		"-C",
		workingDir,
		"push",
		"-u",
		"origin",
		branchName,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := "push from worktree: %w\nOutput: %s"

		return fmt.Errorf(msg, err, string(output))
	}

	return nil
}
