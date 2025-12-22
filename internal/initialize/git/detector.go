// Package git provides utilities for git operations.
package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ChangeDetector detects changes in a git repository.
type ChangeDetector struct {
	projectPath string
}

// NewChangeDetector creates a new ChangeDetector.
func NewChangeDetector(
	projectPath string,
) *ChangeDetector {
	return &ChangeDetector{
		projectPath: projectPath,
	}
}

// IsGitRepo checks if the path is a git repository.
func IsGitRepo(path string) bool {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--is-inside-work-tree",
	)
	cmd.Dir = path

	return cmd.Run() == nil
}

// Snapshot captures the current state of the working tree.
// It returns a string representing the state (stash commit + untracked files).
func (d *ChangeDetector) Snapshot() (string, error) {
	if !IsGitRepo(d.projectPath) {
		return "", errors.New("not a git repository")
	}

	stashHash, err := d.runGit("stash", "create")
	if err != nil {
		return "", err
	}

	untracked, err := d.runGit(
		"ls-files",
		"--others",
		"--exclude-standard",
	)
	if err != nil {
		return "", err
	}

	return stashHash + "|" + untracked, nil
}

// ChangedFiles calculates the files changed between the beforeSnapshot and now.
func (d *ChangeDetector) ChangedFiles(
	beforeSnapshot string,
) ([]string, error) {
	currentSnapshot, err := d.Snapshot()
	if err != nil {
		return nil, err
	}

	beforeParts := strings.SplitN(beforeSnapshot, "|", 2)
	currentParts := strings.SplitN(currentSnapshot, "|", 2)

	if len(beforeParts) != 2 || len(currentParts) != 2 {
		return nil, errors.New("invalid snapshot format")
	}

	beforeStash, beforeUntracked := beforeParts[0], beforeParts[1]
	currentStash, currentUntracked := currentParts[0], currentParts[1]

	changes := make(map[string]bool)

	if err := d.detectStashChanges(beforeStash, currentStash, changes); err != nil {
		return nil, err
	}

	d.detectUntrackedChanges(beforeUntracked, currentUntracked, changes)

	result := make([]string, 0, len(changes))
	for f := range changes {
		result = append(result, f)
	}

	return result, nil
}

func (d *ChangeDetector) detectStashChanges(beforeStash, currentStash string, changes map[string]bool) error {
	if beforeStash == currentStash {
		return nil
	}

	var diffOut string
	var err error

	switch {
	case beforeStash == "":
		diffOut, err = d.runGit("diff", "--name-only", "HEAD", currentStash)
	case currentStash == "":
		diffOut, err = d.runGit("diff", "--name-only", beforeStash, "HEAD")
	default:
		diffOut, err = d.runGit("diff", "--name-only", beforeStash, currentStash)
	}

	if err != nil {
		return err
	}

	for _, f := range strings.Split(diffOut, "\n") {
		if f != "" {
			changes[f] = true
		}
	}

	return nil
}

func (d *ChangeDetector) detectUntrackedChanges(beforeUntracked, currentUntracked string, changes map[string]bool) {
	beforeUntrackedMap := make(map[string]bool)
	for _, f := range strings.Split(beforeUntracked, "\n") {
		if f != "" {
			beforeUntrackedMap[f] = true
		}
	}

	for _, f := range strings.Split(currentUntracked, "\n") {
		if f != "" && !beforeUntrackedMap[f] {
			changes[f] = true
		}
	}
}

func (d *ChangeDetector) runGit(
	args ...string,
) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = d.projectPath
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Some commands exit non-zero on legitimate states?
			// stash create exits 0. ls-files exits 0.
			return "", fmt.Errorf(
				"git %v failed: %s",
				args,
				string(exitErr.Stderr),
			)
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}