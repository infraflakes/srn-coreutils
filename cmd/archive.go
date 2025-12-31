package cmd

import (
	"github.com/spf13/cobra"
	"github.com/infraflakes/srn-coreutils/internal/archive"
)

func init() {
	rootCmd.AddCommand(archiveCmd)

	archiveCmd.AddCommand(archive.UnzipCmd)
	archiveCmd.AddCommand(archive.ZipCmd)
}

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive commands",
	Long:  `Commands for creating and extracting archives.`,
}
