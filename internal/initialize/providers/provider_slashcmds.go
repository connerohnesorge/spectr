package providers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/spf13/afero"
)

// Default frontmatter for slash commands (Markdown format)
const (
	defaultProposalFrontmatter = `---
description: Scaffold a new Spectr change and validate strictly.
---`

	defaultApplyFrontmatter = `---
description: Implement an approved Spectr change and keep tasks in sync.
---`
)

// SlashCommandsInitializer creates or updates slash command files.
// It supports both Markdown (with YAML frontmatter) and TOML formats.
//
// Behavior:
//   - Creates parent directory if it doesn't exist
//   - For each command: creates file with frontmatter (if
//     applicable) + marker-wrapped body
//   - For existing files: updates content between markers, preserves
//     user frontmatter if present
type SlashCommandsInitializer struct {
	// Parent directory (e.g., ".claude/commands/spectr")
	dir string
	// File extension (e.g., ".md" or ".toml")
	ext string
	// Commands to create
	commands []templates.SlashCommand
	// Whether to use globalFs
	isGlobal bool
	// Optional custom frontmatter per command
	frontmatter map[templates.SlashCommand]string
}

// NewSlashCommandsInitializer creates a SlashCommandsInitializer
// for project-relative slash commands.
//
// Parameters:
//   - dir: Directory where slash commands will be created
//     (e.g., ".claude/commands/spectr")
//   - ext: File extension including dot (e.g., ".md" or ".toml")
//   - commands: List of slash commands to create
//     (e.g., []SlashCommand{SlashProposal, SlashApply})
//
// Example:
//
//	NewSlashCommandsInitializer(
//	    ".claude/commands/spectr",
//	    ".md",
//	    []templates.SlashCommand{
//	        templates.SlashProposal,
//	        templates.SlashApply,
//	    },
//	)
func NewSlashCommandsInitializer(
	dir, ext string,
	commands []templates.SlashCommand,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:      dir,
		ext:      ext,
		commands: commands,
		isGlobal: false,
		frontmatter: make(
			map[templates.SlashCommand]string,
		),
	}
}

// NewGlobalSlashCommandsInitializer creates a
// SlashCommandsInitializer for global slash commands.
//
// Parameters:
//   - dir: Directory relative to home directory
//     (e.g., ".config/aider/commands")
//   - ext: File extension including dot (e.g., ".md" or ".toml")
//   - commands: List of slash commands to create
func NewGlobalSlashCommandsInitializer(
	dir, ext string,
	commands []templates.SlashCommand,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:      dir,
		ext:      ext,
		commands: commands,
		isGlobal: true,
		frontmatter: make(
			map[templates.SlashCommand]string,
		),
	}
}

// WithFrontmatter sets custom frontmatter for a specific command.
// This allows providers to customize the frontmatter beyond the defaults.
//
// Example:
//
//	init := NewSlashCommandsInitializer(
//	    ".claude/commands/spectr",
//	    ".md",
//	    commands,
//	)
//	init.WithFrontmatter(
//	    templates.SlashProposal,
//	    "---\ndescription: Custom description\n---",
//	)
//
// Returns the initializer for method chaining.
func (s *SlashCommandsInitializer) WithFrontmatter(
	cmd templates.SlashCommand,
	frontmatter string,
) *SlashCommandsInitializer {
	s.frontmatter[cmd] = frontmatter

	return s
}

// Init creates or updates all slash command files.
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	// Ensure directory exists
	if err := fs.MkdirAll(s.dir, dirPerm); err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to create directory: %w",
			err,
		)
	}

	var result InitResult

	// Create template context
	templateCtx := TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Process each command
	for _, cmd := range s.commands {
		cmdResult, err := s.initCommand(
			fs,
			cmd,
			tm,
			templateCtx,
		)
		if err != nil {
			return result, fmt.Errorf(
				"failed to initialize command %s: %w",
				cmd.String(),
				err,
			)
		}
		result = result.Merge(cmdResult)
	}

	return result, nil
}

// initCommand initializes a single slash command file.
func (s *SlashCommandsInitializer) initCommand(
	fs afero.Fs,
	cmd templates.SlashCommand,
	tm TemplateManager,
	ctx TemplateContext,
) (InitResult, error) {
	// Get template reference and render body
	templateRefRaw := tm.SlashCommand(cmd)

	// Cast to TemplateRef (method returns templates.TemplateRef
	// stored as interface{})
	templateRef, ok := templateRefRaw.(templates.TemplateRef)
	if !ok {
		return InitResult{}, errors.New(
			"SlashCommand() did not return a TemplateRef",
		)
	}

	body, err := templateRef.Render(ctx)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to render template: %w",
			err,
		)
	}

	// Determine file path
	filePath := filepath.Join(
		s.dir,
		cmd.String()+s.ext,
	)

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to check file existence: %w",
			err,
		)
	}

	if !exists {
		return s.createCommand(
			fs,
			filePath,
			cmd,
			body,
		)
	}

	return s.updateCommand(
		fs,
		filePath,
		cmd,
		body,
	)
}

// createCommand creates a new slash command file.
func (s *SlashCommandsInitializer) createCommand(
	fs afero.Fs,
	filePath string,
	cmd templates.SlashCommand,
	body string,
) (InitResult, error) {
	var sections []string

	// Add frontmatter for Markdown files
	if s.ext == ".md" {
		frontmatter := s.getFrontmatter(cmd)
		if frontmatter != "" {
			sections = append(
				sections,
				strings.TrimSpace(frontmatter),
			)
		}
	}

	// Add marker-wrapped body
	markerSection := spectrStartMarker + doubleNewline + body +
		doubleNewline + spectrEndMarker
	sections = append(sections, markerSection)

	// Join sections with double newline
	content := strings.Join(
		sections,
		doubleNewline,
	) + doubleNewline

	// Write file
	err := afero.WriteFile(
		fs,
		filePath,
		[]byte(content),
		configFilePerm,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	return InitResult{
		CreatedFiles: []string{filePath},
	}, nil
}

// updateCommand updates an existing slash command file.
func (s *SlashCommandsInitializer) updateCommand(
	fs afero.Fs,
	filePath string,
	cmd templates.SlashCommand,
	body string,
) (InitResult, error) {
	// Read existing file
	existingBytes, err := afero.ReadFile(
		fs,
		filePath,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to read file: %w",
			err,
		)
	}

	existingContent := string(existingBytes)

	// Find markers
	startIdx := strings.Index(
		existingContent,
		spectrStartMarker,
	)
	if startIdx == -1 {
		return InitResult{}, fmt.Errorf(
			"start marker not found in %s",
			filePath,
		)
	}

	searchOffset := startIdx + len(
		spectrStartMarker,
	)
	relativeEndIdx := strings.Index(
		existingContent[searchOffset:],
		spectrEndMarker,
	)
	if relativeEndIdx == -1 {
		return InitResult{}, fmt.Errorf(
			"end marker not found in %s",
			filePath,
		)
	}
	endIdx := searchOffset + relativeEndIdx

	// Check if file already has frontmatter
	before := existingContent[:startIdx]
	after := existingContent[endIdx+len(spectrEndMarker):]

	// Add frontmatter if it's a Markdown file and doesn't have
	// frontmatter yet
	if s.ext == ".md" &&
		!strings.HasPrefix(
			strings.TrimSpace(before),
			"---",
		) {
		frontmatter := s.getFrontmatter(cmd)
		if frontmatter != "" {
			before = strings.TrimSpace(
				frontmatter,
			) +
				doubleNewline + strings.TrimLeft(
				before,
				"\n\r",
			)
		}
	}

	// Reconstruct file with updated body
	newContent := before + spectrStartMarker + doubleNewline +
		body + doubleNewline + spectrEndMarker + after

	// Write updated file
	err = afero.WriteFile(
		fs,
		filePath,
		[]byte(newContent),
		configFilePerm,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	return InitResult{
		UpdatedFiles: []string{filePath},
	}, nil
}

// getFrontmatter returns the frontmatter for a command.
// Uses custom frontmatter if set, otherwise returns default for
// Markdown files.
func (s *SlashCommandsInitializer) getFrontmatter(
	cmd templates.SlashCommand,
) string {
	// Check for custom frontmatter
	if fm, ok := s.frontmatter[cmd]; ok {
		return fm
	}

	// Return default frontmatter for Markdown files
	if s.ext == ".md" {
		switch cmd {
		case templates.SlashProposal:
			return defaultProposalFrontmatter
		case templates.SlashApply:
			return defaultApplyFrontmatter
		}
	}

	return ""
}

// IsSetup returns true if all slash command files exist.
func (s *SlashCommandsInitializer) IsSetup(
	fs afero.Fs,
	_ *Config,
) bool {
	for _, cmd := range s.commands {
		filePath := filepath.Join(
			s.dir,
			cmd.String()+s.ext,
		)
		exists, err := afero.Exists(fs, filePath)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// Path returns a unique identifier for this initializer.
// Uses pattern "dir/*ext" to distinguish from DirectoryInitializer.
// Example: ".claude/commands/spectr/*.md"
func (s *SlashCommandsInitializer) Path() string {
	return s.dir + "/*" + s.ext
}

// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
func (s *SlashCommandsInitializer) IsGlobal() bool {
	return s.isGlobal
}
