package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Register(NewGeminiProvider())
}

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses ~/.gemini/commands/ for TOML-based slash commands
// (no instruction file).
type GeminiProvider struct {
	BaseProvider
}

// NewGeminiProvider creates a new Gemini CLI provider.
func NewGeminiProvider() *GeminiProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".gemini/commands", ".toml",
	)

	return &GeminiProvider{
		BaseProvider: BaseProvider{
			id:            "gemini",
			name:          "Gemini CLI",
			priority:      PriorityGemini,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatTOML,
			frontmatter:   nil,
		},
	}
}

// Configure overrides BaseProvider.Configure to generate TOML files instead of
// markdown.
func (p *GeminiProvider) Configure(
	projectPath, _ string,
	tm TemplateRenderer,
) error {
	if p.HasSlashCommands() {
		err := p.configureSlashCommands(projectPath, tm)
		if err != nil {
			return fmt.Errorf("failed to configure slash commands: %w", err)
		}
	}

	return nil
}

// configureSlashCommands creates or updates TOML slash command files.
func (p *GeminiProvider) configureSlashCommands(
	projectPath string,
	tm TemplateRenderer,
) error {
	commands := []struct {
		name        string
		description string
	}{
		{
			"proposal",
			"Scaffold a new Spectr change and validate strictly.",
		},
		{
			"apply",
			"Implement an approved Spectr change and keep tasks in sync.",
		},
	}

	var err error
	for _, cmd := range commands {
		err = p.configureTOMLCommand(projectPath, cmd.name, cmd.description, tm)
		if err != nil {
			return err
		}
	}

	return nil
}

// configureTOMLCommand creates or updates a single TOML command file.
func (p *GeminiProvider) configureTOMLCommand(
	projectPath, cmd, description string,
	tm TemplateRenderer,
) error {
	filePath := p.getTOMLCommandPath(projectPath, cmd)

	prompt, err := tm.RenderSlashCommand(cmd, DefaultTemplateContext(), p.id)
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	content := p.generateTOMLContent(description, prompt)

	dir := filepath.Dir(filePath)
	err = EnsureDir(dir)
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
			"failed to write TOML command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// getTOMLCommandPath returns the full path for a TOML command file.
func (p *GeminiProvider) getTOMLCommandPath(projectPath, cmd string) string {
	var relPath string
	switch cmd {
	case "proposal":
		relPath = p.proposalPath
	case "apply":
		relPath = p.applyPath
	}

	return filepath.Join(projectPath, relPath)
}

// generateTOMLContent creates TOML content for a Gemini command.
func (*GeminiProvider) generateTOMLContent(
	description, prompt string,
) string {
	// Escape the prompt for TOML multiline string
	escapedPrompt := strings.ReplaceAll(prompt, `\`, `\\`)
	escapedPrompt = strings.ReplaceAll(escapedPrompt, `"`, `\"`)

	return fmt.Sprintf(`# Spectr command for Gemini CLI
description = "%s"
prompt = """
%s
"""
`, description, escapedPrompt)
}
