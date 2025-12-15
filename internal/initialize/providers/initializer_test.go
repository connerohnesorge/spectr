package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstructionFileInitializerID tests the ID method.
func TestInstructionFileInitializerID(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "CLAUDE.md",
			path:     "CLAUDE.md",
			expected: "instruction:CLAUDE.md",
		},
		{
			name:     "CURSOR.md",
			path:     "CURSOR.md",
			expected: "instruction:CURSOR.md",
		},
		{
			name:     "Global path",
			path:     "~/.codex/CODEX.md",
			expected: "instruction:~/.codex/CODEX.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewInstructionFileInitializer(tt.path)
			if got := init.ID(); got != tt.expected {
				t.Errorf("ID() = %s, want %s", got, tt.expected)
			}
		})
	}
}

// TestInstructionFileInitializerFilePath tests the FilePath method.
func TestInstructionFileInitializerFilePath(t *testing.T) {
	init := NewInstructionFileInitializer("CLAUDE.md")
	if got := init.FilePath(); got != "CLAUDE.md" {
		t.Errorf("FilePath() = %s, want CLAUDE.md", got)
	}
}

// TestInstructionFileInitializerConfigure tests creating and updating instruction files.
func TestInstructionFileInitializerConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewInstructionFileInitializer("CLAUDE.md")
	tm := newMockRenderer()

	// Test create
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, "CLAUDE.md")
	if !FileExists(filePath) {
		t.Error("Instruction file was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(content), SpectrStartMarker) {
		t.Error("File missing start marker")
	}
	if !strings.Contains(string(content), SpectrEndMarker) {
		t.Error("File missing end marker")
	}
	if !strings.Contains(string(content), "Spectr Instructions") {
		t.Error("File missing instruction content")
	}

	// Test update (add content before markers)
	existingContent := "# My Custom Header\n\n" + string(content)
	err = os.WriteFile(filePath, []byte(existingContent), filePerm)
	if err != nil {
		t.Fatalf("Failed to write existing content: %v", err)
	}

	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure (update) failed: %v", err)
	}

	updatedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	if !strings.Contains(string(updatedContent), "# My Custom Header") {
		t.Error("Custom content was not preserved during update")
	}
}

// TestInstructionFileInitializerIsConfigured tests the IsConfigured method.
func TestInstructionFileInitializerIsConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewInstructionFileInitializer("CLAUDE.md")

	// Should not be configured initially
	if init.IsConfigured(tmpDir) {
		t.Error("IsConfigured() = true, want false for non-existent file")
	}

	// Create the file
	tm := newMockRenderer()
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Should be configured now
	if !init.IsConfigured(tmpDir) {
		t.Error("IsConfigured() = false, want true after Configure()")
	}
}

// TestMarkdownSlashCommandInitializerID tests the ID method.
func TestMarkdownSlashCommandInitializerID(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Proposal command",
			path:     ".claude/commands/spectr/proposal.md",
			expected: "markdown-cmd:.claude/commands/spectr/proposal.md",
		},
		{
			name:     "Apply command",
			path:     ".claude/commands/spectr/apply.md",
			expected: "markdown-cmd:.claude/commands/spectr/apply.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewMarkdownSlashCommandInitializer(
				tt.path,
				"proposal",
				FrontmatterProposal,
			)
			if got := init.ID(); got != tt.expected {
				t.Errorf("ID() = %s, want %s", got, tt.expected)
			}
		})
	}
}

// TestMarkdownSlashCommandInitializerFilePath tests the FilePath method.
func TestMarkdownSlashCommandInitializerFilePath(t *testing.T) {
	init := NewMarkdownSlashCommandInitializer(
		".claude/commands/spectr/proposal.md",
		"proposal",
		FrontmatterProposal,
	)
	expected := ".claude/commands/spectr/proposal.md"
	if got := init.FilePath(); got != expected {
		t.Errorf("FilePath() = %s, want %s", got, expected)
	}
}

// TestMarkdownSlashCommandInitializerConfigure tests creating and updating command files.
func TestMarkdownSlashCommandInitializerConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewMarkdownSlashCommandInitializer(
		".claude/commands/spectr/proposal.md",
		"proposal",
		FrontmatterProposal,
	)
	tm := newMockRenderer()

	// Test create
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, ".claude/commands/spectr/proposal.md")
	if !FileExists(filePath) {
		t.Error("Command file was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Check frontmatter is present
	if !strings.HasPrefix(string(content), "---") {
		t.Error("File missing frontmatter")
	}
	if !strings.Contains(string(content), "description:") {
		t.Error("File missing description in frontmatter")
	}
	// Check markers are present
	if !strings.Contains(string(content), SpectrStartMarker) {
		t.Error("File missing start marker")
	}
	if !strings.Contains(string(content), SpectrEndMarker) {
		t.Error("File missing end marker")
	}
	// Check command content
	if !strings.Contains(string(content), "Proposal command content") {
		t.Error("File missing command content")
	}
}

// TestMarkdownSlashCommandInitializerUpdate tests updating existing command files.
func TestMarkdownSlashCommandInitializerUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewMarkdownSlashCommandInitializer(
		".claude/commands/spectr/proposal.md",
		"proposal",
		FrontmatterProposal,
	)
	tm := newMockRenderer()

	// Create initial file
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Modify the mock renderer to return different content
	tm.slashContent["proposal"] = "Updated proposal content"

	// Update the file
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure (update) failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, ".claude/commands/spectr/proposal.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(content), "Updated proposal content") {
		t.Error("File was not updated with new content")
	}
}

// TestMarkdownSlashCommandInitializerIsConfigured tests the IsConfigured method.
func TestMarkdownSlashCommandInitializerIsConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewMarkdownSlashCommandInitializer(
		".claude/commands/spectr/proposal.md",
		"proposal",
		FrontmatterProposal,
	)

	// Should not be configured initially
	if init.IsConfigured(tmpDir) {
		t.Error("IsConfigured() = true, want false for non-existent file")
	}

	// Create the file
	tm := newMockRenderer()
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Should be configured now
	if !init.IsConfigured(tmpDir) {
		t.Error("IsConfigured() = false, want true after Configure()")
	}
}

// TestMarkdownSlashCommandWithoutFrontmatter tests creating files without frontmatter.
func TestMarkdownSlashCommandWithoutFrontmatter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewMarkdownSlashCommandInitializer(
		".test/commands/proposal.md",
		"proposal",
		"", // No frontmatter
	)
	tm := newMockRenderer()

	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, ".test/commands/proposal.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Should not have frontmatter
	if strings.HasPrefix(string(content), "---") {
		t.Error("File should not have frontmatter when none specified")
	}
	// Should still have markers
	if !strings.Contains(string(content), SpectrStartMarker) {
		t.Error("File missing start marker")
	}
}

// TestTOMLSlashCommandInitializerID tests the ID method.
func TestTOMLSlashCommandInitializerID(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Proposal command",
			path:     ".gemini/commands/spectr/proposal.toml",
			expected: "toml-cmd:.gemini/commands/spectr/proposal.toml",
		},
		{
			name:     "Apply command",
			path:     ".gemini/commands/spectr/apply.toml",
			expected: "toml-cmd:.gemini/commands/spectr/apply.toml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewTOMLSlashCommandInitializer(
				tt.path,
				"proposal",
				"Test description",
			)
			if got := init.ID(); got != tt.expected {
				t.Errorf("ID() = %s, want %s", got, tt.expected)
			}
		})
	}
}

// TestTOMLSlashCommandInitializerFilePath tests the FilePath method.
func TestTOMLSlashCommandInitializerFilePath(t *testing.T) {
	init := NewTOMLSlashCommandInitializer(
		".gemini/commands/spectr/proposal.toml",
		"proposal",
		"Test description",
	)
	expected := ".gemini/commands/spectr/proposal.toml"
	if got := init.FilePath(); got != expected {
		t.Errorf("FilePath() = %s, want %s", got, expected)
	}
}

// TestTOMLSlashCommandInitializerConfigure tests creating TOML command files.
func TestTOMLSlashCommandInitializerConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewTOMLSlashCommandInitializer(
		".gemini/commands/spectr/proposal.toml",
		"proposal",
		"Scaffold a new Spectr change.",
	)
	tm := newMockRenderer()

	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, ".gemini/commands/spectr/proposal.toml")
	if !FileExists(filePath) {
		t.Error("TOML command file was not created")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Check TOML format
	if !strings.Contains(string(content), "description =") {
		t.Error("File missing description field")
	}
	if !strings.Contains(string(content), "prompt =") {
		t.Error("File missing prompt field")
	}
	if !strings.Contains(string(content), "Scaffold a new Spectr change.") {
		t.Error("File missing expected description value")
	}
	if !strings.Contains(string(content), "Proposal command content") {
		t.Error("File missing command content in prompt")
	}
}

// TestTOMLSlashCommandInitializerIsConfigured tests the IsConfigured method.
func TestTOMLSlashCommandInitializerIsConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewTOMLSlashCommandInitializer(
		".gemini/commands/spectr/proposal.toml",
		"proposal",
		"Test description",
	)

	// Should not be configured initially
	if init.IsConfigured(tmpDir) {
		t.Error("IsConfigured() = true, want false for non-existent file")
	}

	// Create the file
	tm := newMockRenderer()
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Should be configured now
	if !init.IsConfigured(tmpDir) {
		t.Error("IsConfigured() = false, want true after Configure()")
	}
}

// TestTOMLSlashCommandInitializerOverwrite tests that TOML files are overwritten.
func TestTOMLSlashCommandInitializerOverwrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewTOMLSlashCommandInitializer(
		".gemini/commands/spectr/proposal.toml",
		"proposal",
		"Initial description",
	)
	tm := newMockRenderer()

	// Create initial file
	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Create new initializer with different description
	init2 := NewTOMLSlashCommandInitializer(
		".gemini/commands/spectr/proposal.toml",
		"proposal",
		"Updated description",
	)

	// Update the file
	err = init2.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure (update) failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, ".gemini/commands/spectr/proposal.toml")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(content), "Updated description") {
		t.Error("File was not overwritten with new description")
	}
	if strings.Contains(string(content), "Initial description") {
		t.Error("Old description should not be present")
	}
}

// TestTOMLSlashCommandInitializerEscaping tests proper TOML escaping.
func TestTOMLSlashCommandInitializerEscaping(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	init := NewTOMLSlashCommandInitializer(
		".gemini/commands/spectr/proposal.toml",
		"proposal",
		"Test description",
	)

	// Create mock renderer with content that needs escaping
	tm := &mockTemplateRenderer{
		agentsContent:         "# Test",
		instructionPtrContent: "# Test",
		slashContent: map[string]string{
			"proposal": `Content with "quotes" and \backslashes\`,
		},
	}

	err = init.Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, ".gemini/commands/spectr/proposal.toml")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Verify escaping was applied
	if strings.Contains(string(content), `"quotes"`) {
		t.Error("Quotes should be escaped in TOML")
	}
	if !strings.Contains(string(content), `\"quotes\"`) {
		t.Error("Escaped quotes not found in content")
	}
}

// TestFileInitializerInterface verifies all initializers implement the interface.
func TestFileInitializerInterface(_ *testing.T) {
	// Compile-time check that all types implement FileInitializer
	var _ FileInitializer = (*InstructionFileInitializer)(nil)
	var _ FileInitializer = (*MarkdownSlashCommandInitializer)(nil)
	var _ FileInitializer = (*TOMLSlashCommandInitializer)(nil)
}

// TestGlobalPathHandling tests handling of global paths with ~ expansion.
func TestGlobalPathHandling(t *testing.T) {
	// We can't actually test ~ expansion without writing to home directory,
	// but we can verify the path logic works correctly.

	// Test instruction file with global path
	instrInit := NewInstructionFileInitializer("~/.codex/CODEX.md")
	if instrInit.FilePath() != "~/.codex/CODEX.md" {
		t.Error("FilePath() should return unexpanded path")
	}
	if instrInit.ID() != "instruction:~/.codex/CODEX.md" {
		t.Error("ID() should include unexpanded path")
	}

	// Test markdown command with global path
	mdInit := NewMarkdownSlashCommandInitializer(
		"~/.claude/commands/test.md",
		"test",
		"",
	)
	if mdInit.FilePath() != "~/.claude/commands/test.md" {
		t.Error("FilePath() should return unexpanded path")
	}

	// Test TOML command with global path
	tomlInit := NewTOMLSlashCommandInitializer(
		"~/.gemini/commands/test.toml",
		"test",
		"Test",
	)
	if tomlInit.FilePath() != "~/.gemini/commands/test.toml" {
		t.Error("FilePath() should return unexpanded path")
	}
}

// =============================================================================
// Helper Function Tests
// =============================================================================

// TestConfigureInitializersSuccess tests successful configuration of multiple initializers.
func TestConfigureInitializersSuccess(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	inits := []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
		NewMarkdownSlashCommandInitializer(
			".claude/commands/spectr/proposal.md",
			"proposal",
			FrontmatterProposal,
		),
		NewMarkdownSlashCommandInitializer(
			".claude/commands/spectr/apply.md",
			"apply",
			FrontmatterApply,
		),
	}

	tm := newMockRenderer()

	err = ConfigureInitializers(inits, tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// Verify all files were created
	expectedFiles := []string{
		filepath.Join(tmpDir, "CLAUDE.md"),
		filepath.Join(tmpDir, ".claude/commands/spectr/proposal.md"),
		filepath.Join(tmpDir, ".claude/commands/spectr/apply.md"),
	}

	for _, f := range expectedFiles {
		if !FileExists(f) {
			t.Errorf("Expected file %s was not created", f)
		}
	}
}

// TestConfigureInitializersEmptySlice tests that empty slice is handled correctly.
func TestConfigureInitializersEmptySlice(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	var inits []FileInitializer
	tm := newMockRenderer()

	err = ConfigureInitializers(inits, tmpDir, tm)
	if err != nil {
		t.Errorf("ConfigureInitializers with empty slice should not error: %v", err)
	}
}

// TestConfigureInitializersFailFast tests fail-fast behavior on error.
func TestConfigureInitializersFailFast(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a mock initializer that always fails
	failingInit := &failingInitializer{
		id:   "failing:test.md",
		path: "test.md",
	}

	inits := []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
		failingInit, // This will fail
		NewInstructionFileInitializer("CURSOR.md"), // This should NOT be reached
	}

	tm := newMockRenderer()

	err = ConfigureInitializers(inits, tmpDir, tm)
	if err == nil {
		t.Fatal("ConfigureInitializers should have failed")
	}

	// Verify error message contains the failing initializer ID
	if !strings.Contains(err.Error(), "failing:test.md") {
		t.Errorf("Error should mention failing initializer ID, got: %v", err)
	}

	// Verify first file was created (before failure)
	if !FileExists(filepath.Join(tmpDir, "CLAUDE.md")) {
		t.Error("First file (CLAUDE.md) should have been created before failure")
	}

	// Verify third file was NOT created (fail-fast stopped processing)
	if FileExists(filepath.Join(tmpDir, "CURSOR.md")) {
		t.Error("Third file (CURSOR.md) should NOT have been created due to fail-fast")
	}
}

// TestAreInitializersConfiguredAllConfigured tests when all initializers are configured.
func TestAreInitializersConfiguredAllConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	inits := []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
		NewMarkdownSlashCommandInitializer(
			".claude/commands/spectr/proposal.md",
			"proposal",
			FrontmatterProposal,
		),
	}

	tm := newMockRenderer()

	// Before configuring, should return false
	if AreInitializersConfigured(inits, tmpDir) {
		t.Error("AreInitializersConfigured should return false before configuration")
	}

	// Configure all initializers
	err = ConfigureInitializers(inits, tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// After configuring, should return true
	if !AreInitializersConfigured(inits, tmpDir) {
		t.Error("AreInitializersConfigured should return true after configuration")
	}
}

// TestAreInitializersConfiguredPartiallyConfigured tests when only some are configured.
func TestAreInitializersConfiguredPartiallyConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	inits := []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
		NewInstructionFileInitializer("CURSOR.md"),
	}

	tm := newMockRenderer()

	// Only configure the first one
	err = inits[0].Configure(tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Should return false because not ALL are configured
	if AreInitializersConfigured(inits, tmpDir) {
		t.Error("AreInitializersConfigured should return false when only partially configured")
	}
}

// TestAreInitializersConfiguredEmptySlice tests empty slice handling.
func TestAreInitializersConfiguredEmptySlice(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	var inits []FileInitializer

	// Empty slice should return true (vacuously true)
	if !AreInitializersConfigured(inits, tmpDir) {
		t.Error("AreInitializersConfigured with empty slice should return true")
	}
}

// TestGetInitializerPathsBasic tests basic path collection.
func TestGetInitializerPathsBasic(t *testing.T) {
	inits := []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
		NewMarkdownSlashCommandInitializer(
			".claude/commands/spectr/proposal.md",
			"proposal",
			FrontmatterProposal,
		),
		NewMarkdownSlashCommandInitializer(
			".claude/commands/spectr/apply.md",
			"apply",
			FrontmatterApply,
		),
	}

	paths := GetInitializerPaths(inits)

	expected := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}

	if len(paths) != len(expected) {
		t.Fatalf("Expected %d paths, got %d", len(expected), len(paths))
	}

	for i, e := range expected {
		if paths[i] != e {
			t.Errorf("Path[%d] = %s, want %s", i, paths[i], e)
		}
	}
}

// TestGetInitializerPathsDeduplication tests that duplicate paths are removed.
func TestGetInitializerPathsDeduplication(t *testing.T) {
	// Create initializers with duplicate paths
	inits := []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
		NewInstructionFileInitializer("CURSOR.md"),
		NewInstructionFileInitializer("CLAUDE.md"), // Duplicate
		NewMarkdownSlashCommandInitializer(
			".claude/commands/spectr/proposal.md",
			"proposal",
			FrontmatterProposal,
		),
		NewInstructionFileInitializer("CURSOR.md"), // Duplicate
	}

	paths := GetInitializerPaths(inits)

	// Should only have 3 unique paths
	expected := []string{
		"CLAUDE.md",
		"CURSOR.md",
		".claude/commands/spectr/proposal.md",
	}

	if len(paths) != len(expected) {
		t.Fatalf(
			"Expected %d paths after deduplication, got %d: %v",
			len(expected),
			len(paths),
			paths,
		)
	}

	for i, e := range expected {
		if paths[i] != e {
			t.Errorf("Path[%d] = %s, want %s", i, paths[i], e)
		}
	}
}

// TestGetInitializerPathsEmptySlice tests empty slice handling.
func TestGetInitializerPathsEmptySlice(t *testing.T) {
	var inits []FileInitializer

	paths := GetInitializerPaths(inits)

	if len(paths) != 0 {
		t.Errorf("Expected empty slice, got %v", paths)
	}
}

// TestGetInitializerPathsPreservesOrder tests that order is preserved.
func TestGetInitializerPathsPreservesOrder(t *testing.T) {
	inits := []FileInitializer{
		NewInstructionFileInitializer("C.md"),
		NewInstructionFileInitializer("A.md"),
		NewInstructionFileInitializer("B.md"),
		NewInstructionFileInitializer("A.md"), // Duplicate - should not affect order
	}

	paths := GetInitializerPaths(inits)

	expected := []string{"C.md", "A.md", "B.md"}

	if len(paths) != len(expected) {
		t.Fatalf("Expected %d paths, got %d", len(expected), len(paths))
	}

	for i, e := range expected {
		if paths[i] != e {
			t.Errorf("Path[%d] = %s, want %s (order not preserved)", i, paths[i], e)
		}
	}
}

// failingInitializer is a test helper that always fails on Configure.
type failingInitializer struct {
	id   string
	path string
}

func (f *failingInitializer) ID() string {
	return f.id
}

func (f *failingInitializer) FilePath() string {
	return f.path
}

func (*failingInitializer) Configure(_ string, _ TemplateRenderer) error {
	return &testError{msg: "simulated failure"}
}

func (*failingInitializer) IsConfigured(_ string) bool {
	return false
}

// testError is a simple error type for testing.
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
