# Accept.sh Testing Summary

**Status**: ⚠️ **CRITICAL PARSING BUG FOUND**

## Executive Summary

Comprehensive testing of the `accept.sh` script has identified **1 critical bug** and **2 secondary issues** that affect data quality and user experience.

- **Critical Issue**: Tasks with extra whitespace are silently dropped from output (data loss)
- **Secondary Issues**: Whitespace in headers and descriptions not trimmed (cosmetic but inconsistent)

## Test Framework

| Component | Details |
|-----------|---------|
| **Test Script** | `test_accept_sh.sh` (12 test cases) |
| **Demo Script** | `test_critical_error.sh` (isolated reproduction) |
| **Test Output** | `test_accept_output/` (preserved for debugging) |
| **Current Status** | 11/12 tests passing, 1 failing |

## Critical Bug: Task Parsing Fails with Extra Whitespace

### Bug ID: PARSE_WHITESPACE_001

### Severity: **CRITICAL** (Data Loss)

### Description
The regex pattern used to parse task lines is too strict about whitespace. When markdown contains extra spaces (very common), tasks are silently dropped from the output.

### Demonstration
```markdown
# Input tasks.md
## Phase
- [ ] 1.1 Correct spacing (single space)
- [ ]  1.2 Extra spaces (two spaces) ← WILL BE DROPPED
```

```json
// Generated tasks.jsonc
{
  "tasks": [
    { "id": "1.1", ... },
    // 1.2 is MISSING - silently dropped!
  ]
}
```

### Root Cause
**File**: `internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh`
**Line**: 113

```bash
if [[ "$line" =~ ^-[[:space:]]\[([[:space:]]|x)\][[:space:]]([0-9]+\.[0-9]+)[[:space:]](.+)$ ]]; then
```

The pattern uses `[[:space:]]` which matches **exactly ONE space**:
- After `-`: `^-[[:space:]]` (1 space required)
- After checkbox: `\][[:space:]]` (1 space required)  
- After task ID: `\][0-9]+\.[0-9]+)[[:space:]]` (1 space required)

When markdown has 2+ spaces (common formatting), the entire pattern fails to match.

### Test Results

**Test Case**: `test_critical_error.sh`

```
Input:
- [ ] 1.1 Correct spacing (single space)
- [ ]  1.2 Extra spaces (SHOULD FAIL)
- [x] 2.1 Completed with correct spacing
- [x]  2.2 Completed with extra spaces (SHOULD FAIL)

Expected tasks: 4
Actual tasks: 2

Results:
✓ Task 1.1 parsed
✗ Task 1.2 DROPPED (extra spaces before ID)
✓ Task 2.1 parsed
✗ Task 2.2 DROPPED (extra spaces before ID)
```

### Impact
- **Severity**: HIGH - tasks disappear silently
- **User Awareness**: LOW - no error message, just missing tasks
- **Data Loss**: YES - tasks are lost without warning
- **Frequency**: COMMON - extra whitespace is typical in markdown

### Fix Required
Change `[[:space:]]` to `[[:space:]]*` to allow 0 or more spaces:

```bash
# Current (BROKEN):
^-[[:space:]]\[([[:space:]]|x)\][[:space:]]([0-9]+\.[0-9]+)[[:space:]](.+)$

# Fixed:
^-[[:space:]]*\[([[:space:]]|x)\][[:space:]]*([0-9]+\.[0-9]+)[[:space:]]*(.+)$
```

---

## Secondary Issue #1: Header Whitespace Not Trimmed

### Bug ID: HEADER_WHITESPACE_001
### Severity: **MEDIUM** (Data Quality)

### Description
Section headers capture leading whitespace, resulting in malformed section names in the JSON output.

### Example
```markdown
##  1.   Section with spaces
   ↑↑
   Extra spaces
```

Generated JSON:
```json
{
  "section": " 1.   Section with spaces"  // ← Leading space!
}
```

### Root Cause
Lines 89-100: The regex captures the header text with whitespace, then removes only the numbered prefix but doesn't trim remaining whitespace.

```bash
h2_text="${BASH_REMATCH[1]}"  # Captures with spaces
h2_text=$(echo "$h2_text" | sed -E 's/^[0-9]+\.[[:space:]]*//')  # Removes number but not spaces
```

### Fix Required
Add `xargs` to trim whitespace:
```bash
h2_text=$(echo "$h2_text" | sed -E 's/^[0-9]+\.[[:space:]]*//' | xargs)
```

---

## Secondary Issue #2: Description Whitespace Not Trimmed

### Bug ID: DESC_WHITESPACE_001
### Severity: **MEDIUM** (Data Quality)

### Description
Task descriptions capture leading whitespace, creating inconsistent output.

### Example
```markdown
- [ ] 1.1  Task with extra spaces
          ↑ Extra space here
```

Generated JSON:
```json
{
  "description": " Task with extra spaces"  // ← Leading space!
}
```

### Root Cause
Line 116: The regex group captures everything after the task ID without trimming.

### Fix Required
```bash
# Current:
description="${BASH_REMATCH[3]}"

# Fixed:
description=$(echo "${BASH_REMATCH[3]}" | xargs)
```

---

## Test Coverage

### ✅ Passing Tests (11/12)

1. **Format 1: Basic ## section headers** - PASS
2. **Format 2: # section with ## Tasks** - PASS
3. **Special characters in descriptions** - PASS
4. **Task ID variations (1.0, 10.1, 3.99)** - PASS
5. **Empty sections and non-task lines** - PASS
6. **Whitespace handling** - FAIL (1/2 tasks parsed)
7. **Section transitions** - PASS
8. **Tasks without section** - PASS
9. **Real-world complex tasks** - PASS (counts correct)
10. **Invalid format rejection** - PASS
11. **Orphaned tasks** - PASS
12. **JSON validity** - PASS

### ❌ Failing Tests (1/12)

- **Whitespace handling**: Extra spaces cause task loss
  - Expected: 2 tasks (1.1 with correct spacing, 1.2 with extra spaces)
  - Actual: 1 task parsed (1.2 dropped silently)

---

## Artifacts and Evidence

All test execution artifacts are preserved:

```
test_accept_output/
├── test-format1/           ✅ Passing
├── test-format2/           ✅ Passing
├── test-special/           ✅ Passing
├── test-ids/               ✅ Passing
├── test-mixed/             ✅ Passing
├── test-whitespace/        ❌ FAILING (see tasks.jsonc)
├── test-transitions/       ✅ Passing
├── test-nosection/         ✅ Passing
├── test-complex/           ✅ Passing
├── test-invalid/           ✅ Passing
├── test-orphan/            ✅ Passing
└── test-json/              ✅ Passing

test_critical_demo/
└── spectr/changes/demo/    ← Clear reproduction of critical bug
```

**View failing output**: `test_accept_output/test-whitespace/spectr/changes/test-whitespace/tasks.jsonc`

---

## Recommendations

### Priority 1: Fix Critical Bug
**What**: Update line 113 regex to use `[[:space:]]*` instead of `[[:space:]]`
**Why**: Prevents silent data loss
**Risk**: Very low - only changes whitespace tolerance
**Testing**: All tests should pass after fix

### Priority 2: Trim Whitespace
**What**: Add `| xargs` to lines handling section names and descriptions
**Why**: Improves data consistency
**Risk**: Very low - only cosmetic changes
**Testing**: JSON output will have cleaner values

### Priority 3: Add to CI/CD
**What**: Include whitespace test cases in automated testing
**Why**: Prevents regression
**Implementation**: Run `test_accept_sh.sh` in CI pipeline

---

## How to Reproduce Locally

### Critical Bug Reproduction
```bash
cd /path/to/spectr-src
./test_critical_error.sh
# Output shows tasks 1.2 and 2.2 silently dropped
```

### Full Test Suite
```bash
cd /path/to/spectr-src
./test_accept_sh.sh
# 11 tests pass, 1 test fails (whitespace handling)
```

### Manual Test
```bash
mkdir -p spectr/changes/test-spacing/
cat > spectr/changes/test-spacing/tasks.md << 'EOF'
# Phase
## Tasks
- [ ] 1.1 Single space (works)
- [ ]  1.2 Double space (fails)
EOF

internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh test-spacing
cat spectr/changes/test-spacing/tasks.jsonc
# Shows only 1.1, missing 1.2
```

---

## Next Steps

1. **Review this analysis** - Confirm findings
2. **Fix the regex** - Update line 113
3. **Trim whitespace** - Update lines 89-100, 116
4. **Run full test suite** - Verify 12/12 pass
5. **Update CI/CD** - Add automated testing
6. **Release** - Update version, document fix in changelog

---

## Test Execution Log

```
Test: Format 1: Basic ## section headers
✓ PASS

Test: Format 2: # section with ## Tasks subheader
✓ PASS

Test: Special characters in task descriptions
✓ PASS

Test: Various task ID formats (X.Y)
✓ PASS

Test: Empty sections and non-task lines
✓ PASS

Test: Whitespace handling (multiple spaces, tabs)
✗ FAIL - Only 1 task parsed instead of 2

Test: Section transitions with blank lines
✓ PASS

Test: Tasks without preceding section header
✓ PASS

Test: Real-world complex tasks.md
✓ PASS

Test: Invalid task formats are ignored
✓ PASS

Test: Orphaned task ID variations
✓ PASS

Test: Generated JSONC is valid JSON
✓ PASS

Results: 11 PASS, 1 FAIL
```
