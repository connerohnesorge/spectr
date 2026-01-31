//nolint:revive // Test file uses type assertions without checking - types verified by earlier checks
package providers

import (
	"context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
)

const (
	testClaudeCommandsDir      = ".claude/commands/spectr"
	testGeminiCommandsDir      = ".gemini/commands/spectr"
	testAgentWorkflowsDir      = ".agent/workflows"
	testSpectrPrefix           = "spectr-"
	testCodexPromptsDir        = ".codex/prompts"
	testAcceptSkillName        = "spectr-accept-wo-spectr-bin"
	testValidateSkillName      = "spectr-validate-wo-spectr-bin"
	testProposalSkillPath      = ".agents/skills/spectr-proposal/SKILL.md"
	testApplySkillPath         = ".agents/skills/spectr-apply/SKILL.md"
	testAcceptSkillTargetDir   = ".agents/skills/spectr-accept-wo-spectr-bin"
	testValidateSkillTargetDir = ".agents/skills/spectr-validate-wo-spectr-bin"
	testClaudeAcceptSkillDir   = ".claude/skills/spectr-accept-wo-spectr-bin"
	testClaudeValidateSkillDir = ".claude/skills/spectr-validate-wo-spectr-bin"
	testCodexAcceptSkillDir    = ".codex/skills/spectr-accept-wo-spectr-bin"
	testCodexValidateSkillDir  = ".codex/skills/spectr-validate-wo-spectr-bin"
)

// mockTemplateManager implements TemplateManager for testing
type mockTemplateManager struct{}

func (*mockTemplateManager) InstructionPointer() domain.TemplateRef {
	return domain.TemplateRef{
		Name: "instruction-pointer.md.tmpl",
	}
}

func (*mockTemplateManager) Agents() domain.TemplateRef {
	return domain.TemplateRef{
		Name: "AGENTS.md.tmpl",
	}
}

func (*mockTemplateManager) SlashCommand(
	cmd domain.SlashCommand,
) domain.TemplateRef {
	return domain.TemplateRef{
		Name: fmt.Sprintf(
			"slash-%s.md.tmpl",
			cmd.String(),
		),
		Command: &cmd,
	}
}

func (*mockTemplateManager) SlashCommandWithOverrides(
	cmd domain.SlashCommand,
	overrides *domain.FrontmatterOverride,
) domain.TemplateRef {
	return domain.TemplateRef{
		Name: fmt.Sprintf(
			"slash-%s.md.tmpl",
			cmd.String(),
		),
		Command:   &cmd,
		Overrides: overrides,
	}
}

func (*mockTemplateManager) TOMLSlashCommand(
	cmd domain.SlashCommand,
) domain.TemplateRef {
	return domain.TemplateRef{
		Name: fmt.Sprintf(
			"slash-%s.toml.tmpl",
			cmd.String(),
		),
	}
}

func (*mockTemplateManager) SkillFS(
	skillName string,
) (fs.FS, error) {
	return nil, fmt.Errorf(
		"skill %s not found",
		skillName,
	)
}

func (*mockTemplateManager) ProposalSkill() domain.TemplateRef {
	return domain.TemplateRef{
		Name: "skill-proposal.md.tmpl",
	}
}

func (*mockTemplateManager) ApplySkill() domain.TemplateRef {
	return domain.TemplateRef{
		Name: "skill-apply.md.tmpl",
	}
}

func (*mockTemplateManager) NextSkill() domain.TemplateRef {
	return domain.TemplateRef{
		Name: "skill-next.md.tmpl",
	}
}

// Test each provider returns expected initializers

func TestClaudeProvider_Initializers(
	t *testing.T,
) {
	p := &ClaudeProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Claude should return 6 initializers: Directory (commands), Directory (skills), ConfigFile, SlashCommands, AgentSkills (accept), AgentSkills (validate)
	if len(inits) != 6 {
		t.Fatalf(
			"ClaudeProvider.Initializers() returned %d initializers, want 6",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf(
			"ClaudeProvider.Initializers()[0] is %T, want *DirectoryInitializer",
			inits[0],
		)
	}
	if _, ok := inits[1].(*DirectoryInitializer); !ok {
		t.Errorf(
			"ClaudeProvider.Initializers()[1] is %T, want *DirectoryInitializer",
			inits[1],
		)
	}
	if _, ok := inits[2].(*ConfigFileInitializer); !ok {
		t.Errorf(
			"ClaudeProvider.Initializers()[2] is %T, want *ConfigFileInitializer",
			inits[2],
		)
	}
	if _, ok := inits[3].(*SlashCommandsInitializer); !ok {
		t.Errorf(
			"ClaudeProvider.Initializers()[3] is %T, want *SlashCommandsInitializer",
			inits[3],
		)
	}
	if _, ok := inits[4].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"ClaudeProvider.Initializers()[4] is %T, want *AgentSkillsInitializer",
			inits[4],
		)
	}

	// Check DirectoryInitializer paths
	dirInit := inits[0].(*DirectoryInitializer) //nolint:revive // test code, type checked above
	if len(dirInit.paths) != 1 ||
		dirInit.paths[0] != testClaudeCommandsDir {
		t.Errorf(
			"ClaudeProvider DirectoryInitializer paths = %v, want [\".claude/commands/spectr\"]",
			dirInit.paths,
		)
	}

	// Check second DirectoryInitializer for skills
	skillsDirInit := inits[1].(*DirectoryInitializer) //nolint:revive // test code, type checked above
	if len(skillsDirInit.paths) != 1 ||
		skillsDirInit.paths[0] != ".claude/skills" {
		t.Errorf(
			"ClaudeProvider DirectoryInitializer[1] paths = %v, want [\".claude/skills\"]",
			skillsDirInit.paths,
		)
	}

	// Check ConfigFileInitializer path
	cfgInit := inits[2].(*ConfigFileInitializer) //nolint:revive // test code, type checked above
	if cfgInit.path != "CLAUDE.md" {
		t.Errorf(
			"ClaudeProvider ConfigFileInitializer path = %s, want \"CLAUDE.md\"",
			cfgInit.path,
		)
	}

	// Check SlashCommandsInitializer dir
	slashInit := inits[3].(*SlashCommandsInitializer) //nolint:revive // test code, type checked above
	if slashInit.dir != testClaudeCommandsDir {
		t.Errorf(
			"ClaudeProvider SlashCommandsInitializer dir = %s, want \".claude/commands/spectr\"",
			slashInit.dir,
		)
	}
	if len(slashInit.commands) != 3 {
		t.Errorf(
			"ClaudeProvider SlashCommandsInitializer has %d commands, want 3",
			len(slashInit.commands),
		)
	}

	// Check AgentSkillsInitializer
	skillInit := inits[4].(*AgentSkillsInitializer) //nolint:revive // test code, type checked above
	if skillInit.skillName != testAcceptSkillName {
		t.Errorf(
			"ClaudeProvider AgentSkillsInitializer skillName = %s, want %q",
			skillInit.skillName,
			testAcceptSkillName,
		)
	}
	if skillInit.targetDir != testClaudeAcceptSkillDir {
		t.Errorf(
			"ClaudeProvider AgentSkillsInitializer targetDir = %s, want %q",
			skillInit.targetDir,
			testClaudeAcceptSkillDir,
		)
	}
}

func TestGeminiProvider_Initializers(
	t *testing.T,
) {
	p := &GeminiProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Gemini should return 2 initializers: Directory, TOMLSlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf(
			"GeminiProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf(
			"GeminiProvider.Initializers()[0] is %T, want *DirectoryInitializer",
			inits[0],
		)
	}
	if _, ok := inits[1].(*TOMLSlashCommandsInitializer); !ok {
		t.Errorf(
			"GeminiProvider.Initializers()[1] is %T, want *TOMLSlashCommandsInitializer",
			inits[1],
		)
	}

	// Check DirectoryInitializer paths
	dirInit := inits[0].(*DirectoryInitializer)
	if len(dirInit.paths) != 1 ||
		dirInit.paths[0] != testGeminiCommandsDir {
		t.Errorf(
			"GeminiProvider DirectoryInitializer paths = %v, want [\".gemini/commands/spectr\"]",
			dirInit.paths,
		)
	}

	// Check TOMLSlashCommandsInitializer dir
	slashInit := inits[1].(*TOMLSlashCommandsInitializer)
	if slashInit.dir != testGeminiCommandsDir {
		t.Errorf(
			"GeminiProvider TOMLSlashCommandsInitializer dir = %s, want \".gemini/commands/spectr\"",
			slashInit.dir,
		)
	}
	if len(slashInit.commands) != 3 {
		t.Errorf(
			"GeminiProvider TOMLSlashCommandsInitializer has %d commands, want 3",
			len(slashInit.commands),
		)
	}
}

func TestCostrictProvider_Initializers(
	t *testing.T,
) {
	p := &CostrictProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Costrict should return 3 initializers: Directory, ConfigFile, SlashCommands
	if len(inits) != 3 {
		t.Fatalf(
			"CostrictProvider.Initializers() returned %d initializers, want 3",
			len(inits),
		)
	}

	// Check ConfigFileInitializer path
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "COSTRICT.md" {
		t.Errorf(
			"CostrictProvider ConfigFileInitializer path = %s, want \"COSTRICT.md\"",
			cfgInit.path,
		)
	}
}

func TestQoderProvider_Initializers(
	t *testing.T,
) {
	p := &QoderProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Qoder should return 3 initializers
	if len(inits) != 3 {
		t.Fatalf(
			"QoderProvider.Initializers() returned %d initializers, want 3",
			len(inits),
		)
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "QODER.md" {
		t.Errorf(
			"QoderProvider ConfigFileInitializer path = %s, want \"QODER.md\"",
			cfgInit.path,
		)
	}
}

func TestQwenProvider_Initializers(t *testing.T) {
	p := &QwenProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Qwen should return 3 initializers
	if len(inits) != 3 {
		t.Fatalf(
			"QwenProvider.Initializers() returned %d initializers, want 3",
			len(inits),
		)
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "QWEN.md" {
		t.Errorf(
			"QwenProvider ConfigFileInitializer path = %s, want \"QWEN.md\"",
			cfgInit.path,
		)
	}
}

func TestAntigravityProvider_Initializers(
	t *testing.T,
) {
	p := &AntigravityProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Antigravity should return 3 initializers: Directory, ConfigFile, PrefixedSlashCommands
	if len(inits) != 3 {
		t.Fatalf(
			"AntigravityProvider.Initializers() returned %d initializers, want 3",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[2].(*PrefixedSlashCommandsInitializer); !ok {
		t.Errorf(
			"AntigravityProvider.Initializers()[2] is %T, want *PrefixedSlashCommandsInitializer",
			inits[2],
		)
	}

	// Check PrefixedSlashCommandsInitializer
	slashInit := inits[2].(*PrefixedSlashCommandsInitializer)
	if slashInit.dir != testAgentWorkflowsDir {
		t.Errorf(
			"AntigravityProvider PrefixedSlashCommandsInitializer dir = %s, want \".agent/workflows\"",
			slashInit.dir,
		)
	}
	if slashInit.prefix != testSpectrPrefix {
		t.Errorf(
			"AntigravityProvider PrefixedSlashCommandsInitializer prefix = %s, want \"spectr-\"",
			slashInit.prefix,
		)
	}

	// Check ConfigFileInitializer uses AGENTS.md
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "AGENTS.md" {
		t.Errorf(
			"AntigravityProvider ConfigFileInitializer path = %s, want \"AGENTS.md\"",
			cfgInit.path,
		)
	}
}

func TestClineProvider_Initializers(
	t *testing.T,
) {
	p := &ClineProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 3 {
		t.Fatalf(
			"ClineProvider.Initializers() returned %d initializers, want 3",
			len(inits),
		)
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "CLINE.md" {
		t.Errorf(
			"ClineProvider ConfigFileInitializer path = %s, want \"CLINE.md\"",
			cfgInit.path,
		)
	}
}

func TestCursorProvider_Initializers(
	t *testing.T,
) {
	p := &CursorProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Cursor should return 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf(
			"CursorProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf(
			"CursorProvider.Initializers()[0] is %T, want *DirectoryInitializer",
			inits[0],
		)
	}
	if _, ok := inits[1].(*SlashCommandsInitializer); !ok {
		t.Errorf(
			"CursorProvider.Initializers()[1] is %T, want *SlashCommandsInitializer",
			inits[1],
		)
	}

	dirInit := inits[0].(*DirectoryInitializer)
	if len(dirInit.paths) != 1 ||
		dirInit.paths[0] != ".cursorrules/commands/spectr" {
		t.Errorf(
			"CursorProvider DirectoryInitializer paths = %v, want [\".cursorrules/commands/spectr\"]",
			dirInit.paths,
		)
	}
}

func TestCodexProvider_Initializers(
	t *testing.T,
) {
	p := &CodexProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Codex should return 6 initializers: HomeDirectory, Directory, ConfigFile, HomePrefixedSlashCommands, 2x AgentSkills
	if len(inits) != 6 {
		t.Fatalf(
			"CodexProvider.Initializers() returned %d initializers, want 6",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[0].(*HomeDirectoryInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[0] is %T, want *HomeDirectoryInitializer",
			inits[0],
		)
	}
	if _, ok := inits[1].(*DirectoryInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[1] is %T, want *DirectoryInitializer",
			inits[1],
		)
	}
	if _, ok := inits[2].(*ConfigFileInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[2] is %T, want *ConfigFileInitializer",
			inits[2],
		)
	}
	if _, ok := inits[3].(*HomePrefixedSlashCommandsInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[3] is %T, want *HomePrefixedSlashCommandsInitializer",
			inits[3],
		)
	}
	if _, ok := inits[4].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[4] is %T, want *AgentSkillsInitializer",
			inits[4],
		)
	}
	if _, ok := inits[5].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"CodexProvider.Initializers()[5] is %T, want *AgentSkillsInitializer",
			inits[5],
		)
	}

	// Check HomeDirectoryInitializer paths
	dirInit := inits[0].(*HomeDirectoryInitializer)
	if len(dirInit.paths) != 1 ||
		dirInit.paths[0] != testCodexPromptsDir {
		t.Errorf(
			"CodexProvider HomeDirectoryInitializer paths = %v, want [\".codex/prompts\"]",
			dirInit.paths,
		)
	}

	// Check DirectoryInitializer paths
	skillsDirInit := inits[1].(*DirectoryInitializer)
	if len(skillsDirInit.paths) != 1 ||
		skillsDirInit.paths[0] != ".codex/skills" {
		t.Errorf(
			"CodexProvider DirectoryInitializer paths = %v, want [\".codex/skills\"]",
			skillsDirInit.paths,
		)
	}

	// Check HomePrefixedSlashCommandsInitializer
	slashInit := inits[3].(*HomePrefixedSlashCommandsInitializer)
	if slashInit.dir != testCodexPromptsDir {
		t.Errorf(
			"CodexProvider HomePrefixedSlashCommandsInitializer dir = %s, want \".codex/prompts\"",
			slashInit.dir,
		)
	}
	if slashInit.prefix != testSpectrPrefix {
		t.Errorf(
			"CodexProvider HomePrefixedSlashCommandsInitializer prefix = %s, want \"spectr-\"",
			slashInit.prefix,
		)
	}

	// Check AgentSkillsInitializer instances
	acceptSkill := inits[4].(*AgentSkillsInitializer)
	if acceptSkill.skillName != testAcceptSkillName {
		t.Errorf(
			"CodexProvider AgentSkillsInitializer[4] skillName = %s, want %q",
			acceptSkill.skillName,
			testAcceptSkillName,
		)
	}
	if acceptSkill.targetDir != testCodexAcceptSkillDir {
		t.Errorf(
			"CodexProvider AgentSkillsInitializer[4] targetDir = %s, want %q",
			acceptSkill.targetDir,
			testCodexAcceptSkillDir,
		)
	}

	validateSkill := inits[5].(*AgentSkillsInitializer)
	if validateSkill.skillName != testValidateSkillName {
		t.Errorf(
			"CodexProvider AgentSkillsInitializer[5] skillName = %s, want %q",
			validateSkill.skillName,
			testValidateSkillName,
		)
	}
	if validateSkill.targetDir != testCodexValidateSkillDir {
		t.Errorf(
			"CodexProvider AgentSkillsInitializer[5] targetDir = %s, want %q",
			validateSkill.targetDir,
			testCodexValidateSkillDir,
		)
	}

	// Check ConfigFileInitializer uses AGENTS.md
	cfgInit := inits[2].(*ConfigFileInitializer)
	if cfgInit.path != "AGENTS.md" {
		t.Errorf(
			"CodexProvider ConfigFileInitializer path = %s, want \"AGENTS.md\"",
			cfgInit.path,
		)
	}
}

func TestAiderProvider_Initializers(
	t *testing.T,
) {
	p := &AiderProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Aider should return 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf(
			"AiderProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}
}

func TestWindsurfProvider_Initializers(
	t *testing.T,
) {
	p := &WindsurfProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 2 {
		t.Fatalf(
			"WindsurfProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}
}

func TestKilocodeProvider_Initializers(
	t *testing.T,
) {
	p := &KilocodeProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 2 {
		t.Fatalf(
			"KilocodeProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}
}

func TestContinueProvider_Initializers(
	t *testing.T,
) {
	p := &ContinueProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	if len(inits) != 2 {
		t.Fatalf(
			"ContinueProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}
}

func TestCrushProvider_Initializers(
	t *testing.T,
) {
	p := &CrushProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Crush should return 3 initializers
	if len(inits) != 3 {
		t.Fatalf(
			"CrushProvider.Initializers() returned %d initializers, want 3",
			len(inits),
		)
	}

	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "CRUSH.md" {
		t.Errorf(
			"CrushProvider ConfigFileInitializer path = %s, want \"CRUSH.md\"",
			cfgInit.path,
		)
	}
}

func TestOpencodeProvider_Initializers(
	t *testing.T,
) {
	p := &OpencodeProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Opencode should return 2 initializers: Directory, SlashCommands (no config file)
	if len(inits) != 2 {
		t.Fatalf(
			"OpencodeProvider.Initializers() returned %d initializers, want 2",
			len(inits),
		)
	}
}

func TestKimiProvider_Initializers(
	t *testing.T,
) {
	p := &KimiProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Kimi should return 7 initializers: Directory (skills), ConfigFile, 3x SkillFile, 2x AgentSkills
	if len(inits) != 7 {
		t.Fatalf(
			"KimiProvider.Initializers() returned %d initializers, want 7",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[0] is %T, want *DirectoryInitializer",
			inits[0],
		)
	}
	if _, ok := inits[1].(*ConfigFileInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[1] is %T, want *ConfigFileInitializer",
			inits[1],
		)
	}
	if _, ok := inits[2].(*SkillFileInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[2] is %T, want *SkillFileInitializer",
			inits[2],
		)
	}
	if _, ok := inits[3].(*SkillFileInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[3] is %T, want *SkillFileInitializer",
			inits[3],
		)
	}
	if _, ok := inits[4].(*SkillFileInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[4] is %T, want *SkillFileInitializer",
			inits[4],
		)
	}
	if _, ok := inits[5].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[5] is %T, want *AgentSkillsInitializer",
			inits[5],
		)
	}
	if _, ok := inits[6].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"KimiProvider.Initializers()[6] is %T, want *AgentSkillsInitializer",
			inits[6],
		)
	}

	// Check DirectoryInitializer paths
	dirInit := inits[0].(*DirectoryInitializer)
	if len(dirInit.paths) != 1 ||
		dirInit.paths[0] != ".claude/skills" {
		t.Errorf(
			"KimiProvider DirectoryInitializer paths = %v, want [\".claude/skills\"]",
			dirInit.paths,
		)
	}

	// Check ConfigFileInitializer uses AGENTS.md
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "AGENTS.md" {
		t.Errorf(
			"KimiProvider ConfigFileInitializer path = %s, want \"AGENTS.md\"",
			cfgInit.path,
		)
	}

	// Check SkillFileInitializer paths
	proposalSkill := inits[2].(*SkillFileInitializer)
	if proposalSkill.targetPath != ".claude/skills/spectr-proposal/SKILL.md" {
		t.Errorf(
			"KimiProvider SkillFileInitializer[2] targetPath = %s, want \".claude/skills/spectr-proposal/SKILL.md\"",
			proposalSkill.targetPath,
		)
	}

	applySkill := inits[3].(*SkillFileInitializer)
	if applySkill.targetPath != ".claude/skills/spectr-apply/SKILL.md" {
		t.Errorf(
			"KimiProvider SkillFileInitializer[3] targetPath = %s, want \".claude/skills/spectr-apply/SKILL.md\"",
			applySkill.targetPath,
		)
	}

	nextSkill := inits[4].(*SkillFileInitializer)
	if nextSkill.targetPath != ".claude/skills/spectr-next/SKILL.md" {
		t.Errorf(
			"KimiProvider SkillFileInitializer[4] targetPath = %s, want \".claude/skills/spectr-next/SKILL.md\"",
			nextSkill.targetPath,
		)
	}
}

func TestAmpProvider_Initializers(
	t *testing.T,
) {
	p := &AmpProvider{}
	ctx := context.Background()
	tm := &mockTemplateManager{}

	inits := p.Initializers(ctx, tm)

	// Amp should return 6 initializers: Directory, ConfigFile, 2x SkillFile, 2x AgentSkills
	if len(inits) != 6 {
		t.Fatalf(
			"AmpProvider.Initializers() returned %d initializers, want 6",
			len(inits),
		)
	}

	// Check types
	if _, ok := inits[0].(*DirectoryInitializer); !ok {
		t.Errorf(
			"AmpProvider.Initializers()[0] is %T, want *DirectoryInitializer",
			inits[0],
		)
	}
	if _, ok := inits[1].(*ConfigFileInitializer); !ok {
		t.Errorf(
			"AmpProvider.Initializers()[1] is %T, want *ConfigFileInitializer",
			inits[1],
		)
	}
	if _, ok := inits[2].(*SkillFileInitializer); !ok {
		t.Errorf(
			"AmpProvider.Initializers()[2] is %T, want *SkillFileInitializer",
			inits[2],
		)
	}
	if _, ok := inits[3].(*SkillFileInitializer); !ok {
		t.Errorf(
			"AmpProvider.Initializers()[3] is %T, want *SkillFileInitializer",
			inits[3],
		)
	}
	if _, ok := inits[4].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"AmpProvider.Initializers()[4] is %T, want *AgentSkillsInitializer",
			inits[4],
		)
	}
	if _, ok := inits[5].(*AgentSkillsInitializer); !ok {
		t.Errorf(
			"AmpProvider.Initializers()[5] is %T, want *AgentSkillsInitializer",
			inits[5],
		)
	}

	// Check DirectoryInitializer paths
	dirInit := inits[0].(*DirectoryInitializer)
	if len(dirInit.paths) != 1 ||
		dirInit.paths[0] != ".agents/skills" {
		t.Errorf(
			"AmpProvider DirectoryInitializer paths = %v, want [\".agents/skills\"]",
			dirInit.paths,
		)
	}

	// Check ConfigFileInitializer path
	cfgInit := inits[1].(*ConfigFileInitializer)
	if cfgInit.path != "AMP.md" {
		t.Errorf(
			"AmpProvider ConfigFileInitializer path = %s, want \"AMP.md\"",
			cfgInit.path,
		)
	}

	// Check SkillFileInitializer paths
	proposalSkill := inits[2].(*SkillFileInitializer)
	if proposalSkill.targetPath != testProposalSkillPath {
		t.Errorf(
			"AmpProvider SkillFileInitializer[2] targetPath = %s, want %q",
			proposalSkill.targetPath,
			testProposalSkillPath,
		)
	}

	applySkill := inits[3].(*SkillFileInitializer)
	if applySkill.targetPath != testApplySkillPath {
		t.Errorf(
			"AmpProvider SkillFileInitializer[3] targetPath = %s, want %q",
			applySkill.targetPath,
			testApplySkillPath,
		)
	}

	// Check AgentSkillsInitializer instances
	acceptSkill := inits[4].(*AgentSkillsInitializer)
	if acceptSkill.skillName != testAcceptSkillName {
		t.Errorf(
			"AmpProvider AgentSkillsInitializer[4] skillName = %s, want %q",
			acceptSkill.skillName,
			testAcceptSkillName,
		)
	}
	if acceptSkill.targetDir != testAcceptSkillTargetDir {
		t.Errorf(
			"AmpProvider AgentSkillsInitializer[4] targetDir = %s, want %q",
			acceptSkill.targetDir,
			testAcceptSkillTargetDir,
		)
	}

	validateSkill := inits[5].(*AgentSkillsInitializer)
	if validateSkill.skillName != testValidateSkillName {
		t.Errorf(
			"AmpProvider AgentSkillsInitializer[5] skillName = %s, want %q",
			validateSkill.skillName,
			testValidateSkillName,
		)
	}
	if validateSkill.targetDir != testValidateSkillTargetDir {
		t.Errorf(
			"AmpProvider AgentSkillsInitializer[5] targetDir = %s, want %q",
			validateSkill.targetDir,
			testValidateSkillTargetDir,
		)
	}
}

// Test all 15 providers return expected initializer counts and types

func TestAllProviders_InitializerCounts(
	t *testing.T,
) {
	ctx := context.Background()
	tm := &mockTemplateManager{}

	tests := []struct {
		name          string
		provider      Provider
		expectedCount int
		hasConfigFile bool
		usesTOML      bool
		usesPrefix    bool
		usesHomeFs    bool
	}{
		{
			"claude-code",
			&ClaudeProvider{},
			6,
			true,
			false,
			false,
			false,
		},
		{
			"gemini",
			&GeminiProvider{},
			2,
			false,
			true,
			false,
			false,
		},
		{
			"costrict",
			&CostrictProvider{},
			3,
			true,
			false,
			false,
			false,
		},
		{
			"qoder",
			&QoderProvider{},
			3,
			true,
			false,
			false,
			false,
		},
		{
			"qwen",
			&QwenProvider{},
			3,
			true,
			false,
			false,
			false,
		},
		{
			"antigravity",
			&AntigravityProvider{},
			3,
			true,
			false,
			true,
			false,
		},
		{
			"cline",
			&ClineProvider{},
			3,
			true,
			false,
			false,
			false,
		},
		{
			"cursor",
			&CursorProvider{},
			2,
			false,
			false,
			false,
			false,
		},
		{
			"codex",
			&CodexProvider{},
			6,
			true,
			false,
			true,
			true,
		},
		{
			"aider",
			&AiderProvider{},
			2,
			false,
			false,
			false,
			false,
		},
		{
			"windsurf",
			&WindsurfProvider{},
			2,
			false,
			false,
			false,
			false,
		},
		{
			"kilocode",
			&KilocodeProvider{},
			2,
			false,
			false,
			false,
			false,
		},
		{
			"continue",
			&ContinueProvider{},
			2,
			false,
			false,
			false,
			false,
		},
		{
			"crush",
			&CrushProvider{},
			3,
			true,
			false,
			false,
			false,
		},
		{
			"opencode",
			&OpencodeProvider{},
			2,
			false,
			false,
			false,
			false,
		},
		{
			"amp",
			&AmpProvider{},
			6,
			true,
			false,
			false,
			false,
		},
		{
			"kimi",
			&KimiProvider{},
			7,
			true,
			false,
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inits := tt.provider.Initializers(
				ctx,
				tm,
			)

			// Check count
			if len(inits) != tt.expectedCount {
				t.Errorf(
					"%s returned %d initializers, want %d",
					tt.name,
					len(inits),
					tt.expectedCount,
				)
			}

			// Check for directory initializer (all should have one)
			hasDir := false
			hasHomeDir := false
			hasConfig := false
			hasSlash := false
			hasHomeSlash := false
			hasToml := false
			hasPrefix := false
			hasHomePrefix := false
			hasSkillFile := false

			for _, init := range inits {
				switch init.(type) {
				case *DirectoryInitializer:
					hasDir = true
				case *HomeDirectoryInitializer:
					hasHomeDir = true
				case *ConfigFileInitializer:
					hasConfig = true
				case *SlashCommandsInitializer:
					hasSlash = true
				case *HomeSlashCommandsInitializer:
					hasHomeSlash = true
				case *TOMLSlashCommandsInitializer:
					hasToml = true
				case *PrefixedSlashCommandsInitializer:
					hasPrefix = true
				case *HomePrefixedSlashCommandsInitializer:
					hasHomePrefix = true
				case *SkillFileInitializer:
					hasSkillFile = true
				}
			}

			// All providers should have either project or home directory
			if !hasDir && !hasHomeDir {
				t.Errorf(
					"%s has no directory initializer",
					tt.name,
				)
			}

			// Check config file expectation
			if hasConfig != tt.hasConfigFile {
				t.Errorf(
					"%s hasConfigFile = %v, want %v",
					tt.name,
					hasConfig,
					tt.hasConfigFile,
				)
			}

			// Check TOML expectation
			if hasToml != tt.usesTOML {
				t.Errorf(
					"%s usesTOML = %v, want %v",
					tt.name,
					hasToml,
					tt.usesTOML,
				)
			}

			// Check prefix expectation
			if (hasPrefix || hasHomePrefix) != tt.usesPrefix {
				t.Errorf(
					"%s usesPrefix = %v, want %v",
					tt.name,
					(hasPrefix || hasHomePrefix),
					tt.usesPrefix,
				)
			}

			// Check home filesystem expectation
			if (hasHomeDir || hasHomeSlash || hasHomePrefix) != tt.usesHomeFs {
				t.Errorf(
					"%s usesHomeFs = %v, want %v",
					tt.name,
					(hasHomeDir || hasHomeSlash || hasHomePrefix),
					tt.usesHomeFs,
				)
			}

			// All providers should have slash commands or skills in some form
			if !hasSlash && !hasHomeSlash &&
				!hasToml &&
				!hasPrefix &&
				!hasHomePrefix &&
				!hasSkillFile {
				t.Errorf(
					"%s has no slash command or skill initializer",
					tt.name,
				)
			}
		})
	}
}

// Test provider registration metadata

func TestProviderRegistration_AllProviders(
	t *testing.T,
) {
	// Reset and register all providers
	Reset()
	err := RegisterAllProviders()
	if err != nil {
		t.Fatalf(
			"RegisterAllProviders() failed: %v",
			err,
		)
	}

	// Expected provider metadata
	expected := []struct {
		id       string
		name     string
		priority int
	}{
		{"claude-code", "Claude Code", 1},
		{"gemini", "Gemini CLI", 2},
		{"costrict", "Costrict", 3},
		{"qoder", "Qoder", 4},
		{"qwen", "Qwen Code", 5},
		{"antigravity", "Antigravity", 6},
		{"cline", "Cline", 7},
		{"cursor", "Cursor", 8},
		{"codex", "Codex CLI", 9},
		{"aider", "Aider", 10},
		{"windsurf", "Windsurf", 11},
		{"kilocode", "Kilocode", 12},
		{"continue", "Continue", 13},
		{"crush", "Crush", 14},
		{"amp", "Amp", 15},
		{"opencode", "OpenCode", 16},
		{"kimi", "Kimi", 17},
	}

	// Verify count
	if Count() != 17 {
		t.Fatalf("Count() = %d, want 17", Count())
	}

	// Verify each provider
	for _, exp := range expected {
		reg, ok := Get(exp.id)
		if !ok {
			t.Errorf(
				"Provider %s not found in registry",
				exp.id,
			)

			continue
		}

		if reg.ID != exp.id {
			t.Errorf(
				"Provider %s has ID %s, want %s",
				exp.id,
				reg.ID,
				exp.id,
			)
		}
		if reg.Name != exp.name {
			t.Errorf(
				"Provider %s has Name %s, want %s",
				exp.id,
				reg.Name,
				exp.name,
			)
		}
		if reg.Priority != exp.priority {
			t.Errorf(
				"Provider %s has Priority %d, want %d",
				exp.id,
				reg.Priority,
				exp.priority,
			)
		}
		if reg.Provider == nil {
			t.Errorf(
				"Provider %s has nil Provider",
				exp.id,
			)
		}
	}

	// Verify priority order
	registered := RegisteredProviders()

	// Note: Priorities are no longer strictly sequential due to Amp (priority 15)
	// being inserted between Claude Code and other providers
	expectedOrder := []string{
		"claude-code", "gemini", "costrict", "qoder", "qwen",
		"antigravity", "cline", "cursor", "codex", "aider",
		"windsurf", "kilocode", "continue", "crush", "amp", "opencode", "kimi",
	}

	if len(registered) != len(expectedOrder) {
		t.Fatalf(
			"RegisteredProviders() returned %d providers, want %d",
			len(registered),
			len(expectedOrder),
		)
	}

	for i, reg := range registered {
		if reg.ID != expectedOrder[i] {
			t.Errorf(
				"RegisteredProviders()[%d].ID = %s, want %s (priority order incorrect)",
				i,
				reg.ID,
				expectedOrder[i],
			)
		}
	}
}
