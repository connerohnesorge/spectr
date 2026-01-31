#!/usr/bin/env bash
#
# accept_test.sh - TDD test suite for accept.sh
#
# Comprehensive test cases covering all edge cases and the main bug fixes:
# 1. Indented sub-tasks (critical bug - currently fails)
# 2. Flexible task number formats (currently fails)
# 3. Tasks without numbers (currently fails)
# 4. Section-based hierarchical structure
# 5. V2 JSONC output format
#

set -eo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Get script directory (handle both execution and sourcing)
if [ -n "${BASH_SOURCE[0]:-}" ]; then
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
else
    SCRIPT_DIR="$(pwd)"
fi
ACCEPT_SCRIPT="${SCRIPT_DIR}/accept.sh"

# Ensure accept.sh exists
if [ ! -f "$ACCEPT_SCRIPT" ]; then
    echo "Error: accept.sh not found at $ACCEPT_SCRIPT"
    echo "Script dir: $SCRIPT_DIR"
    exit 1
fi

# Test helper: Run a test case
run_test() {
    local test_name="$1"
    local input_md="$2"
    local expected_task_count="$3"
    local verify_fn="${4:-}"  # Optional custom verification function

    ((TESTS_RUN++))

    # Setup temp directory
    local temp_dir=$(mktemp -d)
    local change_id="test-$$-$TESTS_RUN"
    mkdir -p "$temp_dir/spectr/changes/$change_id"
    echo "$input_md" > "$temp_dir/spectr/changes/$change_id/tasks.md"

    # Run accept.sh (redirect output but don't fail on non-zero exit)
    cd "$temp_dir"
    bash "$ACCEPT_SCRIPT" "$change_id" >/dev/null 2>&1
    local exit_code=$?

    if [ $exit_code -eq 0 ]; then
        # Verify task count
        local tasks_file="spectr/changes/$change_id/tasks.jsonc"
        if [ ! -f "$tasks_file" ]; then
            echo -e "${RED}✗ $test_name${NC} (tasks.jsonc not created)"
            ((TESTS_FAILED++))
            rm -rf "$temp_dir"
            return 1
        fi

        # Strip comments from JSONC before parsing with jq
        local actual_count=$(grep -v '^//' "$tasks_file" | jq '.tasks | length' 2>/dev/null || echo "0")

        # Run custom verification if provided
        if [ -n "$verify_fn" ]; then
            if ! $verify_fn "$temp_dir/spectr/changes/$change_id"; then
                echo -e "${RED}✗ $test_name${NC} (custom verification failed)"
                ((TESTS_FAILED++))
                rm -rf "$temp_dir"
                return 1
            fi
        fi

        if [ "$actual_count" -eq "$expected_task_count" ]; then
            echo -e "${GREEN}✓ $test_name${NC}"
            ((TESTS_PASSED++))
        else
            echo -e "${RED}✗ $test_name${NC} (expected $expected_task_count tasks, got $actual_count)"
            ((TESTS_FAILED++))
        fi
    else
        echo -e "${RED}✗ $test_name${NC} (accept.sh failed)"
        ((TESTS_FAILED++))
    fi

    # Cleanup
    rm -rf "$temp_dir"
}

# Test 1: Basic numbered tasks (existing behavior - should pass)
test_basic_numbered() {
    local input='## Phase 1: Setup
- [ ] 1.1 First task
- [x] 1.2 Second task (completed)
- [ ] 1.3 Third task'

    run_test "Basic numbered tasks" "$input" 3
}

# Test 2: Indented sub-tasks (CRITICAL BUG - currently fails)
test_indented_subtasks() {
    local input='## Phase 1: Setup
- [ ] 1.1 Parent task
  - [ ] 1.2 Sub-task one (2 spaces)
  - [ ] 1.3 Sub-task two (2 spaces)
    - [ ] 1.4 Nested sub-task (4 spaces)'

    run_test "Indented sub-tasks" "$input" 4
}

# Test 3: Tasks without numbers (currently fails)
test_tasks_without_numbers() {
    local input='## Phase 1: Setup
- [ ] Task without number
- [ ] Another task without number'

    run_test "Tasks without numbers" "$input" 2
}

# Test 4: Simple number formats (currently fails)
test_simple_number_formats() {
    local input='## Phase 1: Setup
- [ ] 1 Task with simple number
- [ ] 2. Task with trailing dot
- [ ] 3.1 Task with decimal number'

    run_test "Simple number formats" "$input" 3
}

# Test 5: Mixed indentation levels with tabs (currently fails)
test_mixed_indentation() {
    local input='## Phase 1: Setup
- [ ] 1.1 Level 0
  - [ ] 1.2 Level 1 (2 spaces)
    - [ ] 1.3 Level 2 (4 spaces)
	- [ ] 1.4 Level 1 (tab)'

    run_test "Mixed indentation (spaces + tabs)" "$input" 4
}

# Test 6: Multiple sections with tasks
test_multiple_sections() {
    local input='## Phase 1: Setup
- [ ] 1.1 Setup task

## Phase 2: Implementation
- [ ] 2.1 Implementation task
- [ ] 2.2 Another impl task'

    run_test "Multiple sections" "$input" 3
}

# Test 7: Empty lines and extra whitespace
test_whitespace_handling() {
    local input='## Phase 1: Setup

- [ ] 1.1 Task with blank line above

- [ ] 1.2 Task with blank line below

## Phase 2: More

- [ ] 2.1 Another task'

    run_test "Whitespace and blank lines" "$input" 3
}

# Test 8: Real-world hierarchical structure (simplified from failing-tasks.md)
test_hierarchical_structure() {
    local input='## Phase 1: Add Consolidated Types
- [ ] 1.1 Create internal/errors/types.go
  - [ ] 1.2 Migrate ERROR_TYPE_SHOW_TO_AI
  - [ ] 1.3 Add domain-specific error types
  - [ ] 1.4 Add metadata key constants

## Phase 2: Migrate Internal Usage
- [ ] 2.1 Migrate internal/model/error_model.go
  - [ ] 2.2 Update all references
  - [ ] 2.3 Add deprecation notice'

    run_test "Real-world hierarchical structure" "$input" 7
}

# Test 9: Numbered sections (Phase 1:, Phase 2:)
test_numbered_sections() {
    local input='## 1. First Section
- [ ] 1.1 Task in first section

## 2. Second Section
- [ ] 2.1 Task in second section'

    run_test "Numbered sections" "$input" 2
}

# Test 10: V2 hierarchical output format verification
verify_v2_format() {
    local change_dir="$1"
    local tasks_file="$change_dir/tasks.jsonc"

    # Check version is 2 (strip JSONC comments first)
    local version=$(grep -v '^//' "$tasks_file" | jq -r '.version')
    if [ "$version" != "2" ]; then
        echo "Expected version 2, got $version"
        return 1
    fi

    # Check for includes field (strip JSONC comments first)
    if ! grep -v '^//' "$tasks_file" | jq -e '.includes' >/dev/null 2>&1; then
        echo "Missing includes field in root tasks.jsonc"
        return 1
    fi

    # Check that child files exist (strip JSONC comments first)
    local includes=$(grep -v '^//' "$tasks_file" | jq -r '.includes[]?')
    for child_file in $includes; do
        if [ ! -f "$change_dir/$child_file" ]; then
            echo "Referenced child file $child_file does not exist"
            return 1
        fi

        # Verify child file has parent field (strip JSONC comments first)
        if ! grep -v '^//' "$change_dir/$child_file" | jq -e '.parent' >/dev/null 2>&1; then
            echo "Child file $child_file missing parent field"
            return 1
        fi
    done

    return 0
}

test_v2_hierarchical_output() {
    local input='## Phase 1: Setup
- [ ] 1.1 First task
- [ ] 1.2 Second task

## Phase 2: Implementation
- [ ] 2.1 Impl task
- [ ] 2.2 Another task'

    run_test "V2 hierarchical output format" 4 verify_v2_format
}

# Test 11: Section 0 (preliminary tasks) inline in root
test_section_zero_inline() {
    local input='- [ ] Preliminary task (no section)

## Phase 1: Setup
- [ ] 1.1 Section task'

    # Verify function for section 0
    verify_section_zero() {
        local change_dir="$1"
        local tasks_file="$change_dir/tasks.jsonc"

        # Should have at least one task directly in root tasks array (strip JSONC comments first)
        local root_task_count=$(grep -v '^//' "$tasks_file" | jq '[.tasks[] | select(.children == null or (.children | type) != "string" or (.children | startswith("$ref:") | not))] | length')
        if [ "$root_task_count" -lt 1 ]; then
            echo "Expected at least 1 inline task in root (section 0), got $root_task_count"
            return 1
        fi

        return 0
    }

    run_test "Section 0 tasks inline in root" 2 verify_section_zero
}

# Test 12: Task status preservation
test_status_preservation() {
    local input='## Phase 1: Setup
- [x] 1.1 Completed task
- [ ] 1.2 Pending task'

    verify_status() {
        local change_dir="$1"
        local tasks_file="$change_dir/tasks.jsonc"

        # Check first task is completed (strip JSONC comments first)
        local status1=$(grep -v '^//' "$tasks_file" | jq -r '.tasks[0].status')
        if [ "$status1" != "completed" ]; then
            echo "Expected first task status 'completed', got '$status1'"
            return 1
        fi

        # Check second task is pending (strip JSONC comments first)
        local status2=$(grep -v '^//' "$tasks_file" | jq -r '.tasks[1].status')
        if [ "$status2" != "pending" ]; then
            echo "Expected second task status 'pending', got '$status2'"
            return 1
        fi

        return 0
    }

    run_test "Task status preservation" 2 verify_status
}

# Run all tests
echo "Running accept.sh test suite..."
echo "================================"
echo

echo "DEBUG: About to run test_basic_numbered"
test_basic_numbered
echo "DEBUG: test_basic_numbered completed"
test_indented_subtasks
test_tasks_without_numbers
test_simple_number_formats
test_mixed_indentation
test_multiple_sections
test_whitespace_handling
test_hierarchical_structure
test_numbered_sections
test_v2_hierarchical_output
test_section_zero_inline
test_status_preservation

# Summary
echo
echo "================================"
echo "Test Results:"
echo "  Total:  $TESTS_RUN"
echo -e "  ${GREEN}Passed: $TESTS_PASSED${NC}"
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "  ${RED}Failed: $TESTS_FAILED${NC}"
else
    echo "  Failed: $TESTS_FAILED"
fi
echo

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
