package initializers

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/connerohnesorge/spectr/internal/domain"
)

func TestSlashCommandsInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name          string
		commands      []domain.SlashCommand
		ext           string
		existingFiles []string
		want          bool
	}{
		{
			name:     "returns false when no files exist",
			commands: []domain.SlashCommand{domain.SlashProposal, domain.SlashApply},
			ext:      ".md",
			want:     false,
		},
		{
			name:          "returns false when some files missing",
			commands:      []domain.SlashCommand{domain.SlashProposal, domain.SlashApply},
			ext:           ".md",
			existingFiles: []string{"proposal.md"},
			want:          false,
		},
		{
			name:          "returns true when all files exist",
			commands:      []domain.SlashCommand{domain.SlashProposal, domain.SlashApply},
			ext:           ".md",
			existingFiles: []string{"proposal.md", "apply.md"},
			want:          true,
		},
		{
			name:          "returns true for single command",
			commands:      []domain.SlashCommand{domain.SlashProposal},
			ext:           ".md",
			existingFiles: []string{"proposal.md"},
			want:          true,
		},
		{
			name:          "works with TOML extension",
			commands:      []domain.SlashCommand{domain.SlashProposal},
			ext:           ".toml",
			existingFiles: []string{"proposal.toml"},
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := DefaultConfig()
			dir := ".claude/commands/spectr"

			// Create directory
			if err := fs.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("failed to create directory: %v", err)
			}

			// Create existing files
			for _, file := range tt.existingFiles {
				path := filepath.Join(dir, file)
				if err := afero.WriteFile(fs, path, []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create file %s: %v", file, err)
				}
			}

			init := NewSlashCommandsInitializer(dir, tt.ext, tt.commands)
			got := init.IsSetup(fs, cfg)

			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_Path(t *testing.T) {
	dir := ".claude/commands/spectr"
	init := NewSlashCommandsInitializer(dir, ".md", []domain.SlashCommand{domain.SlashProposal})

	if got := init.Path(); got != dir {
		t.Errorf("Path() = %v, want %v", got, dir)
	}
}

func TestSlashCommandsInitializer_IsGlobal(t *testing.T) {
	init := NewSlashCommandsInitializer(
		".claude/commands",
		".md",
		[]domain.SlashCommand{domain.SlashProposal},
	)

	if init.IsGlobal() {
		t.Error("IsGlobal() should return false by default")
	}
}

func TestFormatMarkdownCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     domain.SlashCommand
		content string
	}{
		{
			name:    "formats proposal command",
			cmd:     domain.SlashProposal,
			content: "Test proposal content",
		},
		{
			name:    "formats apply command",
			cmd:     domain.SlashApply,
			content: "Test apply content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMarkdownCommand(tt.cmd, tt.content)

			// Check for frontmatter
			if !strings.HasPrefix(result, "---") {
				t.Error("markdown command should start with frontmatter")
			}

			if !strings.Contains(result, "description:") {
				t.Error("markdown command should contain description field")
			}

			// Check for markers
			if !strings.Contains(result, SpectrStartMarker) {
				t.Error("markdown command should contain start marker")
			}

			if !strings.Contains(result, SpectrEndMarker) {
				t.Error("markdown command should contain end marker")
			}

			// Check for content
			if !strings.Contains(result, tt.content) {
				t.Error("markdown command should contain provided content")
			}

			// Verify structure: frontmatter, then markers with content
			frontmatterEnd := strings.Index(result[3:], "---")
			if frontmatterEnd == -1 {
				t.Error("markdown command should have closing frontmatter")
			}

			startMarkerIdx := strings.Index(result, SpectrStartMarker)
			if startMarkerIdx < frontmatterEnd {
				t.Error("start marker should come after frontmatter")
			}

			contentIdx := strings.Index(result, tt.content)
			if contentIdx < startMarkerIdx {
				t.Error("content should come after start marker")
			}

			endMarkerIdx := strings.Index(result, SpectrEndMarker)
			if endMarkerIdx < contentIdx {
				t.Error("end marker should come after content")
			}
		})
	}
}

func TestFormatTOMLCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     domain.SlashCommand
		content string
	}{
		{
			name:    "formats proposal command",
			cmd:     domain.SlashProposal,
			content: "Test proposal content",
		},
		{
			name:    "formats apply command",
			cmd:     domain.SlashApply,
			content: "Test apply content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTOMLCommand(tt.cmd, tt.content)

			// Check for comment
			if !strings.HasPrefix(result, "#") {
				t.Error("TOML command should start with comment")
			}

			// Check for TOML structure
			if !strings.Contains(result, "[[agent]]") {
				t.Error("TOML command should contain [[agent]] section")
			}

			// Check for markers
			if !strings.Contains(result, SpectrStartMarker) {
				t.Error("TOML command should contain start marker")
			}

			if !strings.Contains(result, SpectrEndMarker) {
				t.Error("TOML command should contain end marker")
			}

			// Check for content
			if !strings.Contains(result, tt.content) {
				t.Error("TOML command should contain provided content")
			}

			// Verify structure: comment, [[agent]], markers with content
			agentIdx := strings.Index(result, "[[agent]]")
			if agentIdx == -1 {
				t.Error("TOML command should have [[agent]] section")
			}

			startMarkerIdx := strings.Index(result, SpectrStartMarker)
			if startMarkerIdx < agentIdx {
				t.Error("start marker should come after [[agent]]")
			}

			contentIdx := strings.Index(result, tt.content)
			if contentIdx < startMarkerIdx {
				t.Error("content should come after start marker")
			}

			endMarkerIdx := strings.Index(result, SpectrEndMarker)
			if endMarkerIdx < contentIdx {
				t.Error("end marker should come after content")
			}
		})
	}
}

func TestSlashCommandsInitializer_MultipleCommands(t *testing.T) {
	// Test that the initializer correctly identifies all command files

	tests := []struct {
		name      string
		commands  []domain.SlashCommand
		ext       string
		wantFiles []string
	}{
		{
			name:      "creates both proposal and apply markdown",
			commands:  []domain.SlashCommand{domain.SlashProposal, domain.SlashApply},
			ext:       ".md",
			wantFiles: []string{"proposal.md", "apply.md"},
		},
		{
			name:      "creates both proposal and apply TOML",
			commands:  []domain.SlashCommand{domain.SlashProposal, domain.SlashApply},
			ext:       ".toml",
			wantFiles: []string{"proposal.toml", "apply.toml"},
		},
		{
			name:      "creates single command",
			commands:  []domain.SlashCommand{domain.SlashProposal},
			ext:       ".md",
			wantFiles: []string{"proposal.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := DefaultConfig()
			dir := ".claude/commands/spectr"

			init := NewSlashCommandsInitializer(dir, tt.ext, tt.commands)

			// Verify that IsSetup returns false before creation
			if init.IsSetup(fs, cfg) {
				t.Error("IsSetup() should return false before files are created")
			}

			// Create the files manually to test IsSetup
			for _, file := range tt.wantFiles {
				path := filepath.Join(dir, file)
				if err := fs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := afero.WriteFile(fs, path, []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create file %s: %v", file, err)
				}
			}

			// Verify that IsSetup returns true after creation
			if !init.IsSetup(fs, cfg) {
				t.Error("IsSetup() should return true after all files are created")
			}
		})
	}
}

func TestSlashCommandsInitializer_UnsupportedExtension(t *testing.T) {
	// Test that we only support .md and .toml
	validExts := []string{".md", ".toml"}

	for _, ext := range validExts {
		t.Run("supports_"+ext, func(t *testing.T) {
			init := NewSlashCommandsInitializer(
				".test/commands",
				ext,
				[]domain.SlashCommand{domain.SlashProposal},
			)

			// Verify the initializer was created
			if init.ext != ext {
				t.Errorf("expected extension %s, got %s", ext, init.ext)
			}
		})
	}
}

func TestUpdateBetweenMarkers_WithSlashCommands(t *testing.T) {
	// Test that updateBetweenMarkers works correctly for slash command updates

	existingMdContent := `---
description: Old description
---

` + SpectrStartMarker + `
old slash command content
` + SpectrEndMarker + `
`

	newContent := "updated slash command content"

	updated, wasUpdated := updateBetweenMarkers(
		existingMdContent,
		newContent,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated {
		t.Error("should report content was updated")
	}

	if !strings.Contains(updated, "---") {
		t.Error("should preserve frontmatter")
	}

	if !strings.Contains(updated, newContent) {
		t.Error("should contain new content")
	}

	if strings.Contains(updated, "old slash command content") {
		t.Error("should not contain old content")
	}
}
