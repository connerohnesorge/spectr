# Change: Replace full AGENTS.md content with short pointer in instruction files

## Why

Currently, `spectr init` injects the full ~450-line `AGENTS.md.tmpl` content into root-level instruction files (`CLAUDE.md`, `AGENTS.md` at project root, etc.) between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers. This duplicates the same content that already exists in `spectr/AGENTS.md`, wasting AI assistant context window tokens on every request.

## What Changes

- Add new template `instruction-pointer.md.tmpl` containing a short pointer (~15 lines) that directs AI assistants to read `spectr/AGENTS.md` when needed
- Add `RenderInstructionPointer()` method to `TemplateManager`
- Update `TemplateRenderer` interface to include new method
- Update `configureConfigFile()` in `provider.go` to use the pointer template instead of full AGENTS content

## Impact

- Affected specs: cli-interface
- Affected code:
  - `internal/init/templates/spectr/instruction-pointer.md.tmpl` (new)
  - `internal/init/templates.go`
  - `internal/init/providers/provider.go`
