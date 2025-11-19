# Implementation Tasks

## 1. Create Tool Mapping System
- [ ] 1.1 Add mapping structure in `internal/init/registry.go` linking config tool IDs to slash tool IDs
- [ ] 1.2 Define `getSlashToolID(configToolID string) (string, bool)` helper function
- [ ] 1.3 Add mapping entries for all 11 tool pairs (claude-code→claude, cline-config→cline, etc.)

## 2. Update Tool Registry
- [ ] 2.1 Remove slash-only tool entries from `NewRegistry()` in `registry.go`
- [ ] 2.2 Keep config-based tool entries (claude-code, cline, cursor, etc.)
- [ ] 2.3 Verify tool priority values remain sequential after removals
- [ ] 2.4 Update tool count in comments/documentation if present

## 3. Modify Executor Configuration Flow
- [ ] 3.1 Update `configureTools()` in `executor.go` to check for slash command mappings
- [ ] 3.2 After configuring each tool, check if it has a slash command equivalent
- [ ] 3.3 If mapping exists, invoke the slash command configurator
- [ ] 3.4 Ensure both config and slash files are tracked in ExecutionResult

## 4. Handle Slash Command Configuration
- [ ] 4.1 Update `getConfigurator()` to support slash command tool IDs (or create separate method)
- [ ] 4.2 Ensure slash command configurators are instantiated correctly
- [ ] 4.3 Verify file creation works for both config files and slash commands
- [ ] 4.4 Ensure no file overwrites without user intent (respect existing files)

## 5. Update Result Tracking
- [ ] 5.1 Modify `ExecutionResult` to track slash command files separately if needed
- [ ] 5.2 Update completion screen to show both config and slash command files created
- [ ] 5.3 Ensure file counts are accurate in success messages

## 6. Write Tests
- [ ] 6.1 Unit test for tool mapping function (valid mappings, invalid inputs)
- [ ] 6.2 Integration test for executor with config-based tool selection
- [ ] 6.3 Verify slash commands are created when config tool is selected
- [ ] 6.4 Test that removed slash-only tools no longer appear in registry
- [ ] 6.5 Test backward compatibility - existing slash command files are preserved

## 7. Manual Testing
- [ ] 7.1 Run `spectr init` in test project and select `claude-code`
- [ ] 7.2 Verify both `CLAUDE.md` and `.claude/commands/spectr/*.md` are created
- [ ] 7.3 Test with multiple tools selected simultaneously
- [ ] 7.4 Test with existing slash command files (should not overwrite)
- [ ] 7.5 Verify wizard still displays correct tool count and navigation works
