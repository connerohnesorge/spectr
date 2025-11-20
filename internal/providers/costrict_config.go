// Package providers contains individual AI tool provider implementations.
// This file implements the CoStrict provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// CostrictProvider configures CoStrict support by creating the
// COSTRICT.md instruction file. This provider automatically registers
// itself with the global provider registry during initialization.
type CostrictProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "costrict-config",
			Name:           "CoStrict",
			ConfigFilePath: "COSTRICT.md",
			SlashID:        "costrict-slash",
			Priority:       PriorityCostrict,
		}),
		func() providerkit.Provider {
			return &CostrictProvider{}
		},
	)
}

func (*CostrictProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "COSTRICT.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*CostrictProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "COSTRICT.md")

	return providerkit.FileExists(filePath)
}

func (*CostrictProvider) GetName() string {
	return "CoStrict"
}
