## ADDED Requirements

### Requirement: PR Branch Naming Convention
The system SHALL use a mode-specific branch naming convention for PR branches that distinguishes between archive and proposal branches based on the subcommand used.

#### Scenario: Archive branch name format
- **WHEN** user runs `spectr pr archive <change-id>`
- **THEN** the branch is named `spectr/archive/<change-id>`

#### Scenario: Proposal branch name format
- **WHEN** user runs `spectr pr new <change-id>`
- **THEN** the branch is named `spectr/proposal/<change-id>`

#### Scenario: Branch name with special characters
- **WHEN** change ID contains only valid kebab-case characters
- **THEN** the branch name is valid for git

#### Scenario: Branch names clearly indicate PR purpose
- **WHEN** a developer views the branch list
- **THEN** they can distinguish archive PRs from proposal PRs by the branch prefix
- **AND** `spectr/archive/*` indicates a completed change being archived
- **AND** `spectr/proposal/*` indicates a change proposal for review

#### Scenario: Force flag for existing archive branch
- **WHEN** user runs `spectr pr archive <change-id> --force`
- **AND** branch `spectr/archive/<change-id>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Force flag for existing proposal branch
- **WHEN** user runs `spectr pr new <change-id> --force`
- **AND** branch `spectr/proposal/<change-id>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Archive branch conflict without force
- **WHEN** user runs `spectr pr archive <change-id>`
- **AND** branch `spectr/archive/<change-id>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "branch 'spectr/archive/<change-id>' already exists on remote; use --force to delete"
- **AND** the command exits with code 1

#### Scenario: Proposal branch conflict without force
- **WHEN** user runs `spectr pr new <change-id>`
- **AND** branch `spectr/proposal/<change-id>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "branch 'spectr/proposal/<change-id>' already exists on remote; use --force to delete"
- **AND** the command exits with code 1
