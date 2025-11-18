---
title: Quick Start
description: Get started with Spectr in minutes
---

## Initialize a Project

Start by initializing Spectr in your project:

```bash
# Initialize with interactive wizard
spectr init

# Or specify a path
spectr init /path/to/project

# Non-interactive mode with defaults
spectr init --non-interactive
```

This creates the following structure:

```
your-project/
└── spectr/
    ├── project.md        # Project conventions and context
    ├── specs/            # Current specifications (truth)
    │   └── [capability]/ # One directory per capability
    │       ├── spec.md   # Requirements and scenarios
    │       └── design.md # Technical patterns (optional)
    └── changes/          # Proposed changes
        └── archive/      # Completed changes
```

## Create Your First Change

Let's create a simple "Hello World" change:

```bash
# 1. List current state
spectr list              # See active changes
spectr list --specs      # See existing capabilities

# 2. Create a change directory
mkdir -p spectr/changes/add-hello-world/specs/greeting

# 3. Write a proposal
cat > spectr/changes/add-hello-world/proposal.md << 'EOF'
# Change: Add Hello World Greeting

## Why
We need a simple greeting capability to welcome users.

## What Changes
- Add new `greeting` capability with hello world functionality

## Impact
- Affected specs: greeting (new)
- Affected code: None (example)
EOF

# 4. Create delta spec
cat > spectr/changes/add-hello-world/specs/greeting/spec.md << 'EOF'
## ADDED Requirements

### Requirement: Hello World Greeting
The system SHALL provide a greeting function that returns "Hello, World!".

#### Scenario: Greet successfully
- **WHEN** the greeting function is called
- **THEN** it SHALL return "Hello, World!"
EOF

# 5. Create tasks checklist
cat > spectr/changes/add-hello-world/tasks.md << 'EOF'
## 1. Implementation
- [ ] 1.1 Create greeting.go file
- [ ] 1.2 Implement HelloWorld() function
- [ ] 1.3 Write tests for greeting
- [ ] 1.4 Update documentation
EOF

# 6. Validate the change
spectr validate add-hello-world --strict

# 7. After implementation, archive it
spectr archive add-hello-world
```

## File Structure

Understanding the directory structure is crucial:

```
spectr/
├── project.md              # Project-wide conventions
├── specs/                  # CURRENT TRUTH - what IS built
│   └── [capability]/
│       ├── spec.md         # Requirements with scenarios
│       └── design.md       # Technical patterns (optional)
├── changes/                # PROPOSALS - what SHOULD change
│   ├── [change-id]/
│   │   ├── proposal.md     # Why, what, impact
│   │   ├── tasks.md        # Implementation checklist
│   │   ├── design.md       # Technical decisions (optional)
│   │   └── specs/          # Delta changes
│   │       └── [capability]/
│   │           └── spec.md # ADDED/MODIFIED/REMOVED requirements
│   └── archive/            # Completed changes (history)
│       └── YYYY-MM-DD-[change-id]/
```

**Key Concepts:**
- **specs/**: The source of truth for what's currently built
- **changes/**: Proposed modifications, kept separate until approved
- **archive/**: Historical record of all changes with timestamps
- **Delta Specs**: Use `## ADDED`, `## MODIFIED`, `## REMOVED`, or `## RENAMED Requirements` headers
