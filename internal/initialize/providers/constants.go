package providers

import "path/filepath"

// Priority constants for all providers.
// Lower numbers = higher priority (displayed first).
const (
	PriorityClaudeCode  = 1
	PriorityGemini      = 2
	PriorityCostrict    = 3
	PriorityQoder       = 4
	PriorityCodeBuddy   = 5
	PriorityQwen        = 6
	PriorityAntigravity = 7
	PriorityCline       = 8
	PriorityCursor      = 9
	PriorityCodex       = 10
	PriorityAider       = 11
	PriorityTabnine     = 12
	PriorityWindsurf    = 13
	PriorityKilocode    = 14
	PriorityContinue    = 15
)

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

// StandardFrontmatter returns the standard frontmatter templates used by most providers.
// The returned map contains the "proposal" and "apply" keys, each mapped to a YAML frontmatter string for the corresponding command.
func StandardFrontmatter() map[string]string {
	return map[string]string{
		"proposal": FrontmatterProposal,
		"apply":    FrontmatterApply,
	}
}

// StandardCommandPaths returns the standard command paths for a given
// directory and extension.
// Uses subdirectory structure: {dir}/spectr/{command}{ext}
// Example: ".claude/commands", ".md" -> ".claude/commands/spectr/proposal.md"
// StandardCommandPaths constructs standard file paths for the "proposal" and
// "apply" commands located in a "spectr" subdirectory of dir, appending ext
// to each filename and returning the proposal path followed by the apply path.
func StandardCommandPaths(
	dir, ext string,
) (proposalPath, applyPath string) {
	spectrDir := filepath.Join(dir, "spectr")
	proposalPath = filepath.Join(spectrDir, "proposal"+ext)
	applyPath = filepath.Join(spectrDir, "apply"+ext)

	return proposalPath, applyPath
}

// PrefixedCommandPaths returns command paths using a flat prefix pattern.
// Uses flat structure: {dir}/spectr-{command}{ext}
// Example: ".agent/workflows", ".md" -> ".agent/workflows/spectr-proposal.md"
// PrefixedCommandPaths builds filesystem paths for prefixed spectr proposal and apply commands in dir using ext as the file extension.
// The first return is the proposal path (dir/spectr-proposal+ext) and the second is the apply path (dir/spectr-apply+ext).
func PrefixedCommandPaths(
	dir, ext string,
) (proposalPath, applyPath string) {
	proposalPath = filepath.Join(dir, "spectr-proposal"+ext)
	applyPath = filepath.Join(dir, "spectr-apply"+ext)

	return proposalPath, applyPath
}