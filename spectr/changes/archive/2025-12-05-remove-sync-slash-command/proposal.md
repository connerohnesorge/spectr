# Change: Remove Sync Slash Command

## Why

The `/spectr:sync` slash command is being removed because:

1. The sync workflow (detecting spec drift from code) is conceptually complex
  and was never fully implemented
2. The agent-based approach adds overhead without clear user adoption
3. Users who need to sync specs with code can do so manually or through simpler
  mechanisms

## What Changes

- **BREAKING**: Remove `/spectr:sync` slash command from all providers
- Remove sync command template file (`slash-sync.md.tmpl`)
- Remove `GetSyncCommandPath()` method from provider interface
- Remove `FrontmatterSync` constant
- Remove sync-related paths from provider configurations
- Remove "Stage 3: Syncing Specs" section from AGENTS.md template
- Update spec scenarios that reference sync command creation
- Remove `WHY_SYNC.md` documentation file

## Impact

- Affected specs:
  - `cli-framework` - Provider interface
  - `cli-interface` - Init command outputs
  - All `support-*` specs (13 providers)
- Affected code:
  - `internal/initialize/providers/` - All provider files
  - `internal/initialize/templates/tools/slash-sync.md.tmpl`
  - `internal/initialize/templates/spectr/AGENTS.md.tmpl`
  - `.claude/commands/spectr/sync.md`
  - `WHY_SYNC.md`
