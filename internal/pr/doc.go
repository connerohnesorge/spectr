// Package pr provides PR workflow orchestration for spectr change proposals.
//
// This package handles the complete workflow for creating pull requests from
// spectr changes, including:
//
//   - Commit message templating for both archive and new proposal modes
//   - PR body generation with structured markdown
//   - Change metadata extraction for use in templates
//
// The package supports two primary modes:
//
//   - Archive mode: Creates PRs for completed changes that are being archived,
//     including spec delta counts and updated capability information.
//
//   - New mode: Creates PRs for new change proposals ready for team review,
//     including proposal structure and file listings.
//
// Templates use Go's text/template package and are designed to produce
// conventional commit messages and well-structured PR bodies that integrate
// with GitHub, GitLab, Gitea, and other git hosting platforms.
package pr
