# Refactor Providers

## Goal
Redesign the provider system to be more Go idiomatic, separating data (provider definition) from behavior (file generation/configuration), making it easier to add new providers and capabilities.

## Problem
The current `Provider` interface mixes configuration data (ID, Name) with execution logic (`Configure`, `IsConfigured`). This "smart object" pattern makes it harder to:
1. Test provider definitions in isolation.
2. Add new capabilities (like MCP servers) without modifying the interface and all implementations.
3. Reason about what a provider *is* vs what it *does*.

## Solution
1. **Separation of Concerns**: Define `Provider` as a pure data interface (or struct) that describes *what* the provider needs (config files, commands).
2. **Centralized Execution**: Move the "how" (rendering templates, writing files) into a dedicated `Manager` or `Executor` service.
3. **Expandability**: Use optional interfaces or struct fields for additional capabilities, allowing providers to opt-in to features without breaking changes.

## Proposed Changes
- **Refactor `Provider` Interface**: Remove `Configure` and `IsConfigured` methods. Add methods/fields that return data structures describing the desired state.
- **Create `ProviderManager`**: Implement logic to take a `Provider` definition and apply it to the filesystem.
- **Update Existing Providers**: Convert Claude Code, etc., to the new format.

## Risks
- Breaking change for internal APIs (explicitly allowed).
- Need to ensure all existing provider features (custom frontmatter, etc.) are covered by the new data model.
