package pr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/git"
)

// getCLIInstallSuggestion returns installation suggestions for a CLI tool.
func getCLIInstallSuggestion(tool string) string {
	switch tool {
	case "gh":
		return "Install: brew install gh or see https://cli.github.com"
	case "glab":
		return "Install: brew install glab or see " +
			"https://gitlab.com/gitlab-org/cli"
	case "tea":
		return "Install: brew install tea or see https://gitea.com/gitea/tea"
	default:
		return "please install the required CLI tool"
	}
}

// executeArchiveInWorktree runs the archive workflow within the worktree.
func executeArchiveInWorktree(
	config PRConfig,
	worktreePath string,
) (archive.ArchiveResult, error) {
	fmt.Println(
		"Running archive operation in worktree...",
	)

	// Copy the change from source to worktree first
	if err := copyChangeToWorktree(config, worktreePath); err != nil {
		return archive.ArchiveResult{},
			fmt.Errorf(
				"copy change to worktree: %w",
				err,
			)
	}

	// Create archive command
	archiveCmd := &archive.ArchiveCmd{
		ChangeID:  config.ChangeID,
		Yes:       true, // Non-interactive
		SkipSpecs: config.SkipSpecs,
	}

	// Execute archive within the worktree and capture results.
	// The ArchiveResult contains path, operation counts, and capabilities.
	result, err := archive.Archive(
		archiveCmd,
		worktreePath,
	)
	if err != nil {
		return archive.ArchiveResult{}, err
	}

	return result, nil
}

// copyChangeToWorktree copies the change directory from source to worktree.
func copyChangeToWorktree(
	config PRConfig,
	worktreePath string,
) error {
	projectRoot := config.ProjectRoot
	if projectRoot == "" {
		var err error
		projectRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf(
				"get repo root: %w",
				err,
			)
		}
	}

	sourceDir := filepath.Join(
		projectRoot,
		spectrDirName,
		changesDirName,
		config.ChangeID,
	)
	targetDir := filepath.Join(
		worktreePath,
		spectrDirName,
		changesDirName,
		config.ChangeID,
	)

	fmt.Printf(
		"Copying change to worktree: %s\n",
		config.ChangeID,
	)

	// Create target directory structure
	if err := os.MkdirAll(filepath.Dir(targetDir), dirPerm); err != nil {
		return fmt.Errorf(
			"create target directory: %w",
			err,
		)
	}

	// Copy directory recursively
	if err := copyDir(sourceDir, targetDir); err != nil {
		return fmt.Errorf(
			"copy directory: %w",
			err,
		)
	}

	return nil
}

// removeChangeInWorktree removes the change directory within the worktree.
func removeChangeInWorktree(
	config PRConfig,
	worktreePath string,
) error {
	changeDir := filepath.Join(
		worktreePath,
		spectrDirName,
		changesDirName,
		config.ChangeID,
	)

	fmt.Printf(
		"Removing change in worktree: %s\n",
		config.ChangeID,
	)

	// Verify the directory exists before attempting removal
	if _, err := os.Stat(changeDir); os.IsNotExist(
		err,
	) {
		return fmt.Errorf(
			"change directory does not exist in worktree: %s",
			changeDir,
		)
	}

	// Remove the entire change directory
	if err := os.RemoveAll(changeDir); err != nil {
		return fmt.Errorf(
			"remove change directory: %w",
			err,
		)
	}

	return nil
}

// cleanupLocalChange removes the local change directory from the working
// directory.
func cleanupLocalChange(config PRConfig) error {
	projectRoot := config.ProjectRoot
	if projectRoot == "" {
		var err error
		projectRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf(
				"get repo root: %w",
				err,
			)
		}
	}

	changeDir := filepath.Join(
		projectRoot,
		spectrDirName,
		changesDirName,
		config.ChangeID,
	)

	// Remove the entire change directory
	if err := os.RemoveAll(changeDir); err != nil {
		return fmt.Errorf(
			"remove local change directory: %w",
			err,
		)
	}

	return nil
}
