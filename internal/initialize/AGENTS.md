# Initialize Package

Project initialization wizard and AI tool templates. Interactive setup via TUI or --non-interactive mode.

## OVERVIEW
`spectr init` creates `spectr/` directory with templates and configuration. Supports 32+ AI coding assistant providers. Wizard guides users through tool selection, project conventions.

## STRUCTURE
```
internal/initialize/
├── executor.go          # Wizard orchestration and validation
├── filesystem.go       # Directory creation and file operations
├── constants.go          # Initialize constants and paths
├── providers/            # 32+ AI assistant templates (Claude, Cursor, etc.)
└── templates/            # Spectr template scaffolding
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Wizard orchestration | executor.go | Main flow control |
| Provider templates | providers/*.go | 32+ AI tool configs |
| Template scaffolding | templates/*.go | Spectr file templates |

## CONVENTIONS
- **Provider pattern**: Each provider has template file in `providers/` subdirectory
- **Template naming**: `templates/<type>.go` for spectr templates
- **Non-interactive**: Skip all prompts when `--non-interactive` flag set
- **Validation**: Validate project type, detect existing spectr/ directory

## UNIQUE PATTERNS
- **32+ providers**: Each AI tool (Claude, Cursor, Windsurf, etc.) has dedicated provider
- **Template-based initialization**: Uses Go templates for file scaffolding
- **Interactive wizard**: TUI-based selection flows for project types and tools

## ANTI-PATTERNS
- **DON'T skip validation**: Always check for existing spectr/ before init
- **NO hardcoded paths**: Use constants from constants.go
- **DON'T ignore non-interactive**: Respect --non-interactive flag

## KEY FUNCTIONS
- `Executor.Run() error` - Main initialization workflow
- `DetectProjectType() string` - Determine project type from directory structure
- `SelectProviders() []string` - Interactive tool selection
- `WriteProjectFiles() error` - Scaffold spectr/ structure from templates

## PROVIDER TEMPLATES
Supported providers (examples):
- claude - Anthropic Claude Code
- cursor - Cursor AI
- windsurf - Windsurf AI
- aider - Aider
- cline - Cline
- continue - Continue.dev
- codex - Codex
- qoder - Qoder
- antigravity - Antigravity
- costrict - CoStrict
- And many more...

## TEMPLATE STRUCTURE
Each provider template implements:
```go
type Provider interface {
    Name() string
    InitCommand(projectPath string) []string
    AGENTSPath() string
    ProjectMDPath() string
    // ...
}
```

## COMMON PATTERNS
```go
// Create template provider
data := struct {
    ProjectType string   // "cli", "lib", etc.
    Tools      []string // Selected AI tools
}
tmpl, err := template.New("provider").Parse(providerTemplate)
```

## NON-GOALS
- **NOT a generic scaffolding tool**: Focused on spectr-specific initialization
- **NOT cross-platform init**: Assumes Go project structure
- **Limited providers**: Only supports 32+ pre-configured AI tools

## FLOW
1. Detect project type (CLI, lib, web)
2. Select AI tools (interactive or via --tools flag)
3. Validate project structure (no existing spectr/ conflicts)
4. Create directory structure: `spectr/{specs,changes,changes/archive}`
5. Generate provider-specific AGENTS.md files
6. Create project.md with conventions
7. Copy AI tool-specific templates
