package providers

import (
	"context"
	"fmt"
	"testing"
)

// mockProvider is a simple test provider implementation
type mockProvider struct {
	id   string
	name string
}

func (*mockProvider) Initializers(_ context.Context, _ TemplateManager) []Initializer {
	return nil
}

func newMockProvider(id, name string) *mockProvider {
	return &mockProvider{id: id, name: name}
}

func TestRegisterProvider_Valid(t *testing.T) {
	Reset() // Start with clean registry

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider"),
	}

	err := RegisterProvider(reg)
	if err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}

	// Verify it was registered
	if Count() != 1 {
		t.Errorf("Count() = %d, want 1", Count())
	}

	retrieved, ok := Get("test-provider")
	if !ok {
		t.Fatal("Get() returned false, want true")
	}
	if retrieved.ID != "test-provider" {
		t.Errorf("Get().ID = %s, want test-provider", retrieved.ID)
	}
	if retrieved.Name != "Test Provider" {
		t.Errorf("Get().Name = %s, want Test Provider", retrieved.Name)
	}
	if retrieved.Priority != 1 {
		t.Errorf("Get().Priority = %d, want 1", retrieved.Priority)
	}
}

func TestRegisterProvider_EmptyID(t *testing.T) {
	Reset()

	reg := Registration{
		ID:       "", // Empty ID
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test", "Test"),
	}

	err := RegisterProvider(reg)
	if err == nil {
		t.Fatal("RegisterProvider() succeeded with empty ID, want error")
	}
	if err.Error() != "provider ID is required" {
		t.Errorf("RegisterProvider() error = %v, want 'provider ID is required'", err)
	}

	// Verify nothing was registered
	if Count() != 0 {
		t.Errorf("Count() = %d, want 0", Count())
	}
}

func TestRegisterProvider_NilProvider(t *testing.T) {
	Reset()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: nil, // Nil provider
	}

	err := RegisterProvider(reg)
	if err == nil {
		t.Fatal("RegisterProvider() succeeded with nil Provider, want error")
	}
	if err.Error() != "provider implementation is required" {
		t.Errorf("RegisterProvider() error = %v, want 'provider implementation is required'", err)
	}

	// Verify nothing was registered
	if Count() != 0 {
		t.Errorf("Count() = %d, want 0", Count())
	}
}

func TestRegisterProvider_DuplicateID(t *testing.T) {
	Reset()

	// Register first provider
	reg1 := Registration{
		ID:       "test-provider",
		Name:     "Test Provider 1",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider 1"),
	}
	err := RegisterProvider(reg1)
	if err != nil {
		t.Fatalf("First RegisterProvider() failed: %v", err)
	}

	// Try to register duplicate
	reg2 := Registration{
		ID:       "test-provider", // Same ID
		Name:     "Test Provider 2",
		Priority: 2,
		Provider: newMockProvider("test-provider", "Test Provider 2"),
	}
	err = RegisterProvider(reg2)
	if err == nil {
		t.Fatal("RegisterProvider() succeeded with duplicate ID, want error")
	}
	expectedErr := `provider "test-provider" already registered`
	if err.Error() != expectedErr {
		t.Errorf("RegisterProvider() error = %v, want %q", err, expectedErr)
	}

	// Verify only first provider is registered
	if Count() != 1 {
		t.Errorf("Count() = %d, want 1", Count())
	}

	retrieved, ok := Get("test-provider")
	if !ok {
		t.Fatal("Get() returned false, want true")
	}
	if retrieved.Name != "Test Provider 1" {
		t.Errorf(
			"Get().Name = %s, want Test Provider 1 (first registration should be kept)",
			retrieved.Name,
		)
	}
}

func TestRegisteredProviders_SortedByPriority(t *testing.T) {
	Reset()

	// Register providers in random priority order
	providers := []Registration{
		{
			ID:       "provider-c",
			Name:     "Provider C",
			Priority: 3,
			Provider: newMockProvider("provider-c", "Provider C"),
		},
		{
			ID:       "provider-a",
			Name:     "Provider A",
			Priority: 1,
			Provider: newMockProvider("provider-a", "Provider A"),
		},
		{
			ID:       "provider-b",
			Name:     "Provider B",
			Priority: 2,
			Provider: newMockProvider("provider-b", "Provider B"),
		},
		{
			ID:       "provider-e",
			Name:     "Provider E",
			Priority: 5,
			Provider: newMockProvider("provider-e", "Provider E"),
		},
		{
			ID:       "provider-d",
			Name:     "Provider D",
			Priority: 4,
			Provider: newMockProvider("provider-d", "Provider D"),
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%s) failed: %v", reg.ID, err)
		}
	}

	// Get all providers
	registered := RegisteredProviders()

	// Verify count
	if len(registered) != 5 {
		t.Fatalf("RegisteredProviders() returned %d providers, want 5", len(registered))
	}

	// Verify sorted by priority (lower priority number = higher priority)
	expectedOrder := []string{"provider-a", "provider-b", "provider-c", "provider-d", "provider-e"}
	for i, reg := range registered {
		if reg.ID != expectedOrder[i] {
			t.Errorf("RegisteredProviders()[%d].ID = %s, want %s", i, reg.ID, expectedOrder[i])
		}
		if reg.Priority != i+1 {
			t.Errorf("RegisteredProviders()[%d].Priority = %d, want %d", i, reg.Priority, i+1)
		}
	}

	// Verify priorities are in ascending order
	for i := 1; i < len(registered); i++ {
		if registered[i-1].Priority >= registered[i].Priority {
			t.Errorf("RegisteredProviders() not sorted: [%d].Priority=%d >= [%d].Priority=%d",
				i-1, registered[i-1].Priority, i, registered[i].Priority)
		}
	}
}

func TestGet_Found(t *testing.T) {
	Reset()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider"),
	}
	if err := RegisterProvider(reg); err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}

	retrieved, ok := Get("test-provider")
	if !ok {
		t.Fatal("Get() returned false, want true")
	}
	if retrieved.ID != "test-provider" {
		t.Errorf("Get().ID = %s, want test-provider", retrieved.ID)
	}
	if retrieved.Name != "Test Provider" {
		t.Errorf("Get().Name = %s, want Test Provider", retrieved.Name)
	}
	if retrieved.Priority != 1 {
		t.Errorf("Get().Priority = %d, want 1", retrieved.Priority)
	}
	if retrieved.Provider == nil {
		t.Error("Get().Provider is nil, want non-nil")
	}
}

func TestGet_NotFound(t *testing.T) {
	Reset()

	retrieved, ok := Get("nonexistent")
	if ok {
		t.Error("Get() returned true for nonexistent provider, want false")
	}
	if retrieved.ID != "" {
		t.Errorf("Get() returned non-empty Registration for nonexistent provider: %+v", retrieved)
	}
}

func TestCount(t *testing.T) {
	Reset()

	// Initially empty
	if Count() != 0 {
		t.Errorf("Count() = %d, want 0", Count())
	}

	// Register one provider
	reg1 := Registration{
		ID:       "provider-1",
		Name:     "Provider 1",
		Priority: 1,
		Provider: newMockProvider("provider-1", "Provider 1"),
	}
	if err := RegisterProvider(reg1); err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}
	if Count() != 1 {
		t.Errorf("Count() = %d, want 1", Count())
	}

	// Register another provider
	reg2 := Registration{
		ID:       "provider-2",
		Name:     "Provider 2",
		Priority: 2,
		Provider: newMockProvider("provider-2", "Provider 2"),
	}
	if err := RegisterProvider(reg2); err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}
	if Count() != 2 {
		t.Errorf("Count() = %d, want 2", Count())
	}

	// Reset and verify count is 0
	Reset()
	if Count() != 0 {
		t.Errorf("Count() after Reset() = %d, want 0", Count())
	}
}

func TestReset(t *testing.T) {
	Reset()

	// Register some providers
	for i := 1; i <= 3; i++ {
		reg := Registration{
			ID:       fmt.Sprintf("provider-%d", i),
			Name:     fmt.Sprintf("Provider %d", i),
			Priority: i,
			Provider: newMockProvider(fmt.Sprintf("provider-%d", i), fmt.Sprintf("Provider %d", i)),
		}
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider() failed: %v", err)
		}
	}

	if Count() != 3 {
		t.Fatalf("Count() = %d, want 3", Count())
	}

	// Reset
	Reset()

	// Verify empty
	if Count() != 0 {
		t.Errorf("Count() after Reset() = %d, want 0", Count())
	}

	_, ok := Get("provider-1")
	if ok {
		t.Error("Get() returned true after Reset(), want false")
	}

	registered := RegisteredProviders()
	if len(registered) != 0 {
		t.Errorf(
			"RegisteredProviders() after Reset() returned %d providers, want 0",
			len(registered),
		)
	}
}

// Note: TestRegisterAllProviders is skipped until Phase 5 when provider implementations exist
/*
func TestRegisterAllProviders(t *testing.T) {
	Reset()

	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	// Verify all 15 providers are registered
	if Count() != 15 {
		t.Errorf("Count() = %d, want 15", Count())
	}

	// Verify priorities are sequential 1-15
	registered := RegisteredProviders()
	expectedProviders := []struct {
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

	for i, expected := range expectedProviders {
		if i >= len(registered) {
			t.Fatalf("RegisteredProviders() has only %d providers, expected at least %d", len(registered), i+1)
		}
		reg := registered[i]
		if reg.ID != expected.id {
			t.Errorf("RegisteredProviders()[%d].ID = %s, want %s", i, reg.ID, expected.id)
		}
		if reg.Name != expected.name {
			t.Errorf("RegisteredProviders()[%d].Name = %s, want %s", i, reg.Name, expected.name)
		}
		if reg.Priority != expected.priority {
			t.Errorf("RegisteredProviders()[%d].Priority = %d, want %d", i, reg.Priority, expected.priority)
		}
		if reg.Provider == nil {
			t.Errorf("RegisteredProviders()[%d].Provider is nil", i)
		}
	}

	// Verify no errors when retrieving each provider
	for _, expected := range expectedProviders {
		reg, ok := Get(expected.id)
		if !ok {
			t.Errorf("Get(%s) returned false, want true", expected.id)
		}
		if reg.ID != expected.id {
			t.Errorf("Get(%s).ID = %s, want %s", expected.id, reg.ID, expected.id)
		}
	}
}

func TestRegisterAllProviders_NoDuplicates(t *testing.T) {
	Reset()

	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	// Try to register again - should fail with duplicate error
	err = RegisterAllProviders()
	if err == nil {
		t.Fatal("Second RegisterAllProviders() succeeded, want error")
	}

	// Verify error mentions "already registered"
	if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("RegisterAllProviders() error = %v, want error containing 'already registered'", err)
	}
}
*/
