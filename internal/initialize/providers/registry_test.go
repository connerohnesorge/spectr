package providers

import (
	"context"
	"testing"
)

// mockProvider is a simple test provider implementation that implements the Provider interface
type mockProvider struct {
	id   string
	name string
}

func (*mockProvider) Initializers(_ context.Context) []Initializer {
	return nil
}

func newMockProvider(id, name string) Provider {
	return &mockProvider{id: id, name: name}
}

func TestRegisterProvider(t *testing.T) {
	// Reset global registry before test
	Reset()

	tests := []struct {
		name        string
		reg         Registration
		wantErr     bool
		errContains string
	}{
		{
			name: "valid registration",
			reg: Registration{
				ID:       "test-provider",
				Name:     "Test Provider",
				Priority: 1,
				Provider: newMockProvider("test-provider", "Test Provider"),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			reg: Registration{
				ID:       "",
				Name:     "Test Provider",
				Priority: 1,
				Provider: newMockProvider("test", "Test"),
			},
			wantErr:     true,
			errContains: "provider ID is required",
		},
		{
			name: "missing Provider",
			reg: Registration{
				ID:       "test-provider",
				Name:     "Test Provider",
				Priority: 1,
				Provider: nil,
			},
			wantErr:     true,
			errContains: "provider implementation is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset for each test case
			Reset()

			err := RegisterProvider(tt.reg)
			if tt.wantErr {
				if err == nil {
					t.Errorf(
						"RegisterProvider() error = nil, want error containing %q",
						tt.errContains,
					)

					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf(
						"RegisterProvider() error = %q, want error containing %q",
						err.Error(),
						tt.errContains,
					)
				}
			} else if err != nil {
				t.Errorf("RegisterProvider() unexpected error: %v", err)
			}
		})
	}
}

func TestRegisterProvider_Duplicate(t *testing.T) {
	Reset()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider"),
	}

	// First registration should succeed
	err := RegisterProvider(reg)
	if err != nil {
		t.Fatalf("First RegisterProvider() failed: %v", err)
	}

	// Second registration with same ID should fail
	err = RegisterProvider(reg)
	if err == nil {
		t.Fatal("Second RegisterProvider() with duplicate ID should fail")
	}
	if !containsString(err.Error(), "already registered") {
		t.Errorf("Expected error about duplicate registration, got: %v", err)
	}
}

func TestGetRegistration(t *testing.T) {
	Reset()

	// Register test providers
	providers := []Registration{
		{
			ID:       "provider1",
			Name:     "Provider One",
			Priority: 1,
			Provider: newMockProvider("provider1", "Provider One"),
		},
		{
			ID:       "provider2",
			Name:     "Provider Two",
			Priority: 2,
			Provider: newMockProvider("provider2", "Provider Two"),
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%s) failed: %v", reg.ID, err)
		}
	}

	tests := []struct {
		name     string
		id       string
		wantNil  bool
		wantName string
	}{
		{
			name:     "Get existing provider1",
			id:       "provider1",
			wantNil:  false,
			wantName: "Provider One",
		},
		{
			name:     "Get existing provider2",
			id:       "provider2",
			wantNil:  false,
			wantName: "Provider Two",
		},
		{
			name:    "Get non-existent provider",
			id:      "nonexistent",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := GetRegistration(tt.id)
			if tt.wantNil {
				if reg != nil {
					t.Errorf("GetRegistration(%s) = %v, want nil", tt.id, reg)
				}
			} else {
				if reg == nil {
					t.Errorf("GetRegistration(%s) = nil, want non-nil", tt.id)

					return
				}
				if reg.Name != tt.wantName {
					t.Errorf("GetRegistration(%s).Name = %s, want %s", tt.id, reg.Name, tt.wantName)
				}
				if reg.ID != tt.id {
					t.Errorf("GetRegistration(%s).ID = %s, want %s", tt.id, reg.ID, tt.id)
				}
			}
		})
	}
}

func TestAllRegistrations(t *testing.T) {
	Reset()

	// Register providers with different priorities
	providers := []Registration{
		{
			ID:       "high-priority",
			Name:     "High Priority",
			Priority: 1,
			Provider: newMockProvider("high-priority", "High Priority"),
		},
		{
			ID:       "medium-priority",
			Name:     "Medium Priority",
			Priority: 5,
			Provider: newMockProvider("medium-priority", "Medium Priority"),
		},
		{
			ID:       "low-priority",
			Name:     "Low Priority",
			Priority: 10,
			Provider: newMockProvider("low-priority", "Low Priority"),
		},
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%s) failed: %v", reg.ID, err)
		}
	}

	all := AllRegistrations()

	// Check count
	if len(all) != len(providers) {
		t.Errorf("AllRegistrations() returned %d providers, want %d", len(all), len(providers))
	}

	// Verify sorted by priority (ascending)
	for i := 1; i < len(all); i++ {
		if all[i-1].Priority > all[i].Priority {
			t.Errorf(
				"Providers not sorted by priority: %s (priority %d) comes before %s (priority %d)",
				all[i-1].ID,
				all[i-1].Priority,
				all[i].ID,
				all[i].Priority,
			)
		}
	}

	// Verify expected order
	expectedOrder := []string{"high-priority", "medium-priority", "low-priority"}
	for i, expected := range expectedOrder {
		if all[i].ID != expected {
			t.Errorf("Position %d: got ID %s, want %s", i, all[i].ID, expected)
		}
	}
}

func TestIDs(t *testing.T) {
	Reset()

	// Register providers with different priorities
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
	}

	for _, reg := range providers {
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider(%s) failed: %v", reg.ID, err)
		}
	}

	ids := IDs()

	// Check count
	if len(ids) != len(providers) {
		t.Errorf("IDs() returned %d IDs, want %d", len(ids), len(providers))
	}

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

	// Verify sorted by priority
	expectedOrder := []string{"provider-a", "provider-b", "provider-c"}
	for i, expected := range expectedOrder {
		if ids[i] != expected {
			t.Errorf("Position %d: got ID %s, want %s", i, ids[i], expected)
		}
	}
}

func TestCount(t *testing.T) {
	Reset()

	// Initially empty
	if Count() != 0 {
		t.Errorf("Count() = %d, want 0 for empty registry", Count())
	}

	// Add providers one by one
	for i := 1; i <= 3; i++ {
		reg := Registration{
			ID:       "provider" + string(rune('0'+i)),
			Name:     "Provider",
			Priority: i,
			Provider: newMockProvider("provider", "Provider"),
		}
		if err := RegisterProvider(reg); err != nil {
			t.Fatalf("RegisterProvider failed: %v", err)
		}

		if Count() != i {
			t.Errorf("After registering %d providers, Count() = %d, want %d", i, Count(), i)
		}
	}
}

func TestRegistryInstance_RegisterProvider(t *testing.T) {
	r := NewRegistry()

	// Initially empty
	if r.Count() != 0 {
		t.Errorf("New registry should be empty, got %d providers", r.Count())
	}

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider"),
	}

	// Register should succeed
	err := r.registerProvider(reg)
	if err != nil {
		t.Errorf("registerProvider() failed: %v", err)
	}

	if r.Count() != 1 {
		t.Errorf("Expected 1 provider after registration, got %d", r.Count())
	}

	// Duplicate registration should fail
	err = r.registerProvider(reg)
	if err == nil {
		t.Error("Duplicate registration should fail")
	}
}

func TestRegistryInstance_GetRegistration(t *testing.T) {
	r := NewRegistry()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider"),
	}

	// Register provider
	if err := r.registerProvider(reg); err != nil {
		t.Fatalf("registerProvider() failed: %v", err)
	}

	// Get existing provider
	retrieved := r.GetRegistration("test-provider")
	if retrieved == nil {
		t.Fatal("GetRegistration() returned nil for registered provider")
	}
	if retrieved.ID != "test-provider" {
		t.Errorf(
			"GetRegistration() returned wrong provider: got ID %s, want %s",
			retrieved.ID,
			"test-provider",
		)
	}

	// Get non-existent provider
	notFound := r.GetRegistration("nonexistent")
	if notFound != nil {
		t.Errorf("GetRegistration() for non-existent provider should return nil, got %v", notFound)
	}
}

func TestRegistryInstance_AllRegistrations(t *testing.T) {
	r := NewRegistry()

	// Register providers with different priorities
	providers := []Registration{
		{
			ID:       "provider-3",
			Name:     "Provider 3",
			Priority: 3,
			Provider: newMockProvider("provider-3", "Provider 3"),
		},
		{
			ID:       "provider-1",
			Name:     "Provider 1",
			Priority: 1,
			Provider: newMockProvider("provider-1", "Provider 1"),
		},
		{
			ID:       "provider-2",
			Name:     "Provider 2",
			Priority: 2,
			Provider: newMockProvider("provider-2", "Provider 2"),
		},
	}

	for _, reg := range providers {
		if err := r.registerProvider(reg); err != nil {
			t.Fatalf("registerProvider() failed: %v", err)
		}
	}

	all := r.AllRegistrations()

	// Check count
	if len(all) != len(providers) {
		t.Errorf("AllRegistrations() returned %d providers, want %d", len(all), len(providers))
	}

	// Verify sorted by priority
	expectedOrder := []string{"provider-1", "provider-2", "provider-3"}
	for i, expected := range expectedOrder {
		if all[i].ID != expected {
			t.Errorf("Position %d: got ID %s, want %s", i, all[i].ID, expected)
		}
	}
}

func TestRegistryInstance_IDs(t *testing.T) {
	r := NewRegistry()

	// Register providers
	providers := []Registration{
		{
			ID:       "provider-b",
			Name:     "Provider B",
			Priority: 2,
			Provider: newMockProvider("provider-b", "Provider B"),
		},
		{
			ID:       "provider-a",
			Name:     "Provider A",
			Priority: 1,
			Provider: newMockProvider("provider-a", "Provider A"),
		},
	}

	for _, reg := range providers {
		if err := r.registerProvider(reg); err != nil {
			t.Fatalf("registerProvider() failed: %v", err)
		}
	}

	ids := r.IDs()

	// Check count
	if len(ids) != len(providers) {
		t.Errorf("IDs() returned %d IDs, want %d", len(ids), len(providers))
	}

	// Verify sorted by priority
	expectedOrder := []string{"provider-a", "provider-b"}
	for i, expected := range expectedOrder {
		if ids[i] != expected {
			t.Errorf("Position %d: got ID %s, want %s", i, ids[i], expected)
		}
	}
}

func TestRegistryInstance_Count(t *testing.T) {
	r := NewRegistry()

	// Initially empty
	if r.Count() != 0 {
		t.Errorf("Count() = %d, want 0 for new registry", r.Count())
	}

	// Add providers
	for i := 1; i <= 3; i++ {
		reg := Registration{
			ID:       "provider" + string(rune('0'+i)),
			Name:     "Provider",
			Priority: i,
			Provider: newMockProvider("provider", "Provider"),
		}
		if err := r.registerProvider(reg); err != nil {
			t.Fatalf("registerProvider() failed: %v", err)
		}

		if r.Count() != i {
			t.Errorf("After registering %d providers, Count() = %d, want %d", i, r.Count(), i)
		}
	}
}

func TestReset(t *testing.T) {
	Reset()

	// Register a provider
	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider("test-provider", "Test Provider"),
	}
	if err := RegisterProvider(reg); err != nil {
		t.Fatalf("RegisterProvider() failed: %v", err)
	}

	if Count() != 1 {
		t.Errorf("Count() = %d, want 1 after registration", Count())
	}

	// Reset
	Reset()

	if Count() != 0 {
		t.Errorf("Count() = %d, want 0 after Reset()", Count())
	}
}

func TestRegisterAllProviders(t *testing.T) {
	Reset()

	err := RegisterAllProviders()
	if err != nil {
		t.Errorf("RegisterAllProviders() failed: %v", err)
	}

	// Should register all providers
	count := Count()
	if count == 0 {
		t.Error("RegisterAllProviders() should register at least one provider")
	}

	// Verify all registered providers are valid
	for _, reg := range AllRegistrations() {
		if err := reg.Validate(); err != nil {
			t.Errorf("Invalid registration for %s: %v", reg.ID, err)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
