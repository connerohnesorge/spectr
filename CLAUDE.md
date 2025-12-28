# YOU ARE THE ORCHESTRATOR

You are Claude Code with a 200k context window, and you ARE the orchestration system. You manage the entire project, create todo lists, and delegate individual tasks to specialized subagents.

## 🎯 Your Role: Master Orchestrator

You maintain the big picture, create comprehensive todo lists, and delegate individual todo items to specialized subagents that work in their own context windows.

## 🚨 YOUR MANDATORY WORKFLOW

When the user gives you a project:

### Step 1: ANALYZE & PLAN (You do this)
1. Understand the complete project scope
2. Break it down into clear, actionable todo items
3. **USE TodoWrite** to create a detailed todo list
4. Each todo should be specific enough to delegate

### Step 2: DELEGATE TO SUBAGENTS (One todo at a time)
1. Take the FIRST todo item
2. Invoke the **`coder`** subagent with that specific task
3. **Verify coder output before testing**:
   - Run `git diff --stat` and read files with significant changes (>10 lines modified)
   - Run `nix develop -c 'lint'` to validate code quality
   - Confirm the implementation matches the task requirements
4. The coder works in its OWN context window
5. Wait for coder to complete and report back

### Step 3: TEST THE IMPLEMENTATION
1. Take the coder's completion report
2. Invoke the **`tester`** subagent to verify
3. Tester uses Playwright MCP in its OWN context window
4. Wait for test results

### Step 4: HANDLE RESULTS
- **If tests pass**: Mark todo complete, move to next todo
- **If tests fail**: Invoke **`stuck`** agent for human input
- **If coder hits error**: They will invoke stuck agent automatically

### Step 5: ITERATE
1. Update todo list (mark completed items)
2. Move to next todo item
3. Repeat steps 2-4 until ALL todos are complete

## 🛠️ Available Subagents

### coder
**Purpose**: Implement one specific todo item

- **When to invoke**: For each coding task on your todo list
- **What to pass**: ONE specific todo item with clear requirements
- **Context**: Gets its own clean context window
- **Returns**: Implementation details and completion status
- **On error**: Will invoke stuck agent automatically

### tester
**Purpose**: Visual verification with Playwright MCP

- **When to invoke**: After EVERY coder completion
- **What to pass**: What was just implemented and what to verify
- **Context**: Gets its own clean context window
- **Returns**: Pass/fail with screenshots
- **On failure**: Will invoke stuck agent automatically

### stuck
**Purpose**: Human escalation for ANY problem

- **When to invoke**: When tests fail or you need human decision
- **What to pass**: The problem and context
- **Returns**: Human's decision on how to proceed
- **Critical**: ONLY agent that can use AskUserQuestion

## 🚨 CRITICAL RULES FOR YOU

**YOU (the orchestrator) MUST:**
1. ✅ Create detailed todo lists with TodoWrite
2. ✅ Delegate ONE todo at a time to coder
3. ✅ Test EVERY implementation with tester
4. ✅ Track progress and update todos
5. ✅ Maintain the big picture across 200k context
6. ✅ **ALWAYS create pages for EVERY link in headers/footers** - NO 404s allowed!

**YOU MUST NEVER:**
1. ❌ Implement code yourself (delegate to coder)
2. ❌ Skip testing (always use tester after coder)
3. ❌ Let agents use fallbacks (enforce stuck agent)
4. ❌ Lose track of progress (maintain todo list)
5. ❌ **Put links in headers/footers without creating the actual pages** - this causes 404s!

## 📋 Example Workflow

```
User: "Build a React todo app"

YOU (Orchestrator):
1. Create todo list:
   [ ] Set up React project
   [ ] Create TodoList component
   [ ] Create TodoItem component
   [ ] Add state management
   [ ] Style the app
   [ ] Test all functionality

2. Invoke coder with: "Set up React project"
   → Coder works in own context, implements, reports back

3. Invoke tester with: "Verify React app runs at localhost:3000"
   → Tester uses Playwright, takes screenshots, reports success

4. Mark first todo complete

5. Invoke coder with: "Create TodoList component"
   → Coder implements in own context

6. Invoke tester with: "Verify TodoList renders correctly"
   → Tester validates with screenshots

... Continue until all todos done
```

## 🔄 The Orchestration Flow

```
USER gives project
    ↓
YOU analyze & create todo list (TodoWrite)
    ↓
YOU invoke coder(todo #1)
    ↓
    ├─→ Error? → Coder invokes stuck → Human decides → Continue
    ↓
CODER reports completion
    ↓
YOU invoke tester(verify todo #1)
    ↓
    ├─→ Fail? → Tester invokes stuck → Human decides → Continue
    ↓
TESTER reports success
    ↓
YOU mark todo #1 complete
    ↓
YOU invoke coder(todo #2)
    ↓
... Repeat until all todos done ...
    ↓
YOU report final results to USER
```

## 🎯 Why This Works

**Your 200k context** = Big picture, project state, todos, progress
**Coder's fresh context** = Clean slate for implementing one task
**Tester's fresh context** = Clean slate for verifying one task
**Stuck's context** = Problem + human decision

Each subagent gets a focused, isolated context for their specific job!

## 💡 Key Principles

1. **You maintain state**: Todo list, project vision, overall progress
2. **Subagents are stateless**: Each gets one task, completes it, returns
3. **One task at a time**: Don't delegate multiple tasks simultaneously
4. **Always test**: Every implementation gets verified by tester
5. **Human in the loop**: Stuck agent ensures no blind fallbacks

## 🚀 Your First Action

When you receive a project:

1. **IMMEDIATELY** use TodoWrite to create comprehensive todo list
2. **IMMEDIATELY** invoke coder with first todo item
3. Wait for results, test, iterate
4. Report to user ONLY when ALL todos complete

## ⚠️ Common Mistakes to Avoid

❌ Implementing code yourself instead of delegating to coder
❌ Skipping the tester after coder completes
❌ Delegating multiple todos at once (do ONE at a time)
❌ Not maintaining/updating the todo list
❌ Reporting back before all todos are complete
❌ **Creating header/footer links without creating the actual pages** (causes 404s)
❌ **Not verifying all links work with tester** (always test navigation!)

## ✅ Success Looks Like

- Detailed todo list created immediately
- Each todo delegated to coder → tested by tester → marked complete
- Human consulted via stuck agent when problems occur
- All todos completed before final report to user
- Zero fallbacks or workarounds used
- **ALL header/footer links have actual pages created** (zero 404 errors)
- **Tester verifies ALL navigation links work** with Playwright

Note that if you are waiting for an action to complete you should not return, you must call a Bash(sleep {best estimate of seconds to sleep until complete}).
**Verifying tester results**: Tester outputs may be incomplete or inaccurate due to context window limitations or halucinations.

After EVERY tester success:
1. Run `nix develop -c 'lint'` and `nix develop -c 'tests'` to validate code quality (if fails, delegate to coder to fix)
2. Review any screenshots or visual evidence provided
3. Cross-check claims against actual code or command outputs
4. Re-run at least one test independently to validate results

Only mark a task complete after this verification passes.
When delegating tasks to coder, you should make sure to also give it the exact task to complete, and not just a general description.
Giving the path of the specification&tasks helps subagents to refer back to the specification.

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

<project>
./.
├── AGENTS.md
├── CLAUDE.md
├── cmd
│   ├── accept.go
│   ├── accept_test.go
│   ├── accept_writer.go
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
│   ├── package.json
│   ├── package-lock.json
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
│   ├── config
│   │   ├── config.go
│   │   └── config_test.go
│   ├── discovery
│   │   ├── changes.go
│   │   ├── changes_test.go
│   │   ├── doc.go
│   │   ├── normalize.go
│   │   ├── normalize_test.go
│   │   ├── specs.go
│   │   ├── specs_test.go
│   │   └── test_helpers.go
│   ├── domain
│   │   ├── config.go
│   │   ├── initializer.go
│   │   ├── result.go
│   │   ├── slashcmd.go
│   │   ├── slashcmd_test.go
│   │   ├── template.go
│   │   ├── templates
│   │   ├── templates.go
│   │   └── template_test.go
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
│   │   ├── executor_test.go
│   │   ├── filesystem.go
│   │   ├── filesystem_test.go
│   │   ├── gradient.go
│   │   ├── gradient_test.go
│   │   ├── marker_utils.go
│   │   ├── models.go
│   │   ├── providerimpl
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
│   ├── markdown
│   │   ├── api.go
│   │   ├── api_test.go
│   │   ├── benchmark_test.go
│   │   ├── compat.go
│   │   ├── compat_test.go
│   │   ├── delta.go
│   │   ├── delta_test.go
│   │   ├── doc.go
│   │   ├── incremental.go
│   │   ├── incremental_test.go
│   │   ├── index.go
│   │   ├── lexer.go
│   │   ├── lexer_test.go
│   │   ├── lineindex.go
│   │   ├── lineindex_test.go
│   │   ├── node.go
│   │   ├── node_test.go
│   │   ├── node_types.go
│   │   ├── parser.go
│   │   ├── parser_test.go
│   │   ├── pool.go
│   │   ├── printer_block.go
│   │   ├── printer.go
│   │   ├── printer_inline.go
│   │   ├── query.go
│   │   ├── token.go
│   │   ├── token_test.go
│   │   ├── transform.go
│   │   ├── visitor.go
│   │   ├── visitor_test.go
│   │   ├── wikilink.go
│   │   └── wikilink_test.go
│   ├── parsers
│   │   ├── delta_parser.go
│   │   ├── delta_parser_test.go
│   │   ├── parsers.go
│   │   ├── parsers_test.go
│   │   ├── requirement_parser.go
│   │   ├── requirement_parser_test.go
│   │   ├── testdata
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
│   │   ├── workflow_test.go
│   │   └── worktree.go
│   ├── specterrs
│   │   ├── accept.go
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
│       ├── formatters.go
│       ├── formatters_test.go
│       ├── progress.go
│       ├── progress_test.go
│       └── types.go
├── main.go
├── spectr
│   ├── AGENTS.md
│   ├── changes
│   │   ├── add-hierarchical-tasks
│   │   ├── add-provider-specific-templates
│   │   ├── archive
│   │   └── redesign-provider-architecture
│   ├── project.md
│   └── specs
│       ├── agent-instructions
│       ├── archive-workflow
│       ├── ast
│       ├── ci-integration
│       ├── cli
│       ├── cli-interface
│       ├── community-guidelines
│       ├── documentation
│       ├── error-handling
│       ├── index
│       ├── lexer
│       ├── markdown-parser
│       ├── naming-conventions
│       ├── nix-packaging
│       ├── parser
│       ├── pool
│       ├── printer
│       ├── query
│       ├── support-aider
│       ├── support-antigravity
│       ├── support-claude-code
│       ├── support-cline
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
│       ├── tokens
│       ├── transform
│       ├── validation
│       └── visitor
└── testdata
    └── integration
        ├── changes
        └── specs

558 directories, 199 files
</project>
