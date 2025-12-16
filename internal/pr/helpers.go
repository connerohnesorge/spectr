package pr

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// dirPerm is the default permission for created directories.
const dirPerm = 0755

// spectrDirName is the name of the spectr directory.
const spectrDirName = "spectr"

// changesDirName is the name of the changes subdirectory.
const changesDirName = "changes"

// checkCLITool verifies that the required CLI tool is installed.
//
// This function checks if the specified CLI tool is installed in the system.
// If the tool is not found, it suggests installation instructions.
func checkCLITool(tool string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		suggestions := getCLIInstallSuggestion(
			tool,
		)

		return fmt.Errorf(
			"CLI tool '%s' not found in PATH; %s",
			tool,
			suggestions,
		)
	}

	return nil
}

// copyDir recursively copies a directory.
//
// Given a source directory, this function copies the directory and all its
// contents to the destination directory.
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(
			src,
			entry.Name(),
		)
		dstPath := filepath.Join(
			dst,
			entry.Name(),
		)

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(
		dst,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		srcInfo.Mode(),
	)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)

	return err
}

// stageAndCommit stages the spectr/ directory and creates a commit.
func stageAndCommit(
	worktreePath, commitMsg string,
) error {
	fmt.Println("Staging changes...")

	// git add spectr/
	addCmd := exec.Command(
		"git",
		"add",
		"spectr/",
	)
	addCmd.Dir = worktreePath

	output, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git add failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	fmt.Println("Creating commit...")

	// git commit
	commitCmd := exec.Command(
		"git",
		"commit",
		"-m",
		commitMsg,
	)
	commitCmd.Dir = worktreePath

	output, err = commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git commit failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// pushBranch pushes the branch to origin.
func pushBranch(
	worktreePath, branchName string,
) error {
	fmt.Printf("Pushing branch: %s\n", branchName)

	cmd := exec.Command(
		"git",
		"push",
		"-u",
		"origin",
		branchName,
	)
	cmd.Dir = worktreePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"git push failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}
