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
