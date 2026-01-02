# CLI Commands

Thin layer using Kong framework. Entry points to internal/ packages.

## OVERVIEW
CLI commands in cmd/ parse flags, call internal/ business logic, exit with status codes. No business logic here.

## STRUCTURE
```go

cmd/
├── root.go              # Kong CLI struct with all commands
├── init.go              # spectr init
├── list.go              # spectr list
├── validate.go          # spectr validate
├── accept.go            # spectr accept
├── pr.go                # spectr pr archive|new
├── view.go              # spectr view
├── version.go           # spectr version
└── completion.go        # Shell completions
```

## WHERE TO LOOK

| Command | Handler | Internal Package |
|---------|----------|-----------------|
| spectr init | InitCmd.Run() | internal/initialize |
| spectr list | ListCmd.Run() | internal/list |
| spectr validate | ValidateCmd.Run() | internal/validation |
| spectr accept | AcceptCmd.Run() | internal/parsers + internal/discovery |
| spectr archive | ArchiveCmd.Run() | internal/archive |
| spectr pr | PRCmd.Run() | internal/pr |
| spectr view | ViewCmd.Run() | internal/view |

## CONVENTIONS
- **Thin layer**: Delegates to internal/, minimal logic in cmd/
- **Kong tags**: Use `cmd:`, `help:`, `aliases:` struct tags
- **Exit codes**: 0=success, non-zero=error
- **Context**: Pass kong.Context through for flag access

## UNIQUE PATTERNS
- **Kong integration**: root.go defines CLI struct, framework handles parsing/completion
- **Command grouping**: PR commands use internal/pr/ArchiveCmd embedded in PRCmd

## ANTI-PATTERNS
- **NO business logic in cmd/**: Delegate to internal/
- **DON'T bypass Kong**: Use Kong tags, not manual flag parsing
