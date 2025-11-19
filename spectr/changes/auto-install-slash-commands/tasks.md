# Implementation Tasks

## 1. Create Tool Mapping System
- [x] 1.1 Add mapping structure in `internal/init/registry.go` linking config tool IDs to slash tool IDs
- [x] 1.2 Define `GetSlashToolMapping(configToolID string) (string, bool)` helper function
- [x] 1.3 Add mapping entries for all 6 tool pairs (claude-code→claude, cline→cline-slash, etc.)

## 2. Update Tool Registry
- [x] 2.1 Remove slash-only tool entries from `NewRegistry()` in `registry.go`
- [x] 2.2 Keep config-based tool entries (claude-code, cline, costrict-config, qoder-config, codebuddy, qwen)
- [x] 2.3 Verify tool priority values remain sequential after removals
- [x] 2.4 Update tool count in comments/documentation (changed from 17 to 6)

## 3. Modify Executor Configuration Flow
- [x] 3.1 Update `configureTools()` in `executor.go` to check for slash command mappings
- [x] 3.2 After configuring each tool, check if it has a slash command equivalent
- [x] 3.3 If mapping exists, invoke the slash command configurator
- [x] 3.4 Ensure both config and slash files are tracked in ExecutionResult

## 4. Handle Slash Command Configuration
- [x] 4.1 Update `getConfigurator()` to support slash command tool IDs (added cline-slash, windsurf, codebuddy-slash, qwen-slash)
- [x] 4.2 Ensure slash command configurators are instantiated correctly
- [x] 4.3 Verify file creation works for both config files and slash commands
- [x] 4.4 Ensure no file overwrites without user intent (respect existing files)

## 5. Update Result Tracking
- [x] 5.1 Added `getSlashCommandFileInfo()` helper to track slash command files
- [x] 5.2 Update completion screen to show both config and slash command files created
- [x] 5.3 Ensure file counts are accurate in success messages

## 6. Write Tests
- [x] 6.1 Unit test for tool mapping function (valid mappings, invalid inputs)
- [x] 6.2 Integration test for executor with config-based tool selection
- [x] 6.3 Verify slash commands are created when config tool is selected
- [x] 6.4 Test that removed slash-only tools no longer appear in registry
- [x] 6.5 Test backward compatibility - existing slash command files are preserved

## 7. Manual Testing
- [x] 7.1 Run `spectr init` in test project and select `claude-code`
- [x] 7.2 Verify both `CLAUDE.md` and `.claude/commands/spectr/*.md` are created
- [x] 7.3 Test with multiple tools selected simultaneously
- [x] 7.4 Test with existing slash command files (should not overwrite)
- [x] 7.5 Verify wizard still displays correct tool count and navigation works
