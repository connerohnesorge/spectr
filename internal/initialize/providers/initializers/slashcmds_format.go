// Package initializers provides components that initialize spectr files.
package initializers

import (
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// formatMarkdownCommand formats a slash command as a Markdown file.
// It wraps the content with YAML frontmatter and spectr markers.
func formatMarkdownCommand(cmd domain.SlashCommand, content string) string {
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
// It creates a TOML agent block with spectr markers around the content.
func formatTOMLCommand(cmd domain.SlashCommand, content string) string {
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
