//nolint:revive // line-length-limit,file-length-limit,add-constant,unused-receiver - readability over strict formatting
package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Configurator interface for all tool configurators
type Configurator interface {
	// Configure configures a tool for the given project path
	Configure(projectPath, spectrDir string) error
	// IsConfigured checks if a tool is already configured for the given project path
	IsConfigured(projectPath string) bool
	// GetName returns the name of the tool
	GetName() string
}

// spectr_MARKERS for managing config file updates
const (
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
)

// UpdateFileWithMarkers updates a file with content between markers
// If file doesn't exist, creates it with markers
// If file exists, updates content between markers
func UpdateFileWithMarkers(filePath, content, startMarker, endMarker string) error {
	expandedPath, err := ExpandPath(filePath)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	var existingContent string

	if FileExists(expandedPath) {
		// Read existing content
		data, err := os.ReadFile(expandedPath)
		if err != nil {
			return fmt.Errorf("failed to read existing file: %w", err)
		}
		existingContent = string(data)

		// Find markers
		startIndex := findMarkerIndex(existingContent, startMarker, 0)
		var endIndex int
		if startIndex != -1 {
			endIndex = findMarkerIndex(existingContent, endMarker, startIndex+len(startMarker))
		} else {
			endIndex = findMarkerIndex(existingContent, endMarker, 0)
		}

		// Handle different marker states
		switch {
		case startIndex != -1 && endIndex != -1:
			// Both markers found - update content between them
			if endIndex < startIndex {
				return fmt.Errorf(
					"invalid marker state in %s: end marker appears before start marker",
					expandedPath,
				)
			}

			before := existingContent[:startIndex]
			after := existingContent[endIndex+len(endMarker):]
			existingContent = before + startMarker + "\n" + content + "\n" + endMarker + after
		case startIndex == -1 && endIndex == -1:
			// No markers found - prepend with markers
			existingContent = startMarker + "\n" + content + "\n" + endMarker + "\n\n" + existingContent
		default:
			// Only one marker found - error
			return fmt.Errorf("invalid marker state in %s: found start: %t, found end: %t",
				expandedPath, startIndex != -1, endIndex != -1)
		}
	} else {
		// File doesn't exist - create with markers
		existingContent = startMarker + "\n" + content + "\n" + endMarker
	}

	// Ensure parent directory exists
	dir := filepath.Dir(expandedPath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(expandedPath, []byte(existingContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// findMarkerIndex finds the index of a marker that is on its own line
// Returns -1 if not found
func findMarkerIndex(content, marker string, fromIndex int) int {
	currentIndex := strings.Index(content[fromIndex:], marker)
	if currentIndex == -1 {
		return -1
	}
	currentIndex += fromIndex

	for currentIndex != -1 {
		if isMarkerOnOwnLine(content, currentIndex, len(marker)) {
			return currentIndex
		}

		nextIndex := strings.Index(content[currentIndex+len(marker):], marker)
		if nextIndex == -1 {
			return -1
		}
		currentIndex = currentIndex + len(marker) + nextIndex
	}

	return -1
}

// isMarkerOnOwnLine checks if a marker is on its own line (only whitespace around it)
func isMarkerOnOwnLine(content string, markerIndex, markerLength int) bool {
	// Check left side
	leftIndex := markerIndex - 1
	for leftIndex >= 0 && content[leftIndex] != '\n' {
		char := content[leftIndex]
		if char != ' ' && char != '\t' && char != '\r' {
			return false
		}
		leftIndex--
	}

	// Check right side
	rightIndex := markerIndex + markerLength
	for rightIndex < len(content) && content[rightIndex] != '\n' {
		char := content[rightIndex]
		if char != ' ' && char != '\t' && char != '\r' {
			return false
		}
		rightIndex++
	}

	return true
}

// ============================================================================
// Config-based Configurators (create instruction files like CLAUDE.md)
// ============================================================================

// ClaudeCodeConfigurator configures Claude Code
type ClaudeCodeConfigurator struct{}

func (*ClaudeCodeConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*ClaudeCodeConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CLAUDE.md")

	return FileExists(filePath)
}

func (*ClaudeCodeConfigurator) GetName() string {
	return "Claude Code"
}

// ClineConfigurator configures Cline
type ClineConfigurator struct{}

func (*ClineConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*ClineConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CLINE.md")

	return FileExists(filePath)
}

func (*ClineConfigurator) GetName() string {
	return "Cline"
}

// CostrictConfigurator configures CoStrict
type CostrictConfigurator struct{}

func (*CostrictConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*CostrictConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "COSTRICT.md")

	return FileExists(filePath)
}

func (*CostrictConfigurator) GetName() string {
	return "CoStrict"
}

// QoderConfigurator configures Qoder
type QoderConfigurator struct{}

func (*QoderConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*QoderConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "QODER.md")

	return FileExists(filePath)
}

func (*QoderConfigurator) GetName() string {
	return "Qoder"
}

// CodeBuddyConfigurator configures CodeBuddy
type CodeBuddyConfigurator struct{}

func (*CodeBuddyConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*CodeBuddyConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "CODEBUDDY.md")

	return FileExists(filePath)
}

func (*CodeBuddyConfigurator) GetName() string {
	return "CodeBuddy"
}

// QwenConfigurator configures Qwen Code
type QwenConfigurator struct{}

func (*QwenConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*QwenConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "QWEN.md")

	return FileExists(filePath)
}

func (*QwenConfigurator) GetName() string {
	return "Qwen Code"
}

// AntigravityConfigurator configures Antigravity
type AntigravityConfigurator struct{}

func (*AntigravityConfigurator) Configure(projectPath, _spectrDir string) error {
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

func (*AntigravityConfigurator) IsConfigured(projectPath string) bool {
	filePath := filepath.Join(projectPath, "AGENTS.md")

	return FileExists(filePath)
}

func (*AntigravityConfigurator) GetName() string {
	return "Antigravity"
}

// ============================================================================
// Slash Command Configurators
// ============================================================================

// SlashCommandConfig holds configuration for a slash command tool
type SlashCommandConfig struct {
	ToolID      string
	ToolName    string
	Frontmatter map[string]string // proposal, apply, archive frontmatter
	FilePaths   map[string]string // proposal, apply, archive paths
}

// SlashCommandConfigurator configures slash commands for a tool
type SlashCommandConfigurator struct {
	config SlashCommandConfig
}

// NewSlashCommandConfigurator creates a new slash command configurator
func NewSlashCommandConfigurator(config SlashCommandConfig) *SlashCommandConfigurator {
	return &SlashCommandConfigurator{config: config}
}

func (s *SlashCommandConfigurator) Configure(
	projectPath,
	_spectrDir string,
) error {
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
func (s *SlashCommandConfigurator) configureCommand(
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
func (s *SlashCommandConfigurator) updateExistingCommand(
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
func (s *SlashCommandConfigurator) createNewCommand(
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

func (s *SlashCommandConfigurator) IsConfigured(projectPath string) bool {
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

func (s *SlashCommandConfigurator) GetName() string {
	return s.config.ToolName
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

// NewClaudeSlashConfigurator creates a Claude slash command configurator
func NewClaudeSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewKilocodeSlashConfigurator creates a Kilocode slash command configurator
func NewKilocodeSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "kilocode",
		ToolName:    "Kilocode Workflows",
		Frontmatter: make(map[string]string), // No frontmatter for Kilocode
		FilePaths: map[string]string{
			"proposal": ".kilocode/workflows/spectr-proposal.md",
			"apply":    ".kilocode/workflows/spectr-apply.md",
			"archive":  ".kilocode/workflows/spectr-archive.md",
		},
	})
}

// NewQoderSlashConfigurator creates a Qoder slash command configurator
func NewQoderSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewCursorSlashConfigurator creates a Cursor slash command configurator
func NewCursorSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewAiderSlashConfigurator creates an Aider slash command configurator
func NewAiderSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "aider",
		ToolName:    "Aider Commands",
		Frontmatter: make(map[string]string), // No frontmatter for Aider
		FilePaths: map[string]string{
			"proposal": ".aider/commands/spectr-proposal.md",
			"apply":    ".aider/commands/spectr-apply.md",
			"archive":  ".aider/commands/spectr-archive.md",
		},
	})
}

// NewContinueSlashConfigurator creates a Continue slash command configurator
func NewContinueSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewCopilotSlashConfigurator creates a GitHub Copilot slash command configurator
func NewCopilotSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "copilot",
		ToolName:    "GitHub Copilot Instructions",
		Frontmatter: make(map[string]string), // No frontmatter for Copilot
		FilePaths: map[string]string{
			"proposal": ".github/copilot/spectr-proposal.md",
			"apply":    ".github/copilot/spectr-apply.md",
			"archive":  ".github/copilot/spectr-archive.md",
		},
	})
}

// NewMentatSlashConfigurator creates a Mentat slash command configurator
func NewMentatSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "mentat",
		ToolName:    "Mentat Commands",
		Frontmatter: make(map[string]string), // No frontmatter for Mentat
		FilePaths: map[string]string{
			"proposal": ".mentat/commands/spectr-proposal.md",
			"apply":    ".mentat/commands/spectr-apply.md",
			"archive":  ".mentat/commands/spectr-archive.md",
		},
	})
}

// NewTabnineSlashConfigurator creates a Tabnine slash command configurator
func NewTabnineSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "tabnine",
		ToolName:    "Tabnine Commands",
		Frontmatter: make(map[string]string), // No frontmatter for Tabnine
		FilePaths: map[string]string{
			"proposal": ".tabnine/commands/spectr-proposal.md",
			"apply":    ".tabnine/commands/spectr-apply.md",
			"archive":  ".tabnine/commands/spectr-archive.md",
		},
	})
}

// NewSmolSlashConfigurator creates a Smol slash command configurator
func NewSmolSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "smol",
		ToolName:    "Smol Commands",
		Frontmatter: make(map[string]string), // No frontmatter for Smol
		FilePaths: map[string]string{
			"proposal": ".smol/commands/spectr-proposal.md",
			"apply":    ".smol/commands/spectr-apply.md",
			"archive":  ".smol/commands/spectr-archive.md",
		},
	})
}

// NewCostrictSlashConfigurator creates a CoStrict slash command configurator
func NewCostrictSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewClineSlashConfigurator creates a Cline slash command configurator
func NewClineSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewWindsurfSlashConfigurator creates a Windsurf slash command configurator
func NewWindsurfSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewCodeBuddySlashConfigurator creates a CodeBuddy slash command configurator
func NewCodeBuddySlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewQwenSlashConfigurator creates a Qwen slash command configurator
func NewQwenSlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
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
	})
}

// NewAntigravitySlashConfigurator creates an Antigravity slash command configurator
func NewAntigravitySlashConfigurator() *SlashCommandConfigurator {
	return NewSlashCommandConfigurator(SlashCommandConfig{
		ToolID:      "antigravity",
		ToolName:    "Antigravity Workflows",
		Frontmatter: make(map[string]string), // No frontmatter for Antigravity
		FilePaths: map[string]string{
			"proposal": ".agent/workflows/spectr-proposal.md",
			"apply":    ".agent/workflows/spectr-apply.md",
			"archive":  ".agent/workflows/spectr-archive.md",
		},
	})
}
