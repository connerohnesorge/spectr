# Archive Package

Merges change deltas into specs and moves changes to archive. Atomic workflow.

## OVERVIEW
`spectr archive` validates change, merges delta specs into `spectr/specs/`, moves `spectr/changes/<id>/` → `spectr/changes/archive/YYYY-MM-DD-<id>/`. Ensures history preservation and spec consistency.

## STRUCTURE
```
internal/archive/
├── archiver.go          # Main archive orchestration
├── merger.go            # Spec merging logic
├── spec_merger.go       # Requirement-level merge algorithm
├── cmd.go               # CLI command handler
├── interactive_bridge.go # TUI prompts
├── constants.go         # Archive paths and filenames
└── *_test.go            # Integration tests
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Archive workflow | archiver.go | Validation → merge → move |
| Spec merging | merger.go + spec_merger.go | Delta → spec algorithm |
| Merge algorithm | spec_merger.go | Requirement-level merge logic |
| Interactive prompts | interactive_bridge.go | User confirmation |

## CONVENTIONS
- **Atomic operation**: All steps succeed or none do (validate+merge+move)
- **Date prefix**: Archive dirs named `YYYY-MM-DD-<change-id>`
- **Merge order**: ADDED → append, MODIFIED → replace, REMOVED → comment
- **Validation first**: Validate change before merging

## MERGE ALGORITHM (per requirement)
| Operation | Action | Detail |
|-----------|--------|--------|
| ADDED | Append to spec | New requirement added to end |
| MODIFIED | Replace entire requirement | Full content replaces existing (header + all scenarios) |
| REMOVED | Replace with comment | Keep requirement name, add migration info |
| RENAMED | Update header | Change name only |

## UNIQUE PATTERNS
- **SpecMerger**: Reads both source spec and delta, applies operations in sequence
- **Requirement identity**: Matches by name (whitespace-insensitive)
- **Archive move**: Uses os.Rename() for atomic directory move

## ANTI-PATTERNS
- **NEVER merge partial MODIFIED**: MODIFIED must include complete requirement
- **DON'T skip validation**: Always validate before merging
- **NO destructive moves**: Only move after successful merge

## KEY FUNCTIONS
- `Archiver.Archive(changeID string) error` - Main archive workflow
- `NewSpecMerger(specPath, deltaPath) (*SpecMerger, error)` - Create merger
- `SpecMerger.Merge() ([]byte, error)` - Execute merge and return new spec
- `SpecMerger.MoveToArchive() error` - Move change to archive/

## FLOW
1. Validate change structure
2. Parse delta specs from `changes/<id>/specs/`
3. For each affected spec:
   - Read existing `specs/<capability>/spec.md`
   - Create SpecMerger
   - Apply ADDED/MODIFIED/REMOVED/RENAMED operations
   - Write merged spec back to `specs/<capability>/spec.md`
4. Move `changes/<id>/` → `changes/archive/YYYY-MM-DD-<id>/`
5. Validate updated specs

## ERROR HANDLING
- Merge conflict: If multiple changes modify same requirement, manual resolution required
- Missing spec: If ADDED references non-existent spec, create new spec file
- Validation failure: Abort merge, report errors
