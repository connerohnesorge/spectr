package providers

import (
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/spf13/afero"
)

// Mock template manager that returns different content for each command
type mockSlashTemplateManager struct{}

func (*mockSlashTemplateManager) RenderAgents(
	_ TemplateContext,
) (string, error) {
	return "", nil
}

func (*mockSlashTemplateManager) RenderInstructionPointer(
	_ TemplateContext,
) (string, error) {
	return "", nil
}

func (*mockSlashTemplateManager) RenderSlashCommand(
	commandType string,
	_ TemplateContext,
) (string, error) {
	return "Body content for " + commandType + " command", nil
}

func (*mockSlashTemplateManager) InstructionPointer() any {
	tmpl := template.Must(
		template.New("instruction-pointer").
			Parse("instruction pointer"),
	)

	return templates.NewTemplateRef(
		"instruction-pointer",
		tmpl,
	)
}

func (*mockSlashTemplateManager) Agents() any {
	tmpl := template.Must(
		template.New("agents").Parse("agents"),
	)

	return templates.NewTemplateRef(
		"agents",
		tmpl,
	)
}

func (*mockSlashTemplateManager) Project() any {
	tmpl := template.Must(
		template.New("project").Parse("project"),
	)

	return templates.NewTemplateRef(
		"project",
		tmpl,
	)
}

func (*mockSlashTemplateManager) CIWorkflow() any {
	tmpl := template.Must(
		template.New("ci-workflow").
			Parse("ci workflow"),
	)

	return templates.NewTemplateRef(
		"ci-workflow",
		tmpl,
	)
}

func (*mockSlashTemplateManager) SlashCommand(
	cmd any,
) any {
	slashCmd, ok := cmd.(templates.SlashCommand)
	if !ok {
		panic("cmd is not templates.SlashCommand")
	}
	content := "Body content for " + slashCmd.String() + " command"
	tmpl := template.Must(
		template.New(slashCmd.String()).
			Parse(content),
	)

	return templates.NewTemplateRef(
		slashCmd.String(),
		tmpl,
	)
}

func TestSlashCommandsInitializer_Init_CreateMarkdown(
	t *testing.T,
) {
	// Test creating new slash commands in Markdown format
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	commands := []templates.SlashCommand{
		templates.SlashProposal,
		templates.SlashApply,
	}
	init := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		commands,
	)

	result, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should create 2 files
	if len(result.CreatedFiles) != 2 {
		t.Errorf(
			"Init() CreatedFiles count = %d, want 2",
			len(result.CreatedFiles),
		)
	}

	// Check proposal.md
	proposalPath := ".claude/commands/spectr/proposal.md"
	exists, _ := afero.Exists(fs, proposalPath)
	if !exists {
		t.Errorf(
			"File %s should exist",
			proposalPath,
		)
	}

	content, _ := afero.ReadFile(fs, proposalPath)
	contentStr := string(content)

	// Should have frontmatter
	if !strings.HasPrefix(
		strings.TrimSpace(contentStr),
		"---",
	) {
		t.Error(
			"Markdown file should start with frontmatter",
		)
	}
	if !strings.Contains(
		contentStr,
		"description:",
	) {
		t.Error(
			"Frontmatter should contain description",
		)
	}

	// Should have markers and body
	if !strings.Contains(
		contentStr,
		spectrStartMarker,
	) {
		t.Error(
			"File should contain start marker",
		)
	}
	if !strings.Contains(
		contentStr,
		spectrEndMarker,
	) {
		t.Error("File should contain end marker")
	}
	if !strings.Contains(
		contentStr,
		"Body content for proposal command",
	) {
		t.Error(
			"File should contain rendered body",
		)
	}

	// Check apply.md
	applyPath := ".claude/commands/spectr/apply.md"
	exists, _ = afero.Exists(fs, applyPath)
	if !exists {
		t.Errorf(
			"File %s should exist",
			applyPath,
		)
	}

	content, _ = afero.ReadFile(fs, applyPath)
	contentStr = string(content)
	if !strings.Contains(
		contentStr,
		"Body content for apply command",
	) {
		t.Error(
			"Apply file should contain rendered body",
		)
	}
}

func TestSlashCommandsInitializer_Init_CreateTOML(
	t *testing.T,
) {
	// Test creating slash commands in TOML format (no frontmatter)
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	commands := []templates.SlashCommand{
		templates.SlashProposal,
		templates.SlashApply,
	}
	init := NewSlashCommandsInitializer(
		".gemini/commands/spectr",
		".toml",
		commands,
	)

	result, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should create 2 files
	if len(result.CreatedFiles) != 2 {
		t.Errorf(
			"Init() CreatedFiles count = %d, want 2",
			len(result.CreatedFiles),
		)
	}

	// Check proposal.toml
	proposalPath := ".gemini/commands/spectr/proposal.toml"
	exists, _ := afero.Exists(fs, proposalPath)
	if !exists {
		t.Errorf(
			"File %s should exist",
			proposalPath,
		)
	}

	content, _ := afero.ReadFile(fs, proposalPath)
	contentStr := string(content)

	// Should NOT have YAML frontmatter for TOML files
	if strings.HasPrefix(
		strings.TrimSpace(contentStr),
		"---",
	) {
		t.Error(
			"TOML file should not have YAML frontmatter",
		)
	}

	// Should have markers and body
	if !strings.Contains(
		contentStr,
		spectrStartMarker,
	) {
		t.Error(
			"File should contain start marker",
		)
	}
	if !strings.Contains(
		contentStr,
		"Body content for proposal command",
	) {
		t.Error(
			"File should contain rendered body",
		)
	}
}

func TestSlashCommandsInitializer_Init_UpdateExisting(
	t *testing.T,
) {
	// Test updating existing slash command files
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	// Create existing file with markers
	_ = fs.MkdirAll(
		".claude/commands/spectr",
		0o755,
	)
	existingContent := `---
description: Custom description
---

<!-- spectr:START -->

Old body content

<!-- spectr:END -->

Custom footer content
`
	_ = afero.WriteFile(
		fs,
		".claude/commands/spectr/proposal.md",
		[]byte(existingContent),
		0o644,
	)

	commands := []templates.SlashCommand{
		templates.SlashProposal,
	}
	init := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		commands,
	)

	result, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should update 1 file
	if len(result.UpdatedFiles) != 1 {
		t.Errorf(
			"Init() UpdatedFiles count = %d, want 1",
			len(result.UpdatedFiles),
		)
	}

	// Check updated content
	content, _ := afero.ReadFile(
		fs,
		".claude/commands/spectr/proposal.md",
	)
	contentStr := string(content)

	// Should preserve custom frontmatter
	if !strings.Contains(
		contentStr,
		"Custom description",
	) {
		t.Error(
			"File should preserve custom frontmatter",
		)
	}

	// Should preserve footer
	if !strings.Contains(
		contentStr,
		"Custom footer content",
	) {
		t.Error(
			"File should preserve custom footer",
		)
	}

	// Should have new body
	if !strings.Contains(
		contentStr,
		"Body content for proposal command",
	) {
		t.Error(
			"File should have new body content",
		)
	}

	// Should not have old body
	if strings.Contains(
		contentStr,
		"Old body content",
	) {
		t.Error(
			"File should not have old body content",
		)
	}
}

func TestSlashCommandsInitializer_Init_CreatesDirectory(
	t *testing.T,
) {
	// Test that parent directory is created
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	commands := []templates.SlashCommand{
		templates.SlashProposal,
	}
	init := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		commands,
	)

	_, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Check directory exists
	dirExists, _ := afero.DirExists(
		fs,
		".claude/commands/spectr",
	)
	if !dirExists {
		t.Error(
			"Directory .claude/commands/spectr should exist",
		)
	}
}

func TestSlashCommandsInitializer_WithFrontmatter(
	t *testing.T,
) {
	// Test custom frontmatter
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	customFrontmatter := "---\ncustom: value\ndescription: Custom proposal\n---"

	commands := []templates.SlashCommand{
		templates.SlashProposal,
	}
	init := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		commands,
	)
	init.WithFrontmatter(
		templates.SlashProposal,
		customFrontmatter,
	)

	_, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Check file has custom frontmatter
	content, _ := afero.ReadFile(
		fs,
		".claude/commands/spectr/proposal.md",
	)
	contentStr := string(content)

	if !strings.Contains(
		contentStr,
		"custom: value",
	) {
		t.Error(
			"File should contain custom frontmatter",
		)
	}
	if !strings.Contains(
		contentStr,
		"Custom proposal",
	) {
		t.Error(
			"File should contain custom description",
		)
	}
}

func TestSlashCommandsInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name     string
		commands []templates.SlashCommand
		setupFs  func(afero.Fs)
		want     bool
	}{
		{
			name: "all files exist",
			commands: []templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll(
					".claude/commands/spectr",
					0o755,
				)
				_ = afero.WriteFile(
					fs,
					".claude/commands/spectr/proposal.md",
					[]byte("content"),
					0o644,
				)
				_ = afero.WriteFile(
					fs,
					".claude/commands/spectr/apply.md",
					[]byte("content"),
					0o644,
				)
			},
			want: true,
		},
		{
			name: "some files missing",
			commands: []templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll(
					".claude/commands/spectr",
					0o755,
				)
				_ = afero.WriteFile(
					fs,
					".claude/commands/spectr/proposal.md",
					[]byte("content"),
					0o644,
				)
				// apply.md is missing
			},
			want: false,
		},
		{
			name: "no files exist",
			commands: []templates.SlashCommand{
				templates.SlashProposal,
			},
			setupFs: func(_ afero.Fs) {},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if tt.setupFs != nil {
				tt.setupFs(fs)
			}

			init := NewSlashCommandsInitializer(
				".claude/commands/spectr",
				".md",
				tt.commands,
			)
			cfg := NewDefaultConfig()
			got := init.IsSetup(fs, cfg)

			if got != tt.want {
				t.Errorf(
					"IsSetup() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestSlashCommandsInitializer_Path(
	t *testing.T,
) {
	dir := ".claude/commands/spectr"
	ext := ".md"
	commands := []templates.SlashCommand{
		templates.SlashProposal,
	}
	init := NewSlashCommandsInitializer(
		dir,
		ext,
		commands,
	)

	want := dir + "/*" + ext // Path() now returns dir/*ext to distinguish from DirectoryInitializer
	if got := init.Path(); got != want {
		t.Errorf(
			"Path() = %v, want %v",
			got,
			want,
		)
	}
}

func TestSlashCommandsInitializer_IsGlobal(
	t *testing.T,
) {
	tests := []struct {
		name string
		init *SlashCommandsInitializer
		want bool
	}{
		{
			name: "project-relative commands",
			init: NewSlashCommandsInitializer(
				".claude/commands/spectr",
				".md",
				[]templates.SlashCommand{
					templates.SlashProposal,
				},
			),
			want: false,
		},
		{
			name: "global commands",
			init: NewGlobalSlashCommandsInitializer(
				".config/aider/commands",
				".md",
				[]templates.SlashCommand{
					templates.SlashProposal,
				},
			),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.init.IsGlobal(); got != tt.want {
				t.Errorf(
					"IsGlobal() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestSlashCommandsInitializer_Idempotent(
	t *testing.T,
) {
	// Test that running Init multiple times is safe
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	commands := []templates.SlashCommand{
		templates.SlashProposal,
	}
	init := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		commands,
	)

	// First run - creates files
	result1, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("First Init() error = %v", err)
	}
	if len(result1.CreatedFiles) != 1 {
		t.Error(
			"First Init() should create 1 file",
		)
	}

	// Second run - updates files
	result2, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Second Init() error = %v", err)
	}
	if len(result2.UpdatedFiles) != 1 {
		t.Error(
			"Second Init() should update 1 file",
		)
	}

	// Third run - still safe
	result3, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Third Init() error = %v", err)
	}
	if len(result3.UpdatedFiles) != 1 {
		t.Error(
			"Third Init() should update 1 file",
		)
	}

	// Files should still be valid
	content, _ := afero.ReadFile(
		fs,
		".claude/commands/spectr/proposal.md",
	)
	contentStr := string(content)
	if !strings.Contains(
		contentStr,
		spectrStartMarker,
	) {
		t.Error("File should still have markers")
	}
}

func TestSlashCommandsInitializer_AddsFrontmatterToExistingFile(
	t *testing.T,
) {
	// Test that frontmatter is added to existing file without frontmatter
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockSlashTemplateManager{}

	// Create existing file WITHOUT frontmatter
	_ = fs.MkdirAll(
		".claude/commands/spectr",
		0o755,
	)
	existingContent := `<!-- spectr:START -->

Old body content

<!-- spectr:END -->
`
	_ = afero.WriteFile(
		fs,
		".claude/commands/spectr/proposal.md",
		[]byte(existingContent),
		0o644,
	)

	commands := []templates.SlashCommand{
		templates.SlashProposal,
	}
	init := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		commands,
	)

	_, err := init.Init(
		context.Background(),
		fs,
		cfg,
		tm,
	)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Check file now has frontmatter
	content, _ := afero.ReadFile(
		fs,
		".claude/commands/spectr/proposal.md",
	)
	contentStr := string(content)

	if !strings.HasPrefix(
		strings.TrimSpace(contentStr),
		"---",
	) {
		t.Error(
			"File should now have frontmatter",
		)
	}
	if !strings.Contains(
		contentStr,
		"description:",
	) {
		t.Error(
			"Frontmatter should contain description",
		)
	}
}
