# Change: Add CI automation for VHS-generated GIFs

## Why

Demo GIFs in `assets/gifs/` can drift out of sync with CLI changes. Currently,
developers must manually run VHS tape files to regenerate GIFs after updating
commands. Automating this via GitHub Actions ensures documentation stays current
with the codebase.

## What Changes

- Add new GitHub Actions workflow `.github/workflows/vhs.yml`
- Workflow triggers on changes to `*.tape` files in `assets/vhs/`
- Uses `charmbracelet/vhs-action@v2` to process all tape files
- Auto-commits generated GIFs back to the repository
- Generates missing GIF for `partial-match.tape`

## Impact

- Affected specs: `ci-integration`
- Affected code: `.github/workflows/vhs.yml` (new file)
- No breaking changes
