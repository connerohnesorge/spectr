package providers

import (
	"context"
	"strings"
	"testing"
)

// testProvider is a simple provider implementation for testing.
type testProvider struct{}

// Initializers returns an empty list of initializers for testing.
func (*testProvider) Initializers(_ context.Context, _ any) []Initializer {
	return nil
}

// IsConfigured returns false for testing.
func (*testProvider) IsConfigured(_ string) bool {
	return false
}

// TestRegisterProvider_Valid tests that valid registrations succeed.
func TestRegisterProvider_Valid(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RegisterProvider(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &testProvider{},
	})

	if err != nil {
		t.Errorf("RegisterProvider() returned error for valid registration: %v", err)
	}

	if Count() != 1 {
		t.Errorf("Count() = %d, want 1", Count())
	}
}

// TestRegisterProvider_EmptyID tests that empty ID is rejected.
func TestRegisterProvider_EmptyID(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RegisterProvider(Registration{
		ID:       "",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &testProvider{},
	})

	if err == nil {
		t.Error("RegisterProvider() should return error for empty ID")
	}
	if !strings.Contains(err.Error(), "ID is required") {
		t.Errorf("Error message should mention ID requirement, got: %v", err)
	}
}

// TestRegisterProvider_NilProvider tests that nil Provider is rejected.
func TestRegisterProvider_NilProvider(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RegisterProvider(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: nil,
	})

	if err == nil {
		t.Error("RegisterProvider() should return error for nil Provider")
	}
	if !strings.Contains(err.Error(), "implementation is required") {
		t.Errorf("Error message should mention provider implementation, got: %v", err)
	}
}

// TestRegisterProvider_DuplicateID tests that duplicate IDs are rejected with error.
func TestRegisterProvider_DuplicateID(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	// First registration should succeed
	err := RegisterProvider(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &testProvider{},
	})
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Second registration with same ID should fail
	err = RegisterProvider(Registration{
		ID:       "test-provider",
		Name:     "Test Provider 2",
		Priority: 2,
		Provider: &testProvider{},
	})

	if err == nil {
		t.Error("RegisterProvider() should return error for duplicate ID")
	}
	// Check the exact error format as specified in the task
	expectedMsg := `provider "test-provider" already registered`
	if err.Error() != expectedMsg {
		t.Errorf("Error message = %q, want %q", err.Error(), expectedMsg)
	}
}

// TestRegisteredProviders_SortedByPriority tests that providers are returned sorted by priority.
func TestRegisteredProviders_SortedByPriority(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	// Register in non-priority order
	registrations := []Registration{
		{ID: "provider-c", Name: "Provider C", Priority: 10, Provider: &testProvider{}},
		{ID: "provider-a", Name: "Provider A", Priority: 1, Provider: &testProvider{}},
		{ID: "provider-b", Name: "Provider B", Priority: 5, Provider: &testProvider{}},
	}

	for _, reg := range registrations {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%s) failed: %v", reg.ID, err)
		}
	}

	providers := RegisteredProviders()

	if len(providers) != 3 {
		t.Fatalf("RegisteredProviders() returned %d providers, want 3", len(providers))
	}

	// Verify sorting: priority 1, then 5, then 10
	expectedOrder := []struct {
		id       string
		priority int
	}{
		{"provider-a", 1},
		{"provider-b", 5},
		{"provider-c", 10},
	}

	for i, expected := range expectedOrder {
		if providers[i].ID != expected.id {
			t.Errorf("providers[%d].ID = %q, want %q", i, providers[i].ID, expected.id)
		}
		if providers[i].Priority != expected.priority {
			t.Errorf(
				"providers[%d].Priority = %d, want %d",
				i,
				providers[i].Priority,
				expected.priority,
			)
		}
	}
}

// TestGet_Existing tests that Get returns the correct registration for existing providers.
func TestGet_Existing(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RegisterProvider(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 5,
		Provider: &testProvider{},
	})
	if err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}

	reg, found := Get("test-provider")

	if !found {
		t.Error("Get() returned found=false for existing provider")
	}
	if reg.ID != "test-provider" {
		t.Errorf("Get().ID = %q, want %q", reg.ID, "test-provider")
	}
	if reg.Name != "Test Provider" {
		t.Errorf("Get().Name = %q, want %q", reg.Name, "Test Provider")
	}
	if reg.Priority != 5 {
		t.Errorf("Get().Priority = %d, want %d", reg.Priority, 5)
	}
	if reg.Provider == nil {
		t.Error("Get().Provider is nil")
	}
}

// TestGet_NotFound tests that Get returns found=false for non-existent providers.
func TestGet_NotFound(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	reg, found := Get("nonexistent")

	if found {
		t.Error("Get() returned found=true for non-existent provider")
	}
	if reg.ID != "" {
		t.Errorf("Get() for non-existent should return zero Registration, got ID=%q", reg.ID)
	}
}

// TestCount tests that Count returns the correct number of registered providers.
func TestCount(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	if Count() != 0 {
		t.Errorf("Count() on empty registry = %d, want 0", Count())
	}

	for i := 1; i <= 3; i++ {
		err := RegisterProvider(Registration{
			ID:       "provider-" + string(rune('a'+i-1)),
			Name:     "Provider",
			Priority: i,
			Provider: &testProvider{},
		})
		if err != nil {
			t.Fatalf("RegisterProvider() failed: %v", err)
		}
		if Count() != i {
			t.Errorf("Count() after %d registrations = %d, want %d", i, Count(), i)
		}
	}
}

// TestResetRegistry tests that ResetRegistry clears all providers.
func TestResetRegistry(t *testing.T) {
	ResetRegistry()

	err := RegisterProvider(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &testProvider{},
	})
	if err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}

	if Count() != 1 {
		t.Errorf("Count() before reset = %d, want 1", Count())
	}

	ResetRegistry()

	if Count() != 0 {
		t.Errorf("Count() after reset = %d, want 0", Count())
	}
}

// TestRegisterAllProviders tests that all 15 providers are registered correctly.
func TestRegisterAllProviders(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() returned error: %v", err)
	}

	// Verify count
	if Count() != 15 {
		t.Errorf("Count() = %d, want 15", Count())
	}
}

// TestRegisterAllProviders_PrioritiesSequential tests that priorities are 1-15.
func TestRegisterAllProviders_PrioritiesSequential(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providers := RegisteredProviders()

	// Verify priorities are sequential 1-15
	for i, p := range providers {
		expectedPriority := i + 1
		if p.Priority != expectedPriority {
			t.Errorf(
				"providers[%d].Priority = %d, want %d (ID: %s)",
				i,
				p.Priority,
				expectedPriority,
				p.ID,
			)
		}
	}
}

// TestRegisterAllProviders_IDs tests that all expected provider IDs are registered.
func TestRegisterAllProviders_IDs(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	expectedIDs := []string{
		"claude-code",
		"gemini",
		"costrict",
		"qoder",
		"qwen",
		"antigravity",
		"cline",
		"cursor",
		"codex",
		"aider",
		"windsurf",
		"kilocode",
		"continue",
		"crush",
		"opencode",
	}

	for _, id := range expectedIDs {
		reg, found := Get(id)
		if !found {
			t.Errorf("Provider %q not found after RegisterAllProviders()", id)

			continue
		}
		if reg.Provider == nil {
			t.Errorf("Provider %q has nil Provider", id)
		}
	}
}

// TestRegisterAllProviders_Names tests that all providers have correct names.
func TestRegisterAllProviders_Names(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	expectedNames := map[string]string{
		"claude-code": "Claude Code",
		"gemini":      "Gemini CLI",
		"costrict":    "CoStrict",
		"qoder":       "Qoder",
		"qwen":        "Qwen Code",
		"antigravity": "Antigravity",
		"cline":       "Cline",
		"cursor":      "Cursor",
		"codex":       "Codex CLI",
		"aider":       "Aider",
		"windsurf":    "Windsurf",
		"kilocode":    "Kilocode",
		"continue":    "Continue",
		"crush":       "Crush",
		"opencode":    "OpenCode",
	}

	for id, expectedName := range expectedNames {
		reg, found := Get(id)
		if !found {
			t.Errorf("Provider %q not found", id)

			continue
		}
		if reg.Name != expectedName {
			t.Errorf("Provider %q has Name=%q, want %q", id, reg.Name, expectedName)
		}
	}
}

// TestRegisterAllProviders_DoesNotError tests that RegisterAllProviders doesn't error on first call.
func TestRegisterAllProviders_DoesNotError(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	err := RegisterAllProviders()
	if err != nil {
		t.Errorf("RegisterAllProviders() on fresh registry returned error: %v", err)
	}
}

// TestRegisterAllProviders_DoubleCall tests that calling RegisterAllProviders twice fails.
func TestRegisterAllProviders_DoubleCall(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	// First call should succeed
	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("First RegisterAllProviders() failed: %v", err)
	}

	// Second call should fail due to duplicate registration
	err := RegisterAllProviders()
	if err == nil {
		t.Error("Second RegisterAllProviders() should return error for duplicate registration")
	}
}

// TestProviderIDsAreKebabCase tests that all provider IDs are in kebab-case.
func TestProviderIDsAreKebabCase(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providers := RegisteredProviders()

	for _, p := range providers {
		id := p.ID
		for _, char := range id {
			if (char < 'a' || char > 'z') &&
				(char < '0' || char > '9') &&
				char != '-' {
				t.Errorf(
					"Provider ID %s is not in kebab-case (contains invalid character: %c)",
					id,
					char,
				)
			}
		}
	}
}

// TestRegisteredProviders_SortOrder verifies exact sort order after RegisterAllProviders.
func TestRegisteredProviders_SortOrder(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	if err := RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	providers := RegisteredProviders()

	expectedOrder := []string{
		"claude-code", // Priority 1
		"gemini",      // Priority 2
		"costrict",    // Priority 3
		"qoder",       // Priority 4
		"qwen",        // Priority 5
		"antigravity", // Priority 6
		"cline",       // Priority 7
		"cursor",      // Priority 8
		"codex",       // Priority 9
		"aider",       // Priority 10
		"windsurf",    // Priority 11
		"kilocode",    // Priority 12
		"continue",    // Priority 13
		"crush",       // Priority 14
		"opencode",    // Priority 15
	}

	if len(providers) != len(expectedOrder) {
		t.Fatalf("Got %d providers, want %d", len(providers), len(expectedOrder))
	}

	for i, expectedID := range expectedOrder {
		if providers[i].ID != expectedID {
			t.Errorf("providers[%d].ID = %q, want %q", i, providers[i].ID, expectedID)
		}
	}
}
