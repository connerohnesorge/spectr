package providers

import (
	"context"
	"strings"
	"testing"
)

// mockProviderV2 is a minimal ProviderV2 implementation for testing.
type mockProviderV2 struct {
	initializers []Initializer
}

func (m *mockProviderV2) Initializers(_ context.Context) []Initializer {
	return m.initializers
}

// newMockProviderV2 creates a mock ProviderV2 for testing.
func newMockProviderV2() ProviderV2 {
	return &mockProviderV2{}
}

func TestRegistryV2_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		registry := NewRegistryV2()

		err := registry.Register(Registration{
			ID:       "test-provider",
			Name:     "Test Provider",
			Priority: 1,
			Provider: newMockProviderV2(),
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if registry.Count() != 1 {
			t.Errorf("expected count 1, got %d", registry.Count())
		}
	})

	t.Run("multiple registrations", func(t *testing.T) {
		registry := NewRegistryV2()

		providers := []Registration{
			{ID: "provider-a", Name: "Provider A", Priority: 10, Provider: newMockProviderV2()},
			{ID: "provider-b", Name: "Provider B", Priority: 5, Provider: newMockProviderV2()},
			{ID: "provider-c", Name: "Provider C", Priority: 20, Provider: newMockProviderV2()},
		}

		for _, reg := range providers {
			if err := registry.Register(reg); err != nil {
				t.Errorf("failed to register %s: %v", reg.ID, err)
			}
		}

		if registry.Count() != 3 {
			t.Errorf("expected count 3, got %d", registry.Count())
		}
	})
}

func TestRegistryV2_DuplicateRejection(t *testing.T) {
	t.Run("rejects duplicate ID", func(t *testing.T) {
		registry := NewRegistryV2()

		// First registration should succeed
		err := registry.Register(Registration{
			ID:       "duplicate-test",
			Name:     "First Provider",
			Priority: 1,
			Provider: newMockProviderV2(),
		})
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		// Second registration with same ID should fail
		err = registry.Register(Registration{
			ID:       "duplicate-test",
			Name:     "Second Provider",
			Priority: 2,
			Provider: newMockProviderV2(),
		})

		if err == nil {
			t.Error("expected error for duplicate ID, got nil")
		}

		if !strings.Contains(err.Error(), "already registered") {
			t.Errorf("error message should mention 'already registered', got: %v", err)
		}

		if !strings.Contains(err.Error(), "duplicate-test") {
			t.Errorf("error message should mention the duplicate ID, got: %v", err)
		}

		// Count should still be 1
		if registry.Count() != 1 {
			t.Errorf("expected count 1 after duplicate rejection, got %d", registry.Count())
		}
	})

	t.Run("error message includes existing provider name", func(t *testing.T) {
		registry := NewRegistryV2()

		_ = registry.Register(Registration{
			ID:       "my-provider",
			Name:     "Original Name",
			Priority: 1,
			Provider: newMockProviderV2(),
		})

		err := registry.Register(Registration{
			ID:       "my-provider",
			Name:     "New Name",
			Priority: 2,
			Provider: newMockProviderV2(),
		})

		if err == nil {
			t.Fatal("expected error for duplicate ID")
		}

		if !strings.Contains(err.Error(), "Original Name") {
			t.Errorf("error message should include existing provider name, got: %v", err)
		}
	})
}

func TestRegistryV2_ValidationRejection(t *testing.T) {
	tests := []struct {
		name        string
		reg         Registration
		errContains string
	}{
		{
			name: "empty ID",
			reg: Registration{
				ID:       "",
				Name:     "Test",
				Priority: 1,
				Provider: newMockProviderV2(),
			},
			errContains: "ID is required",
		},
		{
			name: "invalid ID format",
			reg: Registration{
				ID:       "InvalidCase",
				Name:     "Test",
				Priority: 1,
				Provider: newMockProviderV2(),
			},
			errContains: "must be kebab-case",
		},
		{
			name: "empty Name",
			reg: Registration{
				ID:       "valid-id",
				Name:     "",
				Priority: 1,
				Provider: newMockProviderV2(),
			},
			errContains: "Name is required",
		},
		{
			name: "negative Priority",
			reg: Registration{
				ID:       "valid-id",
				Name:     "Test",
				Priority: -1,
				Provider: newMockProviderV2(),
			},
			errContains: "Priority must be >= 0",
		},
		{
			name: "nil Provider",
			reg: Registration{
				ID:       "valid-id",
				Name:     "Test",
				Priority: 1,
				Provider: nil,
			},
			errContains: "Provider is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistryV2()
			err := registry.Register(tt.reg)

			if err == nil {
				t.Error("expected validation error, got nil")

				return
			}

			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("error should contain %q, got: %v", tt.errContains, err)
			}

			if registry.Count() != 0 {
				t.Error("invalid registration should not be added to registry")
			}
		})
	}
}

func TestRegistryV2_Get(t *testing.T) {
	t.Run("get existing provider", func(t *testing.T) {
		registry := NewRegistryV2()

		expected := Registration{
			ID:       "test-get",
			Name:     "Test Get Provider",
			Priority: 5,
			Provider: newMockProviderV2(),
		}
		_ = registry.Register(expected)

		got, found := registry.Get("test-get")

		if !found {
			t.Error("expected provider to be found")
		}

		if got == nil {
			t.Fatal("expected non-nil registration")
		}

		if got.ID != expected.ID {
			t.Errorf("ID mismatch: expected %s, got %s", expected.ID, got.ID)
		}

		if got.Name != expected.Name {
			t.Errorf("Name mismatch: expected %s, got %s", expected.Name, got.Name)
		}

		if got.Priority != expected.Priority {
			t.Errorf("Priority mismatch: expected %d, got %d", expected.Priority, got.Priority)
		}
	})

	t.Run("get non-existent provider", func(t *testing.T) {
		registry := NewRegistryV2()

		got, found := registry.Get("non-existent")

		if found {
			t.Error("expected provider not to be found")
		}

		if got != nil {
			t.Error("expected nil registration for non-existent provider")
		}
	})
}

func TestRegistryV2_PrioritySorting(t *testing.T) {
	t.Run("All returns priority-sorted registrations", func(t *testing.T) {
		registry := NewRegistryV2()

		// Register in non-priority order
		_ = registry.Register(
			Registration{
				ID:       "low-priority",
				Name:     "Low",
				Priority: 100,
				Provider: newMockProviderV2(),
			},
		)
		_ = registry.Register(
			Registration{
				ID:       "high-priority",
				Name:     "High",
				Priority: 1,
				Provider: newMockProviderV2(),
			},
		)
		_ = registry.Register(
			Registration{
				ID:       "mid-priority",
				Name:     "Mid",
				Priority: 50,
				Provider: newMockProviderV2(),
			},
		)
		_ = registry.Register(
			Registration{
				ID:       "zero-priority",
				Name:     "Zero",
				Priority: 0,
				Provider: newMockProviderV2(),
			},
		)

		all := registry.All()

		if len(all) != 4 {
			t.Fatalf("expected 4 registrations, got %d", len(all))
		}

		// Verify sorted order by priority
		expectedOrder := []string{"zero-priority", "high-priority", "mid-priority", "low-priority"}
		for i, reg := range all {
			if reg.ID != expectedOrder[i] {
				t.Errorf("position %d: expected %s, got %s", i, expectedOrder[i], reg.ID)
			}
		}
	})

	t.Run("IDs returns priority-sorted IDs", func(t *testing.T) {
		registry := NewRegistryV2()

		_ = registry.Register(
			Registration{ID: "charlie", Name: "C", Priority: 30, Provider: newMockProviderV2()},
		)
		_ = registry.Register(
			Registration{ID: "alpha", Name: "A", Priority: 10, Provider: newMockProviderV2()},
		)
		_ = registry.Register(
			Registration{ID: "bravo", Name: "B", Priority: 20, Provider: newMockProviderV2()},
		)

		ids := registry.IDs()

		expectedOrder := []string{"alpha", "bravo", "charlie"}
		for i, id := range ids {
			if id != expectedOrder[i] {
				t.Errorf("position %d: expected %s, got %s", i, expectedOrder[i], id)
			}
		}
	})

	t.Run("deterministic sort for same priority uses ID as tiebreaker", func(t *testing.T) {
		// When providers have the same priority, they should be sorted alphabetically by ID
		registry := NewRegistryV2()

		// Register in reverse alphabetical order to ensure sorting works
		_ = registry.Register(
			Registration{ID: "provider-c", Name: "C", Priority: 10, Provider: newMockProviderV2()},
		)
		_ = registry.Register(
			Registration{ID: "provider-a", Name: "A", Priority: 10, Provider: newMockProviderV2()},
		)
		_ = registry.Register(
			Registration{ID: "provider-b", Name: "B", Priority: 10, Provider: newMockProviderV2()},
		)

		ids := registry.IDs()

		// Should be sorted alphabetically by ID since all have same priority
		expectedOrder := []string{"provider-a", "provider-b", "provider-c"}
		for i, id := range ids {
			if id != expectedOrder[i] {
				t.Errorf("position %d: expected %s, got %s", i, expectedOrder[i], id)
			}
		}

		// Run multiple times to verify consistency
		for i := range 10 {
			currentIDs := registry.IDs()
			for j, id := range currentIDs {
				if id != expectedOrder[j] {
					t.Errorf(
						"inconsistent sort: run %d position %d expected %s, got %s",
						i,
						j,
						expectedOrder[j],
						id,
					)
				}
			}
		}
	})
}

func TestRegistryV2_Count(t *testing.T) {
	t.Run("empty registry", func(t *testing.T) {
		registry := NewRegistryV2()

		if registry.Count() != 0 {
			t.Errorf("expected 0, got %d", registry.Count())
		}
	})

	t.Run("after registrations", func(t *testing.T) {
		registry := NewRegistryV2()

		for i := range 5 {
			_ = registry.Register(Registration{
				ID:       "provider-" + string(rune('a'+i)),
				Name:     "Test",
				Priority: i,
				Provider: newMockProviderV2(),
			})
		}

		if registry.Count() != 5 {
			t.Errorf("expected 5, got %d", registry.Count())
		}
	})
}

func TestRegistryV2_Reset(t *testing.T) {
	registry := NewRegistryV2()

	_ = registry.Register(
		Registration{ID: "provider-1", Name: "P1", Priority: 1, Provider: newMockProviderV2()},
	)
	_ = registry.Register(
		Registration{ID: "provider-2", Name: "P2", Priority: 2, Provider: newMockProviderV2()},
	)

	if registry.Count() != 2 {
		t.Fatalf("expected count 2 before reset, got %d", registry.Count())
	}

	registry.Reset()

	if registry.Count() != 0 {
		t.Errorf("expected count 0 after reset, got %d", registry.Count())
	}

	// Should be able to register again after reset
	err := registry.Register(
		Registration{ID: "provider-1", Name: "P1", Priority: 1, Provider: newMockProviderV2()},
	)
	if err != nil {
		t.Errorf("should be able to register after reset, got: %v", err)
	}
}

func TestGlobalRegistryV2(t *testing.T) {
	// Reset global registry before and after tests
	ResetV2()
	defer ResetV2()

	t.Run("RegisterV2 and GetV2", func(t *testing.T) {
		err := RegisterV2(Registration{
			ID:       "global-test",
			Name:     "Global Test",
			Priority: 1,
			Provider: newMockProviderV2(),
		})
		if err != nil {
			t.Fatalf("RegisterV2 failed: %v", err)
		}

		reg, found := GetV2("global-test")
		if !found {
			t.Error("expected to find registered provider")
		}
		if reg.Name != "Global Test" {
			t.Errorf("expected name 'Global Test', got %s", reg.Name)
		}
	})

	t.Run("AllV2 returns sorted", func(t *testing.T) {
		ResetV2()

		_ = RegisterV2(
			Registration{ID: "second", Name: "Second", Priority: 20, Provider: newMockProviderV2()},
		)
		_ = RegisterV2(
			Registration{ID: "first", Name: "First", Priority: 10, Provider: newMockProviderV2()},
		)

		all := AllV2()
		if len(all) != 2 {
			t.Fatalf("expected 2, got %d", len(all))
		}
		if all[0].ID != "first" {
			t.Errorf("expected 'first' at index 0, got %s", all[0].ID)
		}
	})

	t.Run("IDsV2 returns sorted IDs", func(t *testing.T) {
		ResetV2()

		_ = RegisterV2(
			Registration{ID: "beta", Name: "Beta", Priority: 2, Provider: newMockProviderV2()},
		)
		_ = RegisterV2(
			Registration{ID: "alpha", Name: "Alpha", Priority: 1, Provider: newMockProviderV2()},
		)

		ids := IDsV2()
		if len(ids) != 2 {
			t.Fatalf("expected 2, got %d", len(ids))
		}
		if ids[0] != "alpha" {
			t.Errorf("expected 'alpha' at index 0, got %s", ids[0])
		}
	})

	t.Run("CountV2", func(t *testing.T) {
		ResetV2()

		if CountV2() != 0 {
			t.Error("expected 0 after reset")
		}

		_ = RegisterV2(
			Registration{ID: "test", Name: "Test", Priority: 1, Provider: newMockProviderV2()},
		)

		if CountV2() != 1 {
			t.Error("expected 1 after registration")
		}
	})

	t.Run("duplicate rejection in global registry", func(t *testing.T) {
		ResetV2()

		_ = RegisterV2(
			Registration{ID: "unique", Name: "First", Priority: 1, Provider: newMockProviderV2()},
		)
		err := RegisterV2(
			Registration{ID: "unique", Name: "Second", Priority: 2, Provider: newMockProviderV2()},
		)

		if err == nil {
			t.Error("expected error for duplicate registration")
		}
	})
}

func TestRegistryV2_ThreadSafety(_ *testing.T) {
	// This test verifies that concurrent operations don't cause races
	// Run with -race flag: go test -race ./...
	registry := NewRegistryV2()

	// Pre-register some providers
	for i := range 10 {
		_ = registry.Register(Registration{
			ID:       "provider-" + string(rune('a'+i)),
			Name:     "Test",
			Priority: i,
			Provider: newMockProviderV2(),
		})
	}

	done := make(chan bool)

	// Concurrent reads
	for range 10 {
		go func() {
			for range 100 {
				_ = registry.All()
				_ = registry.IDs()
				_ = registry.Count()
				_, _ = registry.Get("provider-a")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for range 10 {
		<-done
	}
}
