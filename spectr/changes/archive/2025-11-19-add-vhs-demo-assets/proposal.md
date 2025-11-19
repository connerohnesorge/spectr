# Change: Add VHS Demo Assets and GIFs

## Why

Users need visual demonstrations to quickly understand Spectr's workflow and capabilities. Currently, the README and docs rely entirely on text and code examples, which can be harder to grasp than seeing the actual CLI in action. VHS (Video-to-HTML-Snippet) tape files allow us to create reproducible, version-controlled terminal recordings that demonstrate key features and can be easily regenerated when the CLI changes.

## What Changes

- Add VHS `.tape` files to `assets/vhs/` directory demonstrating core workflows:
  - `init.tape`: Initializing a new Spectr project
  - `list.tape`: Listing changes and specs
  - `validate.tape`: Validating a change with errors and fixes
  - `archive.tape`: Archiving a completed change
  - `workflow.tape`: End-to-end workflow (create, validate, implement, archive)
- Generate GIF files from tapes and store them in `assets/gifs/`
- Update README.md to include demo GIFs in relevant sections (Quick Start, Command Reference)
- Add demo GIFs to docs site pages (`docs/src/content/docs/`)
- Add Makefile target or script to regenerate GIFs from tapes
- Document VHS tape creation and regeneration process in CONTRIBUTING.md or docs

## Impact

- **Affected specs**: `documentation` (enhanced with visual demonstrations)
- **Affected code**: None (documentation and assets only)
- **Affected files**:
  - `README.md` (add GIF embeds)
  - `assets/` (new subdirectories: `vhs/`, `gifs/`)
  - `docs/src/content/docs/*.md` (add GIF embeds)
  - `Makefile` or new script (add `make gifs` target)
  - Optional: `.github/workflows/` (CI to verify tapes still work)
- **Dependencies**: Requires VHS to be installed for developers who want to regenerate GIFs
  - VHS installation is optional (only needed for regenerating demos)
  - Generated GIFs are committed to repo so users don't need VHS
- **Breaking changes**: None
