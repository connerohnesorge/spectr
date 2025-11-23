//nolint:revive // line-length-limit - tool configurations benefit from clarity
package init

// ToolID is a type-safe identifier for AI tools
type ToolID string

// Tool ID constants for all supported AI tools
const (
	// Config-based tools (create instruction files like CLAUDE.md)
	ToolClaudeCode     ToolID = "claude-code"
	ToolCline          ToolID = "cline"
	ToolCostrictConfig ToolID = "costrict-config"
	ToolQoderConfig    ToolID = "qoder-config"
	ToolCodeBuddy      ToolID = "codebuddy"
	ToolQwen           ToolID = "qwen"
	ToolAntigravity    ToolID = "antigravity"

	// Slash command tools (create files in .claude/commands/)
	ToolClaude           ToolID = "claude"
	ToolClineSlash       ToolID = "cline-slash"
	ToolKilocode         ToolID = "kilocode"
	ToolQoderSlash       ToolID = "qoder-slash"
	ToolCursor           ToolID = "cursor"
	ToolAider            ToolID = "aider"
	ToolContinue         ToolID = "continue"
	ToolCopilot          ToolID = "copilot"
	ToolMentat           ToolID = "mentat"
	ToolTabnine          ToolID = "tabnine"
	ToolSmol             ToolID = "smol"
	ToolCostrictSlash    ToolID = "costrict-slash"
	ToolWindsurf         ToolID = "windsurf"
	ToolCodeBuddySlash   ToolID = "codebuddy-slash"
	ToolQwenSlash        ToolID = "qwen-slash"
	ToolAntigravitySlash ToolID = "antigravity-slash"
)

// ToolConfig holds all configuration data for a tool
// This enables data-driven tool registration without individual structs
type ToolConfig struct {
	// ID is the unique type-safe identifier for the tool
	ID ToolID
	// Name is the human-readable name shown in the UI
	Name string
	// Type indicates whether this is a config or slash tool
	Type ToolType
	// ConfigFile is the instruction file path (for config-based tools)
	// Examples: "CLAUDE.md", "CLINE.md", "AGENTS.md"
	ConfigFile string
	// SlashPaths maps command names to file paths (for slash command tools)
	// Keys: "proposal", "apply", "archive"
	// Values: ".claude/commands/spectr-proposal.md", etc.
	SlashPaths map[string]string
	// Frontmatter contains metadata for slash command files
	// Keys match SlashPaths keys
	Frontmatter map[string]string
	// Priority determines display order in the UI (lower numbers first)
	Priority int
}

// toolConfigs is the central registry of all config-based tool configurations
var toolConfigs = map[ToolID]ToolConfig{
	ToolClaudeCode: {
		ID:         ToolClaudeCode,
		Name:       "Claude Code",
		Type:       ToolTypeConfig,
		ConfigFile: "CLAUDE.md",
		Priority:   1,
	},
	ToolCline: {
		ID:         ToolCline,
		Name:       "Cline",
		Type:       ToolTypeConfig,
		ConfigFile: "CLINE.md",
		Priority:   2,
	},
	ToolCostrictConfig: {
		ID:         ToolCostrictConfig,
		Name:       "CoStrict",
		Type:       ToolTypeConfig,
		ConfigFile: "COSTRICT.md",
		Priority:   3,
	},
	ToolQoderConfig: {
		ID:         ToolQoderConfig,
		Name:       "Qoder",
		Type:       ToolTypeConfig,
		ConfigFile: "QODER.md",
		Priority:   4,
	},
	ToolCodeBuddy: {
		ID:         ToolCodeBuddy,
		Name:       "CodeBuddy",
		Type:       ToolTypeConfig,
		ConfigFile: "CODEBUDDY.md",
		Priority:   5,
	},
	ToolQwen: {
		ID:         ToolQwen,
		Name:       "Qwen Code",
		Type:       ToolTypeConfig,
		ConfigFile: "QWEN.md",
		Priority:   6,
	},
	ToolAntigravity: {
		ID:         ToolAntigravity,
		Name:       "Antigravity",
		Type:       ToolTypeConfig,
		ConfigFile: "AGENTS.md",
		Priority:   7,
	},
}

// slashToolConfigs is the central registry of all slash command tool configurations
var slashToolConfigs = map[ToolID]ToolConfig{
	ToolClaude: {
		ID:   ToolClaude,
		Name: "Claude",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".claude/commands/spectr-proposal.md",
			"apply":    ".claude/commands/spectr-apply.md",
			"archive":  ".claude/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly. (project)\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync. (project)\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs. (project)\n---",
		},
		Priority: 1,
	},
	ToolClineSlash: {
		ID:   ToolClineSlash,
		Name: "Cline",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".clinerules/commands/spectr-proposal.md",
			"apply":    ".clinerules/commands/spectr-apply.md",
			"archive":  ".clinerules/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 2,
	},
	ToolKilocode: {
		ID:   ToolKilocode,
		Name: "Kilocode",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".kilocode/commands/spectr-proposal.md",
			"apply":    ".kilocode/commands/spectr-apply.md",
			"archive":  ".kilocode/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 3,
	},
	ToolQoderSlash: {
		ID:   ToolQoderSlash,
		Name: "Qoder",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".qoder/commands/spectr-proposal.md",
			"apply":    ".qoder/commands/spectr-apply.md",
			"archive":  ".qoder/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 4,
	},
	ToolCursor: {
		ID:   ToolCursor,
		Name: "Cursor",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".cursorrules/commands/spectr-proposal.md",
			"apply":    ".cursorrules/commands/spectr-apply.md",
			"archive":  ".cursorrules/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 5,
	},
	ToolAider: {
		ID:   ToolAider,
		Name: "Aider",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".aider/commands/spectr-proposal.md",
			"apply":    ".aider/commands/spectr-apply.md",
			"archive":  ".aider/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 6,
	},
	ToolContinue: {
		ID:   ToolContinue,
		Name: "Continue",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".continue/commands/spectr-proposal.md",
			"apply":    ".continue/commands/spectr-apply.md",
			"archive":  ".continue/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 7,
	},
	ToolCopilot: {
		ID:   ToolCopilot,
		Name: "Copilot",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".github/copilot/commands/spectr-proposal.md",
			"apply":    ".github/copilot/commands/spectr-apply.md",
			"archive":  ".github/copilot/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 8,
	},
	ToolMentat: {
		ID:   ToolMentat,
		Name: "Mentat",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".mentat/commands/spectr-proposal.md",
			"apply":    ".mentat/commands/spectr-apply.md",
			"archive":  ".mentat/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 9,
	},
	ToolTabnine: {
		ID:   ToolTabnine,
		Name: "Tabnine",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".tabnine/commands/spectr-proposal.md",
			"apply":    ".tabnine/commands/spectr-apply.md",
			"archive":  ".tabnine/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 10,
	},
	ToolSmol: {
		ID:   ToolSmol,
		Name: "Smol",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".smol/commands/spectr-proposal.md",
			"apply":    ".smol/commands/spectr-apply.md",
			"archive":  ".smol/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 11,
	},
	ToolCostrictSlash: {
		ID:   ToolCostrictSlash,
		Name: "Costrict",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".costrict/commands/spectr-proposal.md",
			"apply":    ".costrict/commands/spectr-apply.md",
			"archive":  ".costrict/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 12,
	},
	ToolWindsurf: {
		ID:   ToolWindsurf,
		Name: "Windsurf",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".windsurf/commands/spectr-proposal.md",
			"apply":    ".windsurf/commands/spectr-apply.md",
			"archive":  ".windsurf/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 13,
	},
	ToolCodeBuddySlash: {
		ID:   ToolCodeBuddySlash,
		Name: "CodeBuddy",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".codebuddy/commands/spectr-proposal.md",
			"apply":    ".codebuddy/commands/spectr-apply.md",
			"archive":  ".codebuddy/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 14,
	},
	ToolQwenSlash: {
		ID:   ToolQwenSlash,
		Name: "Qwen",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".qwen/commands/spectr-proposal.md",
			"apply":    ".qwen/commands/spectr-apply.md",
			"archive":  ".qwen/commands/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 15,
	},
	ToolAntigravitySlash: {
		ID:   ToolAntigravitySlash,
		Name: "Antigravity",
		Type: ToolTypeSlash,
		SlashPaths: map[string]string{
			"proposal": ".agent/workflows/spectr-proposal.md",
			"apply":    ".agent/workflows/spectr-apply.md",
			"archive":  ".agent/workflows/spectr-archive.md",
		},
		Frontmatter: map[string]string{
			"proposal": "---\ndescription: Scaffold a new Spectr change and validate strictly.\n---",
			"apply":    "---\ndescription: Implement an approved Spectr change and keep tasks in sync.\n---",
			"archive":  "---\ndescription: Archive a deployed Spectr change and update specs.\n---",
		},
		Priority: 16,
	},
}

// configToSlashMapping maps config-based tool IDs to their corresponding slash tool IDs
// This enables automatic slash command installation when a config tool is selected
var configToSlashMapping = map[ToolID]ToolID{
	ToolClaudeCode:     ToolClaude,
	ToolCline:          ToolClineSlash,
	ToolCostrictConfig: ToolCostrictSlash,
	ToolQoderConfig:    ToolQoderSlash,
	ToolCodeBuddy:      ToolCodeBuddySlash,
	ToolQwen:           ToolQwenSlash,
	ToolAntigravity:    ToolAntigravitySlash,
}

// GetToolConfig retrieves the configuration for a tool by its ID
// Returns the config and a boolean indicating whether it was found
func GetToolConfig(id ToolID) (ToolConfig, bool) {
	// Check config-based tools first
	if config, ok := toolConfigs[id]; ok {
		return config, true
	}
	// Check slash command tools
	if config, ok := slashToolConfigs[id]; ok {
		return config, true
	}

	return ToolConfig{}, false
}

// GetSlashToolForConfig returns the slash command tool ID for a config tool ID
// Returns the slash tool ID and a boolean indicating whether a mapping exists
func GetSlashToolForConfig(configID ToolID) (ToolID, bool) {
	slashID, ok := configToSlashMapping[configID]

	return slashID, ok
}

// GetAllConfigTools returns all config-based tool configurations
func GetAllConfigTools() []ToolConfig {
	configs := make([]ToolConfig, 0, len(toolConfigs))
	for _, config := range toolConfigs {
		configs = append(configs, config)
	}

	return configs
}

// GetAllSlashTools returns all slash command tool configurations
func GetAllSlashTools() []ToolConfig {
	configs := make([]ToolConfig, 0, len(slashToolConfigs))
	for _, config := range slashToolConfigs {
		configs = append(configs, config)
	}

	return configs
}
