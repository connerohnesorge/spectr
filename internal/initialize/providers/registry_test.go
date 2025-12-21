package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

type mockProvider struct{}

func (m *mockProvider) Initializers() []types.Initializer { return nil }

func TestRegistryV2(t *testing.T) {
	r := NewRegistry()

	reg1 := Registration{ID: "p1", Name: "Provider 1", Priority: 10, Provider: &mockProvider{}}
	reg2 := Registration{ID: "p2", Name: "Provider 2", Priority: 5, Provider: &mockProvider{}}

	// Register
	assert.NoError(t, r.Register(reg1))
	assert.NoError(t, r.Register(reg2))

	// Duplicate
	assert.Error(t, r.Register(reg1))

	// Get
	p, ok := r.Get("p1")
	assert.True(t, ok)
	assert.Equal(t, "Provider 1", p.Name)

	_, ok = r.Get("none")
	assert.False(t, ok)

	// All (sorted by priority)
	all := r.All()
	assert.Len(t, all, 2)
	assert.Equal(t, "p2", all[0].ID) // Priority 5
	assert.Equal(t, "p1", all[1].ID) // Priority 10
}
