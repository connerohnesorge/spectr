package initializers

import (
	"context"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

func TestNewConfigFileInitializer(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		template string
	}{
		{
			name:     "CLAUDE.md with template",
			path:     "CLAUDE.md",
			template: "# Spectr Instructions",
		},
		{
			name:     "nested path with template",
			path:     ".cursor/rules/spectr.md",
			template: "Some template content",
		},
		{
			name:     "empty template",
			path:     "TEST.md",
			template: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfigFileInitializer(
				tt.path,
				tt.template,
			)

			if c == nil {
				t.Fatal(
					"NewConfigFileInitializer() returned nil",
				)
			}

			if c.Path != tt.path {
				t.Errorf(
					"Path = %s, want %s",
					c.Path,
					tt.path,
				)
			}

			if c.Template != tt.template {
				t.Errorf(
					"Template = %s, want %s",
					c.Template,
					tt.template,
				)
			}
		})
	}
}

func TestConfigFileInitializer_Init_NewFile(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	template := "# Spectr Instructions\nRead spectr/AGENTS.md"
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		template,
	)

	err := c.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify file was created
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("File was not created: %v", err)
	}

	contentStr := string(content)

	// Verify markers are present
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

	// Verify template content is present
	if !strings.Contains(contentStr, template) {
		t.Error(
			"File should contain template content",
		)
	}
}

func TestConfigFileInitializer_Init_NewFileInNestedDirectory(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	template := "Rules content"
	c := NewConfigFileInitializer(
		".cursor/rules/spectr.md",
		template,
	)

	err := c.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify file was created
	exists, err := afero.Exists(
		fs,
		".cursor/rules/spectr.md",
	)
	if err != nil {
		t.Fatalf(
			"Error checking file existence: %v",
			err,
		)
	}

	if !exists {
		t.Error(
			"File was not created in nested directory",
		)
	}

	// Verify directory was created
	info, err := fs.Stat(".cursor/rules")
	if err != nil {
		t.Fatalf(
			"Directory was not created: %v",
			err,
		)
	}

	if !info.IsDir() {
		t.Error("Parent path is not a directory")
	}
}

func TestConfigFileInitializer_Init_AppendsToExistingFileWithoutMarkers(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	// Create existing file without markers
	existingContent := "# Existing Content\nSome instructions"
	err := afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte(existingContent),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create existing file: %v",
			err,
		)
	}

	template := "# Spectr Instructions"
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		template,
	)

	err = c.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify content
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)

	// Verify existing content is preserved
	if !strings.Contains(
		contentStr,
		"# Existing Content",
	) {
		t.Error(
			"Existing content should be preserved",
		)
	}

	// Verify markers and template are appended
	if !strings.Contains(
		contentStr,
		spectrStartMarker,
	) {
		t.Error("Start marker should be appended")
	}

	if !strings.Contains(
		contentStr,
		spectrEndMarker,
	) {
		t.Error("End marker should be appended")
	}

	if !strings.Contains(contentStr, template) {
		t.Error(
			"Template content should be appended",
		)
	}

	// Verify order: existing content comes before markers
	existingIdx := strings.Index(
		contentStr,
		"# Existing Content",
	)
	markerIdx := strings.Index(
		contentStr,
		spectrStartMarker,
	)
	if existingIdx > markerIdx {
		t.Error(
			"Existing content should come before new markers",
		)
	}
}

func TestConfigFileInitializer_Init_UpdatesExistingMarkers(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	// Create existing file with markers and old content
	existingContent := "# Header\n\n" +
		spectrStartMarker + "\nOld spectr content\n" + spectrEndMarker + "\n\n# Footer"
	err := afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte(existingContent),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create existing file: %v",
			err,
		)
	}

	newTemplate := "New spectr content"
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		newTemplate,
	)

	err = c.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify content
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)

	// Verify old content is replaced
	if strings.Contains(
		contentStr,
		"Old spectr content",
	) {
		t.Error(
			"Old content between markers should be replaced",
		)
	}

	// Verify new content is present
	if !strings.Contains(
		contentStr,
		newTemplate,
	) {
		t.Error(
			"New template content should be present",
		)
	}

	// Verify header and footer are preserved
	if !strings.Contains(contentStr, "# Header") {
		t.Error("Header should be preserved")
	}

	if !strings.Contains(contentStr, "# Footer") {
		t.Error("Footer should be preserved")
	}
}

func TestConfigFileInitializer_Init_Idempotent(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	template := "# Spectr Instructions"
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		template,
	)

	// Call Init multiple times
	for i := range 3 {
		err := c.Init(ctx, fs, cfg)
		if err != nil {
			t.Fatalf(
				"Init() call %d failed: %v",
				i+1,
				err,
			)
		}
	}

	// Verify content is correct and not duplicated
	content, err := afero.ReadFile(
		fs,
		"CLAUDE.md",
	)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)

	// Count markers - should be exactly one of each
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

	// Template should appear exactly once
	templateCount := strings.Count(
		contentStr,
		template,
	)
	if templateCount != 1 {
		t.Errorf(
			"Template count = %d, want 1",
			templateCount,
		)
	}
}

func TestConfigFileInitializer_IsSetup_FileNotExists(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	c := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)

	// File doesn't exist
	if c.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when file doesn't exist",
		)
	}
}

func TestConfigFileInitializer_IsSetup_FileExistsWithMarkers(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	template := "# Spectr Instructions"
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		template,
	)

	// Create file with Init
	err := c.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Now IsSetup should return true
	if !c.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return true when file exists with markers",
		)
	}
}

func TestConfigFileInitializer_IsSetup_FileExistsWithoutMarkers(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	// Create file without markers
	err := afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte("# No markers here"),
		0644,
	)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	c := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)

	// IsSetup should return false without markers
	if c.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when file exists without markers",
		)
	}
}

func TestConfigFileInitializer_IsSetup_FileExistsWithOnlyStartMarker(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	// Create file with only start marker
	content := "# Header\n" + spectrStartMarker + "\nContent"
	err := afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	c := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)

	// IsSetup should return false with only start marker
	if c.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when file has only start marker",
		)
	}
}

func TestConfigFileInitializer_IsSetup_FileExistsWithOnlyEndMarker(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	// Create file with only end marker
	content := "Content\n" + spectrEndMarker + "\n# Footer"
	err := afero.WriteFile(
		fs,
		"CLAUDE.md",
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	c := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)

	// IsSetup should return false with only end marker
	if c.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when file has only end marker",
		)
	}
}

func TestConfigFileInitializer_Key_Simple(
	t *testing.T,
) {
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)

	key := c.Key()
	expected := "config:CLAUDE.md"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestConfigFileInitializer_Key_NestedPath(
	t *testing.T,
) {
	c := NewConfigFileInitializer(
		".cursor/rules/spectr.md",
		"template",
	)

	key := c.Key()
	expected := "config:.cursor/rules/spectr.md"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestConfigFileInitializer_Key_Consistent(
	t *testing.T,
) {
	c := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)

	// Key should be consistent across multiple calls
	key1 := c.Key()
	key2 := c.Key()
	key3 := c.Key()

	if key1 != key2 || key2 != key3 {
		t.Errorf(
			"Key() is not consistent: %s, %s, %s",
			key1,
			key2,
			key3,
		)
	}
}

func TestConfigFileInitializer_Key_DifferentPaths(
	t *testing.T,
) {
	c1 := NewConfigFileInitializer(
		"CLAUDE.md",
		"template",
	)
	c2 := NewConfigFileInitializer(
		"AGENTS.md",
		"template",
	)

	// Keys should differ for different paths
	if c1.Key() == c2.Key() {
		t.Errorf(
			"Keys should differ for different paths: %s vs %s",
			c1.Key(),
			c2.Key(),
		)
	}
}

func TestConfigFileInitializer_Key_SamePathDifferentTemplates(
	t *testing.T,
) {
	c1 := NewConfigFileInitializer(
		"CLAUDE.md",
		"template1",
	)
	c2 := NewConfigFileInitializer(
		"CLAUDE.md",
		"template2",
	)

	// Keys should be the same (based only on path, not template)
	if c1.Key() != c2.Key() {
		t.Errorf(
			"Keys should be same for same path: %s vs %s",
			c1.Key(),
			c2.Key(),
		)
	}
}

func TestConfigFileInitializer_ImplementsInterface(
	_ *testing.T,
) {
	// Compile-time check is in configfile.go, but this is a runtime verification
	var _ providers.Initializer = (*ConfigFileInitializer)(nil)
}

func TestFindMarkerIndex(t *testing.T) {
	tests := []struct {
		name    string
		content string
		marker  string
		offset  int
		want    int
	}{
		{
			name:    "marker at start",
			content: spectrStartMarker + "\ncontent",
			marker:  spectrStartMarker,
			offset:  0,
			want:    0,
		},
		{
			name:    "marker in middle",
			content: "before\n" + spectrStartMarker + "\nafter",
			marker:  spectrStartMarker,
			offset:  0,
			want:    7,
		},
		{
			name:    "marker not found",
			content: "no markers here",
			marker:  spectrStartMarker,
			offset:  0,
			want:    -1,
		},
		{
			name:    "search with offset",
			content: spectrStartMarker + "\n" + spectrStartMarker,
			marker:  spectrStartMarker,
			offset:  len(spectrStartMarker),
			want:    len(spectrStartMarker) + 1,
		},
		{
			name:    "offset past content",
			content: spectrStartMarker,
			marker:  spectrStartMarker,
			offset:  100,
			want:    -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMarkerIndex(
				tt.content,
				tt.marker,
				tt.offset,
			)
			if got != tt.want {
				t.Errorf(
					"findMarkerIndex() = %d, want %d",
					got,
					tt.want,
				)
			}
		})
	}
}
