## 1. Create Memory File Providers

- [x] 1.1 Add ClaudeMemoryFileProvider to internal/init/providers.go
- [x] 1.2 Add ClineMemoryFileProvider to internal/init/providers.go
- [x] 1.3 Add QoderMemoryFileProvider to internal/init/providers.go
- [x] 1.4 Add CodeBuddyMemoryFileProvider to internal/init/providers.go
- [x] 1.5 Add QwenMemoryFileProvider to internal/init/providers.go
- [x] 1.6 Add CostrictMemoryFileProvider to internal/init/providers.go
- [x] 1.7 Verify AgentsFileProvider exists and is focused on AGENTS.md only

## 2. Create SpectrAgentsUpdater Provider

- [x] 2.1 Add SpectrAgentsUpdater struct to internal/init/providers.go
- [x] 2.2 Implement ConfigureMemoryFile() to update spectr/AGENTS.md with markers
- [x] 2.3 Implement IsMemoryFileConfigured() to check for markers in spectr/AGENTS.md
- [x] 2.4 Verify it uses existing internal/init/templates/spectr/AGENTS.md.tmpl template

## 3. Rename Slash Command Providers

- [x] 3.1 Rename to ClaudeSlashCommandProvider in internal/init/configurator.go
- [x] 3.2 Rename to ClineSlashCommandProvider in internal/init/configurator.go
- [x] 3.3 Rename to CursorSlashCommandProvider in internal/init/configurator.go
- [x] 3.4 Rename to ContinueSlashCommandProvider in internal/init/configurator.go
- [x] 3.5 Rename to WindsurfSlashCommandProvider in internal/init/configurator.go
- [x] 3.6 Rename to AiderSlashCommandProvider in internal/init/configurator.go
- [x] 3.7 Rename to KilocodeSlashCommandProvider in internal/init/configurator.go
- [x] 3.8 Rename to QoderSlashCommandProvider in internal/init/configurator.go
- [x] 3.9 Rename to CostrictSlashCommandProvider in internal/init/configurator.go
- [x] 3.10 Rename to CopilotSlashCommandProvider in internal/init/configurator.go
- [x] 3.11 Rename to MentatSlashCommandProvider in internal/init/configurator.go
- [x] 3.12 Rename to TabnineSlashCommandProvider in internal/init/configurator.go
- [x] 3.13 Rename to SmolSlashCommandProvider in internal/init/configurator.go
- [x] 3.14 Rename to CodeBuddySlashCommandProvider in internal/init/configurator.go
- [x] 3.15 Rename to QwenSlashCommandProvider in internal/init/configurator.go
- [x] 3.16 Rename to AntigravitySlashCommandProvider in internal/init/configurator.go

## 4. Create Composite Tool Providers

- [x] 4.1 Create internal/init/composite_providers.go file
- [x] 4.2 Add ClaudeCodeToolProvider with embedded fields
- [x] 4.3 Add ClineToolProvider with embedded fields
- [x] 4.4 Add QoderToolProvider with embedded fields
- [x] 4.5 Add CodeBuddyToolProvider with embedded fields
- [x] 4.6 Add QwenToolProvider with embedded fields
- [x] 4.7 Add CostrictToolProvider with embedded fields
- [x] 4.8 Add AntigravityToolProvider with embedded fields
- [x] 4.9 Implement GetName() for each composite provider
- [x] 4.10 Implement GetMemoryFileProvider() for each composite provider
- [x] 4.11 Implement GetSlashCommandProvider() for each composite provider

## 5. Update Registry

- [x] 5.1 Update internal/init/registry.go to remove Configurator references
- [x] 5.2 Add ToolProvider references to ToolDefinition if needed
- [x] 5.3 Update tool registration to work with ToolProvider pattern
- [x] 5.4 Verify auto-installation mapping still works with new pattern

## 6. Update Executor

- [x] 6.1 Remove getConfigurator() method from internal/init/executor.go
- [x] 6.2 Add getToolProvider() method returning ToolProvider
- [x] 6.3 Update Configure workflow to call MemoryFileProvider.ConfigureMemoryFile()
- [x] 6.4 Update Configure workflow to call SlashCommandProvider.ConfigureSlashCommands()
- [x] 6.5 Update Configure workflow to handle SpectrAgentsUpdater for memory file tools
- [x] 6.6 Update file tracking to aggregate files from all providers
- [x] 6.7 Update error handling to collect errors from all provider invocations
- [x] 6.8 Verify auto-installation still works with new getToolProvider() method

## 7. Remove Legacy Code

- [x] 7.1 Delete Configurator interface from internal/init/interfaces.go
- [x] 7.2 Delete ClaudeCodeConfigurator from internal/init/configurator.go
- [x] 7.3 Delete ClineConfigurator from internal/init/configurator.go
- [x] 7.4 Delete CostrictConfigurator from internal/init/configurator.go
- [x] 7.5 Delete QoderConfigurator from internal/init/configurator.go
- [x] 7.6 Delete CodeBuddyConfigurator from internal/init/configurator.go
- [x] 7.7 Delete QwenConfigurator from internal/init/configurator.go
- [x] 7.8 Delete AntigravityConfigurator from internal/init/configurator.go
- [x] 7.9 Delete or clean up internal/init/configurator.go if only slash providers remain
- [x] 7.10 Remove any TODOs referencing legacy Configurator pattern

## 8. Update Tests

- [x] 8.1 Update internal/init/configurator_test.go to test new providers
- [x] 8.2 Update internal/init/executor_test.go to use ToolProvider pattern
- [x] 8.3 Update internal/init/registry_test.go if affected
- [x] 8.4 Add tests for SpectrAgentsUpdater behavior
- [x] 8.5 Add tests for composite provider composition
- [x] 8.6 Verify all existing tests pass with new architecture

## 9. Integration Testing

- [x] 9.1 Test init command with Claude Code tool
- [x] 9.2 Test init command with Cline tool
- [x] 9.3 Test init command with Antigravity tool
- [x] 9.4 Verify CLAUDE.md created correctly in repo root
- [x] 9.5 Verify spectr/AGENTS.md updated correctly by all memory file tools
- [x] 9.6 Verify .claude/commands/spectr/ created correctly
- [x] 9.7 Verify marker-based updates are idempotent
- [x] 9.8 Test auto-installation (config tool triggers slash commands)

## 10. Documentation

- [x] 10.1 Update internal/init package documentation
- [x] 10.2 Add code comments explaining provider composition pattern
- [x] 10.3 Add examples showing how to add new tool providers
- [x] 10.4 Document SpectrAgentsUpdater cross-file update behavior
