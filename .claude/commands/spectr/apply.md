---
name: Spectr: Apply
description: Implement an approved Spectr change and keep tasks in sync.
category: Spectr
tags: [spectr, apply]
---
<!-- spectr:START -->
# Guardrails
- Favor straightforward, minimal implementations first and add complexity only when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `spectr/AGENTS.md` (located inside the `spectr/` directory—run `ls spectr` or `spectr init` if you don't see it) if you need additional Spectr conventions or clarifications.

# Steps
Track these steps as TODOs and complete them one by one.
1. Run `spectr accept <id>` to convert `tasks.md` to `tasks.jsonc` format for stable task tracking.
2. Read `spectr/changes/<id>/proposal.md`, `design.md` (if present), and `tasks.jsonc` to confirm scope and acceptance criteria.
3. Work through tasks sequentially, keeping edits minimal and focused on the requested change. Before starting each task, mark it as `in_progress` in `tasks.jsonc`.
4. After completing and verifying each task, mark it as `completed` IMMEDIATELY in `tasks.jsonc`. Do NOT batch status updates--update each task individually as soon as it is verified.
5. Tasks have status values: `pending`, `in_progress`, `completed`. Transitions should be: `pending` -> `in_progress` (before starting) -> `completed` (immediately after verification).
6. Read `spectr/changes/` and `spectr/specs/` directories when additional context is required.

# Reference
- Read `spectr/changes/<id>/proposal.md` for proposal details.
- Read `spectr/changes/<id>/specs/<capability>/spec.md` for delta specs.

<!-- spectr:END -->
