package providers

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

// mockTemplateManager implements TemplateManager for testing
type mockTemplateManager struct {
	content string
	err     error
}

func (m *mockTemplateManager) RenderAgents(_ TemplateContext) (string, error) {
	return m.content, m.err
}

func (m *mockTemplateManager) RenderInstructionPointer(
	_ TemplateContext,
) (string, error) {
	return m.content, m.err
}

func (m *mockTemplateManager) RenderSlashCommand(
	_ string,
	_ TemplateContext,
) (string, error) {
	return m.content, m.err
}

func TestConfigFileInitializer_Init_CreateNew(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	tm := &mockTemplateManager{
		content: "Test instruction content",
	}

	init := NewConfigFileInitializer("CLAUDE.md", "instruction_pointer")
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should report file as created
	if len(result.CreatedFiles) != 1 {
		t.Errorf("Init() created %d files, want 1", len(result.CreatedFiles))
	}
	if len(result.CreatedFiles) > 0 && result.CreatedFiles[0] != "CLAUDE.md" {
		t.Errorf("Init() created %s, want CLAUDE.md", result.CreatedFiles[0])
	}

	// Verify file exists
	exists, err := afero.Exists(fs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to check file existence: %v", err)
	}
	if !exists {
		t.Fatal("file does not exist")
	}

	// Verify file content
	content, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, SpectrStartMarker) {
		t.Error("file does not contain start marker")
	}
	if !strings.Contains(contentStr, SpectrEndMarker) {
		t.Error("file does not contain end marker")
	}
	if !strings.Contains(contentStr, "Test instruction content") {
		t.Error("file does not contain template content")
	}
}

func TestConfigFileInitializer_Init_UpdateExisting(t *testing.T) {
	tests := []struct {
		name            string
		existingContent string
		templateContent string
		wantUpdated     bool
		checkContent    func(t *testing.T, content string)
	}{
		{
			name: "updates content between existing markers",
			existingContent: `# Existing File

<!-- spectr:START -->
Old content here
<!-- spectr:END -->

More content`,
			templateContent: "New template content",
			wantUpdated:     true,
			checkContent: func(t *testing.T, content string) {
				if !strings.Contains(content, "New template content") {
					t.Error("content not updated")
				}
				if strings.Contains(content, "Old content here") {
					t.Error("old content still present")
				}
				if !strings.Contains(content, "# Existing File") {
					t.Error("existing header removed")
				}
				if !strings.Contains(content, "More content") {
					t.Error("existing footer removed")
				}
			},
		},
		{
			name:            "appends markers when missing",
			existingContent: "# Existing File\n\nNo markers here\n",
			templateContent: "New template content",
			wantUpdated:     true,
			checkContent: func(t *testing.T, content string) {
				if !strings.Contains(content, "New template content") {
					t.Error("template content not added")
				}
				if !strings.Contains(content, SpectrStartMarker) {
					t.Error("start marker not added")
				}
				if !strings.Contains(content, SpectrEndMarker) {
					t.Error("end marker not added")
				}
				if !strings.Contains(content, "# Existing File") {
					t.Error("existing content removed")
				}
			},
		},
		{
			name: "no update when content unchanged",
			existingContent: `<!-- spectr:START -->
Same content
<!-- spectr:END -->
`,
			templateContent: "Same content",
			wantUpdated:     false,
			checkContent: func(t *testing.T, content string) {
				if !strings.Contains(content, "Same content") {
					t.Error("content not preserved")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}
			ctx := context.Background()

			// Create existing file
			if err := afero.WriteFile(fs, "CLAUDE.md", []byte(tt.existingContent), 0644); err != nil {
				t.Fatalf("failed to create existing file: %v", err)
			}

			tm := &mockTemplateManager{
				content: tt.templateContent,
			}

			init := NewConfigFileInitializer("CLAUDE.md", "instruction_pointer")
			result, err := init.Init(ctx, fs, cfg, tm)

			if err != nil {
				t.Fatalf("Init() error = %v", err)
			}

			// Check if file was updated
			if tt.wantUpdated {
				if len(result.UpdatedFiles) != 1 {
					t.Errorf("Init() updated %d files, want 1", len(result.UpdatedFiles))
				}
			} else {
				if len(result.UpdatedFiles) != 0 {
					t.Errorf("Init() updated %d files, want 0", len(result.UpdatedFiles))
				}
			}

			// Verify content
			content, err := afero.ReadFile(fs, "CLAUDE.md")
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			tt.checkContent(t, string(content))
		})
	}
}

func TestConfigFileInitializer_Init_NestedPath(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	ctx := context.Background()

	tm := &mockTemplateManager{
		content: "Test content",
	}

	// Use nested path that requires parent directory creation
	init := NewConfigFileInitializer(".claude/config/CLAUDE.md", "instruction_pointer")
	result, err := init.Init(ctx, fs, cfg, tm)

	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should create file
	if len(result.CreatedFiles) != 1 {
		t.Errorf("Init() created %d files, want 1", len(result.CreatedFiles))
	}

	// Verify file exists
	exists, err := afero.Exists(fs, ".claude/config/CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to check file existence: %v", err)
	}
	if !exists {
		t.Fatal("file does not exist")
	}

	// Verify parent directory exists
	dirExists, err := afero.DirExists(fs, ".claude/config")
	if err != nil {
		t.Fatalf("failed to check directory existence: %v", err)
	}
	if !dirExists {
		t.Fatal("parent directory does not exist")
	}
}

func TestConfigFileInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name        string
		fileExists  bool
		fileContent string
		want        bool
	}{
		{
			name:       "file does not exist",
			fileExists: false,
			want:       false,
		},
		{
			name:        "file exists without markers",
			fileExists:  true,
			fileContent: "# Some content\nNo markers here",
			want:        false,
		},
		{
			name:        "file exists with only start marker",
			fileExists:  true,
			fileContent: "<!-- spectr:START -->\nContent",
			want:        false,
		},
		{
			name:        "file exists with only end marker",
			fileExists:  true,
			fileContent: "Content\n<!-- spectr:END -->",
			want:        false,
		},
		{
			name:        "file exists with both markers",
			fileExists:  true,
			fileContent: "<!-- spectr:START -->\nContent\n<!-- spectr:END -->",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			if tt.fileExists {
				if err := afero.WriteFile(fs, "CLAUDE.md", []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			init := NewConfigFileInitializer("CLAUDE.md", "instruction_pointer")
			if got := init.IsSetup(fs, cfg); got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_Path(t *testing.T) {
	init := NewConfigFileInitializer("CLAUDE.md", "instruction_pointer")
	if got := init.Path(); got != "CLAUDE.md" {
		t.Errorf("Path() = %v, want CLAUDE.md", got)
	}
}

func TestConfigFileInitializer_IsGlobal(t *testing.T) {
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
			init := NewConfigFileInitializer("CLAUDE.md", "instruction_pointer")
			if tt.isGlobal {
				init = init.WithGlobal(true)
			}

			if got := init.IsGlobal(); got != tt.want {
				t.Errorf("IsGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}
