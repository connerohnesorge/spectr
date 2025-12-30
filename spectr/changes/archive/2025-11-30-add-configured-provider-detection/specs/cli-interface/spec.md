# Delta Specification

## ADDED Requirements

### Requirement: Configured Provider Detection in Init Wizard

The initialization wizard SHALL detect which AI tool providers are already
configured for the project and display this status in the tool selection screen.
Already-configured providers SHALL be visually distinguished and pre-selected by
default.

#### Scenario: Display configured indicator for already-configured providers

- **WHEN** user runs `spectr init` on a project with `CLAUDE.md` already present
- **AND** user reaches the tool selection screen
- **THEN** the Claude Code entry displays a "configured" indicator (e.g., dimmed
  text or badge)
- **AND** the indicator is visually distinct from the selection checkbox
- **AND** other unconfigured providers do NOT show the configured indicator

#### Scenario: Pre-select already-configured providers

- **WHEN** user runs `spectr init` on a project with some providers already
  configured
- **AND** user reaches the tool selection screen
- **THEN** already-configured providers have their checkboxes pre-selected
- **AND** users can deselect them if they don't want to update the configuration
- **AND** unconfigured providers remain unselected by default

#### Scenario: Help text explains configured indicator

- **WHEN** user is on the tool selection screen
- **THEN** the help text or screen description explains what the "configured"
  indicator means
- **AND** the explanation clarifies that selecting a configured provider will
  update its files

#### Scenario: No configured providers

- **WHEN** user runs `spectr init` on a fresh project with no providers
  configured
- **AND** user reaches the tool selection screen
- **THEN** no providers show the configured indicator
- **AND** no providers are pre-selected
- **AND** the screen functions as before this change

#### Scenario: All providers configured

- **WHEN** user runs `spectr init` on a project with all available providers
  configured
- **AND** user reaches the tool selection screen
- **THEN** all providers show the configured indicator
- **AND** all providers are pre-selected
- **AND** user can deselect providers they don't want to update

#### Scenario: Configured detection uses provider's IsConfigured method

- **WHEN** the wizard initializes
- **THEN** it calls `IsConfigured(projectPath)` on each provider
- **AND** the result is cached for the wizard session (not re-checked on each
  render)
- **AND** providers with global paths (like Codex) are correctly detected
