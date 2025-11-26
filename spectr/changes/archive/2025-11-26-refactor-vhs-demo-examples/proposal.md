# Change: Refactor VHS Demo Examples

## Why
The current VHS tape files are cluttered with verbose `cat > file << 'EOF'` heredoc blocks that obscure the actual spectr commands being demonstrated. This makes it hard for users to understand what the tool does at a glance. Additionally, there's no easy way to regenerate all demo GIFs at once.

## What Changes
- Create `examples/` directory with separate pre-made spectr projects per demo
- Refactor VHS tapes to use pre-made examples instead of inline heredocs
- Add a `generate-gif` command (via flake.nix) to build all VHS GIFs in one command
- Keep demos focused on showcasing spectr commands, not file creation boilerplate

## Impact
- Affected specs: documentation
- Affected files:
  - `assets/vhs/*.tape` (all 5 tapes: archive, init, list, validate, workflow)
  - New: `examples/{archive,init,list,validate,workflow}/` directories
  - Modified: `flake.nix` (added `generate-gif` script)
- Non-breaking: demo assets only, no CLI changes
