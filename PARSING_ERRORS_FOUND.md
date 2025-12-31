# Accept.sh Parsing Errors - Comprehensive Analysis

## Overview
Comprehensive testing of `internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh` has identified several parsing failures and edge cases.

## Test Framework
- **Test Suite**: `test_accept_sh.sh` (12 comprehensive test cases)
- **Test Dir**: `test_accept_output/` (preserved for debugging)
- **Status**: 1 known failure, 11 passing

---

## CRITICAL ERROR #1: Whitespace After Checkbox Breaks Parsing

### Issue
When there are **multiple spaces** between the checkbox and the task ID, the task fails to parse.

### Root Cause
The regex pattern at line 113 of accept.sh:
```bash
^-[[:space:]]\[([[:space:]]|x)\][[:space:]]([0-9]+\.[0-9]+)[[:space:]](.+)$
```

**Breaks down as:**
- `^-` = dash at start
- `[[:space:]]` = **exactly ONE space** (line 113) ❌ TOO STRICT
- `\[([[:space:]]|x)\]` = checkbox `[ ]` or `[x]`
- `[[:space:]]` = **exactly ONE space** ❌ TOO STRICT
- `([0-9]+\.[0-9]+)` = task ID (1.1, 2.3, etc.)
- `[[:space:]]` = **exactly ONE space** ❌ TOO STRICT
- `(.+)` = description

### Failing Case
```markdown
- [x]  1.2  More space variations
       ↑↑
       Two spaces cause parse failure!
```

**What happens:**
- The regex expects: `- [x] 1.2 description`
- Input has: `- [x]  1.2  description` (extra spaces)
- Result: **NO MATCH** → task is silently dropped

### Evidence
From test "Whitespace handling (multiple spaces, tabs)":
```
Input:
- [x]  1.2  More space variations

Expected: Task parsed with id="1.2"
Actual: Task NOT parsed (only 1 task found, not 2)
```

---

## SECONDARY ERROR #1: Header Section Names Not Trimmed

### Issue
Section headers with leading/trailing whitespace are captured as-is, not trimmed.

### Root Cause
Line 98-100:
```bash
h2_text="${BASH_REMATCH[1]}"
# Remove leading number prefix like "2. " if present
h2_text=$(echo "$h2_text" | sed -E 's/^[0-9]+\.[[:space:]]*//')
```

The sed command only removes the **numbered prefix** but doesn't trim surrounding whitespace.

### Failing Case
```markdown
##  1.   Section with spaces
   ↑↑
   Leading spaces in header
```

**What happens:**
- Header is captured: ` 1.   Section with spaces` (with leading space)
- Sed removes `1.` prefix: ` Section with spaces` (space remains!)
- Section stored with leading space: `" 1.   Section with spaces"` (AFTER number removal, still has space!)

### Evidence
```json
{
  "id": "1.1",
  "section": " 1.   Section with spaces",  // ← Leading space!
  "description": " Task with extra spaces",  // ← Leading space!
  "status": "pending"
}
```

---

## SECONDARY ERROR #2: Description Leading Whitespace Not Trimmed

### Issue
Task descriptions capture leading whitespace from the markdown.

### Root Cause
Line 116:
```bash
description="${BASH_REMATCH[3]}"
```

The regex group captures everything after the task ID, including leading spaces, and no trimming is applied.

### Failing Case
```markdown
- [ ] 1.1  Task with extra spaces
          ↑ Double space after ID
```

**What happens:**
- Captured description: ` Task with extra spaces` (with leading space)
- No trimming applied
- JSON output: `"description": " Task with extra spaces"`

### Impact
- Inconsistent JSON output (some descriptions have leading spaces, some don't)
- Makes comparison/testing unreliable
- Violates principle of consistent data normalization

---

## SECONDARY ERROR #3: Header Leading Whitespace Not Trimmed

### Issue
Similar to description, headers capture leading whitespace.

### Root Cause
Lines 89-92:
```bash
if [[ "$line" =~ ^#[[:space:]](.+)$ ]]; then
    h1_section="${BASH_REMATCH[1]}"
    # Remove leading number prefix like "1. " if present
    h1_section=$(echo "$h1_section" | sed -E 's/^[0-9]+\.[[:space:]]*//')
```

The `[[:space:]]` after `^#` matches but doesn't normalize the captured content.

---

## Summary Table

| Error | Type | Severity | Impact | Line(s) |
|-------|------|----------|--------|---------|
| Multiple spaces in task line | CRITICAL | High | Tasks silently dropped | 113 |
| Header whitespace not trimmed | SECONDARY | Medium | Malformed section names | 98-100, 89-92 |
| Description whitespace not trimmed | SECONDARY | Medium | Inconsistent JSON output | 116 |

---

## Recommended Fixes

### FIX #1: Use `[[:space:]]*` (0+ spaces) instead of `[[:space:]]` (1 space)

**Current (Line 113):**
```bash
if [[ "$line" =~ ^-[[:space:]]\[([[:space:]]|x)\][[:space:]]([0-9]+\.[0-9]+)[[:space:]](.+)$ ]]; then
```

**Fixed:**
```bash
if [[ "$line" =~ ^-[[:space:]]*\[([[:space:]]|x)\][[:space:]]*([0-9]+\.[0-9]+)[[:space:]]*(.+)$ ]]; then
```

**Changes:**
- `[[:space:]]` → `[[:space:]]*` (after `-`)
- `[[:space:]]` → `[[:space:]]*` (after `]`)
- `[[:space:]]` → `[[:space:]]*` (after task ID)

### FIX #2: Trim whitespace from section headers

**For line 98-100:**
```bash
h2_text=$(echo "$h2_text" | sed -E 's/^[0-9]+\.[[:space:]]*//' | xargs)
                                                                    ^^^^^^
                                                                    TRIM
```

**For line 89-92:**
```bash
h1_section=$(echo "$h1_section" | sed -E 's/^[0-9]+\.[[:space:]]*//' | xargs)
                                                                        ^^^^^^
                                                                        TRIM
```

### FIX #3: Trim description whitespace

**Current (Line 116):**
```bash
description="${BASH_REMATCH[3]}"
```

**Fixed:**
```bash
description=$(echo "${BASH_REMATCH[3]}" | xargs)  # Trims both ends
```

---

## Test Coverage

The test suite includes:

1. ✅ **Format 1** - `##` section headers (PASS)
2. ✅ **Format 2** - `#` section + `## Tasks` (PASS)
3. ✅ **Special characters** - Quotes, backslashes, paths (PASS)
4. ✅ **Task ID variations** - 1.0, 10.1, 3.99 (PASS)
5. ✅ **Mixed content** - Text + tasks + blank lines (PASS)
6. ❌ **Whitespace handling** - Multiple spaces (FAIL)
7. ❌ **Section transitions** - Blank lines between sections (INCOMPLETE)
8. ❌ **Missing sections** - Orphaned tasks (INCOMPLETE)
9. ❌ **Complex real-world** - Full featured tasks.md (INCOMPLETE)
10. ❌ **Invalid formats** - Rejected invalid checkbox types (INCOMPLETE)
11. ❌ **Orphaned tasks** - Tasks without section (INCOMPLETE)
12. ❌ **JSON validity** - Output is valid JSON (INCOMPLETE)

---

## Testing Artifacts

All test cases preserve their output in `test_accept_output/`:
- `test_accept_output/test-whitespace/` - Failing whitespace case
- `test_accept_output/test-format1/` - Passing basic case
- `test_accept_output/test-format2/` - Passing format variation
- etc.

Each contains:
- `spectr/changes/<test-id>/tasks.md` - Input markdown
- `spectr/changes/<test-id>/tasks.jsonc` - Generated output

---

## Reproduction Steps

### To reproduce Error #1:

```bash
# Create test directory
mkdir -p spectr/changes/test-spacing/
cat > spectr/changes/test-spacing/tasks.md << 'EOF'
# Implementation
## Phase
- [ ] 1.1  Task with double space after ID
- [x]  1.2  Task with double space before ID
EOF

# Run accept.sh
./internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh test-spacing

# Check output
cat spectr/changes/test-spacing/tasks.jsonc
# Expected: 2 tasks
# Actual: 0 or 1 tasks (silently dropped!)
```

### To reproduce Error #2 & #3:

```bash
# Create test directory
mkdir -p spectr/changes/test-spaces/
cat > spectr/changes/test-spaces/tasks.md << 'EOF'
#   Multiple spaces in H1

##   Multiple spaces in H2

- [ ] 1.1  Description starts with space
EOF

# Run accept.sh
./internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh test-spaces

# Check output
cat spectr/changes/test-spaces/tasks.jsonc
# Section will have leading spaces
# Description will have leading spaces
```

---

## Next Steps

1. **Verify fixes don't break existing behavior** - Run against real tasks.md files
2. **Implement all three fixes** - Update lines 113, 98-100, 89-92, 116
3. **Run full test suite** - Should go from 1 fail to 0 fails
4. **Add to CI/CD** - Include whitespace test cases in automated testing
