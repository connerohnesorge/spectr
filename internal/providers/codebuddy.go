// Package providers contains individual AI tool provider implementations.
// This file implements the CodeBuddy provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// CodeBuddyProvider configures CodeBuddy support by creating the
// CODEBUDDY.md instruction file. This provider automatically registers
// itself with the global provider registry during initialization.
type CodeBuddyProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "codebuddy",
			Name:           "CodeBuddy",
			ConfigFilePath: "CODEBUDDY.md",
			SlashID:        "codebuddy-slash",
			Priority:       PriorityCodeBuddy,
		}),
		func() providerkit.Provider {
			return &CodeBuddyProvider{}
		},
	)
}

func (*CodeBuddyProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "CODEBUDDY.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*CodeBuddyProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CODEBUDDY.md")

	return providerkit.FileExists(filePath)
}

func (*CodeBuddyProvider) GetName() string {
	return "CodeBuddy"
}
