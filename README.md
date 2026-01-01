# Serein Coreutils

Serein is an opinionated CLI that rejects cryptic flags with self-explanatory, English-like sub-commands.

## Features

*   **Music Conversion:** Convert audio files (FLAC, Opus) to MP3 format with embedded cover art, and format M3U playlists for compatibility.
*   **Nix System Management:** Manage NixOS system and Home Manager configurations, including building, listing, and deleting generations, and updating flakes.
*   **Archive Operations:** Compress and extract files using 7z, with support for password protection.
*   **YouTube Audio Download:** Download audio from YouTube URLs using `yt-dlp` with embedded thumbnails and metadata.
*   **File and Directory Search:** Search for words within files, or locate files and directories by name, with optional deletion of matched items.
*   **Todo Management:** An interactive terminal-based application for managing todo lists with contexts, priorities, and more.

## Installation

### Quick Try (Run without Installation)

If you want to quickly try it without installing it permanently:

1.  **Ensure Nix is installed** on your system with flake support enabled.
2.  **Run the CLI directly from GitHub:**
    You can run the stable build:
    ```bash
    nix run github:infraflakes/srn-coreutils -- [args]
    ```
    Or the latest development build:
    ```bash
    nix run github:infraflakes/srn-coreutils/dev -- [args]
    ```

    (Replace `[args]` with any command and its arguments, e.g., `nix run github:infraflakes/srn-coreutils -- music convert mp3 /path/to/dir`)

### For NixOS/Home Manager Configurations

If you manage your system or user environment with NixOS or Home Manager flakes, you can add `srn-coreutils` as an input to your configuration.

1.  **Add `srn-coreutils` as an input in your `flake.nix`:**

    You can choose which branch to follow. For the stable version, use the `main` branch. For the latest development version, use the `dev` branch.

    ```nix
    {
      description = "Your personal NixOS/Home Manager configuration";

      inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        home-manager.url = "github:nix-community/home-manager";
        home-manager.inputs.nixpkgs.follows = "nixpkgs";

        # Add srn-coreutils flake input, tracking the main (stable) branch
        srn-coreutils = {
          url = "github:infraflakes/srn-coreutils";
          # Or, to track the dev branch:
          # url = "github:infraflakes/srn-coreutils/dev";
          inputs.nixpkgs.follows = "nixpkgs";
        };
      };

      outputs = { self, nixpkgs, home-manager, srn-coreutils, ... } @ inputs: {
        # ... your existing outputs
      };
    }
    ```

2.  **Install the CLI in your NixOS or Home Manager configuration:**

    The flake provides a `default` package.

    **Option A: Install System-Wide (NixOS Configuration)**

    ```nix
    # In your configuration.nix (or a NixOS module)
    { config, pkgs, lib, ... }:

    {
      environment.systemPackages = with pkgs; [
        # Reference it from the srn-coreutils flake input
        inputs.srn-coreutils.packages.${pkgs.stdenv.hostPlatform.system}.default

      ];

      # ... other system configurations
    }
    ```

    **Option B: Install via Home Manager (User-Specific)**

    ```nix
    # In your Home Manager configuration (e.g., ~/.config/home-manager/home.nix)
    { config, pkgs, ... }:

    {
      home.packages = [
        # Reference it from the srn-coreutils flake input
        inputs.srn-coreutils.packages.${pkgs.stdenv.hostPlatform.system}.default

      ];

      # ... other Home Manager options
    }
    ```

### Binary Distribution (For Non-Nix Users)

For users not using Nix, the CLI can be downloaded as a single executable binary.

1.  **Download the latest release:**
    Visit the [GitHub Releases page](https://github.com/infraflakes/srn-coreutils/releases) and download the wanted binary.

2.  **Make the binary executable:**
    ```bash
    chmod +x srn
    ```

3.  **Move the binary to your PATH (optional but recommended):**
    ```bash
    sudo mv srn /usr/local/bin/
    ```

### Manual Installation (from source)

If you have a Go environment set up, you can build from source.

1.  **Clone the repo:**
    ```bash
    git clone https://github.com/infraflakes/srn-coreutils
    cd srn-coreutils
    ```

2.  **Build the binary:**
    The included `Makefile` provides an easy way to build the application:
    ```bash
    make build
    ```
    Alternatively, you can use the standard Go command:
    ```bash
    go build -o srn .
    ```

## Prerequisites and Usage:

You can refer to the [documentation](docs/docs.md):

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
