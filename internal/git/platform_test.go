package git

import (
	"testing"
)

func TestDetectPlatformFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected Platform
	}{
		{
			name:     "GitHub HTTPS",
			url:      "https://github.com/user/repo.git",
			expected: PlatformGitHub,
		},
		{
			name:     "GitHub SSH",
			url:      "git@github.com:user/repo.git",
			expected: PlatformGitHub,
		},
		{
			name:     "GitLab HTTPS",
			url:      "https://gitlab.com/user/repo.git",
			expected: PlatformGitLab,
		},
		{
			name:     "GitLab SSH",
			url:      "git@gitlab.com:user/repo.git",
			expected: PlatformGitLab,
		},
		{
			name:     "Self-hosted GitLab",
			url:      "https://gitlab.example.com/user/repo.git",
			expected: PlatformGitLab,
		},
		{
			name:     "Gitea",
			url:      "https://gitea.example.com/user/repo.git",
			expected: PlatformGitea,
		},
		{
			name:     "Forgejo",
			url:      "https://forgejo.example.com/user/repo.git",
			expected: PlatformGitea,
		},
		{
			name:     "Unknown platform",
			url:      "https://unknown.com/user/repo.git",
			expected: PlatformUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectPlatformFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("detectPlatformFromURL(%q) = %v, want %v",
					tt.url, result, tt.expected)
			}
		})
	}
}

func TestGetCLITool(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		expected string
		wantErr  bool
	}{
		{
			name:     "GitHub",
			platform: PlatformGitHub,
			expected: "gh",
			wantErr:  false,
		},
		{
			name:     "GitLab",
			platform: PlatformGitLab,
			expected: "glab",
			wantErr:  false,
		},
		{
			name:     "Gitea",
			platform: PlatformGitea,
			expected: "tea",
			wantErr:  false,
		},
		{
			name:     "Unknown",
			platform: PlatformUnknown,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetCLITool(tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCLITool(%v) error = %v, wantErr %v",
					tt.platform, err, tt.wantErr)

				return
			}
			if result != tt.expected {
				t.Errorf("GetCLITool(%v) = %v, want %v",
					tt.platform, result, tt.expected)
			}
		})
	}
}

func TestGetInstallURL(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		expected string
	}{
		{
			name:     "GitHub",
			platform: PlatformGitHub,
			expected: "https://cli.github.com/",
		},
		{
			name:     "GitLab",
			platform: PlatformGitLab,
			expected: "https://gitlab.com/gitlab-org/cli",
		},
		{
			name:     "Gitea",
			platform: PlatformGitea,
			expected: "https://gitea.com/gitea/tea",
		},
		{
			name:     "Unknown",
			platform: PlatformUnknown,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getInstallURL(tt.platform)
			if result != tt.expected {
				t.Errorf("getInstallURL(%v) = %v, want %v",
					tt.platform, result, tt.expected)
			}
		})
	}
}
