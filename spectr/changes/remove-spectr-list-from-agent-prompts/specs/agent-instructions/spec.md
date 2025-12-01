## ADDED Requirements

### Requirement: Direct File Access for Agents
Agent prompts SHALL instruct AI assistants to use direct file and directory access methods (such as `ls spectr/changes/`, `ls spectr/specs/`, or file reads) instead of CLI commands like `spectr list` to discover changes and specifications.

#### Scenario: Agent discovering active changes
- **WHEN** an agent needs to find active changes in a project
- **THEN** the agent prompt SHALL instruct reading `spectr/changes/` directory directly
- **AND** SHALL NOT instruct running `spectr list`

#### Scenario: Agent discovering specifications
- **WHEN** an agent needs to find existing specifications in a project
- **THEN** the agent prompt SHALL instruct reading `spectr/specs/` directory directly
- **AND** SHALL NOT instruct running `spectr list --specs`

#### Scenario: Agent grounding proposal in current state
- **WHEN** an agent is creating a new change proposal
- **THEN** the agent prompt SHALL instruct reading `spectr/project.md` and exploring directories with `ls` or `rg`
- **AND** SHALL NOT require running `spectr list` commands

### Requirement: User Documentation Preserved
The `spectr list` command references SHALL remain in user-facing documentation since formatted CLI output benefits human users.

#### Scenario: User-facing documentation unchanged
- **WHEN** a user reads README.md or docs/ content
- **THEN** they SHALL still see `spectr list` command examples and documentation
- **AND** the CLI command behavior SHALL remain unchanged
