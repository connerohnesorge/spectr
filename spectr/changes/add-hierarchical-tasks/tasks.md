## 1. Types and Schema
- [ ] 1.1 Update `internal/parsers/types.go` to add `Children` field to Task struct with JSON tag `children,omitempty`
- [ ] 1.2 Add `Summary` struct with Total, Completed, InProgress, Pending fields
- [ ] 1.3 Add `Includes` field to TasksFile struct for glob patterns
- [ ] 1.4 Add `Parent` field to TasksFile struct for child files
- [ ] 1.5 Update TasksFile Version constant to 2
- [ ] 1.6 Add unit tests for new struct fields serialization/deserialization

## 2. Accept Command - Auto-Split Logic
- [ ] 2.1 Add `matchSectionToCapability()` function to normalize section names to kebab-case
- [ ] 2.2 Add `findMatchingDeltaSpec()` function to check if delta spec directory exists
- [ ] 2.3 Update `parseTasksMd()` to track section-to-tasks mapping
- [ ] 2.4 Add `splitTasksByCapability()` function to partition tasks into root vs child files
- [ ] 2.5 Update `writeTasksJSONC()` to handle hierarchical structure
- [ ] 2.6 Add `writeChildTasksJSONC()` function with abbreviated header
- [ ] 2.7 Update `processChange()` to orchestrate auto-split when delta specs exist
- [ ] 2.8 Add unit tests for section name normalization
- [ ] 2.9 Add unit tests for delta spec matching
- [ ] 2.10 Add integration tests for auto-split generation

## 3. Tasks Command Implementation
- [ ] 3.1 Create `cmd/tasks.go` with TasksCmd struct
- [ ] 3.2 Implement `Run()` method with change ID resolution (reuse from accept)
- [ ] 3.3 Add `--flatten` flag for merged view
- [ ] 3.4 Add `--json` flag for JSON output
- [ ] 3.5 Implement `displaySummary()` for section-by-section progress
- [ ] 3.6 Implement `flattenTasks()` to merge hierarchical tasks
- [ ] 3.7 Implement `resolveChildFiles()` to process children refs and includes globs
- [ ] 3.8 Add TasksCmd to CLI struct in `cmd/root.go`
- [ ] 3.9 Add unit tests for summary display logic
- [ ] 3.10 Add unit tests for task flattening
- [ ] 3.11 Add integration tests for tasks command

## 4. Validation Updates
- [ ] 4.1 Add `validateChildReferences()` function to check children refs exist
- [ ] 4.2 Add `validateChildFormat()` function to check child file structure
- [ ] 4.3 Add circular reference detection in validation
- [ ] 4.4 Update archive validation to handle hierarchical tasks
- [ ] 4.5 Add unit tests for child reference validation
- [ ] 4.6 Add unit tests for circular reference detection

## 5. Reading and Parsing
- [ ] 5.1 Add `ReadTasksFile()` function to `internal/parsers/` for reading tasks.jsonc
- [ ] 5.2 Add version detection logic (handle version 1 and 2)
- [ ] 5.3 Add `ExpandHierarchicalTasks()` to resolve refs and includes
- [ ] 5.4 Add summary computation from hierarchical structure
- [ ] 5.5 Add unit tests for version detection
- [ ] 5.6 Add unit tests for hierarchical expansion

## 6. Documentation and Cleanup
- [ ] 6.1 Update AGENTS.md with hierarchical tasks.jsonc documentation
- [ ] 6.2 Add example hierarchical structure to AGENTS.md
- [ ] 6.3 Update tasksJSONHeader constant for version 2 header
- [ ] 6.4 Add abbreviated header constant for child files
- [ ] 6.5 Run `go test ./...` to ensure all tests pass
- [ ] 6.6 Run `golangci-lint run` to ensure code quality
- [ ] 6.7 Manual test with existing `redesign-provider-architecture` change
