---
description: Commit and push all changes to the repository.
allowed-tools: Read, Glob, Grep, Write, Edit, Bash(spectr:*)
subtask: false
context: fork
---
# Guide to Committing and Pushing Changes to a Git Repository

## Basic Workflow: Commit & Push

### 1. Check Current Status
First, see what files have changed:
```bash
git status
```
This shows:
- Untracked files: New files Git isn't tracking yet
- Modified files: Existing files with changes
- Staged files: Changes ready to be committed

### 2. Stage Your Changes
Add specific files to staging:
```bash
git add filename.txt
git add folder/
```

Add all changes at once:
```bash
git add .
```

### 3. Create the Commit
Write a clear, descriptive commit message:
```bash
git commit -m "Add user authentication feature"
```

For more detailed messages, use your default editor:
```bash
git commit
```
This opens an editor where you can write a multi-line message (title + description).

### 4. Push to Remote Repository
Push your commits to the main branch:
```bash
git push origin main
```
*Note: Use `master` instead of `main` for older repositories.*

---

## Common Scenarios

### Scenario 1: First Push (Remote Has No Commits)
```bash
git push -u origin main
```
The `-u` flag sets up tracking, so future pushes can be just `git push`.

### Scenario 2: Committing Everything at Once
Stage and commit in one command (only for tracked files):
```bash
git commit -am "Fix navigation bug"
```

### Scenario 3: Pushing a New Branch
```bash
git checkout -b feature-branch
# ... make changes, commit ...
git push -u origin feature-branch
```

---

## Best Practices

### Commit Messages
- Clear and concise: 50 characters or less for the title
- Imperative mood: "Add feature" not "Added feature"
- Explain the why: Body should explain reasoning if complex

Example:
```
Add password reset functionality

Users were locked out without recovery options.
This adds email-based password reset flow.
```

### Atomic Commits
- One logical change per commit
- Don't mix unrelated changes
- Easier to review and revert

### Commit Frequency
- Commit early and often
- Push when features are complete and tested
- Don't push broken code to shared branches

---

## Troubleshooting

### "Updates were rejected"
```bash
git pull --rebase origin main
git push origin main
```

### "Authentication failed"
- Check your credentials
- For HTTPS: Use personal access token (not password)
- For SSH: Ensure your keys are set up correctly

### Accidental Commit
To undo the last commit (keeps changes staged):
```bash
git reset --soft HEAD~1
```

To completely undo and unstage:
```bash
git reset HEAD~1
```

---

## Complete Example
```bash
# Check what changed
git status

# Stage files
git add src/
git add README.md

# Commit
git commit -m "Implement search with filters"

# Pull latest changes first (good habit)
git pull origin main

# Push your commits
git push origin main
```
