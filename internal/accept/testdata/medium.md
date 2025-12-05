## 1. Foundation

- [x] 1.1 Create `internal/discovery/` package structure
- [x] 1.2 Create `internal/parsers/` package structure
- [x] 1.3 Create `internal/list/` package structure
- [x] 1.4 Define data types in `internal/list/types.go` for ChangeInfo and SpecInfo

## 2. Discovery Implementation

- [x] 2.1 Implement `discovery.GetActiveChanges()` to find changes in `spectr/changes/`
- [x] 2.2 Implement `discovery.GetSpecs()` to find specs in `spectr/specs/`
- [x] 2.3 Add filtering logic to exclude `archive/` directory
- [x] 2.4 Add filtering logic to exclude hidden directories (starting with `.`)
- [x] 2.5 Validate that discovered items contain required files (`proposal.md` or `spec.md`)

## 3. Parser Implementation

- [x] 3.1 Implement `parsers.ExtractTitle()` to parse first H1 heading from markdown
- [x] 3.2 Add logic to remove "Change:" and "Spec:" prefixes from titles
- [x] 3.3 Implement `parsers.CountTasks()` to parse `tasks.md` and count completion
- [x] 3.4 Add regex pattern matching for `- [ ]` and `- [x]` task markers (case-insensitive)
- [x] 3.5 Implement `parsers.CountDeltas()` to count delta sections in change specs
- [x] 3.6 Implement `parsers.CountRequirements()` to count requirements in specs

## 4. List Command Implementation

- [x] 4.1 Add `ListCmd` struct to `cmd/root.go` with Kong tags
- [x] 4.2 Create `cmd/list.go` with `ListCmd.Run()` method
- [x] 4.3 Add `--specs` flag to switch between changes and specs
- [x] 4.4 Add `--long` flag for detailed output
- [x] 4.5 Add `--json` flag for JSON output
- [x] 4.6 Implement logic to collect change information (ID, title, deltaCount, taskStatus)
- [x] 4.7 Implement logic to collect spec information (ID, title, requirementCount)

## 5. Formatter Implementation

- [x] 5.1 Implement `formatters.FormatChangesText()` for default text output (IDs only)
- [x] 5.2 Implement `formatters.FormatChangesLong()` for detailed text output
- [x] 5.3 Implement `formatters.FormatChangesJSON()` for JSON output
- [x] 5.4 Implement `formatters.FormatSpecsText()` for default text output (IDs only)
- [x] 5.5 Implement `formatters.FormatSpecsLong()` for detailed text output
- [x] 5.6 Implement `formatters.FormatSpecsJSON()` for JSON output
- [x] 5.7 Add sorting logic to ensure alphabetical order by ID
- [x] 5.8 Implement "No items found" message for empty results

## 6. Testing

- [x] 6.1 Write unit tests for `discovery.GetActiveChanges()`
- [x] 6.2 Write unit tests for `discovery.GetSpecs()`
- [x] 6.3 Write unit tests for `parsers.ExtractTitle()` with various input formats
- [x] 6.4 Write unit tests for `parsers.CountTasks()` with different task states
- [x] 6.5 Write integration test for `spectr list` (default changes output)
- [x] 6.6 Write integration test for `spectr list --specs`
- [x] 6.7 Write integration test for `spectr list --long`
- [x] 6.8 Write integration test for `spectr list --json`
- [x] 6.9 Test empty results scenarios (no changes/specs found)
- [x] 6.10 Test with archived changes to ensure they're excluded

## 7. Documentation and Polish

- [x] 7.1 Add command help text via Kong struct tags
- [x] 7.2 Test help output with `spectr list --help`
- [x] 7.3 Verify output formatting matches OpenSpec reference behavior
- [x] 7.4 Handle errors gracefully (missing directories, permission issues)
- [x] 7.5 Verify sorting is consistent across all output formats
