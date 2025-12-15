---
title: Installation
description: How to install Spectr
---

## Using Nix Flakes

The recommended way to install Spectr is via Nix flakes:

```bash
# Run directly without installing
nix run github:connerohnesorge/spectr

# Install to your profile
nix profile install github:connerohnesorge/spectr

# Add to your flake.nix inputs
{
  inputs.spectr.url = "github:connerohnesorge/spectr";
}
```

### NixOS Configuration Example

Here is a complete example showing how to add spectr to your NixOS system configuration:

```nix
{
  description = "My NixOS configuration with spectr";

  inputs = {
    # NixOS official package source
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    # Add spectr as an input
    spectr = {
      url = "github:connerohnesorge/spectr";
      # Optional: use the same nixpkgs as your system
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, spectr, ... }: {
    # Define your NixOS system configuration
    nixosConfigurations.my-hostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";

      modules = [
        # Your hardware configuration
        ./hardware-configuration.nix

        # Main system configuration
        ({ pkgs, ... }: {
          # Add spectr to system-wide packages
          environment.systemPackages = [
            # Access spectr package for your system architecture
            spectr.packages.${pkgs.system}.default

            # Other packages...
            pkgs.git
            pkgs.vim
          ];

          # ... rest of your NixOS configuration
        })
      ];
    };
  };
}
```

After adding spectr to your configuration, rebuild your system:

```bash
sudo nixos-rebuild switch --flake .#my-hostname
```

### Development Shell Example

To include spectr in a project's development shell, add it to your `devShells.default` configuration. This is useful when you want spectr available in your project's development environment:

```nix
{
  description = "My project with spectr in dev shell";

  inputs = {
    # Nixpkgs for standard packages
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    # Add spectr as an input
    spectr = {
      url = "github:connerohnesorge/spectr";
      # Optional: share nixpkgs to reduce closure size
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # Optional: flake-utils for multi-system support
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, spectr, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        # Define the default development shell
        devShells.default = pkgs.mkShell {
          # Include spectr and other development tools
          buildInputs = [
            # Add spectr to the development environment
            spectr.packages.${system}.default

            # Add other tools your project needs
            pkgs.git
            pkgs.go
            pkgs.gopls
          ];

          # Optional: Set environment variables
          shellHook = ''
            echo "Development shell loaded with spectr"
            echo "Run 'spectr --help' to get started"
          '';
        };
      }
    );
}
```

Enter the development shell with:

```bash
# Enter the dev shell (runs shellHook automatically)
nix develop

# Or run a command directly in the shell
nix develop --command spectr --help
```

For projects without flake-utils, you can define shells for specific systems:

```nix
{
  description = "Simple project with spectr";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    spectr.url = "github:connerohnesorge/spectr";
  };

  outputs = { self, nixpkgs, spectr, ... }: {
    # Define shell for a specific system
    devShells.x86_64-linux.default =
      let
        pkgs = nixpkgs.legacyPackages.x86_64-linux;
      in
      pkgs.mkShell {
        buildInputs = [
          spectr.packages.x86_64-linux.default
          pkgs.git
        ];
      };

    # Add more systems as needed
    devShells.aarch64-darwin.default =
      let
        pkgs = nixpkgs.legacyPackages.aarch64-darwin;
      in
      pkgs.mkShell {
        buildInputs = [
          spectr.packages.aarch64-darwin.default
          pkgs.git
        ];
      };
  };
}
```

## Building from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/connerohnesorge/spectr.git
cd spectr

# Build with Go
go build -o spectr

# Or use Nix
nix build

# Install to your PATH
mv spectr /usr/local/bin/  # or any directory in your PATH
```

## Requirements

- **Go 1.25+** (if building from source)
- **Nix with flakes enabled** (optional, for Nix installation)
- **Git** (for project version control)
