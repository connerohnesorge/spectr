package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllProvidersRegistered(t *testing.T) {
	all := All()
	assert.Len(t, all, 15)
}

func TestProviderInitializers(t *testing.T) {
	tests := []struct {
		id          string
		expectInits int
	}{
		{"claude-code", 3},
		{"gemini", 2},
		{"cursor", 2},
		{"cline", 2},
		{"aider", 2},
		{"codex", 3},
		{"costrict", 3},
		{"qoder", 3},
		{"qwen", 3},
		{"antigravity", 3},
		{"windsurf", 2},
		{"kilocode", 2},
		{"continue", 2},
		{"crush", 3},
		{"opencode", 2},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			reg, ok := Get(tt.id)
			assert.True(t, ok)
			inits := reg.Provider.Initializers()
			assert.Len(t, inits, tt.expectInits)
		})
	}
}

func TestProviderMetadata(t *testing.T) {
	reg, ok := Get("claude-code")
	assert.True(t, ok)
	assert.Equal(t, "Claude Code", reg.Name)
	assert.Equal(t, 1, reg.Priority)
}
