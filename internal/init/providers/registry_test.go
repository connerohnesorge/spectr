package providers

import (
	"testing"
)

func TestGlobalRegistry(t *testing.T) {
	// All providers are registered via init()
	allProviders := All()

	// Should have all 17 providers registered
	expectedCount := 17
	if len(allProviders) != expectedCount {
		t.Errorf("Expected %d providers, got %d", expectedCount, len(allProviders))
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		wantNil  bool
		wantName string
	}{
		{"Get Claude Code", "claude-code", false, "Claude Code"},
		{"Get Cline", "cline", false, "Cline"},
		{"Get Gemini", "gemini", false, "Gemini CLI"},
		{"Get Cursor", "cursor", false, "Cursor"},
		{"Get Invalid", "nonexistent", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Get(tt.id)
			if tt.wantNil {
				if provider != nil {
					t.Errorf("Get(%s) = %v, want nil", tt.id, provider)
				}
			} else {
				if provider == nil {
					t.Errorf("Get(%s) = nil, want non-nil", tt.id)

					return
				}
				if provider.Name() != tt.wantName {
					t.Errorf("Get(%s).Name() = %s, want %s", tt.id, provider.Name(), tt.wantName)
				}
			}
		})
	}
}

func TestAllSortedByPriority(t *testing.T) {
	allProviders := All()

	// Verify they're sorted by priority
	for i := 1; i < len(allProviders); i++ {
		if allProviders[i-1].Priority() > allProviders[i].Priority() {
			t.Errorf(
				"Providers not sorted by priority: %s (priority %d) comes before %s (priority %d)",
				allProviders[i-1].ID(),
				allProviders[i-1].Priority(),
				allProviders[i].ID(),
				allProviders[i].Priority(),
			)
		}
	}
}

func TestIDs(t *testing.T) {
	ids := IDs()

	// Verify all IDs are non-empty and unique
	seen := make(map[string]bool)
	for _, id := range ids {
		if id == "" {
			t.Error("Found empty ID")
		}
		if seen[id] {
			t.Errorf("Duplicate ID: %s", id)
		}
		seen[id] = true
	}
}

func TestWithConfigFile(t *testing.T) {
	providers := WithConfigFile()

	// All returned providers should have a config file
	for _, p := range providers {
		if !p.HasConfigFile() {
			t.Errorf("Provider %s returned by WithConfigFile() but HasConfigFile() = false", p.ID())
		}
		if p.ConfigFile() == "" {
			t.Errorf("Provider %s has empty ConfigFile()", p.ID())
		}
	}
}

func TestWithSlashCommands(t *testing.T) {
	providers := WithSlashCommands()

	// All returned providers should have slash commands
	for _, p := range providers {
		if !p.HasSlashCommands() {
			t.Errorf(
				"Provider %s returned by WithSlashCommands() but HasSlashCommands() = false",
				p.ID(),
			)
		}
		if p.SlashDir() == "" {
			t.Errorf("Provider %s has empty SlashDir()", p.ID())
		}
	}

	// All 17 providers have slash commands
	expectedCount := 17
	if len(providers) != expectedCount {
		t.Errorf("Expected %d providers with slash commands, got %d", expectedCount, len(providers))
	}
}

func TestProviderIDsAreKebabCase(t *testing.T) {
	allProviders := All()

	for _, p := range allProviders {
		id := p.ID()
		for _, char := range id {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				t.Errorf(
					"Provider ID %s is not in kebab-case (contains invalid character: %c)",
					id,
					char,
				)
			}
		}
	}
}

func TestInstanceRegistry(t *testing.T) {
	r := NewRegistry()

	// Initially empty
	if r.Count() != 0 {
		t.Errorf("New registry should be empty, got %d providers", r.Count())
	}

	// Register a provider
	err := r.Register(NewClaudeProvider())
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	if r.Count() != 1 {
		t.Errorf("Expected 1 provider after registration, got %d", r.Count())
	}

	// Get the provider
	p := r.Get("claude-code")
	if p == nil {
		t.Error("Get returned nil for registered provider")
	}

	// Duplicate registration should fail
	err = r.Register(NewClaudeProvider())
	if err == nil {
		t.Error("Duplicate registration should fail")
	}
}

func TestNewRegistryFromGlobal(t *testing.T) {
	r := NewRegistryFromGlobal()

	// Should have same count as global
	globalCount := Count()
	if r.Count() != globalCount {
		t.Errorf("NewRegistryFromGlobal() has %d providers, global has %d", r.Count(), globalCount)
	}

	// Should be able to get same providers
	p := r.Get("claude-code")
	if p == nil {
		t.Error("Get returned nil for claude-code")
	}
}
