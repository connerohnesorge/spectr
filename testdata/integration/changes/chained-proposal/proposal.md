---
id: chained-proposal
requires:
  - id: prerequisite-feature
    reason: "needs base functionality"
enables:
  - id: future-enhancement
    reason: "unlocks advanced features"
---

# Change: Example Chained Proposal

## Why

This is an example proposal demonstrating the chained proposals feature.
It declares dependencies on other proposals using YAML frontmatter.

## What Changes

- **ADDED**: Example feature with dependencies
- **MODIFIED**: Documentation for chained proposals

## Impact

- Affected specs: `example`
- Affected code: N/A (testdata only)
