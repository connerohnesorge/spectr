# Design: spectr-validate-wo-spectr-bin Skill

## Context

The spectr CLI provides comprehensive validation for specification files and
change proposals. However, in sandboxed environments (Claude Code's execution
context), CI pipelines without spectr installed, or fresh repository checkouts,
users cannot run `spectr validate`. This creates a gap in the development
workflow where validation feedback is unavailable.

The existing `spectr-accept-wo-spectr-bin` skill demonstrates the pattern:
implement core functionality using bash scripts with standard Unix tools (grep,
sed, awk, jq) to provide equivalent behavior without the binary.

### Stakeholders

- **AI Coding Assistants**: Primary users who need to validate specs before/after
  modifications
- **CI/CD Systems**: Automated validation in pipelines without spectr installed
- **Developers**: Manual validation in restricted environments

### Constraints

1. **No external dependencies** beyond standard Unix tools (grep, sed, awk) and
   optionally jq for JSON output
2. **Behavioral parity** with `spectr validate` for core validation rules
3. **Performance**: Must handle validation of 30+ spec files efficiently
4. **Portability**: Must work on Linux, macOS, and Unix-like systems with bash

## Goals / Non-Goals

### Goals

1. Validate spec files for structural correctness (Requirements section,
   scenarios)
2. Validate change delta specs (ADDED/MODIFIED/REMOVED/RENAMED sections)
3. Validate tasks.md files (at least one task item)
4. Provide human-readable error output with file paths and line numbers
5. Provide JSON output for programmatic consumption
6. Exit with appropriate codes for scripting integration

### Non-Goals

1. **Pre-merge validation against base specs** - This requires parsing and
   comparing two files which is complex in bash; we focus on structural
   validation only
2. **Interactive TUI mode** - Out of scope for shell script
3. **Concurrent validation** - Sequential processing is acceptable for skill use
4. **Cross-file duplicate detection across capabilities** - Complex state
   tracking; validate per-file only

## Decisions

### Decision 1: Script Architecture - Single Entry Point with Functions

**What**: Use a single `validate.sh` script with modular functions rather than
multiple scripts.

**Why**: Simpler installation, easier maintenance, single source of truth for
validation logic. The accept skill uses a single script and this pattern works
well.

**Alternatives considered**:
- Multiple scripts (`validate-spec.sh`, `validate-change.sh`, etc.) - Rejected
  due to duplication and complexity
- Python script - Rejected to avoid Python dependency

### Decision 2: Parsing Strategy - Line-by-Line with State Machine

**What**: Process files line-by-line maintaining state for current section,
requirement, etc.

**Why**: Bash's read loop is efficient and allows precise line number tracking
for error reporting. This matches how the Go implementation tracks positions.

```bash
parse_spec_file() {
    local file="$1"
    local line_num=0
    local current_section=""
    local current_requirement=""
    local requirement_line=0
    local has_scenario=false
    local has_shall_must=false

    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))

        # Detect ## sections
        if [[ "$line" =~ ^##[[:space:]](.+)$ ]]; then
            # Flush previous requirement if any
            flush_requirement_validation
            current_section="${BASH_REMATCH[1]}"
            continue
        fi

        # Detect ### Requirement: headers
        if [[ "$line" =~ ^###[[:space:]]Requirement:[[:space:]](.+)$ ]]; then
            flush_requirement_validation
            current_requirement="${BASH_REMATCH[1]}"
            requirement_line=$line_num
            has_scenario=false
            has_shall_must=false
            continue
        fi

        # Detect #### Scenario: headers
        if [[ "$line" =~ ^####[[:space:]]Scenario: ]]; then
            has_scenario=true
            continue
        fi

        # Check for SHALL/MUST in requirement content
        if [[ -n "$current_requirement" ]] && [[ "$line" =~ (SHALL|MUST) ]]; then
            has_shall_must=true
        fi

        # Detect malformed scenarios
        if [[ "$line" =~ ^###[[:space:]]Scenario: ]] || \
           [[ "$line" =~ ^#####[[:space:]]Scenario: ]] || \
           [[ "$line" =~ ^\*\*Scenario: ]] || \
           [[ "$line" =~ ^-[[:space:]]\*\*Scenario: ]]; then
            report_error "$file" "$line_num" "Malformed scenario format"
        fi
    done < "$file"

    # Flush final requirement
    flush_requirement_validation
}
```

**Alternatives considered**:
- Grep-based extraction - Harder to get accurate line numbers
- Awk script - More complex, less readable for maintenance

### Decision 3: Validation Rules - Match Go Implementation Exactly

**What**: Implement the same validation rules as `internal/validation/`:

| Rule | Level | Go Source |
|------|-------|-----------|
| Missing `## Requirements` section | ERROR | spec_rules.go:35-40 |
| Requirement without SHALL/MUST | ERROR (strict) | spec_rules.go:109-117 |
| Requirement without scenario | ERROR (strict) | spec_rules.go:120-127 |
| Malformed scenario format | ERROR | spec_rules.go:130-143 |
| Empty delta section | ERROR | delta_validators.go:65-73 |
| Delta requirement without scenario | ERROR | delta_validators.go:105-114 |
| Delta requirement without SHALL/MUST | ERROR | delta_validators.go:92-102 |
| tasks.md with no task items | ERROR | change_rules.go:797-807 |

**Why**: Users should get identical validation results whether using the binary
or the skill script.

### Decision 4: Error Output Format - Match Human-Readable Output

**What**: Match the format from `internal/validation/formatters.go`:

```
changes/add-feature/specs/auth/spec.md
  [ERROR] line 15: ADDED requirement must have at least one scenario
  [ERROR] line 23: Requirement should contain SHALL or MUST

specs/validation/spec.md
  [ERROR] line 1: Missing required '## Requirements' section

Summary: 2 passed, 2 failed (3 errors), 4 total
```

**Why**: Consistent experience between binary and skill script.

### Decision 5: JSON Output Structure - Match `spectr validate --json`

**What**: Produce JSON matching the binary's output structure:

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
    }
  ],
  "summary": {
    "total": 4,
    "passed": 2,
    "failed": 2,
    "errors": 3,
    "warnings": 0
  }
}
```

**Why**: Programmatic consumers expect consistent JSON structure.

**Note**: Requires `jq` for JSON generation. Script will check for jq and fall
back to human-readable output if unavailable.

### Decision 6: Discovery Functions - Match Go Discovery Logic

**What**: Implement discovery matching `internal/discovery/`:

```bash
discover_specs() {
    # Find all directories under spectr/specs/ containing spec.md
    find spectr/specs -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        if [[ -f "$dir/spec.md" ]]; then
            basename "$dir"
        fi
    done | sort
}

discover_changes() {
    # Find all directories under spectr/changes/ except 'archive'
    find spectr/changes -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        local name
        name=$(basename "$dir")
        if [[ "$name" != "archive" ]]; then
            echo "$name"
        fi
    done | sort
}
```

**Why**: Consistent item discovery between binary and skill.

### Decision 7: Delta Section Detection - Regex Patterns

**What**: Use regex patterns matching `internal/markdown/` matchers:

```bash
# Delta section headers (case-insensitive match on operation)
DELTA_ADDED_PATTERN='^##[[:space:]]+(ADDED|Added)[[:space:]]+Requirements'
DELTA_MODIFIED_PATTERN='^##[[:space:]]+(MODIFIED|Modified)[[:space:]]+Requirements'
DELTA_REMOVED_PATTERN='^##[[:space:]]+(REMOVED|Removed)[[:space:]]+Requirements'
DELTA_RENAMED_PATTERN='^##[[:space:]]+(RENAMED|Renamed)[[:space:]]+Requirements'

# Requirement header
REQUIREMENT_PATTERN='^###[[:space:]]+Requirement:[[:space:]]+(.+)$'

# Scenario header (correct format)
SCENARIO_PATTERN='^####[[:space:]]+Scenario:[[:space:]]+'

# Malformed scenario patterns
MALFORMED_SCENARIO_PATTERNS=(
    '^###[[:space:]]+Scenario:'      # 3 hashtags instead of 4
    '^#####[[:space:]]+Scenario:'    # 5 hashtags
    '^######[[:space:]]+Scenario:'   # 6 hashtags
    '^\*\*Scenario:'                  # Bold instead of header
    '^-[[:space:]]+\*\*Scenario:'    # Bullet + bold
)

# Task item pattern
TASK_PATTERN='^[[:space:]]*-[[:space:]]+\[([ xX])\]'
```

**Why**: These patterns match the Go regex patterns and ensure consistent
parsing.

### Decision 8: Color Output - Conditional on TTY

**What**: Use ANSI color codes only when stdout is a TTY:

```bash
setup_colors() {
    if [[ -t 1 ]]; then
        RED='\033[0;31m'
        YELLOW='\033[0;33m'
        GREEN='\033[0;32m'
        NC='\033[0m' # No Color
    else
        RED=''
        YELLOW=''
        GREEN=''
        NC=''
    fi
}
```

**Why**: Clean output in CI logs and pipes, colored output in interactive use.

### Decision 9: Strict Mode Only - No Warning/Error Distinction

**What**: All validation issues are treated as errors (matching Go's strict-only
mode).

**Why**: The Go implementation removed the strict mode flag and now always
treats warnings as errors. The skill should match this behavior.

### Decision 10: Tasks.md Validation - Flexible Task Pattern

**What**: Use the flexible task pattern from `internal/markdown/`:

```bash
# Matches:
# - [ ] Task description
# - [x] Task description
# - [ ] 1.1 Task description
# - [x] 2.3 Completed task
TASK_PATTERN='^[[:space:]]*-[[:space:]]+\[([ xX])\][[:space:]]+'
```

**Why**: Match the `markdown.MatchFlexibleTask` function behavior.

## Implementation Details

### File Structure

```
internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/
├── SKILL.md                    # AgentSkills metadata and documentation
└── scripts/
    └── validate.sh             # Main validation script (~400-500 lines)
```

### validate.sh Internal Structure

```bash
#!/usr/bin/env bash
#
# validate.sh - Validate Spectr specs and changes without the binary
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
#   2 - Usage error
#

set -euo pipefail

# === Configuration ===
SPECTR_DIR="${SPECTR_DIR:-spectr}"
SPECS_DIR="$SPECTR_DIR/specs"
CHANGES_DIR="$SPECTR_DIR/changes"

# === Color Setup ===
setup_colors() { ... }

# === Pattern Definitions ===
readonly REQUIREMENTS_SECTION_PATTERN='^##[[:space:]]+Requirements[[:space:]]*$'
readonly REQUIREMENT_HEADER_PATTERN='^###[[:space:]]+Requirement:[[:space:]]+(.+)$'
readonly SCENARIO_HEADER_PATTERN='^####[[:space:]]+Scenario:'
# ... more patterns

# === Issue Collection ===
declare -a ISSUES=()
declare -A ISSUE_COUNTS=([errors]=0 [warnings]=0)

add_issue() {
    local level="$1" path="$2" line="$3" message="$4"
    ISSUES+=("$level|$path|$line|$message")
    if [[ "$level" == "ERROR" ]]; then
        ((ISSUE_COUNTS[errors]++))
    fi
}

# === Validation Functions ===

validate_spec_file() {
    local spec_path="$1"
    local line_num=0
    local in_requirements=false
    local current_requirement=""
    local requirement_line=0
    local has_scenario=false
    local has_shall_must=false
    local found_requirements_section=false

    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))

        # Check for ## Requirements section
        if [[ "$line" =~ $REQUIREMENTS_SECTION_PATTERN ]]; then
            found_requirements_section=true
            in_requirements=true
            continue
        fi

        # Check for other ## sections (exit requirements context)
        if [[ "$line" =~ ^##[[:space:]] ]] && [[ ! "$line" =~ $REQUIREMENTS_SECTION_PATTERN ]]; then
            # Flush current requirement before leaving section
            if [[ -n "$current_requirement" ]]; then
                validate_requirement_state "$spec_path" "$current_requirement" \
                    "$requirement_line" "$has_scenario" "$has_shall_must"
            fi
            in_requirements=false
            current_requirement=""
            continue
        fi

        # Only process if in Requirements section
        if ! $in_requirements; then
            continue
        fi

        # Check for ### Requirement: header
        if [[ "$line" =~ $REQUIREMENT_HEADER_PATTERN ]]; then
            # Flush previous requirement
            if [[ -n "$current_requirement" ]]; then
                validate_requirement_state "$spec_path" "$current_requirement" \
                    "$requirement_line" "$has_scenario" "$has_shall_must"
            fi
            current_requirement="${BASH_REMATCH[1]}"
            requirement_line=$line_num
            has_scenario=false
            has_shall_must=false
            continue
        fi

        # Check for #### Scenario: header
        if [[ "$line" =~ $SCENARIO_HEADER_PATTERN ]]; then
            has_scenario=true
            continue
        fi

        # Check for SHALL/MUST in content
        if [[ -n "$current_requirement" ]]; then
            if [[ "$line" =~ (SHALL|MUST) ]]; then
                has_shall_must=true
            fi

            # Check for malformed scenarios
            check_malformed_scenario "$spec_path" "$line_num" "$line"
        fi
    done < "$spec_path"

    # Final flush
    if [[ -n "$current_requirement" ]]; then
        validate_requirement_state "$spec_path" "$current_requirement" \
            "$requirement_line" "$has_scenario" "$has_shall_must"
    fi

    # Check if Requirements section was found
    if ! $found_requirements_section; then
        add_issue "ERROR" "$spec_path" 1 \
            "Missing required '## Requirements' section"
    fi
}

validate_requirement_state() {
    local path="$1" name="$2" line="$3" has_scenario="$4" has_shall_must="$5"

    if [[ "$has_shall_must" != "true" ]]; then
        add_issue "ERROR" "$path: Requirement '$name'" "$line" \
            "Requirement should contain SHALL or MUST to indicate normative requirement"
    fi

    if [[ "$has_scenario" != "true" ]]; then
        add_issue "ERROR" "$path: Requirement '$name'" "$line" \
            "Requirement should have at least one scenario"
    fi
}

check_malformed_scenario() {
    local path="$1" line_num="$2" line="$3"

    # Check each malformed pattern
    if [[ "$line" =~ ^###[[:space:]]+Scenario: ]] || \
       [[ "$line" =~ ^#####[[:space:]]+Scenario: ]] || \
       [[ "$line" =~ ^######[[:space:]]+Scenario: ]] || \
       [[ "$line" =~ ^\*\*Scenario: ]] || \
       [[ "$line" =~ ^-[[:space:]]+\*\*Scenario: ]]; then
        add_issue "ERROR" "$path" "$line_num" \
            "Scenarios must use '#### Scenario:' format (4 hashtags followed by 'Scenario:')"
    fi
}

validate_change() {
    local change_id="$1"
    local change_dir="$CHANGES_DIR/$change_id"
    local specs_dir="$change_dir/specs"

    # Check specs directory exists
    if [[ ! -d "$specs_dir" ]]; then
        add_issue "ERROR" "$specs_dir" 1 \
            "specs directory not found"
        return 1
    fi

    # Find all spec.md files
    local spec_files=()
    while IFS= read -r -d '' file; do
        spec_files+=("$file")
    done < <(find "$specs_dir" -name "spec.md" -print0 2>/dev/null)

    if [[ ${#spec_files[@]} -eq 0 ]]; then
        add_issue "ERROR" "$specs_dir" 1 \
            "no spec.md files found in specs directory"
        return 1
    fi

    local total_deltas=0

    # Validate each delta spec file
    for spec_file in "${spec_files[@]}"; do
        local deltas
        deltas=$(validate_delta_spec_file "$spec_file")
        ((total_deltas += deltas))
    done

    # Check for at least one delta
    if [[ $total_deltas -eq 0 ]]; then
        add_issue "ERROR" "$specs_dir" 1 \
            "Change must have at least one delta (ADDED, MODIFIED, REMOVED, or RENAMED requirement)"
    fi

    # Validate tasks.md if present
    validate_tasks_file "$change_dir"
}

validate_delta_spec_file() {
    local spec_path="$1"
    local delta_count=0
    local line_num=0
    local current_section=""
    local current_requirement=""
    local requirement_line=0
    local has_scenario=false
    local has_shall_must=false
    local section_has_requirements=false
    local section_line=0

    while IFS= read -r line || [[ -n "$line" ]]; do
        ((line_num++))

        # Check for delta section headers
        if [[ "$line" =~ ^##[[:space:]]+(ADDED|MODIFIED|REMOVED|RENAMED)[[:space:]]+Requirements ]]; then
            # Flush previous requirement
            flush_delta_requirement
            # Flush previous section
            flush_delta_section

            current_section="${BASH_REMATCH[1]}"
            section_line=$line_num
            section_has_requirements=false
            ((delta_count++))
            continue
        fi

        # Check for other ## sections (exit delta context)
        if [[ "$line" =~ ^##[[:space:]] ]]; then
            flush_delta_requirement
            flush_delta_section
            current_section=""
            continue
        fi

        # Only process if in a delta section
        if [[ -z "$current_section" ]]; then
            continue
        fi

        # Check for ### Requirement: header
        if [[ "$line" =~ $REQUIREMENT_HEADER_PATTERN ]]; then
            flush_delta_requirement
            section_has_requirements=true
            current_requirement="${BASH_REMATCH[1]}"
            requirement_line=$line_num
            has_scenario=false
            has_shall_must=false
            continue
        fi

        # Check for #### Scenario: header
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
    flush_delta_requirement
    flush_delta_section

    echo "$delta_count"
}

flush_delta_requirement() {
    if [[ -z "$current_requirement" ]]; then
        return
    fi

    local section_upper
    section_upper=$(echo "$current_section" | tr '[:lower:]' '[:upper:]')

    # REMOVED requirements don't need scenarios or SHALL/MUST
    if [[ "$section_upper" == "REMOVED" ]]; then
        current_requirement=""
        return
    fi

    # RENAMED section uses different format, skip normal validation
    if [[ "$section_upper" == "RENAMED" ]]; then
        current_requirement=""
        return
    fi

    # ADDED and MODIFIED require scenarios and SHALL/MUST
    if [[ "$has_shall_must" != "true" ]]; then
        add_issue "ERROR" "$spec_path: $section_upper Requirement '$current_requirement'" \
            "$requirement_line" "$section_upper requirement must contain SHALL or MUST"
    fi

    if [[ "$has_scenario" != "true" ]]; then
        add_issue "ERROR" "$spec_path: $section_upper Requirement '$current_requirement'" \
            "$requirement_line" "$section_upper requirement must have at least one scenario"
    fi

    current_requirement=""
}

flush_delta_section() {
    if [[ -z "$current_section" ]]; then
        return
    fi

    if [[ "$section_has_requirements" != "true" ]]; then
        add_issue "ERROR" "$spec_path" "$section_line" \
            "$current_section Requirements section is empty (no requirements found)"
    fi
}

validate_tasks_file() {
    local change_dir="$1"
    local tasks_file="$change_dir/tasks.md"

    # tasks.md is optional
    if [[ ! -f "$tasks_file" ]]; then
        return
    fi

    local task_count=0
    while IFS= read -r line; do
        if [[ "$line" =~ ^[[:space:]]*-[[:space:]]+\[([ xX])\] ]]; then
            ((task_count++))
        fi
    done < "$tasks_file"

    if [[ $task_count -eq 0 ]]; then
        add_issue "ERROR" "$tasks_file" 1 \
            "tasks.md exists but contains no task items; expected format: '- [ ] Task' or '- [ ] N.N Task'"
    fi
}

# === Discovery Functions ===

discover_specs() {
    if [[ ! -d "$SPECS_DIR" ]]; then
        return
    fi
    find "$SPECS_DIR" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        if [[ -f "$dir/spec.md" ]]; then
            basename "$dir"
        fi
    done | sort
}

discover_changes() {
    if [[ ! -d "$CHANGES_DIR" ]]; then
        return
    fi
    find "$CHANGES_DIR" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | while read -r dir; do
        local name
        name=$(basename "$dir")
        if [[ "$name" != "archive" ]]; then
            echo "$name"
        fi
    done | sort
}

# === Output Functions ===

print_human_results() {
    local items_json="$1"
    # Group issues by file and print with colors
    # ... (detailed implementation)
}

print_json_results() {
    local items_json="$1"
    # Generate JSON output using jq or printf
    # ... (detailed implementation)
}

# === Main Entry Point ===

main() {
    setup_colors

    local mode="" target="" json_output=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --spec)
                mode="spec"
                target="$2"
                shift 2
                ;;
            --change)
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
                echo "Unknown option: $1" >&2
                print_usage
                exit 2
                ;;
        esac
    done

    if [[ -z "$mode" ]]; then
        echo "Error: Must specify --spec, --change, or --all" >&2
        print_usage
        exit 2
    fi

    # Execute validation based on mode
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
        print_json_results
    else
        print_human_results
    fi

    # Exit with appropriate code
    if [[ ${ISSUE_COUNTS[errors]} -gt 0 ]]; then
        exit 1
    fi
    exit 0
}

main "$@"
```

### SKILL.md Content

```markdown
---
name: spectr-validate-wo-spectr-bin
description: Validate Spectr specifications and change proposals without requiring the spectr binary
compatibility:
  requirements:
    - bash (4.0+)
    - grep (GNU or BSD)
    - sed (GNU or BSD)
    - find (GNU or BSD)
    - Optional: jq (for JSON output)
  platforms:
    - Linux
    - macOS
    - Unix-like systems with bash
---

# Spectr Validate (Without Binary)

This skill provides the ability to validate Spectr specifications and change
proposals without requiring the `spectr` binary. This is useful in sandboxed
environments, CI pipelines, or fresh repository checkouts.

## Usage

### Validate a Single Spec

```bash
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --spec validation
```

### Validate a Single Change

```bash
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --change add-feature
```

### Validate Everything

```bash
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all
```

### JSON Output

Add `--json` to any command for machine-readable output:

```bash
bash .claude/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all --json
```

## Validation Rules

The script validates the following rules (matching `spectr validate` behavior):

### Spec Files
- Must have `## Requirements` section
- Requirements must contain SHALL or MUST keywords
- Requirements must have at least one `#### Scenario:` block
- Scenarios must use correct format (4 hashtags)

### Change Delta Specs
- Must have at least one delta section (ADDED, MODIFIED, REMOVED, RENAMED)
- ADDED/MODIFIED requirements must have scenarios
- ADDED/MODIFIED requirements must contain SHALL or MUST
- Delta sections must not be empty

### Tasks Files
- If tasks.md exists, must contain at least one task item (`- [ ]` or `- [x]`)

## Exit Codes

- `0` - All validations passed
- `1` - One or more validations failed
- `2` - Usage error (invalid arguments)

## Limitations

- Does not perform pre-merge validation against base specs
- Does not detect cross-capability duplicate requirements
- Sequential processing (no parallel validation)

For production workflows with full validation, use the `spectr validate` command.
```

## Risks / Trade-offs

### Risk 1: Regex Pattern Drift

**Risk**: Bash regex patterns may drift from Go implementation over time.

**Mitigation**: Document pattern sources in code comments referencing Go file
locations. Include test cases comparing outputs.

### Risk 2: Performance on Large Codebases

**Risk**: Sequential bash processing may be slow with many files.

**Mitigation**: Acceptable trade-off for skill use case. The existing validation
pool uses 6 workers; our sequential approach will be slower but still functional
for typical use (30-50 files).

### Risk 3: Shell Compatibility

**Risk**: Bash-isms may fail on strict POSIX shells.

**Mitigation**: Explicitly require bash 4.0+ in shebang and SKILL.md
compatibility section. Use `#!/usr/bin/env bash`.

### Risk 4: jq Dependency for JSON

**Risk**: jq may not be available in all environments.

**Mitigation**: Make jq optional. Fall back to human-readable output with a
warning if --json requested but jq unavailable. Alternatively, generate JSON
using printf (less elegant but dependency-free).

## Migration Plan

No migration needed - this is a new feature addition.

## Open Questions

1. **Should we support `--type` flag for disambiguation?** The binary supports
   `--type change|spec` when names conflict. For simplicity, we could require
   explicit `--spec` or `--change` flags instead.

   **Proposed Answer**: Use explicit flags (`--spec`, `--change`) only. Simpler
   implementation and clearer API.

2. **Should we support validation of individual delta files?** The binary allows
   `spectr validate changes/foo/specs/bar/spec.md` directly.

   **Proposed Answer**: Not in initial version. Add if requested. Users can
   validate the entire change instead.

3. **Should we support custom spectr directory paths?** The binary uses
   `spectr.yaml` for configuration.

   **Proposed Answer**: Support via `SPECTR_DIR` environment variable for
   flexibility. Default to `spectr/`.
