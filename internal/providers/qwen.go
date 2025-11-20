// Package providers contains individual AI tool provider implementations.
// This file implements the Qwen provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// QwenProvider configures Qwen support by creating the
// configuration file. This provider automatically registers
// itself with the global provider registry during initialization.
type QwenProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "qwen",
			Name:           "Qwen Code",
			ConfigFilePath: "QWEN.md",
			SlashID:        "qwen-slash",
			Priority:       PriorityQwen,
		}),
		func() providerkit.Provider {
			return &QwenProvider{}
		},
	)
}

func (*QwenProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "QWEN.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*QwenProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "QWEN.md")

	return providerkit.FileExists(filePath)
}

func (*QwenProvider) GetName() string {
	return "Qwen Code"
}
