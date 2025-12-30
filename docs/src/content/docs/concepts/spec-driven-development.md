---
title: Spec-Driven Development
description: Understanding the Spectr workflow
---

Spectr implements a **three-stage workflow** for managing changes:

## Stage 1: Creating Changes

Create a proposal when you need to:

- Add features or functionality
- Make breaking changes (API, schema)
- Change architecture or patterns
- Optimize performance (changes behavior)
- Update security patterns

**Skip proposals for:**

- Bug fixes (restore intended behavior)
- Typos, formatting, comments
- Dependency updates (non-breaking)
- Tests for existing behavior

## Stage 2: Implementing Changes

1. Read `proposal.md` - Understand what's being built
2. Read `design.md` (if exists) - Review technical decisions
3. Read `tasks.md` - Get implementation checklist
4. Implement tasks sequentially
5. Mark tasks complete with `- [x]` after implementation
6. **Approval gate**: Do not implement until proposal is approved

## Stage 3: Archiving Changes

After deployment:

1. Run `spectr validate <change>` to ensure quality
2. Run `spectr archive <change>` to merge deltas into specs
3. Changes move to `archive/YYYY-MM-DD-<change>/`
4. Specs in `specs/` are updated with merged requirements
