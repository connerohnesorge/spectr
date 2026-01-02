# Spectr Validate with Binary

## Description
Validates Spectr specifications and change proposals using the `spectr` binary when available. This provides better performance and full feature parity compared to the bash script fallback.

## Usage

### Validate a specific change proposal:
```
spectr-validate-w-spectr-bin <change-id>
```

### Validate all specs and changes:
```
spectr-validate-w-spectr-bin
```

## Requirements
- The `spectr` binary must be available in PATH
- Must be run from within a Spectr-enabled project directory

## Implementation
This skill wraps the `spectr validate` command, passing through all arguments and options to provide native binary performance and feature support.