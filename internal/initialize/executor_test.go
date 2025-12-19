package initialize

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/git"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// TestInitializerIntegration tests the full initialization flow using afero.MemMapFs.
// Task 8.7: Integration test for full initialization flow
func TestInitializerIntegration(t *testing.T) {
	t.Run("initializers create expected files in memory filesystem", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ctx := context.Background()
		cfg := &providers.Config{SpectrDir: "spectr"}

		// Create a mock template manager
		tm := &MockTemplateManager{
			instructionPointer: "# Spectr Instructions\n\nSee spectr/AGENTS.md for details.",
			proposalCommand:    "# Proposal command content",
			applyCommand:       "# Apply command content",
		}

		// Test DirectoryInitializer
		t.Run("directory initializer creates directories", func(t *testing.T) {
			dirInit := providers.NewDirectoryInitializer(".claude/commands/spectr")
			err := dirInit.Init(ctx, fs, cfg, tm)
			if err != nil {
				t.Fatalf("DirectoryInitializer.Init() error = %v", err)
			}

			// Verify directory exists
			exists, err := afero.DirExists(fs, ".claude/commands/spectr")
			if err != nil {
				t.Fatalf("DirExists error = %v", err)
			}
			if !exists {
				t.Error("Directory .claude/commands/spectr was not created")
			}

			// Verify IsSetup returns true
			if !dirInit.IsSetup(fs, cfg) {
				t.Error("IsSetup() returned false after Init()")
			}
		})

		// Test ConfigFileInitializer
		t.Run("config file initializer creates file with markers", func(t *testing.T) {
			configInit := providers.NewConfigFileInitializer("CLAUDE.md")
			err := configInit.Init(ctx, fs, cfg, tm)
			if err != nil {
				t.Fatalf("ConfigFileInitializer.Init() error = %v", err)
			}

			// Verify file exists
			exists, err := afero.Exists(fs, "CLAUDE.md")
			if err != nil {
				t.Fatalf("Exists error = %v", err)
			}
			if !exists {
				t.Error("File CLAUDE.md was not created")
			}

			// Verify file contains markers
			content, err := afero.ReadFile(fs, "CLAUDE.md")
			if err != nil {
				t.Fatalf("ReadFile error = %v", err)
			}
			if !strings.Contains(string(content), "<!-- spectr:START -->") {
				t.Error("File missing start marker")
			}
			if !strings.Contains(string(content), "<!-- spectr:END -->") {
				t.Error("File missing end marker")
			}

			// Verify IsSetup returns true
			if !configInit.IsSetup(fs, cfg) {
				t.Error("IsSetup() returned false after Init()")
			}
		})

		// Test SlashCommandsInitializer
		t.Run("slash commands initializer creates command files", func(t *testing.T) {
			// Ensure directory exists first
			_ = fs.MkdirAll(".claude/commands/spectr", 0755)

			slashInit := providers.NewSlashCommandsInitializer(
				".claude/commands/spectr",
				".md",
				providers.FormatMarkdown,
			)
			err := slashInit.Init(ctx, fs, cfg, tm)
			if err != nil {
				t.Fatalf("SlashCommandsInitializer.Init() error = %v", err)
			}

			// Verify proposal.md exists
			exists, err := afero.Exists(fs, ".claude/commands/spectr/proposal.md")
			if err != nil {
				t.Fatalf("Exists error = %v", err)
			}
			if !exists {
				t.Error("File proposal.md was not created")
			}

			// Verify apply.md exists
			exists, err = afero.Exists(fs, ".claude/commands/spectr/apply.md")
			if err != nil {
				t.Fatalf("Exists error = %v", err)
			}
			if !exists {
				t.Error("File apply.md was not created")
			}

			// Verify IsSetup returns true
			if !slashInit.IsSetup(fs, cfg) {
				t.Error("IsSetup() returned false after Init()")
			}
		})
	})

	t.Run("TOML format slash commands", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ctx := context.Background()
		cfg := &providers.Config{SpectrDir: "spectr"}
		tm := &MockTemplateManager{
			proposalCommand: "proposal content",
			applyCommand:    "apply content",
		}

		// Ensure directory exists
		_ = fs.MkdirAll(".gemini/commands/spectr", 0755)

		slashInit := providers.NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			providers.FormatTOML,
		)
		err := slashInit.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("SlashCommandsInitializer.Init() error = %v", err)
		}

		// Verify proposal.toml exists
		exists, err := afero.Exists(fs, ".gemini/commands/spectr/proposal.toml")
		if err != nil {
			t.Fatalf("Exists error = %v", err)
		}
		if !exists {
			t.Error("File proposal.toml was not created")
		}

		// Verify TOML content structure
		content, err := afero.ReadFile(fs, ".gemini/commands/spectr/proposal.toml")
		if err != nil {
			t.Fatalf("ReadFile error = %v", err)
		}
		if !strings.Contains(string(content), "description =") {
			t.Error("TOML file missing description field")
		}
		if !strings.Contains(string(content), "prompt =") {
			t.Error("TOML file missing prompt field")
		}
	})

	t.Run("global initializers use global flag", func(t *testing.T) {
		globalDirInit := providers.NewGlobalDirectoryInitializer(".config/spectr")
		if !globalDirInit.IsGlobal() {
			t.Error("NewGlobalDirectoryInitializer should return global initializer")
		}

		projectDirInit := providers.NewDirectoryInitializer(".spectr")
		if projectDirInit.IsGlobal() {
			t.Error("NewDirectoryInitializer should return project initializer")
		}
	})

	t.Run("idempotent initialization", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ctx := context.Background()
		cfg := &providers.Config{SpectrDir: "spectr"}
		tm := &MockTemplateManager{
			instructionPointer: "# Instructions v1",
		}

		configInit := providers.NewConfigFileInitializer("TEST.md")

		// First initialization
		err := configInit.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("First Init() error = %v", err)
		}

		content1, _ := afero.ReadFile(fs, "TEST.md")

		// Second initialization should not fail
		err = configInit.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("Second Init() error = %v", err)
		}

		content2, _ := afero.ReadFile(fs, "TEST.md")

		// Content should be the same (idempotent)
		if string(content1) != string(content2) {
			t.Error("Initialization is not idempotent - content changed on second run")
		}
	})

	t.Run("config file update preserves existing content", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ctx := context.Background()
		cfg := &providers.Config{SpectrDir: "spectr"}
		tm := &MockTemplateManager{
			instructionPointer: "# Spectr Instructions",
		}

		// Create existing file with user content
		existingContent := "# My Custom Header\n\nSome user content here.\n"
		_ = afero.WriteFile(fs, "CLAUDE.md", []byte(existingContent), 0644)

		configInit := providers.NewConfigFileInitializer("CLAUDE.md")
		err := configInit.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("Init() error = %v", err)
		}

		content, _ := afero.ReadFile(fs, "CLAUDE.md")
		contentStr := string(content)

		// Should preserve existing content
		if !strings.Contains(contentStr, "# My Custom Header") {
			t.Error("Init() did not preserve existing content")
		}
		if !strings.Contains(contentStr, "Some user content here") {
			t.Error("Init() did not preserve existing content")
		}
		// Should add markers
		if !strings.Contains(contentStr, "<!-- spectr:START -->") {
			t.Error("Init() did not add start marker")
		}
		if !strings.Contains(contentStr, "<!-- spectr:END -->") {
			t.Error("Init() did not add end marker")
		}
	})
}

// TestProviderInitializerIntegration tests that each provider's initializers
// work correctly together.
func TestProviderInitializerIntegration(t *testing.T) {
	ctx := context.Background()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := &MockTemplateManager{
		instructionPointer: "# Spectr\n\nSee spectr/AGENTS.md",
		proposalCommand:    "# Proposal\n\nCreate proposals here.",
		applyCommand:       "# Apply\n\nApply changes here.",
	}

	// Test a subset of providers (the most representative ones)
	testCases := []struct {
		providerID string
		wantFiles  []string
	}{
		{
			providerID: "claude-code",
			wantFiles: []string{
				".claude/commands/spectr",
				"CLAUDE.md",
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
			},
		},
		{
			providerID: "gemini",
			wantFiles: []string{
				".gemini/commands/spectr",
				".gemini/commands/spectr/proposal.toml",
				".gemini/commands/spectr/apply.toml",
			},
		},
		{
			providerID: "cursor",
			wantFiles: []string{
				".cursorrules/commands/spectr",
				".cursorrules/commands/spectr/proposal.md",
				".cursorrules/commands/spectr/apply.md",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.providerID, func(t *testing.T) {
			reg := providers.Get(tc.providerID)
			if reg == nil {
				t.Fatalf("Provider %q not found", tc.providerID)
			}

			fs := afero.NewMemMapFs()
			inits := reg.Provider.Initializers(ctx)

			// Run all initializers
			for _, init := range inits {
				err := init.Init(ctx, fs, cfg, tm)
				if err != nil {
					t.Fatalf("Initializer %q failed: %v", init.Path(), err)
				}
			}

			// Verify expected files/directories exist
			for _, path := range tc.wantFiles {
				exists, err := afero.Exists(fs, path)
				if err != nil {
					t.Fatalf("Exists(%q) error = %v", path, err)
				}
				if !exists {
					// Check if it's a directory
					dirExists, _ := afero.DirExists(fs, path)
					if !dirExists {
						t.Errorf("Expected file/directory %q does not exist", path)
					}
				}
			}
		})
	}
}

// TestGitChangeDetectionIntegration tests git diff change detection.
// Task 8.8: Integration test verifying git diff change detection
func TestGitChangeDetectionIntegration(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available")
	}

	t.Run("detects new files created during initialization", func(t *testing.T) {
		// Create a temporary git repository
		tempDir, err := os.MkdirTemp("", "spectr-git-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Initialize git repo
		cmd := exec.Command("git", "init")
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to init git repo: %v", err)
		}

		// Configure git user
		cmd = exec.Command("git", "config", "user.email", "test@example.com")
		cmd.Dir = tempDir
		_ = cmd.Run()

		cmd = exec.Command("git", "config", "user.name", "Test User")
		cmd.Dir = tempDir
		_ = cmd.Run()

		// Create initial file and commit
		initialFile := filepath.Join(tempDir, "README.md")
		_ = os.WriteFile(initialFile, []byte("# Test"), 0644)

		cmd = exec.Command("git", "add", ".")
		cmd.Dir = tempDir
		_ = cmd.Run()

		cmd = exec.Command("git", "commit", "-m", "initial")
		cmd.Dir = tempDir
		_ = cmd.Run()

		// Create change detector and take snapshot
		detector := git.NewChangeDetector(tempDir)
		snapshot, err := detector.Snapshot()
		if err != nil {
			t.Fatalf("Snapshot() error = %v", err)
		}

		// Create new files (simulating initialization)
		spectrDir := filepath.Join(tempDir, "spectr")
		_ = os.MkdirAll(spectrDir, 0755)
		_ = os.WriteFile(filepath.Join(spectrDir, "project.md"), []byte("# Project"), 0644)
		_ = os.WriteFile(filepath.Join(spectrDir, "AGENTS.md"), []byte("# Agents"), 0644)

		claudeDir := filepath.Join(tempDir, ".claude", "commands", "spectr")
		_ = os.MkdirAll(claudeDir, 0755)
		_ = os.WriteFile(filepath.Join(claudeDir, "proposal.md"), []byte("# Proposal"), 0644)
		_ = os.WriteFile(filepath.Join(claudeDir, "apply.md"), []byte("# Apply"), 0644)

		// Detect changes
		changedFiles, err := detector.ChangedFiles(snapshot)
		if err != nil {
			t.Fatalf("ChangedFiles() error = %v", err)
		}

		// Git status --porcelain may return directories (e.g., "spectr/", ".claude/")
		// instead of individual files when an entire directory is new.
		// We verify that either the specific file path or its parent directory is detected.
		expectedPatterns := []struct {
			file       string
			parentDirs []string // Alternative directory patterns that would also match
		}{
			{
				file:       "spectr/project.md",
				parentDirs: []string{"spectr/", "spectr"},
			},
			{
				file:       "spectr/AGENTS.md",
				parentDirs: []string{"spectr/", "spectr"},
			},
			{
				file:       ".claude/commands/spectr/proposal.md",
				parentDirs: []string{".claude/", ".claude"},
			},
			{
				file:       ".claude/commands/spectr/apply.md",
				parentDirs: []string{".claude/", ".claude"},
			},
		}

		for _, expected := range expectedPatterns {
			found := false
			for _, actual := range changedFiles {
				// Check exact file match
				if actual == expected.file {
					found = true
					break
				}
				// Check if parent directory is reported (untracked directories)
				for _, dir := range expected.parentDirs {
					if actual == dir {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				t.Errorf("Expected file %q (or its parent directory) not in changed files: %v",
					expected.file, changedFiles)
			}
		}
	})

	t.Run("detects modified files", func(t *testing.T) {
		// Create a temporary git repository
		tempDir, err := os.MkdirTemp("", "spectr-git-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Initialize and configure git repo
		cmd := exec.Command("git", "init")
		cmd.Dir = tempDir
		_ = cmd.Run()

		cmd = exec.Command("git", "config", "user.email", "test@example.com")
		cmd.Dir = tempDir
		_ = cmd.Run()

		cmd = exec.Command("git", "config", "user.name", "Test User")
		cmd.Dir = tempDir
		_ = cmd.Run()

		// Create and commit an existing file
		existingFile := filepath.Join(tempDir, "CLAUDE.md")
		_ = os.WriteFile(existingFile, []byte("# Existing content"), 0644)

		cmd = exec.Command("git", "add", ".")
		cmd.Dir = tempDir
		_ = cmd.Run()

		cmd = exec.Command("git", "commit", "-m", "initial")
		cmd.Dir = tempDir
		_ = cmd.Run()

		// Take snapshot
		detector := git.NewChangeDetector(tempDir)
		snapshot, err := detector.Snapshot()
		if err != nil {
			t.Fatalf("Snapshot() error = %v", err)
		}

		// Modify the file (simulating config file update)
		_ = os.WriteFile(
			existingFile,
			[]byte("# Existing content\n\n<!-- spectr:START -->\n# Spectr\n<!-- spectr:END -->\n"),
			0644,
		)

		// Detect changes
		changedFiles, err := detector.ChangedFiles(snapshot)
		if err != nil {
			t.Fatalf("ChangedFiles() error = %v", err)
		}

		// Verify modified file is detected
		found := false
		for _, f := range changedFiles {
			if f == "CLAUDE.md" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Modified file CLAUDE.md not detected in: %v", changedFiles)
		}
	})

	t.Run("IsGitRepo returns false for non-git directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "not-a-repo-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		if git.IsGitRepo(tempDir) {
			t.Error("IsGitRepo() returned true for non-git directory")
		}
	})

	t.Run("ErrNotGitRepo has correct message", func(t *testing.T) {
		expectedMsg := "spectr init requires a git repository. Run 'git init' first."
		if git.ErrNotGitRepo.Error() != expectedMsg {
			t.Errorf("ErrNotGitRepo = %q, want %q", git.ErrNotGitRepo.Error(), expectedMsg)
		}
	})
}

// MockTemplateManager implements providers.TemplateManager for testing.
type MockTemplateManager struct {
	instructionPointer string
	proposalCommand    string
	applyCommand       string
}

func (m *MockTemplateManager) RenderInstructionPointer(
	ctx providers.TemplateContext,
) (string, error) {
	return m.instructionPointer, nil
}

func (m *MockTemplateManager) RenderSlashCommand(
	commandType string,
	ctx providers.TemplateContext,
) (string, error) {
	switch commandType {
	case "proposal":
		return m.proposalCommand, nil
	case "apply":
		return m.applyCommand, nil
	default:
		return "", nil
	}
}
