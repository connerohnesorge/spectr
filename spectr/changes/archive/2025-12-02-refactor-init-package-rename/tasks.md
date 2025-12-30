# Tasks: Rename internal/init to internal/initialize

## Implementation Tasks

- [x] **Rename directory**: Move `internal/init/` to `internal/initialize/`
- [x] **Update package declarations**: Change `package init` to `package
  initialize` in all `.go` files in `internal/initialize/`
- [x] **Update internal self-references**: Update import path in `executor.go`
  and `wizard.go` for providers subpackage
- [x] **Update cmd/init.go imports**: Change import paths and remove unnecessary
  alias
- [x] **Verify build**: Run `go build ./...` to ensure compilation succeeds
- [x] **Run tests**: Run `go test ./...` to ensure all tests pass
- [x] **Lint check**: Run `golangci-lint run` to verify no linting issues

## Verification Checklist

- [x] No files reference old path `internal/init`
- [x] No `initpkg` alias needed in `cmd/init.go`
- [x] All tests pass
- [x] Build succeeds
