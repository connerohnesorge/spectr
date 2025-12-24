package providers

import (
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/spf13/afero"
)

// Mock TemplateManager for testing
type mockTemplateManager struct{}

func (*mockTemplateManager) RenderAgents(
	_ TemplateContext,
) (string, error) {
	return "", nil
}

func (*mockTemplateManager) RenderInstructionPointer(
	_ TemplateContext,
) (string, error) {
	return "", nil
}

func (*mockTemplateManager) RenderSlashCommand(
	_ string,
	_ TemplateContext,
) (string, error) {
	return "", nil
}

func (*mockTemplateManager) InstructionPointer() any {
	return templates.NewTemplateRef(
		"instruction-pointer.md.tmpl",
		nil,
	)
}

func (*mockTemplateManager) Agents() any {
	return templates.NewTemplateRef(
		"AGENTS.md.tmpl",
		nil,
	)
}

func (*mockTemplateManager) Project() any {
	return templates.NewTemplateRef(
		"project.md.tmpl",
		nil,
	)
}

func (*mockTemplateManager) CIWorkflow() any {
	return templates.NewTemplateRef(
		"spectr-ci.yml.tmpl",
		nil,
	)
}

func (*mockTemplateManager) SlashCommand(
	_ any,
) any {
	return templates.NewTemplateRef(
		"slash-proposal.md.tmpl",
		nil,
	)
}

// Mock template getter that returns an actual TemplateRef with a simple template
func mockTemplateGetter(_ TemplateManager) any {
	// Create a simple template
	tmpl := template.Must(
		template.New("test").
			Parse("Test content for instruction pointer"),
	)
	// Return actual TemplateRef using the public constructor
	return templates.NewTemplateRef("test", tmpl)
}

func TestConfigFileInitializer_Init_CreateNew(
	t *testing.T,
) {
	// Test creating a new config file
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockTemplateManager{}

	init := NewConfigFileInitializer(
		"CLAUDE.md",
		mockTemplateGetter,
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

	if len(result.CreatedFiles) != 1 ||
		result.CreatedFiles[0] != "CLAUDE.md" {
		t.Errorf(
			"Init() CreatedFiles = %v, want [CLAUDE.md]",
			result.CreatedFiles,
		)
	}

	if len(result.UpdatedFiles) != 0 {
		t.Errorf(
			"Init() UpdatedFiles = %v, want []",
			result.UpdatedFiles,
		)
	}

	// Check file exists and has correct content
	exists, err := afero.Exists(fs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Fatal("CLAUDE.md should exist")
	}

	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	contentStr := string(content)
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
		"Test content for instruction pointer",
	) {
		t.Error(
			"File should contain rendered content",
		)
	}
}

func TestConfigFileInitializer_Init_UpdateExisting(
	t *testing.T,
) {
	// Test updating an existing file with markers
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockTemplateManager{}

	// Create existing file with markers and different content
	existingContent := `# My Custom Header

<!-- spectr:START -->
Old content here
<!-- spectr:END -->

# My Custom Footer
`
	_ = afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte(existingContent),
		0o644,
	)

	init := NewConfigFileInitializer(
		"CLAUDE.md",
		mockTemplateGetter,
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

	if len(result.UpdatedFiles) != 1 ||
		result.UpdatedFiles[0] != "CLAUDE.md" {
		t.Errorf(
			"Init() UpdatedFiles = %v, want [CLAUDE.md]",
			result.UpdatedFiles,
		)
	}

	if len(result.CreatedFiles) != 0 {
		t.Errorf(
			"Init() CreatedFiles = %v, want []",
			result.CreatedFiles,
		)
	}

	// Check file content was updated
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	contentStr := string(content)

	// Should preserve header and footer
	if !strings.Contains(
		contentStr,
		"# My Custom Header",
	) {
		t.Error(
			"File should preserve custom header",
		)
	}
	if !strings.Contains(
		contentStr,
		"# My Custom Footer",
	) {
		t.Error(
			"File should preserve custom footer",
		)
	}

	// Should have new content
	if !strings.Contains(
		contentStr,
		"Test content for instruction pointer",
	) {
		t.Error("File should contain new content")
	}

	// Should not have old content
	if strings.Contains(
		contentStr,
		"Old content here",
	) {
		t.Error(
			"File should not contain old content",
		)
	}
}

func TestConfigFileInitializer_Init_AppendMarkers(
	t *testing.T,
) {
	// Test appending markers to existing file without markers
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockTemplateManager{}

	// Create existing file without markers
	existingContent := `# My Existing File

Some user content here.
`
	_ = afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte(existingContent),
		0o644,
	)

	init := NewConfigFileInitializer(
		"CLAUDE.md",
		mockTemplateGetter,
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

	if len(result.UpdatedFiles) != 1 {
		t.Errorf(
			"Init() UpdatedFiles = %v, want 1 file",
			result.UpdatedFiles,
		)
	}

	// Check file content
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	contentStr := string(content)

	// Should preserve existing content
	if !strings.Contains(
		contentStr,
		"# My Existing File",
	) {
		t.Error(
			"File should preserve existing content",
		)
	}
	if !strings.Contains(
		contentStr,
		"Some user content here.",
	) {
		t.Error(
			"File should preserve existing content",
		)
	}

	// Should have markers appended
	if !strings.Contains(
		contentStr,
		spectrStartMarker,
	) {
		t.Error("File should have start marker")
	}
	if !strings.Contains(
		contentStr,
		spectrEndMarker,
	) {
		t.Error("File should have end marker")
	}
	if !strings.Contains(
		contentStr,
		"Test content for instruction pointer",
	) {
		t.Error("File should contain new content")
	}
}

func TestConfigFileInitializer_Init_CreatesDirectory(
	t *testing.T,
) {
	// Test that parent directory is created if needed
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockTemplateManager{}

	init := NewConfigFileInitializer(
		"nested/path/CLAUDE.md",
		mockTemplateGetter,
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

	if len(result.CreatedFiles) != 1 {
		t.Errorf(
			"Init() CreatedFiles = %v, want 1 file",
			result.CreatedFiles,
		)
	}

	// Check file exists
	exists, err := afero.Exists(
		fs,
		"nested/path/CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error(
			"File should exist at nested/path/CLAUDE.md",
		)
	}

	// Check directory exists
	dirExists, err := afero.DirExists(
		fs,
		"nested/path",
	)
	if err != nil {
		t.Fatalf("DirExists() error = %v", err)
	}
	if !dirExists {
		t.Error(
			"Directory nested/path should exist",
		)
	}
}

func TestConfigFileInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name    string
		path    string
		setupFs func(afero.Fs)
		want    bool
	}{
		{
			name: "file exists",
			path: "CLAUDE.md",
			setupFs: func(fs afero.Fs) {
				_ = afero.WriteFile(
					fs,
					"CLAUDE.md",
					[]byte("content"),
					0o644,
				)
			},
			want: true,
		},
		{
			name:    "file does not exist",
			path:    "CLAUDE.md",
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

			init := NewConfigFileInitializer(
				tt.path,
				mockTemplateGetter,
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

func TestConfigFileInitializer_Path(
	t *testing.T,
) {
	path := "CLAUDE.md"
	init := NewConfigFileInitializer(
		path,
		mockTemplateGetter,
	)

	if got := init.Path(); got != path {
		t.Errorf(
			"Path() = %v, want %v",
			got,
			path,
		)
	}
}

func TestConfigFileInitializer_IsGlobal(
	t *testing.T,
) {
	tests := []struct {
		name string
		init *ConfigFileInitializer
		want bool
	}{
		{
			name: "project-relative file",
			init: NewConfigFileInitializer(
				"CLAUDE.md",
				mockTemplateGetter,
			),
			want: false,
		},
		{
			name: "global file",
			init: NewGlobalConfigFileInitializer(
				".config/tool/config.md",
				mockTemplateGetter,
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

func TestConfigFileInitializer_Idempotent(
	t *testing.T,
) {
	// Test that running Init multiple times is safe
	fs := afero.NewMemMapFs()
	cfg := NewDefaultConfig()
	tm := &mockTemplateManager{}

	init := NewConfigFileInitializer(
		"CLAUDE.md",
		mockTemplateGetter,
	)

	// First run - creates file
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
		t.Error("First Init() should create file")
	}

	// Second run - updates file
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
			"Second Init() should update file",
		)
	}

	// Third run - updates file again
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
		t.Error("Third Init() should update file")
	}

	// File should still exist and be valid
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(
		contentStr,
		spectrStartMarker,
	) {
		t.Error(
			"File should still have start marker",
		)
	}
	if !strings.Contains(
		contentStr,
		spectrEndMarker,
	) {
		t.Error(
			"File should still have end marker",
		)
	}
}
