---
title: Installation
description: How to install Spectr
---

## Using GitHub Releases

The easiest way to install Spectr is to download a pre-built binary from [GitHub
Releases](https://github.com/connerohnesorge/spectr/releases):

```bash
# Download the latest release for your platform
# Replace {VERSION} with the desired version (e.g., v1.0.0)
# Replace {OS} and {ARCH} with your platform (e.g., Linux_x86_64, Darwin_arm64)

# For Linux x86_64:
curl -LO https://github.com/connerohnesorge/spectr/releases/download/{VERSION}/spectr_{OS}_{ARCH}.tar.gz
tar -xzf spectr_{OS}_{ARCH}.tar.gz
sudo mv spectr /usr/local/bin/

# For macOS arm64 (Apple Silicon):
curl -LO https://github.com/connerohnesorge/spectr/releases/download/v1.0.0/spectr_Darwin_arm64.tar.gz
tar -xzf spectr_Darwin_arm64.tar.gz
sudo mv spectr /usr/local/bin/

# For Windows (PowerShell):
# Download the .zip file from the releases page and extract it
# Then add the directory containing spectr.exe to your PATH
```

### Available Platforms

Pre-built binaries are available for:

- **Linux**: x86_64 (amd64), arm64
- **macOS**: x86_64 (Intel), arm64 (Apple Silicon)
- **Windows**: x86_64 (amd64), arm64

## Using Nix Flakes

Alternatively, you can install Spectr via Nix flakes:

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

Here is a complete example showing how to add spectr to your NixOS system
configuration:

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

To include spectr in a project's development shell, add it to your
`devShells.default` configuration. This is useful when you want spectr available
in your project's development environment:

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
