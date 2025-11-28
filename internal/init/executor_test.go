package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
)

func TestNewInitExecutorWithRootDir(t *testing.T) {
	t.Run("creates executor with custom root dir", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := &InitCmd{Path: tmpDir}
		executor, err := NewInitExecutorWithRootDir(cmd, "my-specs")
		if err != nil {
			t.Fatalf("NewInitExecutorWithRootDir failed: %v", err)
		}

		if executor.rootDir != "my-specs" {
			t.Errorf("expected rootDir 'my-specs', got %s", executor.rootDir)
		}
	})

	t.Run("uses default root dir when empty", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := &InitCmd{Path: tmpDir}
		executor, err := NewInitExecutorWithRootDir(cmd, "")
		if err != nil {
			t.Fatalf("NewInitExecutorWithRootDir failed: %v", err)
		}

		if executor.rootDir != config.DefaultRootDir {
			t.Errorf(
				"expected rootDir '%s', got %s",
				config.DefaultRootDir,
				executor.rootDir,
			)
		}
	})

	t.Run("returns error for non-existent project path", func(t *testing.T) {
		cmd := &InitCmd{Path: "/non/existent/path"}
		_, err := NewInitExecutorWithRootDir(cmd, "spectr")
		if err == nil {
			t.Error("expected error for non-existent path")
		}

		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("returns error for empty project path", func(t *testing.T) {
		cmd := &InitCmd{Path: ""}
		_, err := NewInitExecutorWithRootDir(cmd, "spectr")
		if err == nil {
			t.Error("expected error for empty path")
		}

		if !strings.Contains(err.Error(), "path is required") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestCreateConfigFile(t *testing.T) {
	t.Run("creates config file for custom root dir", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := &InitCmd{Path: tmpDir}
		executor, err := NewInitExecutorWithRootDir(cmd, "my-specs")
		if err != nil {
			t.Fatalf("NewInitExecutorWithRootDir failed: %v", err)
		}

		result := &ExecutionResult{
			CreatedFiles: make([]string, 0),
		}

		err = executor.createConfigFile(result)
		if err != nil {
			t.Fatalf("createConfigFile failed: %v", err)
		}

		// Verify config file was created
		configPath := filepath.Join(tmpDir, config.ConfigFileName)
		if !FileExists(configPath) {
			t.Error("config file was not created")
		}

		// Verify content
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}

		expectedContent := "# Spectr configuration\nroot_dir: my-specs\n"
		if string(content) != expectedContent {
			t.Errorf("expected content:\n%s\ngot:\n%s", expectedContent, content)
		}

		// Verify file was tracked in result
		if len(result.CreatedFiles) != 1 {
			t.Fatalf("expected 1 created file, got %d", len(result.CreatedFiles))
		}
		if result.CreatedFiles[0] != config.ConfigFileName {
			t.Errorf(
				"expected created file '%s', got '%s'",
				config.ConfigFileName,
				result.CreatedFiles[0],
			)
		}
	})

	t.Run("does not create config file for default root dir", func(t *testing.T) {
		tmpDir := t.TempDir()

		cmd := &InitCmd{Path: tmpDir}
		executor, err := NewInitExecutorWithRootDir(cmd, config.DefaultRootDir)
		if err != nil {
			t.Fatalf("NewInitExecutorWithRootDir failed: %v", err)
		}

		result := &ExecutionResult{
			CreatedFiles: make([]string, 0),
		}

		err = executor.createConfigFile(result)
		if err != nil {
			t.Fatalf("createConfigFile failed: %v", err)
		}

		// Verify config file was NOT created
		configPath := filepath.Join(tmpDir, config.ConfigFileName)
		if FileExists(configPath) {
			t.Error("config file should not be created for default root dir")
		}

		// Verify nothing was tracked in result
		if len(result.CreatedFiles) != 0 {
			t.Errorf("expected 0 created files, got %d", len(result.CreatedFiles))
		}
	})

	t.Run("skips existing config file", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create existing config file
		configPath := filepath.Join(tmpDir, config.ConfigFileName)
		err := os.WriteFile(configPath, []byte("existing content"), 0644)
		if err != nil {
			t.Fatalf("failed to create existing config: %v", err)
		}

		cmd := &InitCmd{Path: tmpDir}
		executor, err := NewInitExecutorWithRootDir(cmd, "my-specs")
		if err != nil {
			t.Fatalf("NewInitExecutorWithRootDir failed: %v", err)
		}

		result := &ExecutionResult{
			CreatedFiles: make([]string, 0),
			Errors:       make([]string, 0),
		}

		err = executor.createConfigFile(result)
		if err != nil {
			t.Fatalf("createConfigFile failed: %v", err)
		}

		// Verify original content is preserved
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}

		if string(content) != "existing content" {
			t.Error("config file content was modified")
		}

		// Verify error was tracked
		if len(result.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(result.Errors))
		}
		if !strings.Contains(result.Errors[0], "already exists") {
			t.Errorf("unexpected error message: %s", result.Errors[0])
		}
	})
}

func TestExecuteWithCustomRootDir_CustomRoot(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := &InitCmd{Path: tmpDir}
	executor, err := NewInitExecutorWithRootDir(cmd, "my-specs")
	if err != nil {
		t.Fatalf("NewInitExecutorWithRootDir failed: %v", err)
	}

	result, err := executor.Execute(nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify custom root directory was created
	customRoot := filepath.Join(tmpDir, "my-specs")
	if !FileExists(customRoot) {
		t.Error("custom root directory was not created")
	}

	// Verify subdirectories
	verifySubdirectories(t, customRoot)

	// Verify project files
	verifyProjectFiles(t, customRoot)

	// Verify config file was created
	configPath := filepath.Join(tmpDir, config.ConfigFileName)
	if !FileExists(configPath) {
		t.Error("config file was not created")
	}

	// Verify result contains expected files
	verifyCreatedFiles(t, result, tmpDir, "my-specs")
}

func TestExecuteWithCustomRootDir_DefaultRoot(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := &InitCmd{Path: tmpDir}
	executor, err := NewInitExecutor(cmd)
	if err != nil {
		t.Fatalf("NewInitExecutor failed: %v", err)
	}

	result, err := executor.Execute(nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify default root directory was created
	defaultRoot := filepath.Join(tmpDir, "spectr")
	if !FileExists(defaultRoot) {
		t.Error("default root directory was not created")
	}

	// Verify config file was NOT created for default root
	configPath := filepath.Join(tmpDir, config.ConfigFileName)
	if FileExists(configPath) {
		t.Error("config file should not be created for default root")
	}

	// Verify result contains correct file paths with default root
	hasDefaultRootPath := false
	for _, createdFile := range result.CreatedFiles {
		if strings.HasPrefix(createdFile, "spectr/") {
			hasDefaultRootPath = true

			break
		}
	}
	if !hasDefaultRootPath {
		t.Error("CreatedFiles should contain paths with 'spectr/' prefix")
	}
}

// Helper functions to reduce cyclomatic complexity

func verifySubdirectories(t *testing.T, rootDir string) {
	t.Helper()

	specsDir := filepath.Join(rootDir, "specs")
	changesDir := filepath.Join(rootDir, "changes")

	if !FileExists(specsDir) {
		t.Error("specs directory was not created")
	}
	if !FileExists(changesDir) {
		t.Error("changes directory was not created")
	}
}

func verifyProjectFiles(t *testing.T, rootDir string) {
	t.Helper()

	projectFile := filepath.Join(rootDir, "project.md")
	if !FileExists(projectFile) {
		t.Error("project.md was not created in custom root")
	}

	agentsFile := filepath.Join(rootDir, "AGENTS.md")
	if !FileExists(agentsFile) {
		t.Error("AGENTS.md was not created in custom root")
	}
}

func verifyCreatedFiles(
	t *testing.T,
	result *ExecutionResult,
	tmpDir, rootName string,
) {
	t.Helper()

	expectedRelativeFiles := []string{
		config.ConfigFileName,
		rootName + "/project.md",
		rootName + "/AGENTS.md",
		"README.md",
	}

	for _, expectedFile := range expectedRelativeFiles {
		found := containsFile(result.CreatedFiles, expectedFile)
		if !found {
			t.Errorf(
				"expected file '%s' not found in CreatedFiles: %v",
				expectedFile,
				result.CreatedFiles,
			)
		}
	}

	// Verify directories were created (they use absolute paths)
	expectedDirs := []string{
		filepath.Join(tmpDir, rootName) + "/",
		filepath.Join(tmpDir, rootName, "specs") + "/",
		filepath.Join(tmpDir, rootName, "changes") + "/",
	}

	for _, expectedDir := range expectedDirs {
		found := containsFile(result.CreatedFiles, expectedDir)
		if !found {
			t.Errorf(
				"expected directory '%s' not found in CreatedFiles: %v",
				expectedDir,
				result.CreatedFiles,
			)
		}
	}
}

func containsFile(files []string, target string) bool {
	for _, f := range files {
		if f == target {
			return true
		}
	}

	return false
}
