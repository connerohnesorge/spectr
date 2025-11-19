# Change: Add Antigravity IDE Support

## Why
Antigravity IDE is an existing AI-powered IDE that should be supported by Spectr alongside Claude Code, Cline, Cursor, Qoder, CodeBuddy, and Qwen. This enables Antigravity users to use Spectr for spec-driven development workflows within their IDE environment through auto-installed slash commands and configuration injection.

## What Changes
- Add Antigravity to the tool provider registry as the 7th supported tool
- Create AntigravityConfigurator to handle initialization and AGENTS.md injection
- Update tool mapping to connect config-based tool to slash command equivalents
- Maintain consistency with existing 6-tool provider pattern

## Impact
- **Affected specs:**
  - `cli-interface` - Modified to add Antigravity tool registration
- **Affected code:**
  - `internal/init/registry.go` - Add tool definition and mapping
  - `internal/init/configurator.go` - Create AntigravityConfigurator
  - `internal/init/wizard.go` - Tool selection includes Antigravity (no changes needed, auto-populated from registry)
  - Tests in `internal/init/registry_test.go` and `internal/init/configurator_test.go`
- **User experience:**
  - Users selecting Antigravity during `spectr init` get AGENTS.md injected with spectr instructions
  - Slash commands auto-installed for `/spectr-proposal`, `/spectr-apply`, `/spectr-archive`
  - No changes to user workflows or slash command behavior
