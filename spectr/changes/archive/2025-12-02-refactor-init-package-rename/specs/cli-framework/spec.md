# Cli Framework Specification - Delta

## Purpose

Updates package path references from `internal/init/` to `internal/initialize/` to avoid Go reserved keyword conflict.

## MODIFIED Requirements

### Requirement: Per-Provider File Organization

The init system SHALL organize provider implementations as separate Go files under `internal/initialize/providers/`, with one file per provider.

#### Scenario: Provider file structure

- **WHEN** a provider file is created
- **THEN** it SHALL be named `{provider-id}.go` (e.g., `claude.go`, `gemini.go`)
- **AND** it SHALL contain an `init()` function that registers its provider
- **AND** it SHALL be self-contained with all provider-specific configuration

#### Scenario: Adding a new provider

- **WHEN** a developer adds a new AI CLI provider
- **THEN** they SHALL create a single file under `internal/initialize/providers/`
- **AND** the file SHALL implement the `Provider` interface
- **AND** the file SHALL call `Register()` in its `init()` function
- **AND** no other files SHALL require modification
