## Context

When multiple developers work on separate change proposals in parallel, they may unknowingly target the same specifications or requirements. This creates merge conflicts during archive and wastes implementation effort. A proactive daily check can surface these conflicts early.

## Goals / Non-Goals

**Goals:**
- Detect when two or more pending changes modify the same capability/spec
- Detect when two or more pending changes touch the same requirement (ADDED/MODIFIED/REMOVED/RENAMED)
- Automatically create GitHub issues to notify maintainers of conflicts
- Provide clear, actionable conflict reports

**Non-Goals:**
- Automatic conflict resolution (human decision required)
- Real-time conflict detection on push (use existing validation for that)
- Detecting spec-to-code drift (different feature)
- Detecting broken cross-references (different feature)

## Decisions

### Conflict Detection Algorithm

**Decision:** Use a two-level conflict detection approach:
1. **Capability-level conflicts:** Multiple changes modify the same capability (same `specs/<capability>/` directory)
2. **Requirement-level conflicts:** Multiple changes touch the same requirement name (across any delta operation)

**Rationale:** Capability-level catches broad overlaps; requirement-level catches specific semantic conflicts. Both are valuable signals.

### CLI Integration

**Decision:** Add a new `spectr conflicts` command rather than extending `spectr validate`.

**Rationale:** Conflict detection is conceptually different from validation. Validation checks a single change for correctness; conflict detection compares multiple changes against each other. Separate commands keep responsibilities clear.

### Issue Deduplication

**Decision:** Use issue labels and title patterns to prevent duplicate issues. Before creating a new issue, check for existing open issues with the same conflict signature.

**Rationale:** Daily runs would create duplicate issues without deduplication. Using GitHub's issue search API keeps the issue tracker clean.

**Alternatives considered:**
- State file in repo: Requires commits, clutters history
- GitHub Actions cache: Less visible, harder to debug
- Issue labels + search: Clean, uses native GitHub features

### Output Format

**Decision:** Support both human-readable and JSON output (`--json` flag).

**JSON schema:**
```json
{
  "conflicts": [
    {
      "type": "capability" | "requirement",
      "capability": "auth",
      "requirement": "Two-Factor Authentication",  // null for capability-level
      "changes": ["add-2fa-totp", "add-2fa-sms"],
      "operations": ["ADDED", "ADDED"]  // what each change does
    }
  ],
  "summary": {
    "total_conflicts": 2,
    "affected_capabilities": 1,
    "affected_changes": 2
  }
}
```

## Risks / Trade-offs

**Risk:** False positives when changes intentionally coordinate on same capability.
**Mitigation:** Issues are informational; teams can close if intentional. Clear issue title indicates it's for review, not a hard error.

**Risk:** Workflow fails to create issue due to permissions.
**Mitigation:** Use `GITHUB_TOKEN` with `issues: write` permission. Document required permissions.

**Trade-off:** Daily schedule may miss conflicts created and merged same day.
**Mitigation:** Existing CI validation catches validation errors; this is for early warning on parallel work.

## Open Questions

- Should we support manual trigger (workflow_dispatch) in addition to cron?
- Should issues be assigned to change authors automatically?
