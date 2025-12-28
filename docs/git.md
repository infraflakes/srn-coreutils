## Git Commands

This module provides a collection of convenient aliases for common `git` operations, simplifying daily workflows.

### Staging & Local Changes

*   **Check repository status:**
    ```bash
    srn git status
    ```

*   **Stage one or more files:**
    ```bash
    srn git stage <path1> [path2...]
    ```

*   **Unstage one or more files:**
    ```bash
    srn git unstage <path1> [path2...]
    ```

*   **Show changes for a file:**
    ```bash
    srn git changes <path>
    ```

*   **Discard changes in one or more files:**
    ```bash
    srn git undo <path1> [path2...]
    ```

### Commits

*   **Commit staged changes with a message:**
    ```bash
    srn git commit "<message>"
    ```

*   **Push commits to a remote branch:**
    ```bash
    srn git commit push <branch-name> [force]
    ```

*   **List commits in a condensed view:**
    ```bash
    srn git commit list
    ```

*   **Revert a commit by its SHA:**
    ```bash
    srn git commit undo <SHA>
    ```

*   **Delete previous commits (soft reset, keeps changes staged):**
    ```bash
    srn git commit delete stage <number-of-commits>
    ```

*   **Delete previous commits (mixed reset, keeps changes in working directory):**
    ```bash
    srn git commit delete unstage <number-of-commits>
    ```

### Branches & Syncing

*   **List all local and remote branches:**
    ```bash
    srn git branch list
    ```

*   **Create and switch to a new branch:**
    ```bash
    srn git branch create <branch-name>
    ```

*   **Switch to an existing branch:**
    ```bash
    srn git branch switch <branch-name>
    ```

*   **Delete a local branch:**
    ```bash
    srn git branch delete local <branch-name>
    ```

*   **Delete a remote branch:**
    ```bash
    srn git branch delete remote <branch-name>
    ```

*   **Pull changes from a remote branch:**
    ```bash
    srn git sync <branch-name>
    ```

### Tags

*   **Create an annotated tag:**
    ```bash
    srn git tag create <SHA> <tag-name> "<message>"
    ```

*   **Delete a local tag:**
    ```bash
    srn git tag delete local <tag-name>
    ```

*   **Delete a remote tag:**
    ```bash
    srn git tag delete remote <tag-name>
    ```

*   **Delete a tag from both local and remote (use with caution):**
    ```bash
    srn git tag wipe <tag-name>
    ```
