package initializers

import (
	"context"
	"html/template"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

func newTestTemplate(t *testing.T, name, content string) domain.TemplateRef {
	t.Helper()
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		t.Fatalf("Failed to parse test template: %v", err)
	}

	return domain.TemplateRef{
		Name:     name,
		Template: tmpl,
	}
}

func TestConfigFileInitializer_Init_NewFile(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	tmplRef := newTestTemplate(t, "test.md.tmpl", "# Test Content\nBaseDir: {{.BaseDir}}")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 1 || result.CreatedFiles[0] != "CLAUDE.md" {
		t.Errorf("Init() CreatedFiles = %v, want [CLAUDE.md]", result.CreatedFiles)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "<!-- spectr:start -->") {
		t.Error("Created file should contain start marker")
	}
	if !strings.Contains(contentStr, "<!-- spectr:end -->") {
		t.Error("Created file should contain end marker")
	}
	if !strings.Contains(contentStr, "BaseDir: spectr") {
		t.Error("Created file should contain rendered template with BaseDir")
	}
}

func TestConfigFileInitializer_Init_UpdateBetweenMarkers(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create existing file with markers
	existingContent := `# My Config
Some user content here.

<!-- spectr:start -->
Old content to be replaced
<!-- spectr:end -->

More user content below.`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.UpdatedFiles) != 1 || result.UpdatedFiles[0] != "CLAUDE.md" {
		t.Errorf("Init() UpdatedFiles = %v, want [CLAUDE.md]", result.UpdatedFiles)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Some user content here") {
		t.Error("Content before markers should be preserved")
	}
	if !strings.Contains(contentStr, "More user content below") {
		t.Error("Content after markers should be preserved")
	}
	if !strings.Contains(contentStr, "NEW CONTENT") {
		t.Error("New content should be present")
	}
	if strings.Contains(contentStr, "Old content to be replaced") {
		t.Error("Old content between markers should be replaced")
	}
}

func TestConfigFileInitializer_Init_OrphanedStartWithTrailingEnd(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Orphaned start marker followed eventually by end marker
	existingContent := `# My Config
<!-- spectr:start -->
Some content

<!-- spectr:end -->
Trailing content`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.UpdatedFiles) != 1 {
		t.Errorf("Init() UpdatedFiles = %v, want [CLAUDE.md]", result.UpdatedFiles)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "NEW CONTENT") {
		t.Error("New content should be present")
	}
	if !strings.Contains(contentStr, "Trailing content") {
		t.Error("Content after end marker should be preserved")
	}
}

func TestConfigFileInitializer_Init_OrphanedStartNoEnd(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Orphaned start marker with no end marker
	existingContent := `# My Config
<!-- spectr:start -->
Old incomplete block that never had an end marker
More stuff here`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.UpdatedFiles) != 1 {
		t.Errorf("Init() UpdatedFiles = %v, want [CLAUDE.md]", result.UpdatedFiles)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "<!-- spectr:start -->") {
		t.Error("Should have start marker")
	}
	if !strings.Contains(contentStr, "<!-- spectr:end -->") {
		t.Error("Should have end marker after fix")
	}
	if !strings.Contains(contentStr, "NEW CONTENT") {
		t.Error("New content should be present")
	}
	if strings.Contains(contentStr, "Old incomplete block") {
		t.Error("Old orphaned content should be replaced")
	}
}

func TestConfigFileInitializer_Init_NoDuplicateBlocks(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// File with existing markers
	existingContent := `# My Config
<!-- spectr:start -->
Existing content
<!-- spectr:end -->
`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	// Run init twice
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("First Init() error = %v", err)
	}
	_, err = init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Second Init() error = %v", err)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	startCount := strings.Count(contentStr, "<!-- spectr:start -->")
	endCount := strings.Count(contentStr, "<!-- spectr:end -->")

	if startCount != 1 {
		t.Errorf("Should have exactly 1 start marker, got %d", startCount)
	}
	if endCount != 1 {
		t.Errorf("Should have exactly 1 end marker, got %d", endCount)
	}
}

func TestConfigFileInitializer_Init_TemplateRef(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "myproject"}

	tmplRef := newTestTemplate(t, "test.md.tmpl", `BaseDir: {{.BaseDir}}
SpecsDir: {{.SpecsDir}}
ChangesDir: {{.ChangesDir}}
ProjectFile: {{.ProjectFile}}
AgentsFile: {{.AgentsFile}}`)

	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	expected := []string{
		"BaseDir: myproject",
		"SpecsDir: myproject/specs",
		"ChangesDir: myproject/changes",
		"ProjectFile: myproject/project.md",
		"AgentsFile: myproject/AGENTS.md",
	}

	for _, exp := range expected {
		if !strings.Contains(contentStr, exp) {
			t.Errorf("Content should contain %q", exp)
		}
	}
}

func TestConfigFileInitializer_Init_OrphanedEndMarker(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Orphaned end marker (no start)
	existingContent := `# My Config
Some content
<!-- spectr:end -->
More content`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err == nil {
		t.Error("Init() should return error for orphaned end marker")
	}
	if !strings.Contains(err.Error(), "orphaned end marker") {
		t.Errorf("Error should mention orphaned end marker, got: %v", err)
	}
}

func TestConfigFileInitializer_Init_NestedStartMarkers(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Nested start marker before end
	existingContent := `# My Config
<!-- spectr:start -->
Some content
<!-- spectr:start -->
Nested start
<!-- spectr:end -->`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err == nil {
		t.Error("Init() should return error for nested start markers")
	}
	if !strings.Contains(err.Error(), "nested start marker") {
		t.Errorf("Error should mention nested start marker, got: %v", err)
	}
}

func TestConfigFileInitializer_Init_MultipleStartMarkers(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Multiple start markers without end between them
	existingContent := `# My Config
<!-- spectr:start -->
First block
<!-- spectr:start -->
Second block (no end for first)`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err == nil {
		t.Error("Init() should return error for multiple start markers")
	}
	if !strings.Contains(err.Error(), "multiple start markers") {
		t.Errorf("Error should mention multiple start markers, got: %v", err)
	}
}

func TestConfigFileInitializer_Init_CaseInsensitiveRead(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Uppercase markers (old format)
	existingContent := `# My Config
<!-- spectr:START -->
Old content
<!-- spectr:END -->`
	_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(existingContent), 0o644)

	tmplRef := newTestTemplate(t, "test.md.tmpl", "NEW CONTENT")
	init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	// Should write lowercase markers
	if !strings.Contains(contentStr, "<!-- spectr:start -->") {
		t.Error("Should write lowercase start marker")
	}
	if !strings.Contains(contentStr, "<!-- spectr:end -->") {
		t.Error("Should write lowercase end marker")
	}
	if !strings.Contains(contentStr, "NEW CONTENT") {
		t.Error("Should contain new content")
	}
}

func TestConfigFileInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "returns true when lowercase markers exist",
			content: "# Config\n<!-- spectr:start -->\nContent\n<!-- spectr:end -->",
			want:    true,
		},
		{
			name:    "returns true when uppercase markers exist (case-insensitive)",
			content: "# Config\n<!-- spectr:START -->\nContent\n<!-- spectr:END -->",
			want:    true,
		},
		{
			name:    "returns false when no markers exist",
			content: "# Config\nNo markers here",
			want:    false,
		},
		{
			name:    "returns false when file doesn't exist",
			content: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &domain.Config{SpectrDir: "spectr"}

			if tt.content != "" {
				_ = afero.WriteFile(projectFs, "CLAUDE.md", []byte(tt.content), 0o644)
			}

			tmplRef := newTestTemplate(t, "test.md.tmpl", "TEST")
			init := NewConfigFileInitializer("CLAUDE.md", tmplRef)

			if got := init.IsSetup(projectFs, homeFs, cfg); got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_dedupeKey(t *testing.T) {
	tmplRef := newTestTemplate(t, "test.md.tmpl", "TEST")
	init, ok := NewConfigFileInitializer("CLAUDE.md", tmplRef).(*ConfigFileInitializer)
	if !ok {
		t.Fatal("NewConfigFileInitializer did not return *ConfigFileInitializer")
	}

	want := "ConfigFileInitializer:CLAUDE.md"
	if got := init.DedupeKey(); got != want {
		t.Errorf("dedupeKey() = %q, want %q", got, want)
	}
}

func TestUpdateWithMarkers(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		newContent string
		want       string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "append when no markers exist",
			content:    "# Existing content",
			newContent: "NEW",
			want:       "# Existing content\n\n<!-- spectr:start -->\nNEW\n<!-- spectr:end -->",
		},
		{
			name:       "replace content between markers",
			content:    "Before\n<!-- spectr:start -->\nOLD\n<!-- spectr:end -->\nAfter",
			newContent: "NEW",
			want:       "Before\n<!-- spectr:start -->\nNEW\n<!-- spectr:end -->\nAfter",
		},
		{
			name:       "error on orphaned end marker",
			content:    "Content\n<!-- spectr:end -->\nMore",
			newContent: "NEW",
			wantErr:    true,
			errMsg:     "orphaned end marker",
		},
		{
			name:       "error on nested start markers",
			content:    "<!-- spectr:start -->\nContent\n<!-- spectr:start -->\nNested\n<!-- spectr:end -->",
			newContent: "NEW",
			wantErr:    true,
			errMsg:     "nested start marker",
		},
		{
			name:       "orphaned start replaces to end",
			content:    "Before\n<!-- spectr:start -->\nOrphaned content with no end",
			newContent: "NEW",
			want:       "Before\n<!-- spectr:start -->\nNEW\n<!-- spectr:end -->",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := updateWithMarkers(tt.content, tt.newContent)
			if tt.wantErr {
				if err == nil {
					t.Error("updateWithMarkers() expected error, got nil")

					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf(
						"updateWithMarkers() error = %v, want error containing %q",
						err,
						tt.errMsg,
					)
				}

				return
			}
			if err != nil {
				t.Errorf("updateWithMarkers() unexpected error = %v", err)

				return
			}
			if got != tt.want {
				t.Errorf("updateWithMarkers() = %q, want %q", got, tt.want)
			}
		})
	}
}
