# Implementation Tasks

## 1. Remove ConfigPath from ToolDefinition

- [x] 1.1 Remove `ConfigPath string` field from ToolDefinition struct in `internal/init/models.go`
- [x] 1.2 Remove ConfigPath assignments from all tool registrations in `internal/init/registry.go` (6 tools)
- [x] 1.3 Update `getToolFileInfo()` in `internal/init/executor.go` to return actual configurator file paths
- [x] 1.4 Remove ConfigPath assertions from `internal/init/registry_test.go`
- [x] 1.5 Run tests to verify no regressions: `go test ./internal/init/...`
- [x] 1.6 Run full test suite: `go test ./...`
- [x] 1.7 Verify the build compiles: `go build`
