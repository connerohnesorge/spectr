# Design: JSONC Validation Implementation

## Architecture

### Current Flow
```
tasks.md → parseTasksMd() → []Task → json.MarshalIndent() → JSONC string → Write to disk
```

### Proposed Flow with Validation
```
tasks.md → parseTasksMd() → []Task → json.MarshalIndent() → JSONC string
    → validateJSONCOutput() → Write to disk
                ↓ (if invalid)
              Error with details
```

## Components

### 1. JSONC Validator (`cmd/accept_validator.go`)

**Purpose**: Validate that marshalled JSON is valid JSONC

**Functions**:
- `validateJSONCOutput(jsonData []byte) error`
  - Parse with `StripJSONComments` + `json.Unmarshal`
  - Verify no data loss
  - Return detailed error if invalid

- `validateTasksJSONC(tasksFile TasksFile, jsonData []byte) error`
  - Round-trip validation: TasksFile → JSON → TasksFile → compare
  - Ensures no information loss

### 2. Property-Based Testing (`cmd/accept_validator_test.go`)

**Test Cases**:
1. **Special Characters**: `\n`, `\t`, `\"`, `\\`, `\r`, `\b`, `\f`
2. **Unicode**: Emoji, non-ASCII characters, different languages
3. **Edge Cases**: Very long descriptions, empty strings, only whitespace
4. **JSON Meta-characters**: `{`, `}`, `[`, `]`, `:`, `,`
5. **JSONC Comments**: Descriptions that look like comments `//`, `/*`

**Test Structure**:
```go
func TestJSONCValidation_PropertyBased(t *testing.T) {
    specialChars := []string{
        "\\", "\"", "\n", "\t", "\r",
        "{{", "}}", "//", "/*", "*/",
        "emoji: 🚀", "unicode: 你好",
    }

    for _, char := range specialChars {
        t.Run(fmt.Sprintf("char_%s", sanitize(char)), func(t *testing.T) {
            // Create task with special char
            // Marshal to JSONC
            // Validate round-trip
        })
    }
}
```

### 3. External JSONC Parser Integration

**Options**:
1. Use existing Go JSONC parser (research needed)
2. Shell out to Node.js `jsonc-parser` (heavier dependency)
3. Implement minimal JSONC parser for validation only

**Decision**: Use option 1 if available, fallback to custom validation

### 4. Error Reporting

When validation fails, provide:
- Exact location of invalid escape sequence
- The problematic task ID and description
- Suggested fix or sanitized version

Example error:
```
JSONC validation failed for task 3.2:
  Description: "Update README with \ backslash"
  Error: invalid escape sequence at position 25
  Suggestion: Escape backslash as "\\"
```

## Implementation Plan

### Phase 1: Core Validation (Priority: High)
- Create `validateJSONCOutput()` function
- Add to write paths in `accept_writer.go`
- Basic round-trip testing

### Phase 2: Property-Based Tests (Priority: High)
- Generate test cases with special characters
- Run against existing JSONC files in archive
- Fix any discovered issues

### Phase 3: External Parser (Priority: Medium)
- Research Go JSONC parser libraries
- Integrate if suitable option found
- Add as additional validation layer

### Phase 4: Fuzzing (Priority: Low)
- Set up go-fuzz infrastructure
- Create fuzz target for JSONC generation
- Run extended fuzzing campaigns

## Open Questions

1. **Performance**: Should we validate on every write or only in tests?
   - **Decision**: Validate on every write (fail fast), but make it optional via flag for performance-critical paths

2. **Sanitization vs Rejection**: Should we auto-fix invalid characters or reject?
   - **Decision**: Reject with clear error - don't silently change user data

3. **Which JSONC parser to use?**
   - **Decision**: Research in implementation phase, fallback to custom if needed

## Success Criteria

1. All existing tasks.jsonc files pass validation
2. Property-based tests with 100+ generated inputs all pass
3. Round-trip conversion preserves all data exactly
4. Clear error messages guide users to fix invalid input
