package providers

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// TestClaudeProvider tests the Claude Code provider returns expected initializers
func TestClaudeProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &ClaudeProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Errorf("ClaudeProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[0] = %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*ConfigFileInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[1] = %T, want *ConfigFileInitializer", inits[1])
	}
	if _, ok := inits[2].(*SlashCommandsInitializer); !ok {
		t.Errorf("ClaudeProvider.Initializers()[2] = %T, want *SlashCommandsInitializer", inits[2])
	}

	// Check paths
	if dir, ok := inits[0].(*DirectoryInitializer); ok {
		if len(dir.paths) != 1 || dir.paths[0] != ".claude/commands/spectr" {
			t.Errorf(
				"DirectoryInitializer paths = %v, want [\".claude/commands/spectr\"]",
				dir.paths,
			)
		}
	}
	if config, ok := inits[1].(*ConfigFileInitializer); ok {
		if config.path != "CLAUDE.md" {
			t.Errorf("ConfigFileInitializer path = %s, want CLAUDE.md", config.path)
		}
	}
	slash, ok := inits[2].(*SlashCommandsInitializer)
	if !ok {
		return
	}
	if slash.dir != ".claude/commands/spectr" {
		t.Errorf("SlashCommandsInitializer dir = %s, want .claude/commands/spectr", slash.dir)
	}
	if len(slash.commands) != 2 {
		t.Errorf("SlashCommandsInitializer has %d commands, want 2", len(slash.commands))
	}
}

// TestGeminiProvider tests the Gemini provider returns expected initializers
func TestGeminiProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &GeminiProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, TOMLSlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("GeminiProvider.Initializers() returned %d initializers, want 2", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("GeminiProvider.Initializers()[0] = %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*TOMLSlashCommandsInitializer); !ok {
		t.Errorf(
			"GeminiProvider.Initializers()[1] = %T, want *TOMLSlashCommandsInitializer",
			inits[1],
		)
	}

	// Check paths
	if dir, ok := inits[0].(*DirectoryInitializer); ok {
		if len(dir.paths) != 1 || dir.paths[0] != ".gemini/commands/spectr" {
			t.Errorf(
				"DirectoryInitializer paths = %v, want [\".gemini/commands/spectr\"]",
				dir.paths,
			)
		}
	}
	toml, ok := inits[1].(*TOMLSlashCommandsInitializer)
	if !ok {
		return
	}
	if toml.dir != ".gemini/commands/spectr" {
		t.Errorf(
			"TOMLSlashCommandsInitializer dir = %s, want .gemini/commands/spectr",
			toml.dir,
		)
	}
	if len(toml.commands) != 2 {
		t.Errorf("TOMLSlashCommandsInitializer has %d commands, want 2", len(toml.commands))
	}
}

// TestCostrictProvider tests the Costrict provider returns expected initializers
func TestCostrictProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &CostrictProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Errorf("CostrictProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check config file path
	config, ok := inits[1].(*ConfigFileInitializer)
	if !ok {
		return
	}
	if config.path != "COSTRICT.md" {
		t.Errorf("ConfigFileInitializer path = %s, want COSTRICT.md", config.path)
	}
}

// TestQoderProvider tests the Qoder provider returns expected initializers
func TestQoderProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &QoderProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Errorf("QoderProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check config file path
	config, ok := inits[1].(*ConfigFileInitializer)
	if !ok {
		return
	}
	if config.path != "QODER.md" {
		t.Errorf("ConfigFileInitializer path = %s, want QODER.md", config.path)
	}
}

// TestQwenProvider tests the Qwen provider returns expected initializers
func TestQwenProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &QwenProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Errorf("QwenProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check config file path
	config, ok := inits[1].(*ConfigFileInitializer)
	if !ok {
		return
	}
	if config.path != "QWEN.md" {
		t.Errorf("ConfigFileInitializer path = %s, want QWEN.md", config.path)
	}
}

// TestAntigravityProvider tests the Antigravity provider returns expected initializers
func TestAntigravityProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &AntigravityProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, PrefixedSlashCommands
	if len(inits) != 3 {
		t.Errorf("AntigravityProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("AntigravityProvider.Initializers()[0] = %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*ConfigFileInitializer); !ok {
		t.Errorf(
			"AntigravityProvider.Initializers()[1] = %T, want *ConfigFileInitializer",
			inits[1],
		)
	}
	if _, ok := inits[2].(*PrefixedSlashCommandsInitializer); !ok {
		t.Errorf(
			"AntigravityProvider.Initializers()[2] = %T, want *PrefixedSlashCommandsInitializer",
			inits[2],
		)
	}

	// Check paths
	if dir, ok := inits[0].(*DirectoryInitializer); ok {
		if len(dir.paths) != 1 || dir.paths[0] != ".agent/workflows" {
			t.Errorf("DirectoryInitializer paths = %v, want [\".agent/workflows\"]", dir.paths)
		}
	}
	if config, ok := inits[1].(*ConfigFileInitializer); ok {
		if config.path != "AGENTS.md" {
			t.Errorf("ConfigFileInitializer path = %s, want AGENTS.md", config.path)
		}
	}
	prefixed, ok := inits[2].(*PrefixedSlashCommandsInitializer)
	if !ok {
		return
	}
	if prefixed.dir != ".agent/workflows" {
		t.Errorf(
			"PrefixedSlashCommandsInitializer dir = %s, want .agent/workflows",
			prefixed.dir,
		)
	}
	if prefixed.prefix != "spectr-" {
		t.Errorf("PrefixedSlashCommandsInitializer prefix = %s, want spectr-", prefixed.prefix)
	}
}

// TestClineProvider tests the Cline provider returns expected initializers
func TestClineProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &ClineProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Errorf("ClineProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check config file path
	config, ok := inits[1].(*ConfigFileInitializer)
	if !ok {
		return
	}
	if config.path != "CLINE.md" {
		t.Errorf("ConfigFileInitializer path = %s, want CLINE.md", config.path)
	}
}

// TestCursorProvider tests the Cursor provider returns expected initializers
func TestCursorProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &CursorProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("CursorProvider.Initializers() returned %d initializers, want 2", len(inits))
	}

	// Check types (no ConfigFileInitializer)
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf("CursorProvider.Initializers()[0] = %T, want *DirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*SlashCommandsInitializer); !ok {
		t.Errorf("CursorProvider.Initializers()[1] = %T, want *SlashCommandsInitializer", inits[1])
	}
}

// TestCodexProvider tests the Codex provider returns expected initializers
func TestCodexProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &CodexProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: HomeDirectory, ConfigFile, HomePrefixedSlashCommands
	if len(inits) != 3 {
		t.Errorf("CodexProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check types (uses home filesystem)
	if _, ok := inits[0].(*HomeDirectoryInitializer); !ok {
		t.Errorf("CodexProvider.Initializers()[0] = %T, want *HomeDirectoryInitializer", inits[0])
	}
	if _, ok := inits[1].(*ConfigFileInitializer); !ok {
		t.Errorf("CodexProvider.Initializers()[1] = %T, want *ConfigFileInitializer", inits[1])
	}
	if _, ok := inits[2].(*HomePrefixedSlashCommandsInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[2] = %T, want *HomePrefixedSlashCommandsInitializer",
			inits[2],
		)
	}

	// Check paths
	if homeDir, ok := inits[0].(*HomeDirectoryInitializer); ok {
		if len(homeDir.paths) != 1 || homeDir.paths[0] != ".codex/prompts" {
			t.Errorf(
				"HomeDirectoryInitializer paths = %v, want [\".codex/prompts\"]",
				homeDir.paths,
			)
		}
	}
	if config, ok := inits[1].(*ConfigFileInitializer); ok {
		if config.path != "AGENTS.md" {
			t.Errorf("ConfigFileInitializer path = %s, want AGENTS.md", config.path)
		}
	}
	homePrefixed, ok := inits[2].(*HomePrefixedSlashCommandsInitializer)
	if !ok {
		return
	}
	if homePrefixed.dir != ".codex/prompts" {
		t.Errorf(
			"HomePrefixedSlashCommandsInitializer dir = %s, want .codex/prompts",
			homePrefixed.dir,
		)
	}
	if homePrefixed.prefix != "spectr-" {
		t.Errorf(
			"HomePrefixedSlashCommandsInitializer prefix = %s, want spectr-",
			homePrefixed.prefix,
		)
	}
}

// TestAiderProvider tests the Aider provider returns expected initializers
func TestAiderProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &AiderProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("AiderProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

// TestWindsurfProvider tests the Windsurf provider returns expected initializers
func TestWindsurfProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &WindsurfProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("WindsurfProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

// TestKilocodeProvider tests the Kilocode provider returns expected initializers
func TestKilocodeProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &KilocodeProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("KilocodeProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

// TestContinueProvider tests the Continue provider returns expected initializers
func TestContinueProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &ContinueProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("ContinueProvider.Initializers() returned %d initializers, want 2", len(inits))
	}
}

// TestCrushProvider tests the Crush provider returns expected initializers
func TestCrushProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &CrushProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Errorf("CrushProvider.Initializers() returned %d initializers, want 3", len(inits))
	}

	// Check config file path
	config, ok := inits[1].(*ConfigFileInitializer)
	if !ok {
		return
	}
	if config.path != "CRUSH.md" {
		t.Errorf("ConfigFileInitializer path = %s, want CRUSH.md", config.path)
	}
}

// TestOpencodeProvider tests the OpenCode provider returns expected initializers
func TestOpencodeProvider(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	p := &OpencodeProvider{}
	inits := p.Initializers(ctx, tm)

	// Expect 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Errorf("OpencodeProvider.Initializers() returned %d initializers, want 2", len(inits))
	}

	// Check paths - OpenCode uses .opencode/command/spectr (note: command not commands)
	dir, ok := inits[0].(*DirectoryInitializer)
	if !ok {
		return
	}
	if len(dir.paths) != 1 || dir.paths[0] != ".opencode/command/spectr" {
		t.Errorf(
			"DirectoryInitializer paths = %v, want [\".opencode/command/spectr\"]",
			dir.paths,
		)
	}
}

// TestProviderSlashCommands tests that providers set up slash commands correctly
func TestProviderSlashCommands(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	tests := []struct {
		name     string
		provider Provider
	}{
		{"Claude", &ClaudeProvider{}},
		{"Gemini", &GeminiProvider{}},
		{"Costrict", &CostrictProvider{}},
		{"Qoder", &QoderProvider{}},
		{"Qwen", &QwenProvider{}},
		{"Antigravity", &AntigravityProvider{}},
		{"Cline", &ClineProvider{}},
		{"Cursor", &CursorProvider{}},
		{"Codex", &CodexProvider{}},
		{"Aider", &AiderProvider{}},
		{"Windsurf", &WindsurfProvider{}},
		{"Kilocode", &KilocodeProvider{}},
		{"Continue", &ContinueProvider{}},
		{"Crush", &CrushProvider{}},
		{"OpenCode", &OpencodeProvider{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inits := tt.provider.Initializers(ctx, tm)

			// Every provider should have at least one slash command initializer
			hasSlashCommands := false
			for _, init := range inits {
				switch init.(type) {
				case *SlashCommandsInitializer, *HomeSlashCommandsInitializer,
					*PrefixedSlashCommandsInitializer, *HomePrefixedSlashCommandsInitializer,
					*TOMLSlashCommandsInitializer:
					hasSlashCommands = true
				}
			}

			if !hasSlashCommands {
				t.Errorf("%s provider has no slash command initializer", tt.name)
			}
		})
	}
}

// TestProviderSlashCommandsHaveBothCommands tests that all slash command initializers include both proposal and apply
func TestProviderSlashCommandsHaveBothCommands(t *testing.T) {
	ctx := context.Background()
	tm, err := templates.NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	tests := []struct {
		name     string
		provider Provider
	}{
		{"Claude", &ClaudeProvider{}},
		{"Gemini", &GeminiProvider{}},
		{"Costrict", &CostrictProvider{}},
		{"Qoder", &QoderProvider{}},
		{"Qwen", &QwenProvider{}},
		{"Antigravity", &AntigravityProvider{}},
		{"Cline", &ClineProvider{}},
		{"Cursor", &CursorProvider{}},
		{"Codex", &CodexProvider{}},
		{"Aider", &AiderProvider{}},
		{"Windsurf", &WindsurfProvider{}},
		{"Kilocode", &KilocodeProvider{}},
		{"Continue", &ContinueProvider{}},
		{"Crush", &CrushProvider{}},
		{"OpenCode", &OpencodeProvider{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inits := tt.provider.Initializers(ctx, tm)

			for _, init := range inits {
				var commands map[domain.SlashCommand]domain.TemplateRef

				switch si := init.(type) {
				case *SlashCommandsInitializer:
					commands = si.commands
				case *HomeSlashCommandsInitializer:
					commands = si.commands
				case *PrefixedSlashCommandsInitializer:
					commands = si.commands
				case *HomePrefixedSlashCommandsInitializer:
					commands = si.commands
				case *TOMLSlashCommandsInitializer:
					commands = si.commands
				}

				if commands == nil {
					continue
				}
				// Check for both proposal and apply
				if _, hasProposal := commands[domain.SlashProposal]; !hasProposal {
					t.Errorf("%s slash commands missing SlashProposal", tt.name)
				}
				if _, hasApply := commands[domain.SlashApply]; !hasApply {
					t.Errorf("%s slash commands missing SlashApply", tt.name)
				}
			}
		})
	}
}
