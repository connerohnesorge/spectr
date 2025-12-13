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

func TestClaudeProvider(t *testing.T) {
	p := Get("claude-code")
	if p == nil {
		t.Fatal("Claude provider not registered")
	}

	if p.ID() != "claude-code" {
		t.Errorf(
			"ID() = %s, want claude-code",
			p.ID(),
		)
	}
	if p.Name() != "Claude Code" {
		t.Errorf(
			"Name() = %s, want Claude Code",
			p.Name(),
		)
	}
	if p.Priority() != PriorityClaudeCode {
		t.Errorf(
			"Priority() = %d, want %d",
			p.Priority(),
			PriorityClaudeCode,
		)
	}

	// Check initializers
	inits := p.Initializers()
	if len(inits) != 3 {
		t.Errorf(
			"Expected 3 initializers, got %d",
			len(inits),
		)
	}

	// Check file paths
	paths := p.GetFilePaths()
	expectedPaths := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}
	if len(paths) != len(expectedPaths) {
		t.Errorf(
			"Expected %d paths, got %d",
			len(expectedPaths),
			len(paths),
		)
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
			t.Errorf(
				"Expected path %s not found in GetFilePaths()",
				expected,
			)
		}
	}
}

func TestGeminiProvider(t *testing.T) {
	p := Get("gemini")
	if p == nil {
		t.Fatal("Gemini provider not registered")
	}

	if p.ID() != "gemini" {
		t.Errorf("ID() = %s, want gemini", p.ID())
	}
	if p.Name() != "Gemini CLI" {
		t.Errorf(
			"Name() = %s, want Gemini CLI",
			p.Name(),
		)
	}

	// Check initializers - Gemini has 2 TOML slash commands, no instruction file
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Errorf(
			"Expected 2 initializers, got %d",
			len(inits),
		)
	}

	// Check file paths
	paths := p.GetFilePaths()
	expectedPaths := []string{
		".gemini/commands/spectr/proposal.toml",
		".gemini/commands/spectr/apply.toml",
	}
	if len(paths) != len(expectedPaths) {
		t.Errorf(
			"Expected %d paths, got %d",
			len(expectedPaths),
			len(paths),
		)
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
			t.Errorf(
				"Expected path %s not found in GetFilePaths()",
				expected,
			)
		}
	}
}

func TestCursorProvider(t *testing.T) {
	p := Get("cursor")
	if p == nil {
		t.Fatal("Cursor provider not registered")
	}

	if p.ID() != "cursor" {
		t.Errorf("ID() = %s, want cursor", p.ID())
	}

	// Check initializers - Cursor has only 2 slash commands, no instruction file
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Errorf(
			"Expected 2 initializers, got %d",
			len(inits),
		)
	}

	// Check file paths - should have only slash command files
	paths := p.GetFilePaths()
	expectedPaths := []string{
		".cursorrules/commands/spectr/proposal.md",
		".cursorrules/commands/spectr/apply.md",
	}
	if len(paths) != len(expectedPaths) {
		t.Errorf(
			"Expected %d paths, got %d",
			len(expectedPaths),
			len(paths),
		)
	}
}

func TestProviderConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp(
		"",
		"spectr-test-*",
	)
	if err != nil {
		t.Fatalf(
			"Failed to create temp dir: %v",
			err,
		)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	p := Get("claude-code")
	if p == nil {
		t.Fatal("Claude provider not registered")
	}

	tm := newMockRenderer()

	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Check config file was created
	configPath := filepath.Join(
		tmpDir,
		"CLAUDE.md",
	)
	if !FileExists(configPath) {
		t.Error("Config file was not created")
	}

	// Check slash command files were created
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		cmdPath := filepath.Join(
			tmpDir,
			".claude/commands/spectr",
			cmd+".md",
		)
		if !FileExists(cmdPath) {
			t.Errorf(
				"Slash command file not created: %s",
				cmdPath,
			)
		}
	}
}

func TestProviderIsConfigured(t *testing.T) {
	tmpDir, err := os.MkdirTemp(
		"",
		"spectr-test-*",
	)
	if err != nil {
		t.Fatalf(
			"Failed to create temp dir: %v",
			err,
		)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	p := Get("claude-code")
	if p == nil {
		t.Fatal("Claude provider not registered")
	}

	// Should not be configured initially
	if p.IsConfigured(tmpDir) {
		t.Error(
			"Should not be configured initially",
		)
	}

	// Configure it
	tm := newMockRenderer()
	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Should be configured now
	if !p.IsConfigured(tmpDir) {
		t.Error(
			"Should be configured after Configure()",
		)
	}
}

func TestProviderGetFilePaths(t *testing.T) {
	p := Get("claude-code")
	if p == nil {
		t.Fatal("Claude provider not registered")
	}

	paths := p.GetFilePaths()

	// Should have config file + 2 slash command files
	expectedPaths := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf(
			"Expected %d paths, got %d",
			len(expectedPaths),
			len(paths),
		)
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
			t.Errorf(
				"Expected path %s not found in GetFilePaths()",
				expected,
			)
		}
	}
}

func TestGeminiProviderConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp(
		"",
		"spectr-test-*",
	)
	if err != nil {
		t.Fatalf(
			"Failed to create temp dir: %v",
			err,
		)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	p := Get("gemini")
	if p == nil {
		t.Fatal("Gemini provider not registered")
	}

	tm := newMockRenderer()

	err = ConfigureInitializers(p.Initializers(), tmpDir, tm)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	// Check TOML files were created
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		cmdPath := filepath.Join(
			tmpDir,
			".gemini/commands/spectr",
			cmd+".toml",
		)
		if !FileExists(cmdPath) {
			t.Errorf(
				"TOML command file not created: %s",
				cmdPath,
			)
		}

		// Verify content is TOML format
		content, err := os.ReadFile(cmdPath)
		if err != nil {
			t.Errorf(
				"Failed to read %s: %v",
				cmdPath,
				err,
			)

			continue
		}
		if !strings.Contains(
			string(content),
			"description =",
		) {
			t.Errorf(
				"File %s doesn't look like TOML",
				cmdPath,
			)
		}
		if !strings.Contains(
			string(content),
			"prompt =",
		) {
			t.Errorf(
				"File %s missing prompt field",
				cmdPath,
			)
		}
	}
}

func TestSlashOnlyProviderGetFilePaths(
	t *testing.T,
) {
	p := Get("cursor")
	if p == nil {
		t.Fatal("Cursor provider not registered")
	}

	paths := p.GetFilePaths()

	// Should have only slash command files (no config file)
	expectedPaths := []string{
		".cursorrules/commands/spectr/proposal.md",
		".cursorrules/commands/spectr/apply.md",
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf(
			"Expected %d paths, got %d",
			len(expectedPaths),
			len(paths),
		)
	}
}

func TestAllProvidersHaveRequiredFields(
	t *testing.T,
) {
	allProviders := All()

	for _, p := range allProviders {
		if p.ID() == "" {
			t.Error(
				"Found provider with empty ID",
			)
		}
		if p.Name() == "" {
			t.Errorf(
				"Provider %s has empty Name",
				p.ID(),
			)
		}
		if p.Priority() < 1 {
			t.Errorf(
				"Provider %s has invalid priority: %d",
				p.ID(),
				p.Priority(),
			)
		}

		// All providers should have initializers
		inits := p.Initializers()
		if len(inits) == 0 {
			t.Errorf(
				"Provider %s has no initializers",
				p.ID(),
			)
		}

		// All providers should return file paths
		paths := p.GetFilePaths()
		if len(paths) == 0 {
			t.Errorf(
				"Provider %s has no file paths",
				p.ID(),
			)
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

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf(
			"Failed to get home directory: %v",
			err,
		)
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
			name:  "Path starting with tilde slash",
			input: "~/.config/test",
			expected: filepath.Join(
				homeDir,
				".config/test",
			),
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
				t.Errorf(
					"expandPath(%q) = %q, want %q",
					tt.input,
					result,
					tt.expected,
				)
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
				t.Errorf(
					"isGlobalPath(%q) = %v, want %v",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestCodexProvider(t *testing.T) {
	p := Get("codex")
	if p == nil {
		t.Fatal("Codex provider not registered")
	}

	if p.ID() != "codex" {
		t.Errorf("ID() = %s, want codex", p.ID())
	}
	if p.Name() != "Codex CLI" {
		t.Errorf(
			"Name() = %s, want Codex CLI",
			p.Name(),
		)
	}
	if p.Priority() != PriorityCodex {
		t.Errorf(
			"Priority() = %d, want %d",
			p.Priority(),
			PriorityCodex,
		)
	}

	// Check initializers - Codex has instruction file + 2 slash commands
	inits := p.Initializers()
	if len(inits) != 3 {
		t.Errorf(
			"Expected 3 initializers, got %d",
			len(inits),
		)
	}

	// Check file paths include AGENTS.md and global paths
	paths := p.GetFilePaths()
	if len(paths) != 3 {
		t.Errorf(
			"Expected 3 paths, got %d",
			len(paths),
		)
	}

	// Verify AGENTS.md is included
	foundAgents := false
	for _, path := range paths {
		if path == "AGENTS.md" {
			foundAgents = true

			break
		}
	}
	if !foundAgents {
		t.Error("Expected AGENTS.md in file paths")
	}
}

func TestOpencodeProvider(t *testing.T) {
	p := Get("opencode")
	if p == nil {
		t.Fatal("OpenCode provider not registered")
	}

	if p.ID() != "opencode" {
		t.Errorf(
			"ID() = %s, want opencode",
			p.ID(),
		)
	}
	if p.Name() != "OpenCode" {
		t.Errorf(
			"Name() = %s, want OpenCode",
			p.Name(),
		)
	}
	if p.Priority() != PriorityOpencode {
		t.Errorf(
			"Priority() = %d, want %d",
			p.Priority(),
			PriorityOpencode,
		)
	}

	// Check initializers - OpenCode has only 2 slash commands, no instruction file
	inits := p.Initializers()
	if len(inits) != 2 {
		t.Errorf(
			"Expected 2 initializers, got %d",
			len(inits),
		)
	}

	// Check file paths
	paths := p.GetFilePaths()
	expectedPaths := []string{
		".opencode/command/spectr/proposal.md",
		".opencode/command/spectr/apply.md",
	}
	if len(paths) != len(expectedPaths) {
		t.Errorf(
			"Expected %d paths, got %d",
			len(expectedPaths),
			len(paths),
		)
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
			t.Errorf(
				"Expected path %s not found in GetFilePaths()",
				expected,
			)
		}
	}
}
