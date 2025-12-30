# Implementation Tasks

## 1. Implementation

- [ ] 1.1 Create `internal/init/templates/spectr/instruction-pointer.md.tmpl`
  with short pointer content
- [ ] 1.2 Add `RenderInstructionPointer()` method to `TemplateManager` in
  `templates.go`
- [ ] 1.3 Update `TemplateRenderer` interface in `providers/provider.go` to
  include `RenderInstructionPointer()`
- [ ] 1.4 Update `configureConfigFile()` in `providers/provider.go` to call
  `RenderInstructionPointer()` instead of `RenderAgents()`
- [ ] 1.5 Add unit test for `RenderInstructionPointer()` in `templates_test.go`
- [ ] 1.6 Run `spectr init` on this repo to verify the change works correctly
