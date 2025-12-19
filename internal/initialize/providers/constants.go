package providers

// Priority constants for all providers.
// Lower numbers = higher priority (displayed first).
const (
	PriorityClaudeCode  = 1
	PriorityGemini      = 2
	PriorityCostrict    = 3
	PriorityQoder       = 4
	PriorityQwen        = 5
	PriorityAntigravity = 6
	PriorityCline       = 7
	PriorityCursor      = 8
	PriorityCodex       = 9
	PriorityOpencode    = 10
	PriorityAider       = 11
	PriorityWindsurf    = 13
	PriorityKilocode    = 14
	PriorityContinue    = 15
	PriorityCrush       = 16
)

// Default frontmatter templates for slash commands.
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

// StandardFrontmatter returns the standard frontmatter map for most providers.
func StandardFrontmatter() map[string]string {
	return map[string]string{
		"proposal": FrontmatterProposal,
		"apply":    FrontmatterApply,
	}
}
