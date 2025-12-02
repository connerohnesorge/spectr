# Change: Add tasks.md Structure Validation

## Why
The `tasks.md` file format currently lacks validation, which can lead to inconsistent task organization. Enforcing a stricter structure with numbered section headers (`## 1.`, `## 2.`, etc.) ensures predictable parsing, better task grouping, and easier navigation for both humans and tooling.

## What Changes
- Add validation rules for `tasks.md` file structure
- Require numbered section headers using `## [number]. Section Name` format
- Validate task items are nested under numbered sections
- Add warnings for empty sections and tasks without section grouping
- Update documentation to reflect required format

## Impact
- Affected specs: validation
- Affected code: `internal/validation/`, `internal/parsers/parsers.go`
- Existing changes with non-compliant `tasks.md` files will see validation warnings (not errors, for backward compatibility)
