package providers

import (
	"context"
	"testing"
)

// -----------------------------------------------------------------------------
// Tests for Provider interface (Task 8.2)
// -----------------------------------------------------------------------------

func TestProviderInterface(t *testing.T) {
	t.Run("Provider interface is implemented by all registered providers", func(t *testing.T) {
		// Get all registered providers
		registrations := All()

		if len(registrations) == 0 {
			t.Skip("No providers registered - init() functions may not have run")
		}

		for _, reg := range registrations {
			// Verify Provider field is not nil (it already implements Provider
			// since reg.Provider is typed as Provider in the Registration struct)
			if reg.Provider == nil {
				t.Errorf("provider %s has nil Provider", reg.ID)
			}
		}
	})

	t.Run("Provider.Initializers returns valid slice", func(t *testing.T) {
		registrations := All()

		if len(registrations) == 0 {
			t.Skip("No providers registered")
		}

		ctx := context.Background()

		for _, reg := range registrations {
			initializers := reg.Provider.Initializers(ctx)

			// Initializers should not panic and should return a slice
			// (can be empty, but not nil in practice for our providers)
			if initializers == nil {
				// nil slice is technically valid but we log it
				t.Logf("provider %s returns nil initializers slice", reg.ID)
			}
		}
	})
}

// -----------------------------------------------------------------------------
// Tests for all 17 providers returning expected initializers (Task 8.3)
// -----------------------------------------------------------------------------

// ExpectedProviderConfig describes what initializers a provider should return.
type ExpectedProviderConfig struct {
	ID                    string
	HasDirectoryInit      bool
	HasConfigFileInit     bool
	HasSlashCommandsInit  bool
	ExpectedInitCount     int  // Minimum expected initializer count
	HasGlobalInitializers bool // Whether any initializers use global paths
}

// getExpectedProviderConfigs returns the expected configuration for all 15 providers.
// Note: Some providers don't have config files (instruction files) - they only use slash commands.
func getExpectedProviderConfigs() []ExpectedProviderConfig {
	return []ExpectedProviderConfig{
		// Claude Code: Directory + ConfigFile (CLAUDE.md) + SlashCommands
		{
			ID:                    "claude-code",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Gemini: Directory + SlashCommands (TOML format, no ConfigFile)
		{
			ID:                    "gemini",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},

		// Cursor: Directory + SlashCommands (no ConfigFile)
		{
			ID:                    "cursor",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},

		// Cline: Directory + ConfigFile + SlashCommands
		{
			ID:                    "cline",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Aider: Directory + SlashCommands (no ConfigFile)
		{
			ID:                    "aider",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},

		// Codex: Directory (global) + ConfigFile (AGENTS.md) + SlashCommands (global)
		{
			ID:                    "codex",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: true,
		},

		// Costrict: Directory + ConfigFile + SlashCommands
		{
			ID:                    "costrict",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Qoder: Directory + ConfigFile + SlashCommands
		{
			ID:                    "qoder",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Qwen: Directory + ConfigFile + SlashCommands
		{
			ID:                    "qwen",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Antigravity: Directory + ConfigFile + SlashCommands
		{
			ID:                    "antigravity",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Windsurf: Directory + SlashCommands (no ConfigFile)
		{
			ID:                    "windsurf",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},

		// Kilocode: Directory + SlashCommands (no ConfigFile)
		{
			ID:                    "kilocode",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},

		// Continue: Directory + SlashCommands (no ConfigFile)
		{
			ID:                    "continue",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},

		// Crush: Directory + ConfigFile + SlashCommands
		{
			ID:                    "crush",
			HasDirectoryInit:      true,
			HasConfigFileInit:     true,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     3,
			HasGlobalInitializers: false,
		},

		// Opencode: Directory + SlashCommands (no ConfigFile - uses JSON config)
		{
			ID:                    "opencode",
			HasDirectoryInit:      true,
			HasConfigFileInit:     false,
			HasSlashCommandsInit:  true,
			ExpectedInitCount:     2,
			HasGlobalInitializers: false,
		},
	}
}

func TestAllProvidersReturnExpectedInitializers(t *testing.T) {
	expectedConfigs := getExpectedProviderConfigs()
	ctx := context.Background()

	for _, expected := range expectedConfigs {
		t.Run(expected.ID, func(t *testing.T) {
			reg, found := Get(expected.ID)
			if !found {
				t.Skipf("provider %s not registered", expected.ID)

				return
			}

			initializers := reg.Provider.Initializers(ctx)

			// Check minimum initializer count
			if len(initializers) < expected.ExpectedInitCount {
				t.Errorf("provider %s: expected at least %d initializers, got %d",
					expected.ID, expected.ExpectedInitCount, len(initializers))
			}

			// Categorize initializers by type
			var hasDirectory, hasConfigFile, hasSlashCommands, hasGlobal bool
			for _, init := range initializers {
				if init == nil {
					continue
				}

				switch init.(type) {
				case *DirectoryInitializerBuiltin:
					hasDirectory = true
				case *ConfigFileInitializerBuiltin:
					hasConfigFile = true
				case *SlashCommandsInitializerBuiltin:
					hasSlashCommands = true
				}

				if init.IsGlobal() {
					hasGlobal = true
				}
			}

			// Verify expected initializer types
			if expected.HasDirectoryInit && !hasDirectory {
				t.Errorf("provider %s: expected DirectoryInitializer but none found", expected.ID)
			}
			if expected.HasConfigFileInit && !hasConfigFile {
				t.Errorf("provider %s: expected ConfigFileInitializer but none found", expected.ID)
			}
			if expected.HasSlashCommandsInit && !hasSlashCommands {
				t.Errorf(
					"provider %s: expected SlashCommandsInitializer but none found",
					expected.ID,
				)
			}
			if expected.HasGlobalInitializers && !hasGlobal {
				t.Errorf("provider %s: expected global initializers but none found", expected.ID)
			}
		})
	}
}

func TestProviderInitializerPaths(t *testing.T) {
	// Test that initializers have valid, non-empty paths
	ctx := context.Background()
	registrations := All()

	for _, reg := range registrations {
		t.Run(reg.ID+"_paths", func(t *testing.T) {
			initializers := reg.Provider.Initializers(ctx)

			for i, init := range initializers {
				if init == nil {
					t.Errorf("provider %s: initializer at index %d is nil", reg.ID, i)

					continue
				}

				path := init.Path()
				if path == "" {
					t.Errorf("provider %s: initializer at index %d has empty path", reg.ID, i)
				}
			}
		})
	}
}

func TestProviderInitializerCount(t *testing.T) {
	// Verify the total count of registered providers
	count := Count()

	// We expect at least 15 providers based on the constants
	minExpectedProviders := 15
	if count < minExpectedProviders {
		t.Errorf("expected at least %d providers, got %d", minExpectedProviders, count)
	}

	t.Logf("Total registered providers: %d", count)
}

// -----------------------------------------------------------------------------
// Tests for provider registration metadata (Task 8.4)
// -----------------------------------------------------------------------------

func TestProviderRegistrationMetadata(t *testing.T) {
	registrations := All()

	t.Run("all providers have valid IDs", func(t *testing.T) {
		for _, reg := range registrations {
			if reg.ID == "" {
				t.Error("found provider with empty ID")
			}
			// ID should be kebab-case
			if !kebabCaseRegex.MatchString(reg.ID) {
				t.Errorf("provider ID %q is not valid kebab-case", reg.ID)
			}
		}
	})

	t.Run("all providers have valid Names", func(t *testing.T) {
		for _, reg := range registrations {
			if reg.Name == "" {
				t.Errorf("provider %s has empty Name", reg.ID)
			}
		}
	})

	t.Run("all providers have non-negative Priority", func(t *testing.T) {
		for _, reg := range registrations {
			if reg.Priority < 0 {
				t.Errorf("provider %s has negative Priority: %d", reg.ID, reg.Priority)
			}
		}
	})

	t.Run("no duplicate IDs", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, reg := range registrations {
			if seen[reg.ID] {
				t.Errorf("duplicate provider ID: %s", reg.ID)
			}
			seen[reg.ID] = true
		}
	})

	t.Run("no duplicate Priorities", func(t *testing.T) {
		// Note: This is a warning, not an error - duplicate priorities are allowed
		// but they should be avoided for deterministic ordering
		priorities := make(map[int][]string)
		for _, reg := range registrations {
			priorities[reg.Priority] = append(priorities[reg.Priority], reg.ID)
		}

		for priority, ids := range priorities {
			if len(ids) > 1 {
				t.Logf("Warning: multiple providers share priority %d: %v", priority, ids)
			}
		}
	})
}

// ExpectedMetadata describes the expected metadata for a specific provider.
type ExpectedMetadata struct {
	ID       string
	Name     string
	Priority int
}

func TestSpecificProviderMetadata(t *testing.T) {
	// Test specific expected metadata for known providers
	// Note: Names should match the actual registration names in each provider file
	expectedMetadata := []ExpectedMetadata{
		{ID: "claude-code", Name: "Claude Code", Priority: PriorityClaudeCode},
		{ID: "gemini", Name: "Gemini CLI", Priority: PriorityGemini},
		{ID: "costrict", Name: "CoStrict", Priority: PriorityCostrict},
		{ID: "qoder", Name: "Qoder", Priority: PriorityQoder},
		{ID: "qwen", Name: "Qwen Code", Priority: PriorityQwen},
		{ID: "antigravity", Name: "Antigravity", Priority: PriorityAntigravity},
		{ID: "cline", Name: "Cline", Priority: PriorityCline},
		{ID: "cursor", Name: "Cursor", Priority: PriorityCursor},
		{ID: "codex", Name: "Codex CLI", Priority: PriorityCodex},
		{ID: "opencode", Name: "OpenCode", Priority: PriorityOpencode},
		{ID: "aider", Name: "Aider", Priority: PriorityAider},
		{ID: "windsurf", Name: "Windsurf", Priority: PriorityWindsurf},
		{ID: "kilocode", Name: "Kilocode", Priority: PriorityKilocode},
		{ID: "continue", Name: "Continue", Priority: PriorityContinue},
		{ID: "crush", Name: "Crush", Priority: PriorityCrush},
	}

	for _, expected := range expectedMetadata {
		t.Run(expected.ID, func(t *testing.T) {
			reg, found := Get(expected.ID)
			if !found {
				t.Skipf("provider %s not registered", expected.ID)

				return
			}

			if reg.Name != expected.Name {
				t.Errorf("provider %s: expected Name %q, got %q",
					expected.ID, expected.Name, reg.Name)
			}

			if reg.Priority != expected.Priority {
				t.Errorf("provider %s: expected Priority %d, got %d",
					expected.ID, expected.Priority, reg.Priority)
			}
		})
	}
}

func TestProviderPriorityOrdering(t *testing.T) {
	// Verify that All() returns providers in priority order
	registrations := All()

	if len(registrations) < 2 {
		t.Skip("Not enough providers to test ordering")
	}

	for i := 1; i < len(registrations); i++ {
		prev := registrations[i-1]
		curr := registrations[i]

		// Priority should be ascending (or equal with ID tiebreaker)
		if curr.Priority < prev.Priority {
			t.Errorf(
				"providers not in priority order: %s (priority %d) comes after %s (priority %d)",
				curr.ID,
				curr.Priority,
				prev.ID,
				prev.Priority,
			)
		}

		// If priorities are equal, IDs should be alphabetically ordered
		if curr.Priority == prev.Priority && curr.ID < prev.ID {
			t.Errorf("providers with same priority not in ID order: %s should come before %s",
				prev.ID, curr.ID)
		}
	}
}

// -----------------------------------------------------------------------------
// Helper functions for testing
// -----------------------------------------------------------------------------

// countInitializersByType counts initializers by their concrete type.
func countInitializersByType(initializers []Initializer) (dirs, configs, slashCmds int) {
	for _, init := range initializers {
		if init == nil {
			continue
		}
		switch init.(type) {
		case *DirectoryInitializerBuiltin:
			dirs++
		case *ConfigFileInitializerBuiltin:
			configs++
		case *SlashCommandsInitializerBuiltin:
			slashCmds++
		}
	}

	return dirs, configs, slashCmds
}

func TestInitializerTypesCounting(t *testing.T) {
	// Test that our counting helper works correctly
	ctx := context.Background()

	reg, found := Get("claude-code")
	if !found {
		t.Skip("claude-code provider not registered")
	}

	initializers := reg.Provider.Initializers(ctx)
	dirs, configs, slashCmds := countInitializersByType(initializers)

	if dirs != 1 {
		t.Errorf("expected 1 directory initializer, got %d", dirs)
	}
	if configs != 1 {
		t.Errorf("expected 1 config file initializer, got %d", configs)
	}
	if slashCmds != 1 {
		t.Errorf("expected 1 slash commands initializer, got %d", slashCmds)
	}
}

// TestInitializersImplementInterface verifies all returned initializers
// properly implement the Initializer interface.
func TestInitializersImplementInterface(t *testing.T) {
	ctx := context.Background()
	registrations := All()

	for _, reg := range registrations {
		t.Run(reg.ID, func(t *testing.T) {
			initializers := reg.Provider.Initializers(ctx)

			for i, init := range initializers {
				if init == nil {
					continue
				}

				// Verify all interface methods are callable
				path := init.Path()
				isGlobal := init.IsGlobal()

				// These should not panic
				_ = path
				_ = isGlobal

				t.Logf("provider %s initializer %d: path=%s, isGlobal=%v",
					reg.ID, i, path, isGlobal)
			}
		})
	}
}
