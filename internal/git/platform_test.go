package git

import (
	"strings"
	"testing"
)

func TestDetectPlatform(t *testing.T) {
	tests := []struct {
		name         string
		remoteURL    string
		wantPlatform Platform
		wantCLITool  string
		wantOwner    string
		wantRepo     string
		wantRepoURL  string
	}{
		// GitHub URLs
		{
			name:         "GitHub HTTPS",
			remoteURL:    "https://github.com/owner/repo.git",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.com/owner/repo",
		},
		{
			name:         "GitHub HTTPS without .git",
			remoteURL:    "https://github.com/owner/repo",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.com/owner/repo",
		},
		{
			name:         "GitHub SSH",
			remoteURL:    "git@github.com:owner/repo.git",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.com/owner/repo",
		},
		{
			name:         "GitHub SSH with protocol",
			remoteURL:    "ssh://git@github.com/owner/repo.git",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.com/owner/repo",
		},
		{
			name:         "GitHub Enterprise",
			remoteURL:    "https://github.mycompany.com/owner/repo.git",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.mycompany.com/owner/repo",
		},

		// GitLab URLs
		{
			name:         "GitLab HTTPS",
			remoteURL:    "https://gitlab.com/owner/repo.git",
			wantPlatform: PlatformGitLab,
			wantCLITool:  "glab",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://gitlab.com/owner/repo",
		},
		{
			name:         "GitLab SSH",
			remoteURL:    "git@gitlab.com:owner/repo.git",
			wantPlatform: PlatformGitLab,
			wantCLITool:  "glab",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://gitlab.com/owner/repo",
		},
		{
			name:         "GitLab self-hosted",
			remoteURL:    "https://gitlab.mycompany.com/owner/repo.git",
			wantPlatform: PlatformGitLab,
			wantCLITool:  "glab",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://gitlab.mycompany.com/owner/repo",
		},
		{
			name:         "GitLab with groups",
			remoteURL:    "https://gitlab.com/group/subgroup/repo.git",
			wantPlatform: PlatformGitLab,
			wantCLITool:  "glab",
			wantOwner:    "group/subgroup",
			wantRepo:     "repo",
			wantRepoURL:  "https://gitlab.com/group/subgroup/repo",
		},
		{
			name:         "GitLab with deeply nested groups",
			remoteURL:    "https://gitlab.com/org/team/project/repo.git",
			wantPlatform: PlatformGitLab,
			wantCLITool:  "glab",
			wantOwner:    "org/team/project",
			wantRepo:     "repo",
			wantRepoURL:  "https://gitlab.com/org/team/project/repo",
		},

		// Gitea/Forgejo URLs
		{
			name:         "Gitea HTTPS",
			remoteURL:    "https://gitea.io/owner/repo.git",
			wantPlatform: PlatformGitea,
			wantCLITool:  "tea",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://gitea.io/owner/repo",
		},
		{
			name:         "Gitea self-hosted",
			remoteURL:    "https://git.gitea.example.com/owner/repo.git",
			wantPlatform: PlatformGitea,
			wantCLITool:  "tea",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://git.gitea.example.com/owner/repo",
		},
		{
			name:         "Forgejo HTTPS",
			remoteURL:    "https://forgejo.example.com/owner/repo.git",
			wantPlatform: PlatformGitea,
			wantCLITool:  "tea",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://forgejo.example.com/owner/repo",
		},

		// Bitbucket URLs
		{
			name:         "Bitbucket HTTPS",
			remoteURL:    "https://bitbucket.org/owner/repo.git",
			wantPlatform: PlatformBitbucket,
			wantCLITool:  "",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://bitbucket.org/owner/repo",
		},
		{
			name:         "Bitbucket SSH",
			remoteURL:    "git@bitbucket.org:owner/repo.git",
			wantPlatform: PlatformBitbucket,
			wantCLITool:  "",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://bitbucket.org/owner/repo",
		},
		{
			name:         "Bitbucket self-hosted",
			remoteURL:    "https://bitbucket.mycompany.com/owner/repo.git",
			wantPlatform: PlatformBitbucket,
			wantCLITool:  "",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://bitbucket.mycompany.com/owner/repo",
		},

		// Unknown/custom hosts
		{
			name:         "Unknown host",
			remoteURL:    "https://unknown.example.com/owner/repo.git",
			wantPlatform: PlatformUnknown,
			wantCLITool:  "",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://unknown.example.com/owner/repo",
		},
		{
			name:         "Custom git server",
			remoteURL:    "https://git.mycompany.com/team/project.git",
			wantPlatform: PlatformUnknown,
			wantCLITool:  "",
			wantOwner:    "team",
			wantRepo:     "project",
			wantRepoURL:  "https://git.mycompany.com/team/project",
		},

		// Edge cases
		{
			name:         "Git protocol",
			remoteURL:    "git://github.com/owner/repo.git",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.com/owner/repo",
		},
		{
			name:         "SSH with custom user",
			remoteURL:    "custom-user@github.com:owner/repo.git",
			wantPlatform: PlatformGitHub,
			wantCLITool:  "gh",
			wantOwner:    "owner",
			wantRepo:     "repo",
			wantRepoURL:  "https://github.com/owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectPlatform(
				tt.remoteURL,
			)
			if err != nil {
				t.Fatalf(
					"DetectPlatform(%q) error = %v",
					tt.remoteURL,
					err,
				)
			}

			if got.Platform != tt.wantPlatform {
				t.Errorf(
					"DetectPlatform(%q).Platform = %v, want %v",
					tt.remoteURL,
					got.Platform,
					tt.wantPlatform,
				)
			}
			if got.CLITool != tt.wantCLITool {
				t.Errorf(
					"DetectPlatform(%q).CLITool = %v, want %v",
					tt.remoteURL,
					got.CLITool,
					tt.wantCLITool,
				)
			}
			if got.Owner != tt.wantOwner {
				t.Errorf(
					"DetectPlatform(%q).Owner = %v, want %v",
					tt.remoteURL,
					got.Owner,
					tt.wantOwner,
				)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf(
					"DetectPlatform(%q).Repo = %v, want %v",
					tt.remoteURL,
					got.Repo,
					tt.wantRepo,
				)
			}
			if got.RepoURL != tt.wantRepoURL {
				t.Errorf(
					"DetectPlatform(%q).RepoURL = %v, want %v",
					tt.remoteURL,
					got.RepoURL,
					tt.wantRepoURL,
				)
			}
		})
	}
}

func TestDetectPlatform_InvalidURLs(
	t *testing.T,
) {
	tests := []struct {
		name      string
		remoteURL string
		wantErr   string
	}{
		{
			name:      "Empty string",
			remoteURL: "",
			wantErr:   "empty remote URL",
		},
		{
			name:      "Malformed URL - no path",
			remoteURL: "https://github.com",
			wantErr:   "unrecognized URL format",
		},
		{
			name:      "Malformed URL - only host",
			remoteURL: "github.com",
			wantErr:   "unrecognized URL format",
		},
		{
			name:      "Malformed URL - no protocol no path",
			remoteURL: "github.com/owner",
			wantErr:   "unrecognized URL format",
		},
		{
			name:      "Malformed URL - just owner",
			remoteURL: "https://github.com/owner",
			wantErr:   "failed to extract owner/repo",
		},
		{
			name:      "Malformed URL - random string",
			remoteURL: "not-a-url",
			wantErr:   "unrecognized URL format",
		},
		{
			name:      "Malformed URL - incomplete SSH",
			remoteURL: "git@github.com",
			wantErr:   "unrecognized URL format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DetectPlatform(tt.remoteURL)
			if err == nil {
				t.Fatalf(
					"DetectPlatform(%q) expected error, got nil",
					tt.remoteURL,
				)
			}

			if tt.wantErr != "" &&
				!strings.Contains(
					err.Error(),
					tt.wantErr,
				) {
				t.Errorf(
					"DetectPlatform(%q) error = %v, want error containing %q",
					tt.remoteURL,
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestParseRemoteURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantHost string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "HTTPS URL",
			url:      "https://github.com/owner/repo.git",
			wantHost: "github.com",
			wantPath: "owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "HTTP URL",
			url:      "http://github.com/owner/repo.git",
			wantHost: "github.com",
			wantPath: "owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "SSH URL",
			url:      "git@github.com:owner/repo.git",
			wantHost: "github.com",
			wantPath: "owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "SSH URL with protocol",
			url:      "ssh://git@github.com/owner/repo.git",
			wantHost: "github.com",
			wantPath: "owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "Git protocol URL",
			url:      "git://github.com/owner/repo.git",
			wantHost: "github.com",
			wantPath: "owner/repo.git",
			wantErr:  false,
		},
		{
			name:    "Invalid URL",
			url:     "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, path, err := parseRemoteURL(
				tt.url,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf(
						"parseRemoteURL(%q) expected error, got nil",
						tt.url,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"parseRemoteURL(%q) unexpected error: %v",
					tt.url,
					err,
				)
			}

			if host != tt.wantHost {
				t.Errorf(
					"parseRemoteURL(%q) host = %v, want %v",
					tt.url,
					host,
					tt.wantHost,
				)
			}
			if path != tt.wantPath {
				t.Errorf(
					"parseRemoteURL(%q) path = %v, want %v",
					tt.url,
					path,
					tt.wantPath,
				)
			}
		})
	}
}

func TestExtractOwnerRepo(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantOwner string
		wantRepo  string
	}{
		{
			name:      "Simple owner/repo with .git",
			path:      "owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "Simple owner/repo without .git",
			path:      "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "Nested groups",
			path:      "group/subgroup/repo.git",
			wantOwner: "group/subgroup",
			wantRepo:  "repo",
		},
		{
			name:      "Deeply nested groups",
			path:      "org/team/project/repo.git",
			wantOwner: "org/team/project",
			wantRepo:  "repo",
		},
		{
			name:      "Only repo name",
			path:      "repo.git",
			wantOwner: "",
			wantRepo:  "",
		},
		{
			name:      "Empty path",
			path:      "",
			wantOwner: "",
			wantRepo:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo := extractOwnerRepo(
				tt.path,
			)

			if owner != tt.wantOwner {
				t.Errorf(
					"extractOwnerRepo(%q) owner = %v, want %v",
					tt.path,
					owner,
					tt.wantOwner,
				)
			}
			if repo != tt.wantRepo {
				t.Errorf(
					"extractOwnerRepo(%q) repo = %v, want %v",
					tt.path,
					repo,
					tt.wantRepo,
				)
			}
		})
	}
}

func TestDetectPlatformFromHost(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		wantPlatform Platform
	}{
		// GitHub
		{
			name:         "GitHub.com",
			host:         "github.com",
			wantPlatform: PlatformGitHub,
		},
		{
			name:         "GitHub Enterprise",
			host:         "github.mycompany.com",
			wantPlatform: PlatformGitHub,
		},
		{
			name:         "GitHub uppercase",
			host:         "GITHUB.COM",
			wantPlatform: PlatformGitHub,
		},

		// GitLab
		{
			name:         "GitLab.com",
			host:         "gitlab.com",
			wantPlatform: PlatformGitLab,
		},
		{
			name:         "GitLab self-hosted",
			host:         "gitlab.mycompany.com",
			wantPlatform: PlatformGitLab,
		},
		{
			name:         "GitLab uppercase",
			host:         "GITLAB.COM",
			wantPlatform: PlatformGitLab,
		},

		// Gitea/Forgejo
		{
			name:         "Gitea.io",
			host:         "gitea.io",
			wantPlatform: PlatformGitea,
		},
		{
			name:         "Gitea self-hosted",
			host:         "gitea.mycompany.com",
			wantPlatform: PlatformGitea,
		},
		{
			name:         "Forgejo",
			host:         "forgejo.example.com",
			wantPlatform: PlatformGitea,
		},

		// Bitbucket
		{
			name:         "Bitbucket.org",
			host:         "bitbucket.org",
			wantPlatform: PlatformBitbucket,
		},
		{
			name:         "Bitbucket self-hosted",
			host:         "bitbucket.mycompany.com",
			wantPlatform: PlatformBitbucket,
		},

		// Unknown
		{
			name:         "Unknown host",
			host:         "unknown.example.com",
			wantPlatform: PlatformUnknown,
		},
		{
			name:         "Custom git server",
			host:         "git.mycompany.com",
			wantPlatform: PlatformUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectPlatformFromHost(tt.host)
			if got != tt.wantPlatform {
				t.Errorf(
					"detectPlatformFromHost(%q) = %v, want %v",
					tt.host,
					got,
					tt.wantPlatform,
				)
			}
		})
	}
}

func TestGetCLITool(t *testing.T) {
	tests := []struct {
		platform Platform
		want     string
	}{
		{PlatformGitHub, "gh"},
		{PlatformGitLab, "glab"},
		{PlatformGitea, "tea"},
		{PlatformBitbucket, ""},
		{PlatformUnknown, ""},
	}

	for _, tt := range tests {
		t.Run(
			string(tt.platform),
			func(t *testing.T) {
				got := getCLITool(tt.platform)
				if got != tt.want {
					t.Errorf(
						"getCLITool(%v) = %v, want %v",
						tt.platform,
						got,
						tt.want,
					)
				}
			},
		)
	}
}

func TestBuildRepoURL(t *testing.T) {
	tests := []struct {
		host  string
		owner string
		repo  string
		want  string
	}{
		{
			"github.com",
			"owner",
			"repo",
			"https://github.com/owner/repo",
		},
		{
			"gitlab.com",
			"group/subgroup",
			"repo",
			"https://gitlab.com/group/subgroup/repo",
		},
		{
			"bitbucket.org",
			"workspace",
			"project",
			"https://bitbucket.org/workspace/project",
		},
	}

	for _, tt := range tests {
		name := tt.host + "/" + tt.owner + "/" + tt.repo
		t.Run(name, func(t *testing.T) {
			got := buildRepoURL(
				tt.host,
				tt.owner,
				tt.repo,
			)
			if got != tt.want {
				t.Errorf(
					"buildRepoURL(%q, %q, %q) = %v, want %v",
					tt.host,
					tt.owner,
					tt.repo,
					got,
					tt.want,
				)
			}
		})
	}
}
