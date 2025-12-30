# Tasks

## 1. Create Skill Directory Structure

- [ ] 1.1 Create `internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/` directory
- [ ] 1.2 Create `SKILL.md` with AgentSkills frontmatter and usage documentation
- [ ] 1.3 Create `scripts/` directory for shell scripts

## 2. Implement Core Validation Script

- [ ] 2.1 Create `scripts/validate.sh` with argument parsing (change-id, spec-id, or --all)
- [ ] 2.2 Implement spec file validation (Requirements section detection)
- [ ] 2.3 Implement requirement validation (SHALL/MUST keyword check)
- [ ] 2.4 Implement scenario validation (`#### Scenario:` format detection)
- [ ] 2.5 Implement malformed scenario detection (wrong header levels, bullets)
- [ ] 2.6 Implement change delta validation (ADDED/MODIFIED/REMOVED/RENAMED sections)
- [ ] 2.7 Implement tasks.md validation (at least one task item)
- [ ] 2.8 Implement discovery functions (find all specs, find all changes)

## 3. Implement Output Formatting

- [ ] 3.1 Add human-readable output with file paths and line numbers
- [ ] 3.2 Add color-coded error levels (ERROR in red, WARNING in yellow)
- [ ] 3.3 Add summary line (passed/failed counts)
- [ ] 3.4 Add --json flag for machine-readable output

## 4. Register with Claude Code Provider

- [ ] 4.1 Add `NewAgentSkillsInitializer` call for the new skill in `claude.go`

## 5. Testing and Verification

- [ ] 5.1 Test validation against existing specs in spectr/specs/
- [ ] 5.2 Test validation against existing changes in spectr/changes/
- [ ] 5.3 Verify skill is installed when running `spectr init`
- [ ] 5.4 Compare output with `spectr validate --all` to ensure behavioral parity
