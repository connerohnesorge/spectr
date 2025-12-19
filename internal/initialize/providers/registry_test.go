package providers

import (
	"context"
	"errors"
	"testing"
)

// mockProvider is a simple implementation of Provider for testing.
type mockProvider struct {
	initializers []Initializer
}

func (m *mockProvider) Initializers(_ context.Context) []Initializer {
	return m.initializers
}

// newMockProvider creates a mock provider for testing.
func newMockProvider() *mockProvider {
	return &mockProvider{}
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	// Register a provider
	err := r.Register(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider(),
	})
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	if r.Count() != 1 {
		t.Errorf(
			"Expected 1 provider after registration, got %d",
			r.Count(),
		)
	}
}

func TestRegistry_RegisterRejectsDuplicate(t *testing.T) {
	r := NewRegistry()

	// Register first provider
	err := r.Register(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider(),
	})
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register duplicate ID
	err = r.Register(Registration{
		ID:       "test-provider",
		Name:     "Another Provider",
		Priority: 2,
		Provider: newMockProvider(),
	})

	if err == nil {
		t.Error("Duplicate registration should fail")
	}

	if !errors.Is(err, ErrDuplicateID) {
		t.Errorf(
			"Expected ErrDuplicateID, got: %v",
			err,
		)
	}

	// Count should still be 1
	if r.Count() != 1 {
		t.Errorf(
			"Expected 1 provider after duplicate rejection, got %d",
			r.Count(),
		)
	}
}

func TestRegistry_RegisterValidation(t *testing.T) {
	tests := []struct {
		name    string
		reg     Registration
		wantErr error
	}{
		{
			name: "empty ID",
			reg: Registration{
				ID:       "",
				Name:     "Test",
				Provider: newMockProvider(),
			},
			wantErr: ErrEmptyID,
		},
		{
			name: "empty Name",
			reg: Registration{
				ID:       "test",
				Name:     "",
				Provider: newMockProvider(),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			err := r.Register(tt.reg)

			if err == nil {
				t.Error("Expected error, got nil")
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"Expected %v, got: %v",
					tt.wantErr,
					err,
				)
			}
		})
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	// Register a provider
	_ = r.Register(Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: newMockProvider(),
	})

	// Get existing provider
	reg := r.Get("test-provider")
	if reg == nil {
		t.Error("Get returned nil for registered provider")
		return
	}

	if reg.ID != "test-provider" {
		t.Errorf("Expected ID 'test-provider', got %q", reg.ID)
	}

	if reg.Name != "Test Provider" {
		t.Errorf("Expected Name 'Test Provider', got %q", reg.Name)
	}

	// Get non-existent provider
	reg = r.Get("nonexistent")
	if reg != nil {
		t.Errorf("Get returned non-nil for nonexistent provider: %v", reg)
	}
}

func TestRegistry_All_PrioritySorted(t *testing.T) {
	r := NewRegistry()

	// Register providers in non-priority order
	_ = r.Register(Registration{
		ID:       "low-priority",
		Name:     "Low Priority",
		Priority: 100,
		Provider: newMockProvider(),
	})
	_ = r.Register(Registration{
		ID:       "high-priority",
		Name:     "High Priority",
		Priority: 1,
		Provider: newMockProvider(),
	})
	_ = r.Register(Registration{
		ID:       "medium-priority",
		Name:     "Medium Priority",
		Priority: 50,
		Provider: newMockProvider(),
	})

	all := r.All()

	if len(all) != 3 {
		t.Fatalf("Expected 3 providers, got %d", len(all))
	}

	// Verify sorted by priority (lower first)
	expectedOrder := []string{"high-priority", "medium-priority", "low-priority"}
	for i, expected := range expectedOrder {
		if all[i].ID != expected {
			t.Errorf(
				"Position %d: expected %q, got %q",
				i,
				expected,
				all[i].ID,
			)
		}
	}

	// Verify priorities are in ascending order
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
}

func TestRegistry_IDs_PrioritySorted(t *testing.T) {
	r := NewRegistry()

	// Register providers in non-priority order
	_ = r.Register(Registration{
		ID:       "c-provider",
		Name:     "C Provider",
		Priority: 30,
		Provider: newMockProvider(),
	})
	_ = r.Register(Registration{
		ID:       "a-provider",
		Name:     "A Provider",
		Priority: 10,
		Provider: newMockProvider(),
	})
	_ = r.Register(Registration{
		ID:       "b-provider",
		Name:     "B Provider",
		Priority: 20,
		Provider: newMockProvider(),
	})

	ids := r.IDs()

	if len(ids) != 3 {
		t.Fatalf("Expected 3 IDs, got %d", len(ids))
	}

	// Verify sorted by priority (not alphabetical)
	expectedOrder := []string{"a-provider", "b-provider", "c-provider"}
	for i, expected := range expectedOrder {
		if ids[i] != expected {
			t.Errorf(
				"Position %d: expected %q, got %q",
				i,
				expected,
				ids[i],
			)
		}
	}
}

func TestRegistry_Count(t *testing.T) {
	r := NewRegistry()

	// Initially empty
	if r.Count() != 0 {
		t.Errorf("New registry should be empty, got %d providers", r.Count())
	}

	// Add providers
	for i := 0; i < 5; i++ {
		_ = r.Register(Registration{
			ID:       "provider-" + string(rune('a'+i)),
			Name:     "Provider",
			Priority: i,
			Provider: newMockProvider(),
		})
	}

	if r.Count() != 5 {
		t.Errorf("Expected 5 providers, got %d", r.Count())
	}
}

// Tests for global registry functions

func TestGlobalRegistry_RegisterAndGet(t *testing.T) {
	// Reset global registry for clean test
	Reset()
	defer Reset()

	err := Register(Registration{
		ID:       "global-test-provider",
		Name:     "Global Test Provider",
		Priority: 1,
		Provider: newMockProvider(),
	})
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	reg := Get("global-test-provider")
	if reg == nil {
		t.Error("Get returned nil for registered provider")
		return
	}

	if reg.ID != "global-test-provider" {
		t.Errorf("Expected ID 'global-test-provider', got %q", reg.ID)
	}

	// Test Get for non-existent
	reg = Get("nonexistent")
	if reg != nil {
		t.Errorf("Get returned non-nil for nonexistent provider: %v", reg)
	}
}

func TestGlobalRegistry_RejectsDuplicate(t *testing.T) {
	Reset()
	defer Reset()

	_ = Register(Registration{
		ID:       "dup-test",
		Name:     "Dup Test",
		Priority: 1,
		Provider: newMockProvider(),
	})

	err := Register(Registration{
		ID:       "dup-test",
		Name:     "Another Dup Test",
		Priority: 2,
		Provider: newMockProvider(),
	})

	if err == nil {
		t.Error("Duplicate registration should fail")
	}

	if !errors.Is(err, ErrDuplicateID) {
		t.Errorf("Expected ErrDuplicateID, got: %v", err)
	}
}

func TestGlobalRegistry_AllSorted(t *testing.T) {
	Reset()
	defer Reset()

	_ = Register(Registration{
		ID:       "z-provider",
		Name:     "Z Provider",
		Priority: 99,
		Provider: newMockProvider(),
	})
	_ = Register(Registration{
		ID:       "a-provider",
		Name:     "A Provider",
		Priority: 1,
		Provider: newMockProvider(),
	})
	_ = Register(Registration{
		ID:       "m-provider",
		Name:     "M Provider",
		Priority: 50,
		Provider: newMockProvider(),
	})

	all := All()

	if len(all) != 3 {
		t.Fatalf("Expected 3 registrations, got %d", len(all))
	}

	// Should be sorted by priority, not alphabetically
	if all[0].ID != "a-provider" {
		t.Errorf("First should be 'a-provider' (priority 1), got %q", all[0].ID)
	}
	if all[1].ID != "m-provider" {
		t.Errorf("Second should be 'm-provider' (priority 50), got %q", all[1].ID)
	}
	if all[2].ID != "z-provider" {
		t.Errorf("Third should be 'z-provider' (priority 99), got %q", all[2].ID)
	}
}

func TestGlobalRegistry_IDs(t *testing.T) {
	Reset()
	defer Reset()

	_ = Register(Registration{
		ID:       "third",
		Name:     "Third",
		Priority: 30,
		Provider: newMockProvider(),
	})
	_ = Register(Registration{
		ID:       "first",
		Name:     "First",
		Priority: 10,
		Provider: newMockProvider(),
	})
	_ = Register(Registration{
		ID:       "second",
		Name:     "Second",
		Priority: 20,
		Provider: newMockProvider(),
	})

	ids := IDs()

	expected := []string{"first", "second", "third"}
	for i, exp := range expected {
		if ids[i] != exp {
			t.Errorf("Position %d: expected %q, got %q", i, exp, ids[i])
		}
	}
}

func TestGlobalRegistry_Count(t *testing.T) {
	Reset()
	defer Reset()

	if Count() != 0 {
		t.Errorf("Empty registry should have count 0, got %d", Count())
	}

	_ = Register(Registration{
		ID:       "test1",
		Name:     "Test 1",
		Priority: 1,
		Provider: newMockProvider(),
	})
	_ = Register(Registration{
		ID:       "test2",
		Name:     "Test 2",
		Priority: 2,
		Provider: newMockProvider(),
	})

	if Count() != 2 {
		t.Errorf("Expected count 2, got %d", Count())
	}
}

func TestGlobalRegistry_Reset(t *testing.T) {
	Reset()

	_ = Register(Registration{
		ID:       "to-reset",
		Name:     "To Reset",
		Priority: 1,
		Provider: newMockProvider(),
	})

	if Count() != 1 {
		t.Errorf("Expected 1 provider, got %d", Count())
	}

	Reset()

	if Count() != 0 {
		t.Errorf("After reset, expected 0 providers, got %d", Count())
	}

	if Get("to-reset") != nil {
		t.Error("After reset, Get should return nil")
	}
}

func TestNewRegistryFromGlobal(t *testing.T) {
	Reset()
	defer Reset()

	_ = Register(Registration{
		ID:       "global-provider",
		Name:     "Global Provider",
		Priority: 1,
		Provider: newMockProvider(),
	})

	r := NewRegistryFromGlobal()

	// Should have same count as global
	if r.Count() != Count() {
		t.Errorf(
			"NewRegistryFromGlobal() has %d providers, global has %d",
			r.Count(),
			Count(),
		)
	}

	// Should be able to get same providers
	reg := r.Get("global-provider")
	if reg == nil {
		t.Error("Get returned nil for global-provider")
		return
	}

	if reg.Name != "Global Provider" {
		t.Errorf("Expected name 'Global Provider', got %q", reg.Name)
	}
}

func TestRegistry_ZeroPriority(t *testing.T) {
	r := NewRegistry()

	// Priority 0 is valid and should sort before priority 1
	_ = r.Register(Registration{
		ID:       "priority-one",
		Name:     "Priority One",
		Priority: 1,
		Provider: newMockProvider(),
	})
	_ = r.Register(Registration{
		ID:       "priority-zero",
		Name:     "Priority Zero",
		Priority: 0,
		Provider: newMockProvider(),
	})

	all := r.All()

	if len(all) != 2 {
		t.Fatalf("Expected 2 providers, got %d", len(all))
	}

	if all[0].ID != "priority-zero" {
		t.Errorf("Priority 0 should come first, got %q", all[0].ID)
	}
}

func TestRegistry_NegativePriority(t *testing.T) {
	r := NewRegistry()

	// Negative priority is valid
	_ = r.Register(Registration{
		ID:       "priority-positive",
		Name:     "Priority Positive",
		Priority: 10,
		Provider: newMockProvider(),
	})
	_ = r.Register(Registration{
		ID:       "priority-negative",
		Name:     "Priority Negative",
		Priority: -5,
		Provider: newMockProvider(),
	})

	all := r.All()

	if len(all) != 2 {
		t.Fatalf("Expected 2 providers, got %d", len(all))
	}

	if all[0].ID != "priority-negative" {
		t.Errorf("Negative priority should come first, got %q", all[0].ID)
	}
}
