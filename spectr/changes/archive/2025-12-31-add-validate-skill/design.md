# Design: spectr-validate-wo-spectr-bin AgentSkills Skill

## Context

### Problem Statement

The spectr CLI binary implements comprehensive validation for specification files
and change proposals through the `spectr validate` command. This validation
enforces critical quality standards:

- Structural correctness (required sections, proper markdown formatting)
- Semantic correctness (normative language, scenario coverage)
- Delta operation validity (proper use of ADDED/MODIFIED/REMOVED/RENAMED)
- Cross-file consistency (no duplicate requirements, valid base spec references)

However, this validation is unavailable in environments where the Go binary
cannot be installed or executed:

1. **Sandboxed AI Coding Environments**: Claude Code, Cursor, and similar tools
   run code in isolated execution contexts with restricted binary execution
2. **CI/CD Pipelines**: Automated workflows may run on minimal container images
   without Go runtime or build toolchain
3. **Developer Workstations**: Corporate environments with restricted package
   installation permissions

The existing `spectr-accept-wo-spectr-bin` skill established the solution
pattern: implement equivalent functionality using portable shell scripts with
standard Unix tools. That skill converts `tasks.md` to `tasks.jsonc` using bash,
sed, and jq.

### Current State

**Go Implementation** (`internal/validation/`):

The validation system is organized across multiple packages:

```
internal/validation/
├── validator.go           # Orchestrator (ValidateSpec, ValidateChange)
├── spec_rules.go          # Spec file validation rules
├── change_rules.go        # Change delta validation rules
├── delta_validators.go    # Delta-specific validators
├── parser.go              # Markdown parsing utilities
├── types.go               # ValidationReport, ValidationIssue
├── formatters.go          # Human/JSON output formatting
├── helpers.go             # Type determination, normalization
├── items.go               # Discovery (GetAllItems, GetSpecItems)
└── interactive.go         # Bubbletea TUI (not needed for skill)
```

**Key Validation Rules** (from Go implementation):

| Rule | Level | Source | Line |
|------|-------|--------|------|
| Missing `## Requirements` section | ERROR | spec_rules.go | 34-41 |
| Requirement without SHALL/MUST | ERROR (strict) | spec_rules.go | 109-117 |
| Requirement without scenario | ERROR (strict) | spec_rules.go | 120-127 |
| Malformed scenario format | ERROR | spec_rules.go | 130-146 |
| Empty ADDED section | ERROR | delta_validators.go | 64-74 |
| Empty MODIFIED section | ERROR | delta_validators.go | 195-205 |
| Empty REMOVED section | ERROR | delta_validators.go | 328-337 |
| Empty RENAMED section | ERROR | delta_validators.go | 414-423 |
| ADDED req without scenario | ERROR | delta_validators.go | 105-115 |
| ADDED req without SHALL/MUST | ERROR | delta_validators.go | 92-102 |
| MODIFIED req without scenario | ERROR | delta_validators.go | 236-247 |
| MODIFIED req without SHALL/MUST | ERROR | delta_validators.go | 223-233 |
| Change with zero deltas | ERROR | change_rules.go | 160-170 |
| tasks.md with no task items | ERROR | change_rules.go | 797-807 |

**Markdown Parsing** (`internal/markdown/`):

The Go implementation uses dedicated markdown parsing functions:

- `MatchRequirementHeader(line)` - Matches `### Requirement: Name`
- `MatchScenarioHeader(line)` - Matches `#### Scenario: Name`
- `MatchDeltaSection(line)` - Matches `## ADDED/MODIFIED/REMOVED/RENAMED Requirements`
- `MatchFlexibleTask(line)` - Matches `- [ ]` or `- [x]` with optional task IDs
- `MatchRenamedFrom(line)` - Matches `- FROM: ### Requirement: OldName`
- `MatchRenamedTo(line)` - Matches `- TO: ### Requirement: NewName`

**Discovery** (`internal/discovery/`):

Discovery functions enumerate specs and changes:

```go
// GetSpecIDs returns all spec IDs (directories under spectr/specs/ with spec.md)
func GetSpecIDs(spectrRoot string) ([]string, error)

// GetChangeIDs returns all change IDs (directories under spectr/changes/ except archive)
func GetChangeIDs(spectrRoot string) ([]string, error)
```

### Stakeholders

**Primary Users**:
- **AI Coding Assistants** (Claude Code, Cursor, Windsurf, etc.) - Need real-time
  validation during spec authoring
- **CI/CD Systems** - Automated validation in pipelines without binary
- **Developers** - Validation in restricted corporate environments

**Secondary Stakeholders**:
- **Spectr Maintainers** - Must ensure skill behavior stays synchronized with
  binary implementation
- **AgentSkills Ecosystem** - Skill demonstrates AgentSkills specification
  compliance

### Constraints

**Technical Constraints**:

1. **No External Dependencies**: Cannot require language runtimes (Python, Ruby,
   Node.js) beyond bash
2. **Standard Unix Tools Only**: Must work with tools universally available
   (grep, sed, awk, find)
3. **Optional jq**: JSON output can require jq but must degrade gracefully if
   unavailable
4. **Bash 4.0+**: Can use modern bash features (associative arrays, `[[` tests)
5. **No Pre-merge Validation**: Cannot implement complex delta-to-base-spec
   validation (requires parsing and comparing two specs)
6. **Sequential Processing**: No parallelization (Go uses 6-worker pool)

**Performance Constraints**:

- Must validate 30-50 spec files in reasonable time (<10 seconds)
- Line-by-line parsing acceptable (bash `read` loop is efficient enough)
- No requirement to match Go binary's parallel validation speed

**Compatibility Constraints**:

- Must work on Linux (GNU coreutils)
- Must work on macOS (BSD coreutils with slight differences)
- Regex patterns must be POSIX-compatible or use bash `[[ =~ ]]` extension

**Maintenance Constraints**:

- Regex patterns must stay synchronized with `internal/markdown/` matchers
- Validation rules must stay synchronized with `internal/validation/` logic
- Changes to Go implementation require corresponding skill updates

## Goals / Non-Goals

### Goals

**Primary Goals**:

1. **Structural Validation**: Verify spec files have required sections and proper
   markdown structure
2. **Requirement Validation**: Check requirements contain normative language
   (SHALL/MUST) and have scenario coverage
3. **Delta Validation**: Validate change proposals have proper delta operations
   and non-empty sections
4. **Task Validation**: Ensure tasks.md files contain at least one task item if
   present
5. **Discovery**: Find all specs and changes for bulk validation
6. **Human Output**: Provide readable error messages with file paths and line
   numbers
7. **JSON Output**: Support programmatic consumption with `--json` flag
8. **Exit Codes**: Return appropriate codes for CI integration

**Secondary Goals**:

1. **Color Output**: Color-code error levels in TTY environments
2. **Flexible Invocation**: Support single-item and bulk validation modes
3. **Environment Customization**: Allow custom spectr directory via environment
   variable
4. **Helpful Errors**: Provide actionable error messages with format examples

### Non-Goals

**Explicitly Out of Scope**:

1. **Pre-merge Validation Against Base Specs**: The Go implementation validates
   that MODIFIED/REMOVED/RENAMED requirements exist in base spec and ADDED
   requirements don't. This requires:
   - Parsing base spec to extract requirement names
   - Normalizing requirement names for comparison
   - Cross-referencing delta operations
   - Handling capability-scoped requirement names

   **Rationale**: Too complex for shell script. Would require sophisticated state
   management, parsing two files, and normalization logic. The value/complexity
   ratio doesn't justify implementation.

2. **Cross-Capability Duplicate Detection**: The Go implementation tracks
   requirement names across multiple delta files using composite keys
   (`capability::normalized_name`). This requires maintaining global state across
   file processing.

   **Rationale**: Complex state management in bash. The skill validates per-file
   duplicates but not cross-file duplicates.

3. **Interactive TUI Mode**: The Go implementation provides a Bubbletea-based
   interactive menu (`spectr validate` with no args in TTY).

   **Rationale**: Out of scope for shell script. Users invoke skill with explicit
   arguments.

4. **Concurrent Validation**: The Go implementation uses a 6-worker pool for
   parallel validation.

   **Rationale**: Bash doesn't support true multithreading. Sequential processing
   is acceptable for skill use case.

5. **Strict vs Non-Strict Mode**: The Go implementation originally had a strict
   mode flag but now always validates strictly (warnings → errors).

   **Rationale**: Skill matches current Go behavior (strict only).

6. **Line Number Accuracy for Pre-merge Errors**: When pre-merge validation
   fails, the Go implementation uses `findPreMergeErrorLine()` to locate the
   exact line of the problematic requirement.

   **Rationale**: Not implementing pre-merge validation, so not needed.

7. **Cross-File Conflict Detection Within Changes**: Detecting when a requirement
   appears in both ADDED and MODIFIED sections across different files in a
   change.

   **Rationale**: Complex cross-file state tracking. Per-file validation is
   sufficient.

## Decisions

### Decision 1: Script Architecture - Single Monolithic Script

**What**: Implement all validation logic in a single `validate.sh` script
(~500-600 lines) rather than splitting into multiple scripts or using a modular
approach.

**Why**:
- **Simpler Installation**: Single file to copy, single file to maintain
- **No PATH Issues**: No need to manage multiple script locations or set up PATH
- **Easier Debugging**: All logic in one place for troubleshooting
- **Follows Existing Pattern**: `spectr-accept-wo-spectr-bin` uses single-script
  approach
- **Function Modularization**: Can still organize with bash functions for
  readability

**Alternatives Considered**:

**Alternative 1: Multiple Scripts**
```
scripts/
├── validate-spec.sh
├── validate-change.sh
├── validate-tasks.sh
├── discover.sh
├── format-output.sh
└── validate.sh (orchestrator)
```
**Rejected Because**:
- Complexity in passing state between scripts
- Need to source scripts or manage execution dependencies
- Harder for users to invoke (which script to call?)
- More files to maintain and synchronize

**Alternative 2: Shared Library Approach**
```
scripts/
├── lib/
│   ├── parsers.sh
│   ├── validators.sh
│   └── formatters.sh
└── validate.sh (sources lib/)
```
**Rejected Because**:
- Overkill for ~500 lines of code
- Makes testing harder (need to source library)
- Complicates deployment (must copy multiple files)

**Trade-offs**:
- ✅ Simpler to install and use
- ✅ Easier to understand (all logic visible)
- ✅ No sourcing or PATH management
- ❌ Longer file (but still manageable at 500-600 lines)
- ❌ Requires function organization discipline

### Decision 2: Parsing Strategy - Line-by-Line State Machine

**What**: Process files line-by-line using bash `read` loop, maintaining state
variables for current section, current requirement, etc.

**Implementation Pattern**:
```bash
parse_spec_file() {
    local file="$1"
    local line_num=0
    local in_requirements=false
    local current_requirement=""
    local requirement_line=0
    local has_scenario=false
    local has_shall_must=false
    local found_requirements_section=false

    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))

        # State transitions on section boundaries
        if [[ "$line" =~ ^##[[:space:]](.+)$ ]]; then
            flush_current_requirement  # Validate before transitioning
            section_name="${BASH_REMATCH[1]}"

            if [[ "$section_name" == "Requirements" ]]; then
                found_requirements_section=true
                in_requirements=true
            else
                in_requirements=false
            fi
            continue
        fi

        # Only process if in Requirements section
        if ! $in_requirements; then
            continue
        fi

        # State transitions on requirement boundaries
        if [[ "$line" =~ ^###[[:space:]]Requirement:[[:space:]](.+)$ ]]; then
            flush_current_requirement  # Validate previous requirement
            current_requirement="${BASH_REMATCH[1]}"
            requirement_line=$line_num
            has_scenario=false
            has_shall_must=false
            continue
        fi

        # Update state based on content
        if [[ "$line" =~ ^####[[:space:]]Scenario: ]]; then
            has_scenario=true
        fi

        if [[ "$line" =~ (SHALL|MUST) ]]; then
            has_shall_must=true
        fi

        # Error detection (malformed scenarios)
        detect_malformed_scenario "$file" "$line_num" "$line"
    done < "$file"

    # Flush final requirement
    flush_current_requirement

    # Validate global file state
    if ! $found_requirements_section; then
        add_issue "ERROR" "$file" 1 "Missing required '## Requirements' section"
    fi
}

flush_current_requirement() {
    if [[ -z "$current_requirement" ]]; then
        return
    fi

    if [[ "$has_shall_must" != "true" ]]; then
        add_issue "ERROR" "$file: Requirement '$current_requirement'" \
            "$requirement_line" \
            "Requirement should contain SHALL or MUST to indicate normative requirement"
    fi

    if [[ "$has_scenario" != "true" ]]; then
        add_issue "ERROR" "$file: Requirement '$current_requirement'" \
            "$requirement_line" \
            "Requirement should have at least one scenario"
    fi

    # Reset state
    current_requirement=""
}
```

**Why This Approach**:

1. **Precise Line Numbers**: Tracking `line_num` during read loop gives exact
   error locations
2. **Efficient Memory Usage**: Process line-by-line without loading entire file
3. **Clear State Management**: State variables make parsing logic explicit
4. **Matches Go Pattern**: Go implementation uses similar section/requirement
   tracking
5. **Handles Edge Cases**: Can detect and flush state at boundaries (EOF, section
   changes)

**Alternatives Considered**:

**Alternative 1: Grep-Based Extraction**
```bash
# Extract all requirements
requirements=$(grep -n "^### Requirement:" "$file")

# For each requirement, check for scenarios
for req in $requirements; do
    line_num=$(echo "$req" | cut -d: -f1)
    # Try to find scenario after this line...
done
```
**Rejected Because**:
- Hard to determine requirement boundaries (where does one req end?)
- Difficult to associate scenarios with correct requirement
- Line number arithmetic complex and error-prone
- Can't easily detect malformed scenarios within requirement context

**Alternative 2: Awk Script**
```bash
awk '
/^## Requirements/ { in_req = 1; next }
/^## / && in_req { in_req = 0 }
in_req && /^### Requirement:/ {
    if (current_req != "") validate_requirement()
    current_req = $0
    req_line = NR
    has_scenario = 0
    has_shall = 0
}
in_req && /^#### Scenario:/ { has_scenario = 1 }
in_req && /SHALL|MUST/ { has_shall = 1 }
END { validate_requirement() }
' "$file"
```
**Rejected Because**:
- Awk syntax less familiar to maintainers than bash
- Error reporting from awk script more complex
- Issue collection would need awk arrays or external file
- Harder to integrate with rest of bash script

**Trade-offs**:
- ✅ Accurate line numbers for errors
- ✅ Handles complex state transitions cleanly
- ✅ Efficient (doesn't load entire file into memory)
- ✅ Familiar bash idioms (while read, if statements)
- ❌ Longer code than simple grep (but more robust)
- ❌ Requires careful state management discipline

### Decision 3: Validation Rules - Exact Parity with Go Implementation

**What**: Implement exactly the same validation rules as the Go implementation,
treating all warnings as errors (strict mode only).

**Rule Mapping Table**:

| Rule ID | Description | Go Source | Skill Implementation |
|---------|-------------|-----------|---------------------|
| SPEC-1 | Missing `## Requirements` section | spec_rules.go:34-41 | Check `found_requirements_section` at end |
| SPEC-2 | Requirement without SHALL/MUST | spec_rules.go:109-117 | Check `has_shall_must` in flush |
| SPEC-3 | Requirement without scenario | spec_rules.go:120-127 | Check `has_scenario` in flush |
| SPEC-4 | Malformed scenario format | spec_rules.go:130-146 | Regex check for wrong patterns |
| DELTA-1 | Empty ADDED section | delta_validators.go:64-74 | Track `section_has_requirements` |
| DELTA-2 | Empty MODIFIED section | delta_validators.go:195-205 | Track `section_has_requirements` |
| DELTA-3 | Empty REMOVED section | delta_validators.go:328-337 | Track `section_has_requirements` |
| DELTA-4 | Empty RENAMED section | delta_validators.go:414-423 | Track `section_has_requirements` |
| DELTA-5 | ADDED req without scenario | delta_validators.go:105-115 | Check in delta flush |
| DELTA-6 | ADDED req without SHALL/MUST | delta_validators.go:92-102 | Check in delta flush |
| DELTA-7 | MODIFIED req without scenario | delta_validators.go:236-247 | Check in delta flush |
| DELTA-8 | MODIFIED req without SHALL/MUST | delta_validators.go:223-233 | Check in delta flush |
| DELTA-9 | Change with zero deltas | change_rules.go:160-170 | Check `total_deltas` after all files |
| TASK-1 | tasks.md with no task items | change_rules.go:797-807 | Count tasks, error if zero |

**Why Exact Parity**:

1. **User Expectation**: Users expect same validation results whether using
   binary or skill
2. **CI Integration**: CI pipelines should get identical pass/fail results
3. **Documentation Reuse**: Can reference same validation documentation
4. **Testing Simplification**: Can compare skill output vs binary output for
   correctness

**Alternatives Considered**:

**Alternative 1: Subset of Rules**
- Implement only structural validation (sections exist)
- Skip SHALL/MUST checks, scenario checks
- **Rejected**: Would provide incomplete validation, confusing users about what's
  validated

**Alternative 2: Extended Validation**
- Add bash-specific checks (shellcheck integration, etc.)
- **Rejected**: Scope creep, not relevant to Spectr validation

**Trade-offs**:
- ✅ Consistent user experience
- ✅ Easy to verify correctness (compare outputs)
- ✅ Clear behavior specification
- ❌ Must maintain synchronization with Go code
- ❌ Cannot leverage additional bash capabilities

### Decision 4: Error Output Format - Match Go Human-Readable Format

**What**: Match the exact output format from
`internal/validation/formatters.go:FormatBulkHumanResults()`.

**Format Specification**:

```
<relative-path>
  [<LEVEL>] line <number>: <message>
  [<LEVEL>] line <number>: <message>

<relative-path>
  [<LEVEL>] line <number>: <message>

Summary: <passed> passed, <failed> failed (<errors> errors), <total> total
```

**Example Output**:
```
changes/add-feature/specs/auth/spec.md
  [ERROR] line 15: ADDED requirement must have at least one scenario
  [ERROR] line 23: ADDED requirement must contain SHALL or MUST

specs/validation/spec.md
  [ERROR] line 1: Missing required '## Requirements' section

Summary: 2 passed, 2 failed (3 errors), 4 total
```

**Implementation Details**:

```bash
print_human_results() {
    local total_items=0
    local passed_items=0
    local failed_items=0
    local total_errors=${ISSUE_COUNTS[errors]}

    # Group issues by file
    declare -A issues_by_file
    for issue in "${ISSUES[@]}"; do
        IFS='|' read -r level path line message <<< "$issue"
        issues_by_file["$path"]+="  [${level}] line ${line}: ${message}
"
    done

    # Print grouped issues
    for file in "${!issues_by_file[@]}"; do
        echo "$file"
        echo -n "${issues_by_file[$file]}"
        echo  # Blank line between files
        ((failed_items++))
    done

    # Calculate passed items
    passed_items=$((total_items - failed_items))

    # Print summary
    if [[ $failed_items -gt 0 ]]; then
        echo "Summary: $passed_items passed, $failed_items failed ($total_errors errors), $total_items total"
    else
        echo "Summary: $total_items passed, 0 failed, $total_items total"
    fi
}
```

**Why This Format**:

1. **Consistency**: Same format users see from binary
2. **Readability**: Clear file grouping, indented issues
3. **Copy-Paste Friendly**: File paths in Go-friendly format (file:line syntax)
4. **Summary**: Quick overview of validation results

**Color Enhancement**:

```bash
setup_colors() {
    if [[ -t 1 ]]; then  # stdout is TTY
        RED='\033[0;31m'
        YELLOW='\033[0;33m'
        GREEN='\033[0;32m'
        NC='\033[0m'  # No Color
    else
        RED=''
        YELLOW=''
        GREEN=''
        NC=''
    fi
}

format_error_level() {
    local level="$1"
    case "$level" in
        ERROR)
            echo -e "${RED}[ERROR]${NC}"
            ;;
        WARNING)
            echo -e "${YELLOW}[WARNING]${NC}"
            ;;
        *)
            echo "[$level]"
            ;;
    esac
}
```

**Relative Path Handling**:

```bash
make_relative_path() {
    local full_path="$1"
    local spectr_dir="${SPECTR_DIR:-spectr}"

    # Remove spectr/ prefix if present
    if [[ "$full_path" == "$spectr_dir/"* ]]; then
        echo "${full_path#$spectr_dir/}"
    else
        echo "$full_path"
    fi
}
```

**Trade-offs**:
- ✅ Consistent with binary output
- ✅ Color-coded for readability (when in TTY)
- ✅ Clear file grouping
- ❌ Requires issue grouping logic
- ❌ More complex than simple line-by-line printing

### Decision 5: JSON Output Structure - Match Go JSON Schema

**What**: Produce JSON output matching the structure from
`internal/validation/formatters.go:FormatBulkJSONResults()` when `--json` flag
is provided.

**JSON Schema**:

```json
{
  "version": 1,
  "items": [
    {
      "name": "add-feature",
      "type": "change",
      "valid": false,
      "issues": [
        {
          "level": "ERROR",
          "path": "changes/add-feature/specs/auth/spec.md",
          "line": 15,
          "message": "ADDED requirement must have at least one scenario"
        }
      ]
    },
    {
      "name": "validation",
      "type": "spec",
      "valid": true,
      "issues": []
    }
  ],
  "summary": {
    "total": 2,
    "passed": 1,
    "failed": 1,
    "errors": 1,
    "warnings": 0
  }
}
```

**Implementation Strategy**:

**Option A: Using jq (Preferred)**:
```bash
print_json_results() {
    # Build items array
    local items_json="[]"

    for item_name in "${VALIDATED_ITEMS[@]}"; do
        local item_type="${ITEM_TYPES[$item_name]}"
        local item_issues="[]"
        local is_valid="true"

        # Find issues for this item
        for issue in "${ISSUES[@]}"; do
            IFS='|' read -r level path line message <<< "$issue"

            # Check if issue belongs to this item
            if issue_belongs_to_item "$item_name" "$item_type" "$path"; then
                is_valid="false"
                item_issues=$(echo "$item_issues" | jq \
                    --arg level "$level" \
                    --arg path "$path" \
                    --arg line "$line" \
                    --arg msg "$message" \
                    '. + [{level: $level, path: $path, line: ($line | tonumber), message: $msg}]')
            fi
        done

        # Add item to items array
        items_json=$(echo "$items_json" | jq \
            --arg name "$item_name" \
            --arg type "$item_type" \
            --argjson valid "$is_valid" \
            --argjson issues "$item_issues" \
            '. + [{name: $name, type: $type, valid: $valid, issues: $issues}]')
    done

    # Build final JSON
    jq -n \
        --argjson items "$items_json" \
        --argjson total "${#VALIDATED_ITEMS[@]}" \
        --argjson passed "$passed_count" \
        --argjson failed "$failed_count" \
        --argjson errors "${ISSUE_COUNTS[errors]}" \
        '{
            version: 1,
            items: $items,
            summary: {
                total: $total,
                passed: $passed,
                failed: $failed,
                errors: $errors,
                warnings: 0
            }
        }'
}
```

**Option B: Printf-Based JSON (Fallback if jq unavailable)**:
```bash
print_json_printf() {
    echo "{"
    echo "  \"version\": 1,"
    echo "  \"items\": ["

    local first_item=true
    for item_name in "${VALIDATED_ITEMS[@]}"; do
        if ! $first_item; then
            echo ","
        fi
        first_item=false

        echo "    {"
        echo "      \"name\": \"$(json_escape "$item_name")\","
        echo "      \"type\": \"${ITEM_TYPES[$item_name]}\","
        echo "      \"valid\": $(item_is_valid "$item_name"),"
        echo "      \"issues\": ["

        print_item_issues_json "$item_name"

        echo "      ]"
        echo -n "    }"
    done

    echo
    echo "  ],"
    echo "  \"summary\": {"
    echo "    \"total\": ${#VALIDATED_ITEMS[@]},"
    echo "    \"passed\": $passed_count,"
    echo "    \"failed\": $failed_count,"
    echo "    \"errors\": ${ISSUE_COUNTS[errors]},"
    echo "    \"warnings\": 0"
    echo "  }"
    echo "}"
}

json_escape() {
    local str="$1"
    # Escape special JSON characters
    str="${str//\\/\\\\}"  # Backslash
    str="${str//\"/\\\"}"  # Quote
    str="${str//$'\n'/\\n}"  # Newline
    str="${str//$'\r'/\\r}"  # Carriage return
    str="${str//$'\t'/\\t}"  # Tab
    echo "$str"
}
```

**jq Availability Check**:
```bash
check_json_capability() {
    if command -v jq >/dev/null 2>&1; then
        return 0  # jq available
    else
        return 1  # jq not available
    fi
}

main() {
    # ...argument parsing...

    if $json_output; then
        if ! check_json_capability; then
            echo "Warning: jq not found, falling back to human-readable output" >&2
            json_output=false
        fi
    fi

    # ...validation...

    if $json_output; then
        print_json_results  # Uses jq
    else
        print_human_results
    fi
}
```

**Why This Approach**:

1. **jq is Standard**: Widely available in development environments
2. **Graceful Degradation**: Falls back to human output if jq missing
3. **Correctness**: jq ensures valid JSON (escaping, formatting)
4. **Maintainability**: jq syntax more readable than printf escaping

**Trade-offs**:
- ✅ Valid, well-formatted JSON
- ✅ Matches binary's JSON schema exactly
- ✅ Graceful degradation without jq
- ❌ jq dependency for JSON (but optional)
- ❌ More complex implementation

### Decision 6: Discovery Functions - Match internal/discovery/ Exactly

**What**: Implement discovery functions that produce identical results to
`internal/discovery/GetSpecIDs()` and `internal/discovery/GetChangeIDs()`.

**Spec Discovery**:

```bash
discover_specs() {
    local spectr_dir="${SPECTR_DIR:-spectr}"
    local specs_dir="$spectr_dir/specs"

    # Return empty if directory doesn't exist
    if [[ ! -d "$specs_dir" ]]; then
        return 0
    fi

    # Find all directories containing spec.md
    find "$specs_dir" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        if [[ -f "$dir/spec.md" ]]; then
            basename "$dir"
        fi
    done | sort
}
```

**Change Discovery**:

```bash
discover_changes() {
    local spectr_dir="${SPECTR_DIR:-spectr}"
    local changes_dir="$spectr_dir/changes"

    # Return empty if directory doesn't exist
    if [[ ! -d "$changes_dir" ]]; then
        return 0
    fi

    # Find all directories except 'archive'
    find "$changes_dir" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        local name
        name=$(basename "$dir")
        if [[ "$name" != "archive" ]]; then
            echo "$name"
        fi
    done | sort
}
```

**Why This Approach**:

1. **Exact Match**: Same discovery logic as Go implementation
2. **Error Handling**: Gracefully handles missing directories
3. **Sorting**: Alphabetical ordering for consistent output
4. **Archive Exclusion**: Skips archived changes (matches Go behavior)

**Go Reference** (`internal/discovery/discovery.go`):

```go
func GetSpecIDs(spectrRoot string) ([]string, error) {
    specsDir := filepath.Join(spectrRoot, "specs")
    entries, err := os.ReadDir(specsDir)
    if err != nil {
        if os.IsNotExist(err) {
            return []string{}, nil  // Return empty, not error
        }
        return nil, err
    }

    var specs []string
    for _, entry := range entries {
        if entry.IsDir() {
            specFile := filepath.Join(specsDir, entry.Name(), "spec.md")
            if _, err := os.Stat(specFile); err == nil {
                specs = append(specs, entry.Name())
            }
        }
    }

    sort.Strings(specs)
    return specs, nil
}
```

**Trade-offs**:
- ✅ Exact behavioral match with Go
- ✅ Handles edge cases (missing dirs)
- ✅ Simple, readable implementation
- ❌ Requires find command (standard but worth noting)

### Decision 7: Regex Patterns - Derive from internal/markdown/

**What**: Define regex patterns that match the behavior of
`internal/markdown/matchers.go` functions.

**Pattern Definitions**:

```bash
# Section headers (case-insensitive on level 2 headers)
readonly REQUIREMENTS_SECTION_PATTERN='^##[[:space:]]+Requirements[[:space:]]*$'

# Requirement header (level 3)
# Matches: ### Requirement: User Authentication
# Captures: "User Authentication"
readonly REQUIREMENT_HEADER_PATTERN='^###[[:space:]]+Requirement:[[:space:]]+(.+)$'

# Scenario header (level 4)
# Matches: #### Scenario: Successful login
readonly SCENARIO_HEADER_PATTERN='^####[[:space:]]+Scenario:'

# Delta section headers (case-insensitive on operation type)
readonly DELTA_ADDED_PATTERN='^##[[:space:]]+(ADDED|Added)[[:space:]]+Requirements[[:space:]]*$'
readonly DELTA_MODIFIED_PATTERN='^##[[:space:]]+(MODIFIED|Modified)[[:space:]]+Requirements[[:space:]]*$'
readonly DELTA_REMOVED_PATTERN='^##[[:space:]]+(REMOVED|Removed)[[:space:]]+Requirements[[:space:]]*$'
readonly DELTA_RENAMED_PATTERN='^##[[:space:]]+(RENAMED|Renamed)[[:space:]]+Requirements[[:space:]]*$'

# Malformed scenario patterns (for error detection)
readonly MALFORMED_SCENARIO_3_HASH='^###[[:space:]]+Scenario:'
readonly MALFORMED_SCENARIO_5_HASH='^#####[[:space:]]+Scenario:'
readonly MALFORMED_SCENARIO_6_HASH='^######[[:space:]]+Scenario:'
readonly MALFORMED_SCENARIO_BOLD='^\*\*Scenario:'
readonly MALFORMED_SCENARIO_BULLET_BOLD='^-[[:space:]]+\*\*Scenario:'

# Task item (flexible format)
# Matches: - [ ] Task
#          - [x] Completed task
#          - [ ] 1.1 Numbered task
#          - [X] 2.3 Completed numbered task (uppercase X)
readonly TASK_ITEM_PATTERN='^[[:space:]]*-[[:space:]]+\[([ xX])\][[:space:]]+'

# RENAMED requirement patterns
# Matches: - FROM: ### Requirement: OldName
# Matches: - TO: ### Requirement: NewName
readonly RENAMED_FROM_PATTERN='^-[[:space:]]+FROM:[[:space:]]+###[[:space:]]+Requirement:[[:space:]]+(.+)$'
readonly RENAMED_TO_PATTERN='^-[[:space:]]+TO:[[:space:]]+###[[:space:]]+Requirement:[[:space:]]+(.+)$'

# Alternative RENAMED format (backtick-wrapped)
# Matches: - FROM: `### Requirement: OldName`
readonly RENAMED_FROM_BACKTICK_PATTERN='^-[[:space:]]+FROM:[[:space:]]+`###[[:space:]]+Requirement:[[:space:]]+([^`]+)`'
readonly RENAMED_TO_BACKTICK_PATTERN='^-[[:space:]]+TO:[[:space:]]+`###[[:space:]]+Requirement:[[:space:]]+([^`]+)`'
```

**Go Reference Mapping**:

| Bash Pattern | Go Function | File |
|--------------|-------------|------|
| `REQUIREMENTS_SECTION_PATTERN` | Section detection | parser.go:ExtractSections |
| `REQUIREMENT_HEADER_PATTERN` | `MatchRequirementHeader()` | markdown/matchers.go |
| `SCENARIO_HEADER_PATTERN` | `MatchScenarioHeader()` | markdown/matchers.go |
| `DELTA_*_PATTERN` | `MatchDeltaSection()` | markdown/matchers.go |
| `TASK_ITEM_PATTERN` | `MatchFlexibleTask()` | markdown/matchers.go |
| `RENAMED_*_PATTERN` | `MatchRenamedFrom/To()` | markdown/matchers.go |

**Pattern Testing Strategy**:

```bash
# Test patterns against known valid/invalid inputs
test_patterns() {
    # Valid Requirements header
    [[ "## Requirements" =~ $REQUIREMENTS_SECTION_PATTERN ]] || echo "FAIL: Requirements pattern"

    # Valid requirement header
    [[ "### Requirement: User Auth" =~ $REQUIREMENT_HEADER_PATTERN ]] || echo "FAIL: Requirement pattern"
    [[ "${BASH_REMATCH[1]}" == "User Auth" ]] || echo "FAIL: Requirement capture"

    # Valid scenario header
    [[ "#### Scenario: Login succeeds" =~ $SCENARIO_HEADER_PATTERN ]] || echo "FAIL: Scenario pattern"

    # Invalid (malformed) scenario headers
    [[ "### Scenario: Wrong" =~ $MALFORMED_SCENARIO_3_HASH ]] || echo "FAIL: 3-hash detection"
    [[ "**Scenario: Bold" =~ $MALFORMED_SCENARIO_BOLD ]] || echo "FAIL: Bold detection"

    # Valid task items
    [[ "- [ ] Task" =~ $TASK_ITEM_PATTERN ]] || echo "FAIL: Empty task"
    [[ "- [x] Done" =~ $TASK_ITEM_PATTERN ]] || echo "FAIL: Completed task"
    [[ "  - [ ] Indented" =~ $TASK_ITEM_PATTERN ]] || echo "FAIL: Indented task"
}
```

**Why These Patterns**:

1. **Behavioral Equivalence**: Match Go's markdown parsing behavior
2. **POSIX Compatibility**: Use `[[:space:]]` instead of `\s` for portability
3. **Capture Groups**: Use `(.+)` for extracting requirement names
4. **Case Insensitivity**: Handle both "ADDED" and "Added" (user may vary)

**Trade-offs**:
- ✅ Matches Go behavior exactly
- ✅ Portable across GNU/BSD grep
- ✅ Clear pattern names
- ❌ Must manually sync with Go changes
- ❌ No pattern compilation optimization (acceptable for bash)

### Decision 8: Argument Parsing - Explicit Mode Flags

**What**: Use explicit `--spec`, `--change`, `--all` flags rather than
positional arguments or type inference.

**Argument Structure**:

```bash
Usage: validate.sh [OPTIONS]

Options:
  --spec <spec-id>       Validate a single spec by ID
  --change <change-id>   Validate a single change by ID
  --all                  Validate all specs and changes
  --json                 Output JSON instead of human-readable format
  -h, --help             Show this help message

Examples:
  validate.sh --spec validation
  validate.sh --change add-feature
  validate.sh --all
  validate.sh --all --json
```

**Implementation**:

```bash
main() {
    local mode=""
    local target=""
    local json_output=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --spec)
                if [[ -z "$2" ]]; then
                    echo "Error: --spec requires an argument" >&2
                    print_usage
                    exit 2
                fi
                mode="spec"
                target="$2"
                shift 2
                ;;
            --change)
                if [[ -z "$2" ]]; then
                    echo "Error: --change requires an argument" >&2
                    print_usage
                    exit 2
                fi
                mode="change"
                target="$2"
                shift 2
                ;;
            --all)
                mode="all"
                shift
                ;;
            --json)
                json_output=true
                shift
                ;;
            -h|--help)
                print_usage
                exit 0
                ;;
            *)
                echo "Error: Unknown option: $1" >&2
                print_usage
                exit 2
                ;;
        esac
    done

    # Validate required arguments
    if [[ -z "$mode" ]]; then
        echo "Error: Must specify --spec, --change, or --all" >&2
        print_usage
        exit 2
    fi

    # Execute based on mode
    case "$mode" in
        spec)
            validate_single_spec "$target"
            ;;
        change)
            validate_single_change "$target"
            ;;
        all)
            validate_all
            ;;
    esac

    # Output results
    if $json_output; then
        if check_json_capability; then
            print_json_results
        else
            echo "Warning: jq not found, using human-readable output" >&2
            print_human_results
        fi
    else
        print_human_results
    fi

    # Exit with appropriate code
    if [[ ${ISSUE_COUNTS[errors]} -gt 0 ]]; then
        exit 1
    fi
    exit 0
}
```

**Why Explicit Flags**:

1. **Clarity**: User intent is explicit (`--spec` vs `--change`)
2. **No Ambiguity**: Avoid type inference when name exists as both spec and
   change
3. **Consistency**: Matches common CLI patterns (git, kubectl use explicit flags)
4. **Error Prevention**: Can't accidentally validate wrong type

**Alternatives Considered**:

**Alternative 1: Positional Argument with Type Inference**
```bash
validate.sh validation  # Infer if it's a spec or change
```
**Rejected Because**:
- Ambiguous when same name exists as both spec and change
- Requires directory checking to infer type (slower, more complex)
- Less clear user intent

**Alternative 2: Positional with `--type` Flag**
```bash
validate.sh --type spec validation
validate.sh --type change add-feature
```
**Rejected Because**:
- More verbose than explicit mode flags
- Still requires `--type` disambiguation
- Less intuitive than `--spec` and `--change`

**Trade-offs**:
- ✅ Crystal clear intent
- ✅ No ambiguity resolution needed
- ✅ Simple implementation
- ❌ Slightly more typing than positional args

### Decision 9: Environment Variable Support - SPECTR_DIR Only

**What**: Support customizing spectr directory location via `SPECTR_DIR`
environment variable, defaulting to `"spectr"`.

**Implementation**:

```bash
# Configuration (near top of script)
readonly SPECTR_DIR="${SPECTR_DIR:-spectr}"
readonly SPECS_DIR="$SPECTR_DIR/specs"
readonly CHANGES_DIR="$SPECTR_DIR/changes"

# Usage in discovery
discover_specs() {
    if [[ ! -d "$SPECS_DIR" ]]; then
        return 0
    fi
    find "$SPECS_DIR" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        if [[ -f "$dir/spec.md" ]]; then
            basename "$dir"
        fi
    done | sort
}

# Usage example
$ SPECTR_DIR=my-specs ./validate.sh --all
```

**Why SPECTR_DIR Only**:

1. **Sufficient Flexibility**: Covers the primary customization need
2. **Matches Binary**: Binary uses `spectr.yaml` but this is simpler for script
3. **Standard Pattern**: ENV variables common for directory configuration
4. **No Flag Clutter**: Avoids `--spectr-dir` flag

**Non-Goals**:

- Not implementing full `spectr.yaml` parsing (too complex for bash)
- Not supporting per-directory customization (specs vs changes)
- Not supporting multiple spectr roots

**Trade-offs**:
- ✅ Simple, single-purpose configuration
- ✅ Standard environment variable pattern
- ✅ Easy to use in CI
- ❌ Less flexible than full config file
- ❌ Can't customize specs/changes dirs separately

### Decision 10: Exit Codes - Standard CLI Conventions

**What**: Use standard exit code conventions for success/failure/usage.

**Exit Code Specification**:

| Code | Meaning | When Used |
|------|---------|-----------|
| 0 | Success | All validations passed (no errors) |
| 1 | Validation Failure | One or more items have validation errors |
| 2 | Usage Error | Invalid arguments, missing required flags |

**Implementation**:

```bash
# At end of main()
if [[ ${ISSUE_COUNTS[errors]} -gt 0 ]]; then
    exit 1
fi
exit 0

# For usage errors (throughout argument parsing)
if [[ -z "$mode" ]]; then
    echo "Error: Must specify --spec, --change, or --all" >&2
    print_usage
    exit 2
fi
```

**CI Integration Example**:

```bash
# In CI pipeline
if ! ./validate.sh --all; then
    echo "Validation failed, blocking merge"
    exit 1
fi

# Check exit code explicitly
./validate.sh --all
case $? in
    0)
        echo "✓ All validations passed"
        ;;
    1)
        echo "✗ Validation failures detected"
        exit 1
        ;;
    2)
        echo "✗ Usage error in validation script"
        exit 1
        ;;
esac
```

**Why These Codes**:

1. **Standard Convention**: 0=success, 1=failure, 2=usage error is common
2. **CI Friendly**: Exit 1 fails CI pipelines automatically
3. **Scripting**: Easy to check `if validate.sh; then...`
4. **Matches Binary**: Binary uses same exit code pattern

**Trade-offs**:
- ✅ Standard, expected behavior
- ✅ Works with all CI systems
- ✅ Easy to script against
- ❌ Only 3 codes (but sufficient)

## Implementation Details

### File Structure

```
internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/
├── SKILL.md                    # AgentSkills metadata and documentation (~150 lines)
└── scripts/
    └── validate.sh             # Main validation script (~600 lines)
```

### validate.sh Internal Structure

**Script Organization** (~600 lines total):

```bash
#!/usr/bin/env bash
#
# validate.sh - Validate Spectr specifications and change proposals
#
# This script provides validation capabilities equivalent to `spectr validate`
# for environments where the spectr binary is not available (sandboxed AI
# environments, CI pipelines, fresh checkouts).
#
# Usage:
#   validate.sh --spec <spec-id>      # Validate single spec
#   validate.sh --change <change-id>  # Validate single change
#   validate.sh --all                 # Validate all specs and changes
#   validate.sh --json                # Output JSON (combine with above)
#
# Exit codes:
#   0 - All validations passed
#   1 - One or more validations failed
#   2 - Usage error (invalid arguments)
#
# Environment variables:
#   SPECTR_DIR - Custom spectr directory location (default: "spectr")
#

set -euo pipefail

# === Configuration === (lines 1-30)
readonly SCRIPT_VERSION="1.0.0"
readonly SPECTR_DIR="${SPECTR_DIR:-spectr}"
readonly SPECS_DIR="$SPECTR_DIR/specs"
readonly CHANGES_DIR="$SPECTR_DIR/changes"

# === Color Setup === (lines 31-50)
setup_colors() { ... }

# === Pattern Definitions === (lines 51-100)
readonly REQUIREMENTS_SECTION_PATTERN=...
readonly REQUIREMENT_HEADER_PATTERN=...
# ...all patterns...

# === Issue Collection === (lines 101-120)
declare -a ISSUES=()
declare -A ISSUE_COUNTS=([errors]=0 [warnings]=0)
declare -a VALIDATED_ITEMS=()
declare -A ITEM_TYPES=()

add_issue() { ... }

# === Spec File Validation === (lines 121-250)
validate_spec_file() { ... }
validate_requirement_state() { ... }
check_malformed_scenario() { ... }

# === Change Delta Validation === (lines 251-450)
validate_change() { ... }
validate_delta_spec_file() { ... }
flush_delta_requirement() { ... }
flush_delta_section() { ... }

# === Tasks File Validation === (lines 451-480)
validate_tasks_file() { ... }

# === Discovery Functions === (lines 481-510)
discover_specs() { ... }
discover_changes() { ... }

# === Validation Orchestration === (lines 511-560)
validate_single_spec() { ... }
validate_single_change() { ... }
validate_all() { ... }

# === Output Formatting === (lines 561-650)
print_human_results() { ... }
print_json_results() { ... }
check_json_capability() { ... }

# === Utility Functions === (lines 651-680)
print_usage() { ... }
make_relative_path() { ... }
issue_belongs_to_item() { ... }

# === Main Entry Point === (lines 681-750)
main() { ... }

# Script execution
main "$@"
```

### Key Function Implementations

**validate_spec_file() - Complete Implementation**:

```bash
validate_spec_file() {
    local file="$1"
    local line_num=0
    local in_requirements=false
    local current_requirement=""
    local requirement_line=0
    local has_scenario=false
    local has_shall_must=false
    local found_requirements_section=false

    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))

        # Detect ## Requirements section
        if [[ "$line" =~ $REQUIREMENTS_SECTION_PATTERN ]]; then
            # Flush previous requirement before entering section
            if [[ -n "$current_requirement" ]]; then
                validate_requirement_state "$file" "$current_requirement" \
                    "$requirement_line" "$has_scenario" "$has_shall_must"
                current_requirement=""
            fi
            found_requirements_section=true
            in_requirements=true
            continue
        fi

        # Detect other ## sections (exit Requirements)
        if [[ "$line" =~ ^##[[:space:]] ]]; then
            # Flush current requirement before leaving section
            if [[ -n "$current_requirement" ]]; then
                validate_requirement_state "$file" "$current_requirement" \
                    "$requirement_line" "$has_scenario" "$has_shall_must"
                current_requirement=""
            fi
            in_requirements=false
            continue
        fi

        # Only process lines within Requirements section
        if ! $in_requirements; then
            continue
        fi

        # Detect ### Requirement: header
        if [[ "$line" =~ $REQUIREMENT_HEADER_PATTERN ]]; then
            # Flush previous requirement
            if [[ -n "$current_requirement" ]]; then
                validate_requirement_state "$file" "$current_requirement" \
                    "$requirement_line" "$has_scenario" "$has_shall_must"
            fi
            # Start new requirement
            current_requirement="${BASH_REMATCH[1]}"
            requirement_line=$line_num
            has_scenario=false
            has_shall_must=false
            continue
        fi

        # Detect #### Scenario: header
        if [[ "$line" =~ $SCENARIO_HEADER_PATTERN ]]; then
            has_scenario=true
            continue
        fi

        # Check for SHALL/MUST in requirement content
        if [[ -n "$current_requirement" ]]; then
            if [[ "$line" =~ (SHALL|MUST) ]]; then
                has_shall_must=true
            fi

            # Check for malformed scenarios
            check_malformed_scenario "$file" "$line_num" "$line"
        fi
    done < "$file"

    # Flush final requirement
    if [[ -n "$current_requirement" ]]; then
        validate_requirement_state "$file" "$current_requirement" \
            "$requirement_line" "$has_scenario" "$has_shall_must"
    fi

    # Check if Requirements section was found
    if ! $found_requirements_section; then
        add_issue "ERROR" "$file" 1 \
            "Missing required '## Requirements' section"
    fi
}
```

**validate_delta_spec_file() - Complete Implementation**:

```bash
validate_delta_spec_file() {
    local spec_path="$1"
    local delta_count=0
    local line_num=0
    local current_section=""
    local section_line=0
    local section_has_requirements=false
    local current_requirement=""
    local requirement_line=0
    local has_scenario=false
    local has_shall_must=false

    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))

        # Detect delta section headers
        if [[ "$line" =~ $DELTA_ADDED_PATTERN ]] || \
           [[ "$line" =~ $DELTA_MODIFIED_PATTERN ]] || \
           [[ "$line" =~ $DELTA_REMOVED_PATTERN ]] || \
           [[ "$line" =~ $DELTA_RENAMED_PATTERN ]]; then

            # Flush previous requirement
            flush_delta_requirement "$spec_path" "$current_section" \
                "$current_requirement" "$requirement_line" \
                "$has_scenario" "$has_shall_must"

            # Flush previous section
            flush_delta_section "$spec_path" "$current_section" \
                "$section_line" "$section_has_requirements"

            # Extract section type (ADDED, MODIFIED, etc.)
            if [[ "$line" =~ ([A-Za-z]+)[[:space:]]+Requirements ]]; then
                current_section=$(echo "${BASH_REMATCH[1]}" | tr '[:lower:]' '[:upper:]')
            fi

            section_line=$line_num
            section_has_requirements=false
            current_requirement=""
            ((delta_count++))
            continue
        fi

        # Detect other ## sections (exit delta context)
        if [[ "$line" =~ ^##[[:space:]] ]]; then
            flush_delta_requirement "$spec_path" "$current_section" \
                "$current_requirement" "$requirement_line" \
                "$has_scenario" "$has_shall_must"
            flush_delta_section "$spec_path" "$current_section" \
                "$section_line" "$section_has_requirements"
            current_section=""
            current_requirement=""
            continue
        fi

        # Only process if in delta section
        if [[ -z "$current_section" ]]; then
            continue
        fi

        # Detect ### Requirement: header
        if [[ "$line" =~ $REQUIREMENT_HEADER_PATTERN ]]; then
            # Flush previous requirement
            flush_delta_requirement "$spec_path" "$current_section" \
                "$current_requirement" "$requirement_line" \
                "$has_scenario" "$has_shall_must"

            # Start new requirement
            section_has_requirements=true
            current_requirement="${BASH_REMATCH[1]}"
            requirement_line=$line_num
            has_scenario=false
            has_shall_must=false
            continue
        fi

        # Detect #### Scenario: header
        if [[ "$line" =~ $SCENARIO_HEADER_PATTERN ]]; then
            has_scenario=true
            continue
        fi

        # Check for SHALL/MUST and malformed scenarios
        if [[ -n "$current_requirement" ]]; then
            if [[ "$line" =~ (SHALL|MUST) ]]; then
                has_shall_must=true
            fi
            check_malformed_scenario "$spec_path" "$line_num" "$line"
        fi
    done < "$spec_path"

    # Final flushes
    flush_delta_requirement "$spec_path" "$current_section" \
        "$current_requirement" "$requirement_line" \
        "$has_scenario" "$has_shall_must"
    flush_delta_section "$spec_path" "$current_section" \
        "$section_line" "$section_has_requirements"

    echo "$delta_count"
}

flush_delta_requirement() {
    local spec_path="$1"
    local section="$2"
    local requirement="$3"
    local req_line="$4"
    local has_scenario="$5"
    local has_shall_must="$6"

    if [[ -z "$requirement" ]]; then
        return
    fi

    # REMOVED requirements don't need scenarios or SHALL/MUST
    if [[ "$section" == "REMOVED" ]]; then
        return
    fi

    # RENAMED section uses different format, skip normal validation
    if [[ "$section" == "RENAMED" ]]; then
        return
    fi

    # ADDED and MODIFIED require scenarios and SHALL/MUST
    local req_path="$spec_path: $section Requirement '$requirement'"

    if [[ "$has_shall_must" != "true" ]]; then
        add_issue "ERROR" "$req_path" "$req_line" \
            "$section requirement must contain SHALL or MUST"
    fi

    if [[ "$has_scenario" != "true" ]]; then
        add_issue "ERROR" "$req_path" "$req_line" \
            "$section requirement must have at least one scenario"
    fi
}

flush_delta_section() {
    local spec_path="$1"
    local section="$2"
    local section_line="$3"
    local has_requirements="$4"

    if [[ -z "$section" ]]; then
        return
    fi

    if [[ "$has_requirements" != "true" ]]; then
        add_issue "ERROR" "$spec_path" "$section_line" \
            "$section Requirements section is empty (no requirements found)"
    fi
}
```

### SKILL.md Content Structure

```markdown
---
name: spectr-validate-wo-spectr-bin
description: Validate Spectr specifications and change proposals without requiring the spectr binary
compatibility:
  requirements:
    - bash 4.0+
    - grep (GNU or BSD)
    - sed (GNU or BSD)
    - find (GNU or BSD)
    - Optional: jq 1.5+ (for JSON output)
  platforms:
    - Linux
    - macOS
    - Unix-like systems with bash
---

# Spectr Validate (Without Binary)

This skill provides specification validation capabilities equivalent to `spectr
validate` for environments where the spectr binary cannot be installed:
sandboxed AI coding environments (Claude Code, Cursor), CI pipelines, or fresh
repository checkouts.

## Quick Start

```bash
# Validate a single spec
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --spec validation

# Validate a single change
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --change add-feature

# Validate everything
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all

# JSON output for programmatic use
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all --json
```

## Validation Rules

The script implements the same validation rules as `spectr validate`:

### Spec Files

- **Missing `## Requirements` section**: ERROR
- **Requirement without SHALL or MUST**: ERROR (strict mode)
- **Requirement without `#### Scenario:` block**: ERROR (strict mode)
- **Malformed scenario format** (wrong header level, bullets, bold): ERROR

### Change Delta Specs

- **No delta sections**: ERROR (must have ADDED, MODIFIED, REMOVED, or RENAMED)
- **Empty delta section**: ERROR (section exists but no requirements)
- **ADDED requirement without scenario**: ERROR
- **ADDED requirement without SHALL/MUST**: ERROR
- **MODIFIED requirement without scenario**: ERROR
- **MODIFIED requirement without SHALL/MUST**: ERROR

### Tasks Files

- **tasks.md with no task items**: ERROR (if file exists, must have `- [ ]` or
  `- [x]`)

## Usage Examples

### Validate Before Committing

```bash
# Check your changes before committing
if bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --change my-feature; then
    git add .
    git commit -m "Add feature"
else
    echo "Validation failed, please fix errors"
fi
```

### CI Integration

```bash
# In .github/workflows/validate.yml
- name: Validate Specifications
  run: |
    bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all
```

### JSON Output Processing

```bash
# Parse JSON output with jq
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all --json | \
    jq '.items[] | select(.valid == false) | .name'
```

## Environment Variables

- `SPECTR_DIR`: Custom spectr directory location (default: `spectr`)

```bash
SPECTR_DIR=my-specs bash ./validate.sh --all
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All validations passed |
| 1 | One or more validations failed |
| 2 | Usage error (invalid arguments) |

## Limitations

This skill provides core validation but omits some advanced features:

- **No pre-merge validation**: Does not check if MODIFIED/REMOVED requirements
  exist in base specs
- **No cross-capability duplicate detection**: Validates per-file duplicates
  only
- **Sequential processing**: No parallel validation (slower than binary)

For production workflows requiring full validation, use `spectr validate`.

## Dependencies

- **Required**: bash 4.0+, grep, sed, find (standard on all Unix systems)
- **Optional**: jq (for JSON output; falls back to human-readable if missing)

## Troubleshooting

### jq not found

If you see "Warning: jq not found" when using `--json`:
```bash
# Install jq
brew install jq  # macOS
apt-get install jq  # Debian/Ubuntu
yum install jq  # RHEL/CentOS
```

### Custom spectr directory

If your project uses a non-standard location:
```bash
export SPECTR_DIR=custom-specs
bash ./validate.sh --all
```

### Comparing with binary output

To verify the skill produces identical results:
```bash
# Run both and compare
spectr validate --all > binary-output.txt
bash ./validate.sh --all > skill-output.txt
diff binary-output.txt skill-output.txt
```
```

## Risks / Trade-offs

### Risk 1: Regex Pattern Drift

**Risk**: Bash regex patterns may drift from Go implementation as markdown
parsing evolves.

**Impact**: Medium - Could cause false positives/negatives

**Mitigation**:
- Document pattern sources with Go file references in code comments
- Add test suite comparing skill output vs binary output
- Include validation comparison in CI
- Review patterns on each Go validation update

**Acceptance Criteria**:
- Pattern drift detected within 1 release cycle
- Test failures prevent skill regression

### Risk 2: Performance on Large Codebases

**Risk**: Sequential bash processing significantly slower than Go's concurrent
validation on large projects (100+ specs).

**Impact**: Low - Acceptable for target use cases

**Mitigation**:
- Document performance characteristics in SKILL.md
- Typical usage (30-50 specs) completes in <5 seconds
- AI assistants validate incrementally (single specs/changes)
- Full `--all` validation reserved for CI (acceptable delay)

**Acceptance Criteria**:
- <10s for 50 spec files on typical hardware
- No complaints from users about speed

### Risk 3: Shell Portability Issues

**Risk**: Bash-specific features may fail on strict POSIX shells or older bash
versions.

**Impact**: Medium - Could break in some environments

**Mitigation**:
- Explicitly require bash 4.0+ in shebang and SKILL.md
- Use `#!/usr/bin/env bash` for portability
- Test on macOS (BSD) and Linux (GNU) coreutils
- Avoid GNU-specific grep/sed extensions where possible
- Use `[[ =~ ]]` bash extension (not POSIX but widely supported)

**Acceptance Criteria**:
- Works on macOS and Linux without modification
- Clear error if bash < 4.0

### Risk 4: jq Dependency for JSON

**Risk**: jq may not be available in all environments, limiting JSON output
capability.

**Impact**: Low - Graceful degradation available

**Mitigation**:
- Make jq optional with clear warning
- Fallback to human-readable output if jq missing
- Document jq installation in SKILL.md
- Consider printf-based JSON generation as alternative (more complex but
  dependency-free)

**Acceptance Criteria**:
- JSON output works when jq available
- Clear warning and fallback when jq missing
- No hard failure

### Risk 5: Maintenance Burden

**Risk**: Skill requires updates whenever Go validation logic changes.

**Impact**: Medium - Ongoing maintenance cost

**Mitigation**:
- Clear documentation linking patterns to Go source
- Automated comparison testing (skill vs binary)
- CI checks for validation parity
- Include skill review in Go validation PR checklist

**Acceptance Criteria**:
- Skill updates included in validation PRs
- Comparison tests prevent regressions

## Migration Plan

No migration needed - this is a new feature addition. Existing workflows
continue unchanged.

**Rollout**:
1. Skill automatically installed for Claude Code users on next `spectr init`
2. Users can manually invoke skill scripts as needed
3. No changes to existing `spectr validate` command or behavior

## Open Questions

### Q1: Should we implement printf-based JSON as primary strategy?

**Question**: Should we avoid jq dependency by implementing JSON generation with
printf/echo?

**Options**:
1. **Use jq (current decision)**: Simpler, correct, but optional dependency
2. **Use printf**: Dependency-free but complex escaping logic
3. **Hybrid**: jq preferred, printf fallback

**Recommendation**: Stick with jq (Option 1). Most development environments have
jq, and correctness (escaping, formatting) is important for JSON. Fallback to
human output if jq missing is acceptable.

### Q2: Should we support custom validation rules?

**Question**: Should the skill support user-defined validation rules via config
file?

**Options**:
1. **Hard-coded rules** (current): Match Go exactly, no customization
2. **Config-based**: Load additional rules from `.validate-rules` file

**Recommendation**: Hard-coded (Option 1). Keeps skill simple and aligned with
binary. Custom rules belong in binary implementation, not shell script.

### Q3: Should we validate individual delta files directly?

**Question**: Should `--file` flag support validating
`changes/foo/specs/bar/spec.md` directly?

**Options**:
1. **Change-level only** (current): `--change foo` validates all its deltas
2. **File-level**: `--file changes/foo/specs/bar/spec.md` validates single delta

**Recommendation**: Change-level only (Option 1). Simpler API, and users can
validate entire change easily. Add file-level if requested.

### Q4: Should we implement RENAMED requirements validation?

**Question**: RENAMED sections use special format (`- FROM:`, `- TO:`). Should
we validate these?

**Options**:
1. **Skip RENAMED validation** (current): Just count as delta, don't validate
   format
2. **Full RENAMED validation**: Parse FROM/TO pairs, check format, detect
   malformed

**Recommendation**: Skip (Option 1) for initial version. RENAMED is complex and
rarely used. Can add if requested.

### Q5: Should we support `--type` disambiguation flag?

**Question**: Should we support `--type change|spec` for explicit type
specification?

**Options**:
1. **Explicit mode flags only** (current): `--spec` or `--change` required
2. **Add `--type` flag**: `validate.sh --type change add-feature`

**Recommendation**: Explicit flags (Option 1). Clearer, simpler, no redundancy.
--type would just duplicate --spec/--change functionality.
