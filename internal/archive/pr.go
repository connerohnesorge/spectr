package archive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/conneroisu/spectr/internal/git"
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

	originalBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	branchName := fmt.Sprintf("archive-%s", ctx.ChangeID)
	err = prepareBranchAndCommit(ctx, branchName)
	if err != nil {
		return err
	}

	prURL, err := pushAndCreatePR(ctx, platform, branchName)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Pull request created: %s\n", prURL)

	restoreOriginalBranch(originalBranch, branchName, ctx)

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

// prepareBranchAndCommit creates a branch, stages files, and commits.
// It creates a new git branch, stages the archive files and specs,
// then commits with an automatically generated message.
func prepareBranchAndCommit(ctx PRContext, branchName string) error {
	if err := createBranch(branchName); err != nil {
		return err
	}

	if err := stageArchiveFiles(ctx); err != nil {
		return err
	}

	commitMsg := buildCommitMessage(ctx)
	if err := git.Commit(commitMsg); err != nil {
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
	branchName string,
) (string, error) {
	if err := git.Push(branchName); err != nil {
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

// restoreOriginalBranch attempts to restore the original branch and working
// directory.
//
// After PR creation, this function returns to the branch the user was on
// before the archive and restores the changes directory (and optionally specs).
// If it fails, a warning is printed but no error returned.
func restoreOriginalBranch(originalBranch, branchName string, ctx PRContext) {
	if err := git.CheckoutBranch(originalBranch); err != nil {
		msg := "\nWarning: Failed to restore original branch '%s': %v\n" +
			"You are still on branch '%s'. " +
			"Checkout manually with: git checkout %s\n"
		fmt.Fprintf(
			os.Stderr,
			msg,
			originalBranch,
			err,
			branchName,
			originalBranch,
		)

		return
	}

	fmt.Printf("Restored original branch: %s\n", originalBranch)

	// Restore the changes directory
	changesPath := filepath.Join(ctx.SpectrRoot, "changes")
	if err := git.RestorePath(changesPath); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"\nWarning: Failed to restore changes directory: %v\n",
			err,
		)
	} else {
		fmt.Println("Restored changes directory")
	}

	// Restore specs directory if specs were updated
	if ctx.SkipSpecs {
		return
	}

	specsPath := filepath.Join(ctx.SpectrRoot, "specs")
	if err := git.RestorePath(specsPath); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"\nWarning: Failed to restore specs directory: %v\n",
			err,
		)
	} else {
		fmt.Println("Restored specs directory")
	}
}

// validateGitEnvironment checks if git is properly configured
func validateGitEnvironment() error {
	if err := git.IsGitRepository(); err != nil {
		return err
	}

	return git.HasOriginRemote()
}

// createBranch creates a new git branch for the archive.
// It checks for existing branches with the same name and creates a new
// branch with the pattern "archive-{change-id}".
func createBranch(branchName string) error {
	// Check if branch already exists
	if git.BranchExists(branchName) {
		msg := "branch '%s' already exists. " +
			"Delete it first with: git branch -D %s"

		return fmt.Errorf(msg, branchName, branchName)
	}

	if err := git.CreateBranch(branchName); err != nil {
		return err
	}

	fmt.Printf("Created branch: %s\n", branchName)

	return nil
}

// stageArchiveFiles stages the archived directory and updated specs
func stageArchiveFiles(ctx PRContext) error {
	paths := []string{
		filepath.Join(ctx.SpectrRoot, "changes", "archive", ctx.ArchiveName),
	}

	// Add specs directory if specs were updated
	if !ctx.SkipSpecs {
		paths = append(paths, filepath.Join(ctx.SpectrRoot, "specs"))
	}

	if err := git.StageFiles(paths); err != nil {
		return err
	}

	fmt.Println("Staged files for commit")

	return nil
}
