package track

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/specterrs"
)

const (
	// gitCmd is the git command name.
	gitCmd = "git"

	// gitRepoFlag is the flag for specifying git repository root.
	gitRepoFlag = "-C"

	// commitFooter is appended to all automated commit messages.
	commitFooter = "[Automated by spectr track]"

	// minPorcelainLineLen is the minimum length of a git status porcelain line.
	minPorcelainLineLen = 3
)

// taskFiles lists files that should be excluded from staging.
var taskFiles = []string{
	"tasks.json",
	"tasks.jsonc",
	"tasks.md",
}

// Action represents the type of task status change.
type Action int

const (
	// ActionStart indicates a task transitioned to "in_progress".
	ActionStart Action = iota
	// ActionComplete indicates a task transitioned to "completed".
	ActionComplete
)

// String returns the action verb for commit messages.
func (a Action) String() string {
	switch a {
	case ActionStart:
		return "start"
	case ActionComplete:
		return "complete"
	default:
		return "update"
	}
}

// CommitResult contains the result of a commit operation.
type CommitResult struct {
	// NoFiles is true if no files were staged (only task files changed).
	NoFiles bool
	// CommitHash is the hash of the created commit (empty if NoFiles is true).
	CommitHash string
	// Message is the commit message used.
	Message string
}

// GitExecutor abstracts git operations for testing.
type GitExecutor interface {
	// Status runs `git status --porcelain` and returns the output.
	Status(repoRoot string) (string, error)
	// Add runs `git add` for the specified files.
	Add(repoRoot string, files []string) error
	// Commit runs `git commit` with the given message.
	Commit(repoRoot string, message string) error
	// RevParse runs `git rev-parse` and returns the result.
	RevParse(repoRoot string, ref string) (string, error)
	// DiffNumstat runs `git diff --numstat` for the specified files.
	// Returns output showing lines added/deleted per file.
	// Binary files show as "-\t-\t<filename>".
	DiffNumstat(repoRoot string, files []string) (string, error)
}

// RealGitExecutor implements GitExecutor using actual git commands.
type RealGitExecutor struct{}

// Status runs `git status --porcelain` and returns the output.
func (*RealGitExecutor) Status(repoRoot string) (string, error) {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, repoRoot,
		"status",
		"--porcelain",
	)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git status failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git status: %w", err)
	}

	return string(output), nil
}

// Add runs `git add` for the specified files.
func (*RealGitExecutor) Add(repoRoot string, files []string) error {
	args := []string{gitRepoFlag, repoRoot, "add", "--"}
	args = append(args, files...)

	cmd := exec.Command(gitCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git add failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// Commit runs `git commit` with the given message.
func (*RealGitExecutor) Commit(repoRoot, message string) error {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, repoRoot,
		"commit",
		"-m", message,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git commit failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// RevParse runs `git rev-parse` and returns the result.
func (*RealGitExecutor) RevParse(repoRoot, ref string) (string, error) {
	cmd := exec.Command(
		gitCmd,
		gitRepoFlag, repoRoot,
		"rev-parse",
		ref,
	)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git rev-parse failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// DiffNumstat runs `git diff --numstat` for the specified files.
// This is used to detect binary files, which show as "-\t-\t<filename>".
func (*RealGitExecutor) DiffNumstat(repoRoot string, files []string) (string, error) {
	// Use git diff --numstat to get line counts for each file.
	// Binary files will show "-" for both additions and deletions.
	args := []string{gitRepoFlag, repoRoot, "diff", "--numstat", "--"}
	args = append(args, files...)

	cmd := exec.Command(gitCmd, args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf(
				"git diff --numstat failed: %s",
				strings.TrimSpace(string(exitErr.Stderr)),
			)
		}

		return "", fmt.Errorf("failed to run git diff --numstat: %w", err)
	}

	return string(output), nil
}

// Committer handles git staging and commit operations for task tracking.
type Committer struct {
	changeID    string
	repoRoot    string
	gitExecutor GitExecutor
}

// NewCommitter creates a new Committer for the specified change.
func NewCommitter(changeID, repoRoot string) *Committer {
	return &Committer{
		changeID:    changeID,
		repoRoot:    repoRoot,
		gitExecutor: &RealGitExecutor{},
	}
}

// NewCommitterWithExecutor creates a new Committer with a custom GitExecutor.
// This is primarily used for testing with mock implementations.
func NewCommitterWithExecutor(
	changeID, repoRoot string,
	executor GitExecutor,
) *Committer {
	return &Committer{
		changeID:    changeID,
		repoRoot:    repoRoot,
		gitExecutor: executor,
	}
}

// Commit stages all modified files (excluding task files) and creates a commit.
// Returns CommitResult with NoFiles=true if only task files were modified.
// Returns a GitCommitError if git operations fail.
func (c *Committer) Commit(taskID string, action Action) (CommitResult, error) {
	modifiedFiles, err := c.getModifiedFiles()
	if err != nil {
		return CommitResult{}, &specterrs.GitCommitError{Err: err}
	}

	filesToStage := filterTaskFiles(modifiedFiles)
	if len(filesToStage) == 0 {
		return CommitResult{
			NoFiles: true,
			Message: fmt.Sprintf(
				"spectr(%s): %s task %s",
				c.changeID,
				action.String(),
				taskID,
			),
		}, nil
	}

	if err := c.stageFiles(filesToStage); err != nil {
		return CommitResult{}, &specterrs.GitCommitError{Err: err}
	}

	message := c.buildCommitMessage(taskID, action)
	hash, err := c.createCommit(message)
	if err != nil {
		return CommitResult{}, &specterrs.GitCommitError{Err: err}
	}

	return CommitResult{
		NoFiles:    false,
		CommitHash: hash,
		Message:    message,
	}, nil
}

// getModifiedFiles returns a list of modified files in the working tree.
// This includes both staged and unstaged modifications, as well as
// untracked files.
func (c *Committer) getModifiedFiles() ([]string, error) {
	output, err := c.gitExecutor.Status(c.repoRoot)
	if err != nil {
		return nil, err
	}

	return parseGitStatus(output), nil
}

// parseGitStatus parses git status porcelain output into file paths.
func parseGitStatus(output string) []string {
	var files []string
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if len(line) < minPorcelainLineLen {
			continue
		}

		status := line[0:2]
		filename := strings.TrimSpace(line[minPorcelainLineLen:])
		if filename == "" {
			continue
		}

		// Skip deleted files (D in either position)
		if status[0] == 'D' || status[1] == 'D' {
			continue
		}

		// Include modified (M), added (A), renamed (R), copied (C),
		// and untracked (?) files
		if status[0] != ' ' || status[1] != ' ' {
			files = append(files, filename)
		}
	}

	return files
}

// filterTaskFiles removes task files from the list of files to stage.
func filterTaskFiles(files []string) []string {
	var filtered []string

	for _, file := range files {
		baseName := filepath.Base(file)
		if !isTaskFile(baseName) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// isTaskFile checks if the given filename is a task file.
func isTaskFile(name string) bool {
	for _, taskFile := range taskFiles {
		if name == taskFile {
			return true
		}
	}

	return false
}

// detectBinaryFiles identifies which files from the given list are binary files.
// It uses git diff --numstat output where binary files show as "-\t-\t<filename>".
// Returns a map of binary file paths for O(1) lookup.
func parseBinaryFilesFromNumstat(numstatOutput string) map[string]bool {
	binaryFiles := make(map[string]bool)
	lines := strings.Split(numstatOutput, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Binary files in numstat output format: "-\t-\t<filename>"
		// Text files show: "<added>\t<deleted>\t<filename>"
		parts := strings.Split(line, "\t")
		if len(parts) >= 3 && parts[0] == "-" && parts[1] == "-" {
			// This is a binary file
			filename := parts[2]
			binaryFiles[filename] = true
		}
	}

	return binaryFiles
}

// getBinaryFiles returns a set of binary files from the given file list.
// Uses git diff --numstat to detect binary files.
func (c *Committer) getBinaryFiles(files []string) (map[string]bool, error) {
	if len(files) == 0 {
		return make(map[string]bool), nil
	}

	output, err := c.gitExecutor.DiffNumstat(c.repoRoot, files)
	if err != nil {
		return nil, err
	}

	return parseBinaryFilesFromNumstat(output), nil
}

// stageFiles stages the specified files for commit.
func (c *Committer) stageFiles(files []string) error {
	return c.gitExecutor.Add(c.repoRoot, files)
}

// buildCommitMessage creates the commit message with the standard format.
func (c *Committer) buildCommitMessage(taskID string, action Action) string {
	return fmt.Sprintf(
		"spectr(%s): %s task %s\n\n%s",
		c.changeID,
		action.String(),
		taskID,
		commitFooter,
	)
}

// createCommit creates a git commit with the given message and returns
// the commit hash.
func (c *Committer) createCommit(message string) (string, error) {
	if err := c.gitExecutor.Commit(c.repoRoot, message); err != nil {
		return "", err
	}

	return c.getCommitHash()
}

// getCommitHash returns the hash of the current HEAD commit.
func (c *Committer) getCommitHash() (string, error) {
	return c.gitExecutor.RevParse(c.repoRoot, "HEAD")
}
