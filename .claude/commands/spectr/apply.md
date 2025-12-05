---
name: Spectr: Apply
description: Implement an approved Spectr change and keep tasks in sync.
category: Spectr
tags: [spectr, apply]
---
<!-- spectr:START -->
**Guardrails**
- Favor straightforward, minimal implementations first and add complexity only when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `spectr/AGENTS.md` (located inside the `spectr/` directory—run `ls spectr` or `spectr init` if you don't see it) if you need additional Spectr conventions or clarifications.

**Pre-flight Check**
Before starting implementation, verify the change has been accepted:
1. Check if `spectr/changes/<id>/tasks.json` exists.
2. If `tasks.json` is **missing**, run `spectr accept <change-id>` to convert `tasks.md` to `tasks.json` and mark the change as accepted.
3. If `tasks.json` **exists**, proceed with implementation.

Note: `tasks.json` is the source of truth after acceptance. It contains structured task data with `"completed"` fields to track progress.

**Steps**
Track these steps as TODOs and complete them one by one.
1. Read `spectr/changes/<id>/proposal.md`, `design.md` (if present), and `tasks.json` to confirm scope and acceptance criteria.
2. Work through tasks sequentially, keeping edits minimal and focused on the requested change.
3. Confirm completion before updating statuses—make sure every task is finished before marking it complete.
4. Update `tasks.json` after each task is done: set `"completed": true` for finished tasks. Edit the JSON file directly.
5. Read `spectr/changes/` and `spectr/specs/` directories when additional context is required.

**Reference**
- Read `spectr/changes/<id>/proposal.md` for proposal details.
- Read `spectr/changes/<id>/specs/<capability>/spec.md` for delta specs.
- `tasks.json` structure: each task has `"id"`, `"description"`, `"completed"` (boolean), and optionally `"subtasks"`.

<!-- spectr:END -->
