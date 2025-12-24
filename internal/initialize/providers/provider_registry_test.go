package providers

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// mockProvider is a simple test implementation of the Provider interface
type mockProvider struct {
	id string
}

func (*mockProvider) Initializers(
	_ context.Context,
) []Initializer {
	return nil
}

func TestNewProviderRegistry(t *testing.T) {
	registry := NewProviderRegistry()

	if registry == nil {
		t.Fatal(
			"NewProviderRegistry() returned nil",
		)
	}

	if registry.Count() != 0 {
		t.Errorf(
			"NewProviderRegistry() should start empty, got count %d",
			registry.Count(),
		)
	}
}

func TestRegisterProvider(t *testing.T) {
	registry := NewProviderRegistry()

	reg := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &mockProvider{
			id: "test-provider",
		},
	}

	err := registry.RegisterProvider(reg)
	if err != nil {
		t.Fatalf(
			"RegisterProvider() failed: %v",
			err,
		)
	}

	if registry.Count() != 1 {
		t.Errorf(
			"Expected count 1 after registration, got %d",
			registry.Count(),
		)
	}
}

func TestRegisterProvider_DuplicateID(
	t *testing.T,
) {
	registry := NewProviderRegistry()

	reg1 := Registration{
		ID:       "duplicate-id",
		Name:     "First Provider",
		Priority: 1,
		Provider: &mockProvider{
			id: "duplicate-id",
		},
	}

	reg2 := Registration{
		ID:       "duplicate-id",
		Name:     "Second Provider",
		Priority: 2,
		Provider: &mockProvider{
			id: "duplicate-id",
		},
	}

	// First registration should succeed
	err := registry.RegisterProvider(reg1)
	if err != nil {
		t.Fatalf(
			"First RegisterProvider() failed: %v",
			err,
		)
	}

	// Second registration with same ID should fail
	err = registry.RegisterProvider(reg2)
	if err == nil {
		t.Fatal(
			"RegisterProvider() should have failed for duplicate ID",
		)
	}

	expectedError := `provider "duplicate-id" already registered`
	if err.Error() != expectedError {
		t.Errorf(
			"Expected error %q, got %q",
			expectedError,
			err.Error(),
		)
	}

	// Count should still be 1
	if registry.Count() != 1 {
		t.Errorf(
			"Expected count 1 after duplicate rejection, got %d",
			registry.Count(),
		)
	}
}

func TestProviderRegistry_Get(t *testing.T) {
	registry := NewProviderRegistry()

	reg := Registration{
		ID:       "get-test",
		Name:     "Get Test Provider",
		Priority: 5,
		Provider: &mockProvider{id: "get-test"},
	}

	err := registry.RegisterProvider(reg)
	if err != nil {
		t.Fatalf(
			"RegisterProvider() failed: %v",
			err,
		)
	}

	// Test successful retrieval
	retrieved, ok := registry.Get("get-test")
	if !ok {
		t.Fatal(
			"Get() should have found the provider",
		)
	}

	if retrieved.ID != reg.ID {
		t.Errorf(
			"Expected ID %q, got %q",
			reg.ID,
			retrieved.ID,
		)
	}
	if retrieved.Name != reg.Name {
		t.Errorf(
			"Expected Name %q, got %q",
			reg.Name,
			retrieved.Name,
		)
	}
	if retrieved.Priority != reg.Priority {
		t.Errorf(
			"Expected Priority %d, got %d",
			reg.Priority,
			retrieved.Priority,
		)
	}
}

func TestProviderRegistry_Get_NotFound(
	t *testing.T,
) {
	registry := NewProviderRegistry()

	_, ok := registry.Get("nonexistent")
	if ok {
		t.Error(
			"Get() should have returned false for nonexistent provider",
		)
	}
}

func TestProviderRegistry_All_PrioritySorting(
	t *testing.T,
) {
	registry := NewProviderRegistry()

	// Register providers in non-priority order
	providers := []Registration{
		{
			ID:       "medium",
			Name:     "Medium Priority",
			Priority: 5,
			Provider: &mockProvider{id: "medium"},
		},
		{
			ID:       "high",
			Name:     "High Priority",
			Priority: 1,
			Provider: &mockProvider{id: "high"},
		},
		{
			ID:       "low",
			Name:     "Low Priority",
			Priority: 10,
			Provider: &mockProvider{id: "low"},
		},
		{
			ID:       "higher",
			Name:     "Higher Priority",
			Priority: 2,
			Provider: &mockProvider{id: "higher"},
		},
	}

	for _, p := range providers {
		if err := registry.RegisterProvider(p); err != nil {
			t.Fatalf(
				"RegisterProvider() failed: %v",
				err,
			)
		}
	}

	all := registry.All()

	// Verify count
	if len(all) != 4 {
		t.Fatalf(
			"Expected 4 providers, got %d",
			len(all),
		)
	}

	// Verify sorting by priority (ascending)
	expectedOrder := []string{
		"high",
		"higher",
		"medium",
		"low",
	}
	for i, expected := range expectedOrder {
		if all[i].ID != expected {
			t.Errorf(
				"Position %d: expected ID %q, got %q",
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
				"Priorities not sorted: %d > %d at positions %d and %d",
				all[i-1].Priority,
				all[i].Priority,
				i-1,
				i,
			)
		}
	}
}

func TestProviderRegistry_IDs(t *testing.T) {
	registry := NewProviderRegistry()

	// Register providers with different priorities
	providers := []Registration{
		{
			ID:       "provider-c",
			Name:     "Provider C",
			Priority: 3,
			Provider: &mockProvider{
				id: "provider-c",
			},
		},
		{
			ID:       "provider-a",
			Name:     "Provider A",
			Priority: 1,
			Provider: &mockProvider{
				id: "provider-a",
			},
		},
		{
			ID:       "provider-b",
			Name:     "Provider B",
			Priority: 2,
			Provider: &mockProvider{
				id: "provider-b",
			},
		},
	}

	for _, p := range providers {
		if err := registry.RegisterProvider(p); err != nil {
			t.Fatalf(
				"RegisterProvider() failed: %v",
				err,
			)
		}
	}

	ids := registry.IDs()

	// Verify count
	if len(ids) != 3 {
		t.Fatalf(
			"Expected 3 IDs, got %d",
			len(ids),
		)
	}

	// Verify sorted by priority
	expectedOrder := []string{
		"provider-a",
		"provider-b",
		"provider-c",
	}
	for i, expected := range expectedOrder {
		if ids[i] != expected {
			t.Errorf(
				"Position %d: expected ID %q, got %q",
				i,
				expected,
				ids[i],
			)
		}
	}
}

func TestProviderRegistry_Count(t *testing.T) {
	registry := NewProviderRegistry()

	if registry.Count() != 0 {
		t.Errorf(
			"Expected count 0 for empty registry, got %d",
			registry.Count(),
		)
	}

	// Add providers one by one and verify count
	for i := 1; i <= 5; i++ {
		reg := Registration{
			ID:       string(rune('a' + i - 1)),
			Name:     "Provider",
			Priority: i,
			Provider: &mockProvider{},
		}
		if err := registry.RegisterProvider(reg); err != nil {
			t.Fatalf(
				"RegisterProvider() failed: %v",
				err,
			)
		}

		if registry.Count() != i {
			t.Errorf(
				"Expected count %d, got %d",
				i,
				registry.Count(),
			)
		}
	}
}

func TestGlobalRegistry_RegisterProvider(
	t *testing.T,
) {
	// Reset global registry for test isolation
	ResetProviders()

	reg := Registration{
		ID:       "global-test",
		Name:     "Global Test Provider",
		Priority: 1,
		Provider: &mockProvider{
			id: "global-test",
		},
	}

	err := RegisterProvider(reg)
	if err != nil {
		t.Fatalf(
			"RegisterProvider() failed: %v",
			err,
		)
	}

	if ProviderCount() != 1 {
		t.Errorf(
			"Expected global count 1, got %d",
			ProviderCount(),
		)
	}

	// Cleanup
	ResetProviders()
}

func TestGlobalRegistry_GetProvider(
	t *testing.T,
) {
	ResetProviders()

	reg := Registration{
		ID:       "global-get-test",
		Name:     "Global Get Test",
		Priority: 5,
		Provider: &mockProvider{
			id: "global-get-test",
		},
	}

	if err := RegisterProvider(reg); err != nil {
		t.Fatalf(
			"RegisterProvider() failed: %v",
			err,
		)
	}

	retrieved, ok := GetProvider(
		"global-get-test",
	)
	if !ok {
		t.Fatal(
			"GetProvider() should have found the provider",
		)
	}

	if retrieved.ID != reg.ID {
		t.Errorf(
			"Expected ID %q, got %q",
			reg.ID,
			retrieved.ID,
		)
	}

	// Cleanup
	ResetProviders()
}

func TestGlobalRegistry_AllProviders(
	t *testing.T,
) {
	ResetProviders()

	providers := []Registration{
		{
			ID:       "global-a",
			Name:     "Global A",
			Priority: 2,
			Provider: &mockProvider{},
		},
		{
			ID:       "global-b",
			Name:     "Global B",
			Priority: 1,
			Provider: &mockProvider{},
		},
	}

	for _, p := range providers {
		if err := RegisterProvider(p); err != nil {
			t.Fatalf(
				"RegisterProvider() failed: %v",
				err,
			)
		}
	}

	all := AllProviders()
	if len(all) != 2 {
		t.Fatalf(
			"Expected 2 providers, got %d",
			len(all),
		)
	}

	// Verify sorted by priority
	if all[0].ID != "global-b" ||
		all[1].ID != "global-a" {
		t.Errorf(
			"Expected [global-b, global-a], got [%s, %s]",
			all[0].ID,
			all[1].ID,
		)
	}

	// Cleanup
	ResetProviders()
}

func TestGlobalRegistry_ProviderIDs(
	t *testing.T,
) {
	ResetProviders()

	providers := []Registration{
		{
			ID:       "id-c",
			Name:     "Provider C",
			Priority: 3,
			Provider: &mockProvider{},
		},
		{
			ID:       "id-a",
			Name:     "Provider A",
			Priority: 1,
			Provider: &mockProvider{},
		},
		{
			ID:       "id-b",
			Name:     "Provider B",
			Priority: 2,
			Provider: &mockProvider{},
		},
	}

	for _, p := range providers {
		if err := RegisterProvider(p); err != nil {
			t.Fatalf(
				"RegisterProvider() failed: %v",
				err,
			)
		}
	}

	ids := ProviderIDs()
	expectedOrder := []string{
		"id-a",
		"id-b",
		"id-c",
	}

	if len(ids) != len(expectedOrder) {
		t.Fatalf(
			"Expected %d IDs, got %d",
			len(expectedOrder),
			len(ids),
		)
	}

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

	// Cleanup
	ResetProviders()
}

func TestGlobalRegistry_ResetProviders(
	t *testing.T,
) {
	ResetProviders()

	// Add a provider
	reg := Registration{
		ID:       "reset-test",
		Name:     "Reset Test",
		Priority: 1,
		Provider: &mockProvider{},
	}

	if err := RegisterProvider(reg); err != nil {
		t.Fatalf(
			"RegisterProvider() failed: %v",
			err,
		)
	}

	if ProviderCount() != 1 {
		t.Errorf(
			"Expected count 1 before reset, got %d",
			ProviderCount(),
		)
	}

	// Reset the registry
	ResetProviders()

	if ProviderCount() != 0 {
		t.Errorf(
			"Expected count 0 after reset, got %d",
			ProviderCount(),
		)
	}

	// Verify we can register again after reset
	if err := RegisterProvider(reg); err != nil {
		t.Fatalf(
			"RegisterProvider() failed after reset: %v",
			err,
		)
	}

	if ProviderCount() != 1 {
		t.Errorf(
			"Expected count 1 after re-registration, got %d",
			ProviderCount(),
		)
	}

	// Final cleanup
	ResetProviders()
}

func TestConcurrentRegistration(t *testing.T) {
	registry := NewProviderRegistry()

	// Test thread safety with concurrent registrations
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			reg := Registration{
				ID: fmt.Sprintf(
					"concurrent-%d",
					id,
				),
				Name: fmt.Sprintf(
					"Concurrent Provider %d",
					id,
				),
				Priority: id,
				Provider: &mockProvider{},
			}

			if err := registry.RegisterProvider(reg); err != nil {
				t.Errorf(
					"Concurrent RegisterProvider() failed: %v",
					err,
				)
			}
		}(i)
	}

	wg.Wait()

	if registry.Count() != numGoroutines {
		t.Errorf(
			"Expected count %d after concurrent registrations, got %d",
			numGoroutines,
			registry.Count(),
		)
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	registry := NewProviderRegistry()

	// Register initial provider
	reg := Registration{
		ID:       "read-write-test",
		Name:     "Read Write Test",
		Priority: 1,
		Provider: &mockProvider{},
	}
	if err := registry.RegisterProvider(reg); err != nil {
		t.Fatalf(
			"Initial registration failed: %v",
			err,
		)
	}

	var wg sync.WaitGroup
	numReaders := 5
	numWriters := 5

	// Spawn readers
	for range numReaders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				_ = registry.All()
				_ = registry.IDs()
				_, _ = registry.Get(
					"read-write-test",
				)
				_ = registry.Count()
			}
		}()
	}

	// Spawn writers
	for i := range numWriters {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range 10 {
				reg := Registration{
					ID: fmt.Sprintf(
						"writer-%d-%d",
						id,
						j,
					),
					Name:     "Writer",
					Priority: id*100 + j,
					Provider: &mockProvider{},
				}
				_ = registry.RegisterProvider(reg)
			}
		}(i)
	}

	wg.Wait()

	// Verify registry is still consistent
	all := registry.All()
	if len(all) == 0 {
		t.Error(
			"Registry should not be empty after concurrent operations",
		)
	}

	// Verify sorting is still correct
	for i := 1; i < len(all); i++ {
		if all[i-1].Priority > all[i].Priority {
			t.Errorf(
				"Priorities not sorted after concurrent ops: %d > %d",
				all[i-1].Priority,
				all[i].Priority,
			)
		}
	}
}

func TestProviderRegistry_All_ImmutableReturn(
	t *testing.T,
) {
	registry := NewProviderRegistry()

	reg := Registration{
		ID:       "immutable-test",
		Name:     "Immutable Test",
		Priority: 1,
		Provider: &mockProvider{},
	}

	if err := registry.RegisterProvider(reg); err != nil {
		t.Fatalf(
			"RegisterProvider() failed: %v",
			err,
		)
	}

	// Get all providers
	all1 := registry.All()

	// Modify the returned slice
	all1[0].Priority = 999

	// Get all providers again
	all2 := registry.All()

	// Verify the registry wasn't affected by the modification
	if all2[0].Priority == 999 {
		t.Error(
			"Modifying returned slice should not affect the registry",
		)
	}

	if all2[0].Priority != 1 {
		t.Errorf(
			"Expected priority 1, got %d",
			all2[0].Priority,
		)
	}
}
