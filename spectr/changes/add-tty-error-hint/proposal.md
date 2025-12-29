# Change: Add helpful hint when TTY is unavailable

## Why

When users run `spectr init` in environments without a TTY (CI pipelines, Docker containers, piped commands), they receive a cryptic error message:

```text
spectr: error: wizard failed: could not open a new TTY: open /dev/tty: no such device or address
```

This provides no guidance on how to resolve the issue, leaving users to search for documentation or guess at solutions.

## What Changes

- Detect TTY-related errors from the Bubbletea TUI framework
- Enhance the error message with a helpful hint suggesting `--non-interactive` flag usage
- Provide a concrete example command in the error output

## Impact

- Affected specs: `cli-interface`
- Affected code: `cmd/init.go` (runInteractiveInit function)
