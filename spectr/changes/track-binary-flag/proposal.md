# Change: Add Binary File Filtering to Track Command

## Why

The `spectr track` command currently commits all modified files except task files. This includes binary files (images, compiled artifacts, executables, etc.) which can bloat repository history and are often unintentional additions. Developers need explicit control over whether binary files are tracked and committed automatically.

## What Changes

- Add `--include-binaries` flag to `spectr track` command
- By default, binary files are excluded from automatic commits
- When `--include-binaries` flag is provided, binary files are included as before
- Binary detection uses git's internal binary detection (`git diff --numstat`)
- User-friendly warning messages when binary files are skipped

## Impact

- Affected specs: cli-interface
- Affected code:
  - `cmd/track.go` - Add IncludeBinaries flag to TrackCmd struct
  - `internal/track/committer.go` - Add binary detection and filtering logic
  - `internal/track/committer_test.go` - Add tests for binary filtering
  - `internal/track/tracker.go` - Pass flag through Config
