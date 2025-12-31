package music

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"

	"github.com/infraflakes/srn-libs/cli"
	"github.com/infraflakes/srn-libs/exec"
	"github.com/infraflakes/srn-libs/fs"
)

func init() {
	ConvertCmd.AddCommand(convertMp3NewCmd)
	ConvertCmd.AddCommand(convertPlaylistCmd)
}

var ConvertCmd = cli.NewCommand(
	"convert",
	"Music related conversion utilities",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var convertMp3NewCmd = cli.NewCommand(
	"mp3 [directories...]",
	"Convert opus/flac to mp3 in one or more directories",
	cobra.MinimumNArgs(1),
	func(cmd *cobra.Command, args []string) {
		// Create a single log file in the current directory for all operations
		logPath := "conversion_errors.log"
		logFile, err := fs.CreateFile(logPath)
		if err != nil {
			os.Exit(1)
		}
		defer fs.CloseFile(logFile)

		for _, dir := range args {
			fmt.Printf("--- Processing directory: %s ---", dir)
			_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				ext := strings.ToLower(filepath.Ext(info.Name()))
				if !info.IsDir() && (ext == ".flac" || ext == ".opus") {
					out := strings.TrimSuffix(path, ext) + ".mp3"

					if _, err := os.Stat(out); err == nil {
						fmt.Println("Skipping (exists):", out)
						return nil
					}

					fmt.Println("Converting + embedding cover:", path, "→", out)
					stderr, convErr := exec.ExecuteCommandWithStderr(
						"ffmpeg",
						"-nostdin",
						"-i", path,
						"-map", "0:a",
						"-map", "0:v?",
						"-c:a", "libmp3lame",
						"-q:a", "0",
						"-id3v2_version", "3",
						"-metadata:s:v", "title=Album cover",
						"-metadata:s:v", "comment=Cover (front)",
						out,
					)

					if convErr != nil {
						fs.LogError(logFile, fmt.Sprintf(
							"Conversion error for %s: %v\nFFmpeg Output:\n%s\n",
							path, convErr, stderr,
						))
						fmt.Printf("Conversion error for %s: %v\nFFmpeg Output:\n%s\n",
							path, convErr, stderr,
						)
						return nil
					}

					fi, err := os.Stat(out)
					if err == nil && fi.Size() > 0 {
						fmt.Println("Converted:", out)
						fmt.Println("Deleting source:", path)
						_ = os.Remove(path)
					} else {
						fs.LogError(logFile, fmt.Sprintf(
							"Conversion failed (zero-size or missing output): %s\n", path,
						))
						fmt.Printf("Conversion failed (zero-size or missing output): %s\n", path)
					}
				}
				return nil
			})
		}

		fmt.Println("\nAll done! Check conversion_errors.log for any errors.")
	},
)

var convertPlaylistCmd = cli.NewCommand(
	"playlist [paths/to/.m3u...]",
	"Format one or more playlists",
	cobra.MinimumNArgs(1),
	func(cmd *cobra.Command, args []string) {
		for _, playlist := range args {
			fmt.Printf("--- Formatting playlist: %s ---", playlist)
			f, err := fs.OpenFile(playlist)
			if err != nil {
				fmt.Printf("Skipping %s: %v\n", playlist, err)
				continue
			}

			var lines []string
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			fs.CloseFile(f) // Close file after reading

			formatted := FormatPlaylistLines(lines)
			output := strings.Join(formatted, "\n")

			if err := os.WriteFile(playlist, []byte(output), 0o644); err != nil {
				fmt.Println("Error writing to file:", err)
				continue
			}

			fmt.Printf("'%s' is now Winamp/Ruizu-safe\n", playlist)
		}
	},
)
