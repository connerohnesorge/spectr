# Tasks for Adding Kimi Agent Skills Support

## Implementation Tasks

- [ ] Create Kimi provider implementation (`internal/initialize/providers/kimi.go`)
- [ ] Register Kimi provider in registry (`internal/initialize/providers/registry.go`)
- [ ] Create Kimi AGENTS.md template (`internal/initialize/templates/agents-kimi.md`)
- [ ] Test provider initialization with `spectr init`
- [ ] Validate directory creation: `.kimi/skills` and `.kimi/commands`
- [ ] Validate AGENTS.md file creation
- [ ] Validate agent skills installation
- [ ] Test slash commands functionality

## Validation Tasks

- [ ] Run `spectr validate add-kimi-agentskills-support`
- [ ] Ensure no linting errors: `nix develop -c lint`
- [ ] Run tests: `nix develop -c tests`
- [ ] Verify skills are executable
- [ ] Test in a sample project directory

## Documentation Tasks

- [ ] Update provider documentation if needed
- [ ] Add Kimi to supported providers list
- [ ] Document Kimi-specific configuration

## Success Criteria Checklist

- [ ] `spectr init` creates `.kimi/skills` directory
- [ ] `spectr init` creates `.kimi/commands` directory
- [ ] `AGENTS.md` file is created with Kimi instructions
- [ ] Agent skills are installed and executable
- [ ] Slash commands `spectr-proposal` and `spectr-apply` work
- [ ] All tests pass
- [ ] Validation passes without errors
