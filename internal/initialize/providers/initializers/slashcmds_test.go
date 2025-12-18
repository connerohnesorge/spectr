package initializers

import (
	"context"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

func TestSlashCommandsInitializer_Init_Markdown(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	tm.slashCommands = map[string]string{
		"proposal": "Proposal template content for test",
		"apply":    "Apply template content for test",
	}
	ctx := context.Background()

	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify proposal file was created
	proposalPath := ".claude/commands/spectr/proposal.md"
	exists, err := afero.Exists(fs, proposalPath)
	if err != nil {
		t.Errorf("Exists(proposal) error = %v", err)
		return
	}
	if !exists {
		t.Error("Proposal file was not created")
		return
	}

	// Verify apply file was created
	applyPath := ".claude/commands/spectr/apply.md"
	exists, err = afero.Exists(fs, applyPath)
	if err != nil {
		t.Errorf("Exists(apply) error = %v", err)
		return
	}
	if !exists {
		t.Error("Apply file was not created")
		return
	}

	// Verify proposal content
	proposalContent, err := afero.ReadFile(fs, proposalPath)
	if err != nil {
		t.Errorf("ReadFile(proposal) error = %v", err)
		return
	}

	proposalStr := string(proposalContent)
	if !strings.Contains(proposalStr, "Proposal template content for test") {
		t.Error("Proposal file should contain proposal template content")
	}
	if !strings.Contains(proposalStr, spectrStartMarker) {
		t.Error("Proposal file should contain start marker")
	}
	if !strings.Contains(proposalStr, spectrEndMarker) {
		t.Error("Proposal file should contain end marker")
	}

	// Verify apply content
	applyContent, err := afero.ReadFile(fs, applyPath)
	if err != nil {
		t.Errorf("ReadFile(apply) error = %v", err)
		return
	}

	applyStr := string(applyContent)
	if !strings.Contains(applyStr, "Apply template content for test") {
		t.Error("Apply file should contain apply template content")
	}
	if !strings.Contains(applyStr, spectrStartMarker) {
		t.Error("Apply file should contain start marker")
	}
	if !strings.Contains(applyStr, spectrEndMarker) {
		t.Error("Apply file should contain end marker")
	}
}

func TestSlashCommandsInitializer_Init_TOML(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	tm.slashCommands = map[string]string{
		"proposal": "Proposal TOML prompt content",
		"apply":    "Apply TOML prompt content",
	}
	ctx := context.Background()

	init := NewSlashCommandsInitializer(".gemini/commands/spectr", ".toml", FormatTOML)

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify proposal file was created
	proposalPath := ".gemini/commands/spectr/proposal.toml"
	exists, err := afero.Exists(fs, proposalPath)
	if err != nil {
		t.Errorf("Exists(proposal) error = %v", err)
		return
	}
	if !exists {
		t.Error("Proposal TOML file was not created")
		return
	}

	// Verify apply file was created
	applyPath := ".gemini/commands/spectr/apply.toml"
	exists, err = afero.Exists(fs, applyPath)
	if err != nil {
		t.Errorf("Exists(apply) error = %v", err)
		return
	}
	if !exists {
		t.Error("Apply TOML file was not created")
		return
	}

	// Verify proposal TOML content
	proposalContent, err := afero.ReadFile(fs, proposalPath)
	if err != nil {
		t.Errorf("ReadFile(proposal) error = %v", err)
		return
	}

	proposalStr := string(proposalContent)
	if !strings.Contains(proposalStr, "description = ") {
		t.Error("Proposal TOML should contain description field")
	}
	if !strings.Contains(proposalStr, "prompt = \"\"\"") {
		t.Error("Proposal TOML should contain prompt field with triple quotes")
	}
	if !strings.Contains(proposalStr, "Proposal TOML prompt content") {
		t.Error("Proposal TOML should contain prompt content")
	}
	if !strings.Contains(proposalStr, TomlDescriptionProposal) {
		t.Error("Proposal TOML should contain proposal description")
	}

	// Verify apply TOML content
	applyContent, err := afero.ReadFile(fs, applyPath)
	if err != nil {
		t.Errorf("ReadFile(apply) error = %v", err)
		return
	}

	applyStr := string(applyContent)
	if !strings.Contains(applyStr, "description = ") {
		t.Error("Apply TOML should contain description field")
	}
	if !strings.Contains(applyStr, "Apply TOML prompt content") {
		t.Error("Apply TOML should contain prompt content")
	}
	if !strings.Contains(applyStr, TomlDescriptionApply) {
		t.Error("Apply TOML should contain apply description")
	}
}

func TestSlashCommandsInitializer_Init_WithFrontmatter(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	tm.slashCommands = map[string]string{
		"proposal": "Proposal content",
		"apply":    "Apply content",
	}
	ctx := context.Background()

	frontmatter := map[string]string{
		"proposal": "---\ndescription: Custom proposal description\nallowed_tools: [\"Bash\", \"Read\"]\n---",
		"apply":    "---\ndescription: Custom apply description\n---",
	}

	init := NewSlashCommandsInitializerWithFrontmatter(".claude/commands/spectr", ".md", FormatMarkdown, frontmatter)

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify proposal file has frontmatter
	proposalContent, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Errorf("ReadFile(proposal) error = %v", err)
		return
	}

	proposalStr := string(proposalContent)
	if !strings.Contains(proposalStr, "description: Custom proposal description") {
		t.Error("Proposal should contain custom frontmatter description")
	}
	if !strings.Contains(proposalStr, "allowed_tools:") {
		t.Error("Proposal should contain allowed_tools from frontmatter")
	}

	// Verify frontmatter comes before markers
	frontmatterIndex := strings.Index(proposalStr, "---")
	markerIndex := strings.Index(proposalStr, spectrStartMarker)
	if frontmatterIndex > markerIndex {
		t.Error("Frontmatter should come before markers")
	}

	// Verify apply file has frontmatter
	applyContent, err := afero.ReadFile(fs, ".claude/commands/spectr/apply.md")
	if err != nil {
		t.Errorf("ReadFile(apply) error = %v", err)
		return
	}

	applyStr := string(applyContent)
	if !strings.Contains(applyStr, "description: Custom apply description") {
		t.Error("Apply should contain custom frontmatter description")
	}
}

func TestSlashCommandsInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		ext         string
		createFiles []string // Files to pre-create
		want        bool
	}{
		{
			name:        "returns false when no files exist",
			dir:         ".claude/commands/spectr",
			ext:         ".md",
			createFiles: []string{},
			want:        false,
		},
		{
			name:        "returns false when only proposal exists",
			dir:         ".claude/commands/spectr",
			ext:         ".md",
			createFiles: []string{".claude/commands/spectr/proposal.md"},
			want:        false,
		},
		{
			name:        "returns false when only apply exists",
			dir:         ".claude/commands/spectr",
			ext:         ".md",
			createFiles: []string{".claude/commands/spectr/apply.md"},
			want:        false,
		},
		{
			name: "returns true when both files exist",
			dir:  ".claude/commands/spectr",
			ext:  ".md",
			createFiles: []string{
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
			},
			want: true,
		},
		{
			name: "returns true for TOML format when both files exist",
			dir:  ".gemini/commands/spectr",
			ext:  ".toml",
			createFiles: []string{
				".gemini/commands/spectr/proposal.toml",
				".gemini/commands/spectr/apply.toml",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &providers.Config{SpectrDir: "spectr"}

			// Pre-create files
			for _, f := range tt.createFiles {
				// Create parent directory
				dir := f[:strings.LastIndex(f, "/")]
				if err := fs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
				if err := afero.WriteFile(fs, f, []byte("content"), 0644); err != nil {
					t.Fatalf("Failed to pre-create file %s: %v", f, err)
				}
			}

			format := FormatMarkdown
			if tt.ext == ".toml" {
				format = FormatTOML
			}
			init := NewSlashCommandsInitializer(tt.dir, tt.ext, format)

			got := init.IsSetup(fs, cfg)
			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_Path(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		want string
	}{
		{
			name: "returns directory path",
			dir:  ".claude/commands/spectr",
			want: ".claude/commands/spectr",
		},
		{
			name: "returns gemini path",
			dir:  ".gemini/commands/spectr",
			want: ".gemini/commands/spectr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewSlashCommandsInitializer(tt.dir, ".md", FormatMarkdown)

			got := init.Path()
			if got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_IsGlobal(t *testing.T) {
	tests := []struct {
		name   string
		global bool
		want   bool
	}{
		{
			name:   "project initializer returns false",
			global: false,
			want:   false,
		},
		{
			name:   "global initializer returns true",
			global: true,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var init *SlashCommandsInitializer
			if tt.global {
				init = NewGlobalSlashCommandsInitializer(".config/gemini/commands/spectr", ".toml", FormatTOML)
			} else {
				init = NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
			}

			got := init.IsGlobal()
			if got != tt.want {
				t.Errorf("IsGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_Dir(t *testing.T) {
	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)

	got := init.Dir()
	want := ".claude/commands/spectr"
	if got != want {
		t.Errorf("Dir() = %v, want %v", got, want)
	}
}

func TestSlashCommandsInitializer_Ext(t *testing.T) {
	tests := []struct {
		name string
		ext  string
		want string
	}{
		{
			name: "returns markdown extension",
			ext:  ".md",
			want: ".md",
		},
		{
			name: "returns toml extension",
			ext:  ".toml",
			want: ".toml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewSlashCommandsInitializer(".commands", tt.ext, FormatMarkdown)

			got := init.Ext()
			if got != tt.want {
				t.Errorf("Ext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_Format(t *testing.T) {
	tests := []struct {
		name   string
		format CommandFormat
		want   CommandFormat
	}{
		{
			name:   "returns markdown format",
			format: FormatMarkdown,
			want:   FormatMarkdown,
		},
		{
			name:   "returns toml format",
			format: FormatTOML,
			want:   FormatTOML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewSlashCommandsInitializer(".commands", ".md", tt.format)

			got := init.Format()
			if got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_Init_UpdatesExistingMarkdown(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	// Create directory
	if err := fs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create existing file with markers
	existingContent := `---
description: User's custom frontmatter
---

<!-- spectr:START -->

Old proposal content to be replaced

<!-- spectr:END -->

User notes below that should be preserved
`
	if err := afero.WriteFile(fs, ".claude/commands/spectr/proposal.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing file: %v", err)
	}
	// Also create apply file
	if err := afero.WriteFile(fs, ".claude/commands/spectr/apply.md", []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write apply file: %v", err)
	}

	tm.slashCommands["proposal"] = "New updated proposal content"

	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Read updated content
	content, err := afero.ReadFile(fs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	contentStr := string(content)

	// Verify old content is replaced
	if strings.Contains(contentStr, "Old proposal content to be replaced") {
		t.Error("Old content should be replaced")
	}

	// Verify new content is present
	if !strings.Contains(contentStr, "New updated proposal content") {
		t.Error("New content should be present")
	}

	// Verify user notes are preserved (content after end marker)
	if !strings.Contains(contentStr, "User notes below that should be preserved") {
		t.Error("User notes should be preserved")
	}
}

func TestSlashCommandsInitializer_Init_CreatesDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	init := NewSlashCommandsInitializer(".new/nested/commands/spectr", ".md", FormatMarkdown)

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify directory was created
	dirs := []string{".new", ".new/nested", ".new/nested/commands", ".new/nested/commands/spectr"}
	for _, dir := range dirs {
		exists, err := afero.DirExists(fs, dir)
		if err != nil {
			t.Errorf("DirExists(%s) error = %v", dir, err)
			continue
		}
		if !exists {
			t.Errorf("Directory %s was not created", dir)
		}
	}
}

func TestSlashCommandsInitializer_GlobalWithFrontmatter(t *testing.T) {
	frontmatter := map[string]string{
		"proposal": "---\ndescription: Global proposal\n---",
	}

	init := NewGlobalSlashCommandsInitializerWithFrontmatter(
		".config/tool/commands/spectr",
		".md",
		FormatMarkdown,
		frontmatter,
	)

	if !init.IsGlobal() {
		t.Error("Global initializer with frontmatter should return IsGlobal() = true")
	}

	if init.Dir() != ".config/tool/commands/spectr" {
		t.Errorf("Dir() = %v, want .config/tool/commands/spectr", init.Dir())
	}
}

func TestGenerateTOMLContent(t *testing.T) {
	description := "Test description"
	prompt := "Test prompt content"

	content := generateTOMLContent(description, prompt)

	if !strings.Contains(content, "description = \"Test description\"") {
		t.Error("TOML should contain description")
	}
	if !strings.Contains(content, "prompt = \"\"\"") {
		t.Error("TOML should contain multiline prompt")
	}
	if !strings.Contains(content, "Test prompt content") {
		t.Error("TOML should contain prompt content")
	}
}

func TestGenerateTOMLContent_EscapesQuotes(t *testing.T) {
	description := "Test"
	prompt := `Content with "quotes" and backslash \path`

	content := generateTOMLContent(description, prompt)

	if !strings.Contains(content, `\"quotes\"`) {
		t.Error("TOML should escape double quotes")
	}
	if !strings.Contains(content, `\\path`) {
		t.Error("TOML should escape backslashes")
	}
}
