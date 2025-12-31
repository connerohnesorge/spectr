# Change: Add spectr-validate-wo-spectr-bin AgentSkills Initializer

## Why

Claude Code and other sandboxed AI coding environments cannot install or execute
the spectr binary, creating a validation gap in the development workflow. When
working with Spectr specifications and change proposals, AI assistants need to
validate their work before committing changes, but currently have no way to run
`spectr validate` in restricted execution contexts.

This problem affects several scenarios:

1. **AI Coding Assistants**: Claude Code, Cursor, and other tools running in
   sandboxed environments cannot validate specifications they create or modify
2. **CI/CD Pipelines**: Automated workflows on systems where spectr binary is
   not installed need validation capabilities
3. **Fresh Checkouts**: Developers working in environments where installing Go
   binaries is restricted or complicated

The existing `spectr-accept-wo-spectr-bin` skill demonstrates the solution
pattern: implement core functionality using bash scripts with standard Unix
tools, providing equivalent behavior without requiring the binary.

Without validation capabilities in sandboxed environments, AI assistants cannot:
- Verify spec file structure (Requirements sections, scenario formatting)
- Check change proposals for delta operation validity
- Validate tasks.md files for proper task formatting
- Provide real-time feedback during specification authoring

This creates a disconnect between the workflow and the tooling, forcing manual
validation steps or blind commits that may fail later in CI.

## What Changes

**Add new AgentSkills skill**: `spectr-validate-wo-spectr-bin`

- **Location**: `internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/`
- **Structure**:
  - `SKILL.md` - AgentSkills frontmatter and comprehensive usage documentation
  - `scripts/validate.sh` - Main validation script implementing core validation
    rules

**Register skill with Claude Code provider**:

- Add `NewAgentSkillsInitializer` call in
  `internal/initialize/providers/claude.go`
- Skill will be automatically installed to
  `.claude/skills/spectr-validate-wo-spectr-bin/` when `spectr init` runs

**Validation capabilities** (matching `spectr validate` behavior):

1. **Spec File Validation**:
   - Detect missing `## Requirements` section (ERROR)
   - Validate requirements contain SHALL or MUST keywords (ERROR in strict mode)
   - Verify requirements have at least one `#### Scenario:` block (ERROR in
     strict mode)
   - Detect malformed scenario formatting (wrong header levels, bullets,
     bold) (ERROR)

2. **Change Delta Validation**:
   - Verify at least one delta section exists (ADDED, MODIFIED, REMOVED,
     RENAMED) (ERROR)
   - Check delta sections are not empty (ERROR)
   - Validate ADDED requirements have scenarios and SHALL/MUST (ERROR)
   - Validate MODIFIED requirements have scenarios and SHALL/MUST (ERROR)
   - Check REMOVED requirements format

3. **Tasks File Validation**:
   - If `tasks.md` exists, verify it contains at least one task item (ERROR)
   - Support flexible task formats (`- [ ]`, `- [x]`, with or without task IDs)

4. **Discovery & Bulk Validation**:
   - Discover all specs in `spectr/specs/` (excluding archive)
   - Discover all changes in `spectr/changes/` (excluding archive)
   - Validate individual items or all items with `--all` flag

5. **Output Formats**:
   - Human-readable output with file paths, line numbers, color-coded severity
   - JSON output for programmatic consumption (via `--json` flag)
   - Exit codes for CI integration (0=success, 1=failures, 2=usage error)

**Implementation approach**:

- Line-by-line parsing with state machine tracking current section/requirement
- Regex patterns matching `internal/markdown/` matchers for consistency
- Strict mode only (all warnings treated as errors, matching current Go behavior)
- Conditional color output (TTY detection for clean CI logs)
- Optional jq dependency for JSON (fallback to human output if unavailable)

## Impact

**Affected specs**:
- `spectr/specs/support-claude-code/spec.md` - Add 9 new scenarios for the
  validation skill

**Affected code**:
- `internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/` - New
  skill directory with SKILL.md and scripts/validate.sh (~500 lines total)
- `internal/initialize/providers/claude.go` - Add single
  `NewAgentSkillsInitializer` call (1 line change)

**Breaking changes**: None

**User impact**:
- AI coding assistants gain validation capabilities in sandboxed environments
- Existing workflows continue unchanged
- New validation method available alongside binary
- Skill auto-installs for Claude Code users on next `spectr init`

**Performance considerations**:
- Sequential validation (no parallel workers like binary)
- Acceptable for typical usage (30-50 spec files)
- Bash parsing slower than Go but sufficient for skill use case

**Maintenance implications**:
- Regex patterns must stay synchronized with Go implementation
- Test coverage should compare skill output vs binary output
- Documentation must clarify limitations vs full binary validation
