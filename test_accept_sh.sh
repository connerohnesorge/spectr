#!/usr/bin/env bash
#
# test_accept_sh.sh - Comprehensive test suite for accept.sh parsing
#
# This script tests various edge cases and task.md formatting scenarios
# that could cause parsing failures in accept.sh

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test configuration
REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"
SCRIPT_PATH="${1:-${REPO_ROOT}/internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh}"
TEST_DIR="${REPO_ROOT}/test_accept_output"
PASS_COUNT=0
FAIL_COUNT=0

# Cleanup and setup
cleanup() {
    rm -rf "$TEST_DIR"
}

setup() {
    cleanup
    mkdir -p "$TEST_DIR"
}

# Assert helpers
assert_file_exists() {
    local file="$1"
    local msg="${2:-File should exist}"
    if [ ! -f "$file" ]; then
        echo -e "${RED}FAIL${NC}: $msg ($file)"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        return 1
    fi
    return 0
}

assert_contains() {
    local file="$1"
    local pattern="$2"
    local msg="${3:-Should contain pattern}"
    if ! grep -q "$pattern" "$file"; then
        echo -e "${RED}FAIL${NC}: $msg"
        echo "File: $file"
        echo "Pattern: $pattern"
        echo "Content:"
        cat "$file"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        return 1
    fi
    return 0
}

assert_not_contains() {
    local file="$1"
    local pattern="$2"
    local msg="${3:-Should NOT contain pattern}"
    if grep -q "$pattern" "$file"; then
        echo -e "${RED}FAIL${NC}: $msg"
        echo "File: $file"
        echo "Pattern: $pattern"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        return 1
    fi
    return 0
}

assert_equals() {
    local expected="$1"
    local actual="$2"
    local msg="${3:-Values should match}"
    if [ "$expected" != "$actual" ]; then
        echo -e "${RED}FAIL${NC}: $msg"
        echo "Expected: $expected"
        echo "Actual:   $actual"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        return 1
    fi
    return 0
}

test_case() {
    local name="$1"
    echo -e "\n${YELLOW}Test: $name${NC}"
}

pass() {
    echo -e "${GREEN}PASS${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
}

# Test 1: Basic Format 1 (## for sections)
test_case "Format 1: Basic ## section headers"
setup
mkdir -p "$TEST_DIR/test-format1/spectr/changes/test-format1"
cat > "$TEST_DIR/test-format1/spectr/changes/test-format1/tasks.md" << 'EOF'
# Implementation

## 1. Setup

- [ ] 1.1 Initialize project
- [x] 1.2 Install dependencies

## 2. Development

- [ ] 2.1 Write core logic
- [ ] 2.2 Add tests
EOF

cd "$TEST_DIR/test-format1"
if bash "$SCRIPT_PATH" test-format1; then
    assert_file_exists "spectr/changes/test-format1/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-format1/tasks.jsonc" '"id": "1.1"' "Should parse task 1.1"
    assert_contains "spectr/changes/test-format1/tasks.jsonc" '"id": "1.2"' "Should parse task 1.2"
    assert_contains "spectr/changes/test-format1/tasks.jsonc" '"status": "pending"' "Should have pending status"
    assert_contains "spectr/changes/test-format1/tasks.jsonc" '"status": "completed"' "Should have completed status"
    assert_contains "spectr/changes/test-format1/tasks.jsonc" '"section": "Setup"' "Should extract Setup section"
    assert_contains "spectr/changes/test-format1/tasks.jsonc" '"section": "Development"' "Should extract Development section"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 2: Format 2 (# for sections, ## Tasks)
test_case "Format 2: # section with ## Tasks subheader"
setup
mkdir -p "$TEST_DIR/test-format2/spectr/changes/test-format2"
cat > "$TEST_DIR/test-format2/spectr/changes/test-format2/tasks.md" << 'EOF'
# 1. Backend

## Tasks

- [ ] 1.1 Create API endpoint
- [ ] 1.2 Add database migration

# 2. Frontend

## Tasks

- [ ] 2.1 Build UI component
- [ ] 2.2 Wire up API calls
EOF

cd "$TEST_DIR/test-format2"
if bash "$SCRIPT_PATH" test-format2; then
    assert_file_exists "spectr/changes/test-format2/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-format2/tasks.jsonc" '"section": "Backend"' "Should extract Backend section"
    assert_contains "spectr/changes/test-format2/tasks.jsonc" '"section": "Frontend"' "Should extract Frontend section"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 3: Special characters in descriptions
test_case "Special characters in task descriptions"
setup
mkdir -p "$TEST_DIR/test-special/spectr/changes/test-special"
cat > "$TEST_DIR/test-special/spectr/changes/test-special/tasks.md" << 'EOF'
# Implementation

## 1. Phase 1

- [ ] 1.1 Handle "quoted" text properly
- [ ] 1.2 Support backslash \ character
- [ ] 1.3 Test: foo/bar/baz paths
- [ ] 1.4 Handle JSON with {nested: "objects"}
EOF

cd "$TEST_DIR/test-special"
if bash "$SCRIPT_PATH" test-special; then
    assert_file_exists "spectr/changes/test-special/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-special/tasks.jsonc" 'quoted' "Should preserve quoted text"
    assert_contains "spectr/changes/test-special/tasks.jsonc" 'backslash' "Should preserve backslash reference"
    assert_contains "spectr/changes/test-special/tasks.jsonc" 'paths' "Should preserve path"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 4: Task ID variations
test_case "Various task ID formats (X.Y)"
setup
mkdir -p "$TEST_DIR/test-ids/spectr/changes/test-ids"
cat > "$TEST_DIR/test-ids/spectr/changes/test-ids/tasks.md" << 'EOF'
# Implementation

## Section

- [ ] 1.0 First level zero
- [ ] 1.1 Normal task
- [ ] 10.1 Double digit section
- [ ] 3.99 Large subtask number
EOF

cd "$TEST_DIR/test-ids"
if bash "$SCRIPT_PATH" test-ids; then
    assert_file_exists "spectr/changes/test-ids/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-ids/tasks.jsonc" '"id": "1.0"' "Should parse task ID 1.0"
    assert_contains "spectr/changes/test-ids/tasks.jsonc" '"id": "10.1"' "Should parse task ID 10.1"
    assert_contains "spectr/changes/test-ids/tasks.jsonc" '"id": "3.99"' "Should parse task ID 3.99"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 5: Empty sections and mixed content
test_case "Empty sections and non-task lines"
setup
mkdir -p "$TEST_DIR/test-mixed/spectr/changes/test-mixed"
cat > "$TEST_DIR/test-mixed/spectr/changes/test-mixed/tasks.md" << 'EOF'
# Implementation

## Section 1

Some explanation text that isn't a task.

- [ ] 1.1 Actual task

More filler text

- [ ] 1.2 Another task

## Section 2

- [ ] 2.1 First task in section 2

EOF

cd "$TEST_DIR/test-mixed"
if bash "$SCRIPT_PATH" test-mixed; then
    assert_file_exists "spectr/changes/test-mixed/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-mixed/tasks.jsonc" '"id": "1.1"' "Should parse task 1.1"
    assert_contains "spectr/changes/test-mixed/tasks.jsonc" '"id": "1.2"' "Should parse task 1.2"
    assert_contains "spectr/changes/test-mixed/tasks.jsonc" '"id": "2.1"' "Should parse task 2.1"
    assert_not_contains "spectr/changes/test-mixed/tasks.jsonc" "explanation" "Should not include description text"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 6: Whitespace variations
test_case "Whitespace handling (multiple spaces, tabs)"
setup
mkdir -p "$TEST_DIR/test-whitespace/spectr/changes/test-whitespace"
cat > "$TEST_DIR/test-whitespace/spectr/changes/test-whitespace/tasks.md" << 'EOF'
#   Multiple spaces in header

##  1.   Section with spaces

- [ ] 1.1  Task with extra spaces
- [x]  1.2  More space variations
EOF

cd "$TEST_DIR/test-whitespace"
if bash "$SCRIPT_PATH" test-whitespace; then
    assert_file_exists "spectr/changes/test-whitespace/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-whitespace/tasks.jsonc" '"id": "1.1"' "Should parse despite extra spaces"
    assert_contains "spectr/changes/test-whitespace/tasks.jsonc" '"id": "1.2"' "Should parse 1.2"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 7: Section transitions and empty lines
test_case "Section transitions with blank lines"
setup
mkdir -p "$TEST_DIR/test-transitions/spectr/changes/test-transitions"
cat > "$TEST_DIR/test-transitions/spectr/changes/test-transitions/tasks.md" << 'EOF'
# Implementation

## Phase 1

- [ ] 1.1 First task



- [ ] 1.2 Second task

## Phase 2

- [ ] 2.1 Task in phase 2
EOF

cd "$TEST_DIR/test-transitions"
if bash "$SCRIPT_PATH" test-transitions; then
    assert_file_exists "spectr/changes/test-transitions/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-transitions/tasks.jsonc" '"section": "Phase 1"' "Should have Phase 1"
    assert_contains "spectr/changes/test-transitions/tasks.jsonc" '"section": "Phase 2"' "Should have Phase 2"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 8: Missing section (no header before tasks)
test_case "Tasks without preceding section header"
setup
mkdir -p "$TEST_DIR/test-nosection/spectr/changes/test-nosection"
cat > "$TEST_DIR/test-nosection/spectr/changes/test-nosection/tasks.md" << 'EOF'
- [ ] 1.1 Orphan task without section
- [ ] 1.2 Another orphan task

# Implementation

## Section

- [ ] 2.1 Task with section
EOF

cd "$TEST_DIR/test-nosection"
if bash "$SCRIPT_PATH" test-nosection; then
    assert_file_exists "spectr/changes/test-nosection/tasks.jsonc" "JSONC should be generated"
    assert_contains "spectr/changes/test-nosection/tasks.jsonc" '"id": "1.1"' "Should parse orphan tasks"
    assert_contains "spectr/changes/test-nosection/tasks.jsonc" '"section": ""' "Orphan should have empty section or default"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 9: Real-world complex format
test_case "Real-world complex tasks.md"
setup
mkdir -p "$TEST_DIR/test-complex/spectr/changes/test-complex"
cat > "$TEST_DIR/test-complex/spectr/changes/test-complex/tasks.md" << 'EOF'
# Change: Add User Authentication

## 1. Backend Setup

- [ ] 1.1 Create user database schema with migrations
- [x] 1.2 Implement JWT token generation and validation
- [ ] 1.3 Add password hashing (bcrypt) implementation

## 2. API Endpoints

- [ ] 2.1 Create POST /auth/login endpoint
- [ ] 2.2 Create POST /auth/signup endpoint  
- [ ] 2.3 Create POST /auth/logout endpoint
- [ ] 2.4 Add middleware for protected routes

## 3. Frontend Implementation

- [ ] 3.1 Build login form component
- [ ] 3.2 Add form validation and error messages
- [ ] 3.3 Implement logout button in navbar
- [x] 3.4 Style authentication pages with CSS

## 4. Testing

- [ ] 4.1 Write unit tests for JWT validation
- [ ] 4.2 Write integration tests for auth endpoints
- [ ] 4.3 Write E2E tests for login flow
EOF

cd "$TEST_DIR/test-complex"
if bash "$SCRIPT_PATH" test-complex; then
    assert_file_exists "spectr/changes/test-complex/tasks.jsonc" "JSONC should be generated"
    # Check task counts
    local task_count=$(grep -c '"id"' "spectr/changes/test-complex/tasks.jsonc")
    echo "Total tasks parsed: $task_count"
    if [ "$task_count" -ge 12 ]; then
        echo -e "${GREEN}PASS${NC}: All tasks parsed ($task_count found)"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}FAIL${NC}: Expected 13 tasks, found $task_count"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 10: Invalid checkbox formats (should NOT match)
test_case "Invalid task formats are ignored"
setup
mkdir -p "$TEST_DIR/test-invalid/spectr/changes/test-invalid"
cat > "$TEST_DIR/test-invalid/spectr/changes/test-invalid/tasks.md" << 'EOF'
# Implementation

## Section

- [ ] 1.1 Valid task
- [*] 1.2 Invalid checkbox (should not parse)
- () 1.3 Wrong bracket type (should not parse)
* [ ] 1.4 Wrong bullet style (should not parse)
- [ ] 1.5 Another valid task
EOF

cd "$TEST_DIR/test-invalid"
if bash "$SCRIPT_PATH" test-invalid; then
    assert_file_exists "spectr/changes/test-invalid/tasks.jsonc" "JSONC should be generated"
    local task_count=$(grep -c '"id"' "spectr/changes/test-invalid/tasks.jsonc")
    if [ "$task_count" -eq 2 ]; then
        echo -e "${GREEN}PASS${NC}: Correctly ignored invalid formats (2 tasks parsed)"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}FAIL${NC}: Expected 2 valid tasks, found $task_count"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 11: Task ID without section
test_case "Orphaned task ID variations"
setup
mkdir -p "$TEST_DIR/test-orphan/spectr/changes/test-orphan"
cat > "$TEST_DIR/test-orphan/spectr/changes/test-orphan/tasks.md" << 'EOF'
- [ ] 1.1 First orphan
- [ ] 1.2 Second orphan

# Section A

## Subsection

- [ ] 2.1 With section
EOF

cd "$TEST_DIR/test-orphan"
if bash "$SCRIPT_PATH" test-orphan; then
    assert_file_exists "spectr/changes/test-orphan/tasks.jsonc" "JSONC should be generated"
    pass
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Test 12: JSON validity of output
test_case "Generated JSONC is valid JSON (after comment removal)"
setup
mkdir -p "$TEST_DIR/test-json/spectr/changes/test-json"
cat > "$TEST_DIR/test-json/spectr/changes/test-json/tasks.md" << 'EOF'
# Implementation

## Tasks

- [ ] 1.1 Test JSON validity
- [x] 1.2 Check parsing
EOF

cd "$TEST_DIR/test-json"
if bash "$SCRIPT_PATH" test-json; then
    assert_file_exists "spectr/changes/test-json/tasks.jsonc" "JSONC file should exist"
    # Remove comments and validate JSON
    if grep -v '^[[:space:]]*\/\/' "spectr/changes/test-json/tasks.jsonc" | jq . > /dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}: Generated JSONC is valid JSON"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}FAIL${NC}: Generated JSONC is not valid JSON"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
else
    echo -e "${RED}FAIL${NC}: Script failed to run"
    FAIL_COUNT=$((FAIL_COUNT + 1))
fi
cd - > /dev/null

# Summary
echo -e "\n${YELLOW}=====================================${NC}"
echo -e "${YELLOW}Test Summary${NC}"
echo -e "${YELLOW}=====================================${NC}"
echo -e "Passed: ${GREEN}$PASS_COUNT${NC}"
echo -e "Failed: ${RED}$FAIL_COUNT${NC}"

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed!${NC}"
    cleanup
    exit 0
else
    echo -e "\n${RED}Some tests failed. Test output preserved in $TEST_DIR${NC}"
    exit 1
fi
