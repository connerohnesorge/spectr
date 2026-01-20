# Tasks for Adding Gemini Agent Skills Support

## 1. Implementation

- [ ] 1.1 Create GEMINI.md template (`internal/initialize/templates/gemini.md`)
      with Spectr instruction pointer and managed markers
- [ ] 1.2 Update Gemini provider to add `.gemini/skills/` directory initializer
- [ ] 1.3 Update Gemini provider to add GEMINI.md config file initializer
- [ ] 1.4 Update Gemini provider to add spectr-accept-wo-spectr-bin skill
      initializer
- [ ] 1.5 Update Gemini provider to add spectr-validate-wo-spectr-bin skill
      initializer
- [ ] 1.6 Ensure initializers are ordered correctly (directories before skills)

## 2. Testing

- [ ] 2.1 Test `spectr init` creates `.gemini/skills/` directory
- [ ] 2.2 Test skills are installed with correct structure (SKILL.md +
      scripts/)
- [ ] 2.3 Test scripts have executable permissions (0755)
- [ ] 2.4 Test GEMINI.md is created with correct content
- [ ] 2.5 Test existing TOML slash commands still work
- [ ] 2.6 Test idempotent initialization (run twice, no errors)
- [ ] 2.7 Run full test suite: `nix develop -c tests`

## 3. Validation

- [ ] 3.1 Run `spectr validate add-gemini-agentskills-support`
- [ ] 3.2 Run linter: `nix develop -c lint`
- [ ] 3.3 Verify SKILL.md frontmatter matches Agent Skills spec (name,
      description fields)
- [ ] 3.4 Test in a fresh project directory

## 4. Success Criteria Verification

- [ ] 4.1 `spectr init` creates `.gemini/skills/` with both skills
- [ ] 4.2 `GEMINI.md` file is created in project root
- [ ] 4.3 Existing TOML commands in `.gemini/commands/spectr/` preserved
- [ ] 4.4 Skills include SKILL.md with valid frontmatter
- [ ] 4.5 Scripts are executable (scripts/accept.sh, scripts/validate.sh)
- [ ] 4.6 All tests pass
- [ ] 4.7 Validation passes without errors
