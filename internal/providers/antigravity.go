// Package providers contains individual AI tool provider implementations.
// This file implements the Antigravity IDE provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// AntigravityProvider configures Antigravity IDE support by creating
// the AGENTS.md instruction file. This provider automatically registers
// itself with the global provider registry during initialization.
type AntigravityProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "antigravity",
			Name:           "Antigravity",
			ConfigFilePath: "AGENTS.md",
			SlashID:        "antigravity-slash",
			Priority:       PriorityAntigravity,
		}),
		func() providerkit.Provider {
			return &AntigravityProvider{}
		},
	)
}

func (*AntigravityProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "AGENTS.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*AntigravityProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "AGENTS.md")

	return providerkit.FileExists(filePath)
}

func (*AntigravityProvider) GetName() string {
	return "Antigravity"
}
