# Change: Add Uncommitted Filter Hotkey

## Summary

Add an 'h' hotkey to the `spectr list -I` interactive TUI that filters the list to show only:
1. Changes with uncommitted git modifications (files in the change directory have uncommitted changes)
2. Changes with fully completed tasks.md files (all tasks marked complete)

This helps users quickly identify changes that are ready to be committed or archived.

## Motivation

When working with multiple active changes, it can be difficult to identify which changes:
- Have local modifications that need to be committed to git
- Have all tasks completed and are ready for archive

The 'h' hotkey ("harvestable" or "ready to handle") provides a quick filter to surface these actionable changes.

## Scope

- Add 'h' hotkey to the interactive list TUI
- Implement git status checking for change directories
- Add filtering logic for uncommitted changes with complete tasks
- Toggle behavior: press 'h' to enable filter, press again to disable
- Update help text to include the new hotkey

## Non-Goals

- Modifying the non-interactive list output
- Adding similar functionality to specs list mode
- Automatic git operations (commit, push, etc.)
