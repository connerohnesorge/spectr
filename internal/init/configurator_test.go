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

// Test ClaudeMemoryFileProvider
func TestClaudeMemoryFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &ClaudeMemoryFileProvider{}

	// Test IsMemoryFileConfigured (should be false initially)
	if provider.IsMemoryFileConfigured(tmpDir) {
		t.Error("Expected IsMemoryFileConfigured to return false for unconfigured project")
	}

	// Test ConfigureMemoryFile
	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	// Test IsMemoryFileConfigured (should be true now)
	if !provider.IsMemoryFileConfigured(tmpDir) {
		t.Error("Expected IsMemoryFileConfigured to return true after configuration")
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

// Test ClineMemoryFileProvider
func TestClineMemoryFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &ClineMemoryFileProvider{}

	if provider.IsMemoryFileConfigured(tmpDir) {
		t.Error("Expected IsMemoryFileConfigured to return false initially")
	}

	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	if !provider.IsMemoryFileConfigured(tmpDir) {
		t.Error("Expected IsMemoryFileConfigured to return true after configuration")
	}

	filePath := filepath.Join(tmpDir, "CLINE.md")
	if !FileExists(filePath) {
		t.Error("CLINE.md was not created")
	}
}

// Test CostrictMemoryFileProvider
func TestCostrictMemoryFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &CostrictMemoryFileProvider{}

	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "COSTRICT.md")
	if !FileExists(filePath) {
		t.Error("COSTRICT.md was not created")
	}
}

// Test QoderMemoryFileProvider
func TestQoderMemoryFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &QoderMemoryFileProvider{}

	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "QODER.md")
	if !FileExists(filePath) {
		t.Error("QODER.md was not created")
	}
}

// Test CodeBuddyMemoryFileProvider
func TestCodeBuddyMemoryFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &CodeBuddyMemoryFileProvider{}

	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "CODEBUDDY.md")
	if !FileExists(filePath) {
		t.Error("CODEBUDDY.md was not created")
	}
}

// Test QwenMemoryFileProvider
func TestQwenMemoryFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &QwenMemoryFileProvider{}

	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "QWEN.md")
	if !FileExists(filePath) {
		t.Error("QWEN.md was not created")
	}
}

// Test SlashCommandProvider - Claude
//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestClaudeSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewClaudeSlashCommandProvider()

	// Test AreSlashCommandsConfigured (should be false initially)
	if provider.AreSlashCommandsConfigured(tmpDir) {
		t.Error("Expected AreSlashCommandsConfigured to return false initially")
	}

	// Test ConfigureSlashCommands
	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
	}

	// Test AreSlashCommandsConfigured (should be true now)
	if !provider.AreSlashCommandsConfigured(tmpDir) {
		t.Error("Expected AreSlashCommandsConfigured to return true after configuration")
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

// Test SlashCommandProvider - Kilocode
func TestKilocodeSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewKilocodeSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - Qoder
func TestQoderSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewQoderSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - Cursor
func TestCursorSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewCursorSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - Cline
func TestClineSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewClineSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - Windsurf
func TestWindsurfSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewWindsurfSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - CoStrict
func TestCostrictSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewCostrictSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - CodeBuddy
func TestCodeBuddySlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewCodeBuddySlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - Qwen
func TestQwenSlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewQwenSlashCommandProvider()

	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
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

// Test SlashCommandProvider - Update existing files
func TestSlashCommandProvider_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewClaudeSlashCommandProvider()

	// First configuration
	err := provider.ConfigureSlashCommands(tmpDir)
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
	err = provider.ConfigureSlashCommands(tmpDir)
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

// Test all remaining slash command providers exist and work
//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestAllSlashCommandProviders(t *testing.T) {
	tests := []struct {
		name     string
		provider SlashCommandProvider
		basePath string
	}{
		{"Aider", NewAiderSlashCommandProvider(), ".aider/commands"},
		{"Continue", NewContinueSlashCommandProvider(), ".continue/commands"},
		{"Copilot", NewCopilotSlashCommandProvider(), ".github/copilot"},
		{"Mentat", NewMentatSlashCommandProvider(), ".mentat/commands"},
		{"Tabnine", NewTabnineSlashCommandProvider(), ".tabnine/commands"},
		{"Smol", NewSmolSlashCommandProvider(), ".smol/commands"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			err := tt.provider.ConfigureSlashCommands(tmpDir)
			if err != nil {
				t.Fatalf("ConfigureSlashCommands failed: %v", err)
			}

			if !tt.provider.AreSlashCommandsConfigured(tmpDir) {
				t.Error("Expected AreSlashCommandsConfigured to return true after configuration")
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

// Test AgentsFileProvider (for Antigravity)
func TestAgentsFileProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := &AgentsFileProvider{}

	err := provider.ConfigureMemoryFile(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureMemoryFile failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "AGENTS.md")
	if !FileExists(filePath) {
		t.Error("AGENTS.md was not created")
	}
}

// Test SlashCommandProvider - Antigravity
//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestAntigravitySlashCommandProvider(t *testing.T) {
	tmpDir := t.TempDir()
	provider := NewAntigravitySlashCommandProvider()

	// Test AreSlashCommandsConfigured (should be false initially)
	if provider.AreSlashCommandsConfigured(tmpDir) {
		t.Error("Expected AreSlashCommandsConfigured to return false initially")
	}

	// Test ConfigureSlashCommands
	err := provider.ConfigureSlashCommands(tmpDir)
	if err != nil {
		t.Fatalf("ConfigureSlashCommands failed: %v", err)
	}

	// Test AreSlashCommandsConfigured (should be true now)
	if !provider.AreSlashCommandsConfigured(tmpDir) {
		t.Error("Expected AreSlashCommandsConfigured to return true after configuration")
	}

	// Verify all three command files exist
	commands := []string{"proposal", "apply", "archive"}
	expectedPaths := map[string]string{
		"proposal": ".agent/workflows/spectr-proposal.md",
		"apply":    ".agent/workflows/spectr-apply.md",
		"archive":  ".agent/workflows/spectr-archive.md",
	}

	for _, cmd := range commands {
		relPath := expectedPaths[cmd]
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
}
