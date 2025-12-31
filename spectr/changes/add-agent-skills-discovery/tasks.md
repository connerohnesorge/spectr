# Tasks: Add Agent Skills Discovery

## Phase 1: Core Discovery Package

- [ ] 1.1 Create `internal/discovery/doc.go` with package documentation
- [ ] 1.2 Create `internal/discovery/skills.go` with SkillMetadata struct
- [ ] 1.3 Implement `ParseSkillMetadata()` function with YAML frontmatter parsing
- [ ] 1.4 Implement `ListEmbeddedSkills()` function with directory walking
- [ ] 1.5 Implement `IsSkillInstalled()` function with filesystem checks
- [ ] 1.6 Create `internal/discovery/testdata/` directory structure
- [ ] 1.7 Add test fixtures: valid_skill, minimal_skill, invalid_yaml, no_frontmatter, missing_name
- [ ] 1.8 Create `internal/discovery/skills_test.go`
- [ ] 1.9 Write test: `TestParseSkillMetadata_ValidFrontmatter`
- [ ] 1.10 Write test: `TestParseSkillMetadata_MinimalFrontmatter`
- [ ] 1.11 Write test: `TestParseSkillMetadata_MissingRequiredField`
- [ ] 1.12 Write test: `TestParseSkillMetadata_InvalidYAML`
- [ ] 1.13 Write test: `TestParseSkillMetadata_NoFrontmatter`
- [ ] 1.14 Write test: `TestListEmbeddedSkills_MultipleSkills`
- [ ] 1.15 Write test: `TestListEmbeddedSkills_MissingSkillMd`
- [ ] 1.16 Write test: `TestListEmbeddedSkills_EmptyDirectory`
- [ ] 1.17 Write test: `TestIsSkillInstalled_Installed`
- [ ] 1.18 Write test: `TestIsSkillInstalled_NotInstalled`
- [ ] 1.19 Run discovery tests: `go test ./internal/discovery/...`
- [ ] 1.20 Verify test coverage >80% for discovery package

## Phase 2: List Integration

- [ ] 2.1 Add `SkillInfo` struct to `internal/list/types.go` with JSON tags
- [ ] 2.2 Add `ListSkills(skillFS fs.FS) ([]SkillInfo, error)` method to `internal/list/lister.go`
- [ ] 2.3 Implement SkillMetadata → SkillInfo conversion with installation status
- [ ] 2.4 Add test: `TestLister_ListSkills` with mock skillFS
- [ ] 2.5 Add test: `TestLister_ListSkills_EmptySkills` for empty skillFS
- [ ] 2.6 Add test: `TestLister_ListSkills_InstallationStatus` to verify installed field
- [ ] 2.7 Run list tests: `go test ./internal/list/...`

## Phase 3: CLI Command

- [ ] 3.1 Add `Skills bool` field to `ListCmd` struct in `cmd/list.go`
- [ ] 3.2 Add mutual exclusivity validation for `--skills` flag in `Run()` method
- [ ] 3.3 Implement `listSkills(lister *list.Lister) error` method
- [ ] 3.4 Implement `formatSkillsText(skills []SkillInfo) string` formatter
- [ ] 3.5 Implement `formatSkillsLong(skills []SkillInfo) string` formatter
- [ ] 3.6 Implement `formatSkillsJSON(skills []SkillInfo) (string, error)` formatter
- [ ] 3.7 Wire up `--skills` flag routing in `Run()` method
- [ ] 3.8 Add test: `TestListCmd_Skills_TextFormat`
- [ ] 3.9 Add test: `TestListCmd_Skills_LongFormat`
- [ ] 3.10 Add test: `TestListCmd_Skills_JSONFormat`
- [ ] 3.11 Add test: `TestListCmd_Skills_IncompatibleWithSpecs`
- [ ] 3.12 Add test: `TestListCmd_Skills_IncompatibleWithAll`
- [ ] 3.13 Add test: `TestListCmd_Skills_IncompatibleWithInteractive`
- [ ] 3.14 Run cmd tests: `go test ./cmd/...`

## Phase 4: Specification Delta

- [ ] 4.1 Create delta spec file: `spectr/changes/add-agent-skills-discovery/specs/provider-system/spec.md`
- [ ] 4.2 Add `## ADDED Requirements` section header
- [ ] 4.3 Add `### Requirement: Skill Discovery` with 4 scenarios:
  - [ ] 4.3.1 Scenario: List embedded skills with metadata
  - [ ] 4.3.2 Scenario: Parse skill frontmatter
  - [ ] 4.3.3 Scenario: Discover skills from embedded filesystem
  - [ ] 4.3.4 Scenario: Check skill installation status
- [ ] 4.4 Add `### Requirement: CLI Skill Listing` with 4 scenarios:
  - [ ] 4.4.1 Scenario: List skills in text format
  - [ ] 4.4.2 Scenario: List skills in long format
  - [ ] 4.4.3 Scenario: List skills in JSON format
  - [ ] 4.4.4 Scenario: Skills flag mutual exclusivity
- [ ] 4.5 Ensure all scenarios use WHEN/THEN format
- [ ] 4.6 Cross-reference existing AgentSkillsInitializer requirements (line 714)

## Phase 5: Validation and Testing

- [ ] 5.1 Run `spectr validate add-agent-skills-discovery`
- [ ] 5.2 Fix any validation errors in delta spec
- [ ] 5.3 Run full test suite: `nix develop -c tests`
- [ ] 5.4 Run linter: `nix develop -c lint`
- [ ] 5.5 Fix any lint errors or warnings
- [ ] 5.6 Manual test: `spectr list --skills` with embedded skills
- [ ] 5.7 Manual test: `spectr list --skills --long` output formatting
- [ ] 5.8 Manual test: `spectr list --skills --json | jq` valid JSON
- [ ] 5.9 Manual test: `spectr list --skills --specs` error handling
- [ ] 5.10 Verify total code size <500 lines (excluding tests)
- [ ] 5.11 Verify test coverage >80% for new code
- [ ] 5.12 Run acceptance test: Install a skill, verify `Installed: true` in output

## Dependencies

- Ensure `gopkg.in/yaml.v3` is in go.mod (already present)
- Ensure `github.com/spf13/afero` is in go.mod (already present)
- No new dependencies required

## Notes

- Tasks can be parallelized: 1.x, 2.x, and 4.x can be done concurrently
- Phase 3 depends on Phase 1 and 2 completion
- Phase 5 depends on all previous phases
- Each test task should verify actual behavior, not just pass
- Follow existing code style and patterns from internal/list and cmd/list packages
