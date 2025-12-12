## 1. FileInitializer Interface & Implementations

- [ ] 1.1 Create `FileInitializer` interface in `internal/initialize/providers/initializer.go`
- [ ] 1.2 Create `InstructionFileInitializer` in `internal/initialize/providers/instruction_initializer.go`
- [ ] 1.3 Create `MarkdownSlashCommandInitializer` in `internal/initialize/providers/markdown_command_initializer.go`
- [ ] 1.4 Create `TOMLSlashCommandInitializer` in `internal/initialize/providers/toml_command_initializer.go`
- [ ] 1.5 Add unit tests for each initializer type (create, update, IsConfigured scenarios)

## 2. Helper Functions

- [ ] 2.1 Add `ConfigureInitializers()` helper function (fail-fast iteration)
- [ ] 2.2 Add `AreInitializersConfigured()` helper function
- [ ] 2.3 Add `GetInitializerPaths()` helper function (with deduplication)
- [ ] 2.4 Add unit tests for helper functions (including fail-fast behavior)

## 3. Provider Interface Update

- [ ] 3.1 Update `Provider` interface to 6 methods (ID, Name, Priority, Initializers, IsConfigured, GetFilePaths)
- [ ] 3.2 Remove `Configure()` method from Provider interface
- [ ] 3.3 Remove `HasConfigFile()` and `HasSlashCommands()` methods from Provider interface
- [ ] 3.4 Remove `ConfigFile()`, `GetProposalCommandPath()`, `GetApplyCommandPath()` methods
- [ ] 3.5 Remove `CommandFormat()` method and `CommandFormat` type
- [ ] 3.6 Remove `BaseProvider` struct entirely

## 4. Registry Updates

- [ ] 4.1 Remove `WithConfigFile()` filter function from registry
- [ ] 4.2 Remove `WithSlashCommands()` filter function from registry
- [ ] 4.3 Update any registry code that depended on removed Provider methods

## 5. Migrate All 23 Providers (Single Atomic Change)

- [ ] 5.1 Migrate `ClaudeProvider` (reference implementation with instruction file + markdown commands)
- [ ] 5.2 Migrate `GeminiProvider` (TOML commands only, no instruction file)
- [ ] 5.3 Migrate providers with instruction file + markdown commands (Cline, Cursor, Aider, Codebuddy, etc.)
- [ ] 5.4 Migrate providers with markdown commands only (Continue, etc.)
- [ ] 5.5 Migrate providers with global paths (Codex with ~/.codex/)
- [ ] 5.6 Ensure all providers implement 6 methods and use helper functions

## 6. Executor & Wizard Updates

- [ ] 6.1 Update executor to call `ConfigureInitializers(provider.Initializers(), projectPath, tm)`
- [ ] 6.2 Remove wizard filtering logic that used `HasConfigFile()`/`HasSlashCommands()`
- [ ] 6.3 Update wizard to show all providers equally (no capability-based filtering)
- [ ] 6.4 Update any other code that depended on removed Provider methods

## 7. Constants & Existing Helpers

- [ ] 7.1 Keep `StandardCommandPaths()` helper if still useful for initializer construction
- [ ] 7.2 Keep `StandardFrontmatter()` or convert to constants for initializer use
- [ ] 7.3 Keep `expandPath()` and `isGlobalPath()` helpers (used internally by initializers)
- [ ] 7.4 Remove unused helpers from old BaseProvider implementation

## 8. Test Updates

- [ ] 8.1 Update existing provider tests to work with new 6-method interface
- [ ] 8.2 Update executor tests to use `ConfigureInitializers()` pattern
- [ ] 8.3 Remove tests for removed methods (HasConfigFile, HasSlashCommands, etc.)
- [ ] 8.4 Add integration tests for full provider configuration flow

## 9. Validation & Manual Testing

- [ ] 9.1 Run `spectr validate refactor-provider-initializers --strict`
- [ ] 9.2 Run `go test ./internal/initialize/...`
- [ ] 9.3 Run `go build ./...` to ensure no compilation errors
- [ ] 9.4 Run `golangci-lint run` to check for lint issues
- [ ] 9.5 Manual test: `spectr init` with multiple providers on fresh project
- [ ] 9.6 Manual test: `spectr init` on project with existing configs (update scenario)
- [ ] 9.7 Manual test: Verify all generated files have correct content
