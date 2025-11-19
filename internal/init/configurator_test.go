package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test UpdateFileWithMarkers with new file creation
func TestUpdateFileWithMarkers_NewFile(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	content := "This is test content"

	err := UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
	if err != nil {
		t.Fatalf("UpdateFileWithMarkers failed: %v", err)
	}

	// Verify file exists
	if !FileExists(filePath) {
		t.Fatal("File was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	result := string(data)
	expected := SpectrStartMarker +
		"\n" + content + "\n" + SpectrEndMarker

	if result != expected {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

// Test UpdateFileWithMarkers with existing file and markers
func TestUpdateFileWithMarkers_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create initial file with markers
	initialContent := "Initial content"
	initial := SpectrStartMarker + "\n" + initialContent + "\n" + SpectrEndMarker
	if err := os.WriteFile(filePath, []byte(initial), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Update with new content
	newContent := "Updated content"
	err := UpdateFileWithMarkers(filePath, newContent, SpectrStartMarker, SpectrEndMarker)
	if err != nil {
		t.Fatalf("UpdateFileWithMarkers failed: %v", err)
	}

	// Read and verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	result := string(data)
	expected := SpectrStartMarker + "\n" + newContent + "\n" + SpectrEndMarker

	if result != expected {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

// Test UpdateFileWithMarkers preserves content outside markers
func TestUpdateFileWithMarkers_PreservesOutsideContent(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create initial file with content before and after markers
	before := "# Header\n\nSome intro text\n\n"
	after := "\n\n## Footer\n\nSome footer text"
	markedContent := "Initial managed content"
	initial := before + SpectrStartMarker + "\n" + markedContent + "\n" + SpectrEndMarker + after

	if err := os.WriteFile(filePath, []byte(initial), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Update with new content
	newContent := "Updated managed content"
	err := UpdateFileWithMarkers(filePath, newContent, SpectrStartMarker, SpectrEndMarker)
	if err != nil {
		t.Fatalf("UpdateFileWithMarkers failed: %v", err)
	}

	// Read and verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	result := string(data)
	expected := before + SpectrStartMarker + "\n" + newContent + "\n" + SpectrEndMarker + after

	if result != expected {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

// Test UpdateFileWithMarkers with file without markers
func TestUpdateFileWithMarkers_PrependToExisting(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create file without markers
	existing := "Existing content without markers"
	if err := os.WriteFile(filePath, []byte(existing), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Update with markers
	newContent := "New managed content"
	err := UpdateFileWithMarkers(filePath, newContent, SpectrStartMarker, SpectrEndMarker)
	if err != nil {
		t.Fatalf("UpdateFileWithMarkers failed: %v", err)
	}

	// Read and verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	result := string(data)
	expected := SpectrStartMarker + "\n" + newContent + "\n" + SpectrEndMarker + "\n\n" + existing

	if result != expected {
		t.Errorf("Content mismatch.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

// Test UpdateFileWithMarkers with invalid marker state (only start marker)
func TestUpdateFileWithMarkers_InvalidMarkerState(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")

	// Create file with only start marker
	initial := SpectrStartMarker + "\nSome content"
	if err := os.WriteFile(filePath, []byte(initial), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Should fail with invalid marker state
	err := UpdateFileWithMarkers(filePath, "New content", SpectrStartMarker, SpectrEndMarker)
	if err == nil {
		t.Fatal("Expected error for invalid marker state, got nil")
	}

	if !strings.Contains(err.Error(), "invalid marker state") {
		t.Errorf("Expected 'invalid marker state' error, got: %v", err)
	}
}

// Test ClaudeCodeConfigurator
func TestClaudeCodeConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &ClaudeCodeConfigurator{}

	// Test GetName
	if configurator.GetName() != "Claude Code" {
		t.Errorf("Expected name 'Claude Code', got '%s'", configurator.GetName())
	}

	// Test IsConfigured (should be false initially)
	if configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return false for unconfigured project")
	}

	// Test Configure
	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Test IsConfigured (should be true now)
	if !configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return true after configuration")
	}

	// Verify file exists and has correct structure
	filePath := filepath.Join(tmpDir, "CLAUDE.md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, SpectrStartMarker) {
		t.Error("Generated file missing start marker")
	}
	if !strings.Contains(content, SpectrEndMarker) {
		t.Error("Generated file missing end marker")
	}
	if !strings.Contains(content, "Spectr Instructions") {
		t.Error("Generated file missing expected content")
	}
}

// Test ClineConfigurator
func TestClineConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &ClineConfigurator{}

	if configurator.GetName() != "Cline" {
		t.Errorf("Expected name 'Cline', got '%s'", configurator.GetName())
	}

	if configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return false initially")
	}

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	if !configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return true after configuration")
	}

	filePath := filepath.Join(tmpDir, "CLINE.md")
	if !FileExists(filePath) {
		t.Error("CLINE.md was not created")
	}
}

// Test CostrictConfigurator
func TestCostrictConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &CostrictConfigurator{}

	if configurator.GetName() != "CoStrict" {
		t.Errorf("Expected name 'CoStrict', got '%s'", configurator.GetName())
	}

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "COSTRICT.md")
	if !FileExists(filePath) {
		t.Error("COSTRICT.md was not created")
	}
}

// Test QoderConfigurator
func TestQoderConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &QoderConfigurator{}

	if configurator.GetName() != "Qoder" {
		t.Errorf("Expected name 'Qoder', got '%s'", configurator.GetName())
	}

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "QODER.md")
	if !FileExists(filePath) {
		t.Error("QODER.md was not created")
	}
}

// Test CodeBuddyConfigurator
func TestCodeBuddyConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &CodeBuddyConfigurator{}

	if configurator.GetName() != "CodeBuddy" {
		t.Errorf("Expected name 'CodeBuddy', got '%s'", configurator.GetName())
	}

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "CODEBUDDY.md")
	if !FileExists(filePath) {
		t.Error("CODEBUDDY.md was not created")
	}
}

// Test QwenConfigurator
func TestQwenConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &QwenConfigurator{}

	if configurator.GetName() != "Qwen Code" {
		t.Errorf("Expected name 'Qwen Code', got '%s'", configurator.GetName())
	}

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "QWEN.md")
	if !FileExists(filePath) {
		t.Error("QWEN.md was not created")
	}
}

// Test SlashCommandConfigurator - Claude
//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestClaudeSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewClaudeSlashConfigurator()

	// Test GetName
	if configurator.GetName() != "Claude Slash Commands" {
		t.Errorf("Expected name 'Claude Slash Commands', got '%s'", configurator.GetName())
	}

	// Test IsConfigured (should be false initially)
	if configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return false initially")
	}

	// Test Configure
	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Test IsConfigured (should be true now)
	if !configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return true after configuration")
	}

	// Verify all three command files exist
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".claude", "commands", "spectr", cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Command file not created: %s", filePath)

			continue
		}

		// Verify file structure
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read command file %s: %v", filePath, err)

			continue
		}

		content := string(data)

		// Should have frontmatter
		if !strings.Contains(content, "---") {
			t.Errorf("Command file %s missing frontmatter", filePath)
		}

		// Should have markers
		if !strings.Contains(content, SpectrStartMarker) {
			t.Errorf("Command file %s missing start marker", filePath)
		}
		if !strings.Contains(content, SpectrEndMarker) {
			t.Errorf("Command file %s missing end marker", filePath)
		}

		// Should have spectr instructions
		if !strings.Contains(content, "spectr") {
			t.Errorf("Command file %s missing spectr content", filePath)
		}
	}
}

// Test SlashCommandConfigurator - Kilocode
func TestKilocodeSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewKilocodeSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".kilocode", "workflows", "spectr-"+cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Workflow file not created: %s", filePath)
		}
	}
}

// Test SlashCommandConfigurator - Qoder
func TestQoderSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewQoderSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".qoder", "commands", "spectr", cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Command file not created: %s", filePath)
		}
	}
}

// Test SlashCommandConfigurator - Cursor
func TestCursorSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewCursorSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".cursor", "commands", "spectr-"+cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Command file not created: %s", filePath)
		}

		// Verify frontmatter format
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)

			continue
		}

		content := string(data)
		if !strings.Contains(content, "name: /spectr-"+cmd) {
			t.Error("File missing correct frontmatter name format")
		}
	}
}

// Test SlashCommandConfigurator - Cline
func TestClineSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewClineSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".clinerules", "spectr-"+cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Rule file not created: %s", filePath)
		}

		// Verify markdown header format
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)

			continue
		}

		content := string(data)
		if !strings.Contains(content, "# Spectr") {
			t.Error("File missing markdown header format")
		}
	}
}

// Test SlashCommandConfigurator - Windsurf
func TestWindsurfSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewWindsurfSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".windsurf", "workflows", "spectr-"+cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Workflow file not created: %s", filePath)
		}

		// Verify auto_execution_mode in frontmatter
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)

			continue
		}

		content := string(data)
		if !strings.Contains(content, "auto_execution_mode: 3") {
			t.Error("File missing auto_execution_mode in frontmatter")
		}
	}
}

// Test SlashCommandConfigurator - CoStrict
func TestCostrictSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewCostrictSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".cospec", "spectr", "commands", "spectr-"+cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Command file not created: %s", filePath)
		}

		// Verify argument-hint in frontmatter
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)

			continue
		}

		content := string(data)
		if !strings.Contains(content, "argument-hint:") {
			t.Error("File missing argument-hint in frontmatter")
		}
	}
}

// Test SlashCommandConfigurator - CodeBuddy
func TestCodeBuddySlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewCodeBuddySlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".codebuddy", "commands", "spectr", cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Command file not created: %s", filePath)
		}
	}
}

// Test SlashCommandConfigurator - Qwen
func TestQwenSlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewQwenSlashConfigurator()

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Verify file paths
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		filePath := filepath.Join(tmpDir, ".qwen", "commands", "spectr-"+cmd+".md")
		if !FileExists(filePath) {
			t.Errorf("Command file not created: %s", filePath)
		}
	}
}

// Test SlashCommandConfigurator - Update existing files
func TestSlashCommandConfigurator_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewClaudeSlashConfigurator()

	// First configuration
	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Initial configure failed: %v", err)
	}

	// Modify a file manually (simulate user changes outside markers)
	filePath := filepath.Join(tmpDir, ".claude", "commands", "spectr", "proposal.md")
	data, _ := os.ReadFile(filePath)
	modified := "# My Custom Header\n\n" + string(data) + "\n\n# My Custom Footer"
	if err := os.WriteFile(filePath, []byte(modified), 0644); err != nil {
		t.Fatalf("Failed to write modified file: %v", err)
	}

	// Second configuration (should update only content between markers)
	err = configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Second configure failed: %v", err)
	}

	// Verify custom content is preserved
	data, _ = os.ReadFile(filePath)
	content := string(data)

	if !strings.Contains(content, "# My Custom Header") {
		t.Error("Custom header was not preserved")
	}
	if !strings.Contains(content, "# My Custom Footer") {
		t.Error("Custom footer was not preserved")
	}
	if !strings.Contains(content, SpectrStartMarker) {
		t.Error("Start marker is missing")
	}
	if !strings.Contains(content, SpectrEndMarker) {
		t.Error("End marker is missing")
	}
}

// Test all remaining slash configurators exist and work
//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestAllSlashConfigurators(t *testing.T) {
	tests := []struct {
		name         string
		configurator *SlashCommandConfigurator
		basePath     string
	}{
		{"Aider", NewAiderSlashConfigurator(), ".aider/commands"},
		{"Continue", NewContinueSlashConfigurator(), ".continue/commands"},
		{"Copilot", NewCopilotSlashConfigurator(), ".github/copilot"},
		{"Mentat", NewMentatSlashConfigurator(), ".mentat/commands"},
		{"Tabnine", NewTabnineSlashConfigurator(), ".tabnine/commands"},
		{"Smol", NewSmolSlashConfigurator(), ".smol/commands"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			err := tt.configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
			if err != nil {
				t.Fatalf("Configure failed: %v", err)
			}

			if !tt.configurator.IsConfigured(tmpDir) {
				t.Error("Expected IsConfigured to return true after configuration")
			}

			// Verify at least one command file exists
			commands := []string{"proposal", "apply", "archive"}
			foundFiles := 0
			for _, cmd := range commands {
				// Check if file exists in base path
				pattern := filepath.Join(tmpDir, tt.basePath, "*"+cmd+"*")
				matches, _ := filepath.Glob(pattern)
				if len(matches) > 0 {
					foundFiles++
				}
			}

			if foundFiles == 0 {
				t.Errorf("No command files found for %s in %s", tt.name, tt.basePath)
			}
		})
	}
}

// Test isMarkerOnOwnLine helper function
func TestIsMarkerOnOwnLine(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		marker   string
		expected bool
	}{
		{
			name:     "marker on own line",
			content:  "Some text\n<!-- MARKER -->\nMore text",
			marker:   "<!-- MARKER -->",
			expected: true,
		},
		{
			name:     "marker with text before",
			content:  "Some text <!-- MARKER -->\nMore text",
			marker:   "<!-- MARKER -->",
			expected: false,
		},
		{
			name:     "marker with text after",
			content:  "Some text\n<!-- MARKER --> more text\nMore text",
			marker:   "<!-- MARKER -->",
			expected: false,
		},
		{
			name:     "marker with whitespace",
			content:  "Some text\n  <!-- MARKER -->  \nMore text",
			marker:   "<!-- MARKER -->",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := strings.Index(tt.content, tt.marker)
			if index == -1 {
				t.Fatal("Marker not found in content")
			}

			result := isMarkerOnOwnLine(tt.content, index, len(tt.marker))
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test findMarkerIndex helper function
func TestFindMarkerIndex(t *testing.T) {
	content := "Some text\n<!-- MARKER -->\nMore text\nInline <!-- MARKER --> text\n<!-- MARKER -->\nEnd"
	marker := "<!-- MARKER -->"

	// Should find first marker on own line
	index := findMarkerIndex(content, marker, 0)
	if index == -1 {
		t.Fatal("First marker not found")
	}

	// Verify it's the first occurrence on its own line
	if !isMarkerOnOwnLine(content, index, len(marker)) {
		t.Error("Found marker is not on its own line")
	}

	// Find second marker on own line
	index2 := findMarkerIndex(content, marker, index+len(marker))
	if index2 == -1 {
		t.Fatal("Second marker not found")
	}

	// Should skip the inline marker
	if index2 <= strings.Index(content, "Inline <!-- MARKER -->") {
		t.Error("Should have skipped inline marker")
	}
}

// Test AntigravityConfigurator
func TestAntigravityConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := &AntigravityConfigurator{}

	if configurator.GetName() != "Antigravity" {
		t.Errorf("Expected name 'Antigravity', got '%s'", configurator.GetName())
	}

	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "AGENTS.md")
	if !FileExists(filePath) {
		t.Error("AGENTS.md was not created")
	}
}

// Test SlashCommandConfigurator - Antigravity
//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestAntigravitySlashConfigurator(t *testing.T) {
	tmpDir := t.TempDir()
	configurator := NewAntigravitySlashConfigurator()

	// Test GetName
	if configurator.GetName() != "Antigravity Workflows" {
		t.Errorf("Expected name 'Antigravity Workflows', got '%s'", configurator.GetName())
	}

	// Test IsConfigured (should be false initially)
	if configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return false initially")
	}

	// Test Configure
	err := configurator.Configure(tmpDir, filepath.Join(tmpDir, "spectr"))
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Test IsConfigured (should be true now)
	if !configurator.IsConfigured(tmpDir) {
		t.Error("Expected IsConfigured to return true after configuration")
	}

	// Verify all three command files exist
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		relPath := configurator.config.FilePaths[cmd]
		cmdPath := filepath.Join(tmpDir, relPath)
		if !FileExists(cmdPath) {
			t.Errorf("Command file %s was not created at %s", cmd, cmdPath)
		}

		// Verify file contains Spectr markers
		content, err := os.ReadFile(cmdPath)
		if err != nil {
			t.Errorf("Failed to read %s: %v", cmdPath, err)

			continue
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, SpectrStartMarker) {
			t.Errorf("Command file %s missing start marker", cmd)
		}
		if !strings.Contains(contentStr, SpectrEndMarker) {
			t.Errorf("Command file %s missing end marker", cmd)
		}
	}

	// Test file paths
	expectedPaths := map[string]string{
		"proposal": ".agent/workflows/spectr-proposal.md",
		"apply":    ".agent/workflows/spectr-apply.md",
		"archive":  ".agent/workflows/spectr-archive.md",
	}

	for cmd, expectedPath := range expectedPaths {
		actualPath := configurator.config.FilePaths[cmd]
		if actualPath != expectedPath {
			t.Errorf("Command %s: expected path %s, got %s", cmd, expectedPath, actualPath)
		}
	}
}
