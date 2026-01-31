# Slash Commands: Timeline Addition

## ADDED Requirements

### Requirement: Timeline Slash Command

The system SHALL provide a `/spectr:timeline` slash command that analyzes all
active changes and generates a structured implementation timeline with
dependency analysis and phase-based ordering.

#### Scenario: Discover all active changes

- **WHEN** `/spectr:timeline` is executed
- **THEN** the command SHALL discover all directories in `spectr/changes/`
  (excluding `archive/` subdirectory)
- AND extract proposal metadata from each `proposal.md` file
- AND include change ID, title, description, and task counts

#### Scenario: Parse proposal frontmatter

- **WHEN** analyzing proposals with chained-proposal metadata
- **THEN** the command SHALL extract `requires` and `enables` relationships
- AND parse reason fields for each dependency
- AND include this metadata in the timeline

#### Scenario: Build and validate dependency graph

- **WHEN** constructing the dependency graph from all proposals
- **THEN** the command SHALL detect circular dependencies
- AND report them as errors with detailed cycle information
- AND fail gracefully without generating incomplete timeline

#### Scenario: Calculate implementation phases

- **WHEN** analyzing dependencies
- **THEN** the command SHALL identify changes that can run in parallel
- AND group them into phases (sequential batches)
- AND assign each change a phase number in the timeline

#### Scenario: Generate timeline.json

- **GIVEN** all changes analyzed and phases calculated
- **WHEN** generating output
- **THEN** the command SHALL create `./spectr/timeline.json` with structure:

```json
{
  "generated": "ISO-8601 timestamp",
  "summary": {
    "total_changes": number,
    "total_phases": number,
    "active_count": number,
    "archived_count": number
  },
  "dependency_graph": {
    "<change-id>": {
      "requires": [
        { "id": "<id>", "reason": "..." }
      ],
      "enables": [
        { "id": "<id>", "reason": "..." }
      ]
    }
  },
  "timeline": [
    {
      "phase": number,
      "parallel": true/false,
      "changes": [
        {
          "id": "<change-id>",
          "title": "...",
          "description": "...",
          "tasks": { "total": number, "completed": number },
          "blocked_by": ["<id1>", "<id2>"],
          "notes": "..."
        }
      ]
    }
  ],
  "notes": "implementation recommendations"
}
```

#### Scenario: Handle empty changes directory

- **WHEN** no active changes exist
- **THEN** the command SHALL generate timeline.json with empty timeline array
- AND set summary counts to zero
- AND include explanatory notes

#### Scenario: Include human-readable formatting

- **WHEN** generating timeline.json
- **THEN** JSON SHALL be pretty-printed with 2-space indentation
- AND include descriptive field names (no abbreviations)
- AND include comments/notes explaining blockers and parallelization
- AND be suitable for both human reading and machine parsing

### Requirement: Timeline Output Content

The timeline.json file SHALL contain comprehensive metadata for implementation
planning without duplicating spec capability information.

#### Scenario: Phase-based ordering

- **GIVEN** a set of interdependent changes
- **WHEN** computing phases
- **THEN** phase 0 SHALL contain only changes with no dependencies
- AND each subsequent phase SHALL contain changes whose requirements are all
  completed in earlier phases
- AND changes in the same phase can be implemented in parallel

#### Scenario: Dependency information

- **WHEN** including dependency metadata
- **THEN** each change SHALL list what blocks it (requires)
- AND what it enables (enables)
- AND each dependency SHALL include the reason from the proposal

#### Scenario: Task status reflection

- **WHEN** a change has associated tasks
- **THEN** timeline SHALL include task counts (total and completed)
- AND indicate completion percentage or status

#### Scenario: Implementation notes

- **GIVEN** the dependency graph and phases
- **WHEN** generating the timeline
- **THEN** it SHALL include notes on:
  - Changes that create critical path bottlenecks
  - Opportunities for parallelization
  - Risk factors (high dependencies, large scope)
  - Recommended implementation order

### Requirement: Skill Definition

The system SHALL define the `/spectr:timeline` skill with clear instructions
for AI agents to execute timeline generation.

#### Scenario: Skill metadata

- **WHEN** defining the skill
- **THEN** SKILL.md SHALL include name, description, and usage guidelines
- AND specify the command input (optional change-id filter)
- AND document expected output (timeline.json location and format)

#### Scenario: Step-by-step instructions

- **GIVEN** a skill definition
- **WHEN** an AI agent reads it
- **THEN** it SHALL provide clear steps:
  1. Discover active changes
  2. Parse proposal metadata
  3. Build dependency graph
  4. Calculate phases
  5. Generate timeline.json
  6. Validate output quality

#### Scenario: Error handling guidance

- **WHEN** defining the skill
- **THEN** it SHALL include guidance on handling:
  - Circular dependencies (error with details)
  - Malformed frontmatter (error with location)
  - Missing proposal files (warning, skip change)
  - Empty changes directory (success, empty timeline)
