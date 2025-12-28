## Find Operations

The `find` command provides convenient wrappers for `grep` and `find` to locate words, files, or directories. Each command includes an optional `delete` mode.

### Usage

*   **Find words within files:**
    ```bash
    srn find word <path> <terms...>
    ```

*   **Find and delete files containing specific words:**
    ```bash
    srn find word delete <path> <terms...>
    ```

*   **Find files by name:**
    ```bash
    srn find file <path> <terms...>
    ```

*   **Find and delete files by name:**
    ```bash
    srn find file delete <path> <terms...>
    ```

*   **Find directories by name:**
    ```bash
    srn find dir <path> <terms...>
    ```

*   **Find and delete directories by name:**
    ```bash
    srn find dir delete <path> <terms...>
    ```

### Examples

*   **Search for all occurrences of the word "error" in the current directory:**
    ```bash
    srn find word . "error"
    ```

*   **Find all files in `src/` containing "TODO" or "FIXME" and be prompted to delete them:**
    ```bash
    srn find word delete ./src "TODO" "FIXME"
    ```

*   **Find all files with "config" in their name:**
    ```bash
    srn find file . "config"
    ```

*   **Find and be prompted to delete all directories named `__pycache__`:**
    ```bash
    srn find dir delete . "__pycache__"
    ```
