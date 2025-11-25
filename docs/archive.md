## Archive Operations

This module provides commands for creating and extracting archives using `7z`.

### Usage

*   **Zip files:**
    ```bash
    srn archive zip [archive-name] [target-to-archive...]
    ```

*   **Zip files with a password:**
    ```bash
    srn archive zip password [archive-name] [target-to-archive...]
    ```
    *Note: You will be prompted to enter a password.*

*   **Unzip files:**
    ```bash
    srn archive unzip [target-to-unarchive]
    ```

*   **Unzip files with a password:**
    ```bash
    srn archive unzip password [password] [target-to-unarchive]
    ```

### Examples

*   **Create a standard zip archive:**
    ```bash
    srn archive zip my_archive.zip file1.txt my_folder/
    ```

*   **Create a password-protected archive:**
    ```bash
    srn archive zip password my_secret_archive.7z sensitive_data/
    # You will be prompted to enter a password in the terminal
    ```

*   **Extract a standard archive:**
    ```bash
    srn archive unzip my_archive.zip
    ```

*   **Extract a password-protected archive:**
    ```bash
    srn archive unzip password my_secret_archive.7z
    # You will be prompted to enter the password
    ```
