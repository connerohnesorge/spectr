## 1. Domain Types and Parsing

- [ ] 1.1 Create `internal/domain/proposal.go` with ProposalMetadata and
      Dependency types
- [ ] 1.2 Implement ParseProposalFrontmatter function to extract YAML from
      proposal.md
- [ ] 1.3 Add unit tests for frontmatter parsing (valid, empty, malformed YAML)
- [ ] 1.4 Handle edge cases: no frontmatter, empty requires/enables, missing
      reason

## 2. Archive Detection

- [ ] 2.1 Add IsChangeArchived function to `internal/discovery/changes.go`
- [ ] 2.2 Implement glob-based archive lookup with date prefix handling
- [ ] 2.3 Add GetArchivedChangeIDs function to list all archived change IDs
- [ ] 2.4 Add unit tests for archive detection

## 3. Dependency Validation

- [ ] 3.1 Create `internal/validation/deps.go` with dependency validation
      functions
- [ ] 3.2 Implement BuildDependencyGraph to construct DAG from all proposals
- [ ] 3.3 Implement DetectCycles using topological sort or DFS coloring
- [ ] 3.4 Implement ValidateDependencies to check requires against archive
- [ ] 3.5 Add unit tests for cycle detection and dependency validation

## 4. Integrate with Validate Command

- [ ] 4.1 Modify `cmd/validate.go` to parse proposal frontmatter
- [ ] 4.2 Call dependency validation during change validation
- [ ] 4.3 Emit warnings for unmet dependencies (non-blocking)
- [ ] 4.4 Emit errors for circular dependencies (blocking)
- [ ] 4.5 Add integration test for validate with dependencies

## 5. Integrate with Accept Command

- [ ] 5.1 Modify `cmd/accept.go` to check dependencies before accepting
- [ ] 5.2 Hard fail if any requires entry is not archived
- [ ] 5.3 Display clear error message listing unmet dependencies
- [ ] 5.4 Add integration test for accept with unmet dependencies

## 6. Graph Command

- [ ] 6.1 Create `cmd/graph.go` with GraphCmd struct and flags
- [ ] 6.2 Implement ASCII tree output (default)
- [ ] 6.3 Implement DOT format output (`--dot` flag)
- [ ] 6.4 Implement JSON output (`--json` flag)
- [ ] 6.5 Register graph command in `cmd/root.go`
- [ ] 6.6 Add unit and integration tests for graph command

## 7. Documentation and Specs

- [ ] 7.1 Add delta spec for cli-interface (new graph command, modified
      accept/validate)
- [ ] 7.2 Update AGENTS.md with chained proposals documentation
- [ ] 7.3 Add example proposal.md with frontmatter to testdata
