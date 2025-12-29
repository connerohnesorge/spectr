# Change: Remove Track Command

## Summary

Fully remove the `spectr track` subcommand from the CLI. This is a hard breaking
change that removes the command entirely along with its supporting
infrastructure.

## Motivation

The `spectr track` command provides automated git commits when task statuses
change in `tasks.jsonc`. This feature is being removed to simplify the CLI and
reduce maintenance burden.

## Scope

### Code Removal

- Remove `cmd/track.go` - Track command implementation
- Remove `cmd/track_test.go` - Track command tests
- Remove `internal/track/` package - All track-related infrastructure:
  - `tracker.go`, `tracker_test.go`
  - `watcher.go`, `watcher_test.go`
  - `committer.go`, `committer_test.go`
  - `doc.go`
- Remove `internal/specterrs/track.go` - Track-specific error types
- Remove `TrackCmd` from `cmd/root.go` CLI struct

### Spec Updates

- Remove "Track Command" requirement from `cli-interface` spec
- Remove "Track Command Flags" requirement from `cli-interface` spec

## Breaking Change

This is an intentional breaking change. Users who depend on `spectr track` will
need to implement their own automation if required.

## Non-Goals

- Backward compatibility documentation
- Migration path
- Deprecation warnings
