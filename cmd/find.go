package cmd

import (
	"github.com/spf13/cobra"
	"srn/internal/find"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Search for words, files, or directories",
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.AddCommand(find.WordCmd)
	findCmd.AddCommand(find.FileCmd)
	findCmd.AddCommand(find.DirCmd)
}
