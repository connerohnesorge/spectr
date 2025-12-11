# Change: Add path normalization for item commands

## Why

When users or agents run `spectr validate spectr/changes/replace-regex-with-blackfriday`, the command fails with "item not found" because the command expects just the change ID (`replace-regex-with-blackfriday`), not the full path. This is a common user expectation, especially when tab-completing paths or copy-pasting from file explorers.

## What Changes

- All commands accepting item names (validate, archive, accept) will normalize path arguments
- If the input matches `spectr/changes/<id>` or `spectr/changes/<id>/...`, extract `<id>` as the change ID and infer type as "change"
- If the input matches `spectr/specs/<id>` or `spectr/specs/<id>/...`, extract `<id>` as the spec ID and infer type as "spec"
- The inferred type from path structure takes precedence over auto-detection (avoids ambiguity errors)
- Existing behavior for simple IDs (e.g., `replace-regex-with-blackfriday`) remains unchanged

## Impact

- Affected specs: `cli-framework`
- Affected code:
  - `internal/discovery/normalize.go` (new file with shared normalization function)
  - `cmd/validate.go` (call normalization)
  - `cmd/archive.go` (call normalization)
  - `cmd/accept.go` (call normalization)
