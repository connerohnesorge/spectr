## REMOVED Requirements

### Requirement: Archive PR Automation Flag

**Reason**: Feature complexity outweighs benefits; standard git workflows are more transparent and controllable.

**Migration**: Users can run `git add . && git commit -m "Archive: <change-id>" && gh pr create` after archiving.

### Requirement: Archive PR Branch Naming

**Reason**: Branch naming logic only exists to support the removed `--pr` flag.

**Migration**: Users create branches manually with their preferred naming convention.

### Requirement: Archive PR Commit Strategy

**Reason**: Commit strategy logic only exists to support the removed `--pr` flag.

**Migration**: Users control commit granularity and message format directly.

### Requirement: Archive PR Platform Detection

**Reason**: Platform detection for GitHub/GitLab/Gitea only exists to support the removed `--pr` flag.

**Migration**: Users invoke their platform's CLI directly (gh, glab, tea).

### Requirement: Archive PR Title and Body

**Reason**: PR title/body formatting only exists to support the removed `--pr` flag.

**Migration**: Users write PR descriptions suited to their workflow.

### Requirement: Archive PR Error Handling

**Reason**: Error handling for git operations only exists to support the removed `--pr` flag.

**Migration**: Standard git error messages are visible when running commands directly.

### Requirement: Archive PR Success Reporting

**Reason**: PR URL reporting only exists to support the removed `--pr` flag.

**Migration**: PR CLI tools display URLs directly when creating PRs.
