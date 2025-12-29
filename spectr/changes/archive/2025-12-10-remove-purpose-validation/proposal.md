# Change: Remove Purpose Section Requirement

## Why

When `spectr pr archive` archives a change that introduces a new spec, the
generated skeleton includes a short placeholder Purpose (`"TODO: Add purpose
description"`). This causes `spectr validate --all --strict` to fail post-merge
because the Purpose is under 50 characters. This creates a workflow friction
where archiving a legitimate new-spec change produces a PR that fails CI
validation.

Removing the Purpose section requirement entirely eliminates this friction.
Purpose sections provide minimal value compared to well-written Requirements
with scenarios, which are the actual source of truth for behavior.

## What Changes

- Remove `## Purpose` section validation (both missing and length checks)
- Remove `## Purpose` from spec skeleton generation during archive
- Update project.md to remove Purpose length constraint
- Simplify spec file structure to just `# Title` + `## Requirements`

## Impact

- Affected specs: `validation`
- Affected code:
  - `internal/validation/spec_rules.go` - Remove Purpose checks
  - `internal/archive/spec_merger.go` - Remove Purpose from skeleton
  - `spectr/project.md` - Remove Purpose length constraint
  - Various test files that include Purpose in fixtures
