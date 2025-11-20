# Proposal: Refactor internal/init Package for Maintainability and Type Safety

## Why

The `internal/init` package has accumulated significant technical debt with ~1,400 lines of code containing massive duplication and poor separation of concerns. The current implementation has 18 nearly identical configurator structs, 15 duplicated slash command factory functions, and multiple giant switch statements that make adding new tools error-prone and tedious. This refactoring will reduce code by ~60%, improve type safety, and make adding new AI tool integrations a simple data declaration instead of writing hundreds of lines of boilerplate.

## What Changes

- **BREAKING**: Replace string-based tool IDs with type-safe constants
- Replace 18 identical configurator implementations with single data-driven configurator using tool config registry
- Replace 15 slash command factory functions with single factory using configuration data
- Consolidate 3 giant switch statements (getConfigurator, getToolFileInfo, configToSlashMapping) into registry-based lookups
- Extract marker-based file update logic into dedicated utility
- Unify template rendering patterns to always use TemplateManager
- Consolidate duplicate constants (filePerm/filePerms, dirPerm/dirPerms)
- Reorganize code into focused files: tool_definitions.go, file_operations.go, marker_utils.go

## Impact

- **Affected Specs**: cli-interface (initialization wizard and tool configuration)
- **Affected Code**:
  - internal/init/configurator.go (875 lines → ~150 lines)
  - internal/init/registry.go (146 lines → ~300 lines, includes tool definitions)
  - internal/init/executor.go (510 lines → ~350 lines)
  - internal/init/templates.go (106 lines → ~80 lines)
  - internal/init/constants.go (new file for consolidated constants)
  - internal/init/filesystem.go (minor cleanup)
- **Test Impact**: All existing tests should continue passing; refactor preserves behavior
- **User Impact**: No changes to CLI interface or user experience
- **Future Benefit**: Adding new AI tool support reduces from ~100 lines to ~10 lines of config
