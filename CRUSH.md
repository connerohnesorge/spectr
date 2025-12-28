<!-- spectr:start -->
# Spectr Instructions

These instructions are for AI assistants working in this project.

## Critical: Before Creating Delta Specs

**MANDATORY PRE-FLIGHT CHECKLIST** - Follow this BEFORE writing any `## ADDED/MODIFIED/REMOVED Requirements`:

### 1. Read the Base Spec First
- If using `## MODIFIED Requirements`, you MUST read `spectr/specs/<capability>/spec.md` FIRST
- Verify the exact requirement name exists in the base spec
- Copy the FULL requirement block (requirement + all scenarios)
- Only then paste into your delta spec and modify

### 2. Choose ADDED vs MODIFIED Correctly
- **ADDED**: New requirement that doesn't exist in base spec
- **MODIFIED**: Existing requirement you're changing
- **Rule**: If you haven't read the base spec yet, you CANNOT use MODIFIED

### 3. Validate Before Submission
- Every `## ADDED/MODIFIED/REMOVED Requirements` section MUST have at least one requirement
- Every requirement MUST have at least one `#### Scenario:` (4 hashtags, not bullets)
- MODIFIED requirements MUST match names in base spec exactly (case-insensitive)
- Run `spectr validate <change-id>` before marking complete

### Common Validation Errors to Avoid

**"requirement does not exist in base spec"**
- **Cause**: Used `## MODIFIED Requirements` for a requirement that doesn't exist
- **Fix**: Use `## ADDED Requirements` instead, OR verify spelling matches base spec exactly

**"Requirements section is empty (no requirements found)"**
- **Cause**: Created section header (`## ADDED Requirements`) but forgot to add requirements
- **Fix**: Remove empty sections OR add at least one requirement with scenario

**"Requirement must have at least one scenario"**
- **Cause**: Requirement exists but has no `#### Scenario:` blocks
- **Fix**: Add at least one scenario with WHEN/THEN structure

## ADDED vs MODIFIED Decision Tree

Before writing delta specs:

1. **Does this requirement exist in the base spec?**
   - YES → Read `spectr/specs/<capability>/spec.md` and find it
     - Found exact match? → Use `## MODIFIED Requirements`
     - Not found? → Check spelling, then use `## ADDED Requirements`
   - NO/UNSURE → Use `## ADDED Requirements` (safer default)

2. **If using MODIFIED:**
   - Read base spec: `spectr/specs/<capability>/spec.md`
   - Copy FULL requirement block (header + description + all scenarios)
   - Paste into delta spec under `## MODIFIED Requirements`
   - Edit to reflect new behavior
   - Keep at least one `#### Scenario:`

## Opening the Full Agent Guide

Always open `@/spectr/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/spectr/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

## Task Delegation Paths

When delegating tasks from a change proposal to subagents:
- Provide the proposal path: `spectr/changes/<id>/proposal.md`
- Include task context: `spectr/changes/<id>/tasks.jsonc`
- Reference delta specs: `spectr/changes/<id>/specs/<capability>/spec.md`

## Quick Validation Commands

Before delegating tasks or marking complete:

```bash
# Validate a specific change
spectr validate <change-id>

# Read base spec to verify requirement names
cat spectr/specs/<capability>/spec.md | grep "### Requirement:"

# List all capabilities to find the right one
ls spectr/specs/
```

**Remember**: MODIFIED requires the requirement to exist in base spec. When in doubt, use ADDED.

<!-- spectr:end -->

./.
├── cmd
│   ├── accept.go
│   ├── accept_test.go
│   ├── completion.go
│   ├── init.go
│   ├── init_test.go
│   ├── list.go
│   ├── list_test.go
│   ├── pr.go
│   ├── pr_test.go
│   ├── root.go
│   ├── validate.go
│   ├── version.go
│   ├── view.go
│   └── view_test.go
├── CODE_OF_CONDUCT.md
├── CRUSH.md
├── docs
│   ├── astro.config.mjs
│   ├── bun.lock
│   ├── dist
│   │   ├── 404.html
│   │   ├── _astro
│   │   ├── changelog
│   │   ├── concepts
│   │   ├── favicon.svg
│   │   ├── getting-started
│   │   ├── guides
│   │   ├── index.html
│   │   ├── index.md
│   │   ├── llms-full.txt
│   │   ├── llms-small.txt
│   │   ├── llms.txt
│   │   ├── pagefind
│   │   ├── reference
│   │   ├── sitegraph
│   │   ├── sitemap-0.xml
│   │   ├── sitemap-index.xml
│   │   ├── warp
│   │   └── warp.xml
│   ├── public
│   │   └── favicon.svg
│   ├── README.md
│   ├── spectr
│   │   └── changes
│   ├── src
│   │   ├── assets
│   │   ├── content
│   │   └── content.config.ts
│   ├── tsconfig.json
│   └── uno.config.ts
├── examples
│   ├── archive
│   │   └── spectr
│   ├── init
│   ├── list
│   │   └── spectr
│   ├── partial-match
│   │   └── spectr
│   └── validate
│       ├── broken
│       ├── fixed
│       └── spectr
├── flake.lock
├── flake.nix
├── go.mod
├── go.sum
├── internal
│   ├── archive
│   │   ├── archiver.go
│   │   ├── archiver_test.go
│   │   ├── cmd.go
│   │   ├── constants.go
│   │   ├── interactive_bridge.go
│   │   ├── spec_merger.go
│   │   ├── spec_merger_test.go
│   │   ├── types.go
│   │   ├── validator.go
│   │   └── validator_test.go
│   ├── discovery
│   │   ├── changes.go
│   │   ├── changes_test.go
│   │   ├── doc.go
│   │   ├── specs.go
│   │   ├── specs_test.go
│   │   └── test_helpers.go
│   ├── git
│   │   ├── branch.go
│   │   ├── doc.go
│   │   ├── platform.go
│   │   ├── platform_test.go
│   │   ├── worktree.go
│   │   └── worktree_test.go
│   ├── initialize
│   │   ├── constants.go
│   │   ├── executor.go
│   │   ├── filesystem.go
│   │   ├── filesystem_test.go
│   │   ├── gradient.go
│   │   ├── gradient_test.go
│   │   ├── marker_utils.go
│   │   ├── models.go
│   │   ├── providers
│   │   ├── templates
│   │   ├── templates.go
│   │   ├── templates_test.go
│   │   ├── wizard.go
│   │   └── wizard_test.go
│   ├── list
│   │   ├── formatters.go
│   │   ├── formatters_test.go
│   │   ├── formatters_unified.go
│   │   ├── interactive.go
│   │   ├── interactive_test.go
│   │   ├── lister.go
│   │   ├── lister_test.go
│   │   ├── types.go
│   │   └── types_test.go
│   ├── parsers
│   │   ├── delta_parser.go
│   │   ├── delta_parser_test.go
│   │   ├── parsers.go
│   │   ├── parsers_test.go
│   │   ├── requirement_parser.go
│   │   ├── requirement_parser_test.go
│   │   └── types.go
│   ├── pr
│   │   ├── doc.go
│   │   ├── dryrun.go
│   │   ├── helpers.go
│   │   ├── integration_test.go
│   │   ├── platforms.go
│   │   ├── templates.go
│   │   ├── templates_test.go
│   │   ├── workflow.go
│   │   └── workflow_test.go
│   ├── specterrs
│   │   ├── archive.go
│   │   ├── doc.go
│   │   ├── environment.go
│   │   ├── git.go
│   │   ├── initialize.go
│   │   ├── list.go
│   │   ├── pr.go
│   │   └── validation.go
│   ├── tui
│   │   ├── helpers.go
│   │   ├── helpers_test.go
│   │   ├── menu.go
│   │   ├── menu_test.go
│   │   ├── styles.go
│   │   ├── table.go
│   │   ├── table_test.go
│   │   └── types.go
│   ├── validation
│   │   ├── change_rules.go
│   │   ├── change_rules_test.go
│   │   ├── constants.go
│   │   ├── delta_validators.go
│   │   ├── formatters.go
│   │   ├── formatters_test.go
│   │   ├── helpers.go
│   │   ├── helpers_test.go
│   │   ├── integration_base_spec_test.go
│   │   ├── interactive.go
│   │   ├── interactive_test.go
│   │   ├── items.go
│   │   ├── items_test.go
│   │   ├── parser.go
│   │   ├── parser_test.go
│   │   ├── spec_rules.go
│   │   ├── spec_rules_test.go
│   │   ├── test_line_numbers_test.go
│   │   ├── types.go
│   │   ├── types_test.go
│   │   ├── validator.go
│   │   └── validator_test.go
│   ├── version
│   │   └── version.go
│   └── view
│       ├── dashboard.go
│       ├── dashboard_test.go
│       ├── formatters_demo_test.go
│       ├── formatters.go
│       ├── formatters_test.go
│       ├── progress_demo_test.go
│       ├── progress.go
│       ├── progress_test.go
│       └── types.go
├── LICENSE
├── main.go
├── README.md
├── spectr
│   ├── AGENTS.md
│   ├── changes
│   │   ├── add-crush-support
│   │   ├── add-provider-search-filter
│   │   ├── add-stdout-flag-list
│   │   ├── archive
│   │   └── refactor-agents-md-injection
│   ├── project.md
│   └── specs
│       ├── agent-instructions
│       ├── archive-workflow
│       ├── ci-integration
│       ├── cli-framework
│       ├── cli-interface
│       ├── community-guidelines
│       ├── documentation
│       ├── error-handling
│       ├── naming-conventions
│       ├── nix-packaging
│       ├── support-aider
│       ├── support-antigravity
│       ├── support-claude-code
│       ├── support-cline
│       ├── support-codebuddy
│       ├── support-codex
│       ├── support-continue
│       ├── support-costrict
│       ├── support-crush
│       ├── support-cursor
│       ├── support-gemini
│       ├── support-kilocode
│       ├── support-opencode
│       ├── support-qoder
│       ├── support-qwen
│       ├── support-windsurf
│       └── validation
└── testdata
    └── integration
        ├── changes
        └── specs

543 directories, 153 files
