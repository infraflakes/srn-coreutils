## Nix System Management

This module provides helper commands for managing Nix, NixOS, and Home Manager configurations.

### General Commands

*   **Update all flake inputs:**
    ```bash
    srn nix update
    ```

*   **Search for a package in nixpkgs:**
    ```bash
    srn nix search <package-name>
    ```

*   **Fetch a package and check its hash:**
    ```bash
    srn nix hash <url-to-pkg>
    ```

*   **Run garbage collection to clean the Nix store:**
    ```bash
    srn nix clean
    ```

*   **Format all `.nix` files in the current directory with Alejandra:**
    ```bash
    srn nix lint
    ```

### NixOS System Management (`sys`)

These commands manage the system-level configuration for a NixOS machine.

*   **Build and switch to a new NixOS configuration:**
    ```bash
    sudo srn nix sys build <path/to/flake#config>
    ```

*   **List all system generations:**
    ```bash
    srn nix sys gen
    ```

*   **Delete specific system generations:**
    ```bash
    sudo srn nix sys delete <generation-number>
    # Or delete a range of generations
    sudo srn nix sys delete <start-number>-<end-number>
    ```

### Home Manager Management (`home`)

These commands manage the user-level configuration for Home Manager.

*   **Build and switch to a new Home Manager configuration:**
    ```bash
    srn nix home build <path/to/flake#home>
    ```

*   **List all Home Manager generations:**
    ```bash
    srn nix home gen
    ```

*   **Delete specific Home Manager generations:**
    ```bash
    srn nix home delete <generation-number>
    # Or delete a range of generations
    srn nix home delete <start-number>-<end-number>
    ```

### Examples

*   **Update all flake inputs for a project:**
    ```bash
    srn nix update
    ```
*   **Fetch a package and check its hash:**
    ```bash
    srn nix hash https://github.com/infraflakes/srn-coreutils/releases/download/v3.0.0/srn_3.0.0_linux_amd64.tar.gz
    ```

*   **Build the NixOS configuration from `/etc/nixos`:**
    ```bash
    sudo srn nix sys build /etc/nixos/#my-nixos-config
    ```

*   **Build a home-manager configuration:**
    ```bash
    srn nix home build .#my-home-config
    ```

*   **Delete system generations 10 through 15:**
    ```bash
    sudo srn nix sys delete 10-15
    ```
