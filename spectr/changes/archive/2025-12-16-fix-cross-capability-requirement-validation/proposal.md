# Change: Fix Cross-Capability Requirement Name Validation

## Why

The current validator incorrectly flags an error when the same requirement name
appears in delta specs for different capabilities. For example, if
`support-aider` and `support-cursor` both have a requirement named "No
Instruction File", and a change modifies both, the validator reports:

```text
REMOVED Requirement 'No Instruction File': Requirement
'No Instruction File' is REMOVED in multiple files
```

This is incorrect because these are **different requirements** that happen to
share a name. Each capability has its own namespace, so `support-aider::No
Instruction File` is distinct from `support-cursor::No Instruction File`.

This limitation blocks legitimate multi-capability changes like provider
architecture redesigns.

## What Changes

- **FIX**: Validator should scope requirement name uniqueness to the capability
  level, not globally across all delta files
- **FIX**: When checking for duplicate MODIFIED/REMOVED requirements, compare
  `(capability, requirement_name)` tuples, not just `requirement_name`

## Impact

- Affected specs: `validation` (if it exists)
- Affected code:
  - `internal/validation/` - Delta validation logic
  - Specifically the code that checks for duplicate requirement modifications
