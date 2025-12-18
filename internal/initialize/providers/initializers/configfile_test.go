package initializers

import (
	"context"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

func TestConfigFileInitializer_Init_NewFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	tm.instructionPointer = "# Spectr Instructions\n\nTest content."
	ctx := context.Background()

	init := NewConfigFileInitializer("CLAUDE.md")

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify file was created
	exists, err := afero.Exists(fs, "CLAUDE.md")
	if err != nil {
		t.Errorf("Exists() error = %v", err)
		return
	}
	if !exists {
		t.Error("File CLAUDE.md was not created")
		return
	}

	// Verify content
	content, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	contentStr := string(content)

	// Check for markers
	if !strings.Contains(contentStr, spectrStartMarker) {
		t.Error("File should contain start marker")
	}
	if !strings.Contains(contentStr, spectrEndMarker) {
		t.Error("File should contain end marker")
	}

	// Check for content
	if !strings.Contains(contentStr, "# Spectr Instructions") {
		t.Error("File should contain instruction pointer content")
	}
}

func TestConfigFileInitializer_Init_UpdateExistingWithMarkers(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	// Create existing file with markers
	existingContent := `# My Custom Instructions

Some user content above.

<!-- spectr:START -->
Old spectr content that should be replaced.
<!-- spectr:END -->

User content below that should be preserved.
`
	if err := afero.WriteFile(fs, "CLAUDE.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing file: %v", err)
	}

	tm.instructionPointer = "New instruction pointer content"

	init := NewConfigFileInitializer("CLAUDE.md")

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Read updated content
	content, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	contentStr := string(content)

	// Verify user content is preserved
	if !strings.Contains(contentStr, "My Custom Instructions") {
		t.Error("User content above markers should be preserved")
	}
	if !strings.Contains(contentStr, "User content below that should be preserved") {
		t.Error("User content below markers should be preserved")
	}

	// Verify old content is replaced
	if strings.Contains(contentStr, "Old spectr content that should be replaced") {
		t.Error("Old spectr content should be replaced")
	}

	// Verify new content is present
	if !strings.Contains(contentStr, "New instruction pointer content") {
		t.Error("New instruction pointer content should be present")
	}
}

func TestConfigFileInitializer_Init_AppendToFileWithoutMarkers(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	// Create existing file without markers
	existingContent := `# My Project README

This is my project documentation.
It does not have spectr markers.
`
	if err := afero.WriteFile(fs, "README.md", []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing file: %v", err)
	}

	tm.instructionPointer = "Spectr instructions appended"

	init := NewConfigFileInitializer("README.md")

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Read updated content
	content, err := afero.ReadFile(fs, "README.md")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	contentStr := string(content)

	// Verify original content is preserved
	if !strings.Contains(contentStr, "My Project README") {
		t.Error("Original content should be preserved")
	}
	if !strings.Contains(contentStr, "It does not have spectr markers") {
		t.Error("Original content should be preserved")
	}

	// Verify markers and new content are appended
	if !strings.Contains(contentStr, spectrStartMarker) {
		t.Error("Start marker should be appended")
	}
	if !strings.Contains(contentStr, spectrEndMarker) {
		t.Error("End marker should be appended")
	}
	if !strings.Contains(contentStr, "Spectr instructions appended") {
		t.Error("New spectr content should be appended")
	}

	// Verify markers come after original content
	originalIndex := strings.Index(contentStr, "My Project README")
	markerIndex := strings.Index(contentStr, spectrStartMarker)
	if markerIndex < originalIndex {
		t.Error("Markers should be appended after original content")
	}
}

func TestConfigFileInitializer_Init_CreatesParentDirectories(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	init := NewConfigFileInitializer("docs/nested/INSTRUCTIONS.md")

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify parent directories were created
	dirs := []string{"docs", "docs/nested"}
	for _, dir := range dirs {
		exists, err := afero.DirExists(fs, dir)
		if err != nil {
			t.Errorf("DirExists(%s) error = %v", dir, err)
			continue
		}
		if !exists {
			t.Errorf("Parent directory %s was not created", dir)
		}
	}

	// Verify file exists
	exists, err := afero.Exists(fs, "docs/nested/INSTRUCTIONS.md")
	if err != nil {
		t.Errorf("Exists() error = %v", err)
		return
	}
	if !exists {
		t.Error("File was not created")
	}
}

func TestConfigFileInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		fileContent string // Empty means file doesn't exist
		want        bool
	}{
		{
			name:        "returns false when file does not exist",
			path:        "CLAUDE.md",
			fileContent: "",
			want:        false,
		},
		{
			name:        "returns false when file exists but no markers",
			path:        "CLAUDE.md",
			fileContent: "# Some content without markers",
			want:        false,
		},
		{
			name:        "returns false when file has only start marker",
			path:        "CLAUDE.md",
			fileContent: "# Content\n<!-- spectr:START -->\nMore content",
			want:        false,
		},
		{
			name:        "returns false when file has only end marker",
			path:        "CLAUDE.md",
			fileContent: "# Content\n<!-- spectr:END -->\nMore content",
			want:        false,
		},
		{
			name:        "returns true when file exists with both markers",
			path:        "CLAUDE.md",
			fileContent: "# Content\n<!-- spectr:START -->\nSpectr content\n<!-- spectr:END -->\nMore content",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &providers.Config{SpectrDir: "spectr"}

			// Create file if content is provided
			if tt.fileContent != "" {
				if err := afero.WriteFile(fs, tt.path, []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			}

			init := NewConfigFileInitializer(tt.path)

			got := init.IsSetup(fs, cfg)
			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_Path(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "returns file path",
			path: "CLAUDE.md",
			want: "CLAUDE.md",
		},
		{
			name: "returns nested path",
			path: "docs/AI_INSTRUCTIONS.md",
			want: "docs/AI_INSTRUCTIONS.md",
		},
		{
			name: "returns cursorrules path",
			path: ".cursorrules",
			want: ".cursorrules",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewConfigFileInitializer(tt.path)

			got := init.Path()
			if got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_IsGlobal(t *testing.T) {
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
			var init *ConfigFileInitializer
			if tt.global {
				init = NewGlobalConfigFileInitializer(".config/aider/aider.conf.yml")
			} else {
				init = NewConfigFileInitializer("CLAUDE.md")
			}

			got := init.IsGlobal()
			if got != tt.want {
				t.Errorf("IsGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigFileInitializer_Init_Idempotent(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	tm.instructionPointer = "Idempotent test content"
	ctx := context.Background()

	init := NewConfigFileInitializer("CLAUDE.md")

	// Run Init multiple times
	for i := 0; i < 3; i++ {
		err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Errorf("Init() run %d: error = %v", i+1, err)
		}
	}

	// Verify file has correct content (not duplicated)
	content, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	contentStr := string(content)

	// Count marker occurrences - should be exactly one pair
	startCount := strings.Count(contentStr, spectrStartMarker)
	endCount := strings.Count(contentStr, spectrEndMarker)

	if startCount != 1 {
		t.Errorf("Should have exactly 1 start marker, got %d", startCount)
	}
	if endCount != 1 {
		t.Errorf("Should have exactly 1 end marker, got %d", endCount)
	}

	// Content should appear exactly once
	contentCount := strings.Count(contentStr, "Idempotent test content")
	if contentCount != 1 {
		t.Errorf("Content should appear exactly once, appeared %d times", contentCount)
	}
}

func TestConfigFileInitializer_Init_PreservesMarkerOrder(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	tm.instructionPointer = "Test content"
	ctx := context.Background()

	init := NewConfigFileInitializer("CLAUDE.md")

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	content, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	contentStr := string(content)

	startIndex := strings.Index(contentStr, spectrStartMarker)
	endIndex := strings.Index(contentStr, spectrEndMarker)

	if startIndex >= endIndex {
		t.Error("Start marker should come before end marker")
	}
}
