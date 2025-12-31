// Package providers implements the composable initializer architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider returns a list of initializers that handle specific aspects
// of configuration (directories, instruction files, slash commands).
//
// # Adding a New Provider
//
// To add a new AI CLI provider:
//
// 1. Create a new file (e.g., providers/mytool.go) with a type implementing Provider: //nolint:lll
//
//	package providers
//
//	import (
//		"context"
//		"github.com/connerohnesorge/spectr/internal/domain"
//	)
//
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) Initializers(ctx context.Context, tm TemplateManager) []Initializer { //nolint:lll
//		return []Initializer{
//			NewDirectoryInitializer(".mytool/commands/spectr"),
//			NewConfigFileInitializer("MYTOOL.md", tm.InstructionPointer()),
//			NewSlashCommandsInitializer(".mytool/commands/spectr", map[domain.SlashCommand]domain.TemplateRef{ //nolint:lll
//				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
//				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
//			}),
//		}
//	}
//
// 2. Register the provider in RegisterAllProviders() in registry.go:
//
//	{ID: "mytool", Name: "MyTool", Priority: 100, Provider: &MyToolProvider{}}
//
// The initializers handle all common logic - no boilerplate needed!
package providers

import (
	"context"
	"io/fs"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// TemplateManager is defined in internal/initialize/templates.go.
// We define it as an interface here to avoid import cycles.
// The actual implementation is passed from the executor.
type TemplateManager interface {
	// InstructionPointer returns the template for instruction pointer files
	// (CLAUDE.md, CLINE.md, etc.)
	InstructionPointer() domain.TemplateRef

	// Agents returns the template for the AGENTS.md file
	Agents() domain.TemplateRef

	// SlashCommand returns the template for a Markdown slash command
	SlashCommand(cmd domain.SlashCommand) domain.TemplateRef

	// ProviderSlashCommand returns a provider-aware Markdown slash command template.
	// Providers can opt in to custom overrides by specifying a provider ID.
	ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef

	// TOMLSlashCommand returns the template for a TOML slash command
	TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef

	// ProviderTOMLSlashCommand returns a provider-aware TOML slash command template.
	// Providers can opt in to custom overrides by specifying a provider ID.
	ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef

	// SkillFS returns an fs.FS rooted at the skill directory for the given
	// skill name. Returns an error if the skill does not exist.
	// The filesystem contains all files under templates/skills/<skillName>/
	// with paths relative to the skill root (e.g., SKILL.md, scripts/accept.sh).
	SkillFS(skillName string) (fs.FS, error)
}

// Provider represents an AI CLI tool (Claude Code, Gemini, Cline, etc.).
// Each provider returns a list of initializers that configure the tool.
//
// Provider metadata (ID, Name, Priority) is specified at registration time
// via the Registration struct, not as methods on this interface.
type Provider interface {
	// Initializers returns the list of initializers for this provider.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout
	//   - tm: TemplateManager for type-safe template access
	//
	// The TemplateManager provides type-safe accessors like:
	//   - tm.InstructionPointer() - for CLAUDE.md, CLINE.md, etc.
	//   - tm.Agents() - for AGENTS.md
	//   - tm.SlashCommand(domain.SlashProposal) - for Markdown
	//   - tm.TOMLSlashCommand(domain.SlashProposal) - for TOML (Gemini)
	//
	// Example:
	//	func (p *ClaudeProvider) Initializers(
	//		ctx context.Context,
	//		tm TemplateManager,
	//	) []Initializer {
	//		return []Initializer{
	//			NewDirectoryInitializer(".claude/commands/spectr"),
	//			NewConfigFileInitializer(
	//				"CLAUDE.md",
	//				tm.InstructionPointer(),
	//			),
	//			NewSlashCommandsInitializer(
	//				".claude/commands/spectr",
	//				map[domain.SlashCommand]domain.TemplateRef{
	//					domain.SlashProposal: tm.SlashCommand(
	//						domain.SlashProposal,
	//					),
	//					domain.SlashApply: tm.SlashCommand(
	//						domain.SlashApply,
	//					),
	//				},
	//			),
	//		}
	//	}
	Initializers(ctx context.Context, tm TemplateManager) []Initializer
}
