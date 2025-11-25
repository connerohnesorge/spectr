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
// Gemini uses ~/.gemini/commands/ for TOML-based slash commands (no instruction file).
type GeminiProvider struct {
	BaseProvider
}

// NewGeminiProvider creates a new Gemini CLI provider.
func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{
		BaseProvider: BaseProvider{
			id:            "gemini",
			name:          "Gemini CLI",
			priority:      PriorityGemini,
			configFile:    "",
			slashDir:      ".gemini/commands",
			commandFormat: FormatTOML,
			frontmatter:   nil,
		},
	}
}

// Configure overrides BaseProvider.Configure to generate TOML files instead of markdown.
func (p *GeminiProvider) Configure(projectPath, spectrDir string, tm TemplateRenderer) error {
	if p.HasSlashCommands() {
		if err := p.configureSlashCommands(projectPath, tm); err != nil {
			return fmt.Errorf("failed to configure slash commands: %w", err)
		}
	}

	return nil
}

// configureSlashCommands creates or updates TOML slash command files.
func (p *GeminiProvider) configureSlashCommands(projectPath string, tm TemplateRenderer) error {
	commands := []struct {
		name        string
		description string
	}{
		{"proposal", "Scaffold a new Spectr change and validate strictly."},
		{"apply", "Implement an approved Spectr change and keep tasks in sync."},
		{"archive", "Archive a deployed Spectr change and update specs."},
	}

	for _, cmd := range commands {
		if err := p.configureTOMLCommand(projectPath, cmd.name, cmd.description, tm); err != nil {
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

	prompt, err := tm.RenderSlashCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	content := p.generateTOMLContent(description, prompt)

	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write TOML command file %s: %w", filePath, err)
	}

	return nil
}

// getTOMLCommandPath returns the full path for a TOML command file.
func (p *GeminiProvider) getTOMLCommandPath(projectPath, cmd string) string {
	filename := fmt.Sprintf("spectr-%s.toml", cmd)

	return filepath.Join(projectPath, p.slashDir, filename)
}

// generateTOMLContent creates TOML content for a Gemini command.
func (p *GeminiProvider) generateTOMLContent(description, prompt string) string {
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

// getSlashCommandPath returns paths with .toml extension for TOML format.
func (p *GeminiProvider) getSlashCommandPath(projectPath, cmd string) string {
	filename := fmt.Sprintf("spectr-%s.toml", cmd)

	return filepath.Join(projectPath, p.slashDir, filename)
}

// IsConfigured checks if all TOML command files exist.
func (p *GeminiProvider) IsConfigured(projectPath string) bool {
	if p.HasSlashCommands() {
		commands := []string{"proposal", "apply", "archive"}
		for _, cmd := range commands {
			filePath := p.getSlashCommandPath(projectPath, cmd)
			if !FileExists(filePath) {
				return false
			}
		}
	}

	return true
}

// GetFilePaths returns the TOML file paths for Gemini.
func (p *GeminiProvider) GetFilePaths() []string {
	var paths []string

	if p.HasSlashCommands() {
		commands := []string{"proposal", "apply", "archive"}
		for _, cmd := range commands {
			filename := fmt.Sprintf("spectr-%s.toml", cmd)
			paths = append(paths, filepath.Join(p.slashDir, filename))
		}
	}

	return paths
}
