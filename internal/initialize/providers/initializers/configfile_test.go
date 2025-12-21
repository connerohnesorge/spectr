package initializers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func TestConfigFileInitializer(t *testing.T) {
	fs := afero.NewMemMapFs()
	initializer := NewConfigFileInitializer("CLAUDE.md")

	mockTM := &MockTemplateRenderer{
		RenderInstructionPointerFunc: func(ctx types.TemplateContext) (string, error) {
			return "Check AGENTS.md", nil
		},
	}

	// Test Init - Create
	err := initializer.Init(context.Background(), fs, fs, nil, mockTM)
	assert.NoError(t, err)

	exists, _ := afero.Exists(fs, "CLAUDE.md")
	assert.True(t, exists)

	content, _ := afero.ReadFile(fs, "CLAUDE.md")
	assert.Contains(t, string(content), types.SpectrStartMarker)
	assert.Contains(t, string(content), "Check AGENTS.md")

	// Test IsSetup
	setup, err := initializer.IsSetup(fs, fs, nil)
	assert.NoError(t, err)
	assert.True(t, setup)
}