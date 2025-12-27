package providers

import (
	"context"
	"strings"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

func TestConfigFileInitializer_Init_NewFile(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create a simple template
	tmpl := template.Must(template.New("test.tmpl").Parse("Test content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

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

	expectedContent := "<!-- spectr:start -->\nTest content\n<!-- spectr:end -->\n"
	if string(content) != expectedContent {
		t.Errorf("file content = %q, want %q", string(content), expectedContent)
	}
}

func TestConfigFileInitializer_Init_UpdateBetweenMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with markers
	existingContent := `Some content before
<!-- spectr:start -->
Old content
<!-- spectr:end -->
Some content after`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

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

	expectedContent := `Some content before
<!-- spectr:start -->
New content
<!-- spectr:end -->
Some content after`

	if string(content) != expectedContent {
		t.Errorf("file content = %q, want %q", string(content), expectedContent)
	}
}

func TestConfigFileInitializer_Init_CaseInsensitiveMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with uppercase markers
	existingContent := `Some content before
<!-- SPECTR:START -->
Old content
<!-- SPECTR:END -->
Some content after`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify file content - should write lowercase markers
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expectedContent := `Some content before
<!-- spectr:start -->
New content
<!-- spectr:end -->
Some content after`

	if string(content) != expectedContent {
		t.Errorf("file content = %q, want %q", string(content), expectedContent)
	}
}

func TestConfigFileInitializer_Init_OrphanedStartWithTrailingEnd(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with orphaned start and trailing end
	existingContent := `Some content before
<!-- spectr:start -->
Some middle content
More middle content
<!-- spectr:end -->
Some content after`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify - should succeed and use trailing end marker
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify file content
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expectedContent := `Some content before
<!-- spectr:start -->
New content
<!-- spectr:end -->
Some content after`

	if string(content) != expectedContent {
		t.Errorf("file content = %q, want %q", string(content), expectedContent)
	}
}

func TestConfigFileInitializer_Init_OrphanedStartNoEnd(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with orphaned start marker (no end)
	existingContent := `Some content before
<!-- spectr:start -->
Old content that should be replaced
More old content`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify - should succeed and replace from start marker onward
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify file content
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expectedContent := `Some content before
<!-- spectr:start -->
New content
<!-- spectr:end -->`

	if string(content) != expectedContent {
		t.Errorf("file content = %q, want %q", string(content), expectedContent)
	}
}

func TestConfigFileInitializer_Init_OrphanedEndMarker(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with orphaned end marker (no start)
	existingContent := `Some content before
<!-- spectr:end -->
Some content after`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify - should return error
	if err == nil {
		t.Fatal("Init() expected error for orphaned end marker, got nil")
	}

	if !strings.Contains(err.Error(), "orphaned end marker") {
		t.Errorf("Init() error = %v, want error containing 'orphaned end marker'", err)
	}
}

func TestConfigFileInitializer_Init_NestedStartMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with nested start markers
	existingContent := `Some content before
<!-- spectr:start -->
Some content
<!-- spectr:start -->
Nested content
<!-- spectr:end -->
Some content after`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify - should return error
	if err == nil {
		t.Fatal("Init() expected error for nested start markers, got nil")
	}

	if !strings.Contains(err.Error(), "nested start marker") {
		t.Errorf("Init() error = %v, want error containing 'nested start marker'", err)
	}
}

func TestConfigFileInitializer_Init_MultipleStartMarkers(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with multiple start markers (no end)
	existingContent := `Some content before
<!-- spectr:start -->
First block
<!-- spectr:start -->
Second block`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify - should return error
	if err == nil {
		t.Fatal("Init() expected error for multiple start markers, got nil")
	}

	if !strings.Contains(err.Error(), "multiple start markers") {
		t.Errorf("Init() error = %v, want error containing 'multiple start markers'", err)
	}
}

func TestConfigFileInitializer_Init_NoDuplicateBlocks(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create existing file with markers
	existingContent := `Some content before
<!-- spectr:start -->
Old content
<!-- spectr:end -->
Some content after`

	if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create template
	tmpl := template.Must(template.New("test.tmpl").Parse("New content"))
	templateRef := domain.TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	// Test - run twice
	init := NewConfigFileInitializer("CLAUDE.md", templateRef)
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() first run error = %v", err)
	}

	_, err = init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() second run error = %v", err)
	}

	// Verify file content - should only have one block
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// Count occurrences of start marker
	startCount := strings.Count(string(content), "<!-- spectr:start -->")
	if startCount != 1 {
		t.Errorf("found %d start markers, want 1", startCount)
	}

	// Count occurrences of end marker
	endCount := strings.Count(string(content), "<!-- spectr:end -->")
	if endCount != 1 {
		t.Errorf("found %d end markers, want 1", endCount)
	}
}

func TestConfigFileInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name       string
		fileExists bool
		want       bool
	}{
		{
			name:       "returns true if file exists",
			fileExists: true,
			want:       true,
		},
		{
			name:       "returns false if file does not exist",
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
				if err := afero.WriteFile(projectFs, "CLAUDE.md", []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			// Create template
			tmpl := template.Must(template.New("test.tmpl").Parse("Test content"))
			templateRef := domain.TemplateRef{
				Name:     "test.tmpl",
				Template: tmpl,
			}

			// Test
			init := NewConfigFileInitializer("CLAUDE.md", templateRef)
			got := init.IsSetup(projectFs, homeFs, cfg)

			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_DedupeKey(t *testing.T) {
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
			name: "path with slashes",
			path: "docs/CLAUDE.md",
			want: "ConfigFileInitializer:docs/CLAUDE.md",
		},
		{
			name: "normalizes path",
			path: "./CLAUDE.md",
			want: "ConfigFileInitializer:CLAUDE.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := template.Must(template.New("test.tmpl").Parse("Test content"))
			templateRef := domain.TemplateRef{
				Name:     "test.tmpl",
				Template: tmpl,
			}

			init := NewConfigFileInitializer(tt.path, templateRef)
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf("dedupeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
