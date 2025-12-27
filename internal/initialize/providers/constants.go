package providers

// This file is intentionally minimal.
// All old helper functions and priority constants have been removed
// as part of the provider architecture redesign.
//
// - Priority values are now defined inline in RegisterAllProviders()
// - Helper functions like StandardCommandPaths() have been replaced
//   by initializers
// - Frontmatter templates are now managed by the TemplateManager
