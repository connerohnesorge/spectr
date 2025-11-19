---
title: Creating Changes
description: Learn how to create and propose changes in Spectr with change proposals and spec deltas.
---

Creating a change proposal is the first step in the Spectr workflow. This guide walks you through proposing a new feature, breaking change, or architectural update.

## When to Create a Change

Create a change proposal when you need to:

- **Add features or functionality** - New capabilities your users need
- **Make breaking changes** - Updates to APIs, schemas, or behavior
- **Change architecture** - New patterns, services, or structural updates
- **Optimize performance** - Changes that affect how the system behaves
- **Update security patterns** - New security measures or requirements

### Skip the proposal for:

- Bug fixes that restore intended behavior
- Typos, formatting, or comment changes
- Non-breaking dependency updates
- Configuration adjustments
- Tests for existing behavior

## Creating a Change Proposal

### 1. Choose a Unique Change ID

Use a kebab-case, verb-led identifier:

- ✅ `add-two-factor-auth` - Verb-led, descriptive
- ✅ `update-api-versioning` - Verb-led, clear
- ✅ `refactor-payment-service` - Verb-led, focused
- ❌ `two-factor-auth` - Missing verb
- ❌ `twoFactorAuth` - Not kebab-case

### 2. Scaffold the Change Directory

Create the directory structure under `spectr/changes/<change-id>/`:

```bash
mkdir -p spectr/changes/add-two-factor-auth/specs/auth
```

This creates:
- `proposal.md` - Why and what
- `tasks.md` - Implementation checklist
- `design.md` (optional) - Technical decisions
- `specs/<capability>/spec.md` - Spec deltas

### 3. Write the Proposal

Create `proposal.md` explaining the change:

```markdown
# Change: Two-Factor Authentication Support

## Why
Users require additional security beyond passwords. Two-factor authentication reduces unauthorized access risk.

## What Changes
- **BREAKING** - Auth API requires second factor during login
- New OTP delivery mechanism
- Updated session token format
- New user settings for 2FA preferences

## Impact
- Affected specs: `auth`, `notifications`, `user-settings`
- Affected code: `services/auth`, `api/auth`, `components/LoginFlow`
```

### 4. Create Implementation Checklist

Create `tasks.md` with implementation steps:

```markdown
## 1. API & Core Logic
- [ ] 1.1 Add OTP generation logic
- [ ] 1.2 Update authentication endpoints
- [ ] 1.3 Create OTP storage schema
- [ ] 1.4 Add session validation for 2FA

## 2. User Interface
- [ ] 2.1 Create OTP input component
- [ ] 2.2 Update login flow
- [ ] 2.3 Add 2FA settings page
- [ ] 2.4 Add OTP delivery options

## 3. Testing & Deployment
- [ ] 3.1 Write unit tests
- [ ] 3.2 Write integration tests
- [ ] 3.3 Create migration script
- [ ] 3.4 Update documentation
```

### 5. Write Spec Deltas

Create specification changes under `specs/<capability>/spec.md`:

```markdown
## ADDED Requirements

### Requirement: OTP Generation
The system SHALL generate a time-based one-time password during login.

#### Scenario: OTP successfully generated
- **WHEN** user submits valid credentials
- **THEN** a 6-digit OTP is generated and delivered
- **AND** a 5-minute expiration is set

#### Scenario: OTP validation required
- **WHEN** user receives an OTP
- **THEN** they must submit it before session creation
- **AND** invalid OTPs are rejected

## MODIFIED Requirements

### Requirement: User Authentication
Users SHALL provide credentials and a second factor to authenticate.

#### Scenario: Password and OTP required
- **WHEN** user submits credentials
- **THEN** the system verifies credentials
- **AND** requests an OTP
- **AND** validates the OTP before granting access
```

**Important**: Every requirement must have at least one scenario using the format:

```markdown
#### Scenario: Descriptive scenario name
- **WHEN** condition
- **THEN** result
```

### 6. Optional: Create Design Document

Create `design.md` if your change involves:

- Cross-cutting changes (multiple services/modules)
- New architectural patterns
- External dependencies or data model changes
- Security or performance complexity
- Ambiguities that need technical decisions

```markdown
## Context
OTP delivery requires integration with notification service.

## Goals
- Provide flexible OTP delivery (email, SMS, authenticator)
- Maintain user security during delivery
- Support recovery codes for account recovery

## Decisions
- **OTP standard**: Time-based (TOTP) for compatibility
- **Delivery**: Via existing notification service
- **Recovery**: Support 8 backup codes per user

## Risks
- Network delays in OTP delivery
- **Mitigation**: Increased timeout window, retry mechanism
```

### 7. Validate Your Change

Before requesting approval, validate your change:

```bash
spectr validate <change-id> --strict
```

This checks:
- All scenarios properly formatted (#### Scenario:)
- Every requirement has at least one scenario
- Spec deltas follow conventions
- Required files exist

Fix any issues and validate again.

## Multi-Capability Changes

For changes affecting multiple capabilities, create separate spec files:

```
spectr/changes/add-2fa/
├── proposal.md
├── tasks.md
├── specs/
│   ├── auth/spec.md          # Auth requirements
│   ├── notifications/spec.md # Notification requirements
│   └── user-settings/spec.md # Settings requirements
```

## Common Patterns

### Feature Addition

```markdown
## ADDED Requirements
### Requirement: New Feature Name
[Description of the feature]

#### Scenario: Success case
- **WHEN** user performs action
- **THEN** expected result
```

### Behavior Change

```markdown
## MODIFIED Requirements
### Requirement: Existing Feature Name
[Updated requirement with new behavior]

#### Scenario: Updated behavior
- **WHEN** condition
- **THEN** new expected result
```

### Removing Functionality

```markdown
## REMOVED Requirements
### Requirement: Deprecated Feature
**Reason**: No longer needed, replaced by XYZ
**Migration**: Users should use ABC instead
```

## Next Steps

Once your proposal is created and validated:

1. Share the proposal with your team
2. Get approval before implementation
3. Move to the [Archiving Workflow](./archiving-workflow) when complete

## Tips for Good Change Proposals

- **Keep it focused** - One capability per change when possible
- **Write clear scenarios** - Every requirement needs concrete examples
- **Be specific** - "Adds login" is vague; "Adds OTP-based 2FA" is clear
- **Think about impact** - Explicitly list affected systems
- **Plan for testing** - Use tasks.md to capture test requirements

## Further Reading

- [Spec-Driven Development](/concepts/spec-driven-development)
- [Delta Specifications](/concepts/delta-specifications)
- [CLI Commands](/reference/cli-commands)
