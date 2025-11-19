# Serein Coreutils

Serein is an opinionated CLI wrapper that replaces cryptic flags with self-explanatory, English-like sub-commands.

## Features

*   **Container Management:** Easily build, delete, list, and run Podman containers, including options for interactive shells, background execution, volume mounts, and device passthrough.
*   **Music Conversion:** Convert audio files (FLAC, Opus) to MP3 format with embedded cover art, and format M3U playlists for compatibility.
*   **Nix System Management:** Manage NixOS system and Home Manager configurations, including building, listing, and deleting generations, and updating flakes.
*   **Archive Operations:** Compress and extract files using 7z, with support for password protection.
*   **YouTube Audio Download:** Download audio from YouTube URLs using `yt-dlp` with embedded thumbnails and metadata.
*   **Git Helper:** A collection of aliases for common Git commands.
*   **File and Directory Search:** Search for words within files, or locate files and directories by name, with optional deletion of matched items.
*   **Todo Management:** An interactive terminal-based application for managing todo lists with contexts, priorities, and more.

## Installation

### Quick Try (Run without Installation)

If you want to quickly try `serein` without installing it permanently:

1.  **Ensure Nix is installed** on your system with flake support enabled.
2.  **Run Serein directly:**

    **Stable binary release:**
    ```bash
    nix run github:nixuris/serein-cli#stable -- [args]
    ```
    **Latest development build:**
    ```bash
    nix run github:nixuris/serein-cli#test -- [args]
    ```
    (Replace `[args]` with any `serein` command and its arguments, e.g., `nix run github:nixuris/serein-cli -- music convert mp3 /path/to/dir`)

### For NixOS/Home Manager Configurations

If you manage your system or user environment with NixOS or Home Manager flakes, you can add `serein-cli` as an input to your configuration:

1.  **Add `serein-cli` as an input in your `flake.nix`:**

    ```nix
    {
      description = "Your personal NixOS/Home Manager configuration";

      inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        home-manager.url = "github:nix-community/home-manager";
        home-manager.inputs.nixpkgs.follows = "nixpkgs";

        # Add serein-cli flake as an input
        serein-cli.url = "github:nixuris/serein-cli";
        serein-cli.inputs.nixpkgs.follows = "nixpkgs"; # Ensure consistent nixpkgs
      };

      outputs = { self, nixpkgs, home-manager, serein-cli, ... } @ inputs: {
        # ... your existing outputs
      };
    }
    ```

2.  **Install `serein` in your NixOS or Home Manager configuration:**

    **Option A: Install System-Wide (NixOS Configuration)**

    ```nix
    # In your configuration.nix (or a NixOS module)
    { config, pkgs, lib, ... }:

    {
      environment.systemPackages = with pkgs; [
        # Reference serein from the serein-cli flake input
        inputs.serein-cli.packages.${pkgs.system}.test # bleeding edge
        inputs.serein-cli.packages.${pkgs.system}.stable # follow release
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
        # Reference serein from the serein-cli flake input
        inputs.serein-cli.packages.${pkgs.system}.test # bleeding edge
        inputs.serein-cli.packages.${pkgs.system}.default # follow release
      ];

      # ... other Home Manager options
    }
    ```

### Binary Distribution (For Non-Nix Users)

For users not using Nix, `serein` can be downloaded as a single executable binary.

1.  **Download the latest release:**
    Visit the [GitHub Releases page](https://github.com/nixuris/serein-cli/releases) and download the wanted binary.

2.  **Make the binary executable:**
    ```bash
    chmod +x ./serein
    ```

3.  **Move the binary to your PATH (optional but recommended):**
    ```bash
    sudo mv ./serein /usr/local/bin/
    ```

### Manual Installation (from source)

If you have a Go environment set up, you can build from source.

1.  **Clone the repo:**
    ```bash
    git clone https://github.com/nixuris/serein-cli
    cd serein-cli
    ```

2.  **Build the binary:**
    The included `Makefile` provides an easy way to build the application:
    ```bash
    make build
    ```
    Alternatively, you can use the standard Go command:
    ```bash
    go build -o serein .
    ```

## Prerequisites and Usage:

Serein acts as a wrapper for several external tools. Using which ultimately depends on your own use case, you can either use --dry-run flag to see the backend tools used, or refer to the [documentation](docs/docs.md):

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the GPL-3.0 License - see the LICENSE file for details.
