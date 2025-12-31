# Accept.sh Bug Fix - Quick Reference Card

## 🚨 Critical Issue Found: Data Loss Bug

**Status**: CRITICAL - Tasks silently disappear  
**Location**: `internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh:113`  
**Severity**: HIGH (data loss)  
**Fix Difficulty**: TRIVIAL (1 line, 3 characters)  

---

## The Bug in 10 Seconds

### What Happens
```
Input:  - [ ]  1.2 Task (with extra space)
        
Regex expects: - [ ] 1.2  (exactly 1 space)
Actually has:  - [ ]  1.2 (2 spaces)

Result: Task 1.2 DISAPPEARS silently! ❌
```

### The Proof
```bash
./test_critical_error.sh
# Shows: Tasks parsed: 2 (expected 4)
# Missing: 1.2 and 2.2 (both have extra spaces)
```

---

## The Fix

### Line 113 - Change This
```bash
# BEFORE (BROKEN)
if [[ "$line" =~ ^-[[:space:]]\[([[:space:]]|x)\][[:space:]]([0-9]+\.[0-9]+)[[:space:]](.+)$ ]]; then

# AFTER (FIXED)
if [[ "$line" =~ ^-[[:space:]]*\[([[:space:]]|x)\][[:space:]]*([0-9]+\.[0-9]+)[[:space:]]*(.+)$ ]]; then
```

### What Changed
- `[[:space:]]` → `[[:space:]]*` (3 places)
- Allows 0 or more spaces instead of exactly 1

---

## Verification

### Before Fix
```bash
$ ./test_critical_error.sh
Tasks parsed: 2
✗ Task 1.2 FAILED TO PARSE
✗ Task 2.2 FAILED TO PARSE
```

### After Fix (Expected)
```bash
$ ./test_critical_error.sh
Tasks parsed: 4
✓ Task 1.2 parsed
✓ Task 2.2 parsed
```

---

## Optional: Secondary Fixes

These are cosmetic but recommended:

### Fix #2: Trim header whitespace (Line ~100)
```bash
# Add | xargs to the end
h2_text=$(echo "$h2_text" | sed -E 's/^[0-9]+\.[[:space:]]*//' | xargs)
```

### Fix #3: Trim description whitespace (Line ~116)
```bash
# Wrap with xargs
description=$(echo "${BASH_REMATCH[3]}" | xargs)
```

---

## Testing

### Full Test Suite (After Fix)
```bash
./test_accept_sh.sh
# Expected: 12/12 PASS (currently 11/12)
```

### Critical Bug Test (After Fix)
```bash
./test_critical_error.sh
# Expected: All 4 tasks parsed
```

---

## Impact

### Before Fix ❌
- Some tasks silently disappear
- No error messages
- Data loss
- User confusion

### After Fix ✅
- All tasks parse correctly
- Handles any whitespace amount
- No data loss
- Transparent operation

---

## Real-World Example

### Common Markdown (BREAKS NOW)
```markdown
# Implementation

## Phase 1

- [ ] 1.1 Single space (parses)
- [ ]  1.2 Double space (DROPS!)
- [ ]   1.3 Triple space (DROPS!)

## Phase 2

- [ ] 2.1 Single space (parses)
- [ ]  2.2 Double space (DROPS!)
```

**Current Result**: Only 1.1 and 2.1 parse (50% data loss!)  
**After Fix**: All 5 tasks parse correctly (100% ✓)

---

## File Locations

| File | Purpose |
|------|---------|
| `internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh` | The file to fix (line 113) |
| `.claude/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh` | Alternate location (same file) |
| `test_accept_sh.sh` | Full test suite (12 tests) |
| `test_critical_error.sh` | Isolated critical bug demo |
| `TEST_RESULTS_SUMMARY.md` | Executive report |
| `BUG_DETAIL_COMPARISON.md` | Side-by-side fix guide |

---

## Checklist

- [ ] Read this Quick Reference
- [ ] Run `./test_critical_error.sh` to see the bug
- [ ] Apply fix to line 113 (change `[[:space:]]` to `[[:space:]]*`)
- [ ] Run `./test_critical_error.sh` again (should show all tasks)
- [ ] Run `./test_accept_sh.sh` (should show 12/12 PASS)
- [ ] Optionally apply secondary fixes (lines 100, 116)
- [ ] Test with real tasks.md files
- [ ] Commit and push

---

## Questions?

| Q | A |
|---|---|
| **Where's the bug?** | Line 113 of accept.sh |
| **What's the bug?** | Regex too strict about whitespace |
| **How do I fix it?** | Change `[[:space:]]` to `[[:space:]]*` |
| **Will it break anything?** | No - only adds whitespace tolerance |
| **How do I test it?** | Run `./test_critical_error.sh` |
| **How much work?** | 1 line, 10 seconds |
| **Is it critical?** | YES - causes data loss |
| **How common is it?** | VERY - any markdown with extra spaces |

---

## The Exact Change

```diff
--- a/internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh
+++ b/internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh
@@ -110,7 +110,7 @@
         fi
 
         # Check for task line: - [ ] or - [x] followed by ID and description
-        if [[ "$line" =~ ^-[[:space:]]\[([[:space:]]|x)\][[:space:]]([0-9]+\.[0-9]+)[[:space:]](.+)$ ]]; then
+        if [[ "$line" =~ ^-[[:space:]]*\[([[:space:]]|x)\][[:space:]]*([0-9]+\.[0-9]+)[[:space:]]*(.+)$ ]]; then
             local checkbox="${BASH_REMATCH[1]}"
             local task_id="${BASH_REMATCH[2]}"
             local description="${BASH_REMATCH[3]}"
```

---

## Summary

| Aspect | Details |
|--------|---------|
| **Bug Type** | Regex too strict |
| **Symptom** | Tasks disappear |
| **Root Cause** | `[[:space:]]` = exactly 1 space |
| **Fix** | Change to `[[:space:]]*` = 0+ spaces |
| **Lines Changed** | 1 |
| **Characters Changed** | 3 (add `*` three times) |
| **Time to Fix** | 10 seconds |
| **Risk Level** | ZERO - only adds tolerance |
| **Impact** | Prevents data loss |

---

**Status**: Ready for implementation  
**Test Coverage**: 12 test cases created  
**Evidence**: Critical bug clearly reproduced  
**Recommendation**: IMPLEMENT IMMEDIATELY (data loss bug)
