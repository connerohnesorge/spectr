## 1. Update Provider Interface and BaseProvider

- [ ] 1.1 Add `GetProposalCommandPath() string` to Provider interface
- [ ] 1.2 Add `GetArchiveCommandPath() string` to Provider interface
- [ ] 1.3 Add `GetApplyCommandPath() string` to Provider interface
- [ ] 1.4 Remove `SlashDir() string` from Provider interface
- [ ] 1.5 Replace `slashDir` field with `proposalPath`, `archivePath`, `applyPath` in BaseProvider
- [ ] 1.6 Implement the three new methods on BaseProvider
- [ ] 1.7 Update `HasSlashCommands()` to check any path is non-empty
- [ ] 1.8 Add `StandardCommandPaths(dir, ext string)` helper function

## 2. Update BaseProvider Internal Methods

- [ ] 2.1 Update `Configure()` to use new path methods
- [ ] 2.2 Update `configureSlashCommands()` to iterate over path methods
- [ ] 2.3 Update `configureSlashCommand()` to accept full path instead of building it
- [ ] 2.4 Remove `getSlashCommandPath()` private method
- [ ] 2.5 Update `IsConfigured()` to use new path methods
- [ ] 2.6 Update `GetFilePaths()` to use new path methods

## 3. Update Provider Implementations

- [ ] 3.1 Update ClaudeProvider with standard markdown paths
- [ ] 3.2 Update ClineProvider with standard markdown paths
- [ ] 3.3 Update CursorProvider with standard markdown paths
- [ ] 3.4 Update WindsurfProvider with standard markdown paths
- [ ] 3.5 Update AiderProvider with standard markdown paths
- [ ] 3.6 Update ContinueProvider with standard markdown paths
- [ ] 3.7 Update CodebuddyProvider with standard markdown paths
- [ ] 3.8 Update CoStrictProvider with standard markdown paths
- [ ] 3.9 Update KilocodeProvider with standard markdown paths
- [ ] 3.10 Update TabnineProvider with standard markdown paths
- [ ] 3.11 Update QwenProvider with standard markdown paths
- [ ] 3.12 Update QoderProvider with standard markdown paths
- [ ] 3.13 Update AntigravityProvider with standard markdown paths

## 4. Simplify GeminiProvider

- [ ] 4.1 Update GeminiProvider to use TOML paths directly via new fields
- [ ] 4.2 Remove overridden `Configure()` method (use BaseProvider)
- [ ] 4.3 Remove `configureSlashCommands()` private method
- [ ] 4.4 Remove `configureTOMLCommand()` private method
- [ ] 4.5 Remove `getTOMLCommandPath()` private method
- [ ] 4.6 Remove overridden `IsConfigured()` method
- [ ] 4.7 Remove overridden `GetFilePaths()` method
- [ ] 4.8 Keep `generateTOMLContent()` for TOML-specific content generation

## 5. Update Tests

- [ ] 5.1 Update `provider_test.go` to test new path methods
- [ ] 5.2 Remove tests for deprecated `SlashDir()` method
- [ ] 5.3 Update `registry_test.go` to verify path methods instead of SlashDir
- [ ] 5.4 Add test for `StandardCommandPaths()` helper
- [ ] 5.5 Verify all providers return valid paths in tests

## 6. Validation

- [ ] 6.1 Run `go build ./...` to verify compilation
- [ ] 6.2 Run `go test ./internal/init/providers/...` to verify tests pass
- [ ] 6.3 Run `golangci-lint run` to verify no linting errors
- [ ] 6.4 Run `spectr init` manually to verify end-to-end functionality
