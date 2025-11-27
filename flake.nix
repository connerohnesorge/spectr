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
        generate-gif = {
          exec = rooted ''
            cd "$REPO_ROOT"
            DEMOS="init list validate archive workflow"

            if [[ "''${1:-}" == "-h" ]] || [[ "''${1:-}" == "--help" ]]; then
              echo "Usage: generate-gif [demo]"
              echo "  generate-gif        # Generate all GIFs"
              echo "  generate-gif init   # Generate single GIF"
              echo "Available demos: $DEMOS"
              exit 0
            fi

            mkdir -p "$REPO_ROOT/assets/gifs"

            if [[ -n "''${1:-}" ]]; then
              if [[ ! -f "$REPO_ROOT/assets/vhs/$1.tape" ]]; then
                echo "Error: Unknown demo '$1'. Available: $DEMOS" >&2
                exit 1
              fi
              echo "==> Generating $1.gif..."
              vhs "$REPO_ROOT/assets/vhs/$1.tape"
            else
              echo "==> Generating all demo GIFs..."
              for demo in $DEMOS; do
                echo "==> Generating $demo.gif..."
                vhs "$REPO_ROOT/assets/vhs/$demo.tape"
              done
              echo "==> All GIFs generated successfully!"
            fi

            rm -rf "$REPO_ROOT"/_demo
          '';
          description = "Generate VHS demo GIFs";
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

            # For docs site
            biome
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
        default = pkgs.buildGoModule rec {
          pname = "spectr";
          version = "0.0.1";
          src = self;
          vendorHash = "sha256-5mlBRTA+olCXHyKHb7zkLe0ySTv9Y2Qq4GbylzqSyjY=";
          ldflags = [
            "-s"
            "-w"
            "-X github.com/connerohnesorge/spectr/internal/version.Version=${version}"
          ];
          meta = with pkgs.lib; {
            description = "A CLI tool for spec-driven development workflow with change proposals, validation, and archiving";
            homepage = "https://github.com/connerohnesorge/spectr";
            license = licenses.asl20;
            maintainers = with maintainers; [connerohnesorge];
          };
        };
      };

      formatter = treefmt-nix.lib.mkWrapper pkgs treefmtModule;
    });
}
