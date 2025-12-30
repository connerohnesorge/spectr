# Proposal: Rename `spectr pr new` to `spectr pr proposal`

## Summary

Rename the `spectr pr new` subcommand to `spectr pr proposal` to align with the
slash command naming convention (`/spectr:proposal`). This creates a consistent
vocabulary across the CLI and IDE integrations.

## Motivation

The slash commands in `.claude/commands/spectr/` use the term "proposal" for
creating change proposals:

- `/spectr:proposal` - Scaffold a new Spectr change
- `/spectr:apply` - Implement an approved change
- `/spectr:sync` - Detect spec drift

However, the CLI uses `spectr pr new` for creating PR-based proposals, which
creates cognitive dissonance:

- Users familiar with slash commands expect "proposal" terminology
- "new" is generic and doesn't convey the purpose (creating a proposal for
  review)
- "proposal" matches the workflow concept documented in `spectr/AGENTS.md`

## Proposed Change

Rename `spectr pr new` to `spectr pr proposal`:

```text
Before: spectr pr new <change-id>
After:  spectr pr proposal <change-id>
```text

The command's behavior remains identical:

- Creates a PR containing a Spectr change proposal for review
- Copies the change to an isolated git worktree without archiving
- Original change remains in `spectr/changes/<change-id>/`

## Impact

- **CLI breaking change**: Users must update scripts/documentation using `spectr
  pr new`
- **Internal refactor**: Rename `ModeNew` → `ModeProposal`, `PRNewCmd` →
  `PRProposalCmd`
- **No functional change**: All behavior remains identical

## Alternatives Considered

1. **Add alias `proposal` for `new`**: Would allow both names but adds confusion
2. **Keep `new`**: Inconsistent with slash command naming
3. **Rename slash command to `/spectr:new`**: "proposal" better describes the
  action
