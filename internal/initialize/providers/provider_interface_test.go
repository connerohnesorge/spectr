package providers

import (
	"context"
	"testing"
)

// TestAllProvidersImplementProvider verifies that all registered providers implement the Provider interface.
func TestAllProvidersImplementProvider(t *testing.T) {
	// Use global registry which has all providers registered via init() functions
	regs := All()

	if len(regs) == 0 {
		t.Fatal("No providers registered in global registry")
	}

	t.Logf("Testing %d registered providers", len(regs))

	for _, reg := range regs {
		t.Run(reg.ID, func(t *testing.T) {
			if reg.Provider == nil {
				t.Fatal("Provider is nil")
			}

			// Verify it implements Provider interface by calling Initializers
			ctx := context.Background()
			initializers := reg.Provider.Initializers(ctx)

			// Should return a non-nil slice (empty is okay, but nil is not)
			if initializers == nil {
				t.Errorf("Provider %s returned nil initializers slice", reg.ID)
			}
		})
	}
}

// TestAllProvidersReturnValidInitializers verifies that all providers return non-nil and valid initializers.
func TestAllProvidersReturnValidInitializers(t *testing.T) {
	regs := All()

	for _, reg := range regs {
		t.Run(reg.ID, func(t *testing.T) {
			ctx := context.Background()
			initializers := reg.Provider.Initializers(ctx)

			// Verify each initializer is valid
			for i, init := range initializers {
				if init == nil {
					t.Errorf("Provider %s returned nil initializer at index %d", reg.ID, i)

					continue
				}

				// Verify Path() returns non-empty string
				path := init.Path()
				if path == "" {
					t.Errorf("Provider %s initializer at index %d has empty Path()", reg.ID, i)
				}

				// Verify IsGlobal() doesn't panic
				_ = init.IsGlobal()
			}
		})
	}
}

// TestProviderInitializerCounts verifies that each provider returns the expected number and types of initializers.
func TestProviderInitializerCounts(t *testing.T) {
	tests := []struct {
		providerID         string
		wantInitCount      int
		wantDirInit        bool
		wantConfigFileInit bool
		wantSlashCmdsInit  bool
	}{
		// Providers with directory + config file + slash commands (3 initializers)
		{
			providerID:         "claude-code",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "cline",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "costrict",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "codex",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "qwen",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "qoder",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "antigravity",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "crush",
			wantInitCount:      3,
			wantDirInit:        true,
			wantConfigFileInit: true,
			wantSlashCmdsInit:  true,
		},

		// Providers with directory + slash commands only (2 initializers, no config file)
		{
			providerID:         "gemini",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "cursor",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "aider",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "windsurf",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "kilocode",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "continue",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
		{
			providerID:         "opencode",
			wantInitCount:      2,
			wantDirInit:        true,
			wantConfigFileInit: false,
			wantSlashCmdsInit:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.providerID, func(t *testing.T) {
			reg := Get(tt.providerID)
			if reg == nil {
				t.Fatalf("Provider %s not found in registry", tt.providerID)
			}

			ctx := context.Background()
			initializers := reg.Provider.Initializers(ctx)

			// Verify count
			if len(initializers) != tt.wantInitCount {
				t.Errorf(
					"Provider %s returned %d initializers, want %d",
					tt.providerID,
					len(initializers),
					tt.wantInitCount,
				)
			}

			// Count initializer types
			var hasDirInit, hasConfigFileInit, hasSlashCmdsInit bool
			for _, init := range initializers {
				switch init.(type) {
				case *DirectoryInitializer:
					hasDirInit = true
				case *ConfigFileInitializer:
					hasConfigFileInit = true
				case *SlashCommandsInitializer:
					hasSlashCmdsInit = true
				}
			}

			// Verify expected types
			if hasDirInit != tt.wantDirInit {
				t.Errorf(
					"Provider %s has DirectoryInitializer=%v, want %v",
					tt.providerID,
					hasDirInit,
					tt.wantDirInit,
				)
			}
			if hasConfigFileInit != tt.wantConfigFileInit {
				t.Errorf(
					"Provider %s has ConfigFileInitializer=%v, want %v",
					tt.providerID,
					hasConfigFileInit,
					tt.wantConfigFileInit,
				)
			}
			if hasSlashCmdsInit != tt.wantSlashCmdsInit {
				t.Errorf(
					"Provider %s has SlashCommandsInitializer=%v, want %v",
					tt.providerID,
					hasSlashCmdsInit,
					tt.wantSlashCmdsInit,
				)
			}
		})
	}
}

// TestProviderInitializerPaths verifies that each provider's initializers have expected paths.
func TestProviderInitializerPaths(t *testing.T) {
	tests := []struct {
		providerID         string
		expectedDirPath    string
		expectedConfigPath string // empty if provider has no config file
		expectedSlashPath  string
	}{
		{
			providerID:         "claude-code",
			expectedDirPath:    ".claude/commands/spectr",
			expectedConfigPath: "CLAUDE.md",
			expectedSlashPath:  ".claude/commands/spectr",
		},
		{
			providerID:         "gemini",
			expectedDirPath:    ".gemini/commands/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".gemini/commands/spectr",
		},
		{
			providerID:         "cline",
			expectedDirPath:    ".clinerules/commands/spectr",
			expectedConfigPath: "CLINE.md",
			expectedSlashPath:  ".clinerules/commands/spectr",
		},
		{
			providerID:         "cursor",
			expectedDirPath:    ".cursorrules/commands/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".cursorrules/commands/spectr",
		},
		{
			providerID:         "aider",
			expectedDirPath:    ".aider/prompts/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".aider/prompts/spectr",
		},
		{
			providerID:         "continue",
			expectedDirPath:    ".continue/commands/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".continue/commands/spectr",
		},
		{
			providerID:         "windsurf",
			expectedDirPath:    ".windsurf/commands/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".windsurf/commands/spectr",
		},
		{
			providerID:         "costrict",
			expectedDirPath:    ".costrict/commands/spectr",
			expectedConfigPath: "COSTRICT.md",
			expectedSlashPath:  ".costrict/commands/spectr",
		},
		{
			providerID:         "kilocode",
			expectedDirPath:    ".kilocode/commands/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".kilocode/commands/spectr",
		},
		{
			providerID:         "codex",
			expectedDirPath:    ".codex/commands/spectr",
			expectedConfigPath: "AGENTS.md",
			expectedSlashPath:  ".codex/commands/spectr",
		},
		{
			providerID:         "qwen",
			expectedDirPath:    ".qwen/commands/spectr",
			expectedConfigPath: "QWEN.md",
			expectedSlashPath:  ".qwen/commands/spectr",
		},
		{
			providerID:         "qoder",
			expectedDirPath:    ".qoder/commands/spectr",
			expectedConfigPath: "QODER.md",
			expectedSlashPath:  ".qoder/commands/spectr",
		},
		{
			providerID:         "antigravity",
			expectedDirPath:    ".antigravity/commands/spectr",
			expectedConfigPath: "ANTIGRAVITY.md",
			expectedSlashPath:  ".antigravity/commands/spectr",
		},
		{
			providerID:         "opencode",
			expectedDirPath:    ".opencode/command/spectr",
			expectedConfigPath: "",
			expectedSlashPath:  ".opencode/command/spectr",
		},
		{
			providerID:         "crush",
			expectedDirPath:    ".crush/commands/spectr",
			expectedConfigPath: "CRUSH.md",
			expectedSlashPath:  ".crush/commands/spectr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.providerID, func(t *testing.T) {
			reg := Get(tt.providerID)
			if reg == nil {
				t.Fatalf("Provider %s not found in registry", tt.providerID)
			}

			ctx := context.Background()
			initializers := reg.Provider.Initializers(ctx)

			var foundDirPath, foundConfigPath, foundSlashPath bool

			for _, init := range initializers {
				path := init.Path()
				switch init.(type) {
				case *DirectoryInitializer:
					if path == tt.expectedDirPath {
						foundDirPath = true
					} else {
						t.Errorf("Provider %s DirectoryInitializer has path %q, want %q", tt.providerID, path, tt.expectedDirPath)
					}
				case *ConfigFileInitializer:
					if path == tt.expectedConfigPath {
						foundConfigPath = true
					} else {
						t.Errorf("Provider %s ConfigFileInitializer has path %q, want %q", tt.providerID, path, tt.expectedConfigPath)
					}
				case *SlashCommandsInitializer:
					if path == tt.expectedSlashPath {
						foundSlashPath = true
					} else {
						t.Errorf("Provider %s SlashCommandsInitializer has path %q, want %q", tt.providerID, path, tt.expectedSlashPath)
					}
				}
			}

			// Verify we found the expected initializers
			if !foundDirPath {
				t.Errorf(
					"Provider %s missing DirectoryInitializer with path %q",
					tt.providerID,
					tt.expectedDirPath,
				)
			}
			if tt.expectedConfigPath != "" && !foundConfigPath {
				t.Errorf(
					"Provider %s missing ConfigFileInitializer with path %q",
					tt.providerID,
					tt.expectedConfigPath,
				)
			}
			if !foundSlashPath {
				t.Errorf(
					"Provider %s missing SlashCommandsInitializer with path %q",
					tt.providerID,
					tt.expectedSlashPath,
				)
			}
		})
	}
}

// TestProviderSlashCommandFormats verifies that each provider uses the expected command format.
func TestProviderSlashCommandFormats(t *testing.T) {
	tests := []struct {
		providerID     string
		expectedFormat CommandFormat
		expectedExt    string
	}{
		// Gemini uses TOML format
		{providerID: "gemini", expectedFormat: FormatTOML, expectedExt: ".toml"},

		// All others use Markdown format
		{providerID: "claude-code", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "cline", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "cursor", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "aider", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "continue", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "windsurf", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "costrict", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "kilocode", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "codex", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "qwen", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "qoder", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "antigravity", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "opencode", expectedFormat: FormatMarkdown, expectedExt: ".md"},
		{providerID: "crush", expectedFormat: FormatMarkdown, expectedExt: ".md"},
	}

	for _, tt := range tests {
		t.Run(tt.providerID, func(t *testing.T) {
			reg := Get(tt.providerID)
			if reg == nil {
				t.Fatalf("Provider %s not found in registry", tt.providerID)
			}

			ctx := context.Background()
			initializers := reg.Provider.Initializers(ctx)

			var foundSlashCmd *SlashCommandsInitializer
			for _, init := range initializers {
				if slashCmd, ok := init.(*SlashCommandsInitializer); ok {
					foundSlashCmd = slashCmd

					break
				}
			}

			if foundSlashCmd == nil {
				t.Fatalf("Provider %s has no SlashCommandsInitializer", tt.providerID)
			}

			// Verify format
			if foundSlashCmd.Format != tt.expectedFormat {
				t.Errorf(
					"Provider %s slash command format = %v, want %v",
					tt.providerID,
					foundSlashCmd.Format,
					tt.expectedFormat,
				)
			}

			// Verify extension
			if foundSlashCmd.Ext != tt.expectedExt {
				t.Errorf(
					"Provider %s slash command extension = %q, want %q",
					tt.providerID,
					foundSlashCmd.Ext,
					tt.expectedExt,
				)
			}
		})
	}
}

// TestProviderRegistrationMetadata verifies that each provider has valid registration metadata.
func TestProviderRegistrationMetadata(t *testing.T) {
	tests := []struct {
		providerID       string
		expectedName     string
		expectedPriority int
	}{
		{
			providerID:       "claude-code",
			expectedName:     "Claude Code",
			expectedPriority: PriorityClaudeCode,
		},
		{providerID: "gemini", expectedName: "Gemini CLI", expectedPriority: PriorityGemini},
		{providerID: "costrict", expectedName: "Costrict", expectedPriority: PriorityCostrict},
		{providerID: "qoder", expectedName: "Qoder", expectedPriority: PriorityQoder},
		{providerID: "qwen", expectedName: "Qwen", expectedPriority: PriorityQwen},
		{
			providerID:       "antigravity",
			expectedName:     "Antigravity",
			expectedPriority: PriorityAntigravity,
		},
		{providerID: "cline", expectedName: "Cline", expectedPriority: PriorityCline},
		{providerID: "cursor", expectedName: "Cursor", expectedPriority: PriorityCursor},
		{providerID: "codex", expectedName: "Codex", expectedPriority: PriorityCodex},
		{providerID: "opencode", expectedName: "OpenCode", expectedPriority: PriorityOpencode},
		{providerID: "aider", expectedName: "Aider", expectedPriority: PriorityAider},
		{providerID: "windsurf", expectedName: "Windsurf", expectedPriority: PriorityWindsurf},
		{providerID: "kilocode", expectedName: "Kilocode", expectedPriority: PriorityKilocode},
		{providerID: "continue", expectedName: "Continue", expectedPriority: PriorityContinue},
		{providerID: "crush", expectedName: "Crush", expectedPriority: PriorityCrush},
	}

	for _, tt := range tests {
		t.Run(tt.providerID, func(t *testing.T) {
			reg := Get(tt.providerID)
			if reg == nil {
				t.Fatalf("Provider %s not found in registry", tt.providerID)
			}

			// Verify ID is non-empty and matches expected
			if reg.ID == "" {
				t.Error("Provider ID is empty")
			}
			if reg.ID != tt.providerID {
				t.Errorf("Provider ID = %q, want %q", reg.ID, tt.providerID)
			}

			// Verify ID is kebab-case (lowercase with hyphens, no other special chars)
			for i, ch := range reg.ID {
				if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') && ch != '-' {
					t.Errorf(
						"Provider ID %q contains invalid character %q at position %d (must be lowercase alphanumeric or hyphen)",
						reg.ID,
						ch,
						i,
					)
				}
			}

			// Verify Name is non-empty and matches expected
			if reg.Name == "" {
				t.Error("Provider Name is empty")
			}
			if reg.Name != tt.expectedName {
				t.Errorf("Provider Name = %q, want %q", reg.Name, tt.expectedName)
			}

			// Verify Priority is in expected range (1-20 is reasonable)
			if reg.Priority < 1 || reg.Priority > 20 {
				t.Errorf("Provider Priority = %d, want value between 1 and 20", reg.Priority)
			}
			if reg.Priority != tt.expectedPriority {
				t.Errorf("Provider Priority = %d, want %d", reg.Priority, tt.expectedPriority)
			}

			// Verify Provider is non-nil
			if reg.Provider == nil {
				t.Error("Provider instance is nil")
			}
		})
	}
}

// TestAllProvidersHaveUniqueIDs verifies that all provider IDs are unique.
func TestAllProvidersHaveUniqueIDs(t *testing.T) {
	regs := All()
	seen := make(map[string]bool)

	for _, reg := range regs {
		if seen[reg.ID] {
			t.Errorf("Duplicate provider ID found: %s", reg.ID)
		}
		seen[reg.ID] = true
	}
}

// TestProviderPrioritySorting verifies that providers are returned in priority order.
func TestProviderPrioritySorting(t *testing.T) {
	regs := All()

	if len(regs) < 2 {
		t.Skip("Need at least 2 providers to test sorting")
	}

	// Verify they are sorted by priority (ascending)
	for i := range len(regs) - 1 {
		if regs[i].Priority > regs[i+1].Priority {
			t.Errorf(
				"Providers not sorted by priority: %s (priority %d) comes before %s (priority %d)",
				regs[i].ID,
				regs[i].Priority,
				regs[i+1].ID,
				regs[i+1].Priority,
			)
		}
	}
}
