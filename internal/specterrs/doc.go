// Package specterrs provides centralized error types for the spectr CLI.
//
// All custom error types in this package:
//   - Use pointer receivers for the Error() method
//   - Include structured fields for contextual information
//   - Implement Unwrap() when wrapping underlying errors
//
// Error types are organized by domain:
//   - git.go: Git repository and branch errors
//   - archive.go: Archive workflow errors
//   - validation.go: Spec/change validation errors
//   - initialize.go: Project initialization errors
//   - list.go: List command errors
//   - environment.go: Environment configuration errors
//   - pr.go: Pull request workflow errors
package specterrs
