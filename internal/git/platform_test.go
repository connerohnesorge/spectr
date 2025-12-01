package git

import (
	"testing"
)

func TestDetectPlatform_GitHub(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected Platform
		cliTool  string
	}{
		{
			name:     "GitHub SSH URL",
			url:      "git@github.com:owner/repo.git",
			expected: PlatformGitHub,
			cliTool:  "gh",
		},
		{
			name:     "GitHub SSH URL without .git",
			url:      "git@github.com:owner/repo",
			expected: PlatformGitHub,
			cliTool:  "gh",
		},
		{
			name:     "GitHub HTTPS URL",
			url:      "https://github.com/owner/repo.git",
			expected: PlatformGitHub,
			cliTool:  "gh",
		},
		{
			name:     "GitHub HTTPS URL without .git",
			url:      "https://github.com/owner/repo",
			expected: PlatformGitHub,
			cliTool:  "gh",
		},
		{
			name:     "GitHub Enterprise SSH",
			url:      "git@github.mycompany.com:org/repo.git",
			expected: PlatformGitHub,
			cliTool:  "gh",
		},
		{
			name:     "GitHub Enterprise HTTPS",
			url:      "https://github.mycompany.com/org/repo.git",
			expected: PlatformGitHub,
			cliTool:  "gh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := DetectPlatform(tt.url)
			if info.Platform != tt.expected {
				t.Errorf("DetectPlatform(%q) platform = %v, want %v",
					tt.url, info.Platform, tt.expected)
			}
			if info.CLITool != tt.cliTool {
				t.Errorf("DetectPlatform(%q) cliTool = %v, want %v",
					tt.url, info.CLITool, tt.cliTool)
			}
		})
	}
}

func TestDetectPlatform_GitLab(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected Platform
		cliTool  string
	}{
		{
			name:     "GitLab SSH URL",
			url:      "git@gitlab.com:owner/repo.git",
			expected: PlatformGitLab,
			cliTool:  "glab",
		},
		{
			name:     "GitLab HTTPS URL",
			url:      "https://gitlab.com/owner/repo.git",
			expected: PlatformGitLab,
			cliTool:  "glab",
		},
		{
			name:     "GitLab self-hosted SSH",
			url:      "git@gitlab.mycompany.com:org/repo.git",
			expected: PlatformGitLab,
			cliTool:  "glab",
		},
		{
			name:     "GitLab self-hosted HTTPS",
			url:      "https://gitlab.internal.example.com/org/repo.git",
			expected: PlatformGitLab,
			cliTool:  "glab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := DetectPlatform(tt.url)
			if info.Platform != tt.expected {
				t.Errorf("DetectPlatform(%q) platform = %v, want %v",
					tt.url, info.Platform, tt.expected)
			}
			if info.CLITool != tt.cliTool {
				t.Errorf("DetectPlatform(%q) cliTool = %v, want %v",
					tt.url, info.CLITool, tt.cliTool)
			}
		})
	}
}

func TestDetectPlatform_Gitea(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected Platform
		cliTool  string
	}{
		{
			name:     "Gitea SSH URL",
			url:      "git@gitea.example.com:owner/repo.git",
			expected: PlatformGitea,
			cliTool:  "tea",
		},
		{
			name:     "Gitea HTTPS URL",
			url:      "https://gitea.example.com/owner/repo.git",
			expected: PlatformGitea,
			cliTool:  "tea",
		},
		{
			name:     "Forgejo SSH URL",
			url:      "git@forgejo.example.com:owner/repo.git",
			expected: PlatformGitea,
			cliTool:  "tea",
		},
		{
			name:     "Forgejo HTTPS URL",
			url:      "https://forgejo.example.com/owner/repo.git",
			expected: PlatformGitea,
			cliTool:  "tea",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := DetectPlatform(tt.url)
			if info.Platform != tt.expected {
				t.Errorf("DetectPlatform(%q) platform = %v, want %v",
					tt.url, info.Platform, tt.expected)
			}
			if info.CLITool != tt.cliTool {
				t.Errorf("DetectPlatform(%q) cliTool = %v, want %v",
					tt.url, info.CLITool, tt.cliTool)
			}
		})
	}
}

func TestDetectPlatform_Bitbucket(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected Platform
		cliTool  string
	}{
		{
			name:     "Bitbucket SSH URL",
			url:      "git@bitbucket.org:owner/repo.git",
			expected: PlatformBitbucket,
			cliTool:  "bb",
		},
		{
			name:     "Bitbucket HTTPS URL",
			url:      "https://bitbucket.org/owner/repo.git",
			expected: PlatformBitbucket,
			cliTool:  "bb",
		},
		{
			name:     "Bitbucket Server SSH",
			url:      "git@bitbucket.mycompany.com:org/repo.git",
			expected: PlatformBitbucket,
			cliTool:  "bb",
		},
		{
			name:     "Bitbucket Server HTTPS",
			url:      "https://bitbucket.internal.example.com/org/repo.git",
			expected: PlatformBitbucket,
			cliTool:  "bb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := DetectPlatform(tt.url)
			if info.Platform != tt.expected {
				t.Errorf("DetectPlatform(%q) platform = %v, want %v",
					tt.url, info.Platform, tt.expected)
			}
			if info.CLITool != tt.cliTool {
				t.Errorf("DetectPlatform(%q) cliTool = %v, want %v",
					tt.url, info.CLITool, tt.cliTool)
			}
		})
	}
}

func TestDetectPlatform_Unknown(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Custom host SSH",
			url:  "git@custom.example.com:owner/repo.git",
		},
		{
			name: "Custom host HTTPS",
			url:  "https://custom.example.com/owner/repo.git",
		},
		{
			name: "Empty URL",
			url:  "",
		},
		{
			name: "Whitespace only",
			url:  "   ",
		},
		{
			name: "Invalid format",
			url:  "not-a-valid-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := DetectPlatform(tt.url)
			if info.Platform != PlatformUnknown {
				t.Errorf("DetectPlatform(%q) platform = %v, want %v",
					tt.url, info.Platform, PlatformUnknown)
			}
		})
	}
}

func TestPlatform_String(t *testing.T) {
	tests := []struct {
		platform Platform
		expected string
	}{
		{PlatformGitHub, "GitHub"},
		{PlatformGitLab, "GitLab"},
		{PlatformGitea, "Gitea"},
		{PlatformBitbucket, "Bitbucket"},
		{PlatformUnknown, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.platform.String(); got != tt.expected {
				t.Errorf("Platform.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectPlatform_RepoURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectedURL string
	}{
		{
			name:        "preserves original URL",
			url:         "git@github.com:owner/repo.git",
			expectedURL: "git@github.com:owner/repo.git",
		},
		{
			name:        "trims whitespace",
			url:         "  git@github.com:owner/repo.git  ",
			expectedURL: "git@github.com:owner/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := DetectPlatform(tt.url)
			if info.RepoURL != tt.expectedURL {
				t.Errorf("DetectPlatform(%q) RepoURL = %v, want %v",
					tt.url, info.RepoURL, tt.expectedURL)
			}
		})
	}
}
