# Tasks

## 1. Create Skill Directory Structure

- [ ] 1.1 Create `internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/` directory
- [ ] 1.2 Create `scripts/` subdirectory for shell scripts
- [ ] 1.3 Create `SKILL.md` with AgentSkills frontmatter (name, description, compatibility requirements)
- [ ] 1.4 Document usage examples in SKILL.md (--spec, --change, --all, --json flags)
- [ ] 1.5 Document validation rules in SKILL.md (matching spectr validate behavior)
- [ ] 1.6 Document exit codes in SKILL.md (0=success, 1=failure, 2=usage)
- [ ] 1.7 Document limitations in SKILL.md (no pre-merge validation, no cross-capability duplicates)

## 2. Implement Core Script Infrastructure

- [ ] 2.1 Create `scripts/validate.sh` with proper shebang (`#!/usr/bin/env bash`)
- [ ] 2.2 Add set flags (`set -euo pipefail`) for robust error handling
- [ ] 2.3 Implement configuration section (SPECTR_DIR environment variable support)
- [ ] 2.4 Implement color setup function with TTY detection (ANSI codes only if interactive)
- [ ] 2.5 Define regex patterns matching internal/markdown/ (Requirements, Requirement, Scenario headers)
- [ ] 2.6 Define malformed scenario patterns (3/5/6 hashtags, bold, bullets)
- [ ] 2.7 Define delta section patterns (ADDED, MODIFIED, REMOVED, RENAMED)
- [ ] 2.8 Define task item pattern (flexible format with optional IDs)
- [ ] 2.9 Implement issue collection arrays and counters (ISSUES array, ISSUE_COUNTS associative array)
- [ ] 2.10 Implement add_issue() function for standardized issue tracking

## 3. Implement Spec File Validation

- [ ] 3.1 Implement validate_spec_file() function with line-by-line reading
- [ ] 3.2 Implement state tracking (in_requirements flag, current_requirement, requirement_line)
- [ ] 3.3 Implement Requirements section detection (## Requirements pattern match)
- [ ] 3.4 Implement requirement header parsing (### Requirement: pattern)
- [ ] 3.5 Implement scenario header detection (#### Scenario: pattern)
- [ ] 3.6 Implement SHALL/MUST keyword detection in requirement content
- [ ] 3.7 Implement malformed scenario detection (check_malformed_scenario function)
- [ ] 3.8 Implement requirement state validation (validate_requirement_state function)
- [ ] 3.9 Implement requirement flushing on section boundaries
- [ ] 3.10 Implement missing Requirements section error reporting

## 4. Implement Change Delta Validation

- [ ] 4.1 Implement validate_change() function with specs directory detection
- [ ] 4.2 Implement delta spec file discovery (find all spec.md under specs/)
- [ ] 4.3 Implement validate_delta_spec_file() function with delta counting
- [ ] 4.4 Implement delta section detection (ADDED/MODIFIED/REMOVED/RENAMED Requirements)
- [ ] 4.5 Implement requirement parsing within delta sections
- [ ] 4.6 Implement delta requirement validation (scenarios, SHALL/MUST for ADDED/MODIFIED)
- [ ] 4.7 Implement empty delta section detection
- [ ] 4.8 Implement flush_delta_requirement() function with section-specific rules
- [ ] 4.9 Implement flush_delta_section() function for empty section checks
- [ ] 4.10 Implement total delta count check (error if zero deltas)
- [ ] 4.11 Skip scenario/SHALL validation for REMOVED requirements
- [ ] 4.12 Skip normal validation for RENAMED section (different format)

## 5. Implement Tasks File Validation

- [ ] 5.1 Implement validate_tasks_file() function
- [ ] 5.2 Check if tasks.md exists (skip if not present)
- [ ] 5.3 Count task items using flexible pattern (`- [ ]`, `- [x]`, with/without IDs)
- [ ] 5.4 Report error if tasks.md exists but has zero tasks
- [ ] 5.5 Provide helpful error message with expected format examples

## 6. Implement Discovery Functions

- [ ] 6.1 Implement discover_specs() function using find command
- [ ] 6.2 Filter for directories under spectr/specs/ containing spec.md
- [ ] 6.3 Sort spec IDs alphabetically
- [ ] 6.4 Implement discover_changes() function using find command
- [ ] 6.5 Filter for directories under spectr/changes/ excluding 'archive'
- [ ] 6.6 Sort change IDs alphabetically
- [ ] 6.7 Handle missing directories gracefully (return empty if directory doesn't exist)

## 7. Implement Validation Orchestration

- [ ] 7.1 Implement validate_single_spec() function
- [ ] 7.2 Verify spec directory exists before validation
- [ ] 7.3 Call validate_spec_file() for spec.md
- [ ] 7.4 Implement validate_single_change() function
- [ ] 7.5 Verify change directory exists before validation
- [ ] 7.6 Call validate_change() for change directory
- [ ] 7.7 Implement validate_all() function
- [ ] 7.8 Discover and validate all specs
- [ ] 7.9 Discover and validate all changes
- [ ] 7.10 Aggregate results across all items

## 8. Implement Output Formatting

- [ ] 8.1 Implement print_human_results() function
- [ ] 8.2 Group issues by file path for cleaner output
- [ ] 8.3 Print file headers followed by indented issues
- [ ] 8.4 Apply color codes to error levels ([ERROR] in red, [WARNING] in yellow)
- [ ] 8.5 Format issue lines with "line X: message" pattern
- [ ] 8.6 Implement summary line (X passed, Y failed (E errors), Z total)
- [ ] 8.7 Add blank line separators between failed items
- [ ] 8.8 Implement print_json_results() function
- [ ] 8.9 Generate JSON structure with version, items, summary fields
- [ ] 8.10 Use jq for JSON generation if available
- [ ] 8.11 Fallback to printf-based JSON if jq unavailable
- [ ] 8.12 Warn user if --json requested but jq not found

## 9. Implement Argument Parsing and Main Entry Point

- [ ] 9.1 Implement print_usage() function with flag documentation
- [ ] 9.2 Implement main() function with argument parsing loop
- [ ] 9.3 Parse --spec <id> flag and store mode/target
- [ ] 9.4 Parse --change <id> flag and store mode/target
- [ ] 9.5 Parse --all flag and set mode
- [ ] 9.6 Parse --json flag and set json_output boolean
- [ ] 9.7 Parse -h/--help flag and print usage
- [ ] 9.8 Validate required arguments (error if mode not specified)
- [ ] 9.9 Execute validation based on mode (spec/change/all)
- [ ] 9.10 Output results based on json_output flag
- [ ] 9.11 Exit with code 0 if no errors, 1 if errors, 2 if usage error

## 10. Register Skill with Claude Code Provider

- [ ] 10.1 Open `internal/initialize/providers/claude.go`
- [ ] 10.2 Add `NewAgentSkillsInitializer` call in Initializers() method
- [ ] 10.3 Pass skill name `spectr-validate-wo-spectr-bin`
- [ ] 10.4 Pass target directory `.claude/skills/spectr-validate-wo-spectr-bin`
- [ ] 10.5 Pass TemplateManager parameter

## 11. Testing and Verification

- [ ] 11.1 Test skill installation via `spectr init` (verify .claude/skills/ created)
- [ ] 11.2 Verify SKILL.md has valid frontmatter
- [ ] 11.3 Verify scripts/validate.sh is executable (0755 permissions)
- [ ] 11.4 Test --spec validation against existing specs in spectr/specs/validation/spec.md
- [ ] 11.5 Test --change validation against existing changes (pick one with deltas)
- [ ] 11.6 Test --all validation to ensure discovery works
- [ ] 11.7 Create test spec with missing Requirements section, verify ERROR reported
- [ ] 11.8 Create test spec with requirement missing scenario, verify ERROR reported
- [ ] 11.9 Create test spec with malformed scenario (3 hashtags), verify ERROR reported
- [ ] 11.10 Create test change with empty ADDED section, verify ERROR reported
- [ ] 11.11 Create test change with ADDED requirement missing scenario, verify ERROR reported
- [ ] 11.12 Create test tasks.md with no tasks, verify ERROR reported
- [ ] 11.13 Test --json output format (verify valid JSON with jq if available)
- [ ] 11.14 Test color output in TTY (verify ANSI codes present)
- [ ] 11.15 Test color output in pipe (verify no ANSI codes)
- [ ] 11.16 Compare skill output vs `spectr validate` output for consistency
- [ ] 11.17 Verify exit codes (0 for pass, 1 for fail, 2 for usage error)

## 12. Documentation and Polish

- [ ] 12.1 Add inline comments to validate.sh explaining key sections
- [ ] 12.2 Add function docstrings for major validation functions
- [ ] 12.3 Add examples to SKILL.md for each validation mode
- [ ] 12.4 Document SPECTR_DIR environment variable override
- [ ] 12.5 Document jq optional dependency and fallback behavior
- [ ] 12.6 Add troubleshooting section to SKILL.md (common errors)
- [ ] 12.7 Update SKILL.md with performance notes (sequential vs parallel)
