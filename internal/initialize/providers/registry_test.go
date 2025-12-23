package providers

import (
	"context"
	"sync"
	"testing"
)

// mockProvider is a simple test implementation of Provider
type mockProvider struct {
	id string
}

func (m *mockProvider) Initializers(_ context.Context) []Initializer {
	_ = m

	return nil
}

func TestRegistry_Register(t *testing.T) {
	tests := []struct {
		name        string
		reg         Registration
		wantErr     bool
		errContains string
	}{
		{
			name: "successful registration",
			reg: Registration{
				ID:       "test-provider",
				Name:     "Test Provider",
				Priority: 1,
				Provider: &mockProvider{id: "test-provider"},
			},
			wantErr: false,
		},
		{
			name: "registration with priority 10",
			reg: Registration{
				ID:       "another-provider",
				Name:     "Another Provider",
				Priority: 10,
				Provider: &mockProvider{id: "another-provider"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			err := r.Register(tt.reg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegistry_DuplicateIDRejection(t *testing.T) {
	r := NewRegistry()

	reg1 := Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: 1,
		Provider: &mockProvider{id: "claude-code"},
	}

	// First registration should succeed
	err := r.Register(reg1)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Second registration with same ID should fail
	reg2 := Registration{
		ID:       "claude-code",
		Name:     "Claude Code 2",
		Priority: 2,
		Provider: &mockProvider{id: "claude-code"},
	}

	err = r.Register(reg2)
	if err == nil {
		t.Fatal("Expected error for duplicate ID, got nil")
	}

	expectedErrMsg := `provider with ID "claude-code" already registered`
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 5,
		Provider: &mockProvider{id: "test-provider"},
	}

	err := r.Register(reg)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	tests := []struct {
		name     string
		id       string
		wantNil  bool
		wantName string
	}{
		{
			name:     "get existing provider",
			id:       "test-provider",
			wantNil:  false,
			wantName: "Test Provider",
		},
		{
			name:    "get non-existent provider",
			id:      "does-not-exist",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Get(tt.id)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Get() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Fatal("Get() returned nil, want Registration")
				}
				if got.Name != tt.wantName {
					t.Errorf("Get().Name = %q, want %q", got.Name, tt.wantName)
				}
			}
		})
	}
}

func TestRegistry_All_PrioritySorting(t *testing.T) {
	r := NewRegistry()

	// Register providers in non-priority order
	providers := []Registration{
		{
			ID:       "low-priority",
			Name:     "Low Priority",
			Priority: 10,
			Provider: &mockProvider{id: "low-priority"},
		},
		{
			ID:       "high-priority",
			Name:     "High Priority",
			Priority: 1,
			Provider: &mockProvider{id: "high-priority"},
		},
		{
			ID:       "medium-priority",
			Name:     "Medium Priority",
			Priority: 5,
			Provider: &mockProvider{id: "medium-priority"},
		},
	}

	for _, p := range providers {
		if err := r.Register(p); err != nil {
			t.Fatalf("Failed to register %q: %v", p.ID, err)
		}
	}

	// Get all registrations
	all := r.All()

	// Should be sorted by priority (lower first)
	expectedOrder := []string{"high-priority", "medium-priority", "low-priority"}
	if len(all) != len(expectedOrder) {
		t.Fatalf("All() returned %d registrations, want %d", len(all), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if all[i].ID != expected {
			t.Errorf("All()[%d].ID = %q, want %q", i, all[i].ID, expected)
		}
	}
}

func TestRegistry_IDs_PrioritySorting(t *testing.T) {
	r := NewRegistry()

	// Register providers with various priorities
	providers := []Registration{
		{ID: "provider-a", Name: "Provider A", Priority: 100, Provider: &mockProvider{}},
		{ID: "provider-b", Name: "Provider B", Priority: 1, Provider: &mockProvider{}},
		{ID: "provider-c", Name: "Provider C", Priority: 50, Provider: &mockProvider{}},
		{ID: "provider-d", Name: "Provider D", Priority: 25, Provider: &mockProvider{}},
	}

	for _, p := range providers {
		if err := r.Register(p); err != nil {
			t.Fatalf("Failed to register %q: %v", p.ID, err)
		}
	}

	ids := r.IDs()
	expectedOrder := []string{"provider-b", "provider-d", "provider-c", "provider-a"}

	if len(ids) != len(expectedOrder) {
		t.Fatalf("IDs() returned %d IDs, want %d", len(ids), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if ids[i] != expected {
			t.Errorf("IDs()[%d] = %q, want %q", i, ids[i], expected)
		}
	}
}

func TestRegistry_Count(t *testing.T) {
	r := NewRegistry()

	if r.Count() != 0 {
		t.Errorf("Count() = %d, want 0 for empty registry", r.Count())
	}

	// Add providers one by one
	for i := 1; i <= 5; i++ {
		reg := Registration{
			ID:       string(rune('a' + i - 1)),
			Name:     "Provider",
			Priority: i,
			Provider: &mockProvider{},
		}
		if err := r.Register(reg); err != nil {
			t.Fatalf("Failed to register provider: %v", err)
		}

		if r.Count() != i {
			t.Errorf("Count() = %d, want %d", r.Count(), i)
		}
	}
}

func TestRegistry_GlobalFunctions(t *testing.T) {
	// Reset global registry before test
	Reset()

	reg1 := Registration{
		ID:       "global-test-1",
		Name:     "Global Test 1",
		Priority: 1,
		Provider: &mockProvider{},
	}
	reg2 := Registration{
		ID:       "global-test-2",
		Name:     "Global Test 2",
		Priority: 2,
		Provider: &mockProvider{},
	}

	// Test global Register
	if err := Register(reg1); err != nil {
		t.Fatalf("Global Register() failed: %v", err)
	}
	if err := Register(reg2); err != nil {
		t.Fatalf("Global Register() failed: %v", err)
	}

	// Test global Get
	got := Get("global-test-1")
	if got == nil {
		t.Fatal("Global Get() returned nil")
	}
	if got.Name != "Global Test 1" {
		t.Errorf("Global Get().Name = %q, want %q", got.Name, "Global Test 1")
	}

	// Test global Count
	if Count() != 2 {
		t.Errorf("Global Count() = %d, want 2", Count())
	}

	// Test global All
	all := All()
	if len(all) != 2 {
		t.Errorf("Global All() returned %d registrations, want 2", len(all))
	}

	// Test global IDs
	ids := IDs()
	expectedIDs := []string{"global-test-1", "global-test-2"}
	if len(ids) != len(expectedIDs) {
		t.Fatalf("Global IDs() returned %d IDs, want %d", len(ids), len(expectedIDs))
	}
	for i, expected := range expectedIDs {
		if ids[i] != expected {
			t.Errorf("Global IDs()[%d] = %q, want %q", i, ids[i], expected)
		}
	}

	// Clean up
	Reset()
}

func TestRegistry_ThreadSafety(t *testing.T) {
	r := NewRegistry()
	const numGoroutines = 10
	const numOpsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent registrations, gets, and counts
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()

			// Each goroutine registers its own provider
			reg := Registration{
				ID:       string(rune('a' + id)),
				Name:     "Provider",
				Priority: id,
				Provider: &mockProvider{},
			}
			if err := r.Register(reg); err != nil {
				// Ignore duplicate errors in concurrent test
				return
			}

			// Perform multiple reads
			for range numOpsPerGoroutine {
				_ = r.Get(string(rune('a' + id)))
				_ = r.Count()
				_ = r.All()
				_ = r.IDs()
			}
		}(i)
	}

	wg.Wait()

	// Verify final count
	count := r.Count()
	if count != numGoroutines {
		t.Errorf("Count() = %d, want %d after concurrent operations", count, numGoroutines)
	}
}

func TestRegistry_EmptyRegistry(t *testing.T) {
	r := NewRegistry()

	if r.Count() != 0 {
		t.Errorf("Count() = %d, want 0 for new registry", r.Count())
	}

	if r.Get("does-not-exist") != nil {
		t.Error("Get() should return nil for empty registry")
	}

	all := r.All()
	if len(all) != 0 {
		t.Errorf("All() returned %d registrations, want 0 for empty registry", len(all))
	}

	ids := r.IDs()
	if len(ids) != 0 {
		t.Errorf("IDs() returned %d IDs, want 0 for empty registry", len(ids))
	}
}
