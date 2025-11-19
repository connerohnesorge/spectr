---
title: Archiving Workflow
description: Learn how to archive completed changes and update specifications after deployment.
---

The archiving workflow is the final step in Spectr's change management process. After your implementation is deployed, you archive the change and update the master specifications.

## When to Archive

Archive a change when:

- ✅ Implementation is complete and tested
- ✅ All tasks in `tasks.md` are finished
- ✅ Changes are merged to main/master
- ✅ Features are deployed to production

## Archiving a Change

### 1. Prepare Your Implementation

Ensure all work is complete:

```bash
# Verify all tasks are done
cat spectr/changes/<change-id>/tasks.md

# Example - all tasks should be checked:
# - [x] 1.1 Task completed
# - [x] 1.2 Task completed
# - [x] 2.1 Task completed
```

### 2. Archive the Change

Use the archive command:

```bash
# Interactive mode - choose options
spectr archive <change-id>

# Non-interactive mode
spectr archive <change-id> --yes
```

For tooling-only changes (no spec updates needed):

```bash
spectr archive <change-id> --skip-specs --yes
```

### 3. What Gets Archived

The archiver automatically:

1. **Moves** `changes/<change-id>/` to `changes/archive/YYYY-MM-DD-<change-id>/`
2. **Merges** spec deltas into `specs/` (the source of truth)
3. **Updates** requirement content, adds new requirements, removes deprecated ones
4. **Validates** all changes

## Archive Workflow Steps

### Step 1: Review Your Change

Before archiving, verify:

```bash
# See what's in your change
spectr show <change-id> --json

# Validate the change passes strict checks
spectr validate <change-id> --strict
```

### Step 2: Run Archive

```bash
spectr archive add-two-factor-auth --yes
```

This:
- Moves the change directory to archive
- Processes all spec deltas
- Merges into `specs/` (master specs)
- Updates requirements with new/modified/removed items

### Step 3: Verify Results

Check that specs were updated correctly:

```bash
# View the updated spec
cat spectr/specs/auth/spec.md

# Validate all specs pass
spectr validate --strict
```

### Step 4: Commit Changes

Commit the archived change and updated specs:

```bash
git add spectr/
git commit -m "Archive: add-two-factor-auth

- Move change proposal to archive
- Update auth spec with 2FA requirements
- Update notifications spec with OTP delivery
"
```

## Understanding Spec Merging

When you archive, deltas are merged into master specs. Here's what happens:

### ADDED Requirements

```markdown
# Delta
## ADDED Requirements
### Requirement: Two-Factor Authentication
The system SHALL support OTP-based 2FA.
...

# Result
# specs/auth/spec.md now contains:
### Requirement: Two-Factor Authentication
The system SHALL support OTP-based 2FA.
...
```

### MODIFIED Requirements

```markdown
# Delta
## MODIFIED Requirements
### Requirement: User Authentication
Users SHALL provide credentials and OTP.
...

# Result
# The entire requirement is replaced with the modified version
```

### REMOVED Requirements

```markdown
# Delta
## REMOVED Requirements
### Requirement: Deprecated Legacy Login
**Reason**: Replaced by OAuth
**Migration**: Use OAuth flow instead

# Result
# The requirement is removed from specs/
```

## Example: Complete Archiving

### Before Archiving

```
spectr/
├── changes/
│   └── add-2fa/
│       ├── proposal.md
│       ├── tasks.md       # All [x] checked
│       └── specs/
│           ├── auth/spec.md
│           └── notifications/spec.md
└── specs/
    ├── auth/spec.md       # Current auth requirements
    └── notifications/spec.md
```

### Archive Command

```bash
spectr archive add-2fa --yes
```

### After Archiving

```
spectr/
├── changes/
│   └── archive/
│       └── 2025-01-15-add-2fa/
│           ├── proposal.md
│           ├── tasks.md
│           └── specs/
│               ├── auth/spec.md
│               └── notifications/spec.md
└── specs/
    ├── auth/spec.md       # Now includes 2FA requirements
    └── notifications/spec.md  # Now includes OTP notifications
```

## Troubleshooting Archive Issues

### Issue: Requirement not found when archiving

**Problem**: You modified a requirement that doesn't exist in the current spec.

**Solution**: Check the requirement name matches exactly:

```bash
# View current spec requirements
spectr show auth --json | jq '.requirements[].title'

# Verify your modified requirement name matches
```

### Issue: Validation fails after archive

**Problem**: The merged specs have validation errors.

**Solution**:

```bash
# See validation errors
spectr validate --strict

# Check the archived spec for issues
cat spectr/specs/auth/spec.md

# Fix manually if needed
```

### Issue: Archive didn't move directory

**Problem**: Change directory wasn't archived.

**Solution**: Check if archive succeeded:

```bash
# Should exist
ls spectr/changes/archive/2025-01-15-add-2fa/

# Should not exist
ls spectr/changes/add-2fa/  # Should fail
```

## Special Cases

### Tooling-Only Changes

For changes that don't affect specs (e.g., CI/CD, infrastructure):

```bash
spectr archive <change-id> --skip-specs --yes
```

This archives the change without attempting to merge specs.

### Multi-Spec Changes

If your change modified multiple specs, the archiver handles all of them:

```
spectr/changes/add-webhook-auth/
└── specs/
    ├── auth/spec.md           # Auth changes
    ├── api-design/spec.md     # API changes
    └── webhooks/spec.md       # Webhook changes
```

All three are merged into their respective `specs/` files automatically.

## Checklist: Before Archiving

- [ ] All tasks marked complete in `tasks.md`
- [ ] Implementation merged to main
- [ ] All tests passing
- [ ] No validation errors: `spectr validate <change-id> --strict`
- [ ] Changes deployed to production
- [ ] Team notified of new features

## Checklist: After Archiving

- [ ] Change moved to `changes/archive/`
- [ ] Specs updated in `specs/`
- [ ] No validation errors: `spectr validate --strict`
- [ ] All changes committed to git
- [ ] Team can reference new spec in future changes

## Next Steps

After archiving:

1. **Use updated specs** - Future changes reference the new requirements
2. **Plan next changes** - Create new change proposals as needed
3. **Monitor deployment** - Ensure features work as specified

## Key Concepts

### Why Archive?

- **Auditing**: Track what changes were made and when
- **Context**: See the reasoning in proposals
- **Evolution**: Understand how requirements evolved
- **Rollback**: Recover change details if needed

### Master Spec as Truth

After archiving, the `specs/` directory is the source of truth. Older specs are replaced with merged updates from the change.

## Further Reading

- [Creating Changes](./creating-changes)
- [Spec-Driven Development](/concepts/spec-driven-development)
- [CLI Commands](/reference/cli-commands)
