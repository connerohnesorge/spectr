# Proposal: Rename internal/init to internal/initialize

## Summary

Rename the `internal/init` package to `internal/initialize` to avoid Go reserved
keyword conflict and simplify import aliasing in consumers.

## Motivation

The current package path `github.com/connerohnesorge/spectr/internal/init`
conflicts with Go's reserved `init` keyword, requiring awkward import aliases
like `initpkg` in consuming code:

```go
initpkg "github.com/connerohnesorge/spectr/internal/init"
```

Renaming to `initialize` allows clean imports:

```go
"github.com/connerohnesorge/spectr/internal/initialize"
```

## Scope

- Rename directory: `internal/init` -> `internal/initialize`
- Update package declaration in all files from `package init` to `package
  initialize`
- Update import paths in all consumers
- Simplify import aliases where they become unnecessary

## Affected Files

### Directory Rename

- `internal/init/` -> `internal/initialize/`
- `internal/init/providers/` -> `internal/initialize/providers/`

### Import Updates

Files importing the package:

1. `cmd/init.go` - Uses both `internal/init` and `internal/init/providers`
2. `internal/init/executor.go` - Self-references `internal/init/providers`
3. `internal/init/wizard.go` - Self-references `internal/init/providers`

## Risk Assessment

- **Low risk**: Pure refactoring with no behavioral changes
- **Automated verification**: `go build ./...` and `go test ./...` will catch
  any missed updates
- All changes are mechanical find-and-replace operations

## Success Criteria

1. `go build ./...` succeeds
2. `go test ./...` passes
3. No import aliases needed for the `initialize` package
4. No references to old `internal/init` path remain
