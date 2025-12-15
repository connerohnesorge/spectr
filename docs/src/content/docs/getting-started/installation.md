---
title: Installation
description: How to install Spectr
---

## Using GitHub Releases

The easiest way to install Spectr is to download a pre-built binary from [GitHub Releases](https://github.com/connerohnesorge/spectr/releases):

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
