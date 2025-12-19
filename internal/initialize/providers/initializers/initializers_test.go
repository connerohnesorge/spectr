// Package initializers provides unit tests for the built-in Initializer
// implementations.
//
// All tests use afero.MemMapFs for in-memory filesystem testing.
//
//nolint:revive // bool-literal-in-expr, cyclomatic - test code
package initializers

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// mockTemplateRenderer implements providers.TemplateRenderer for testing.
// It returns predictable content for each template type.
type mockTemplateRenderer struct {
	agentsContent             string
	instructionPointerContent string
	slashCommandContent       map[string]string
}

func newMockTemplateRenderer() *mockTemplateRenderer {
	return &mockTemplateRenderer{
		agentsContent:             "# AGENTS.md Content\nTest agents content",
		instructionPointerContent: "# Spectr Instructions\nRead spectr/AGENTS.md for details.",
		slashCommandContent: map[string]string{
			"proposal": "Create a new change proposal for the requested feature.",
			"apply":    "Apply the specified change proposal.",
		},
	}
}

func (m *mockTemplateRenderer) RenderAgents(ctx providers.TemplateContext) (string, error) {
	return m.agentsContent, nil
}

func (m *mockTemplateRenderer) RenderInstructionPointer(
	ctx providers.TemplateContext,
) (string, error) {
	return m.instructionPointerContent, nil
}

func (m *mockTemplateRenderer) RenderSlashCommand(
	command string,
	ctx providers.TemplateContext,
) (string, error) {
	if content, ok := m.slashCommandContent[command]; ok {
		return content, nil
	}

	return "", nil
}

// newTestConfig creates a Config for testing.
func newTestConfig() *providers.Config {
	return providers.NewConfig("spectr")
}

// ============================================================================
// DirectoryInitializer Tests
// ============================================================================

func TestNewDirectoryInitializer(t *testing.T) {
	t.Run("returns nil for empty paths", func(t *testing.T) {
		di := NewDirectoryInitializer(false)
		if di != nil {
			t.Error("expected nil for empty paths")
		}
	})

	t.Run("creates initializer for single path", func(t *testing.T) {
		di := NewDirectoryInitializer(false, ".claude/commands")
		if di == nil {
			t.Fatal("expected non-nil initializer")
		}
		if di.Path() != ".claude/commands" {
			t.Errorf("Path() = %q, want %q", di.Path(), ".claude/commands")
		}
		if di.IsGlobal() != false {
			t.Error("IsGlobal() should be false")
		}
	})

	t.Run("creates initializer for multiple paths", func(t *testing.T) {
		di := NewDirectoryInitializer(true, ".config/tool/commands", ".config/tool/data")
		if di == nil {
			t.Fatal("expected non-nil initializer")
		}
		if len(di.Paths()) != 2 {
			t.Errorf("Paths() length = %d, want 2", len(di.Paths()))
		}
		if di.IsGlobal() != true {
			t.Error("IsGlobal() should be true")
		}
	})
}

func TestDirectoryInitializer_Init(t *testing.T) {
	t.Run("creates single directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, ".claude/commands/spectr")

		err := di.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Verify directory was created
		info, err := fs.Stat(".claude/commands/spectr")
		if err != nil {
			t.Fatalf("directory not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected directory, got file")
		}
	})

	t.Run("creates multiple directories", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, ".claude/commands/spectr", ".claude/data")

		err := di.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Verify both directories were created
		for _, path := range []string{".claude/commands/spectr", ".claude/data"} {
			info, err := fs.Stat(path)
			if err != nil {
				t.Errorf("directory %q not created: %v", path, err)

				continue
			}
			if !info.IsDir() {
				t.Errorf("expected directory at %q, got file", path)
			}
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, "a/b/c/d/e")

		err := di.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Verify all parent directories were created
		for _, path := range []string{"a", "a/b", "a/b/c", "a/b/c/d", "a/b/c/d/e"} {
			info, err := fs.Stat(path)
			if err != nil {
				t.Errorf("directory %q not created: %v", path, err)

				continue
			}
			if !info.IsDir() {
				t.Errorf("expected directory at %q, got file", path)
			}
		}
	})

	t.Run("idempotency - run twice same result", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, ".claude/commands/spectr")

		// Run Init twice
		err := di.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("first Init failed: %v", err)
		}

		err = di.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("second Init failed: %v", err)
		}

		// Verify directory still exists and is valid
		info, err := fs.Stat(".claude/commands/spectr")
		if err != nil {
			t.Fatalf("directory not found after second Init: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected directory, got file")
		}
	})
}

func TestDirectoryInitializer_IsSetup(t *testing.T) {
	t.Run("returns false when directory does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, ".claude/commands/spectr")

		if di.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return false for non-existent directory")
		}
	})

	t.Run("returns true when directory exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, ".claude/commands/spectr")

		// Create the directory
		if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}

		if !di.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return true for existing directory")
		}
	})

	t.Run("returns false when path is a file not directory", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, "test-path")

		// Create a file at the path
		if err := afero.WriteFile(fs, "test-path", []byte("file content"), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}

		if di.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return false when path is a file")
		}
	})

	t.Run("returns false when any directory is missing", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, "dir1", "dir2", "dir3")

		// Create only dir1 and dir3
		if err := fs.MkdirAll("dir1", 0755); err != nil {
			t.Fatalf("failed to create dir1: %v", err)
		}
		if err := fs.MkdirAll("dir3", 0755); err != nil {
			t.Fatalf("failed to create dir3: %v", err)
		}

		if di.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return false when dir2 is missing")
		}
	})

	t.Run("returns true when all directories exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		di := NewDirectoryInitializer(false, "dir1", "dir2", "dir3")

		// Create all directories
		for _, dir := range []string{"dir1", "dir2", "dir3"} {
			if err := fs.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("failed to create %s: %v", dir, err)
			}
		}

		if !di.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return true when all directories exist")
		}
	})
}

// ============================================================================
// ConfigFileInitializer Tests
// ============================================================================

func TestNewConfigFileInitializer(t *testing.T) {
	t.Run("returns nil for empty path", func(t *testing.T) {
		ci := NewConfigFileInitializer("", "instruction-pointer", false)
		if ci != nil {
			t.Error("expected nil for empty path")
		}
	})

	t.Run("creates initializer with valid path", func(t *testing.T) {
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)
		if ci == nil {
			t.Fatal("expected non-nil initializer")
		}
		if ci.Path() != "CLAUDE.md" {
			t.Errorf("Path() = %q, want %q", ci.Path(), "CLAUDE.md")
		}
		if ci.TemplateName() != "instruction-pointer" {
			t.Errorf("TemplateName() = %q, want %q", ci.TemplateName(), "instruction-pointer")
		}
		if ci.IsGlobal() != false {
			t.Error("IsGlobal() should be false")
		}
	})

	t.Run("creates global initializer", func(t *testing.T) {
		ci := NewConfigFileInitializer(".config/tool/config.md", "instruction-pointer", true)
		if ci == nil {
			t.Fatal("expected non-nil initializer")
		}
		if ci.IsGlobal() != true {
			t.Error("IsGlobal() should be true")
		}
	})
}

func TestConfigFileInitializer_Init(t *testing.T) {
	t.Run("creates new file with markers", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)

		err := ci.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Read the created file
		content, err := afero.ReadFile(fs, "CLAUDE.md")
		if err != nil {
			t.Fatalf("failed to read created file: %v", err)
		}

		contentStr := string(content)

		// Check for markers
		if !strings.Contains(contentStr, "<!-- spectr:START -->") {
			t.Error("file missing start marker")
		}
		if !strings.Contains(contentStr, "<!-- spectr:END -->") {
			t.Error("file missing end marker")
		}

		// Check for rendered content
		if !strings.Contains(contentStr, "# Spectr Instructions") {
			t.Error("file missing rendered template content")
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("path/to/nested/CLAUDE.md", "instruction-pointer", false)

		err := ci.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Verify file was created in nested directory
		exists, err := afero.Exists(fs, "path/to/nested/CLAUDE.md")
		if err != nil {
			t.Fatalf("failed to check file existence: %v", err)
		}
		if !exists {
			t.Error("file was not created in nested directory")
		}
	})

	t.Run("updates existing file with markers", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)

		// Create initial file with markers and old content
		initialContent := `# My Project

Some existing content here.

<!-- spectr:START -->
Old spectr content that should be replaced.
<!-- spectr:END -->

More content after markers.
`
		if err := afero.WriteFile(fs, "CLAUDE.md", []byte(initialContent), 0644); err != nil {
			t.Fatalf("failed to create initial file: %v", err)
		}

		// Run Init to update
		err := ci.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Read updated file
		content, err := afero.ReadFile(fs, "CLAUDE.md")
		if err != nil {
			t.Fatalf("failed to read updated file: %v", err)
		}

		contentStr := string(content)

		// Check that existing content is preserved
		if !strings.Contains(contentStr, "# My Project") {
			t.Error("existing content before markers was not preserved")
		}
		if !strings.Contains(contentStr, "Some existing content here.") {
			t.Error("existing content before markers was not preserved")
		}
		if !strings.Contains(contentStr, "More content after markers.") {
			t.Error("existing content after markers was not preserved")
		}

		// Check that old content between markers was replaced
		if strings.Contains(contentStr, "Old spectr content that should be replaced.") {
			t.Error("old content between markers was not replaced")
		}

		// Check that new content is present
		if !strings.Contains(contentStr, "# Spectr Instructions") {
			t.Error("new template content was not added")
		}
	})

	t.Run("appends markers to file without them", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)

		// Create initial file without markers
		initialContent := `# My Project

Some existing content without spectr markers.
`
		if err := afero.WriteFile(fs, "CLAUDE.md", []byte(initialContent), 0644); err != nil {
			t.Fatalf("failed to create initial file: %v", err)
		}

		// Run Init to append markers
		err := ci.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Read updated file
		content, err := afero.ReadFile(fs, "CLAUDE.md")
		if err != nil {
			t.Fatalf("failed to read updated file: %v", err)
		}

		contentStr := string(content)

		// Check that original content is preserved
		if !strings.Contains(contentStr, "# My Project") {
			t.Error("original content was not preserved")
		}
		if !strings.Contains(contentStr, "Some existing content without spectr markers.") {
			t.Error("original content was not preserved")
		}

		// Check that markers and new content were appended
		if !strings.Contains(contentStr, "<!-- spectr:START -->") {
			t.Error("start marker was not appended")
		}
		if !strings.Contains(contentStr, "<!-- spectr:END -->") {
			t.Error("end marker was not appended")
		}
		if !strings.Contains(contentStr, "# Spectr Instructions") {
			t.Error("template content was not appended")
		}

		// Verify markers appear after original content
		originalContentIndex := strings.Index(contentStr, "# My Project")
		startMarkerIndex := strings.Index(contentStr, "<!-- spectr:START -->")
		if startMarkerIndex < originalContentIndex {
			t.Error("markers should be appended after original content")
		}
	})

	t.Run("idempotency - run twice same result", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)

		// Run Init twice
		err := ci.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("first Init failed: %v", err)
		}

		firstContent, err := afero.ReadFile(fs, "CLAUDE.md")
		if err != nil {
			t.Fatalf("failed to read file after first Init: %v", err)
		}

		err = ci.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("second Init failed: %v", err)
		}

		secondContent, err := afero.ReadFile(fs, "CLAUDE.md")
		if err != nil {
			t.Fatalf("failed to read file after second Init: %v", err)
		}

		// Content should be the same after running twice
		if string(firstContent) != string(secondContent) {
			t.Error("file content changed after second Init")
		}
	})
}

func TestConfigFileInitializer_IsSetup(t *testing.T) {
	t.Run("returns false when file does not exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)

		if ci.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return false for non-existent file")
		}
	})

	t.Run("returns true when file exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		ci := NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false)

		// Create the file
		if err := afero.WriteFile(fs, "CLAUDE.md", []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}

		if !ci.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return true for existing file")
		}
	})
}

// ============================================================================
// SlashCommandsInitializer Tests
// ============================================================================

func TestNewSlashCommandsInitializer(t *testing.T) {
	t.Run("returns nil for empty dir", func(t *testing.T) {
		si := NewSlashCommandsInitializer("", ".md", providers.FormatMarkdown, nil, false)
		if si != nil {
			t.Error("expected nil for empty dir")
		}
	})

	t.Run("creates initializer with valid dir", func(t *testing.T) {
		frontmatter := map[string]string{
			"proposal": "---\ndescription: Create proposal\n---",
			"apply":    "---\ndescription: Apply proposal\n---",
		}
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			frontmatter,
			false,
		)
		if si == nil {
			t.Fatal("expected non-nil initializer")
		}
		if si.Path() != ".claude/commands/spectr" {
			t.Errorf("Path() = %q, want %q", si.Path(), ".claude/commands/spectr")
		}
		if si.Dir() != ".claude/commands/spectr" {
			t.Errorf("Dir() = %q, want %q", si.Dir(), ".claude/commands/spectr")
		}
		if si.Ext() != ".md" {
			t.Errorf("Ext() = %q, want %q", si.Ext(), ".md")
		}
		if si.Format() != providers.FormatMarkdown {
			t.Errorf("Format() = %v, want FormatMarkdown", si.Format())
		}
		if si.IsGlobal() != false {
			t.Error("IsGlobal() should be false")
		}
	})

	t.Run("creates global initializer with TOML format", func(t *testing.T) {
		si := NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			providers.FormatTOML,
			nil,
			true,
		)
		if si == nil {
			t.Fatal("expected non-nil initializer")
		}
		if si.Format() != providers.FormatTOML {
			t.Errorf("Format() = %v, want FormatTOML", si.Format())
		}
		if si.IsGlobal() != true {
			t.Error("IsGlobal() should be true")
		}
	})
}

func TestSlashCommandsInitializer_Init_Markdown(t *testing.T) {
	t.Run("creates markdown command files", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		frontmatter := map[string]string{
			"proposal": "---\ndescription: Create a new change proposal\nallowed-tools: Read, Write, Edit, Glob, Grep\n---",
			"apply":    "---\ndescription: Apply a change proposal\nallowed-tools: Read, Write, Edit, Bash\n---",
		}
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			frontmatter,
			false,
		)

		err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Check both command files were created
		for _, cmd := range []string{"proposal", "apply"} {
			filePath := filepath.Join(".claude/commands/spectr", cmd+".md")
			content, err := afero.ReadFile(fs, filePath)
			if err != nil {
				t.Errorf("failed to read %s: %v", filePath, err)

				continue
			}

			contentStr := string(content)

			// Check for frontmatter
			if !strings.Contains(contentStr, "---") {
				t.Errorf("%s missing frontmatter", filePath)
			}
			if !strings.Contains(contentStr, "description:") {
				t.Errorf("%s missing description in frontmatter", filePath)
			}

			// Check for markers
			if !strings.Contains(contentStr, "<!-- spectr:START -->") {
				t.Errorf("%s missing start marker", filePath)
			}
			if !strings.Contains(contentStr, "<!-- spectr:END -->") {
				t.Errorf("%s missing end marker", filePath)
			}
		}
	})

	t.Run("creates directory if not exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			nil,
			false,
		)

		err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Check directory was created
		info, err := fs.Stat(".claude/commands/spectr")
		if err != nil {
			t.Fatalf("directory not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected directory")
		}
	})
}

func TestSlashCommandsInitializer_Init_TOML(t *testing.T) {
	t.Run("creates toml command files", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		frontmatter := map[string]string{
			"proposal": "description = \"Create a new change proposal\"",
			"apply":    "description = \"Apply a change proposal\"",
		}
		si := NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			providers.FormatTOML,
			frontmatter,
			false,
		)

		err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Check both command files were created
		for _, cmd := range []string{"proposal", "apply"} {
			filePath := filepath.Join(".gemini/commands/spectr", cmd+".toml")
			content, err := afero.ReadFile(fs, filePath)
			if err != nil {
				t.Errorf("failed to read %s: %v", filePath, err)

				continue
			}

			contentStr := string(content)

			// Check for TOML frontmatter (description = "...")
			if !strings.Contains(contentStr, "description = ") {
				t.Errorf("%s missing TOML description", filePath)
			}

			// Check for markers
			if !strings.Contains(contentStr, "<!-- spectr:START -->") {
				t.Errorf("%s missing start marker", filePath)
			}
			if !strings.Contains(contentStr, "<!-- spectr:END -->") {
				t.Errorf("%s missing end marker", filePath)
			}
		}
	})
}

func TestSlashCommandsInitializer_Init_Update(t *testing.T) {
	t.Run("updates existing command files", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		// Create existing proposal file with old content
		existingContent := `---
description: Create a new change proposal
allowed-tools: Read, Write, Edit, Glob, Grep
---

<!-- spectr:START -->

Old proposal content that should be replaced.

<!-- spectr:END -->
`
		if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := afero.WriteFile(fs, ".claude/commands/spectr/proposal.md", []byte(existingContent), 0644); err != nil {
			t.Fatalf("failed to create existing file: %v", err)
		}

		// Also create apply file
		existingApply := `---
description: Apply a change proposal
---

<!-- spectr:START -->

Old apply content.

<!-- spectr:END -->
`
		if err := afero.WriteFile(fs, ".claude/commands/spectr/apply.md", []byte(existingApply), 0644); err != nil {
			t.Fatalf("failed to create apply file: %v", err)
		}

		frontmatter := map[string]string{
			"proposal": "---\ndescription: Create a new change proposal\nallowed-tools: Read, Write, Edit, Glob, Grep\n---",
			"apply":    "---\ndescription: Apply a change proposal\n---",
		}
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			frontmatter,
			false,
		)

		err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Check proposal file was updated
		content, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
		if err != nil {
			t.Fatalf("failed to read proposal file: %v", err)
		}

		contentStr := string(content)

		// Old content should be replaced
		if strings.Contains(contentStr, "Old proposal content that should be replaced.") {
			t.Error("old content was not replaced")
		}

		// Frontmatter should be preserved
		if !strings.Contains(contentStr, "description: Create a new change proposal") {
			t.Error("frontmatter was not preserved")
		}

		// New content from template should be present
		if !strings.Contains(
			contentStr,
			"Create a new change proposal for the requested feature.",
		) {
			t.Error("new template content was not added")
		}
	})

	t.Run("adds frontmatter to existing file without it", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		// Create existing proposal file WITHOUT frontmatter
		existingContent := `<!-- spectr:START -->

Old proposal content without frontmatter.

<!-- spectr:END -->
`
		if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := afero.WriteFile(fs, ".claude/commands/spectr/proposal.md", []byte(existingContent), 0644); err != nil {
			t.Fatalf("failed to create existing file: %v", err)
		}

		// Also create apply file without frontmatter
		existingApply := `<!-- spectr:START -->
Old apply content.
<!-- spectr:END -->
`
		if err := afero.WriteFile(fs, ".claude/commands/spectr/apply.md", []byte(existingApply), 0644); err != nil {
			t.Fatalf("failed to create apply file: %v", err)
		}

		frontmatter := map[string]string{
			"proposal": "---\ndescription: Create a new change proposal\n---",
			"apply":    "---\ndescription: Apply a change proposal\n---",
		}
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			frontmatter,
			false,
		)

		err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Check proposal file has frontmatter added
		content, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
		if err != nil {
			t.Fatalf("failed to read proposal file: %v", err)
		}

		contentStr := string(content)

		// Frontmatter should be added
		if !strings.Contains(contentStr, "---") {
			t.Error("frontmatter was not added")
		}
		if !strings.Contains(contentStr, "description: Create a new change proposal") {
			t.Error("frontmatter description was not added")
		}

		// Markers should still be present
		if !strings.Contains(contentStr, "<!-- spectr:START -->") {
			t.Error("start marker was removed")
		}
		if !strings.Contains(contentStr, "<!-- spectr:END -->") {
			t.Error("end marker was removed")
		}
	})
}

func TestSlashCommandsInitializer_IsSetup(t *testing.T) {
	t.Run("returns false when no files exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			nil,
			false,
		)

		if si.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return false when no files exist")
		}
	})

	t.Run("returns false when only one file exists", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			nil,
			false,
		)

		// Create only proposal file
		if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := afero.WriteFile(fs, ".claude/commands/spectr/proposal.md", []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}

		if si.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return false when apply file is missing")
		}
	})

	t.Run("returns true when both files exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			nil,
			false,
		)

		// Create both files
		if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := afero.WriteFile(fs, ".claude/commands/spectr/proposal.md", []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create proposal file: %v", err)
		}
		if err := afero.WriteFile(fs, ".claude/commands/spectr/apply.md", []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create apply file: %v", err)
		}

		if !si.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return true when both files exist")
		}
	})

	t.Run("returns true for TOML files", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		si := NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			providers.FormatTOML,
			nil,
			false,
		)

		// Create both TOML files
		if err := fs.MkdirAll(".gemini/commands/spectr", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := afero.WriteFile(fs, ".gemini/commands/spectr/proposal.toml", []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create proposal file: %v", err)
		}
		if err := afero.WriteFile(fs, ".gemini/commands/spectr/apply.toml", []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create apply file: %v", err)
		}

		if !si.IsSetup(fs, newTestConfig()) {
			t.Error("IsSetup() should return true for TOML files")
		}
	})
}

func TestSlashCommandsInitializer_Idempotency(t *testing.T) {
	t.Run("running init twice produces same functional result", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		frontmatter := map[string]string{
			"proposal": "---\ndescription: Create a new change proposal\n---",
			"apply":    "---\ndescription: Apply a change proposal\n---",
		}
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			frontmatter,
			false,
		)

		// Run Init twice
		err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("first Init failed: %v", err)
		}

		err = si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("second Init failed: %v", err)
		}

		// Run a third time to ensure stability
		err = si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("third Init failed: %v", err)
		}

		// Read files after third Init
		proposal3, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
		if err != nil {
			t.Fatalf("failed to read proposal after third Init: %v", err)
		}
		apply3, err := afero.ReadFile(fs, ".claude/commands/spectr/apply.md")
		if err != nil {
			t.Fatalf("failed to read apply after third Init: %v", err)
		}

		err = si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
		if err != nil {
			t.Fatalf("fourth Init failed: %v", err)
		}

		// Read files after fourth Init
		proposal4, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
		if err != nil {
			t.Fatalf("failed to read proposal after fourth Init: %v", err)
		}
		apply4, err := afero.ReadFile(fs, ".claude/commands/spectr/apply.md")
		if err != nil {
			t.Fatalf("failed to read apply after fourth Init: %v", err)
		}

		// After initial normalization (first run creates, second normalizes),
		// subsequent runs should be stable
		if string(proposal3) != string(proposal4) {
			t.Errorf(
				"proposal content not stable\n--- THIRD:\n%s\n--- FOURTH:\n%s",
				string(proposal3),
				string(proposal4),
			)
		}
		if string(apply3) != string(apply4) {
			t.Errorf(
				"apply content not stable\n--- THIRD:\n%s\n--- FOURTH:\n%s",
				string(apply3),
				string(apply4),
			)
		}
	})

	t.Run("essential content preserved across runs", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		frontmatter := map[string]string{
			"proposal": "---\ndescription: Create a new change proposal\n---",
			"apply":    "---\ndescription: Apply a change proposal\n---",
		}
		si := NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			providers.FormatMarkdown,
			frontmatter,
			false,
		)

		// Run Init multiple times
		for i := range 3 {
			err := si.Init(context.Background(), fs, newTestConfig(), newMockTemplateRenderer())
			if err != nil {
				t.Fatalf("Init %d failed: %v", i+1, err)
			}

			// Verify essential content is present after each run
			proposal, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
			if err != nil {
				t.Fatalf("failed to read proposal after Init %d: %v", i+1, err)
			}

			proposalStr := string(proposal)
			if !strings.Contains(proposalStr, "<!-- spectr:START -->") {
				t.Errorf("run %d: proposal missing start marker", i+1)
			}
			if !strings.Contains(proposalStr, "<!-- spectr:END -->") {
				t.Errorf("run %d: proposal missing end marker", i+1)
			}
			if !strings.Contains(proposalStr, "description: Create a new change proposal") {
				t.Errorf("run %d: proposal missing frontmatter", i+1)
			}
			if !strings.Contains(
				proposalStr,
				"Create a new change proposal for the requested feature.",
			) {
				t.Errorf("run %d: proposal missing template content", i+1)
			}
		}
	})
}
