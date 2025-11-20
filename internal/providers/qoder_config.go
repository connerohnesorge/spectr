// Package providers contains individual AI tool provider implementations.
// This file implements the Qoder provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// QoderProvider configures Qoder support by creating the
// configuration file. This provider automatically registers
// itself with the global provider registry during initialization.
type QoderProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "qoder-config",
			Name:           "Qoder",
			ConfigFilePath: "QODER.md",
			SlashID:        "qoder-slash",
			Priority:       PriorityQoder,
		}),
		func() providerkit.Provider {
			return &QoderProvider{}
		},
	)
}

func (*QoderProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "QODER.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*QoderProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "QODER.md")

	return providerkit.FileExists(filePath)
}

func (*QoderProvider) GetName() string {
	return "Qoder"
}
