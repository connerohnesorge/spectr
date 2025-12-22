package providers

import "path/filepath"

// default frontmatter templates for slash commands.
var (
	// FrontmatterProposal is the YAML frontmatter for proposal commands.
	FrontmatterProposal = `---
description: Scaffold a new Spectr change and validate strictly.
---`

	// FrontmatterApply is the YAML frontmatter for apply commands.
	FrontmatterApply = `---
description: Implement an approved Spectr change and keep tasks in sync.
---`
)

// StandardCommandPaths returns the standard command paths for a given
// directory and extension.
// Uses subdirectory structure: {dir}/spectr/{command}{ext}
// Example: ".claude/commands", ".md" -> ".claude/commands/spectr/proposal.md"
// Returns proposalPath, applyPath.
func StandardCommandPaths(
	dir, ext string,
) (proposalPath, applyPath string) {
	spectrDir := filepath.Join(dir, "spectr")
	proposalPath = filepath.Join(
		spectrDir,
		"proposal"+ext,
	)
	applyPath = filepath.Join(
		spectrDir,
		"apply"+ext,
	)

	return proposalPath, applyPath
}

// PrefixedCommandPaths returns command paths using a flat prefix pattern.
// Uses flat structure: {dir}/spectr-{command}{ext}
// Example: ".agent/workflows", ".md" -> ".agent/workflows/spectr-proposal.md"
// Returns proposalPath, applyPath.
func PrefixedCommandPaths(
	dir, ext string,
) (proposalPath, applyPath string) {
	proposalPath = filepath.Join(
		dir,
		"spectr-proposal"+ext,
	)
	applyPath = filepath.Join(
		dir,
		"spectr-apply"+ext,
	)

	return proposalPath, applyPath
}
