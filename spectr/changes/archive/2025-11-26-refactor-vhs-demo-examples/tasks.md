# Implementation Tasks

## 1. Create Example Projects

- [x] 1.1 Create `examples/archive/` with a complete spectr project and a change
  ready to archive
- [x] 1.2 Create `examples/validate/` with spectr project containing
  intentionally broken and fixed specs
- [x] 1.3 Create `examples/workflow/` with a spectr project for end-to-end
  workflow demo
- [x] 1.4 Create `examples/list/` with a spectr project containing multiple
  specs and changes
- [x] 1.5 Create `examples/init/` as an empty directory (init demo creates fresh
  project)

## 2. Refactor VHS Tapes

- [x] 2.1 Refactor `archive.tape` to use `examples/archive/` instead of inline
  heredocs
- [x] 2.2 Refactor `validate.tape` to use `examples/validate/` instead of inline
  heredocs
- [x] 2.3 Refactor `workflow.tape` to use `examples/workflow/` instead of inline
  heredocs
- [x] 2.4 Refactor `list.tape` to use `examples/list/` instead of hardcoded
  paths
- [x] 2.5 Review `init.tape` - may only need minor cleanup (init creates fresh
  projects)

## 3. Create GIF Generation Script

- [x] 3.1 Create `scripts/generate-gifs.sh` that runs VHS on all tape files
- [x] 3.2 Add error handling and progress output
- [x] 3.3 Add option to generate single GIF or all GIFs
- [x] 3.4 Document usage in script header

## 4. Validation

- [x] 4.1 Run generate script and verify all GIFs are created
- [x] 4.2 Visually verify GIFs show clean, focused demos
- [x] 4.3 Update README if needed with new asset generation instructions
