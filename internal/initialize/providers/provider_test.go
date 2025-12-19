package providers

import (
	"context"
	"testing"
)

// TestProviderInterface verifies that Provider interface works correctly.
func TestProviderInterface(t *testing.T) {
	t.Run("ProviderFunc implements Provider", func(t *testing.T) {
		var p Provider = ProviderFunc(func(ctx context.Context) []Initializer {
			return nil
		})
		if p == nil {
			t.Error("ProviderFunc should implement Provider interface")
		}
	})

	t.Run("ProviderFunc returns initializers", func(t *testing.T) {
		expected := []Initializer{
			NewDirectoryInitializer("test"),
		}
		p := ProviderFunc(func(ctx context.Context) []Initializer {
			return expected
		})

		result := p.Initializers(context.Background())
		if len(result) != len(expected) {
			t.Errorf("Initializers() returned %d items, want %d", len(result), len(expected))
		}
	})
}

// TestAllProvidersImplementProvider verifies all registered providers implement Provider.
func TestAllProvidersImplementProvider(t *testing.T) {
	// Get all providers from the global registry
	allRegs := All()

	if len(allRegs) == 0 {
		t.Skip("No providers registered - run tests with 'go test' to trigger init()")
	}

	for _, reg := range allRegs {
		t.Run(reg.ID, func(t *testing.T) {
			// Check that Provider implements Provider
			if reg.Provider == nil {
				t.Errorf("Provider %q has nil Provider", reg.ID)
				return
			}

			// Verify Initializers method works
			ctx := context.Background()
			inits := reg.Provider.Initializers(ctx)

			// Every provider should return at least one initializer
			if len(inits) == 0 {
				t.Errorf("Provider %q returned 0 initializers, expected at least 1", reg.ID)
			}

			// Verify each initializer is not nil
			for i, init := range inits {
				if init == nil {
					t.Errorf("Provider %q returned nil initializer at index %d", reg.ID, i)
				}
			}
		})
	}
}

// ExpectedProviderInitializers defines the expected initializers for each provider.
type ExpectedProviderInitializers struct {
	ID                 string
	Name               string
	Priority           int
	HasDirectory       bool
	HasConfigFile      bool
	HasSlashCommands   bool
	DirectoryPath      string
	ConfigFilePath     string
	SlashCommandsPath  string
	SlashCommandFormat CommandFormat
	IsGlobal           bool
}

// GetExpectedProviders returns the expected configuration for all providers.
// This is based on the spec and actual provider implementations.
func GetExpectedProviders() []ExpectedProviderInitializers {
	return []ExpectedProviderInitializers{
		{
			ID:                 "claude-code",
			Name:               "Claude Code",
			Priority:           PriorityClaudeCode,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".claude/commands/spectr",
			ConfigFilePath:     "CLAUDE.md",
			SlashCommandsPath:  ".claude/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "gemini",
			Name:               "Gemini CLI",
			Priority:           PriorityGemini,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".gemini/commands/spectr",
			SlashCommandsPath:  ".gemini/commands/spectr",
			SlashCommandFormat: FormatTOML,
		},
		{
			ID:                 "costrict",
			Name:               "CoStrict",
			Priority:           PriorityCostrict,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".costrict/commands/spectr",
			ConfigFilePath:     "COSTRICT.md",
			SlashCommandsPath:  ".costrict/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "qoder",
			Name:               "Qoder",
			Priority:           PriorityQoder,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".qoder/commands/spectr",
			ConfigFilePath:     "QODER.md",
			SlashCommandsPath:  ".qoder/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "qwen",
			Name:               "Qwen Code",
			Priority:           PriorityQwen,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".qwen/commands/spectr",
			ConfigFilePath:     "QWEN.md",
			SlashCommandsPath:  ".qwen/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "antigravity",
			Name:               "Antigravity",
			Priority:           PriorityAntigravity,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".agent/workflows",
			ConfigFilePath:     "AGENTS.md",
			SlashCommandsPath:  ".agent/workflows",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "cline",
			Name:               "Cline",
			Priority:           PriorityCline,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".clinerules/commands/spectr",
			ConfigFilePath:     "CLINE.md",
			SlashCommandsPath:  ".clinerules/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "cursor",
			Name:               "Cursor",
			Priority:           PriorityCursor,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".cursorrules/commands/spectr",
			SlashCommandsPath:  ".cursorrules/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "codex",
			Name:               "Codex CLI",
			Priority:           PriorityCodex,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".codex/prompts",
			ConfigFilePath:     "AGENTS.md",
			SlashCommandsPath:  ".codex/prompts",
			SlashCommandFormat: FormatMarkdown,
			IsGlobal:           true, // Codex uses global directory
		},
		{
			ID:                 "opencode",
			Name:               "OpenCode",
			Priority:           PriorityOpencode,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".opencode/command/spectr",
			SlashCommandsPath:  ".opencode/command/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "aider",
			Name:               "Aider",
			Priority:           PriorityAider,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".aider/commands/spectr",
			SlashCommandsPath:  ".aider/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "windsurf",
			Name:               "Windsurf",
			Priority:           PriorityWindsurf,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".windsurf/commands/spectr",
			SlashCommandsPath:  ".windsurf/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "kilocode",
			Name:               "Kilocode",
			Priority:           PriorityKilocode,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".kilocode/commands/spectr",
			SlashCommandsPath:  ".kilocode/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "continue",
			Name:               "Continue",
			Priority:           PriorityContinue,
			HasDirectory:       true,
			HasConfigFile:      false,
			HasSlashCommands:   true,
			DirectoryPath:      ".continue/commands/spectr",
			SlashCommandsPath:  ".continue/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
		{
			ID:                 "crush",
			Name:               "Crush",
			Priority:           PriorityCrush,
			HasDirectory:       true,
			HasConfigFile:      true,
			HasSlashCommands:   true,
			DirectoryPath:      ".crush/commands/spectr",
			ConfigFilePath:     "CRUSH.md",
			SlashCommandsPath:  ".crush/commands/spectr",
			SlashCommandFormat: FormatMarkdown,
		},
	}
}

// TestAllProvidersReturnExpectedInitializers verifies all providers return expected initializers (Task 8.3).
func TestAllProvidersReturnExpectedInitializers(t *testing.T) {
	expectedProviders := GetExpectedProviders()
	ctx := context.Background()

	for _, expected := range expectedProviders {
		t.Run(expected.ID, func(t *testing.T) {
			reg := Get(expected.ID)
			if reg == nil {
				t.Fatalf("Provider %q not registered", expected.ID)
			}

			inits := reg.Provider.Initializers(ctx)

			// Check initializer count
			expectedCount := 0
			if expected.HasDirectory {
				expectedCount++
			}
			if expected.HasConfigFile {
				expectedCount++
			}
			if expected.HasSlashCommands {
				expectedCount++
			}

			if len(inits) != expectedCount {
				t.Errorf("Provider %q returned %d initializers, expected %d",
					expected.ID, len(inits), expectedCount)
			}

			// Verify initializer types and paths
			hasDirectory := false
			hasConfigFile := false
			hasSlashCommands := false

			for _, init := range inits {
				switch i := init.(type) {
				case *directoryInitializer:
					hasDirectory = true
					if expected.DirectoryPath != "" && i.Path() != expected.DirectoryPath {
						t.Errorf("Provider %q directory path = %q, want %q",
							expected.ID, i.Path(), expected.DirectoryPath)
					}
					// For codex, check if global is set correctly
					if expected.IsGlobal && !i.IsGlobal() {
						t.Errorf("Provider %q directory should be global", expected.ID)
					}
				case *configFileInitializer:
					hasConfigFile = true
					if expected.ConfigFilePath != "" && i.Path() != expected.ConfigFilePath {
						t.Errorf("Provider %q config file path = %q, want %q",
							expected.ID, i.Path(), expected.ConfigFilePath)
					}
				case *slashCommandsInitializer:
					hasSlashCommands = true
					if expected.SlashCommandsPath != "" && i.Path() != expected.SlashCommandsPath {
						t.Errorf("Provider %q slash commands path = %q, want %q",
							expected.ID, i.Path(), expected.SlashCommandsPath)
					}
				}
			}

			if expected.HasDirectory && !hasDirectory {
				t.Errorf("Provider %q missing DirectoryInitializer", expected.ID)
			}
			if expected.HasConfigFile && !hasConfigFile {
				t.Errorf("Provider %q missing ConfigFileInitializer", expected.ID)
			}
			if expected.HasSlashCommands && !hasSlashCommands {
				t.Errorf("Provider %q missing SlashCommandsInitializer", expected.ID)
			}
		})
	}
}

// TestProviderRegistrationMetadata verifies provider registration metadata (Task 8.4).
func TestProviderRegistrationMetadata(t *testing.T) {
	allRegs := All()

	if len(allRegs) == 0 {
		t.Skip("No providers registered")
	}

	t.Run("all providers have valid metadata", func(t *testing.T) {
		for _, reg := range allRegs {
			// ID must not be empty
			if reg.ID == "" {
				t.Error("Registration has empty ID")
			}

			// Name must not be empty
			if reg.Name == "" {
				t.Errorf("Registration %q has empty Name", reg.ID)
			}

			// Provider must not be nil
			if reg.Provider == nil {
				t.Errorf("Registration %q has nil Provider", reg.ID)
			}
		}
	})

	t.Run("provider IDs are unique", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, reg := range allRegs {
			if seen[reg.ID] {
				t.Errorf("Duplicate provider ID: %q", reg.ID)
			}
			seen[reg.ID] = true
		}
	})

	t.Run("providers are sorted by priority", func(t *testing.T) {
		for i := 1; i < len(allRegs); i++ {
			if allRegs[i-1].Priority > allRegs[i].Priority {
				t.Errorf(
					"Providers not sorted by priority: %q (priority %d) comes before %q (priority %d)",
					allRegs[i-1].ID,
					allRegs[i-1].Priority,
					allRegs[i].ID,
					allRegs[i].Priority,
				)
			}
		}
	})

	t.Run("expected providers are registered", func(t *testing.T) {
		expectedProviders := GetExpectedProviders()
		for _, expected := range expectedProviders {
			reg := Get(expected.ID)
			if reg == nil {
				t.Errorf("Expected provider %q not registered", expected.ID)
				continue
			}

			if reg.Name != expected.Name {
				t.Errorf("Provider %q Name = %q, want %q",
					expected.ID, reg.Name, expected.Name)
			}

			if reg.Priority != expected.Priority {
				t.Errorf("Provider %q Priority = %d, want %d",
					expected.ID, reg.Priority, expected.Priority)
			}
		}
	})
}

// TestProviderCount verifies the expected number of providers are registered.
func TestProviderCount(t *testing.T) {
	count := Count()
	expectedProviders := GetExpectedProviders()
	expectedCount := len(expectedProviders)

	if count != expectedCount {
		t.Errorf("Count() = %d, want %d", count, expectedCount)
		t.Log("Registered providers:")
		for _, reg := range All() {
			t.Logf("  - %s (%s)", reg.ID, reg.Name)
		}
	}
}

// TestRegistrationValidate tests the Validate method of Registration.
func TestRegistrationValidate(t *testing.T) {
	tests := []struct {
		name    string
		reg     Registration
		wantErr error
	}{
		{
			name: "valid registration",
			reg: Registration{
				ID:       "test-provider",
				Name:     "Test Provider",
				Priority: 1,
				Provider: ProviderFunc(func(ctx context.Context) []Initializer {
					return nil
				}),
			},
			wantErr: nil,
		},
		{
			name: "empty ID",
			reg: Registration{
				ID:       "",
				Name:     "Test Provider",
				Provider: ProviderFunc(func(ctx context.Context) []Initializer { return nil }),
			},
			wantErr: ErrEmptyID,
		},
		{
			name: "empty Name",
			reg: Registration{
				ID:       "test",
				Name:     "",
				Provider: ProviderFunc(func(ctx context.Context) []Initializer { return nil }),
			},
			wantErr: ErrEmptyName,
		},
		{
			name: "nil Provider",
			reg: Registration{
				ID:       "test",
				Name:     "Test",
				Provider: nil,
			},
			wantErr: ErrNilProvider,
		},
		{
			name: "zero priority is valid",
			reg: Registration{
				ID:       "test-zero",
				Name:     "Test Zero",
				Priority: 0,
				Provider: ProviderFunc(func(ctx context.Context) []Initializer { return nil }),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.reg.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestInitializerPath verifies initializers return correct paths.
func TestInitializerPath(t *testing.T) {
	tests := []struct {
		name     string
		init     Initializer
		wantPath string
	}{
		{
			name:     "directory initializer",
			init:     NewDirectoryInitializer(".test/dir"),
			wantPath: ".test/dir",
		},
		{
			name:     "config file initializer",
			init:     NewConfigFileInitializer("TEST.md"),
			wantPath: "TEST.md",
		},
		{
			name:     "slash commands initializer",
			init:     NewSlashCommandsInitializer(".test/commands", ".md", FormatMarkdown),
			wantPath: ".test/commands",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.init.Path()
			if path != tt.wantPath {
				t.Errorf("Path() = %q, want %q", path, tt.wantPath)
			}
		})
	}
}

// TestInitializerIsGlobal verifies the IsGlobal flag is set correctly.
func TestInitializerIsGlobal(t *testing.T) {
	tests := []struct {
		name       string
		init       Initializer
		wantGlobal bool
	}{
		{
			name:       "project directory initializer",
			init:       NewDirectoryInitializer(".test/dir"),
			wantGlobal: false,
		},
		{
			name:       "global directory initializer",
			init:       NewGlobalDirectoryInitializer(".config/test"),
			wantGlobal: true,
		},
		{
			name:       "project config file initializer",
			init:       NewConfigFileInitializer("TEST.md"),
			wantGlobal: false,
		},
		{
			name:       "global config file initializer",
			init:       NewGlobalConfigFileInitializer(".config/test.md"),
			wantGlobal: true,
		},
		{
			name:       "project slash commands initializer",
			init:       NewSlashCommandsInitializer(".test/commands", ".md", FormatMarkdown),
			wantGlobal: false,
		},
		{
			name: "global slash commands initializer",
			init: NewGlobalSlashCommandsInitializer(
				".config/commands",
				".md",
				FormatMarkdown,
			),
			wantGlobal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isGlobal := tt.init.IsGlobal()
			if isGlobal != tt.wantGlobal {
				t.Errorf("IsGlobal() = %v, want %v", isGlobal, tt.wantGlobal)
			}
		})
	}
}
