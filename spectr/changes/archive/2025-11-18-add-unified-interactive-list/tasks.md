# Implementation Tasks: Unified Interactive List Mode

## 1. Type-Safe Data Structures

- [x] 1.1 Create unified `Item` type that wraps either `ChangeInfo` or `SpecInfo`
- [x] 1.2 Add `ItemType` enum (change, spec) for runtime type tracking
- [x] 1.3 Define `ItemList` type for managing mixed change and spec collections

## 2. Lister Extension

- [x] 2.1 Add `ListAll()` method to Lister that returns unified `[]Item`
- [x] 2.2 Implement sorting/ordering logic for mixed item lists
- [x] 2.3 Add optional filtering parameters (by type, status, etc.)

## 3. Interactive Model Expansion

- [x] 3.1 Refactor `interactiveModel` to support multiple item types
- [x] 3.2 Add item type indicator in table (e.g., [SPEC] or [CHANGE] prefix)
- [x] 3.3 Implement type-specific keyboard handlers (e.g., edit only for specs)
- [x] 3.4 Add mode toggle key for filtering by item type
- [x] 3.5 Update help text to reflect new capabilities

## 4. Table Display Layer

- [x] 4.1 Create unified `RunInteractiveAll()` function
- [x] 4.2 Implement dual-column table: ID, Type, Title, Details (varies by type)
- [x] 4.3 Add visual differentiation for item types (colors, icons, or prefixes)
- [x] 4.4 Handle variable-width columns based on item content

## 5. Command-Line Integration

- [x] 5.1 Add `--all` flag to ListCmd for unified mode
- [x] 5.2 Update validation logic to exclude incompatible flag combinations
- [x] 5.3 Route to new interactive mode when `--all` and `--interactive` are combined
- [x] 5.4 Maintain backward compatibility with existing flags

## 6. Editor Integration

- [x] 6.1 Extend handleEdit() to support spec editing in unified mode
- [x] 6.2 Skip edit option gracefully for change items
- [x] 6.3 Verify file paths are correct for both types

## 7. Clipboard & Selection

- [x] 7.1 Ensure selected ID works for both changes and specs
- [x] 7.2 Test clipboard copy for unified selections
- [x] 7.3 Verify error handling for missing files in unified mode

## 8. Testing

- [x] 8.1 Write unit tests for unified Item type and conversion
- [x] 8.2 Write unit tests for ListAll() method
- [x] 8.3 Write unit tests for mixed item table rendering
- [x] 8.4 Write integration tests for interactive mode with mixed items
- [x] 8.5 Test type-specific behaviors (edit, selection) in unified mode
- [x] 8.6 Verify all existing tests still pass with changes

## 9. Documentation

- [x] 9.1 Update help text for list command
- [x] 9.2 Add examples of unified mode usage
- [x] 9.3 Document new keyboard shortcuts/behavior
- [x] 9.4 Update AGENTS.md if needed for workflow changes

## Notes

- Maintain single responsibility: each module handles its concern
- Use minimal changes to existing interactive.go structure
- Ensure table column widths scale appropriately for unified view
- Consider performance impact of loading both changes and specs simultaneously
