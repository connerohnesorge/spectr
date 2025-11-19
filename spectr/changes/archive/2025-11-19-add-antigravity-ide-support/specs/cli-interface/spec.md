## MODIFIED Requirements

### Requirement: Automatic Slash Command Installation

When a config-based AI tool is selected during initialization, the system SHALL automatically install the corresponding slash command files for that tool without requiring separate user selection.

Config-based tools include those that create instruction files (e.g., `claude-code` creates `CLAUDE.md`). Slash command files are the workflow command files (e.g., `.claude/commands/spectr/proposal.md`).

This automatic installation provides users with complete Spectr integration in a single selection, eliminating the need for redundant tool entries in the wizard.

#### Scenario: Claude Code auto-installs slash commands

- **WHEN** user selects `claude-code` in the init wizard
- **THEN** the system creates `CLAUDE.md` in the project root
- **AND** the system creates `.claude/commands/spectr/proposal.md`
- **AND** the system creates `.claude/commands/spectr/apply.md`
- **AND** the system creates `.claude/commands/spectr/archive.md`
- **AND** all files are tracked in the execution result
- **AND** the completion screen shows all 4 files created

#### Scenario: Antigravity auto-installs slash commands

- **WHEN** user selects `antigravity` in the init wizard
- **THEN** the system creates `ANTIGRAVITY.md` in the project root
- **AND** the system creates `.antigravity/commands/spectr/proposal.md`
- **AND** the system creates `.antigravity/commands/spectr/apply.md`
- **AND** the system creates `.antigravity/commands/spectr/archive.md`
- **AND** all files are tracked in the execution result
- **AND** the completion screen shows all 4 files created

#### Scenario: Multiple tools with slash commands selected

- **WHEN** user selects both `claude-code` and `cursor` in the init wizard
- **THEN** the system creates `CLAUDE.md` and both config + slash commands for Claude
- **AND** the system creates `.cursor/commands/spectr-proposal.md` and slash commands for Cursor
- **AND** all files from both tools are created and tracked separately
- **AND** the completion screen lists all created files grouped by tool

#### Scenario: Slash command files already exist

- **WHEN** user runs init and selects `claude-code`
- **AND** `.claude/commands/spectr/proposal.md` already exists
- **THEN** the existing file's content between `<!-- spectr:START -->` and `<!-- spectr:END -->` is updated
- **AND** the file's YAML frontmatter is preserved
- **AND** no error occurs
- **AND** the file is marked as "updated" rather than "created" in execution result

#### Scenario: Config-based tool without slash mapping

- **WHEN** a config-based tool has no slash command equivalent in the mapping
- **THEN** only the config file is created
- **AND** no error occurs
- **AND** the system continues with remaining tool configurations

#### Scenario: Tool mapping is explicit and centralized

- **WHEN** a developer reviews the mapping logic
- **THEN** they find a `configToSlashMapping` map in `internal/init/registry.go`
- **AND** the map contains explicit entries for each tool pair
- **AND** the mapping includes Antigravity with slash command variants
- **AND** the map can be extended for new tools
