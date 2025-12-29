# Implementation Tasks

## 1. Update Provider Interface and BaseProvider

- [x] 1.1 Add `GetProposalCommandPath() string` to Provider interface
- [x] 1.2 Add `GetArchiveCommandPath() string` to Provider interface
- [x] 1.3 Add `GetApplyCommandPath() string` to Provider interface
- [x] 1.4 Remove `SlashDir() string` from Provider interface
- [x] 1.5 Replace `slashDir` field with `proposalPath`, `archivePath`,
  `applyPath` in BaseProvider
- [x] 1.6 Implement the three new methods on BaseProvider
- [x] 1.7 Update `HasSlashCommands()` to check any path is non-empty
- [x] 1.8 Add `StandardCommandPaths(dir, ext string)` helper function

## 2. Update BaseProvider Internal Methods

- [x] 2.1 Update `Configure()` to use new path methods
- [x] 2.2 Update `configureSlashCommands()` to iterate over path methods
- [x] 2.3 Update `configureSlashCommand()` to accept full path instead of
  building it
- [x] 2.4 Remove `getSlashCommandPath()` private method
- [x] 2.5 Update `IsConfigured()` to use new path methods
- [x] 2.6 Update `GetFilePaths()` to use new path methods

## 3. Update Provider Implementations

- [x] 3.1 Update ClaudeProvider with standard markdown paths
- [x] 3.2 Update ClineProvider with standard markdown paths
- [x] 3.3 Update CursorProvider with standard markdown paths
- [x] 3.4 Update WindsurfProvider with standard markdown paths
- [x] 3.5 Update AiderProvider with standard markdown paths
- [x] 3.6 Update ContinueProvider with standard markdown paths
- [x] 3.7 Update CodebuddyProvider with standard markdown paths
- [x] 3.8 Update CoStrictProvider with standard markdown paths
- [x] 3.9 Update KilocodeProvider with standard markdown paths
- [x] 3.10 Update TabnineProvider with standard markdown paths
- [x] 3.11 Update QwenProvider with standard markdown paths
- [x] 3.12 Update QoderProvider with standard markdown paths
- [x] 3.13 Update AntigravityProvider with standard markdown paths

## 4. Simplify GeminiProvider

- [x] 4.1 Update GeminiProvider to use TOML paths directly via new fields
- [x] 4.2 Remove overridden `Configure()` method (use BaseProvider) - KEPT:
  needed for TOML-specific logic
- [x] 4.3 Remove `configureSlashCommands()` private method - KEPT: needed for
  TOML-specific logic
- [x] 4.4 Remove `configureTOMLCommand()` private method - KEPT: needed for
  TOML-specific logic
- [x] 4.5 Remove `getTOMLCommandPath()` private method - SIMPLIFIED: uses path
  fields directly
- [x] 4.6 Remove overridden `IsConfigured()` method - DONE
- [x] 4.7 Remove overridden `GetFilePaths()` method - DONE
- [x] 4.8 Keep `generateTOMLContent()` for TOML-specific content generation

## 5. Update Tests

- [x] 5.1 Update `provider_test.go` to test new path methods
- [x] 5.2 Remove tests for deprecated `SlashDir()` method
- [x] 5.3 Update `registry_test.go` to verify path methods instead of SlashDir
- [x] 5.4 Add test for `StandardCommandPaths()` helper
- [x] 5.5 Verify all providers return valid paths in tests

## 6. Validation

- [x] 6.1 Run `go build ./...` to verify compilation
- [x] 6.2 Run `go test ./internal/init/providers/...` to verify tests pass
- [x] 6.3 Run `golangci-lint run` to verify no linting errors
- [x] 6.4 Run `spectr init` manually to verify end-to-end functionality
