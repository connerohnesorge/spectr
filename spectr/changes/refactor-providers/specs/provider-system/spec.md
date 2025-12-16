# Provider System

## Purpose
Defines the architecture for AI tool integrations in a decoupled, data-driven way.

## Requirements

### Requirement: Provider Data Interface
Providers SHALL be defined as data structures or data-returning interfaces, not behavior-driven "smart objects".

#### Scenario: Metadata access
- **WHEN** accessing provider info
- **THEN** it returns ID, Name, and Priority without performing side effects or IO

### Requirement: Capabilities
Providers SHALL declare capabilities via data structs, enabling individual expandability.

#### Scenario: Instruction File Capability
- **WHEN** a provider needs an instruction file (e.g. `CLAUDE.md`)
- **THEN** it exposes an `InstructionFile` capability definition
- **AND** the definition includes the filename and template data/rendering logic

#### Scenario: Slash Command Capability
- **WHEN** a provider needs slash commands
- **THEN** it exposes a `SlashCommands` capability definition
- **AND** the definition includes a list of commands with paths and descriptions

### Requirement: Enactor Pattern
The application of provider configurations SHALL be handled by a central Enactor, not by the providers themselves.

#### Scenario: Configuration application
- **WHEN** `spectr init` runs
- **THEN** the Enactor iterates over the selected provider's capabilities
- **AND** performs the necessary filesystem operations (writing files, rendering templates)
