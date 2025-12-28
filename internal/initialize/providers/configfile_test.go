package providers

import (
	"context"
	"html/template"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

func TestConfigFileInitializer_Init_NewFile(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create template
	tmpl := createTestTemplate(t, "test content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.CreatedFiles) != 1 || result.CreatedFiles[0] != "CLAUDE.md" {
		t.Errorf("CreatedFiles = %v, want [CLAUDE.md]", result.CreatedFiles)
	}
	if len(result.UpdatedFiles) != 0 {
		t.Errorf("UpdatedFiles = %v, want []", result.UpdatedFiles)
	}

	// Verify file content
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "<!-- spectr:start -->\ntest content\n<!-- spectr:end -->"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestConfigFileInitializer_Init_UpdateBetweenMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with markers
	existing := "Header content\n<!-- spectr:start -->\nold content\n<!-- spectr:end -->\nFooter content"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.UpdatedFiles) != 1 || result.UpdatedFiles[0] != "CLAUDE.md" {
		t.Errorf("UpdatedFiles = %v, want [CLAUDE.md]", result.UpdatedFiles)
	}
	if len(result.CreatedFiles) != 0 {
		t.Errorf("CreatedFiles = %v, want []", result.CreatedFiles)
	}

	// Verify file content
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "Header content\n<!-- spectr:start -->\nnew content\n<!-- spectr:end -->\nFooter content"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestConfigFileInitializer_Init_CaseInsensitiveMarkers(t *testing.T) {
	tests := []struct {
		name     string
		existing string
		want     string
	}{
		{
			name:     "uppercase markers",
			existing: "Header\n<!-- SPECTR:START -->\nold\n<!-- SPECTR:END -->\nFooter",
			want:     "Header\n<!-- spectr:start -->\nnew\n<!-- spectr:end -->\nFooter",
		},
		{
			name:     "mixed case markers",
			existing: "Header\n<!-- Spectr:Start -->\nold\n<!-- Spectr:End -->\nFooter",
			want:     "Header\n<!-- spectr:start -->\nnew\n<!-- spectr:end -->\nFooter",
		},
		{
			name:     "lowercase markers",
			existing: "Header\n<!-- spectr:start -->\nold\n<!-- spectr:end -->\nFooter",
			want:     "Header\n<!-- spectr:start -->\nnew\n<!-- spectr:end -->\nFooter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing file
			_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(tt.existing), 0644)

			// Create template
			tmpl := createTestTemplate(t, "new")

			// Create initializer
			init := NewConfigFileInitializer("CLAUDE.md", tmpl)

			// Execute
			_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
			if err != nil {
				t.Fatalf("Init() failed: %v", err)
			}

			// Verify file content
			content, err := afero.ReadFile(projectFs, "CLAUDE.md")
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if string(content) != tt.want {
				t.Errorf("file content = %q, want %q", string(content), tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_Init_OrphanedStartWithTrailingEnd(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with orphaned start marker and trailing end marker
	existing := "Header\n<!-- spectr:start -->\nold content\nmore content\n<!-- spectr:end -->"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.UpdatedFiles) != 1 {
		t.Errorf("UpdatedFiles = %v, want 1 file", result.UpdatedFiles)
	}

	// Verify file content - should use the trailing end marker
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "Header\n<!-- spectr:start -->\nnew content\n<!-- spectr:end -->"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestConfigFileInitializer_Init_OrphanedStartWithNoEnd(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with orphaned start marker (no end marker)
	existing := "Header\n<!-- spectr:start -->\nold content\nmore content"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.UpdatedFiles) != 1 {
		t.Errorf("UpdatedFiles = %v, want 1 file", result.UpdatedFiles)
	}

	// Verify file content - should replace everything from start marker onward
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "Header\n<!-- spectr:start -->\nnew content\n<!-- spectr:end -->"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestConfigFileInitializer_Init_NoMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file without markers
	existing := "Existing content without markers"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.UpdatedFiles) != 1 {
		t.Errorf("UpdatedFiles = %v, want 1 file", result.UpdatedFiles)
	}

	// Verify file content - should append block at end
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "Existing content without markers\n\n<!-- spectr:start -->\nnew content\n<!-- spectr:end -->"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestConfigFileInitializer_Init_TemplateContextUsage(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "myspectr"}

	// Create template that uses context variables
	tmplText := "Base: {{.BaseDir}}, Specs: {{.SpecsDir}}, Changes: {{.ChangesDir}}"
	tmpl := createTestTemplate(t, tmplText)

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify file content - template should be rendered with correct context
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "<!-- spectr:start -->\nBase: myspectr, Specs: myspectr/specs, Changes: myspectr/changes\n<!-- spectr:end -->"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestConfigFileInitializer_Init_ErrorOrphanedEndMarker(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with orphaned end marker (no start)
	existing := "Header\n<!-- spectr:end -->\nFooter"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err == nil {
		t.Fatal("Init() should fail with orphaned end marker")
	}

	if !strings.Contains(err.Error(), "orphaned end marker") {
		t.Errorf("error message = %q, want to contain 'orphaned end marker'", err.Error())
	}
}

func TestConfigFileInitializer_Init_ErrorNestedStartMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with nested start markers
	existing := "Header\n<!-- spectr:start -->\nContent\n<!-- spectr:start -->\nNested\n<!-- spectr:end -->"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err == nil {
		t.Fatal("Init() should fail with nested start markers")
	}

	if !strings.Contains(err.Error(), "nested start marker") {
		t.Errorf("error message = %q, want to contain 'nested start marker'", err.Error())
	}
}

func TestConfigFileInitializer_Init_ErrorMultipleStartMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with multiple start markers without end
	existing := "Header\n<!-- spectr:start -->\nContent\n<!-- spectr:start -->\nMore content"
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existing), 0644)

	// Create template
	tmpl := createTestTemplate(t, "new content")

	// Create initializer
	init := NewConfigFileInitializer("CLAUDE.md", tmpl)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err == nil {
		t.Fatal("Init() should fail with multiple start markers")
	}

	if !strings.Contains(err.Error(), "multiple start markers") {
		t.Errorf("error message = %q, want to contain 'multiple start markers'", err.Error())
	}
}

func TestConfigFileInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name       string
		fileExists bool
		want       bool
	}{
		{
			name:       "returns true when file exists",
			fileExists: true,
			want:       true,
		},
		{
			name:       "returns false when file does not exist",
			fileExists: false,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			if tt.fileExists {
				_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte("content"), 0644)
			}

			// Create template
			tmpl := createTestTemplate(t, "content")

			// Create initializer
			init := NewConfigFileInitializer("CLAUDE.md", tmpl)

			// Execute
			got := init.IsSetup(projectFs, homeFs, cfg)

			// Check result
			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_dedupeKey(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "simple path",
			path: "CLAUDE.md",
			want: "ConfigFileInitializer:CLAUDE.md",
		},
		{
			name: "path with directory",
			path: "docs/CLAUDE.md",
			want: "ConfigFileInitializer:docs/CLAUDE.md",
		},
		{
			name: "path with trailing slash",
			path: "CLAUDE.md/",
			want: "ConfigFileInitializer:CLAUDE.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createTestTemplate(t, "content")
			init := &ConfigFileInitializer{path: tt.path, template: tmpl}
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf("dedupeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create a test template
func createTestTemplate(t *testing.T, content string) domain.TemplateRef {
	t.Helper()
	tmpl, err := template.New("test").Parse(content)
	if err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	return domain.TemplateRef{
		Name:     "test",
		Template: tmpl,
	}
}
