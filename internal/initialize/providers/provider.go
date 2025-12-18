// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface and registration system that all
// AI CLI tools (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider handles both its instruction file (e.g., CLAUDE.md) and slash
// commands (e.g., .claude/commands/) through composable Initializers.
//
// # Adding a New Provider
//
// To add a new AI CLI provider, create a new file (e.g., providers/mytools.go):
//
// Example:
//
//	package providers
//
//	func init() {
//	    err := RegisterV2(Registration{
//	        ID:       "mytool",
//	        Name:     "MyTool",
//	        Priority: 100,
//	        Provider: &MyToolProvider{},
//	    })
//	    if err != nil {
//	        panic(err)
//	    }
//	}
//
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) Initializers(ctx context.Context) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".mytool/commands/spectr"),
//	        NewConfigFileInitializer("MYTOOL.md"),
//	        NewSlashCommandsInitializerWithFrontmatter(
//	            ".mytool/commands/spectr",
//	            ".md",
//	            FormatMarkdown,
//	            map[string]string{
//	                "proposal": FrontmatterProposal,
//	                "apply":    FrontmatterApply,
//	            },
//	        ),
//	    }
//	}
//
// Each provider returns a list of Initializers that handle directory creation,
// config file management, and slash command generation.
package providers

// CommandFormat specifies the format for slash command files.
type CommandFormat int

const (
	// FormatMarkdown uses markdown files with
	// YAML frontmatter (Claude, Cline, etc.)
	FormatMarkdown CommandFormat = iota
	// FormatTOML uses TOML files (Gemini CLI)
	FormatTOML
)
