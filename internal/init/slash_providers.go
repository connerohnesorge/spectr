//nolint:revive // line-length-limit,file-length-limit - readability over strict formatting
package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// Slash Command Providers
// ============================================================================
// Each slash command provider configures slash commands for a specific tool.
// Slash commands are invoked conditionally (unlike memory files which are
// included in every agent invocation).

// SlashCommandConfig holds configuration for a slash command tool
type SlashCommandConfig struct {
	ToolID      string
	ToolName    string
	Frontmatter map[string]string // proposal, apply, archive frontmatter
	FilePaths   map[string]string // proposal, apply, archive paths
}

// BaseSlashCommandProvider provides common slash command configuration logic
type BaseSlashCommandProvider struct {
	config SlashCommandConfig
}

// ConfigureSlashCommands implements SlashCommandProvider interface
func (s *BaseSlashCommandProvider) ConfigureSlashCommands(projectPath string) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	commands := []string{"proposal", "apply", "archive"}

	for _, cmd := range commands {
		if err := s.configureCommand(tm, projectPath, cmd); err != nil {
			return err
		}
	}

	return nil
}

// configureCommand configures a single slash command
func (s *BaseSlashCommandProvider) configureCommand(
	tm *TemplateManager,
	projectPath, cmd string,
) error {
	relPath, ok := s.config.FilePaths[cmd]
	if !ok {
		return fmt.Errorf("missing file path for command: %s", cmd)
	}

	filePath := filepath.Join(projectPath, relPath)

	body, err := tm.RenderSlashCommand(cmd)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			cmd,
			err,
		)
	}

	if FileExists(filePath) {
		return s.updateExistingCommand(filePath, body)
	}

	return s.createNewCommand(filePath, cmd, body)
}

// updateExistingCommand updates an existing slash command file
func (s *BaseSlashCommandProvider) updateExistingCommand(
	filePath, body string,
) error {
	if err := updateSlashCommandBody(filePath, body); err != nil {
		return fmt.Errorf(
			"failed to update slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// createNewCommand creates a new slash command file
func (s *BaseSlashCommandProvider) createNewCommand(
	filePath, cmd, body string,
) error {
	var sections []string

	if frontmatter, ok := s.config.Frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	sections = append(
		sections,
		SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker,
	)

	content := strings.Join(sections, newlineDouble) + newlineDouble

	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			filePath,
			err,
		)
	}

	if err := os.WriteFile(filePath, []byte(content), defaultFilePerm); err != nil {
		return fmt.Errorf(
			"failed to write slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// AreSlashCommandsConfigured implements SlashCommandProvider interface
func (s *BaseSlashCommandProvider) AreSlashCommandsConfigured(projectPath string) bool {
	// Check if all three slash command files exist
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		relPath, ok := s.config.FilePaths[cmd]
		if !ok {
			return false
		}

		filePath := filepath.Join(projectPath, relPath)
		if !FileExists(filePath) {
			return false
		}
	}

	return true
}

// updateSlashCommandBody updates the body of a slash command file between markers
func updateSlashCommandBody(filePath, body string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	startIndex := strings.Index(contentStr, SpectrStartMarker)
	endIndex := strings.Index(contentStr, SpectrEndMarker)

	if startIndex == -1 || endIndex == -1 || endIndex <= startIndex {
		return fmt.Errorf("missing Spectr markers in %s", filePath)
	}

	before := contentStr[:startIndex+len(SpectrStartMarker)]
	after := contentStr[endIndex:]
	updatedContent := before + "\n" + body + "\n" + after

	if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ============================================================================
// Specific Slash Command Provider Implementations
// ============================================================================

// ClaudeSlashCommandProvider configures Claude slash commands
type ClaudeSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewClaudeSlashCommandProvider creates a new Claude slash command provider
func NewClaudeSlashCommandProvider() *ClaudeSlashCommandProvider {
	return &ClaudeSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "claude",
				ToolName: "Claude Slash Commands",
				Frontmatter: map[string]string{
					"proposal": `---
name: Spectr: Proposal
description: Scaffold a new Spectr change and validate strictly.
category: Spectr
tags: [spectr, change]
---`,
					"apply": `---
name: Spectr: Apply
description: Implement an approved Spectr change and keep tasks in sync.
category: Spectr
tags: [spectr, apply]
---`,
					"archive": `---
name: Spectr: Archive
description: Archive a deployed Spectr change and update specs.
category: Spectr
tags: [spectr, archive]
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".claude/commands/spectr/proposal.md",
					"apply":    ".claude/commands/spectr/apply.md",
					"archive":  ".claude/commands/spectr/archive.md",
				},
			},
		},
	}
}

// ClineSlashCommandProvider configures Cline slash commands
type ClineSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewClineSlashCommandProvider creates a new Cline slash command provider
func NewClineSlashCommandProvider() *ClineSlashCommandProvider {
	return &ClineSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "cline",
				ToolName: "Cline Rules",
				Frontmatter: map[string]string{
					"proposal": "# Spectr: Proposal\n\nScaffold a new Spectr change and validate strictly.",
					"apply":    "# Spectr: Apply\n\nImplement an approved Spectr change and keep tasks in sync.",
					"archive":  "# Spectr: Archive\n\nArchive a deployed Spectr change and update specs.",
				},
				FilePaths: map[string]string{
					"proposal": ".clinerules/spectr-proposal.md",
					"apply":    ".clinerules/spectr-apply.md",
					"archive":  ".clinerules/spectr-archive.md",
				},
			},
		},
	}
}

// CursorSlashCommandProvider configures Cursor slash commands
type CursorSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewCursorSlashCommandProvider creates a new Cursor slash command provider
func NewCursorSlashCommandProvider() *CursorSlashCommandProvider {
	return &CursorSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "cursor",
				ToolName: "Cursor Commands",
				Frontmatter: map[string]string{
					"proposal": `---
name: /spectr-proposal
id: spectr-proposal
category: Spectr
description: Scaffold a new Spectr change and validate strictly.
---`,
					"apply": `---
name: /spectr-apply
id: spectr-apply
category: Spectr
description: Implement an approved Spectr change and keep tasks in sync.
---`,
					"archive": `---
name: /spectr-archive
id: spectr-archive
category: Spectr
description: Archive a deployed Spectr change and update specs.
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".cursor/commands/spectr-proposal.md",
					"apply":    ".cursor/commands/spectr-apply.md",
					"archive":  ".cursor/commands/spectr-archive.md",
				},
			},
		},
	}
}

// ContinueSlashCommandProvider configures Continue slash commands
type ContinueSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewContinueSlashCommandProvider creates a new Continue slash command provider
func NewContinueSlashCommandProvider() *ContinueSlashCommandProvider {
	return &ContinueSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "continue",
				ToolName: "Continue Commands",
				Frontmatter: map[string]string{
					"proposal": `---
name: spectr-proposal
description: Scaffold a new Spectr change and validate strictly.
---`,
					"apply": `---
name: spectr-apply
description: Implement an approved Spectr change and keep tasks in sync.
---`,
					"archive": `---
name: spectr-archive
description: Archive a deployed Spectr change and update specs.
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".continue/commands/spectr-proposal.md",
					"apply":    ".continue/commands/spectr-apply.md",
					"archive":  ".continue/commands/spectr-archive.md",
				},
			},
		},
	}
}

// WindsurfSlashCommandProvider configures Windsurf slash commands
type WindsurfSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewWindsurfSlashCommandProvider creates a new Windsurf slash command provider
func NewWindsurfSlashCommandProvider() *WindsurfSlashCommandProvider {
	return &WindsurfSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "windsurf",
				ToolName: "Windsurf Workflows",
				Frontmatter: map[string]string{
					"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\nauto_execution_mode: 3\n---",
					"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\nauto_execution_mode: 3\n---",
					"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\nauto_execution_mode: 3\n---",
				},
				FilePaths: map[string]string{
					"proposal": ".windsurf/workflows/spectr-proposal.md",
					"apply":    ".windsurf/workflows/spectr-apply.md",
					"archive":  ".windsurf/workflows/spectr-archive.md",
				},
			},
		},
	}
}

// AiderSlashCommandProvider configures Aider slash commands
type AiderSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewAiderSlashCommandProvider creates a new Aider slash command provider
func NewAiderSlashCommandProvider() *AiderSlashCommandProvider {
	return &AiderSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "aider",
				ToolName:    "Aider Commands",
				Frontmatter: make(map[string]string), // No frontmatter for Aider
				FilePaths: map[string]string{
					"proposal": ".aider/commands/spectr-proposal.md",
					"apply":    ".aider/commands/spectr-apply.md",
					"archive":  ".aider/commands/spectr-archive.md",
				},
			},
		},
	}
}

// KilocodeSlashCommandProvider configures Kilocode slash commands
type KilocodeSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewKilocodeSlashCommandProvider creates a new Kilocode slash command provider
func NewKilocodeSlashCommandProvider() *KilocodeSlashCommandProvider {
	return &KilocodeSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "kilocode",
				ToolName:    "Kilocode Workflows",
				Frontmatter: make(map[string]string), // No frontmatter for Kilocode
				FilePaths: map[string]string{
					"proposal": ".kilocode/workflows/spectr-proposal.md",
					"apply":    ".kilocode/workflows/spectr-apply.md",
					"archive":  ".kilocode/workflows/spectr-archive.md",
				},
			},
		},
	}
}

// QoderSlashCommandProvider configures Qoder slash commands
type QoderSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewQoderSlashCommandProvider creates a new Qoder slash command provider
func NewQoderSlashCommandProvider() *QoderSlashCommandProvider {
	return &QoderSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "qoder",
				ToolName: "Qoder Commands",
				Frontmatter: map[string]string{
					"proposal": `---
name: Spectr: Proposal
description: Scaffold a new Spectr change and validate strictly.
category: Spectr
tags: [spectr, change]
---`,
					"apply": `---
name: Spectr: Apply
description: Implement an approved Spectr change and keep tasks in sync.
category: Spectr
tags: [spectr, apply]
---`,
					"archive": `---
name: Spectr: Archive
description: Archive a deployed Spectr change and update specs.
category: Spectr
tags: [spectr, archive]
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".qoder/commands/spectr/proposal.md",
					"apply":    ".qoder/commands/spectr/apply.md",
					"archive":  ".qoder/commands/spectr/archive.md",
				},
			},
		},
	}
}

// CostrictSlashCommandProvider configures CoStrict slash commands
type CostrictSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewCostrictSlashCommandProvider creates a new CoStrict slash command provider
func NewCostrictSlashCommandProvider() *CostrictSlashCommandProvider {
	return &CostrictSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "costrict",
				ToolName: "CoStrict Commands",
				Frontmatter: map[string]string{
					"proposal": `---
description: "Scaffold a new Spectr change and validate strictly."
argument-hint: feature description or request
---`,
					"apply": `---
description: "Implement an approved Spectr change and keep tasks in sync."
argument-hint: change-id
---`,
					"archive": `---
description: "Archive a deployed Spectr change and update specs."
argument-hint: change-id
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".cospec/spectr/commands/spectr-proposal.md",
					"apply":    ".cospec/spectr/commands/spectr-apply.md",
					"archive":  ".cospec/spectr/commands/spectr-archive.md",
				},
			},
		},
	}
}

// CopilotSlashCommandProvider configures GitHub Copilot slash commands
type CopilotSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewCopilotSlashCommandProvider creates a new GitHub Copilot slash command provider
func NewCopilotSlashCommandProvider() *CopilotSlashCommandProvider {
	return &CopilotSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "copilot",
				ToolName:    "GitHub Copilot Instructions",
				Frontmatter: make(map[string]string), // No frontmatter for Copilot
				FilePaths: map[string]string{
					"proposal": ".github/copilot/spectr-proposal.md",
					"apply":    ".github/copilot/spectr-apply.md",
					"archive":  ".github/copilot/spectr-archive.md",
				},
			},
		},
	}
}

// MentatSlashCommandProvider configures Mentat slash commands
type MentatSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewMentatSlashCommandProvider creates a new Mentat slash command provider
func NewMentatSlashCommandProvider() *MentatSlashCommandProvider {
	return &MentatSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "mentat",
				ToolName:    "Mentat Commands",
				Frontmatter: make(map[string]string), // No frontmatter for Mentat
				FilePaths: map[string]string{
					"proposal": ".mentat/commands/spectr-proposal.md",
					"apply":    ".mentat/commands/spectr-apply.md",
					"archive":  ".mentat/commands/spectr-archive.md",
				},
			},
		},
	}
}

// TabnineSlashCommandProvider configures Tabnine slash commands
type TabnineSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewTabnineSlashCommandProvider creates a new Tabnine slash command provider
func NewTabnineSlashCommandProvider() *TabnineSlashCommandProvider {
	return &TabnineSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "tabnine",
				ToolName:    "Tabnine Commands",
				Frontmatter: make(map[string]string), // No frontmatter for Tabnine
				FilePaths: map[string]string{
					"proposal": ".tabnine/commands/spectr-proposal.md",
					"apply":    ".tabnine/commands/spectr-apply.md",
					"archive":  ".tabnine/commands/spectr-archive.md",
				},
			},
		},
	}
}

// SmolSlashCommandProvider configures Smol slash commands
type SmolSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewSmolSlashCommandProvider creates a new Smol slash command provider
func NewSmolSlashCommandProvider() *SmolSlashCommandProvider {
	return &SmolSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "smol",
				ToolName:    "Smol Commands",
				Frontmatter: make(map[string]string), // No frontmatter for Smol
				FilePaths: map[string]string{
					"proposal": ".smol/commands/spectr-proposal.md",
					"apply":    ".smol/commands/spectr-apply.md",
					"archive":  ".smol/commands/spectr-archive.md",
				},
			},
		},
	}
}

// CodeBuddySlashCommandProvider configures CodeBuddy slash commands
type CodeBuddySlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewCodeBuddySlashCommandProvider creates a new CodeBuddy slash command provider
func NewCodeBuddySlashCommandProvider() *CodeBuddySlashCommandProvider {
	return &CodeBuddySlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "codebuddy",
				ToolName: "CodeBuddy Commands",
				Frontmatter: map[string]string{
					"proposal": `---
name: Spectr: Proposal
description: Scaffold a new Spectr change and validate strictly.
category: Spectr
tags: [spectr, change]
---`,
					"apply": `---
name: Spectr: Apply
description: Implement an approved Spectr change and keep tasks in sync.
category: Spectr
tags: [spectr, apply]
---`,
					"archive": `---
name: Spectr: Archive
description: Archive a deployed Spectr change and update specs.
category: Spectr
tags: [spectr, archive]
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".codebuddy/commands/spectr/proposal.md",
					"apply":    ".codebuddy/commands/spectr/apply.md",
					"archive":  ".codebuddy/commands/spectr/archive.md",
				},
			},
		},
	}
}

// QwenSlashCommandProvider configures Qwen slash commands
type QwenSlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewQwenSlashCommandProvider creates a new Qwen slash command provider
func NewQwenSlashCommandProvider() *QwenSlashCommandProvider {
	return &QwenSlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:   "qwen",
				ToolName: "Qwen Commands",
				Frontmatter: map[string]string{
					"proposal": `---
name: /spectr-proposal
id: spectr-proposal
category: Spectr
description: Scaffold a new Spectr change and validate strictly.
---`,
					"apply": `---
name: /spectr-apply
id: spectr-apply
category: Spectr
description: Implement an approved Spectr change and keep tasks in sync.
---`,
					"archive": `---
name: /spectr-archive
id: spectr-archive
category: Spectr
description: Archive a deployed Spectr change and update specs.
---`,
				},
				FilePaths: map[string]string{
					"proposal": ".qwen/commands/spectr-proposal.md",
					"apply":    ".qwen/commands/spectr-apply.md",
					"archive":  ".qwen/commands/spectr-archive.md",
				},
			},
		},
	}
}

// AntigravitySlashCommandProvider configures Antigravity slash commands
type AntigravitySlashCommandProvider struct {
	*BaseSlashCommandProvider
}

// NewAntigravitySlashCommandProvider creates a new Antigravity slash command provider
func NewAntigravitySlashCommandProvider() *AntigravitySlashCommandProvider {
	return &AntigravitySlashCommandProvider{
		BaseSlashCommandProvider: &BaseSlashCommandProvider{
			config: SlashCommandConfig{
				ToolID:      "antigravity",
				ToolName:    "Antigravity Workflows",
				Frontmatter: make(map[string]string), // No frontmatter for Antigravity
				FilePaths: map[string]string{
					"proposal": ".agent/workflows/spectr-proposal.md",
					"apply":    ".agent/workflows/spectr-apply.md",
					"archive":  ".agent/workflows/spectr-archive.md",
				},
			},
		},
	}
}
