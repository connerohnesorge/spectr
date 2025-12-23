package providers

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestSlashCommandsInitializer_Init_Markdown(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	tm := &mockTemplateManager{
		content: "Run the proposal workflow",
	}

	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should create both proposal and apply files
	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Verify proposal.md exists
	exists, err := afero.Exists(fs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Fatalf("failed to check proposal.md existence: %v", err)
	}
	if !exists {
		t.Fatal("proposal.md does not exist")
	}

	// Verify apply.md exists
	exists, err = afero.Exists(fs, ".claude/commands/spectr/apply.md")
	if err != nil {
		t.Fatalf("failed to check apply.md existence: %v", err)
	}
	if !exists {
		t.Fatal("apply.md does not exist")
	}

	// Verify proposal.md content has YAML frontmatter
	content, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Fatalf("failed to read proposal.md: %v", err)
	}

	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---\n") {
		t.Error("proposal.md does not start with YAML frontmatter")
	}
	if !strings.Contains(contentStr, "description:") {
		t.Error("proposal.md frontmatter missing description field")
	}
	if !strings.Contains(contentStr, "Run the proposal workflow") {
		t.Error("proposal.md missing template content")
	}
}

func TestSlashCommandsInitializer_Init_TOML(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	tm := &mockTemplateManager{
		content: "Run the proposal workflow",
	}

	init := NewSlashCommandsInitializer(".gemini/commands/spectr", ".toml", FormatTOML)
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should create both proposal and apply files
	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Verify proposal.toml exists
	exists, err := afero.Exists(fs, ".gemini/commands/spectr/proposal.toml")
	if err != nil {
		t.Fatalf("failed to check proposal.toml existence: %v", err)
	}
	if !exists {
		t.Fatal("proposal.toml does not exist")
	}

	// Verify apply.toml exists
	exists, err = afero.Exists(fs, ".gemini/commands/spectr/apply.toml")
	if err != nil {
		t.Fatalf("failed to check apply.toml existence: %v", err)
	}
	if !exists {
		t.Fatal("apply.toml does not exist")
	}

	// Verify proposal.toml content
	content, err := afero.ReadFile(fs, ".gemini/commands/spectr/proposal.toml")
	if err != nil {
		t.Fatalf("failed to read proposal.toml: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "description =") {
		t.Error("proposal.toml missing description field")
	}
	if !strings.Contains(contentStr, "prompt =") {
		t.Error("proposal.toml missing prompt field")
	}
	if !strings.Contains(contentStr, "Run the proposal workflow") {
		t.Error("proposal.toml missing template content")
	}
	if !strings.Contains(contentStr, "# Spectr command") {
		t.Error("proposal.toml missing comment header")
	}
}

func TestSlashCommandsInitializer_Init_TOML_Escaping(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	// Template content with special characters that need escaping
	tm := &mockTemplateManager{
		content: `Test "quoted" content with \ backslash`,
	}

	init := NewSlashCommandsInitializer(".gemini/commands/spectr", ".toml", FormatTOML)
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify proposal.toml escapes special characters
	content, err := afero.ReadFile(fs, ".gemini/commands/spectr/proposal.toml")
	if err != nil {
		t.Fatalf("failed to read proposal.toml: %v", err)
	}

	contentStr := string(content)
	// Should escape backslashes and quotes
	if !strings.Contains(contentStr, `\\`) {
		t.Error("proposal.toml did not escape backslash")
	}
	if !strings.Contains(contentStr, `\"`) {
		t.Error("proposal.toml did not escape quotes")
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}
}

func TestSlashCommandsInitializer_Init_UpdateExisting(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	// Create existing files
	if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := afero.WriteFile(fs, ".claude/commands/spectr/proposal.md", []byte("old content"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}
	if err := afero.WriteFile(fs, ".claude/commands/spectr/apply.md", []byte("old content"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	tm := &mockTemplateManager{
		content: "New content",
	}

	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should report files as updated, not created
	if len(result.UpdatedFiles) != 2 {
		t.Errorf("Init() updated %d files, want 2", len(result.UpdatedFiles))
	}
	if len(result.CreatedFiles) != 0 {
		t.Errorf("Init() created %d files, want 0", len(result.CreatedFiles))
	}

	// Verify content was updated
	content, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Fatalf("failed to read proposal.md: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "New content") {
		t.Error("proposal.md content not updated")
	}
	if strings.Contains(contentStr, "old content") {
		t.Error("proposal.md old content still present")
	}
}

func TestSlashCommandsInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name        string
		createFiles []string
		want        bool
	}{
		{
			name:        "no files exist",
			createFiles: nil,
			want:        false,
		},
		{
			name:        "only proposal exists",
			createFiles: []string{".claude/commands/spectr/proposal.md"},
			want:        false,
		},
		{
			name:        "only apply exists",
			createFiles: []string{".claude/commands/spectr/apply.md"},
			want:        false,
		},
		{
			name: "both files exist",
			createFiles: []string{
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create files
			for _, file := range tt.createFiles {
				if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := afero.WriteFile(fs, file, []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create file %s: %v", file, err)
				}
			}

			init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
			if got := init.IsSetup(fs, cfg); got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_Path(t *testing.T) {
	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
	if got := init.Path(); got != ".claude/commands/spectr" {
		t.Errorf("Path() = %v, want .claude/commands/spectr", got)
	}
}

func TestSlashCommandsInitializer_IsGlobal(t *testing.T) {
	tests := []struct {
		name     string
		isGlobal bool
		want     bool
	}{
		{
			name:     "default is not global",
			isGlobal: false,
			want:     false,
		},
		{
			name:     "configured as global",
			isGlobal: true,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
			if tt.isGlobal {
				init = init.WithGlobal(true)
			}

			if got := init.IsGlobal(); got != tt.want {
				t.Errorf("IsGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_DirectoryCreation(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	tm := &mockTemplateManager{
		content: "Test content",
	}

	// Use a nested directory that doesn't exist yet
	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify directory was created
	dirExists, err := afero.DirExists(fs, ".claude/commands/spectr")
	if err != nil {
		t.Fatalf("failed to check directory existence: %v", err)
	}
	if !dirExists {
		t.Fatal("directory was not created")
	}

	// Verify files were created
	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}
}
