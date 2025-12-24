package initializers

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestConfigFileInitializer_Init_CreateNew(
	t *testing.T,
) {
	// Test the basic file creation logic through helper function
	// We test the full Init separately with actual templates

	contentStr := "existing content"
	newContent := "new spectr content"

	updated, wasUpdated := updateBetweenMarkers(
		contentStr,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated {
		t.Error(
			"should report content was updated",
		)
	}

	if !strings.Contains(
		updated,
		SpectrStartMarker,
	) {
		t.Error(
			"updated content should contain start marker",
		)
	}

	if !strings.Contains(
		updated,
		SpectrEndMarker,
	) {
		t.Error(
			"updated content should contain end marker",
		)
	}

	if !strings.Contains(updated, newContent) {
		t.Error(
			"updated content should contain new content",
		)
	}

	if !strings.Contains(
		updated,
		"existing content",
	) {
		t.Error(
			"updated content should preserve existing content",
		)
	}
}

func TestConfigFileInitializer_UpdateExisting(
	t *testing.T,
) {
	// Test updating content between existing markers
	existingContent := `# Configuration

` + SpectrStartMarker + `
old content
` + SpectrEndMarker + `

More content here`

	newContent := "updated content"

	updated, wasUpdated := updateBetweenMarkers(
		existingContent,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated {
		t.Error(
			"should report content was updated",
		)
	}

	if strings.Contains(updated, "old content") {
		t.Error(
			"updated content should not contain old content",
		)
	}

	if !strings.Contains(updated, newContent) {
		t.Error(
			"updated content should contain new content",
		)
	}

	if !strings.Contains(
		updated,
		"# Configuration",
	) {
		t.Error(
			"updated content should preserve content before markers",
		)
	}

	if !strings.Contains(
		updated,
		"More content here",
	) {
		t.Error(
			"updated content should preserve content after markers",
		)
	}
}

func TestConfigFileInitializer_NoMarkersInExisting(
	t *testing.T,
) {
	// Test appending markers when they don't exist
	existingContent := "# Configuration\n\nExisting content"
	newContent := "spectr content"

	updated, wasUpdated := updateBetweenMarkers(
		existingContent,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated {
		t.Error(
			"should report content was updated",
		)
	}

	if !strings.Contains(
		updated,
		existingContent,
	) {
		t.Error(
			"updated content should preserve existing content",
		)
	}

	if !strings.Contains(
		updated,
		SpectrStartMarker,
	) {
		t.Error(
			"updated content should contain start marker",
		)
	}

	if !strings.Contains(
		updated,
		SpectrEndMarker,
	) {
		t.Error(
			"updated content should contain end marker",
		)
	}

	if !strings.Contains(updated, newContent) {
		t.Error(
			"updated content should contain new content",
		)
	}

	// Markers should be appended at the end
	startIdx := strings.Index(
		updated,
		SpectrStartMarker,
	)
	existingIdx := strings.Index(
		updated,
		existingContent,
	)
	if startIdx < existingIdx {
		t.Error(
			"markers should be appended after existing content",
		)
	}
}

func TestConfigFileInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name        string
		fileExists  bool
		fileContent string
		wantSetup   bool
	}{
		{
			name:       "file does not exist",
			fileExists: false,
			wantSetup:  false,
		},
		{
			name:        "file exists without markers",
			fileExists:  true,
			fileContent: "# Configuration\n\nNo markers here",
			wantSetup:   false,
		},
		{
			name:        "file exists with start marker only",
			fileExists:  true,
			fileContent: "# Config\n" + SpectrStartMarker + "\ncontent",
			wantSetup:   false,
		},
		{
			name:        "file exists with end marker only",
			fileExists:  true,
			fileContent: "# Config\ncontent\n" + SpectrEndMarker,
			wantSetup:   false,
		},
		{
			name:       "file exists with both markers",
			fileExists: true,
			fileContent: "# Config\n" + SpectrStartMarker +
				"\ncontent\n" + SpectrEndMarker,
			wantSetup: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := DefaultConfig()

			if tt.fileExists {
				if err := afero.WriteFile(fs, "test.md", []byte(tt.fileContent), 0644); err != nil {
					t.Fatalf(
						"failed to create test file: %v",
						err,
					)
				}
			}

			init := &ConfigFileInitializer{
				path: "test.md",
			}

			got := init.IsSetup(fs, cfg)
			if got != tt.wantSetup {
				t.Errorf(
					"IsSetup() = %v, want %v",
					got,
					tt.wantSetup,
				)
			}
		})
	}
}

func TestConfigFileInitializer_Path(
	t *testing.T,
) {
	init := &ConfigFileInitializer{
		path: "CLAUDE.md",
	}

	if got := init.Path(); got != "CLAUDE.md" {
		t.Errorf(
			"Path() = %v, want CLAUDE.md",
			got,
		)
	}
}

func TestConfigFileInitializer_IsGlobal(
	t *testing.T,
) {
	init := &ConfigFileInitializer{
		path:     "CLAUDE.md",
		isGlobal: false,
	}

	if init.IsGlobal() {
		t.Error(
			"IsGlobal() should return false by default",
		)
	}

	init.isGlobal = true
	if !init.IsGlobal() {
		t.Error(
			"IsGlobal() should return true when set",
		)
	}
}

func TestUpdateBetweenMarkers_Idempotent(
	t *testing.T,
) {
	initialContent := "# Config\n\nInitial content"
	newContent := "spectr content"

	// First update - adds markers
	updated1, wasUpdated1 := updateBetweenMarkers(
		initialContent,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated1 {
		t.Error(
			"first update should report change",
		)
	}

	// Second update with same content - should be idempotent
	updated2, wasUpdated2 := updateBetweenMarkers(
		updated1,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if wasUpdated2 {
		t.Error(
			"second update should not report change (idempotent)",
		)
	}

	if updated1 != updated2 {
		t.Error(
			"content should be identical after idempotent update",
		)
	}
}

func TestUpdateBetweenMarkers_PreservesStructure(
	t *testing.T,
) {
	existingContent := `# Configuration File

Some preamble text.

` + SpectrStartMarker + `
old spectr content
` + SpectrEndMarker + `

Some epilogue text.

## Another Section

More content.`

	newContent := "new spectr content"

	updated, _ := updateBetweenMarkers(
		existingContent,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	// Check structure is preserved
	if !strings.Contains(
		updated,
		"# Configuration File",
	) {
		t.Error("should preserve title")
	}

	if !strings.Contains(
		updated,
		"Some preamble text.",
	) {
		t.Error("should preserve preamble")
	}

	if !strings.Contains(
		updated,
		"Some epilogue text.",
	) {
		t.Error("should preserve epilogue")
	}

	if !strings.Contains(
		updated,
		"## Another Section",
	) {
		t.Error("should preserve other sections")
	}

	if !strings.Contains(updated, newContent) {
		t.Error("should contain new content")
	}

	if strings.Contains(
		updated,
		"old spectr content",
	) {
		t.Error("should not contain old content")
	}
}

func TestUpdateBetweenMarkers_OrphanedStartMarker(
	t *testing.T,
) {
	// Start marker exists but no end marker - should replace from start to EOF
	existingContent := `# Configuration

Some preamble text.

` + SpectrStartMarker + `
orphaned content without end marker
more orphaned content`

	newContent := "new spectr content"

	updated, wasUpdated := updateBetweenMarkers(
		existingContent,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated {
		t.Error(
			"should report content was updated",
		)
	}

	// Should preserve content before start marker
	if !strings.Contains(
		updated,
		"Some preamble text.",
	) {
		t.Error("should preserve preamble")
	}

	// Should contain new content with proper markers
	if !strings.Contains(updated, newContent) {
		t.Error("should contain new content")
	}

	// Should have both markers
	if !strings.Contains(
		updated,
		SpectrStartMarker,
	) {
		t.Error("should contain start marker")
	}
	if !strings.Contains(
		updated,
		SpectrEndMarker,
	) {
		t.Error("should contain end marker")
	}

	// Should NOT contain orphaned content
	if strings.Contains(
		updated,
		"orphaned content",
	) {
		t.Error(
			"should not contain orphaned content",
		)
	}

	// Should have exactly one start marker
	if strings.Count(
		updated,
		SpectrStartMarker,
	) != 1 {
		t.Errorf(
			"should have exactly one start marker, got %d",
			strings.Count(
				updated,
				SpectrStartMarker,
			),
		)
	}
}
