# Change: Fix VHS Tape Execution Environment

## Why

The current VHS tape files use an inefficient pattern of copying example
directories to a temporary `_demo` folder, running commands with `cd _demo &&
...`, then cleaning up. This adds unnecessary setup commands and makes the tapes
harder to read. Additionally, there are multiple `echo ''` statements that only
add blank lines without providing value.

## What Changes

- Modify VHS tapes to run directly in the `examples/<name>/` directories instead
  of copying to `_demo`
- Remove useless `echo ''` statements that only print blank lines
- Use VHS `Hide`/`Show` directives to hide necessary setup commands from the
  recording
- Simplify the init.tape to use a proper temporary directory pattern
- Update tape comments to reflect the new execution model

## Impact

- Affected specs: `documentation` (VHS tapes are part of documentation assets)
- Affected code: `assets/vhs/*.tape` (5 tape files)
- Affected directories: `examples/` structure remains unchanged
- No breaking changes - tapes will produce the same visual output with cleaner
  source
