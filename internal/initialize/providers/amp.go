package providers

import (
	"context"
)

// AmpProvider implements the Provider interface for Amp.
// Amp uses .agents/skills/ for agent skills with SKILL.md frontmatter.
type AmpProvider struct{}

// Initializers returns the list of initializers for Amp.
func (*AmpProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".agents/skills",
		),
		NewConfigFileInitializer(
			"AMP.md",
			tm.InstructionPointer(),
		),
		NewSkillFileInitializer(
			".agents/skills/spectr-proposal/SKILL.md",
			tm.ProposalSkill(),
		),
		NewSkillFileInitializer(
			".agents/skills/spectr-apply/SKILL.md",
			tm.ApplySkill(),
		),
		NewAgentSkillsInitializer(
			"spectr-accept-wo-spectr-bin",
			".agents/skills/spectr-accept-wo-spectr-bin",
			tm,
		),
		NewAgentSkillsInitializer(
			"spectr-validate-wo-spectr-bin",
			".agents/skills/spectr-validate-wo-spectr-bin",
			tm,
		),
	}
}
