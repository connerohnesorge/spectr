// Package pr provides PR workflow orchestration for spectr changes.
// It includes commit/PR message templating and platform CLI integration
// for creating pull requests that propose or archive spectr changes.
//
// The package supports two main workflows:
//
//   - Archive PR: Creates a PR that archives a completed change by
//     applying its delta specs to main specifications and moving
//     the change to the archive directory.
//
//   - New PR: Creates a PR to propose a new change for review,
//     copying the change directory to a new branch.
//
// Platform support includes GitHub (gh), GitLab (glab), Gitea (tea),
// and Bitbucket (manual URL output). Each platform uses its native
// CLI tool for PR creation when available.
package pr
