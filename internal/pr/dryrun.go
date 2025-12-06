package pr

// This file provides dry-run simulation for PR workflow operations.

import (
	"fmt"
	"strings"
)

// executeDryRun simulates the PR workflow without making changes.
func executeDryRun(
	config PRConfig,
	ctx *workflowContext,
) (*PRResult, error) {
	fmt.Println()
	fmt.Println("=== DRY RUN MODE ===")
	fmt.Println("The following actions would be performed:")
	fmt.Println()

	printWorktreeStep(ctx)
	printOperationStep(config)
	printStageStep()
	printCommitStep(config)
	printPushStep(ctx)
	printPRStep(config, ctx)
	printCleanupStep()

	fmt.Println()
	fmt.Println("=== END DRY RUN ===")

	return &PRResult{
		BranchName: ctx.branchName,
		Platform:   ctx.platformInfo.Platform,
	}, nil
}

// printWorktreeStep prints the worktree creation step.
func printWorktreeStep(ctx *workflowContext) {
	fmt.Printf("1. Create worktree on branch: %s (based on %s)\n",
		ctx.branchName, ctx.baseBranch)
}

// printOperationStep prints the operation step (archive or copy).
func printOperationStep(config PRConfig) {
	switch config.Mode {
	case ModeArchive:
		fmt.Printf("2. Copy change '%s' to worktree\n", config.ChangeID)
		fmt.Printf("3. Run archive (skip-specs: %v)\n", config.SkipSpecs)
	case ModeProposal:
		fmt.Printf("2. Copy change '%s' to worktree\n", config.ChangeID)
	case ModeRemove:
		fmt.Printf("2. Copy change '%s' to worktree\n", config.ChangeID)
		fmt.Println("3. Remove change directory")
	}
}

// printStageStep prints the staging step.
func printStageStep() {
	fmt.Println("4. Stage changes: git add spectr/")
}

// printCommitStep prints the commit step with message preview.
func printCommitStep(config PRConfig) {
	fmt.Println("5. Create commit with message:")

	commitData := CommitTemplateData{
		ChangeID: config.ChangeID,
		Mode:     config.Mode,
	}

	commitMsg, _ := RenderCommitMessage(commitData)

	for _, line := range strings.Split(commitMsg, "\n") {
		fmt.Printf("   | %s\n", line)
	}
}

// printPushStep prints the push step.
func printPushStep(ctx *workflowContext) {
	fmt.Printf("\n6. Push branch: git push -u origin %s\n", ctx.branchName)
}

// printPRStep prints the PR creation step.
func printPRStep(config PRConfig, ctx *workflowContext) {
	prTitle := GetPRTitle(config.ChangeID, config.Mode)

	fmt.Println("7. Create PR:")
	fmt.Printf("   Platform: %s\n", ctx.platformInfo.Platform)
	fmt.Printf("   CLI tool: %s\n", ctx.platformInfo.CLITool)
	fmt.Printf("   Title: %s\n", prTitle)
	fmt.Printf("   Base: %s\n", strings.TrimPrefix(ctx.baseBranch, "origin/"))
	fmt.Printf("   Draft: %v\n", config.Draft)
}

// printCleanupStep prints the cleanup step.
func printCleanupStep() {
	fmt.Println("\n8. Cleanup worktree")
}
