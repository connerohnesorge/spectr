package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) string {
	tmpDir, err := os.MkdirTemp(
		"",
		"spectr-git-test-*",
	)
	require.NoError(t, err)

	run := func(name string, args ...string) {
		cmd := exec.Command(name, args...)
		cmd.Dir = tmpDir
		err := cmd.Run()
		require.NoError(
			t,
			err,
			"failed to run %s %v",
			name,
			args,
		)
	}

	run("git", "init")
	run(
		"git",
		"config",
		"user.email",
		"test@example.com",
	)
	run("git", "config", "user.name", "Test User")

	// Create an initial commit so HEAD exists
	initialFile := filepath.Join(
		tmpDir,
		"initial.txt",
	)
	err = os.WriteFile(
		initialFile,
		[]byte("initial"),
		0644,
	)
	require.NoError(t, err)

	run("git", "add", "initial.txt")
	run("git", "commit", "-m", "initial commit")

	return tmpDir
}

func TestChangeDetector(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	detector := NewChangeDetector(tmpDir)

	assert.True(t, IsGitRepo(tmpDir))
	assert.False(t, IsGitRepo(os.TempDir()))

	// 1. Take initial snapshot
	snapshot1, err := detector.Snapshot()
	assert.NoError(t, err)

	// 2. Modify a tracked file and add an untracked file
	err = os.WriteFile(
		filepath.Join(tmpDir, "initial.txt"),
		[]byte("modified"),
		0644,
	)
	assert.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(tmpDir, "new.txt"),
		[]byte("new file"),
		0644,
	)
	assert.NoError(t, err)

	// 3. Get changed files
	changed, err := detector.ChangedFiles(
		snapshot1,
	)
	assert.NoError(t, err)

	assert.Contains(t, changed, "initial.txt")
	assert.Contains(t, changed, "new.txt")
	assert.Len(t, changed, 2)
}
