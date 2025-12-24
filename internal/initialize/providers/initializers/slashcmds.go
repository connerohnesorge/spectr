package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// SlashCommandsInitializer creates slash command files for a provider.
// Supports both Markdown (.md) and TOML (.toml) formats.
type SlashCommandsInitializer struct {
	dir      string
	ext      string
	commands []domain.SlashCommand
	isGlobal bool
}

// fileUpdateCtx holds context for updating a file.
type fileUpdateCtx struct {
	existing    []byte
	content     string // raw template content for marker updates
	fileContent string // fully formatted file content
}

// NewSlashCommandsInitializer creates a SlashCommandsInitializer.
//
// The dir is the base directory (e.g., ".claude/commands/spectr").
// The ext is the file extension (".md" or ".toml").
// The commands is a list of slash commands to create.
func NewSlashCommandsInitializer(
	dir, ext string,
	commands []domain.SlashCommand,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:      dir,
		ext:      ext,
		commands: commands,
		isGlobal: false,
	}
}

// Init creates or updates slash command files.
// Files are created in the format: {dir}/{command}{ext}
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm any,
) (InitResult, error) {
	var result InitResult

	templateProvider, ok := tm.(TemplateProvider)
	if !ok {
		return result, fmt.Errorf(
			"expected TemplateProvider, got %T",
			tm,
		)
	}

	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	if err := fs.MkdirAll(s.dir, DirPerm); err != nil {
		return result, fmt.Errorf(
			"failed to create directory: %w",
			err,
		)
	}

	tp := templateProvider
	ctx := templateCtx

	for _, cmd := range s.commands {
		r, err := s.processCommand(
			fs,
			tp,
			ctx,
			cmd,
		)
		if err != nil {
			return result, err
		}

		result = result.Merge(r)
	}

	return result, nil
}

// processCommand handles a single command file.
func (s *SlashCommandsInitializer) processCommand(
	fs afero.Fs,
	tp TemplateProvider,
	ctx domain.TemplateContext,
	cmd domain.SlashCommand,
) (InitResult, error) {
	var result InitResult

	filePath := filepath.Join(
		s.dir,
		cmd.String()+s.ext,
	)
	content, fileContent, err := s.renderCommand(
		tp,
		ctx,
		cmd,
	)
	if err != nil {
		return result, err
	}

	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return result, err
	}

	if !exists {
		return s.createFile(
			fs,
			filePath,
			fileContent,
		)
	}

	return s.maybeUpdateFile(
		fs,
		filePath,
		content,
		fileContent,
	)
}

// renderCommand renders a command template and formats the content.
func (s *SlashCommandsInitializer) renderCommand(
	tp TemplateProvider,
	ctx domain.TemplateContext,
	cmd domain.SlashCommand,
) (content, fileContent string, err error) {
	templateName, err := cmd.TemplateName()
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to get template name for %s: %w",
			cmd.String(),
			err,
		)
	}

	templateRef := domain.TemplateRef{
		Name:     templateName,
		Template: tp.GetTemplates(),
	}

	content, err = templateRef.Render(ctx)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to render template for %s: %w",
			cmd.String(),
			err,
		)
	}

	fileContent, err = s.formatContent(
		cmd,
		content,
	)
	if err != nil {
		return "", "", err
	}

	return content, fileContent, nil
}

// createFile creates a new command file.
func (*SlashCommandsInitializer) createFile(
	fs afero.Fs,
	filePath, fileContent string,
) (InitResult, error) {
	var result InitResult

	err := afero.WriteFile(
		fs,
		filePath,
		[]byte(fileContent),
		FilePerm,
	)
	if err != nil {
		return result, fmt.Errorf(
			"failed to write file %s: %w",
			filePath,
			err,
		)
	}

	result.CreatedFiles = append(
		result.CreatedFiles,
		filePath,
	)

	return result, nil
}

// maybeUpdateFile updates a file if its content differs.
func (s *SlashCommandsInitializer) maybeUpdateFile(
	fs afero.Fs,
	filePath, content, fileContent string,
) (InitResult, error) {
	var result InitResult

	existing, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return result, fmt.Errorf(
			"failed to read file %s: %w",
			filePath,
			err,
		)
	}

	if string(existing) == fileContent {
		return result, nil
	}

	updateCtx := fileUpdateCtx{
		existing:    existing,
		content:     content,
		fileContent: fileContent,
	}

	updated, err := s.updateFile(
		fs,
		filePath,
		updateCtx,
	)
	if err != nil {
		return result, err
	}

	if updated {
		result.UpdatedFiles = append(
			result.UpdatedFiles,
			filePath,
		)
	}

	return result, nil
}

// formatContent formats the content based on file extension.
func (s *SlashCommandsInitializer) formatContent(
	cmd domain.SlashCommand,
	content string,
) (string, error) {
	switch s.ext {
	case ".md":
		return formatMarkdownCommand(
			cmd,
			content,
		), nil
	case ".toml":
		return formatTOMLCommand(
			cmd,
			content,
		), nil
	default:
		return "", fmt.Errorf(
			"unsupported extension: %s",
			s.ext,
		)
	}
}

// IsSetup returns true if all command files exist.
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

// updateFile updates an existing command file.
// Returns true if the file was updated.
func (s *SlashCommandsInitializer) updateFile(
	fs afero.Fs,
	filePath string,
	ctx fileUpdateCtx,
) (bool, error) {
	if s.ext == ".md" {
		return s.updateMarkdownFile(
			fs,
			filePath,
			ctx,
		)
	}

	// For TOML, just replace the whole file
	err := afero.WriteFile(
		fs,
		filePath,
		[]byte(ctx.fileContent),
		FilePerm,
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to write file %s: %w",
			filePath,
			err,
		)
	}

	return true, nil
}

// updateMarkdownFile updates a markdown file between markers.
func (*SlashCommandsInitializer) updateMarkdownFile(
	fs afero.Fs,
	filePath string,
	ctx fileUpdateCtx,
) (bool, error) {
	contentStr := string(ctx.existing)
	newContentStr, wasUpdated := updateBetweenMarkers(
		contentStr,
		ctx.content,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if !wasUpdated {
		return false, nil
	}

	err := afero.WriteFile(
		fs,
		filePath,
		[]byte(newContentStr),
		FilePerm,
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to write file %s: %w",
			filePath,
			err,
		)
	}

	return true, nil
}

// Path returns the command directory path.
func (s *SlashCommandsInitializer) Path() string {
	return s.dir
}

// IsGlobal returns whether this initializer operates on global files.
func (s *SlashCommandsInitializer) IsGlobal() bool {
	return s.isGlobal
}

// formatMarkdownCommand formats a slash command as a Markdown file.
func formatMarkdownCommand(
	cmd domain.SlashCommand,
	content string,
) string {
	var frontmatter string

	const (
		proposalDesc = "description: Scaffold a new Spectr change.\n"
		applyDesc    = "description: Implement an approved Spectr change.\n"
	)

	switch cmd {
	case domain.SlashProposal:
		frontmatter = "---\n" + proposalDesc + "---"
	case domain.SlashApply:
		frontmatter = "---\n" + applyDesc + "---"
	}

	return frontmatter + "\n\n" +
		SpectrStartMarker + "\n" +
		content + "\n" +
		SpectrEndMarker + "\n"
}

// formatTOMLCommand formats a slash command as a TOML file.
func formatTOMLCommand(
	cmd domain.SlashCommand,
	content string,
) string {
	var description string

	switch cmd {
	case domain.SlashProposal:
		description = "Scaffold a new Spectr change and validate strictly."
	case domain.SlashApply:
		description = "Implement an approved Spectr change."
	}

	var sb strings.Builder

	sb.WriteString("# " + description + "\n\n")
	sb.WriteString("[[agent]]\n")
	sb.WriteString(SpectrStartMarker + "\n")
	sb.WriteString(content + "\n")
	sb.WriteString(SpectrEndMarker + "\n")

	return sb.String()
}
