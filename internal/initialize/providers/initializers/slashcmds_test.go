package initializers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// mockRenderer implements providers.TemplateRenderer for testing
type mockRenderer struct {
	content map[string]string // command name -> rendered content
	err     error             // error to return, if any
}

func newMockRenderer() *mockRenderer {
	return &mockRenderer{
		content: map[string]string{
			"proposal": "This is the proposal command content.",
			"apply":    "This is the apply command content.",
		},
	}
}

func (m *mockRenderer) RenderAgents(
	_ providers.TemplateContext,
) (string, error) {
	return "# AGENTS content", m.err
}

func (m *mockRenderer) RenderInstructionPointer(
	_ providers.TemplateContext,
) (string, error) {
	return "# Instruction pointer", m.err
}

func (m *mockRenderer) RenderSlashCommand(
	command string,
	_ providers.TemplateContext,
) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if content, ok := m.content[command]; ok {
		return content, nil
	}

	return "", fmt.Errorf(
		"unknown command: %s",
		command,
	)
}

func TestNewSlashCommandsInitializer(
	t *testing.T,
) {
	tests := []struct {
		name   string
		dir    string
		ext    string
		format providers.CommandFormat
	}{
		{
			name:   "markdown format",
			dir:    ".claude/commands/spectr",
			ext:    ".md",
			format: providers.FormatMarkdown,
		},
		{
			name:   "toml format",
			dir:    ".gemini/commands/spectr",
			ext:    ".toml",
			format: providers.FormatTOML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := newMockRenderer()
			s := NewSlashCommandsInitializer(
				tt.dir,
				tt.ext,
				tt.format,
				renderer,
			)

			if s == nil {
				t.Fatal(
					"NewSlashCommandsInitializer() returned nil",
				)
			}

			if s.Dir != tt.dir {
				t.Errorf(
					"Dir = %s, want %s",
					s.Dir,
					tt.dir,
				)
			}

			if s.Extension != tt.ext {
				t.Errorf(
					"Extension = %s, want %s",
					s.Extension,
					tt.ext,
				)
			}

			if s.Format != tt.format {
				t.Errorf(
					"Format = %d, want %d",
					s.Format,
					tt.format,
				)
			}

			if s.Renderer != renderer {
				t.Error(
					"Renderer was not set correctly",
				)
			}
		})
	}
}

func TestSlashCommandsInitializer_Init_MarkdownFormat(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	err := s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify both proposal and apply files were created
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		filePath := ".claude/commands/spectr/" + cmd + ".md"
		exists, err := afero.Exists(fs, filePath)
		if err != nil {
			t.Errorf(
				"Error checking file %s: %v",
				filePath,
				err,
			)

			continue
		}

		if !exists {
			t.Errorf(
				"File %s was not created",
				filePath,
			)

			continue
		}

		// Read and verify content
		content, err := afero.ReadFile(
			fs,
			filePath,
		)
		if err != nil {
			t.Errorf(
				"Error reading file %s: %v",
				filePath,
				err,
			)

			continue
		}

		contentStr := string(content)

		// Verify frontmatter is present
		if !strings.Contains(contentStr, "---") {
			t.Errorf(
				"File %s should contain YAML frontmatter",
				filePath,
			)
		}

		// Verify markers are present
		if !strings.Contains(
			contentStr,
			spectrStartMarker,
		) {
			t.Errorf(
				"File %s should contain start marker",
				filePath,
			)
		}

		if !strings.Contains(
			contentStr,
			spectrEndMarker,
		) {
			t.Errorf(
				"File %s should contain end marker",
				filePath,
			)
		}

		// Verify command content is present
		expectedContent := renderer.content[cmd]
		if !strings.Contains(
			contentStr,
			expectedContent,
		) {
			t.Errorf(
				"File %s should contain rendered content",
				filePath,
			)
		}
	}
}

func TestSlashCommandsInitializer_Init_TOMLFormat(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".gemini/commands/spectr",
		".toml",
		providers.FormatTOML,
		renderer,
	)

	err := s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify both proposal and apply files were created in TOML format
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		filePath := ".gemini/commands/spectr/" + cmd + ".toml"
		exists, err := afero.Exists(fs, filePath)
		if err != nil {
			t.Errorf(
				"Error checking file %s: %v",
				filePath,
				err,
			)

			continue
		}

		if !exists {
			t.Errorf(
				"File %s was not created",
				filePath,
			)

			continue
		}

		// Read and verify content
		content, err := afero.ReadFile(
			fs,
			filePath,
		)
		if err != nil {
			t.Errorf(
				"Error reading file %s: %v",
				filePath,
				err,
			)

			continue
		}

		contentStr := string(content)

		// Verify TOML format
		if !strings.Contains(
			contentStr,
			"description =",
		) {
			t.Errorf(
				"File %s should contain description field",
				filePath,
			)
		}

		if !strings.Contains(
			contentStr,
			"prompt =",
		) {
			t.Errorf(
				"File %s should contain prompt field",
				filePath,
			)
		}

		// Verify command content is present in prompt
		expectedContent := renderer.content[cmd]
		if !strings.Contains(
			contentStr,
			expectedContent,
		) {
			t.Errorf(
				"File %s should contain rendered content in prompt",
				filePath,
			)
		}
	}
}

func TestSlashCommandsInitializer_Init_CreatesDirectory(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr/nested",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	err := s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify directory was created
	info, err := fs.Stat(
		".claude/commands/spectr/nested",
	)
	if err != nil {
		t.Fatalf(
			"Directory was not created: %v",
			err,
		)
	}

	if !info.IsDir() {
		t.Error("Path should be a directory")
	}
}

func TestSlashCommandsInitializer_Init_Idempotent(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// Call Init multiple times
	for i := range 3 {
		err := s.Init(ctx, fs, cfg)
		if err != nil {
			t.Fatalf(
				"Init() call %d failed: %v",
				i+1,
				err,
			)
		}
	}

	// Verify files still exist and are correct
	proposalPath := ".claude/commands/spectr/proposal.md"
	content, err := afero.ReadFile(
		fs,
		proposalPath,
	)
	if err != nil {
		t.Fatalf(
			"Failed to read file after multiple Init calls: %v",
			err,
		)
	}

	contentStr := string(content)

	// Should have exactly one set of markers
	startCount := strings.Count(
		contentStr,
		spectrStartMarker,
	)
	endCount := strings.Count(
		contentStr,
		spectrEndMarker,
	)

	if startCount != 1 {
		t.Errorf(
			"Start marker count = %d, want 1",
			startCount,
		)
	}

	if endCount != 1 {
		t.Errorf(
			"End marker count = %d, want 1",
			endCount,
		)
	}
}

func TestSlashCommandsInitializer_Init_UpdatesExistingFile(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// First Init
	err := s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("First Init() failed: %v", err)
	}

	// Update renderer content
	renderer.content["proposal"] = "Updated proposal content"

	// Second Init should update
	err = s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Second Init() failed: %v", err)
	}

	// Verify content was updated
	proposalPath := ".claude/commands/spectr/proposal.md"
	content, err := afero.ReadFile(
		fs,
		proposalPath,
	)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(
		string(content),
		"Updated proposal content",
	) {
		t.Error(
			"File should contain updated content",
		)
	}
}

func TestSlashCommandsInitializer_Init_WithConfig(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()

	// Use custom config with different spectr dir
	cfg := &providers.Config{
		SpectrDir: "custom-spectr",
	}

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	err := s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf(
			"Init() with custom config failed: %v",
			err,
		)
	}

	// Files should be created (config affects template context, not file paths)
	proposalPath := ".claude/commands/spectr/proposal.md"
	exists, err := afero.Exists(fs, proposalPath)
	if err != nil {
		t.Fatalf("Error checking file: %v", err)
	}

	if !exists {
		t.Error(
			"File should be created even with custom config",
		)
	}
}

func TestSlashCommandsInitializer_Init_RenderError(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := &mockRenderer{
		content: make(map[string]string),
		err:     errors.New("render error"),
	}
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	err := s.Init(ctx, fs, cfg)
	if err == nil {
		t.Error(
			"Init() should fail when renderer returns error",
		)
	}
}

func TestSlashCommandsInitializer_IsSetup_NotExists(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// Files don't exist
	if s.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when files don't exist",
		)
	}
}

func TestSlashCommandsInitializer_IsSetup_BothExist(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// Create files
	err := s.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Now IsSetup should return true
	if !s.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return true when both files exist",
		)
	}
}

func TestSlashCommandsInitializer_IsSetup_OnlyProposalExists(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// Create only proposal file
	err := fs.MkdirAll(
		".claude/commands/spectr",
		0755,
	)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	err = afero.WriteFile(
		fs,
		".claude/commands/spectr/proposal.md",
		[]byte("content"),
		0644,
	)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// IsSetup should return false (apply is missing)
	if s.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when only proposal exists",
		)
	}
}

func TestSlashCommandsInitializer_IsSetup_OnlyApplyExists(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// Create only apply file
	err := fs.MkdirAll(
		".claude/commands/spectr",
		0755,
	)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	err = afero.WriteFile(
		fs,
		".claude/commands/spectr/apply.md",
		[]byte("content"),
		0644,
	)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// IsSetup should return false (proposal is missing)
	if s.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when only apply exists",
		)
	}
}

func TestSlashCommandsInitializer_Key_Markdown(
	t *testing.T,
) {
	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	key := s.Key()
	expected := "slashcmds:.claude/commands/spectr:.md:0"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestSlashCommandsInitializer_Key_TOML(
	t *testing.T,
) {
	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".gemini/commands/spectr",
		".toml",
		providers.FormatTOML,
		renderer,
	)

	key := s.Key()
	expected := "slashcmds:.gemini/commands/spectr:.toml:1"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestSlashCommandsInitializer_Key_Consistent(
	t *testing.T,
) {
	renderer := newMockRenderer()
	s := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	// Key should be consistent across multiple calls
	key1 := s.Key()
	key2 := s.Key()
	key3 := s.Key()

	if key1 != key2 || key2 != key3 {
		t.Errorf(
			"Key() is not consistent: %s, %s, %s",
			key1,
			key2,
			key3,
		)
	}
}

func TestSlashCommandsInitializer_Key_DifferentConfigs(
	t *testing.T,
) {
	renderer := newMockRenderer()

	s1 := NewSlashCommandsInitializer(
		".claude/commands/spectr",
		".md",
		providers.FormatMarkdown,
		renderer,
	)

	s2 := NewSlashCommandsInitializer(
		".gemini/commands/spectr",
		".toml",
		providers.FormatTOML,
		renderer,
	)

	// Keys should differ for different configurations
	if s1.Key() == s2.Key() {
		t.Errorf(
			"Keys should differ for different configs: %s vs %s",
			s1.Key(),
			s2.Key(),
		)
	}
}

func TestSlashCommandsInitializer_ImplementsInterface(
	_ *testing.T,
) {
	// Compile-time check is in slashcmds.go, but this is a runtime verification
	var _ providers.Initializer = (*SlashCommandsInitializer)(nil)
}

func TestSlashCommands(t *testing.T) {
	commands := slashCommands()

	// Should return proposal and apply commands
	if len(commands) != 2 {
		t.Errorf(
			"slashCommands() returned %d commands, want 2",
			len(commands),
		)
	}

	// Verify proposal command
	foundProposal := false
	foundApply := false

	for _, cmd := range commands {
		switch cmd.name {
		case "proposal":
			foundProposal = true
			if cmd.description == "" {
				t.Error(
					"proposal command should have a description",
				)
			}
		case "apply":
			foundApply = true
			if cmd.description == "" {
				t.Error(
					"apply command should have a description",
				)
			}
		}
	}

	if !foundProposal {
		t.Error(
			"slashCommands() should include proposal command",
		)
	}

	if !foundApply {
		t.Error(
			"slashCommands() should include apply command",
		)
	}
}

func TestGetFrontmatter(t *testing.T) {
	tests := []struct {
		command  string
		wantDesc bool
		wantYAML bool
	}{
		{
			command:  "proposal",
			wantDesc: true,
			wantYAML: true,
		},
		{
			command:  "apply",
			wantDesc: true,
			wantYAML: true,
		},
		{
			command:  "unknown",
			wantDesc: false,
			wantYAML: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			frontmatter := getFrontmatter(
				tt.command,
			)

			if tt.wantYAML {
				if !strings.HasPrefix(
					frontmatter,
					"---",
				) {
					t.Error(
						"Frontmatter should start with ---",
					)
				}
				if !strings.HasSuffix(
					frontmatter,
					"---",
				) {
					t.Error(
						"Frontmatter should end with ---",
					)
				}
			} else if frontmatter != "" {
				t.Errorf("Unknown command should return empty string, got %s", frontmatter)
			}

			if !tt.wantDesc {
				return
			}

			if !strings.Contains(
				frontmatter,
				"description:",
			) {
				t.Error(
					"Frontmatter should contain description",
				)
			}
		})
	}
}

func TestGenerateTOMLContent(t *testing.T) {
	s := &SlashCommandsInitializer{}

	description := "Test description"
	prompt := "Test prompt content"

	content := s.generateTOMLContent(
		description,
		prompt,
	)

	// Verify TOML structure
	if !strings.Contains(
		content,
		"description =",
	) {
		t.Error(
			"Content should contain description field",
		)
	}

	if !strings.Contains(
		content,
		`"`+description+`"`,
	) {
		t.Error(
			"Content should contain the description value",
		)
	}

	if !strings.Contains(content, "prompt =") {
		t.Error(
			"Content should contain prompt field",
		)
	}

	if !strings.Contains(content, prompt) {
		t.Error(
			"Content should contain the prompt value",
		)
	}

	// Verify multiline string format
	if !strings.Contains(content, `"""`) {
		t.Error(
			"Prompt should use multiline string format",
		)
	}
}

func TestCreateMarkdownContent(t *testing.T) {
	frontmatter := "---\ndescription: Test\n---"
	body := "Command body content"

	content := createMarkdownContent(
		frontmatter,
		body,
	)

	// Verify frontmatter is at the start
	if !strings.HasPrefix(content, "---") {
		t.Error(
			"Content should start with frontmatter",
		)
	}

	// Verify markers are present
	if !strings.Contains(
		content,
		spectrStartMarker,
	) {
		t.Error(
			"Content should contain start marker",
		)
	}

	if !strings.Contains(
		content,
		spectrEndMarker,
	) {
		t.Error(
			"Content should contain end marker",
		)
	}

	// Verify body is present
	if !strings.Contains(content, body) {
		t.Error("Content should contain body")
	}

	// Verify frontmatter comes before markers
	frontmatterIdx := strings.Index(
		content,
		"description: Test",
	)
	markerIdx := strings.Index(
		content,
		spectrStartMarker,
	)
	if frontmatterIdx > markerIdx {
		t.Error(
			"Frontmatter should come before markers",
		)
	}
}

func TestCreateMarkdownContent_EmptyFrontmatter(
	t *testing.T,
) {
	frontmatter := ""
	body := "Command body content"

	content := createMarkdownContent(
		frontmatter,
		body,
	)

	// Should start with markers (no frontmatter)
	if !strings.HasPrefix(
		content,
		spectrStartMarker,
	) {
		t.Error(
			"Content without frontmatter should start with start marker",
		)
	}

	// Verify body is present
	if !strings.Contains(content, body) {
		t.Error("Content should contain body")
	}
}
