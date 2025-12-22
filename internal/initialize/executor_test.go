package initialize

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutor_Execute(t *testing.T) {
	tmpDir, err := os.MkdirTemp(
		"",
		"spectr-exec-test-*",
	)
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = tmpDir
		err := cmd.Run()
		require.NoError(t, err)
	}

	runGit("init")
	runGit(
		"config",
		"user.email",
		"test@example.com",
	)
	runGit("config", "user.name", "test")
	err = os.WriteFile(
		filepath.Join(tmpDir, "README.md"),
		[]byte("# Test"),
		0644,
	)
	require.NoError(t, err)
	runGit("add", "README.md")
	runGit("commit", "-m", "initial")

	cmd := &InitCmd{
		Path: tmpDir,
	}
	executor, err := NewInitExecutor(cmd)
	require.NoError(t, err)

	result, err := executor.Execute(
		[]string{"claude-code"},
		false,
	)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify files created (physical check)
	assert.FileExists(
		t,
		filepath.Join(
			tmpDir,
			"spectr",
			"project.md",
		),
	)
	assert.FileExists(
		t,
		filepath.Join(
			tmpDir,
			"spectr",
			"AGENTS.md",
		),
	)
	assert.FileExists(
		t,
		filepath.Join(tmpDir, "CLAUDE.md"),
	)
	assert.FileExists(
		t,
		filepath.Join(
			tmpDir,
			".claude",
			"commands",
			"spectr",
			"proposal.md",
		),
	)

	// Verify git-detected changes in result
	// Note: Git might return CLAUDE.md as a change if it's new.
	assert.Contains(
		t,
		result.CreatedFiles,
		"CLAUDE.md",
	)
	assert.Contains(
		t,
		result.CreatedFiles,
		"spectr/project.md",
	)
}
