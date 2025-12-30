## 1. Core Implementation

- [ ] 1.1 Create `AgentSkillsInitializer` type in
  `internal/initialize/providers/agentskills.go`
  - Implement `Initializer` interface
  - Accept skill name and target directory as constructor parameters
  - Copy all files from embedded skill directory to target
  - Implement `dedupeKey()` for deduplication
  - Implement `IsSetup()` to check if skill already exists

- [ ] 1.2 Add embedded skill directory structure at
  `internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/`
  - Create `SKILL.md` with valid AgentSkills frontmatter
  - Create `scripts/accept.sh` with tasks.md to tasks.jsonc conversion logic

- [ ] 1.3 Update `TemplateManager` in `internal/initialize/templates.go`
  - Add `//go:embed templates/skills/**/*` directive
  - Add `SkillFS(name string) (fs.FS, error)` method to retrieve skill files
  - Handle skill not found error case

## 2. Provider Integration

- [ ] 2.1 Update `ClaudeProvider.Initializers()` in
  `internal/initialize/providers/claude.go`
  - Add `NewDirectoryInitializer(".claude/skills")` for skills directory
  - Add `NewAgentSkillsInitializer("spectr-accept-wo-spectr-bin",
    ".claude/skills/spectr-accept-wo-spectr-bin", tm)`
  - Maintain existing initializer order (directories first)

## 3. Testing

- [ ] 3.1 Write unit tests for `AgentSkillsInitializer`
  - Test successful skill copy to empty directory
  - Test idempotency (re-running produces same result)
  - Test `IsSetup()` returns true when skill exists
  - Test `dedupeKey()` format

- [ ] 3.2 Write integration test for Claude provider
  - Verify skill directory created at correct path
  - Verify SKILL.md content matches embedded template
  - Verify scripts/accept.sh is executable

- [ ] 3.3 Test accept.sh script manually
  - Create sample tasks.md
  - Run script and verify valid tasks.jsonc output
  - Test edge cases (empty tasks, nested sections)

## 4. Validation & Cleanup

- [ ] 4.1 Run `nix develop -c lint` and fix any issues
- [ ] 4.2 Run `nix develop -c tests` and ensure all tests pass
- [ ] 4.3 Update provider-system spec delta if design changes during
  implementation
