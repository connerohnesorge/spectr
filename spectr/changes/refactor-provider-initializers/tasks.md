## 1. FileInitializer Interface & Implementations

- [ ] 1.1 Create `FileInitializer` interface in `internal/initialize/providers/initializer.go`
- [ ] 1.2 Create `InstructionFileInitializer` in `internal/initialize/providers/instruction_initializer.go`
- [ ] 1.3 Create `MarkdownSlashCommandInitializer` in `internal/initialize/providers/markdown_command_initializer.go`
- [ ] 1.4 Create `TOMLSlashCommandInitializer` in `internal/initialize/providers/toml_command_initializer.go`
- [ ] 1.5 Add unit tests for each initializer type

## 2. Helper Functions

- [ ] 2.1 Add `ConfigureInitializers()` helper function (fail-fast iteration)
- [ ] 2.2 Add `AreInitializersConfigured()` helper function
- [ ] 2.3 Add `GetInitializerPaths()` helper function (with deduplication)
- [ ] 2.4 Add unit tests for helper functions

## 3. Provider Interface Update

- [ ] 3.1 Update `Provider` interface to 6 methods (remove Configure, HasConfigFile, HasSlashCommands, ConfigFile, GetProposalCommandPath, GetApplyCommandPath, CommandFormat)
- [ ] 3.2 Remove `CommandFormat` type
- [ ] 3.3 Remove `BaseProvider` struct entirely
- [ ] 3.4 Update `Registry` if needed for new interface

## 4. Migrate All 23 Providers (Single Atomic Change)

- [ ] 4.1 Migrate `ClaudeProvider` (reference implementation)
- [ ] 4.2 Migrate `GeminiProvider` (TOML commands)
- [ ] 4.3 Migrate remaining 21 providers to new pattern
- [ ] 4.4 Ensure all providers use helper functions for IsConfigured/GetFilePaths

## 5. Executor & Wizard Updates

- [ ] 5.1 Update executor to call `ConfigureInitializers(provider.Initializers(), ...)` instead of `provider.Configure()`
- [ ] 5.2 Remove wizard filtering logic that used HasConfigFile/HasSlashCommands
- [ ] 5.3 Update any other code that depended on removed Provider methods

## 6. Cleanup & Validation

- [ ] 6.1 Remove old BaseProvider methods (configureConfigFile, configureSlashCommands, etc.)
- [ ] 6.2 Remove any unused helper functions from old implementation
- [ ] 6.3 Run `spectr validate refactor-provider-initializers --strict`
- [ ] 6.4 Run `go test ./internal/initialize/...`
- [ ] 6.5 Run `go build ./...` to ensure no compilation errors
- [ ] 6.6 Manual test: `spectr init` with multiple providers
- [ ] 6.7 Manual test: `spectr init` on project with existing configs (update scenario)
