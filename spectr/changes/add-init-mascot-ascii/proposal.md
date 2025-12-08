# Change: Add ASCII Ghost Mascot to Init Wizard Banner

## Why

The init wizard currently displays only the "SPECTR" text logo with a gradient. Adding the official Spectr mascot (ghost gopher) as ASCII art alongside the text creates a more memorable and branded experience during initialization. This reinforces Spectr's visual identity and makes the init wizard more visually engaging.

## What Changes

- Add ASCII art representation of the ghost gopher mascot from `assets/logo.png`
- Display mascot alongside the existing "SPECTR" gradient text in the init wizard intro screen
- Apply the same purple-to-pink gradient styling to the mascot ASCII art
- Ensure the combined banner fits well in standard terminal widths (80+ columns)

## Impact

- Affected specs: `cli-interface` (init wizard visual presentation)
- Affected code: `internal/initialize/wizard.go` (asciiArt constant and renderIntro function)
