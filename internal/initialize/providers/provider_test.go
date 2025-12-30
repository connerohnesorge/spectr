//nolint:revive // Test file uses type assertions without checking - types verified by earlier checks
package providers

import (
	"context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
)

const (
	testClaudeCommandsDir = ".claude/commands/spectr"
	testGeminiCommandsDir = ".gemini/commands/spectr"
	testAgentWorkflowsDir = ".agent/workflows"
	testSpectrPrefix      = "spectr-"
	testCodexPromptsDir   = ".codex/prompts"
)

// mockTemplateManager implements TemplateManager for testing
type mockTemplateManager struct{}

func (*mockTemplateManager) InstructionPointer() domain.TemplateRef {
	return domain.TemplateRef{Name: "instruction-pointer.md.tmpl"}
}

func (*mockTemplateManager) Agents() domain.TemplateRef {
	return domain.TemplateRef{Name: "AGENTS.md.tmpl"}
}

func (*mockTemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	return domain.TemplateRef{Name: fmt.Sprintf("slash-%s.md.tmpl", cmd.String())}
}

func (*mockTemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	return domain.TemplateRef{Name: fmt.Sprintf("slash-%s.toml.tmpl", cmd.String())}
}

func (*mockTemplateManager) SkillFS(skillName string) (fs.FS, error) {
	return nil, fmt.Errorf("skill %s not found", skillName)
}

// Test each provider returns expected initializers

func TestClaudeProvider_Initializers(t *testing.T) {
	p := &ClaudeProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Claude should return 5 initializers: Directory (commands), Directory (skills), ConfigFile, SlashCommands, AgentSkills
	if len(inits) != 5 {
		t.Fatalf("ClaudeProvider.Initializers() returned %d initializers, want 5", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[0] is %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*DirectoryInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[1] is %T, want *DirectoryInitializer", inits[1])
	}
	if _, ok := inits[2].(*ConfigFileInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[2] is %T, want *ConfigFileInitializer", inits[2])
	}
	if _, ok := inits[3].(*SlashCommandsInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[3] is %T, want *SlashCommandsInitializer", inits[3])
	}
	if _, ok := inits[4].(*AgentSkillsInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[4] is %T, want *AgentSkillsInitializer", inits[4])
	}

	// Check DirectoryInitializer paths
	dirInit := inits[0].(*DirectoryInitializer) //nolint:revive // test code, type checked above
	if len(dirInit.paths) != 1 || dirInit.paths[0] != testClaudeCommandsDir {
		t.Errorf(
			"ClaudeProvider DirectoryInitializer paths = %v, want [\".claude/commands/spectr\"]",
			dirInit.paths,
		)
	}

	// Check second DirectoryInitializer for skills
	skillsDirInit := inits[1].(*DirectoryInitializer) //nolint:revive // test code, type checked above
	if len(skillsDirInit.paths) != 1 || skillsDirInit.paths[0] != ".claude/skills" {
		t.Errorf(
			"ClaudeProvider DirectoryInitializer[1] paths = %v, want [\".claude/skills\"]",
			skillsDirInit.paths,
		)
	}

	// Check ConfigFileInitializer path
	cfgInit := inits[2].(*ConfigFileInitializer) //nolint:revive // test code, type checked above
	if cfgInit.path != "CLAUDE.md" {
		t.Errorf("ClaudeProvider ConfigFileInitializer path = %s, want \"CLAUDE.md\"", cfgInit.path)
	}

	// Check SlashCommandsInitializer dir
	slashInit := inits[3].(*SlashCommandsInitializer) //nolint:revive // test code, type checked above
	if slashInit.dir != testClaudeCommandsDir {
		t.Errorf(
			"ClaudeProvider SlashCommandsInitializer dir = %s, want \".claude/commands/spectr\"",
			slashInit.dir,
		)
	}
	if len(slashInit.commands) != 2 {
		t.Errorf(
			"ClaudeProvider SlashCommandsInitializer has %d commands, want 2",
			len(slashInit.commands),
		)
	}

	// Check AgentSkillsInitializer
	skillInit := inits[4].(*AgentSkillsInitializer) //nolint:revive // test code, type checked above
	if skillInit.skillName != "spectr-accept-wo-spectr-bin" {
		t.Errorf(
			"ClaudeProvider AgentSkillsInitializer skillName = %s, want \"spectr-accept-wo-spectr-bin\"",
			skillInit.skillName,
		)
	}
	if skillInit.targetDir != ".claude/skills/spectr-accept-wo-spectr-bin" {
		t.Errorf(
			"ClaudeProvider AgentSkillsInitializer targetDir = %s, want \".claude/skills/spectr-accept-wo-spectr-bin\"",
			skillInit.targetDir,
		)
	}
}

func TestGeminiProvider_Initializers(t *testing.T) {
	p := &GeminiProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Gemini should return 2 initializers: Directory, TOMLSlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf("GeminiProvider.Initializers() returned %d initializers, want 2", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("GeminiProvider.Initializers()[0] is %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*TOMLSlashCommandsInitializer); !ok {
		t.Errorf(
			"GeminiProvider.Initializers()[1] is %T, want *TOMLSlashCommandsInitializer",
			inits[1],
		)
	}

	// Check DirectoryInitializer paths
	dirInit := inits[0].(*DirectoryInitializer)
	if len(dirInit.paths) != 1 || dirInit.paths[0] != testGeminiCommandsDir {
		t.Errorf(
			"GeminiProvider DirectoryInitializer paths = %v, want [\".gemini/commands/spectr\"]",
			dirInit.paths,
		)
	}

	// Check TOMLSlashCommandsInitializer dir
	slashInit := inits[1].(*TOMLSlashCommandsInitializer)
	if slashInit.dir != testGeminiCommandsDir {
		t.Errorf(
			"GeminiProvider TOMLSlashCommandsInitializer dir = %s, want \".gemini/commands/spectr\"",
			slashInit.dir,
		)
	}
	if len(slashInit.commands) != 2 {
		t.Errorf(
			"GeminiProvider TOMLSlashCommandsInitializer has %d commands, want 2",
			len(slashInit.commands),
		)
	}
}

func TestCostrictProvider_Initializers(t *testing.T) {
	p := &CostrictProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Costrict should return 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Fatalf("CostrictProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check ConfigFileInitializer path
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "COSTRICT.md" {
		t.Errorf(
			"CostrictProvider ConfigFileInitializer path = %s, want \"COSTRICT.md\"",
			cfgInit.path,
		)
	}
}

func TestQoderProvider_Initializers(t *testing.T) {
	p := &QoderProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Qoder should return 3 initializers
	if len(inits) != 3 {
		t.Fatalf("QoderProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "QODER.md" {
		t.Errorf("QoderProvider ConfigFileInitializer path = %s, want \"QODER.md\"", cfgInit.path)
	}
}

func TestQwenProvider_Initializers(t *testing.T) {
	p := &QwenProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Qwen should return 3 initializers
	if len(inits) != 3 {
		t.Fatalf("QwenProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "QWEN.md" {
		t.Errorf("QwenProvider ConfigFileInitializer path = %s, want \"QWEN.md\"", cfgInit.path)
	}
}

func TestAntigravityProvider_Initializers(t *testing.T) {
	p := &AntigravityProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Antigravity should return 3 initializers: Directory, ConfigFile, PrefixedSlashCommands
	if len(inits) != 3 {
		t.Fatalf("AntigravityProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check types
	if _, ok := inits[2].(*PrefixedSlashCommandsInitializer); !ok {
		t.Errorf(
			"AntigravityProvider.Initializers()[2] is %T, want *PrefixedSlashCommandsInitializer",
			inits[2],
		)
	}

	// Check PrefixedSlashCommandsInitializer
	slashInit := inits[2].(*PrefixedSlashCommandsInitializer)
	if slashInit.dir != testAgentWorkflowsDir {
		t.Errorf(
			"AntigravityProvider PrefixedSlashCommandsInitializer dir = %s, want \".agent/workflows\"",
			slashInit.dir,
		)
	}
	if slashInit.prefix != testSpectrPrefix {
		t.Errorf(
			"AntigravityProvider PrefixedSlashCommandsInitializer prefix = %s, want \"spectr-\"",
			slashInit.prefix,
		)
	}

	// Check ConfigFileInitializer uses AGENTS.md
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "AGENTS.md" {
		t.Errorf(
			"AntigravityProvider ConfigFileInitializer path = %s, want \"AGENTS.md\"",
			cfgInit.path,
		)
	}
}

func TestClineProvider_Initializers(t *testing.T) {
	p := &ClineProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 3 {
		t.Fatalf("ClineProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "CLINE.md" {
		t.Errorf("ClineProvider ConfigFileInitializer path = %s, want \"CLINE.md\"", cfgInit.path)
	}
}

func TestCursorProvider_Initializers(t *testing.T) {
	p := &CursorProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Cursor should return 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf("CursorProvider.Initializers() returned %d initializers, want 2", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("CursorProvider.Initializers()[0] is %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*SlashCommandsInitializer); !ok {
		t.Errorf("CursorProvider.Initializers()[1] is %T, want *SlashCommandsInitializer", inits[1])
	}

	dirInit := inits[0].(*DirectoryInitializer)
	if len(dirInit.paths) != 1 || dirInit.paths[0] != ".cursorrules/commands/spectr" {
		t.Errorf(
			"CursorProvider DirectoryInitializer paths = %v, want [\".cursorrules/commands/spectr\"]",
			dirInit.paths,
		)
	}
}

func TestCodexProvider_Initializers(t *testing.T) {
	p := &CodexProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Codex should return 3 initializers: HomeDirectory, ConfigFile, HomePrefixedSlashCommands
	if len(inits) != 3 {
		t.Fatalf("CodexProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*HomeDirectoryInitializer); !ok {
		t.Errorf("CodexProvider.Initializers()[0] is %T, want *HomeDirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*ConfigFileInitializer); !ok {
		t.Errorf("CodexProvider.Initializers()[1] is %T, want *ConfigFileInitializer", inits[1])
	}
	if _, ok := inits[2].(*HomePrefixedSlashCommandsInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[2] is %T, want *HomePrefixedSlashCommandsInitializer",
			inits[2],
		)
	}

	// Check HomeDirectoryInitializer paths
	dirInit := inits[0].(*HomeDirectoryInitializer)
	if len(dirInit.paths) != 1 || dirInit.paths[0] != testCodexPromptsDir {
		t.Errorf(
			"CodexProvider HomeDirectoryInitializer paths = %v, want [\".codex/prompts\"]",
			dirInit.paths,
		)
	}

	// Check HomePrefixedSlashCommandsInitializer
	slashInit := inits[2].(*HomePrefixedSlashCommandsInitializer)
	if slashInit.dir != testCodexPromptsDir {
		t.Errorf(
			"CodexProvider HomePrefixedSlashCommandsInitializer dir = %s, want \".codex/prompts\"",
			slashInit.dir,
		)
	}
	if slashInit.prefix != testSpectrPrefix {
		t.Errorf(
			"CodexProvider HomePrefixedSlashCommandsInitializer prefix = %s, want \"spectr-\"",
			slashInit.prefix,
		)
	}

	// Check ConfigFileInitializer uses AGENTS.md
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "AGENTS.md" {
		t.Errorf("CodexProvider ConfigFileInitializer path = %s, want \"AGENTS.md\"", cfgInit.path)
	}
}

func TestAiderProvider_Initializers(t *testing.T) {
	p := &AiderProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Aider should return 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf("AiderProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

func TestWindsurfProvider_Initializers(t *testing.T) {
	p := &WindsurfProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 2 {
		t.Fatalf("WindsurfProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

func TestKilocodeProvider_Initializers(t *testing.T) {
	p := &KilocodeProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 2 {
		t.Fatalf("KilocodeProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

func TestContinueProvider_Initializers(t *testing.T) {
	p := &ContinueProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 2 {
		t.Fatalf("ContinueProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

func TestCrushProvider_Initializers(t *testing.T) {
	p := &CrushProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Crush should return 3 initializers
	if len(inits) != 3 {
		t.Fatalf("CrushProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "CRUSH.md" {
		t.Errorf("CrushProvider ConfigFileInitializer path = %s, want \"CRUSH.md\"", cfgInit.path)
	}
}

func TestOpencodeProvider_Initializers(t *testing.T) {
	p := &OpencodeProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Opencode should return 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf("OpencodeProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

// Test all 15 providers return expected initializer counts and types

func TestAllProviders_InitializerCounts(t *testing.T) {
	ctx := context.Background()
	tm := &mockTemplateManager{}

	tests := []struct {
		name          string
		provider      Provider
		expectedCount int
		hasConfigFile bool
		usesTOML      bool
		usesPrefix    bool
		usesHomeFs    bool
	}{
		{"claude-code", &ClaudeProvider{}, 5, true, false, false, false},
		{"gemini", &GeminiProvider{}, 2, false, true, false, false},
		{"costrict", &CostrictProvider{}, 3, true, false, false, false},
		{"qoder", &QoderProvider{}, 3, true, false, false, false},
		{"qwen", &QwenProvider{}, 3, true, false, false, false},
		{"antigravity", &AntigravityProvider{}, 3, true, false, true, false},
		{"cline", &ClineProvider{}, 3, true, false, false, false},
		{"cursor", &CursorProvider{}, 2, false, false, false, false},
		{"codex", &CodexProvider{}, 3, true, false, true, true},
		{"aider", &AiderProvider{}, 2, false, false, false, false},
		{"windsurf", &WindsurfProvider{}, 2, false, false, false, false},
		{"kilocode", &KilocodeProvider{}, 2, false, false, false, false},
		{"continue", &ContinueProvider{}, 2, false, false, false, false},
		{"crush", &CrushProvider{}, 3, true, false, false, false},
		{"opencode", &OpencodeProvider{}, 2, false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inits := tt.provider.Initializers(ctx, tm)

			// Check count
			if len(inits) != tt.expectedCount {
				t.Errorf(
					"%s returned %d initializers, want %d",
					tt.name,
					len(inits),
					tt.expectedCount,
				)
			}

			// Check for directory initializer (all should have one)
			hasDir := false
			hasHomeDir := false
			hasConfig := false
			hasSlash := false
			hasHomeSlash := false
			hasToml := false
			hasPrefix := false
			hasHomePrefix := false

			for _, init := range inits {
				switch init.(type) {
				case *DirectoryInitializer:
					hasDir = true
				case *HomeDirectoryInitializer:
					hasHomeDir = true
				case *ConfigFileInitializer:
					hasConfig = true
				case *SlashCommandsInitializer:
					hasSlash = true
				case *HomeSlashCommandsInitializer:
					hasHomeSlash = true
				case *TOMLSlashCommandsInitializer:
					hasToml = true
				case *PrefixedSlashCommandsInitializer:
					hasPrefix = true
				case *HomePrefixedSlashCommandsInitializer:
					hasHomePrefix = true
				}
			}

			// All providers should have either project or home directory
			if !hasDir && !hasHomeDir {
				t.Errorf("%s has no directory initializer", tt.name)
			}

			// Check config file expectation
			if hasConfig != tt.hasConfigFile {
				t.Errorf("%s hasConfigFile = %v, want %v", tt.name, hasConfig, tt.hasConfigFile)
			}

			// Check TOML expectation
			if hasToml != tt.usesTOML {
				t.Errorf("%s usesTOML = %v, want %v", tt.name, hasToml, tt.usesTOML)
			}

			// Check prefix expectation
			if (hasPrefix || hasHomePrefix) != tt.usesPrefix {
				t.Errorf(
					"%s usesPrefix = %v, want %v",
					tt.name,
					(hasPrefix || hasHomePrefix),
					tt.usesPrefix,
				)
			}

			// Check home filesystem expectation
			if (hasHomeDir || hasHomeSlash || hasHomePrefix) != tt.usesHomeFs {
				t.Errorf(
					"%s usesHomeFs = %v, want %v",
					tt.name,
					(hasHomeDir || hasHomeSlash || hasHomePrefix),
					tt.usesHomeFs,
				)
			}

			// All providers should have slash commands in some form
			if !hasSlash && !hasHomeSlash && !hasToml && !hasPrefix && !hasHomePrefix {
				t.Errorf("%s has no slash command initializer", tt.name)
			}
		})
	}
}

// Test provider registration metadata

func TestProviderRegistration_AllProviders(t *testing.T) {
	// Reset and register all providers
	Reset()
	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	// Expected provider metadata
	expected := []struct {
		id       string
		name     string
		priority int
	}{
		{"claude-code", "Claude Code", 1},
		{"gemini", "Gemini CLI", 2},
		{"costrict", "Costrict", 3},
		{"qoder", "Qoder", 4},
		{"qwen", "Qwen Code", 5},
		{"antigravity", "Antigravity", 6},
		{"cline", "Cline", 7},
		{"cursor", "Cursor", 8},
		{"codex", "Codex CLI", 9},
		{"aider", "Aider", 10},
		{"windsurf", "Windsurf", 11},
		{"kilocode", "Kilocode", 12},
		{"continue", "Continue", 13},
		{"crush", "Crush", 14},
		{"opencode", "OpenCode", 15},
	}

	// Verify count
	if Count() != 15 {
		t.Fatalf("Count() = %d, want 15", Count())
	}

	// Verify each provider
	for _, exp := range expected {
		reg, ok := Get(exp.id)
		if !ok {
			t.Errorf("Provider %s not found in registry", exp.id)

			continue
		}

		if reg.ID != exp.id {
			t.Errorf("Provider %s has ID %s, want %s", exp.id, reg.ID, exp.id)
		}
		if reg.Name != exp.name {
			t.Errorf("Provider %s has Name %s, want %s", exp.id, reg.Name, exp.name)
		}
		if reg.Priority != exp.priority {
			t.Errorf("Provider %s has Priority %d, want %d", exp.id, reg.Priority, exp.priority)
		}
		if reg.Provider == nil {
			t.Errorf("Provider %s has nil Provider", exp.id)
		}
	}

	// Verify priority order
	registered := RegisteredProviders()
	for i := range registered {
		if registered[i].Priority != i+1 {
			t.Errorf(
				"RegisteredProviders()[%d].Priority = %d, want %d (priorities should be sequential 1-15)",
				i,
				registered[i].Priority,
				i+1,
			)
		}
		if i < len(expected) && registered[i].ID != expected[i].id {
			t.Errorf(
				"RegisteredProviders()[%d].ID = %s, want %s (priority order incorrect)",
				i,
				registered[i].ID,
				expected[i].id,
			)
		}
	}
}
