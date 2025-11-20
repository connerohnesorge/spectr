package init

// ToolProvider is the main interface for all tool providers.
// It composes memory file and slash command provider capabilities.
type ToolProvider interface {
	// GetName returns the human-readable name of the tool
	GetName() string
	// GetMemoryFileProvider returns the memory file provider if this tool
	// uses one, nil otherwise
	GetMemoryFileProvider() MemoryFileProvider
	// GetSlashCommandProvider returns the slash command provider if this
	// tool uses one, nil otherwise
	GetSlashCommandProvider() SlashCommandProvider
}

// MemoryFileProvider handles configuration of memory files (like CLAUDE.md,
// CLINE.md, etc.) that are included in every agent invocation.
type MemoryFileProvider interface {
	// ConfigureMemoryFile configures the memory file for the tool
	ConfigureMemoryFile(projectPath string) error
	// IsMemoryFileConfigured checks if the memory file is already configured
	IsMemoryFileConfigured(projectPath string) bool
}

// SlashCommandProvider handles configuration of slash commands
// (like .claude/commands/spectr/*.md) that are invoked conditionally.
type SlashCommandProvider interface {
	// ConfigureSlashCommands configures the slash commands for the tool
	ConfigureSlashCommands(projectPath string) error
	// AreSlashCommandsConfigured checks if the slash commands are already
	// configured
	AreSlashCommandsConfigured(projectPath string) bool
}
