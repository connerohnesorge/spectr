# Tasks

## 1. Domain Types and Parsing

- [x] 1.1 Create `internal/domain/proposal.go` with ProposalMetadata and
      Dependency types

- [x] 1.2 Implement ParseProposalFrontmatter function to extract YAML from
      proposal.md

- [x] 1.3 Add unit tests for frontmatter parsing (valid, empty, malformed YAML)

- [x] 1.4 Handle edge cases: no frontmatter, empty requires/enables, missing
      reason

## 2. Archive Detection

- [x] 2.1 Add IsChangeArchived function to `internal/discovery/changes.go`

- [x] 2.2 Implement glob-based archive lookup with date prefix handling

- [x] 2.3 Add GetArchivedChangeIDs function to list all archived change IDs

- [x] 2.4 Add unit tests for archive detection

## 3. Dependency Validation

- [x] 3.1 Create `internal/validation/deps.go` with dependency validation
      functions

- [x] 3.2 Implement BuildDependencyGraph to construct DAG from all proposals

- [x] 3.3 Implement DetectCycles using topological sort or DFS coloring

- [x] 3.4 Implement ValidateDependencies to check requires against archive

- [x] 3.5 Add unit tests for cycle detection and dependency validation

## 4. Integrate with Validate Command

- [x] 4.1 Modify `cmd/validate.go` to parse proposal frontmatter

- [x] 4.2 Call dependency validation during change validation

- [x] 4.3 Emit warnings for unmet dependencies (non-blocking)

- [x] 4.4 Emit errors for circular dependencies (blocking)

- [x] 4.5 Add integration test for validate with dependencies

## 5. Integrate with Accept Command

- [x] 5.1 Modify `cmd/accept.go` to check dependencies before accepting

- [x] 5.2 Hard fail if any requires entry is not archived

- [x] 5.3 Display clear error message listing unmet dependencies

- [x] 5.4 Add integration test for accept with unmet dependencies

## 6. Graph Command

- [x] 6.1 Create `cmd/graph.go` with GraphCmd struct and flags

- [x] 6.2 Implement ASCII tree output (default)

- [x] 6.3 Implement DOT format output (`--dot` flag)

- [x] 6.4 Implement JSON output (`--json` flag)

- [x] 6.5 Register graph command in `cmd/root.go`

- [x] 6.6 Add unit and integration tests for graph command

## 7. Documentation and Specs

- [x] 7.1 Add delta spec for cli-interface (new graph command, modified
      accept/validate)

- [x] 7.2 Update AGENTS.md with chained proposals documentation

- [x] 7.3 Add example proposal.md with frontmatter to testdata
