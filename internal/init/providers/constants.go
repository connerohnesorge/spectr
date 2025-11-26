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

	// FrontmatterSync is the YAML frontmatter for sync commands.
	FrontmatterSync = `---
description: Detect spec drift from code and update specs interactively.
---`
)

// StandardFrontmatter returns the standard frontmatter map for most providers.
func StandardFrontmatter() map[string]string {
	return map[string]string{
		"proposal": FrontmatterProposal,
		"apply":    FrontmatterApply,
		"sync":     FrontmatterSync,
	}
}

// StandardCommandPaths returns the standard command paths for a given
// directory and extension.
// Returns proposalPath, syncPath, applyPath.
func StandardCommandPaths(
	dir, ext string,
) (proposalPath, syncPath, applyPath string) {
	spectrDir := filepath.Join(dir, "spectr")
	proposalPath = filepath.Join(spectrDir, "proposal"+ext)
	syncPath = filepath.Join(spectrDir, "sync"+ext)
	applyPath = filepath.Join(spectrDir, "apply"+ext)

	return proposalPath, syncPath, applyPath
}
