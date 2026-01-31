package providers

import (
	"context"
)

// KimiProvider implements the Provider interface for Kimi CLI.
// Kimi uses AGENTS.md and .claude/skills/**/ for skills (Kimi loads
// skills from the same location as Claude Code).
type KimiProvider struct{}

// Initializers returns the list of initializers for Kimi CLI.
func (*KimiProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".claude/skills"),
		NewConfigFileInitializer(
			"AGENTS.md",
			tm.Agents(),
		),
		NewSkillFileInitializer(
			".claude/skills/spectr-proposal/SKILL.md",
			tm.ProposalSkill(),
		),
		NewSkillFileInitializer(
			".claude/skills/spectr-apply/SKILL.md",
			tm.ApplySkill(),
		),
		NewSkillFileInitializer(
			".claude/skills/spectr-next/SKILL.md",
			tm.NextSkill(),
		),
		NewAgentSkillsInitializer(
			"spectr-accept-wo-spectr-bin",
			".claude/skills/spectr-accept-wo-spectr-bin",
			tm,
		),
		NewAgentSkillsInitializer(
			"spectr-validate-wo-spectr-bin",
			".claude/skills/spectr-validate-wo-spectr-bin",
			tm,
		),
	}
}
