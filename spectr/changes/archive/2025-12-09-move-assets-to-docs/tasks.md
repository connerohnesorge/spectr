# Implementation Tasks

## 1. Move Assets Directory

- [ ] 1.1 Move `assets/` directory to `docs/src/assets/` (preserving logo.png,
  gifs/, vhs/)
- [ ] 1.2 Remove empty `assets/` directory from root

## 2. Update README.md References

- [ ] 2.1 Update logo reference (line 3): `assets/logo.png` ->
  `docs/src/assets/logo.png`
- [ ] 2.2 Update init.gif reference (line 289): `assets/gifs/init.gif` ->
  `docs/src/assets/gifs/init.gif`
- [ ] 2.3 Update list.gif reference (line 326): `assets/gifs/list.gif` ->
  `docs/src/assets/gifs/list.gif`
- [ ] 2.4 Update validate.gif reference (line 367): `assets/gifs/validate.gif`
  -> `docs/src/assets/gifs/validate.gif`
- [ ] 2.5 Update archive.gif reference (line 471): `assets/gifs/archive.gif` ->
  `docs/src/assets/gifs/archive.gif`

## 3. Update VHS Tape Output Paths

- [ ] 3.1 Update `init.tape` output path to `docs/src/assets/gifs/init.gif`
- [ ] 3.2 Update `list.tape` output path to `docs/src/assets/gifs/list.gif`
- [ ] 3.3 Update `validate.tape` output path to
  `docs/src/assets/gifs/validate.gif`
- [ ] 3.4 Update `archive.tape` output path to
  `docs/src/assets/gifs/archive.gif`
- [ ] 3.5 Update `partial-match.tape` output path to
  `docs/src/assets/gifs/partial-match.gif`
- [ ] 3.6 Update `pr-hotkey.tape` output path to
  `docs/src/assets/gifs/pr-hotkey.gif`

## 4. Update flake.nix

- [ ] 4.1 Update generate-gif script paths from `assets/` to `docs/src/assets/`

## 5. Update Docs Site References

- [ ] 5.1 Update `docs/src/content/docs/index.mdx` hero image path (now same
  directory)

## 6. Validation

- [ ] 6.1 Run `spectr validate move-assets-to-docs --strict`
- [ ] 6.2 Verify README GIF images display correctly
- [ ] 6.3 Verify docs site builds successfully with `npm run build` in docs/
- [ ] 6.4 Verify generate-gif command works in nix develop shell
