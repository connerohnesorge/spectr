# Change: Clarify design.md Purpose - Implementation Details

## Why

The current documentation describes `design.md` as containing "Technical decisions" or "Technical patterns," which is too abstract. In practice, effective `design.md` files contain **specific implementation details** such as code structures, API signatures, data models, and concrete examples. This disconnect leads to:

1. **Vague design documents** - Authors write high-level principles instead of actionable implementation details
2. **Poor AI agent guidance** - Agents receive insufficient specificity for implementation
3. **Inconsistent content** - Some design.md files are detailed (like `add-provider-specific-templates/design.md` with Go code examples), while others are abstract

## What Changes

Update documentation to clarify that `design.md` contains **specific implementation details**:

| File | Current Description | Updated Description |
|------|---------------------|---------------------|
| `spectr/AGENTS.md` | "Technical decisions" | "Implementation details (code structures, APIs, data models)" |
| `spectr/project.md` | (mentions design.md but not its purpose) | Add clarification section |
| `README.md` | "Technical patterns" / "Technical decisions" | "Implementation details with code examples" |
| `docs/src/content/docs/guides/creating-changes.md` | "Technical decisions" | "Implementation details with code examples" |

### Additional Clarifications to Add

1. **Examples of what to include**:
   - Data structures and type definitions
   - API signatures and interfaces
   - File and directory structures
   - Code snippets showing patterns
   - Configuration schemas

2. **Explicit distinction from other files**:
   - `proposal.md`: High-level why/what/impact
   - `tasks.md`: Implementation checklist
   - `spec.md`: Requirements and acceptance criteria
   - `design.md`: Specific implementation details for developers

3. **Before/after examples** showing what detailed vs. vague design documents look like

## Impact

- **Affected documentation files**: 4 files need updates
- **No code changes** - purely documentation
- **Backward compatible** - existing design.md files remain valid
- **Impact on existing users**: Agents will receive more specific guidance

## Files to Update

1. `spectr/AGENTS.md` - Multiple locations (directory structure, skeleton, descriptions)
2. `spectr/project.md` - Add clarification section about file purposes
3. `README.md` - Directory structure comments and FAQ
4. `docs/src/content/docs/guides/creating-changes.md` - Design document section
