package discovery

// skipDirsSet is a pre-computed set for O(1) directory skip lookups.
// Includes common large directories that should not be traversed.
var skipDirsSet = map[string]struct{}{
	gitDirName:       {},
	"node_modules":   {},
	"vendor":         {},
	"target":         {},
	"dist":           {},
	"build":          {},
	".cache":         {},
	".local":         {},
	".npm":           {},
	".pnpm":          {},
	".yarn":          {},
	".cargo":         {},
	".rustup":        {},
	"__pycache__":    {},
	".venv":          {},
	"venv":           {},
	".tox":           {},
	".nox":           {},
	".eggs":          {},
	"*.egg-info":     {},
	".pytest_cache":  {},
	".mypy_cache":    {},
	".ruff_cache":    {},
	"coverage":       {},
	".coverage":      {},
	".gradle":        {},
	".m2":            {},
	".ivy2":          {},
	"bin":            {},
	"obj":            {},
	"out":            {},
	".next":          {},
	".nuxt":          {},
	".svelte-kit":    {},
	".vercel":        {},
	".netlify":       {},
	"_build":         {},
	"site-packages":  {},
	".terraform":     {},
	".pulumi":        {},
	".serverless":    {},
	"testdata":       {},
	"fixtures":       {},
	".direnv":        {},
	".devenv":        {},
	"result":         {}, // Nix build output symlink
	".nix-defexpr":   {},
	".nix-profile":   {},
	"zig-cache":      {},
	"zig-out":        {},
	".zig-cache":     {},
	"bazel-bin":      {},
	"bazel-out":      {},
	"bazel-testlogs": {},
}

// shouldSkipDirectory returns true if the directory should be skipped during downward discovery.
func shouldSkipDirectory(dirName string) bool {
	// Fast path: check the pre-computed set
	if _, skip := skipDirsSet[dirName]; skip {
		return true
	}

	// Skip hidden directories (except .git which is handled separately)
	if len(dirName) > 1 && dirName[0] == '.' && dirName != gitDirName {
		return true
	}

	return false
}
