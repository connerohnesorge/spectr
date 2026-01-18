# Change: Fix JSONC File Validation

## Summary

The current JSONC file creation via `json.MarshalIndent` does not guarantee valid JSONC output when task descriptions contain special characters or escape sequences. This change introduces robust JSONC validation through property-based testing and round-trip verification to ensure all generated JSONC files are parseable.

## Problem

When tasks.md files contain task descriptions with:
- Backslashes, quotes, or other special characters that need escaping
- Unicode characters or emoji
- Newlines, tabs, or control characters

The `json.Marshal` function may produce output that:
1. Has invalid escape sequences for JSONC format
2. Fails to parse correctly when read back via `StripJSONComments` + `json.Unmarshal`
3. Loses information during round-trip conversion (write → read → write)

This breaks the reliability of the tasks.jsonc format and can cause failures during `spectr accept` regeneration.

## Solution

1. **Add JSONC validation function**: Create `validateJSONCOutput()` that verifies generated JSONC can be successfully parsed back
2. **Implement round-trip testing**: Ensure Task → JSON → Task produces identical results
3. **Add property-based tests**: Generate random task descriptions with special characters to find edge cases
4. **Add escape sequence sanitization**: Pre-process task descriptions before marshalling if needed
5. **Use external JSONC parser**: Integrate a JSONC parser library to validate output beyond basic JSON parsing

## Scope

- **In scope**:
  - Validation of JSONC output from `writeTasksJSONC`, `writeRootTasksJSONC`, `writeChildTasksJSONC`
  - Round-trip testing for Task struct serialization
  - Property-based testing with special characters
  - Integration of JSONC parser for validation

- **Out of scope**:
  - Changing the overall tasks.jsonc format or schema
  - Modifying how tasks.md is parsed
  - Performance optimizations

## Risks

- **Low risk**: Adding validation is a defensive measure that shouldn't break existing functionality
- **Test complexity**: Property-based tests may find edge cases that require careful handling

## Testing Strategy

1. Unit tests for JSONC validation function
2. Property-based tests with fuzzing-like random inputs
3. Round-trip tests for all Task field combinations
4. Integration tests with real tasks.md files from archived changes
