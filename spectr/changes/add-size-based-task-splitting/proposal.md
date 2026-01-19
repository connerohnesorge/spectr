# Change: Add Size-Based Task Splitting for Large tasks.jsonc Files

## Why

Large change proposals with 50+ tasks produce `tasks.jsonc` files that exceed
typical AI agent Read operation limits (~100-150 lines). This forces agents to
either:

1. Request partial reads with offset/limit (losing context)
2. Hit token/line truncation limits (missing tasks)
3. Fail to understand the full scope of work

Claude Code and other AI agents work best when they can read entire task files in
one operation. By automatically splitting large `tasks.md` files into multiple
smaller `tasks.jsonc` files at accept time, we ensure optimal agent readability
while maintaining human-friendly single-file authoring.

This replaces the previous `add-hierarchical-tasks` proposal, which tied
splitting to delta spec directories. The new approach is simpler: split based
purely on file size, using section boundaries for clean breaks.

## What Changes

- **ADDED**: Automatic detection of large tasks.md files (>100 lines) during
  `spectr accept`
- **ADDED**: Smart splitting logic that preserves section boundaries and
  subsection groupings
- **ADDED**: Root `tasks.jsonc` with references to split files (`tasks-1.jsonc`,
  `tasks-2.jsonc`, etc.)
- **ADDED**: Child task file format with full headers documenting origin (parent
  change ID, parent task ID)
- **ADDED**: Task reference syntax using `"children": "$ref:tasks-N.jsonc"` for
  hierarchical navigation
- **ADDED**: Version 2 schema for tasks.jsonc with `children`, `parent`, and
  `includes` fields
- **MODIFIED**: `spectr accept` command to implement splitting logic with
  100-line threshold
- **MODIFIED**: Task ID schema to support hierarchical IDs (e.g., `5.1`, `5.1.1`)

## Impact

- **Affected specs**: `cli-interface`
- **Affected code**:
  - `cmd/accept.go` - add splitting detection and orchestration
  - `cmd/accept_writer.go` - add JSONC marshaller and splitting functions
  - `internal/parsers/types.go` - add Children, Parent, Includes fields
  - `internal/parsers/parsers.go` - add ReadTasksFile function for hierarchical
    reading
  - `cmd/jsonc.go` (NEW) - custom JSONC marshaller handling quotes, trailing
    commas, comments
- **Breaking changes**: None - Version 1 flat files remain fully supported
- **Replaces**: The `add-hierarchical-tasks` change proposal (will be archived
  without implementation)

## Critical Issue: JSONC Marshalling

The current implementation uses `json.MarshalIndent()` which produces strict JSON.
This causes issues when combined with JSONC headers:

1. **Quote Escaping**: Task descriptions with quotes get double-escaped, breaking
   JSONC syntax
2. **JSONC Incompatibility**: Standard JSON marshalling doesn't support trailing
   commas or inline comments
3. **File Validity**: Generated tasks.jsonc files may be invalid JSONC, causing
   parsing failures in agents

**Solution**: Implement a custom JSONC marshaller that:

- Properly handles quote escaping in task descriptions
- Supports trailing commas in JSON arrays/objects
- Allows inline comments for documentation
- Validates output before writing to disk

This ensures all generated tasks.jsonc files (both root and child) are valid,
agent-readable JSONC.

## Splitting Algorithm Specification

The `spectr accept` command SHALL implement the following splitting algorithm for
`tasks.md` files exceeding 100 lines:

### Detection & Thresholds

- **Trigger**: `tasks.md` file size exceeds 100 lines (default threshold, configurable)
- **Minimum split size**: 25 lines (configurable) - prevents generation of tiny files
  with <25 lines
- **Output location**: All files written to the change directory
  (e.g., `spectr/changes/[change-id]/`)

### Section & Subsection Detection

Section boundaries are detected using Markdown header hierarchy:

1. **Primary sections** (H2 headers): Pattern `## N. [Section Name]` where `N` is
   a digit (e.g., `## 1. Implementation`)
2. **Subsections** (H3 headers): Pattern `### N.M. [Subsection Name]` where `N.M`
   indicates nesting (e.g., `### 1.1 API Endpoints`)
3. **Task items** (unordered lists): Pattern `- [ ] N.M.K Task description` under
   subsections

Each primary section (H2) becomes a candidate for splitting; subsections (H3)
within a section are the atomic grouping units.

### Splitting Strategy

**Algorithm**:

1. **Parse `tasks.md`** into logical segments:
   - Group consecutive lines by section header (H2)
   - Within each section, identify subsection boundaries (H3)
   - Preserve ordering of all H2 and H3 headers and tasks

2. **Calculate segment sizes** (in lines, including headers):
   - Each (H2 section + its H3 subsections + tasks) is one candidate segment
   - Record cumulative line count per segment

3. **Decide split points**:
   - If a single primary section (H2) exceeds 100 lines:
     - Split that section's H3 subsections across multiple child files
     - Prefer breaking at H3 boundaries; if a single H3 subsection exceeds 100
       lines, break at the nearest complete task item (after `- [ ]` line)
   - If cumulative segments fit within 100 lines, keep together in one child file
   - Apply **minimum-size bundling rule**: If a segment or subsection group falls
     below 25 lines, attempt to bundle it with the next adjacent segment
     (maintaining H2 group coherence where possible)

4. **Avoid overflow**: When bundling would violate the 100-line limit, create a new
   child file at the last valid H2 or H3 boundary

### Output Structure

**Root `tasks.jsonc`**:

- Contains only top-level tasks (one per primary section H2)
- Each task has an `id` matching its section number (e.g., "1", "2")
- Each task includes a `"children"` field: `"$ref:tasks-N.jsonc"`
- Root also includes `"includes": ["tasks-*.jsonc"]` for glob matching

**Child files** (e.g., `tasks-1.jsonc`, `tasks-2.jsonc`):

- Written by `cmd/accept_writer.go` using custom JSONC marshaller (`cmd/jsonc.go`)
- Include JSONC header comments documenting origin:

  ```jsonc
  // Generated by: spectr accept [change-id]
  // Parent change: [change-id]
  // Parent task: [top-level section ID]
  ```

- Contains `"version": 2`, `"parent": "[section-id]"`, and `"tasks": [...]`
- Task IDs use hierarchical format (e.g., `1.1`, `1.2`) matching parsed structure
- All task descriptions properly escaped for JSONC validity

### Implementation Locations

The following files SHALL implement this algorithm:

- **`cmd/accept.go`**:
  - Detect `tasks.md` size; if >100 lines, invoke splitting logic
  - Orchestrate root file creation and child file coordination

- **`cmd/accept_writer.go`**:
  - Implement parsing of `tasks.md` into sections/subsections
  - Calculate split points per algorithm above
  - Call custom JSONC marshaller for each child file

- **`cmd/jsonc.go`** (NEW):
  - Implement custom JSONC marshaller handling quote escaping, trailing commas,
    JSONC comments
  - Validate output JSONC before writing to disk
- **`internal/parsers/types.go`**:
  - Add `Children string` field to task struct (holds `$ref:tasks-N.jsonc`)
  - Add `Parent string` field (holds parent section ID)
  - Add `Includes []string` field (holds `["tasks-*.jsonc"]`)
  
- **`internal/parsers/parsers.go`**:
  - Add `ReadTasksFile(path string) (RootTasks, []ChildTasks, error)` function
  - Recursively read root and all child `tasks-*.jsonc` files
  - Return unified task list or per-file structures as needed

### Configuration

Support `.spectr/config.yaml` or flag overrides for:

```textyaml
splitConfig:
  maxLinesPerFile: 100          # Trigger threshold (default 100)
  minLinesPerFile: 25           # Minimum before bundling (default 25)
```text

Command-line equivalent: `spectr accept [change-id] --max-lines-per-file 100
--min-lines-per-file 25`

### Regeneration & Status Preservation

When `spectr accept` is re-run on the same change:

- Re-parse `tasks.md` and regenerate all split files
- Preserve task status values (pending/in_progress/completed) from existing
  `tasks-*.jsonc` files
- Match tasks by ID (e.g., "1.1" from old `tasks-1.jsonc` to "1.1" in newly
  generated file)

## Examples

### Before (flat, 150-line tasks.md)

```text
spectr/changes/my-big-change/
├── proposal.md
├── tasks.md          (150 lines, 45 tasks)
└── tasks.jsonc       (generated: 150 lines, hard to read)
```text

### After (split across 2 files)

```text
spectr/changes/my-big-change/
├── proposal.md
├── tasks.md          (150 lines, source of truth)
├── tasks.jsonc       (root: 40 lines with refs)
├── tasks-1.jsonc     (section 1: 60 lines)
└── tasks-2.jsonc     (section 2: 50 lines)
```text

Root `tasks.jsonc`:

```jsonc
{
  "version": 2,
  "tasks": [
    {
      "id": "1",
      "section": "Implementation",
      "description": "Implement core features",
      "status": "pending",
      "children": "$ref:tasks-1.jsonc"
    },
    {
      "id": "2",
      "section": "Testing",
      "description": "Add test coverage",
      "status": "pending",
      "children": "$ref:tasks-2.jsonc"
    }
  ],
  "includes": ["tasks-*.jsonc"]
}
```text

Child `tasks-1.jsonc`:

```jsonc
// Generated by: spectr accept add-size-based-task-splitting
// Parent change: add-size-based-task-splitting
// Parent task: 1
// Status values: "pending" | "in_progress" | "completed"
{
  "version": 2,
  "parent": "1",
  "tasks": [
    {
      "id": "1.1",
      "section": "Implementation",
      "description": "Create database schema",
      "status": "pending"
    },
    {
      "id": "1.2",
      "section": "Implementation",
      "description": "Implement API handlers",
      "status": "pending"
    }
  ]
}
```text
