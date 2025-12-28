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
	_ *TemplateContext,
) (string, error) {
	return m.agentsContent, nil
}

func (m *mockTemplateRenderer) RenderInstructionPointer(
	_ *TemplateContext,
) (string, error) {
	return m.instructionPtrContent, nil
}

func (m *mockTemplateRenderer) RenderSlashCommand(
	command string,
	_ *TemplateContext,
) (string, error) {
	return m.slashContent[command], nil
}

func TestClaudeProvider(t *testing.T) {
	p := NewClaudeProvider()

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
	if p.ConfigFile() != "CLAUDE.md" {
		t.Errorf(
			"ConfigFile() = %s, want CLAUDE.md",
			p.ConfigFile(),
		)
	}
	if p.GetProposalCommandPath() != ".claude/commands/spectr/proposal.md" {
		t.Errorf(
			"GetProposalCommandPath() = %s, want .claude/commands/spectr/proposal.md",
			p.GetProposalCommandPath(),
		)
	}
	if p.GetApplyCommandPath() != ".claude/commands/spectr/apply.md" {
		t.Errorf(
			"GetApplyCommandPath() = %s, want .claude/commands/spectr/apply.md",
			p.GetApplyCommandPath(),
		)
	}
	if p.CommandFormat() != FormatMarkdown {
		t.Errorf(
			"CommandFormat() = %d, want FormatMarkdown",
			p.CommandFormat(),
		)
	}
	if !p.HasConfigFile() {
		t.Error(
			"HasConfigFile() = false, want true",
		)
	}
	if !p.HasSlashCommands() {
		t.Error(
			"HasSlashCommands() = false, want true",
		)
	}
}

func TestGeminiProvider(t *testing.T) {
	p := NewGeminiProvider()

	if p.ID() != "gemini" {
		t.Errorf("ID() = %s, want gemini", p.ID())
	}
	if p.Name() != "Gemini CLI" {
		t.Errorf(
			"Name() = %s, want Gemini CLI",
			p.Name(),
		)
	}
	if p.ConfigFile() != "" {
		t.Errorf(
			"ConfigFile() = %s, want empty",
			p.ConfigFile(),
		)
	}
	if p.GetProposalCommandPath() != ".gemini/commands/spectr/proposal.toml" {
		t.Errorf(
			"GetProposalCommandPath() = %s, want .gemini/commands/spectr/proposal.toml",
			p.GetProposalCommandPath(),
		)
	}
	if p.GetApplyCommandPath() != ".gemini/commands/spectr/apply.toml" {
		t.Errorf(
			"GetApplyCommandPath() = %s, want .gemini/commands/spectr/apply.toml",
			p.GetApplyCommandPath(),
		)
	}
	if p.CommandFormat() != FormatTOML {
		t.Errorf(
			"CommandFormat() = %d, want FormatTOML",
			p.CommandFormat(),
		)
	}
	if p.HasConfigFile() {
		t.Error(
			"HasConfigFile() = true, want false",
		)
	}
	if !p.HasSlashCommands() {
		t.Error(
			"HasSlashCommands() = false, want true",
		)
	}
}

func TestCursorProvider(t *testing.T) {
	p := NewCursorProvider()

	if p.ID() != "cursor" {
		t.Errorf("ID() = %s, want cursor", p.ID())
	}
	if p.ConfigFile() != "" {
		t.Errorf(
			"ConfigFile() should be empty for cursor, got %s",
			p.ConfigFile(),
		)
	}
	if !p.HasSlashCommands() {
		t.Error(
			"Cursor should have slash commands",
		)
	}
	if p.HasConfigFile() {
		t.Error(
			"Cursor should not have config file",
		)
	}
}

func TestBaseProviderConfigure(t *testing.T) {
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

	p := NewClaudeProvider()
	tm := newMockRenderer()

	err = p.Configure(
		tmpDir,
		filepath.Join(tmpDir, "spectr"),
		tm,
	)
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

func TestBaseProviderIsConfigured(t *testing.T) {
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

	p := NewClaudeProvider()

	// Should not be configured initially
	if p.IsConfigured(tmpDir) {
		t.Error(
			"Should not be configured initially",
		)
	}

	// Configure it
	tm := newMockRenderer()
	err = p.Configure(
		tmpDir,
		filepath.Join(tmpDir, "spectr"),
		tm,
	)
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

func TestBaseProviderGetFilePaths(t *testing.T) {
	p := NewClaudeProvider()
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

	p := NewGeminiProvider()
	tm := newMockRenderer()

	err = p.Configure(
		tmpDir,
		filepath.Join(tmpDir, "spectr"),
		tm,
	)
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
	p := NewCursorProvider()
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

		// All providers should have slash commands
		if !p.HasSlashCommands() {
			t.Errorf(
				"Provider %s has no slash commands",
				p.ID(),
			)
		}
		// Check that at least one command path is set
		if p.GetProposalCommandPath() == "" &&
			p.GetApplyCommandPath() == "" {
			t.Errorf(
				"Provider %s has no command paths set",
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
	p := NewCodexProvider()

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
	if p.ConfigFile() != "AGENTS.md" {
		t.Errorf(
			"ConfigFile() = %s, want AGENTS.md",
			p.ConfigFile(),
		)
	}
	if !p.HasConfigFile() {
		t.Error(
			"HasConfigFile() = false, want true",
		)
	}
	if !p.HasSlashCommands() {
		t.Error(
			"HasSlashCommands() = false, want true",
		)
	}
	if p.CommandFormat() != FormatMarkdown {
		t.Errorf(
			"CommandFormat() = %d, want FormatMarkdown",
			p.CommandFormat(),
		)
	}
}

func TestOpencodeProvider(t *testing.T) {
	p := NewOpencodeProvider()

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
	if p.ConfigFile() != "" {
		t.Errorf(
			"ConfigFile() = %s, want empty string",
			p.ConfigFile(),
		)
	}
	if p.HasConfigFile() {
		t.Error(
			"HasConfigFile() = true, want false",
		)
	}
	if !p.HasSlashCommands() {
		t.Error(
			"HasSlashCommands() = false, want true",
		)
	}
	if p.CommandFormat() != FormatMarkdown {
		t.Errorf(
			"CommandFormat() = %d, want FormatMarkdown",
			p.CommandFormat(),
		)
	}
	if p.GetProposalCommandPath() != ".opencode/command/spectr/proposal.md" {
		t.Errorf(
			"GetProposalCommandPath() = %s, want .opencode/command/spectr/proposal.md",
			p.GetProposalCommandPath(),
		)
	}
	if p.GetApplyCommandPath() != ".opencode/command/spectr/apply.md" {
		t.Errorf(
			"GetApplyCommandPath() = %s, want .opencode/command/spectr/apply.md",
			p.GetApplyCommandPath(),
		)
	}
}
