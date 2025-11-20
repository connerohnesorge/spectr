// Package providers contains individual AI tool provider implementations.
// This file implements the Claude Code provider.
package providers

import (
	"path/filepath"

	"github.com/conneroisu/spectr/internal/providerkit"
)

// ClaudeCodeProvider configures Claude Code support by creating the
// CLAUDE.md instruction file. This provider automatically registers
// itself with the global provider registry during initialization.
type ClaudeCodeProvider struct{}

func init() {
	MustRegister(
		NewConfigMetadata(ConfigParams{
			ID:             "claude-code",
			Name:           "Claude Code",
			ConfigFilePath: "CLAUDE.md",
			SlashID:        "claude",
			Priority:       PriorityClaudeCode,
		}),
		func() providerkit.Provider {
			return &ClaudeCodeProvider{}
		},
	)
}

func (*ClaudeCodeProvider) Configure(projectPath, _spectrDir string) error {
	tm, err := providerkit.NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "CLAUDE.md")

	return providerkit.UpdateFileWithMarkers(
		filePath,
		content,
		providerkit.SpectrStartMarker,
		providerkit.SpectrEndMarker,
	)
}

func (*ClaudeCodeProvider) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CLAUDE.md")

	return providerkit.FileExists(filePath)
}

func (*ClaudeCodeProvider) GetName() string {
	return "Claude Code"
}
