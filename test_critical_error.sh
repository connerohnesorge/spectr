#!/usr/bin/env bash
#
# test_critical_error.sh - Isolated reproduction of critical parsing error
#
# This script clearly demonstrates the whitespace parsing bug in accept.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")" && pwd)"
SCRIPT_PATH="${REPO_ROOT}/internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh"
TEST_DIR="${REPO_ROOT}/test_critical_demo"

echo "=========================================="
echo "CRITICAL ERROR DEMONSTRATION"
echo "=========================================="
echo ""
echo "Script: $SCRIPT_PATH"
echo "Test Dir: $TEST_DIR"
echo ""

# Cleanup and setup
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR/spectr/changes/demo"

# Create task file with the problematic whitespace pattern
cat > "$TEST_DIR/spectr/changes/demo/tasks.md" << 'EOF'
# Implementation

## Phase

- [ ] 1.1 Correct spacing (single space after bracket)
- [ ]  1.2 Extra spaces (double space - WILL FAIL TO PARSE!)
- [x] 2.1 Completed task with correct spacing
- [x]  2.2 Completed task with extra spaces (WILL FAIL!)
EOF

echo "Input tasks.md:"
echo "==============="
cat "$TEST_DIR/spectr/changes/demo/tasks.md"
echo ""

# Run the script
cd "$TEST_DIR"
bash "$SCRIPT_PATH" demo 2>&1

echo ""
echo "Generated tasks.jsonc:"
echo "======================"
cat "spectr/changes/demo/tasks.jsonc" | grep -A 100 '"tasks"' | head -40

echo ""
echo "Analysis:"
echo "========="

# Count tasks that were parsed
TASK_COUNT=$(grep -c '"id"' "spectr/changes/demo/tasks.jsonc")
echo "Tasks that were PARSED: $TASK_COUNT"
echo "Tasks EXPECTED: 4 (1.1, 1.2, 2.1, 2.2)"
echo ""

if grep -q '"id": "1.1"' "spectr/changes/demo/tasks.jsonc"; then
    echo "✓ Task 1.1 parsed (correct spacing)"
else
    echo "✗ Task 1.1 NOT parsed"
fi

if grep -q '"id": "1.2"' "spectr/changes/demo/tasks.jsonc"; then
    echo "✓ Task 1.2 parsed"
else
    echo "✗ Task 1.2 FAILED TO PARSE (extra spaces!) - CRITICAL BUG"
fi

if grep -q '"id": "2.1"' "spectr/changes/demo/tasks.jsonc"; then
    echo "✓ Task 2.1 parsed"
else
    echo "✗ Task 2.1 NOT parsed"
fi

if grep -q '"id": "2.2"' "spectr/changes/demo/tasks.jsonc"; then
    echo "✓ Task 2.2 parsed"
else
    echo "✗ Task 2.2 FAILED TO PARSE (extra spaces!) - CRITICAL BUG"
fi

echo ""
echo "ROOT CAUSE:"
echo "==========="
echo "The regex at line 113 uses [[:space:]] which matches exactly ONE space."
echo "When there are multiple spaces (common in markdown), the pattern fails."
echo ""
echo "Current pattern:"
echo "  ^-[[:space:]]\\[([[:space:]]|x)\\][[:space:]]([0-9]+\\.[0-9]+)[[:space:]](.+)$"
echo "            ↑                    ↑                          ↑"
echo "            Exactly 1 space     Exactly 1 space          Exactly 1 space"
echo ""
echo "Should be:"
echo "  ^-[[:space:]]*\\[([[:space:]]|x)\\][[:space:]]*([0-9]+\\.[0-9]+)[[:space:]]*(.+)$"
echo "            ↑                    ↑                          ↑"
echo "           0+ spaces            0+ spaces                0+ spaces"
echo ""
echo "This silently drops tasks from the output - a data loss bug!"
