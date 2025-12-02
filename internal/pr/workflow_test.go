package pr

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/git"
)

// TestPRConfig_Struct tests basic PRConfig struct construction.
func TestPRConfig_Struct(t *testing.T) {
	tests := []struct {
		name   string
		config PRConfig
	}{
		{
			name: "archive mode config",
			config: PRConfig{
				ChangeID:    "add-feature-x",
				Mode:        ModeArchive,
				BaseBranch:  "main",
				Draft:       true,
				Force:       false,
				DryRun:      false,
				SkipSpecs:   true,
				ProjectRoot: "/tmp/project",
			},
		},
		{
			name: "new mode config",
			config: PRConfig{
				ChangeID:    "new-proposal",
				Mode:        ModeNew,
				BaseBranch:  "",
				Draft:       false,
				Force:       true,
				DryRun:      true,
				SkipSpecs:   false,
				ProjectRoot: "",
			},
		},
		{
			name: "minimal config",
			config: PRConfig{
				ChangeID: "minimal",
				Mode:     ModeNew,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify struct fields are set correctly
			if tt.config.ChangeID == "" {
				t.Error("ChangeID should not be empty")
			}

			// Test that struct is properly initialized
			cfg := PRConfig{
				ChangeID:    tt.config.ChangeID,
				Mode:        tt.config.Mode,
				BaseBranch:  tt.config.BaseBranch,
				Draft:       tt.config.Draft,
				Force:       tt.config.Force,
				DryRun:      tt.config.DryRun,
				SkipSpecs:   tt.config.SkipSpecs,
				ProjectRoot: tt.config.ProjectRoot,
			}

			if cfg.ChangeID != tt.config.ChangeID {
				t.Errorf("ChangeID = %q, want %q", cfg.ChangeID, tt.config.ChangeID)
			}
			if cfg.Mode != tt.config.Mode {
				t.Errorf("Mode = %q, want %q", cfg.Mode, tt.config.Mode)
			}
			if cfg.BaseBranch != tt.config.BaseBranch {
				t.Errorf("BaseBranch = %q, want %q", cfg.BaseBranch, tt.config.BaseBranch)
			}
			if cfg.Draft != tt.config.Draft {
				t.Errorf("Draft = %v, want %v", cfg.Draft, tt.config.Draft)
			}
			if cfg.Force != tt.config.Force {
				t.Errorf("Force = %v, want %v", cfg.Force, tt.config.Force)
			}
			if cfg.DryRun != tt.config.DryRun {
				t.Errorf("DryRun = %v, want %v", cfg.DryRun, tt.config.DryRun)
			}
			if cfg.SkipSpecs != tt.config.SkipSpecs {
				t.Errorf("SkipSpecs = %v, want %v", cfg.SkipSpecs, tt.config.SkipSpecs)
			}
			if cfg.ProjectRoot != tt.config.ProjectRoot {
				t.Errorf("ProjectRoot = %q, want %q", cfg.ProjectRoot, tt.config.ProjectRoot)
			}
		})
	}
}

// TestPRResult_Struct tests basic PRResult struct construction.
func TestPRResult_Struct(t *testing.T) {
	tests := []struct {
		name   string
		result PRResult
	}{
		{
			name: "full result - archive mode",
			result: PRResult{
				PRURL:       "https://github.com/owner/repo/pull/123",
				BranchName:  "spectr/archive/add-feature-x",
				ArchivePath: "spectr/changes/archive/2024-01-15-add-feature-x/",
				Counts: archive.OperationCounts{
					Added:    3,
					Modified: 2,
					Removed:  1,
					Renamed:  0,
				},
				Platform:  git.PlatformGitHub,
				ManualURL: "",
			},
		},
		{
			name: "bitbucket result with manual URL - proposal mode",
			result: PRResult{
				PRURL:       "",
				BranchName:  "spectr/proposal/new-feature",
				ArchivePath: "",
				Counts:      archive.OperationCounts{},
				Platform:    git.PlatformBitbucket,
				ManualURL:   "https://bitbucket.org/owner/repo/pull-requests/new",
			},
		},
		{
			name: "minimal result - proposal mode",
			result: PRResult{
				BranchName: "spectr/proposal/minimal",
				Platform:   git.PlatformGitLab,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that struct is properly initialized
			res := PRResult{
				PRURL:       tt.result.PRURL,
				BranchName:  tt.result.BranchName,
				ArchivePath: tt.result.ArchivePath,
				Counts:      tt.result.Counts,
				Platform:    tt.result.Platform,
				ManualURL:   tt.result.ManualURL,
			}

			if res.PRURL != tt.result.PRURL {
				t.Errorf("PRURL = %q, want %q", res.PRURL, tt.result.PRURL)
			}
			if res.BranchName != tt.result.BranchName {
				t.Errorf("BranchName = %q, want %q", res.BranchName, tt.result.BranchName)
			}
			if res.ArchivePath != tt.result.ArchivePath {
				t.Errorf("ArchivePath = %q, want %q", res.ArchivePath, tt.result.ArchivePath)
			}
			if res.Platform != tt.result.Platform {
				t.Errorf("Platform = %v, want %v", res.Platform, tt.result.Platform)
			}
			if res.ManualURL != tt.result.ManualURL {
				t.Errorf("ManualURL = %q, want %q", res.ManualURL, tt.result.ManualURL)
			}
			if res.Counts.Added != tt.result.Counts.Added {
				t.Errorf("Counts.Added = %d, want %d", res.Counts.Added, tt.result.Counts.Added)
			}
			if res.Counts.Modified != tt.result.Counts.Modified {
				t.Errorf(
					"Counts.Modified = %d, want %d",
					res.Counts.Modified,
					tt.result.Counts.Modified,
				)
			}
			if res.Counts.Removed != tt.result.Counts.Removed {
				t.Errorf(
					"Counts.Removed = %d, want %d",
					res.Counts.Removed,
					tt.result.Counts.Removed,
				)
			}
			if res.Counts.Renamed != tt.result.Counts.Renamed {
				t.Errorf(
					"Counts.Renamed = %d, want %d",
					res.Counts.Renamed,
					tt.result.Counts.Renamed,
				)
			}
		})
	}
}

// TestCheckCLITool tests CLI tool detection.
func TestCheckCLITool(t *testing.T) {
	tests := []struct {
		name      string
		tool      string
		wantError bool
	}{
		{
			name:      "git exists",
			tool:      "git",
			wantError: false,
		},
		{
			name:      "nonexistent tool",
			tool:      "nonexistent-tool-xyz-12345",
			wantError: true,
		},
		{
			name:      "another nonexistent tool",
			tool:      "fake-cli-tool-that-does-not-exist",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCLITool(tt.tool)
			if tt.wantError && err == nil {
				t.Errorf("checkCLITool(%q) expected error, got nil", tt.tool)
			}
			if !tt.wantError && err != nil {
				t.Errorf("checkCLITool(%q) unexpected error: %v", tt.tool, err)
			}
		})
	}
}

// TestGetCLIInstallSuggestion tests installation suggestion generation.
func TestGetCLIInstallSuggestion(t *testing.T) {
	tests := []struct {
		name             string
		tool             string
		expectedContains []string
	}{
		{
			name: "gh suggestions",
			tool: "gh",
			expectedContains: []string{
				"brew install gh",
				"https://cli.github.com",
			},
		},
		{
			name: "glab suggestions",
			tool: "glab",
			expectedContains: []string{
				"brew install glab",
				"https://gitlab.com/gitlab-org/cli",
			},
		},
		{
			name: "tea suggestions",
			tool: "tea",
			expectedContains: []string{
				"brew install tea",
				"https://gitea.com/gitea/tea",
			},
		},
		{
			name: "unknown tool generic message",
			tool: "unknown-tool",
			expectedContains: []string{
				"please install the required CLI tool",
			},
		},
		{
			name: "empty tool name",
			tool: "",
			expectedContains: []string{
				"please install the required CLI tool",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCLIInstallSuggestion(tt.tool)
			for _, exp := range tt.expectedContains {
				if !contains(result, exp) {
					t.Errorf("getCLIInstallSuggestion(%q) = %q, expected to contain %q",
						tt.tool, result, exp)
				}
			}
		})
	}
}

// TestCopyDir tests directory copying functionality.
func TestCopyDir(t *testing.T) {
	t.Run("copy directory with files", func(t *testing.T) {
		// Create source directory
		srcDir := t.TempDir()
		dstDir := filepath.Join(t.TempDir(), "destination")

		// Create test files in source
		testFiles := map[string]string{
			"file1.txt":         "content of file1",
			"file2.md":          "# Markdown content",
			"subdir/nested.txt": "nested file content",
			"subdir/deep/a.txt": "deeply nested",
			"another/file.json": `{"key": "value"}`,
		}

		for path, content := range testFiles {
			fullPath := filepath.Join(srcDir, path)
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}
			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write file %s: %v", path, err)
			}
		}

		// Copy directory
		if err := copyDir(srcDir, dstDir); err != nil {
			t.Fatalf("copyDir() error = %v", err)
		}

		// Verify all files were copied
		for path, expectedContent := range testFiles {
			fullPath := filepath.Join(dstDir, path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read copied file %s: %v", path, err)

				continue
			}
			if string(content) != expectedContent {
				t.Errorf("File %s content = %q, want %q", path, string(content), expectedContent)
			}
		}
	})

	t.Run("copy empty directory", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := filepath.Join(t.TempDir(), "empty-dest")

		if err := copyDir(srcDir, dstDir); err != nil {
			t.Fatalf("copyDir() error = %v", err)
		}

		// Verify destination exists
		info, err := os.Stat(dstDir)
		if err != nil {
			t.Fatalf("Destination directory should exist: %v", err)
		}
		if !info.IsDir() {
			t.Error("Destination should be a directory")
		}
	})

	t.Run("non-existent source directory", func(t *testing.T) {
		srcDir := "/nonexistent/source/directory/xyz123"
		dstDir := t.TempDir()

		err := copyDir(srcDir, dstDir)
		if err == nil {
			t.Error("copyDir() expected error for non-existent source, got nil")
		}
	})
}

// TestExtractURLFromOutput tests URL extraction from CLI output.
func TestExtractURLFromOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "typical gh output - single line",
			output:   "https://github.com/owner/repo/pull/123",
			expected: "https://github.com/owner/repo/pull/123",
		},
		{
			name: "gh output with leading text",
			output: `Creating pull request for spectr/add-feature against main
https://github.com/owner/repo/pull/456`,
			expected: "https://github.com/owner/repo/pull/456",
		},
		{
			name: "glab output with MR info",
			output: `Creating merge request from spectr/feature to main
https://gitlab.com/owner/repo/-/merge_requests/789
Merge request created!`,
			expected: "https://gitlab.com/owner/repo/-/merge_requests/789",
		},
		{
			name: "gitea output",
			output: `Pull request created
https://gitea.example.com/owner/repo/pulls/10
Done!`,
			expected: "https://gitea.example.com/owner/repo/pulls/10",
		},
		{
			name:     "http URL",
			output:   "http://localhost:3000/owner/repo/pull/5",
			expected: "http://localhost:3000/owner/repo/pull/5",
		},
		{
			name: "output with no URL",
			output: `Some output
that does not contain
any URLs at all`,
			expected: "Some output\nthat does not contain\nany URLs at all",
		},
		{
			name: "multiple URLs on separate lines - returns first URL line",
			output: `Some intro text
https://example.com/first
https://example.com/second`,
			expected: "https://example.com/first",
		},
		{
			name: "URLs not at line start - returns full output",
			output: `First: https://example.com/first
Second: https://example.com/second`,
			expected: "First: https://example.com/first\nSecond: https://example.com/second",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
		{
			name:     "whitespace only",
			output:   "   \n  \t  \n  ",
			expected: "",
		},
		{
			name:     "URL in middle of line",
			output:   `PR created at https://github.com/owner/repo/pull/999 successfully`,
			expected: "PR created at https://github.com/owner/repo/pull/999 successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractURLFromOutput(tt.output)
			if result != tt.expected {
				t.Errorf("extractURLFromOutput() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestWriteTempBodyFile tests temporary file creation for PR bodies.
func TestWriteTempBodyFile(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "simple body",
			body: "## Summary\n\nThis is a test PR body.",
		},
		{
			name: "complex markdown",
			body: `## Summary

This PR includes:

- Feature A
- Feature B

## Details

| Column | Value |
|--------|-------|
| One    | 1     |
| Two    | 2     |

**Bold text** and *italic text*

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```",
		},
		{
			name: "empty body",
			body: "",
		},
		{
			name: "unicode content",
			body: "Unicode: \u00e9\u00e8\u00ea \u4e2d\u6587 \U0001F600",
		},
		{
			name: "special characters",
			body: "Special: $PATH ${VAR} `backticks` \"quotes\" 'single'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, err := writeTempBodyFile(tt.body)
			if err != nil {
				t.Fatalf("writeTempBodyFile() error = %v", err)
			}

			// Ensure cleanup
			defer func() { _ = os.Remove(filePath) }()

			// Verify file exists
			info, err := os.Stat(filePath)
			if err != nil {
				t.Fatalf("Temp file should exist: %v", err)
			}
			if info.IsDir() {
				t.Error("Should be a file, not a directory")
			}

			// Verify content
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}
			if string(content) != tt.body {
				t.Errorf("File content = %q, want %q", string(content), tt.body)
			}
		})
	}
}

// TestWorkflowContext_Struct tests workflowContext struct fields.
func TestWorkflowContext_Struct(t *testing.T) {
	tests := []struct {
		name string
		ctx  workflowContext
	}{
		{
			name: "github context - archive mode",
			ctx: workflowContext{
				platformInfo: git.PlatformInfo{
					Platform: git.PlatformGitHub,
					CLITool:  "gh",
					RepoURL:  "https://github.com/owner/repo",
				},
				baseBranch: "origin/main",
				branchName: "spectr/archive/add-feature",
			},
		},
		{
			name: "gitlab context - proposal mode",
			ctx: workflowContext{
				platformInfo: git.PlatformInfo{
					Platform: git.PlatformGitLab,
					CLITool:  "glab",
					RepoURL:  "https://gitlab.com/group/project",
				},
				baseBranch: "origin/develop",
				branchName: "spectr/proposal/new-proposal",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := workflowContext{
				platformInfo: tt.ctx.platformInfo,
				baseBranch:   tt.ctx.baseBranch,
				branchName:   tt.ctx.branchName,
			}

			if ctx.platformInfo.Platform != tt.ctx.platformInfo.Platform {
				t.Errorf("Platform = %v, want %v",
					ctx.platformInfo.Platform, tt.ctx.platformInfo.Platform)
			}
			if ctx.platformInfo.CLITool != tt.ctx.platformInfo.CLITool {
				t.Errorf("CLITool = %q, want %q",
					ctx.platformInfo.CLITool, tt.ctx.platformInfo.CLITool)
			}
			if ctx.baseBranch != tt.ctx.baseBranch {
				t.Errorf("baseBranch = %q, want %q",
					ctx.baseBranch, tt.ctx.baseBranch)
			}
			if ctx.branchName != tt.ctx.branchName {
				t.Errorf("branchName = %q, want %q",
					ctx.branchName, tt.ctx.branchName)
			}
		})
	}
}

// TestCopyFile tests single file copying functionality.
func TestCopyFile(t *testing.T) {
	t.Run("copy regular file", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		srcPath := filepath.Join(srcDir, "source.txt")
		dstPath := filepath.Join(dstDir, "dest.txt")

		content := "This is the file content\nwith multiple lines"
		if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write source file: %v", err)
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			t.Fatalf("copyFile() error = %v", err)
		}

		// Verify content
		result, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("Failed to read dest file: %v", err)
		}
		if string(result) != content {
			t.Errorf("copyFile() content = %q, want %q", string(result), content)
		}
	})

	t.Run("copy file preserves permissions", func(t *testing.T) {
		srcDir := t.TempDir()
		dstDir := t.TempDir()

		srcPath := filepath.Join(srcDir, "source.sh")
		dstPath := filepath.Join(dstDir, "dest.sh")

		content := "#!/bin/bash\necho hello"
		if err := os.WriteFile(srcPath, []byte(content), 0755); err != nil {
			t.Fatalf("Failed to write source file: %v", err)
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			t.Fatalf("copyFile() error = %v", err)
		}

		// Verify permissions (mode bits might differ slightly by umask)
		srcInfo, _ := os.Stat(srcPath)
		dstInfo, _ := os.Stat(dstPath)

		// Check executable bit is preserved
		srcExec := srcInfo.Mode() & 0111
		dstExec := dstInfo.Mode() & 0111
		if srcExec != dstExec {
			t.Errorf("copyFile() exec bits = %o, want %o", dstExec, srcExec)
		}
	})

	t.Run("copy non-existent file", func(t *testing.T) {
		srcPath := "/nonexistent/file.txt"
		dstPath := filepath.Join(t.TempDir(), "dest.txt")

		err := copyFile(srcPath, dstPath)
		if err == nil {
			t.Error("copyFile() expected error for non-existent source, got nil")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
