package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockTemplateRenderer implements TemplateRenderer for testing
type mockTemplateRenderer struct {
	agentsContent         string
	instructionPtrContent string
	slashContent          map[string]string
}

func newMockRenderer() *mockTemplateRenderer {
	return &mockTemplateRenderer{
		agentsContent:         "# Test AGENTS content",
		instructionPtrContent: "# Spectr Instructions\nRead spectr/AGENTS.md",
		slashContent: map[string]string{
			"proposal": "Proposal command content",
			"apply":    "Apply command content",
		},
	}
}

func (m *mockTemplateRenderer) RenderAgents(
	_ TemplateContext,
) (string, error) {
	return m.agentsContent, nil
}

func (m *mockTemplateRenderer) RenderInstructionPointer(
	_ TemplateContext,
) (string, error) {
	return m.instructionPtrContent, nil
}

func (m *mockTemplateRenderer) RenderSlashCommand(
	command string,
	_ TemplateContext,
) (string, error) {
	return m.slashContent[command], nil
}

// =============================================================================
// Provider Interface Tests (6-method interface)
// =============================================================================

func TestClaudeProviderInterface(t *testing.T) {
	p := Get("claude-code")
	if p == nil {
		t.Fatal("claude-code provider not found in registry")
	}

	// Test ID()
	if p.ID() != "claude-code" {
		t.Errorf("ID() = %s, want claude-code", p.ID())
	}

	// Test Name()
	if p.Name() != "Claude Code" {
		t.Errorf("Name() = %s, want Claude Code", p.Name())
	}

	// Test Priority()
	if p.Priority() != PriorityClaudeCode {
		t.Errorf("Priority() = %d, want %d", p.Priority(), PriorityClaudeCode)
	}

	// Test Initializers()
	inits := p.Initializers()
	if len(inits) != 3 {
		t.Errorf("Initializers() returned %d items, want 3", len(inits))
	}

	// Verify initializer types and paths
	expectedPaths := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}
	for i, expectedPath := range expectedPaths {
		if i >= len(inits) {
			break
		}
		if inits[i].FilePath() != expectedPath {
			t.Errorf(
				"Initializer[%d].FilePath() = %s, want %s",
				i,
				inits[i].FilePath(),
				expectedPath,
			)
		}
	}

	// Test GetFilePaths()
	paths := p.GetFilePaths()
	if len(paths) != len(expectedPaths) {
		t.Errorf("GetFilePaths() returned %d paths, want %d", len(paths), len(expectedPaths))
	}
	for i, expectedPath := range expectedPaths {
		if i >= len(paths) {
			break
		}
		if paths[i] != expectedPath {
			t.Errorf("GetFilePaths()[%d] = %s, want %s", i, paths[i], expectedPath)
		}
	}
}

func TestGeminiProviderInterface(t *testing.T) {
	p := Get("gemini")
	if p == nil {
		t.Fatal("gemini provider not found in registry")
	}

	// Test ID()
	if p.ID() != "gemini" {
		t.Errorf("ID() = %s, want gemini", p.ID())
	}

	// Test Name()
	if p.Name() != "Gemini CLI" {
		t.Errorf("Name() = %s, want Gemini CLI", p.Name())
	}

	// Test Priority()
	if p.Priority() != PriorityGemini {
		t.Errorf("Priority() = %d, want %d", p.Priority(), PriorityGemini)
	}

	// Test Initializers() - Gemini has no instruction file, only TOML commands
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Errorf("Initializers() returned %d items, want 2 (TOML commands only)", len(inits))
	}

	// Verify initializer types and paths
	expectedPaths := []string{
		".gemini/commands/spectr/proposal.toml",
		".gemini/commands/spectr/apply.toml",
	}
	for i, expectedPath := range expectedPaths {
		if i >= len(inits) {
			break
		}
		if inits[i].FilePath() != expectedPath {
			t.Errorf(
				"Initializer[%d].FilePath() = %s, want %s",
				i,
				inits[i].FilePath(),
				expectedPath,
			)
		}
	}

	// Test GetFilePaths()
	paths := p.GetFilePaths()
	if len(paths) != len(expectedPaths) {
		t.Errorf("GetFilePaths() returned %d paths, want %d", len(paths), len(expectedPaths))
	}
}

func TestCursorProviderInterface(t *testing.T) {
	p := Get("cursor")
	if p == nil {
		t.Fatal("cursor provider not found in registry")
	}

	// Test ID()
	if p.ID() != "cursor" {
		t.Errorf("ID() = %s, want cursor", p.ID())
	}

	// Test Initializers() - Cursor has no instruction file, only markdown commands
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Errorf("Initializers() returned %d items, want 2 (commands only)", len(inits))
	}

	// Verify paths
	expectedPaths := []string{
		".cursorrules/commands/spectr/proposal.md",
		".cursorrules/commands/spectr/apply.md",
	}
	paths := p.GetFilePaths()
	if len(paths) != len(expectedPaths) {
		t.Errorf("GetFilePaths() returned %d paths, want %d", len(paths), len(expectedPaths))
	}
}

func TestCodexProviderInterface(t *testing.T) {
	p := Get("codex")
	if p == nil {
		t.Fatal("codex provider not found in registry")
	}

	// Test ID()
	if p.ID() != "codex" {
		t.Errorf("ID() = %s, want codex", p.ID())
	}

	// Test Name()
	if p.Name() != "Codex CLI" {
		t.Errorf("Name() = %s, want Codex CLI", p.Name())
	}

	// Test Priority()
	if p.Priority() != PriorityCodex {
		t.Errorf("Priority() = %d, want %d", p.Priority(), PriorityCodex)
	}

	// Test Initializers() - Codex has instruction file + global command paths
	inits := p.Initializers()
	if len(inits) != 3 {
		t.Errorf("Initializers() returned %d items, want 3", len(inits))
	}

	// Verify first initializer is instruction file
	if inits[0].FilePath() != "AGENTS.md" {
		t.Errorf("First initializer FilePath() = %s, want AGENTS.md", inits[0].FilePath())
	}

	// Verify command paths are global (start with ~/)
	for i := 1; i < len(inits); i++ {
		path := inits[i].FilePath()
		if !strings.HasPrefix(path, "~/") {
			t.Errorf("Codex command path should be global (start with ~/), got %s", path)
		}
	}
}

func TestOpencodeProviderInterface(t *testing.T) {
	p := Get("opencode")
	if p == nil {
		t.Fatal("opencode provider not found in registry")
	}

	// Test ID()
	if p.ID() != "opencode" {
		t.Errorf("ID() = %s, want opencode", p.ID())
	}

	// Test Name()
	if p.Name() != "OpenCode" {
		t.Errorf("Name() = %s, want OpenCode", p.Name())
	}

	// Test Priority()
	if p.Priority() != PriorityOpencode {
		t.Errorf("Priority() = %d, want %d", p.Priority(), PriorityOpencode)
	}

	// Test Initializers() - OpenCode has no instruction file
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Errorf("Initializers() returned %d items, want 2", len(inits))
	}

	// Verify paths
	expectedPaths := []string{
		".opencode/command/spectr/proposal.md",
		".opencode/command/spectr/apply.md",
	}
	for i, expectedPath := range expectedPaths {
		if i >= len(inits) {
			break
		}
		if inits[i].FilePath() != expectedPath {
			t.Errorf(
				"Initializer[%d].FilePath() = %s, want %s",
				i,
				inits[i].FilePath(),
				expectedPath,
			)
		}
	}
}

// =============================================================================
// IsConfigured Tests
// =============================================================================

func TestClaudeProviderIsConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	p := Get("claude-code")
	if p == nil {
		t.Fatal("claude-code provider not found")
	}

	// Should not be configured initially
	if p.IsConfigured(tmpDir) {
		t.Error("Should not be configured initially")
	}

	// Configure using ConfigureInitializers helper
	tm := newMockRenderer()
	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// Should be configured now
	if !p.IsConfigured(tmpDir) {
		t.Error("Should be configured after ConfigureInitializers()")
	}
}

func TestGeminiProviderIsConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	p := Get("gemini")
	if p == nil {
		t.Fatal("gemini provider not found")
	}

	// Should not be configured initially
	if p.IsConfigured(tmpDir) {
		t.Error("Should not be configured initially")
	}

	// Configure using ConfigureInitializers helper
	tm := newMockRenderer()
	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// Should be configured now
	if !p.IsConfigured(tmpDir) {
		t.Error("Should be configured after ConfigureInitializers()")
	}
}

// =============================================================================
// GetFilePaths Tests
// =============================================================================

func TestClaudeProviderGetFilePaths(t *testing.T) {
	p := Get("claude-code")
	if p == nil {
		t.Fatal("claude-code provider not found")
	}

	paths := p.GetFilePaths()

	// Should have instruction file + 2 slash command files
	expectedPaths := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
	}

	for _, expected := range expectedPaths {
		found := false
		for _, path := range paths {
			if path == expected {
				found = true

				break
			}
		}
		if !found {
			t.Errorf("Expected path %s not found in GetFilePaths()", expected)
		}
	}
}

func TestCursorProviderGetFilePaths(t *testing.T) {
	p := Get("cursor")
	if p == nil {
		t.Fatal("cursor provider not found")
	}

	paths := p.GetFilePaths()

	// Should have only slash command files (no instruction file)
	expectedPaths := []string{
		".cursorrules/commands/spectr/proposal.md",
		".cursorrules/commands/spectr/apply.md",
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
	}
}

// =============================================================================
// All Providers Validation Tests
// =============================================================================

func TestAllProvidersHaveRequiredFields(t *testing.T) {
	allProviders := All()

	for _, p := range allProviders {
		if p.ID() == "" {
			t.Error("Found provider with empty ID")
		}
		if p.Name() == "" {
			t.Errorf("Provider %s has empty Name", p.ID())
		}
		if p.Priority() < 1 {
			t.Errorf("Provider %s has invalid priority: %d", p.ID(), p.Priority())
		}

		// All providers should have at least one initializer
		inits := p.Initializers()
		if len(inits) == 0 {
			t.Errorf("Provider %s has no initializers", p.ID())
		}

		// GetFilePaths should return at least one path
		paths := p.GetFilePaths()
		if len(paths) == 0 {
			t.Errorf("Provider %s has no file paths", p.ID())
		}

		// Verify initializers count matches paths count (after deduplication)
		if len(paths) > len(inits) {
			t.Errorf("Provider %s: GetFilePaths() returned more paths than initializers", p.ID())
		}
	}
}

func TestAllProvidersHaveInitializers(t *testing.T) {
	allProviders := All()

	for _, p := range allProviders {
		inits := p.Initializers()

		if len(inits) == 0 {
			t.Errorf("Provider %s has no initializers", p.ID())

			continue
		}

		// Every initializer should have valid ID and FilePath
		for i, init := range inits {
			if init.ID() == "" {
				t.Errorf("Provider %s: Initializer[%d] has empty ID", p.ID(), i)
			}
			if init.FilePath() == "" {
				t.Errorf("Provider %s: Initializer[%d] has empty FilePath", p.ID(), i)
			}
		}
	}
}

func TestPrioritiesAreUnique(t *testing.T) {
	allProviders := All()
	priorities := make(map[int]string)

	for _, p := range allProviders {
		if existingID, exists := priorities[p.Priority()]; exists {
			t.Errorf(
				"Duplicate priority %d for providers %s and %s",
				p.Priority(),
				existingID,
				p.ID(),
			)
		}
		priorities[p.Priority()] = p.ID()
	}
}

// =============================================================================
// Helper Function Tests (expandPath, isGlobalPath)
// =============================================================================

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Path not starting with tilde",
			input:    ".config/test",
			expected: ".config/test",
		},
		{
			name:     "Path starting with tilde slash",
			input:    "~/.config/test",
			expected: filepath.Join(homeDir, ".config/test"),
		},
		{
			name:     "Absolute path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "Tilde only without slash",
			input:    "~",
			expected: "~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsGlobalPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Path starting with tilde slash",
			input:    "~/.config/test",
			expected: true,
		},
		{
			name:     "Path starting with absolute slash",
			input:    "/absolute/path",
			expected: true,
		},
		{
			name:     "Relative path with dot",
			input:    ".foo/bar",
			expected: false,
		},
		{
			name:     "Simple relative path",
			input:    "foo/bar",
			expected: false,
		},
		{
			name:     "Current directory dot",
			input:    "./foo",
			expected: false,
		},
		{
			name:     "Parent directory",
			input:    "../foo",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGlobalPath(tt.input)
			if result != tt.expected {
				t.Errorf("isGlobalPath(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Integration Tests - Full Provider Configuration Flow
// =============================================================================

func TestClaudeProviderFullConfigurationFlow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// 1. Get the provider
	p := Get("claude-code")
	if p == nil {
		t.Fatal("claude-code provider not found")
	}

	// 2. Verify not configured initially
	if p.IsConfigured(tmpDir) {
		t.Error("Provider should not be configured initially")
	}

	// 3. Get initializers
	inits := p.Initializers()
	if len(inits) != 3 {
		t.Fatalf("Expected 3 initializers, got %d", len(inits))
	}

	// 4. Configure using ConfigureInitializers helper
	tm := newMockRenderer()
	err = ConfigureInitializers(inits, tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// 5. Verify IsConfigured returns true
	if !p.IsConfigured(tmpDir) {
		t.Error("Provider should be configured after ConfigureInitializers()")
	}

	// 6. Verify GetFilePaths returns correct paths
	paths := p.GetFilePaths()
	expectedPaths := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}
	if len(paths) != len(expectedPaths) {
		t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
	}

	// 7. Verify files were actually created
	for _, path := range expectedPaths {
		fullPath := filepath.Join(tmpDir, path)
		if !FileExists(fullPath) {
			t.Errorf("File %s was not created", path)
		}
	}

	// 8. Verify instruction file content
	configPath := filepath.Join(tmpDir, "CLAUDE.md")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	if !strings.Contains(string(content), SpectrStartMarker) {
		t.Error("Config file missing start marker")
	}
	if !strings.Contains(string(content), SpectrEndMarker) {
		t.Error("Config file missing end marker")
	}

	// 9. Verify slash command files content
	for _, cmd := range []string{"proposal", "apply"} {
		cmdPath := filepath.Join(tmpDir, ".claude/commands/spectr", cmd+".md")
		cmdContent, err := os.ReadFile(cmdPath)
		if err != nil {
			t.Fatalf("Failed to read %s command file: %v", cmd, err)
		}
		if !strings.Contains(string(cmdContent), "---") {
			t.Errorf("Command file %s missing frontmatter", cmd)
		}
		if !strings.Contains(string(cmdContent), SpectrStartMarker) {
			t.Errorf("Command file %s missing start marker", cmd)
		}
	}
}

func TestGeminiProviderFullConfigurationFlow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// 1. Get the provider
	p := Get("gemini")
	if p == nil {
		t.Fatal("gemini provider not found")
	}

	// 2. Verify not configured initially
	if p.IsConfigured(tmpDir) {
		t.Error("Provider should not be configured initially")
	}

	// 3. Get initializers
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Fatalf("Expected 2 initializers, got %d", len(inits))
	}

	// 4. Configure using ConfigureInitializers helper
	tm := newMockRenderer()
	err = ConfigureInitializers(inits, tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// 5. Verify IsConfigured returns true
	if !p.IsConfigured(tmpDir) {
		t.Error("Provider should be configured after ConfigureInitializers()")
	}

	// 6. Verify TOML files were created with correct format
	for _, cmd := range []string{"proposal", "apply"} {
		cmdPath := filepath.Join(tmpDir, ".gemini/commands/spectr", cmd+".toml")
		if !FileExists(cmdPath) {
			t.Errorf("TOML command file %s was not created", cmdPath)

			continue
		}

		cmdContent, err := os.ReadFile(cmdPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", cmdPath, err)
		}

		// Verify TOML format
		if !strings.Contains(string(cmdContent), "description =") {
			t.Errorf("File %s missing description field", cmdPath)
		}
		if !strings.Contains(string(cmdContent), "prompt =") {
			t.Errorf("File %s missing prompt field", cmdPath)
		}
	}
}

func TestCursorProviderFullConfigurationFlow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// 1. Get the provider
	p := Get("cursor")
	if p == nil {
		t.Fatal("cursor provider not found")
	}

	// 2. Verify not configured initially
	if p.IsConfigured(tmpDir) {
		t.Error("Provider should not be configured initially")
	}

	// 3. Configure using ConfigureInitializers helper
	tm := newMockRenderer()
	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("ConfigureInitializers failed: %v", err)
	}

	// 4. Verify IsConfigured returns true
	if !p.IsConfigured(tmpDir) {
		t.Error("Provider should be configured after ConfigureInitializers()")
	}

	// 5. Verify only slash command files were created (no instruction file)
	paths := p.GetFilePaths()
	if len(paths) != 2 {
		t.Errorf("Expected 2 paths (commands only), got %d", len(paths))
	}

	// Verify all paths are command paths
	for _, path := range paths {
		if !strings.Contains(path, "commands") {
			t.Errorf("Path %s should be a command path", path)
		}
	}
}

// =============================================================================
// Configuration Update Tests
// =============================================================================

func TestProviderConfigurationUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	p := Get("claude-code")
	if p == nil {
		t.Fatal("claude-code provider not found")
	}

	tm := newMockRenderer()

	// First configuration
	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("First ConfigureInitializers failed: %v", err)
	}

	// Add custom content to config file before markers
	configPath := filepath.Join(tmpDir, "CLAUDE.md")
	existingContent, _ := os.ReadFile(configPath)
	newContent := "# My Custom Header\n\n" + string(existingContent)
	err = os.WriteFile(configPath, []byte(newContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to write custom content: %v", err)
	}

	// Second configuration (update)
	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("Second ConfigureInitializers failed: %v", err)
	}

	// Verify custom content was preserved
	updatedContent, _ := os.ReadFile(configPath)
	if !strings.Contains(string(updatedContent), "# My Custom Header") {
		t.Error("Custom content was not preserved during update")
	}

	// Verify still configured
	if !p.IsConfigured(tmpDir) {
		t.Error("Provider should still be configured after update")
	}
}

// =============================================================================
// Provider Interface Compliance Test
// =============================================================================

func TestProviderInterfaceCompliance(t *testing.T) {
	// This is a compile-time check that all providers implement the interface
	// The test just verifies all registered providers can be used polymorphically
	allProviders := All()

	if len(allProviders) == 0 {
		t.Fatal("No providers registered")
	}

	for _, p := range allProviders {
		// Verify each provider can be used through the interface
		var _ Provider = p //nolint:staticcheck // Compile-time interface check

		// Verify all methods can be called
		_ = p.ID()
		_ = p.Name()
		_ = p.Priority()
		_ = p.Initializers()
		_ = p.GetFilePaths()

		// IsConfigured requires a path, test with empty string (should not panic)
		_ = p.IsConfigured("")
	}
}
