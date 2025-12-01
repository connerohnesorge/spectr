# Why Replace `/archive` with `/sync`?

The `/spectr:archive` slash command is being replaced by `/spectr:sync` to support a more flexible, code-first workflow and leverage the agent's reasoning capabilities.

## 1. Supporting Code-First Development
The `archive` command assumed a strict Spec-Driven Development (SDD) lifecycle:
1. Define Spec Change
2. Implement Change
3. Archive (Merge) Change to Spec

However, many developers work "Code First":
1. Modify Code
2. Update Spec to reflect Code

The `/sync` command acknowledges that **Code can be understood as ultimate Source of Truth**. It allows the agent to look at the current state of the codebase and "sync" the specs to match, regardless of whether a formal "Change" object was used. This bridges the gap for users who may not want to strictly follow the formal SDD lifecycle for every modification.

## 2. `archive` is a Simple CLI Wrapper
The old `/archive` slash command essentially just ran `spectr archive <id>`. This is a trivial operation that doesn't require a complex prompt. A user can simply ask the agent "archive change 123" or run the CLI command themselves.

## 3. `sync` Requires Agent Intelligence
"Syncing" specs with code is a complex, qualitative task that fits the agent's strengths:
- **Drift Detection:** The agent analyzes the implementation to see where it diverges from the documentation.
- **Interactive Review:** The agent can propose spec updates based on code behavior, which is harder to automate with a simple CLI tool.
- **Safety:** It ensures specs remain accurate documentation without enforcing a rigid process.

## Summary
Moving to `/sync` decouples the agent's assistance from the strict "Change" object lifecycle, making Spectr useful for maintaining documentation in *any* development workflow, while removing a redundant wrapper around a simple CLI command.
