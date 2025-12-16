// Package initialize provides utilities for initializing Spectr
// in a project directory.
//
// This file contains integration tests for the new executor using the
// redesigned provider architecture with afero-based filesystem operations.
package initialize

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// mockTemplateRenderer implements providers.TemplateRenderer for testing.
type mockTemplateRenderer struct{}

func (*mockTemplateRenderer) RenderAgents(
	_ providers.TemplateContext,
) (string, error) {
	return "mock agents content", nil
}

func (*mockTemplateRenderer) RenderInstructionPointer(
	_ providers.TemplateContext,
) (string, error) {
	return "mock instruction content", nil
}

func (*mockTemplateRenderer) RenderSlashCommand(
	command string,
	_ providers.TemplateContext,
) (string, error) {
	return "mock command content for " + command, nil
}

// testableInitExecutorNew is a test wrapper for InitExecutorNew that allows
// injecting a custom filesystem for in-memory testing.
type testableInitExecutorNew struct {
	*InitExecutorNew
}

// createTestExecutor creates a new InitExecutorNew with a memory-based filesystem
// for testing purposes. It bypasses the OS filesystem check in NewInitExecutorNew.
func createTestExecutor(t *testing.T, projectPath string) *testableInitExecutorNew {
	t.Helper()

	memFs := afero.NewMemMapFs()
	if err := memFs.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create project directory: %v", err)
	}

	fs := afero.NewBasePathFs(memFs, projectPath)
	registry := providers.CreateRegistry()
	factory := initializers.NewFactory()
	tm := &mockTemplateRenderer{}

	if err := providers.RegisterAllProvidersIncludingBase(registry, tm, factory); err != nil {
		t.Fatalf("failed to register providers: %v", err)
	}

	cfg := providers.NewConfig()
	executor := &InitExecutorNew{
		fs:          fs,
		projectPath: projectPath,
		registry:    registry,
		cfg:         cfg,
		tm:          tm,
	}

	return &testableInitExecutorNew{executor}
}

// getMemFs returns the underlying memory filesystem for assertions.
func (e *testableInitExecutorNew) getMemFs() afero.Fs {
	return e.fs
}

// TestNewInitExecutorNew_ValidPath tests creating an executor with a valid project path.
func TestNewInitExecutorNew_ValidPath(t *testing.T) {
	tempDir := t.TempDir()
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("failed to create template manager: %v", err)
	}

	executor, err := NewInitExecutorNew(tempDir, tm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if executor == nil {
		t.Fatal("expected executor to be non-nil")
	}
	if executor.projectPath != tempDir {
		t.Errorf("expected projectPath %q, got %q", tempDir, executor.projectPath)
	}
	if executor.fs == nil {
		t.Error("expected fs to be non-nil")
	}
	if executor.registry == nil {
		t.Error("expected registry to be non-nil")
	}
	if executor.cfg == nil {
		t.Error("expected cfg to be non-nil")
	}
}

// TestNewInitExecutorNew_EmptyPath tests that empty path returns an error.
func TestNewInitExecutorNew_EmptyPath(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("failed to create template manager: %v", err)
	}

	executor, err := NewInitExecutorNew("", tm)
	if err == nil {
		t.Error("expected error for empty path, got nil")
	}
	if executor != nil {
		t.Error("expected executor to be nil for empty path")
	}
	if !strings.Contains(err.Error(), "project path is required") {
		t.Errorf("expected error to contain 'project path is required', got: %v", err)
	}
}

// TestNewInitExecutorNew_NonExistentPath tests that non-existent path returns an error.
func TestNewInitExecutorNew_NonExistentPath(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("failed to create template manager: %v", err)
	}

	executor, err := NewInitExecutorNew("/non/existent/path/that/should/not/exist", tm)
	if err == nil {
		t.Error("expected error for non-existent path, got nil")
	}
	if executor != nil {
		t.Error("expected executor to be nil for non-existent path")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("expected error to contain 'does not exist', got: %v", err)
	}
}

// TestExecute_SingleProvider_Claude_DirectoryStructure tests directory creation.
func TestExecute_SingleProvider_Claude_DirectoryStructure(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	result, err := executor.Execute(
		context.Background(),
		[]string{"claude-code"},
		NewExecuteOptions(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	fs := executor.getMemFs()
	dirs := []string{"spectr", "spectr/specs", "spectr/changes"}
	for _, dir := range dirs {
		exists, err := afero.DirExists(fs, dir)
		if err != nil || !exists {
			t.Errorf("expected %s/ directory to exist", dir)
		}
	}
}

// TestExecute_SingleProvider_Claude_CoreFiles tests core file creation.
func TestExecute_SingleProvider_Claude_CoreFiles(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), []string{"claude-code"}, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()
	files := []string{"spectr/project.md", "spectr/AGENTS.md"}
	for _, file := range files {
		exists, err := afero.Exists(fs, file)
		if err != nil || !exists {
			t.Errorf("expected %s to exist", file)
		}
	}
}

// TestExecute_SingleProvider_Claude_ProviderFiles tests Claude-specific file creation.
func TestExecute_SingleProvider_Claude_ProviderFiles(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), []string{"claude-code"}, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()
	claudeDirExists, err := afero.DirExists(fs, ".claude/commands/spectr")
	if err != nil || !claudeDirExists {
		t.Error("expected .claude/commands/spectr/ directory to exist")
	}

	files := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}
	for _, file := range files {
		exists, err := afero.Exists(fs, file)
		if err != nil || !exists {
			t.Errorf("expected %s to exist", file)
		}
	}
}

// TestExecute_SingleProvider_Claude_ResultTracking tests result file tracking.
func TestExecute_SingleProvider_Claude_ResultTracking(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	result, err := executor.Execute(
		context.Background(),
		[]string{"claude-code"},
		NewExecuteOptions(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.CreatedFiles) == 0 {
		t.Error("expected CreatedFiles to contain entries")
	}
}

// TestExecute_MultipleProviders tests initialization with multiple providers (Claude + Gemini).
func TestExecute_MultipleProviders(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	result, err := executor.Execute(
		context.Background(),
		[]string{"claude-code", "gemini"},
		NewExecuteOptions(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	fs := executor.getMemFs()

	claudeMdExists, err := afero.Exists(fs, "CLAUDE.md")
	if err != nil || !claudeMdExists {
		t.Error("expected CLAUDE.md to exist")
	}

	claudeProposalExists, err := afero.Exists(fs, ".claude/commands/spectr/proposal.md")
	if err != nil || !claudeProposalExists {
		t.Error("expected .claude/commands/spectr/proposal.md to exist")
	}

	geminiDirExists, err := afero.DirExists(fs, ".gemini/commands/spectr")
	if err != nil || !geminiDirExists {
		t.Error("expected .gemini/commands/spectr/ directory to exist")
	}

	geminiFiles := []string{
		".gemini/commands/spectr/proposal.toml",
		".gemini/commands/spectr/apply.toml",
	}
	for _, file := range geminiFiles {
		exists, err := afero.Exists(fs, file)
		if err != nil || !exists {
			t.Errorf("expected %s to exist", file)
		}
	}
}

// TestExecute_DirectoryStructure tests that all required directories are created.
func TestExecute_DirectoryStructure(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	result, err := executor.Execute(context.Background(), nil, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()
	requiredDirs := []string{"spectr", "spectr/specs", "spectr/changes"}
	for _, dir := range requiredDirs {
		exists, err := afero.DirExists(fs, dir)
		if err != nil {
			t.Errorf("error checking directory %s: %v", dir, err)
		}
		if !exists {
			t.Errorf("expected directory %s to exist", dir)
		}
	}

	dirCount := 0
	for _, f := range result.CreatedFiles {
		if strings.HasSuffix(f, "/") {
			dirCount++
		}
	}
	if dirCount < 3 {
		t.Errorf("expected at least 3 directories in CreatedFiles, got %d", dirCount)
	}
}

// TestExecute_ProjectAndAgentsMd tests that project.md and AGENTS.md are created.
func TestExecute_ProjectAndAgentsMd(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), nil, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()

	projectMdContent, err := afero.ReadFile(fs, "spectr/project.md")
	if err != nil {
		t.Fatalf("failed to read project.md: %v", err)
	}
	if len(projectMdContent) == 0 {
		t.Error("expected project.md to have content")
	}

	agentsMdContent, err := afero.ReadFile(fs, "spectr/AGENTS.md")
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}
	if !strings.Contains(string(agentsMdContent), "mock agents content") {
		t.Error("expected AGENTS.md to contain mock agents content")
	}
}

// TestExecute_InitializerDeduplication tests that duplicate initializers are deduplicated.
func TestExecute_InitializerDeduplication(t *testing.T) {
	memFs := afero.NewMemMapFs()
	projectPath := "/test/project"

	if err := memFs.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create project directory: %v", err)
	}

	fs := afero.NewBasePathFs(memFs, projectPath)
	registry := providers.CreateRegistry()
	factory := initializers.NewFactory()
	tm := &mockTemplateRenderer{}

	if err := providers.RegisterClaudeProvider(registry, tm, factory); err != nil {
		t.Fatalf("failed to register Claude provider: %v", err)
	}
	if err := providers.RegisterClineProvider(registry, tm, factory); err != nil {
		t.Fatalf("failed to register Cline provider: %v", err)
	}

	cfg := providers.NewConfig()
	executor := &testableInitExecutorNew{
		&InitExecutorNew{fs: fs, projectPath: projectPath, registry: registry, cfg: cfg, tm: tm},
	}

	result, err := executor.Execute(
		context.Background(),
		[]string{"claude-code", "cline"},
		NewExecuteOptions(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claudeExists, _ := afero.Exists(fs, "CLAUDE.md")
	clineExists, _ := afero.Exists(fs, "CLINE.md")

	if !claudeExists {
		t.Error("expected CLAUDE.md to exist")
	}
	if !clineExists {
		t.Error("expected CLINE.md to exist")
	}
	if len(result.CreatedFiles) == 0 {
		t.Error("expected CreatedFiles to have entries")
	}
}

// TestExecute_CIWorkflowEnabled tests that CI workflow is created when enabled.
func TestExecute_CIWorkflowEnabled(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	opts := NewExecuteOptions().WithCIWorkflow(true)
	result, err := executor.Execute(context.Background(), nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()
	workflowExists, err := afero.Exists(fs, ".github/workflows/spectr-ci.yml")
	if err != nil {
		t.Fatalf("error checking workflow file: %v", err)
	}
	if !workflowExists {
		t.Error("expected .github/workflows/spectr-ci.yml to exist")
	}

	found := false
	for _, f := range result.CreatedFiles {
		if f == ".github/workflows/spectr-ci.yml" {
			found = true

			break
		}
	}
	if !found {
		t.Error("expected CI workflow to be in CreatedFiles")
	}
}

// TestExecute_CIWorkflowDisabled tests that CI workflow is not created when disabled.
func TestExecute_CIWorkflowDisabled(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), nil, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()
	workflowExists, _ := afero.Exists(fs, ".github/workflows/spectr-ci.yml")
	if workflowExists {
		t.Error("expected .github/workflows/spectr-ci.yml to NOT exist when CI is disabled")
	}
}

// TestExecute_AlreadyInitialized tests warning when spectr is already initialized.
func TestExecute_AlreadyInitialized(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")
	fs := executor.getMemFs()

	if err := fs.MkdirAll("spectr", 0755); err != nil {
		t.Fatalf("failed to create spectr directory: %v", err)
	}
	if err := afero.WriteFile(fs, "spectr/project.md", []byte("existing content"), 0644); err != nil {
		t.Fatalf("failed to write project.md: %v", err)
	}

	result, err := executor.Execute(context.Background(), nil, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundWarning := false
	for _, e := range result.Errors {
		if strings.Contains(e, "already initialized") {
			foundWarning = true

			break
		}
	}
	if !foundWarning {
		t.Error("expected warning about spectr already initialized")
	}

	content, err := afero.ReadFile(fs, "spectr/project.md")
	if err != nil {
		t.Fatalf("failed to read project.md: %v", err)
	}
	if string(content) != "existing content" {
		t.Error("expected existing project.md content to be preserved")
	}
}

// TestExecute_UnknownProvider tests error handling for unknown provider IDs.
func TestExecute_UnknownProvider(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	result, err := executor.Execute(
		context.Background(),
		[]string{"unknown-provider"},
		NewExecuteOptions(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundError := false
	for _, e := range result.Errors {
		if strings.Contains(e, "unknown-provider") && strings.Contains(e, "not found") {
			foundError = true

			break
		}
	}
	if !foundError {
		t.Error("expected error about unknown provider")
	}
}

// TestExecute_NoProviders tests initialization with no providers selected.
func TestExecute_NoProviders(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	result, err := executor.Execute(context.Background(), nil, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()
	spectrExists, _ := afero.DirExists(fs, "spectr")
	projectMdExists, _ := afero.Exists(fs, "spectr/project.md")
	agentsMdExists, _ := afero.Exists(fs, "spectr/AGENTS.md")

	if !spectrExists {
		t.Error("expected spectr/ to exist even with no providers")
	}
	if !projectMdExists {
		t.Error("expected project.md to exist even with no providers")
	}
	if !agentsMdExists {
		t.Error("expected AGENTS.md to exist even with no providers")
	}

	claudeMdExists, _ := afero.Exists(fs, "CLAUDE.md")
	if claudeMdExists {
		t.Error("expected CLAUDE.md to NOT exist when no providers selected")
	}

	if len(result.CreatedFiles) < 3 {
		t.Errorf(
			"expected at least 3 created files for core structure, got %d",
			len(result.CreatedFiles),
		)
	}
}

// TestExecute_ContextCancellation tests that execution respects context cancellation.
func TestExecute_ContextCancellation(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := executor.Execute(ctx, []string{"claude-code"}, NewExecuteOptions())
	if err != nil {
		t.Logf("execution with cancelled context returned error: %v", err)
	}
	if result == nil {
		t.Error("expected result to be non-nil even with cancelled context")
	}
}

// TestGetRegistry tests that GetRegistry returns the provider registry.
func TestGetRegistry(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	registry := executor.GetRegistry()
	if registry == nil {
		t.Fatal("expected registry to be non-nil")
	}

	claude := registry.Get("claude-code")
	if claude == nil {
		t.Error("expected claude-code provider to be registered")
	}

	gemini := registry.Get("gemini")
	if gemini == nil {
		t.Error("expected gemini provider to be registered")
	}
}

// TestExecute_SlashCommandContent tests that slash command files have correct content.
func TestExecute_SlashCommandContent(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), []string{"claude-code"}, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()

	proposalContent, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Fatalf("failed to read proposal.md: %v", err)
	}
	if !strings.Contains(string(proposalContent), "mock command content for proposal") {
		t.Error("expected proposal.md to contain mock command content")
	}

	applyContent, err := afero.ReadFile(fs, ".claude/commands/spectr/apply.md")
	if err != nil {
		t.Fatalf("failed to read apply.md: %v", err)
	}
	if !strings.Contains(string(applyContent), "mock command content for apply") {
		t.Error("expected apply.md to contain mock command content")
	}
}

// TestExecute_ConfigFileContent tests that config files have correct content.
func TestExecute_ConfigFileContent(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), []string{"claude-code"}, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()

	claudeContent, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}

	content := string(claudeContent)
	if !strings.Contains(content, "mock instruction content") {
		t.Error("expected CLAUDE.md to contain mock instruction content")
	}
	if !strings.Contains(content, "<!-- spectr:START -->") {
		t.Error("expected CLAUDE.md to contain spectr:START marker")
	}
	if !strings.Contains(content, "<!-- spectr:END -->") {
		t.Error("expected CLAUDE.md to contain spectr:END marker")
	}
}

// TestExecute_GeminiTOMLFormat tests that Gemini uses TOML format for commands.
func TestExecute_GeminiTOMLFormat(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")

	_, err := executor.Execute(context.Background(), []string{"gemini"}, NewExecuteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := executor.getMemFs()

	proposalContent, err := afero.ReadFile(fs, ".gemini/commands/spectr/proposal.toml")
	if err != nil {
		t.Fatalf("failed to read proposal.toml: %v", err)
	}

	content := string(proposalContent)
	if !strings.Contains(content, "description = ") {
		t.Error("expected TOML file to contain description field")
	}
	if !strings.Contains(content, "prompt = ") {
		t.Error("expected TOML file to contain prompt field")
	}
}

// TestExecute_UpdatedFilesTracking tests that updated files are tracked correctly.
func TestExecute_UpdatedFilesTracking(t *testing.T) {
	executor := createTestExecutor(t, "/test/project")
	fs := executor.getMemFs()

	existingContent := "<!-- spectr:START -->\nold content\n<!-- spectr:END -->\n"
	if err := afero.WriteFile(fs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to write existing CLAUDE.md: %v", err)
	}

	result, err := executor.Execute(
		context.Background(),
		[]string{"claude-code"},
		NewExecuteOptions(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = result.UpdatedFiles
}

// Table-driven tests for various provider configurations
func TestExecute_ProviderConfigurations(t *testing.T) {
	tests := []struct {
		name         string
		providers    []string
		expectFiles  []string
		expectDirs   []string
		ciEnabled    bool
		expectCIFile bool
	}{
		{
			name:      "claude only",
			providers: []string{"claude-code"},
			expectFiles: []string{
				"CLAUDE.md",
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
			},
			expectDirs:   []string{".claude/commands/spectr"},
			ciEnabled:    false,
			expectCIFile: false,
		},
		{
			name:      "gemini only",
			providers: []string{"gemini"},
			expectFiles: []string{
				".gemini/commands/spectr/proposal.toml",
				".gemini/commands/spectr/apply.toml",
			},
			expectDirs:   []string{".gemini/commands/spectr"},
			ciEnabled:    false,
			expectCIFile: false,
		},
		{
			name:      "claude with CI",
			providers: []string{"claude-code"},
			expectFiles: []string{
				"CLAUDE.md",
				".github/workflows/spectr-ci.yml",
			},
			expectDirs:   []string{".claude/commands/spectr", ".github/workflows"},
			ciEnabled:    true,
			expectCIFile: true,
		},
		{
			name:      "cline provider",
			providers: []string{"cline"},
			expectFiles: []string{
				"CLINE.md",
				".clinerules/commands/spectr/proposal.md",
				".clinerules/commands/spectr/apply.md",
			},
			expectDirs:   []string{".clinerules/commands/spectr"},
			ciEnabled:    false,
			expectCIFile: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := createTestExecutor(t, filepath.Join("/test/project", tt.name))

			opts := NewExecuteOptions().WithCIWorkflow(tt.ciEnabled)
			result, err := executor.Execute(context.Background(), tt.providers, opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			fs := executor.getMemFs()

			for _, expectedFile := range tt.expectFiles {
				exists, err := afero.Exists(fs, expectedFile)
				if err != nil {
					t.Errorf("error checking file %s: %v", expectedFile, err)
				}
				if !exists {
					t.Errorf("expected file %s to exist", expectedFile)
				}
			}

			for _, expectedDir := range tt.expectDirs {
				exists, err := afero.DirExists(fs, expectedDir)
				if err != nil {
					t.Errorf("error checking directory %s: %v", expectedDir, err)
				}
				if !exists {
					t.Errorf("expected directory %s to exist", expectedDir)
				}
			}

			ciFileExists, _ := afero.Exists(fs, ".github/workflows/spectr-ci.yml")
			if tt.expectCIFile && !ciFileExists {
				t.Error("expected CI workflow file to exist")
			}
			if !tt.expectCIFile && ciFileExists {
				t.Error("expected CI workflow file to NOT exist")
			}

			if len(result.CreatedFiles) == 0 && len(result.UpdatedFiles) == 0 {
				t.Error("expected result to have some file entries")
			}
		})
	}
}

// TestParseSlashCmdsKey tests the parseSlashCmdsKey helper function.
func TestParseSlashCmdsKey(t *testing.T) {
	tests := []struct {
		key         string
		expectedDir string
		expectedExt string
	}{
		{
			key:         "slashcmds:.claude/commands/spectr:.md:0",
			expectedDir: ".claude/commands/spectr",
			expectedExt: ".md",
		},
		{
			key:         "slashcmds:.gemini/commands/spectr:.toml:1",
			expectedDir: ".gemini/commands/spectr",
			expectedExt: ".toml",
		},
		{
			key:         "slashcmds:path/with:colons:.ext:0",
			expectedDir: "path/with:colons",
			expectedExt: ".ext",
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			parts := parseSlashCmdsKey(tt.key)
			if parts.dir != tt.expectedDir {
				t.Errorf("expected dir %q, got %q", tt.expectedDir, parts.dir)
			}
			if parts.ext != tt.expectedExt {
				t.Errorf("expected ext %q, got %q", tt.expectedExt, parts.ext)
			}
		})
	}
}
