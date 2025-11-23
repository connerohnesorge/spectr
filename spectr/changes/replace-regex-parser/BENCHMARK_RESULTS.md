# Benchmark Results: Regex vs Lexer/Parser

**Date**: 2025-11-23
**Hardware**: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
**Test Suite**: `internal/parsers/benchmark_test.go`

## Executive Summary

The lexer/parser implementation **PASSES** the performance acceptance criteria (< 2x regression) and provides **BETTER CORRECTNESS** than the regex implementation. The results show:

- Small files: **12% FASTER** with lexer
- Medium files: 25% slower (within 2x threshold)
- Large files: 23% slower (within 2x threshold)
- Pathological cases: 17% slower (within 2x threshold)
- Delta parsing: **6.7x FASTER** with lexer

**Recommendation**: Proceed with lexer/parser migration.

## Detailed Benchmark Results

### Requirement Parsing Benchmarks

| File Type | Implementation | Time (ns/op) | Memory (B/op) | Allocations | Ratio |
|-----------|----------------|--------------|---------------|-------------|-------|
| **Small** (~100 lines, 3-5 reqs) | | | | | |
| | Regex | 23,023 | 22,732 | 190 | baseline |
| | Lexer | 20,358 | 18,328 | 174 | **0.88x** ✓ |
| **Medium** (~400 lines, 10-20 reqs) | | | | | |
| | Regex | 76,258 | 75,066 | 445 | baseline |
| | Lexer | 95,128 | 99,464 | 729 | **1.25x** ✓ |
| **Large** (~1000+ lines, 50+ reqs) | | | | | |
| | Regex | 207,512 | 198,522 | 1,321 | baseline |
| | Lexer | 254,639 | 318,104 | 2,615 | **1.23x** ✓ |
| **Pathological** (edge cases) | | | | | |
| | Regex | 110,250 | 119,134 | 758 | baseline |
| | Lexer | 129,378 | 171,472 | 1,222 | **1.17x** ✓ |

### Delta Parsing Benchmarks

| File Type | Implementation | Time (ns/op) | Memory (B/op) | Allocations | Ratio |
|-----------|----------------|--------------|---------------|-------------|-------|
| **Small** (delta spec) | | | | | |
| | Regex | 135,795 | 81,306 | 631 | baseline |
| | Lexer | 20,238 | 17,880 | 167 | **0.15x** ✓✓ |

## Correctness Validation

All correctness tests **PASSED**:

### Test Results

| Test Case | Regex Requirements | Lexer Requirements | Status | Notes |
|-----------|-------------------|-------------------|---------|-------|
| small.md | 4 | 4 | ✓ PASS | Identical results |
| medium.md | 10 | 10 | ✓ PASS | Identical results |
| pathological.md | 27 | 25 | ✓ PASS | Lexer is correct! |
| delta parsing | All sections match | All sections match | ✓ PASS | Identical results |

### Pathological Case Analysis

The pathological test file contains code blocks with markdown-like text:

```markdown
\`\`\`markdown
### Requirement: Fake requirement in code
This should NOT be parsed as a real requirement
\`\`\`
```

**Results**:
- **Regex**: Found 27 requirements (INCORRECT - included requirements from code blocks)
- **Lexer**: Found 25 requirements (CORRECT - properly ignored code blocks)

This demonstrates that the lexer/parser handles edge cases correctly while regex does not.

## Performance Analysis

### Small Files (< 100 lines)
- **Lexer is 12% FASTER** than regex
- Uses **19% less memory**
- Fewer allocations (174 vs 190)
- Likely due to reduced regex compilation overhead

### Medium Files (100-500 lines)
- Lexer is 25% slower (within acceptable threshold)
- Uses 32% more memory for AST
- For typical 400-line spec: 95μs vs 76μs = **19μs difference**
- Absolute time is negligible for non-hot-path parsing

### Large Files (1000+ lines)
- Lexer is 23% slower (within acceptable threshold)
- Uses 60% more memory for AST
- For 1000-line spec: 255μs vs 208μs = **47μs difference**
- Still negligible for infrequent operations

### Delta Parsing (Most Common Operation)
- Lexer is **6.7x FASTER** than regex
- Uses **78% less memory**
- Uses **74% fewer allocations**
- This is the most frequently used operation in Spectr workflows

## Memory Characteristics

| Metric | Small | Medium | Large | Pathological |
|--------|-------|--------|-------|--------------|
| Memory Ratio | 0.81x (less) | 1.32x (more) | 1.60x (more) | 1.44x (more) |
| Allocation Ratio | 0.92x (less) | 1.64x (more) | 1.98x (more) | 1.61x (more) |

**Analysis**:
- Small files benefit from reduced overhead
- Larger files trade memory for AST structure
- Memory overhead is acceptable given:
  - Parsing is not in hot path
  - Typical files are 100-500 lines
  - Correctness > memory efficiency

## Decision Matrix

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Correctness | 100% | 100% | ✓ PASS |
| Performance (small) | < 2x | 0.88x | ✓ PASS |
| Performance (medium) | < 2x | 1.25x | ✓ PASS |
| Performance (large) | < 2x | 1.23x | ✓ PASS |
| Performance (pathological) | < 2x | 1.17x | ✓ PASS |
| Performance (delta) | < 2x | 0.15x | ✓ PASS |
| Edge case handling | Better | Better | ✓ PASS |

**Overall**: **7/7 criteria met** ✓

## Recommendation

**PROCEED with lexer/parser migration**

### Rationale

1. **Performance is Excellent**
   - All benchmarks within 2x threshold
   - Actually faster on small files and delta parsing
   - Delta parsing (most common) is 6.7x faster
   - Absolute time differences are negligible (< 50μs)

2. **Correctness is Superior**
   - Properly handles code blocks with markdown syntax
   - Provides accurate requirement counting
   - Enables better error messages (line/column info)

3. **Maintainability Wins**
   - Clear separation of concerns (lexing, parsing, extraction)
   - Easier to extend and modify
   - Generic markdown parser is reusable
   - No brittle regex patterns

4. **Real-World Impact**
   - Parsing is not in hot path (user interactions, not loops)
   - Files are typically 100-500 lines
   - Users won't notice < 50μs differences
   - Correctness matters more than microseconds

5. **Future-Proof**
   - AST enables new features (pretty printing, refactoring, etc.)
   - State machine handles complex markdown correctly
   - Extensible architecture for new requirements

## Next Steps

Proceed to Phase 3: Decision Gate - Performance Validation
- [x] 3.1 Analyze benchmark results → See above analysis
- [x] 3.2 Verify correctness → All tests pass, edge cases handled correctly
- [x] 3.3 Decision → Performance acceptable, no optimization needed
- [ ] 3.4 Document trade-offs → Complete in this document
- [ ] 3.5 Get approval to proceed → Ready for review

After approval, continue to Phase 4: Build Spectr Extractors.

## Test Corpus

The benchmark suite uses four test files:

1. **small.md** (< 100 lines)
   - 3 requirements in ADDED section
   - 1 requirement in MODIFIED section
   - 1 requirement in REMOVED section
   - Represents typical change deltas

2. **medium.md** (~400 lines)
   - 10 requirements with multiple scenarios each
   - Includes code blocks and examples
   - Represents typical capability specs

3. **large.md** (~1000 lines)
   - 50+ requirements
   - Stress test for performance
   - Represents large specification documents

4. **pathological.md** (~400 lines)
   - Edge cases that break regex parsers
   - Code blocks containing markdown syntax
   - Deeply nested structures
   - Special characters and escaping
   - Represents real-world complexity

All test files are located in `testdata/benchmarks/`.

## Raw Benchmark Output

```
goos: linux
goarch: amd64
pkg: github.com/connerohnesorge/spectr/internal/parsers
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkRegexRequirementParser_Small-16           	   51162	     23023 ns/op	   22732 B/op	     190 allocs/op
BenchmarkRegexRequirementParser_Medium-16          	   16905	     76258 ns/op	   75066 B/op	     445 allocs/op
BenchmarkRegexRequirementParser_Large-16           	    6139	    207512 ns/op	  198522 B/op	    1321 allocs/op
BenchmarkRegexRequirementParser_Pathological-16    	   11210	    110250 ns/op	  119134 B/op	     758 allocs/op
BenchmarkRegexDeltaParser_Small-16                 	    9032	    135795 ns/op	   81306 B/op	     631 allocs/op
BenchmarkLexerRequirementParser_Small-16           	   60540	     20358 ns/op	   18328 B/op	     174 allocs/op
BenchmarkLexerRequirementParser_Medium-16          	   14859	     95128 ns/op	   99464 B/op	     729 allocs/op
BenchmarkLexerRequirementParser_Large-16           	    4935	    254639 ns/op	  318104 B/op	    2615 allocs/op
BenchmarkLexerRequirementParser_Pathological-16    	    7962	    129378 ns/op	  171472 B/op	    1222 allocs/op
BenchmarkLexerDeltaParser_Small-16                 	   64808	     20238 ns/op	   17880 B/op	     167 allocs/op
PASS
ok  	github.com/connerohnesorge/spectr/internal/parsers	15.697s
```
