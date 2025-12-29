# Delta Specification

## ADDED Requirements

### Requirement: Partial Change ID Resolution for Archive Command

The `spectr archive` command SHALL support intelligent partial ID matching when
a non-exact change ID is provided as an argument. The resolution algorithm SHALL
prioritize prefix matches over substring matches and require a unique match to
proceed.

#### Scenario: Exact ID match takes precedence

- **WHEN** user runs `spectr archive add-feature`
- **AND** a change with ID `add-feature` exists
- **THEN** the archive proceeds with `add-feature`
- **AND** no resolution message is displayed

#### Scenario: Unique prefix match resolves successfully

- **WHEN** user runs `spectr archive refactor`
- **AND** only one change ID starts with `refactor` (e.g.,
  `refactor-unified-interactive-tui`)
- **THEN** a message is displayed: "Resolved 'refactor' ->
  'refactor-unified-interactive-tui'"
- **AND** the archive proceeds with the resolved ID

#### Scenario: Unique substring match resolves successfully

- **WHEN** user runs `spectr archive unified`
- **AND** no change ID starts with `unified`
- **AND** only one change ID contains `unified` (e.g.,
  `refactor-unified-interactive-tui`)
- **THEN** a message is displayed: "Resolved 'unified' ->
  'refactor-unified-interactive-tui'"
- **AND** the archive proceeds with the resolved ID

#### Scenario: Multiple prefix matches cause error

- **WHEN** user runs `spectr archive add`
- **AND** multiple change IDs start with `add` (e.g., `add-feature`,
  `add-hotkey`)
- **THEN** an error is displayed: "Ambiguous ID 'add' matches multiple changes:
  add-feature, add-hotkey"
- **AND** the command exits with error code 1
- **AND** no archive operation is performed

#### Scenario: Multiple substring matches cause error

- **WHEN** user runs `spectr archive search`
- **AND** no change ID starts with `search`
- **AND** multiple change IDs contain `search` (e.g., `add-search-hotkey`,
  `update-search-ui`)
- **THEN** an error is displayed: "Ambiguous ID 'search' matches multiple
  changes: add-search-hotkey, update-search-ui"
- **AND** the command exits with error code 1
- **AND** no archive operation is performed

#### Scenario: No match found

- **WHEN** user runs `spectr archive nonexistent`
- **AND** no change ID matches `nonexistent` (neither prefix nor substring)
- **THEN** an error is displayed: "No change found matching 'nonexistent'"
- **AND** the command exits with error code 1
- **AND** no archive operation is performed

#### Scenario: Case-insensitive matching

- **WHEN** user runs `spectr archive REFACTOR`
- **AND** a change ID `refactor-unified-interactive-tui` exists
- **THEN** the partial match succeeds (case-insensitive)
- **AND** the archive proceeds with the resolved ID

#### Scenario: Prefix match preferred over substring match

- **WHEN** user runs `spectr archive add`
- **AND** change ID `add-feature` exists (prefix match)
- **AND** change ID `update-add-button` exists (substring match only)
- **THEN** the prefix match `add-feature` is selected
- **AND** the substring-only match is ignored in preference calculation
