package providers

import (
	"testing"
)

// resetForTest clears the registry before each test
func resetForTest() {
	ResetRegistry()
}

func TestRegisterProvider_Valid(t *testing.T) {
	resetForTest()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &ClaudeProvider{},
	}

	err := RegisterProvider(reg)
	if err != nil {
		t.Fatalf("RegisterProvider() error = %v, want nil", err)
	}

	if Count() != 1 {
		t.Errorf("Count() = %d, want 1", Count())
	}
}

func TestRegisterProvider_EmptyID(t *testing.T) {
	resetForTest()

	reg := Registration{
		ID:       "",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &ClaudeProvider{},
	}

	err := RegisterProvider(reg)
	if err == nil {
		t.Fatal("RegisterProvider() error = nil, want error for empty ID")
	}

	expectedMsg := "provider ID is required"
	if err.Error() != expectedMsg {
		t.Errorf("RegisterProvider() error = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestRegisterProvider_NilProvider(t *testing.T) {
	resetForTest()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: nil,
	}

	err := RegisterProvider(reg)
	if err == nil {
		t.Fatal("RegisterProvider() error = nil, want error for nil Provider")
	}

	expectedMsg := "provider implementation is required"
	if err.Error() != expectedMsg {
		t.Errorf("RegisterProvider() error = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestRegisterProvider_DuplicateID(t *testing.T) {
	resetForTest()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &ClaudeProvider{},
	}

	// First registration should succeed
	err := RegisterProvider(reg)
	if err != nil {
		t.Fatalf("First RegisterProvider() error = %v, want nil", err)
	}

	// Second registration with same ID should fail
	err = RegisterProvider(reg)
	if err == nil {
		t.Fatal("RegisterProvider() error = nil, want error for duplicate ID")
	}

	expectedMsg := `provider "test-provider" already registered`
	if err.Error() != expectedMsg {
		t.Errorf("RegisterProvider() error = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestRegisteredProviders_Sorted(t *testing.T) {
	resetForTest()

	// Register providers in non-priority order
	providers := []Registration{
		{ID: "provider3", Name: "Provider 3", Priority: 3, Provider: &ClaudeProvider{}},
		{ID: "provider1", Name: "Provider 1", Priority: 1, Provider: &GeminiProvider{}},
		{ID: "provider2", Name: "Provider 2", Priority: 2, Provider: &ClineProvider{}},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%q) error = %v, want nil", reg.ID, err)
		}
	}

	result := RegisteredProviders()

	// Check count
	if len(result) != 3 {
		t.Fatalf("RegisteredProviders() returned %d providers, want 3", len(result))
	}

	// Check sorted by priority (lower first)
	if result[0].Priority != 1 {
		t.Errorf("result[0].Priority = %d, want 1", result[0].Priority)
	}
	if result[1].Priority != 2 {
		t.Errorf("result[1].Priority = %d, want 2", result[1].Priority)
	}
	if result[2].Priority != 3 {
		t.Errorf("result[2].Priority = %d, want 3", result[2].Priority)
	}

	// Check IDs match expected order
	if result[0].ID != "provider1" {
		t.Errorf("result[0].ID = %q, want %q", result[0].ID, "provider1")
	}
	if result[1].ID != "provider2" {
		t.Errorf("result[1].ID = %q, want %q", result[1].ID, "provider2")
	}
	if result[2].ID != "provider3" {
		t.Errorf("result[2].ID = %q, want %q", result[2].ID, "provider3")
	}
}

func TestGet_Found(t *testing.T) {
	resetForTest()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &ClaudeProvider{},
	}

	if err := RegisterProvider(reg); err != nil {
		t.Fatalf("RegisterProvider() error = %v, want nil", err)
	}

	result, found := Get("test-provider")
	if !found {
		t.Fatal("Get(test-provider) found = false, want true")
	}

	if result.ID != "test-provider" {
		t.Errorf("Get(test-provider).ID = %q, want %q", result.ID, "test-provider")
	}
	if result.Name != "Test Provider" {
		t.Errorf("Get(test-provider).Name = %q, want %q", result.Name, "Test Provider")
	}
	if result.Priority != 1 {
		t.Errorf("Get(test-provider).Priority = %d, want 1", result.Priority)
	}
}

func TestGet_NotFound(t *testing.T) {
	resetForTest()

	result, found := Get("nonexistent")
	if found {
		t.Fatal("Get(nonexistent) found = true, want false")
	}

	// Empty registration should be returned
	if result.ID != "" {
		t.Errorf("Get(nonexistent).ID = %q, want empty string", result.ID)
	}
}

func TestCount(t *testing.T) {
	resetForTest()

	if Count() != 0 {
		t.Errorf("Count() = %d, want 0 for empty registry", Count())
	}

	// Register 3 providers
	providers := []Registration{
		{ID: "provider1", Name: "Provider 1", Priority: 1, Provider: &ClaudeProvider{}},
		{ID: "provider2", Name: "Provider 2", Priority: 2, Provider: &GeminiProvider{}},
		{ID: "provider3", Name: "Provider 3", Priority: 3, Provider: &ClineProvider{}},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%q) error = %v, want nil", reg.ID, err)
		}
	}

	if Count() != 3 {
		t.Errorf("Count() = %d, want 3", Count())
	}
}

func TestRegisterAllProviders_Success(t *testing.T) {
	resetForTest()

	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() error = %v, want nil", err)
	}

	// Check all 15 providers are registered
	if Count() != 15 {
		t.Errorf("Count() = %d, want 15", Count())
	}

	// Check specific provider IDs and priorities
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			reg, found := Get(tt.id)
			if !found {
				t.Fatalf("Get(%q) found = false, want true", tt.id)
			}

			if reg.ID != tt.id {
				t.Errorf("Get(%q).ID = %q, want %q", tt.id, reg.ID, tt.id)
			}

			if reg.Name != tt.name {
				t.Errorf("Get(%q).Name = %q, want %q", tt.id, reg.Name, tt.name)
			}

			if reg.Priority != tt.priority {
				t.Errorf("Get(%q).Priority = %d, want %d", tt.id, reg.Priority, tt.priority)
			}

			if reg.Provider == nil {
				t.Errorf("Get(%q).Provider = nil, want non-nil", tt.id)
			}
		})
	}
}

func TestRegisterAllProviders_PrioritiesSequential(t *testing.T) {
	resetForTest()

	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() error = %v, want nil", err)
	}

	providers := RegisteredProviders()

	// Check priorities are sequential 1-15
	for i, reg := range providers {
		expectedPriority := i + 1
		if reg.Priority != expectedPriority {
			t.Errorf("providers[%d].Priority = %d, want %d", i, reg.Priority, expectedPriority)
		}
	}
}

func TestRegisterAllProviders_ProviderInitializers(t *testing.T) {
	resetForTest()

	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf("RegisterAllProviders() error = %v, want nil", err)
	}

	// Verify each provider is registered correctly
	// (actual initializer testing is done in provider_test.go with real TemplateManager)
	providers := RegisteredProviders()

	for _, reg := range providers {
		t.Run(reg.ID, func(t *testing.T) {
			// Provider should not be nil
			if reg.Provider == nil {
				t.Errorf("%s provider is nil", reg.ID)
			}
		})
	}
}
