// Package providers contains individual AI tool provider implementations.
// This file implements the Cline provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// ClineProvider configures Cline support by creating the
// configuration file. This provider automatically registers
// itself with the global provider registry during initialization.
type ClineProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "cline",
			Name:           "Cline",
			ConfigFilePath: "CLINE.md",
			SlashID:        "cline-slash",
			Priority:       PriorityCline,
		}),
		func() providerkit.Provider {
			return &ClineProvider{}
		},
	)
}

func (*ClineProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "CLINE.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*ClineProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CLINE.md")

	return providerkit.FileExists(filePath)
}

func (*ClineProvider) GetName() string {
	return "Cline"
}
