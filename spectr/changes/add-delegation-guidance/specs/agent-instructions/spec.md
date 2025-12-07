## ADDED Requirements

### Requirement: Delegation Path References
When an AI orchestrator delegates implementation tasks to subagents, the instruction pointer SHALL include guidance to provide explicit paths to change proposal resources (`{{ .ChangesDir }}/<change-id>/proposal.md`, `tasks.json`, and delta specs) so subagents can reference the authoritative specification.

#### Scenario: Orchestrator delegates coding task
- **WHEN** an orchestrator agent delegates a task to a coder subagent
- **THEN** the instruction pointer content SHALL advise including the path to `{{ .ChangesDir }}/<change-id>/` directory
- **AND** the subagent SHALL be able to read proposal.md, tasks.json, and spec deltas for context

#### Scenario: Subagent needs implementation context
- **WHEN** a subagent begins work on a delegated task
- **THEN** the subagent SHALL have access to the change proposal path
- **AND** the subagent SHALL reference the spec deltas under `{{ .ChangesDir }}/<change-id>/specs/`
