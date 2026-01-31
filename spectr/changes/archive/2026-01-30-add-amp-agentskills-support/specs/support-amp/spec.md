# Amp Support Specification (Delta)

## MODIFIED Requirements

### Requirement: Amp Provider

The system SHALL provide an `AmpProvider` that generates Amp-compatible agent skills in `.agents/skills/`.

#### Scenario: Provider returns skill initializers

- **WHEN** `AmpProvider.Initializers(ctx, tm)` is called
- **THEN** it SHALL return initializers for:
  - `.agents/skills/spectr-proposal/` directory creation
  - `.agents/skills/spectr-apply/` directory creation
  - `.agents/skills/spectr-proposal/SKILL.md` file creation
  - `.agents/skills/spectr-apply/SKILL.md` file creation
  - `.agents/skills/spectr-accept-wo-spectr-bin/` embedded skill
  - `.agents/skills/spectr-validate-wo-spectr-bin/` embedded skill

#### Scenario: Skill directory structure

- **WHEN** Amp initializers execute
- **THEN** they SHALL create the following structure:

  ```text
  .agents/skills/
  ├── spectr-proposal/
  │   └── SKILL.md
  ├── spectr-apply/
  │   └── SKILL.md
  ├── spectr-accept-wo-spectr-bin/
  │   ├── SKILL.md
  │   └── scripts/accept.sh
  └── spectr-validate-wo-spectr-bin/
      ├── SKILL.md
      └── scripts/validate.sh
  ```

#### Scenario: Provider initialization is idempotent

- **WHEN** Amp initializers execute multiple times
- **THEN** they SHALL only create files that do not already exist
- **AND** existing files SHALL NOT be overwritten
- **AND** each execution SHALL produce the same final state
