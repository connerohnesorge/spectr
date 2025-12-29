# Change: Move assets/ directory into docs/

## Why

The `assets/` directory at the repository root contains media files (logo, GIFs,
VHS tapes) that are primarily used by documentation. Moving them into
`docs/src/assets/` consolidates documentation-related files and aligns asset
location with the docs site structure.

## What Changes

- Move `assets/` (logo.png, gifs/, vhs/) to `docs/src/assets/`
- Update README.md references from `assets/` to `docs/src/assets/`
- Update VHS tape output paths from `assets/gifs/` to `docs/src/assets/gifs/`
- Update flake.nix generate-gif command paths
- Update all spec file path references

## Impact

- Affected specs: `documentation`, `ci-integration`
- Affected code:
  - `README.md` (5 asset references)
  - `assets/vhs/*.tape` (6 tape files - output paths)
  - `flake.nix` (generate-gif script)
  - `docs/src/content/docs/index.mdx` (hero image path)
- Affected archive: Historical references in `spectr/changes/archive/` remain
  unchanged (historical accuracy)
