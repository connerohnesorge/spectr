# CLI Interface Delta Spec

## ADDED Requirements

### Requirement: Project Configuration File

The system SHALL support an optional `spectr.yaml` configuration file at the
project root.

#### Scenario: Config file present

- **WHEN** `spectr.yaml` exists in project root
- **THEN** load and parse the configuration

#### Scenario: Config file absent

- **WHEN** `spectr.yaml` does not exist
- **THEN** proceed with default behavior (no appended tasks)

#### Scenario: Config file malformed

- **WHEN** `spectr.yaml` contains invalid YAML
- **THEN** display error message and exit non-zero

### Requirement: Append Tasks Configuration

The system SHALL support an `append_tasks` section in `spectr.yaml` with a
configurable section name and list of tasks.

#### Scenario: Valid append_tasks configuration

- **WHEN** config contains `append_tasks.section` and `append_tasks.tasks`
- **THEN** parse section name as string and tasks as list of strings

#### Scenario: Missing section name

- **WHEN** `append_tasks.tasks` exists but `append_tasks.section` is missing
- **THEN** use default section name "Automated Tasks"

#### Scenario: Empty tasks list

- **WHEN** `append_tasks.tasks` is empty or missing
- **THEN** do not append any tasks

### Requirement: Auto-Append Tasks on Accept

The system SHALL append configured tasks to `tasks.jsonc` during
`spectr accept`.

#### Scenario: Append tasks with configured section

- **WHEN** `spectr accept <id>` runs with valid `append_tasks` config
- **THEN** append tasks to `tasks.jsonc` under the configured section name
- **AND** generate sequential task IDs continuing from the last task

#### Scenario: Task ID generation for appended tasks

- **WHEN** appending tasks after existing tasks (e.g., last ID was 3.2)
- **THEN** start appended tasks at next section number (e.g., 4.1, 4.2)

#### Scenario: No config present during accept

- **WHEN** `spectr accept <id>` runs without `spectr.yaml`
- **THEN** produce identical output to current behavior (no appended tasks)
