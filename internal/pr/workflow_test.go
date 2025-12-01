package pr

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/git"
)

func TestCheckCLIAvailable(t *testing.T) {
	tests := []struct {
		name    string
		tool    string
		wantErr bool
	}{
		{
			name:    "empty tool name",
			tool:    "",
			wantErr: true,
		},
		{
			name:    "nonexistent tool",
			tool:    "nonexistent-cli-tool-12345",
			wantErr: true,
		},
		// Note: We don't test for existing tools as they may not be installed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCLIAvailable(tt.tool)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"checkCLIAvailable() error = %v, wantErr %v",
					err, tt.wantErr,
				)
			}
		})
	}
}

func TestParsePRURL(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name: "GitHub PR URL",
			output: "Creating pull request...\n" +
				"https://github.com/owner/repo/pull/123\nDone",
			expected: "https://github.com/owner/repo/pull/123",
		},
		{
			name: "GitLab MR URL",
			output: "Creating merge request...\n" +
				"https://gitlab.com/owner/repo/-/merge_requests/456\nDone",
			expected: "https://gitlab.com/owner/repo/-/merge_requests/456",
		},
		{
			name: "Gitea PR URL",
			output: "Creating pull request...\n" +
				"https://gitea.example.com/owner/repo/pulls/789\nDone",
			expected: "https://gitea.example.com/owner/repo/pulls/789",
		},
		{
			name:     "no URL in output",
			output:   "Some error message without URL",
			expected: "",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
		{
			name:     "GitHub URL with org",
			output:   "https://github.com/my-org/my-repo/pull/1",
			expected: "https://github.com/my-org/my-repo/pull/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePRURL(tt.output)
			if result != tt.expected {
				t.Errorf("parsePRURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "spectr-copy-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create source structure
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), dirPerm); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	// Create test files
	files := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
	}
	for name, content := range files {
		path := filepath.Join(srcDir, name)
		if err := os.WriteFile(path, []byte(content), filePerm); err != nil {
			t.Fatalf("Failed to write file %s: %v", name, err)
		}
	}

	// Copy directory
	dstDir := filepath.Join(tmpDir, "dst")
	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir() error = %v", err)
	}

	// Verify copied files
	for name, expectedContent := range files {
		path := filepath.Join(dstDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read copied file %s: %v", name, err)

			continue
		}
		if string(content) != expectedContent {
			t.Errorf(
				"File %s content = %q, want %q",
				name, string(content), expectedContent,
			)
		}
	}
}

func TestCopyDir_NonExistentSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-copy-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	err = copyDir(
		filepath.Join(tmpDir, "nonexistent"),
		filepath.Join(tmpDir, "dst"),
	)
	if err == nil {
		t.Error("copyDir() expected error for nonexistent source")
	}
}

func TestConfig_Validation(t *testing.T) {
	// Test that Config fields have expected defaults (zero values)
	cfg := Config{}

	if cfg.ChangeID != "" {
		t.Error("Config.ChangeID should be empty by default")
	}
	if cfg.BaseBranch != "" {
		t.Error("Config.BaseBranch should be empty by default")
	}
	if cfg.Draft {
		t.Error("Config.Draft should be false by default")
	}
	if cfg.Force {
		t.Error("Config.Force should be false by default")
	}
	if cfg.DryRun {
		t.Error("Config.DryRun should be false by default")
	}
	if cfg.SkipSpecs {
		t.Error("Config.SkipSpecs should be false by default")
	}
}

func TestResult_Fields(t *testing.T) {
	result := Result{
		PRURL:       "https://github.com/owner/repo/pull/123",
		BranchName:  "spectr/archive/test-change",
		ArchivePath: "spectr/changes/archive/2024-01-15-test-change",
	}

	wantURL := "https://github.com/owner/repo/pull/123"
	if result.PRURL != wantURL {
		t.Errorf("Result.PRURL = %v, want %v", result.PRURL, wantURL)
	}

	wantBranch := "spectr/archive/test-change"
	if result.BranchName != wantBranch {
		t.Errorf("Result.BranchName = %v, want %v", result.BranchName, wantBranch)
	}

	wantPath := "spectr/changes/archive/2024-01-15-test-change"
	if result.ArchivePath != wantPath {
		t.Errorf("Result.ArchivePath = %v, want %v", result.ArchivePath, wantPath)
	}
}

func TestBuildArchiveCommand(t *testing.T) {
	tests := []struct {
		name      string
		changeID  string
		skipSpecs bool
		wantArgs  []string
	}{
		{
			name:      "basic archive command",
			changeID:  "test-change",
			skipSpecs: false,
			wantArgs:  []string{"archive", "test-change", "--yes"},
		},
		{
			name:      "archive with skip specs",
			changeID:  "test-change",
			skipSpecs: true,
			wantArgs: []string{
				"archive", "test-change", "--yes", "--skip-specs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := buildArchiveCommand(tt.changeID, tt.skipSpecs)

			// Check that args match (excluding the executable path)
			if len(cmd.Args) < len(tt.wantArgs)+1 {
				t.Errorf(
					"buildArchiveCommand() args = %v, want at least %d args",
					cmd.Args, len(tt.wantArgs)+1,
				)

				return
			}

			// Compare args (skip first which is executable)
			for i, arg := range tt.wantArgs {
				if cmd.Args[i+1] != arg {
					t.Errorf(
						"buildArchiveCommand() arg[%d] = %v, want %v",
						i+1, cmd.Args[i+1], arg,
					)
				}
			}
		})
	}
}

func TestCreatePR_Bitbucket(t *testing.T) {
	// Bitbucket should return ErrNoPRSupport
	platform := git.PlatformInfo{
		Platform: git.PlatformBitbucket,
		CLITool:  "bb",
	}

	_, err := createPR(platform, "test title", "test body", "main", false)
	if err != ErrNoPRSupport {
		t.Errorf(
			"createPR() for Bitbucket error = %v, want ErrNoPRSupport", err,
		)
	}
}

func TestCreatePR_Unknown(t *testing.T) {
	// Unknown platform should return ErrNoPRSupport
	platform := git.PlatformInfo{
		Platform: git.PlatformUnknown,
		CLITool:  "",
	}

	_, err := createPR(platform, "test title", "test body", "main", false)
	if err != ErrNoPRSupport {
		t.Errorf(
			"createPR() for Unknown platform error = %v, want ErrNoPRSupport",
			err,
		)
	}
}

func TestErrors(t *testing.T) {
	// Verify error messages are meaningful
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrNoCLI",
			err:  ErrNoCLI,
			want: "required CLI tool is not installed",
		},
		{
			name: "ErrNotAuthenticated",
			err:  ErrNotAuthenticated,
			want: "CLI tool is not authenticated",
		},
		{
			name: "ErrBranchExists",
			err:  ErrBranchExists,
			want: "branch already exists (use --force to overwrite)",
		},
		{
			name: "ErrChangeNotFound",
			err:  ErrChangeNotFound,
			want: "change not found",
		},
		{
			name: "ErrNoPRSupport",
			err:  ErrNoPRSupport,
			want: "platform does not support CLI PR creation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf(
					"%s.Error() = %v, want %v",
					tt.name, tt.err.Error(), tt.want,
				)
			}
		})
	}
}

func TestValidatePlatformCLI(t *testing.T) {
	// Test that Bitbucket and Unknown skip CLI validation
	tests := []struct {
		name     string
		platform git.PlatformInfo
		wantErr  bool
	}{
		{
			name: "Bitbucket skips validation",
			platform: git.PlatformInfo{
				Platform: git.PlatformBitbucket,
				CLITool:  "bb",
			},
			wantErr: false,
		},
		{
			name: "Unknown skips validation",
			platform: git.PlatformInfo{
				Platform: git.PlatformUnknown,
				CLITool:  "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlatformCLI(tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validatePlatformCLI() error = %v, wantErr %v",
					err, tt.wantErr,
				)
			}
		})
	}
}

func TestResolveBaseBranch(t *testing.T) {
	// Test that configured branch is returned as-is
	branch, err := resolveBaseBranch("develop")
	if err != nil {
		t.Errorf("resolveBaseBranch() unexpected error: %v", err)
	}
	if branch != "develop" {
		t.Errorf("resolveBaseBranch() = %v, want develop", branch)
	}
}

func TestResolveChangePath_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-change-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cfg := Config{
		ChangeID:   "nonexistent-change",
		WorkingDir: tmpDir,
	}

	_, err = resolveChangePath(cfg)
	if err != ErrChangeNotFound {
		t.Errorf("resolveChangePath() error = %v, want ErrChangeNotFound", err)
	}
}
