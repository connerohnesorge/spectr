# Implementation Tasks

## 1. Requirement Parsing Infrastructure

- [x] 1.1 Create `internal/parsers/requirement_parser.go`
  - Parse `### Requirement: <name>` headers
  - Extract requirement name and full block content
  - Parse `#### Scenario:` blocks within requirements
  - Handle multi-line requirement text
- [x] 1.2 Create `internal/parsers/delta_parser.go`
  - Parse `## ADDED Requirements` sections
  - Parse `## MODIFIED Requirements` sections
  - Parse `## REMOVED Requirements` sections (requirement names only)
  - Parse `## RENAMED Requirements` sections (FROM/TO pairs)
  - Normalize requirement names (trim, case-insensitive matching)
- [x] 1.3 Add tests for parsers
  - Test requirement block extraction
  - Test delta section parsing
  - Test edge cases (empty sections, malformed headers)
  - Test scenario parsing

## 2. Spec Merging Logic

- [x] 2.1 Create `internal/archive/spec_merger.go`
  - Implement requirement name normalization
  - Build requirement name-to-block map from base spec
  - Apply RENAMED operations (update map keys and headers)
  - Apply REMOVED operations (delete from map)
  - Apply MODIFIED operations (replace blocks)
  - Apply ADDED operations (insert new blocks)
  - Preserve original requirement ordering when possible
  - Append new requirements at end
- [x] 2.2 Implement spec reconstruction
  - Split spec into preamble, ## Requirements header, body, and after
  - Rebuild requirements section from processed blocks
  - Maintain blank line spacing (collapse 3+ newlines to 2)
  - Ensure proper markdown structure
- [x] 2.3 Create spec skeleton generator for new specs
  - Generate `# <Capability> Specification` header
  - Add purpose section with placeholder text
  - Add `## Requirements` section
- [x] 2.4 Add tests for spec merger
  - Test each delta operation type
  - Test operation ordering
  - Test duplicate detection
  - Test cross-section conflict detection
  - Test requirement ordering preservation

## 3. Validation Layer

- [x] 3.1 Create `internal/archive/validator.go`
  - Validate proposal.md structure (informational only)
  - Validate delta spec files have at least one operation
  - Validate each requirement has at least one scenario
  - Validate scenario format (`#### Scenario: <name>`)
  - Detect duplicate requirements within sections
  - Detect cross-section conflicts (same requirement in multiple sections)
  - Validate RENAMED interplay with MODIFIED (must reference new name)
  - Validate ADDED doesn't collide with RENAMED TO
- [x] 3.2 Implement pre-merge validation
  - Check source requirements exist for MODIFIED/REMOVED/RENAMED
  - Check target doesn't exist for ADDED/RENAMED
  - Validate delta operations are valid for new vs existing specs
- [x] 3.3 Implement post-merge validation
  - Validate rebuilt spec structure
  - Ensure no duplicate requirement names in result
  - Validate all scenarios are properly formatted
- [x] 3.4 Add tests for validation
  - Test duplicate detection
  - Test conflict detection
  - Test scenario validation
  - Test pre/post merge validation

## 4. Archive Command Implementation

- [x] 4.1 Create `cmd/archive.go`
  - Define ArchiveCmd struct with flags (Yes, SkipSpecs, NoValidate)
  - Implement Run() method
  - Add command to CLI struct in `cmd/root.go`
- [x] 4.2 Create `internal/archive/archiver.go`
  - Implement change selection (interactive if no ID provided)
  - Check if change directory exists
  - Get current working directory and locate `spectr/changes`
- [x] 4.3 Implement validation workflow
  - Run proposal validation (informational warnings only)
  - Run delta spec validation (blocking errors)
  - Skip validation if --no-validate flag (with confirmation)
  - Display validation errors with chalk-style formatting
- [x] 4.4 Implement task checking
  - Use existing `parsers.CountTasks()` for task status
  - Display task progress
  - Warn if incomplete tasks exist
  - Require confirmation if incomplete (unless --yes flag)
- [x] 4.5 Implement spec update workflow
  - Find delta specs in `changes/<id>/specs/*/spec.md`
  - Find corresponding main specs in `spectr/specs/*/spec.md`
  - Display specs to update (create vs update status)
  - Prompt for confirmation (unless --yes flag)
  - Skip spec updates if --skip-specs flag
- [x] 4.6 Implement spec application
  - For each delta spec, build updated spec via merger
  - Validate rebuilt spec before writing
  - Write updated spec to main specs directory
  - Display operation counts per spec (+ added, ~ modified, - removed, →
    renamed)
  - Display totals across all specs
- [x] 4.7 Implement archive operation
  - Create archive directory: `spectr/changes/archive/`
  - Generate archive name: `YYYY-MM-DD-<change-id>`
  - Check if archive already exists (error if it does)
  - Move change directory to archive
  - Display success message with archive name
- [x] 4.8 Add tests for archive command
  - Test interactive change selection
  - Test validation workflow
  - Test task checking
  - Test spec update workflow
  - Test archive operation
  - Test --yes, --skip-specs, --no-validate flags

## 5. Integration and Documentation

- [x] 5.1 Update go.mod if new dependencies needed
  - Add any required libraries (e.g., for colored output, prompts)
- [x] 5.2 Run tests and ensure all pass
  - Unit tests for all new packages
  - Integration test for full archive workflow
- [x] 5.3 Test command manually
  - Create test change with delta specs
  - Run `spectr archive` interactively
  - Run `spectr archive <id> --yes` non-interactively
  - Verify specs are correctly updated
  - Verify archive is created with correct date
- [x] 5.4 Update tasks.md checkboxes to completed (- [x])
  - Mark all items as completed after verification

## 6. Final Validation

- [x] 6.1 Run `spectr validate add-archive-command --strict`
  - Ensure no validation errors
  - Fix any issues found
- [x] 6.2 Verify proposal completeness
  - All tasks completed
  - All delta specs valid
  - Design decisions documented
