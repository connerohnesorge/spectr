package providers

import (
	"context"
	"os/exec"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/ralph"
)

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/spectr/ for slash commands.
type ClaudeProvider struct{}

// Initializers returns the list of initializers for Claude Code.
func (*ClaudeProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	// Claude Code needs special frontmatter for the proposal command:
	// - Add "context: fork" to run in forked sub-agent context
	// - Remove "agent" field (not supported by Claude Code slash commands)
	proposalOverrides := &domain.FrontmatterOverride{
		// Set:    map[string]any{"context": "fork"},
		Remove: []string{"agent"},
	}

	return []Initializer{
		NewDirectoryInitializer(
			".claude/commands/spectr",
		),
		NewDirectoryInitializer(".claude/skills"),
		NewConfigFileInitializer(
			"CLAUDE.md",
			tm.InstructionPointer(),
		),
		NewSlashCommandsInitializer(
			".claude/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommandWithOverrides(
					domain.SlashProposal,
					proposalOverrides,
				),
				domain.SlashApply: tm.SlashCommand(
					domain.SlashApply,
				),
			},
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

// Binary returns the CLI binary name for the Claude Code provider.
// Claude Code is invoked via the "claude" binary.
func (*ClaudeProvider) Binary() string {
	return "claude"
}

// InvokeTask creates an exec.Cmd configured to run the Claude Code CLI with the given task.
// Claude Code accepts prompts via stdin using the pattern: echo "prompt" | claude
//
// The command is configured but not started. The orchestrator will attach a PTY and
// start the command to enable full terminal emulation.
//
//nolint:gocritic // task parameter required by ralph.Ralpher interface, unused in this implementation
func (*ClaudeProvider) InvokeTask(
	ctx context.Context,
	_ *ralph.Task,
	prompt string,
) (*exec.Cmd, error) {
	// Create the command using the "claude" binary
	cmd := exec.CommandContext(ctx, "claude")

	// Configure stdin to provide the prompt
	cmd.Stdin = strings.NewReader(prompt)

	// Return the command ready for PTY attachment (not started)
	return cmd, nil
}
