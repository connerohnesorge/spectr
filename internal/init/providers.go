//nolint:revive // line-length-limit - readability over strict formatting
package init

import (
	"fmt"
	"path/filepath"
)

// ============================================================================
// Memory File Providers
// ============================================================================
// Each memory file provider manages a specific memory file in the repo root
// that is included in every agent invocation for its respective tool.

// ClaudeMemoryFileProvider manages CLAUDE.md in the repo root
type ClaudeMemoryFileProvider struct{}

func (*ClaudeMemoryFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "CLAUDE.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*ClaudeMemoryFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CLAUDE.md")

	return FileExists(filePath)
}

// ClineMemoryFileProvider manages CLINE.md in the repo root
type ClineMemoryFileProvider struct{}

func (*ClineMemoryFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "CLINE.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*ClineMemoryFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CLINE.md")

	return FileExists(filePath)
}

// QoderMemoryFileProvider manages QODER.md in the repo root
type QoderMemoryFileProvider struct{}

func (*QoderMemoryFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "QODER.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*QoderMemoryFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "QODER.md")

	return FileExists(filePath)
}

// CodeBuddyMemoryFileProvider manages CODEBUDDY.md in the repo root
type CodeBuddyMemoryFileProvider struct{}

func (*CodeBuddyMemoryFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "CODEBUDDY.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*CodeBuddyMemoryFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CODEBUDDY.md")

	return FileExists(filePath)
}

// QwenMemoryFileProvider manages QWEN.md in the repo root
type QwenMemoryFileProvider struct{}

func (*QwenMemoryFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "QWEN.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*QwenMemoryFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "QWEN.md")

	return FileExists(filePath)
}

// CostrictMemoryFileProvider manages COSTRICT.md in the repo root
type CostrictMemoryFileProvider struct{}

func (*CostrictMemoryFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "COSTRICT.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*CostrictMemoryFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "COSTRICT.md")

	return FileExists(filePath)
}

// AgentsFileProvider manages AGENTS.md in the repo root (for Antigravity)
// This provider already exists and is kept focused on managing AGENTS.md
type AgentsFileProvider struct{}

func (*AgentsFileProvider) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	filePath := filepath.Join(projectPath, "AGENTS.md")

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*AgentsFileProvider) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "AGENTS.md")

	return FileExists(filePath)
}

// ============================================================================
// SpectrAgentsUpdater Provider
// ============================================================================
// SpectrAgentsUpdater ensures all memory file-based tools update
// spectr/AGENTS.md with generic Spectr usage instructions.
// This is a cross-cutting concern that applies to all memory file tools.

// SpectrAgentsUpdater updates spectr/AGENTS.md with Spectr instructions
type SpectrAgentsUpdater struct{}

func (*SpectrAgentsUpdater) ConfigureMemoryFile(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	content, err := tm.RenderAgents()
	if err != nil {
		return err
	}

	spectrDir := filepath.Join(projectPath, "spectr")
	filePath := filepath.Join(spectrDir, "AGENTS.md")

	// Ensure spectr directory exists
	if err := EnsureDir(spectrDir); err != nil {
		return fmt.Errorf("failed to create spectr directory: %w", err)
	}

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

func (*SpectrAgentsUpdater) IsMemoryFileConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "spectr", "AGENTS.md")

	return FileExists(filePath)
}
