package initializers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func TestSlashCommandsInitializer(t *testing.T) {
	fs := afero.NewMemMapFs()
	initializer := NewSlashCommandsInitializer("proposal", "commands/proposal.md", "---\nfront: matter\n---")

	mockTM := &MockTemplateRenderer{
		RenderSlashCommandFunc: func(cmd string, ctx types.TemplateContext) (string, error) {
			return "Slash Command Content", nil
		},
	}

	// Test Init - Create
	err := initializer.Init(context.Background(), fs, fs, nil, mockTM)
	assert.NoError(t, err)

	exists, _ := afero.Exists(fs, "commands/proposal.md")
	assert.True(t, exists)

	content, _ := afero.ReadFile(fs, "commands/proposal.md")
	assert.Contains(t, string(content), "front: matter")
	assert.Contains(t, string(content), types.SpectrStartMarker)
	assert.Contains(t, string(content), "Slash Command Content")

	// Test IsSetup
	setup, err := initializer.IsSetup(fs, fs, nil)
	assert.NoError(t, err)
	assert.True(t, setup)
}