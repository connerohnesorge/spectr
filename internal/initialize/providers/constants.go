package providers

import "strings"

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

// File and directory permission constants.
const (
	dirPerm  = 0o755
	filePerm = 0o644
)

// Marker constants for managing config file updates.
const (
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
)

// Common strings.
const (
	newline       = "\n"
	newlineDouble = "\n\n"
)

// findMarkerIndex finds the index of a marker in content, starting from offset.
func findMarkerIndex(content, marker string, offset int) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}

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
