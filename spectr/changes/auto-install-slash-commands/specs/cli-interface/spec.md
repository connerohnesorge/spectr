## ADDED Requirements

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
- **AND** the mapping includes all 11 tools with slash command variants
- **AND** the map can be extended for new tools

## MODIFIED Requirements

### Requirement: Flat Tool List in Initialization Wizard

The initialization wizard SHALL present all AI tool options in a single unified flat list without visual grouping by tool type. Slash-only tool entries SHALL be removed from the registry as their functionality is now provided via automatic installation when the corresponding config-based tool is selected.

#### Scenario: Display only config-based tools in wizard

- **WHEN** user runs `spectr init` and reaches the tool selection screen
- **THEN** only config-based AI tools are displayed (e.g., `claude-code`, `cline`, `cursor`)
- **AND** slash-only tool entries (e.g., `claude`, `kilocode`) are not shown
- **AND** tools are sorted by priority
- **AND** no section headers (e.g., "Config-Based Tools", "Slash Command Tools") are shown
- **AND** each tool appears as a single checkbox item with its name

#### Scenario: Keyboard navigation across displayed tools

- **WHEN** user navigates with arrow keys (↑/↓)
- **THEN** the cursor moves through all displayed config-based tools sequentially
- **AND** navigation is continuous without group boundaries
- **AND** the first tool is selected by default on screen load

#### Scenario: Tool selection works uniformly

- **WHEN** user presses space to toggle any tool
- **THEN** the checkbox state changes (checked/unchecked)
- **AND** selection state is preserved when navigating
- **AND** both config file and slash commands will be installed when confirmed

#### Scenario: Bulk selection operations

- **WHEN** user presses 'a' to select all
- **THEN** all displayed config-based tools are checked
- **AND** WHEN user presses 'n' to select none
- **THEN** all tools are unchecked
- **AND** operations work across all displayed tools

#### Scenario: Help text clarity

- **WHEN** the tool selection screen is displayed
- **THEN** the help text shows keyboard controls (↑/↓, space, a, n, enter, q)
- **AND** the help text does NOT reference tool groupings or categories
- **AND** the screen title clearly indicates "Select AI Tools to Configure"

#### Scenario: Reduced tool count in wizard

- **WHEN** the wizard displays the tool list
- **THEN** fewer total tools are shown compared to the previous implementation
- **AND** the count reflects only config-based tools (not slash-only duplicates)
- **AND** navigation and selection work correctly with the reduced count
