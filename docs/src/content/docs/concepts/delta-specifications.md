---
title: Delta Specifications
description: How to define changes in Spectr
---

**Delta specs** describe proposed changes using operation headers:

```markdown
## ADDED Requirements
### Requirement: New Feature
The system SHALL provide new functionality.

#### Scenario: Success case
- **WHEN** condition occurs
- **THEN** expected result

## MODIFIED Requirements
### Requirement: Existing Feature
[Complete modified requirement with all scenarios]

## REMOVED Requirements
### Requirement: Deprecated Feature
**Reason**: Why removing
**Migration**: How to handle existing usage

## RENAMED Requirements
- FROM: `### Requirement: Old Name`
- TO: `### Requirement: New Name`
```

## Key Rules

- **ADDED**: New capabilities that stand alone
- **MODIFIED**: Changes to existing requirements (include FULL updated content)
- **REMOVED**: Deprecated features (provide reason and migration path)
- **RENAMED**: Name-only changes (use with MODIFIED if behavior changes too)
