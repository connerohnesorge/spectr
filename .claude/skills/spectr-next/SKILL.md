---
name: spectr-next
description: Execute the next pending task from a change proposal
---

# Spectr Next Task Execution

## Guardrails

- Favor straightforward, minimal implementations first and add complexity only
  when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `spectr/AGENTS.md` and `spectr/project.md` (located inside
  the `spectr/` directory—run `ls spectr`) for project and
  Spectr conventions or clarifications.

## Steps

1. Discover the change proposal directory in `spectr/changes/`.
2. Parse `tasks.jsonc` to find the first task with `status: "pending"`.
3. Execute the task based on its description.
4. Update the task status: `pending` → `in_progress` → `completed`.
5. Report progress and suggest next steps.

## Reference

- Read `spectr/changes/<id>/proposal.md` for proposal details.
- Read `spectr/changes/<id>/tasks.jsonc` for task status and descriptions.
- Read `spectr/changes/<id>/specs/<capability>/spec.md` for delta specs.
