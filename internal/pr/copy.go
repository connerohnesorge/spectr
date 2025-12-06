// Package pr provides pull request creation and management functionality.
// This file contains functions for copying files and directories during
// the PR workflow, particularly for worktree operations.
package pr

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/git"
)

// copyChangeToWorktree copies the change directory from the source project
// to the worktree. This is necessary because worktrees are created from
// a clean branch and need the change files copied over.
func copyChangeToWorktree(config PRConfig, worktreePath string) error {
	projectRoot := config.ProjectRoot
	if projectRoot == "" {
		var err error
		projectRoot, err = git.GetRepoRoot()
		if err != nil {
			return fmt.Errorf("get repo root: %w", err)
		}
	}

	sourceDir := filepath.Join(
		projectRoot, "spectr", "changes", config.ChangeID,
	)
	targetDir := filepath.Join(
		worktreePath, "spectr", "changes", config.ChangeID,
	)

	fmt.Printf("Copying change to worktree: %s\n", config.ChangeID)

	if err := os.MkdirAll(filepath.Dir(targetDir), dirPerm); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	if err := copyDir(sourceDir, targetDir); err != nil {
		return fmt.Errorf("copy directory: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory and all its contents.
// It preserves file permissions from the source directory.
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

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

// copyFile copies a single file, preserving its permissions.
// It reads the entire source file and writes to the destination.
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
