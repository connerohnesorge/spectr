# Design: Auto-Sync Tasks (tasks.jsonc → tasks.md)

## Overview

This design document describes the implementation of automatic task status
synchronization from `tasks.jsonc` to `tasks.md` before every spectr command.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|---------|
| Hook timing | `AfterApply` | Runs after parsing, before Run() |
| Package location | `internal/sync/` | Clean separation |
| Missing tasks.md | Generate from jsonc | Full sync capability |
| Source of truth | `tasks.jsonc` | Machine-readable |
| Status mapping | `pending`/`in_progress` → `[ ]` | Convention match |

## Architecture

```text
┌─────────────────────────────────────────────────────────────────┐
│                         main.go                                  │
│  ctx, _ := app.Parse(os.Args[1:])  ←── AfterApply hook fires    │
│  ctx.Run()                                                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      cmd/root.go                                 │
│  type CLI struct {                                               │
│      NoSync  bool `help:"Skip task sync" name:"no-sync"`         │
│      Verbose bool `help:"Verbose output" name:"verbose"`         │
│      // ... existing commands                                    │
│  }                                                               │
│                                                                  │
│  func (c *CLI) AfterApply() error {                              │
│      if c.NoSync { return nil }                                  │
│      return sync.SyncAllActiveChanges(c.Verbose)                 │
│  }                                                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    internal/sync/sync.go                         │
│                                                                  │
│  SyncAllActiveChanges(verbose bool) error                        │
│      │                                                           │
│      ├── discovery.GetActiveChanges(projectPath)                 │
│      │                                                           │
│      └── for each change:                                        │
│              SyncTasksToMarkdown(changeDir, verbose)             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                internal/sync/markdown.go                         │
│                                                                  │
│  SyncTasksToMarkdown(changeDir string, verbose bool) error       │
│      │                                                           │
│      ├── Read tasks.jsonc (source of truth)                      │
│      │       └── parsers.ReadTasksJson()                         │
│      │                                                           │
│      ├── Check if tasks.md exists                                │
│      │       ├── YES: Update status markers in-place             │
│      │       └── NO:  Generate tasks.md from scratch             │
│      │                                                           │
│      └── Write updated tasks.md                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Implementation Details

### 1. Global Flags (cmd/root.go)

Add two global flags to the CLI struct:

```go
type CLI struct {
    // Global flags (apply to all commands)
    NoSync  bool `help:"Skip automatic task sync" name:"no-sync" short:"S"`
    Verbose bool `help:"Enable verbose output"    name:"verbose" short:"v"`

    // Existing commands...
    Init       InitCmd                   `cmd:""`
    List       ListCmd                   `cmd:""`
    // ...
}
```

### 2. AfterApply Hook (cmd/root.go)

Kong automatically calls `AfterApply()` after parsing but before `Run()`:

```go
func (c *CLI) AfterApply() error {
    if c.NoSync {
        return nil
    }

    projectRoot, err := os.Getwd()
    if err != nil {
        // Log error but don't block command
        fmt.Fprintf(os.Stderr, "sync: failed to get working directory: %v\n", err)
        return nil
    }

    // Check if spectr/ directory exists (not initialized = skip)
    spectrDir := filepath.Join(projectRoot, "spectr")
    if _, err := os.Stat(spectrDir); os.IsNotExist(err) {
        return nil
    }

    return sync.SyncAllActiveChanges(projectRoot, c.Verbose)
}
```

### 3. Sync Package (internal/sync/)

#### sync.go - Entry Point

```go
package sync

// SyncAllActiveChanges synchronizes task statuses from tasks.jsonc to tasks.md
// for all active changes in the project.
func SyncAllActiveChanges(projectRoot string, verbose bool) error {
    changeIDs, err := discovery.GetActiveChanges(projectRoot)
    if err != nil {
        // Directory doesn't exist or other error - skip silently
        return nil
    }

    var totalSynced int
    for _, id := range changeIDs {
        changeDir := filepath.Join(projectRoot, "spectr", "changes", id)

        synced, err := SyncTasksToMarkdown(changeDir)
        if err != nil {
            // Log error but continue with other changes
            fmt.Fprintf(os.Stderr, "sync: %s: %v\n", id, err)
            continue
        }

        if verbose && synced > 0 {
            fmt.Printf("Synced %d task statuses in %s\n", synced, id)
        }
        totalSynced += synced
    }

    return nil
}
```

#### markdown.go - Core Sync Logic

```go
package sync

// SyncTasksToMarkdown updates tasks.md checkbox statuses from tasks.jsonc.
// Returns the number of tasks whose status was updated.
func SyncTasksToMarkdown(changeDir string) (int, error) {
    tasksJsoncPath := filepath.Join(changeDir, "tasks.jsonc")
    tasksMdPath := filepath.Join(changeDir, "tasks.md")

    // Skip if no tasks.jsonc (not yet accepted)
    if _, err := os.Stat(tasksJsoncPath); os.IsNotExist(err) {
        return 0, nil
    }

    // Read source of truth
    tasksFile, err := parsers.ReadTasksJson(tasksJsoncPath)
    if err != nil {
        return 0, fmt.Errorf("read tasks.jsonc: %w", err)
    }

    // Build ID -> status map
    statusMap := buildStatusMap(tasksFile.Tasks)

    // Check if tasks.md exists
    if _, err := os.Stat(tasksMdPath); os.IsNotExist(err) {
        // Generate tasks.md from scratch
        return generateTasksMd(tasksMdPath, tasksFile.Tasks)
    }

    // Update existing tasks.md in-place
    return updateTasksMd(tasksMdPath, statusMap)
}

// buildStatusMap creates a map from task ID to checkbox character
func buildStatusMap(tasks []parsers.Task) map[string]rune {
    m := make(map[string]rune, len(tasks))
    for _, t := range tasks {
        if t.Status == parsers.TaskStatusCompleted {
            m[t.ID] = 'x'
        } else {
            m[t.ID] = ' ' // pending and in_progress both map to unchecked
        }
    }
    return m
}
```

#### update.go - In-Place Update Logic

```go
package sync

// updateTasksMd reads tasks.md, updates checkbox statuses, writes back.
// Preserves all formatting, comments, and structure.
func updateTasksMd(path string, statusMap map[string]rune) (int, error) {
    file, err := os.Open(path)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    var lines []string
    var updated int
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        newLine, changed := updateTaskLine(line, statusMap)
        lines = append(lines, newLine)
        if changed {
            updated++
        }
    }

    if err := scanner.Err(); err != nil {
        return 0, err
    }

    // Only write if changes were made
    if updated > 0 {
        content := strings.Join(lines, "\n")
        if !strings.HasSuffix(content, "\n") {
            content += "\n"
        }
        if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
            return 0, err
        }
    }

    return updated, nil
}

// updateTaskLine updates a single line's checkbox if it's a task line.
// Returns the (possibly modified) line and whether it was changed.
func updateTaskLine(line string, statusMap map[string]rune) (string, bool) {
    match, ok := markdown.MatchFlexibleTask(line)
    if !ok {
        return line, false
    }

    taskID := match.Number
    if taskID == "" {
        return line, false
    }

    desiredStatus, exists := statusMap[taskID]
    if !exists {
        return line, false
    }

    if match.Status == desiredStatus {
        return line, false
    }

    // Update checkbox character
    idx := strings.Index(line, "- [")
    if idx == -1 {
        return line, false
    }
    checkboxIdx := idx + 3

    newLine := line[:checkboxIdx] + string(desiredStatus) + line[checkboxIdx+1:]
    return newLine, true
}
```

#### generate.go - Generate tasks.md from tasks.jsonc

```go
package sync

// generateTasksMd creates a new tasks.md file from tasks.jsonc data.
func generateTasksMd(path string, tasks []parsers.Task) (int, error) {
    var sb strings.Builder
    sb.WriteString("# Tasks\n\n")

    var currentSection string
    sectionNum := 0

    for _, task := range tasks {
        if task.Section != currentSection {
            currentSection = task.Section
            sectionNum++
            if currentSection != "" {
                sb.WriteString(fmt.Sprintf("## %d. %s\n\n", sectionNum, currentSection))
            }
        }

        checkbox := ' '
        if task.Status == parsers.TaskStatusCompleted {
            checkbox = 'x'
        }
        sb.WriteString(fmt.Sprintf("- [%c] %s %s\n", checkbox, task.ID, task.Description))
    }

    if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
        return 0, err
    }

    return len(tasks), nil
}
```

## Task ID Matching Strategy

The existing `markdown.MatchFlexibleTask()` extracts task IDs in these formats:

- `1.1` (decimal)
- `1.` (dot-suffixed)
- `1` (number only)

Matching algorithm:

1. Parse tasks.jsonc to get `{ID: status}` map
2. For each line in tasks.md, use `MatchFlexibleTask()` to extract ID
3. Look up ID in status map
4. If found and status differs, update the checkbox character

## Edge Cases

| Case | Handling |
|------|----------|
| No spectr/ directory | Skip sync silently (not initialized) |
| No tasks.jsonc | Skip that change (not yet accepted) |
| No tasks.md but tasks.jsonc exists | Generate tasks.md from tasks.jsonc |
| Task in tasks.md not in tasks.jsonc | Leave unchanged |
| Task in tasks.jsonc not in tasks.md | Ignored (would need regeneration) |
| tasks.md has comments/links | Preserved (only checkbox updated) |
| Sync fails for one change | Log error, continue with others |
| Read-only filesystem | Log error, don't block command |

## Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `cmd/root.go` | Modify | Add flags and AfterApply() hook |
| `internal/sync/sync.go` | Create | Entry point, iterate active changes |
| `internal/sync/markdown.go` | Create | Core sync logic, status map building |
| `internal/sync/update.go` | Create | In-place tasks.md update |
| `internal/sync/generate.go` | Create | Generate tasks.md from tasks.jsonc |
| `internal/sync/sync_test.go` | Create | Unit tests for sync logic |

## Testing Strategy

**Unit Tests** (`internal/sync/sync_test.go`):

- Test `updateTaskLine()` with various task formats
- Test `buildStatusMap()` with all status values
- Test `generateTasksMd()` output format
- Test sync with mock filesystem (using temp dirs)

**Table-Driven Test Cases**:

```go
tests := []struct {
    name           string
    tasksJsonc     string
    tasksMd        string
    wantTasksMd    string
    wantSyncCount  int
}{
    {
        name: "update single task from pending to completed",
        tasksJsonc: `{"version":1,"tasks":[{"id":"1.1","status":"completed"}]}`,
        tasksMd:    "- [ ] 1.1 Task one\n",
        wantTasksMd: "- [x] 1.1 Task one\n",
        wantSyncCount: 1,
    },
}
```

## Performance Considerations

- **File I/O**: Only read/write files that need updating
- **Memory**: Stream tasks.md line-by-line, don't load entire file
- **Disk writes**: Skip write if no changes detected
- **Discovery**: Reuse existing `GetActiveChanges()` (already optimized)

Expected overhead: <10ms for typical projects with <10 active changes.
