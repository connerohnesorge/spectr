# Remove --strict Flag: Always Validate Strictly

## Summary

Remove the `--strict` flag from the `spectr validate` command and make strict
validation (warnings treated as errors) the default and only behavior. This
simplifies the CLI and ensures consistent validation quality across all use
cases.

## Motivation

1. **Simplicity**: The `--strict` flag adds cognitive overhead without providing
  significant value. Users rarely need lenient validation.
2. **Quality Enforcement**: Strict validation catches more issues upfront,
  improving spec quality.
3. **CI/CD Consistency**: Removes the risk of different validation behavior
  between local development and CI pipelines.
4. **Reduced API Surface**: One less flag to document, test, and maintain.

## Changes

### Behavioral Changes

1. Validation always treats warnings as errors (strict mode)
2. Exit code is 1 if any warnings or errors exist
3. The `--strict` flag is removed from the CLI

### Affected Components

- `cmd/validate.go`: Remove `Strict` field and always pass `true` to validator
- `internal/validation/validator.go`: Simplify or remove strict mode handling
- `internal/validation/spec_rules.go`: Always apply warning-to-error conversion
- `internal/validation/change_rules.go`: Always apply warning-to-error
  conversion
- `internal/validation/interactive.go`: Remove strict parameter
- `spectr/specs/cli-interface/spec.md`: Update specification
- `spectr/specs/validation/spec.md`: Update specification
- `spectr/specs/documentation/spec.md`: Update specification to remove --strict
  reference

## Backward Compatibility

This is a **breaking change** for users who:

1. Rely on lenient validation to pass with warnings
2. Have CI/CD pipelines using `--strict` flag (will fail with "unknown flag"
  error)

**Migration path**: Users should fix warnings in their specs before upgrading,
or remove `--strict` from their CI commands.

## Alternatives Considered

1. **Keep as-is**: Rejected because strict mode is the recommended practice and
  should be the default.
2. **Add `--lenient` flag instead**: Rejected as it increases complexity without
  strong use case.
3. **Default to strict with `--lenient` opt-out**: Rejected as unnecessary
  complexity.
