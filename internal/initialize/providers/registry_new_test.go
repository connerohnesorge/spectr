package providers

import (
	"context"
	"testing"
)

// mockProvider is a simple mock implementation of NewProvider for testing.
type mockProvider struct{}

func (*mockProvider) Initializers(
	_ context.Context,
) []Initializer {
	return nil // Simple mock - no initializers needed for registry tests
}

func TestCreateRegistry(t *testing.T) {
	reg := CreateRegistry()

	if reg == nil {
		t.Fatal("CreateRegistry() returned nil")
	}

	if reg.Count() != 0 {
		t.Errorf(
			"CreateRegistry() should return empty registry, got %d providers",
			reg.Count(),
		)
	}

	if len(reg.All()) != 0 {
		t.Errorf(
			"CreateRegistry().All() should return empty slice, got %d items",
			len(reg.All()),
		)
	}

	if len(reg.IDs()) != 0 {
		t.Errorf(
			"CreateRegistry().IDs() should return empty slice, got %d items",
			len(reg.IDs()),
		)
	}
}

func TestProviderRegistry_Register(t *testing.T) {
	reg := CreateRegistry()

	registration := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &mockProvider{},
	}

	err := reg.Register(registration)
	if err != nil {
		t.Errorf("Register() failed: %v", err)
	}

	if reg.Count() != 1 {
		t.Errorf(
			"Count() = %d, want 1",
			reg.Count(),
		)
	}
}

func TestProviderRegistry_Register_Duplicate(
	t *testing.T,
) {
	reg := CreateRegistry()

	registration := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &mockProvider{},
	}

	// First registration should succeed
	err := reg.Register(registration)
	if err != nil {
		t.Fatalf(
			"First Register() failed: %v",
			err,
		)
	}

	// Second registration with same ID should fail
	err = reg.Register(registration)
	if err == nil {
		t.Error(
			"Register() should return error for duplicate ID",
		)
	}

	// Count should still be 1
	if reg.Count() != 1 {
		t.Errorf(
			"Count() = %d, want 1 (duplicate should not be added)",
			reg.Count(),
		)
	}
}

func TestProviderRegistry_Get(t *testing.T) {
	reg := CreateRegistry()

	registration := Registration{
		ID:       "test-provider",
		Name:     "Test Provider",
		Priority: 1,
		Provider: &mockProvider{},
	}

	err := reg.Register(registration)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	// Get existing provider
	got := reg.Get("test-provider")
	if got == nil {
		t.Fatal(
			"Get() returned nil for registered provider",
		)
	}

	if got.ID != "test-provider" {
		t.Errorf(
			"Get().ID = %s, want test-provider",
			got.ID,
		)
	}

	if got.Name != "Test Provider" {
		t.Errorf(
			"Get().Name = %s, want Test Provider",
			got.Name,
		)
	}

	if got.Priority != 1 {
		t.Errorf(
			"Get().Priority = %d, want 1",
			got.Priority,
		)
	}
}

func TestProviderRegistry_Get_Unknown(
	t *testing.T,
) {
	reg := CreateRegistry()

	// Get non-existent provider should return nil
	got := reg.Get("nonexistent")
	if got != nil {
		t.Errorf(
			"Get() for unknown ID should return nil, got %+v",
			got,
		)
	}
}

func TestProviderRegistry_All_PrioritySorting(
	t *testing.T,
) {
	reg := CreateRegistry()

	// Register providers in non-priority order
	registrations := []Registration{
		{
			ID:       "low-priority",
			Name:     "Low Priority",
			Priority: 100,
			Provider: &mockProvider{},
		},
		{
			ID:       "high-priority",
			Name:     "High Priority",
			Priority: 1,
			Provider: &mockProvider{},
		},
		{
			ID:       "medium-priority",
			Name:     "Medium Priority",
			Priority: 50,
			Provider: &mockProvider{},
		},
	}

	for _, r := range registrations {
		if err := reg.Register(r); err != nil {
			t.Fatalf(
				"Register() failed for %s: %v",
				r.ID,
				err,
			)
		}
	}

	all := reg.All()

	if len(all) != 3 {
		t.Fatalf(
			"All() returned %d items, want 3",
			len(all),
		)
	}

	// Verify sorted by priority (ascending - lower number = higher priority)
	expectedOrder := []string{
		"high-priority",
		"medium-priority",
		"low-priority",
	}
	for i, expected := range expectedOrder {
		if all[i].ID != expected {
			t.Errorf(
				"All()[%d].ID = %s, want %s",
				i,
				all[i].ID,
				expected,
			)
		}
	}

	// Verify priority values are ascending
	for i := 1; i < len(all); i++ {
		if all[i-1].Priority > all[i].Priority {
			t.Errorf(
				"All() not sorted by priority: %s (priority %d) should come before %s (priority %d)",
				all[i].ID,
				all[i].Priority,
				all[i-1].ID,
				all[i-1].Priority,
			)
		}
	}
}

func TestProviderRegistry_IDs_PrioritySorting(
	t *testing.T,
) {
	reg := CreateRegistry()

	// Register providers in non-priority order
	registrations := []Registration{
		{
			ID:       "third",
			Name:     "Third",
			Priority: 30,
			Provider: &mockProvider{},
		},
		{
			ID:       "first",
			Name:     "First",
			Priority: 10,
			Provider: &mockProvider{},
		},
		{
			ID:       "second",
			Name:     "Second",
			Priority: 20,
			Provider: &mockProvider{},
		},
	}

	for _, r := range registrations {
		if err := reg.Register(r); err != nil {
			t.Fatalf(
				"Register() failed for %s: %v",
				r.ID,
				err,
			)
		}
	}

	ids := reg.IDs()

	if len(ids) != 3 {
		t.Fatalf(
			"IDs() returned %d items, want 3",
			len(ids),
		)
	}

	// Verify sorted by priority
	expectedOrder := []string{
		"first",
		"second",
		"third",
	}
	for i, expected := range expectedOrder {
		if ids[i] != expected {
			t.Errorf(
				"IDs()[%d] = %s, want %s",
				i,
				ids[i],
				expected,
			)
		}
	}
}

func TestProviderRegistry_Count(t *testing.T) {
	reg := CreateRegistry()

	// Empty registry
	if reg.Count() != 0 {
		t.Errorf(
			"Count() on empty registry = %d, want 0",
			reg.Count(),
		)
	}

	// Add one provider
	err := reg.Register(Registration{
		ID:       "provider1",
		Name:     "Provider 1",
		Priority: 1,
		Provider: &mockProvider{},
	})
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	if reg.Count() != 1 {
		t.Errorf(
			"Count() after one registration = %d, want 1",
			reg.Count(),
		)
	}

	// Add another provider
	err = reg.Register(Registration{
		ID:       "provider2",
		Name:     "Provider 2",
		Priority: 2,
		Provider: &mockProvider{},
	})
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	if reg.Count() != 2 {
		t.Errorf(
			"Count() after two registrations = %d, want 2",
			reg.Count(),
		)
	}
}

func TestProviderRegistry_EmptyRegistry(
	t *testing.T,
) {
	reg := CreateRegistry()

	// All operations should work on empty registry
	if reg.Count() != 0 {
		t.Errorf(
			"Count() = %d, want 0",
			reg.Count(),
		)
	}

	all := reg.All()
	if all == nil {
		t.Error(
			"All() returned nil, want empty slice",
		)
	}
	if len(all) != 0 {
		t.Errorf(
			"All() returned %d items, want 0",
			len(all),
		)
	}

	ids := reg.IDs()
	if ids == nil {
		t.Error(
			"IDs() returned nil, want empty slice",
		)
	}
	if len(ids) != 0 {
		t.Errorf(
			"IDs() returned %d items, want 0",
			len(ids),
		)
	}

	got := reg.Get("anything")
	if got != nil {
		t.Errorf(
			"Get() on empty registry should return nil, got %+v",
			got,
		)
	}
}

func TestProviderRegistry_SingleRegistration(
	t *testing.T,
) {
	reg := CreateRegistry()

	registration := Registration{
		ID:       "only-provider",
		Name:     "Only Provider",
		Priority: 42,
		Provider: &mockProvider{},
	}

	err := reg.Register(registration)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	// Verify count
	if reg.Count() != 1 {
		t.Errorf(
			"Count() = %d, want 1",
			reg.Count(),
		)
	}

	// Verify All returns the single registration
	all := reg.All()
	if len(all) != 1 {
		t.Fatalf(
			"All() returned %d items, want 1",
			len(all),
		)
	}
	if all[0].ID != "only-provider" {
		t.Errorf(
			"All()[0].ID = %s, want only-provider",
			all[0].ID,
		)
	}

	// Verify IDs returns the single ID
	ids := reg.IDs()
	if len(ids) != 1 {
		t.Fatalf(
			"IDs() returned %d items, want 1",
			len(ids),
		)
	}
	if ids[0] != "only-provider" {
		t.Errorf(
			"IDs()[0] = %s, want only-provider",
			ids[0],
		)
	}

	// Verify Get returns the registration
	got := reg.Get("only-provider")
	if got == nil {
		t.Fatal("Get() returned nil")
	}
	if got.ID != "only-provider" {
		t.Errorf(
			"Get().ID = %s, want only-provider",
			got.ID,
		)
	}
	if got.Name != "Only Provider" {
		t.Errorf(
			"Get().Name = %s, want Only Provider",
			got.Name,
		)
	}
	if got.Priority != 42 {
		t.Errorf(
			"Get().Priority = %d, want 42",
			got.Priority,
		)
	}
}

func TestProviderRegistry_MultipleRegistrations(
	t *testing.T,
) {
	reg := CreateRegistry()

	registrations := []Registration{
		{
			ID:       "claude-code",
			Name:     "Claude Code",
			Priority: 1,
			Provider: &mockProvider{},
		},
		{
			ID:       "cursor",
			Name:     "Cursor",
			Priority: 2,
			Provider: &mockProvider{},
		},
		{
			ID:       "cline",
			Name:     "Cline",
			Priority: 3,
			Provider: &mockProvider{},
		},
		{
			ID:       "gemini",
			Name:     "Gemini CLI",
			Priority: 4,
			Provider: &mockProvider{},
		},
		{
			ID:       "roo-code",
			Name:     "Roo Code",
			Priority: 5,
			Provider: &mockProvider{},
		},
	}

	for _, r := range registrations {
		if err := reg.Register(r); err != nil {
			t.Fatalf(
				"Register() failed for %s: %v",
				r.ID,
				err,
			)
		}
	}

	// Verify count
	if reg.Count() != 5 {
		t.Errorf(
			"Count() = %d, want 5",
			reg.Count(),
		)
	}

	// Verify all are retrievable
	for _, r := range registrations {
		got := reg.Get(r.ID)
		if got == nil {
			t.Errorf("Get(%s) returned nil", r.ID)

			continue
		}
		if got.ID != r.ID {
			t.Errorf(
				"Get(%s).ID = %s, want %s",
				r.ID,
				got.ID,
				r.ID,
			)
		}
		if got.Name != r.Name {
			t.Errorf(
				"Get(%s).Name = %s, want %s",
				r.ID,
				got.Name,
				r.Name,
			)
		}
		if got.Priority != r.Priority {
			t.Errorf(
				"Get(%s).Priority = %d, want %d",
				r.ID,
				got.Priority,
				r.Priority,
			)
		}
	}

	// Verify All returns all in priority order
	all := reg.All()
	if len(all) != 5 {
		t.Fatalf(
			"All() returned %d items, want 5",
			len(all),
		)
	}

	expectedIDs := []string{
		"claude-code",
		"cursor",
		"cline",
		"gemini",
		"roo-code",
	}
	for i, expectedID := range expectedIDs {
		if all[i].ID != expectedID {
			t.Errorf(
				"All()[%d].ID = %s, want %s",
				i,
				all[i].ID,
				expectedID,
			)
		}
	}

	// Verify IDs returns all IDs in priority order
	ids := reg.IDs()
	if len(ids) != 5 {
		t.Fatalf(
			"IDs() returned %d items, want 5",
			len(ids),
		)
	}
	for i, expectedID := range expectedIDs {
		if ids[i] != expectedID {
			t.Errorf(
				"IDs()[%d] = %s, want %s",
				i,
				ids[i],
				expectedID,
			)
		}
	}
}

func TestProviderRegistry_DuplicateRejectionPreservesOriginal(
	t *testing.T,
) {
	reg := CreateRegistry()

	// Register original
	original := Registration{
		ID:       "my-provider",
		Name:     "Original Name",
		Priority: 1,
		Provider: &mockProvider{},
	}
	if err := reg.Register(original); err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	// Attempt to register duplicate with different values
	duplicate := Registration{
		ID:       "my-provider",
		Name:     "Duplicate Name",
		Priority: 99,
		Provider: &mockProvider{},
	}
	err := reg.Register(duplicate)
	if err == nil {
		t.Error(
			"Register() should return error for duplicate ID",
		)
	}

	// Verify original values are preserved
	got := reg.Get("my-provider")
	if got == nil {
		t.Fatal("Get() returned nil")
	}
	if got.Name != "Original Name" {
		t.Errorf(
			"Get().Name = %s, want Original Name (original should be preserved)",
			got.Name,
		)
	}
	if got.Priority != 1 {
		t.Errorf(
			"Get().Priority = %d, want 1 (original should be preserved)",
			got.Priority,
		)
	}
}
