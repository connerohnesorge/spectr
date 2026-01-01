# PR Package

Creates pull requests from changes using git worktree isolation. Multi-platform support.

## OVERVIEW
`spectr pr archive` creates PR for completed change. Creates isolated git worktree, commits archive changes, pushes, creates PR via platform CLI (gh, glab, tea, or manual URL). Cleans up worktree after completion.

## STRUCTURE
```text

internal/pr/
├── platforms.go         # Platform detection and CLI invocation
├── helpers.go           # Git worktree operations
├── dryrun.go           # Preview mode logic
├── doc.go              # Package documentation
└── *_test.go           # Integration tests
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Archive PR flow | cmd/pr.go (embeds ArchiveCmd) | Main orchestration |
| New proposal PR | cmd/pr.go (embeds NewCmd) | Proposal review PR |
| Platform detection | platforms.go | GitHub, GitLab, Gitea, Bitbucket |
| Worktree operations | helpers.go | Create, cleanup, commit |

## CONVENTIONS
- **Isolated worktree**: Never modify user's working directory
- **Auto cleanup**: Always remove worktree (even on error)
- **Platform auto-detect**: Use available CLI (gh > glab > tea > manual)
- **Structured commits**: Conventional commits with spectr metadata

## SUPPORTED PLATFORMS

| Platform | CLI | Commands |
|----------|-----|----------|
| GitHub | gh | `gh pr create --base <base> --title <title> --body <body>` |
| GitLab | glab | `glab mr create --source <branch> --target <base>` |
| Gitea/Forgejo | tea | `tea pr create --head <branch> --base <base>` |
| Bitbucket | Manual | Print instructions, user creates PR in UI |

## UNIQUE PATTERNS
- **Git worktree isolation**: `git worktree add` creates separate working copy
- **Commit message templating**: Uses Go templates for both archive and proposal modes
- **PR body generation**: Structured markdown with change details, affected specs, task summary

## ANTI-PATTERNS
- **NO manual git operations**: Use helpers.go for worktree management
- **DON'T leave worktrees**: Always cleanup after PR creation (or on error)
- **NO force pushes**: Use safe git operations, avoid `--force`

## KEY FUNCTIONS
- `ArchiveCmd.Run() error` - Archive + create PR workflow
- `NewCmd.Run() error` - Proposal review PR workflow
- `CreateWorktree(base, branch) (worktreePath, cleanup, error)` - Create isolated worktree
- `DetectPlatform() Platform` - Find available Git platform CLI
- `CreatePR(platform, options) error` - Invoke platform CLI or print manual instructions

## COMMIT MESSAGE TEMPLATE
```
<type>: <short summary>

<optional body>

<optional footer>

Spectr-Change: <change-id>
Spectr-Specs: <affected-specs>
```

## FLOW (Archive)
1. Validate change with `spectr validate`
2. Execute archive in worktree: `spectr archive <id>`
3. Commit with structured message
4. Push branch to remote
5. Detect platform, invoke CLI for PR creation
6. Cleanup worktree
7. Print PR URL

## FLOW (New Proposal)
1. Validate change with `spectr validate`
2. Create worktree
3. Copy change files to worktree
4. Commit with "feat: Proposal for <title>"
5. Push branch
6. Create PR with proposal details
7. Cleanup worktree
