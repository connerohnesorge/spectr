# Implementation Tasks

## 1. Update Interactive Model

- [x] 1.1 Add `itemType` field to `interactiveModel` struct to distinguish between specs and changes mode
- [x] 1.2 Update `RunInteractiveSpecs()` to set `itemType` to "spec" when creating the model
- [x] 1.3 Update `RunInteractiveChanges()` to set `itemType` to "change" when creating the model
- [x] 1.4 Update help text in `RunInteractiveSpecs()` to include "e: edit spec"

## 2. Implement Editor Launch Logic

- [x] 2.1 Create `handleEdit()` method on `interactiveModel` that:
  - Checks if `itemType` is "spec" (return early if not)
  - Gets the selected row ID
  - Checks if $EDITOR environment variable is set
  - Constructs the spec file path: `<projectPath>/spectr/specs/<spec-id>/spec.md`
  - Verifies the file exists
  - Launches the editor process and waits for it to complete
  - Returns updated model with any errors captured
- [x] 2.2 Add `projectPath` field to `interactiveModel` struct to support file path construction
- [x] 2.3 Update `RunInteractiveSpecs()` to pass `projectPath` when creating the model

## 3. Add Keyboard Event Handling

- [x] 3.1 Add case for "e" key in the `Update()` method's key handler
- [x] 3.2 Call `handleEdit()` when 'e' is pressed
- [x] 3.3 Return the updated model without quitting (TUI stays active)
- [x] 3.4 Ensure error messages are displayed if editor launch fails

## 4. Error Handling

- [x] 4.1 Handle $EDITOR not set - display clear error message
- [x] 4.2 Handle spec file not found - display path in error message
- [x] 4.3 Handle editor process launch failure - display underlying error
- [x] 4.4 Ensure TUI remains active after errors (don't quit)

## 5. Testing

- [x] 5.1 Write unit tests for `handleEdit()` method
- [x] 5.2 Test with various editors (vim, nano, emacs, VS Code)
- [x] 5.3 Test error cases: $EDITOR not set, file not found, editor binary missing
- [x] 5.4 Test that 'e' key is ignored in changes mode
- [x] 5.5 Manual integration test: navigate specs list, press 'e', verify editor opens correct file
- [x] 5.6 Verify help text is updated and displayed correctly

## 6. Documentation

- [x] 6.1 Update spec delta is already created in this change
- [x] 6.2 Verify all scenarios in spec delta are covered by implementation
