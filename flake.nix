/**
# Go Development Shell Template

## Description
Complete Go development environment with modern tooling for building, testing,
and maintaining Go applications. Includes the Go toolchain, linting, formatting,
live reloading, and testing utilities for productive Go development.

## Platform Support
- ✅ x86_64-linux
- ✅ aarch64-linux (ARM64 Linux)
- ✅ x86_64-darwin (Intel macOS)
- ✅ aarch64-darwin (Apple Silicon macOS)

## What This Provides
- **Go Toolchain**: Go 1.25 compiler and runtime
- **Development Tools**: air (live reload), delve (debugger), gopls (language server)
- **Code Quality**: golangci-lint, revive, gofmt, goimports
- **Testing**: gotestfmt for enhanced test output
- **Documentation**: gomarkdoc for generating markdown from Go code
- **Formatting**: gofumpt for stricter Go formatting

## Usage
```bash
# Create new project from template
nix flake init -t github:conneroisu/dotfiles#go-shell

# Enter development shell
nix develop

# Start live reload development
air

# Run tests with formatting
go test ./... | gotestfmt

# Format code
nix fmt
```

## Development Workflow
- Use air for automatic recompilation during development
- golangci-lint provides comprehensive linting
- gopls enables rich IDE integration
- All tools configured for optimal Go development experience
*/
{
  description = "A development shell for go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix.url = "github:numtide/treefmt-nix";
    treefmt-nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    treefmt-nix,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          (final: prev: {
            # Add your overlays here
            # Example:
            # my-overlay = final: prev: {
            #   my-package = prev.callPackage ./my-package { };
            # };
            final.buildGoModule = prev.buildGo125Module;
            buildGoModule = prev.buildGo125Module;
          })
        ];
      };

      rooted = exec:
        builtins.concatStringsSep "\n"
        [
          ''REPO_ROOT="$(git rev-parse --show-toplevel)"''
          exec
        ];

      scripts = {
        dx = {
          exec = rooted ''$EDITOR "$REPO_ROOT"/flake.nix'';
          description = "Edit flake.nix";
        };
        lint = {
          exec = rooted ''
            cd "$REPO_ROOT"
            golangci-lint run
            cd -
          '';
          description = "Run golangci-lint";
        };
        tests = {
          exec = rooted ''
            gotestsum --format testname -- -race "$REPO_ROOT"/... -timeout=2m
          '';
          description = "Run tests";
        };
      };

      scriptPackages =
        pkgs.lib.mapAttrs
        (
          name: script:
            pkgs.writeShellApplication {
              inherit name;
              text = script.exec;
              runtimeInputs = script.deps or [];
            }
        )
        scripts;

      treefmtModule = {
        projectRootFile = "flake.nix";
        programs = {
          alejandra.enable = true; # Nix formatter
          gofmt.enable = true; # Go formatter
          golines.enable = true; # Go formatter (Shorter lines)
          goimports.enable = true; # Go formatter (Organize/Clean imports)
        };
      };
    in {
      devShells.default = pkgs.mkShell {
        name = "dev";

        # Available packages on https://search.nixos.org/packages
        packages = with pkgs;
          [
            alejandra # Nix
            nixd
            statix
            deadnix

            go_1_25 # Go Tools
            air
            golangci-lint
            golangci-lint-langserver
            gopls
            revive
            golines
            gomarkdoc
            gotests
            gotestsum
            gotools
            reftools
            goreleaser
            vhs
          ]
          ++ builtins.attrValues scriptPackages;
      };

      devShells.ci = pkgs.mkShell {
        name = "ci";

        # Minimal CI/CD tooling - no interactive development tools
        packages = with pkgs; [
          # Nix formatting
          alejandra

          # Go toolchain
          go_1_25

          # Testing
          gotestsum
          gotools

          # Linting
          golangci-lint

          # Build
          goreleaser
        ];
      };

      packages = {
        default = pkgs.buildGoModule {
          pname = "spectr";
          version = "0.0.1";
          src = self;
          vendorHash = "sha256-6bE9HNbebJ4ivHF7YynZwL6mu+T3wEfESjQdyR8q59M=";
          meta = with pkgs.lib; {
            description = "A CLI tool for spec-driven development workflow with change proposals, validation, and archiving";
            homepage = "https://github.com/conneroisu/spectr";
            license = licenses.asl20;
            maintainers = with maintainers; [connerohnesorge];
          };
        };
      };

      formatter = treefmt-nix.lib.mkWrapper pkgs treefmtModule;
    });
}
