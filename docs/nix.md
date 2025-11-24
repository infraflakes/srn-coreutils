## Nix System Management

This module provides helper commands for managing Nix, NixOS, and Home Manager configurations.

### General Commands

*   **Update all flake inputs:**
    ```bash
    serein nix update
    ```

*   **Search for a package in nixpkgs:**
    ```bash
    serein nix search <package-name>
    ```

*   **Fetch a package and check its hash:**
    ```bash
    serein nix hash <url-to-pkg>
    ```

*   **Run garbage collection to clean the Nix store:**
    ```bash
    serein nix clean
    ```

*   **Format all `.nix` files in the current directory with Alejandra:**
    ```bash
    serein nix lint
    ```

### NixOS System Management (`sys`)

These commands manage the system-level configuration for a NixOS machine.

*   **Build and switch to a new NixOS configuration:**
    ```bash
    sudo serein nix sys build <path/to/flake#config>
    ```

*   **List all system generations:**
    ```bash
    serein nix sys gen
    ```

*   **Delete specific system generations:**
    ```bash
    sudo serein nix sys delete <generation-number>
    # Or delete a range of generations
    sudo serein nix sys delete <start-number>-<end-number>
    ```

### Home Manager Management (`home`)

These commands manage the user-level configuration for Home Manager.

*   **Build and switch to a new Home Manager configuration:**
    ```bash
    serein nix home build <path/to/flake#home>
    ```

*   **List all Home Manager generations:**
    ```bash
    serein nix home gen
    ```

*   **Delete specific Home Manager generations:**
    ```bash
    serein nix home delete <generation-number>
    # Or delete a range of generations
    serein nix home delete <start-number>-<end-number>
    ```

### Examples

*   **Update all flake inputs for a project:**
    ```bash
    serein nix update
    ```
*   **Fetch a package and check its hash:**
    ```bash
    serein nix hash https://github.com/nixuris/srn-coreutils/releases/download/v3.0.0/serein_3.0.0_linux_amd64.tar.gz
    ```

*   **Build the NixOS configuration from `/etc/nixos`:**
    ```bash
    sudo serein nix sys build /etc/nixos/#my-nixos-config
    ```

*   **Build a home-manager configuration:**
    ```bash
    serein nix home build .#my-home-config
    ```

*   **Delete system generations 10 through 15:**
    ```bash
    sudo serein nix sys delete 10-15
    ```
