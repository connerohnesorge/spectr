# Implementation Tasks

## 1. Template Infrastructure

- [ ] 1.1 Create `internal/domain/templates/skill-proposal.md.tmpl` with Amp frontmatter and proposal creation instructions
- [ ] 1.2 Create `internal/domain/templates/skill-apply.md.tmpl` with Amp frontmatter and apply instructions
- [ ] 1.3 Update `internal/initialize/templates.go` embed directive to include `skill-*.md.tmpl`
- [ ] 1.4 Add `ProposalSkill() domain.TemplateRef` method to TemplateManager
- [ ] 1.5 Add `ApplySkill() domain.TemplateRef` method to TemplateManager

## 2. SkillFileInitializer

- [ ] 2.1 Create `internal/initialize/providers/skillfile.go` with `SkillFileInitializer` struct
- [ ] 2.2 Implement `NewSkillFileInitializer(targetPath, template)` constructor
- [ ] 2.3 Implement `Init(ctx, projectFs, homeFs, cfg, tm)` method
- [ ] 2.4 Implement `IsSetup(projectFs, homeFs, cfg)` method
- [ ] 2.5 Implement `dedupeKey()` method returning `SkillFileInitializer:<targetPath>`
- [ ] 2.6 Add SkillFileInitializer to initializer type priority ordering

## 3. Amp Provider

- [ ] 3.1 Create `internal/initialize/providers/amp.go` with `AmpProvider` struct
- [ ] 3.2 Implement `Initializers(ctx, tm)` method returning:
  - DirectoryInitializer for `.agents/skills/spectr-proposal/`
  - DirectoryInitializer for `.agents/skills/spectr-apply/`
  - SkillFileInitializer for `.agents/skills/spectr-proposal/SKILL.md`
  - SkillFileInitializer for `.agents/skills/spectr-apply/SKILL.md`
  - AgentSkillsInitializer for `spectr-accept-wo-spectr-bin`
  - AgentSkillsInitializer for `spectr-validate-wo-spectr-bin`
- [ ] 3.3 Register Amp provider in `internal/initialize/providers/registry.go` with priority 15

## 4. Testing

- [ ] 4.1 Add unit tests for `SkillFileInitializer` in `skillfile_test.go`
- [ ] 4.2 Add unit tests for `AmpProvider.Initializers()` in `amp_test.go`
- [ ] 4.3 Add integration test verifying `.agents/skills/` directory structure
- [ ] 4.4 Test SKILL.md frontmatter parsing and validation
- [ ] 4.5 Test embedded skill copying (spectr-accept-wo-spectr-bin, spectr-validate-wo-spectr-bin)
- [ ] 4.6 Test template variable substitution in skill content
- [ ] 4.7 Test deduplication when multiple providers generate same skills

## 5. Documentation

- [ ] 5.1 Create `spectr/specs/support-amp/spec.md` (merge from delta)
- [ ] 5.2 Update `spectr/specs/provider-system/spec.md` (merge delta changes)
- [ ] 5.3 Update `spectr/specs/agent-instructions/spec.md` (merge delta changes)
- [ ] 5.4 Add Amp to provider list in README.md
- [ ] 5.5 Add Amp example to initialization wizard documentation

## 6. Validation

- [ ] 6.1 Run `spectr validate add-amp-agentskills-support`
- [ ] 6.2 Run `nix develop -c lint` to check Go linting
- [ ] 6.3 Run `nix develop -c tests` to verify all tests pass
- [ ] 6.4 Run `spectr init` and select Amp provider to verify end-to-end workflow
- [ ] 6.5 Verify generated SKILL.md files have correct frontmatter format
- [ ] 6.6 Verify `/spectr:proposal` can be invoked in Amp (manual test)
