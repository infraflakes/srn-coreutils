package cmd

import (
	"github.com/spf13/cobra"
	"serein/internal/asus"
)

func init() {
	rootCmd.AddCommand(asusCmd)

	asusCmd.AddCommand(asus.AsusProfileCmd)
}

var asusCmd = &cobra.Command{
	Use:   "asus",
	Short: "Asus related commands",
	Long:  `A collection of commands to manage Asus machines via asusctl.`,
}
