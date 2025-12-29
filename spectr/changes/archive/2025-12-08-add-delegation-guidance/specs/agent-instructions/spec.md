# Delta Specification

## ADDED Requirements

### Requirement: Delegation Context for Subagents

When orchestrators delegate implementation tasks to subagents or when agents
complete tasks from a change proposal, the instruction pointer SHALL include
guidance to provide change directory paths so subagents can reference the
authoritative specification.

#### Scenario: Orchestrator delegating task to coder subagent

- **WHEN** an orchestrator delegates a task from an active change proposal to a
  coder subagent
- **THEN** the instruction pointer SHALL guide the orchestrator to include the
  path to `<changes-dir>/<id>/proposal.md`
- **AND** SHALL guide inclusion of `<changes-dir>/<id>/tasks.json` for task
  context
- **AND** SHALL guide inclusion of relevant delta spec paths
  `<changes-dir>/<id>/specs/<capability>/spec.md`

#### Scenario: Agent completing tasks from change proposal

- **WHEN** an agent is completing tasks defined in a change proposal
- **THEN** the instruction pointer SHALL instruct the agent to read the proposal
  and tasks files for authoritative context
- **AND** SHALL reference the change directory using template variables for
  dynamic paths
