// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider handles both its instruction file (e.g., CLAUDE.md) and slash
// commands (e.g., .claude/commands/) in a single implementation.
//
// # Adding a New Provider
//
// To add a new AI CLI provider, create a new file
// (e.g., providers/mytools.go) with:
//
// Example:
//
//	package providers
//
//	func init() {
//		Register(&MyToolProvider{})
//	}
//
//	type MyToolProvider struct {
//		BaseProvider
//	}
//
//	func NewMyToolProvider() *MyToolProvider {
//		proposalPath, syncPath, applyPath := StandardCommandPaths(
//			".mytool/commands", ".md",
//		)
//
//		return &MyToolProvider{
//			BaseProvider: BaseProvider{
//				id:            "mytool",
//				name:          "MyTool",
//				priority:      100,
//				// Empty if no instruction file
//				configFile:    "MYTOOL.md",
//				// Empty if no slash commands
//				proposalPath:  proposalPath,
//				syncPath:      syncPath,
//				applyPath:     applyPath,
//				commandFormat: FormatMarkdown,
//				frontmatter: map[string]string{
//					"proposal":
//					"---\ndescription: Scaffold a new Spectr change.\n---",
//					"apply":
//					"---\ndescription: Implement an \n---",
//					"sync":
//					"---\ndescription: Sync a deployed .\n---",
//				},
//			},
//		}
//	}
//
// The BaseProvider handles all common logic.
// Override Configure() only for special formats.
//
//nolint:revive // File length acceptable for provider interface definition
package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CommandFormat specifies the format for slash command files.
type CommandFormat int

const (
	// FormatMarkdown uses markdown files with
	// YAML frontmatter (Claude, Cline, etc.)
	FormatMarkdown CommandFormat = iota
	// FormatTOML uses TOML files (Gemini CLI)
	FormatTOML
)

// Provider represents an AI CLI tool (Claude Code, Gemini, Cline, etc.).
// Each provider handles both its instruction file AND slash commands.
type Provider interface {
	// ID returns the unique provider identifier (kebab-case).
	// Example: "claude-code", "gemini", "cline"
	ID() string

	// Name returns the human-readable provider name for display.
	// Example: "Claude Code", "Gemini CLI", "Cline"
	Name() string

	// Priority returns the display order (lower = higher priority).
	// Claude Code should be 1, other major tools 2-10, etc.
	Priority() int

	// ConfigFile returns the instruction file path (e.g., "CLAUDE.md").
	// Returns empty string if the provider has no instruction file.
	ConfigFile() string

	// GetProposalCommandPath returns the relative path for the
	// proposal command.
	// Example: ".claude/commands/spectr-proposal.md"
	// Returns empty string if the provider has no proposal command.
	GetProposalCommandPath() string

	// GetSyncCommandPath returns the relative path for the sync command.
	// Example: ".claude/commands/spectr-sync.md"
	// Returns empty string if the provider has no sync command.
	GetSyncCommandPath() string

	// GetApplyCommandPath returns the relative path for the apply command.
	// Example: ".claude/commands/spectr-apply.md"
	// Returns empty string if the provider has no apply command.
	GetApplyCommandPath() string

	// CommandFormat returns Markdown or TOML for slash command files.
	CommandFormat() CommandFormat

	// Configure applies all configuration (instruction file + slash commands).
	// projectPath is the root project directory.
	// spectrDir is the path to the spectr/ directory.
	Configure(projectPath, spectrDir string, tm TemplateRenderer) error

	// IsConfigured checks if the provider is fully configured.
	// Returns true if all expected files exist.
	IsConfigured(projectPath string) bool

	// GetFilePaths returns the file paths that this provider creates/updates.
	GetFilePaths() []string

	// HasConfigFile returns true if this provider creates an instruction file.
	HasConfigFile() bool

	// HasSlashCommands returns true if this provider creates slash commands.
	HasSlashCommands() bool
}

// TemplateRenderer provides template rendering capabilities.
//
// This interface allows providers to render templates without depending on the
// full TemplateManager.
type TemplateRenderer interface {
	// RenderAgents renders the AGENTS.md template content.
	RenderAgents() (string, error)
	// RenderSlashCommand renders a slash command template
	// IE. proposal, apply, or sync.
	RenderSlashCommand(command string) (string, error)
}

// BaseProvider provides a default implementation of the Provider
// interface. Embed this in your provider struct for common functionality.
type BaseProvider struct {
	id            string
	name          string
	priority      int
	configFile    string // Empty if no instruction file
	proposalPath  string // e.g., ".claude/commands/spectr-proposal.md"
	syncPath      string // e.g., ".claude/commands/spectr-sync.md"
	applyPath     string // e.g., ".claude/commands/spectr-apply.md"
	commandFormat CommandFormat
	frontmatter   map[string]string // Command name -> frontmatter content
}

// ID returns the provider identifier.
func (p *BaseProvider) ID() string {
	return p.id
}

// Name returns the human-readable name.
func (p *BaseProvider) Name() string {
	return p.name
}

// Priority returns the display order.
func (p *BaseProvider) Priority() int {
	return p.priority
}

// ConfigFile returns the instruction file path.
func (p *BaseProvider) ConfigFile() string {
	return p.configFile
}

// GetProposalCommandPath returns the relative path for the proposal command.
func (p *BaseProvider) GetProposalCommandPath() string {
	return p.proposalPath
}

// GetSyncCommandPath returns the relative path for the sync command.
func (p *BaseProvider) GetSyncCommandPath() string {
	return p.syncPath
}

// GetApplyCommandPath returns the relative path for the apply command.
func (p *BaseProvider) GetApplyCommandPath() string {
	return p.applyPath
}

// CommandFormat returns the command file format.
func (p *BaseProvider) CommandFormat() CommandFormat {
	return p.commandFormat
}

// HasConfigFile returns true if this provider has an instruction file.
func (p *BaseProvider) HasConfigFile() bool {
	return p.configFile != ""
}

// HasSlashCommands returns true if this provider has slash commands.
func (p *BaseProvider) HasSlashCommands() bool {
	return p.proposalPath != "" || p.syncPath != "" || p.applyPath != ""
}

// Configure applies all configuration for the provider.
func (p *BaseProvider) Configure(
	projectPath, _ string,
	tm TemplateRenderer,
) error {
	// Configure instruction file if provider has one
	if p.HasConfigFile() {
		if err := p.configureConfigFile(projectPath, tm); err != nil {
			return fmt.Errorf("failed to configure instruction file: %w", err)
		}
	}

	// Configure slash commands if provider has them
	if p.HasSlashCommands() {
		if err := p.configureSlashCommands(projectPath, tm); err != nil {
			return fmt.Errorf("failed to configure slash commands: %w", err)
		}
	}

	return nil
}

// configureConfigFile creates or updates the instruction file.
func (p *BaseProvider) configureConfigFile(
	projectPath string,
	tm TemplateRenderer,
) error {
	content, err := tm.RenderAgents()
	if err != nil {
		return fmt.Errorf("failed to render agents template: %w", err)
	}

	filePath := filepath.Join(projectPath, p.configFile)

	return UpdateFileWithMarkers(
		filePath,
		content,
		SpectrStartMarker,
		SpectrEndMarker,
	)
}

// configureSlashCommands creates or updates the slash command files.
func (p *BaseProvider) configureSlashCommands(
	projectPath string,
	tm TemplateRenderer,
) error {
	paths := map[string]string{
		"proposal": p.proposalPath,
		"apply":    p.applyPath,
		"sync":     p.syncPath,
	}

	for cmd, relPath := range paths {
		if relPath == "" {
			continue
		}
		fullPath := filepath.Join(projectPath, relPath)
		if err := p.configureSlashCommand(fullPath, cmd, tm); err != nil {
			return err
		}
	}

	return nil
}

// configureSlashCommand configures a single slash command file.
func (p *BaseProvider) configureSlashCommand(
	filePath, cmd string,
	tm TemplateRenderer,
) error {
	body, err := tm.RenderSlashCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	if FileExists(filePath) {
		return p.updateExistingSlashCommand(filePath, body, cmd)
	}

	return p.createNewSlashCommand(filePath, cmd, body)
}

// updateExistingSlashCommand updates an existing slash command file.
func (p *BaseProvider) updateExistingSlashCommand(
	filePath, body, cmd string,
) error {
	frontmatter := p.frontmatter[cmd]
	err := updateSlashCommandBody(filePath, body, frontmatter)
	if err != nil {
		return fmt.Errorf(
			"failed to update slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// createNewSlashCommand creates a new slash command file.
func (p *BaseProvider) createNewSlashCommand(filePath, cmd, body string) error {
	var sections []string

	if frontmatter, ok := p.frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	sections = append(
		sections,
		SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker,
	)
	content := strings.Join(sections, newlineDouble) + newlineDouble

	dir := filepath.Dir(filePath)
	err := EnsureDir(dir)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			filePath,
			err,
		)
	}

	err = os.WriteFile(filePath, []byte(content), filePerm)
	if err != nil {
		return fmt.Errorf(
			"failed to write slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// IsConfigured checks if the provider is fully configured.
func (p *BaseProvider) IsConfigured(projectPath string) bool {
	// Check instruction file if provider has one
	if p.HasConfigFile() {
		filePath := filepath.Join(projectPath, p.configFile)
		if !FileExists(filePath) {
			return false
		}
	}

	// Check slash commands using path methods
	paths := []string{p.proposalPath, p.syncPath, p.applyPath}
	for _, relPath := range paths {
		if relPath == "" {
			continue
		}
		filePath := filepath.Join(projectPath, relPath)
		if !FileExists(filePath) {
			return false
		}
	}

	return true
}

// GetFilePaths returns the file paths that this provider creates/updates.
func (p *BaseProvider) GetFilePaths() []string {
	var paths []string

	if p.HasConfigFile() {
		paths = append(paths, p.configFile)
	}

	// Add command paths that are non-empty
	for _, relPath := range []string{
		p.proposalPath, p.syncPath, p.applyPath,
	} {
		if relPath != "" {
			paths = append(paths, relPath)
		}
	}

	return paths
}
