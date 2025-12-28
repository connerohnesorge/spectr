package providers

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// mockTemplateManager provides mock implementations for template manager methods.
// It implements the interface expected by providers.
type mockTemplateManager struct{}

// InstructionPointer returns a mock TemplateRef for instruction pointer.
func (*mockTemplateManager) InstructionPointer() domain.TemplateRef {
	return domain.TemplateRef{Name: "instruction-pointer.md.tmpl", Template: nil}
}

// SlashCommand returns a mock TemplateRef for markdown slash commands.
func (*mockTemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	names := map[domain.SlashCommand]string{
		domain.SlashProposal: "slash-proposal.md.tmpl",
		domain.SlashApply:    "slash-apply.md.tmpl",
	}

	return domain.TemplateRef{Name: names[cmd], Template: nil}
}

// TOMLSlashCommand returns a mock TemplateRef for TOML slash commands.
func (*mockTemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	names := map[domain.SlashCommand]string{
		domain.SlashProposal: "slash-proposal.toml.tmpl",
		domain.SlashApply:    "slash-apply.toml.tmpl",
	}

	return domain.TemplateRef{Name: names[cmd], Template: nil}
}

// providerExpectation describes what a provider should return.
type providerExpectation struct {
	providerID       string
	name             string
	priority         int
	initializerCount int
	configFile       string // empty if no config file
	commandsDir      string
	usesTOML         bool
	usesPrefix       bool
	prefix           string
	usesHomeDir      bool
}

// TestAllProvidersReturnExpectedInitializers tests that all 15 providers return the expected
// initializers with correct counts, types, and paths.
func TestAllProvidersReturnExpectedInitializers(t *testing.T) {
	expectations := []providerExpectation{
		// Priority 1: Claude Code
		{
			providerID:       "claude-code",
			name:             "Claude Code",
			priority:         1,
			initializerCount: 3,
			configFile:       "CLAUDE.md",
			commandsDir:      ".claude/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 2: Gemini (TOML, no config file)
		{
			providerID:       "gemini",
			name:             "Gemini CLI",
			priority:         2,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".gemini/commands/spectr",
			usesTOML:         true,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 3: CoStrict
		{
			providerID:       "costrict",
			name:             "CoStrict",
			priority:         3,
			initializerCount: 3,
			configFile:       "COSTRICT.md",
			commandsDir:      ".costrict/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 4: Qoder
		{
			providerID:       "qoder",
			name:             "Qoder",
			priority:         4,
			initializerCount: 3,
			configFile:       "QODER.md",
			commandsDir:      ".qoder/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 5: Qwen
		{
			providerID:       "qwen",
			name:             "Qwen Code",
			priority:         5,
			initializerCount: 3,
			configFile:       "QWEN.md",
			commandsDir:      ".qwen/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 6: Antigravity (uses prefix, AGENTS.md)
		{
			providerID:       "antigravity",
			name:             "Antigravity",
			priority:         6,
			initializerCount: 3,
			configFile:       "AGENTS.md",
			commandsDir:      ".agent/workflows",
			usesTOML:         false,
			usesPrefix:       true,
			prefix:           "spectr-",
			usesHomeDir:      false,
		},
		// Priority 7: Cline
		{
			providerID:       "cline",
			name:             "Cline",
			priority:         7,
			initializerCount: 3,
			configFile:       "CLINE.md",
			commandsDir:      ".clinerules/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 8: Cursor (no config file)
		{
			providerID:       "cursor",
			name:             "Cursor",
			priority:         8,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".cursorrules/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 9: Codex (uses home directory, prefix, AGENTS.md)
		{
			providerID:       "codex",
			name:             "Codex CLI",
			priority:         9,
			initializerCount: 3,
			configFile:       "AGENTS.md",
			commandsDir:      ".codex/prompts",
			usesTOML:         false,
			usesPrefix:       true,
			prefix:           "spectr-",
			usesHomeDir:      true,
		},
		// Priority 10: Aider (no config file)
		{
			providerID:       "aider",
			name:             "Aider",
			priority:         10,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".aider/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 11: Windsurf (no config file)
		{
			providerID:       "windsurf",
			name:             "Windsurf",
			priority:         11,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".windsurf/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 12: Kilocode (no config file)
		{
			providerID:       "kilocode",
			name:             "Kilocode",
			priority:         12,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".kilocode/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 13: Continue (no config file)
		{
			providerID:       "continue",
			name:             "Continue",
			priority:         13,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".continue/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 14: Crush
		{
			providerID:       "crush",
			name:             "Crush",
			priority:         14,
			initializerCount: 3,
			configFile:       "CRUSH.md",
			commandsDir:      ".crush/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
		// Priority 15: OpenCode (no config file)
		{
			providerID:       "opencode",
			name:             "OpenCode",
			priority:         15,
			initializerCount: 2,
			configFile:       "",
			commandsDir:      ".opencode/commands/spectr",
			usesTOML:         false,
			usesPrefix:       false,
			usesHomeDir:      false,
		},
	}

	// Reset and register all providers
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	mockTM := &mockTemplateManager{}

	for _, exp := range expectations {
		t.Run(exp.providerID, func(t *testing.T) {
			// Get provider from registry
			reg, found := Get(exp.providerID)
			if !found {
				t.Fatalf("Provider %q not found in registry", exp.providerID)
			}

			// Verify registration metadata
			if reg.Name != exp.name {
				t.Errorf("Provider %q Name = %q, want %q", exp.providerID, reg.Name, exp.name)
			}
			if reg.Priority != exp.priority {
				t.Errorf(
					"Provider %q Priority = %d, want %d",
					exp.providerID,
					reg.Priority,
					exp.priority,
				)
			}

			// Get initializers
			inits := reg.Provider.Initializers(context.Background(), mockTM)
			if len(inits) != exp.initializerCount {
				t.Errorf(
					"Provider %q returned %d initializers, want %d",
					exp.providerID,
					len(inits),
					exp.initializerCount,
				)
			}

			// Verify initializer types and configurations
			verifyInitializerTypes(t, exp, inits)
		})
	}
}

// verifyInitializerTypes checks that the returned initializers have the correct types.
//
//nolint:revive,gocritic // cyclomatic - test helper needs to check all initializer types; hugeParam - test code
func verifyInitializerTypes(t *testing.T, exp providerExpectation, inits []domain.Initializer) {
	t.Helper()

	var hasDirectory, hasHomeDirectory, hasConfigFile, hasSlashCommands bool

	for _, init := range inits {
		switch v := init.(type) {
		case *initializers.DirectoryInitializer:
			hasDirectory = true
			// Verify directory path via dedupeKey
			expectedKey := "DirectoryInitializer:" + exp.commandsDir
			if v.DedupeKey() != expectedKey {
				t.Errorf("DirectoryInitializer path = %q, want %q", v.DedupeKey(), expectedKey)
			}

		case *initializers.HomeDirectoryInitializer:
			hasHomeDirectory = true
			if !exp.usesHomeDir {
				t.Errorf("Provider %q has HomeDirectoryInitializer but usesHomeDir=false", exp.providerID)
			}

		case *initializers.ConfigFileInitializer:
			hasConfigFile = true
			if exp.configFile == "" {
				t.Errorf("Provider %q has ConfigFileInitializer but expected no config file", exp.providerID)
			}
			// Verify path via dedupeKey
			expectedKey := "ConfigFileInitializer:" + exp.configFile
			if v.DedupeKey() != expectedKey {
				t.Errorf("ConfigFileInitializer path = %q, want %q", v.DedupeKey(), expectedKey)
			}

		case *initializers.SlashCommandsInitializer:
			hasSlashCommands = true
			if exp.usesTOML {
				t.Errorf("Provider %q has SlashCommandsInitializer but expected TOML", exp.providerID)
			}
			if exp.usesPrefix {
				t.Errorf("Provider %q has SlashCommandsInitializer but expected prefixed", exp.providerID)
			}
			if exp.usesHomeDir {
				t.Errorf("Provider %q has SlashCommandsInitializer but expected home directory", exp.providerID)
			}

		case *initializers.HomeSlashCommandsInitializer:
			hasSlashCommands = true
			if !exp.usesHomeDir {
				t.Errorf("Provider %q has HomeSlashCommandsInitializer but usesHomeDir=false", exp.providerID)
			}

		case *initializers.PrefixedSlashCommandsInitializer:
			hasSlashCommands = true
			if !exp.usesPrefix {
				t.Errorf("Provider %q has PrefixedSlashCommandsInitializer but usesPrefix=false", exp.providerID)
			}
			if exp.usesHomeDir {
				t.Errorf("Provider %q has PrefixedSlashCommandsInitializer but expected home directory", exp.providerID)
			}
			// Verify prefix in dedupeKey
			expectedKey := "PrefixedSlashCommandsInitializer:" + exp.commandsDir + ":" + exp.prefix
			if v.DedupeKey() != expectedKey {
				t.Errorf("PrefixedSlashCommandsInitializer key = %q, want %q", v.DedupeKey(), expectedKey)
			}

		case *initializers.HomePrefixedSlashCommandsInitializer:
			hasSlashCommands = true
			if !exp.usesHomeDir {
				t.Errorf("Provider %q has HomePrefixedSlashCommandsInitializer but usesHomeDir=false", exp.providerID)
			}
			if !exp.usesPrefix {
				t.Errorf("Provider %q has HomePrefixedSlashCommandsInitializer but usesPrefix=false", exp.providerID)
			}
			// Verify prefix in dedupeKey
			expectedKey := "HomePrefixedSlashCommandsInitializer:" + exp.commandsDir + ":" + exp.prefix
			if v.DedupeKey() != expectedKey {
				t.Errorf("HomePrefixedSlashCommandsInitializer key = %q, want %q", v.DedupeKey(), expectedKey)
			}

		case *initializers.TOMLSlashCommandsInitializer:
			hasSlashCommands = true
			if !exp.usesTOML {
				t.Errorf("Provider %q has TOMLSlashCommandsInitializer but usesTOML=false", exp.providerID)
			}
			// Verify path via dedupeKey
			expectedKey := "TOMLSlashCommandsInitializer:" + exp.commandsDir
			if v.DedupeKey() != expectedKey {
				t.Errorf("TOMLSlashCommandsInitializer path = %q, want %q", v.DedupeKey(), expectedKey)
			}

		default:
			t.Errorf("Provider %q has unexpected initializer type: %T", exp.providerID, v)
		}
	}

	// Verify required initializers are present
	if !hasDirectory && !hasHomeDirectory {
		t.Errorf(
			"Provider %q is missing DirectoryInitializer or HomeDirectoryInitializer",
			exp.providerID,
		)
	}
	if exp.configFile != "" && !hasConfigFile {
		t.Errorf(
			"Provider %q is missing ConfigFileInitializer for %q",
			exp.providerID,
			exp.configFile,
		)
	}
	if exp.configFile == "" && hasConfigFile {
		t.Errorf("Provider %q has ConfigFileInitializer but expected none", exp.providerID)
	}
	if !hasSlashCommands {
		t.Errorf("Provider %q is missing slash commands initializer", exp.providerID)
	}
}

// TestProviderInitializersWithNilTemplateManager tests that providers handle nil/invalid
// template manager gracefully.
func TestProviderInitializersWithNilTemplateManager(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providers := RegisteredProviders()
	for _, reg := range providers {
		t.Run(reg.ID+"_nil_tm", func(t *testing.T) {
			// Should return nil when template manager doesn't implement required interface
			inits := reg.Provider.Initializers(context.Background(), nil)
			if inits != nil {
				t.Errorf(
					"Provider %q returned non-nil initializers with nil template manager",
					reg.ID,
				)
			}
		})

		t.Run(reg.ID+"_wrong_type_tm", func(t *testing.T) {
			// Should return nil when template manager is wrong type
			inits := reg.Provider.Initializers(context.Background(), "wrong type")
			if inits != nil {
				t.Errorf(
					"Provider %q returned non-nil initializers with wrong type template manager",
					reg.ID,
				)
			}
		})
	}
}

// TestProviderRegistrationMetadata verifies that all providers have correct registration metadata.
func TestProviderRegistrationMetadata(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	expectedMetadata := map[string]struct {
		name     string
		priority int
	}{
		"claude-code": {name: "Claude Code", priority: 1},
		"gemini":      {name: "Gemini CLI", priority: 2},
		"costrict":    {name: "CoStrict", priority: 3},
		"qoder":       {name: "Qoder", priority: 4},
		"qwen":        {name: "Qwen Code", priority: 5},
		"antigravity": {name: "Antigravity", priority: 6},
		"cline":       {name: "Cline", priority: 7},
		"cursor":      {name: "Cursor", priority: 8},
		"codex":       {name: "Codex CLI", priority: 9},
		"aider":       {name: "Aider", priority: 10},
		"windsurf":    {name: "Windsurf", priority: 11},
		"kilocode":    {name: "Kilocode", priority: 12},
		"continue":    {name: "Continue", priority: 13},
		"crush":       {name: "Crush", priority: 14},
		"opencode":    {name: "OpenCode", priority: 15},
	}

	for id, expected := range expectedMetadata {
		t.Run(id, func(t *testing.T) {
			reg, found := Get(id)
			if !found {
				t.Fatalf("Provider %q not found", id)
			}

			if reg.Name != expected.name {
				t.Errorf("Provider %q Name = %q, want %q", id, reg.Name, expected.name)
			}
			if reg.Priority != expected.priority {
				t.Errorf("Provider %q Priority = %d, want %d", id, reg.Priority, expected.priority)
			}
			if reg.Provider == nil {
				t.Errorf("Provider %q has nil Provider", id)
			}
		})
	}
}

// TestProviderPrioritiesAreSequential verifies priorities are 1-15 with no gaps.
func TestProviderPrioritiesAreSequential(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providers := RegisteredProviders()
	if len(providers) != 15 {
		t.Fatalf("Expected 15 providers, got %d", len(providers))
	}

	// Since RegisteredProviders returns sorted by priority
	for i, p := range providers {
		expectedPriority := i + 1
		if p.Priority != expectedPriority {
			t.Errorf(
				"providers[%d] has Priority=%d, want %d (ID: %s)",
				i,
				p.Priority,
				expectedPriority,
				p.ID,
			)
		}
	}
}

// TestProvidersWithConfigFile verifies providers that should have config files.
func TestProvidersWithConfigFile(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providersWithConfig := map[string]string{
		"claude-code": "CLAUDE.md",
		"costrict":    "COSTRICT.md",
		"qoder":       "QODER.md",
		"qwen":        "QWEN.md",
		"antigravity": "AGENTS.md",
		"cline":       "CLINE.md",
		"codex":       "AGENTS.md",
		"crush":       "CRUSH.md",
	}

	mockTM := &mockTemplateManager{}

	for id, expectedConfig := range providersWithConfig {
		t.Run(id, func(t *testing.T) {
			reg, found := Get(id)
			if !found {
				t.Fatalf("Provider %q not found", id)
			}

			inits := reg.Provider.Initializers(context.Background(), mockTM)

			foundConfig := false
			for _, init := range inits {
				cfg, ok := init.(*initializers.ConfigFileInitializer)
				if !ok {
					continue
				}
				foundConfig = true
				expectedKey := "ConfigFileInitializer:" + expectedConfig
				if cfg.DedupeKey() != expectedKey {
					t.Errorf(
						"Provider %q config file = %q, want %q",
						id,
						cfg.DedupeKey(),
						expectedKey,
					)
				}
			}

			if !foundConfig {
				t.Errorf("Provider %q should have ConfigFileInitializer for %q", id, expectedConfig)
			}
		})
	}
}

// TestProvidersWithoutConfigFile verifies providers that should NOT have config files.
func TestProvidersWithoutConfigFile(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providersWithoutConfig := []string{
		"gemini",
		"cursor",
		"aider",
		"windsurf",
		"kilocode",
		"continue",
		"opencode",
	}

	mockTM := &mockTemplateManager{}

	for _, id := range providersWithoutConfig {
		t.Run(id, func(t *testing.T) {
			reg, found := Get(id)
			if !found {
				t.Fatalf("Provider %q not found", id)
			}

			inits := reg.Provider.Initializers(context.Background(), mockTM)

			for _, init := range inits {
				if _, ok := init.(*initializers.ConfigFileInitializer); ok {
					t.Errorf("Provider %q should NOT have ConfigFileInitializer", id)
				}
			}
		})
	}
}

// TestGeminiProviderUsesTOML verifies that Gemini uses TOML format.
func TestGeminiProviderUsesTOML(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	reg, found := Get("gemini")
	if !found {
		t.Fatal("Gemini provider not found")
	}

	mockTM := &mockTemplateManager{}
	inits := reg.Provider.Initializers(context.Background(), mockTM)

	foundTOML := false
	for _, init := range inits {
		if _, ok := init.(*initializers.TOMLSlashCommandsInitializer); ok {
			foundTOML = true
		}
		// Should NOT have regular SlashCommandsInitializer
		if _, ok := init.(*initializers.SlashCommandsInitializer); ok {
			t.Error("Gemini should not have SlashCommandsInitializer")
		}
	}

	if !foundTOML {
		t.Error("Gemini should have TOMLSlashCommandsInitializer")
	}
}

// TestCodexProviderUsesHomeDirectory verifies that Codex uses home directory.
func TestCodexProviderUsesHomeDirectory(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	reg, found := Get("codex")
	if !found {
		t.Fatal("Codex provider not found")
	}

	mockTM := &mockTemplateManager{}
	inits := reg.Provider.Initializers(context.Background(), mockTM)

	foundHomeDir := false
	foundHomePrefixed := false

	for _, init := range inits {
		switch v := init.(type) {
		case *initializers.HomeDirectoryInitializer:
			foundHomeDir = true
			expectedKey := "HomeDirectoryInitializer:.codex/prompts"
			if v.DedupeKey() != expectedKey {
				t.Errorf("Codex HomeDirectoryInitializer key = %q, want %q", v.DedupeKey(), expectedKey)
			}
		case *initializers.HomePrefixedSlashCommandsInitializer:
			foundHomePrefixed = true
			expectedKey := "HomePrefixedSlashCommandsInitializer:.codex/prompts:spectr-"
			if v.DedupeKey() != expectedKey {
				t.Errorf("Codex HomePrefixedSlashCommandsInitializer key = %q, want %q", v.DedupeKey(), expectedKey)
			}
		case *initializers.DirectoryInitializer:
			t.Error("Codex should not have DirectoryInitializer")
		case *initializers.SlashCommandsInitializer:
			t.Error("Codex should not have SlashCommandsInitializer")
		}
	}

	if !foundHomeDir {
		t.Error("Codex should have HomeDirectoryInitializer")
	}
	if !foundHomePrefixed {
		t.Error("Codex should have HomePrefixedSlashCommandsInitializer")
	}
}

// TestAntigravityProviderUsesPrefixedSlashCommands verifies Antigravity uses prefixed slash commands.
func TestAntigravityProviderUsesPrefixedSlashCommands(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	reg, found := Get("antigravity")
	if !found {
		t.Fatal("Antigravity provider not found")
	}

	mockTM := &mockTemplateManager{}
	inits := reg.Provider.Initializers(context.Background(), mockTM)

	foundPrefixed := false
	for _, init := range inits {
		if v, ok := init.(*initializers.PrefixedSlashCommandsInitializer); ok {
			foundPrefixed = true
			expectedKey := "PrefixedSlashCommandsInitializer:.agent/workflows:spectr-"
			if v.DedupeKey() != expectedKey {
				t.Errorf(
					"Antigravity PrefixedSlashCommandsInitializer key = %q, want %q",
					v.DedupeKey(),
					expectedKey,
				)
			}
		}
		// Should NOT have regular SlashCommandsInitializer
		if _, ok := init.(*initializers.SlashCommandsInitializer); ok {
			t.Error("Antigravity should not have SlashCommandsInitializer")
		}
	}

	if !foundPrefixed {
		t.Error("Antigravity should have PrefixedSlashCommandsInitializer")
	}
}

// TestAllProvidersHaveDirectoryInitializer verifies all providers create directories.
func TestAllProvidersHaveDirectoryInitializer(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	mockTM := &mockTemplateManager{}
	providers := RegisteredProviders()

	for _, reg := range providers {
		t.Run(reg.ID, func(t *testing.T) {
			inits := reg.Provider.Initializers(context.Background(), mockTM)

			hasDirectoryInit := false
			for _, init := range inits {
				switch init.(type) {
				case *initializers.DirectoryInitializer, *initializers.HomeDirectoryInitializer:
					hasDirectoryInit = true
				}
			}

			if !hasDirectoryInit {
				t.Errorf("Provider %q should have a directory initializer", reg.ID)
			}
		})
	}
}

// TestAllProvidersHaveSlashCommandsInitializer verifies all providers create slash commands.
func TestAllProvidersHaveSlashCommandsInitializer(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	mockTM := &mockTemplateManager{}
	providers := RegisteredProviders()

	for _, reg := range providers {
		t.Run(reg.ID, func(t *testing.T) {
			inits := reg.Provider.Initializers(context.Background(), mockTM)

			hasSlashCmdInit := false
			for _, init := range inits {
				switch init.(type) {
				case *initializers.SlashCommandsInitializer,
					*initializers.HomeSlashCommandsInitializer,
					*initializers.PrefixedSlashCommandsInitializer,
					*initializers.HomePrefixedSlashCommandsInitializer,
					*initializers.TOMLSlashCommandsInitializer:
					hasSlashCmdInit = true
				}
			}

			if !hasSlashCmdInit {
				t.Errorf("Provider %q should have a slash commands initializer", reg.ID)
			}
		})
	}
}

// TestProviderInitializerCountsMatchExpectation ensures the exact count of initializers.
func TestProviderInitializerCountsMatchExpectation(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	expectedCounts := map[string]int{
		"claude-code": 3, // DirectoryInitializer + ConfigFileInitializer + SlashCommandsInitializer
		"gemini":      2, // DirectoryInitializer + TOMLSlashCommandsInitializer
		"costrict":    3,
		"qoder":       3,
		"qwen":        3,
		"antigravity": 3, // DirectoryInitializer + ConfigFileInitializer + PrefixedSlashCommandsInitializer
		"cline":       3,
		"cursor":      2, // DirectoryInitializer + SlashCommandsInitializer
		"codex":       3, // HomeDirectoryInitializer + ConfigFileInitializer + HomePrefixedSlashCommandsInitializer
		"aider":       2,
		"windsurf":    2,
		"kilocode":    2,
		"continue":    2,
		"crush":       3,
		"opencode":    2,
	}

	mockTM := &mockTemplateManager{}

	for id, expectedCount := range expectedCounts {
		t.Run(id, func(t *testing.T) {
			reg, found := Get(id)
			if !found {
				t.Fatalf("Provider %q not found", id)
			}

			inits := reg.Provider.Initializers(context.Background(), mockTM)
			if len(inits) != expectedCount {
				t.Errorf("Provider %q has %d initializers, want %d", id, len(inits), expectedCount)
			}
		})
	}
}
