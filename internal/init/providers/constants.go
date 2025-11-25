package providers

import "path/filepath"

// Priority constants for all providers.
// Lower numbers = higher priority (displayed first).
const (
	PriorityClaudeCode  = 1
	PriorityCline       = 2
	PriorityCostrict    = 3
	PriorityQoder       = 4
	PriorityCodeBuddy   = 5
	PriorityQwen        = 6
	PriorityAntigravity = 7
	PriorityGemini      = 8
	PriorityCursor      = 10
	PriorityCopilot     = 11
	PriorityAider       = 12
	PriorityContinue    = 13
	PriorityMentat      = 14
	PriorityTabnine     = 15
	PrioritySmol        = 16
	PriorityWindsurf    = 17
	PriorityKilocode    = 18
)

// Frontmatter templates for slash commands.
var (
	// FrontmatterProposal is the YAML frontmatter for proposal commands.
	FrontmatterProposal = `---
description: Scaffold a new Spectr change and validate strictly.
---`

	// FrontmatterApply is the YAML frontmatter for apply commands.
	FrontmatterApply = `---
description: Implement an approved Spectr change and keep tasks in sync.
---`

	// FrontmatterArchive is the YAML frontmatter for archive commands.
	FrontmatterArchive = `---
description: Archive a deployed Spectr change and update specs.
---`
)

// StandardFrontmatter returns the standard frontmatter map for most providers.
func StandardFrontmatter() map[string]string {
	return map[string]string{
		"proposal": FrontmatterProposal,
		"apply":    FrontmatterApply,
		"archive":  FrontmatterArchive,
	}
}

// StandardCommandPaths returns the standard command paths for a given
// directory and extension.
// Returns proposalPath, archivePath, applyPath.
func StandardCommandPaths(
	dir, ext string,
) (proposalPath, archivePath, applyPath string) {
	proposalPath = filepath.Join(dir, "spectr-proposal"+ext)
	archivePath = filepath.Join(dir, "spectr-archive"+ext)
	applyPath = filepath.Join(dir, "spectr-apply"+ext)

	return proposalPath, archivePath, applyPath
}
