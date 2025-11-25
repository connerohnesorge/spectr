package providers

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
	FrontmatterProposal = "---\n" +
		"description: Scaffold a new Spectr change and validate strictly.\n" +
		"---"

	// FrontmatterApply is the YAML frontmatter for apply commands.
	FrontmatterApply = "---\n" +
		"description: Implement an approved Spectr change and keep tasks in sync.\n" +
		"---"

	// FrontmatterArchive is the YAML frontmatter for archive commands.
	FrontmatterArchive = "---\n" +
		"description: Archive a deployed Spectr change and update specs.\n" +
		"---"

	// FrontmatterProposalProject includes (project) marker for Claude Code.
	FrontmatterProposalProject = "---\n" +
		"description: Scaffold a new Spectr change and validate strictly. (project)\n" +
		"---"

	// FrontmatterApplyProject includes (project) marker for Claude Code.
	FrontmatterApplyProject = "---\n" +
		"description: Implement an approved Spectr change and keep tasks in sync. (project)\n" +
		"---"

	// FrontmatterArchiveProject includes (project) marker for Claude Code.
	FrontmatterArchiveProject = "---\n" +
		"description: Archive a deployed Spectr change and update specs. (project)\n" +
		"---"
)

// StandardFrontmatter returns the standard frontmatter map for most providers.
func StandardFrontmatter() map[string]string {
	return map[string]string{
		"proposal": FrontmatterProposal,
		"apply":    FrontmatterApply,
		"archive":  FrontmatterArchive,
	}
}

// ProjectFrontmatter returns frontmatter with (project) markers for Claude Code.
func ProjectFrontmatter() map[string]string {
	return map[string]string{
		"proposal": FrontmatterProposalProject,
		"apply":    FrontmatterApplyProject,
		"archive":  FrontmatterArchiveProject,
	}
}
